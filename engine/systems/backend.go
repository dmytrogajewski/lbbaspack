package systems

import (
	"fmt"
	"lbbaspack/engine/events"
)

const SystemTypeBackend SystemType = "backend"

type BackendSystem struct {
	BaseSystem
	backendCounters map[int]int // backend ID -> packet count
	totalPackets    int
}

func NewBackendSystem() *BackendSystem {
	return &BackendSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"BackendAssignment",
			},
		},
		backendCounters: make(map[int]int),
		totalPackets:    0,
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (bs *BackendSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeBackend,
		System:       bs,
		Dependencies: []SystemType{SystemTypeCollision},
		Conflicts:    []SystemType{},
		Provides:     []string{"backend_assignment", "load_balancing"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     false,
	}
}

func (bs *BackendSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// First, read current entity counters to sync our internal state
	// But only if we haven't already assigned packets (to avoid overwriting our internal state)
	for _, entity := range bs.FilterEntities(entities) {
		backendComp := entity.GetBackendAssignment()
		if backendComp == nil {
			continue
		}
		backend := backendComp

		// Read the entity's counter to sync our internal state
		backendID := backend.GetBackendID()
		entityCount := backend.GetAssignedPackets()

		// Only update internal counter if it's not already higher (preserve packet assignments)
		if currentCount, exists := bs.backendCounters[backendID]; !exists || entityCount > currentCount {
			bs.backendCounters[backendID] = entityCount
		}
	}

	// Then, update entity counters to match our internal tracking
	for _, entity := range bs.FilterEntities(entities) {
		backendComp := entity.GetBackendAssignment()
		if backendComp == nil {
			continue
		}
		backend := backendComp

		// Update the entity's counter to match our internal tracking
		backendID := backend.GetBackendID()
		if count, exists := bs.backendCounters[backendID]; exists {
			currentCount := backend.GetAssignedPackets()
			// Increment the entity counter to reach the target
			for i := currentCount; i < count; i++ {
				backend.IncrementAssignedPackets()
			}
		}
	}
}

func (bs *BackendSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Listen for packet caught events to assign to backends
	eventDispatcher.Subscribe(events.EventPacketCaught, func(event *events.Event) {
		// We need entities to update, but we don't have them in the event
		// So we'll just update internal counters and let Update method sync them
		bs.assignPacketToBackend()
	})

	// Listen for game start events to reset backend counters
	eventDispatcher.Subscribe(events.EventGameStart, func(event *events.Event) {
		bs.resetBackendCounters()
	})
}

// InitializeBackendCounters initializes the backend counters from existing entities
func (bs *BackendSystem) InitializeBackendCounters(entities []Entity) {
	for _, entity := range entities {
		if backendComp := entity.GetBackendAssignment(); backendComp != nil {
			backendID := backendComp.GetBackendID()
			bs.backendCounters[backendID] = 0
			fmt.Printf("[BackendSystem] Initialized counter for backend %d\n", backendID)
		}
	}
}

func (bs *BackendSystem) assignPacketToBackend() {
	// Check if there are any backends
	if len(bs.backendCounters) == 0 {
		// No backends available, don't assign packet
		return
	}

	// Find backend with least packets (load balancing)
	minCount := 999999
	selectedBackend := 0

	for backendID, count := range bs.backendCounters {
		if count < minCount {
			minCount = count
			selectedBackend = backendID
		}
	}

	// Increment the selected backend's counter
	bs.backendCounters[selectedBackend]++
	bs.totalPackets++

	fmt.Printf("Packet assigned to backend %d (total: %d, backend count: %d)\n",
		selectedBackend, bs.totalPackets, bs.backendCounters[selectedBackend])
}

func (bs *BackendSystem) GetBackendStats() map[int]int {
	return bs.backendCounters
}

func (bs *BackendSystem) GetTotalPackets() int {
	return bs.totalPackets
}

// resetBackendCounters resets all backend counters when starting a new game
func (bs *BackendSystem) resetBackendCounters() {
	bs.backendCounters = make(map[int]int)
	bs.totalPackets = 0
	fmt.Println("[BackendSystem] Reset backend counters")
}
