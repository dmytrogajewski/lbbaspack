package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// DrawableSystem defines the interface that all drawable systems must implement
// This provides compile-time type checking for the Draw method signature
type DrawableSystem interface {
	System
	Draw(screen *ebiten.Image, entities []Entity)
}

// DrawableSystemGeneric is a generic wrapper that ensures type safety
// This allows us to catch signature mismatches at compile time
type DrawableSystemGeneric[T any] interface {
	System
	Draw(screen T, entities []Entity)
}

// DrawableSystemInfo extends SystemInfo with drawable-specific information
type DrawableSystemInfo struct {
	*SystemInfo
	IsDrawable bool
}

// NewDrawableSystemInfo creates a new DrawableSystemInfo with proper validation
func NewDrawableSystemInfo(info *SystemInfo, isDrawable bool) *DrawableSystemInfo {
	return &DrawableSystemInfo{
		SystemInfo: info,
		IsDrawable: isDrawable,
	}
}

// ValidateDrawableSystem checks if a system implements the correct Draw method
// This function can be used for runtime validation if needed
func ValidateDrawableSystem(system System) (DrawableSystem, bool) {
	if drawable, ok := system.(DrawableSystem); ok {
		return drawable, true
	}
	return nil, false
}

// DrawableSystemRegistry provides type-safe registration of drawable systems
type DrawableSystemRegistry struct {
	systems map[SystemType]DrawableSystem
}

// NewDrawableSystemRegistry creates a new registry for drawable systems
func NewDrawableSystemRegistry() *DrawableSystemRegistry {
	return &DrawableSystemRegistry{
		systems: make(map[SystemType]DrawableSystem),
	}
}

// RegisterDrawableSystem registers a drawable system with compile-time type checking
func (r *DrawableSystemRegistry) RegisterDrawableSystem(systemType SystemType, system DrawableSystem) error {
	r.systems[systemType] = system
	return nil
}

// GetDrawableSystem retrieves a drawable system with type safety
func (r *DrawableSystemRegistry) GetDrawableSystem(systemType SystemType) (DrawableSystem, bool) {
	system, exists := r.systems[systemType]
	return system, exists
}

// GetAllDrawableSystems returns all registered drawable systems
func (r *DrawableSystemRegistry) GetAllDrawableSystems() map[SystemType]DrawableSystem {
	return r.systems
}
