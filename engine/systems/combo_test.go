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

	// Test initial values
	if cs.currentCombo != 0 {
		t.Errorf("Expected initial currentCombo to be 0, got %d", cs.currentCombo)
	}

	if cs.comboTimer != 0.0 {
		t.Errorf("Expected initial comboTimer to be 0.0, got %f", cs.comboTimer)
	}

	if cs.comboTimeout != 3.0 {
		t.Errorf("Expected comboTimeout to be 3.0, got %f", cs.comboTimeout)
	}

	if cs.lastComboTime != 0.0 {
		t.Errorf("Expected initial lastComboTime to be 0.0, got %f", cs.lastComboTime)
	}
}

func TestComboSystem_Update_NoEntities(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test with no entities
	entities := []Entity{}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify timer increased
	if cs.comboTimer != 0.016 {
		t.Errorf("Expected comboTimer to be 0.016, got %f", cs.comboTimer)
	}

	// Verify combo remains 0
	if cs.currentCombo != 0 {
		t.Errorf("Expected currentCombo to remain 0, got %d", cs.currentCombo)
	}
}

func TestComboSystem_Update_WithEntities(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity with combo component
	entity := createComboEntity(1)
	entities := []Entity{entity}

	// Set some combo state
	cs.currentCombo = 5
	cs.comboTimer = 10.0
	cs.lastComboTime = 8.0

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify timer increased
	if cs.comboTimer != 10.016 {
		t.Errorf("Expected comboTimer to be 10.016, got %f", cs.comboTimer)
	}

	// Verify combo component was updated
	comboComp := entity.GetCombo()
	if comboComp == nil {
		t.Fatal("Expected combo component to exist")
	}

	comboObj := comboComp.(*components.Combo)
	if comboObj.Streak != 5 {
		t.Errorf("Expected combo streak to be 5, got %d", comboObj.Streak)
	}

	if comboObj.Timer != 2.016 { // 10.016 - 8.0
		t.Errorf("Expected combo timer to be 2.016, got %f", comboObj.Timer)
	}
}

func TestComboSystem_Update_ComboExpiration(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Set up combo that should expire
	cs.currentCombo = 5
	cs.comboTimer = 10.0
	cs.lastComboTime = 6.0 // 4 seconds ago, should expire (timeout is 3.0)

	// Create entity with combo component
	entity := createComboEntity(1)
	entities := []Entity{entity}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify combo was reset
	if cs.currentCombo != 0 {
		t.Errorf("Expected currentCombo to be reset to 0, got %d", cs.currentCombo)
	}

	// Verify combo component was updated
	comboComp := entity.GetCombo()
	comboObj := comboComp.(*components.Combo)
	if comboObj.Streak != 0 {
		t.Errorf("Expected combo streak to be reset to 0, got %d", comboObj.Streak)
	}
}

func TestComboSystem_Update_ComboNotExpired(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Set up combo that should not expire
	cs.currentCombo = 3
	cs.comboTimer = 10.0
	cs.lastComboTime = 8.0 // 2 seconds ago, should not expire (timeout is 3.0)

	// Create entity with combo component
	entity := createComboEntity(1)
	entities := []Entity{entity}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify combo was not reset
	if cs.currentCombo != 3 {
		t.Errorf("Expected currentCombo to remain 3, got %d", cs.currentCombo)
	}

	// Verify combo component was updated
	comboComp := entity.GetCombo()
	comboObj := comboComp.(*components.Combo)
	if comboObj.Streak != 3 {
		t.Errorf("Expected combo streak to remain 3, got %d", comboObj.Streak)
	}
}

func TestComboSystem_Update_ComboExpirationSingleCombo(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Set up single combo that expires (should not print message)
	cs.currentCombo = 1
	cs.comboTimer = 10.0
	cs.lastComboTime = 6.0 // 4 seconds ago, should expire

	// Create entity with combo component
	entity := createComboEntity(1)
	entities := []Entity{entity}

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify combo was reset
	if cs.currentCombo != 0 {
		t.Errorf("Expected currentCombo to be reset to 0, got %d", cs.currentCombo)
	}
}

func TestComboSystem_Initialize(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	cs.Initialize(eventDispatcher)

	// Verify initial state
	if cs.currentCombo != 0 {
		t.Errorf("Expected initial currentCombo to be 0, got %d", cs.currentCombo)
	}

	if cs.lastComboTime != 0.0 {
		t.Errorf("Expected initial lastComboTime to be 0.0, got %f", cs.lastComboTime)
	}
}

func TestComboSystem_EventHandling_FirstPacket(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	cs.Initialize(eventDispatcher)

	// Set up some timer state
	cs.comboTimer = 10.0

	// Publish packet caught event
	event := events.NewEvent(events.EventPacketCaught, nil)
	eventDispatcher.Publish(event)

	// Verify combo was incremented
	if cs.currentCombo != 1 {
		t.Errorf("Expected currentCombo to be 1, got %d", cs.currentCombo)
	}

	// Verify last combo time was updated
	if cs.lastComboTime != 10.0 {
		t.Errorf("Expected lastComboTime to be 10.0, got %f", cs.lastComboTime)
	}
}

func TestComboSystem_EventHandling_MultiplePackets(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	cs.Initialize(eventDispatcher)

	// Set up some timer state
	cs.comboTimer = 10.0

	// Publish multiple packet caught events
	for i := 0; i < 3; i++ {
		event := events.NewEvent(events.EventPacketCaught, nil)
		eventDispatcher.Publish(event)
		cs.comboTimer += 0.5 // Simulate time passing
	}

	// Verify combo was incremented correctly
	if cs.currentCombo != 3 {
		t.Errorf("Expected currentCombo to be 3, got %d", cs.currentCombo)
	}

	// Verify last combo time was updated to the last event
	if cs.lastComboTime != 11.0 {
		t.Errorf("Expected lastComboTime to be 11.0, got %f", cs.lastComboTime)
	}
}

func TestComboSystem_calculateComboBonus_NoBonus(t *testing.T) {
	cs := NewComboSystem()

	// Test combos that don't qualify for bonus
	testCases := []int{0, 1, 2}

	for _, combo := range testCases {
		cs.currentCombo = combo
		bonus := cs.calculateComboBonus()
		if bonus != 0 {
			t.Errorf("Expected bonus for combo %d to be 0, got %d", combo, bonus)
		}
	}
}

func TestComboSystem_calculateComboBonus_3xCombo(t *testing.T) {
	cs := NewComboSystem()

	// Test 3x combo
	cs.currentCombo = 3
	bonus := cs.calculateComboBonus()
	if bonus != 10 {
		t.Errorf("Expected bonus for 3x combo to be 10, got %d", bonus)
	}
}

func TestComboSystem_calculateComboBonus_5xCombo(t *testing.T) {
	cs := NewComboSystem()

	// Test 5x combo
	cs.currentCombo = 5
	bonus := cs.calculateComboBonus()
	if bonus != 20 {
		t.Errorf("Expected bonus for 5x combo to be 20, got %d", bonus)
	}
}

func TestComboSystem_calculateComboBonus_7xCombo(t *testing.T) {
	cs := NewComboSystem()

	// Test 7x combo
	cs.currentCombo = 7
	bonus := cs.calculateComboBonus()
	if bonus != 30 {
		t.Errorf("Expected bonus for 7x combo to be 30, got %d", bonus)
	}
}

func TestComboSystem_calculateComboBonus_10xCombo(t *testing.T) {
	cs := NewComboSystem()

	// Test 10x combo
	cs.currentCombo = 10
	bonus := cs.calculateComboBonus()
	if bonus != 50 {
		t.Errorf("Expected bonus for 10x combo to be 50, got %d", bonus)
	}
}

func TestComboSystem_calculateComboBonus_HigherCombo(t *testing.T) {
	cs := NewComboSystem()

	// Test combo higher than 10x (should still give 50 bonus)
	cs.currentCombo = 15
	bonus := cs.calculateComboBonus()
	if bonus != 50 {
		t.Errorf("Expected bonus for 15x combo to be 50, got %d", bonus)
	}
}

func TestComboSystem_GetCurrentCombo(t *testing.T) {
	cs := NewComboSystem()

	// Set combo
	cs.currentCombo = 7

	// Get combo
	result := cs.GetCurrentCombo()

	// Verify result
	if result != 7 {
		t.Errorf("Expected GetCurrentCombo to return 7, got %d", result)
	}
}

func TestComboSystem_GetComboTimer(t *testing.T) {
	cs := NewComboSystem()

	// Set up timer state
	cs.comboTimer = 15.0
	cs.lastComboTime = 10.0

	// Get timer
	result := cs.GetComboTimer()

	// Verify result (15.0 - 10.0 = 5.0)
	if result != 5.0 {
		t.Errorf("Expected GetComboTimer to return 5.0, got %f", result)
	}
}

func TestComboSystem_GetComboTimer_Zero(t *testing.T) {
	cs := NewComboSystem()

	// Set up timer state (no time passed since last combo)
	cs.comboTimer = 10.0
	cs.lastComboTime = 10.0

	// Get timer
	result := cs.GetComboTimer()

	// Verify result (10.0 - 10.0 = 0.0)
	if result != 0.0 {
		t.Errorf("Expected GetComboTimer to return 0.0, got %f", result)
	}
}

func TestComboSystem_Integration(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	cs.Initialize(eventDispatcher)

	// Create entity with combo component
	entity := createComboEntity(1)
	entities := []Entity{entity}

	// Simulate game loop with packet catches
	cs.comboTimer = 0.0

	// First packet caught
	event1 := events.NewEvent(events.EventPacketCaught, nil)
	eventDispatcher.Publish(event1)
	cs.comboTimer += 0.5

	// Update system
	cs.Update(0.5, entities, eventDispatcher)

	// Verify combo state
	if cs.currentCombo != 1 {
		t.Errorf("Expected currentCombo to be 1, got %d", cs.currentCombo)
	}

	// Second packet caught (within timeout)
	event2 := events.NewEvent(events.EventPacketCaught, nil)
	eventDispatcher.Publish(event2)
	cs.comboTimer += 0.5

	// Update system
	cs.Update(0.5, entities, eventDispatcher)

	// Verify combo state
	if cs.currentCombo != 2 {
		t.Errorf("Expected currentCombo to be 2, got %d", cs.currentCombo)
	}

	// Wait for combo to expire
	cs.comboTimer += 4.0 // More than 3.0 second timeout

	// Update system
	cs.Update(4.0, entities, eventDispatcher)

	// Verify combo was reset
	if cs.currentCombo != 0 {
		t.Errorf("Expected currentCombo to be reset to 0, got %d", cs.currentCombo)
	}

	// Verify combo component was updated
	comboComp := entity.GetCombo()
	comboObj := comboComp.(*components.Combo)
	if comboObj.Streak != 0 {
		t.Errorf("Expected combo streak to be reset to 0, got %d", comboObj.Streak)
	}
}

func TestComboSystem_EntityWithoutComboComponent(t *testing.T) {
	cs := NewComboSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Create entity without combo component
	entity := entities.NewEntity(1)
	entities := []Entity{entity}

	// Set some combo state
	cs.currentCombo = 5
	cs.comboTimer = 10.0
	cs.lastComboTime = 8.0

	// Run update
	cs.Update(0.016, entities, eventDispatcher)

	// Verify timer increased
	if cs.comboTimer != 10.016 {
		t.Errorf("Expected comboTimer to be 10.016, got %f", cs.comboTimer)
	}

	// Verify combo remains unchanged (no combo component to update)
	if cs.currentCombo != 5 {
		t.Errorf("Expected currentCombo to remain 5, got %d", cs.currentCombo)
	}
}

// Helper function to create test entities

func createComboEntity(id uint64) Entity {
	entity := entities.NewEntity(id)
	combo := components.NewCombo()
	entity.AddComponent(combo)
	return entity
}
