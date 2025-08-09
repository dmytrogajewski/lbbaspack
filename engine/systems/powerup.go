package systems

import (
	"fmt"
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
		// If no holder, nothing to update
		return
	}

	// Update active power-ups in component
	for name, rem := range state.RemainingByName {
		newRem := rem - deltaTime
		if newRem <= 0 {
			delete(state.RemainingByName, name)
			fmt.Printf("Power-up %s expired\n", name)
		} else {
			state.RemainingByName[name] = newRem
		}
	}
}

func (pus *PowerUpSystem) Initialize(eventDispatcher *events.EventDispatcher) {}

func (pus *PowerUpSystem) activatePowerUp(powerUpName string, eventDispatcher *events.EventDispatcher) {
	// Set default duration for power-ups
	duration := 10.0 // 10 seconds default

	// Set power-up specific durations
	switch powerUpName {
	case "SpeedBoost":
		duration = 15.0
	case "DoublePoints":
		duration = 20.0
	case "SlowMotion":
		duration = 12.0
	}

	// Note: In the stateless model, activation should be performed by a caller that has access to PowerUpState.
	fmt.Printf("Power-up %s activated for %.1f seconds\n", powerUpName, duration)

	// Publish power-up activated event
	eventDispatcher.Publish(events.NewEvent(events.EventPowerUpActivated, &events.EventData{
		Powerup:  &powerUpName,
		Duration: &duration,
	}))
}

// Removed stateful getters
