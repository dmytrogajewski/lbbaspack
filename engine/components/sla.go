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

// Added helpers for ECS-stateless systems to update counters on the component

// GetTotals returns total, caught, lost counters
func (s *SLA) GetTotals() (int, int, int) {
	return s.Total, s.Caught, s.Lost
}

// SetTarget updates the SLA target percentage
func (s *SLA) SetTarget(target float64) {
	s.Target = target
}

// SetErrorBudget updates the allowed error budget and resets remaining accordingly
func (s *SLA) SetErrorBudget(budget int) {
	s.ErrorBudget = budget
	if s.RemainingErrors > budget {
		s.RemainingErrors = budget
	}
}

// IncrementCaught increments caught and total counters
func (s *SLA) IncrementCaught() {
	s.Caught++
	s.Total++
}

// IncrementLost increments lost and total counters and decreases remaining
func (s *SLA) IncrementLost() {
	s.Lost++
	s.Total++
	s.RemainingErrors = s.ErrorBudget - s.Lost
	if s.RemainingErrors < 0 {
		s.RemainingErrors = 0
	}
}

// ResetCounters clears derived counters but preserves configuration
func (s *SLA) ResetCounters() {
	s.Total = 0
	s.Caught = 0
	s.Lost = 0
	s.Current = 100.0
	s.RemainingErrors = s.ErrorBudget
}
