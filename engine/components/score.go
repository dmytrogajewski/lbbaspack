package components

type Score struct {
	Value int
}

func NewScore() *Score {
	return &Score{Value: 0}
}
