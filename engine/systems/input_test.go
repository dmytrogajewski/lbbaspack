package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"
)

func TestNewInputSystem(t *testing.T) {
	is := NewInputSystem()

	// Test that the system is properly initialized
	if is == nil {
		t.Fatal("NewInputSystem returned nil")
	}

	// Test required components
	expectedComponents := []string{"Transform", "State"}
	if len(is.RequiredComponents) != len(expectedComponents) {
		t.Errorf("Expected %d required components, got %d", len(expectedComponents), len(is.RequiredComponents))
	}

	for i, component := range expectedComponents {
		if is.RequiredComponents[i] != component {
			t.Errorf("Expected required component %s at index %d, got %s", component, i, is.RequiredComponents[i])
		}
	}

	// Test initial values
	if is.lastMouseX != 0 {
		t.Errorf("Expected initial lastMouseX to be 0, got %f", is.lastMouseX)
	}

	if is.lastMouseY != 0 {
		t.Errorf("Expected initial lastMouseY to be 0, got %f", is.lastMouseY)
	}
}

func TestInputSystem_Update_NoEntities(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test with no entities
	entities := []Entity{}

	// Run update
	is.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// Since this system doesn't have any state that changes without entities,
	// we just verify it doesn't panic or error
}

func TestInputSystem_Update_EntityWithoutRequiredComponents(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with only Transform (missing State)
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 100)
	entity.AddComponent(transform)

	entities := []Entity{entity}

	// Run update
	is.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// The system should skip entities without required components
}

func TestInputSystem_Update_EntityWithTransformOnly(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with only Transform (missing State)
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 100)
	entity.AddComponent(transform)

	entities := []Entity{entity}

	// Run update
	is.Update(0.016, entities, eventDispatcher)

	// Verify transform position remains unchanged
	// (since no state component, input should not be processed)
	transformComp := entity.GetTransform()
	if transformComp == nil {
		t.Fatal("Expected transform component to exist")
	}

	transformObj := transformComp
	x, y := transformObj.GetX(), transformObj.GetY()
	if x != 100 || y != 100 {
		t.Errorf("Expected position to remain (100, 100), got (%f, %f)", x, y)
	}
}

func TestInputSystem_Update_EntityWithStateOnly(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with only State (missing Transform)
	entity := entities.NewEntity(1)
	state := components.NewState(components.StatePlaying)
	entity.AddComponent(state)

	entities := []Entity{entity}

	// Run update
	is.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// The system should skip entities without required components
}

func TestInputSystem_Update_EntityWithBothComponents_NotPlaying(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with both components but not in playing state
	entity := createInputEntity(1, components.StateMenu, 100, 100)

	entities := []Entity{entity}

	// Run update
	is.Update(0.016, entities, eventDispatcher)

	// Verify transform position remains unchanged
	// (since not in playing state, input should not be processed)
	transformComp := entity.GetTransform()
	transformObj := transformComp
	x, y := transformObj.GetX(), transformObj.GetY()
	if x != 100 || y != 100 {
		t.Errorf("Expected position to remain (100, 100), got (%f, %f)", x, y)
	}
}

func TestInputSystem_Update_EntityWithBothComponents_GameOver(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with both components but in game over state
	entity := createInputEntity(1, components.StateGameOver, 100, 100)

	entities := []Entity{entity}

	// Run update
	is.Update(0.016, entities, eventDispatcher)

	// Verify transform position remains unchanged
	// (since not in playing state, input should not be processed)
	transformComp := entity.GetTransform()
	transformObj := transformComp
	x, y := transformObj.GetX(), transformObj.GetY()
	if x != 100 || y != 100 {
		t.Errorf("Expected position to remain (100, 100), got (%f, %f)", x, y)
	}
}

func TestInputSystem_Update_EntityWithBothComponents_Playing(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with both components in playing state
	entity := createInputEntity(1, components.StatePlaying, 100, 100)

	entities := []Entity{entity}

	// Run update
	is.Update(0.016, entities, eventDispatcher)

	// Verify transform position was updated
	// Note: In a real test environment, we would mock ebiten.CursorPosition()
	// For now, we just verify the system doesn't crash and processes the entity
	transformComp := entity.GetTransform()
	if transformComp == nil {
		t.Fatal("Expected transform component to exist")
	}

	// The actual position update depends on ebiten.CursorPosition() which we can't easily mock
	// in this test environment, so we just verify the component exists and the system runs
}

func TestInputSystem_Update_MultipleEntities(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create multiple entities with different states
	entity1 := createInputEntity(1, components.StatePlaying, 100, 100)  // Should process input
	entity2 := createInputEntity(2, components.StateMenu, 200, 200)     // Should not process input
	entity3 := createInputEntity(3, components.StateGameOver, 300, 300) // Should not process input

	entities := []Entity{entity1, entity2, entity3}

	// Run update
	is.Update(0.016, entities, eventDispatcher)

	// Verify only the playing entity was processed
	// The other entities should remain unchanged
	transform2 := entity2.GetTransform()
	x2, y2 := transform2.GetX(), transform2.GetY()
	if x2 != 200 || y2 != 200 {
		t.Errorf("Expected entity2 position to remain (200, 200), got (%f, %f)", x2, y2)
	}

	transform3 := entity3.GetTransform()
	x3, y3 := transform3.GetX(), transform3.GetY()
	if x3 != 300 || y3 != 300 {
		t.Errorf("Expected entity3 position to remain (300, 300), got (%f, %f)", x3, y3)
	}
}

func TestInputSystem_Update_InvalidComponentTypes(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with invalid component types
	entity := entities.NewEntity(1)

	// Add components that don't implement the required interfaces
	entity.AddComponent(&mockComponent{componentType: "Transform"})
	entity.AddComponent(&mockComponent{componentType: "State"})

	entities := []Entity{entity}

	// Run update
	is.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// The system should skip entities with invalid component types
}

func TestInputSystem_Update_EntityWithNullComponents(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with null components (simulated by not adding any)
	entity := entities.NewEntity(1)

	entities := []Entity{entity}

	// Run update
	is.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// The system should handle null components gracefully
}

func TestInputSystem_Update_EntityStateTransition(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity in menu state
	entity := createInputEntity(1, components.StateMenu, 100, 100)

	entities := []Entity{entity}

	// Run update in menu state
	is.Update(0.016, entities, eventDispatcher)

	// Verify position unchanged
	transformComp := entity.GetTransform()
	transformObj := transformComp
	x1, y1 := transformObj.GetX(), transformObj.GetY()
	if x1 != 100 || y1 != 100 {
		t.Errorf("Expected position to remain (100, 100) in menu state, got (%f, %f)", x1, y1)
	}

	// Change to playing state
	stateComp := entity.GetState()
	stateObj := stateComp
	stateObj.SetState("playing")

	// Run update in playing state
	is.Update(0.016, entities, eventDispatcher)

	// Verify position was processed (though actual value depends on ebiten)
	// We just verify the system doesn't crash and processes the entity
	transformObj2 := entity.GetTransform()
	if transformObj2 == nil {
		t.Fatal("Expected transform component to still exist after state transition")
	}
}

func TestInputSystem_Update_EntityStateTransitionBack(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity in playing state
	entity := createInputEntity(1, components.StatePlaying, 100, 100)

	entities := []Entity{entity}

	// Run update in playing state
	is.Update(0.016, entities, eventDispatcher)

	// Change back to menu state
	stateComp := entity.GetState()
	stateObj := stateComp
	stateObj.SetState("menu")

	// Run update in menu state
	is.Update(0.016, entities, eventDispatcher)

	// Verify the system handles the transition gracefully
	transformComp := entity.GetTransform()
	if transformComp == nil {
		t.Fatal("Expected transform component to still exist after state transition")
	}
}

func TestInputSystem_Update_ZeroDeltaTime(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity in playing state
	entity := createInputEntity(1, components.StatePlaying, 100, 100)

	entities := []Entity{entity}

	// Run update with zero delta time
	is.Update(0.0, entities, eventDispatcher)

	// Verify no errors occurred
	// The system should handle zero delta time gracefully
}

func TestInputSystem_Update_LargeDeltaTime(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity in playing state
	entity := createInputEntity(1, components.StatePlaying, 100, 100)

	entities := []Entity{entity}

	// Run update with large delta time
	is.Update(1.0, entities, eventDispatcher)

	// Verify no errors occurred
	// The system should handle large delta time gracefully
}

func TestInputSystem_Update_NegativeDeltaTime(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity in playing state
	entity := createInputEntity(1, components.StatePlaying, 100, 100)

	entities := []Entity{entity}

	// Run update with negative delta time
	is.Update(-0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// The system should handle negative delta time gracefully
}

func TestInputSystem_Integration(t *testing.T) {
	is := NewInputSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create multiple entities with different states
	playingEntity := createInputEntity(1, components.StatePlaying, 100, 100)
	menuEntity := createInputEntity(2, components.StateMenu, 200, 200)
	gameOverEntity := createInputEntity(3, components.StateGameOver, 300, 300)

	entities := []Entity{playingEntity, menuEntity, gameOverEntity}

	// Run multiple updates
	for i := 0; i < 3; i++ {
		is.Update(0.016, entities, eventDispatcher)
	}

	// Verify menu and game over entities remain unchanged
	menuTransform := menuEntity.GetTransform()
	mx, my := menuTransform.GetX(), menuTransform.GetY()
	if mx != 200 || my != 200 {
		t.Errorf("Expected menu entity position to remain (200, 200), got (%f, %f)", mx, my)
	}

	gameOverTransform := gameOverEntity.GetTransform()
	gx, gy := gameOverTransform.GetX(), gameOverTransform.GetY()
	if gx != 300 || gy != 300 {
		t.Errorf("Expected game over entity position to remain (300, 300), got (%f, %f)", gx, gy)
	}

	// Verify playing entity was processed (though actual position depends on ebiten)
	playingTransform := playingEntity.GetTransform()
	if playingTransform == nil {
		t.Fatal("Expected playing entity transform to still exist")
	}
}

// Helper function to create test entities

func createInputEntity(id uint64, stateType components.StateType, x, y float64) Entity {
	entity := entities.NewEntity(id)
	transform := components.NewTransform(x, y)
	state := components.NewState(stateType)
	entity.AddComponent(transform)
	entity.AddComponent(state)
	return entity
}
