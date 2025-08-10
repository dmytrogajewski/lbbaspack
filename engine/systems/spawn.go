package systems

import (
	"fmt"
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
	"math/rand"
)

const SystemTypeSpawn SystemType = "spawn"

// SpawnSystem manages the spawning of packets and power-ups in the game.
// It handles spawn timing, level progression, DDoS attacks, and entity creation.
type SpawnSystem struct {
	BaseSystem
	spawnCallback func() Entity
}

// NewSpawnSystem creates a new spawn system with default configuration.
// The spawnCallback function is used to create new entities when spawning is needed.
func NewSpawnSystem(spawnCallback func() Entity) *SpawnSystem {
	return &SpawnSystem{
		BaseSystem:    BaseSystem{RequiredComponents: []string{"Spawner"}},
		spawnCallback: spawnCallback,
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (ss *SpawnSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeSpawn,
		System:       ss,
		Dependencies: []SystemType{},
		Conflicts:    []SystemType{},
		Provides:     []string{"entity_spawning", "packet_generation"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     false,
	}
}

// IncreasePacketSpeed increases the packet speed by a percentage.
func (ss *SpawnSystem) IncreasePacketSpeed(percent float64) {}

// IncreaseLevel increases the level and adjusts spawn rate for higher density.
// Higher levels result in faster packet spawning to increase difficulty.
func (ss *SpawnSystem) IncreaseLevel(newLevel int) {
	// No-op; levels should be set on Spawner component externally
}

// Initialize sets up event listeners for the spawn system.
func (ss *SpawnSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	eventDispatcher.Subscribe(events.EventLevelUp, func(event *events.Event) {
		// Level up events are handled in Update via Spawner component
	})
}

// Update processes the spawn system for the given delta time.
// This includes DDoS attack management, timer updates, and entity spawning.
func (ss *SpawnSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Find singleton Spawner component holder; assume entity 1 has it or scan
	var spawner *components.Spawner
	for _, e := range entities {
		if comp := e.GetComponentByName("Spawner"); comp != nil {
			if s, ok := comp.(*components.Spawner); ok {
				spawner = s
				break
			}
		}
	}
	if spawner == nil {
		return
	}

	// DDoS state update
	if spawner.IsDDoSActive {
		spawner.DDOSTimer += deltaTime
		if spawner.DDOSTimer >= spawner.DDoSDuration {
			spawner.IsDDoSActive = false
			spawner.DDOSTimer = 0
			eventDispatcher.Publish(events.NewEvent(events.EventDDoSEnd, nil))
		}
	} else {
		spawner.DDoSCooldown -= deltaTime
		if spawner.DDoSCooldown <= 0 {
			if rand.Float64() < 0.01 {
				spawner.IsDDoSActive = true
				spawner.DDOSTimer = 0
				spawner.DDoSCooldown = 20.0 + rand.Float64()*20.0
				eventDispatcher.Publish(events.NewEvent(events.EventDDoSStart, &events.EventData{Duration: &spawner.DDoSDuration, Level: &spawner.Level}))
			}
		}
	}

	// Timers
	spawner.PacketSpawnElapsed += deltaTime
	spawner.PowerUpSpawnElapsed += deltaTime

	// Spawn packet (more frequent during DDoS)
	effectivePacketSpawnRate := spawner.PacketSpawnRate
	if spawner.IsDDoSActive && spawner.DDoSMult > 1 {
		effectivePacketSpawnRate = spawner.PacketSpawnRate / (spawner.DDoSMult * 1.5)
		if effectivePacketSpawnRate < 0.05 {
			effectivePacketSpawnRate = 0.05
		}
	}
	if spawner.PacketSpawnElapsed >= effectivePacketSpawnRate {
		spawner.PacketSpawnElapsed = 0
		ss.spawnPacketWithConfig(spawner)
	}

	// Spawn powerup
	if spawner.PowerUpSpawnElapsed >= spawner.PowerUpSpawnRate {
		spawner.PowerUpSpawnElapsed = 0
		ss.spawnPowerUp()
	}
}

func (ss *SpawnSystem) spawnPacketWithConfig(spawner *components.Spawner) {
	if ss.spawnCallback == nil {
		return
	}
	packet := ss.spawnCallback()
	entity, ok := packet.(interface {
		AddComponent(components.Component)
		GetComponentNames() []string
	})
	if !ok {
		return
	}
	x := float64(rand.Intn(800 - 15))
	y := -15.0
	entity.AddComponent(components.NewTransform(x, y))
	entity.AddComponent(components.NewSprite(15, 15, components.RandomPacketColor()))
	entity.AddComponent(components.NewCollider(15, 15, "packet"))
	physics := components.NewPhysics()
	speed := spawner.PacketSpeed
	// During DDoS, packets fall a bit slower than normal
	if spawner.IsDDoSActive && spawner.DDoSMult > 1 {
		speed = spawner.PacketSpeed * 0.7
	}
	physics.SetVelocity(0, speed)
	entity.AddComponent(physics)
	entity.AddComponent(components.NewPacketType(components.RandomPacketName(), 10))
}

// logPacketSpawn logs information about the spawned packet for debugging.
// logging helpers removed

// logComponentInfo logs component information for debugging.
func (ss *SpawnSystem) logComponentInfo(entity interface{ GetComponentNames() []string }, entityType string) {
	componentNames := entity.GetComponentNames()
	fmt.Printf("[SpawnSystem] Spawned %s with components: %v\n", entityType, componentNames)

	hasTransform := false
	hasSprite := false
	for _, name := range componentNames {
		if name == "Transform" {
			hasTransform = true
		}
		if name == "Sprite" {
			hasSprite = true
		}
	}
	fmt.Printf("[SpawnSystem] %s has Transform: %v, Sprite: %v\n", entityType, hasTransform, hasSprite)
}

// Test-only legacy helpers are defined in legacy_test_shims_test.go

// spawnPowerUp creates a new power-up entity with all required components.
func (ss *SpawnSystem) spawnPowerUp() {
	if ss.spawnCallback == nil {
		fmt.Println("[SpawnSystem] Spawn callback is nil!")
		return
	}

	fmt.Println("[SpawnSystem] Spawn callback is not nil, spawning powerup...")
	powerup := ss.spawnCallback()

	// Entity interface for adding components and getting component names
	entity, ok := powerup.(interface {
		AddComponent(components.Component)
		GetComponentNames() []string
	})
	if !ok {
		fmt.Println("[SpawnSystem] Failed to cast spawned powerup to interface")
		return
	}

	ss.addPowerUpComponents(entity)
	ss.logPowerUpSpawn(entity)
}

// addPowerUpComponents adds all required components to a power-up entity.
func (ss *SpawnSystem) addPowerUpComponents(entity interface {
	AddComponent(components.Component)
}) {
	name, col := randomPowerUpNameAndColor()
	entity.AddComponent(components.NewTransform(float64(rand.Intn(800-15)), -15))
	entity.AddComponent(components.NewSprite(15, 15, col))
	entity.AddComponent(components.NewCollider(15, 15, "powerup"))

	physics := components.NewPhysics()
	physics.SetVelocity(0, 50)
	entity.AddComponent(physics)

	entity.AddComponent(components.NewPowerUpType(name, 10.0))
}

// logPowerUpSpawn logs information about the spawned power-up for debugging.
func (ss *SpawnSystem) logPowerUpSpawn(entity interface{ GetComponentNames() []string }) {
	ss.logComponentInfo(entity, "powerup")
}

// randomPowerUpNameAndColor returns a random power-up name and color.
func randomPowerUpNameAndColor() (string, color.RGBA) {
	types := []struct {
		name string
		col  color.RGBA
	}{
		{"Speed Boost", color.RGBA{255, 255, 0, 255}},
		{"Wide Catch", color.RGBA{0, 255, 255, 255}},
		{"Multi-Catch", color.RGBA{255, 0, 255, 255}},
		{"Time Slow", color.RGBA{0, 0, 255, 255}},
		{"Shield", color.RGBA{0, 255, 0, 255}},
		{"Auto-Balancer", color.RGBA{255, 165, 0, 255}},
	}
	idx := rand.Intn(len(types))
	return types[idx].name, types[idx].col
}
