package systems

import (
	"lbbaspack/engine/components"
	"lbbaspack/engine/entities"
	"lbbaspack/engine/events"
	"testing"
)

func TestNewSLASystem(t *testing.T) {
	ss := NewSLASystem()
	if ss == nil {
		t.Fatal("NewSLASystem returned nil")
	}
	expected := []string{"SLA"}
	if len(ss.RequiredComponents) != len(expected) {
		t.Fatalf("expected %d required components, got %d", len(expected), len(ss.RequiredComponents))
	}
	for i, c := range expected {
		if ss.RequiredComponents[i] != c {
			t.Fatalf("expected %s, got %s", c, ss.RequiredComponents[i])
		}
	}
}

func TestSLASystem_Update_ComputesFromComponent(t *testing.T) {
	ss := NewSLASystem()
	ed := events.NewEventDispatcher()
	e := entities.NewEntity(1)
	sla := components.NewSLA(95.0, 10)
	sla.Total = 10
	sla.Caught = 8
	sla.Lost = 2
	e.AddComponent(sla)
	ss.Update(0.016, []Entity{e}, ed)
	if sla.GetCurrent() != 80.0 {
		t.Errorf("expected current 80.0, got %.2f", sla.GetCurrent())
	}
	if sla.GetErrorsRemaining() != 8 {
		t.Errorf("expected remaining 8, got %d", sla.GetErrorsRemaining())
	}
}

func TestSLASystem_Initialize_NoPanic(t *testing.T) {
	ss := NewSLASystem()
	ss.Initialize(events.NewEventDispatcher())
}
