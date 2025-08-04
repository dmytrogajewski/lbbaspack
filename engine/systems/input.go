package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
)

type InputSystem struct {
	BaseSystem
	lastMouseX float64
	lastMouseY float64
}

func NewInputSystem() *InputSystem {
	return &InputSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Transform",
				"State",
			},
		},
		lastMouseX: 0,
		lastMouseY: 0,
	}
}

func (is *InputSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Handle mouse input for load balancer movement
	for _, entity := range is.FilterEntities(entities) {
		transformComp := entity.GetTransform()
		stateComp := entity.GetState()
		if transformComp == nil || stateComp == nil {
			continue
		}

		transform, ok1 := transformComp.(components.TransformComponent)
		state, ok2 := stateComp.(components.StateComponent)
		if !ok1 || !ok2 {
			continue
		}

		// Only process input if game is in playing state
		if state.GetState() == "playing" {
			is.handleMouseInput(transform, eventDispatcher)
		}
	}
}

func (is *InputSystem) handleMouseInput(transform components.TransformComponent, eventDispatcher *events.EventDispatcher) {
	// Get mouse position
	mouseX, _ := ebiten.CursorPosition()

	// Update load balancer position
	transform.SetPosition(float64(mouseX), transform.GetY())
}
