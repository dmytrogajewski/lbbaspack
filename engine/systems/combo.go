package systems

import (
	"fmt"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

const SystemTypeCombo SystemType = "combo"

// comboTimeoutSeconds is a constant threshold for combo expiration
const comboTimeoutSeconds = 3.0

type ComboSystem struct {
	BaseSystem
}

func NewComboSystem() *ComboSystem {
	return &ComboSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Combo",
			},
		},
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (cs *ComboSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeCombo,
		System:       cs,
		Dependencies: []SystemType{SystemTypeCollision},
		Conflicts:    []SystemType{},
		Provides:     []string{"combo_tracking", "score_multiplier"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     true,
	}
}

func (cs *ComboSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Update per-entity combo timers and handle expiration
	for _, entity := range cs.FilterEntities(entities) {
		comboComp := entity.GetCombo()
		if comboComp == nil {
			continue
		}
		combo := comboComp.(*components.Combo)

		// Advance timer since last combo increment
		combo.Timer += deltaTime

		// Expire combo if timeout exceeded
		if combo.Timer > comboTimeoutSeconds && combo.Streak > 0 {
			if combo.Streak > 1 {
				fmt.Printf("Combo expired! Final combo: %d\n", combo.Streak)
			}
			combo.Reset()
		}
	}
}

func (cs *ComboSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Stateless: no internal state; combo increments are handled by collision when a packet is caught
}
