package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"
)

func TestSpawnSystem_SpawnsPacketUsingSpawnerComponent(t *testing.T) {
	var created []*entities.Entity
	ss := NewSpawnSystem(func() Entity {
		e := entities.NewEntity(uint64(len(created) + 1))
		created = append(created, e)
		return e
	})

	spawnerHolder := entities.NewEntity(99)
	spawner := components.NewSpawner()
	spawner.PacketSpawnRate = 0.0
	spawnerHolder.AddComponent(spawner)

	ss.Update(0.016, []Entity{spawnerHolder}, events.NewEventDispatcher())

	if len(created) == 0 {
		t.Fatalf("expected at least one spawned entity")
	}
	got := created[0]
	if !got.HasComponent("Transform") || !got.HasComponent("Sprite") || !got.HasComponent("Collider") || !got.HasComponent("Physics") || !got.HasComponent("PacketType") {
		t.Fatalf("spawned entity missing required components")
	}
}

func TestSpawnSystem_RespectsSpawnerSpeed(t *testing.T) {
	var created []*entities.Entity
	ss := NewSpawnSystem(func() Entity {
		e := entities.NewEntity(uint64(len(created) + 1))
		created = append(created, e)
		return e
	})

	spawnerHolder := entities.NewEntity(1)
	spawner := components.NewSpawner()
	spawner.PacketSpawnRate = 0.0
	spawner.PacketSpeed = 321.0
	spawnerHolder.AddComponent(spawner)

	ss.Update(0.016, []Entity{spawnerHolder}, events.NewEventDispatcher())
	if len(created) == 0 {
		t.Fatalf("no spawn")
	}
	p := created[0]
	phys := p.GetPhysics()
	if phys == nil {
		t.Fatalf("missing physics")
	}
	if phys.GetVelocityY() != 321.0 {
		t.Fatalf("expected velocityY 321, got %.1f", phys.GetVelocityY())
	}
}
