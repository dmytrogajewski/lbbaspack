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
	expectedComponents := []string{"State", "GameSession"}
	if len(gss.RequiredComponents) != len(expectedComponents) {
		t.Errorf("Expected %d required components, got %d", len(expectedComponents), len(gss.RequiredComponents))
	}

	for i, component := range expectedComponents {
		if gss.RequiredComponents[i] != component {
			t.Errorf("Expected required component %s at index %d, got %s", component, i, gss.RequiredComponents[i])
		}
	}
}

func TestGameStateSystem_Update_NoEntities(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test with no entities
	entities := []Entity{}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// System is stateless, so no internal state to verify
}

func TestGameStateSystem_Update_WithEntities(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with state and game session components
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 10.0
	gameSession.Score = 50
	gameSession.Level = 3
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify game session was updated
	if gameSession.GameTime != 10.016 {
		t.Errorf("Expected game time to be 10.016, got %f", gameSession.GameTime)
	}

	// Verify other values remain unchanged
	if gameSession.Score != 50 {
		t.Errorf("Expected score to remain 50, got %d", gameSession.Score)
	}

	if gameSession.Level != 3 {
		t.Errorf("Expected level to remain 3, got %d", gameSession.Level)
	}
}

func TestGameStateSystem_Update_PlayingState(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with playing state
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 10.0
	gameSession.Score = 50
	gameSession.Level = 3
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify game session was updated
	if gameSession.GameTime != 10.016 {
		t.Errorf("Expected game time to be 10.016, got %f", gameSession.GameTime)
	}

	// Verify state remains playing
	if entity.GetState().GetState() != "playing" {
		t.Errorf("Expected state to remain playing, got %s", entity.GetState().GetState())
	}
}

func TestGameStateSystem_Update_GameOverState(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with game over state
	entity := createStateEntity(1, components.StateGameOver)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 10.0
	gameSession.Score = 50
	gameSession.Level = 3
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify game session was updated
	if gameSession.GameTime != 10.016 {
		t.Errorf("Expected game time to be 10.016, got %f", gameSession.GameTime)
	}

	// Verify state remains game over
	if entity.GetState().GetState() != "gameover" {
		t.Errorf("Expected state to remain gameover, got %s", entity.GetState().GetState())
	}
}

func TestGameStateSystem_getStateString(t *testing.T) {
	gss := NewGameStateSystem()

	// Test the helper function
	result := gss.getStateString()
	if result != "menu" {
		t.Errorf("Expected getStateString to return 'menu', got %s", result)
	}
}

func TestGameStateSystem_Initialize(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Verify no errors occurred
	// System is stateless, so no internal state to verify
}

func TestGameStateSystem_EventHandling_GameStart(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Create entity with menu state
	entity := createStateEntity(1, components.StateMenu)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 10.0
	gameSession.Score = 50
	gameSession.Level = 3
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Publish game start event
	event := events.NewEvent(events.EventGameStart, nil)
	eventDispatcher.Publish(event)

	// Run update to process any state changes
	gss.Update(0.016, entities, eventDispatcher)

	// Verify state remains menu (no transition from menu to playing via event)
	if entity.GetState().GetState() != "menu" {
		t.Errorf("Expected state to remain menu, got %s", entity.GetState().GetState())
	}

	// Verify game session values remain unchanged
	if gameSession.GameTime != 10.016 {
		t.Errorf("Expected game time to be 10.016, got %f", gameSession.GameTime)
	}

	if gameSession.Score != 50 {
		t.Errorf("Expected score to remain 50, got %d", gameSession.Score)
	}

	if gameSession.Level != 3 {
		t.Errorf("Expected level to remain 3, got %d", gameSession.Level)
	}
}

func TestGameStateSystem_EventHandling_GameStart_AlreadyPlaying(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Create entity with playing state
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 10.0
	gameSession.Score = 50
	gameSession.Level = 3
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Publish game start event
	event := events.NewEvent(events.EventGameStart, nil)
	eventDispatcher.Publish(event)

	// Run update to process any state changes
	gss.Update(0.016, entities, eventDispatcher)

	// Verify state remains playing (no transition from playing to playing)
	if entity.GetState().GetState() != "playing" {
		t.Errorf("Expected state to remain playing, got %s", entity.GetState().GetState())
	}

	// Verify game session values remain unchanged
	if gameSession.GameTime != 10.016 {
		t.Errorf("Expected game time to be 10.016, got %f", gameSession.GameTime)
	}

	if gameSession.Score != 50 {
		t.Errorf("Expected score to remain 50, got %d", gameSession.Score)
	}

	if gameSession.Level != 3 {
		t.Errorf("Expected level to remain 3, got %d", gameSession.Level)
	}
}

func TestGameStateSystem_EventHandling_GameOver(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Create entity with playing state
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 15.5
	gameSession.Score = 75
	gameSession.Level = 2
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Publish game over event
	event := events.NewEvent(events.EventGameOver, nil)
	eventDispatcher.Publish(event)

	// Run update to process any state changes
	gss.Update(0.016, entities, eventDispatcher)

	// Verify state remains playing (no automatic transition via event)
	if entity.GetState().GetState() != "playing" {
		t.Errorf("Expected state to remain playing, got %s", entity.GetState().GetState())
	}

	// Verify game session values remain unchanged
	if gameSession.GameTime != 15.516 {
		t.Errorf("Expected game time to be 15.516, got %f", gameSession.GameTime)
	}

	if gameSession.Score != 75 {
		t.Errorf("Expected score to remain 75, got %d", gameSession.Score)
	}

	if gameSession.Level != 2 {
		t.Errorf("Expected level to remain 2, got %d", gameSession.Level)
	}
}

func TestGameStateSystem_EventHandling_GameOver_NotPlaying(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Create entity with menu state
	entity := createStateEntity(1, components.StateMenu)
	gameSession := components.NewGameSession()
	gameSession.Score = 25
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Publish game over event
	event := events.NewEvent(events.EventGameOver, nil)
	eventDispatcher.Publish(event)

	// Run update to process any state changes
	gss.Update(0.016, entities, eventDispatcher)

	// Verify state remains menu (no transition from menu to game over)
	if entity.GetState().GetState() != "menu" {
		t.Errorf("Expected state to remain menu, got %s", entity.GetState().GetState())
	}

	// Verify score remains unchanged
	if gameSession.Score != 25 {
		t.Errorf("Expected score to remain 25, got %d", gameSession.Score)
	}
}

func TestGameStateSystem_EventHandling_PacketCaught_Playing(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Create entity with playing state
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.Score = 25
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Publish packet caught event
	event := events.NewEvent(events.EventPacketCaught, nil)
	eventDispatcher.Publish(event)

	// Run update to process any state changes
	gss.Update(0.016, entities, eventDispatcher)

	// Verify score remains unchanged (packet caught doesn't affect score in this system)
	if gameSession.Score != 25 {
		t.Errorf("Expected score to remain 25, got %d", gameSession.Score)
	}
}

func TestGameStateSystem_EventHandling_PacketCaught_NotPlaying(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Create entity with menu state
	entity := createStateEntity(1, components.StateMenu)
	gameSession := components.NewGameSession()
	gameSession.Score = 25
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Publish packet caught event
	event := events.NewEvent(events.EventPacketCaught, nil)
	eventDispatcher.Publish(event)

	// Run update to process any state changes
	gss.Update(0.016, entities, eventDispatcher)

	// Verify score remains unchanged
	if gameSession.Score != 25 {
		t.Errorf("Expected score to remain 25, got %d", gameSession.Score)
	}
}

func TestGameStateSystem_checkLevelUp_ScoreBased(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with playing state
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 10.0
	gameSession.Score = 100 // Exactly 100, would trigger level up if score-based checks were implemented
	gameSession.Level = 1
	gameSession.LastLevelUpTime = 5.0 // More than 1 second ago
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update - note: current implementation only checks time-based level ups, not score-based
	gss.Update(0.016, entities, eventDispatcher)

	// Verify level was NOT increased (since score-based level ups are not implemented)
	if gameSession.Level != 1 {
		t.Errorf("Expected level to remain 1, got %d", gameSession.Level)
	}

	// Verify last level up time was NOT updated (since no level up occurred)
	if gameSession.LastLevelUpTime != 5.0 {
		t.Errorf("Expected last level up time to remain 5.0, got %f", gameSession.LastLevelUpTime)
	}

	// Verify game time was updated (this always happens)
	expectedGameTime := 10.0 + 0.016
	if gameSession.GameTime != expectedGameTime {
		t.Errorf("Expected game time to be %f, got %f", expectedGameTime, gameSession.GameTime)
	}
}

func TestGameStateSystem_checkLevelUp_ScoreBased_TooSoon(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with playing state
	gameSession := components.NewGameSession()
	gameSession.GameTime = 10.0
	gameSession.Score = 100
	gameSession.Level = 1
	gameSession.LastLevelUpTime = 9.5 // Less than 1 second ago
	entity := createStateEntity(1, components.StatePlaying)
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update - note: current implementation only checks time-based level ups, not score-based
	gss.Update(0.016, entities, eventDispatcher)

	// Verify level was not increased (since score-based level ups are not implemented)
	if gameSession.Level != 1 {
		t.Errorf("Expected level to remain 1, got %d", gameSession.Level)
	}

	// Verify last level up time was not updated (since no level up occurred)
	if gameSession.LastLevelUpTime != 9.5 {
		t.Errorf("Expected last level up time to remain 9.5, got %f", gameSession.LastLevelUpTime)
	}

	// Verify game time was updated (this always happens)
	expectedGameTime := 10.0 + 0.016
	if gameSession.GameTime != expectedGameTime {
		t.Errorf("Expected game time to be %f, got %f", expectedGameTime, gameSession.GameTime)
	}
}

func TestGameStateSystem_checkLevelUp_ScoreBased_ZeroScore(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with playing state
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 10.0
	gameSession.Score = 0 // Zero score, would not trigger level up even if score-based checks were implemented
	gameSession.Level = 1
	gameSession.LastLevelUpTime = 5.0
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update - note: current implementation only checks time-based level ups, not score-based
	gss.Update(0.016, entities, eventDispatcher)

	// Verify level was not increased (since score-based level ups are not implemented)
	if gameSession.Level != 1 {
		t.Errorf("Expected level to remain 1, got %d", gameSession.Level)
	}

	// Verify game time was updated (this always happens)
	expectedGameTime := 10.0 + 0.016
	if gameSession.GameTime != expectedGameTime {
		t.Errorf("Expected game time to be %f, got %f", expectedGameTime, gameSession.GameTime)
	}
}

func TestGameStateSystem_updatePlayingState_TimeBasedLevelUp(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with playing state
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 30.0 // Exactly 30 seconds
	gameSession.Score = 50
	gameSession.Level = 1
	gameSession.LastLevelUpTime = 5.0 // More than 1 second ago
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update to trigger time-based level up
	gss.Update(0.016, entities, eventDispatcher)

	// Verify level was increased
	if gameSession.Level != 2 {
		t.Errorf("Expected level to be 2, got %d", gameSession.Level)
	}

	// Verify last level up time was updated (should be the new GameTime after delta time addition)
	expectedLastLevelUpTime := 30.0 + 0.016
	if gameSession.LastLevelUpTime != expectedLastLevelUpTime {
		t.Errorf("Expected last level up time to be %f, got %f", expectedLastLevelUpTime, gameSession.LastLevelUpTime)
	}
}

func TestGameStateSystem_updatePlayingState_TimeBasedLevelUp_TooSoon(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with playing state
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 30.0 // Exactly 30 seconds
	gameSession.Score = 50
	gameSession.Level = 1
	gameSession.LastLevelUpTime = 29.5 // Less than 1 second ago
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify level was not increased
	if gameSession.Level != 1 {
		t.Errorf("Expected level to remain 1, got %d", gameSession.Level)
	}

	// Verify last level up time was not updated
	if gameSession.LastLevelUpTime != 29.5 {
		t.Errorf("Expected last level up time to remain 29.5, got %f", gameSession.LastLevelUpTime)
	}
}

func TestGameStateSystem_updatePlayingState_TimeBasedLevelUp_NotThirtySeconds(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with playing state
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 25.0 // Not 30 seconds
	gameSession.Score = 50
	gameSession.Level = 1
	gameSession.LastLevelUpTime = 5.0
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify level was not increased
	if gameSession.Level != 1 {
		t.Errorf("Expected level to remain 1, got %d", gameSession.Level)
	}
}

func TestGameStateSystem_GetCurrentState(t *testing.T) {
	gss := NewGameStateSystem()

	// Test with menu state entity
	menuEntity := createStateEntity(1, components.StateMenu)
	entities := []Entity{menuEntity}

	// Run update to determine state
	gss.Update(0.016, entities, nil)

	// Verify state is determined from components
	if menuEntity.GetState().GetState() != "menu" {
		t.Errorf("Expected state to be menu, got %s", menuEntity.GetState().GetState())
	}

	// Test with playing state entity
	playingEntity := createStateEntity(2, components.StatePlaying)
	entities = []Entity{playingEntity}

	// Run update to determine state
	gss.Update(0.016, entities, nil)

	// Verify state is determined from components
	if playingEntity.GetState().GetState() != "playing" {
		t.Errorf("Expected state to be playing, got %s", playingEntity.GetState().GetState())
	}
}

func TestGameStateSystem_GetScore(t *testing.T) {
	gss := NewGameStateSystem()

	// Create entity with game session
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.Score = 150
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, nil)

	// Verify score is accessible through game session component
	if gameSession.Score != 150 {
		t.Errorf("Expected score to be 150, got %d", gameSession.Score)
	}
}

func TestGameStateSystem_GetLevel(t *testing.T) {
	gss := NewGameStateSystem()

	// Create entity with game session
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.Level = 5
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, nil)

	// Verify level is accessible through game session component
	if gameSession.Level != 5 {
		t.Errorf("Expected level to be 5, got %d", gameSession.Level)
	}
}

func TestGameStateSystem_GetGameTime(t *testing.T) {
	gss := NewGameStateSystem()

	// Create entity with game session
	entity := createStateEntity(1, components.StatePlaying)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 45.7
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, nil)

	// Verify game time is accessible through game session component
	if gameSession.GameTime != 45.716 {
		t.Errorf("Expected game time to be 45.716, got %f", gameSession.GameTime)
	}
}

func TestGameStateSystem_Integration(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	gss.Initialize(eventDispatcher)

	// Create entity with initial state
	entity := createStateEntity(1, components.StateMenu)
	gameSession := components.NewGameSession()
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Simulate game loop with state transitions
	// First update - should be in menu state
	gss.Update(0.5, entities, eventDispatcher)

	// Verify initial state
	if entity.GetState().GetState() != "menu" {
		t.Errorf("Expected initial state to be menu, got %s", entity.GetState().GetState())
	}

	// Change to playing state
	entity.GetState().SetState("playing")

	// Update system
	gss.Update(0.5, entities, eventDispatcher)

	// Verify state change
	if entity.GetState().GetState() != "playing" {
		t.Errorf("Expected state to be playing, got %s", entity.GetState().GetState())
	}

	// Simulate level up conditions
	gameSession.Score = 100
	gameSession.GameTime = 30.0
	gameSession.LastLevelUpTime = 5.0

	// Update system to trigger level up
	gss.Update(0.016, entities, eventDispatcher)

	// Verify level up occurred
	if gameSession.Level != 2 {
		t.Errorf("Expected level to be 2, got %d", gameSession.Level)
	}

	// Change to game over state
	entity.GetState().SetState("gameover")

	// Update system
	gss.Update(0.5, entities, eventDispatcher)

	// Verify state change
	if entity.GetState().GetState() != "gameover" {
		t.Errorf("Expected state to be gameover, got %s", entity.GetState().GetState())
	}

	// Verify game session was updated (should be the previous time plus the last delta time)
	expectedGameTime := 30.0 + 0.5 + 0.016
	if gameSession.GameTime != expectedGameTime {
		t.Errorf("Expected game time to be %f, got %f", expectedGameTime, gameSession.GameTime)
	}
}

func TestGameStateSystem_EntityWithoutStateComponent(t *testing.T) {
	gss := NewGameStateSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity without state component
	entity := entities.NewEntity(1)
	gameSession := components.NewGameSession()
	gameSession.GameTime = 0.0
	entity.AddComponent(gameSession)
	entities := []Entity{entity}

	// Run update
	gss.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// System should handle entities without state components gracefully
	if gameSession.GameTime != 0.016 {
		t.Errorf("Expected game time to be 0.016, got %f", gameSession.GameTime)
	}
}

// Helper function to create test entities
func createStateEntity(id uint64, initialState components.StateType) Entity {
	entity := entities.NewEntity(id)
	state := components.NewState(initialState)
	entity.AddComponent(state)
	return entity
}
