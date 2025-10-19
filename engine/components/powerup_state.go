package components

// PowerUpState tracks active power-ups and timers at the entity/world level
type PowerUpState struct {
	// Simple single-active map of power-up name to remaining time
	RemainingByName map[string]float64
}

func NewPowerUpState() *PowerUpState {
	return &PowerUpState{RemainingByName: make(map[string]float64)}
}

func (p *PowerUpState) GetType() string { return "PowerUpState" }
