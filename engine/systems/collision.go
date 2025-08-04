package systems

import (
	"fmt"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

type CollisionSystem struct {
	BaseSystem
	score int
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Transform",
				"Collider",
			},
		},
		score: 0,
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
		} else if entity.HasComponent("PacketType") {
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
				// Packet caught!
				cs.score += 10
				fmt.Printf("Packet caught! Score: %d\n", cs.score)

				// Deactivate packet
				packet.(interface{ SetActive(bool) }).SetActive(false)

				// Publish packet caught event
				eventDispatcher.Publish(events.NewEvent(events.EventPacketCaught, &events.EventData{
					Score:  &cs.score,
					Packet: packet,
				}))
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
					}))
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
			fmt.Printf("Packet missed! Score: %d\n", cs.score)

			// Publish packet lost event
			eventDispatcher.Publish(events.NewEvent(events.EventPacketLost, &events.EventData{
				Score: &cs.score,
			}))
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
