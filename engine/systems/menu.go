package systems

import (
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const SystemTypeMenu SystemType = "menu"

type MenuSystem struct {
	BaseSystem
	Screen *ebiten.Image
}

func NewMenuSystem(screen *ebiten.Image) *MenuSystem {
	return &MenuSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"MenuState", // Requires MenuState component
			},
		},
		Screen: screen,
	}
}

func (ms *MenuSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Read/write MenuState component
	var menuState *components.MenuState
	for _, e := range entities {
		if comp := e.GetComponentByName("MenuState"); comp != nil {
			if st, ok := comp.(*components.MenuState); ok {
				menuState = st
				break
			}
		}
	}
	if menuState == nil {
		return
	}

	// Handle input with key latch stored in component
	if ebiten.IsKeyPressed(ebiten.KeyUp) && !menuState.KeyLatch {
		menuState.SelectedMode = (menuState.SelectedMode - 1 + len(menuOptions)) % len(menuOptions)
		menuState.KeyLatch = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && !menuState.KeyLatch {
		menuState.SelectedMode = (menuState.SelectedMode + 1) % len(menuOptions)
		menuState.KeyLatch = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && !menuState.KeyLatch {
		ms.startGame(eventDispatcher, menuState.SelectedMode)
		menuState.KeyLatch = true
	}
	if !ebiten.IsKeyPressed(ebiten.KeyUp) && !ebiten.IsKeyPressed(ebiten.KeyDown) && !ebiten.IsKeyPressed(ebiten.KeyEnter) {
		menuState.KeyLatch = false
	}
}

var menuOptions = []string{
	"Mission Critical (99.95% SLA, 3 errors)",
	"Business Critical (99.5% SLA, 10 errors)",
	"Business Operational (99% SLA, 25 errors)",
	"Office Productivity (95% SLA, 50 errors)",
	"Best Effort (90% SLA, 100 errors)",
}

var menuSLA = []float64{99.95, 99.5, 99.0, 95.0, 90.0}
var menuErrors = []int{3, 10, 25, 50, 100}

func (ms *MenuSystem) startGame(eventDispatcher *events.EventDispatcher, selectedMode int) {
	// Publish game start event with selected mode from component
	sla := menuSLA[selectedMode]
	errs := menuErrors[selectedMode]
	eventDispatcher.Publish(events.NewEvent(events.EventGameStart, &events.EventData{
		Mode:   &selectedMode,
		SLA:    &sla,
		Errors: &errs,
	}))
}

func (ms *MenuSystem) Draw(screen *ebiten.Image, entities []Entity) {
	// Draw menu background
	screen.Fill(color.RGBA{20, 20, 40, 255})

	// Draw title
	title := "LBaaS Packet Catcher - ECS Edition"
	text.Draw(screen, title, basicfont.Face7x13, 200, 100, color.White)

	// Draw subtitle
	subtitle := "Select Game Mode:"
	text.Draw(screen, subtitle, basicfont.Face7x13, 200, 130, color.White)

	// Read MenuState to get selected mode
	selected := 0
	for _, e := range entities {
		if comp := e.GetComponentByName("MenuState"); comp != nil {
			if st, ok := comp.(*components.MenuState); ok {
				selected = st.SelectedMode
				break
			}
		}
	}

	// Draw menu options
	for i, option := range menuOptions {
		y := 160 + i*30
		var col color.Color = color.White
		if i == selected {
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
