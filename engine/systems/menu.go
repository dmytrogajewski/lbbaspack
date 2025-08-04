package systems

import (
	"image/color"
	"lbbaspack/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const SystemTypeMenu SystemType = "menu"

type MenuSystem struct {
	BaseSystem
	Screen       *ebiten.Image
	selectedMode int
	menuOptions  []string
	menuSLA      []float64
	menuErrors   []int
	keyPressed   bool
}

func NewMenuSystem(screen *ebiten.Image) *MenuSystem {
	return &MenuSystem{
		BaseSystem:   BaseSystem{},
		Screen:       screen,
		selectedMode: 0,
		menuOptions: []string{
			"Mission Critical (99.95% SLA, 3 errors)",
			"Business Critical (99.5% SLA, 10 errors)",
			"Business Operational (99% SLA, 25 errors)",
			"Office Productivity (95% SLA, 50 errors)",
			"Best Effort (90% SLA, 100 errors)",
		},
		menuSLA:    []float64{99.95, 99.5, 99.0, 95.0, 90.0},
		menuErrors: []int{3, 10, 25, 50, 100},
		keyPressed: false,
	}
}

func (ms *MenuSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Handle menu navigation with key state tracking
	if ebiten.IsKeyPressed(ebiten.KeyUp) && !ms.keyPressed {
		ms.selectedMode = (ms.selectedMode - 1 + len(ms.menuOptions)) % len(ms.menuOptions)
		ms.keyPressed = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && !ms.keyPressed {
		ms.selectedMode = (ms.selectedMode + 1) % len(ms.menuOptions)
		ms.keyPressed = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && !ms.keyPressed {
		// Start game with selected mode
		ms.startGame(eventDispatcher)
		ms.keyPressed = true
	}

	// Reset key pressed state when no keys are pressed
	if !ebiten.IsKeyPressed(ebiten.KeyUp) && !ebiten.IsKeyPressed(ebiten.KeyDown) && !ebiten.IsKeyPressed(ebiten.KeyEnter) {
		ms.keyPressed = false
	}
}

func (ms *MenuSystem) startGame(eventDispatcher *events.EventDispatcher) {
	// Publish game start event with selected mode
	eventDispatcher.Publish(events.NewEvent(events.EventGameStart, &events.EventData{
		Mode:   &ms.selectedMode,
		SLA:    &ms.menuSLA[ms.selectedMode],
		Errors: &ms.menuErrors[ms.selectedMode],
	}))
}

func (ms *MenuSystem) Draw(screen *ebiten.Image) {
	// Draw menu background
	screen.Fill(color.RGBA{20, 20, 40, 255})

	// Draw title
	title := "LBaaS Packet Catcher - ECS Edition"
	text.Draw(screen, title, basicfont.Face7x13, 200, 100, color.White)

	// Draw subtitle
	subtitle := "Select Game Mode:"
	text.Draw(screen, subtitle, basicfont.Face7x13, 200, 130, color.White)

	// Draw menu options
	for i, option := range ms.menuOptions {
		y := 160 + i*30
		var col color.Color = color.White
		if i == ms.selectedMode {
			col = color.RGBA{255, 255, 0, 255} // Yellow for selected
		}
		text.Draw(screen, option, basicfont.Face7x13, 150, y, col)
	}

	// Draw instructions
	instructions := []string{
		"Use UP/DOWN arrows to select mode",
		"Press ENTER to start game",
		"",
		"Game Controls:",
		"A/D or Arrow Keys - Move load balancer",
		"Mouse Click - Move to mouse position",
		"Ctrl+X - Exit game",
	}

	for i, instruction := range instructions {
		y := 350 + i*15
		text.Draw(screen, instruction, basicfont.Face7x13, 150, y, color.White)
	}
}
