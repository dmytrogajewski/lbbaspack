package main

import (
	"fmt"
	"log"

	"lbbaspack/engine/components"
	"lbbaspack/engine/ecs"
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
	systems         []systems.System
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

	// Create systems
	spawnSys := systems.NewSpawnSystem(func() systems.Entity {
		return world.NewEntity()
	})
	inputSys := systems.NewInputSystem()
	movementSys := systems.NewMovementSystem()
	collisionSys := systems.NewCollisionSystem()
	powerUpSys := systems.NewPowerUpSystem()
	backendSys := systems.NewBackendSystem()
	slaSys := systems.NewSLASystem(spawnSys)
	comboSys := systems.NewComboSystem()
	gameStateSys := systems.NewGameStateSystem()
	renderSys := systems.NewRenderSystem()
	particleSys := systems.NewParticleSystem()
	routingSys := systems.NewRoutingSystem()
	uiSys := systems.NewUISystem(nil)     // Will be set in Draw
	menuSys := systems.NewMenuSystem(nil) // Will be set in Draw

	// Initialize all systems
	spawnSys.Initialize(eventDispatcher)
	backendSys.Initialize(eventDispatcher)
	slaSys.Initialize(eventDispatcher)
	comboSys.Initialize(eventDispatcher)
	gameStateSys.Initialize(eventDispatcher)
	particleSys.Initialize(eventDispatcher)
	uiSys.Initialize(eventDispatcher)
	routingSys.Initialize(eventDispatcher)

	game := &Game{
		World:           world,
		initialized:     true,
		RenderSys:       renderSys,
		UISys:           uiSys,
		MenuSys:         menuSys,
		gameState:       components.StateMenu,
		eventDispatcher: eventDispatcher,
		systems: []systems.System{
			spawnSys,
			inputSys,
			movementSys,
			collisionSys,
			powerUpSys,
			backendSys,
			slaSys,
			comboSys,
			gameStateSys,
			particleSys,
			routingSys,
		},
	}

	// Set up event handlers after game is created
	eventDispatcher.Subscribe(events.EventGameStart, func(event *events.Event) {
		fmt.Println("Game start event received, transitioning to playing state")
		// Set SLA parameters based on selected mode
		if event.Data.SLA != nil {
			slaSys.SetTargetSLA(*event.Data.SLA)
		}
		if event.Data.Errors != nil {
			slaSys.SetErrorBudget(*event.Data.Errors)
		}
		game.gameState = components.StatePlaying
	})

	eventDispatcher.Subscribe(events.EventGameOver, func(event *events.Event) {
		game.gameState = components.StateGameOver
		fmt.Println("Game over event received, transitioning to game over state")
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
		// Update all systems
		entitiesInterface := make([]systems.Entity, len(g.World.Entities))
		for i, entity := range g.World.Entities {
			entitiesInterface[i] = entity
		}

		for i, system := range g.systems {
			fmt.Printf("[Game] Updating system %d: %T\n", i, system)
			system.Update(deltaTime, entitiesInterface, g.eventDispatcher)
		}
	case components.StateGameOver:
		fmt.Println("[Game] In game over state")
		// Game over state - only handle input for restart
		// For now, just keep the game in this state
		// TODO: Add restart functionality
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

		// Draw particles
		if particleSys, ok := g.systems[len(g.systems)-2].(*systems.ParticleSystem); ok {
			particleSys.Draw(screen)
		}

		// Draw routing
		if routingSys, ok := g.systems[len(g.systems)-1].(*systems.RoutingSystem); ok {
			routingSys.Draw(screen)
		}

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
