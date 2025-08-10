package systems

import (
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const SystemTypeParticle SystemType = "particle"

type ParticleSystem struct {
	BaseSystem
}

func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		BaseSystem: BaseSystem{},
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (ps *ParticleSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeParticle,
		System:       ps,
		Dependencies: []SystemType{SystemTypeCollision},
		Conflicts:    []SystemType{},
		Provides:     []string{"visual_effects", "particle_rendering"},
		Requires:     []string{},
		Drawable:     true,
		Optional:     true,
	}
}

func (ps *ParticleSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Update particles from ParticleState component
	var state *components.ParticleState
	for _, e := range entities {
		if comp := e.GetComponentByName("ParticleState"); comp != nil {
			if s, ok := comp.(*components.ParticleState); ok {
				state = s
				break
			}
		}
	}
	if state == nil {
		return
	}
	for i := len(state.Particles) - 1; i >= 0; i-- {
		p := state.Particles[i]
		p.Update(deltaTime)
		if !p.Active {
			state.Particles = append(state.Particles[:i], state.Particles[i+1:]...)
		}
	}

	// consume queued requests if any
	if len(state.Requests) > 0 {
		for _, req := range state.Requests {
			switch req.Kind {
			case "packet":
				ps.CreatePacketCatchEffect(req.X, req.Y, req.Color, state)
			case "powerup":
				ps.CreatePowerUpEffect(req.X, req.Y, req.Color, state)
			}
		}
		state.Requests = state.Requests[:0]
	}
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image, entities []Entity) {
	// Draw all active particles from component state
	var state *components.ParticleState
	for _, e := range entities {
		if comp := e.GetComponentByName("ParticleState"); comp != nil {
			if s, ok := comp.(*components.ParticleState); ok {
				state = s
				break
			}
		}
	}
	if state == nil {
		return
	}
	for _, particle := range state.Particles {
		if particle.Active {
			alpha := particle.GetAlpha()
			particleColor := color.RGBA{R: particle.Color.R, G: particle.Color.G, B: particle.Color.B, A: alpha}
			vector.DrawFilledCircle(screen, float32(particle.X), float32(particle.Y), float32(particle.Size), particleColor, false)
		}
	}
}

// CreatePacketCatchEffect creates particle effect when packet is caught
func (ps *ParticleSystem) CreatePacketCatchEffect(x, y float64, packetColor color.RGBA, state *components.ParticleState) {
	// Create multiple particles in a burst
	for range 8 {
		speed := 50.0 + rand.Float64()*50.0
		vx := speed * float64(rand.Float64()-0.5)
		vy := speed * float64(rand.Float64()-0.5)
		life := 0.5 + rand.Float64()*0.5
		size := 2.0 + rand.Float64()*3.0

		particle := components.NewParticle(x, y, vx, vy, life, packetColor, size)
		state.Particles = append(state.Particles, particle)
	}
}

// CreatePowerUpEffect creates particle effect when power-up is collected
func (ps *ParticleSystem) CreatePowerUpEffect(x, y float64, powerUpColor color.RGBA, state *components.ParticleState) {
	// Create sparkle effect
	for i := 0; i < 12; i++ {
		speed := 30.0 + rand.Float64()*40.0
		vx := speed * float64(rand.Float64()-0.5)
		vy := speed * float64(rand.Float64()-0.5)
		life := 1.0 + rand.Float64()*1.0
		size := 1.0 + rand.Float64()*2.0

		particle := components.NewParticle(x, y, vx, vy, life, powerUpColor, size)
		state.Particles = append(state.Particles, particle)
	}
}

func (ps *ParticleSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	eventDispatcher.Subscribe(events.EventCollisionDetected, func(event *events.Event) {
		if event == nil || event.Data == nil {
			return
		}
		if event.Data.TagA == nil || event.Data.TagB == nil {
			return
		}
		// Spawn packet or powerup effect at collision point by adding a request to ParticleState holder
		kind := ""
		if *event.Data.TagA == "packet" || *event.Data.TagB == "packet" {
			kind = "packet"
		} else if *event.Data.TagA == "powerup" || *event.Data.TagB == "powerup" {
			kind = "powerup"
		}
		if kind == "" {
			return
		}
	})
}
