package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"
)

func TestNewCollisionSystem(t *testing.T) {
	cs := NewCollisionSystem()

	// Test that the system is properly initialized
	if cs == nil {
		t.Fatal("NewCollisionSystem returned nil")
	}

	// Test required components
	expectedComponents := []string{"Transform", "Collider"}
	if len(cs.RequiredComponents) != len(expectedComponents) {
		t.Errorf("Expected %d required components, got %d", len(expectedComponents), len(cs.RequiredComponents))
	}

	for i, component := range expectedComponents {
		if cs.RequiredComponents[i] != component {
			t.Errorf("Expected required component %s at index %d, got %s", component, i, cs.RequiredComponents[i])
		}
	}

	// No internal score anymore
}

func TestCollisionSystem_Update_NoEntities(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test with no entities
	entities := []Entity{}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// No internal score anymore
}

func TestCollisionSystem_Update_NoLoadBalancer(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create packet entities without load balancer
	packet1 := createPacketEntity(1, 100, 100)
	packet2 := createPacketEntity(2, 200, 200)

	entities := []Entity{packet1, packet2}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// No internal score anymore
}

func TestCollisionSystem_Update_PacketCollision(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create load balancer at position (100, 100)
	loadBalancer := createLoadBalancerEntity(1, 100, 100)

	// Create packet that will collide with load balancer
	packet := createPacketEntity(2, 105, 105) // Overlapping position

	entities := []Entity{loadBalancer, packet}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// The collision system only detects collisions and publishes events
	// It doesn't deactivate entities - that's handled by other systems
	// Verify that the packet remains active (collision system doesn't modify entities)
	if !packet.IsActive() {
		t.Error("Expected packet to remain active - collision system only detects collisions")
	}
}

func TestCollisionSystem_Update_PacketNoCollision(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create load balancer at position (100, 100)
	loadBalancer := createLoadBalancerEntity(1, 100, 100)

	// Create packet that won't collide with load balancer
	packet := createPacketEntity(2, 200, 200) // Non-overlapping position

	entities := []Entity{loadBalancer, packet}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// No internal score anymore

	// Verify packet remains active
	if !packet.IsActive() {
		t.Error("Expected packet to remain active when no collision")
	}
}

func TestCollisionSystem_Update_PowerUpCollision(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create load balancer at position (100, 100)
	loadBalancer := createLoadBalancerEntity(1, 100, 100)

	// Create power-up that will collide with load balancer
	powerUp := createPowerUpEntity(2, 105, 105, "SpeedBoost")

	entities := []Entity{loadBalancer, powerUp}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// The collision system only detects collisions and publishes events
	// It doesn't deactivate entities - that's handled by other systems
	// Verify that the power-up remains active (collision system doesn't modify entities)
	if !powerUp.IsActive() {
		t.Error("Expected power-up to remain active - collision system only detects collisions")
	}
}

func TestCollisionSystem_Update_PowerUpNoCollision(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create load balancer at position (100, 100)
	loadBalancer := createLoadBalancerEntity(1, 100, 100)

	// Create power-up that won't collide with load balancer
	powerUp := createPowerUpEntity(2, 200, 200, "SpeedBoost")

	entities := []Entity{loadBalancer, powerUp}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify power-up remains active
	if !powerUp.IsActive() {
		t.Error("Expected power-up to remain active when no collision")
	}
}

func TestCollisionSystem_Update_PacketMissed(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create packet that falls off screen (Y > 600)
	packet := createPacketEntity(1, 100, 650)

	entities := []Entity{packet}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// The collision system only detects collisions between entities
	// It doesn't handle off-screen detection - that's handled by the OffscreenSystem
	// Verify that the packet remains active (collision system doesn't check screen bounds)
	if !packet.IsActive() {
		t.Error("Expected packet to remain active - collision system doesn't handle off-screen detection")
	}
}

func TestCollisionSystem_Update_PacketNotMissed(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create packet that's still on screen (Y <= 600)
	packet := createPacketEntity(1, 100, 500)

	entities := []Entity{packet}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify packet remains active
	if !packet.IsActive() {
		t.Error("Expected packet to remain active when still on screen")
	}
}

func TestCollisionSystem_Update_MultiplePackets(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create load balancer
	loadBalancer := createLoadBalancerEntity(1, 100, 100)

	// Create multiple packets - some colliding, some not
	packet1 := createPacketEntity(2, 105, 105) // Will collide
	packet2 := createPacketEntity(3, 200, 200) // Won't collide
	packet3 := createPacketEntity(4, 110, 110) // Will collide

	entities := []Entity{loadBalancer, packet1, packet2, packet3}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// The collision system only detects collisions and publishes events
	// It doesn't deactivate entities - that's handled by other systems
	// Verify that all packets remain active (collision system doesn't modify entities)
	if !packet1.IsActive() {
		t.Error("Expected packet1 to remain active - collision system only detects collisions")
	}
	if !packet3.IsActive() {
		t.Error("Expected packet3 to remain active - collision system only detects collisions")
	}

	// Verify non-colliding packet remains active
	if !packet2.IsActive() {
		t.Error("Expected packet2 to remain active when no collision")
	}
}

func TestCollisionSystem_Update_EntityWithoutRequiredComponents(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with only Transform (missing Collider)
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 100)
	entity.AddComponent(transform)

	entities := []Entity{entity}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// No internal score anymore
}

func TestCollisionSystem_Update_InactiveEntity(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create load balancer
	loadBalancer := createLoadBalancerEntity(1, 100, 100)

	// Create inactive packet
	packet := createPacketEntity(2, 105, 105)
	packet.(interface{ SetActive(bool) }).SetActive(false)

	entities := []Entity{loadBalancer, packet}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// No internal score anymore
}

func TestCollisionSystem_checkCollision_Overlapping(t *testing.T) {
	cs := NewCollisionSystem()

	// Create overlapping colliders
	transform1 := components.NewTransform(100, 100)
	collider1 := components.NewCollider(50, 50, "test1")

	transform2 := components.NewTransform(120, 120)
	collider2 := components.NewCollider(50, 50, "test2")

	// Check collision
	result := cs.checkCollision(transform1, collider1, transform2, collider2)

	// Should collide (overlapping)
	if !result {
		t.Error("Expected collision detection for overlapping objects")
	}
}

func TestCollisionSystem_checkCollision_NotOverlapping(t *testing.T) {
	cs := NewCollisionSystem()

	// Create non-overlapping colliders
	transform1 := components.NewTransform(100, 100)
	collider1 := components.NewCollider(50, 50, "test1")

	transform2 := components.NewTransform(200, 200)
	collider2 := components.NewCollider(50, 50, "test2")

	// Check collision
	result := cs.checkCollision(transform1, collider1, transform2, collider2)

	// Should not collide (not overlapping)
	if result {
		t.Error("Expected no collision detection for non-overlapping objects")
	}
}

func TestCollisionSystem_checkCollision_Touching(t *testing.T) {
	cs := NewCollisionSystem()

	// Create touching colliders (edge case)
	transform1 := components.NewTransform(100, 100)
	collider1 := components.NewCollider(50, 50, "test1")

	transform2 := components.NewTransform(150, 100) // Touching at x=150
	collider2 := components.NewCollider(50, 50, "test2")

	// Check collision
	result := cs.checkCollision(transform1, collider1, transform2, collider2)

	// Should collide (AABB considers touching as collision)
	if !result {
		t.Error("Expected collision detection for touching objects in AABB collision")
	}
}

func TestCollisionSystem_Integration(t *testing.T) {
	cs := NewCollisionSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create load balancer
	loadBalancer := createLoadBalancerEntity(1, 100, 100)

	// Create various entities
	packet1 := createPacketEntity(2, 105, 105)                // Will collide
	packet2 := createPacketEntity(3, 200, 200)                // Won't collide
	packet3 := createPacketEntity(4, 100, 700)                // Will fall off screen
	powerUp := createPowerUpEntity(5, 110, 110, "SpeedBoost") // Will collide

	entities := []Entity{loadBalancer, packet1, packet2, packet3, powerUp}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// The collision system only detects collisions and publishes events
	// It doesn't deactivate entities - that's handled by other systems
	// Verify that all entities remain active (collision system doesn't modify entities)
	if !packet1.IsActive() {
		t.Error("Expected packet1 to remain active - collision system only detects collisions")
	}
	if !powerUp.IsActive() {
		t.Error("Expected powerUp to remain active - collision system only detects collisions")
	}

	// Verify non-colliding packet remains active
	if !packet2.IsActive() {
		t.Error("Expected packet2 to remain active when no collision")
	}

	// Verify packet that would fall off screen remains active
	// (collision system doesn't handle off-screen detection)
	if !packet3.IsActive() {
		t.Error("Expected packet3 to remain active - collision system doesn't handle off-screen detection")
	}
}

// Helper functions to create test entities

func createLoadBalancerEntity(id uint64, x, y float64) Entity {
	entity := entities.NewEntity(id)
	transform := components.NewTransform(x, y)
	collider := components.NewCollider(50, 50, "loadbalancer")
	entity.AddComponent(transform)
	entity.AddComponent(collider)
	return entity
}

func createPacketEntity(id uint64, x, y float64) Entity {
	entity := entities.NewEntity(id)
	transform := components.NewTransform(x, y)
	collider := components.NewCollider(20, 20, "packet")
	packetType := components.NewPacketType("HTTP", 1)
	entity.AddComponent(transform)
	entity.AddComponent(collider)
	entity.AddComponent(packetType)
	return entity
}

func createPowerUpEntity(id uint64, x, y float64, name string) Entity {
	entity := entities.NewEntity(id)
	transform := components.NewTransform(x, y)
	collider := components.NewCollider(30, 30, "powerup")
	powerUpType := components.NewPowerUpType(name, 5.0)
	entity.AddComponent(transform)
	entity.AddComponent(collider)
	entity.AddComponent(powerUpType)
	return entity
}
