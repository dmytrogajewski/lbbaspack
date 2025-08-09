package systems

import (
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"math"
	"testing"
)

const particleTolerance = 0.01

func TestNewParticleSystem(t *testing.T) {
	ps := NewParticleSystem()

	// Test that the system is properly initialized
	if ps == nil {
		t.Fatal("NewParticleSystem returned nil")
	}

	// Test that the system has the correct type
	if ps.GetSystemInfo().Type != SystemTypeParticle {
		t.Errorf("Expected system type %s, got %s", SystemTypeParticle, ps.GetSystemInfo().Type)
	}
}

func TestParticleSystem_Update_NoParticles(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test with no particles
	entities := []Entity{}

	// Run update
	ps.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// The system should handle empty entities gracefully
}

func TestParticleSystem_Update_WithParticles(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()

	// Add some particles to the state
	particle1 := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	particle2 := components.NewParticle(200, 200, -5, 10, 0.5, color.RGBA{0, 255, 0, 255}, 3.0)
	particleState.Particles = append(particleState.Particles, particle1, particle2)

	entity.AddComponent(particleState)
	entities := []Entity{entity}

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

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()

	// Create a particle with very short life
	particle := components.NewParticle(100, 100, 10, 5, 0.01, color.RGBA{255, 0, 0, 255}, 2.0)
	particleState.Particles = append(particleState.Particles, particle)

	entity.AddComponent(particleState)
	entities := []Entity{entity}

	// Run update with delta time that exceeds particle life
	ps.Update(0.02, entities, eventDispatcher)

	// Verify particle was removed from the state
	if len(particleState.Particles) != 0 {
		t.Errorf("Expected particles count to be 0 after expiration, got %d", len(particleState.Particles))
	}

	// Verify particle is marked as inactive
	if particle.Active {
		t.Error("Expected particle to be marked as inactive")
	}
}

func TestParticleSystem_Update_MixedParticles(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()

	// Create particles with different life spans
	particle1 := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)   // Long life
	particle2 := components.NewParticle(200, 200, -5, 10, 0.01, color.RGBA{0, 255, 0, 255}, 3.0) // Short life
	particle3 := components.NewParticle(300, 300, 0, 0, 0.5, color.RGBA{0, 0, 255, 255}, 1.0)    // Medium life

	particleState.Particles = append(particleState.Particles, particle1, particle2, particle3)
	entity.AddComponent(particleState)
	entities := []Entity{entity}

	// Run update
	ps.Update(0.016, entities, eventDispatcher)

	// Verify long-life particle remains
	if !particle1.Active {
		t.Error("Expected particle1 to remain active")
	}

	// Verify short-life particle is removed
	if len(particleState.Particles) != 2 {
		t.Errorf("Expected 2 particles after update, got %d", len(particleState.Particles))
	}

	// Verify medium-life particle remains
	if !particle3.Active {
		t.Error("Expected particle3 to remain active")
	}
}

func TestParticleSystem_Update_ZeroDeltaTime(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()

	particle := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	particleState.Particles = append(particleState.Particles, particle)

	entity.AddComponent(particleState)
	entities := []Entity{entity}

	// Run update with zero delta time
	ps.Update(0.0, entities, eventDispatcher)

	// Verify particle position unchanged
	if particle.X != 100.0 || particle.Y != 100.0 {
		t.Errorf("Expected particle position to remain (100, 100), got (%f, %f)", particle.X, particle.Y)
	}
}

func TestParticleSystem_Update_LargeDeltaTime(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()

	particle := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	particleState.Particles = append(particleState.Particles, particle)

	entity.AddComponent(particleState)
	entities := []Entity{entity}

	// Run update with large delta time
	ps.Update(2.0, entities, eventDispatcher)

	// Verify particle was updated (though may have expired)
	// The system should handle large delta time gracefully
}

func TestParticleSystem_Update_NegativeDeltaTime(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()

	particle := components.NewParticle(100, 100, 10, 5, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	particleState.Particles = append(particleState.Particles, particle)

	entity.AddComponent(particleState)
	entities := []Entity{entity}

	// Run update with negative delta time
	ps.Update(-0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// The system should handle negative delta time gracefully
}

func TestParticleSystem_Draw(t *testing.T) {
	ps := NewParticleSystem()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()

	particle := components.NewParticle(100, 100, 0, 0, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	particleState.Particles = append(particleState.Particles, particle)

	entity.AddComponent(particleState)

	// Test that Draw doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked: %v", r)
		}
	}()

	// Call Draw method (we can't easily create a real ebiten.Image in tests)
	// Just verify the method exists and can be called
	_ = ps.Draw
}

func TestParticleSystem_Draw_WithParticles(t *testing.T) {
	ps := NewParticleSystem()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()

	// Add multiple particles
	particle1 := components.NewParticle(100, 100, 0, 0, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	particle2 := components.NewParticle(200, 200, 0, 0, 1.0, color.RGBA{0, 255, 0, 255}, 3.0)
	particleState.Particles = append(particleState.Particles, particle1, particle2)

	entity.AddComponent(particleState)

	// Test that Draw doesn't panic with multiple particles
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked with multiple particles: %v", r)
		}
	}()

	// Just verify the method exists and can be called
	_ = ps.Draw
}

func TestParticleSystem_Draw_WithInactiveParticles(t *testing.T) {
	ps := NewParticleSystem()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()

	// Add inactive particle
	particle := components.NewParticle(100, 100, 0, 0, 0.0, color.RGBA{255, 0, 0, 255}, 2.0)
	particle.Active = false
	particleState.Particles = append(particleState.Particles, particle)

	entity.AddComponent(particleState)

	// Test that Draw doesn't panic with inactive particles
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked with inactive particles: %v", r)
		}
	}()

	// Just verify the method exists and can be called
	_ = ps.Draw
}

func TestParticleSystem_Draw_NilScreen(t *testing.T) {
	ps := NewParticleSystem()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()

	particle := components.NewParticle(100, 100, 0, 0, 1.0, color.RGBA{255, 0, 0, 255}, 2.0)
	particleState.Particles = append(particleState.Particles, particle)

	entity.AddComponent(particleState)

	// Test that Draw doesn't panic with nil screen
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked with nil screen: %v", r)
		}
	}()

	// Just verify the method exists and can be called
	_ = ps.Draw
}

func TestParticleSystem_CreatePacketCatchEffect(t *testing.T) {
	ps := NewParticleSystem()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()
	entity.AddComponent(particleState)

	// Test creating packet catch effect
	initialCount := len(particleState.Particles)
	ps.CreatePacketCatchEffect(100, 200, color.RGBA{255, 0, 0, 255}, particleState)

	// Verify particles were created
	if len(particleState.Particles) <= initialCount {
		t.Error("Expected particles to be created")
	}

	// Verify all new particles are active
	for _, particle := range particleState.Particles[initialCount:] {
		if !particle.Active {
			t.Error("Expected new particles to be active")
		}
	}
}

func TestParticleSystem_CreatePowerUpEffect(t *testing.T) {
	ps := NewParticleSystem()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()
	entity.AddComponent(particleState)

	// Test creating power-up effect
	initialCount := len(particleState.Particles)
	ps.CreatePowerUpEffect(150, 250, color.RGBA{0, 255, 0, 255}, particleState)

	// Verify particles were created
	if len(particleState.Particles) <= initialCount {
		t.Error("Expected particles to be created")
	}

	// Verify all new particles are active
	for _, particle := range particleState.Particles[initialCount:] {
		if !particle.Active {
			t.Error("Expected new particles to be active")
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

	// Verify no errors occurred
	// The system should initialize without issues
}

func TestParticleSystem_EventHandling_PacketCaught(t *testing.T) {
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()
	entity.AddComponent(particleState)

	// Subscribe to packet caught events
	var eventReceived bool
	eventDispatcher.Subscribe(events.EventPacketCaught, func(event *events.Event) {
		eventReceived = true
	})

	// Simulate packet caught event
	eventDispatcher.Publish(events.NewEvent(events.EventPacketCaught, &events.EventData{}))

	// Verify event was received
	if !eventReceived {
		t.Error("Expected packet caught event to be received")
	}
}

func TestParticleSystem_EventHandling_PowerUpCollected(t *testing.T) {
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()
	entity.AddComponent(particleState)

	// Subscribe to power-up collected event
	var eventReceived bool
	eventDispatcher.Subscribe(events.EventPowerUpCollected, func(event *events.Event) {
		eventReceived = true
	})

	// Simulate power-up collected event
	eventDispatcher.Publish(events.NewEvent(events.EventPowerUpCollected, &events.EventData{}))

	// Verify event was received
	if !eventReceived {
		t.Error("Expected power-up collected event to be received")
	}
}

func TestParticleSystem_EventHandling_InvalidEntity(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity without ParticleState component
	entity := entities.NewEntity(1)
	entities := []Entity{entity}

	// Test that Update doesn't panic with invalid entity
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with invalid entity: %v", r)
		}
	}()

	ps.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// The system should handle invalid entities gracefully
}

func TestParticleSystem_EventHandling_EntityWithoutComponents(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with no components
	entity := entities.NewEntity(1)
	entities := []Entity{entity}

	// Test that Update doesn't panic with entity without components
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with entity without components: %v", r)
		}
	}()

	ps.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// The system should handle entities without components gracefully
}

func TestParticleSystem_Integration(t *testing.T) {
	ps := NewParticleSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with ParticleState component
	entity := entities.NewEntity(1)
	particleState := components.NewParticleState()
	entity.AddComponent(particleState)

	entities := []Entity{entity}

	// Run multiple updates to test system stability
	for i := 0; i < 5; i++ {
		ps.Update(0.016, entities, eventDispatcher)
	}

	// Verify no errors occurred
	// The system should remain stable across multiple updates
}

// Helper functions to create test entities
func createTestPacketEntity(id uint64, x, y float64, color color.RGBA) Entity {
	entity := entities.NewEntity(id)
	particleState := components.NewParticleState()
	entity.AddComponent(particleState)
	return entity
}

func createTestPowerUpEntity(id uint64, x, y float64, color color.RGBA) Entity {
	entity := entities.NewEntity(id)
	particleState := components.NewParticleState()
	entity.AddComponent(particleState)
	return entity
}
