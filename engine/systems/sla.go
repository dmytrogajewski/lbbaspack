package systems

import (
	"fmt"
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

const SystemTypeSLA SystemType = "sla"

type SLASystem struct {
	BaseSystem
}

func NewSLASystem() *SLASystem {
	return &SLASystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"SLA",
			},
		},
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
	for _, entity := range ss.FilterEntities(entities) {
		slaComp := entity.GetSLA()
		if slaComp == nil {
			continue
		}
		if sla, ok := slaComp.(*components.SLA); ok {
			if sla.Total > 0 {
				currentSLA := float64(sla.Caught) / float64(sla.Total) * 100.0
				sla.SetCurrent(currentSLA)
				remainingErrors := sla.ErrorBudget - sla.Lost
				if remainingErrors < 0 {
					remainingErrors = 0
				}
				sla.SetErrorsRemaining(remainingErrors)
				if currentSLA < sla.GetTarget() {
					fmt.Printf("SLA violation! Current: %.2f%%, Target: %.2f%%\n", currentSLA, sla.GetTarget())
				}
			} else {
				sla.SetCurrent(100.0)
				sla.SetErrorsRemaining(sla.ErrorBudget)
			}
		}
	}
}

func (ss *SLASystem) Initialize(eventDispatcher *events.EventDispatcher) {}

// updateSLA removed; SLA events should be published by systems that mutate SLA components if needed

func (ss *SLASystem) SetTargetSLA(target float64) {}

func (ss *SLASystem) SetErrorBudget(budget int) {}

// No internal state; keep empty methods for interface stability if used elsewhere
func (ss *SLASystem) Reset() {}
