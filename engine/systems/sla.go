package systems

import (
	"fmt"
	"lbbaspack/engine/events"
)

const SystemTypeSLA SystemType = "sla"

type SLASystem struct {
	BaseSystem
	totalPackets  int
	caughtPackets int
	lostPackets   int
	errorBudget   int
	spawnSys      *SpawnSystem // Reference to SpawnSystem
}

func NewSLASystem(spawnSys *SpawnSystem) *SLASystem {
	return &SLASystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"SLA",
			},
		},
		totalPackets:  0,
		caughtPackets: 0,
		lostPackets:   0,
		errorBudget:   10, // Default error budget
		spawnSys:      spawnSys,
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (slas *SLASystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeSLA,
		System:       slas,
		Dependencies: []SystemType{SystemTypeCollision}, // Depends on collision for packet events
		Conflicts:    []SystemType{},
		Provides:     []string{"sla_monitoring", "performance_tracking"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     false,
	}
}

func (ss *SLASystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Update SLA components
	for _, entity := range ss.FilterEntities(entities) {
		slaComp := entity.GetSLA()
		if slaComp == nil {
			continue
		}
		sla := slaComp

		// Calculate current SLA percentage
		if ss.totalPackets > 0 {
			currentSLA := float64(ss.caughtPackets) / float64(ss.totalPackets) * 100.0
			sla.SetCurrent(currentSLA)

			// Update remaining errors
			remainingErrors := ss.errorBudget - ss.lostPackets
			if remainingErrors < 0 {
				remainingErrors = 0
			}
			sla.SetErrorsRemaining(remainingErrors)

			// Check for SLA violations
			if currentSLA < sla.GetTarget() {
				fmt.Printf("SLA violation! Current: %.2f%%, Target: %.2f%%\n", currentSLA, sla.GetTarget())
			}
		}
	}
}

func (ss *SLASystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Listen for packet events to update SLA
	eventDispatcher.Subscribe(events.EventPacketCaught, func(event *events.Event) {
		ss.caughtPackets++
		ss.totalPackets++
		ss.updateSLA(eventDispatcher)
	})

	eventDispatcher.Subscribe(events.EventPacketLost, func(event *events.Event) {
		ss.lostPackets++
		ss.totalPackets++
		// Increase packet speed by 5% on each lost packet
		if ss.spawnSys != nil {
			ss.spawnSys.IncreasePacketSpeed(5.0)
		}
		ss.updateSLA(eventDispatcher)
	})
}

func (ss *SLASystem) updateSLA(eventDispatcher *events.EventDispatcher) {
	if ss.totalPackets > 0 {
		currentSLA := float64(ss.caughtPackets) / float64(ss.totalPackets) * 100.0
		remainingErrors := ss.errorBudget - ss.lostPackets

		// Only print "Packet lost!" message when packets are actually lost
		if ss.lostPackets > 0 {
			fmt.Printf("Packet lost! SLA: %.2f%%, Errors remaining: %d/%d\n", currentSLA, remainingErrors, ss.errorBudget)
		}

		// Publish SLA update event for UI (for both caught and lost packets)
		eventDispatcher.Publish(events.NewEvent(events.EventSLAUpdated, &events.EventData{
			Current:   &currentSLA,
			Caught:    &ss.caughtPackets,
			Lost:      &ss.lostPackets,
			Remaining: &remainingErrors,
			Budget:    &ss.errorBudget,
		}))

		// Check if error budget has been exceeded
		if remainingErrors <= 0 {
			fmt.Printf("ERROR BUDGET EXCEEDED! Game Over!\n")
			// Publish game over event
			eventDispatcher.Publish(events.NewEvent(events.EventGameOver, &events.EventData{
				Score: &ss.caughtPackets,
				Lost:  &ss.lostPackets,
			}))
		}
	}
}

func (ss *SLASystem) SetTargetSLA(target float64) {
	// This sets the target SLA for all entities with an SLA component
	fmt.Printf("SLA target set to %.2f%%\n", target)
	// Optionally, you could update all SLA components here
}

func (ss *SLASystem) SetErrorBudget(budget int) {
	ss.errorBudget = budget
	fmt.Printf("Error budget set to %d errors\n", budget)
	// Optionally, you could update all SLA components here
}
