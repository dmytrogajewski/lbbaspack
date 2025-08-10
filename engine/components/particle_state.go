package components

import "image/color"

// ParticleState stores particles for visual effects
type ParticleState struct {
	Particles []*Particle
	Requests  []*ParticleEffectRequest
}

func NewParticleState() *ParticleState {
	return &ParticleState{Particles: make([]*Particle, 0), Requests: make([]*ParticleEffectRequest, 0)}
}

func (p *ParticleState) GetType() string { return "ParticleState" }

// Route visuals component state
type Route struct {
	StartX, StartY float64
	EndX, EndY     float64
	Progress       float64
	Speed          float64
	Color          color.RGBA
	Active         bool
}

type RouteState struct {
	Routes []*Route
}

func NewRouteState() *RouteState { return &RouteState{Routes: make([]*Route, 0)} }

func (r *RouteState) GetType() string { return "RouteState" }
