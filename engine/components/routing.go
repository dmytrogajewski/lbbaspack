package components

// Routing component tracks packets that are being routed to backends
type Routing struct {
	TargetBackendID int
	IsRouted        bool
	RouteProgress   float64
	OriginalSpeed   float64 // Store the original packet speed
}

func NewRouting(targetBackendID int, originalSpeed float64) *Routing {
	return &Routing{
		TargetBackendID: targetBackendID,
		IsRouted:        true,
		RouteProgress:   0.0,
		OriginalSpeed:   originalSpeed,
	}
}

// GetType implements Component interface
func (r *Routing) GetType() string {
	return "Routing"
}

// GetTargetBackendID returns the target backend ID
func (r *Routing) GetTargetBackendID() int {
	return r.TargetBackendID
}

// IsPacketRouted returns whether the packet is being routed
func (r *Routing) IsPacketRouted() bool {
	return r.IsRouted
}

// GetRouteProgress returns the routing progress (0.0 to 1.0)
func (r *Routing) GetRouteProgress() float64 {
	return r.RouteProgress
}

// SetRouteProgress sets the routing progress
func (r *Routing) SetRouteProgress(progress float64) {
	r.RouteProgress = progress
}

// GetOriginalSpeed returns the original packet speed
func (r *Routing) GetOriginalSpeed() float64 {
	return r.OriginalSpeed
}
