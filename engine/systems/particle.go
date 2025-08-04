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
	particles []*components.Particle
}

func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		BaseSystem: BaseSystem{},
		particles:  make([]*components.Particle, 0),
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
	// Update existing particles
	for i := len(ps.particles) - 1; i >= 0; i-- {
		particle := ps.particles[i]
		particle.Update(deltaTime)

		// Remove dead particles
		if !particle.Active {
			ps.particles = append(ps.particles[:i], ps.particles[i+1:]...)
		}
	}
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	// Draw all active particles
	for _, particle := range ps.particles {
		if particle.Active {
			alpha := particle.GetAlpha()
			particleColor := color.RGBA{
				R: particle.Color.R,
				G: particle.Color.G,
				B: particle.Color.B,
				A: alpha,
			}
			vector.DrawFilledCircle(screen, float32(particle.X), float32(particle.Y), float32(particle.Size), particleColor, false)
		}
	}
}

// CreatePacketCatchEffect creates particle effect when packet is caught
func (ps *ParticleSystem) CreatePacketCatchEffect(x, y float64, packetColor color.RGBA) {
	// Create multiple particles in a burst
	for range 8 {
		speed := 50.0 + rand.Float64()*50.0
		vx := speed * float64(rand.Float64()-0.5)
		vy := speed * float64(rand.Float64()-0.5)
		life := 0.5 + rand.Float64()*0.5
		size := 2.0 + rand.Float64()*3.0

		particle := components.NewParticle(x, y, vx, vy, life, packetColor, size)
		ps.particles = append(ps.particles, particle)
	}
}

// CreatePowerUpEffect creates particle effect when power-up is collected
func (ps *ParticleSystem) CreatePowerUpEffect(x, y float64, powerUpColor color.RGBA) {
	// Create sparkle effect
	for i := 0; i < 12; i++ {
		speed := 30.0 + rand.Float64()*40.0
		vx := speed * float64(rand.Float64()-0.5)
		vy := speed * float64(rand.Float64()-0.5)
		life := 1.0 + rand.Float64()*1.0
		size := 1.0 + rand.Float64()*2.0

		particle := components.NewParticle(x, y, vx, vy, life, powerUpColor, size)
		ps.particles = append(ps.particles, particle)
	}
}

func (ps *ParticleSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Listen for collision events to create particle effects
	eventDispatcher.Subscribe(events.EventPacketCaught, func(event *events.Event) {
		if event.Data.Packet != nil {
			if packetEntity, ok := event.Data.Packet.(Entity); ok {
				transformComp := packetEntity.GetTransform()
				spriteComp := packetEntity.GetSprite()
				if transformComp != nil && spriteComp != nil {
					ps.CreatePacketCatchEffect(transformComp.GetX()+7.5, transformComp.GetY()+7.5, spriteComp.GetColor())
				}
			}
		}
	})

	eventDispatcher.Subscribe(events.EventPowerUpCollected, func(event *events.Event) {
		if event.Data.Packet != nil {
			if powerupEntity, ok := event.Data.Packet.(Entity); ok {
				transformComp := powerupEntity.GetTransform()
				spriteComp := powerupEntity.GetSprite()
				if transformComp != nil && spriteComp != nil {
					ps.CreatePowerUpEffect(transformComp.GetX()+7.5, transformComp.GetY()+7.5, spriteComp.GetColor())
				}
			}
		}
	})
}
