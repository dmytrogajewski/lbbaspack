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
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Transform",
				"Sprite",
			},
		},
	}
}

// Add a new method to set the screen for each frame
func (rs *RenderSystem) UpdateWithScreen(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher, screen *ebiten.Image) {
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

	_ = activeCount
	_ = visibleCount

	// Render all active and visible entities
	for _, entity := range filteredEntities {
		if activeEntity, ok := entity.(interface{ IsActive() bool }); ok && activeEntity.IsActive() {
			transformComp := entity.GetTransform()
			spriteComp := entity.GetSprite()
			if transformComp != nil && spriteComp != nil {
				if spriteComp.IsVisible() {
					// debug: optional detailed logging can be added via a debug flag
					vector.DrawFilledRect(screen, float32(transformComp.GetX()), float32(transformComp.GetY()), float32(spriteComp.GetWidth()), float32(spriteComp.GetHeight()), spriteComp.GetColor(), false)

					// Draw labels
					label := ""

					// debug: component list logging removed for stateless compliance

					if colliderComp := entity.GetCollider(); colliderComp != nil {
						if colliderComp.GetTag() == "packet" {
							// Get packet type for proper label
							if packetTypeComp := entity.GetPacketType(); packetTypeComp != nil {
								label = packetTypeComp.GetName() // Show actual packet type (HTTP, TCP, etc.)
								// label derived from component
							} else {
								label = "packet"
							}
						} else if colliderComp.GetTag() == "loadbalancer" {
							label = "LBaaS"
						}
					}
					if powerUpComp := entity.GetPowerUpType(); powerUpComp != nil {
						label = powerUpComp.GetName()
					}
					if backendComp := entity.GetBackendAssignment(); backendComp != nil {
						label = "Backend " + fmt.Sprint(backendComp.GetBackendID())
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
					}
				}
			}
		}
	}
}
