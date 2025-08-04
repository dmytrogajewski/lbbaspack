package components

// Transform component represents position, rotation, and scale
type Transform struct {
	X, Y           float64
	Rotation       float64
	ScaleX, ScaleY float64
}

// NewTransform creates a new transform component
func NewTransform(x, y float64) *Transform {
	return &Transform{
		X: x, Y: y,
		Rotation: 0,
		ScaleX:   1, ScaleY: 1,
	}
}

// GetType implements Component interface
func (t *Transform) GetType() string {
	return "Transform"
}

// GetX implements TransformComponent interface
func (t *Transform) GetX() float64 {
	return t.X
}

// GetY implements TransformComponent interface
func (t *Transform) GetY() float64 {
	return t.Y
}

// SetPosition sets the position
func (t *Transform) SetPosition(x, y float64) {
	t.X = x
	t.Y = y
}

// GetPosition returns the current position
func (t *Transform) GetPosition() (float64, float64) {
	return t.X, t.Y
}
