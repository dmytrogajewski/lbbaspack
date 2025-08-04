package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"math"
	"testing"
)

const tolerance = 0.0001

func TestNewMovementSystem(t *testing.T) {
	ms := NewMovementSystem()

	// Test that the system is properly initialized
	if ms == nil {
		t.Fatal("NewMovementSystem returned nil")
	}

	// Test required components
	expectedComponents := []string{"Transform", "Physics"}
	if len(ms.RequiredComponents) != len(expectedComponents) {
		t.Errorf("Expected %d required components, got %d", len(expectedComponents), len(ms.RequiredComponents))
	}

	for i, component := range expectedComponents {
		if ms.RequiredComponents[i] != component {
			t.Errorf("Expected required component %s at index %d, got %s", component, i, ms.RequiredComponents[i])
		}
	}

	// Test initial call count
	if ms.callCount != 0 {
		t.Errorf("Expected initial callCount to be 0, got %d", ms.callCount)
	}
}

func TestMovementSystem_Update_NoEntities(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test with no entities
	entities := []Entity{}

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify call count was incremented
	if ms.callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", ms.callCount)
	}
}

func TestMovementSystem_Update_EntityWithoutRequiredComponents(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with only Transform (missing Physics)
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 100)
	entity.AddComponent(transform)

	entities := []Entity{entity}

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify call count was incremented
	if ms.callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", ms.callCount)
	}

	// Verify transform position remains unchanged
	// (since no physics component, movement should not be processed)
	transformComp := entity.GetTransform()
	if transformComp == nil {
		t.Fatal("Expected transform component to exist")
	}

	transformObj := transformComp.(components.TransformComponent)
	x, y := transformObj.GetX(), transformObj.GetY()
	if math.Abs(x-100) > tolerance || math.Abs(y-100) > tolerance {
		t.Errorf("Expected position to remain (100, 100), got (%f, %f)", x, y)
	}
}

func TestMovementSystem_Update_EntityWithPhysicsOnly(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with only Physics (missing Transform)
	entity := entities.NewEntity(1)
	physics := components.NewPhysics()
	entity.AddComponent(physics)

	entities := []Entity{entity}

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify call count was incremented
	if ms.callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", ms.callCount)
	}

	// The system should skip entities without required components
}

func TestMovementSystem_Update_EntityWithBothComponents(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with both Transform and Physics
	entity := createMovementEntity(1, 100, 100, 10, 5)

	entities := []Entity{entity}

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify call count was incremented
	if ms.callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", ms.callCount)
	}

	// Verify transform position was updated based on physics
	transformComp := entity.GetTransform()
	transformObj := transformComp.(components.TransformComponent)
	x, y := transformObj.GetX(), transformObj.GetY()

	// Expected position: (100 + 10*0.016, 100 + 5*0.016) = (100.16, 100.08)
	expectedX := 100.0 + 10.0*0.016
	expectedY := 100.0 + 5.0*0.016

	if math.Abs(x-expectedX) > tolerance {
		t.Errorf("Expected X position to be %f, got %f", expectedX, x)
	}

	if math.Abs(y-expectedY) > tolerance {
		t.Errorf("Expected Y position to be %f, got %f", expectedY, y)
	}
}

func TestMovementSystem_Update_EntityWithZeroVelocity(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with zero velocity
	entity := createMovementEntity(1, 100, 100, 0, 0)

	entities := []Entity{entity}

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify transform position remains unchanged
	transformComp := entity.GetTransform()
	transformObj := transformComp.(components.TransformComponent)
	x, y := transformObj.GetX(), transformObj.GetY()

	if math.Abs(x-100) > tolerance || math.Abs(y-100) > tolerance {
		t.Errorf("Expected position to remain (100, 100), got (%f, %f)", x, y)
	}
}

func TestMovementSystem_Update_EntityWithNegativeVelocity(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with negative velocity
	entity := createMovementEntity(1, 100, 100, -10, -5)

	entities := []Entity{entity}

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify transform position was updated correctly
	transformComp := entity.GetTransform()
	transformObj := transformComp.(components.TransformComponent)
	x, y := transformObj.GetX(), transformObj.GetY()

	// Expected position: (100 + (-10)*0.016, 100 + (-5)*0.016) = (99.84, 99.92)
	expectedX := 100.0 + (-10.0)*0.016
	expectedY := 100.0 + (-5.0)*0.016

	if math.Abs(x-expectedX) > tolerance {
		t.Errorf("Expected X position to be %f, got %f", expectedX, x)
	}

	if math.Abs(y-expectedY) > tolerance {
		t.Errorf("Expected Y position to be %f, got %f", expectedY, y)
	}
}

func TestMovementSystem_Update_MultipleEntities(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create multiple entities with different velocities
	entity1 := createMovementEntity(1, 100, 100, 10, 5)  // Moving right and down
	entity2 := createMovementEntity(2, 200, 200, -5, 10) // Moving left and down
	entity3 := createMovementEntity(3, 300, 300, 0, 0)   // Stationary

	entities := []Entity{entity1, entity2, entity3}

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify all entities were processed correctly
	transform1 := entity1.GetTransform().(components.TransformComponent)
	x1, y1 := transform1.GetX(), transform1.GetY()
	expectedX1 := 100.0 + 10.0*0.016
	expectedY1 := 100.0 + 5.0*0.016
	if math.Abs(x1-expectedX1) > tolerance || math.Abs(y1-expectedY1) > tolerance {
		t.Errorf("Entity1: Expected position (%f, %f), got (%f, %f)", expectedX1, expectedY1, x1, y1)
	}

	transform2 := entity2.GetTransform().(components.TransformComponent)
	x2, y2 := transform2.GetX(), transform2.GetY()
	expectedX2 := 200.0 + (-5.0)*0.016
	expectedY2 := 200.0 + 10.0*0.016
	if math.Abs(x2-expectedX2) > tolerance || math.Abs(y2-expectedY2) > tolerance {
		t.Errorf("Entity2: Expected position (%f, %f), got (%f, %f)", expectedX2, expectedY2, x2, y2)
	}

	transform3 := entity3.GetTransform().(components.TransformComponent)
	x3, y3 := transform3.GetX(), transform3.GetY()
	if math.Abs(x3-300) > tolerance || math.Abs(y3-300) > tolerance {
		t.Errorf("Entity3: Expected position (300, 300), got (%f, %f)", x3, y3)
	}
}

func TestMovementSystem_Update_EntityWithPacketType(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with Transform, Physics, and PacketType
	entity := createMovementEntity(1, 100, 100, 10, 5)
	packetType := components.NewPacketType("TestPacket", 1)
	entityObj := entity.(*entities.Entity)
	entityObj.AddComponent(packetType)

	entities := []Entity{entity}

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify call count was incremented
	if ms.callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", ms.callCount)
	}

	// Verify transform position was updated
	transformComp := entity.GetTransform()
	transformObj := transformComp.(components.TransformComponent)
	x, y := transformObj.GetX(), transformObj.GetY()

	expectedX := 100.0 + 10.0*0.016
	expectedY := 100.0 + 5.0*0.016

	if math.Abs(x-expectedX) > tolerance {
		t.Errorf("Expected X position to be %f, got %f", expectedX, x)
	}

	if math.Abs(y-expectedY) > tolerance {
		t.Errorf("Expected Y position to be %f, got %f", expectedY, y)
	}
}

func TestMovementSystem_Update_EntityWithInvalidComponentTypes(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with invalid component types
	entity := entities.NewEntity(1)

	// Add components that don't implement the required interfaces
	entity.AddComponent(&mockComponent{componentType: "Transform"})
	entity.AddComponent(&mockComponent{componentType: "Physics"})

	entities := []Entity{entity}

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify call count was incremented
	if ms.callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", ms.callCount)
	}

	// The system should skip entities with invalid component types
}

func TestMovementSystem_Update_EntityWithNullComponents(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with null components (simulated by not adding any)
	entity := entities.NewEntity(1)

	entities := []Entity{entity}

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify call count was incremented
	if ms.callCount != 1 {
		t.Errorf("Expected callCount to be 1, got %d", ms.callCount)
	}

	// The system should handle null components gracefully
}

func TestMovementSystem_Update_ZeroDeltaTime(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with velocity
	entity := createMovementEntity(1, 100, 100, 10, 5)

	entities := []Entity{entity}

	// Run update with zero delta time
	ms.Update(0.0, entities, eventDispatcher)

	// Verify transform position remains unchanged
	transformComp := entity.GetTransform()
	transformObj := transformComp.(components.TransformComponent)
	x, y := transformObj.GetX(), transformObj.GetY()

	if math.Abs(x-100) > tolerance || math.Abs(y-100) > tolerance {
		t.Errorf("Expected position to remain (100, 100), got (%f, %f)", x, y)
	}
}

func TestMovementSystem_Update_LargeDeltaTime(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with velocity
	entity := createMovementEntity(1, 100, 100, 10, 5)

	entities := []Entity{entity}

	// Run update with large delta time
	ms.Update(1.0, entities, eventDispatcher)

	// Verify transform position was updated with large movement
	transformComp := entity.GetTransform()
	transformObj := transformComp.(components.TransformComponent)
	x, y := transformObj.GetX(), transformObj.GetY()

	// Expected position: (100 + 10*1.0, 100 + 5*1.0) = (110, 105)
	expectedX := 100.0 + 10.0*1.0
	expectedY := 100.0 + 5.0*1.0

	if math.Abs(x-expectedX) > tolerance {
		t.Errorf("Expected X position to be %f, got %f", expectedX, x)
	}

	if math.Abs(y-expectedY) > tolerance {
		t.Errorf("Expected Y position to be %f, got %f", expectedY, y)
	}
}

func TestMovementSystem_Update_NegativeDeltaTime(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with velocity
	entity := createMovementEntity(1, 100, 100, 10, 5)

	entities := []Entity{entity}

	// Run update with negative delta time
	ms.Update(-0.016, entities, eventDispatcher)

	// Verify transform position was updated (negative delta time should still move)
	transformComp := entity.GetTransform()
	transformObj := transformComp.(components.TransformComponent)
	x, y := transformObj.GetX(), transformObj.GetY()

	// Expected position: (100 + 10*(-0.016), 100 + 5*(-0.016)) = (99.84, 99.92)
	expectedX := 100.0 + 10.0*(-0.016)
	expectedY := 100.0 + 5.0*(-0.016)

	if math.Abs(x-expectedX) > tolerance {
		t.Errorf("Expected X position to be %f, got %f", expectedX, x)
	}

	if math.Abs(y-expectedY) > tolerance {
		t.Errorf("Expected Y position to be %f, got %f", expectedY, y)
	}
}

func TestMovementSystem_Update_MultipleUpdates(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with velocity
	entity := createMovementEntity(1, 100, 100, 10, 5)

	entities := []Entity{entity}

	// Run multiple updates
	for i := 0; i < 5; i++ {
		ms.Update(0.016, entities, eventDispatcher)
	}

	// Verify call count was incremented correctly
	if ms.callCount != 5 {
		t.Errorf("Expected callCount to be 5, got %d", ms.callCount)
	}

	// Verify transform position was updated correctly
	transformComp := entity.GetTransform()
	transformObj := transformComp.(components.TransformComponent)
	x, y := transformObj.GetX(), transformObj.GetY()

	// Expected position after 5 updates: (100 + 10*0.016*5, 100 + 5*0.016*5) = (100.8, 100.4)
	expectedX := 100.0 + 10.0*0.016*5
	expectedY := 100.0 + 5.0*0.016*5

	if math.Abs(x-expectedX) > tolerance {
		t.Errorf("Expected X position to be %f, got %f", expectedX, x)
	}

	if math.Abs(y-expectedY) > tolerance {
		t.Errorf("Expected Y position to be %f, got %f", expectedY, y)
	}
}

func TestMovementSystem_Update_PhysicsUpdateCalled(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with velocity
	entity := createMovementEntity(1, 100, 100, 10, 5)

	entities := []Entity{entity}

	// Get initial physics state
	physicsComp := entity.GetPhysics()
	physicsObj := physicsComp.(components.PhysicsComponent)
	initialVX := physicsObj.GetVelocityX()
	initialVY := physicsObj.GetVelocityY()

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify physics was updated (velocity should remain the same since no acceleration)
	// The physics.Update() method should have been called
	finalVX := physicsObj.GetVelocityX()
	finalVY := physicsObj.GetVelocityY()

	if math.Abs(finalVX-initialVX) > tolerance {
		t.Errorf("Expected velocity X to remain %f, got %f", initialVX, finalVX)
	}

	if math.Abs(finalVY-initialVY) > tolerance {
		t.Errorf("Expected velocity Y to remain %f, got %f", initialVY, finalVY)
	}
}

func TestMovementSystem_Update_PhysicsWithAcceleration(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with velocity and apply acceleration
	entity := createMovementEntity(1, 100, 100, 10, 5)
	physicsComp := entity.GetPhysics()
	physicsObj := physicsComp.(components.PhysicsComponent)

	// Apply force to create acceleration
	physicsObj.(*components.Physics).ApplyForce(5, 2)

	entities := []Entity{entity}

	// Run update
	ms.Update(0.016, entities, eventDispatcher)

	// Verify physics was updated with acceleration
	// Initial velocity: (10, 5)
	// Applied force: (5, 2) with mass 1.0
	// Expected final velocity: (10 + 5*0.016, 5 + 2*0.016) = (10.08, 5.032)
	expectedVX := 10.0 + 5.0*0.016
	expectedVY := 5.0 + 2.0*0.016

	finalVX := physicsObj.GetVelocityX()
	finalVY := physicsObj.GetVelocityY()

	if math.Abs(finalVX-expectedVX) > tolerance {
		t.Errorf("Expected velocity X to be %f, got %f", expectedVX, finalVX)
	}

	if math.Abs(finalVY-expectedVY) > tolerance {
		t.Errorf("Expected velocity Y to be %f, got %f", expectedVY, finalVY)
	}
}

func TestMovementSystem_Update_CallCountLogging(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with velocity
	entity := createMovementEntity(1, 100, 100, 10, 5)

	entities := []Entity{entity}

	// Run updates to reach call count that triggers logging (every 60 calls)
	for i := 0; i < 60; i++ {
		ms.Update(0.016, entities, eventDispatcher)
	}

	// Verify call count is 60
	if ms.callCount != 60 {
		t.Errorf("Expected callCount to be 60, got %d", ms.callCount)
	}

	// The system should have logged movement information
	// We can't easily test the actual logging output, but we can verify the system didn't crash
}

func TestMovementSystem_Update_PacketLogging(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with Transform, Physics, and PacketType
	entity := createMovementEntity(1, 100, 100, 10, 5)
	packetType := components.NewPacketType("TestPacket", 1)
	entityObj := entity.(*entities.Entity)
	entityObj.AddComponent(packetType)

	entities := []Entity{entity}

	// Run updates to reach call count that triggers packet logging (every 60 calls)
	for i := 0; i < 60; i++ {
		ms.Update(0.016, entities, eventDispatcher)
	}

	// Verify call count is 60
	if ms.callCount != 60 {
		t.Errorf("Expected callCount to be 60, got %d", ms.callCount)
	}

	// The system should have logged packet information
	// We can't easily test the actual logging output, but we can verify the system didn't crash
}

func TestMovementSystem_Integration(t *testing.T) {
	ms := NewMovementSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create multiple entities with different configurations
	entity1 := createMovementEntity(1, 100, 100, 10, 5)  // Moving entity
	entity2 := createMovementEntity(2, 200, 200, 0, 0)   // Stationary entity
	entity3 := createMovementEntity(3, 300, 300, -5, 10) // Moving entity with packet
	packetType := components.NewPacketType("TestPacket", 1)
	entity3Obj := entity3.(*entities.Entity)
	entity3Obj.AddComponent(packetType)

	entities := []Entity{entity1, entity2, entity3}

	// Run multiple updates
	for i := 0; i < 10; i++ {
		ms.Update(0.016, entities, eventDispatcher)
	}

	// Verify call count
	if ms.callCount != 10 {
		t.Errorf("Expected callCount to be 10, got %d", ms.callCount)
	}

	// Verify all entities were processed correctly
	transform1 := entity1.GetTransform().(components.TransformComponent)
	x1, y1 := transform1.GetX(), transform1.GetY()
	expectedX1 := 100.0 + 10.0*0.016*10
	expectedY1 := 100.0 + 5.0*0.016*10
	if math.Abs(x1-expectedX1) > tolerance || math.Abs(y1-expectedY1) > tolerance {
		t.Errorf("Entity1: Expected position (%f, %f), got (%f, %f)", expectedX1, expectedY1, x1, y1)
	}

	transform2 := entity2.GetTransform().(components.TransformComponent)
	x2, y2 := transform2.GetX(), transform2.GetY()
	if math.Abs(x2-200) > tolerance || math.Abs(y2-200) > tolerance {
		t.Errorf("Entity2: Expected position (200, 200), got (%f, %f)", x2, y2)
	}

	transform3 := entity3.GetTransform().(components.TransformComponent)
	x3, y3 := transform3.GetX(), transform3.GetY()
	expectedX3 := 300.0 + (-5.0)*0.016*10
	expectedY3 := 300.0 + 10.0*0.016*10
	if math.Abs(x3-expectedX3) > tolerance || math.Abs(y3-expectedY3) > tolerance {
		t.Errorf("Entity3: Expected position (%f, %f), got (%f, %f)", expectedX3, expectedY3, x3, y3)
	}
}

// Helper function to create test entities

func createMovementEntity(id uint64, x, y, vx, vy float64) Entity {
	entity := entities.NewEntity(id)
	transform := components.NewTransform(x, y)
	physics := components.NewPhysics()
	physics.SetVelocity(vx, vy)
	entity.AddComponent(transform)
	entity.AddComponent(physics)
	return entity
}
