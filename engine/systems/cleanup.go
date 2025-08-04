package systems

import (
	"fmt"
	"lbbaspack/engine/events"
)

const SystemTypeCleanup SystemType = "cleanup"

// CleanupSystem removes inactive entities from the world
type CleanupSystem struct {
	BaseSystem
}

func NewCleanupSystem() *CleanupSystem {
	return &CleanupSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{}, // No required components, processes all entities
		},
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (cs *CleanupSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeCleanup,
		System:       cs,
		Dependencies: []SystemType{SystemTypeCollision},
		Conflicts:    []SystemType{},
		Provides:     []string{"entity_cleanup", "memory_management"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     false,
	}
}

func (cs *CleanupSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Count inactive entities before cleanup
	inactiveCount := 0
	for _, entity := range entities {
		if !entity.IsActive() {
			inactiveCount++
		}
	}

	if inactiveCount > 0 {
		fmt.Printf("[CleanupSystem] Found %d inactive entities to remove\n", inactiveCount)

		// Log inactive entities for debugging
		for _, entity := range entities {
			if !entity.IsActive() {
				if entityInterface, ok := entity.(interface{ GetComponentNames() []string }); ok {
					fmt.Printf("[CleanupSystem] Found inactive entity with components: %v\n", entityInterface.GetComponentNames())
				}
			}
		}

		// Note: Actual removal is handled by the main game loop calling World.RemoveInactiveEntities()
		// This system just identifies and logs inactive entities
	}
}

func (cs *CleanupSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Subscribe to game state change events to trigger cleanup
	eventDispatcher.Subscribe(events.EventGameStart, func(event *events.Event) {
		fmt.Println("[CleanupSystem] Game start event received, cleanup will be handled by main game loop")
	})
}
