package systems

import (
	"fmt"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

const SystemTypeCombo SystemType = "combo"

type ComboSystem struct {
	BaseSystem
	currentCombo  int
	comboTimer    float64
	comboTimeout  float64
	lastComboTime float64
}

func NewComboSystem() *ComboSystem {
	return &ComboSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Combo",
			},
		},
		currentCombo:  0,
		comboTimer:    0.0,
		comboTimeout:  3.0, // 3 seconds to maintain combo
		lastComboTime: 0.0,
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
	cs.comboTimer += deltaTime

	// Check if combo has expired
	if cs.comboTimer-cs.lastComboTime > cs.comboTimeout {
		if cs.currentCombo > 1 {
			fmt.Printf("Combo expired! Final combo: %d\n", cs.currentCombo)
		}
		cs.currentCombo = 0
	}

	// Update combo components for any entities that have them
	for _, entity := range cs.FilterEntities(entities) {
		comboComp := entity.GetCombo()
		if comboComp == nil {
			continue
		}

		// Update combo data
		comboObj := comboComp.(*components.Combo)
		comboObj.Streak = cs.currentCombo
		comboObj.Timer = cs.comboTimer - cs.lastComboTime
	}
}

func (cs *ComboSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Listen for packet caught events to update combo
	eventDispatcher.Subscribe(events.EventPacketCaught, func(event *events.Event) {
		cs.currentCombo++
		cs.lastComboTime = cs.comboTimer

		// Calculate bonus points based on combo
		bonusPoints := cs.calculateComboBonus()

		if cs.currentCombo > 1 {
			fmt.Printf("Combo! x%d (+%d bonus points)\n", cs.currentCombo, bonusPoints)
		}

		// Publish combo event with bonus points
		comboEvent := events.NewEvent(events.EventType("combo_achieved"), &events.EventData{
			ComboCount:  &cs.currentCombo,
			BonusPoints: &bonusPoints,
		})
		eventDispatcher.Publish(comboEvent)
	})
}

func (cs *ComboSystem) calculateComboBonus() int {
	// Bonus points based on combo multiplier
	switch {
	case cs.currentCombo >= 10:
		return 50 // 10x combo = 50 bonus points
	case cs.currentCombo >= 7:
		return 30 // 7x combo = 30 bonus points
	case cs.currentCombo >= 5:
		return 20 // 5x combo = 20 bonus points
	case cs.currentCombo >= 3:
		return 10 // 3x combo = 10 bonus points
	default:
		return 0
	}
}

func (cs *ComboSystem) GetCurrentCombo() int {
	return cs.currentCombo
}

func (cs *ComboSystem) GetComboTimer() float64 {
	return cs.comboTimer - cs.lastComboTime
}
