package components

// RoutingRequest is a marker component indicating a packet should be assigned to a backend
type RoutingRequest struct{}

func NewRoutingRequest() *RoutingRequest { return &RoutingRequest{} }

func (r *RoutingRequest) GetType() string { return "RoutingRequest" }
