package systems

import (
	"fmt"
	"lbbaspack/engine/events"
)

const SystemTypePowerUp SystemType = "powerup"

type PowerUpSystem struct {
	BaseSystem
	activePowerUps map[string]float64 // powerup name -> remaining time
}

func NewPowerUpSystem() *PowerUpSystem {
	return &PowerUpSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"PowerUpType",
			},
		},
		activePowerUps: make(map[string]float64),
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
	// Update active power-ups
	for powerUpName, remainingTime := range pus.activePowerUps {
		pus.activePowerUps[powerUpName] = remainingTime - deltaTime
		if pus.activePowerUps[powerUpName] <= 0 {
			delete(pus.activePowerUps, powerUpName)
			fmt.Printf("Power-up %s expired\n", powerUpName)
		}
	}
}

func (pus *PowerUpSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Listen for power-up collected events
	eventDispatcher.Subscribe(events.EventPowerUpCollected, func(event *events.Event) {
		if event.Data.Powerup != nil {
			pus.activatePowerUp(*event.Data.Powerup, eventDispatcher)
		}
	})
}

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

	pus.activePowerUps[powerUpName] = duration
	fmt.Printf("Power-up %s activated for %.1f seconds\n", powerUpName, duration)

	// Publish power-up activated event
	eventDispatcher.Publish(events.NewEvent(events.EventPowerUpActivated, &events.EventData{
		Powerup:  &powerUpName,
		Duration: &duration,
	}))
}

func (pus *PowerUpSystem) IsPowerUpActive(powerUpName string) bool {
	_, active := pus.activePowerUps[powerUpName]
	return active
}

func (pus *PowerUpSystem) GetActivePowerUps() map[string]float64 {
	return pus.activePowerUps
}
