package systems

import (
	"fmt"
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
	// Find all routed packets
	var routedPackets []Entity
	var backends []Entity

	// Separate entities
	for _, entity := range entities {
		if entity.HasComponent("Routing") {
			routedPackets = append(routedPackets, entity)
		} else if entity.HasComponent("BackendAssignment") {
			backends = append(backends, entity)
		}
	}

	// Process each routed packet
	for _, packet := range routedPackets {
		prs.processRoutedPacket(packet, backends, deltaTime, eventDispatcher)
	}
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
}
