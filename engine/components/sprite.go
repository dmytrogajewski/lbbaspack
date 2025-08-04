package components

import "image/color"

// Sprite component represents visual representation
type Sprite struct {
	Width, Height float64
	Color         color.RGBA
	Visible       bool
	Layer         int // Rendering layer (higher = on top)
}

// NewSprite creates a new sprite component
func NewSprite(width, height float64, color color.RGBA) *Sprite {
	return &Sprite{
		Width:   width,
		Height:  height,
		Color:   color,
		Visible: true,
		Layer:   0,
	}
}

// GetType implements Component interface
func (s *Sprite) GetType() string {
	return "Sprite"
}

// GetWidth implements SpriteComponent interface
func (s *Sprite) GetWidth() float64 {
	return s.Width
}

// GetHeight implements SpriteComponent interface
func (s *Sprite) GetHeight() float64 {
	return s.Height
}

// GetColor implements SpriteComponent interface
func (s *Sprite) GetColor() color.RGBA {
	return s.Color
}

// IsVisible implements SpriteComponent interface
func (s *Sprite) IsVisible() bool {
	return s.Visible
}

// SetColor sets the sprite color
func (s *Sprite) SetColor(color color.RGBA) {
	s.Color = color
}

// SetVisible sets the sprite visibility
func (s *Sprite) SetVisible(visible bool) {
	s.Visible = visible
}
