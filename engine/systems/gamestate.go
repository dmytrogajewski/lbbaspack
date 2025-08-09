package systems

import (
	"fmt"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

const SystemTypeGameState SystemType = "gamestate"

type GameStateSystem struct {
	BaseSystem
}

func NewGameStateSystem() *GameStateSystem {
	return &GameStateSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"State",
				"GameSession",
			},
		},
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
	// Read session component
	var session *components.GameSession
	for _, e := range entities {
		if comp := e.GetComponentByName("GameSession"); comp != nil {
			if s, ok := comp.(*components.GameSession); ok {
				session = s
				break
			}
		}
	}
	if session == nil {
		return
	}

	session.GameTime += deltaTime

	// Determine current state from State component(s)
	currentState := components.StateMenu
	for _, e := range gss.FilterEntities(entities) {
		if st := e.GetState(); st != nil {
			switch st.GetState() {
			case "playing":
				currentState = components.StatePlaying
			case "gameover":
				currentState = components.StateGameOver
			}
			break
		}
	}

	// Normalize state string back to components (optional)
	for _, entity := range gss.FilterEntities(entities) {
		if st := entity.GetState(); st != nil {
			switch currentState {
			case components.StateMenu:
				st.SetState("menu")
			case components.StatePlaying:
				st.SetState("playing")
			case components.StateGameOver:
				st.SetState("gameover")
			}
		}
	}

	switch currentState {
	case components.StatePlaying:
		gss.updatePlayingState(deltaTime, eventDispatcher, session)
	case components.StateGameOver:
		gss.updateGameOverState(deltaTime, eventDispatcher, session)
	}
}

func (gss *GameStateSystem) getStateString() string {
	// This helper is now used only to normalize state strings when writing to State components
	// We don't track state internally; default to menu
	switch components.StateMenu {
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
		// no internal state to change
	})

	eventDispatcher.Subscribe(events.EventGameOver, func(event *events.Event) {})

	eventDispatcher.Subscribe(events.EventPacketCaught, func(event *events.Event) {})
}

// No internal transitions; state is kept in components and events

// No internal transitions; handled by components/UI

func (gss *GameStateSystem) updatePlayingState(deltaTime float64, eventDispatcher *events.EventDispatcher, session *components.GameSession) {
	if int(session.GameTime)%30 == 0 && int(session.GameTime) > 0 && session.GameTime-session.LastLevelUpTime > 1.0 {
		gss.levelUp(eventDispatcher, session)
	}
}

func (gss *GameStateSystem) updateGameOverState(deltaTime float64, eventDispatcher *events.EventDispatcher, session *components.GameSession) {
	// Handle restart logic or return to menu
	// This could be expanded with restart functionality
}

func (gss *GameStateSystem) checkLevelUp(eventDispatcher *events.EventDispatcher, session *components.GameSession) {
	if session.Score%100 == 0 && session.Score > 0 && session.GameTime-session.LastLevelUpTime > 1.0 {
		gss.levelUp(eventDispatcher, session)
	}
}

// levelUp handles the actual level-up logic
func (gss *GameStateSystem) levelUp(eventDispatcher *events.EventDispatcher, session *components.GameSession) {
	session.Level++
	session.LastLevelUpTime = session.GameTime
	fmt.Printf("Level up! Now level %d (Score: %d, Time: %.1fs)\n", session.Level, session.Score, session.GameTime)
	level := session.Level
	time := session.GameTime
	eventDispatcher.Publish(events.NewEvent(events.EventLevelUp, &events.EventData{Level: &level, Time: &time}))
}

// Accessors removed; data is in components
