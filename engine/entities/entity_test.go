package entities

import (
	"fmt"
	"image/color"
	"lbbaspack/engine/components"
	"testing"
)

// Mock component for testing
type testMockComponent struct {
	componentType string
}

func (mc *testMockComponent) GetType() string {
	return mc.componentType
}

// TestNewEntity tests the NewEntity constructor
func TestNewEntity(t *testing.T) {
	id := uint64(123)
	entity := NewEntity(id)

	if entity == nil {
		t.Fatal("Expected entity to be created")
	}

	if entity.ID != id {
		t.Errorf("Expected entity ID %d, got %d", id, entity.ID)
	}

	if entity.Components == nil {
		t.Error("Expected components map to be initialized")
	}

	if len(entity.Components) != 0 {
		t.Error("Expected components map to be empty initially")
	}

	if !entity.Active {
		t.Error("Expected entity to be active initially")
	}
}

// TestEntity_AddComponent tests the AddComponent method
func TestEntity_AddComponent(t *testing.T) {
	entity := NewEntity(1)

	// Create a mock component
	mockComponent := &testMockComponent{
		componentType: "TestComponent",
	}

	entity.AddComponent(mockComponent)

	if len(entity.Components) != 1 {
		t.Errorf("Expected 1 component, got %d", len(entity.Components))
	}

	if entity.Components["TestComponent"] != mockComponent {
		t.Error("Expected component to be stored correctly")
	}

	// Test adding another component
	mockComponent2 := &testMockComponent{
		componentType: "TestComponent2",
	}

	entity.AddComponent(mockComponent2)

	if len(entity.Components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(entity.Components))
	}

	if entity.Components["TestComponent2"] != mockComponent2 {
		t.Error("Expected second component to be stored correctly")
	}

	// Test overwriting existing component
	mockComponent3 := &testMockComponent{
		componentType: "TestComponent",
	}

	entity.AddComponent(mockComponent3)

	if len(entity.Components) != 2 {
		t.Errorf("Expected 2 components after overwrite, got %d", len(entity.Components))
	}

	if entity.Components["TestComponent"] != mockComponent3 {
		t.Error("Expected component to be overwritten")
	}
}

// TestEntity_GetComponent tests the GetComponent method
func TestEntity_GetComponent(t *testing.T) {
	entity := NewEntity(1)

	// Test getting non-existent component
	component := entity.GetComponent("NonExistent")
	if component != nil {
		t.Error("Expected nil for non-existent component")
	}

	// Test getting existing component
	mockComponent := &testMockComponent{
		componentType: "TestComponent",
	}
	entity.AddComponent(mockComponent)

	retrievedComponent := entity.GetComponent("TestComponent")
	if retrievedComponent != mockComponent {
		t.Error("Expected to retrieve the correct component")
	}
}

// TestEntity_GetComponentByName tests the GetComponentByName method
func TestEntity_GetComponentByName(t *testing.T) {
	entity := NewEntity(1)

	// Test getting non-existent component
	component := entity.GetComponentByName("NonExistent")
	if component != nil {
		t.Error("Expected nil for non-existent component")
	}

	// Test getting existing component
	mockComponent := &testMockComponent{
		componentType: "TestComponent",
	}
	entity.AddComponent(mockComponent)

	retrievedComponent := entity.GetComponentByName("TestComponent")
	if retrievedComponent != mockComponent {
		t.Error("Expected to retrieve the correct component")
	}
}

// TestEntity_HasComponent tests the HasComponent method
func TestEntity_HasComponent(t *testing.T) {
	entity := NewEntity(1)

	// Test non-existent component
	if entity.HasComponent("NonExistent") {
		t.Error("Expected false for non-existent component")
	}

	// Test existing component
	mockComponent := &testMockComponent{
		componentType: "TestComponent",
	}
	entity.AddComponent(mockComponent)

	if !entity.HasComponent("TestComponent") {
		t.Error("Expected true for existing component")
	}
}

// TestEntity_RemoveComponent tests the RemoveComponent method
func TestEntity_RemoveComponent(t *testing.T) {
	entity := NewEntity(1)

	// Add components
	mockComponent1 := &testMockComponent{componentType: "Component1"}
	mockComponent2 := &testMockComponent{componentType: "Component2"}
	entity.AddComponent(mockComponent1)
	entity.AddComponent(mockComponent2)

	if len(entity.Components) != 2 {
		t.Errorf("Expected 2 components before removal, got %d", len(entity.Components))
	}

	// Remove one component
	entity.RemoveComponent("Component1")

	if len(entity.Components) != 1 {
		t.Errorf("Expected 1 component after removal, got %d", len(entity.Components))
	}

	if entity.HasComponent("Component1") {
		t.Error("Expected Component1 to be removed")
	}

	if !entity.HasComponent("Component2") {
		t.Error("Expected Component2 to still exist")
	}

	// Test removing non-existent component (should not panic)
	entity.RemoveComponent("NonExistent")

	if len(entity.Components) != 1 {
		t.Errorf("Expected 1 component after removing non-existent, got %d", len(entity.Components))
	}
}

// TestEntity_SetActive tests the SetActive method
func TestEntity_SetActive(t *testing.T) {
	entity := NewEntity(1)

	// Test initial state
	if !entity.Active {
		t.Error("Expected entity to be active initially")
	}

	// Test setting to inactive
	entity.SetActive(false)
	if entity.Active {
		t.Error("Expected entity to be inactive after SetActive(false)")
	}

	// Test setting back to active
	entity.SetActive(true)
	if !entity.Active {
		t.Error("Expected entity to be active after SetActive(true)")
	}
}

// TestEntity_IsActive tests the IsActive method
func TestEntity_IsActive(t *testing.T) {
	entity := NewEntity(1)

	// Test initial state
	if !entity.IsActive() {
		t.Error("Expected entity to be active initially")
	}

	// Test after setting inactive
	entity.SetActive(false)
	if entity.IsActive() {
		t.Error("Expected entity to be inactive")
	}

	// Test after setting active
	entity.SetActive(true)
	if !entity.IsActive() {
		t.Error("Expected entity to be active")
	}
}

// TestEntity_GetComponentNames tests the GetComponentNames method
func TestEntity_GetComponentNames(t *testing.T) {
	entity := NewEntity(1)

	// Test empty entity
	names := entity.GetComponentNames()
	if len(names) != 0 {
		t.Errorf("Expected empty names list, got %d names", len(names))
	}

	// Add components
	mockComponent1 := &testMockComponent{componentType: "Component1"}
	mockComponent2 := &testMockComponent{componentType: "Component2"}
	mockComponent3 := &testMockComponent{componentType: "Component3"}

	entity.AddComponent(mockComponent1)
	entity.AddComponent(mockComponent2)
	entity.AddComponent(mockComponent3)

	names = entity.GetComponentNames()
	if len(names) != 3 {
		t.Errorf("Expected 3 component names, got %d", len(names))
	}

	// Check that all expected names are present
	expectedNames := map[string]bool{
		"Component1": true,
		"Component2": true,
		"Component3": true,
	}

	for _, name := range names {
		if !expectedNames[name] {
			t.Errorf("Unexpected component name: %s", name)
		}
	}
}

// TestEntity_TypeSafeGetters tests all the type-safe component getters
func TestEntity_TypeSafeGetters(t *testing.T) {
	entity := NewEntity(1)

	t.Run("GetTransform", func(t *testing.T) {
		// Test with no component
		transform := entity.GetTransform()
		if transform != nil {
			t.Error("Expected nil transform when no component exists")
		}

		// Test with wrong component type
		mockComponent := &testMockComponent{componentType: "Transform"}
		entity.AddComponent(mockComponent)
		transform = entity.GetTransform()
		if transform != nil {
			t.Error("Expected nil transform when component is wrong type")
		}

		// Test with correct component type
		entity.RemoveComponent("Transform")
		realTransform := components.NewTransform(0, 0)
		entity.AddComponent(realTransform)
		transform = entity.GetTransform()
		if transform == nil {
			t.Error("Expected transform when correct component exists")
		}
	})

	t.Run("GetSprite", func(t *testing.T) {
		// Test with no component
		sprite := entity.GetSprite()
		if sprite != nil {
			t.Error("Expected nil sprite when no component exists")
		}

		// Test with correct component type
		realSprite := components.NewSprite(32, 32, color.RGBA{255, 255, 255, 255})
		entity.AddComponent(realSprite)
		sprite = entity.GetSprite()
		if sprite == nil {
			t.Error("Expected sprite when correct component exists")
		}
	})

	t.Run("GetCollider", func(t *testing.T) {
		// Test with no component
		collider := entity.GetCollider()
		if collider != nil {
			t.Error("Expected nil collider when no component exists")
		}

		// Test with correct component type
		realCollider := components.NewCollider(32, 32, "test")
		entity.AddComponent(realCollider)
		collider = entity.GetCollider()
		if collider == nil {
			t.Error("Expected collider when correct component exists")
		}
	})

	t.Run("GetPhysics", func(t *testing.T) {
		// Test with no component
		physics := entity.GetPhysics()
		if physics != nil {
			t.Error("Expected nil physics when no component exists")
		}

		// Test with correct component type
		realPhysics := components.NewPhysics()
		entity.AddComponent(realPhysics)
		physics = entity.GetPhysics()
		if physics == nil {
			t.Error("Expected physics when correct component exists")
		}
	})

	t.Run("GetPacketType", func(t *testing.T) {
		// Test with no component
		packetType := entity.GetPacketType()
		if packetType != nil {
			t.Error("Expected nil packetType when no component exists")
		}

		// Test with correct component type
		realPacketType := components.NewPacketType("HTTP", 1)
		entity.AddComponent(realPacketType)
		packetType = entity.GetPacketType()
		if packetType == nil {
			t.Error("Expected packetType when correct component exists")
		}
	})

	t.Run("GetState", func(t *testing.T) {
		// Test with no component
		state := entity.GetState()
		if state != nil {
			t.Error("Expected nil state when no component exists")
		}

		// Test with correct component type
		realState := components.NewState(components.StatePlaying)
		entity.AddComponent(realState)
		state = entity.GetState()
		if state == nil {
			t.Error("Expected state when correct component exists")
		}
	})

	t.Run("GetCombo", func(t *testing.T) {
		// Test with no component
		combo := entity.GetCombo()
		if combo != nil {
			t.Error("Expected nil combo when no component exists")
		}

		// Test with correct component type
		realCombo := components.NewCombo()
		entity.AddComponent(realCombo)
		combo = entity.GetCombo()
		if combo == nil {
			t.Error("Expected combo when correct component exists")
		}
	})

	t.Run("GetSLA", func(t *testing.T) {
		// Test with no component
		sla := entity.GetSLA()
		if sla != nil {
			t.Error("Expected nil sla when no component exists")
		}

		// Test with correct component type
		realSLA := components.NewSLA(99.5, 10)
		entity.AddComponent(realSLA)
		sla = entity.GetSLA()
		if sla == nil {
			t.Error("Expected sla when correct component exists")
		}
	})

	t.Run("GetBackendAssignment", func(t *testing.T) {
		// Test with no component
		backend := entity.GetBackendAssignment()
		if backend != nil {
			t.Error("Expected nil backend when no component exists")
		}

		// Test with correct component type
		realBackend := components.NewBackendAssignment(1)
		entity.AddComponent(realBackend)
		backend = entity.GetBackendAssignment()
		if backend == nil {
			t.Error("Expected backend when correct component exists")
		}
	})

	t.Run("GetPowerUpType", func(t *testing.T) {
		// Test with no component
		powerUp := entity.GetPowerUpType()
		if powerUp != nil {
			t.Error("Expected nil powerUp when no component exists")
		}

		// Test with correct component type
		realPowerUp := components.NewPowerUpType("Shield", 5.0)
		entity.AddComponent(realPowerUp)
		powerUp = entity.GetPowerUpType()
		if powerUp == nil {
			t.Error("Expected powerUp when correct component exists")
		}
	})
}

// TestEntity_Concurrency tests thread safety of entity operations
func TestEntity_Concurrency(t *testing.T) {
	entity := NewEntity(1)
	done := make(chan bool, 10)

	// Test concurrent AddComponent operations
	for i := 0; i < 5; i++ {
		go func(id int) {
			component := &testMockComponent{
				componentType: fmt.Sprintf("Component%d", id),
			}
			entity.AddComponent(component)
			done <- true
		}(i)
	}

	// Test concurrent GetComponent operations
	for i := 0; i < 5; i++ {
		go func(id int) {
			entity.GetComponent(fmt.Sprintf("Component%d", id))
			entity.HasComponent(fmt.Sprintf("Component%d", id))
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify final state
	names := entity.GetComponentNames()
	if len(names) != 5 {
		t.Errorf("Expected 5 components after concurrent operations, got %d", len(names))
	}
}

// TestEntity_EdgeCases tests edge cases and error conditions
func TestEntity_EdgeCases(t *testing.T) {
	t.Run("Nil Component", func(t *testing.T) {
		entity := NewEntity(1)

		// This should not panic
		entity.AddComponent(nil)

		if len(entity.Components) != 0 {
			t.Error("Expected no components when adding nil component")
		}
	})

	t.Run("Empty Component Type", func(t *testing.T) {
		entity := NewEntity(1)

		component := &testMockComponent{componentType: ""}
		entity.AddComponent(component)

		if len(entity.Components) != 1 {
			t.Error("Expected 1 component when adding component with empty type")
		}

		if entity.Components[""] != component {
			t.Error("Expected component to be stored with empty key")
		}
	})

	t.Run("Zero ID", func(t *testing.T) {
		entity := NewEntity(0)

		if entity.ID != 0 {
			t.Errorf("Expected entity ID 0, got %d", entity.ID)
		}
	})

	t.Run("Large ID", func(t *testing.T) {
		largeID := uint64(18446744073709551615) // Max uint64
		entity := NewEntity(largeID)

		if entity.ID != largeID {
			t.Errorf("Expected entity ID %d, got %d", largeID, entity.ID)
		}
	})
}

// TestEntity_Integration tests integration scenarios
func TestEntity_Integration(t *testing.T) {
	entity := NewEntity(1)

	// Add multiple components of different types
	transform := components.NewTransform(100, 200)
	sprite := components.NewSprite(64, 64, color.RGBA{255, 255, 255, 255})
	collider := components.NewCollider(64, 64, "player")
	physics := components.NewPhysics()
	packetType := components.NewPacketType("HTTP", 1)
	state := components.NewState(components.StatePlaying)
	combo := components.NewCombo()
	sla := components.NewSLA(99.5, 10)
	backend := components.NewBackendAssignment(1)
	powerUp := components.NewPowerUpType("Shield", 5.0)

	entity.AddComponent(transform)
	entity.AddComponent(sprite)
	entity.AddComponent(collider)
	entity.AddComponent(physics)
	entity.AddComponent(packetType)
	entity.AddComponent(state)
	entity.AddComponent(combo)
	entity.AddComponent(sla)
	entity.AddComponent(backend)
	entity.AddComponent(powerUp)

	// Verify all components are present
	expectedComponents := []string{
		"Transform", "Sprite", "Collider", "Physics", "PacketType",
		"State", "Combo", "SLA", "BackendAssignment", "PowerUpType",
	}

	for _, componentName := range expectedComponents {
		if !entity.HasComponent(componentName) {
			t.Errorf("Expected component %s to be present", componentName)
		}
	}

	// Test type-safe getters
	if entity.GetTransform() == nil {
		t.Error("Expected GetTransform to return component")
	}
	if entity.GetSprite() == nil {
		t.Error("Expected GetSprite to return component")
	}
	if entity.GetCollider() == nil {
		t.Error("Expected GetCollider to return component")
	}
	if entity.GetPhysics() == nil {
		t.Error("Expected GetPhysics to return component")
	}
	if entity.GetPacketType() == nil {
		t.Error("Expected GetPacketType to return component")
	}
	if entity.GetState() == nil {
		t.Error("Expected GetState to return component")
	}
	if entity.GetCombo() == nil {
		t.Error("Expected GetCombo to return component")
	}
	if entity.GetSLA() == nil {
		t.Error("Expected GetSLA to return component")
	}
	if entity.GetBackendAssignment() == nil {
		t.Error("Expected GetBackendAssignment to return component")
	}
	if entity.GetPowerUpType() == nil {
		t.Error("Expected GetPowerUpType to return component")
	}

	// Test component removal
	entity.RemoveComponent("Transform")
	if entity.HasComponent("Transform") {
		t.Error("Expected Transform component to be removed")
	}

	// Test state changes
	entity.SetActive(false)
	if entity.IsActive() {
		t.Error("Expected entity to be inactive")
	}

	entity.SetActive(true)
	if !entity.IsActive() {
		t.Error("Expected entity to be active")
	}

	// Verify final component count
	names := entity.GetComponentNames()
	if len(names) != 9 { // 10 - 1 removed
		t.Errorf("Expected 9 components after removal, got %d", len(names))
	}
}

// Benchmark tests for performance
func BenchmarkNewEntity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewEntity(uint64(i))
	}
}

func BenchmarkEntity_AddComponent(b *testing.B) {
	entity := NewEntity(1)
	component := &testMockComponent{componentType: "TestComponent"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entity.AddComponent(component)
	}
}

func BenchmarkEntity_GetComponent(b *testing.B) {
	entity := NewEntity(1)
	component := &testMockComponent{componentType: "TestComponent"}
	entity.AddComponent(component)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entity.GetComponent("TestComponent")
	}
}

func BenchmarkEntity_HasComponent(b *testing.B) {
	entity := NewEntity(1)
	component := &testMockComponent{componentType: "TestComponent"}
	entity.AddComponent(component)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		entity.HasComponent("TestComponent")
	}
}
