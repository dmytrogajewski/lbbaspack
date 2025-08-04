package systems

import (
	"fmt"
	"image/color"
	"lbbaspack/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const SystemTypeRouting SystemType = "routing"

type RoutingSystem struct {
	BaseSystem
	routes []*Route
}

type Route struct {
	StartX, StartY float64
	EndX, EndY     float64
	Progress       float64
	Speed          float64
	Color          color.RGBA
	Active         bool
}

func NewRoutingSystem() *RoutingSystem {
	return &RoutingSystem{
		BaseSystem: BaseSystem{},
		routes:     make([]*Route, 0),
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
	// Update existing routes
	for i := len(rs.routes) - 1; i >= 0; i-- {
		route := rs.routes[i]
		if route.Active {
			route.Progress += route.Speed * deltaTime
			if route.Progress >= 1.0 {
				route.Active = false
			}
		}
	}

	// Remove completed routes in a separate loop to avoid index issues
	for i := len(rs.routes) - 1; i >= 0; i-- {
		if !rs.routes[i].Active {
			rs.routes = append(rs.routes[:i], rs.routes[i+1:]...)
		}
	}
}

func (rs *RoutingSystem) Draw(screen *ebiten.Image) {
	// Check if screen is nil to avoid panic
	if screen == nil {
		return
	}

	// Draw all active routes
	for _, route := range rs.routes {
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

func (rs *RoutingSystem) CreateRoute(startX, startY, endX, endY float64, packetColor color.RGBA) {
	route := &Route{
		StartX:   startX,
		StartY:   startY,
		EndX:     endX,
		EndY:     endY,
		Progress: 0.0,
		Speed:    1.5, // Slower speed for better visibility
		Color:    packetColor,
		Active:   true,
	}
	rs.routes = append(rs.routes, route)
}

func (rs *RoutingSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Listen for packet caught events to create routing visualization
	eventDispatcher.Subscribe(events.EventPacketCaught, func(event *events.Event) {
		fmt.Println("[RoutingSystem] Packet caught event received")
		if event.Data.Packet != nil {
			if packetEntity, ok := event.Data.Packet.(Entity); ok {
				// Get packet position and color
				transformComp := packetEntity.GetTransform()
				spriteComp := packetEntity.GetSprite()
				if transformComp == nil || spriteComp == nil {
					fmt.Println("[RoutingSystem] Packet entity missing transform or sprite")
					return
				}

				transform := transformComp
				sprite := spriteComp

				// Find target backend using round-robin
				// We'll need to track which backend to send to next
				// For now, create routes to different backends based on packet position
				startX := transform.GetX() + 7.5 // Center of packet
				startY := transform.GetY() + 7.5

				// Route to different backends based on packet X position
				backendIndex := int(startX/200) % 4 // 4 backends
				if backendIndex < 0 {
					backendIndex = 0
				}
				if backendIndex >= 4 {
					backendIndex = 3
				}

				// Calculate backend position
				backendWidth := 120.0
				backendSpacing := (800.0 - backendWidth*4.0) / 5.0
				backendX := backendSpacing + float64(backendIndex)*(backendWidth+backendSpacing)
				backendY := 550.0 // Backend Y position

				endX := backendX + backendWidth/2
				endY := backendY + 20.0 // Center of backend

				fmt.Printf("[RoutingSystem] Creating route from (%.1f, %.1f) to (%.1f, %.1f) for backend %d\n",
					startX, startY, endX, endY, backendIndex)
				rs.CreateRoute(startX, startY, endX, endY, sprite.GetColor())
			} else {
				fmt.Println("[RoutingSystem] Packet entity is not of correct type")
			}
		} else {
			fmt.Println("[RoutingSystem] Packet data is nil")
		}
	})
}
