package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"testing"
)

// systemTestEntity implements the Entity interface for testing
type systemTestEntity struct {
	entity          *entities.Entity
	componentsAdded []string
	isActive        bool
}

func newSystemTestEntity(id uint64) *systemTestEntity {
	return &systemTestEntity{
		entity:          entities.NewEntity(id),
		componentsAdded: make([]string, 0),
		isActive:        true,
	}
}

func (ste *systemTestEntity) GetID() uint64 {
	return ste.entity.GetID()
}

func (ste *systemTestEntity) GetComponent(componentType string) components.Component {
	return ste.entity.GetComponent(componentType)
}

func (ste *systemTestEntity) HasComponent(componentType string) bool {
	return ste.entity.HasComponent(componentType)
}

func (ste *systemTestEntity) IsActive() bool {
	return ste.isActive
}

func (ste *systemTestEntity) GetComponentByName(typeName string) components.Component {
	return ste.entity.GetComponentByName(typeName)
}

func (ste *systemTestEntity) GetTransform() components.TransformComponent {
	if comp := ste.entity.GetComponent("Transform"); comp != nil {
		return comp.(components.TransformComponent)
	}
	return nil
}

func (ste *systemTestEntity) GetSprite() components.SpriteComponent {
	if comp := ste.entity.GetComponent("Sprite"); comp != nil {
		return comp.(components.SpriteComponent)
	}
	return nil
}

func (ste *systemTestEntity) GetCollider() components.ColliderComponent {
	if comp := ste.entity.GetComponent("Collider"); comp != nil {
		return comp.(components.ColliderComponent)
	}
	return nil
}

func (ste *systemTestEntity) GetPhysics() components.PhysicsComponent {
	if comp := ste.entity.GetComponent("Physics"); comp != nil {
		return comp.(components.PhysicsComponent)
	}
	return nil
}

func (ste *systemTestEntity) GetPacketType() components.PacketTypeComponent {
	if comp := ste.entity.GetComponent("PacketType"); comp != nil {
		return comp.(components.PacketTypeComponent)
	}
	return nil
}

func (ste *systemTestEntity) GetState() components.StateComponent {
	if comp := ste.entity.GetComponent("State"); comp != nil {
		return comp.(components.StateComponent)
	}
	return nil
}

func (ste *systemTestEntity) GetCombo() components.ComboComponent {
	if comp := ste.entity.GetComponent("Combo"); comp != nil {
		return comp.(components.ComboComponent)
	}
	return nil
}

func (ste *systemTestEntity) GetSLA() components.SLAComponent {
	if comp := ste.entity.GetComponent("SLA"); comp != nil {
		return comp.(components.SLAComponent)
	}
	return nil
}

func (ste *systemTestEntity) GetBackendAssignment() components.BackendAssignmentComponent {
	if comp := ste.entity.GetComponent("BackendAssignment"); comp != nil {
		return comp.(components.BackendAssignmentComponent)
	}
	return nil
}

func (ste *systemTestEntity) GetPowerUpType() components.PowerUpTypeComponent {
	if comp := ste.entity.GetComponent("PowerUpType"); comp != nil {
		return comp.(components.PowerUpTypeComponent)
	}
	return nil
}

func (ste *systemTestEntity) AddComponent(component components.Component) {
	ste.componentsAdded = append(ste.componentsAdded, component.GetType())
	ste.entity.AddComponent(component)
}

func (ste *systemTestEntity) GetComponentNames() []string {
	return ste.componentsAdded
}

func (ste *systemTestEntity) SetActive(active bool) {
	ste.isActive = active
}

// TestBaseSystem_GetRequiredComponents tests the GetRequiredComponents method
func TestBaseSystem_GetRequiredComponents(t *testing.T) {
	tests := []struct {
		name               string
		requiredComponents []string
		expectedComponents []string
	}{
		{
			name:               "Empty required components",
			requiredComponents: []string{},
			expectedComponents: []string{},
		},
		{
			name:               "Single required component",
			requiredComponents: []string{"Transform"},
			expectedComponents: []string{"Transform"},
		},
		{
			name:               "Multiple required components",
			requiredComponents: []string{"Transform", "Sprite", "Physics"},
			expectedComponents: []string{"Transform", "Sprite", "Physics"},
		},
		{
			name:               "Duplicate components",
			requiredComponents: []string{"Transform", "Transform", "Sprite"},
			expectedComponents: []string{"Transform", "Transform", "Sprite"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := &BaseSystem{
				RequiredComponents: tt.requiredComponents,
			}

			result := bs.GetRequiredComponents()

			if len(result) != len(tt.expectedComponents) {
				t.Errorf("GetRequiredComponents() returned %d components, expected %d", len(result), len(tt.expectedComponents))
			}

			for i, component := range result {
				if component != tt.expectedComponents[i] {
					t.Errorf("GetRequiredComponents()[%d] = %s, expected %s", i, component, tt.expectedComponents[i])
				}
			}
		})
	}
}

// TestBaseSystem_FilterEntities tests the FilterEntities method
func TestBaseSystem_FilterEntities(t *testing.T) {
	tests := []struct {
		name               string
		requiredComponents []string
		entities           []Entity
		expectedCount      int
		description        string
	}{
		{
			name:               "No entities",
			requiredComponents: []string{"Transform"},
			entities:           []Entity{},
			expectedCount:      0,
			description:        "Should return empty slice when no entities provided",
		},
		{
			name:               "No required components",
			requiredComponents: []string{},
			entities: func() []Entity {
				entity1 := newSystemTestEntity(1)
				entity1.AddComponent(components.NewTransform(0, 0))
				entity2 := newSystemTestEntity(2)
				entity2.AddComponent(components.NewSprite(10, 10, components.RandomPacketColor()))
				return []Entity{entity1, entity2}
			}(),
			expectedCount: 2,
			description:   "Should return all entities when no components required",
		},
		{
			name:               "All entities have required components",
			requiredComponents: []string{"Transform"},
			entities: func() []Entity {
				entity1 := newSystemTestEntity(1)
				entity1.AddComponent(components.NewTransform(0, 0))
				entity2 := newSystemTestEntity(2)
				entity2.AddComponent(components.NewTransform(10, 10))
				return []Entity{entity1, entity2}
			}(),
			expectedCount: 2,
			description:   "Should return all entities when they all have required components",
		},
		{
			name:               "Some entities have required components",
			requiredComponents: []string{"Transform", "Sprite"},
			entities: func() []Entity {
				entity1 := newSystemTestEntity(1)
				entity1.AddComponent(components.NewTransform(0, 0))
				entity1.AddComponent(components.NewSprite(10, 10, components.RandomPacketColor()))

				entity2 := newSystemTestEntity(2)
				entity2.AddComponent(components.NewTransform(10, 10))
				// Missing Sprite component

				entity3 := newSystemTestEntity(3)
				entity3.AddComponent(components.NewTransform(20, 20))
				entity3.AddComponent(components.NewSprite(15, 15, components.RandomPacketColor()))

				return []Entity{entity1, entity2, entity3}
			}(),
			expectedCount: 2,
			description:   "Should return only entities with all required components",
		},
		{
			name:               "No entities have required components",
			requiredComponents: []string{"Transform", "Sprite", "Physics"},
			entities: func() []Entity {
				entity1 := newSystemTestEntity(1)
				entity1.AddComponent(components.NewTransform(0, 0))
				// Missing Sprite and Physics

				entity2 := newSystemTestEntity(2)
				entity2.AddComponent(components.NewSprite(10, 10, components.RandomPacketColor()))
				// Missing Transform and Physics

				return []Entity{entity1, entity2}
			}(),
			expectedCount: 0,
			description:   "Should return empty slice when no entities have all required components",
		},
		{
			name:               "Inactive entities are filtered out",
			requiredComponents: []string{"Transform"},
			entities: func() []Entity {
				entity1 := newSystemTestEntity(1)
				entity1.AddComponent(components.NewTransform(0, 0))
				entity1.SetActive(true)

				entity2 := newSystemTestEntity(2)
				entity2.AddComponent(components.NewTransform(10, 10))
				entity2.SetActive(false) // Inactive

				entity3 := newSystemTestEntity(3)
				entity3.AddComponent(components.NewTransform(20, 20))
				entity3.SetActive(true)

				return []Entity{entity1, entity2, entity3}
			}(),
			expectedCount: 2,
			description:   "Should filter out inactive entities",
		},
		{
			name:               "Multiple required components",
			requiredComponents: []string{"Transform", "Sprite", "Physics"},
			entities: func() []Entity {
				entity1 := newSystemTestEntity(1)
				entity1.AddComponent(components.NewTransform(0, 0))
				entity1.AddComponent(components.NewSprite(10, 10, components.RandomPacketColor()))
				entity1.AddComponent(components.NewPhysics())

				entity2 := newSystemTestEntity(2)
				entity2.AddComponent(components.NewTransform(10, 10))
				entity2.AddComponent(components.NewSprite(15, 15, components.RandomPacketColor()))
				// Missing Physics

				entity3 := newSystemTestEntity(3)
				entity3.AddComponent(components.NewTransform(20, 20))
				entity3.AddComponent(components.NewSprite(20, 20, components.RandomPacketColor()))
				entity3.AddComponent(components.NewPhysics())

				return []Entity{entity1, entity2, entity3}
			}(),
			expectedCount: 2,
			description:   "Should return entities with all multiple required components",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := &BaseSystem{
				RequiredComponents: tt.requiredComponents,
			}

			result := bs.FilterEntities(tt.entities)

			if len(result) != tt.expectedCount {
				t.Errorf("FilterEntities() returned %d entities, expected %d. %s",
					len(result), tt.expectedCount, tt.description)
			}

			// Verify that returned entities have all required components
			for _, entity := range result {
				if !entity.IsActive() {
					t.Errorf("FilterEntities() returned inactive entity")
				}

				for _, componentType := range tt.requiredComponents {
					if !entity.HasComponent(componentType) {
						t.Errorf("FilterEntities() returned entity missing required component: %s", componentType)
					}
				}
			}
		})
	}
}

// TestBaseSystem_FilterEntities_EdgeCases tests edge cases for FilterEntities
func TestBaseSystem_FilterEntities_EdgeCases(t *testing.T) {
	t.Run("Nil entities slice", func(t *testing.T) {
		bs := &BaseSystem{
			RequiredComponents: []string{"Transform"},
		}

		result := bs.FilterEntities(nil)

		// When no entities are provided, result can be nil (empty slice)
		if len(result) != 0 {
			t.Errorf("FilterEntities() returned %d entities, expected 0", len(result))
		}
	})

	t.Run("Entity with nil components", func(t *testing.T) {
		bs := &BaseSystem{
			RequiredComponents: []string{"Transform"},
		}

		entity := newSystemTestEntity(1)
		// Don't add any components

		entities := []Entity{entity}
		result := bs.FilterEntities(entities)

		if len(result) != 0 {
			t.Errorf("FilterEntities() returned %d entities, expected 0 for entity without required components", len(result))
		}
	})

	t.Run("Mixed active and inactive entities", func(t *testing.T) {
		bs := &BaseSystem{
			RequiredComponents: []string{"Transform"},
		}

		entity1 := newSystemTestEntity(1)
		entity1.AddComponent(components.NewTransform(0, 0))
		entity1.SetActive(true)

		entity2 := newSystemTestEntity(2)
		entity2.AddComponent(components.NewTransform(10, 10))
		entity2.SetActive(false)

		entity3 := newSystemTestEntity(3)
		entity3.AddComponent(components.NewTransform(20, 20))
		entity3.SetActive(true)

		entities := []Entity{entity1, entity2, entity3}
		result := bs.FilterEntities(entities)

		if len(result) != 2 {
			t.Errorf("FilterEntities() returned %d entities, expected 2 (only active ones)", len(result))
		}

		// Verify only active entities are returned
		for _, entity := range result {
			if !entity.IsActive() {
				t.Error("FilterEntities() returned inactive entity")
			}
		}
	})
}

// TestEntityInterface tests the Entity interface methods
func TestEntityInterface(t *testing.T) {
	t.Run("Component getters", func(t *testing.T) {
		entity := newSystemTestEntity(1)
		entity.AddComponent(components.NewTransform(10, 20))
		entity.AddComponent(components.NewSprite(15, 15, components.RandomPacketColor()))
		entity.AddComponent(components.NewPhysics())
		entity.AddComponent(components.NewCollider(15, 15, "test"))

		// Test HasComponent
		if !entity.HasComponent("Transform") {
			t.Error("HasComponent() should return true for existing component")
		}
		if entity.HasComponent("NonExistent") {
			t.Error("HasComponent() should return false for non-existent component")
		}

		// Test GetComponent
		transform := entity.GetComponent("Transform")
		if transform == nil {
			t.Error("GetComponent() should return component for existing component")
		}
		if transform.GetType() != "Transform" {
			t.Errorf("GetComponent() returned component with wrong type: %s", transform.GetType())
		}

		// Test GetComponentByName
		sprite := entity.GetComponentByName("Sprite")
		if sprite == nil {
			t.Error("GetComponentByName() should return component for existing component")
		}

		// Test type-safe getters
		transformComp := entity.GetTransform()
		if transformComp == nil {
			t.Error("GetTransform() should return TransformComponent")
		}

		spriteComp := entity.GetSprite()
		if spriteComp == nil {
			t.Error("GetSprite() should return SpriteComponent")
		}

		physicsComp := entity.GetPhysics()
		if physicsComp == nil {
			t.Error("GetPhysics() should return PhysicsComponent")
		}

		colliderComp := entity.GetCollider()
		if colliderComp == nil {
			t.Error("GetCollider() should return ColliderComponent")
		}

		// Test non-existent components
		if entity.GetPacketType() != nil {
			t.Error("GetPacketType() should return nil for non-existent component")
		}
	})

	t.Run("IsActive", func(t *testing.T) {
		entity := newSystemTestEntity(1)

		if !entity.IsActive() {
			t.Error("New entity should be active by default")
		}

		entity.SetActive(false)
		if entity.IsActive() {
			t.Error("Entity should be inactive after SetActive(false)")
		}

		entity.SetActive(true)
		if !entity.IsActive() {
			t.Error("Entity should be active after SetActive(true)")
		}
	})
}

// TestBaseSystem_Integration tests integration scenarios
func TestBaseSystem_Integration(t *testing.T) {
	t.Run("Complex filtering scenario", func(t *testing.T) {
		bs := &BaseSystem{
			RequiredComponents: []string{"Transform", "Sprite", "Physics"},
		}

		// Create various entities with different component combinations
		entities := []Entity{
			func() Entity {
				entity := newSystemTestEntity(1)
				entity.AddComponent(components.NewTransform(0, 0))
				entity.AddComponent(components.NewSprite(10, 10, components.RandomPacketColor()))
				entity.AddComponent(components.NewPhysics())
				return entity
			}(),
			func() Entity {
				entity := newSystemTestEntity(2)
				entity.AddComponent(components.NewTransform(10, 10))
				entity.AddComponent(components.NewSprite(15, 15, components.RandomPacketColor()))
				// Missing Physics
				return entity
			}(),
			func() Entity {
				entity := newSystemTestEntity(3)
				entity.AddComponent(components.NewTransform(20, 20))
				entity.AddComponent(components.NewSprite(20, 20, components.RandomPacketColor()))
				entity.AddComponent(components.NewPhysics())
				entity.SetActive(false) // Inactive
				return entity
			}(),
			func() Entity {
				entity := newSystemTestEntity(4)
				entity.AddComponent(components.NewTransform(30, 30))
				entity.AddComponent(components.NewSprite(25, 25, components.RandomPacketColor()))
				entity.AddComponent(components.NewPhysics())
				return entity
			}(),
			func() Entity {
				entity := newSystemTestEntity(5)
				// No components at all
				return entity
			}(),
		}

		result := bs.FilterEntities(entities)

		// Should only return entities 1 and 4 (have all components and are active)
		if len(result) != 2 {
			t.Errorf("FilterEntities() returned %d entities, expected 2", len(result))
		}

		// Verify the returned entities have all required components and are active
		for _, entity := range result {
			if !entity.IsActive() {
				t.Error("FilterEntities() returned inactive entity")
			}

			for _, componentType := range bs.RequiredComponents {
				if !entity.HasComponent(componentType) {
					t.Errorf("FilterEntities() returned entity missing required component: %s", componentType)
				}
			}
		}
	})
}
