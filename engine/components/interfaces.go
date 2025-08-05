package components

import "image/color"

// Component represents a component interface
type Component interface {
	GetType() string
}

// TransformComponent represents transform functionality
type TransformComponent interface {
	Component
	GetX() float64
	GetY() float64
	SetPosition(x, y float64)
}

// SpriteComponent represents sprite functionality
type SpriteComponent interface {
	Component
	GetWidth() float64
	GetHeight() float64
	GetColor() color.RGBA
	IsVisible() bool
	SetVisible(visible bool)
}

// ColliderComponent represents collider functionality
type ColliderComponent interface {
	Component
	GetTag() string
	GetWidth() float64
	GetHeight() float64
	SetTag(tag string)
}

// PhysicsComponent represents physics functionality
type PhysicsComponent interface {
	Component
	GetVelocityX() float64
	GetVelocityY() float64
	SetVelocity(x, y float64)
}

// PacketTypeComponent represents packet type functionality
type PacketTypeComponent interface {
	Component
	GetName() string
	GetPriority() int
}

// StateComponent represents state functionality
type StateComponent interface {
	Component
	GetState() string
	SetState(state string)
}

// ComboComponent represents combo functionality
type ComboComponent interface {
	Component
	GetCount() int
	GetMultiplier() float64
	Increment()
	Reset()
}

// SLAComponent represents SLA functionality
type SLAComponent interface {
	Component
	GetCurrent() float64
	GetTarget() float64
	GetErrorsRemaining() int
	SetCurrent(current float64)
	SetErrorsRemaining(errors int)
}

// BackendAssignmentComponent represents backend assignment functionality
type BackendAssignmentComponent interface {
	Component
	GetBackendID() int
	GetAssignedPackets() int
	SetBackendID(id int)
	IncrementAssignedPackets()
}

// PowerUpTypeComponent represents power-up functionality
type PowerUpTypeComponent interface {
	Component
	GetName() string
	GetDuration() float64
	GetEffect() string
}

// RoutingComponent represents routing functionality
type RoutingComponent interface {
	Component
	GetTargetBackendID() int
	IsPacketRouted() bool
	GetRouteProgress() float64
	SetRouteProgress(progress float64)
	GetOriginalSpeed() float64
}
