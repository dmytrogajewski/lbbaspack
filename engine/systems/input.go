package systems

import (
	"fmt"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
)

const SystemTypeInput SystemType = "input"

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
	// Handle mouse input for load balancer movement
	for _, entity := range is.FilterEntities(entities) {
		transformComp := entity.GetTransform()
		stateComp := entity.GetState()
		if transformComp == nil || stateComp == nil {
			continue
		}

		transform := transformComp
		state := stateComp

		// Debug: Log entity state
		if entityInterface, ok := entity.(interface{ GetComponentNames() []string }); ok {
			componentNames := entityInterface.GetComponentNames()
			if len(componentNames) > 0 && componentNames[0] == "Transform" {
				fmt.Printf("[InputSystem] Entity state: %s\n", state.GetState())
			}
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
