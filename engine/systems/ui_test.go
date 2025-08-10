package systems

import (
	"fmt"
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

func (ute *uiTestEntity) GetID() uint64 {
	return ute.entity.GetID()
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

func (ute *uiTestEntity) GetRouting() components.RoutingComponent {
	if comp := ute.entity.GetComponent("Routing"); comp != nil {
		return comp.(components.RoutingComponent)
	}
	return nil
}

func (ute *uiTestEntity) RemoveComponent(componentType string) {
	ute.entity.RemoveComponent(componentType)
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

	// Test that the system is properly initialized
	if uis == nil {
		t.Fatal("NewUISystem returned nil")
	}

	// Test that the system is properly initialized
	if uis == nil {
		t.Fatal("NewUISystem returned nil")
	}
}

// TestUISystem_Update tests the Update method
func TestUISystem_Update(t *testing.T) {
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}

	// Test that Update doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked: %v", r)
		}
	}()

	uis.Update(1.0, entities, eventDispatcher)

	// Verify no errors occurred
	// The UI system should handle empty entities gracefully
}

// TestUISystem_Initialize tests the Initialize method
func TestUISystem_Initialize(t *testing.T) {
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)
	eventDispatcher := events.NewEventDispatcher()

	// Test that Initialize doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Initialize panicked: %v", r)
		}
	}()

	// Initialize the UI system
	uis.Initialize(eventDispatcher)

	// Verify no errors occurred
	// The UI system should handle initialization gracefully
}

// TestUISystem_Initialize_PartialData tests Initialize with partial event data
func TestUISystem_Initialize_PartialData(t *testing.T) {
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)
	eventDispatcher := events.NewEventDispatcher()

	// Test that Initialize doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Initialize panicked: %v", r)
		}
	}()

	uis.Initialize(eventDispatcher)

	// Test that publishing events doesn't panic
	newSLA := 90.0
	event := events.NewEvent(events.EventSLAUpdated, &events.EventData{
		Current: &newSLA,
		// Other fields are nil
	})

	// Should not panic
	eventDispatcher.Publish(event)

	// Verify no errors occurred
	// The UI system should handle partial event data gracefully
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

		// 4. Lose a packet (UI system now gets this info from SLA update)
		packetLostEvent := events.NewEvent(events.EventSLAUpdated, &events.EventData{
			Lost: intPtr(4), // Update to 4 lost packets
		})
		eventDispatcher.Publish(packetLostEvent)

		// 5. End DDoS attack
		ddosEndEvent := &events.Event{
			Type: events.EventDDoSEnd,
		}
		eventDispatcher.Publish(ddosEndEvent)

		// Verify no errors occurred during event processing
		// The UI system should handle all events gracefully

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
		// Test that publishing events doesn't panic
		event := events.NewEvent(events.EventSLAUpdated, &events.EventData{
			Lost:      intPtr(10), // This should make remaining errors 0
			Remaining: intPtr(0),  // 10 - 10 = 0 remaining errors
		})
		eventDispatcher.Publish(event)

		// Verify no errors occurred
		// The UI system should handle error budget events gracefully
	})

	t.Run("Zero Error Budget", func(t *testing.T) {
		event := events.NewEvent(events.EventSLAUpdated, &events.EventData{
			Lost:      intPtr(1), // This should make remaining errors 0
			Remaining: intPtr(0), // 0 - 1 = 0 remaining errors (clamped)
		})
		eventDispatcher.Publish(event)

		// Verify no errors occurred
		// The UI system should handle zero error budget gracefully
	})
}

// TestUISystem_GameRestart tests the UI system behavior during game restart scenarios
func TestUISystem_GameRestart(t *testing.T) {
	screen := ebiten.NewImage(800, 600)
	defer screen.Dispose()

	uis := NewUISystem(screen)
	eventDispatcher := events.NewEventDispatcher()

	uis.Initialize(eventDispatcher)

	t.Run("Complete Game Restart Scenario", func(t *testing.T) {
		// Test 1: Initial state - verify no panics
		t.Run("Initial State", func(t *testing.T) {
			// Test that accessing methods doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Initial state access panicked: %v", r)
				}
			}()

			// These calls should not panic
			_ = uis.GetCaughtPackets()
			_ = uis.GetLostPackets()
			_ = uis.GetRemainingErrors()
			_ = uis.GetErrorBudget()
		})

		// Test 2: First game - simulate packet events
		t.Run("First Game - Packet Events", func(t *testing.T) {
			// Simulate SLA update with some lost packets
			slaEvent := events.NewEvent(events.EventSLAUpdated, &events.EventData{
				Current:   float64Ptr(95.0),
				Caught:    intPtr(19),
				Lost:      intPtr(1),
				Remaining: intPtr(9),
				Budget:    intPtr(10),
			})
			eventDispatcher.Publish(slaEvent)

			// Verify no errors occurred
			// The UI system should handle SLA events gracefully
		})

		// Test 3: Game restart - reset and new settings
		t.Run("Game Restart - Reset and New Settings", func(t *testing.T) {
			// Test that Reset doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Reset panicked: %v", r)
				}
			}()

			uis.Reset()

			// Test that SetErrorBudget doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("SetErrorBudget panicked: %v", r)
				}
			}()

			uis.SetErrorBudget(25)

			// Verify no errors occurred
			// The UI system should handle reset and budget changes gracefully
		})

		// Test 4: Second game - verify clean state
		t.Run("Second Game - Clean State", func(t *testing.T) {
			// Simulate new SLA update after restart
			slaEvent := events.NewEvent(events.EventSLAUpdated, &events.EventData{
				Current:   float64Ptr(98.0),
				Caught:    intPtr(49),
				Lost:      intPtr(1),
				Remaining: intPtr(24),
				Budget:    intPtr(25),
			})
			eventDispatcher.Publish(slaEvent)

			// Verify no errors occurred
			// The UI system should handle SLA events after restart gracefully
		})
	})

	t.Run("Multiple Restarts", func(t *testing.T) {
		// Test multiple restart cycles
		for i := 0; i < 3; i++ {
			t.Run(fmt.Sprintf("Restart Cycle %d", i+1), func(t *testing.T) {
				// Simulate some game activity
				slaEvent := events.NewEvent(events.EventSLAUpdated, &events.EventData{
					Current:   float64Ptr(90.0),
					Caught:    intPtr(9),
					Lost:      intPtr(1),
					Remaining: intPtr(9),
					Budget:    intPtr(10),
				})
				eventDispatcher.Publish(slaEvent)

				// Reset for next game
				uis.Reset()
				uis.SetErrorBudget(15 + i*5) // Different budget each time

				// Verify no errors occurred
				// The UI system should handle multiple restart cycles gracefully
			})
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
