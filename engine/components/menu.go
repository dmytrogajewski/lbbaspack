package components

// MenuState holds menu selection and key latch
type MenuState struct {
	SelectedMode int
	KeyLatch     bool
}

func NewMenuState() *MenuState { return &MenuState{SelectedMode: 0, KeyLatch: false} }

func (m *MenuState) GetType() string { return "MenuState" }
