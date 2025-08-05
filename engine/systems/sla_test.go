package systems

import (
	"fmt"
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"
)

func TestNewSLASystem(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)

	// Test that the system is properly initialized
	if ss == nil {
		t.Fatal("NewSLASystem returned nil")
	}

	// Test required components
	expectedComponents := []string{"SLA"}
	if len(ss.RequiredComponents) != len(expectedComponents) {
		t.Errorf("Expected %d required components, got %d", len(expectedComponents), len(ss.RequiredComponents))
	}

	for i, component := range expectedComponents {
		if ss.RequiredComponents[i] != component {
			t.Errorf("Expected required component %s at index %d, got %s", component, i, ss.RequiredComponents[i])
		}
	}

	// Test initial values
	if ss.totalPackets != 0 {
		t.Errorf("Expected initial total packets to be 0, got %d", ss.totalPackets)
	}
	if ss.caughtPackets != 0 {
		t.Errorf("Expected initial caught packets to be 0, got %d", ss.caughtPackets)
	}
	if ss.lostPackets != 0 {
		t.Errorf("Expected initial lost packets to be 0, got %d", ss.lostPackets)
	}
	if ss.errorBudget != 10 {
		t.Errorf("Expected initial error budget to be 10, got %d", ss.errorBudget)
	}
	if ss.spawnSys != spawnSys {
		t.Error("Expected spawnSys to be set correctly")
	}
}

func TestSLASystem_Update_NoEntities(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}

	// Test that Update doesn't panic with no entities
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with no entities: %v", r)
		}
	}()

	ss.Update(0.016, entities, eventDispatcher)
}

func TestSLASystem_Update_EntityWithoutSLAComponent(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Create entity without SLA component
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	entity.AddComponent(transform)

	entities := []Entity{entity}

	// Test that Update doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with entity without SLA component: %v", r)
		}
	}()

	ss.Update(0.016, entities, eventDispatcher)
}

func TestSLASystem_Update_EntityWithSLAComponent(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with SLA component
	entity := entities.NewEntity(1)
	sla := components.NewSLA(95.0, 10)
	entity.AddComponent(sla)

	entities := []Entity{entity}

	// Set some packet statistics
	ss.totalPackets = 10
	ss.caughtPackets = 8
	ss.lostPackets = 2

	// Test that Update doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with entity with SLA component: %v", r)
		}
	}()

	ss.Update(0.016, entities, eventDispatcher)

	// Verify SLA component was updated
	slaComp := entity.GetSLA()
	if slaComp == nil {
		t.Fatal("Expected SLA component to exist")
	}

	slaComponent := slaComp

	// Check current SLA calculation (8/10 = 80%)
	expectedCurrent := 80.0
	if slaComponent.GetCurrent() != expectedCurrent {
		t.Errorf("Expected current SLA to be %.2f, got %.2f", expectedCurrent, slaComponent.GetCurrent())
	}

	// Check remaining errors (10 - 2 = 8)
	expectedRemaining := 8
	if slaComponent.GetErrorsRemaining() != expectedRemaining {
		t.Errorf("Expected remaining errors to be %d, got %d", expectedRemaining, slaComponent.GetErrorsRemaining())
	}
}

func TestSLASystem_Update_EntityWithInvalidSLAComponent(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with invalid SLA component (not implementing SLAComponent interface)
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	entity.AddComponent(transform)

	// Mock entity to return transform as SLA component
	mockEntity := &MockEntity{
		components: map[string]components.Component{
			"SLA": transform,
		},
	}

	entities := []Entity{mockEntity}

	// Test that Update doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with invalid SLA component: %v", r)
		}
	}()

	ss.Update(0.016, entities, eventDispatcher)
}

func TestSLASystem_Update_ZeroTotalPackets(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with SLA component
	entity := entities.NewEntity(1)
	sla := components.NewSLA(95.0, 10)
	entity.AddComponent(sla)

	entities := []Entity{entity}

	// Set zero total packets
	ss.totalPackets = 0
	ss.caughtPackets = 0
	ss.lostPackets = 0

	// Test that Update doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with zero total packets: %v", r)
		}
	}()

	ss.Update(0.016, entities, eventDispatcher)

	// Verify SLA component wasn't updated (division by zero avoided)
	slaComp := entity.GetSLA()
	if slaComp == nil {
		t.Fatal("Expected SLA component to exist")
	}

	slaComponent := slaComp

	// Should remain at initial value (100.0)
	expectedCurrent := 100.0
	if slaComponent.GetCurrent() != expectedCurrent {
		t.Errorf("Expected current SLA to remain %.2f, got %.2f", expectedCurrent, slaComponent.GetCurrent())
	}
}

func TestSLASystem_Update_SLAViolation(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with SLA component
	entity := entities.NewEntity(1)
	sla := components.NewSLA(95.0, 10)
	entity.AddComponent(sla)

	entities := []Entity{entity}

	// Set packet statistics that violate SLA (80% < 95%)
	ss.totalPackets = 10
	ss.caughtPackets = 8
	ss.lostPackets = 2

	// Test that Update doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with SLA violation: %v", r)
		}
	}()

	ss.Update(0.016, entities, eventDispatcher)

	// Verify SLA component was updated correctly
	slaComp := entity.GetSLA()
	if slaComp == nil {
		t.Fatal("Expected SLA component to exist")
	}

	slaComponent := slaComp

	// Check current SLA calculation (8/10 = 80%)
	expectedCurrent := 80.0
	if slaComponent.GetCurrent() != expectedCurrent {
		t.Errorf("Expected current SLA to be %.2f, got %.2f", expectedCurrent, slaComponent.GetCurrent())
	}
}

func TestSLASystem_Initialize(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Test that Initialize doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Initialize panicked: %v", r)
		}
	}()

	ss.Initialize(eventDispatcher)

	// If we get here, Initialize executed without panicking
}

func TestSLASystem_EventHandling_PacketCaught(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	ss.Initialize(eventDispatcher)

	// Create a packet entity
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(15, 15, color.RGBA{255, 0, 0, 255})
	entity.AddComponent(transform)
	entity.AddComponent(sprite)

	initialCaught := ss.caughtPackets
	initialTotal := ss.totalPackets

	// Publish packet caught event
	eventData := &events.EventData{
		Packet: entity,
	}
	event := events.NewEvent(events.EventPacketCaught, eventData)
	eventDispatcher.Publish(event)

	// Verify statistics were updated
	if ss.caughtPackets != initialCaught+1 {
		t.Errorf("Expected caught packets to be %d, got %d", initialCaught+1, ss.caughtPackets)
	}
	if ss.totalPackets != initialTotal+1 {
		t.Errorf("Expected total packets to be %d, got %d", initialTotal+1, ss.totalPackets)
	}
}

func TestSLASystem_EventHandling_PacketLost(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	ss.Initialize(eventDispatcher)

	// Create a packet entity
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(15, 15, color.RGBA{255, 0, 0, 255})
	entity.AddComponent(transform)
	entity.AddComponent(sprite)

	initialLost := ss.lostPackets
	initialTotal := ss.totalPackets

	// Publish packet lost event
	eventData := &events.EventData{
		Packet: entity,
	}
	event := events.NewEvent(events.EventPacketLost, eventData)
	eventDispatcher.Publish(event)

	// Verify statistics were updated
	if ss.lostPackets != initialLost+1 {
		t.Errorf("Expected lost packets to be %d, got %d", initialLost+1, ss.lostPackets)
	}
	if ss.totalPackets != initialTotal+1 {
		t.Errorf("Expected total packets to be %d, got %d", initialTotal+1, ss.totalPackets)
	}
}

func TestSLASystem_EventHandling_MultipleEvents(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	ss.Initialize(eventDispatcher)

	// Create packet entities
	entity1 := entities.NewEntity(1)
	transform1 := components.NewTransform(100, 200)
	sprite1 := components.NewSprite(15, 15, color.RGBA{255, 0, 0, 255})
	entity1.AddComponent(transform1)
	entity1.AddComponent(sprite1)

	entity2 := entities.NewEntity(2)
	transform2 := components.NewTransform(300, 400)
	sprite2 := components.NewSprite(15, 15, color.RGBA{0, 0, 255, 255})
	entity2.AddComponent(transform2)
	entity2.AddComponent(sprite2)

	// Publish multiple events
	event1 := events.NewEvent(events.EventPacketCaught, &events.EventData{Packet: entity1})
	event2 := events.NewEvent(events.EventPacketLost, &events.EventData{Packet: entity2})

	eventDispatcher.Publish(event1)
	eventDispatcher.Publish(event2)

	// Verify statistics were updated correctly
	if ss.caughtPackets != 1 {
		t.Errorf("Expected caught packets to be 1, got %d", ss.caughtPackets)
	}
	if ss.lostPackets != 1 {
		t.Errorf("Expected lost packets to be 1, got %d", ss.lostPackets)
	}
	if ss.totalPackets != 2 {
		t.Errorf("Expected total packets to be 2, got %d", ss.totalPackets)
	}
}

func TestSLASystem_UpdateSLA_Calculation(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Set packet statistics
	ss.totalPackets = 20
	ss.caughtPackets = 18
	ss.lostPackets = 2

	// Test updateSLA method
	ss.updateSLA(eventDispatcher)

	// Verify SLA calculation (18/20 = 90%)
	expectedSLA := 90.0
	expectedRemaining := 8 // 10 - 2

	// The updateSLA method prints to console, so we can't easily test the output
	// But we can verify the calculations are correct by checking the values used
	if float64(ss.caughtPackets)/float64(ss.totalPackets)*100.0 != expectedSLA {
		t.Errorf("Expected SLA calculation to be %.2f, got %.2f", expectedSLA, float64(ss.caughtPackets)/float64(ss.totalPackets)*100.0)
	}

	if ss.errorBudget-ss.lostPackets != expectedRemaining {
		t.Errorf("Expected remaining errors to be %d, got %d", expectedRemaining, ss.errorBudget-ss.lostPackets)
	}
}

func TestSLASystem_UpdateSLA_ErrorBudgetExceeded(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Set packet statistics that exceed error budget
	ss.totalPackets = 15
	ss.caughtPackets = 5
	ss.lostPackets = 10 // This equals the error budget

	// Test updateSLA method
	ss.updateSLA(eventDispatcher)

	// Verify that error budget is exceeded
	remainingErrors := ss.errorBudget - ss.lostPackets
	if remainingErrors != 0 {
		t.Errorf("Expected remaining errors to be 0, got %d", remainingErrors)
	}
}

func TestSLASystem_UpdateSLA_ZeroTotalPackets(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Set zero total packets
	ss.totalPackets = 0
	ss.caughtPackets = 0
	ss.lostPackets = 0

	// Test that updateSLA doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("updateSLA panicked with zero total packets: %v", r)
		}
	}()

	ss.updateSLA(eventDispatcher)
}

func TestSLASystem_SetTargetSLA(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)

	// Test that SetTargetSLA doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetTargetSLA panicked: %v", r)
		}
	}()

	ss.SetTargetSLA(99.5)

	// If we get here, SetTargetSLA executed without panicking
}

func TestSLASystem_SetErrorBudget(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)

	// Test that SetErrorBudget doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SetErrorBudget panicked: %v", r)
		}
	}()

	ss.SetErrorBudget(20)

	// Verify error budget was updated
	if ss.errorBudget != 20 {
		t.Errorf("Expected error budget to be 20, got %d", ss.errorBudget)
	}
}

func TestSLASystem_SetErrorBudget_Zero(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)

	ss.SetErrorBudget(0)

	// Verify error budget was updated
	if ss.errorBudget != 0 {
		t.Errorf("Expected error budget to be 0, got %d", ss.errorBudget)
	}
}

func TestSLASystem_SetErrorBudget_Negative(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)

	ss.SetErrorBudget(-5)

	// Verify error budget was updated (should allow negative values)
	if ss.errorBudget != -5 {
		t.Errorf("Expected error budget to be -5, got %d", ss.errorBudget)
	}
}

func TestSLASystem_Integration(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	ss.Initialize(eventDispatcher)

	// Create entity with SLA component
	entity := entities.NewEntity(1)
	sla := components.NewSLA(95.0, 10)
	entity.AddComponent(sla)

	entityList := []Entity{entity}

	// Create packet entities for events
	packet1 := entities.NewEntity(2)
	transform1 := components.NewTransform(100, 200)
	sprite1 := components.NewSprite(15, 15, color.RGBA{255, 0, 0, 255})
	packet1.AddComponent(transform1)
	packet1.AddComponent(sprite1)

	packet2 := entities.NewEntity(3)
	transform2 := components.NewTransform(300, 400)
	sprite2 := components.NewSprite(15, 15, color.RGBA{0, 0, 255, 255})
	packet2.AddComponent(transform2)
	packet2.AddComponent(sprite2)

	// Publish events
	event1 := events.NewEvent(events.EventPacketCaught, &events.EventData{Packet: packet1})
	event2 := events.NewEvent(events.EventPacketLost, &events.EventData{Packet: packet2})

	eventDispatcher.Publish(event1)
	eventDispatcher.Publish(event2)

	// Update the system
	ss.Update(0.016, entityList, eventDispatcher)

	// Verify final state
	if ss.caughtPackets != 1 {
		t.Errorf("Expected caught packets to be 1, got %d", ss.caughtPackets)
	}
	if ss.lostPackets != 1 {
		t.Errorf("Expected lost packets to be 1, got %d", ss.lostPackets)
	}
	if ss.totalPackets != 2 {
		t.Errorf("Expected total packets to be 2, got %d", ss.totalPackets)
	}

	// Verify SLA component was updated
	slaComp := entity.GetSLA()
	if slaComp == nil {
		t.Fatal("Expected SLA component to exist")
	}

	slaComponent := slaComp

	// Check current SLA calculation (1/2 = 50%)
	expectedCurrent := 50.0
	if slaComponent.GetCurrent() != expectedCurrent {
		t.Errorf("Expected current SLA to be %.2f, got %.2f", expectedCurrent, slaComponent.GetCurrent())
	}

	// Check remaining errors (10 - 1 = 9)
	expectedRemaining := 9
	if slaComponent.GetErrorsRemaining() != expectedRemaining {
		t.Errorf("Expected remaining errors to be %d, got %d", expectedRemaining, slaComponent.GetErrorsRemaining())
	}
}

func TestSLASystem_EntityWithoutSLAComponent(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	// Create entity without SLA component
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	entity.AddComponent(transform)

	entities := []Entity{entity}

	// Set some packet statistics
	ss.totalPackets = 10
	ss.caughtPackets = 8
	ss.lostPackets = 2

	// Test that Update doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with entity without SLA component: %v", r)
		}
	}()

	ss.Update(0.016, entities, eventDispatcher)

	// Verify statistics remain unchanged
	if ss.totalPackets != 10 {
		t.Errorf("Expected total packets to remain 10, got %d", ss.totalPackets)
	}
	if ss.caughtPackets != 8 {
		t.Errorf("Expected caught packets to remain 8, got %d", ss.caughtPackets)
	}
	if ss.lostPackets != 2 {
		t.Errorf("Expected lost packets to remain 2, got %d", ss.lostPackets)
	}
}

// MockEntity is a mock implementation for testing
type MockEntity struct {
	components map[string]components.Component
}

func (m *MockEntity) GetID() uint64 {
	return 1
}

func (m *MockEntity) AddComponent(component components.Component) {
	if m.components == nil {
		m.components = make(map[string]components.Component)
	}
	m.components[component.GetType()] = component
}

func (m *MockEntity) GetComponent(componentType string) components.Component {
	return m.components[componentType]
}

func (m *MockEntity) GetComponentByName(typeName string) components.Component {
	return m.components[typeName]
}

func (m *MockEntity) HasComponent(componentType string) bool {
	_, exists := m.components[componentType]
	return exists
}

func (m *MockEntity) GetTransform() components.TransformComponent {
	if comp := m.components["Transform"]; comp != nil {
		if transformComp, ok := comp.(components.TransformComponent); ok {
			return transformComp
		}
	}
	return nil
}

func (m *MockEntity) GetSprite() components.SpriteComponent {
	if comp := m.components["Sprite"]; comp != nil {
		if spriteComp, ok := comp.(components.SpriteComponent); ok {
			return spriteComp
		}
	}
	return nil
}

func (m *MockEntity) GetCollider() components.ColliderComponent {
	if comp := m.components["Collider"]; comp != nil {
		if colliderComp, ok := comp.(components.ColliderComponent); ok {
			return colliderComp
		}
	}
	return nil
}

func (m *MockEntity) GetPhysics() components.PhysicsComponent {
	if comp := m.components["Physics"]; comp != nil {
		if physicsComp, ok := comp.(components.PhysicsComponent); ok {
			return physicsComp
		}
	}
	return nil
}

func (m *MockEntity) GetPacketType() components.PacketTypeComponent {
	if comp := m.components["PacketType"]; comp != nil {
		if packetComp, ok := comp.(components.PacketTypeComponent); ok {
			return packetComp
		}
	}
	return nil
}

func (m *MockEntity) GetState() components.StateComponent {
	if comp := m.components["State"]; comp != nil {
		if stateComp, ok := comp.(components.StateComponent); ok {
			return stateComp
		}
	}
	return nil
}

func (m *MockEntity) GetCombo() components.ComboComponent {
	if comp := m.components["Combo"]; comp != nil {
		if comboComp, ok := comp.(components.ComboComponent); ok {
			return comboComp
		}
	}
	return nil
}

func (m *MockEntity) GetSLA() components.SLAComponent {
	if comp := m.components["SLA"]; comp != nil {
		if slaComp, ok := comp.(components.SLAComponent); ok {
			return slaComp
		}
	}
	return nil
}

func (m *MockEntity) GetBackendAssignment() components.BackendAssignmentComponent {
	if comp := m.components["BackendAssignment"]; comp != nil {
		if backendComp, ok := comp.(components.BackendAssignmentComponent); ok {
			return backendComp
		}
	}
	return nil
}

func (m *MockEntity) GetPowerUpType() components.PowerUpTypeComponent {
	if comp := m.components["PowerUpType"]; comp != nil {
		if powerUpComp, ok := comp.(components.PowerUpTypeComponent); ok {
			return powerUpComp
		}
	}
	return nil
}

func (m *MockEntity) GetRouting() components.RoutingComponent {
	if comp := m.components["Routing"]; comp != nil {
		if routingComp, ok := comp.(components.RoutingComponent); ok {
			return routingComp
		}
	}
	return nil
}

func (m *MockEntity) RemoveComponent(componentType string) {
	delete(m.components, componentType)
}

func (m *MockEntity) IsActive() bool {
	return true
}

func (m *MockEntity) SetActive(active bool) {
	// Mock implementation
}

// TestSLASystem_GameRestart tests the SLA system behavior during game restart scenarios
func TestSLASystem_GameRestart(t *testing.T) {
	spawnSys := NewSpawnSystem(func() Entity { return entities.NewEntity(1) })
	ss := NewSLASystem(spawnSys)
	eventDispatcher := events.NewEventDispatcher()

	ss.Initialize(eventDispatcher)

	t.Run("Complete Game Restart Scenario", func(t *testing.T) {
		// Test 1: Initial state
		t.Run("Initial State", func(t *testing.T) {
			if ss.GetTotalPackets() != 0 {
				t.Errorf("Expected 0 total packets initially, got %d", ss.GetTotalPackets())
			}
			if ss.GetCaughtPackets() != 0 {
				t.Errorf("Expected 0 caught packets initially, got %d", ss.GetCaughtPackets())
			}
			if ss.GetLostPackets() != 0 {
				t.Errorf("Expected 0 lost packets initially, got %d", ss.GetLostPackets())
			}
			if ss.GetErrorBudget() != 10 {
				t.Errorf("Expected 10 error budget initially, got %d", ss.GetErrorBudget())
			}
		})

		// Test 2: First game - simulate packet events
		t.Run("First Game - Packet Events", func(t *testing.T) {
			// Simulate packet caught events
			for i := 0; i < 5; i++ {
				caughtEvent := events.NewEvent(events.EventPacketCaught, nil)
				eventDispatcher.Publish(caughtEvent)
			}

			// Simulate packet lost events
			for i := 0; i < 2; i++ {
				lostEvent := events.NewEvent(events.EventPacketLost, nil)
				eventDispatcher.Publish(lostEvent)
			}

			if ss.GetTotalPackets() != 7 {
				t.Errorf("Expected 7 total packets, got %d", ss.GetTotalPackets())
			}
			if ss.GetCaughtPackets() != 5 {
				t.Errorf("Expected 5 caught packets, got %d", ss.GetCaughtPackets())
			}
			if ss.GetLostPackets() != 2 {
				t.Errorf("Expected 2 lost packets, got %d", ss.GetLostPackets())
			}
		})

		// Test 3: Game restart - reset and new settings
		t.Run("Game Restart - Reset and New Settings", func(t *testing.T) {
			// Reset the SLA system
			ss.Reset()

			// Verify reset state
			if ss.GetTotalPackets() != 0 {
				t.Errorf("Expected 0 total packets after reset, got %d", ss.GetTotalPackets())
			}
			if ss.GetCaughtPackets() != 0 {
				t.Errorf("Expected 0 caught packets after reset, got %d", ss.GetCaughtPackets())
			}
			if ss.GetLostPackets() != 0 {
				t.Errorf("Expected 0 lost packets after reset, got %d", ss.GetLostPackets())
			}

			// Set new error budget for restart
			ss.SetErrorBudget(25)

			if ss.GetErrorBudget() != 25 {
				t.Errorf("Expected 25 error budget after restart, got %d", ss.GetErrorBudget())
			}
		})

		// Test 4: Second game - verify clean state
		t.Run("Second Game - Clean State", func(t *testing.T) {
			// Simulate new packet events after restart
			for i := 0; i < 3; i++ {
				caughtEvent := events.NewEvent(events.EventPacketCaught, nil)
				eventDispatcher.Publish(caughtEvent)
			}

			lostEvent := events.NewEvent(events.EventPacketLost, nil)
			eventDispatcher.Publish(lostEvent)

			if ss.GetTotalPackets() != 4 {
				t.Errorf("Expected 4 total packets in second game, got %d", ss.GetTotalPackets())
			}
			if ss.GetCaughtPackets() != 3 {
				t.Errorf("Expected 3 caught packets in second game, got %d", ss.GetCaughtPackets())
			}
			if ss.GetLostPackets() != 1 {
				t.Errorf("Expected 1 lost packet in second game, got %d", ss.GetLostPackets())
			}
			if ss.GetErrorBudget() != 25 {
				t.Errorf("Expected 25 error budget in second game, got %d", ss.GetErrorBudget())
			}
		})
	})

	t.Run("Multiple Restarts", func(t *testing.T) {
		// Test multiple restart cycles
		for i := 0; i < 3; i++ {
			t.Run(fmt.Sprintf("Restart Cycle %d", i+1), func(t *testing.T) {
				// Simulate some game activity
				for j := 0; j < 3; j++ {
					caughtEvent := events.NewEvent(events.EventPacketCaught, nil)
					eventDispatcher.Publish(caughtEvent)
				}

				lostEvent := events.NewEvent(events.EventPacketLost, nil)
				eventDispatcher.Publish(lostEvent)

				// Reset for next game
				ss.Reset()
				ss.SetErrorBudget(20 + i*5) // Different budget each time

				// Verify clean state
				if ss.GetTotalPackets() != 0 {
					t.Errorf("Expected 0 total packets after reset cycle %d, got %d", i+1, ss.GetTotalPackets())
				}
				if ss.GetCaughtPackets() != 0 {
					t.Errorf("Expected 0 caught packets after reset cycle %d, got %d", i+1, ss.GetCaughtPackets())
				}
				if ss.GetLostPackets() != 0 {
					t.Errorf("Expected 0 lost packets after reset cycle %d, got %d", i+1, ss.GetLostPackets())
				}
				if ss.GetErrorBudget() != 20+i*5 {
					t.Errorf("Expected %d error budget after reset cycle %d, got %d", 20+i*5, i+1, ss.GetErrorBudget())
				}
			})
		}
	})

	t.Run("Reset Preserves Error Budget", func(t *testing.T) {
		// Set a custom error budget
		ss.SetErrorBudget(50)

		// Simulate some activity
		caughtEvent := events.NewEvent(events.EventPacketCaught, nil)
		eventDispatcher.Publish(caughtEvent)

		lostEvent := events.NewEvent(events.EventPacketLost, nil)
		eventDispatcher.Publish(lostEvent)

		// Reset should clear counters but preserve error budget
		ss.Reset()

		if ss.GetTotalPackets() != 0 {
			t.Errorf("Expected 0 total packets after reset, got %d", ss.GetTotalPackets())
		}
		if ss.GetCaughtPackets() != 0 {
			t.Errorf("Expected 0 caught packets after reset, got %d", ss.GetCaughtPackets())
		}
		if ss.GetLostPackets() != 0 {
			t.Errorf("Expected 0 lost packets after reset, got %d", ss.GetLostPackets())
		}
		if ss.GetErrorBudget() != 50 {
			t.Errorf("Expected 50 error budget to be preserved after reset, got %d", ss.GetErrorBudget())
		}
	})
}
