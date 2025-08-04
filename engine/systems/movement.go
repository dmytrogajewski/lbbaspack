package systems

import (
	"fmt"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

const SystemTypeMovement SystemType = "movement"

type MovementSystem struct {
	BaseSystem
	callCount int
}

func NewMovementSystem() *MovementSystem {
	return &MovementSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Transform",
				"Physics",
			},
		},
		callCount: 0,
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (ms *MovementSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeMovement,
		System:       ms,
		Dependencies: []SystemType{SystemTypeSpawn},
		Conflicts:    []SystemType{},
		Provides:     []string{"entity_movement", "physics_update"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     false,
	}
}

func (ms *MovementSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	ms.callCount++

	packetCount := 0
	for _, entity := range ms.FilterEntities(entities) {
		transformComp := entity.GetTransform()
		physicsComp := entity.GetPhysics()

		if transformComp == nil || physicsComp == nil {
			continue
		}

		transform := transformComp
		physics := physicsComp

		// Update physics
		physicsObj := physicsComp.(*components.Physics)
		physicsObj.Update(deltaTime)

		// Update position
		oldX, oldY := transform.GetX(), transform.GetY()
		transform.SetPosition(transform.GetX()+physics.GetVelocityX()*deltaTime,
			transform.GetY()+physics.GetVelocityY()*deltaTime)

		if ms.callCount%60 == 0 {
			if entityInterface, ok := entity.(interface{ GetComponentNames() []string }); ok {
				fmt.Printf("[MovementSystem] Entity moved from (%.1f, %.1f) to (%.1f, %.1f), components: %v\n", oldX, oldY, transform.GetX(), transform.GetY(), entityInterface.GetComponentNames())
			}
		}

		// Check if this is a packet (has PacketType component)
		if entity.HasComponent("PacketType") {
			packetCount++
			if ms.callCount%60 == 0 { // Print every second
				fmt.Printf("[MovementSystem] Packet at (%.1f, %.1f) -> (%.1f, %.1f), vel(%.1f, %.1f)\n",
					oldX, oldY, transform.GetX(), transform.GetY(), physics.GetVelocityX(), physics.GetVelocityY())
			}
		}
	}

	if ms.callCount%60 == 0 && packetCount > 0 {
		fmt.Printf("[MovementSystem] Processing %d packets\n", packetCount)
	}
}
