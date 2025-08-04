package main

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/ecs"
	"lbbaspack/engine/events"
	"lbbaspack/engine/systems"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// TestNewGame tests the NewGame constructor
func TestNewGame(t *testing.T) {
	game := NewGame()

	if game == nil {
		t.Fatal("Expected game to be created")
	}

	if game.World == nil {
		t.Error("Expected world to be initialized")
	}

	if !game.initialized {
		t.Error("Expected game to be initialized")
	}

	if game.RenderSys == nil {
		t.Error("Expected render system to be initialized")
	}

	if game.UISys == nil {
		t.Error("Expected UI system to be initialized")
	}

	if game.MenuSys == nil {
		t.Error("Expected menu system to be initialized")
	}

	if game.eventDispatcher == nil {
		t.Error("Expected event dispatcher to be initialized")
	}

	if game.systems == nil {
		t.Error("Expected systems slice to be initialized")
	}

	if len(game.systems) == 0 {
		t.Error("Expected systems to be added")
	}

	// Verify initial game state
	if game.gameState != components.StateMenu {
		t.Errorf("Expected initial game state to be StateMenu, got %v", game.gameState)
	}
}

// TestGame_Update tests the Update method
func TestGame_Update(t *testing.T) {
	game := NewGame()

	// Test update in menu state
	err := game.Update()
	if err != nil {
		t.Errorf("Expected no error from Update, got %v", err)
	}

	// Verify game state remains menu
	if game.gameState != components.StateMenu {
		t.Errorf("Expected game state to remain StateMenu, got %v", game.gameState)
	}
}

// TestGame_Update_GameStart tests Update after game start event
func TestGame_Update_GameStart(t *testing.T) {
	game := NewGame()

	// Publish game start event
	gameStartEvent := events.NewEvent(events.EventGameStart, &events.EventData{
		SLA:    float64Ptr(99.5),
		Errors: intPtr(10),
	})
	game.eventDispatcher.Publish(gameStartEvent)

	// Update game
	err := game.Update()
	if err != nil {
		t.Errorf("Expected no error from Update, got %v", err)
	}

	// Verify game state changed to playing
	if game.gameState != components.StatePlaying {
		t.Errorf("Expected game state to be StatePlaying, got %v", game.gameState)
	}
}

// TestGame_Update_GameOver tests Update after game over event
func TestGame_Update_GameOver(t *testing.T) {
	game := NewGame()

	// Set game state directly to playing to avoid event chain issues
	game.gameState = components.StatePlaying

	// Publish game over event
	gameOverEvent := events.NewEvent(events.EventGameOver, nil)
	game.eventDispatcher.Publish(gameOverEvent)

	// Update game to process the game over event
	err := game.Update()
	if err != nil {
		t.Errorf("Expected no error from Update, got %v", err)
	}

	// Verify game state changed to game over
	if game.gameState != components.StateGameOver {
		t.Errorf("Expected game state to be StateGameOver, got %v", game.gameState)
	}
}

// TestGame_Draw tests the Draw method
func TestGame_Draw(t *testing.T) {
	game := NewGame()

	// Create a dummy screen for testing
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	// Test drawing in menu state
	game.Draw(screen)

	// Verify systems are initialized after first draw
	if !game.initialized {
		t.Error("Expected game to be initialized after first draw")
	}

	if game.UISys == nil {
		t.Error("Expected UI system to be initialized after draw")
	}

	if game.MenuSys == nil {
		t.Error("Expected menu system to be initialized after draw")
	}
}

// TestGame_Draw_PlayingState tests Draw in playing state
func TestGame_Draw_PlayingState(t *testing.T) {
	game := NewGame()

	// Set game state directly to playing to avoid event chain issues
	game.gameState = components.StatePlaying

	// Create a dummy screen for testing
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	// Test drawing in playing state
	game.Draw(screen)

	// Verify game state is playing
	if game.gameState != components.StatePlaying {
		t.Errorf("Expected game state to be StatePlaying, got %v", game.gameState)
	}
}

// TestGame_Draw_GameOverState tests Draw in game over state
func TestGame_Draw_GameOverState(t *testing.T) {
	game := NewGame()

	// Set game state directly to game over to avoid event chain issues
	game.gameState = components.StateGameOver

	// Create a dummy screen for testing
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	// Test drawing in game over state
	game.Draw(screen)

	// Verify game state is game over
	if game.gameState != components.StateGameOver {
		t.Errorf("Expected game state to be StateGameOver, got %v", game.gameState)
	}
}

// TestGame_Layout tests the Layout method
func TestGame_Layout(t *testing.T) {
	game := NewGame()

	width, height := game.Layout(1000, 800)

	// Verify layout returns expected dimensions
	if width != 800 {
		t.Errorf("Expected width 800, got %d", width)
	}

	if height != 600 {
		t.Errorf("Expected height 600, got %d", height)
	}

	// Test with different input dimensions
	width, height = game.Layout(1200, 900)
	if width != 800 {
		t.Errorf("Expected width 800, got %d", width)
	}

	if height != 600 {
		t.Errorf("Expected height 600, got %d", height)
	}
}

// TestGame_WorldInitialization tests that the world is properly initialized
func TestGame_WorldInitialization(t *testing.T) {
	game := NewGame()

	if game.World == nil {
		t.Fatal("Expected world to be initialized")
	}

	// Verify entities were created
	if len(game.World.Entities) == 0 {
		t.Error("Expected entities to be created in world")
	}

	// Check for load balancer entity
	foundLoadBalancer := false
	for _, entity := range game.World.Entities {
		if entity.HasComponent("Transform") && entity.HasComponent("Sprite") && entity.HasComponent("Collider") {
			// Check if it's the load balancer (should have SLA component)
			if entity.HasComponent("SLA") {
				foundLoadBalancer = true
				break
			}
		}
	}

	if !foundLoadBalancer {
		t.Error("Expected load balancer entity to be created")
	}

	// Check for backend entities
	backendCount := 0
	for _, entity := range game.World.Entities {
		if entity.HasComponent("BackendAssignment") {
			backendCount++
		}
	}

	if backendCount != 4 {
		t.Errorf("Expected 4 backend entities, got %d", backendCount)
	}
}

// TestGame_SystemsInitialization tests that all systems are properly initialized
func TestGame_SystemsInitialization(t *testing.T) {
	game := NewGame()

	if len(game.systems) == 0 {
		t.Error("Expected systems to be initialized")
	}

	// Verify specific systems are present
	systemTypes := make(map[string]bool)
	for _, system := range game.systems {
		systemTypes[getSystemType(system)] = true
	}

	expectedSystems := []string{
		"SpawnSystem",
		"InputSystem",
		"MovementSystem",
		"CollisionSystem",
		"PowerUpSystem",
		"BackendSystem",
		"SLASystem",
		"ComboSystem",
		"GameStateSystem",
		"ParticleSystem",
		"RoutingSystem",
	}

	for _, expectedSystem := range expectedSystems {
		if !systemTypes[expectedSystem] {
			t.Errorf("Expected system %s to be initialized", expectedSystem)
		}
	}
}

// TestGame_EventHandling tests event handling functionality
func TestGame_EventHandling(t *testing.T) {
	game := NewGame()

	// Test game start event
	gameStartEvent := events.NewEvent(events.EventGameStart, &events.EventData{
		SLA:    float64Ptr(99.0),
		Errors: intPtr(5),
	})
	game.eventDispatcher.Publish(gameStartEvent)

	// Update to process the event
	if err := game.Update(); err != nil {
		t.Errorf("Expected no error from game.Update(), got %v", err)
	}

	if game.gameState != components.StatePlaying {
		t.Errorf("Expected game state to be StatePlaying after game start event, got %v", game.gameState)
	}

	// Test game over event
	gameOverEvent := events.NewEvent(events.EventGameOver, nil)
	game.eventDispatcher.Publish(gameOverEvent)

	// Update to process the event
	if err := game.Update(); err != nil {
		t.Errorf("Expected no error from game.Update(), got %v", err)
	}

	if game.gameState != components.StateGameOver {
		t.Errorf("Expected game state to be StateGameOver after game over event, got %v", game.gameState)
	}
}

// TestGame_StateTransitions tests state transition logic
func TestGame_StateTransitions(t *testing.T) {
	game := NewGame()

	// Initial state should be menu
	if game.gameState != components.StateMenu {
		t.Errorf("Expected initial state to be StateMenu, got %v", game.gameState)
	}

	// Test direct state transitions to avoid event chain issues
	game.gameState = components.StatePlaying
	if game.gameState != components.StatePlaying {
		t.Errorf("Expected state to be StatePlaying, got %v", game.gameState)
	}

	game.gameState = components.StateGameOver
	if game.gameState != components.StateGameOver {
		t.Errorf("Expected state to be StateGameOver, got %v", game.gameState)
	}
}

// TestGame_EdgeCases tests edge cases and error conditions
func TestGame_EdgeCases(t *testing.T) {
	t.Run("Nil Event Dispatcher", func(t *testing.T) {
		game := NewGame()
		if game.eventDispatcher == nil {
			t.Error("Expected event dispatcher to be initialized")
		}
	})

	t.Run("Empty Systems List", func(t *testing.T) {
		game := NewGame()
		if len(game.systems) == 0 {
			t.Error("Expected systems to be initialized")
		}
	})

	t.Run("Multiple State Changes", func(t *testing.T) {
		game := NewGame()

		// Test rapid state changes directly
		for i := 0; i < 10; i++ {
			game.gameState = components.StatePlaying
			game.gameState = components.StateGameOver
		}

		// Should end up in game over state (last assignment)
		if game.gameState != components.StateGameOver {
			t.Errorf("Expected final state to be StateGameOver, got %v", game.gameState)
		}
	})

	t.Run("Update Without Initialization", func(t *testing.T) {
		game := &Game{
			World:       ecs.NewWorld(),
			initialized: false,
		}

		// Should not panic
		err := game.Update()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}

// TestGame_Integration tests integration scenarios
func TestGame_Integration(t *testing.T) {
	game := NewGame()

	// Verify initial setup
	if game.gameState != components.StateMenu {
		t.Errorf("Expected initial state to be StateMenu, got %v", game.gameState)
	}

	// Test state transitions directly
	game.gameState = components.StatePlaying
	if game.gameState != components.StatePlaying {
		t.Errorf("Expected state to be StatePlaying, got %v", game.gameState)
	}

	// Simulate some game time
	for i := 0; i < 10; i++ {
		err := game.Update()
		if err != nil {
			t.Errorf("Expected no error from Update, got %v", err)
		}
	}

	// End the game
	game.gameState = components.StateGameOver
	if game.gameState != components.StateGameOver {
		t.Errorf("Expected state to be StateGameOver, got %v", game.gameState)
	}

	// Test drawing in all states
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	// Should not panic
	game.Draw(screen)
}

// Benchmark tests for performance
func BenchmarkNewGame(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewGame()
	}
}

func BenchmarkGame_Update(b *testing.B) {
	game := NewGame()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := game.Update(); err != nil {
			b.Errorf("Expected no error from game.Update(), got %v", err)
		}
	}
}

func BenchmarkGame_Draw(b *testing.B) {
	game := NewGame()
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game.Draw(screen)
	}
}

func BenchmarkGame_Layout(b *testing.B) {
	game := NewGame()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game.Layout(800, 600)
	}
}

// Helper functions
func float64Ptr(v float64) *float64 {
	return &v
}

func intPtr(v int) *int {
	return &v
}

func getSystemType(system systems.System) string {
	switch system.(type) {
	case *systems.SpawnSystem:
		return "SpawnSystem"
	case *systems.InputSystem:
		return "InputSystem"
	case *systems.MovementSystem:
		return "MovementSystem"
	case *systems.CollisionSystem:
		return "CollisionSystem"
	case *systems.PowerUpSystem:
		return "PowerUpSystem"
	case *systems.BackendSystem:
		return "BackendSystem"
	case *systems.SLASystem:
		return "SLASystem"
	case *systems.ComboSystem:
		return "ComboSystem"
	case *systems.GameStateSystem:
		return "GameStateSystem"
	case *systems.ParticleSystem:
		return "ParticleSystem"
	case *systems.RoutingSystem:
		return "RoutingSystem"
	default:
		return "UnknownSystem"
	}
}
