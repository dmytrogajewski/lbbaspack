package systems

import (
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"
)

func TestNewRoutingSystem(t *testing.T) {
	rs := NewRoutingSystem()

	// Test that the system is properly initialized
	if rs == nil {
		t.Fatal("NewRoutingSystem returned nil")
	}

	// Test routes slice initialization
	if rs.routes == nil {
		t.Fatal("Routes slice should not be nil")
	}

	if len(rs.routes) != 0 {
		t.Errorf("Expected initial routes count to be 0, got %d", len(rs.routes))
	}
}

func TestRoutingSystem_Update_NoRoutes(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}

	// Test that Update doesn't panic with no routes
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with no routes: %v", r)
		}
	}()

	rs.Update(0.016, entities, eventDispatcher)

	// Verify routes slice remains empty
	if len(rs.routes) != 0 {
		t.Errorf("Expected routes count to remain 0, got %d", len(rs.routes))
	}
}

func TestRoutingSystem_Update_WithActiveRoutes(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}

	// Create some active routes
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})
	rs.CreateRoute(150, 250, 350, 450, color.RGBA{0, 255, 0, 255})

	initialRouteCount := len(rs.routes)

	// Test that Update doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with active routes: %v", r)
		}
	}()

	rs.Update(0.016, entities, eventDispatcher)

	// Verify routes are still active (not completed yet)
	if len(rs.routes) != initialRouteCount {
		t.Errorf("Expected routes count to remain %d, got %d", initialRouteCount, len(rs.routes))
	}

	// Verify routes are still active
	for _, route := range rs.routes {
		if !route.Active {
			t.Error("Expected route to still be active")
		}
	}
}

func TestRoutingSystem_Update_RouteCompletion(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}

	// Create a route
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})

	// Update until route completes (progress >= 1.0)
	// With speed 1.5, it should complete in about 0.67 seconds
	for i := 0; i < 70; i++ {
		rs.Update(0.01, entities, eventDispatcher)
	}

	// Verify route was removed after completion
	if len(rs.routes) != 0 {
		t.Errorf("Expected routes count to be 0 after completion, got %d", len(rs.routes))
	}
}

func TestRoutingSystem_Update_MixedRoutes(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}

	// Create routes
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})
	rs.CreateRoute(150, 250, 350, 450, color.RGBA{0, 255, 0, 255})

	initialRouteCount := len(rs.routes)

	// Update a few times
	for i := 0; i < 10; i++ {
		rs.Update(0.016, entities, eventDispatcher)
	}

	// Verify routes are still active (not completed yet)
	if len(rs.routes) != initialRouteCount {
		t.Errorf("Expected routes count to remain %d, got %d", initialRouteCount, len(rs.routes))
	}

	// Verify all routes are still active
	for _, route := range rs.routes {
		if !route.Active {
			t.Error("Expected route to still be active")
		}
	}
}

func TestRoutingSystem_Update_ZeroDeltaTime(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}

	// Create a route
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})

	initialProgress := rs.routes[0].Progress

	// Test that Update doesn't panic with zero delta time
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with zero delta time: %v", r)
		}
	}()

	rs.Update(0.0, entities, eventDispatcher)

	// Verify progress remains unchanged
	if rs.routes[0].Progress != initialProgress {
		t.Errorf("Expected progress to remain %f, got %f", initialProgress, rs.routes[0].Progress)
	}
}

func TestRoutingSystem_Update_LargeDeltaTime(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}

	// Create a route
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})

	// Test that Update doesn't panic with large delta time
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with large delta time: %v", r)
		}
	}()

	rs.Update(10.0, entities, eventDispatcher)

	// Verify route was completed and removed
	if len(rs.routes) != 0 {
		t.Errorf("Expected routes count to be 0 after large delta time, got %d", len(rs.routes))
	}
}

func TestRoutingSystem_Update_NegativeDeltaTime(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}

	// Create a route
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})

	initialProgress := rs.routes[0].Progress

	// Test that Update doesn't panic with negative delta time
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Update panicked with negative delta time: %v", r)
		}
	}()

	rs.Update(-0.016, entities, eventDispatcher)

	// Verify progress decreased (negative delta time)
	if rs.routes[0].Progress >= initialProgress {
		t.Errorf("Expected progress to decrease with negative delta time, got %f", rs.routes[0].Progress)
	}
}

func TestRoutingSystem_CreateRoute(t *testing.T) {
	rs := NewRoutingSystem()

	// Test route creation
	startX, startY := 100.0, 200.0
	endX, endY := 300.0, 400.0
	routeColor := color.RGBA{255, 0, 0, 255}

	rs.CreateRoute(startX, startY, endX, endY, routeColor)

	// Verify route was created
	if len(rs.routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(rs.routes))
	}

	route := rs.routes[0]
	if route.StartX != startX {
		t.Errorf("Expected StartX to be %f, got %f", startX, route.StartX)
	}
	if route.StartY != startY {
		t.Errorf("Expected StartY to be %f, got %f", startY, route.StartY)
	}
	if route.EndX != endX {
		t.Errorf("Expected EndX to be %f, got %f", endX, route.EndX)
	}
	if route.EndY != endY {
		t.Errorf("Expected EndY to be %f, got %f", endY, route.EndY)
	}
	if route.Progress != 0.0 {
		t.Errorf("Expected Progress to be 0.0, got %f", route.Progress)
	}
	if route.Speed != 1.5 {
		t.Errorf("Expected Speed to be 1.5, got %f", route.Speed)
	}
	if route.Color != routeColor {
		t.Errorf("Expected Color to be %v, got %v", routeColor, route.Color)
	}
	if !route.Active {
		t.Error("Expected route to be active")
	}
}

func TestRoutingSystem_CreateRoute_MultipleRoutes(t *testing.T) {
	rs := NewRoutingSystem()

	// Create multiple routes
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})
	rs.CreateRoute(150, 250, 350, 450, color.RGBA{0, 255, 0, 255})
	rs.CreateRoute(200, 300, 400, 500, color.RGBA{0, 0, 255, 255})

	// Verify all routes were created
	if len(rs.routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(rs.routes))
	}

	// Verify all routes are active
	for i, route := range rs.routes {
		if !route.Active {
			t.Errorf("Expected route %d to be active", i)
		}
	}
}

func TestRoutingSystem_CreateRoute_ZeroCoordinates(t *testing.T) {
	rs := NewRoutingSystem()

	// Test route creation with zero coordinates
	rs.CreateRoute(0, 0, 0, 0, color.RGBA{255, 0, 0, 255})

	// Verify route was created
	if len(rs.routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(rs.routes))
	}

	route := rs.routes[0]
	if route.StartX != 0 {
		t.Errorf("Expected StartX to be 0, got %f", route.StartX)
	}
	if route.StartY != 0 {
		t.Errorf("Expected StartY to be 0, got %f", route.StartY)
	}
	if route.EndX != 0 {
		t.Errorf("Expected EndX to be 0, got %f", route.EndX)
	}
	if route.EndY != 0 {
		t.Errorf("Expected EndY to be 0, got %f", route.EndY)
	}
}

func TestRoutingSystem_CreateRoute_NegativeCoordinates(t *testing.T) {
	rs := NewRoutingSystem()

	// Test route creation with negative coordinates
	rs.CreateRoute(-100, -200, -300, -400, color.RGBA{255, 0, 0, 255})

	// Verify route was created
	if len(rs.routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(rs.routes))
	}

	route := rs.routes[0]
	if route.StartX != -100 {
		t.Errorf("Expected StartX to be -100, got %f", route.StartX)
	}
	if route.StartY != -200 {
		t.Errorf("Expected StartY to be -200, got %f", route.StartY)
	}
	if route.EndX != -300 {
		t.Errorf("Expected EndX to be -300, got %f", route.EndX)
	}
	if route.EndY != -400 {
		t.Errorf("Expected EndY to be -400, got %f", route.EndY)
	}
}

func TestRoutingSystem_Draw_NilScreen(t *testing.T) {
	rs := NewRoutingSystem()

	// Create a route
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})

	// Test that Draw doesn't panic with nil screen
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked with nil screen: %v", r)
		}
	}()

	rs.Draw(nil)
}

func TestRoutingSystem_Draw_NoRoutes(t *testing.T) {
	rs := NewRoutingSystem()

	// Test that Draw doesn't panic with no routes
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked with no routes: %v", r)
		}
	}()

	rs.Draw(nil)
}

func TestRoutingSystem_Draw_WithRoutes(t *testing.T) {
	rs := NewRoutingSystem()

	// Create routes
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})
	rs.CreateRoute(150, 250, 350, 450, color.RGBA{0, 255, 0, 255})

	// Test that Draw doesn't panic with routes
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked with routes: %v", r)
		}
	}()

	rs.Draw(nil)
}

func TestRoutingSystem_Draw_CompletedRoutes(t *testing.T) {
	rs := NewRoutingSystem()

	// Create a route and complete it
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})
	rs.routes[0].Progress = 1.0
	rs.routes[0].Active = false

	// Test that Draw doesn't panic with completed routes
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Draw panicked with completed routes: %v", r)
		}
	}()

	rs.Draw(nil)
}

func TestRoutingSystem_Initialize(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Test that Initialize doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Initialize panicked: %v", r)
		}
	}()

	rs.Initialize(eventDispatcher)

	// If we get here, Initialize executed without panicking
}

func TestRoutingSystem_EventHandling_PacketCaught(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	rs.Initialize(eventDispatcher)

	// Create a packet entity
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(15, 15, color.RGBA{255, 0, 0, 255})
	entity.AddComponent(transform)
	entity.AddComponent(sprite)

	// Publish packet caught event
	eventData := &events.EventData{
		Packet: entity,
	}
	event := events.NewEvent(events.EventPacketCaught, eventData)
	eventDispatcher.Publish(event)

	// Verify route was created
	if len(rs.routes) != 1 {
		t.Errorf("Expected 1 route after packet caught event, got %d", len(rs.routes))
	}

	route := rs.routes[0]
	if !route.Active {
		t.Error("Expected created route to be active")
	}
	if route.Progress != 0.0 {
		t.Errorf("Expected route progress to be 0.0, got %f", route.Progress)
	}
}

func TestRoutingSystem_EventHandling_PacketCaught_NilPacket(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	rs.Initialize(eventDispatcher)

	initialRouteCount := len(rs.routes)

	// Publish packet caught event with nil packet
	eventData := &events.EventData{
		Packet: nil,
	}
	event := events.NewEvent(events.EventPacketCaught, eventData)
	eventDispatcher.Publish(event)

	// Verify no route was created
	if len(rs.routes) != initialRouteCount {
		t.Errorf("Expected routes count to remain %d, got %d", initialRouteCount, len(rs.routes))
	}
}

func TestRoutingSystem_EventHandling_PacketCaught_InvalidEntity(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	rs.Initialize(eventDispatcher)

	initialRouteCount := len(rs.routes)

	// Publish packet caught event with invalid entity (not Entity interface)
	eventData := &events.EventData{
		Packet: "not an entity",
	}
	event := events.NewEvent(events.EventPacketCaught, eventData)
	eventDispatcher.Publish(event)

	// Verify no route was created
	if len(rs.routes) != initialRouteCount {
		t.Errorf("Expected routes count to remain %d, got %d", initialRouteCount, len(rs.routes))
	}
}

func TestRoutingSystem_EventHandling_PacketCaught_MissingTransform(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	rs.Initialize(eventDispatcher)

	// Create entity without Transform component
	entity := entities.NewEntity(1)
	sprite := components.NewSprite(15, 15, color.RGBA{255, 0, 0, 255})
	entity.AddComponent(sprite)

	initialRouteCount := len(rs.routes)

	// Publish packet caught event
	eventData := &events.EventData{
		Packet: entity,
	}
	event := events.NewEvent(events.EventPacketCaught, eventData)
	eventDispatcher.Publish(event)

	// Verify no route was created
	if len(rs.routes) != initialRouteCount {
		t.Errorf("Expected routes count to remain %d, got %d", initialRouteCount, len(rs.routes))
	}
}

func TestRoutingSystem_EventHandling_PacketCaught_MissingSprite(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	rs.Initialize(eventDispatcher)

	// Create entity without Sprite component
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	entity.AddComponent(transform)

	initialRouteCount := len(rs.routes)

	// Publish packet caught event
	eventData := &events.EventData{
		Packet: entity,
	}
	event := events.NewEvent(events.EventPacketCaught, eventData)
	eventDispatcher.Publish(event)

	// Verify no route was created
	if len(rs.routes) != initialRouteCount {
		t.Errorf("Expected routes count to remain %d, got %d", initialRouteCount, len(rs.routes))
	}
}

func TestRoutingSystem_EventHandling_MultiplePackets(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	rs.Initialize(eventDispatcher)

	// Create multiple packet entities
	entity1 := entities.NewEntity(1)
	transform1 := components.NewTransform(100, 200)
	sprite1 := components.NewSprite(15, 15, color.RGBA{255, 0, 0, 255})
	entity1.AddComponent(transform1)
	entity1.AddComponent(sprite1)

	entity2 := entities.NewEntity(2)
	transform2 := components.NewTransform(300, 400)
	sprite2 := components.NewSprite(15, 15, color.RGBA{0, 255, 0, 255})
	entity2.AddComponent(transform2)
	entity2.AddComponent(sprite2)

	// Publish multiple packet caught events
	event1 := events.NewEvent(events.EventPacketCaught, &events.EventData{Packet: entity1})
	event2 := events.NewEvent(events.EventPacketCaught, &events.EventData{Packet: entity2})

	eventDispatcher.Publish(event1)
	eventDispatcher.Publish(event2)

	// Verify routes were created
	if len(rs.routes) != 2 {
		t.Errorf("Expected 2 routes after multiple packet caught events, got %d", len(rs.routes))
	}

	// Verify all routes are active
	for i, route := range rs.routes {
		if !route.Active {
			t.Errorf("Expected route %d to be active", i)
		}
	}
}

func TestRoutingSystem_Integration(t *testing.T) {
	rs := NewRoutingSystem()
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	rs.Initialize(eventDispatcher)

	// Create a packet entity
	entity := entities.NewEntity(1)
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(15, 15, color.RGBA{255, 0, 0, 255})
	entity.AddComponent(transform)
	entity.AddComponent(sprite)

	// Publish packet caught event
	eventData := &events.EventData{
		Packet: entity,
	}
	event := events.NewEvent(events.EventPacketCaught, eventData)
	eventDispatcher.Publish(event)

	// Verify route was created
	if len(rs.routes) != 1 {
		t.Errorf("Expected 1 route after packet caught event, got %d", len(rs.routes))
	}

	// Update the route
	entities := []Entity{}
	rs.Update(0.016, entities, eventDispatcher)

	// Verify route is still active
	if len(rs.routes) != 1 {
		t.Errorf("Expected 1 route after update, got %d", len(rs.routes))
	}

	route := rs.routes[0]
	if !route.Active {
		t.Error("Expected route to still be active")
	}

	// Update until route completes
	for i := 0; i < 100; i++ {
		rs.Update(0.01, entities, eventDispatcher)
		if len(rs.routes) == 0 {
			break
		}
	}

	// Verify route was completed and removed
	if len(rs.routes) != 0 {
		t.Errorf("Expected routes count to be 0 after completion, got %d", len(rs.routes))
	}
}

func TestRoutingSystem_RouteProperties(t *testing.T) {
	rs := NewRoutingSystem()

	// Create a route
	startX, startY := 100.0, 200.0
	endX, endY := 300.0, 400.0
	routeColor := color.RGBA{255, 0, 0, 255}

	rs.CreateRoute(startX, startY, endX, endY, routeColor)

	route := rs.routes[0]

	// Test route properties
	if route.StartX != startX {
		t.Errorf("Expected StartX to be %f, got %f", startX, route.StartX)
	}
	if route.StartY != startY {
		t.Errorf("Expected StartY to be %f, got %f", startY, route.StartY)
	}
	if route.EndX != endX {
		t.Errorf("Expected EndX to be %f, got %f", endX, route.EndX)
	}
	if route.EndY != endY {
		t.Errorf("Expected EndY to be %f, got %f", endY, route.EndY)
	}
	if route.Progress != 0.0 {
		t.Errorf("Expected Progress to be 0.0, got %f", route.Progress)
	}
	if route.Speed != 1.5 {
		t.Errorf("Expected Speed to be 1.5, got %f", route.Speed)
	}
	if route.Color != routeColor {
		t.Errorf("Expected Color to be %v, got %v", routeColor, route.Color)
	}
	if !route.Active {
		t.Error("Expected route to be active")
	}
}

func TestRoutingSystem_RouteProgress(t *testing.T) {
	rs := NewRoutingSystem()

	// Create a route
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})

	route := rs.routes[0]
	initialProgress := route.Progress

	// Update the route
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}
	rs.Update(0.016, entities, eventDispatcher)

	// Verify progress increased
	if route.Progress <= initialProgress {
		t.Errorf("Expected progress to increase, got %f", route.Progress)
	}

	// Verify progress is less than 1.0 (not completed)
	if route.Progress >= 1.0 {
		t.Errorf("Expected progress to be less than 1.0, got %f", route.Progress)
	}
}

func TestRoutingSystem_RouteCompletion(t *testing.T) {
	rs := NewRoutingSystem()

	// Create a route
	rs.CreateRoute(100, 200, 300, 400, color.RGBA{255, 0, 0, 255})

	// Manually set progress to completion
	rs.routes[0].Progress = 1.0

	// Update the route
	eventDispatcher := events.NewEventDispatcher()
	entities := []Entity{}
	rs.Update(0.016, entities, eventDispatcher)

	// Verify route was removed after completion
	if len(rs.routes) != 0 {
		t.Errorf("Expected routes count to be 0 after completion, got %d", len(rs.routes))
	}
}
