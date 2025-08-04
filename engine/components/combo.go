package components

type Combo struct {
	Streak int
	Timer  float64
}

func NewCombo() *Combo {
	return &Combo{Streak: 0, Timer: 0}
}

// GetType implements Component interface
func (c *Combo) GetType() string {
	return "Combo"
}

// GetCount implements ComboComponent interface
func (c *Combo) GetCount() int {
	return c.Streak
}

// GetMultiplier implements ComboComponent interface
func (c *Combo) GetMultiplier() float64 {
	if c.Streak <= 1 {
		return 1.0
	}
	return float64(c.Streak) * 0.5
}

// Increment implements ComboComponent interface
func (c *Combo) Increment() {
	c.Streak++
}

// Reset implements ComboComponent interface
func (c *Combo) Reset() {
	c.Streak = 0
	c.Timer = 0
}
