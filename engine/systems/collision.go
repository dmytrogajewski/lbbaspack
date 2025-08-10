package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/events"
)

const SystemTypeCollision SystemType = "collision"

type CollisionSystem struct {
	BaseSystem
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{
		BaseSystem: BaseSystem{
			RequiredComponents: []string{
				"Transform",
				"Collider",
			},
		},
	}
}

// GetSystemInfo returns the system metadata for dependency resolution
func (cs *CollisionSystem) GetSystemInfo() *SystemInfo {
	return &SystemInfo{
		Type:         SystemTypeCollision,
		System:       cs,
		Dependencies: []SystemType{}, // No dependencies - runs independently and checks collisions
		Conflicts:    []SystemType{},
		Provides:     []string{"collision_detection", "packet_catching"},
		Requires:     []string{},
		Drawable:     false,
		Optional:     false,
	}
}

func (cs *CollisionSystem) Update(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	// Build list of active collider-bearing entities
	type coll struct {
		e Entity
		t components.TransformComponent
		c components.ColliderComponent
	}
	var colliders []coll
	for _, entity := range cs.FilterEntities(entities) {
		if !entity.IsActive() {
			continue
		}
		t := entity.GetTransform()
		c := entity.GetCollider()
		if t == nil || c == nil {
			continue
		}
		colliders = append(colliders, coll{e: entity, t: t, c: c})
	}

	// Pairwise collision check and event emit
	for i := 0; i < len(colliders); i++ {
		for j := i + 1; j < len(colliders); j++ {
			a := colliders[i]
			b := colliders[j]
			if cs.checkCollision(a.t, a.c, b.t, b.c) {
				tagA := a.c.GetTag()
				tagB := b.c.GetTag()
				ax, ay := a.t.GetX(), a.t.GetY()
				bx, by := b.t.GetX(), b.t.GetY()
				eventDispatcher.Publish(events.NewEvent(events.EventCollisionDetected, &events.EventData{
					EntityA: a.e,
					EntityB: b.e,
					TagA:    &tagA,
					TagB:    &tagB,
					PosAX:   &ax,
					PosAY:   &ay,
					PosBX:   &bx,
					PosBY:   &by,
				}))
			}
		}
	}
}

func (cs *CollisionSystem) checkCollision(transform1 components.TransformComponent, collider1 components.ColliderComponent,
	transform2 components.TransformComponent, collider2 components.ColliderComponent) bool {

	// Simple AABB collision detection
	left1 := transform1.GetX()
	right1 := transform1.GetX() + collider1.GetWidth()
	top1 := transform1.GetY()
	bottom1 := transform1.GetY() + collider1.GetHeight()

	left2 := transform2.GetX()
	right2 := transform2.GetX() + collider2.GetWidth()
	top2 := transform2.GetY()
	bottom2 := transform2.GetY() + collider2.GetHeight()

	return !(right1 < left2 || left1 > right2 || bottom1 < top2 || top1 > bottom2)
}

// checkCollisionLBEnhanced removed; CollisionSystem stays pure
