package systems

import (
	"fmt"
	"image/color"
	"lbbaspack/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

type RenderSystem struct {
	BaseSystem
	callCount int
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Transform",
				"Sprite",
			},
		},
		callCount: 0,
	}
}

// Add a new method to set the screen for each frame
func (rs *RenderSystem) UpdateWithScreen(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher, screen *ebiten.Image) {
	rs.callCount++

	// Clear screen with solid background first
	screen.Fill(color.RGBA{20, 20, 40, 255})

	activeCount := 0
	visibleCount := 0
	filteredEntities := rs.FilterEntities(entities)
	for _, entity := range filteredEntities {
		if activeEntity, ok := entity.(interface{ IsActive() bool }); ok && activeEntity.IsActive() {
			activeCount++
			if spriteComp := entity.GetSprite(); spriteComp != nil {
				if spriteComp.IsVisible() {
					visibleCount++
				}
			}
		}
	}

	fmt.Printf("[RenderSystem] Frame %d: %d filtered, %d active, %d visible\n", rs.callCount, len(filteredEntities), activeCount, visibleCount)

	// Render all active and visible entities
	for _, entity := range filteredEntities {
		if activeEntity, ok := entity.(interface{ IsActive() bool }); ok && activeEntity.IsActive() {
			transformComp := entity.GetTransform()
			spriteComp := entity.GetSprite()
			if transformComp != nil && spriteComp != nil {
				if spriteComp.IsVisible() {
					if entityInterface, ok := entity.(interface{ GetComponentNames() []string }); ok {
						fmt.Printf("[RenderSystem] Rendering entity at (%.1f, %.1f) with components: %v\n", transformComp.GetX(), transformComp.GetY(), entityInterface.GetComponentNames())
					}
					vector.DrawFilledRect(screen, float32(transformComp.GetX()), float32(transformComp.GetY()), float32(spriteComp.GetWidth()), float32(spriteComp.GetHeight()), spriteComp.GetColor(), false)

					// Draw labels
					label := ""

					// Debug: List all component names for this entity
					if entityInterface, ok := entity.(interface{ GetComponentNames() []string }); ok {
						componentNames := entityInterface.GetComponentNames()
						fmt.Printf("[RenderSystem] Entity at (%.1f, %.1f) has components: %v\n", transformComp.GetX(), transformComp.GetY(), componentNames)
					}

					if colliderComp := entity.GetCollider(); colliderComp != nil {
						fmt.Printf("[RenderSystem] Entity has collider with tag: %s\n", colliderComp.GetTag())
						if colliderComp.GetTag() == "packet" {
							// Get packet type for proper label
							if packetTypeComp := entity.GetPacketType(); packetTypeComp != nil {
								label = packetTypeComp.GetName() // Show actual packet type (HTTP, TCP, etc.)
								fmt.Printf("[RenderSystem] Drawing packet label: %s at (%.1f, %.1f)\n", label, transformComp.GetX(), transformComp.GetY())
							} else {
								label = "packet"
								fmt.Printf("[RenderSystem] Drawing generic packet label at (%.1f, %.1f)\n", transformComp.GetX(), transformComp.GetY())
							}
						} else if colliderComp.GetTag() == "loadbalancer" {
							label = "LBaaS"
							fmt.Printf("[RenderSystem] Drawing LBaaS label at (%.1f, %.1f)\n", transformComp.GetX(), transformComp.GetY())
						} else if colliderComp.GetTag() == "backend" {
							// Backend labels are handled by the BackendAssignment component
							fmt.Printf("[RenderSystem] Found backend entity at (%.1f, %.1f)\n", transformComp.GetX(), transformComp.GetY())
						}
					} else {
						if entityInterface, ok := entity.(interface{ GetComponentNames() []string }); ok {
							fmt.Printf("[RenderSystem] Entity at (%.1f, %.1f) has no collider component, components: %v\n", transformComp.GetX(), transformComp.GetY(), entityInterface.GetComponentNames())
						}
						fmt.Printf("[RenderSystem] Entity at (%.1f, %.1f) has no collider component\n", transformComp.GetX(), transformComp.GetY())
					}
					if powerUpComp := entity.GetPowerUpType(); powerUpComp != nil {
						label = powerUpComp.GetName()
						fmt.Printf("[RenderSystem] Drawing power-up label: %s at (%.1f, %.1f)\n", label, transformComp.GetX(), transformComp.GetY())
					}
					if backendComp := entity.GetBackendAssignment(); backendComp != nil {
						label = fmt.Sprintf("Backend %d", backendComp.GetBackendID())
						fmt.Printf("[RenderSystem] Drawing backend label: %s at (%.1f, %.1f)\n", label, transformComp.GetX(), transformComp.GetY())
					}
					if label != "" {
						textX := int(transformComp.GetX())
						textY := int(transformComp.GetY()) - 5
						if textY < 0 {
							textY = int(transformComp.GetY()) + int(spriteComp.GetHeight()) + 12
						}
						// Use a more visible color and ensure text is drawn
						textColor := color.RGBA{255, 255, 255, 255} // Bright white
						text.Draw(screen, label, basicfont.Face7x13, textX, textY, textColor)
						fmt.Printf("[RenderSystem] Actually drawing text '%s' at (%d, %d)\n", label, textX, textY)
					}
				}
			}
		}
	}
}
