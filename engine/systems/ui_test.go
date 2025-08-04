package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// uiTestEntity implements the Entity interface for UI testing
type uiTestEntity struct {
	entity          *entities.Entity
	componentsAdded []string
	isActive        bool
}

func newUITestEntity(id uint64) *uiTestEntity {
	return &uiTestEntity{
		entity:          entities.NewEntity(id),
		componentsAdded: make([]string, 0),
		isActive:        true,
	}
}

func (ute *uiTestEntity) GetComponent(componentType string) components.Component {
	return ute.entity.GetComponent(componentType)
}

func (ute *uiTestEntity) HasComponent(componentType string) bool {
	return ute.entity.HasComponent(componentType)
}

func (ute *uiTestEntity) IsActive() bool {
	return ute.isActive
}

func (ute *uiTestEntity) GetComponentByName(typeName string) components.Component {
	return ute.entity.GetComponentByName(typeName)
}

func (ute *uiTestEntity) GetTransform() components.TransformComponent {
	if comp := ute.entity.GetComponent("Transform"); comp != nil {
		return comp.(components.TransformComponent)
	}
	return nil
}

func (ute *uiTestEntity) GetSprite() components.SpriteComponent {
	if comp := ute.entity.GetComponent("Sprite"); comp != nil {
		return comp.(components.SpriteComponent)
	}
	return nil
}

func (ute *uiTestEntity) GetCollider() components.ColliderComponent {
	if comp := ute.entity.GetComponent("Collider"); comp != nil {
		return comp.(components.ColliderComponent)
	}
	return nil
}

func (ute *uiTestEntity) GetPhysics() components.PhysicsComponent {
	if comp := ute.entity.GetComponent("Physics"); comp != nil {
		return comp.(components.PhysicsComponent)
	}
	return nil
}

func (ute *uiTestEntity) GetPacketType() components.PacketTypeComponent {
	if comp := ute.entity.GetComponent("PacketType"); comp != nil {
		return comp.(components.PacketTypeComponent)
	}
	return nil
}

func (ute *uiTestEntity) GetState() components.StateComponent {
	if comp := ute.entity.GetComponent("State"); comp != nil {
		return comp.(components.StateComponent)
	}
	return nil
}

func (ute *uiTestEntity) GetCombo() components.ComboComponent {
	if comp := ute.entity.GetComponent("Combo"); comp != nil {
		return comp.(components.ComboComponent)
	}
	return nil
}

func (ute *uiTestEntity) GetSLA() components.SLAComponent {
	if comp := ute.entity.GetComponent("SLA"); comp != nil {
		return comp.(components.SLAComponent)
	}
	return nil
}

func (ute *uiTestEntity) GetBackendAssignment() components.BackendAssignmentComponent {
	if comp := ute.entity.GetComponent("BackendAssignment"); comp != nil {
		return comp.(components.BackendAssignmentComponent)
	}
	return nil
}

func (ute *uiTestEntity) GetPowerUpType() components.PowerUpTypeComponent {
	if comp := ute.entity.GetComponent("PowerUpType"); comp != nil {
		return comp.(components.PowerUpTypeComponent)
	}
	return nil
}

func (ute *uiTestEntity) AddComponent(component components.Component) {
	ute.componentsAdded = append(ute.componentsAdded, component.GetType())
	ute.entity.AddComponent(component)
}

func (ute *uiTestEntity) GetComponentNames() []string {
	return ute.componentsAdded
}

func (ute *uiTestEntity) SetActive(active bool) {
	ute.isActive = active
}

// TestNewUISystem tests the UISystem constructor
func TestNewUISystem(t *testing.T) {
	// Create a dummy screen for testing
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)

	// Test initial values
	if uis.score != 0 {
		t.Errorf("Expected initial score to be 0, got %d", uis.score)
	}

	if uis.currentSLA != 100.0 {
		t.Errorf("Expected initial currentSLA to be 100.0, got %f", uis.currentSLA)
	}

	if uis.targetSLA != 99.5 {
		t.Errorf("Expected initial targetSLA to be 99.5, got %f", uis.targetSLA)
	}

	if uis.caughtPackets != 0 {
		t.Errorf("Expected initial caughtPackets to be 0, got %d", uis.caughtPackets)
	}

	if uis.lostPackets != 0 {
		t.Errorf("Expected initial lostPackets to be 0, got %d", uis.lostPackets)
	}

	if uis.remainingErrors != 10 {
		t.Errorf("Expected initial remainingErrors to be 10, got %d", uis.remainingErrors)
	}

	if uis.errorBudget != 10 {
		t.Errorf("Expected initial errorBudget to be 10, got %d", uis.errorBudget)
	}

	if uis.level != 1 {
		t.Errorf("Expected initial level to be 1, got %d", uis.level)
	}

	if uis.isDDoSActive {
		t.Error("Expected initial isDDoSActive to be false")
	}
}

// TestUISystem_Update tests the Update method
func TestUISystem_Update(t *testing.T) {
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}

	// Test that Update doesn't panic and doesn't change state
	initialScore := uis.score
	initialSLA := uis.currentSLA

	uis.Update(1.0, entities, eventDispatcher)

	if uis.score != initialScore {
		t.Errorf("Update should not change score, expected %d, got %d", initialScore, uis.score)
	}

	if uis.currentSLA != initialSLA {
		t.Errorf("Update should not change currentSLA, expected %f, got %f", initialSLA, uis.currentSLA)
	}
}

// TestUISystem_Initialize tests the Initialize method
func TestUISystem_Initialize(t *testing.T) {
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the UI system
	uis.Initialize(eventDispatcher)

	// Test that event handlers are registered by dispatching events
	t.Run("SLA Updated Event", func(t *testing.T) {
		newSLA := 95.5
		newTarget := 98.0
		newCaught := 15
		newLost := 2
		newRemaining := 8
		newBudget := 10

		event := events.NewEvent(events.EventSLAUpdated, &events.EventData{
			Current:   &newSLA,
			Target:    &newTarget,
			Caught:    &newCaught,
			Lost:      &newLost,
			Remaining: &newRemaining,
			Budget:    &newBudget,
		})

		eventDispatcher.Publish(event)

		if uis.currentSLA != newSLA {
			t.Errorf("Expected currentSLA to be %f, got %f", newSLA, uis.currentSLA)
		}

		if uis.targetSLA != newTarget {
			t.Errorf("Expected targetSLA to be %f, got %f", newTarget, uis.targetSLA)
		}

		if uis.caughtPackets != newCaught {
			t.Errorf("Expected caughtPackets to be %d, got %d", newCaught, uis.caughtPackets)
		}

		if uis.lostPackets != newLost {
			t.Errorf("Expected lostPackets to be %d, got %d", newLost, uis.lostPackets)
		}

		if uis.remainingErrors != newRemaining {
			t.Errorf("Expected remainingErrors to be %d, got %d", newRemaining, uis.remainingErrors)
		}

		if uis.errorBudget != newBudget {
			t.Errorf("Expected errorBudget to be %d, got %d", newBudget, uis.errorBudget)
		}
	})

	t.Run("Packet Lost Event", func(t *testing.T) {
		initialLost := uis.lostPackets

		event := events.NewEvent(events.EventPacketLost, nil)

		eventDispatcher.Publish(event)

		if uis.lostPackets != initialLost+1 {
			t.Errorf("Expected lostPackets to be %d, got %d", initialLost+1, uis.lostPackets)
		}

		expectedRemaining := uis.errorBudget - uis.lostPackets
		if expectedRemaining < 0 {
			expectedRemaining = 0
		}

		if uis.remainingErrors != expectedRemaining {
			t.Errorf("Expected remainingErrors to be %d, got %d", expectedRemaining, uis.remainingErrors)
		}
	})

	t.Run("Level Up Event", func(t *testing.T) {
		newLevel := 5

		event := events.NewEvent(events.EventLevelUp, &events.EventData{
			Level: &newLevel,
		})

		eventDispatcher.Publish(event)

		if uis.level != newLevel {
			t.Errorf("Expected level to be %d, got %d", newLevel, uis.level)
		}
	})

	t.Run("DDoS Start Event", func(t *testing.T) {
		event := &events.Event{
			Type: events.EventDDoSStart,
		}

		eventDispatcher.Publish(event)

		if !uis.isDDoSActive {
			t.Error("Expected isDDoSActive to be true after DDoS start event")
		}
	})

	t.Run("DDoS End Event", func(t *testing.T) {
		event := &events.Event{
			Type: events.EventDDoSEnd,
		}

		eventDispatcher.Publish(event)

		if uis.isDDoSActive {
			t.Error("Expected isDDoSActive to be false after DDoS end event")
		}
	})
}

// TestUISystem_Initialize_PartialData tests Initialize with partial event data
func TestUISystem_Initialize_PartialData(t *testing.T) {
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)
	eventDispatcher := events.NewEventDispatcher()

	uis.Initialize(eventDispatcher)

	// Test with only some fields set
	newSLA := 90.0
	event := events.NewEvent(events.EventSLAUpdated, &events.EventData{
		Current: &newSLA,
		// Other fields are nil
	})

	eventDispatcher.Publish(event)

	if uis.currentSLA != newSLA {
		t.Errorf("Expected currentSLA to be %f, got %f", newSLA, uis.currentSLA)
	}

	// Other fields should remain unchanged
	if uis.targetSLA != 99.5 {
		t.Errorf("Expected targetSLA to remain 99.5, got %f", uis.targetSLA)
	}

	if uis.caughtPackets != 0 {
		t.Errorf("Expected caughtPackets to remain 0, got %d", uis.caughtPackets)
	}
}

// TestUISystem_Draw tests the Draw method
func TestUISystem_Draw(t *testing.T) {
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)
	entities := []Entity{}

	// Test that Draw doesn't panic
	uis.Draw(screen, entities)

	// Test with entities containing different components
	t.Run("With Combo Component", func(t *testing.T) {
		entity := newUITestEntity(1)
		combo := components.NewCombo()
		combo.Streak = 3
		entity.AddComponent(combo)

		entities := []Entity{entity}

		// Should not panic
		uis.Draw(screen, entities)
	})

	t.Run("With State Component", func(t *testing.T) {
		entity := newUITestEntity(1)
		state := components.NewState(components.StatePlaying)
		entity.AddComponent(state)

		entities := []Entity{entity}

		// Should not panic
		uis.Draw(screen, entities)
	})

	t.Run("With Backend Assignment Component", func(t *testing.T) {
		entity := newUITestEntity(1)
		backend := components.NewBackendAssignment(1)
		backend.Counter = 5
		entity.AddComponent(backend)

		entities := []Entity{entity}

		// Should not panic
		uis.Draw(screen, entities)
	})

	t.Run("With Multiple Components", func(t *testing.T) {
		entity1 := newUITestEntity(1)
		combo := components.NewCombo()
		combo.Streak = 2
		entity1.AddComponent(combo)

		entity2 := newUITestEntity(2)
		state := components.NewState(components.StatePlaying)
		entity2.AddComponent(state)

		entity3 := newUITestEntity(3)
		backend := components.NewBackendAssignment(2)
		backend.Counter = 10
		entity3.AddComponent(backend)

		entities := []Entity{entity1, entity2, entity3}

		// Should not panic
		uis.Draw(screen, entities)
	})
}

// TestUISystem_Draw_EdgeCases tests edge cases for the Draw method
func TestUISystem_Draw_EdgeCases(t *testing.T) {
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)

	t.Run("Nil Screen", func(t *testing.T) {
		// This should panic with nil screen
		entities := []Entity{}
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic with nil screen")
			}
		}()
		uis.Draw(nil, entities)
	})

	t.Run("Nil Entities", func(t *testing.T) {
		// This should not panic
		uis.Draw(screen, nil)
	})

	t.Run("Empty Entities", func(t *testing.T) {
		// This should not panic
		entities := []Entity{}
		uis.Draw(screen, entities)
	})

	t.Run("Entity Without Required Components", func(t *testing.T) {
		entity := newUITestEntity(1)
		// Don't add any components

		entities := []Entity{entity}

		// Should not panic
		uis.Draw(screen, entities)
	})

	t.Run("Inactive Entity", func(t *testing.T) {
		entity := newUITestEntity(1)
		entity.SetActive(false)
		combo := components.NewCombo()
		combo.Streak = 5
		entity.AddComponent(combo)

		entities := []Entity{entity}

		// Should not panic
		uis.Draw(screen, entities)
	})
}

// TestUISystem_Integration tests integration scenarios
func TestUISystem_Integration(t *testing.T) {
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	uis.Initialize(eventDispatcher)

	t.Run("Complete Game Flow", func(t *testing.T) {
		// Simulate game progression
		// 1. Update SLA
		slaEvent := events.NewEvent(events.EventSLAUpdated, &events.EventData{
			Current:   float64Ptr(95.0),
			Target:    float64Ptr(99.5),
			Caught:    intPtr(20),
			Lost:      intPtr(3),
			Remaining: intPtr(7),
			Budget:    intPtr(10),
		})
		eventDispatcher.Publish(slaEvent)

		// 2. Level up
		levelEvent := events.NewEvent(events.EventLevelUp, &events.EventData{
			Level: intPtr(3),
		})
		eventDispatcher.Publish(levelEvent)

		// 3. Start DDoS attack
		ddosStartEvent := &events.Event{
			Type: events.EventDDoSStart,
		}
		eventDispatcher.Publish(ddosStartEvent)

		// 4. Lose a packet
		packetLostEvent := &events.Event{
			Type: events.EventPacketLost,
		}
		eventDispatcher.Publish(packetLostEvent)

		// 5. End DDoS attack
		ddosEndEvent := &events.Event{
			Type: events.EventDDoSEnd,
		}
		eventDispatcher.Publish(ddosEndEvent)

		// Verify final state
		if uis.currentSLA != 95.0 {
			t.Errorf("Expected currentSLA to be 95.0, got %f", uis.currentSLA)
		}

		if uis.targetSLA != 99.5 {
			t.Errorf("Expected targetSLA to be 99.5, got %f", uis.targetSLA)
		}

		if uis.caughtPackets != 20 {
			t.Errorf("Expected caughtPackets to be 20, got %d", uis.caughtPackets)
		}

		if uis.lostPackets != 4 { // 3 from SLA event + 1 from packet lost event
			t.Errorf("Expected lostPackets to be 4, got %d", uis.lostPackets)
		}

		if uis.level != 3 {
			t.Errorf("Expected level to be 3, got %d", uis.level)
		}

		if uis.isDDoSActive {
			t.Error("Expected isDDoSActive to be false after DDoS end")
		}

		// Test drawing with entities
		entity1 := newUITestEntity(1)
		combo := components.NewCombo()
		combo.Streak = 5
		entity1.AddComponent(combo)

		entity2 := newUITestEntity(2)
		state := components.NewState(components.StatePlaying)
		entity2.AddComponent(state)

		entity3 := newUITestEntity(3)
		backend := components.NewBackendAssignment(1)
		backend.Counter = 15
		entity3.AddComponent(backend)

		entities := []Entity{entity1, entity2, entity3}

		// Should not panic
		uis.Draw(screen, entities)
	})
}

// TestUISystem_ErrorBudgetEdgeCases tests error budget edge cases
func TestUISystem_ErrorBudgetEdgeCases(t *testing.T) {
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)
	eventDispatcher := events.NewEventDispatcher()

	uis.Initialize(eventDispatcher)

	t.Run("Error Budget Exceeded", func(t *testing.T) {
		// Set up initial state with 1 error remaining
		uis.remainingErrors = 1
		uis.errorBudget = 10
		uis.lostPackets = 9

		// Lose another packet
		event := &events.Event{
			Type: events.EventPacketLost,
		}
		eventDispatcher.Publish(event)

		// Should have 0 remaining errors
		if uis.remainingErrors != 0 {
			t.Errorf("Expected remainingErrors to be 0, got %d", uis.remainingErrors)
		}

		// Lose more packets
		eventDispatcher.Publish(event)
		eventDispatcher.Publish(event)

		// Should still be 0 (not negative)
		if uis.remainingErrors != 0 {
			t.Errorf("Expected remainingErrors to remain 0, got %d", uis.remainingErrors)
		}
	})

	t.Run("Zero Error Budget", func(t *testing.T) {
		uis.errorBudget = 0
		uis.remainingErrors = 0
		uis.lostPackets = 0

		event := &events.Event{
			Type: events.EventPacketLost,
		}
		eventDispatcher.Publish(event)

		// Should remain 0
		if uis.remainingErrors != 0 {
			t.Errorf("Expected remainingErrors to be 0, got %d", uis.remainingErrors)
		}
	})
}

// Helper functions for creating pointers
func float64Ptr(v float64) *float64 {
	return &v
}

func intPtr(v int) *int {
	return &v
}
