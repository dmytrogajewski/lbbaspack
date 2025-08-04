package components

// Physics component represents velocity and acceleration
type Physics struct {
	VelocityX, VelocityY         float64
	AccelerationX, AccelerationY float64
	Mass                         float64
	Friction                     float64
}

// NewPhysics creates a new physics component
func NewPhysics() *Physics {
	return &Physics{
		VelocityX: 0, VelocityY: 0,
		AccelerationX: 0, AccelerationY: 0,
		Mass:     1.0,
		Friction: 1.0, // No friction for packets
	}
}

// GetType implements Component interface
func (p *Physics) GetType() string {
	return "Physics"
}

// GetVelocityX implements PhysicsComponent interface
func (p *Physics) GetVelocityX() float64 {
	return p.VelocityX
}

// GetVelocityY implements PhysicsComponent interface
func (p *Physics) GetVelocityY() float64 {
	return p.VelocityY
}

// SetVelocity sets the velocity
func (p *Physics) SetVelocity(vx, vy float64) {
	p.VelocityX = vx
	p.VelocityY = vy
}

// GetVelocity returns the current velocity
func (p *Physics) GetVelocity() (float64, float64) {
	return p.VelocityX, p.VelocityY
}

// ApplyForce applies a force to the physics component
func (p *Physics) ApplyForce(fx, fy float64) {
	p.AccelerationX += fx / p.Mass
	p.AccelerationY += fy / p.Mass
}

// Update updates the physics component
func (p *Physics) Update(deltaTime float64) {
	// Apply acceleration
	p.VelocityX += p.AccelerationX * deltaTime
	p.VelocityY += p.AccelerationY * deltaTime

	// Apply friction (only if not 1.0 to avoid unnecessary multiplication)
	if p.Friction != 1.0 {
		p.VelocityX *= p.Friction
		p.VelocityY *= p.Friction
	}

	// Reset acceleration
	p.AccelerationX = 0
	p.AccelerationY = 0
}
