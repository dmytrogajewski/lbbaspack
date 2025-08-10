package components

// Collider component represents collision detection
type Collider struct {
	Width, Height float64
	Active        bool
	Tag           string // For identifying collision types
}

// NewCollider creates a new collider component
func NewCollider(width, height float64, tag string) *Collider {
	return &Collider{
		Width:  width,
		Height: height,
		Active: true,
		Tag:    tag,
	}
}

// GetType implements Component interface
func (c *Collider) GetType() string {
	return "Collider"
}

// GetTag implements ColliderComponent interface
func (c *Collider) GetTag() string {
	return c.Tag
}

// GetWidth implements ColliderComponent interface
func (c *Collider) GetWidth() float64 {
	return c.Width
}

// GetHeight implements ColliderComponent interface
func (c *Collider) GetHeight() float64 {
	return c.Height
}

// SetTag implements ColliderComponent interface
func (c *Collider) SetTag(tag string) {
	c.Tag = tag
}

// SetWidth updates collider width
func (c *Collider) SetWidth(width float64) { c.Width = width }

// SetHeight updates collider height
func (c *Collider) SetHeight(height float64) { c.Height = height }

// GetBounds returns the collision bounds
func (c *Collider) GetBounds(x, y float64) (float64, float64, float64, float64) {
	return x, y, x + c.Width, y + c.Height
}

// CheckCollision checks collision with another collider
func (c *Collider) CheckCollision(x1, y1 float64, other *Collider, x2, y2 float64) bool {
	if !c.Active || !other.Active {
		return false
	}

	left1, top1, right1, bottom1 := c.GetBounds(x1, y1)
	left2, top2, right2, bottom2 := other.GetBounds(x2, y2)

	return left1 < right2 && right1 > left2 && top1 < bottom2 && bottom1 > top2
}
