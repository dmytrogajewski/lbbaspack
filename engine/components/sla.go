package components

type SLA struct {
	Target          float64
	Current         float64
	ErrorBudget     int
	RemainingErrors int
	Total           int
	Caught          int
	Lost            int
}

func NewSLA(target float64, errorBudget int) *SLA {
	return &SLA{
		Target:          target,
		Current:         100.0,
		ErrorBudget:     errorBudget,
		RemainingErrors: errorBudget,
	}
}

// GetType implements Component interface
func (s *SLA) GetType() string {
	return "SLA"
}

// GetCurrent implements SLAComponent interface
func (s *SLA) GetCurrent() float64 {
	return s.Current
}

// GetTarget implements SLAComponent interface
func (s *SLA) GetTarget() float64 {
	return s.Target
}

// GetErrorsRemaining implements SLAComponent interface
func (s *SLA) GetErrorsRemaining() int {
	return s.RemainingErrors
}

// SetCurrent implements SLAComponent interface
func (s *SLA) SetCurrent(current float64) {
	s.Current = current
}

// SetErrorsRemaining implements SLAComponent interface
func (s *SLA) SetErrorsRemaining(errors int) {
	s.RemainingErrors = errors
}
