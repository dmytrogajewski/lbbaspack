package systems

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// TestDrawableSystemInterface demonstrates compile-time type checking
func TestDrawableSystemInterface(t *testing.T) {
	t.Run("Compile-time Type Checking", func(t *testing.T) {
		// Create drawable systems
		routingSys := NewRoutingSystem()
		particleSys := NewParticleSystem()

		// These should compile without errors if they implement DrawableSystem correctly
		var _ DrawableSystem = routingSys
		var _ DrawableSystem = particleSys

		// Test that they can be used in a generic context
		testDrawableSystem := func(ds DrawableSystem) {
			// This function can only accept systems that implement DrawableSystem
			// If a system doesn't implement the correct Draw method signature,
			// this will fail at compile time
		}

		testDrawableSystem(routingSys)
		testDrawableSystem(particleSys)

		t.Log("All drawable systems implement the correct interface")
	})
}

// TestDrawableSystemRegistry demonstrates type-safe registration
func TestDrawableSystemRegistry(t *testing.T) {
	t.Run("Type-Safe Registration", func(t *testing.T) {
		registry := NewDrawableSystemRegistry()

		// Create drawable systems
		routingSys := NewRoutingSystem()
		particleSys := NewParticleSystem()

		// Register systems - this will fail at compile time if they don't implement DrawableSystem
		err := registry.RegisterDrawableSystem(SystemTypeRouting, routingSys)
		if err != nil {
			t.Errorf("Failed to register routing system: %v", err)
		}

		err = registry.RegisterDrawableSystem(SystemTypeParticle, particleSys)
		if err != nil {
			t.Errorf("Failed to register particle system: %v", err)
		}

		// Retrieve systems with type safety
		routing, exists := registry.GetDrawableSystem(SystemTypeRouting)
		if !exists {
			t.Error("Routing system not found in registry")
		}
		if routing == nil {
			t.Error("Retrieved routing system is nil")
		}

		particle, exists := registry.GetDrawableSystem(SystemTypeParticle)
		if !exists {
			t.Error("Particle system not found in registry")
		}
		if particle == nil {
			t.Error("Retrieved particle system is nil")
		}

		// Test that retrieved systems can be used as DrawableSystem
		screen := ebiten.NewImage(800, 600)
		defer screen.Dispose()

		entities := []Entity{}
		routing.Draw(screen, entities)
		particle.Draw(screen, entities)

		t.Log("All drawable systems work correctly with type safety")
	})
}

// TestValidateDrawableSystem demonstrates runtime validation
func TestValidateDrawableSystem(t *testing.T) {
	t.Run("Runtime Validation", func(t *testing.T) {
		// Test with valid drawable systems
		routingSys := NewRoutingSystem()
		particleSys := NewParticleSystem()

		drawable, ok := ValidateDrawableSystem(routingSys)
		if !ok {
			t.Error("Routing system should be a valid drawable system")
		}
		if drawable == nil {
			t.Error("Validated routing system should not be nil")
		}

		drawable, ok = ValidateDrawableSystem(particleSys)
		if !ok {
			t.Error("Particle system should be a valid drawable system")
		}
		if drawable == nil {
			t.Error("Validated particle system should not be nil")
		}

		// Test with non-drawable system
		spawnSys := NewSpawnSystem(func() Entity { return nil })
		drawable, ok = ValidateDrawableSystem(spawnSys)
		if ok {
			t.Error("Spawn system should not be a valid drawable system")
		}
		if drawable != nil {
			t.Error("Non-drawable system should return nil")
		}

		t.Log("Runtime validation works correctly")
	})
}

// TestDrawableSystemGeneric demonstrates generic type safety
func TestDrawableSystemGeneric(t *testing.T) {
	t.Run("Generic Type Safety", func(t *testing.T) {
		// This test demonstrates how generics can be used for additional type safety
		// The generic interface ensures that the Draw method signature is correct

		// This would fail at compile time if the Draw method signature is wrong
		routingSys := NewRoutingSystem()
		particleSys := NewParticleSystem()

		// These type assertions will fail at compile time if the systems don't implement
		// the correct Draw method signature
		var _ DrawableSystemGeneric[*ebiten.Image] = routingSys
		var _ DrawableSystemGeneric[*ebiten.Image] = particleSys

		t.Log("Generic type safety demonstrated - all systems implement correct interface")
	})
}
