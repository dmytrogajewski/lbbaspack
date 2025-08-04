package systems

import (
	"lbbaspack/engine/events"
	"testing"
)

func TestNewPowerUpSystem(t *testing.T) {
	pus := NewPowerUpSystem()

	// Test that the system is properly initialized
	if pus == nil {
		t.Fatal("NewPowerUpSystem returned nil")
	}

	// Test required components
	expectedComponents := []string{"PowerUpType"}
	if len(pus.RequiredComponents) != len(expectedComponents) {
		t.Errorf("Expected %d required components, got %d", len(expectedComponents), len(pus.RequiredComponents))
	}

	for i, component := range expectedComponents {
		if pus.RequiredComponents[i] != component {
			t.Errorf("Expected required component %s at index %d, got %s", component, i, pus.RequiredComponents[i])
		}
	}

	// Test activePowerUps map initialization
	if pus.activePowerUps == nil {
		t.Fatal("ActivePowerUps map should not be nil")
	}

	if len(pus.activePowerUps) != 0 {
		t.Errorf("Expected initial activePowerUps count to be 0, got %d", len(pus.activePowerUps))
	}
}

func TestPowerUpSystem_Update_NoActivePowerUps(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test with no active power-ups
	entities := []Entity{}

	// Run update
	pus.Update(0.016, entities, eventDispatcher)

	// Verify activePowerUps map remains empty
	if len(pus.activePowerUps) != 0 {
		t.Errorf("Expected activePowerUps count to remain 0, got %d", len(pus.activePowerUps))
	}
}

func TestPowerUpSystem_Update_WithActivePowerUps(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Add some active power-ups
	pus.activePowerUps["SpeedBoost"] = 15.0
	pus.activePowerUps["DoublePoints"] = 20.0

	entities := []Entity{}

	// Run update
	pus.Update(0.016, entities, eventDispatcher)

	// Verify power-ups were updated
	if pus.activePowerUps["SpeedBoost"] != 15.0-0.016 {
		t.Errorf("Expected SpeedBoost remaining time to be %f, got %f", 15.0-0.016, pus.activePowerUps["SpeedBoost"])
	}

	if pus.activePowerUps["DoublePoints"] != 20.0-0.016 {
		t.Errorf("Expected DoublePoints remaining time to be %f, got %f", 20.0-0.016, pus.activePowerUps["DoublePoints"])
	}

	// Verify both power-ups are still active
	if !pus.IsPowerUpActive("SpeedBoost") {
		t.Error("Expected SpeedBoost to still be active")
	}

	if !pus.IsPowerUpActive("DoublePoints") {
		t.Error("Expected DoublePoints to still be active")
	}
}

func TestPowerUpSystem_Update_PowerUpExpiration(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Add a power-up with very short remaining time
	pus.activePowerUps["SpeedBoost"] = 0.01

	entities := []Entity{}

	// Run update with delta time that exceeds remaining time
	pus.Update(0.02, entities, eventDispatcher)

	// Verify power-up was removed
	if len(pus.activePowerUps) != 0 {
		t.Errorf("Expected activePowerUps count to be 0 after expiration, got %d", len(pus.activePowerUps))
	}

	// Verify power-up is no longer active
	if pus.IsPowerUpActive("SpeedBoost") {
		t.Error("Expected SpeedBoost to no longer be active")
	}
}

func TestPowerUpSystem_Update_MixedPowerUps(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Add power-ups with different remaining times
	pus.activePowerUps["SpeedBoost"] = 15.0   // Long duration
	pus.activePowerUps["DoublePoints"] = 0.01 // Short duration
	pus.activePowerUps["SlowMotion"] = 12.0   // Medium duration

	entities := []Entity{}

	// Run update
	pus.Update(0.02, entities, eventDispatcher)

	// Verify only DoublePoints was removed (expired)
	if len(pus.activePowerUps) != 2 {
		t.Errorf("Expected 2 power-ups to remain, got %d", len(pus.activePowerUps))
	}

	// Verify SpeedBoost and SlowMotion are still active
	if !pus.IsPowerUpActive("SpeedBoost") {
		t.Error("Expected SpeedBoost to still be active")
	}
	if !pus.IsPowerUpActive("SlowMotion") {
		t.Error("Expected SlowMotion to still be active")
	}

	// Verify DoublePoints is no longer active
	if pus.IsPowerUpActive("DoublePoints") {
		t.Error("Expected DoublePoints to no longer be active")
	}
}

func TestPowerUpSystem_Update_ZeroDeltaTime(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Add a power-up
	pus.activePowerUps["SpeedBoost"] = 15.0

	entities := []Entity{}

	// Run update with zero delta time
	pus.Update(0.0, entities, eventDispatcher)

	// Verify power-up remaining time remains unchanged
	if pus.activePowerUps["SpeedBoost"] != 15.0 {
		t.Errorf("Expected SpeedBoost remaining time to remain 15.0, got %f", pus.activePowerUps["SpeedBoost"])
	}
}

func TestPowerUpSystem_Update_LargeDeltaTime(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Add a power-up
	pus.activePowerUps["SpeedBoost"] = 15.0

	entities := []Entity{}

	// Run update with large delta time
	pus.Update(20.0, entities, eventDispatcher)

	// Verify power-up was removed due to expiration
	if len(pus.activePowerUps) != 0 {
		t.Errorf("Expected activePowerUps count to be 0 after large delta time, got %d", len(pus.activePowerUps))
	}

	// Verify power-up is no longer active
	if pus.IsPowerUpActive("SpeedBoost") {
		t.Error("Expected SpeedBoost to no longer be active")
	}
}

func TestPowerUpSystem_Update_NegativeDeltaTime(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Add a power-up
	pus.activePowerUps["SpeedBoost"] = 15.0

	entities := []Entity{}

	// Run update with negative delta time
	pus.Update(-0.016, entities, eventDispatcher)

	// Verify power-up remaining time increases (negative delta time)
	if pus.activePowerUps["SpeedBoost"] != 15.0-(-0.016) {
		t.Errorf("Expected SpeedBoost remaining time to be %f, got %f", 15.0-(-0.016), pus.activePowerUps["SpeedBoost"])
	}
}

func TestPowerUpSystem_Initialize(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test that Initialize doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Initialize panicked: %v", r)
		}
	}()

	pus.Initialize(eventDispatcher)

	// If we get here, Initialize executed without panicking
}

func TestPowerUpSystem_activatePowerUp_SpeedBoost(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Activate SpeedBoost power-up
	pus.activatePowerUp("SpeedBoost", eventDispatcher)

	// Verify power-up was activated
	if !pus.IsPowerUpActive("SpeedBoost") {
		t.Error("Expected SpeedBoost to be active")
	}

	// Verify correct duration
	if pus.activePowerUps["SpeedBoost"] != 15.0 {
		t.Errorf("Expected SpeedBoost duration to be 15.0, got %f", pus.activePowerUps["SpeedBoost"])
	}
}

func TestPowerUpSystem_activatePowerUp_DoublePoints(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Activate DoublePoints power-up
	pus.activatePowerUp("DoublePoints", eventDispatcher)

	// Verify power-up was activated
	if !pus.IsPowerUpActive("DoublePoints") {
		t.Error("Expected DoublePoints to be active")
	}

	// Verify correct duration
	if pus.activePowerUps["DoublePoints"] != 20.0 {
		t.Errorf("Expected DoublePoints duration to be 20.0, got %f", pus.activePowerUps["DoublePoints"])
	}
}

func TestPowerUpSystem_activatePowerUp_SlowMotion(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Activate SlowMotion power-up
	pus.activatePowerUp("SlowMotion", eventDispatcher)

	// Verify power-up was activated
	if !pus.IsPowerUpActive("SlowMotion") {
		t.Error("Expected SlowMotion to be active")
	}

	// Verify correct duration
	if pus.activePowerUps["SlowMotion"] != 12.0 {
		t.Errorf("Expected SlowMotion duration to be 12.0, got %f", pus.activePowerUps["SlowMotion"])
	}
}

func TestPowerUpSystem_activatePowerUp_UnknownPowerUp(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Activate unknown power-up
	pus.activatePowerUp("UnknownPowerUp", eventDispatcher)

	// Verify power-up was activated with default duration
	if !pus.IsPowerUpActive("UnknownPowerUp") {
		t.Error("Expected UnknownPowerUp to be active")
	}

	// Verify default duration
	if pus.activePowerUps["UnknownPowerUp"] != 10.0 {
		t.Errorf("Expected UnknownPowerUp duration to be 10.0, got %f", pus.activePowerUps["UnknownPowerUp"])
	}
}

func TestPowerUpSystem_activatePowerUp_Reactivation(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Activate SpeedBoost power-up
	pus.activatePowerUp("SpeedBoost", eventDispatcher)

	// Verify initial activation
	if pus.activePowerUps["SpeedBoost"] != 15.0 {
		t.Errorf("Expected initial SpeedBoost duration to be 15.0, got %f", pus.activePowerUps["SpeedBoost"])
	}

	// Update to reduce remaining time
	pus.Update(5.0, []Entity{}, eventDispatcher)

	// Verify remaining time was reduced
	if pus.activePowerUps["SpeedBoost"] != 10.0 {
		t.Errorf("Expected SpeedBoost remaining time to be 10.0, got %f", pus.activePowerUps["SpeedBoost"])
	}

	// Reactivate the same power-up
	pus.activatePowerUp("SpeedBoost", eventDispatcher)

	// Verify duration was reset to full duration
	if pus.activePowerUps["SpeedBoost"] != 15.0 {
		t.Errorf("Expected SpeedBoost duration to be reset to 15.0, got %f", pus.activePowerUps["SpeedBoost"])
	}
}

func TestPowerUpSystem_IsPowerUpActive(t *testing.T) {
	pus := NewPowerUpSystem()

	// Test with no active power-ups
	if pus.IsPowerUpActive("SpeedBoost") {
		t.Error("Expected SpeedBoost to not be active")
	}

	// Add a power-up
	pus.activePowerUps["SpeedBoost"] = 15.0

	// Test with active power-up
	if !pus.IsPowerUpActive("SpeedBoost") {
		t.Error("Expected SpeedBoost to be active")
	}

	// Test with non-existent power-up
	if pus.IsPowerUpActive("NonExistent") {
		t.Error("Expected NonExistent to not be active")
	}
}

func TestPowerUpSystem_GetActivePowerUps(t *testing.T) {
	pus := NewPowerUpSystem()

	// Test with no active power-ups
	activePowerUps := pus.GetActivePowerUps()
	if len(activePowerUps) != 0 {
		t.Errorf("Expected 0 active power-ups, got %d", len(activePowerUps))
	}

	// Add some power-ups
	pus.activePowerUps["SpeedBoost"] = 15.0
	pus.activePowerUps["DoublePoints"] = 20.0

	// Test with active power-ups
	activePowerUps = pus.GetActivePowerUps()
	if len(activePowerUps) != 2 {
		t.Errorf("Expected 2 active power-ups, got %d", len(activePowerUps))
	}

	// Verify specific power-ups
	if activePowerUps["SpeedBoost"] != 15.0 {
		t.Errorf("Expected SpeedBoost remaining time to be 15.0, got %f", activePowerUps["SpeedBoost"])
	}

	if activePowerUps["DoublePoints"] != 20.0 {
		t.Errorf("Expected DoublePoints remaining time to be 20.0, got %f", activePowerUps["DoublePoints"])
	}
}

func TestPowerUpSystem_EventHandling_PowerUpCollected(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	pus.Initialize(eventDispatcher)

	// Publish power-up collected event
	powerUpName := "SpeedBoost"
	eventData := &events.EventData{
		Powerup: &powerUpName,
	}
	event := events.NewEvent(events.EventPowerUpCollected, eventData)
	eventDispatcher.Publish(event)

	// Verify power-up was activated
	if !pus.IsPowerUpActive("SpeedBoost") {
		t.Error("Expected SpeedBoost to be active after event")
	}

	// Verify correct duration
	if pus.activePowerUps["SpeedBoost"] != 15.0 {
		t.Errorf("Expected SpeedBoost duration to be 15.0, got %f", pus.activePowerUps["SpeedBoost"])
	}
}

func TestPowerUpSystem_EventHandling_PowerUpCollected_NilPowerUp(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	pus.Initialize(eventDispatcher)

	initialCount := len(pus.activePowerUps)

	// Publish power-up collected event with nil powerup
	eventData := &events.EventData{
		Powerup: nil,
	}
	event := events.NewEvent(events.EventPowerUpCollected, eventData)
	eventDispatcher.Publish(event)

	// Verify no power-up was activated
	if len(pus.activePowerUps) != initialCount {
		t.Errorf("Expected activePowerUps count to remain %d, got %d", initialCount, len(pus.activePowerUps))
	}
}

func TestPowerUpSystem_EventHandling_MultiplePowerUps(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	pus.Initialize(eventDispatcher)

	// Publish multiple power-up collected events
	powerUp1 := "SpeedBoost"
	powerUp2 := "DoublePoints"
	powerUp3 := "SlowMotion"

	event1 := events.NewEvent(events.EventPowerUpCollected, &events.EventData{Powerup: &powerUp1})
	event2 := events.NewEvent(events.EventPowerUpCollected, &events.EventData{Powerup: &powerUp2})
	event3 := events.NewEvent(events.EventPowerUpCollected, &events.EventData{Powerup: &powerUp3})

	eventDispatcher.Publish(event1)
	eventDispatcher.Publish(event2)
	eventDispatcher.Publish(event3)

	// Verify all power-ups were activated
	if !pus.IsPowerUpActive("SpeedBoost") {
		t.Error("Expected SpeedBoost to be active")
	}
	if !pus.IsPowerUpActive("DoublePoints") {
		t.Error("Expected DoublePoints to be active")
	}
	if !pus.IsPowerUpActive("SlowMotion") {
		t.Error("Expected SlowMotion to be active")
	}

	// Verify correct durations
	if pus.activePowerUps["SpeedBoost"] != 15.0 {
		t.Errorf("Expected SpeedBoost duration to be 15.0, got %f", pus.activePowerUps["SpeedBoost"])
	}
	if pus.activePowerUps["DoublePoints"] != 20.0 {
		t.Errorf("Expected DoublePoints duration to be 20.0, got %f", pus.activePowerUps["DoublePoints"])
	}
	if pus.activePowerUps["SlowMotion"] != 12.0 {
		t.Errorf("Expected SlowMotion duration to be 12.0, got %f", pus.activePowerUps["SlowMotion"])
	}
}

func TestPowerUpSystem_Integration(t *testing.T) {
	pus := NewPowerUpSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	pus.Initialize(eventDispatcher)

	// Activate multiple power-ups
	pus.activatePowerUp("SpeedBoost", eventDispatcher)
	pus.activatePowerUp("DoublePoints", eventDispatcher)

	// Verify initial state
	if len(pus.activePowerUps) != 2 {
		t.Errorf("Expected 2 active power-ups, got %d", len(pus.activePowerUps))
	}

	// Update power-ups
	pus.Update(5.0, []Entity{}, eventDispatcher)

	// Verify remaining times
	if pus.activePowerUps["SpeedBoost"] != 10.0 {
		t.Errorf("Expected SpeedBoost remaining time to be 10.0, got %f", pus.activePowerUps["SpeedBoost"])
	}
	if pus.activePowerUps["DoublePoints"] != 15.0 {
		t.Errorf("Expected DoublePoints remaining time to be 15.0, got %f", pus.activePowerUps["DoublePoints"])
	}

	// Activate another power-up
	pus.activatePowerUp("SlowMotion", eventDispatcher)

	// Verify new power-up was added
	if len(pus.activePowerUps) != 3 {
		t.Errorf("Expected 3 active power-ups, got %d", len(pus.activePowerUps))
	}

	// Update until one power-up expires
	pus.Update(10.0, []Entity{}, eventDispatcher)

	// Verify SpeedBoost expired
	if len(pus.activePowerUps) != 2 {
		t.Errorf("Expected 2 active power-ups after expiration, got %d", len(pus.activePowerUps))
	}

	if pus.IsPowerUpActive("SpeedBoost") {
		t.Error("Expected SpeedBoost to have expired")
	}

	// Verify remaining power-ups are still active
	if !pus.IsPowerUpActive("DoublePoints") {
		t.Error("Expected DoublePoints to still be active")
	}
	if !pus.IsPowerUpActive("SlowMotion") {
		t.Error("Expected SlowMotion to still be active")
	}

	// Get active power-ups
	activePowerUps := pus.GetActivePowerUps()
	if len(activePowerUps) != 2 {
		t.Errorf("Expected GetActivePowerUps to return 2 power-ups, got %d", len(activePowerUps))
	}
}
