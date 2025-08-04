package systems

import (
	"fmt"
	"image/color"
	"lbbaspack/engine/events"

	"lbbaspack/engine/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const SystemTypeUI SystemType = "ui"

type UISystem struct {
	BaseSystem
	score           int
	currentSLA      float64
	targetSLA       float64
	caughtPackets   int
	lostPackets     int
	remainingErrors int
	errorBudget     int
	level           int  // Current game level
	isDDoSActive    bool // Show DDoS warning
}

func NewUISystem(screen *ebiten.Image) *UISystem {
	return &UISystem{
		BaseSystem:      BaseSystem{},
		score:           0,
		currentSLA:      100.0,
		targetSLA:       99.5,
		caughtPackets:   0,
		lostPackets:     0,
		remainingErrors: 10,
		errorBudget:     10,
		level:           1, // Start at level 1
		isDDoSActive:    false,
	}
}

func (uis *UISystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// UI system doesn't draw in Update - it should be called from Draw method
	// This method is for processing UI logic only
}

func (uis *UISystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Listen for SLA update events
	eventDispatcher.Subscribe(events.EventSLAUpdated, func(event *events.Event) {
		if event.Data.Current != nil {
			uis.currentSLA = *event.Data.Current
		}
		if event.Data.Target != nil {
			uis.targetSLA = *event.Data.Target
		}
		if event.Data.Caught != nil {
			uis.caughtPackets = *event.Data.Caught
		}
		if event.Data.Lost != nil {
			uis.lostPackets = *event.Data.Lost
		}
		if event.Data.Remaining != nil {
			uis.remainingErrors = *event.Data.Remaining
		}
		if event.Data.Budget != nil {
			uis.errorBudget = *event.Data.Budget
		}
	})

	// Note: Packet lost events are handled by the SLA system which publishes EventSLAUpdated
	// The UI system gets the updated counts from EventSLAUpdated events

	// Listen for level-up events
	eventDispatcher.Subscribe(events.EventLevelUp, func(event *events.Event) {
		if event.Data.Level != nil {
			uis.level = *event.Data.Level
		}
	})

	// Listen for DDoS events
	eventDispatcher.Subscribe(events.EventDDoSStart, func(event *events.Event) {
		uis.isDDoSActive = true
	})
	eventDispatcher.Subscribe(events.EventDDoSEnd, func(event *events.Event) {
		uis.isDDoSActive = false
	})
}

// Getter methods for testing
func (uis *UISystem) GetCaughtPackets() int {
	return uis.caughtPackets
}

func (uis *UISystem) GetLostPackets() int {
	return uis.lostPackets
}

func (uis *UISystem) GetRemainingErrors() int {
	return uis.remainingErrors
}

func (uis *UISystem) GetErrorBudget() int {
	return uis.errorBudget
}

// Reset method to clear all counters for new game
func (uis *UISystem) Reset() {
	uis.score = 0
	uis.currentSLA = 100.0
	uis.caughtPackets = 0
	uis.lostPackets = 0
	uis.remainingErrors = uis.errorBudget // Reset to current error budget
	uis.level = 1
	uis.isDDoSActive = false
	fmt.Printf("UI system reset - counters cleared, error budget: %d, remaining errors: %d\n", uis.errorBudget, uis.remainingErrors)
}

// SetErrorBudget updates the error budget and remaining errors
func (uis *UISystem) SetErrorBudget(budget int) {
	uis.errorBudget = budget
	uis.remainingErrors = budget // Reset remaining errors to new budget
	fmt.Printf("UI system error budget set to %d, remaining errors: %d\n", uis.errorBudget, uis.remainingErrors)
}

func (uis *UISystem) Draw(screen *ebiten.Image, entities []Entity) {
	// Draw UI elements
	text.Draw(screen, "LBaaS Packet Catcher - ECS Edition", basicfont.Face7x13, 10, 20, color.White)
	text.Draw(screen, "Use A/D or Arrow Keys to move", basicfont.Face7x13, 10, 35, color.White)
	text.Draw(screen, "Mouse click to move load balancer", basicfont.Face7x13, 10, 50, color.White)
	text.Draw(screen, "Catch falling network packets!", basicfont.Face7x13, 10, 65, color.White)

	// Draw dynamic SLA stats
	slaText := fmt.Sprintf("SLA: %.2f%% (Target: %.2f%%)", uis.currentSLA, uis.targetSLA)
	errorBudgetText := fmt.Sprintf("Errors: %d/%d left", uis.remainingErrors, uis.errorBudget)
	scoreText := fmt.Sprintf("Score: %d", uis.caughtPackets*10)

	// Draw SLA stats
	text.Draw(screen, slaText, basicfont.Face7x13, 300, 20, color.RGBA{200, 255, 200, 255})
	text.Draw(screen, errorBudgetText, basicfont.Face7x13, 300, 35, color.RGBA{255, 200, 200, 255})

	// Draw score
	text.Draw(screen, scoreText, basicfont.Face7x13, 10, 80, color.White)

	// Draw DDoS warning banner
	if uis.isDDoSActive {
		text.Draw(screen, "!!! DDoS ATTACK !!!", basicfont.Face7x13, 300, 60, color.RGBA{255, 50, 50, 255})
	}

	// Find combo component
	var comboText string
	for _, entity := range entities {
		if combo := entity.GetComponentByName("Combo"); combo != nil {
			comboComp := combo.(*components.Combo)
			if comboComp.Streak > 1 {
				comboText = fmt.Sprintf("Combo: x%d", comboComp.Streak)
			}
			break
		}
	}

	// Find level from game state
	var levelText string
	for _, entity := range entities {
		if state := entity.GetComponentByName("State"); state != nil {
			stateComp := state.(*components.State)
			if stateComp.Current == components.StatePlaying {
				levelText = fmt.Sprintf("Level: %d", uis.level) // Use tracked level
				break
			}
		}
	}

	// Draw combo
	if comboText != "" {
		text.Draw(screen, comboText, basicfont.Face7x13, 10, 110, color.RGBA{255, 255, 100, 255})
	}

	// Draw level
	if levelText != "" {
		text.Draw(screen, levelText, basicfont.Face7x13, 10, 125, color.RGBA{200, 200, 255, 255})
	}

	// Draw backend stats
	backendY := 150
	for _, entity := range entities {
		if backend := entity.GetComponentByName("BackendAssignment"); backend != nil {
			ba := backend.(*components.BackendAssignment)
			backendText := fmt.Sprintf("Backend %d: %d packets", ba.BackendID, ba.Counter)
			text.Draw(screen, backendText, basicfont.Face7x13, 10, backendY, color.RGBA{100, 255, 100, 255})
			backendY += 15
		}
	}

	// Draw instructions
	text.Draw(screen, "Ctrl+X to exit", basicfont.Face7x13, 10, backendY+10, color.White)
}
