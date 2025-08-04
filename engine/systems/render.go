package systems

import (
	"fmt"
	"image/color"
	"lbbaspack/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"

	"lbbaspack/engine/components"
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
				if sprite, ok := spriteComp.(components.SpriteComponent); ok && sprite.IsVisible() {
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
				transform, ok1 := transformComp.(components.TransformComponent)
				sprite, ok2 := spriteComp.(components.SpriteComponent)
				if ok1 && ok2 && sprite.IsVisible() {
					if entityInterface, ok := entity.(interface{ GetComponentNames() []string }); ok {
						fmt.Printf("[RenderSystem] Rendering entity at (%.1f, %.1f) with components: %v\n", transform.GetX(), transform.GetY(), entityInterface.GetComponentNames())
					}
					ebitenutil.DrawRect(screen, transform.GetX(), transform.GetY(), sprite.GetWidth(), sprite.GetHeight(), sprite.GetColor())

					// Draw labels
					label := ""

					// Debug: List all component names for this entity
					if entityInterface, ok := entity.(interface{ GetComponentNames() []string }); ok {
						componentNames := entityInterface.GetComponentNames()
						fmt.Printf("[RenderSystem] Entity at (%.1f, %.1f) has components: %v\n", transform.GetX(), transform.GetY(), componentNames)
					}

					if colliderComp := entity.GetCollider(); colliderComp != nil {
						if collider, ok := colliderComp.(components.ColliderComponent); ok {
							fmt.Printf("[RenderSystem] Entity has collider with tag: %s\n", collider.GetTag())
							if collider.GetTag() == "packet" {
								// Get packet type for proper label
								if packetTypeComp := entity.GetPacketType(); packetTypeComp != nil {
									if packetType, ok := packetTypeComp.(components.PacketTypeComponent); ok {
										label = packetType.GetName() // Show actual packet type (HTTP, TCP, etc.)
										fmt.Printf("[RenderSystem] Drawing packet label: %s at (%.1f, %.1f)\n", label, transform.GetX(), transform.GetY())
									}
								} else {
									label = "packet"
									fmt.Printf("[RenderSystem] Drawing generic packet label at (%.1f, %.1f)\n", transform.GetX(), transform.GetY())
								}
							} else if collider.GetTag() == "loadbalancer" {
								label = "LBaaS"
								fmt.Printf("[RenderSystem] Drawing LBaaS label at (%.1f, %.1f)\n", transform.GetX(), transform.GetY())
							} else if collider.GetTag() == "backend" {
								// Backend labels are handled by the BackendAssignment component
								fmt.Printf("[RenderSystem] Found backend entity at (%.1f, %.1f)\n", transform.GetX(), transform.GetY())
							}
						}
					} else {
						if entityInterface, ok := entity.(interface{ GetComponentNames() []string }); ok {
							fmt.Printf("[RenderSystem] Entity at (%.1f, %.1f) has no collider component, components: %v\n", transform.GetX(), transform.GetY(), entityInterface.GetComponentNames())
						}
						fmt.Printf("[RenderSystem] Entity at (%.1f, %.1f) has no collider component\n", transform.GetX(), transform.GetY())
					}
					if powerUpComp := entity.GetPowerUpType(); powerUpComp != nil {
						if powerUp, ok := powerUpComp.(components.PowerUpTypeComponent); ok {
							label = powerUp.GetName()
							fmt.Printf("[RenderSystem] Drawing power-up label: %s at (%.1f, %.1f)\n", label, transform.GetX(), transform.GetY())
						}
					}
					if backendComp := entity.GetBackendAssignment(); backendComp != nil {
						if backend, ok := backendComp.(components.BackendAssignmentComponent); ok {
							label = fmt.Sprintf("Backend %d", backend.GetBackendID())
							fmt.Printf("[RenderSystem] Drawing backend label: %s at (%.1f, %.1f)\n", label, transform.GetX(), transform.GetY())
						}
					}
					if label != "" {
						textX := int(transform.GetX())
						textY := int(transform.GetY()) - 5
						if textY < 0 {
							textY = int(transform.GetY()) + int(sprite.GetHeight()) + 12
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
