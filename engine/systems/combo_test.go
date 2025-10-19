package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"
)

func TestNewComboSystem(t *testing.T) {
	cs := NewComboSystem()

	// Test that the system is properly initialized
	if cs == nil {
		t.Fatal("NewComboSystem returned nil")
	}

	// Test required components
	expectedComponents := []string{"Combo"}
	if len(cs.RequiredComponents) != len(expectedComponents) {
		t.Errorf("Expected %d required components, got %d", len(expectedComponents), len(cs.RequiredComponents))
	}

	for i, component := range expectedComponents {
		if cs.RequiredComponents[i] != component {
			t.Errorf("Expected required component %s at index %d, got %s", component, i, cs.RequiredComponents[i])
		}
	}
}

func TestComboSystem_Update_NoEntities(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test with no entities
	entities := []Entity{}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// System is stateless, so no internal state to verify
}

func TestComboSystem_Update_WithEntities(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with combo component
	entity := createComboEntity(1)
	entities := []Entity{entity}

	// Set some combo state in the component
	comboComp := entity.GetCombo()
	comboObj := comboComp.(*components.Combo)
	comboObj.Streak = 5
	comboObj.Timer = 2.0

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify combo component timer was updated
	if comboObj.Timer != 2.016 {
		t.Errorf("Expected combo timer to be 2.016, got %f", comboObj.Timer)
	}

	// Verify combo streak remains unchanged
	if comboObj.Streak != 5 {
		t.Errorf("Expected combo streak to remain 5, got %d", comboObj.Streak)
	}
}

func TestComboSystem_Update_ComboExpiration(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with combo component that should expire
	entity := createComboEntity(1)
	entities := []Entity{entity}

	// Set up combo that should expire (timer > 3.0 seconds)
	comboComp := entity.GetCombo()
	comboObj := comboComp.(*components.Combo)
	comboObj.Streak = 5
	comboObj.Timer = 4.0 // More than 3.0 second timeout

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify combo was reset
	if comboObj.Streak != 0 {
		t.Errorf("Expected combo streak to be reset to 0, got %d", comboObj.Streak)
	}

	if comboObj.Timer != 0 {
		t.Errorf("Expected combo timer to be reset to 0, got %f", comboObj.Timer)
	}
}

func TestComboSystem_Update_ComboNotExpired(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with combo component that should not expire
	entity := createComboEntity(1)
	entities := []Entity{entity}

	// Set up combo that should not expire (timer < 3.0 seconds)
	comboComp := entity.GetCombo()
	comboObj := comboComp.(*components.Combo)
	comboObj.Streak = 3
	comboObj.Timer = 2.0 // Less than 3.0 second timeout

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify combo was not reset
	if comboObj.Streak != 3 {
		t.Errorf("Expected combo streak to remain 3, got %d", comboObj.Streak)
	}

	// Verify timer was incremented
	if comboObj.Timer != 2.016 {
		t.Errorf("Expected combo timer to be 2.016, got %f", comboObj.Timer)
	}
}

func TestComboSystem_Update_ComboExpirationSingleCombo(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with single combo that expires
	entity := createComboEntity(1)
	entities := []Entity{entity}

	// Set up single combo that expires (should not print message)
	comboComp := entity.GetCombo()
	comboObj := comboComp.(*components.Combo)
	comboObj.Streak = 1
	comboObj.Timer = 4.0 // More than 3.0 second timeout

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify combo was reset
	if comboObj.Streak != 0 {
		t.Errorf("Expected combo streak to be reset to 0, got %d", comboObj.Streak)
	}
}

func TestComboSystem_Initialize(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system (should be a no-op for stateless system)
	cs.Initialize(eventDispatcher)

	// Verify no errors occurred
	// System is stateless, so no internal state to verify
}

func TestComboSystem_EntityWithoutComboComponent(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity without combo component
	entity := entities.NewEntity(1)
	entities := []Entity{entity}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify no errors occurred
	// System should handle entities without combo components gracefully
}

func TestComboSystem_Integration(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	cs.Initialize(eventDispatcher)

	// Create entity with combo component
	entity := createComboEntity(1)
	entities := []Entity{entity}

	// Simulate game loop with combo updates
	comboComp := entity.GetCombo()
	comboObj := comboComp.(*components.Combo)

	// First update - combo should not expire
	cs.Update(0.5, entities, eventDispatcher)

	// Verify combo state
	if comboObj.Streak != 0 {
		t.Errorf("Expected combo streak to be 0, got %d", comboObj.Streak)
	}

	// Simulate combo increment (this would happen in collision system)
	comboObj.Increment()
	comboObj.Timer = 0

	// Update system
	cs.Update(0.5, entities, eventDispatcher)

	// Verify combo state
	if comboObj.Streak != 1 {
		t.Errorf("Expected combo streak to be 1, got %d", comboObj.Streak)
	}

	// Simulate another combo increment
	comboObj.Increment()
	comboObj.Timer = 0

	// Update system
	cs.Update(0.5, entities, eventDispatcher)

	// Verify combo state
	if comboObj.Streak != 2 {
		t.Errorf("Expected combo streak to be 2, got %d", comboObj.Streak)
	}

	// Wait for combo to expire
	cs.Update(4.0, entities, eventDispatcher)

	// Verify combo was reset
	if comboObj.Streak != 0 {
		t.Errorf("Expected combo streak to be reset to 0, got %d", comboObj.Streak)
	}

	// Verify combo component was updated
	if comboObj.Timer != 0 {
		t.Errorf("Expected combo timer to be reset to 0, got %f", comboObj.Timer)
	}
}

// Helper function to create test entities
func createComboEntity(id uint64) Entity {
	entity := entities.NewEntity(id)
	combo := components.NewCombo()
	entity.AddComponent(combo)
	return entity
}
