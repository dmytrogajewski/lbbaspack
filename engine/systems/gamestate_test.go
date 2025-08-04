package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"
)

func TestNewGameStateSystem(t *testing.T) {
	gss := NewGameStateSystem()

	// Test that the system is properly initialized
	if gss == nil {
		t.Fatal("NewGameStateSystem returned nil")
	}

	// Test required components
	expectedComponents := []string{"State"}
	if len(gss.RequiredComponents) != len(expectedComponents) {
		t.Errorf("Expected %d required components, got %d", len(expectedComponents), len(gss.RequiredComponents))
	}

	for i, component := range expectedComponents {
		if gss.RequiredComponents[i] != component {
			t.Errorf("Expected required component %s at index %d, got %s", component, i, gss.RequiredComponents[i])
		}
	}

	// Test initial values
	if gss.currentState != components.StateMenu {
		t.Errorf("Expected initial currentState to be StateMenu, got %v", gss.currentState)
	}

	if gss.gameTime != 0.0 {
		t.Errorf("Expected initial gameTime to be 0.0, got %f", gss.gameTime)
	}

	if gss.score != 0 {
		t.Errorf("Expected initial score to be 0, got %d", gss.score)
	}

	if gss.level != 1 {
		t.Errorf("Expected initial level to be 1, got %d", gss.level)
	}

	if gss.lastLevelUpTime != 0.0 {
		t.Errorf("Expected initial lastLevelUpTime to be 0.0, got %f", gss.lastLevelUpTime)
	}
}

func TestGameStateSystem_Update_NoEntities(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test with no entities
	entities := []Entity{}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify game time increased
	if gss.gameTime != 0.016 {
		t.Errorf("Expected gameTime to be 0.016, got %f", gss.gameTime)
	}

	// Verify state remains unchanged
	if gss.currentState != components.StateMenu {
		t.Errorf("Expected currentState to remain StateMenu, got %v", gss.currentState)
	}
}

func TestGameStateSystem_Update_WithEntities(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with state component
	entity := createStateEntity(1, components.StateMenu)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify game time increased
	if gss.gameTime != 0.016 {
		t.Errorf("Expected gameTime to be 0.016, got %f", gss.gameTime)
	}

	// Verify state component was updated
	stateComp := entity.GetState()
	if stateComp == nil {
		t.Fatal("Expected state component to exist")
	}

	state := stateComp
	if state.GetState() != "menu" {
		t.Errorf("Expected state to be 'menu', got %s", state.GetState())
	}
}

func TestGameStateSystem_Update_PlayingState(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Set to playing state
	gss.currentState = components.StatePlaying
	gss.gameTime = 10.0

	// Create entity with state component
	entity := createStateEntity(1, components.StatePlaying)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify game time increased
	if gss.gameTime != 10.016 {
		t.Errorf("Expected gameTime to be 10.016, got %f", gss.gameTime)
	}

	// Verify state component was updated
	stateComp := entity.GetState()
	state := stateComp
	if state.GetState() != "playing" {
		t.Errorf("Expected state to be 'playing', got %s", state.GetState())
	}
}

func TestGameStateSystem_Update_GameOverState(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Set to game over state
	gss.currentState = components.StateGameOver
	gss.gameTime = 20.0

	// Create entity with state component
	entity := createStateEntity(1, components.StateGameOver)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify game time increased
	if gss.gameTime != 20.016 {
		t.Errorf("Expected gameTime to be 20.016, got %f", gss.gameTime)
	}

	// Verify state component was updated
	stateComp := entity.GetState()
	state := stateComp
	if state.GetState() != "gameover" {
		t.Errorf("Expected state to be 'gameover', got %s", state.GetState())
	}
}

func TestGameStateSystem_getStateString(t *testing.T) {
	gss := NewGameStateSystem()

	// Test menu state
	gss.currentState = components.StateMenu
	if gss.getStateString() != "menu" {
		t.Errorf("Expected state string to be 'menu', got %s", gss.getStateString())
	}

	// Test playing state
	gss.currentState = components.StatePlaying
	if gss.getStateString() != "playing" {
		t.Errorf("Expected state string to be 'playing', got %s", gss.getStateString())
	}

	// Test game over state
	gss.currentState = components.StateGameOver
	if gss.getStateString() != "gameover" {
		t.Errorf("Expected state string to be 'gameover', got %s", gss.getStateString())
	}

	// Test unknown state
	gss.currentState = 999 // Invalid state
	if gss.getStateString() != "unknown" {
		t.Errorf("Expected state string to be 'unknown', got %s", gss.getStateString())
	}
}

func TestGameStateSystem_Initialize(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Verify initial state remains unchanged
	if gss.currentState != components.StateMenu {
		t.Errorf("Expected currentState to remain StateMenu, got %v", gss.currentState)
	}

	if gss.score != 0 {
		t.Errorf("Expected score to remain 0, got %d", gss.score)
	}

	if gss.level != 1 {
		t.Errorf("Expected level to remain 1, got %d", gss.level)
	}
}

func TestGameStateSystem_EventHandling_GameStart(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Verify initial state
	if gss.currentState != components.StateMenu {
		t.Errorf("Expected initial currentState to be StateMenu, got %v", gss.currentState)
	}

	// Publish game start event
	event := events.NewEvent(events.EventGameStart, nil)
	eventDispatcher.Publish(event)

	// Verify transition to playing state
	if gss.currentState != components.StatePlaying {
		t.Errorf("Expected currentState to be StatePlaying, got %v", gss.currentState)
	}

	// Verify game state was reset
	if gss.gameTime != 0.0 {
		t.Errorf("Expected gameTime to be reset to 0.0, got %f", gss.gameTime)
	}

	if gss.score != 0 {
		t.Errorf("Expected score to be reset to 0, got %d", gss.score)
	}

	if gss.level != 1 {
		t.Errorf("Expected level to be reset to 1, got %d", gss.level)
	}
}

func TestGameStateSystem_EventHandling_GameStart_AlreadyPlaying(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Set to playing state
	gss.currentState = components.StatePlaying
	gss.gameTime = 10.0
	gss.score = 50
	gss.level = 3

	// Publish game start event
	event := events.NewEvent(events.EventGameStart, nil)
	eventDispatcher.Publish(event)

	// Verify state remains playing (no transition from playing to playing)
	if gss.currentState != components.StatePlaying {
		t.Errorf("Expected currentState to remain StatePlaying, got %v", gss.currentState)
	}

	// Verify game state was not reset
	if gss.gameTime != 10.0 {
		t.Errorf("Expected gameTime to remain 10.0, got %f", gss.gameTime)
	}

	if gss.score != 50 {
		t.Errorf("Expected score to remain 50, got %d", gss.score)
	}

	if gss.level != 3 {
		t.Errorf("Expected level to remain 3, got %d", gss.level)
	}
}

func TestGameStateSystem_EventHandling_GameOver(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Set to playing state
	gss.currentState = components.StatePlaying
	gss.gameTime = 15.5
	gss.score = 75
	gss.level = 2

	// Publish game over event
	event := events.NewEvent(events.EventGameOver, nil)
	eventDispatcher.Publish(event)

	// Verify transition to game over state
	if gss.currentState != components.StateGameOver {
		t.Errorf("Expected currentState to be StateGameOver, got %v", gss.currentState)
	}

	// Verify game state was preserved
	if gss.gameTime != 15.5 {
		t.Errorf("Expected gameTime to remain 15.5, got %f", gss.gameTime)
	}

	if gss.score != 75 {
		t.Errorf("Expected score to remain 75, got %d", gss.score)
	}

	if gss.level != 2 {
		t.Errorf("Expected level to remain 2, got %d", gss.level)
	}
}

func TestGameStateSystem_EventHandling_GameOver_NotPlaying(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Set to menu state (not playing)
	gss.currentState = components.StateMenu

	// Publish game over event
	event := events.NewEvent(events.EventGameOver, nil)
	eventDispatcher.Publish(event)

	// Verify state remains menu (no transition from menu to game over)
	if gss.currentState != components.StateMenu {
		t.Errorf("Expected currentState to remain StateMenu, got %v", gss.currentState)
	}
}

func TestGameStateSystem_EventHandling_PacketCaught_Playing(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Set to playing state
	gss.currentState = components.StatePlaying
	gss.score = 25

	// Publish packet caught event
	event := events.NewEvent(events.EventPacketCaught, nil)
	eventDispatcher.Publish(event)

	// Verify score increased
	if gss.score != 35 {
		t.Errorf("Expected score to be 35, got %d", gss.score)
	}
}

func TestGameStateSystem_EventHandling_PacketCaught_NotPlaying(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Set to menu state (not playing)
	gss.currentState = components.StateMenu
	gss.score = 25

	// Publish packet caught event
	event := events.NewEvent(events.EventPacketCaught, nil)
	eventDispatcher.Publish(event)

	// Verify score did not increase
	if gss.score != 25 {
		t.Errorf("Expected score to remain 25, got %d", gss.score)
	}
}

func TestGameStateSystem_checkLevelUp_ScoreBased(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Set to playing state
	gss.currentState = components.StatePlaying
	gss.gameTime = 10.0
	gss.lastLevelUpTime = 5.0

	// Set score to 90, so when packet is caught (+10), it becomes 100 and triggers level up
	gss.score = 90

	// Publish packet caught event
	event := events.NewEvent(events.EventPacketCaught, nil)
	eventDispatcher.Publish(event)

	// Verify level increased
	if gss.level != 2 {
		t.Errorf("Expected level to be 2, got %d", gss.level)
	}

	// Verify last level up time was updated
	if gss.lastLevelUpTime != 10.0 {
		t.Errorf("Expected lastLevelUpTime to be 10.0, got %f", gss.lastLevelUpTime)
	}
}

func TestGameStateSystem_checkLevelUp_ScoreBased_TooSoon(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Set to playing state
	gss.currentState = components.StatePlaying
	gss.gameTime = 10.0
	gss.lastLevelUpTime = 9.5 // Less than 1 second ago

	// Set score to trigger level up (100 points)
	gss.score = 100

	// Publish packet caught event
	event := events.NewEvent(events.EventPacketCaught, nil)
	eventDispatcher.Publish(event)

	// Verify level did not increase (too soon)
	if gss.level != 1 {
		t.Errorf("Expected level to remain 1, got %d", gss.level)
	}
}

func TestGameStateSystem_checkLevelUp_ScoreBased_ZeroScore(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Set to playing state
	gss.currentState = components.StatePlaying
	gss.gameTime = 10.0
	gss.lastLevelUpTime = 5.0

	// Set score to 0 (should not trigger level up)
	gss.score = 0

	// Publish packet caught event
	event := events.NewEvent(events.EventPacketCaught, nil)
	eventDispatcher.Publish(event)

	// Verify level did not increase
	if gss.level != 1 {
		t.Errorf("Expected level to remain 1, got %d", gss.level)
	}
}

func TestGameStateSystem_updatePlayingState_TimeBasedLevelUp(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Set to playing state
	gss.currentState = components.StatePlaying
	gss.gameTime = 30.0        // Exactly 30 seconds
	gss.lastLevelUpTime = 25.0 // More than 1 second ago

	// Create entity with state component
	entity := createStateEntity(1, components.StatePlaying)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify level increased
	if gss.level != 2 {
		t.Errorf("Expected level to be 2, got %d", gss.level)
	}

	// Verify last level up time was updated
	if gss.lastLevelUpTime != 30.016 {
		t.Errorf("Expected lastLevelUpTime to be 30.016, got %f", gss.lastLevelUpTime)
	}
}

func TestGameStateSystem_updatePlayingState_TimeBasedLevelUp_TooSoon(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Set to playing state
	gss.currentState = components.StatePlaying
	gss.gameTime = 30.0        // Exactly 30 seconds
	gss.lastLevelUpTime = 29.5 // Less than 1 second ago

	// Create entity with state component
	entity := createStateEntity(1, components.StatePlaying)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify level did not increase (too soon)
	if gss.level != 1 {
		t.Errorf("Expected level to remain 1, got %d", gss.level)
	}
}

func TestGameStateSystem_updatePlayingState_TimeBasedLevelUp_NotThirtySeconds(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Set to playing state
	gss.currentState = components.StatePlaying
	gss.gameTime = 25.0 // Not 30 seconds
	gss.lastLevelUpTime = 20.0

	// Create entity with state component
	entity := createStateEntity(1, components.StatePlaying)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify level did not increase
	if gss.level != 1 {
		t.Errorf("Expected level to remain 1, got %d", gss.level)
	}
}

func TestGameStateSystem_GetCurrentState(t *testing.T) {
	gss := NewGameStateSystem()

	// Test menu state
	gss.currentState = components.StateMenu
	if gss.GetCurrentState() != components.StateMenu {
		t.Errorf("Expected GetCurrentState to return StateMenu, got %v", gss.GetCurrentState())
	}

	// Test playing state
	gss.currentState = components.StatePlaying
	if gss.GetCurrentState() != components.StatePlaying {
		t.Errorf("Expected GetCurrentState to return StatePlaying, got %v", gss.GetCurrentState())
	}

	// Test game over state
	gss.currentState = components.StateGameOver
	if gss.GetCurrentState() != components.StateGameOver {
		t.Errorf("Expected GetCurrentState to return StateGameOver, got %v", gss.GetCurrentState())
	}
}

func TestGameStateSystem_GetScore(t *testing.T) {
	gss := NewGameStateSystem()

	// Set score
	gss.score = 150

	// Get score
	result := gss.GetScore()

	// Verify result
	if result != 150 {
		t.Errorf("Expected GetScore to return 150, got %d", result)
	}
}

func TestGameStateSystem_GetLevel(t *testing.T) {
	gss := NewGameStateSystem()

	// Set level
	gss.level = 5

	// Get level
	result := gss.GetLevel()

	// Verify result
	if result != 5 {
		t.Errorf("Expected GetLevel to return 5, got %d", result)
	}
}

func TestGameStateSystem_GetGameTime(t *testing.T) {
	gss := NewGameStateSystem()

	// Set game time
	gss.gameTime = 45.7

	// Get game time
	result := gss.GetGameTime()

	// Verify result
	if result != 45.7 {
		t.Errorf("Expected GetGameTime to return 45.7, got %f", result)
	}
}

func TestGameStateSystem_Integration(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Create entity with state component
	entity := createStateEntity(1, components.StateMenu)
	entities := []Entity{entity}

	// Verify initial state
	if gss.GetCurrentState() != components.StateMenu {
		t.Errorf("Expected initial state to be StateMenu, got %v", gss.GetCurrentState())
	}

	// Start game
	startEvent := events.NewEvent(events.EventGameStart, nil)
	eventDispatcher.Publish(startEvent)

	// Verify transition to playing
	if gss.GetCurrentState() != components.StatePlaying {
		t.Errorf("Expected state to be StatePlaying after game start, got %v", gss.GetCurrentState())
	}

	// Simulate some gameplay
	for i := 0; i < 5; i++ {
		packetEvent := events.NewEvent(events.EventPacketCaught, nil)
		eventDispatcher.Publish(packetEvent)
		gss.Update(0.016, entities, eventDispatcher)
	}

	// Verify score increased
	if gss.GetScore() != 50 {
		t.Errorf("Expected score to be 50, got %d", gss.GetScore())
	}

	// Simulate time passing to trigger level up
	gss.gameTime = 30.0
	gss.Update(0.016, entities, eventDispatcher)

	// Verify level increased
	if gss.GetLevel() != 2 {
		t.Errorf("Expected level to be 2, got %d", gss.GetLevel())
	}

	// End game
	gameOverEvent := events.NewEvent(events.EventGameOver, nil)
	eventDispatcher.Publish(gameOverEvent)

	// Verify transition to game over
	if gss.GetCurrentState() != components.StateGameOver {
		t.Errorf("Expected state to be StateGameOver after game over, got %v", gss.GetCurrentState())
	}
}

func TestGameStateSystem_EntityWithoutStateComponent(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity without state component
	entity := entities.NewEntity(1)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify game time increased
	if gss.gameTime != 0.016 {
		t.Errorf("Expected gameTime to be 0.016, got %f", gss.gameTime)
	}

	// Verify state remains unchanged
	if gss.currentState != components.StateMenu {
		t.Errorf("Expected currentState to remain StateMenu, got %v", gss.currentState)
	}
}

// Helper function to create test entities

func createStateEntity(id uint64, initialState components.StateType) Entity {
	entity := entities.NewEntity(id)
	state := components.NewState(initialState)
	entity.AddComponent(state)
	return entity
}
