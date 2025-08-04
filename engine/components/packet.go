package components

import (
	"image/color"
	"math/rand"
)

// PacketType component identifies the type of packet and its value
type PacketType struct {
	Name  string
	Value int
}

func NewPacketType(name string, value int) *PacketType {
	return &PacketType{Name: name, Value: value}
}

// GetType implements Component interface
func (pt *PacketType) GetType() string {
	return "PacketType"
}

// GetName implements PacketTypeComponent interface
func (pt *PacketType) GetName() string {
	return pt.Name
}

// GetPriority implements PacketTypeComponent interface
func (pt *PacketType) GetPriority() int {
	return pt.Value
}

// RandomPacketColor returns a random packet color
func RandomPacketColor() color.RGBA {
	packetColors := []color.RGBA{
		{255, 0, 0, 255},   // Red for HTTP
		{0, 255, 0, 255},   // Green for HTTPS
		{0, 0, 255, 255},   // Blue for TCP
		{255, 255, 0, 255}, // Yellow for UDP
		{255, 0, 255, 255}, // Magenta for WebSocket
	}
	return packetColors[rand.Intn(len(packetColors))]
}

// RandomPacketName returns a random packet type name
func RandomPacketName() string {
	packetNames := []string{"HTTP", "HTTPS", "TCP", "UDP", "WebSocket"}
	return packetNames[rand.Intn(len(packetNames))]
}
