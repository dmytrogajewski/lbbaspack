package systems

import (
	"fmt"
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

const SystemTypePowerUp SystemType = "powerup"

type PowerUpSystem struct {
	BaseSystem
}

func NewPowerUpSystem() *PowerUpSystem {
	return &PowerUpSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{},
		},
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (pus *PowerUpSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypePowerUp,
		System:       pus,
		Dependencies: []SystemType{SystemTypeCollision},
		Conflicts:    []SystemType{},
		Provides:     []string{"powerup_management", "effect_activation"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     true,
	}
}

func (pus *PowerUpSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Find PowerUpState component holder
	var state *components.PowerUpState
	for _, e := range entities {
		if comp := e.GetComponentByName("PowerUpState"); comp != nil {
			if ps, ok := comp.(*components.PowerUpState); ok {
				state = ps
				break
			}
		}
	}
	if state == nil {
		return
	}

	// Apply new activations and deactivate collected entity; queue particle effect
	for _, e := range entities {
		if e.HasComponent("PowerUpActivation") && e.HasComponent("PowerUpType") {
			if put := e.GetPowerUpType(); put != nil {
				name := put.GetName()
				duration := put.GetDuration()
				if state.RemainingByName == nil {
					state.RemainingByName = make(map[string]float64)
				}
				state.RemainingByName[name] = duration
			}
			// Find ParticleState holder and append request
			if t := e.GetTransform(); t != nil {
				x := t.GetX() + 7.5
				y := t.GetY() + 7.5
				var col color.RGBA
				if s := e.GetSprite(); s != nil {
					col = s.GetColor()
				}
				for _, ent := range entities {
					if comp := ent.GetComponentByName("ParticleState"); comp != nil {
						if pstate, ok := comp.(*components.ParticleState); ok {
							pstate.Requests = append(pstate.Requests, components.NewParticleEffectRequest(x, y, col, "powerup"))
							break
						}
					}
				}
			}
			// deactivate the powerup entity; Cleanup will remove it
			e.(interface{ SetActive(bool) }).SetActive(false)
			e.RemoveComponent("PowerUpActivation")
		}
	}

	// Update timers
	for name, rem := range state.RemainingByName {
		newRem := rem - deltaTime
		if newRem <= 0 {
			delete(state.RemainingByName, name)
			fmt.Printf("Power-up %s expired\n", name)
			eventDispatcher.Publish(events.NewEvent(events.EventType("powerup_expired"), &events.EventData{Powerup: &name}))
		} else {
			state.RemainingByName[name] = newRem
		}
	}
}

func (pus *PowerUpSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Subscribe to generalized collision to detect powerup collections
	eventDispatcher.Subscribe(events.EventCollisionDetected, func(event *events.Event) {
		if event == nil || event.Data == nil {
			return
		}
		if event.Data.TagA == nil || event.Data.TagB == nil {
			return
		}
		// Powerup collected when loadbalancer collides with powerup
		aIsLB := *event.Data.TagA == "loadbalancer"
		bIsLB := *event.Data.TagB == "loadbalancer"
		aIsPU := *event.Data.TagA == "powerup"
		bIsPU := *event.Data.TagB == "powerup"
		if (aIsLB && bIsPU) || (bIsLB && aIsPU) {
			var pe Entity
			if aIsPU {
				if e, ok := event.Data.EntityA.(Entity); ok {
					pe = e
				}
			} else {
				if e, ok := event.Data.EntityB.(Entity); ok {
					pe = e
				}
			}
			if pe != nil {
				if put := pe.GetPowerUpType(); put != nil {
					name := put.GetName()
					duration := put.GetDuration()
					// Guard against duplicate activation markers
					if !pe.HasComponent("PowerUpActivation") {
						pe.AddComponent(components.NewPowerUpActivation(name, duration))
					}
				}
			}
		}
	})
}

// Test-only legacy helpers are defined in legacy_test_shims_test.go

// Removed stateful getters
