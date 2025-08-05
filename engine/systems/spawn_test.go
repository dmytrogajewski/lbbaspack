package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"math"
	"testing"
)

// spawnTestEntity is a test entity for spawn system testing
type spawnTestEntity struct {
	entity          *entities.Entity
	componentsAdded []string
}

func newSpawnTestEntity(id uint64) *spawnTestEntity {
	return &spawnTestEntity{
		entity:          entities.NewEntity(id),
		componentsAdded: make([]string, 0),
	}
}

func (ste *spawnTestEntity) AddComponent(component components.Component) {
	ste.componentsAdded = append(ste.componentsAdded, component.GetType())
	ste.entity.AddComponent(component)
}

func (ste *spawnTestEntity) GetComponentNames() []string {
	return ste.componentsAdded
}

func (ste *spawnTestEntity) GetID() uint64 {
	return ste.entity.GetID()
}

// Implement Entity interface methods
func (ste *spawnTestEntity) GetComponent(componentType string) components.Component {
	return ste.entity.GetComponent(componentType)
}

func (ste *spawnTestEntity) GetComponentByName(typeName string) components.Component {
	return ste.entity.GetComponentByName(typeName)
}

func (ste *spawnTestEntity) HasComponent(componentType string) bool {
	return ste.entity.HasComponent(componentType)
}

func (ste *spawnTestEntity) RemoveComponent(componentType string) {
	ste.entity.RemoveComponent(componentType)
}

func (ste *spawnTestEntity) SetActive(active bool) {
	ste.entity.SetActive(active)
}

func (ste *spawnTestEntity) IsActive() bool {
	return ste.entity.IsActive()
}

func (ste *spawnTestEntity) GetTransform() components.TransformComponent {
	return ste.entity.GetTransform()
}

func (ste *spawnTestEntity) GetSprite() components.SpriteComponent {
	return ste.entity.GetSprite()
}

func (ste *spawnTestEntity) GetCollider() components.ColliderComponent {
	return ste.entity.GetCollider()
}

func (ste *spawnTestEntity) GetPhysics() components.PhysicsComponent {
	return ste.entity.GetPhysics()
}

func (ste *spawnTestEntity) GetPacketType() components.PacketTypeComponent {
	return ste.entity.GetPacketType()
}

func (ste *spawnTestEntity) GetState() components.StateComponent {
	return ste.entity.GetState()
}

func (ste *spawnTestEntity) GetCombo() components.ComboComponent {
	return ste.entity.GetCombo()
}

func (ste *spawnTestEntity) GetSLA() components.SLAComponent {
	return ste.entity.GetSLA()
}

func (ste *spawnTestEntity) GetBackendAssignment() components.BackendAssignmentComponent {
	return ste.entity.GetBackendAssignment()
}

func (ste *spawnTestEntity) GetPowerUpType() components.PowerUpTypeComponent {
	return ste.entity.GetPowerUpType()
}

func (ste *spawnTestEntity) GetRouting() components.RoutingComponent {
	return ste.entity.GetRouting()
}

func TestNewSpawnSystem(t *testing.T) {
	spawnCallback := func() Entity {
		return newSpawnTestEntity(1)
	}

	ss := NewSpawnSystem(spawnCallback)

	// Test that the system is properly initialized
	if ss == nil {
		t.Fatal("NewSpawnSystem returned nil")
	}

	// Test initial values
	if ss.lastPacketSpawn != 0 {
		t.Errorf("Expected lastPacketSpawn to be 0, got %f", ss.lastPacketSpawn)
	}

	if ss.packetSpawnRate != 1.0 {
		t.Errorf("Expected packetSpawnRate to be 1.0, got %f", ss.packetSpawnRate)
	}

	if ss.lastPowerUpSpawn != 0 {
		t.Errorf("Expected lastPowerUpSpawn to be 0, got %f", ss.lastPowerUpSpawn)
	}

	if ss.powerUpSpawnRate != 10.0 {
		t.Errorf("Expected powerUpSpawnRate to be 10.0, got %f", ss.powerUpSpawnRate)
	}

	if ss.packetSpeed != 100 {
		t.Errorf("Expected packetSpeed to be 100, got %f", ss.packetSpeed)
	}

	if ss.level != 1 {
		t.Errorf("Expected level to be 1, got %d", ss.level)
	}

	// Test DDoS attack state
	if ss.isDDoSActive {
		t.Error("Expected isDDoSActive to be false initially")
	}

	if ss.ddosTimer != 0 {
		t.Errorf("Expected ddosTimer to be 0, got %f", ss.ddosTimer)
	}

	if ss.ddosDuration != 5.0 {
		t.Errorf("Expected ddosDuration to be 5.0, got %f", ss.ddosDuration)
	}

	if ss.ddosMultiplier != 10.0 {
		t.Errorf("Expected ddosMultiplier to be 10.0, got %f", ss.ddosMultiplier)
	}

	if ss.ddosCooldown != 10.0 {
		t.Errorf("Expected ddosCooldown to be 10.0, got %f", ss.ddosCooldown)
	}

	// Test spawn callback
	if ss.spawnCallback == nil {
		t.Error("Expected spawnCallback to be set")
	}
}

func TestSpawnSystem_IncreasePacketSpeed(t *testing.T) {
	ss := NewSpawnSystem(func() Entity { return newSpawnTestEntity(1) })

	initialSpeed := ss.packetSpeed

	// Test 20% increase
	ss.IncreasePacketSpeed(20.0)
	expectedSpeed := initialSpeed * 1.2
	if math.Abs(ss.packetSpeed-expectedSpeed) > 0.0001 {
		t.Errorf("Expected packet speed to be %f after 20%% increase, got %f", expectedSpeed, ss.packetSpeed)
	}

	// Test 50% increase
	ss.IncreasePacketSpeed(50.0)
	expectedSpeed = expectedSpeed * 1.5
	if math.Abs(ss.packetSpeed-expectedSpeed) > 0.0001 {
		t.Errorf("Expected packet speed to be %f after 50%% increase, got %f", expectedSpeed, ss.packetSpeed)
	}

	// Test negative increase (decrease)
	ss.IncreasePacketSpeed(-25.0)
	expectedSpeed = expectedSpeed * 0.75
	if math.Abs(ss.packetSpeed-expectedSpeed) > 0.0001 {
		t.Errorf("Expected packet speed to be %f after -25%% increase, got %f", expectedSpeed, ss.packetSpeed)
	}
}

func TestSpawnSystem_IncreaseLevel(t *testing.T) {
	ss := NewSpawnSystem(func() Entity { return newSpawnTestEntity(1) })

	// Test level 1 (initial)
	if ss.level != 1 {
		t.Errorf("Expected initial level to be 1, got %d", ss.level)
	}
	if math.Abs(ss.packetSpawnRate-1.0) > 0.0001 {
		t.Errorf("Expected initial spawn rate to be 1.0, got %f", ss.packetSpawnRate)
	}

	// Test level 2
	ss.IncreaseLevel(2)
	if ss.level != 2 {
		t.Errorf("Expected level to be 2, got %d", ss.level)
	}
	expectedRate := 1.0 / (1.0 + float64(2-1)*0.2) // 1.0 / 1.2 = 0.833...
	if math.Abs(ss.packetSpawnRate-expectedRate) > 0.0001 {
		t.Errorf("Expected spawn rate to be %f for level 2, got %f", expectedRate, ss.packetSpawnRate)
	}

	// Test level 5
	ss.IncreaseLevel(5)
	if ss.level != 5 {
		t.Errorf("Expected level to be 5, got %d", ss.level)
	}
	expectedRate = 1.0 / (1.0 + float64(5-1)*0.2) // 1.0 / 1.8 = 0.556...
	if math.Abs(ss.packetSpawnRate-expectedRate) > 0.0001 {
		t.Errorf("Expected spawn rate to be %f for level 5, got %f", expectedRate, ss.packetSpawnRate)
	}

	// Test level 10
	ss.IncreaseLevel(10)
	if ss.level != 10 {
		t.Errorf("Expected level to be 10, got %d", ss.level)
	}
	expectedRate = 1.0 / (1.0 + float64(10-1)*0.2) // 1.0 / 2.8 = 0.357...
	if math.Abs(ss.packetSpawnRate-expectedRate) > 0.0001 {
		t.Errorf("Expected spawn rate to be %f for level 10, got %f", expectedRate, ss.packetSpawnRate)
	}
}

func TestSpawnSystem_Initialize(t *testing.T) {
	ss := NewSpawnSystem(func() Entity { return newSpawnTestEntity(1) })
	eventDispatcher := events.NewEventDispatcher()

	// Initialize the system
	ss.Initialize(eventDispatcher)

	// Test that level-up event handler is registered
	// We can't directly test the subscription, but we can test that the system
	// responds to level-up events
	level := 3
	event := events.NewEvent(events.EventLevelUp, &events.EventData{Level: &level})
	eventDispatcher.Publish(event)

	// The event should trigger the level increase
	if ss.level != 3 {
		t.Errorf("Expected level to be updated to 3 after level-up event, got %d", ss.level)
	}
}

func TestSpawnSystem_Update_NoSpawn(t *testing.T) {
	spawnCount := 0
	spawnCallback := func() Entity {
		spawnCount++
		return newSpawnTestEntity(uint64(spawnCount))
	}

	ss := NewSpawnSystem(spawnCallback)
	eventDispatcher := events.NewEventDispatcher()

	// Update with small delta time (should not trigger spawn)
	ss.Update(0.5, []Entity{}, eventDispatcher)

	if spawnCount != 0 {
		t.Errorf("Expected no spawns with 0.5s delta time, got %d", spawnCount)
	}

	if ss.lastPacketSpawn != 0.5 {
		t.Errorf("Expected lastPacketSpawn to be 0.5, got %f", ss.lastPacketSpawn)
	}
}

func TestSpawnSystem_Update_PacketSpawn(t *testing.T) {
	spawnCount := 0
	spawnedEntities := make([]Entity, 0)
	spawnCallback := func() Entity {
		spawnCount++
		entity := newSpawnTestEntity(uint64(spawnCount))
		spawnedEntities = append(spawnedEntities, entity)
		return entity
	}

	ss := NewSpawnSystem(spawnCallback)
	eventDispatcher := events.NewEventDispatcher()

	// Update with enough time to trigger packet spawn
	ss.Update(1.1, []Entity{}, eventDispatcher)

	if spawnCount != 1 {
		t.Errorf("Expected 1 spawn with 1.1s delta time, got %d", spawnCount)
	}

	if ss.lastPacketSpawn != 0 {
		t.Errorf("Expected lastPacketSpawn to be reset to 0, got %f", ss.lastPacketSpawn)
	}

	// Check that the spawned entity has the required components
	if len(spawnedEntities) == 0 {
		t.Fatal("No entities were spawned")
	}

	entity := spawnedEntities[0]

	// Check required components exist
	if !entity.HasComponent("Transform") {
		t.Error("Spawned entity missing Transform component")
	}
	if !entity.HasComponent("Sprite") {
		t.Error("Spawned entity missing Sprite component")
	}
	if !entity.HasComponent("Collider") {
		t.Error("Spawned entity missing Collider component")
	}
	if !entity.HasComponent("Physics") {
		t.Error("Spawned entity missing Physics component")
	}
	if !entity.HasComponent("PacketType") {
		t.Error("Spawned entity missing PacketType component")
	}

	// Check that physics component has correct velocity
	physics := entity.GetPhysics()
	if physics == nil {
		t.Fatal("Physics component not found")
	}

	vx := physics.GetVelocityX()
	vy := physics.GetVelocityY()
	if vx != 0 {
		t.Errorf("Expected packet velocity X to be 0, got %f", vx)
	}
	if vy != ss.packetSpeed {
		t.Errorf("Expected packet velocity Y to be %f, got %f", ss.packetSpeed, vy)
	}
}

func TestSpawnSystem_Update_PowerUpSpawn(t *testing.T) {
	spawnCount := 0
	spawnedEntities := make([]Entity, 0)
	spawnCallback := func() Entity {
		spawnCount++
		entity := newSpawnTestEntity(uint64(spawnCount))
		spawnedEntities = append(spawnedEntities, entity)
		return entity
	}

	ss := NewSpawnSystem(spawnCallback)
	eventDispatcher := events.NewEventDispatcher()

	// Update with enough time to trigger powerup spawn
	ss.Update(10.1, []Entity{}, eventDispatcher)

	if spawnCount != 2 { // 1 packet + 1 powerup
		t.Errorf("Expected 2 spawns with 10.1s delta time, got %d", spawnCount)
	}

	// Check that the powerup entity has the required components
	if len(spawnedEntities) < 2 {
		t.Fatal("Not enough entities were spawned")
	}

	powerupEntity := spawnedEntities[1] // Second entity should be powerup

	// Check required components exist
	if !powerupEntity.HasComponent("Transform") {
		t.Error("Spawned powerup missing Transform component")
	}
	if !powerupEntity.HasComponent("Sprite") {
		t.Error("Spawned powerup missing Sprite component")
	}
	if !powerupEntity.HasComponent("Collider") {
		t.Error("Spawned powerup missing Collider component")
	}
	if !powerupEntity.HasComponent("Physics") {
		t.Error("Spawned powerup missing Physics component")
	}
	if !powerupEntity.HasComponent("PowerUpType") {
		t.Error("Spawned powerup missing PowerUpType component")
	}

	// Check that physics component has correct velocity (slower than packets)
	physics := powerupEntity.GetPhysics()
	if physics == nil {
		t.Fatal("Physics component not found on powerup")
	}

	vx := physics.GetVelocityX()
	vy := physics.GetVelocityY()
	if vx != 0 {
		t.Errorf("Expected powerup velocity X to be 0, got %f", vx)
	}
	if vy != 50 {
		t.Errorf("Expected powerup velocity Y to be 50, got %f", vy)
	}
}

func TestSpawnSystem_Update_DDoSAttack(t *testing.T) {
	spawnCount := 0
	spawnCallback := func() Entity {
		spawnCount++
		return newSpawnTestEntity(uint64(spawnCount))
	}

	ss := NewSpawnSystem(spawnCallback)
	eventDispatcher := events.NewEventDispatcher()

	// Manually trigger DDoS attack
	ss.isDDoSActive = true
	ss.ddosTimer = 0
	ss.ddosCooldown = 0

	// Update with time less than DDoS duration
	ss.Update(2.0, []Entity{}, eventDispatcher)

	if !ss.isDDoSActive {
		t.Error("Expected DDoS attack to still be active after 2.0s")
	}

	if ss.ddosTimer != 2.0 {
		t.Errorf("Expected ddosTimer to be 2.0, got %f", ss.ddosTimer)
	}

	// Update to end DDoS attack
	ss.Update(3.1, []Entity{}, eventDispatcher)

	if ss.isDDoSActive {
		t.Error("Expected DDoS attack to end after duration exceeded")
	}

	if ss.ddosTimer != 0 {
		t.Errorf("Expected ddosTimer to be reset to 0, got %f", ss.ddosTimer)
	}
}

func TestSpawnSystem_Update_DDoSAttackSpawnRate(t *testing.T) {
	spawnCount := 0
	spawnCallback := func() Entity {
		spawnCount++
		return newSpawnTestEntity(uint64(spawnCount))
	}

	ss := NewSpawnSystem(spawnCallback)
	eventDispatcher := events.NewEventDispatcher()

	// Set level to 2 for testing
	ss.IncreaseLevel(2)
	normalSpawnRate := ss.packetSpawnRate

	// Manually trigger DDoS attack and set the spawn rate
	ss.isDDoSActive = true
	ss.ddosTimer = 0
	ss.packetSpawnRate = normalSpawnRate / ss.ddosMultiplier

	// Update to trigger spawn during DDoS
	ss.Update(0.1, []Entity{}, eventDispatcher) // Should spawn much faster during DDoS

	if spawnCount == 0 {
		t.Error("Expected packet to spawn during DDoS attack")
	}

	// Check that spawn rate is reduced during DDoS
	expectedDDoSSpawnRate := normalSpawnRate / ss.ddosMultiplier
	if math.Abs(ss.packetSpawnRate-expectedDDoSSpawnRate) > 0.0001 {
		t.Errorf("Expected DDoS spawn rate to be %f, got %f", expectedDDoSSpawnRate, ss.packetSpawnRate)
	}
}

func TestSpawnSystem_Update_NilSpawnCallback(t *testing.T) {
	ss := NewSpawnSystem(nil)
	eventDispatcher := events.NewEventDispatcher()

	// This should not panic
	ss.Update(1.1, []Entity{}, eventDispatcher)
}

func TestSpawnSystem_Update_MultipleUpdates(t *testing.T) {
	spawnCount := 0
	spawnCallback := func() Entity {
		spawnCount++
		return newSpawnTestEntity(uint64(spawnCount))
	}

	ss := NewSpawnSystem(spawnCallback)
	eventDispatcher := events.NewEventDispatcher()

	// Multiple updates that should trigger multiple spawns
	ss.Update(1.1, []Entity{}, eventDispatcher) // First packet
	ss.Update(1.1, []Entity{}, eventDispatcher) // Second packet
	ss.Update(1.1, []Entity{}, eventDispatcher) // Third packet

	if spawnCount != 3 {
		t.Errorf("Expected 3 spawns after 3 updates, got %d", spawnCount)
	}
}

func TestSpawnSystem_Update_LevelProgression(t *testing.T) {
	spawnCount := 0
	spawnCallback := func() Entity {
		spawnCount++
		return newSpawnTestEntity(uint64(spawnCount))
	}

	ss := NewSpawnSystem(spawnCallback)
	eventDispatcher := events.NewEventDispatcher()

	// Increase level and test spawn rate
	ss.IncreaseLevel(3)
	expectedSpawnRate := 1.0 / (1.0 + float64(3-1)*0.2)

	// Update with new spawn rate
	ss.Update(expectedSpawnRate+0.1, []Entity{}, eventDispatcher)

	if spawnCount != 1 {
		t.Errorf("Expected 1 spawn with level 3 spawn rate, got %d", spawnCount)
	}
}

func TestSpawnSystem_Update_ZeroDeltaTime(t *testing.T) {
	spawnCount := 0
	spawnCallback := func() Entity {
		spawnCount++
		return newSpawnTestEntity(uint64(spawnCount))
	}

	ss := NewSpawnSystem(spawnCallback)
	eventDispatcher := events.NewEventDispatcher()

	// Update with zero delta time
	ss.Update(0.0, []Entity{}, eventDispatcher)

	if spawnCount != 0 {
		t.Errorf("Expected no spawns with zero delta time, got %d", spawnCount)
	}

	if ss.lastPacketSpawn != 0 {
		t.Errorf("Expected lastPacketSpawn to remain 0, got %f", ss.lastPacketSpawn)
	}
}

func TestSpawnSystem_Update_NegativeDeltaTime(t *testing.T) {
	spawnCount := 0
	spawnCallback := func() Entity {
		spawnCount++
		return newSpawnTestEntity(uint64(spawnCount))
	}

	ss := NewSpawnSystem(spawnCallback)
	eventDispatcher := events.NewEventDispatcher()

	// Update with negative delta time
	ss.Update(-1.0, []Entity{}, eventDispatcher)

	if spawnCount != 0 {
		t.Errorf("Expected no spawns with negative delta time, got %d", spawnCount)
	}

	if ss.lastPacketSpawn != -1.0 {
		t.Errorf("Expected lastPacketSpawn to be -1.0, got %f", ss.lastPacketSpawn)
	}
}

func TestSpawnSystem_Update_LargeDeltaTime(t *testing.T) {
	spawnCount := 0
	spawnCallback := func() Entity {
		spawnCount++
		return newSpawnTestEntity(uint64(spawnCount))
	}

	ss := NewSpawnSystem(spawnCallback)
	eventDispatcher := events.NewEventDispatcher()

	// Update with very large delta time
	ss.Update(100.0, []Entity{}, eventDispatcher)

	// The spawn system only spawns one packet and one powerup per update call
	// With 100 seconds, it should spawn at least 1 packet and 1 powerup
	if spawnCount < 2 {
		t.Errorf("Expected at least 2 spawns (1 packet + 1 powerup) with 100s delta time, got %d", spawnCount)
	}
}

func TestRandomPowerUpNameAndColor(t *testing.T) {
	// Test that the function returns valid powerup types
	validNames := []string{"Speed Boost", "Wide Catch", "Multi-Catch", "Time Slow", "Shield", "Auto-Balancer"}
	validColors := []struct {
		r, g, b, a uint8
	}{
		{255, 255, 0, 255}, // Yellow
		{0, 255, 255, 255}, // Cyan
		{255, 0, 255, 255}, // Magenta
		{0, 0, 255, 255},   // Blue
		{0, 255, 0, 255},   // Green
		{255, 165, 0, 255}, // Orange
	}

	// Test multiple calls to ensure randomness
	for i := 0; i < 10; i++ {
		name, color := randomPowerUpNameAndColor()

		// Check name is valid
		nameValid := false
		for _, validName := range validNames {
			if name == validName {
				nameValid = true
				break
			}
		}
		if !nameValid {
			t.Errorf("Invalid powerup name returned: %s", name)
		}

		// Check color is valid
		colorValid := false
		for _, validColor := range validColors {
			if color.R == validColor.r && color.G == validColor.g &&
				color.B == validColor.b && color.A == validColor.a {
				colorValid = true
				break
			}
		}
		if !colorValid {
			t.Errorf("Invalid powerup color returned: R=%d, G=%d, B=%d, A=%d",
				color.R, color.G, color.B, color.A)
		}
	}
}

func TestSpawnSystem_Integration(t *testing.T) {
	spawnCount := 0
	spawnedEntities := make([]Entity, 0)
	spawnCallback := func() Entity {
		spawnCount++
		entity := newSpawnTestEntity(uint64(spawnCount))
		spawnedEntities = append(spawnedEntities, entity)
		return entity
	}

	ss := NewSpawnSystem(spawnCallback)
	eventDispatcher := events.NewEventDispatcher()

	// Initialize system
	ss.Initialize(eventDispatcher)

	// Simulate game progression
	ss.IncreaseLevel(2)
	ss.IncreasePacketSpeed(25.0)

	// Update multiple times to test full integration
	ss.Update(1.1, []Entity{}, eventDispatcher)  // First packet
	ss.Update(1.1, []Entity{}, eventDispatcher)  // Second packet
	ss.Update(10.1, []Entity{}, eventDispatcher) // Third packet + powerup

	if spawnCount != 4 {
		t.Errorf("Expected 4 spawns in integration test (3 packets + 1 powerup), got %d", spawnCount)
	}

	// Verify packet speed was increased
	if ss.packetSpeed != 125.0 {
		t.Errorf("Expected packet speed to be 125.0 after 25%% increase, got %f", ss.packetSpeed)
	}

	// Verify level was increased
	if ss.level != 2 {
		t.Errorf("Expected level to be 2, got %d", ss.level)
	}

	// Verify spawn rate was adjusted for level
	expectedSpawnRate := 1.0 / (1.0 + float64(2-1)*0.2)
	if math.Abs(ss.packetSpawnRate-expectedSpawnRate) > 0.0001 {
		t.Errorf("Expected spawn rate to be %f for level 2, got %f", expectedSpawnRate, ss.packetSpawnRate)
	}
}
