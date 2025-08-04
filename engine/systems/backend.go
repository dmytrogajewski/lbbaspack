package systems

import (
	"fmt"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

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

func (bs *BackendSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Update backend counters and handle packet assignments
	for _, entity := range bs.FilterEntities(entities) {
		backendComp := entity.GetBackendAssignment()
		if backendComp == nil {
			continue
		}
		backend, ok := backendComp.(components.BackendAssignmentComponent)
		if !ok {
			continue
		}

		// Update the global counter for this backend
		bs.backendCounters[backend.GetBackendID()] = backend.GetAssignedPackets()
	}
}

func (bs *BackendSystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Listen for packet caught events to assign to backends
	eventDispatcher.Subscribe(events.EventPacketCaught, func(event *events.Event) {
		bs.assignPacketToBackend()
	})
}

func (bs *BackendSystem) assignPacketToBackend() {
	// Simple round-robin assignment
	backendCount := len(bs.backendCounters)
	if backendCount == 0 {
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
