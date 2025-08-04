package main

import (
	"fmt"
	"log"

	"lbbaspack/engine/components"
	"lbbaspack/engine/ecs"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"lbbaspack/engine/systems"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Game struct {
	World           *ecs.World
	initialized     bool
	RenderSys       *systems.RenderSystem
	UISys           *systems.UISystem
	MenuSys         *systems.MenuSystem
	gameState       components.StateType
	eventDispatcher *events.EventDispatcher
	systemManager   *systems.SystemManager
}

func NewGame() *Game {
	fmt.Println("Initializing game...")
	world := ecs.NewWorld()

	// --- Entity/Component Initialization ---
	// Spawn Load Balancer
	loadBalancer := world.NewEntity()
	loadBalancer.AddComponent(components.NewTransform(350, 480))
	loadBalancer.AddComponent(components.NewSprite(100, 20, color.RGBA{100, 100, 255, 255}))
	loadBalancer.AddComponent(components.NewCollider(100, 20, "loadbalancer"))
	loadBalancer.AddComponent(&components.State{Current: components.StateMenu}) // Start in menu state
	loadBalancer.AddComponent(&components.Combo{})                              // Add combo component
	loadBalancer.AddComponent(components.NewSLA(99.5, 10))                      // Add SLA component to load balancer

	// Spawn Backends
	backendCount := 4
	backendWidth := 120
	backendSpacing := (800 - backendWidth*backendCount) / (backendCount + 1)
	backendY := 600 - 50
	for i := 0; i < backendCount; i++ {
		backend := world.NewEntity()
		x := float64(backendSpacing + i*(backendWidth+backendSpacing))
		backend.AddComponent(components.NewTransform(x, float64(backendY)))
		backend.AddComponent(components.NewSprite(float64(backendWidth), 40, color.RGBA{0, 255, 0, 255}))
		backend.AddComponent(components.NewCollider(float64(backendWidth), 40, "backend")) // Add collider for labels
		backend.AddComponent(components.NewBackendAssignment(i))                           // Add backend assignment
		backend.AddComponent(components.NewSLA(99.5, 10))                                  // Add SLA component
	}

	// --- System Initialization ---
	eventDispatcher := events.NewEventDispatcher()

	// Create system factory and manager
	systemFactory := systems.NewSystemFactory(func() systems.Entity {
		return world.NewEntity()
	}, eventDispatcher)

	systemManager, err := systemFactory.CreateSystemManager()
	if err != nil {
		log.Fatalf("Failed to create system manager: %v", err)
	}

	// Get individual systems for special handling
	uiSys := systems.NewUISystem(nil)     // Will be set in Draw
	menuSys := systems.NewMenuSystem(nil) // Will be set in Draw
	renderSys := systems.NewRenderSystem()

	// Initialize UI system
	uiSys.Initialize(eventDispatcher)

	game := &Game{
		World:           world,
		initialized:     true,
		RenderSys:       renderSys,
		UISys:           uiSys,
		MenuSys:         menuSys,
		gameState:       components.StateMenu,
		eventDispatcher: eventDispatcher,
		systemManager:   systemManager,
	}

	// Set up event handlers after game is created
	eventDispatcher.Subscribe(events.EventGameStart, func(event *events.Event) {
		fmt.Println("Game start event received, transitioning to playing state")

		// Clean up any existing game entities (packets, power-ups, etc.)
		game.cleanupGameEntities()

		// Set SLA parameters based on selected mode
		if slaSys, err := systemFactory.GetSystemByType(game.systemManager, systems.SystemTypeSLA); err == nil {
			if slaSystem, ok := slaSys.(*systems.SLASystem); ok {
				if event.Data.SLA != nil {
					slaSystem.SetTargetSLA(*event.Data.SLA)
				}
				if event.Data.Errors != nil {
					slaSystem.SetErrorBudget(*event.Data.Errors)
				}
			}
		}
		game.gameState = components.StatePlaying

		// Ensure load balancer state is updated to playing
		for _, entity := range game.World.Entities {
			if entity.HasComponent("State") {
				stateComp := entity.GetState()
				if stateComp != nil {
					stateComp.SetState("playing")
					fmt.Printf("[Game] Updated entity %d state to playing\n", entity.ID)
				}
			}
		}
	})

	eventDispatcher.Subscribe(events.EventGameOver, func(event *events.Event) {
		game.gameState = components.StateGameOver
		fmt.Println("Game over event received, transitioning to game over state")
	})

	// Add handler for returning to menu
	eventDispatcher.Subscribe(events.EventReturnToMenu, func(event *events.Event) {
		game.gameState = components.StateMenu
		fmt.Println("Return to menu event received, transitioning to menu state")

		// Reset load balancer state to menu
		for _, entity := range game.World.Entities {
			if entity.HasComponent("State") {
				stateComp := entity.GetState()
				if stateComp != nil {
					stateComp.SetState("menu")
					fmt.Printf("[Game] Reset entity %d state to menu\n", entity.ID)
				}
			}
		}
	})

	return game
}

func (g *Game) Update() error {
	deltaTime := 1.0 / 60.0 // 60 FPS

	// Debug: Print current game state
	fmt.Printf("[Game] Current game state: %v\n", g.gameState)

	// Handle game state
	switch g.gameState {
	case components.StateMenu:
		fmt.Println("[Game] In menu state - updating menu system only")
		// Menu system handles its own update
		if g.MenuSys != nil {
			// Convert entities to interface type
			entitiesInterface := make([]systems.Entity, len(g.World.Entities))
			for i, entity := range g.World.Entities {
				entitiesInterface[i] = entity
			}
			g.MenuSys.Update(deltaTime, entitiesInterface, g.eventDispatcher)
		}
	case components.StatePlaying:
		fmt.Println("[Game] In playing state - updating all systems")
		// Update all systems using the system manager
		entitiesInterface := make([]systems.Entity, len(g.World.Entities))
		for i, entity := range g.World.Entities {
			entitiesInterface[i] = entity
		}

		g.systemManager.UpdateAll(deltaTime, entitiesInterface, g.eventDispatcher)

		// Update UI system separately since it's not in the system manager
		g.UISys.Update(deltaTime, entitiesInterface, g.eventDispatcher)

		// Clean up inactive entities after all systems have updated
		g.World.RemoveInactiveEntities()
	case components.StateGameOver:
		fmt.Println("[Game] In game over state")
		// Handle escape key to return to menu
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			fmt.Println("[Game] Escape pressed in game over state, returning to menu")
			// Clean up game entities when returning to menu
			g.cleanupGameEntities()
			// Publish return to menu event
			g.eventDispatcher.Publish(events.NewEvent(events.EventReturnToMenu, nil))
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.initialized {
		fmt.Println("Initializing render systems...")
		g.UISys = systems.NewUISystem(screen)
		g.MenuSys = systems.NewMenuSystem(screen)
		g.initialized = true
		fmt.Println("Render systems added successfully")
	}

	// Handle different game states
	switch g.gameState {
	case components.StateMenu:
		g.MenuSys.Draw(screen)
	case components.StatePlaying:
		// Convert []*entities.Entity to []systems.Entity
		entitiesInterface := make([]systems.Entity, len(g.World.Entities))
		for i, entity := range g.World.Entities {
			entitiesInterface[i] = entity
		}
		g.RenderSys.UpdateWithScreen(1.0/60.0, entitiesInterface, g.eventDispatcher, screen)

		// Draw particles and routing using system manager
		g.systemManager.DrawAll(screen, entitiesInterface)

		g.UISys.Draw(screen, entitiesInterface)
	case components.StateGameOver:
		// Draw game over screen
		screen.Fill(color.RGBA{20, 20, 40, 255})
		text.Draw(screen, "GAME OVER", basicfont.Face7x13, 350, 280, color.White)
		text.Draw(screen, "Press ESC to return to menu", basicfont.Face7x13, 300, 320, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

// cleanupGameEntities removes all game-related entities (packets, power-ups) but keeps the load balancer and backends
func (g *Game) cleanupGameEntities() {
	fmt.Println("[Game] Cleaning up game entities...")

	// Keep track of entities to remove
	entitiesToRemove := make([]*entities.Entity, 0)

	for _, entity := range g.World.Entities {
		// Check if this is a game entity (packet or power-up) that should be removed
		if entity.HasComponent("PacketType") || entity.HasComponent("PowerUpType") {
			entitiesToRemove = append(entitiesToRemove, entity)
			fmt.Printf("[Game] Marking entity %d for removal (has PacketType or PowerUpType)\n", entity.ID)
		}
	}

	// Remove the marked entities
	for _, entity := range entitiesToRemove {
		g.World.RemoveEntity(entity)
		fmt.Printf("[Game] Removed entity %d\n", entity.ID)
	}

	fmt.Printf("[Game] Cleanup complete. Removed %d entities. World now has %d entities.\n",
		len(entitiesToRemove), len(g.World.Entities))
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("LBaaS Packet Catcher - ECS Edition")
	ebiten.SetFullscreen(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
