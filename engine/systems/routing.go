package systems

import (
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const SystemTypeRouting SystemType = "routing"

type RoutingSystem struct {
	BaseSystem
}

type Route = components.Route

func NewRoutingSystem() *RoutingSystem {
	return &RoutingSystem{
		BaseSystem: BaseSystem{},
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (rs *RoutingSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeRouting,
		System:       rs,
		Dependencies: []SystemType{SystemTypeCollision},
		Conflicts:    []SystemType{},
		Provides:     []string{"route_visualization", "network_paths"},
		Requires:     []string{},
		Drawable:     true,
		Optional:     true,
	}
}

func (rs *RoutingSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Update existing routes from RouteState component
	var state *components.RouteState
	for _, e := range entities {
		if comp := e.GetComponentByName("RouteState"); comp != nil {
			if s, ok := comp.(*components.RouteState); ok {
				state = s
				break
			}
		}
	}
	if state == nil {
		return
	}
	for i := len(state.Routes) - 1; i >= 0; i-- {
		route := state.Routes[i]
		if route.Active {
			route.Progress += route.Speed * deltaTime
			if route.Progress >= 1.0 {
				route.Active = false
			}
		}
	}
	for i := len(state.Routes) - 1; i >= 0; i-- {
		if !state.Routes[i].Active {
			state.Routes = append(state.Routes[:i], state.Routes[i+1:]...)
		}
	}
}

func (rs *RoutingSystem) Draw(screen *ebiten.Image, entities []Entity) {
	// Check if screen is nil to avoid panic
	if screen == nil {
		return
	}

	// Draw all active routes from RouteState
	var state *components.RouteState
	for _, e := range entities {
		if comp := e.GetComponentByName("RouteState"); comp != nil {
			if s, ok := comp.(*components.RouteState); ok {
				state = s
				break
			}
		}
	}
	if state == nil {
		return
	}
	for _, route := range state.Routes {
		if route.Active {
			currentX := route.StartX + (route.EndX-route.StartX)*route.Progress
			currentY := route.StartY + (route.EndY-route.StartY)*route.Progress

			// Draw thicker line from start to current position
			vector.StrokeLine(screen, float32(route.StartX), float32(route.StartY), float32(currentX), float32(currentY), 2, route.Color, false)

			// Draw packet at current position with larger size
			vector.DrawFilledCircle(screen, float32(currentX), float32(currentY), 4, route.Color, false)

			// Draw a small trail effect
			if route.Progress > 0.1 {
				trailX := route.StartX + (route.EndX-route.StartX)*(route.Progress-0.1)
				trailY := route.StartY + (route.EndY-route.StartY)*(route.Progress-0.1)
				vector.DrawFilledCircle(screen, float32(trailX), float32(trailY), 2, route.Color, false)
			}
		}
	}
}

func (rs *RoutingSystem) CreateRoute(startX, startY, endX, endY float64, packetColor color.RGBA, state *components.RouteState) {
	route := &components.Route{StartX: startX, StartY: startY, EndX: endX, EndY: endY, Progress: 0.0, Speed: 1.5, Color: packetColor, Active: true}
	state.Routes = append(state.Routes, route)
}

// GetRoutes returns the current routes for testing
func (rs *RoutingSystem) GetRoutes() []*Route { return nil }

func (rs *RoutingSystem) Initialize(eventDispatcher *events.EventDispatcher) {}
