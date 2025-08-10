package components

// SLAEvents is a world-level queue for SLA-related increments
// Systems write increments here; SLASystem consumes and applies them to SLA components
type SLAEvents struct {
	CaughtIncrements int
	LostIncrements   int
}

func NewSLAEvents() *SLAEvents { return &SLAEvents{} }

func (e *SLAEvents) GetType() string { return "SLAEvents" }
