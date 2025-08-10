package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
)

const SystemTypeInput SystemType = "input"

type InputSystem struct {
	BaseSystem
}

func NewInputSystem() *InputSystem {
	return &InputSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Transform",
				"State",
			},
		},
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

	// Read active power-ups (Speed Boost) from any holder entity
	moveSpeed := 300.0
	for _, e := range entities {
		if comp := e.GetComponentByName("PowerUpState"); comp != nil {
			if ps, ok := comp.(*components.PowerUpState); ok {
				if ps.RemainingByName != nil {
					if rem, ok := ps.RemainingByName["Speed Boost"]; ok && rem > 0 {
						moveSpeed = 450.0
						break
					}
					if rem, ok := ps.RemainingByName["SpeedBoost"]; ok && rem > 0 {
						moveSpeed = 450.0
						break
					}
				}
			}
		}
	}

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
			is.handleLoadBalancerInputWithSpeed(transform, eventDispatcher, deltaTime, moveSpeed)
		}
	}
}

func (is *InputSystem) handleKeyboardInput(eventDispatcher *events.EventDispatcher) {
	// Handle Ctrl+X to exit game
	if ebiten.IsKeyPressed(ebiten.KeyX) && (ebiten.IsKeyPressed(ebiten.KeyControl) || ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight)) {
		eventDispatcher.Publish(events.NewEvent(events.EventExit, nil))
	}
}

// Legacy signature kept for tests; uses default moveSpeed
func (is *InputSystem) handleLoadBalancerInput(transform components.TransformComponent, eventDispatcher *events.EventDispatcher, deltaTime float64) {
	is.handleLoadBalancerInputWithSpeed(transform, eventDispatcher, deltaTime, 300.0)
}

func (is *InputSystem) handleLoadBalancerInputWithSpeed(transform components.TransformComponent, eventDispatcher *events.EventDispatcher, deltaTime float64, moveSpeed float64) {
	// Apply both inputs with keyboard taking precedence if moved this frame
	keyboardMoved := is.handleKeyboardMovementWithSpeed(transform, deltaTime, moveSpeed)
	if !keyboardMoved {
		is.handleMouseInput(transform, eventDispatcher)
	}
}

// Legacy signature kept for tests; uses default moveSpeed
func (is *InputSystem) handleKeyboardMovement(transform components.TransformComponent, deltaTime float64) bool {
	return is.handleKeyboardMovementWithSpeed(transform, deltaTime, 300.0)
}

func (is *InputSystem) handleKeyboardMovementWithSpeed(transform components.TransformComponent, deltaTime float64, moveSpeed float64) bool {
	// moveSpeed provided by caller taking into account any active power-ups

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
	// Only apply mouse movement while left button is pressed
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return
	}

	// Get mouse position
	mouseX, _ := ebiten.CursorPosition()

	// Update load balancer position
	transform.SetPosition(float64(mouseX), transform.GetY())
}
