package components

// PowerUpType component identifies the type of power-up and its duration
type PowerUpType struct {
	Name     string
	Duration float64 // seconds
}

func NewPowerUpType(name string, duration float64) *PowerUpType {
	return &PowerUpType{Name: name, Duration: duration}
}

// GetType implements Component interface
func (p *PowerUpType) GetType() string {
	return "PowerUpType"
}

// GetName implements PowerUpTypeComponent interface
func (p *PowerUpType) GetName() string {
	return p.Name
}

// GetDuration implements PowerUpTypeComponent interface
func (p *PowerUpType) GetDuration() float64 {
	return p.Duration
}

// GetEffect implements PowerUpTypeComponent interface
func (p *PowerUpType) GetEffect() string {
	return p.Name // For now, effect is the same as name
}
