package systems

import (
	"fmt"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

const SystemTypeCollision SystemType = "collision"

type CollisionSystem struct {
	BaseSystem
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Transform",
				"Collider",
			},
		},
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (cs *CollisionSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeCollision,
		System:       cs,
		Dependencies: []SystemType{}, // No dependencies - runs independently and checks collisions
		Conflicts:    []SystemType{},
		Provides:     []string{"collision_detection", "packet_catching"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     false,
	}
}

func (cs *CollisionSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Get load balancer
	var loadBalancer Entity
	var packets []Entity
	var powerUps []Entity

	// Separate entities by type
	for _, entity := range cs.FilterEntities(entities) {
		colliderComp := entity.GetCollider()
		if colliderComp == nil {
			continue
		}
		collider := colliderComp

		if entityInterface, ok := entity.(interface{ GetComponentNames() []string }); ok {
			fmt.Printf("[CollisionSystem] Checking entity with components: %v, collider tag: %s\n", entityInterface.GetComponentNames(), collider.GetTag())
		}

		// Check if entity is active and visible
		if !entity.IsActive() {
			continue
		}

		// Categorize entities
		if collider.GetTag() == "loadbalancer" {
			loadBalancer = entity
		} else if entity.HasComponent("PacketType") && !entity.HasComponent("Routing") {
			// Only process packets that are not already being routed
			packets = append(packets, entity)
		} else if entity.HasComponent("PowerUpType") {
			powerUps = append(powerUps, entity)
		}
	}

	// Check for packet collisions with load balancer
	if loadBalancer != nil {
		lbTransformComp := loadBalancer.GetTransform()
		lbColliderComp := loadBalancer.GetCollider()
		if lbTransformComp == nil || lbColliderComp == nil {
			return
		}

		lbTransform := lbTransformComp
		lbCollider := lbColliderComp

		for _, packet := range packets {
			packetTransformComp := packet.GetTransform()
			packetColliderComp := packet.GetCollider()
			if packetTransformComp == nil || packetColliderComp == nil {
				continue
			}

			packetTransform := packetTransformComp
			packetCollider := packetColliderComp

			// Check collision
			if cs.checkCollision(lbTransform, lbCollider, packetTransform, packetCollider) {
				// Packet caught by load balancer - route it instead of destroying
				cs.routePacket(packet, loadBalancer, entities, eventDispatcher)

				// Update Combo on the load balancer (component-based)
				if comboComp := loadBalancer.GetCombo(); comboComp != nil {
					if combo, ok := comboComp.(*components.Combo); ok {
						combo.Increment()
						combo.Timer = 0
						// Optional: publish combo event for UI/score systems
						if combo.Streak > 1 {
							bonus := 0
							// basic bonus mapping similar to previous logic
							switch {
							case combo.Streak >= 10:
								bonus = 50
							case combo.Streak >= 7:
								bonus = 30
							case combo.Streak >= 5:
								bonus = 20
							case combo.Streak >= 3:
								bonus = 10
							}
							cc := combo.Streak
							eventDispatcher.Publish(events.NewEvent(events.EventType("combo_achieved"), &events.EventData{ComboCount: &cc, BonusPoints: &bonus}))
						}
					}
				}

				// Update SLA counters on the load balancer (or any entity holding SLA)
				if slaComp := loadBalancer.GetSLA(); slaComp != nil {
					if sla, ok := slaComp.(*components.SLA); ok {
						sla.IncrementCaught()
						// Publish packet caught
						score := sla.Caught * 10
						eventDispatcher.Publish(events.NewEvent(events.EventPacketCaught, &events.EventData{Score: &score, Packet: packet}))

						// Publish SLA updated for UI
						current := 100.0
						if sla.Total > 0 {
							current = float64(sla.Caught) / float64(sla.Total) * 100.0
							sla.SetCurrent(current)
						}
						remaining := sla.ErrorBudget - sla.Lost
						if remaining < 0 {
							remaining = 0
						}
						eventDispatcher.Publish(events.NewEvent(events.EventSLAUpdated, &events.EventData{
							Current:   &current,
							Caught:    &sla.Caught,
							Lost:      &sla.Lost,
							Remaining: &remaining,
							Budget:    &sla.ErrorBudget,
						}))
					}
				}
			}
		}

		// Check for power-up collisions
		for _, powerUp := range powerUps {
			powerUpTransformComp := powerUp.GetTransform()
			powerUpColliderComp := powerUp.GetCollider()
			if powerUpTransformComp == nil || powerUpColliderComp == nil {
				continue
			}

			powerUpTransform := powerUpTransformComp
			powerUpCollider := powerUpColliderComp

			if cs.checkCollision(lbTransform, lbCollider, powerUpTransform, powerUpCollider) {
				// Power-up collected!
				powerUpTypeComp := powerUp.GetPowerUpType()
				if powerUpTypeComp != nil {
					powerUpType := powerUpTypeComp
					fmt.Printf("Power-up collected: %s\n", powerUpType.GetName())

					// Deactivate power-up
					powerUp.(interface{ SetActive(bool) }).SetActive(false)

					// Publish power-up collected event
					powerupName := powerUpType.GetName()
					eventDispatcher.Publish(events.NewEvent(events.EventPowerUpCollected, &events.EventData{
						Powerup: &powerupName,
						Packet:  powerUp,
					}))

					// Activate timer in PowerUpState component if present
					for _, e := range entities {
						if comp := e.GetComponentByName("PowerUpState"); comp != nil {
							if pstate, ok := comp.(*components.PowerUpState); ok {
								duration := 10.0
								switch powerupName {
								case "SpeedBoost":
									duration = 15.0
								case "DoublePoints":
									duration = 20.0
								case "SlowMotion":
									duration = 12.0
								}
								if pstate.RemainingByName == nil {
									pstate.RemainingByName = make(map[string]float64)
								}
								pstate.RemainingByName[powerupName] = duration
							}
						}
					}

					// Also add a particle effect into ParticleState if present
					lbX, lbY := lbTransform.GetX()+7.5, lbTransform.GetY()+7.5
					if sprite := powerUp.GetSprite(); sprite != nil {
						for _, e := range entities {
							if comp := e.GetComponentByName("ParticleState"); comp != nil {
								if pstate, ok := comp.(*components.ParticleState); ok {
									(&ParticleSystem{}).CreatePowerUpEffect(lbX, lbY, sprite.GetColor(), pstate)
								}
							}
						}
					}
				}
			}
		}
	}

	// Check for packets that fell off screen
	for _, packet := range packets {
		transformComp := packet.GetTransform()
		if transformComp == nil {
			continue
		}
		transform := transformComp

		if transform.GetY() > 600 {
			// Packet missed
			packet.(interface{ SetActive(bool) }).SetActive(false)

			// Update SLA counters on the load balancer (or any entity holding SLA)
			if loadBalancer != nil {
				if slaComp := loadBalancer.GetSLA(); slaComp != nil {
					if sla, ok := slaComp.(*components.SLA); ok {
						sla.IncrementLost()
						score := sla.Caught * 10
						eventDispatcher.Publish(events.NewEvent(events.EventPacketLost, &events.EventData{Score: &score}))

						// Publish SLA updated for UI
						current := 100.0
						if sla.Total > 0 {
							current = float64(sla.Caught) / float64(sla.Total) * 100.0
							sla.SetCurrent(current)
						}
						remaining := sla.ErrorBudget - sla.Lost
						if remaining < 0 {
							remaining = 0
						}
						eventDispatcher.Publish(events.NewEvent(events.EventSLAUpdated, &events.EventData{
							Current:   &current,
							Caught:    &sla.Caught,
							Lost:      &sla.Lost,
							Remaining: &remaining,
							Budget:    &sla.ErrorBudget,
						}))

						// If error budget exhausted, signal game over
						if remaining <= 0 {
							eventDispatcher.Publish(events.NewEvent(events.EventGameOver, &events.EventData{
								Score: &score,
								Lost:  &sla.Lost,
							}))
						}
					}
				}
			}
		}
	}
}

func (cs *CollisionSystem) checkCollision(transform1 components.TransformComponent, collider1 components.ColliderComponent,
	transform2 components.TransformComponent, collider2 components.ColliderComponent) bool {

	// Simple AABB collision detection
	left1 := transform1.GetX()
	right1 := transform1.GetX() + collider1.GetWidth()
	top1 := transform1.GetY()
	bottom1 := transform1.GetY() + collider1.GetHeight()

	left2 := transform2.GetX()
	right2 := transform2.GetX() + collider2.GetWidth()
	top2 := transform2.GetY()
	bottom2 := transform2.GetY() + collider2.GetHeight()

	return !(right1 < left2 || left1 > right2 || bottom1 < top2 || top1 > bottom2)
}

// routePacket routes a packet from the load balancer to a backend
func (cs *CollisionSystem) routePacket(packet Entity, loadBalancer Entity, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Get packet position
	packetTransform := packet.GetTransform()
	if packetTransform == nil {
		return
	}

	// Find available backends
	var backends []Entity
	for _, entity := range entities {
		if entity.HasComponent("BackendAssignment") {
			backends = append(backends, entity)
		}
	}

	if len(backends) == 0 {
		// No backends available, destroy packet
		packet.(interface{ SetActive(bool) }).SetActive(false)
		return
	}

	// Use round-robin to select backend
	// For now, use packet position to determine backend
	backendIndex := int(packetTransform.GetX()/200) % len(backends)
	if backendIndex < 0 {
		backendIndex = 0
	}
	if backendIndex >= len(backends) {
		backendIndex = len(backends) - 1
	}

	selectedBackend := backends[backendIndex]
	backendAssignment := selectedBackend.GetBackendAssignment()
	if backendAssignment != nil {
		backendAssignment.IncrementAssignedPackets()
	}

	// Get original packet speed
	originalSpeed := 150.0 // Default speed
	if physicsComp := packet.GetPhysics(); physicsComp != nil {
		// Calculate the magnitude of the velocity vector
		vx := physicsComp.GetVelocityX()
		vy := physicsComp.GetVelocityY()
		originalSpeed = float64(vx*vx + vy*vy)
		if originalSpeed > 0 {
			originalSpeed = float64(originalSpeed)
		}
	}

	// Add routing component to packet with original speed
	packet.AddComponent(components.NewRouting(backendIndex, originalSpeed))

	// Update packet to route to backend
	cs.updatePacketForRouting(packet, selectedBackend)

	// Create route visual via RouteState if available
	if ptx := packetTransform.GetX(); true {
		startX := ptx + 7.5
		startY := packetTransform.GetY() + 7.5
		if bt := selectedBackend.GetTransform(); bt != nil {
			endX := bt.GetX() + 60.0
			endY := bt.GetY() + 20.0
			for _, e := range entities {
				if comp := e.GetComponentByName("RouteState"); comp != nil {
					if rstate, ok := comp.(*components.RouteState); ok {
						color := packet.GetSprite().GetColor()
						(&RoutingSystem{}).CreateRoute(startX, startY, endX, endY, color, rstate)
					}
				}
			}
		}
	}

	fmt.Printf("Packet routed to backend %d!\n", backendIndex)
}

// updatePacketForRouting updates the packet to move toward the backend
func (cs *CollisionSystem) updatePacketForRouting(packet Entity, backend Entity) {
	// Get packet and backend positions
	packetTransform := packet.GetTransform()
	backendTransform := backend.GetTransform()
	if packetTransform == nil || backendTransform == nil {
		return
	}

	// Calculate direction to backend
	backendX := backendTransform.GetX() + 60.0 // Center of backend (assuming 120 width)
	backendY := backendTransform.GetY() + 20.0 // Center of backend (assuming 40 height)

	packetX := packetTransform.GetX()
	packetY := packetTransform.GetY()

	// Calculate direction vector
	dx := backendX - packetX
	dy := backendY - packetY

	// Normalize and set velocity
	distance := float64(dx*dx + dy*dy)
	if distance > 0 {
		distance = float64(dx*dx + dy*dy)
		speed := 150.0 // pixels per second
		dx = dx / distance * speed
		dy = dy / distance * speed
	}

	// Update packet physics to move toward backend
	physicsComp := packet.GetPhysics()
	if physicsComp != nil {
		physicsComp.SetVelocity(dx, dy)
	}

	// Add routing component to track packet state
	// We'll need to create a routing component or use existing components
	// For now, we'll modify the packet's behavior through physics
}
