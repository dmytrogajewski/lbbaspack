package components

// PowerUpActivation marks a collected power-up to be applied by PowerUpSystem
type PowerUpActivation struct {
	Name     string
	Duration float64
}

func NewPowerUpActivation(name string, duration float64) *PowerUpActivation {
	return &PowerUpActivation{Name: name, Duration: duration}
}

func (p *PowerUpActivation) GetType() string { return "PowerUpActivation" }
