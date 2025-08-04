package systems

import (
	"fmt"
	"image/color"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
	"math/rand"
)

// SpawnSystem manages the spawning of packets and power-ups in the game.
// It handles spawn timing, level progression, DDoS attacks, and entity creation.
type SpawnSystem struct {
	BaseSystem
	lastPacketSpawn  float64
	packetSpawnRate  float64
	lastPowerUpSpawn float64
	powerUpSpawnRate float64
	spawnCallback    func() Entity
	packetSpeed      float64
	level            int

	// DDoS attack state
	isDDoSActive   bool
	ddosTimer      float64
	ddosDuration   float64
	ddosMultiplier float64
	ddosCooldown   float64
}

// NewSpawnSystem creates a new spawn system with default configuration.
// The spawnCallback function is used to create new entities when spawning is needed.
func NewSpawnSystem(spawnCallback func() Entity) *SpawnSystem {
	return &SpawnSystem{
		BaseSystem:       BaseSystem{},
		lastPacketSpawn:  0,
		packetSpawnRate:  1.0,
		lastPowerUpSpawn: 0,
		powerUpSpawnRate: 10.0,
		spawnCallback:    spawnCallback,
		packetSpeed:      100,
		level:            1,
		isDDoSActive:     false,
		ddosTimer:        0,
		ddosDuration:     5.0,
		ddosMultiplier:   10.0,
		ddosCooldown:     10.0,
	}
}

// IncreasePacketSpeed increases the packet speed by a percentage.
func (ss *SpawnSystem) IncreasePacketSpeed(percent float64) {
	ss.packetSpeed *= (1.0 + percent/100.0)
	fmt.Printf("[SpawnSystem] Packet speed increased to %.2f\n", ss.packetSpeed)
}

// IncreaseLevel increases the level and adjusts spawn rate for higher density.
// Higher levels result in faster packet spawning to increase difficulty.
func (ss *SpawnSystem) IncreaseLevel(newLevel int) {
	ss.level = newLevel
	// Formula: base rate / (1 + level * 0.2) - increases density by ~20% per level
	ss.packetSpawnRate = 1.0 / (1.0 + float64(ss.level-1)*0.2)
	fmt.Printf("[SpawnSystem] Level increased to %d, spawn rate: %.3f seconds\n", ss.level, ss.packetSpawnRate)
}

// Initialize sets up event listeners for the spawn system.
func (ss *SpawnSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	eventDispatcher.Subscribe(events.EventLevelUp, func(event *events.Event) {
		if event.Data.Level != nil {
			ss.IncreaseLevel(*event.Data.Level)
		}
	})
}

// Update processes the spawn system for the given delta time.
// This includes DDoS attack management, timer updates, and entity spawning.
func (ss *SpawnSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	fmt.Printf("[SpawnSystem] Update called with deltaTime: %.3f, entities count: %d\n", deltaTime, len(entities))

	ss.updateDDoSAttack(deltaTime, eventDispatcher)
	ss.updateTimers(deltaTime)
	ss.trySpawnPacket(eventDispatcher)
	ss.trySpawnPowerUp(eventDispatcher)
}

// updateDDoSAttack manages DDoS attack state and timing.
func (ss *SpawnSystem) updateDDoSAttack(deltaTime float64, eventDispatcher *events.EventDispatcher) {
	if ss.isDDoSActive {
		ss.ddosTimer += deltaTime
		if ss.ddosTimer >= ss.ddosDuration {
			ss.endDDoSAttack(eventDispatcher)
		}
	} else {
		ss.ddosCooldown -= deltaTime
		if ss.ddosCooldown <= 0 {
			ss.tryStartDDoSAttack(eventDispatcher)
		}
	}
}

// endDDoSAttack terminates the current DDoS attack and restores normal spawn rates.
func (ss *SpawnSystem) endDDoSAttack(eventDispatcher *events.EventDispatcher) {
	ss.isDDoSActive = false
	ss.ddosTimer = 0
	ss.packetSpawnRate = 1.0 / (1.0 + float64(ss.level-1)*0.2)
	fmt.Println("[SpawnSystem] DDoS attack ended. Spawn rate restored.")
	eventDispatcher.Publish(events.NewEvent(events.EventDDoSEnd, nil))
}

// tryStartDDoSAttack attempts to start a DDoS attack with a random chance.
func (ss *SpawnSystem) tryStartDDoSAttack(eventDispatcher *events.EventDispatcher) {
	if rand.Float64() < 0.01 {
		ss.startDDoSAttack(eventDispatcher)
	}
}

// startDDoSAttack initiates a DDoS attack and publishes the start event.
func (ss *SpawnSystem) startDDoSAttack(eventDispatcher *events.EventDispatcher) {
	ss.isDDoSActive = true
	ss.ddosTimer = 0
	ss.ddosCooldown = 20.0 + rand.Float64()*20.0
	ss.packetSpawnRate = (1.0 / (1.0 + float64(ss.level-1)*0.2)) / ss.ddosMultiplier
	fmt.Println("[SpawnSystem] DDoS attack started! Spawn rate massively increased.")
	eventDispatcher.Publish(events.NewEvent(events.EventDDoSStart, &events.EventData{
		Duration: &ss.ddosDuration,
		Level:    &ss.level,
	}))
}

// updateTimers increments the spawn timers by the given delta time.
func (ss *SpawnSystem) updateTimers(deltaTime float64) {
	ss.lastPacketSpawn += deltaTime
	ss.lastPowerUpSpawn += deltaTime
	fmt.Printf("[SpawnSystem] lastPacketSpawn: %.3f, packetSpawnRate: %.3f\n", ss.lastPacketSpawn, ss.packetSpawnRate)
}

// trySpawnPacket attempts to spawn a packet if enough time has passed.
func (ss *SpawnSystem) trySpawnPacket(eventDispatcher *events.EventDispatcher) {
	if ss.lastPacketSpawn >= ss.packetSpawnRate {
		fmt.Printf("[SpawnSystem] SPAWNING PACKET! lastPacketSpawn: %.3f >= %.3f\n", ss.lastPacketSpawn, ss.packetSpawnRate)
		ss.lastPacketSpawn = 0
		ss.spawnPacket()
	} else {
		fmt.Printf("[SpawnSystem] Not spawning yet, need %.3f more time\n", ss.packetSpawnRate-ss.lastPacketSpawn)
	}
}

// spawnPacket creates a new packet entity with all required components.
func (ss *SpawnSystem) spawnPacket() {
	if ss.spawnCallback == nil {
		fmt.Println("[SpawnSystem] Spawn callback is nil!")
		return
	}

	fmt.Println("[SpawnSystem] Spawn callback is not nil, spawning packet...")
	packet := ss.spawnCallback()

	// Entity interface for adding components and getting component names
	entity, ok := packet.(interface {
		AddComponent(components.Component)
		GetComponentNames() []string
	})
	if !ok {
		fmt.Println("[SpawnSystem] Failed to cast spawned entity to interface")
		return
	}

	ss.addPacketComponents(entity)
	ss.logPacketSpawn(packet, entity)
}

// addPacketComponents adds all required components to a packet entity.
func (ss *SpawnSystem) addPacketComponents(entity interface {
	AddComponent(components.Component)
}) {
	x := float64(rand.Intn(800 - 15))
	y := -15.0
	fmt.Printf("[SpawnSystem] Creating packet at position (%.1f, %.1f)\n", x, y)

	entity.AddComponent(components.NewTransform(x, y))
	entity.AddComponent(components.NewSprite(15, 15, components.RandomPacketColor()))
	entity.AddComponent(components.NewCollider(15, 15, "packet"))

	physics := components.NewPhysics()
	physics.SetVelocity(0, ss.packetSpeed)
	entity.AddComponent(physics)

	entity.AddComponent(components.NewPacketType(components.RandomPacketName(), 10))
}

// logPacketSpawn logs information about the spawned packet for debugging.
func (ss *SpawnSystem) logPacketSpawn(packet Entity, entity interface{ GetComponentNames() []string }) {
	ss.logEntityStatus(packet)
	ss.logComponentInfo(entity, "packet")
}

// logEntityStatus logs the active status of an entity.
func (ss *SpawnSystem) logEntityStatus(packet Entity) {
	if activeEntity, ok := packet.(interface{ IsActive() bool }); ok {
		fmt.Printf("[SpawnSystem] Spawned entity is active: %v\n", activeEntity.IsActive())
	}
}

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

// trySpawnPowerUp attempts to spawn a power-up if enough time has passed.
func (ss *SpawnSystem) trySpawnPowerUp(eventDispatcher *events.EventDispatcher) {
	if ss.lastPowerUpSpawn >= ss.powerUpSpawnRate {
		ss.lastPowerUpSpawn = 0
		ss.spawnPowerUp()
	}
}

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
