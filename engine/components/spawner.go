package components

// Spawner holds spawning configuration and runtime timers
type Spawner struct {
	PacketSpawnElapsed  float64
	PacketSpawnRate     float64
	PowerUpSpawnElapsed float64
	PowerUpSpawnRate    float64
	PacketSpeed         float64
	Level               int
	// DDoS state
	IsDDoSActive bool
	DDOSTimer    float64
	DDoSDuration float64
	DDoSMult     float64
	DDoSCooldown float64
}

func NewSpawner() *Spawner {
	return &Spawner{
		PacketSpawnElapsed:  0,
		PacketSpawnRate:     1.0,
		PowerUpSpawnElapsed: 0,
		PowerUpSpawnRate:    10.0,
		PacketSpeed:         100.0,
		Level:               1,
		IsDDoSActive:        false,
		DDOSTimer:           0,
		DDoSDuration:        5.0,
		DDoSMult:            10.0,
		DDoSCooldown:        10.0,
	}
}

func (s *Spawner) GetType() string { return "Spawner" }

func (s *Spawner) IncreasePacketSpeed(percent float64) {
	s.PacketSpeed *= (1.0 + percent/100.0)
}

func (s *Spawner) SetLevel(level int) {
	s.Level = level
	s.PacketSpawnRate = 1.0 / (1.0 + float64(level-1)*0.2)
}
