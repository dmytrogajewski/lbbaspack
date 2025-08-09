package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

const SystemTypeMovement SystemType = "movement"

type MovementSystem struct {
	BaseSystem
}

func NewMovementSystem() *MovementSystem {
	return &MovementSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Transform",
				"Physics",
			},
		},
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

		_ = oldX
		_ = oldY
	}
}
