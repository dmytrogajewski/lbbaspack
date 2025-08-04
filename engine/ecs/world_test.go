package ecs

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"lbbaspack/engine/systems"
	"testing"
)

// Mock system for testing
type mockSystem struct {
	updateCalled bool
	deltaTime    float64
	entities     []systems.Entity
	dispatcher   *events.EventDispatcher
}

func (ms *mockSystem) Update(deltaTime float64, entities []systems.Entity, eventDispatcher *events.EventDispatcher) {
	ms.updateCalled = true
	ms.deltaTime = deltaTime
	ms.entities = entities
	ms.dispatcher = eventDispatcher
}

func (ms *mockSystem) GetRequiredComponents() []string {
	return []string{} // No required components for mock system
}

// TestNewWorld tests the NewWorld constructor
func TestNewWorld(t *testing.T) {
	world := NewWorld()

	if world == nil {
		t.Fatal("Expected world to be created")
	}

	if world.Entities == nil {
		t.Error("Expected entities slice to be initialized")
	}

	if len(world.Entities) != 0 {
		t.Error("Expected entities slice to be empty initially")
	}

	if world.Systems == nil {
		t.Error("Expected systems slice to be initialized")
	}

	if len(world.Systems) != 0 {
		t.Error("Expected systems slice to be empty initially")
	}

	if world.EventDispatcher == nil {
		t.Error("Expected event dispatcher to be initialized")
	}

	if world.nextEntityID != 1 {
		t.Errorf("Expected nextEntityID to be 1, got %d", world.nextEntityID)
	}
}

// TestWorld_AddEntity tests the AddEntity method
func TestWorld_AddEntity(t *testing.T) {
	world := NewWorld()

	// Create test entities
	entity1 := entities.NewEntity(1)
	entity2 := entities.NewEntity(2)

	// Add entities
	world.AddEntity(entity1)
	world.AddEntity(entity2)

	if len(world.Entities) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(world.Entities))
	}

	if world.Entities[0] != entity1 {
		t.Error("Expected first entity to be entity1")
	}

	if world.Entities[1] != entity2 {
		t.Error("Expected second entity to be entity2")
	}

	// Test adding nil entity
	world.AddEntity(nil)
	if len(world.Entities) != 3 {
		t.Errorf("Expected 3 entities after adding nil, got %d", len(world.Entities))
	}

	if world.Entities[2] != nil {
		t.Error("Expected third entity to be nil")
	}
}

// TestWorld_NewEntity tests the NewEntity method
func TestWorld_NewEntity(t *testing.T) {
	world := NewWorld()

	// Create first entity
	entity1 := world.NewEntity()

	if entity1 == nil {
		t.Fatal("Expected entity to be created")
	}

	if entity1.ID != 1 {
		t.Errorf("Expected entity ID 1, got %d", entity1.ID)
	}

	if len(world.Entities) != 1 {
		t.Errorf("Expected 1 entity in world, got %d", len(world.Entities))
	}

	if world.Entities[0] != entity1 {
		t.Error("Expected entity to be added to world")
	}

	if world.nextEntityID != 2 {
		t.Errorf("Expected nextEntityID to be 2, got %d", world.nextEntityID)
	}

	// Create second entity
	entity2 := world.NewEntity()

	if entity2 == nil {
		t.Fatal("Expected second entity to be created")
	}

	if entity2.ID != 2 {
		t.Errorf("Expected entity ID 2, got %d", entity2.ID)
	}

	if len(world.Entities) != 2 {
		t.Errorf("Expected 2 entities in world, got %d", len(world.Entities))
	}

	if world.Entities[1] != entity2 {
		t.Error("Expected second entity to be added to world")
	}

	if world.nextEntityID != 3 {
		t.Errorf("Expected nextEntityID to be 3, got %d", world.nextEntityID)
	}
}

// TestWorld_AddSystem tests the AddSystem method
func TestWorld_AddSystem(t *testing.T) {
	world := NewWorld()

	// Create mock systems
	system1 := &mockSystem{}
	system2 := &mockSystem{}

	// Add systems
	world.AddSystem(system1)
	world.AddSystem(system2)

	if len(world.Systems) != 2 {
		t.Errorf("Expected 2 systems, got %d", len(world.Systems))
	}

	if world.Systems[0] != system1 {
		t.Error("Expected first system to be system1")
	}

	if world.Systems[1] != system2 {
		t.Error("Expected second system to be system2")
	}

	// Test adding nil system
	world.AddSystem(nil)
	if len(world.Systems) != 3 {
		t.Errorf("Expected 3 systems after adding nil, got %d", len(world.Systems))
	}

	if world.Systems[2] != nil {
		t.Error("Expected third system to be nil")
	}
}

// TestWorld_Update tests the Update method
func TestWorld_Update(t *testing.T) {
	world := NewWorld()

	// Create mock system
	mockSys := &mockSystem{}
	world.AddSystem(mockSys)

	// Create test entities
	entity1 := world.NewEntity()
	entity2 := world.NewEntity()

	// Add components to entities
	transform1 := components.NewTransform(100, 200)
	sprite1 := components.NewSprite(32, 32, components.RandomPacketColor())
	entity1.AddComponent(transform1)
	entity1.AddComponent(sprite1)

	transform2 := components.NewTransform(300, 400)
	collider2 := components.NewCollider(64, 64, "test")
	entity2.AddComponent(transform2)
	entity2.AddComponent(collider2)

	// Update world
	deltaTime := 0.016 // 60 FPS
	world.Update(deltaTime)

	// Verify system was called
	if !mockSys.updateCalled {
		t.Error("Expected system Update to be called")
	}

	if mockSys.deltaTime != deltaTime {
		t.Errorf("Expected deltaTime %f, got %f", deltaTime, mockSys.deltaTime)
	}

	if len(mockSys.entities) != 2 {
		t.Errorf("Expected 2 entities passed to system, got %d", len(mockSys.entities))
	}

	if mockSys.dispatcher != world.EventDispatcher {
		t.Error("Expected event dispatcher to be passed to system")
	}

	// Verify entities are correctly converted to interface
	if mockSys.entities[0] != entity1 {
		t.Error("Expected first entity in system to be entity1")
	}

	if mockSys.entities[1] != entity2 {
		t.Error("Expected second entity in system to be entity2")
	}
}

// TestWorld_Update_EmptyWorld tests Update with no entities or systems
func TestWorld_Update_EmptyWorld(t *testing.T) {
	world := NewWorld()

	// Should not panic
	world.Update(0.016)
}

// TestWorld_Update_NoSystems tests Update with entities but no systems
func TestWorld_Update_NoSystems(t *testing.T) {
	world := NewWorld()

	// Add entities but no systems
	entity1 := world.NewEntity()
	entity2 := world.NewEntity()

	// Should not panic
	world.Update(0.016)

	// Verify entities still exist
	if len(world.Entities) != 2 {
		t.Errorf("Expected 2 entities, got %d", len(world.Entities))
	}

	if world.Entities[0] != entity1 {
		t.Error("Expected entity1 to still exist")
	}

	if world.Entities[1] != entity2 {
		t.Error("Expected entity2 to still exist")
	}
}

// TestWorld_Update_NoEntities tests Update with systems but no entities
func TestWorld_Update_NoEntities(t *testing.T) {
	world := NewWorld()

	// Add system but no entities
	mockSys := &mockSystem{}
	world.AddSystem(mockSys)

	// Should not panic
	world.Update(0.016)

	// Verify system was called with empty entities slice
	if !mockSys.updateCalled {
		t.Error("Expected system Update to be called")
	}

	if len(mockSys.entities) != 0 {
		t.Errorf("Expected 0 entities passed to system, got %d", len(mockSys.entities))
	}
}

// TestWorld_Update_MultipleSystems tests Update with multiple systems
func TestWorld_Update_MultipleSystems(t *testing.T) {
	world := NewWorld()

	// Create multiple mock systems
	system1 := &mockSystem{}
	system2 := &mockSystem{}
	system3 := &mockSystem{}

	world.AddSystem(system1)
	world.AddSystem(system2)
	world.AddSystem(system3)

	// Add entities
	world.NewEntity()
	world.NewEntity()

	// Update world
	deltaTime := 0.033
	world.Update(deltaTime)

	// Verify all systems were called
	if !system1.updateCalled {
		t.Error("Expected system1 Update to be called")
	}
	if !system2.updateCalled {
		t.Error("Expected system2 Update to be called")
	}
	if !system3.updateCalled {
		t.Error("Expected system3 Update to be called")
	}

	// Verify all systems received correct parameters
	systems := []*mockSystem{system1, system2, system3}
	for i, system := range systems {
		if system.deltaTime != deltaTime {
			t.Errorf("Expected system%d deltaTime %f, got %f", i+1, deltaTime, system.deltaTime)
		}

		if len(system.entities) != 2 {
			t.Errorf("Expected system%d to receive 2 entities, got %d", i+1, len(system.entities))
		}

		if system.dispatcher != world.EventDispatcher {
			t.Errorf("Expected system%d to receive correct event dispatcher", i+1)
		}
	}
}

// TestWorld_EntityIDSequence tests that entity IDs are assigned sequentially
func TestWorld_EntityIDSequence(t *testing.T) {
	world := NewWorld()

	// Create multiple entities
	entities := make([]*entities.Entity, 10)
	for i := 0; i < 10; i++ {
		entities[i] = world.NewEntity()
	}

	// Verify IDs are sequential starting from 1
	for i, entity := range entities {
		expectedID := uint64(i + 1)
		if entity.ID != expectedID {
			t.Errorf("Expected entity %d to have ID %d, got %d", i, expectedID, entity.ID)
		}
	}

	// Verify nextEntityID is correct
	if world.nextEntityID != 11 {
		t.Errorf("Expected nextEntityID to be 11, got %d", world.nextEntityID)
	}
}

// TestWorld_EventDispatcher tests that the event dispatcher is properly initialized
func TestWorld_EventDispatcher(t *testing.T) {
	world := NewWorld()

	if world.EventDispatcher == nil {
		t.Fatal("Expected event dispatcher to be initialized")
	}

	// Test that we can subscribe to events
	eventReceived := false
	world.EventDispatcher.Subscribe(events.EventGameStart, func(event *events.Event) {
		eventReceived = true
	})

	// Publish an event
	gameStartEvent := events.NewEvent(events.EventGameStart, nil)
	world.EventDispatcher.Publish(gameStartEvent)

	if !eventReceived {
		t.Error("Expected event to be received")
	}
}

// TestWorld_Integration tests integration scenarios
func TestWorld_Integration(t *testing.T) {
	world := NewWorld()

	// Create multiple systems
	system1 := &mockSystem{}
	system2 := &mockSystem{}
	world.AddSystem(system1)
	world.AddSystem(system2)

	// Create multiple entities with components
	entity1 := world.NewEntity()
	transform1 := components.NewTransform(100, 200)
	sprite1 := components.NewSprite(32, 32, components.RandomPacketColor())
	entity1.AddComponent(transform1)
	entity1.AddComponent(sprite1)

	entity2 := world.NewEntity()
	transform2 := components.NewTransform(300, 400)
	collider2 := components.NewCollider(64, 64, "player")
	physics2 := components.NewPhysics()
	entity2.AddComponent(transform2)
	entity2.AddComponent(collider2)
	entity2.AddComponent(physics2)

	entity3 := world.NewEntity()
	packetType := components.NewPacketType("HTTP", 1)
	state := components.NewState(components.StatePlaying)
	entity3.AddComponent(packetType)
	entity3.AddComponent(state)

	// Verify entities were added to world
	if len(world.Entities) != 3 {
		t.Errorf("Expected 3 entities in world, got %d", len(world.Entities))
	}

	// Verify entity IDs
	if entity1.ID != 1 {
		t.Errorf("Expected entity1 ID 1, got %d", entity1.ID)
	}
	if entity2.ID != 2 {
		t.Errorf("Expected entity2 ID 2, got %d", entity2.ID)
	}
	if entity3.ID != 3 {
		t.Errorf("Expected entity3 ID 3, got %d", entity3.ID)
	}

	// Update world multiple times
	for i := 0; i < 5; i++ {
		world.Update(0.016)
	}

	// Verify systems were called
	if !system1.updateCalled {
		t.Error("Expected system1 to be called")
	}
	if !system2.updateCalled {
		t.Error("Expected system2 to be called")
	}

	// Verify systems received correct entity count
	if len(system1.entities) != 3 {
		t.Errorf("Expected system1 to receive 3 entities, got %d", len(system1.entities))
	}
	if len(system2.entities) != 3 {
		t.Errorf("Expected system2 to receive 3 entities, got %d", len(system2.entities))
	}

	// Verify nextEntityID is correct
	if world.nextEntityID != 4 {
		t.Errorf("Expected nextEntityID to be 4, got %d", world.nextEntityID)
	}
}

// TestWorld_EdgeCases tests edge cases and error conditions
func TestWorld_EdgeCases(t *testing.T) {
	t.Run("Nil Entity Addition", func(t *testing.T) {
		world := NewWorld()
		world.AddEntity(nil)

		if len(world.Entities) != 1 {
			t.Errorf("Expected 1 entity after adding nil, got %d", len(world.Entities))
		}

		if world.Entities[0] != nil {
			t.Error("Expected nil entity to be stored")
		}
	})

	t.Run("Nil System Addition", func(t *testing.T) {
		world := NewWorld()
		world.AddSystem(nil)

		if len(world.Systems) != 1 {
			t.Errorf("Expected 1 system after adding nil, got %d", len(world.Systems))
		}

		if world.Systems[0] != nil {
			t.Error("Expected nil system to be stored")
		}
	})

	t.Run("Large Number of Entities", func(t *testing.T) {
		world := NewWorld()

		// Create many entities
		for i := 0; i < 1000; i++ {
			entity := world.NewEntity()
			if entity.ID != uint64(i+1) {
				t.Errorf("Expected entity ID %d, got %d", i+1, entity.ID)
			}
		}

		if len(world.Entities) != 1000 {
			t.Errorf("Expected 1000 entities, got %d", len(world.Entities))
		}

		if world.nextEntityID != 1001 {
			t.Errorf("Expected nextEntityID to be 1001, got %d", world.nextEntityID)
		}
	})

	t.Run("Large Number of Systems", func(t *testing.T) {
		world := NewWorld()

		// Create many systems
		for i := 0; i < 100; i++ {
			system := &mockSystem{}
			world.AddSystem(system)
		}

		if len(world.Systems) != 100 {
			t.Errorf("Expected 100 systems, got %d", len(world.Systems))
		}

		// Add an entity and update
		world.NewEntity()
		world.Update(0.016)

		// Verify all systems were called
		for i, system := range world.Systems {
			mockSys := system.(*mockSystem)
			if !mockSys.updateCalled {
				t.Errorf("Expected system %d to be called", i)
			}
		}
	})

	t.Run("Zero Delta Time", func(t *testing.T) {
		world := NewWorld()
		mockSys := &mockSystem{}
		world.AddSystem(mockSys)
		world.NewEntity()

		world.Update(0.0)

		if !mockSys.updateCalled {
			t.Error("Expected system to be called with zero delta time")
		}

		if mockSys.deltaTime != 0.0 {
			t.Errorf("Expected deltaTime 0.0, got %f", mockSys.deltaTime)
		}
	})

	t.Run("Negative Delta Time", func(t *testing.T) {
		world := NewWorld()
		mockSys := &mockSystem{}
		world.AddSystem(mockSys)
		world.NewEntity()

		world.Update(-0.016)

		if !mockSys.updateCalled {
			t.Error("Expected system to be called with negative delta time")
		}

		if mockSys.deltaTime != -0.016 {
			t.Errorf("Expected deltaTime -0.016, got %f", mockSys.deltaTime)
		}
	})
}

func TestWorld_RemoveEntity(t *testing.T) {
	world := NewWorld()

	// Add some entities
	entity1 := world.NewEntity()
	entity2 := world.NewEntity()
	entity3 := world.NewEntity()

	initialCount := len(world.Entities)
	if initialCount != 3 {
		t.Errorf("Expected 3 entities, got %d", initialCount)
	}

	// Remove entity2
	world.RemoveEntity(entity2)

	if len(world.Entities) != 2 {
		t.Errorf("Expected 2 entities after removal, got %d", len(world.Entities))
	}

	// Check that entity1 and entity3 are still there
	found1 := false
	found3 := false
	for _, entity := range world.Entities {
		if entity.ID == entity1.ID {
			found1 = true
		}
		if entity.ID == entity3.ID {
			found3 = true
		}
	}

	if !found1 {
		t.Error("Entity1 should still be in the world")
	}
	if !found3 {
		t.Error("Entity3 should still be in the world")
	}

	// Try to remove non-existent entity (should not panic)
	nonExistentEntity := entities.NewEntity(999)
	world.RemoveEntity(nonExistentEntity)

	if len(world.Entities) != 2 {
		t.Errorf("Expected 2 entities after removing non-existent entity, got %d", len(world.Entities))
	}
}

func TestWorld_RemoveInactiveEntities(t *testing.T) {
	world := NewWorld()

	// Add some entities
	entity1 := world.NewEntity() // Active by default
	entity2 := world.NewEntity()
	entity3 := world.NewEntity()

	// Make entity2 inactive
	entity2.SetActive(false)

	// Verify initial state
	if !entity1.IsActive() || !entity3.IsActive() {
		t.Error("Entity1 and Entity3 should be active initially")
	}

	initialCount := len(world.Entities)
	if initialCount != 3 {
		t.Errorf("Expected 3 entities, got %d", initialCount)
	}

	// Remove inactive entities
	world.RemoveInactiveEntities()

	if len(world.Entities) != 2 {
		t.Errorf("Expected 2 entities after removing inactive ones, got %d", len(world.Entities))
	}

	// Check that only active entities remain
	for _, entity := range world.Entities {
		if !entity.IsActive() {
			t.Errorf("Entity %d should be active", entity.ID)
		}
	}
}

func TestWorld_ClearAllEntities(t *testing.T) {
	world := NewWorld()

	// Add some entities
	world.NewEntity()
	world.NewEntity()
	world.NewEntity()

	if len(world.Entities) != 3 {
		t.Errorf("Expected 3 entities, got %d", len(world.Entities))
	}

	// Clear all entities
	world.ClearAllEntities()

	if len(world.Entities) != 0 {
		t.Errorf("Expected 0 entities after clearing, got %d", len(world.Entities))
	}
}

func TestWorld_EntityCleanup_Integration(t *testing.T) {
	world := NewWorld()

	// Create a load balancer (should persist)
	loadBalancer := world.NewEntity()
	loadBalancer.AddComponent(components.NewTransform(350, 480))
	loadBalancer.AddComponent(components.NewSprite(100, 20, components.RandomPacketColor()))
	loadBalancer.AddComponent(components.NewCollider(100, 20, "loadbalancer"))

	// Create some packets (should be removed when inactive)
	packet1 := world.NewEntity()
	packet1.AddComponent(components.NewTransform(100, 100))
	packet1.AddComponent(components.NewPacketType("HTTP", 10))

	packet2 := world.NewEntity()
	packet2.AddComponent(components.NewTransform(200, 200))
	packet2.AddComponent(components.NewPacketType("HTTPS", 15))

	// Create a power-up (should be removed when inactive)
	powerUp := world.NewEntity()
	powerUp.AddComponent(components.NewTransform(300, 300))
	powerUp.AddComponent(components.NewPowerUpType("SpeedBoost", 15.0))

	initialCount := len(world.Entities)
	if initialCount != 4 {
		t.Errorf("Expected 4 entities initially, got %d", initialCount)
	}

	// Deactivate packets and power-up
	packet1.SetActive(false)
	packet2.SetActive(false)
	powerUp.SetActive(false)

	// Remove inactive entities
	world.RemoveInactiveEntities()

	if len(world.Entities) != 1 {
		t.Errorf("Expected 1 entity after cleanup (load balancer), got %d", len(world.Entities))
	}

	// Check that only the load balancer remains
	remainingEntity := world.Entities[0]
	if remainingEntity.ID != loadBalancer.ID {
		t.Errorf("Expected load balancer to remain, got entity %d", remainingEntity.ID)
	}

	if !remainingEntity.IsActive() {
		t.Error("Load balancer should still be active")
	}
}

// Benchmark tests for performance
func BenchmarkNewWorld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewWorld()
	}
}

func BenchmarkWorld_NewEntity(b *testing.B) {
	world := NewWorld()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.NewEntity()
	}
}

func BenchmarkWorld_AddEntity(b *testing.B) {
	world := NewWorld()
	entity := entities.NewEntity(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.AddEntity(entity)
	}
}

func BenchmarkWorld_AddSystem(b *testing.B) {
	world := NewWorld()
	system := &mockSystem{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.AddSystem(system)
	}
}

func BenchmarkWorld_Update(b *testing.B) {
	world := NewWorld()

	// Add some systems
	for i := 0; i < 5; i++ {
		world.AddSystem(&mockSystem{})
	}

	// Add some entities
	for i := 0; i < 10; i++ {
		entity := world.NewEntity()
		entity.AddComponent(components.NewTransform(100, 200))
		entity.AddComponent(components.NewSprite(32, 32, components.RandomPacketColor()))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update(0.016)
	}
}
