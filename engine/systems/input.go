package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
)

const SystemTypeInput SystemType = "input"

type InputSystem struct {
	BaseSystem
	lastMouseX        float64
	lastMouseY        float64
	activeInputMethod string // "keyboard" or "mouse"
	keyboardLastUsed  bool   // Track if keyboard was used in the last frame
}

func NewInputSystem() *InputSystem {
	return &InputSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Transform",
				"State",
			},
		},
		lastMouseX:        0,
		lastMouseY:        0,
		activeInputMethod: "mouse", // Start with mouse as default
		keyboardLastUsed:  false,
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (is *InputSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeInput,
		System:       is,
		Dependencies: []SystemType{},
		Conflicts:    []SystemType{},
		Provides:     []string{"user_input", "load_balancer_control"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     false,
	}
}

func (is *InputSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Handle keyboard input for game control
	is.handleKeyboardInput(eventDispatcher)

	// Handle input for load balancer movement
	filteredEntities := is.FilterEntities(entities)

	for _, entity := range filteredEntities {
		transformComp := entity.GetTransform()
		stateComp := entity.GetState()
		if transformComp == nil || stateComp == nil {
			continue
		}

		transform := transformComp
		state := stateComp

		// Only process input if game is in playing state
		if state.GetState() == "playing" {
			is.handleLoadBalancerInput(transform, eventDispatcher, deltaTime)
		}
	}
}

func (is *InputSystem) handleKeyboardInput(eventDispatcher *events.EventDispatcher) {
	// Handle Ctrl+X to exit game
	if ebiten.IsKeyPressed(ebiten.KeyX) && (ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)) {
		eventDispatcher.Publish(events.NewEvent(events.EventExit, nil))
	}
}

func (is *InputSystem) handleLoadBalancerInput(transform components.TransformComponent, eventDispatcher *events.EventDispatcher, deltaTime float64) {
	// Check for keyboard input
	keyboardMoved := is.handleKeyboardMovement(transform, deltaTime)

	// Check for mouse input (but don't apply it yet)
	mouseX, _ := ebiten.CursorPosition()
	mouseMoved := mouseX != int(is.lastMouseX)

	// Update input method tracking
	if keyboardMoved {
		// Keyboard was used - switch to keyboard mode
		is.activeInputMethod = "keyboard"
		is.keyboardLastUsed = true
	} else if mouseMoved && is.activeInputMethod == "mouse" {
		// Mouse was moved and we're already in mouse mode - stay in mouse mode
		is.keyboardLastUsed = false
	} else if mouseMoved && is.activeInputMethod == "keyboard" && !is.keyboardLastUsed {
		// Mouse was moved and we're in keyboard mode but keyboard hasn't been used recently
		// Switch to mouse mode
		is.activeInputMethod = "mouse"
		is.keyboardLastUsed = false
	} else if !keyboardMoved && !mouseMoved {
		// No input detected - maintain current mode
		is.keyboardLastUsed = false
	}

	// Apply input based on active method
	if is.activeInputMethod == "keyboard" {
		// In keyboard mode - only apply keyboard input
		// (keyboard input was already applied in handleKeyboardMovement)
	} else {
		// In mouse mode - apply mouse input
		is.handleMouseInput(transform, eventDispatcher)
	}

	// Update last mouse position
	is.lastMouseX = float64(mouseX)
}

func (is *InputSystem) handleKeyboardMovement(transform components.TransformComponent, deltaTime float64) bool {
	const moveSpeed = 300.0 // pixels per second

	currentX := transform.GetX()
	newX := currentX

	// Check for A/D keys
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		newX -= moveSpeed * deltaTime
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		newX += moveSpeed * deltaTime
	}

	// Update position if there was movement
	if newX != currentX {
		transform.SetPosition(newX, transform.GetY())
		return true
	}

	return false
}

func (is *InputSystem) handleMouseInput(transform components.TransformComponent, eventDispatcher *events.EventDispatcher) {
	// Get mouse position
	mouseX, _ := ebiten.CursorPosition()

	// Update load balancer position
	transform.SetPosition(float64(mouseX), transform.GetY())
}
