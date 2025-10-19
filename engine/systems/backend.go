package systems

import (
	"lbbaspack/engine/events"
)

const SystemTypeBackend SystemType = "backend"

type BackendSystem struct {
	BaseSystem
}

func NewBackendSystem() *BackendSystem {
	return &BackendSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"BackendAssignment",
			},
		},
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
	// Stateless: no-op. Routing/assignment is handled by collision/packet routing
}

func (bs *BackendSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Stateless: no subscriptions required
}

// InitializeBackendCounters initializes the backend counters from existing entities
func (bs *BackendSystem) InitializeBackendCounters(entities []Entity) {}

func (bs *BackendSystem) GetBackendStats() map[int]int { return map[int]int{} }

func (bs *BackendSystem) GetTotalPackets() int { return 0 }
