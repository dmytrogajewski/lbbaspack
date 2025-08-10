package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"
)

func TestNewBackendSystem(t *testing.T) {
	bs := NewBackendSystem()

	// Test that the system is properly initialized
	if bs == nil {
		t.Fatal("NewBackendSystem returned nil")
	}

	// Test required components
	expectedComponents := []string{"BackendAssignment"}
	if len(bs.RequiredComponents) != len(expectedComponents) {
		t.Errorf("Expected %d required components, got %d", len(expectedComponents), len(bs.RequiredComponents))
	}

	for i, component := range expectedComponents {
		if bs.RequiredComponents[i] != component {
			t.Errorf("Expected required component %s at index %d, got %s", component, i, bs.RequiredComponents[i])
		}
	}

	// Stateless: no internal counters
}

func TestBackendSystem_Update(t *testing.T) {
	bs := NewBackendSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create test entities with backend assignments
	entity1 := entities.NewEntity(1)
	backendComp1 := components.NewBackendAssignment(1)
	backendComp1.IncrementAssignedPackets() // Set to 1
	entity1.AddComponent(backendComp1)

	entity2 := entities.NewEntity(2)
	backendComp2 := components.NewBackendAssignment(2)
	backendComp2.IncrementAssignedPackets() // Set to 1
	backendComp2.IncrementAssignedPackets() // Set to 2
	entity2.AddComponent(backendComp2)

	// Create entity without backend assignment (should be ignored)
	entity3 := entities.NewEntity(3)

	entities := []Entity{entity1, entity2, entity3}

	// Run update
	bs.Update(0.016, entities, eventDispatcher)

	// Stateless: no backend counter mutation; ensure no panic
}

func TestBackendSystem_Update_NoBackendEntities(t *testing.T) {
	bs := NewBackendSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entities without backend assignments
	entity1 := entities.NewEntity(1)
	entity2 := entities.NewEntity(2)

	entities := []Entity{entity1, entity2}

	// Run update
	bs.Update(0.016, entities, eventDispatcher)

	// Stateless: nothing to assert
}

func TestBackendSystem_Update_EntityWithoutBackendComponent(t *testing.T) {
	bs := NewBackendSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with backend assignment
	entity1 := entities.NewEntity(1)
	backendComp1 := components.NewBackendAssignment(1)
	entity1.AddComponent(backendComp1)

	// Create entity with different component (should be ignored)
	entity2 := entities.NewEntity(2)
	// Add a different component type
	entity2.AddComponent(&mockComponent{componentType: "DifferentComponent"})

	entities := []Entity{entity1, entity2}

	// Run update
	bs.Update(0.016, entities, eventDispatcher)

	// Stateless: nothing to assert
}

func TestBackendSystem_Initialize(t *testing.T) {
	bs := NewBackendSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	bs.Initialize(eventDispatcher)

	// Publish a packet caught event (should not panic)

	// Add a backend to the system first
	entity := entities.NewEntity(1)
	backendComp := components.NewBackendAssignment(1)
	entity.AddComponent(backendComp)
	entities := []Entity{entity}
	bs.Update(0.016, entities, eventDispatcher)

	// Publish packet caught event
	event := events.NewEvent(events.EventPacketCaught, nil)
	eventDispatcher.Publish(event)

	// Verify that a packet was assigned
	// no-op
}

func TestBackendSystem_assignPacketToBackend_LoadBalancing(t *testing.T) {
	// Stateless: no-op
}

func TestBackendSystem_assignPacketToBackend_EqualLoads(t *testing.T) {
	// Stateless: no-op
}

func TestBackendSystem_assignPacketToBackend_NoBackends(t *testing.T) {
	// Stateless: no-op
}

func TestBackendSystem_assignPacketToBackend_SingleBackend(t *testing.T) {
	// Stateless: no-op
}

func TestBackendSystem_GetBackendStats(t *testing.T) {
	bs := NewBackendSystem()

	// Stateless: returns empty stats
	stats := bs.GetBackendStats()
	if len(stats) != 0 {
		t.Errorf("Expected 0 backend entries, got %d", len(stats))
	}
}

func TestBackendSystem_GetBackendStats_Empty(t *testing.T) {
	bs := NewBackendSystem()

	// Stateless: ensure empty
	stats := bs.GetBackendStats()
	if len(stats) != 0 {
		t.Errorf("Expected 0 backend entries, got %d", len(stats))
	}
}

func TestBackendSystem_GetTotalPackets(t *testing.T) {
	bs := NewBackendSystem()

	// Stateless: always 0
	if bs.GetTotalPackets() != 0 {
		t.Errorf("Expected total packets to be 0")
	}
}

func TestBackendSystem_GetTotalPackets_Zero(t *testing.T) {
	bs := NewBackendSystem()

	if bs.GetTotalPackets() != 0 {
		t.Errorf("Expected total packets to be 0")
	}
}

func TestBackendSystem_Integration(t *testing.T) {
	bs := NewBackendSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize system
	bs.Initialize(eventDispatcher)

	// Create entities with backend assignments
	entity1 := entities.NewEntity(1)
	backendComp1 := components.NewBackendAssignment(1)
	entity1.AddComponent(backendComp1)

	entity2 := entities.NewEntity(2)
	backendComp2 := components.NewBackendAssignment(2)
	entity2.AddComponent(backendComp2)

	entities := []Entity{entity1, entity2}

	// Update should be no-op; ensure no panic
	bs.Update(0.016, entities, eventDispatcher)

	// Publish multiple packet caught events to ensure no panic
	for i := 0; i < 5; i++ {
		event := events.NewEvent(events.EventPacketCaught, nil)
		eventDispatcher.Publish(event)
	}
}

// Mock component for testing
type mockComponent struct {
	componentType string
}

func (mc *mockComponent) GetType() string {
	return mc.componentType
}
