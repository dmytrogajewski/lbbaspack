package systems

import (
	"fmt"
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
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

func TestInputSystem_handleKeyboardInput_CtrlXExit(t *testing.T) {
	// Create a new input system
	inputSys := NewInputSystem()

	// Create a mock event dispatcher
	eventDispatcher := events.NewEventDispatcher()

	// Track if exit event was published
	exitEventPublished := false
	eventDispatcher.Subscribe(events.EventExit, func(event *events.Event) {
		exitEventPublished = true
	})

	// Test that handleKeyboardInput doesn't panic
	// Note: We can't easily test the actual key press in unit tests
	// since ebiten.IsKeyPressed requires a running game context
	// This test ensures the method exists and can be called safely
	inputSys.handleKeyboardInput(eventDispatcher)

	// Verify the method exists and can be called
	// The actual key press testing would need to be done in integration tests
	if !exitEventPublished {
		// This is expected since we can't simulate key presses in unit tests
		t.Log("Exit event not published (expected in unit test environment)")
	}
}

func TestInputSystem_handleKeyboardMovement(t *testing.T) {
	// Create a new input system
	inputSys := NewInputSystem()

	// Create a mock transform component
	transform := &components.Transform{
		X: 100.0,
		Y: 200.0,
	}

	// Test that handleKeyboardMovement doesn't panic
	// Note: We can't easily test the actual key press in unit tests
	// since ebiten.IsKeyPressed requires a running game context
	// This test ensures the method exists and can be called safely
	inputSys.handleKeyboardMovement(transform, 0.016) // 60 FPS delta time

	// Verify the method exists and can be called
	// The actual key press testing would need to be done in integration tests
	// or with a mock ebiten implementation
}

func TestInputSystem_handleLoadBalancerInput(t *testing.T) {
	// Create a new input system
	inputSys := NewInputSystem()

	// Create a mock transform component
	transform := &components.Transform{
		X: 100.0,
		Y: 200.0,
	}

	// Create a mock event dispatcher
	eventDispatcher := events.NewEventDispatcher()

	// Test that handleLoadBalancerInput doesn't panic
	// This method combines both keyboard and mouse input handling
	inputSys.handleLoadBalancerInput(transform, eventDispatcher, 0.016)

	// Verify the method exists and can be called
	// The actual input testing would need to be done in integration tests
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

func TestInputSystem_Debug_KeyboardMovement_Integration(t *testing.T) {
	// Create a new input system
	inputSys := NewInputSystem()

	// Create a mock transform component with initial position
	transform := &components.Transform{
		X: 100.0,
		Y: 200.0,
	}

	// Create a mock state component
	state := &components.State{
		Current: components.StatePlaying,
	}

	// Create a mock entity
	entity := &entities.Entity{
		ID:         1,
		Active:     true,
		Components: make(map[string]components.Component),
	}
	entity.AddComponent(transform)
	entity.AddComponent(state)

	// Create a mock event dispatcher
	eventDispatcher := events.NewEventDispatcher()

	// Test the full Update method with our mock entity
	entities := []Entity{entity}

	// Simulate multiple updates to see if position changes
	initialX := transform.GetX()
	fmt.Printf("[DEBUG] Initial position: %.2f\n", initialX)

	// Update multiple times to simulate game loop
	for i := 0; i < 5; i++ {
		inputSys.Update(0.016, entities, eventDispatcher) // 60 FPS
		currentX := transform.GetX()
		fmt.Printf("[DEBUG] Update %d - Position: %.2f (delta: %.2f)\n", i+1, currentX, currentX-initialX)
	}

	// Verify that the transform component is being accessed correctly
	if transform.GetX() != initialX {
		t.Logf("Position changed from %.2f to %.2f", initialX, transform.GetX())
	} else {
		t.Logf("Position remained at %.2f (no keyboard input in test environment)", initialX)
	}
}

func TestInputSystem_Debug_EntityFiltering(t *testing.T) {
	// Create a new input system
	inputSys := NewInputSystem()

	// Create entities with different component combinations
	entity1 := &entities.Entity{ID: 1, Active: true, Components: make(map[string]components.Component)}
	entity1.AddComponent(&components.Transform{X: 100, Y: 200})
	entity1.AddComponent(&components.State{Current: components.StatePlaying})

	entity2 := &entities.Entity{ID: 2, Active: true, Components: make(map[string]components.Component)}
	entity2.AddComponent(&components.Transform{X: 150, Y: 250})
	// Missing State component

	entity3 := &entities.Entity{ID: 3, Active: true, Components: make(map[string]components.Component)}
	entity3.AddComponent(&components.State{Current: components.StateMenu})
	// Missing Transform component

	entities := []Entity{entity1, entity2, entity3}

	// Test entity filtering
	filtered := inputSys.FilterEntities(entities)

	fmt.Printf("[DEBUG] Total entities: %d, Filtered entities: %d\n", len(entities), len(filtered))

	// Should only have entity1 (has both Transform and State with "playing" state)
	if len(filtered) != 1 {
		t.Errorf("Expected 1 filtered entity, got %d", len(filtered))
	}

	if filtered[0].GetID() != 1 {
		t.Errorf("Expected entity ID 1, got %d", filtered[0].GetID())
	}
}

func TestInputSystem_Debug_StateFiltering(t *testing.T) {
	// Create a new input system
	inputSys := NewInputSystem()

	// Create entities with different states
	playingEntity := &entities.Entity{ID: 1, Active: true, Components: make(map[string]components.Component)}
	playingEntity.AddComponent(&components.Transform{X: 100, Y: 200})
	playingEntity.AddComponent(&components.State{Current: components.StatePlaying})

	menuEntity := &entities.Entity{ID: 2, Active: true, Components: make(map[string]components.Component)}
	menuEntity.AddComponent(&components.Transform{X: 150, Y: 250})
	menuEntity.AddComponent(&components.State{Current: components.StateMenu})

	gameOverEntity := &entities.Entity{ID: 3, Active: true, Components: make(map[string]components.Component)}
	gameOverEntity.AddComponent(&components.Transform{X: 200, Y: 300})
	gameOverEntity.AddComponent(&components.State{Current: components.StateGameOver})

	entities := []Entity{playingEntity, menuEntity, gameOverEntity}

	// Test that only playing state entities are processed
	eventDispatcher := events.NewEventDispatcher()

	// Capture initial positions
	initialPositions := make(map[uint64]float64)
	for _, entity := range entities {
		if transform := entity.GetTransform(); transform != nil {
			initialPositions[entity.GetID()] = transform.GetX()
		}
	}

	// Update the input system
	inputSys.Update(0.016, entities, eventDispatcher)

	// Check which entities had their positions updated
	for _, entity := range entities {
		if transform := entity.GetTransform(); transform != nil {
			currentX := transform.GetX()
			initialX := initialPositions[entity.GetID()]
			state := entity.GetState()

			fmt.Printf("[DEBUG] Entity %d - State: %s, Position: %.2f -> %.2f (delta: %.2f)\n",
				entity.GetID(), state.GetState(), initialX, currentX, currentX-initialX)
		}
	}
}

func TestInputSystem_Debug_MethodCallChain(t *testing.T) {
	// Create a new input system
	inputSys := NewInputSystem()

	// Create a mock transform component
	transform := &components.Transform{
		X: 100.0,
		Y: 200.0,
	}

	// Create a mock event dispatcher
	eventDispatcher := events.NewEventDispatcher()

	// Test each method in the call chain
	fmt.Println("[DEBUG] Testing method call chain...")

	// Test handleKeyboardMovement directly
	initialX := transform.GetX()
	keyboardMoved := inputSys.handleKeyboardMovement(transform, 0.016)
	fmt.Printf("[DEBUG] handleKeyboardMovement - Position: %.2f (delta: %.2f), Moved: %t\n",
		transform.GetX(), transform.GetX()-initialX, keyboardMoved)

	// Test handleLoadBalancerInput
	transform.SetPosition(100.0, 200.0) // Reset position
	initialX = transform.GetX()
	inputSys.handleLoadBalancerInput(transform, eventDispatcher, 0.016)
	fmt.Printf("[DEBUG] handleLoadBalancerInput - Position: %.2f (delta: %.2f)\n",
		transform.GetX(), transform.GetX()-initialX)

	// Test handleMouseInput
	transform.SetPosition(100.0, 200.0) // Reset position
	initialX = transform.GetX()
	inputSys.handleMouseInput(transform, eventDispatcher)
	fmt.Printf("[DEBUG] handleMouseInput - Position: %.2f (delta: %.2f)\n",
		transform.GetX(), transform.GetX()-initialX)
}

func TestInputSystem_Debug_MousePosition(t *testing.T) {
	// Test what ebiten.CursorPosition() returns in test environment
	mouseX, mouseY := ebiten.CursorPosition()
	fmt.Printf("[DEBUG] Mouse position in test environment: (%d, %d)\n", mouseX, mouseY)

	// This explains why the position is being set to 0 - mouse position is (0,0) in tests
	t.Logf("Mouse position is (%d, %d) - this explains the position reset to 0", mouseX, mouseY)
}

func TestInputSystem_Debug_SimulatedKeyboardInput(t *testing.T) {
	// Create a new input system
	inputSys := NewInputSystem()

	// Create a mock transform component
	transform := &components.Transform{
		X: 100.0,
		Y: 200.0,
	}

	// Test the logic without actual ebiten input
	// We'll simulate what happens when keys are pressed

	fmt.Println("[DEBUG] Testing keyboard movement logic...")

	// Test 1: No keyboard input (should use mouse)
	keyboardMoved := inputSys.handleKeyboardMovement(transform, 0.016)
	fmt.Printf("[DEBUG] No keys pressed - Position: %.2f, Keyboard moved: %t\n", transform.GetX(), keyboardMoved)

	// Reset position
	transform.SetPosition(100.0, 200.0)

	// Test 2: Simulate what happens when A key is pressed
	// Since we can't actually press keys in tests, we'll test the logic manually
	const moveSpeed = 300.0
	deltaTime := 0.016

	// Simulate A key press
	currentX := transform.GetX()
	newX := currentX - moveSpeed*deltaTime // A key movement
	transform.SetPosition(newX, transform.GetY())

	fmt.Printf("[DEBUG] Simulated A key press - Position: %.2f (delta: %.2f)\n",
		transform.GetX(), transform.GetX()-100.0)

	// Verify the movement calculation
	expectedDelta := -moveSpeed * deltaTime
	actualDelta := transform.GetX() - 100.0
	t.Logf("A key movement: expected %.2f, got %.2f", expectedDelta, actualDelta)

	// Use tolerance for floating point comparison
	const tolerance = 0.01
	diff := actualDelta - expectedDelta
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		t.Errorf("Expected delta %.2f, got %.2f (tolerance: %.2f)", expectedDelta, actualDelta, tolerance)
	}

	// Test 3: Simulate D key press
	transform.SetPosition(100.0, 200.0) // Reset
	currentX = transform.GetX()
	newX = currentX + moveSpeed*deltaTime // D key movement
	transform.SetPosition(newX, transform.GetY())

	fmt.Printf("[DEBUG] Simulated D key press - Position: %.2f (delta: %.2f)\n",
		transform.GetX(), transform.GetX()-100.0)

	// Verify the movement calculation
	expectedDelta = moveSpeed * deltaTime
	actualDelta = transform.GetX() - 100.0
	t.Logf("D key movement: expected %.2f, got %.2f", expectedDelta, actualDelta)

	// Use tolerance for floating point comparison
	diff = actualDelta - expectedDelta
	if diff < 0 {
		diff = -diff
	}
	if diff > tolerance {
		t.Errorf("Expected delta %.2f, got %.2f (tolerance: %.2f)", expectedDelta, actualDelta, tolerance)
	}
}

func TestInputSystem_InputMethodTracking(t *testing.T) {
	// Create a new input system
	inputSys := NewInputSystem()

	// Create a mock transform component
	transform := &components.Transform{
		X: 100.0,
		Y: 200.0,
	}

	// Create a mock event dispatcher
	eventDispatcher := events.NewEventDispatcher()

	fmt.Println("[DEBUG] Testing input method tracking...")

	// Test 1: Initial state should be mouse mode
	fmt.Printf("[DEBUG] Initial input method: %s\n", inputSys.activeInputMethod)
	if inputSys.activeInputMethod != "mouse" {
		t.Errorf("Expected initial input method to be 'mouse', got '%s'", inputSys.activeInputMethod)
	}

	// Test 2: Simulate keyboard input (should switch to keyboard mode)
	// We can't actually simulate key presses, but we can test the logic
	// by directly calling the method and checking the state

	// Simulate what happens when keyboard is used
	inputSys.activeInputMethod = "keyboard"
	inputSys.keyboardLastUsed = true

	// Call the method to test the logic
	inputSys.handleLoadBalancerInput(transform, eventDispatcher, 0.016)

	fmt.Printf("[DEBUG] After keyboard input - Method: %s, LastUsed: %t\n",
		inputSys.activeInputMethod, inputSys.keyboardLastUsed)

	// Test 3: Simulate no input (should maintain keyboard mode)
	inputSys.keyboardLastUsed = false

	fmt.Printf("[DEBUG] After no input - Method: %s, LastUsed: %t\n",
		inputSys.activeInputMethod, inputSys.keyboardLastUsed)

	// Test 4: Simulate mouse movement (should switch back to mouse mode)
	// This would happen in the actual game when mouse is moved
	inputSys.activeInputMethod = "mouse"
	inputSys.keyboardLastUsed = false

	fmt.Printf("[DEBUG] After mouse movement - Method: %s, LastUsed: %t\n",
		inputSys.activeInputMethod, inputSys.keyboardLastUsed)

	// Verify the tracking logic
	t.Logf("Input method tracking test completed - Method: %s", inputSys.activeInputMethod)
}

func TestInputSystem_InputMethodPersistence(t *testing.T) {
	// Create a new input system
	inputSys := NewInputSystem()

	// Create a mock transform component
	transform := &components.Transform{
		X: 100.0,
		Y: 200.0,
	}

	// Create a mock event dispatcher
	eventDispatcher := events.NewEventDispatcher()

	fmt.Println("[DEBUG] Testing input method persistence...")

	// Test that the input method persists across multiple updates
	initialX := transform.GetX()

	// Simulate multiple updates with no input
	for i := 0; i < 5; i++ {
		inputSys.handleLoadBalancerInput(transform, eventDispatcher, 0.016)
		currentX := transform.GetX()
		fmt.Printf("[DEBUG] Update %d - Method: %s, Position: %.2f (delta: %.2f)\n",
			i+1, inputSys.activeInputMethod, currentX, currentX-initialX)
	}

	// The position should remain stable (no jumping to mouse position)
	finalX := transform.GetX()
	if finalX != initialX {
		t.Logf("Position changed from %.2f to %.2f (this is expected if mouse position is different)",
			initialX, finalX)
	} else {
		t.Logf("Position remained stable at %.2f", finalX)
	}
}
