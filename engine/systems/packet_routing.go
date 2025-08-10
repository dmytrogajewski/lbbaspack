package systems

import (
	"fmt"
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

const SystemTypePacketRouting SystemType = "packet_routing"

type PacketRoutingSystem struct {
	BaseSystem
}

func NewPacketRoutingSystem() *PacketRoutingSystem {
	return &PacketRoutingSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Transform",
				"Routing",
			},
		},
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (prs *PacketRoutingSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypePacketRouting,
		System:       prs,
		Dependencies: []SystemType{SystemTypeCollision},
		Conflicts:    []SystemType{},
		Provides:     []string{"packet_routing", "backend_delivery"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     false,
	}
}

func (prs *PacketRoutingSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Gather backends once
	var backends []Entity
	for _, entity := range entities {
		if entity.HasComponent("BackendAssignment") {
			backends = append(backends, entity)
		}
	}

	// Assign routing for packets that have a RoutingRequest but no Routing yet
	for _, entity := range entities {
		if entity.HasComponent("RoutingRequest") && !entity.HasComponent("Routing") {
			prs.assignRouting(entity, backends, entities)
			// remove request marker
			entity.RemoveComponent("RoutingRequest")
		}
	}

	// Update routed packets
	for _, entity := range entities {
		if entity.HasComponent("Routing") {
			prs.processRoutedPacket(entity, backends, deltaTime, eventDispatcher)
		}
	}
}

// assignRouting chooses a backend (respects Auto-Balancer) and adds Routing
func (prs *PacketRoutingSystem) assignRouting(packet Entity, backends []Entity, entities []Entity) {
	if len(backends) == 0 {
		packet.(interface{ SetActive(bool) }).SetActive(false)
		return
	}

	// Check Auto-Balancer power-up
	autoBalancerActive := false
	for _, e := range entities {
		if comp := e.GetComponentByName("PowerUpState"); comp != nil {
			if ps, ok := comp.(*components.PowerUpState); ok && ps.RemainingByName != nil {
				if rem, ok := ps.RemainingByName["Auto-Balancer"]; ok && rem > 0 {
					autoBalancerActive = true
					break
				}
			}
		}
	}

	backendIndex := 0
	if autoBalancerActive {
		minCount := int(^uint(0) >> 1)
		for i, be := range backends {
			if ba := be.GetBackendAssignment(); ba != nil {
				if ba.GetAssignedPackets() < minCount {
					minCount = ba.GetAssignedPackets()
					backendIndex = i
				}
			}
		}
	} else {
		if tx := packet.GetTransform(); tx != nil {
			backendIndex = int(tx.GetX()/200) % len(backends)
			if backendIndex < 0 {
				backendIndex = 0
			}
			if backendIndex >= len(backends) {
				backendIndex = len(backends) - 1
			}
		}
	}

	selectedBackend := backends[backendIndex]
	if ba := selectedBackend.GetBackendAssignment(); ba != nil {
		ba.IncrementAssignedPackets()
	}

	// Queue particle effect for packet catch
	if pt := packet.GetTransform(); pt != nil {
		x := pt.GetX() + 7.5
		y := pt.GetY() + 7.5
		var col color.RGBA
		if s := packet.GetSprite(); s != nil {
			col = s.GetColor()
		}
		for _, e := range entities {
			if comp := e.GetComponentByName("ParticleState"); comp != nil {
				if state, ok := comp.(*components.ParticleState); ok {
					state.Requests = append(state.Requests, components.NewParticleEffectRequest(x, y, col, "packet"))
					break
				}
			}
		}
	}

	// original speed
	originalSpeed := 150.0
	if physicsComp := packet.GetPhysics(); physicsComp != nil {
		vx := physicsComp.GetVelocityX()
		vy := physicsComp.GetVelocityY()
		s := float64(vx*vx + vy*vy)
		if s > 0 {
			originalSpeed = s
		}
	}
	packet.AddComponent(components.NewRouting(selectedBackend.GetBackendAssignment().GetBackendID(), originalSpeed))
}

func (prs *PacketRoutingSystem) processRoutedPacket(packet Entity, backends []Entity, deltaTime float64, eventDispatcher *events.EventDispatcher) {
	routingComp := packet.GetRouting()
	transformComp := packet.GetTransform()
	physicsComp := packet.GetPhysics()

	if routingComp == nil || transformComp == nil {
		return
	}

	targetBackendID := routingComp.GetTargetBackendID()

	// Find target backend
	var targetBackend Entity
	for _, backend := range backends {
		backendAssignment := backend.GetBackendAssignment()
		if backendAssignment != nil && backendAssignment.GetBackendID() == targetBackendID {
			targetBackend = backend
			break
		}
	}

	if targetBackend == nil {
		// Target backend not found, destroy packet
		packet.(interface{ SetActive(bool) }).SetActive(false)
		return
	}

	// Get backend position
	backendTransform := targetBackend.GetTransform()
	if backendTransform == nil {
		return
	}

	// Calculate distance to backend
	packetX := transformComp.GetX()
	packetY := transformComp.GetY()
	backendX := backendTransform.GetX() + 60.0 // Center of backend
	backendY := backendTransform.GetY() + 20.0 // Center of backend

	dx := backendX - packetX
	dy := backendY - packetY
	distance := dx*dx + dy*dy

	// Check if packet has reached the backend
	if distance < 100.0 { // Within 10 pixels of backend center
		// Packet delivered to backend
		fmt.Printf("Packet delivered to backend %d!\n", targetBackendID)

		// Remove routing component and destroy packet
		packet.RemoveComponent("Routing")
		packet.(interface{ SetActive(bool) }).SetActive(false)

		// Publish packet delivered event
		eventDispatcher.Publish(events.NewEvent(events.EventPacketDelivered, &events.EventData{
			BackendID: &targetBackendID,
		}))

		return
	}

	// Update packet movement toward backend
	if physicsComp != nil {
		// Calculate direction vector
		if distance > 0 {
			// Use the original packet speed
			speed := routingComp.GetOriginalSpeed()
			if speed <= 0 {
				speed = 150.0 // Fallback speed
			}
			normalizedDx := dx / distance * speed
			normalizedDy := dy / distance * speed
			physicsComp.SetVelocity(normalizedDx, normalizedDy)
		}
	}

	// Update routing progress
	initialDistance := 400.0 // Approximate initial distance
	currentProgress := 1.0 - (distance / (initialDistance * initialDistance))
	if currentProgress < 0 {
		currentProgress = 0
	}
	if currentProgress > 1 {
		currentProgress = 1
	}
	routingComp.SetRouteProgress(currentProgress)
}

func (prs *PacketRoutingSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Initialize packet routing system
	fmt.Println("[PacketRoutingSystem] Initialized")
	// Subscribe to collision: if a loadbalancer collides with a packet, mark RoutingRequest
	eventDispatcher.Subscribe(events.EventCollisionDetected, func(event *events.Event) {
		if event == nil || event.Data == nil {
			return
		}
		if event.Data.TagA == nil || event.Data.TagB == nil {
			return
		}
		aIsLB := *event.Data.TagA == "loadbalancer"
		bIsLB := *event.Data.TagB == "loadbalancer"
		aIsPacket := *event.Data.TagA == "packet"
		bIsPacket := *event.Data.TagB == "packet"
		if (aIsLB && bIsPacket) || (bIsLB && aIsPacket) {
			var packet Entity
			if aIsPacket {
				if p, ok := event.Data.EntityA.(Entity); ok {
					packet = p
				}
			} else {
				if p, ok := event.Data.EntityB.(Entity); ok {
					packet = p
				}
			}
			if packet != nil && !packet.HasComponent("Routing") && !packet.HasComponent("RoutingRequest") {
				packet.AddComponent(components.NewRoutingRequest())
			}
		}
	})
}
