package systems

import (
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewRenderSystem(t *testing.T) {
	rs := NewRenderSystem()

	// Test that the system is properly initialized
	if rs == nil {
		t.Fatal("NewRenderSystem returned nil")
	}

	// Test required components
	expectedComponents := []string{"Transform", "Sprite"}
	if len(rs.RequiredComponents) != len(expectedComponents) {
		t.Errorf("Expected %d required components, got %d", len(expectedComponents), len(rs.RequiredComponents))
	}

	for i, component := range expectedComponents {
		if rs.RequiredComponents[i] != component {
			t.Errorf("Expected required component %s at index %d, got %s", component, i, rs.RequiredComponents[i])
		}
	}

	// Test call count initialization
	if rs.callCount != 0 {
		t.Errorf("Expected initial call count to be 0, got %d", rs.callCount)
	}
}

func TestRenderSystem_UpdateWithScreen_NoEntities(t *testing.T) {
	rs := NewRenderSystem()
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}

	// Create a mock screen (nil for testing purposes)
	var screen *ebiten.Image = nil

	// Test that UpdateWithScreen panics with nil screen (expected behavior)
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil screen
			t.Logf("Expected panic with nil screen: %v", r)
		} else {
			t.Error("Expected panic with nil screen, but no panic occurred")
		}
	}()

	rs.UpdateWithScreen(0.016, entities, eventDispatcher, screen)
}

func TestRenderSystem_UpdateWithScreen_EntityWithoutRequiredComponents(t *testing.T) {
	rs := NewRenderSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with only Transform component
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	entity.AddComponent(transform)

	entities := []Entity{entity}
	var screen *ebiten.Image = nil

	// Test that UpdateWithScreen panics with nil screen (expected behavior)
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil screen
			t.Logf("Expected panic with nil screen: %v", r)
		} else {
			t.Error("Expected panic with nil screen, but no panic occurred")
		}
	}()

	rs.UpdateWithScreen(0.016, entities, eventDispatcher, screen)
}

func TestRenderSystem_UpdateWithScreen_EntityWithBothComponents(t *testing.T) {
	rs := NewRenderSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with both required components
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(50, 30, color.RGBA{255, 0, 0, 255})
	entity.AddComponent(transform)
	entity.AddComponent(sprite)

	entities := []Entity{entity}
	var screen *ebiten.Image = nil

	// Test that UpdateWithScreen panics with nil screen (expected behavior)
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil screen
			t.Logf("Expected panic with nil screen: %v", r)
		} else {
			t.Error("Expected panic with nil screen, but no panic occurred")
		}
	}()

	rs.UpdateWithScreen(0.016, entities, eventDispatcher, screen)
}

func TestRenderSystem_UpdateWithScreen_EntityWithPacketType(t *testing.T) {
	rs := NewRenderSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with packet type
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(50, 30, color.RGBA{255, 0, 0, 255})
	collider := components.NewCollider(50, 30, "packet")
	packetType := components.NewPacketType("HTTP", 1)
	entity.AddComponent(transform)
	entity.AddComponent(sprite)
	entity.AddComponent(collider)
	entity.AddComponent(packetType)

	entities := []Entity{entity}
	var screen *ebiten.Image = nil

	// Test that UpdateWithScreen panics with nil screen (expected behavior)
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil screen
			t.Logf("Expected panic with nil screen: %v", r)
		} else {
			t.Error("Expected panic with nil screen, but no panic occurred")
		}
	}()

	rs.UpdateWithScreen(0.016, entities, eventDispatcher, screen)
}

func TestRenderSystem_UpdateWithScreen_EntityWithLoadBalancer(t *testing.T) {
	rs := NewRenderSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with load balancer tag
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(50, 30, color.RGBA{255, 0, 0, 255})
	collider := components.NewCollider(50, 30, "loadbalancer")
	entity.AddComponent(transform)
	entity.AddComponent(sprite)
	entity.AddComponent(collider)

	entities := []Entity{entity}
	var screen *ebiten.Image = nil

	// Test that UpdateWithScreen panics with nil screen (expected behavior)
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil screen
			t.Logf("Expected panic with nil screen: %v", r)
		} else {
			t.Error("Expected panic with nil screen, but no panic occurred")
		}
	}()

	rs.UpdateWithScreen(0.016, entities, eventDispatcher, screen)
}

func TestRenderSystem_UpdateWithScreen_EntityWithPowerUp(t *testing.T) {
	rs := NewRenderSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with power-up
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(50, 30, color.RGBA{255, 0, 0, 255})
	powerUpType := components.NewPowerUpType("SpeedBoost", 15.0)
	entity.AddComponent(transform)
	entity.AddComponent(sprite)
	entity.AddComponent(powerUpType)

	entities := []Entity{entity}
	var screen *ebiten.Image = nil

	// Test that UpdateWithScreen panics with nil screen (expected behavior)
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil screen
			t.Logf("Expected panic with nil screen: %v", r)
		} else {
			t.Error("Expected panic with nil screen, but no panic occurred")
		}
	}()

	rs.UpdateWithScreen(0.016, entities, eventDispatcher, screen)
}

func TestRenderSystem_UpdateWithScreen_EntityWithBackendAssignment(t *testing.T) {
	rs := NewRenderSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with backend assignment
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(50, 30, color.RGBA{255, 0, 0, 255})
	backendAssignment := components.NewBackendAssignment(1)
	entity.AddComponent(transform)
	entity.AddComponent(sprite)
	entity.AddComponent(backendAssignment)

	entities := []Entity{entity}
	var screen *ebiten.Image = nil

	// Test that UpdateWithScreen panics with nil screen (expected behavior)
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil screen
			t.Logf("Expected panic with nil screen: %v", r)
		} else {
			t.Error("Expected panic with nil screen, but no panic occurred")
		}
	}()

	rs.UpdateWithScreen(0.016, entities, eventDispatcher, screen)
}

func TestRenderSystem_UpdateWithScreen_MultipleEntities(t *testing.T) {
	rs := NewRenderSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create multiple entities
	entity1 := entities.NewEntity(1)
	transform1 := components.NewTransform(100, 200)
	sprite1 := components.NewSprite(50, 30, color.RGBA{255, 0, 0, 255})
	entity1.AddComponent(transform1)
	entity1.AddComponent(sprite1)

	entity2 := entities.NewEntity(2)
	transform2 := components.NewTransform(300, 400)
	sprite2 := components.NewSprite(60, 40, color.RGBA{0, 255, 0, 255})
	entity2.AddComponent(transform2)
	entity2.AddComponent(sprite2)

	entities := []Entity{entity1, entity2}
	var screen *ebiten.Image = nil

	// Test that UpdateWithScreen panics with nil screen (expected behavior)
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil screen
			t.Logf("Expected panic with nil screen: %v", r)
		} else {
			t.Error("Expected panic with nil screen, but no panic occurred")
		}
	}()

	rs.UpdateWithScreen(0.016, entities, eventDispatcher, screen)
}

func TestRenderSystem_UpdateWithScreen_ZeroDeltaTime(t *testing.T) {
	rs := NewRenderSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with both required components
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(50, 30, color.RGBA{255, 0, 0, 255})
	entity.AddComponent(transform)
	entity.AddComponent(sprite)

	entities := []Entity{entity}
	var screen *ebiten.Image = nil

	// Test that UpdateWithScreen panics with nil screen (expected behavior)
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil screen
			t.Logf("Expected panic with nil screen: %v", r)
		} else {
			t.Error("Expected panic with nil screen, but no panic occurred")
		}
	}()

	rs.UpdateWithScreen(0.0, entities, eventDispatcher, screen)
}

func TestRenderSystem_UpdateWithScreen_NilScreen(t *testing.T) {
	rs := NewRenderSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with both required components
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(50, 30, color.RGBA{255, 0, 0, 255})
	entity.AddComponent(transform)
	entity.AddComponent(sprite)

	entities := []Entity{entity}

	// Test that UpdateWithScreen panics with nil screen (expected behavior)
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil screen
			t.Logf("Expected panic with nil screen: %v", r)
		} else {
			t.Error("Expected panic with nil screen, but no panic occurred")
		}
	}()

	rs.UpdateWithScreen(0.016, entities, eventDispatcher, nil)
}

func TestRenderSystem_Integration(t *testing.T) {
	rs := NewRenderSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create multiple entities with different configurations
	entity1 := entities.NewEntity(1) // Packet entity
	transform1 := components.NewTransform(100, 200)
	sprite1 := components.NewSprite(50, 30, color.RGBA{255, 0, 0, 255})
	collider1 := components.NewCollider(50, 30, "packet")
	packetType1 := components.NewPacketType("HTTP", 1)
	entity1.AddComponent(transform1)
	entity1.AddComponent(sprite1)
	entity1.AddComponent(collider1)
	entity1.AddComponent(packetType1)

	entity2 := entities.NewEntity(2) // Load balancer entity
	transform2 := components.NewTransform(300, 400)
	sprite2 := components.NewSprite(60, 40, color.RGBA{0, 255, 0, 255})
	collider2 := components.NewCollider(60, 40, "loadbalancer")
	entity2.AddComponent(transform2)
	entity2.AddComponent(sprite2)
	entity2.AddComponent(collider2)

	entity3 := entities.NewEntity(3) // Power-up entity
	transform3 := components.NewTransform(500, 600)
	sprite3 := components.NewSprite(40, 40, color.RGBA{0, 0, 255, 255})
	powerUpType3 := components.NewPowerUpType("SpeedBoost", 15.0)
	entity3.AddComponent(transform3)
	entity3.AddComponent(sprite3)
	entity3.AddComponent(powerUpType3)

	entity4 := entities.NewEntity(4) // Backend entity
	transform4 := components.NewTransform(700, 800)
	sprite4 := components.NewSprite(45, 35, color.RGBA{255, 255, 0, 255})
	backendAssignment4 := components.NewBackendAssignment(1)
	entity4.AddComponent(transform4)
	entity4.AddComponent(sprite4)
	entity4.AddComponent(backendAssignment4)

	entities := []Entity{entity1, entity2, entity3, entity4}
	var screen *ebiten.Image = nil

	// Test that UpdateWithScreen panics with nil screen (expected behavior)
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil screen
			t.Logf("Expected panic with nil screen: %v", r)
		} else {
			t.Error("Expected panic with nil screen, but no panic occurred")
		}
	}()

	rs.UpdateWithScreen(0.016, entities, eventDispatcher, screen)
}
