package events

import (
	"testing"
	"time"
)

// TestEventTypeConstants tests that all event type constants are defined
func TestEventTypeConstants(t *testing.T) {
	expectedTypes := []EventType{
		EventPacketCaught,
		EventPacketLost,
		EventPowerUpCollected,
		EventPowerUpActivated,
		EventGameOver,
		EventGameStart,
		EventSLAUpdated,
		EventLevelUp,
		EventDDoSStart,
		EventDDoSEnd,
	}

	for _, eventType := range expectedTypes {
		if eventType == "" {
			t.Errorf("Event type constant should not be empty")
		}
	}
}

// TestNewEvent tests the NewEvent function
func TestNewEvent(t *testing.T) {
	t.Run("With Data", func(t *testing.T) {
		score := 100
		level := 5
		data := &EventData{
			Score: &score,
			Level: &level,
		}

		event := NewEvent(EventPacketCaught, data)

		if event.Type != EventPacketCaught {
			t.Errorf("Expected event type %s, got %s", EventPacketCaught, event.Type)
		}

		if event.Data != data {
			t.Errorf("Expected event data to match provided data")
		}

		if event.Timestamp.IsZero() {
			t.Error("Expected timestamp to be set")
		}

		// Check that timestamp is recent (within 1 second)
		if time.Since(event.Timestamp) > time.Second {
			t.Error("Expected timestamp to be recent")
		}
	})

	t.Run("Without Data", func(t *testing.T) {
		event := NewEvent(EventGameStart, nil)

		if event.Type != EventGameStart {
			t.Errorf("Expected event type %s, got %s", EventGameStart, event.Type)
		}

		if event.Data != nil {
			t.Error("Expected event data to be nil")
		}

		if event.Timestamp.IsZero() {
			t.Error("Expected timestamp to be set")
		}
	})
}

// TestNewEventWithMap tests the NewEventWithMap function
func TestNewEventWithMap(t *testing.T) {
	t.Run("Complete Data Map", func(t *testing.T) {
		data := map[string]interface{}{
			"score":        100,
			"packet":       "test_packet",
			"powerup":      "shield",
			"level":        3,
			"time":         45.5,
			"current":      95.2,
			"target":       99.5,
			"caught":       25,
			"lost":         2,
			"remaining":    8,
			"budget":       10,
			"combo_count":  5,
			"bonus_points": 50,
			"duration":     30.0,
			"mode":         1,
			"sla":          98.5,
			"errors":       3,
		}

		event := NewEventWithMap(EventSLAUpdated, data)

		if event.Type != EventSLAUpdated {
			t.Errorf("Expected event type %s, got %s", EventSLAUpdated, event.Type)
		}

		if event.Data == nil {
			t.Fatal("Expected event data to be created")
		}

		// Check all fields
		if *event.Data.Score != 100 {
			t.Errorf("Expected score 100, got %d", *event.Data.Score)
		}
		if event.Data.Packet != "test_packet" {
			t.Errorf("Expected packet 'test_packet', got %v", event.Data.Packet)
		}
		if *event.Data.Powerup != "shield" {
			t.Errorf("Expected powerup 'shield', got %s", *event.Data.Powerup)
		}
		if *event.Data.Level != 3 {
			t.Errorf("Expected level 3, got %d", *event.Data.Level)
		}
		if *event.Data.Time != 45.5 {
			t.Errorf("Expected time 45.5, got %f", *event.Data.Time)
		}
		if *event.Data.Current != 95.2 {
			t.Errorf("Expected current 95.2, got %f", *event.Data.Current)
		}
		if *event.Data.Target != 99.5 {
			t.Errorf("Expected target 99.5, got %f", *event.Data.Target)
		}
		if *event.Data.Caught != 25 {
			t.Errorf("Expected caught 25, got %d", *event.Data.Caught)
		}
		if *event.Data.Lost != 2 {
			t.Errorf("Expected lost 2, got %d", *event.Data.Lost)
		}
		if *event.Data.Remaining != 8 {
			t.Errorf("Expected remaining 8, got %d", *event.Data.Remaining)
		}
		if *event.Data.Budget != 10 {
			t.Errorf("Expected budget 10, got %d", *event.Data.Budget)
		}
		if *event.Data.ComboCount != 5 {
			t.Errorf("Expected combo count 5, got %d", *event.Data.ComboCount)
		}
		if *event.Data.BonusPoints != 50 {
			t.Errorf("Expected bonus points 50, got %d", *event.Data.BonusPoints)
		}
		if *event.Data.Duration != 30.0 {
			t.Errorf("Expected duration 30.0, got %f", *event.Data.Duration)
		}
		if *event.Data.Mode != 1 {
			t.Errorf("Expected mode 1, got %d", *event.Data.Mode)
		}
		if *event.Data.SLA != 98.5 {
			t.Errorf("Expected SLA 98.5, got %f", *event.Data.SLA)
		}
		if *event.Data.Errors != 3 {
			t.Errorf("Expected errors 3, got %d", *event.Data.Errors)
		}
	})

	t.Run("Partial Data Map", func(t *testing.T) {
		data := map[string]interface{}{
			"score": 150,
			"level": 7,
		}

		event := NewEventWithMap(EventLevelUp, data)

		if event.Type != EventLevelUp {
			t.Errorf("Expected event type %s, got %s", EventLevelUp, event.Type)
		}

		if event.Data == nil {
			t.Fatal("Expected event data to be created")
		}

		// Check set fields
		if *event.Data.Score != 150 {
			t.Errorf("Expected score 150, got %d", *event.Data.Score)
		}
		if *event.Data.Level != 7 {
			t.Errorf("Expected level 7, got %d", *event.Data.Level)
		}

		// Check unset fields are nil
		if event.Data.Powerup != nil {
			t.Error("Expected powerup to be nil")
		}
		if event.Data.Time != nil {
			t.Error("Expected time to be nil")
		}
		if event.Data.Current != nil {
			t.Error("Expected current to be nil")
		}
	})

	t.Run("Empty Data Map", func(t *testing.T) {
		data := map[string]interface{}{}

		event := NewEventWithMap(EventGameStart, data)

		if event.Type != EventGameStart {
			t.Errorf("Expected event type %s, got %s", EventGameStart, event.Type)
		}

		if event.Data == nil {
			t.Fatal("Expected event data to be created")
		}

		// All fields should be nil
		if event.Data.Score != nil {
			t.Error("Expected score to be nil")
		}
		if event.Data.Level != nil {
			t.Error("Expected level to be nil")
		}
		if event.Data.Powerup != nil {
			t.Error("Expected powerup to be nil")
		}
	})

	t.Run("Invalid Type Conversions", func(t *testing.T) {
		data := map[string]interface{}{
			"score": "not_an_int",
			"level": "not_an_int",
			"time":  "not_a_float",
		}

		event := NewEventWithMap(EventPacketCaught, data)

		if event.Type != EventPacketCaught {
			t.Errorf("Expected event type %s, got %s", EventPacketCaught, event.Type)
		}

		if event.Data == nil {
			t.Fatal("Expected event data to be created")
		}

		// Invalid conversions should result in nil values
		if event.Data.Score != nil {
			t.Error("Expected score to be nil for invalid type")
		}
		if event.Data.Level != nil {
			t.Error("Expected level to be nil for invalid type")
		}
		if event.Data.Time != nil {
			t.Error("Expected time to be nil for invalid type")
		}
	})
}

// TestNewEventDispatcher tests the NewEventDispatcher function
func TestNewEventDispatcher(t *testing.T) {
	dispatcher := NewEventDispatcher()

	if dispatcher == nil {
		t.Fatal("Expected dispatcher to be created")
	}

	if dispatcher.handlers == nil {
		t.Error("Expected handlers map to be initialized")
	}

	if len(dispatcher.handlers) != 0 {
		t.Error("Expected handlers map to be empty initially")
	}
}

// TestEventDispatcher_Subscribe tests the Subscribe method
func TestEventDispatcher_Subscribe(t *testing.T) {
	dispatcher := NewEventDispatcher()

	// Test subscribing a handler
	handler := func(event *Event) {
		// Empty handler for testing
	}

	dispatcher.Subscribe(EventPacketCaught, handler)

	if len(dispatcher.handlers[EventPacketCaught]) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(dispatcher.handlers[EventPacketCaught]))
	}

	// Test subscribing multiple handlers for the same event type
	handler2 := func(event *Event) {
		// Empty handler for testing
	}

	dispatcher.Subscribe(EventPacketCaught, handler2)

	if len(dispatcher.handlers[EventPacketCaught]) != 2 {
		t.Errorf("Expected 2 handlers, got %d", len(dispatcher.handlers[EventPacketCaught]))
	}

	// Test subscribing handlers for different event types
	dispatcher.Subscribe(EventGameStart, handler)

	if len(dispatcher.handlers[EventGameStart]) != 1 {
		t.Errorf("Expected 1 handler for GameStart, got %d", len(dispatcher.handlers[EventGameStart]))
	}

	// Verify handlers are stored correctly
	if len(dispatcher.handlers[EventPacketCaught]) != 2 {
		t.Errorf("Expected 2 handlers for PacketCaught, got %d", len(dispatcher.handlers[EventPacketCaught]))
	}
}

// TestEventDispatcher_Publish tests the Publish method
func TestEventDispatcher_Publish(t *testing.T) {
	t.Run("Single Handler", func(t *testing.T) {
		dispatcher := NewEventDispatcher()
		handlerCalled := false
		receivedEvent := (*Event)(nil)

		handler := func(event *Event) {
			handlerCalled = true
			receivedEvent = event
		}

		dispatcher.Subscribe(EventPacketCaught, handler)

		event := NewEvent(EventPacketCaught, nil)
		dispatcher.Publish(event)

		if !handlerCalled {
			t.Error("Expected handler to be called")
		}

		if receivedEvent != event {
			t.Error("Expected handler to receive the published event")
		}
	})

	t.Run("Multiple Handlers", func(t *testing.T) {
		dispatcher := NewEventDispatcher()
		handler1Called := false
		handler2Called := false

		handler1 := func(event *Event) {
			handler1Called = true
		}

		handler2 := func(event *Event) {
			handler2Called = true
		}

		dispatcher.Subscribe(EventPacketCaught, handler1)
		dispatcher.Subscribe(EventPacketCaught, handler2)

		event := NewEvent(EventPacketCaught, nil)
		dispatcher.Publish(event)

		if !handler1Called {
			t.Error("Expected handler1 to be called")
		}

		if !handler2Called {
			t.Error("Expected handler2 to be called")
		}
	})

	t.Run("No Handlers", func(t *testing.T) {
		dispatcher := NewEventDispatcher()

		// Should not panic when publishing to event type with no handlers
		event := NewEvent(EventPacketCaught, nil)
		dispatcher.Publish(event)
	})

	t.Run("Different Event Types", func(t *testing.T) {
		dispatcher := NewEventDispatcher()
		packetHandlerCalled := false
		gameHandlerCalled := false

		packetHandler := func(event *Event) {
			packetHandlerCalled = true
		}

		gameHandler := func(event *Event) {
			gameHandlerCalled = true
		}

		dispatcher.Subscribe(EventPacketCaught, packetHandler)
		dispatcher.Subscribe(EventGameStart, gameHandler)

		// Publish packet event
		packetEvent := NewEvent(EventPacketCaught, nil)
		dispatcher.Publish(packetEvent)

		if !packetHandlerCalled {
			t.Error("Expected packet handler to be called")
		}

		if gameHandlerCalled {
			t.Error("Expected game handler to not be called")
		}

		// Reset and publish game event
		packetHandlerCalled = false
		gameHandlerCalled = false

		gameEvent := NewEvent(EventGameStart, nil)
		dispatcher.Publish(gameEvent)

		if packetHandlerCalled {
			t.Error("Expected packet handler to not be called")
		}

		if !gameHandlerCalled {
			t.Error("Expected game handler to be called")
		}
	})

	t.Run("Handler with Event Data", func(t *testing.T) {
		dispatcher := NewEventDispatcher()
		receivedData := (*EventData)(nil)

		handler := func(event *Event) {
			receivedData = event.Data
		}

		dispatcher.Subscribe(EventSLAUpdated, handler)

		score := 100
		level := 5
		data := &EventData{
			Score: &score,
			Level: &level,
		}

		event := NewEvent(EventSLAUpdated, data)
		dispatcher.Publish(event)

		if receivedData != data {
			t.Error("Expected handler to receive the correct event data")
		}

		if *receivedData.Score != 100 {
			t.Errorf("Expected score 100, got %d", *receivedData.Score)
		}

		if *receivedData.Level != 5 {
			t.Errorf("Expected level 5, got %d", *receivedData.Level)
		}
	})
}

// TestEventDispatcher_Integration tests integration scenarios
func TestEventDispatcher_Integration(t *testing.T) {
	t.Run("Multiple Event Types and Handlers", func(t *testing.T) {
		dispatcher := NewEventDispatcher()
		eventsReceived := make(map[EventType]int)

		// Create handlers for multiple event types
		handler := func(event *Event) {
			eventsReceived[event.Type]++
		}

		dispatcher.Subscribe(EventPacketCaught, handler)
		dispatcher.Subscribe(EventPacketLost, handler)
		dispatcher.Subscribe(EventGameStart, handler)
		dispatcher.Subscribe(EventGameOver, handler)

		// Publish various events
		dispatcher.Publish(NewEvent(EventPacketCaught, nil))
		dispatcher.Publish(NewEvent(EventPacketCaught, nil))
		dispatcher.Publish(NewEvent(EventPacketLost, nil))
		dispatcher.Publish(NewEvent(EventGameStart, nil))
		dispatcher.Publish(NewEvent(EventGameOver, nil))
		dispatcher.Publish(NewEvent(EventPacketCaught, nil))

		// Verify counts
		expected := map[EventType]int{
			EventPacketCaught: 3,
			EventPacketLost:   1,
			EventGameStart:    1,
			EventGameOver:     1,
		}

		for eventType, expectedCount := range expected {
			if eventsReceived[eventType] != expectedCount {
				t.Errorf("Expected %d events of type %s, got %d", expectedCount, eventType, eventsReceived[eventType])
			}
		}
	})

	t.Run("Handler Order", func(t *testing.T) {
		dispatcher := NewEventDispatcher()
		callOrder := []int{}

		handler1 := func(event *Event) {
			callOrder = append(callOrder, 1)
		}

		handler2 := func(event *Event) {
			callOrder = append(callOrder, 2)
		}

		handler3 := func(event *Event) {
			callOrder = append(callOrder, 3)
		}

		// Subscribe in order
		dispatcher.Subscribe(EventPacketCaught, handler1)
		dispatcher.Subscribe(EventPacketCaught, handler2)
		dispatcher.Subscribe(EventPacketCaught, handler3)

		event := NewEvent(EventPacketCaught, nil)
		dispatcher.Publish(event)

		expectedOrder := []int{1, 2, 3}
		if len(callOrder) != len(expectedOrder) {
			t.Errorf("Expected %d handler calls, got %d", len(expectedOrder), len(callOrder))
		}

		for i, expected := range expectedOrder {
			if callOrder[i] != expected {
				t.Errorf("Expected handler %d to be called at position %d, got %d", expected, i, callOrder[i])
			}
		}
	})
}

// TestEventData_EdgeCases tests edge cases for EventData
func TestEventData_EdgeCases(t *testing.T) {
	t.Run("Nil EventData", func(t *testing.T) {
		event := NewEvent(EventGameStart, nil)

		if event.Data != nil {
			t.Error("Expected event data to be nil")
		}
	})

	t.Run("Empty EventData", func(t *testing.T) {
		data := &EventData{}
		event := NewEvent(EventGameStart, data)

		if event.Data != data {
			t.Error("Expected event data to match provided data")
		}

		// All fields should be nil
		if data.Score != nil {
			t.Error("Expected score to be nil")
		}
		if data.Level != nil {
			t.Error("Expected level to be nil")
		}
		if data.Powerup != nil {
			t.Error("Expected powerup to be nil")
		}
	})

	t.Run("Zero Values in EventData", func(t *testing.T) {
		score := 0
		level := 0
		time := 0.0
		data := &EventData{
			Score: &score,
			Level: &level,
			Time:  &time,
		}

		event := NewEvent(EventLevelUp, data)

		if *event.Data.Score != 0 {
			t.Errorf("Expected score 0, got %d", *event.Data.Score)
		}
		if *event.Data.Level != 0 {
			t.Errorf("Expected level 0, got %d", *event.Data.Level)
		}
		if *event.Data.Time != 0.0 {
			t.Errorf("Expected time 0.0, got %f", *event.Data.Time)
		}
	})
}

// TestEvent_EdgeCases tests edge cases for Event
func TestEvent_EdgeCases(t *testing.T) {
	t.Run("Nil Event", func(t *testing.T) {
		dispatcher := NewEventDispatcher()
		handlerCalled := false

		handler := func(event *Event) {
			handlerCalled = true
		}

		dispatcher.Subscribe(EventPacketCaught, handler)

		// Should not panic when publishing nil event
		dispatcher.Publish(nil)

		if handlerCalled {
			t.Error("Expected handler to not be called with nil event")
		}
	})

	t.Run("Event with Zero Timestamp", func(t *testing.T) {
		event := &Event{
			Type:      EventGameStart,
			Timestamp: time.Time{},
			Data:      nil,
		}

		if !event.Timestamp.IsZero() {
			t.Error("Expected timestamp to be zero")
		}
	})
}

// Benchmark tests for performance
func BenchmarkNewEvent(b *testing.B) {
	score := 100
	level := 5
	data := &EventData{
		Score: &score,
		Level: &level,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewEvent(EventPacketCaught, data)
	}
}

func BenchmarkNewEventWithMap(b *testing.B) {
	data := map[string]interface{}{
		"score":        100,
		"level":        5,
		"time":         45.5,
		"current":      95.2,
		"target":       99.5,
		"caught":       25,
		"lost":         2,
		"remaining":    8,
		"budget":       10,
		"combo_count":  5,
		"bonus_points": 50,
		"duration":     30.0,
		"mode":         1,
		"sla":          98.5,
		"errors":       3,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewEventWithMap(EventSLAUpdated, data)
	}
}

func BenchmarkEventDispatcher_Publish(b *testing.B) {
	dispatcher := NewEventDispatcher()
	handler := func(event *Event) {
		// Empty handler
	}

	dispatcher.Subscribe(EventPacketCaught, handler)
	event := NewEvent(EventPacketCaught, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dispatcher.Publish(event)
	}
}
