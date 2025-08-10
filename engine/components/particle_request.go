package components

import "image/color"

// ParticleEffectRequest asks ParticleSystem to spawn an effect
type ParticleEffectRequest struct {
	X, Y  float64
	Color color.RGBA
	Kind  string // "packet" or "powerup"
}

func NewParticleEffectRequest(x, y float64, col color.RGBA, kind string) *ParticleEffectRequest {
	return &ParticleEffectRequest{X: x, Y: y, Color: col, Kind: kind}
}

func (r *ParticleEffectRequest) GetType() string { return "ParticleEffectRequest" }
