package ecs

import (
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"lbbaspack/engine/systems"
)

type World struct {
	Entities        []*entities.Entity
	Systems         []systems.System
	EventDispatcher *events.EventDispatcher
	nextEntityID    uint64
}

func NewWorld() *World {
	return &World{
		Entities:        make([]*entities.Entity, 0),
		Systems:         make([]systems.System, 0),
		EventDispatcher: events.NewEventDispatcher(),
		nextEntityID:    1,
	}
}

func (w *World) AddEntity(entity *entities.Entity) {
	w.Entities = append(w.Entities, entity)
}

func (w *World) NewEntity() *entities.Entity {
	entity := entities.NewEntity(w.nextEntityID)
	w.nextEntityID++
	w.AddEntity(entity)
	return entity
}

func (w *World) AddSystem(system systems.System) {
	w.Systems = append(w.Systems, system)
}

func (w *World) Update(deltaTime float64) {
	// Convert entities to the interface type for systems
	entitiesInterface := make([]systems.Entity, len(w.Entities))
	for i, entity := range w.Entities {
		entitiesInterface[i] = entity
	}

	for _, system := range w.Systems {
		system.Update(deltaTime, entitiesInterface, w.EventDispatcher)
	}
}
