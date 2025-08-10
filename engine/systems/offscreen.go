package systems

import (
	"lbbaspack/engine/events"
)

const SystemTypeOffscreen SystemType = "offscreen"

// OffscreenSystem detects entities that left the screen and emits events
type OffscreenSystem struct{ BaseSystem }

func NewOffscreenSystem() *OffscreenSystem {
	return &OffscreenSystem{BaseSystem: BaseSystem{RequiredComponents: []string{"Transform", "Collider"}}}
}

func (os *OffscreenSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeOffscreen,
		System:       os,
		Dependencies: []SystemType{},
		Conflicts:    []SystemType{},
		Provides:     []string{"offscreen_detection"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     false,
	}
}

func (os *OffscreenSystem) Initialize(eventDispatcher *events.EventDispatcher) {}

func (os *OffscreenSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Deactivate entities that leave the screen; emit a generic offscreen/loss event
	// Identify load balancer entity (if any) to include in event context
	var lb Entity
	for _, ent := range os.FilterEntities(entities) {
		if c := ent.GetCollider(); c != nil && c.GetTag() == "loadbalancer" {
			lb = ent
			break
		}
	}
	for _, e := range os.FilterEntities(entities) {
		t := e.GetTransform()
		c := e.GetCollider()
		if t == nil || c == nil {
			continue
		}
		if t.GetY() > 600 {
			if pe, ok := e.(interface{ SetActive(bool) }); ok {
				pe.SetActive(false)
			}
			tag := c.GetTag()
			x, y := t.GetX(), t.GetY()
			// Emit generic collider-offscreen event
			// Include load balancer entity if found to allow SLA system to update LB counters
			var tagA *string
			if lb != nil {
				if lbC := lb.GetCollider(); lbC != nil {
					tga := lbC.GetTag()
					tagA = &tga
				}
			}
			eventDispatcher.Publish(events.NewEvent(events.EventColliderOffscreen, &events.EventData{
				EntityA: lb,
				TagA:    tagA,
				EntityB: e,
				TagB:    &tag,
				PosBX:   &x,
				PosBY:   &y,
			}))
		}
	}
}
