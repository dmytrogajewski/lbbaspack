package systems

import (
	"fmt"
	"lbbaspack/engine/events"
)

// SystemFactory creates and configures all game systems with proper dependencies
type SystemFactory struct {
	entityFactory   func() Entity
	eventDispatcher *events.EventDispatcher
}

// NewSystemFactory creates a new system factory
func NewSystemFactory(entityFactory func() Entity, eventDispatcher *events.EventDispatcher) *SystemFactory {
	return &SystemFactory{
		entityFactory:   entityFactory,
		eventDispatcher: eventDispatcher,
	}
}

// CreateSystemManager creates and configures all systems with proper dependencies
func (sf *SystemFactory) CreateSystemManager() (*SystemManager, error) {
	manager := NewSystemManager()

	// Create all systems
	spawnSys := NewSpawnSystem(sf.entityFactory)
	inputSys := NewInputSystem()
	movementSys := NewMovementSystem()
	collisionSys := NewCollisionSystem()
	powerUpSys := NewPowerUpSystem()
	backendSys := NewBackendSystem()
	slaSys := NewSLASystem(spawnSys)
	comboSys := NewComboSystem()
	gameStateSys := NewGameStateSystem()
	particleSys := NewParticleSystem()
	routingSys := NewRoutingSystem()
	cleanupSys := NewCleanupSystem()

	// Register all ECS systems using SystemInfoer
	systems := []SystemInfoer{
		spawnSys,
		inputSys,
		movementSys,
		collisionSys,
		powerUpSys,
		backendSys,
		slaSys,
		comboSys,
		gameStateSys,
		particleSys,
		routingSys,
		cleanupSys,
	}

	for _, sys := range systems {
		systemInfo := sys.GetSystemInfo()
		if err := manager.RegisterSystem(systemInfo); err != nil {
			return nil, fmt.Errorf("failed to register system %s: %w", systemInfo.Type, err)
		}
	}

	// Build execution order
	if err := manager.BuildExecutionOrder(); err != nil {
		return nil, fmt.Errorf("failed to build execution order: %w", err)
	}

	// Initialize all systems
	spawnSys.Initialize(sf.eventDispatcher)
	backendSys.Initialize(sf.eventDispatcher)
	slaSys.Initialize(sf.eventDispatcher)
	comboSys.Initialize(sf.eventDispatcher)
	gameStateSys.Initialize(sf.eventDispatcher)
	particleSys.Initialize(sf.eventDispatcher)
	routingSys.Initialize(sf.eventDispatcher)
	cleanupSys.Initialize(sf.eventDispatcher)

	// Print execution order for debugging
	manager.PrintExecutionOrder()

	return manager, nil
}

// GetSystemByType is a helper function to get a specific system from the manager
func (sf *SystemFactory) GetSystemByType(manager *SystemManager, systemType SystemType) (System, error) {
	system, exists := manager.GetSystem(systemType)
	if !exists {
		return nil, fmt.Errorf("system %s not found", systemType)
	}
	return system, nil
}
