package systems

import (
	"fmt"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

const SystemTypeGameState SystemType = "gamestate"

type GameStateSystem struct {
	BaseSystem
	currentState    components.StateType
	gameTime        float64
	score           int
	level           int
	lastLevelUpTime float64 // Track when we last leveled up to prevent multiple level-ups
}

func NewGameStateSystem() *GameStateSystem {
	return &GameStateSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"State",
			},
		},
		currentState:    components.StateMenu,
		gameTime:        0.0,
		score:           0,
		level:           1,
		lastLevelUpTime: 0.0,
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (gss *GameStateSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeGameState,
		System:       gss,
		Dependencies: []SystemType{}, // No dependencies - runs independently and receives events
		Conflicts:    []SystemType{},
		Provides:     []string{"game_state_management", "level_progression"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     false,
	}
}

func (gss *GameStateSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	gss.gameTime += deltaTime

	// Update all entities with state components
	for _, entity := range gss.FilterEntities(entities) {
		stateComp := entity.GetState()
		if stateComp == nil {
			continue
		}
		state := stateComp
		state.SetState(gss.getStateString())
	}

	// Handle state-specific logic
	switch gss.currentState {
	case components.StatePlaying:
		gss.updatePlayingState(deltaTime, eventDispatcher)
	case components.StateGameOver:
		gss.updateGameOverState(deltaTime, eventDispatcher)
	}
}

func (gss *GameStateSystem) getStateString() string {
	switch gss.currentState {
	case components.StateMenu:
		return "menu"
	case components.StatePlaying:
		return "playing"
	case components.StateGameOver:
		return "gameover"
	default:
		return "unknown"
	}
}

func (gss *GameStateSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Handle state transitions based on events
	eventDispatcher.Subscribe(events.EventGameStart, func(event *events.Event) {
		gss.transitionToPlaying()
	})

	eventDispatcher.Subscribe(events.EventGameOver, func(event *events.Event) {
		gss.transitionToGameOver()
	})

	eventDispatcher.Subscribe(events.EventPacketCaught, func(event *events.Event) {
		if gss.currentState == components.StatePlaying {
			gss.score += 10
			gss.checkLevelUp(eventDispatcher)
		}
	})
}

func (gss *GameStateSystem) transitionToPlaying() {
	if gss.currentState == components.StateMenu || gss.currentState == components.StateGameOver {
		gss.currentState = components.StatePlaying
		gss.gameTime = 0.0
		gss.score = 0
		gss.level = 1
		fmt.Println("Game started! Good luck!")
	}
}

func (gss *GameStateSystem) transitionToGameOver() {
	if gss.currentState == components.StatePlaying {
		gss.currentState = components.StateGameOver
		fmt.Printf("Game Over! Final Score: %d, Level: %d, Time: %.1fs\n",
			gss.score, gss.level, gss.gameTime)
	}
}

func (gss *GameStateSystem) updatePlayingState(deltaTime float64, eventDispatcher *events.EventDispatcher) {
	// Check for level progression every 30 seconds (time-based)
	if int(gss.gameTime)%30 == 0 && int(gss.gameTime) > 0 && gss.gameTime-gss.lastLevelUpTime > 1.0 {
		gss.levelUp(eventDispatcher)
	}
}

func (gss *GameStateSystem) updateGameOverState(deltaTime float64, eventDispatcher *events.EventDispatcher) {
	// Handle restart logic or return to menu
	// This could be expanded with restart functionality
}

func (gss *GameStateSystem) checkLevelUp(eventDispatcher *events.EventDispatcher) {
	// Level up every 100 points (score-based)
	if gss.score%100 == 0 && gss.score > 0 && gss.gameTime-gss.lastLevelUpTime > 1.0 {
		gss.levelUp(eventDispatcher)
	}
}

// levelUp handles the actual level-up logic
func (gss *GameStateSystem) levelUp(eventDispatcher *events.EventDispatcher) {
	gss.level++
	gss.lastLevelUpTime = gss.gameTime
	fmt.Printf("Level up! Now level %d (Score: %d, Time: %.1fs)\n", gss.level, gss.score, gss.gameTime)

	levelEvent := events.NewEvent(events.EventLevelUp, &events.EventData{
		Level: &gss.level,
		Time:  &gss.gameTime,
	})
	eventDispatcher.Publish(levelEvent)
}

func (gss *GameStateSystem) GetCurrentState() components.StateType {
	return gss.currentState
}

func (gss *GameStateSystem) GetScore() int {
	return gss.score
}

func (gss *GameStateSystem) GetLevel() int {
	return gss.level
}

func (gss *GameStateSystem) GetGameTime() float64 {
	return gss.gameTime
}
