package components

import "image/color"

// Particle component for visual effects
type Particle struct {
	X, Y      float64
	VelocityX float64
	VelocityY float64
	Life      float64
	MaxLife   float64
	Color     color.RGBA
	Size      float64
	Active    bool
}

// NewParticle creates a new particle
func NewParticle(x, y, vx, vy, life float64, color color.RGBA, size float64) *Particle {
	return &Particle{
		X:         x,
		Y:         y,
		VelocityX: vx,
		VelocityY: vy,
		Life:      life,
		MaxLife:   life,
		Color:     color,
		Size:      size,
		Active:    true,
	}
}

// Update updates the particle
func (p *Particle) Update(deltaTime float64) {
	p.X += p.VelocityX * deltaTime
	p.Y += p.VelocityY * deltaTime
	p.Life -= deltaTime

	if p.Life <= 0 {
		p.Active = false
	}
}

// GetAlpha returns the alpha value based on remaining life
func (p *Particle) GetAlpha() uint8 {
	if p.MaxLife <= 0 {
		return 255
	}
	alpha := uint8((p.Life / p.MaxLife) * 255)
	return alpha
}
