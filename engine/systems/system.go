package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

// Entity interface defines what an entity must provide
// This is defined where it's used (in systems package)
type Entity interface {
	GetComponent(componentType string) components.Component
	HasComponent(componentType string) bool
	IsActive() bool
	GetComponentByName(typeName string) components.Component
	// Type-safe component getters
	GetTransform() components.TransformComponent
	GetSprite() components.SpriteComponent
	GetCollider() components.ColliderComponent
	GetPhysics() components.PhysicsComponent
	GetPacketType() components.PacketTypeComponent
	GetState() components.StateComponent
	GetCombo() components.ComboComponent
	GetSLA() components.SLAComponent
	GetBackendAssignment() components.BackendAssignmentComponent
	GetPowerUpType() components.PowerUpTypeComponent
}

// System represents a game system that processes entities
type System interface {
	Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher)
	GetRequiredComponents() []string
}

// BaseSystem provides common functionality for systems
type BaseSystem struct {
	RequiredComponents []string
}

// GetRequiredComponents returns the components required by this system
func (bs *BaseSystem) GetRequiredComponents() []string {
	return bs.RequiredComponents
}

// FilterEntities returns entities that have all required components
func (bs *BaseSystem) FilterEntities(entities []Entity) []Entity {
	var filtered []Entity

	for _, entity := range entities {
		if !entity.IsActive() {
			continue
		}

		hasAllComponents := true
		for _, componentType := range bs.RequiredComponents {
			if !entity.HasComponent(componentType) {
				hasAllComponents = false
				break
			}
		}

		if hasAllComponents {
			filtered = append(filtered, entity)
		}
	}

	return filtered
}
