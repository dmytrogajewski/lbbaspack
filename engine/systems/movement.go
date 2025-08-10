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
	// Detect global slow motion power-up once per frame
	slowMotionFactor := 1.0
	for _, e := range entities {
		if comp := e.GetComponentByName("PowerUpState"); comp != nil {
			if ps, ok := comp.(*components.PowerUpState); ok {
				if ps.RemainingByName != nil {
					if rem, ok := ps.RemainingByName["SlowMotion"]; ok && rem > 0 {
						slowMotionFactor = 0.5
						break
					}
					if rem, ok := ps.RemainingByName["Time Slow"]; ok && rem > 0 {
						slowMotionFactor = 0.5
						break
					}
				}
			}
		}
	}

	effectiveDelta := deltaTime * slowMotionFactor
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
		physicsObj.Update(effectiveDelta)

		// Update position
		oldX, oldY := transform.GetX(), transform.GetY()
		transform.SetPosition(transform.GetX()+physics.GetVelocityX()*effectiveDelta,
			transform.GetY()+physics.GetVelocityY()*effectiveDelta)

		_ = oldX
		_ = oldY
	}
}
