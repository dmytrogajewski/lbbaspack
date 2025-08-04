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

	// Test that backendCounters is initialized
	if bs.backendCounters == nil {
		t.Error("backendCounters map was not initialized")
	}

	// Test initial values
	if bs.totalPackets != 0 {
		t.Errorf("Expected initial totalPackets to be 0, got %d", bs.totalPackets)
	}
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

	// Verify backend counters were updated
	stats := bs.GetBackendStats()
	if len(stats) != 2 {
		t.Errorf("Expected 2 backend entries, got %d", len(stats))
	}

	if stats[1] != 1 {
		t.Errorf("Expected backend 1 to have 1 packet, got %d", stats[1])
	}

	if stats[2] != 2 {
		t.Errorf("Expected backend 2 to have 2 packets, got %d", stats[2])
	}
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

	// Verify no backend counters were created
	stats := bs.GetBackendStats()
	if len(stats) != 0 {
		t.Errorf("Expected 0 backend entries, got %d", len(stats))
	}
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

	// Verify only the backend entity was processed
	stats := bs.GetBackendStats()
	if len(stats) != 1 {
		t.Errorf("Expected 1 backend entry, got %d", len(stats))
	}

	if stats[1] != 0 {
		t.Errorf("Expected backend 1 to have 0 packets, got %d", stats[1])
	}
}

func TestBackendSystem_Initialize(t *testing.T) {
	bs := NewBackendSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	bs.Initialize(eventDispatcher)

	// Verify event subscription by publishing a packet caught event
	// and checking if the backend counter is updated
	initialTotal := bs.GetTotalPackets()

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
	finalStats := bs.GetBackendStats()
	finalTotal := bs.GetTotalPackets()

	if finalTotal != initialTotal+1 {
		t.Errorf("Expected total packets to increase by 1, got %d -> %d", initialTotal, finalTotal)
	}

	if len(finalStats) == 0 {
		t.Error("Expected backend stats to be populated after packet assignment")
	}
}

func TestBackendSystem_assignPacketToBackend_LoadBalancing(t *testing.T) {
	bs := NewBackendSystem()

	// Add multiple backends with different packet counts
	bs.backendCounters[1] = 5 // Most loaded
	bs.backendCounters[2] = 2 // Least loaded
	bs.backendCounters[3] = 4 // Medium loaded

	// Assign a packet
	bs.assignPacketToBackend()

	// Verify packet was assigned to backend with least packets (backend 2)
	if bs.backendCounters[2] != 3 {
		t.Errorf("Expected backend 2 to have 3 packets after assignment, got %d", bs.backendCounters[2])
	}

	// Verify other backends unchanged
	if bs.backendCounters[1] != 5 {
		t.Errorf("Expected backend 1 to still have 5 packets, got %d", bs.backendCounters[1])
	}

	if bs.backendCounters[3] != 4 {
		t.Errorf("Expected backend 3 to still have 4 packets, got %d", bs.backendCounters[3])
	}

	// Verify total packets increased
	if bs.totalPackets != 1 {
		t.Errorf("Expected total packets to be 1, got %d", bs.totalPackets)
	}
}

func TestBackendSystem_assignPacketToBackend_EqualLoads(t *testing.T) {
	bs := NewBackendSystem()

	// Add multiple backends with equal packet counts
	bs.backendCounters[1] = 3
	bs.backendCounters[2] = 3
	bs.backendCounters[3] = 3

	// Assign a packet
	bs.assignPacketToBackend()

	// Verify one of the backends got the packet
	totalPackets := 0
	for _, count := range bs.backendCounters {
		totalPackets += count
	}

	expectedTotal := 9 + 1 // 3 backends * 3 packets + 1 new packet
	if totalPackets != expectedTotal {
		t.Errorf("Expected total packets across all backends to be %d, got %d", expectedTotal, totalPackets)
	}

	// Verify total packets increased
	if bs.totalPackets != 1 {
		t.Errorf("Expected total packets to be 1, got %d", bs.totalPackets)
	}
}

func TestBackendSystem_assignPacketToBackend_NoBackends(t *testing.T) {
	bs := NewBackendSystem()

	// Don't add any backends
	initialTotal := bs.totalPackets

	// Assign a packet
	bs.assignPacketToBackend()

	// Verify nothing changed
	if bs.totalPackets != initialTotal {
		t.Errorf("Expected total packets to remain unchanged, got %d", bs.totalPackets)
	}

	stats := bs.GetBackendStats()
	if len(stats) != 0 {
		t.Errorf("Expected no backend stats, got %d entries", len(stats))
	}
}

func TestBackendSystem_assignPacketToBackend_SingleBackend(t *testing.T) {
	bs := NewBackendSystem()

	// Add single backend
	bs.backendCounters[1] = 5

	// Assign a packet
	bs.assignPacketToBackend()

	// Verify packet was assigned to the only backend
	if bs.backendCounters[1] != 6 {
		t.Errorf("Expected backend 1 to have 6 packets after assignment, got %d", bs.backendCounters[1])
	}

	// Verify total packets increased
	if bs.totalPackets != 1 {
		t.Errorf("Expected total packets to be 1, got %d", bs.totalPackets)
	}
}

func TestBackendSystem_GetBackendStats(t *testing.T) {
	bs := NewBackendSystem()

	// Add some backend data
	bs.backendCounters[1] = 5
	bs.backendCounters[2] = 3
	bs.backendCounters[3] = 7

	// Get stats
	stats := bs.GetBackendStats()

	// Verify stats are returned correctly
	if len(stats) != 3 {
		t.Errorf("Expected 3 backend entries, got %d", len(stats))
	}

	if stats[1] != 5 {
		t.Errorf("Expected backend 1 to have 5 packets, got %d", stats[1])
	}

	if stats[2] != 3 {
		t.Errorf("Expected backend 2 to have 3 packets, got %d", stats[2])
	}

	if stats[3] != 7 {
		t.Errorf("Expected backend 3 to have 7 packets, got %d", stats[3])
	}
}

func TestBackendSystem_GetBackendStats_Empty(t *testing.T) {
	bs := NewBackendSystem()

	// Get stats without any backends
	stats := bs.GetBackendStats()

	// Verify empty stats
	if len(stats) != 0 {
		t.Errorf("Expected 0 backend entries, got %d", len(stats))
	}
}

func TestBackendSystem_GetTotalPackets(t *testing.T) {
	bs := NewBackendSystem()

	// Set total packets
	bs.totalPackets = 42

	// Get total
	total := bs.GetTotalPackets()

	// Verify total is returned correctly
	if total != 42 {
		t.Errorf("Expected total packets to be 42, got %d", total)
	}
}

func TestBackendSystem_GetTotalPackets_Zero(t *testing.T) {
	bs := NewBackendSystem()

	// Get total without any packets
	total := bs.GetTotalPackets()

	// Verify zero total
	if total != 0 {
		t.Errorf("Expected total packets to be 0, got %d", total)
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

	// Update system to register backends
	bs.Update(0.016, entities, eventDispatcher)

	// Verify initial state
	initialStats := bs.GetBackendStats()
	if len(initialStats) != 2 {
		t.Errorf("Expected 2 backend entries after update, got %d", len(initialStats))
	}

	// Publish multiple packet caught events
	for i := 0; i < 5; i++ {
		event := events.NewEvent(events.EventPacketCaught, nil)
		eventDispatcher.Publish(event)
	}

	// Verify load balancing worked
	finalStats := bs.GetBackendStats()
	finalTotal := bs.GetTotalPackets()

	if finalTotal != 5 {
		t.Errorf("Expected total packets to be 5, got %d", finalTotal)
	}

	// Verify packets were distributed (should be roughly balanced)
	totalBackend1 := finalStats[1]
	totalBackend2 := finalStats[2]

	if totalBackend1 < 2 || totalBackend1 > 3 {
		t.Errorf("Expected backend 1 to have 2-3 packets, got %d", totalBackend1)
	}

	if totalBackend2 < 2 || totalBackend2 > 3 {
		t.Errorf("Expected backend 2 to have 2-3 packets, got %d", totalBackend2)
	}

	// Verify total matches sum of individual backends
	if totalBackend1+totalBackend2 != finalTotal {
		t.Errorf("Expected sum of backend packets (%d) to equal total (%d)", totalBackend1+totalBackend2, finalTotal)
	}
}

// Mock component for testing
type mockComponent struct {
	componentType string
}

func (mc *mockComponent) GetType() string {
	return mc.componentType
}
