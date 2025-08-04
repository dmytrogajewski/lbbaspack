package entities

import (
	"fmt"
	"lbbaspack/engine/components"
	"sync"
)

// Entity represents a game object with components
type Entity struct {
	ID         uint64
	Components map[string]components.Component
	Active     bool
	mu         sync.RWMutex
}

// NewEntity creates a new entity
func NewEntity(id uint64) *Entity {
	return &Entity{
		ID:         id,
		Components: make(map[string]components.Component),
		Active:     true,
	}
}

// AddComponent adds a component to the entity
func (e *Entity) AddComponent(component components.Component) {
	if component == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Components[component.GetType()] = component
	fmt.Printf("[Entity] Added component %s to entity %d\n", component.GetType(), e.ID)
}

// GetComponent retrieves a component by type
func (e *Entity) GetComponent(componentType string) components.Component {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Components[componentType]
}

// GetComponentByName implements systems.Entity interface for backward compatibility
func (e *Entity) GetComponentByName(typeName string) components.Component {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Components[typeName]
}

// HasComponent checks if entity has a specific component
func (e *Entity) HasComponent(componentType string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	_, exists := e.Components[componentType]
	return exists
}

// RemoveComponent removes a component from the entity
func (e *Entity) RemoveComponent(componentType string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.Components, componentType)
}

// SetActive sets the entity's active state
func (e *Entity) SetActive(active bool) {
	e.Active = active
}

// IsActive returns the entity's active state
func (e *Entity) IsActive() bool {
	return e.Active
}

// GetID returns the entity's ID
func (e *Entity) GetID() uint64 {
	return e.ID
}

// GetComponentNames returns all component type names for debugging
func (e *Entity) GetComponentNames() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	names := make([]string, 0, len(e.Components))
	for t := range e.Components {
		names = append(names, t)
	}
	return names
}

// Type-safe component getters
func (e *Entity) GetTransform() components.TransformComponent {
	if comp := e.GetComponent("Transform"); comp != nil {
		if transform, ok := comp.(components.TransformComponent); ok {
			return transform
		}
	}
	return nil
}

func (e *Entity) GetSprite() components.SpriteComponent {
	if comp := e.GetComponent("Sprite"); comp != nil {
		if sprite, ok := comp.(components.SpriteComponent); ok {
			return sprite
		}
	}
	return nil
}

func (e *Entity) GetCollider() components.ColliderComponent {
	if comp := e.GetComponent("Collider"); comp != nil {
		if collider, ok := comp.(components.ColliderComponent); ok {
			return collider
		}
	}
	return nil
}

func (e *Entity) GetPhysics() components.PhysicsComponent {
	if comp := e.GetComponent("Physics"); comp != nil {
		if physics, ok := comp.(components.PhysicsComponent); ok {
			return physics
		}
	}
	return nil
}

func (e *Entity) GetPacketType() components.PacketTypeComponent {
	if comp := e.GetComponent("PacketType"); comp != nil {
		if packetType, ok := comp.(components.PacketTypeComponent); ok {
			return packetType
		}
	}
	return nil
}

func (e *Entity) GetState() components.StateComponent {
	if comp := e.GetComponent("State"); comp != nil {
		if state, ok := comp.(components.StateComponent); ok {
			return state
		}
	}
	return nil
}

func (e *Entity) GetCombo() components.ComboComponent {
	if comp := e.GetComponent("Combo"); comp != nil {
		if combo, ok := comp.(components.ComboComponent); ok {
			return combo
		}
	}
	return nil
}

func (e *Entity) GetSLA() components.SLAComponent {
	if comp := e.GetComponent("SLA"); comp != nil {
		if sla, ok := comp.(components.SLAComponent); ok {
			return sla
		}
	}
	return nil
}

func (e *Entity) GetBackendAssignment() components.BackendAssignmentComponent {
	if comp := e.GetComponent("BackendAssignment"); comp != nil {
		if backend, ok := comp.(components.BackendAssignmentComponent); ok {
			return backend
		}
	}
	return nil
}

func (e *Entity) GetPowerUpType() components.PowerUpTypeComponent {
	if comp := e.GetComponent("PowerUpType"); comp != nil {
		if powerUp, ok := comp.(components.PowerUpTypeComponent); ok {
			return powerUp
		}
	}
	return nil
}
