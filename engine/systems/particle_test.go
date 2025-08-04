package systems

import (
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

const particleTolerance = 0.0001

func TestNewParticleSystem(t *testing.T) {
	ps := NewParticleSystem()

	// Test that the system is properly initialized
	if ps == nil {
		t.Fatal("NewParticleSystem returned nil")
	}

	// Test particles slice initialization
	if ps.particles == nil {
		t.Fatal("Particles slice should not be nil")
	}

	if len(ps.particles) != 0 {
		t.Errorf("Expected initial particles count to be 0, got %d", len(ps.particles))
	}
}

func TestParticleSystem_Update_NoParticles(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test with no particles
	entities := []Entity{}

	// Run update
	ps.Update(0.016, entities, eventDispatcher)

	// Verify particles slice remains empty
	if len(ps.particles) != 0 {
		t.Errorf("Expected particles count to remain 0, got %d", len(ps.particles))
	}
}

func TestParticleSystem_Update_WithParticles(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create some particles
	particle1 := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	particle2 := components.NewParticle(200, 200, -5, 10, 0.5, color.RGBA{0, 255, 0, 255}, 3.0)
	ps.particles = append(ps.particles, particle1, particle2)

	entities := []Entity{}

	// Run update
	ps.Update(0.016, entities, eventDispatcher)

	// Verify particles were updated
	// Particle1 should still be active (life > 0)
	if !particle1.Active {
		t.Error("Expected particle1 to still be active")
	}

	// Particle2 should still be active (life > 0)
	if !particle2.Active {
		t.Error("Expected particle2 to still be active")
	}

	// Verify positions were updated
	expectedX1 := 100.0 + 10.0*0.016
	expectedY1 := 100.0 + 5.0*0.016
	if math.Abs(particle1.X-expectedX1) > particleTolerance {
		t.Errorf("Expected particle1 X to be %f, got %f", expectedX1, particle1.X)
	}
	if math.Abs(particle1.Y-expectedY1) > particleTolerance {
		t.Errorf("Expected particle1 Y to be %f, got %f", expectedY1, particle1.Y)
	}

	expectedX2 := 200.0 + (-5.0)*0.016
	expectedY2 := 200.0 + 10.0*0.016
	if math.Abs(particle2.X-expectedX2) > particleTolerance {
		t.Errorf("Expected particle2 X to be %f, got %f", expectedX2, particle2.X)
	}
	if math.Abs(particle2.Y-expectedY2) > particleTolerance {
		t.Errorf("Expected particle2 Y to be %f, got %f", expectedY2, particle2.Y)
	}
}

func TestParticleSystem_Update_ParticleExpiration(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create a particle with very short life
	particle := components.NewParticle(100, 100, 10, 5, 0.01, color.RGBA{255, 0, 0, 255}, 2.0)
	ps.particles = append(ps.particles, particle)

	// Run update with delta time that exceeds particle life
	ps.Update(0.02, []Entity{}, eventDispatcher)

	// Verify particle was removed
	if len(ps.particles) != 0 {
		t.Errorf("Expected particles count to be 0 after expiration, got %d", len(ps.particles))
	}

	// Verify particle is marked as inactive
	if particle.Active {
		t.Error("Expected particle to be marked as inactive")
	}
}

func TestParticleSystem_Update_MixedParticles(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create particles with different life spans
	particle1 := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)   // Long life
	particle2 := components.NewParticle(200, 200, -5, 10, 0.01, color.RGBA{0, 255, 0, 255}, 3.0) // Short life
	particle3 := components.NewParticle(300, 300, 0, 0, 0.5, color.RGBA{0, 0, 255, 255}, 1.0)    // Medium life
	ps.particles = append(ps.particles, particle1, particle2, particle3)

	// Run update
	ps.Update(0.02, []Entity{}, eventDispatcher)

	// Verify only particle2 was removed (expired)
	if len(ps.particles) != 2 {
		t.Errorf("Expected 2 particles to remain, got %d", len(ps.particles))
	}

	// Verify particle1 and particle3 are still active
	if !particle1.Active {
		t.Error("Expected particle1 to still be active")
	}
	if !particle3.Active {
		t.Error("Expected particle3 to still be active")
	}

	// Verify particle2 is inactive
	if particle2.Active {
		t.Error("Expected particle2 to be inactive")
	}
}

func TestParticleSystem_Update_ZeroDeltaTime(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create a particle
	particle := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	ps.particles = append(ps.particles, particle)

	initialX := particle.X
	initialY := particle.Y
	initialLife := particle.Life

	// Run update with zero delta time
	ps.Update(0.0, []Entity{}, eventDispatcher)

	// Verify particle position and life remain unchanged
	if math.Abs(particle.X-initialX) > particleTolerance {
		t.Errorf("Expected particle X to remain %f, got %f", initialX, particle.X)
	}
	if math.Abs(particle.Y-initialY) > particleTolerance {
		t.Errorf("Expected particle Y to remain %f, got %f", initialY, particle.Y)
	}
	if math.Abs(particle.Life-initialLife) > particleTolerance {
		t.Errorf("Expected particle life to remain %f, got %f", initialLife, particle.Life)
	}
}

func TestParticleSystem_Update_LargeDeltaTime(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create a particle
	particle := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	ps.particles = append(ps.particles, particle)

	// Run update with large delta time
	ps.Update(2.0, []Entity{}, eventDispatcher)

	// Verify particle was removed due to expiration
	if len(ps.particles) != 0 {
		t.Errorf("Expected particles count to be 0 after large delta time, got %d", len(ps.particles))
	}

	// Verify particle is inactive
	if particle.Active {
		t.Error("Expected particle to be inactive after large delta time")
	}
}

func TestParticleSystem_Update_NegativeDeltaTime(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create a particle
	particle := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	ps.particles = append(ps.particles, particle)

	initialX := particle.X
	initialY := particle.Y
	initialLife := particle.Life

	// Run update with negative delta time
	ps.Update(-0.016, []Entity{}, eventDispatcher)

	// Verify particle position and life are updated correctly with negative delta time
	expectedX := initialX + 10.0*(-0.016)
	expectedY := initialY + 5.0*(-0.016)
	expectedLife := initialLife - (-0.016) // Life increases with negative delta time

	if math.Abs(particle.X-expectedX) > particleTolerance {
		t.Errorf("Expected particle X to be %f, got %f", expectedX, particle.X)
	}
	if math.Abs(particle.Y-expectedY) > particleTolerance {
		t.Errorf("Expected particle Y to be %f, got %f", expectedY, particle.Y)
	}
	if math.Abs(particle.Life-expectedLife) > particleTolerance {
		t.Errorf("Expected particle life to be %f, got %f", expectedLife, particle.Life)
	}
}

func TestParticleSystem_Draw(t *testing.T) {
	ps := NewParticleSystem()

	// Create a screen for testing
	screen := ebiten.NewImage(800, 600)

	// Test that Draw doesn't panic with no particles
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked with no particles: %v", r)
		}
	}()

	ps.Draw(screen)

	// If we get here, Draw executed without panicking
}

func TestParticleSystem_Draw_WithParticles(t *testing.T) {
	ps := NewParticleSystem()

	// Create a screen for testing
	screen := ebiten.NewImage(800, 600)

	// Create some particles
	particle1 := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	particle2 := components.NewParticle(200, 200, -5, 10, 0.5, color.RGBA{0, 255, 0, 255}, 3.0)
	ps.particles = append(ps.particles, particle1, particle2)

	// Test that Draw doesn't panic with particles
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked with particles: %v", r)
		}
	}()

	ps.Draw(screen)

	// If we get here, Draw executed without panicking
}

func TestParticleSystem_Draw_WithInactiveParticles(t *testing.T) {
	ps := NewParticleSystem()

	// Create a screen for testing
	screen := ebiten.NewImage(800, 600)

	// Create particles and mark one as inactive
	particle1 := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	particle2 := components.NewParticle(200, 200, -5, 10, 0.5, color.RGBA{0, 255, 0, 255}, 3.0)
	particle2.Active = false
	ps.particles = append(ps.particles, particle1, particle2)

	// Test that Draw doesn't panic with inactive particles
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked with inactive particles: %v", r)
		}
	}()

	ps.Draw(screen)

	// If we get here, Draw executed without panicking
}

func TestParticleSystem_Draw_NilScreen(t *testing.T) {
	ps := NewParticleSystem()

	// Test that Draw handles nil screen gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked with nil screen: %v", r)
		}
	}()

	ps.Draw(nil)

	// If we get here, Draw handled nil screen gracefully (which is good)
}

func TestParticleSystem_CreatePacketCatchEffect(t *testing.T) {
	ps := NewParticleSystem()

	initialCount := len(ps.particles)
	packetColor := color.RGBA{255, 0, 0, 255}

	// Create packet catch effect
	ps.CreatePacketCatchEffect(100, 100, packetColor)

	// Verify particles were created
	if len(ps.particles) != initialCount+8 {
		t.Errorf("Expected %d particles, got %d", initialCount+8, len(ps.particles))
	}

	// Verify all new particles are active
	for i := initialCount; i < len(ps.particles); i++ {
		particle := ps.particles[i]
		if !particle.Active {
			t.Errorf("Expected particle %d to be active", i)
		}
		if particle.Color != packetColor {
			t.Errorf("Expected particle %d to have packet color, got %v", i, particle.Color)
		}
		if particle.Size < 2.0 || particle.Size > 5.0 {
			t.Errorf("Expected particle %d size to be between 2.0 and 5.0, got %f", i, particle.Size)
		}
		if particle.Life < 0.5 || particle.Life > 1.0 {
			t.Errorf("Expected particle %d life to be between 0.5 and 1.0, got %f", i, particle.Life)
		}
	}
}

func TestParticleSystem_CreatePowerUpEffect(t *testing.T) {
	ps := NewParticleSystem()

	initialCount := len(ps.particles)
	powerUpColor := color.RGBA{0, 255, 0, 255}

	// Create power-up effect
	ps.CreatePowerUpEffect(200, 200, powerUpColor)

	// Verify particles were created
	if len(ps.particles) != initialCount+12 {
		t.Errorf("Expected %d particles, got %d", initialCount+12, len(ps.particles))
	}

	// Verify all new particles are active
	for i := initialCount; i < len(ps.particles); i++ {
		particle := ps.particles[i]
		if !particle.Active {
			t.Errorf("Expected particle %d to be active", i)
		}
		if particle.Color != powerUpColor {
			t.Errorf("Expected particle %d to have power-up color, got %v", i, particle.Color)
		}
		if particle.Size < 1.0 || particle.Size > 3.0 {
			t.Errorf("Expected particle %d size to be between 1.0 and 3.0, got %f", i, particle.Size)
		}
		if particle.Life < 1.0 || particle.Life > 2.0 {
			t.Errorf("Expected particle %d life to be between 1.0 and 2.0, got %f", i, particle.Life)
		}
	}
}

func TestParticleSystem_Initialize(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test that Initialize doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Initialize panicked: %v", r)
		}
	}()

	ps.Initialize(eventDispatcher)

	// If we get here, Initialize executed without panicking
}

func TestParticleSystem_EventHandling_PacketCaught(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	ps.Initialize(eventDispatcher)

	// Create a mock packet entity
	packetEntity := createTestPacketEntity(1, 100, 100, color.RGBA{255, 0, 0, 255})

	initialCount := len(ps.particles)

	// Publish packet caught event
	eventData := &events.EventData{
		Packet: packetEntity,
	}
	event := events.NewEvent(events.EventPacketCaught, eventData)
	eventDispatcher.Publish(event)

	// Verify particles were created
	if len(ps.particles) != initialCount+8 {
		t.Errorf("Expected %d particles after packet caught event, got %d", initialCount+8, len(ps.particles))
	}
}

func TestParticleSystem_EventHandling_PowerUpCollected(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	ps.Initialize(eventDispatcher)

	// Create a mock power-up entity
	powerUpEntity := createTestPowerUpEntity(1, 200, 200, color.RGBA{0, 255, 0, 255})

	initialCount := len(ps.particles)

	// Publish power-up collected event
	eventData := &events.EventData{
		Packet: powerUpEntity,
	}
	event := events.NewEvent(events.EventPowerUpCollected, eventData)
	eventDispatcher.Publish(event)

	// Verify particles were created
	if len(ps.particles) != initialCount+12 {
		t.Errorf("Expected %d particles after power-up collected event, got %d", initialCount+12, len(ps.particles))
	}
}

func TestParticleSystem_EventHandling_InvalidEntity(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	ps.Initialize(eventDispatcher)

	initialCount := len(ps.particles)

	// Publish event with nil packet
	eventData := &events.EventData{
		Packet: nil,
	}
	event := events.NewEvent(events.EventPacketCaught, eventData)
	eventDispatcher.Publish(event)

	// Verify no particles were created
	if len(ps.particles) != initialCount {
		t.Errorf("Expected particles count to remain %d, got %d", initialCount, len(ps.particles))
	}
}

func TestParticleSystem_EventHandling_EntityWithoutComponents(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	ps.Initialize(eventDispatcher)

	// Create entity without required components
	entity := entities.NewEntity(1)

	initialCount := len(ps.particles)

	// Publish event with entity without components
	eventData := &events.EventData{
		Packet: entity,
	}
	event := events.NewEvent(events.EventPacketCaught, eventData)
	eventDispatcher.Publish(event)

	// Verify no particles were created
	if len(ps.particles) != initialCount {
		t.Errorf("Expected particles count to remain %d, got %d", initialCount, len(ps.particles))
	}
}

func TestParticleSystem_Integration(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	ps.Initialize(eventDispatcher)

	// Create entities
	packetEntity := createTestPacketEntity(1, 100, 100, color.RGBA{255, 0, 0, 255})
	powerUpEntity := createTestPowerUpEntity(2, 200, 200, color.RGBA{0, 255, 0, 255})

	// Publish events
	packetEvent := events.NewEvent(events.EventPacketCaught, &events.EventData{Packet: packetEntity})
	powerUpEvent := events.NewEvent(events.EventPowerUpCollected, &events.EventData{Packet: powerUpEntity})

	eventDispatcher.Publish(packetEvent)
	eventDispatcher.Publish(powerUpEvent)

	// Verify particles were created
	expectedParticles := 8 + 12 // 8 from packet catch + 12 from power-up
	if len(ps.particles) != expectedParticles {
		t.Errorf("Expected %d particles, got %d", expectedParticles, len(ps.particles))
	}

	// Update particles
	ps.Update(0.016, []Entity{}, eventDispatcher)

	// Verify particles are still active
	activeCount := 0
	for _, particle := range ps.particles {
		if particle.Active {
			activeCount++
		}
	}

	if activeCount != expectedParticles {
		t.Errorf("Expected %d active particles, got %d", expectedParticles, activeCount)
	}

	// Test drawing
	screen := ebiten.NewImage(800, 600)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked in integration test: %v", r)
		}
	}()

	ps.Draw(screen)
}

// Helper functions to create test entities

func createTestPacketEntity(id uint64, x, y float64, color color.RGBA) Entity {
	entity := entities.NewEntity(id)
	transform := components.NewTransform(x, y)
	sprite := components.NewSprite(15, 15, color)
	entity.AddComponent(transform)
	entity.AddComponent(sprite)
	return entity
}

func createTestPowerUpEntity(id uint64, x, y float64, color color.RGBA) Entity {
	entity := entities.NewEntity(id)
	transform := components.NewTransform(x, y)
	sprite := components.NewSprite(15, 15, color)
	entity.AddComponent(transform)
	entity.AddComponent(sprite)
	return entity
}
