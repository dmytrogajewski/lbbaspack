package components

type StateType int

const (
	StateMenu StateType = iota
	StatePlaying
	StateGameOver
)

type State struct {
	Current StateType
}

func NewState(initial StateType) *State {
	return &State{Current: initial}
}

// GetType implements Component interface
func (s *State) GetType() string {
	return "State"
}

// GetState implements StateComponent interface
func (s *State) GetState() string {
	switch s.Current {
	case StateMenu:
		return "menu"
	case StatePlaying:
		return "playing"
	case StateGameOver:
		return "gameover"
	default:
		return "unknown"
	}
}

// SetState implements StateComponent interface
func (s *State) SetState(state string) {
	switch state {
	case "menu":
		s.Current = StateMenu
	case "playing":
		s.Current = StatePlaying
	case "gameover":
		s.Current = StateGameOver
	}
}
