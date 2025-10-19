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

func (ss *SLASystem) Initialize(eventDispatcher *events.EventDispatcher) {
	// Subscribe and update SLA counters directly on components
	eventDispatcher.Subscribe(events.EventCollisionDetected, func(event *events.Event) {
		if event == nil || event.Data == nil {
			return
		}
		if event.Data.TagA == nil || event.Data.TagB == nil {
			return
		}
		aIsLB := *event.Data.TagA == "loadbalancer"
		bIsLB := *event.Data.TagB == "loadbalancer"
		aIsPacket := *event.Data.TagA == "packet"
		bIsPacket := *event.Data.TagB == "packet"
		if (aIsLB && bIsPacket) || (bIsLB && aIsPacket) {
			var target Entity
			if aIsLB {
				if e, ok := event.Data.EntityA.(Entity); ok {
					target = e
				}
			} else if bIsLB {
				if e, ok := event.Data.EntityB.(Entity); ok {
					target = e
				}
			}
			if target != nil {
				if sc := target.GetSLA(); sc != nil {
					if sla, ok := sc.(*components.SLA); ok {
						sla.IncrementCaught()
						current := 100.0
						if sla.Total > 0 {
							current = float64(sla.Caught) / float64(sla.Total) * 100.0
							sla.SetCurrent(current)
						}
						remaining := sla.ErrorBudget - sla.Lost
						if remaining < 0 {
							remaining = 0
						}
						eventDispatcher.Publish(events.NewEvent(events.EventSLAUpdated, &events.EventData{Current: &current, Caught: &sla.Caught, Lost: &sla.Lost, Remaining: &remaining, Budget: &sla.ErrorBudget}))
					}
				}
			}
		}
	})
	eventDispatcher.Subscribe(events.EventColliderOffscreen, func(event *events.Event) {
		if event == nil || event.Data == nil || event.Data.TagB == nil {
			return
		}
		if *event.Data.TagB != "packet" {
			return
		}
		// increment lost on any entity that has SLA (prefer one present in event)
		for _, ent := range []interface{}{event.Data.EntityA, event.Data.EntityB} {
			if e, ok := ent.(Entity); ok && e != nil {
				if sc := e.GetSLA(); sc != nil {
					if sla, ok := sc.(*components.SLA); ok {
						sla.IncrementLost()
						current := 100.0
						if sla.Total > 0 {
							current = float64(sla.Caught) / float64(sla.Total) * 100.0
							sla.SetCurrent(current)
						}
						remaining := sla.ErrorBudget - sla.Lost
						if remaining < 0 {
							remaining = 0
						}
						eventDispatcher.Publish(events.NewEvent(events.EventSLAUpdated, &events.EventData{Current: &current, Caught: &sla.Caught, Lost: &sla.Lost, Remaining: &remaining, Budget: &sla.ErrorBudget}))
						if remaining <= 0 {
							eventDispatcher.Publish(events.NewEvent(events.EventGameOver, &events.EventData{Score: &sla.Caught, Lost: &sla.Lost}))
						}
						break
					}
				}
			}
		}
	})
}

// React to offscreen events in Update by scanning SLA components; keep Initialize empty to remain stateless

// SLA counters now updated by reacting to events in Update only

// updateSLA removed; SLA events should be published by systems that mutate SLA components if needed

func (ss *SLASystem) SetTargetSLA(target float64) {}

func (ss *SLASystem) SetErrorBudget(budget int) {}

// No internal state; keep empty methods for interface stability if used elsewhere
func (ss *SLASystem) Reset() {}
