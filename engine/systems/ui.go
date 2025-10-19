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
}

func NewUISystem(screen *ebiten.Image) *UISystem {
	return &UISystem{
		BaseSystem: BaseSystem{},
	}
}

func (uis *UISystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// UI system doesn't draw in Update - it should be called from Draw method
	// This method is for processing UI logic only
}

func (uis *UISystem) Initialize(eventDispatcher *events.EventDispatcher) {}

// Getter methods for testing
func (uis *UISystem) GetCaughtPackets() int { return 0 }

func (uis *UISystem) GetLostPackets() int { return 0 }

func (uis *UISystem) GetRemainingErrors() int { return 0 }

func (uis *UISystem) GetErrorBudget() int { return 0 }

// Reset method to clear all counters for new game
func (uis *UISystem) Reset() {}

// SetErrorBudget updates the error budget and remaining errors
func (uis *UISystem) SetErrorBudget(budget int) {}

func (uis *UISystem) Draw(screen *ebiten.Image, entities []Entity) {
	// Compute UI model from components
	// SLA (take from load balancer entity that has State)
	currentSLA := 100.0
	targetSLA := 0.0
	remaining := 0
	budget := 0
	caught := 0
	isDDoS := false
	level := 1

	// Read GameSession and Spawner
	for _, e := range entities {
		if comp := e.GetComponentByName("GameSession"); comp != nil {
			if gs, ok := comp.(*components.GameSession); ok {
				level = gs.Level
			}
		}
		if comp := e.GetComponentByName("Spawner"); comp != nil {
			if sp, ok := comp.(*components.Spawner); ok {
				isDDoS = sp.IsDDoSActive
			}
		}
	}
	// SLA from entity that has State (load balancer)
	for _, e := range entities {
		if e.HasComponent("State") {
			if slaComp := e.GetSLA(); slaComp != nil {
				currentSLA = slaComp.GetCurrent()
				targetSLA = slaComp.GetTarget()
				remaining = slaComp.GetErrorsRemaining()
				// Budget is not in interface; read from component when concrete
				if s, ok := slaComp.(*components.SLA); ok {
					budget = s.ErrorBudget
					caught = s.Caught
				}
			}
			break
		}
	}

	// Draw UI elements
	text.Draw(screen, "LBaaS Packet Catcher - ECS Edition", basicfont.Face7x13, 10, 20, color.White)
	text.Draw(screen, "Use A/D or Arrow Keys to move", basicfont.Face7x13, 10, 35, color.White)
	text.Draw(screen, "Mouse click to move load balancer", basicfont.Face7x13, 10, 50, color.White)
	text.Draw(screen, "Catch falling network packets!", basicfont.Face7x13, 10, 65, color.White)

	// Draw dynamic SLA stats
	slaText := fmt.Sprintf("SLA: %.2f%% (Target: %.2f%%)", currentSLA, targetSLA)
	errorBudgetText := fmt.Sprintf("Errors: %d/%d left", remaining, budget)
	scoreText := fmt.Sprintf("Score: %d", caught*10)

	// Draw SLA stats
	text.Draw(screen, slaText, basicfont.Face7x13, 300, 20, color.RGBA{200, 255, 200, 255})
	text.Draw(screen, errorBudgetText, basicfont.Face7x13, 300, 35, color.RGBA{255, 200, 200, 255})

	// Draw score
	text.Draw(screen, scoreText, basicfont.Face7x13, 10, 80, color.White)

	// Draw DDoS warning banner
	if isDDoS {
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

	levelText := fmt.Sprintf("Level: %d", level)

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
