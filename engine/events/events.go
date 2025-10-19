package events

import "time"

// EventType represents different types of events
type EventType string

const (
	EventPacketCaught      EventType = "packet_caught"
	EventPacketLost        EventType = "packet_lost"
	EventPowerUpCollected  EventType = "powerup_collected"
	EventPowerUpActivated  EventType = "powerup_activated"
	EventGameOver          EventType = "game_over"
	EventGameStart         EventType = "game_start"
	EventReturnToMenu      EventType = "return_to_menu"
	EventExit              EventType = "exit"
	EventSLAUpdated        EventType = "sla_updated"
	EventLevelUp           EventType = "level_up"
	EventDDoSStart         EventType = "ddos_start"
	EventDDoSEnd           EventType = "ddos_end"
	EventPacketDelivered   EventType = "packet_delivered"
	EventCollisionDetected EventType = "collision_detected"
	EventColliderOffscreen EventType = "collider_offscreen"
)

// EventData represents typed event data
type EventData struct {
	Score       *int
	Packet      interface{} // Entity type
	Powerup     *string
	Level       *int
	Time        *float64
	Current     *float64
	Target      *float64
	Caught      *int
	Lost        *int
	Remaining   *int
	Budget      *int
	ComboCount  *int
	BonusPoints *int
	Duration    *float64
	Mode        *int
	SLA         *float64
	Errors      *int
	BackendID   *int
	// Generalized collision context
	EntityA interface{}
	EntityB interface{}
	TagA    *string
	TagB    *string
	PosAX   *float64
	PosAY   *float64
	PosBX   *float64
	PosBY   *float64
}

// Event represents a game event
type Event struct {
	Type      EventType
	Timestamp time.Time
	Data      *EventData
}

// NewEvent creates a new event
func NewEvent(eventType EventType, data *EventData) *Event {
	return &Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
}

// NewEventWithMap creates a new event from a map (for backward compatibility)
func NewEventWithMap(eventType EventType, data map[string]interface{}) *Event {
	eventData := &EventData{}

	// Extract typed data from map
	if score, ok := data["score"].(int); ok {
		eventData.Score = &score
	}
	if packet, ok := data["packet"]; ok {
		eventData.Packet = packet
	}
	if powerup, ok := data["powerup"].(string); ok {
		eventData.Powerup = &powerup
	}
	if level, ok := data["level"].(int); ok {
		eventData.Level = &level
	}
	if time, ok := data["time"].(float64); ok {
		eventData.Time = &time
	}
	if current, ok := data["current"].(float64); ok {
		eventData.Current = &current
	}
	if target, ok := data["target"].(float64); ok {
		eventData.Target = &target
	}
	if caught, ok := data["caught"].(int); ok {
		eventData.Caught = &caught
	}
	if lost, ok := data["lost"].(int); ok {
		eventData.Lost = &lost
	}
	if remaining, ok := data["remaining"].(int); ok {
		eventData.Remaining = &remaining
	}
	if budget, ok := data["budget"].(int); ok {
		eventData.Budget = &budget
	}
	if comboCount, ok := data["combo_count"].(int); ok {
		eventData.ComboCount = &comboCount
	}
	if bonusPoints, ok := data["bonus_points"].(int); ok {
		eventData.BonusPoints = &bonusPoints
	}
	if duration, ok := data["duration"].(float64); ok {
		eventData.Duration = &duration
	}
	if mode, ok := data["mode"].(int); ok {
		eventData.Mode = &mode
	}
	if sla, ok := data["sla"].(float64); ok {
		eventData.SLA = &sla
	}
	if errors, ok := data["errors"].(int); ok {
		eventData.Errors = &errors
	}

	return NewEvent(eventType, eventData)
}

// EventHandler is a function that handles events
type EventHandler func(*Event)

// EventDispatcher manages event handling
type EventDispatcher struct {
	handlers map[EventType][]EventHandler
}

// NewEventDispatcher creates a new event dispatcher
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[EventType][]EventHandler),
	}
}

// Subscribe adds an event handler
func (ed *EventDispatcher) Subscribe(eventType EventType, handler EventHandler) {
	ed.handlers[eventType] = append(ed.handlers[eventType], handler)
}

// Publish sends an event to all subscribers
func (ed *EventDispatcher) Publish(event *Event) {
	if event == nil {
		return
	}
	if handlers, exists := ed.handlers[event.Type]; exists {
		for _, handler := range handlers {
			handler(event)
		}
	}
}
