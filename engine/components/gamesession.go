package components

// GameSession holds global game progression data
type GameSession struct {
	GameTime        float64
	Score           int
	Level           int
	LastLevelUpTime float64
}

func NewGameSession() *GameSession {
	return &GameSession{GameTime: 0, Score: 0, Level: 1, LastLevelUpTime: 0}
}

func (g *GameSession) GetType() string { return "GameSession" }
