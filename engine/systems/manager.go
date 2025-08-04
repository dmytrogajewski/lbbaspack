package systems

import (
	"fmt"
	"lbbaspack/engine/events"

	"github.com/hajimehoshi/ebiten/v2"
)

// SystemType represents different types of systems
type SystemType string

// SystemInfo contains metadata about a system
type SystemInfo struct {
	Type         SystemType
	System       System
	Dependencies []SystemType
	Conflicts    []SystemType // Systems that cannot run together
	Provides     []string     // Capabilities this system provides
	Requires     []string     // Capabilities this system requires
	Drawable     bool         // Whether this system has a Draw method
	Optional     bool         // Whether this system is optional
}

// DependencyResolver implements robust dependency resolution based on libsolv practices
type DependencyResolver struct {
	systems map[SystemType]*SystemInfo
	// Capability mapping: capability -> list of systems that provide it
	capabilityMap map[string][]SystemType
	// Dependency graph for topological sorting
	dependencyGraph map[SystemType][]SystemType
	// Reverse dependency graph for impact analysis
	reverseGraph map[SystemType][]SystemType
	verbose      bool // Enable verbose logging
}

// SystemManager manages system dependencies and execution order
type SystemManager struct {
	systems     map[SystemType]*SystemInfo
	updateOrder []SystemType
	drawOrder   []SystemType
	resolver    *DependencyResolver
	verbose     bool // Enable verbose logging
}

// NewDependencyResolver creates a new dependency resolver
func NewDependencyResolver() *DependencyResolver {
	return &DependencyResolver{
		systems:         make(map[SystemType]*SystemInfo),
		capabilityMap:   make(map[string][]SystemType),
		dependencyGraph: make(map[SystemType][]SystemType),
		reverseGraph:    make(map[SystemType][]SystemType),
		verbose:         true, // Enable verbose logging by default
	}
}

// NewSystemManager creates a new system manager
func NewSystemManager() *SystemManager {
	return &SystemManager{
		systems:     make(map[SystemType]*SystemInfo),
		updateOrder: make([]SystemType, 0),
		drawOrder:   make([]SystemType, 0),
		resolver:    NewDependencyResolver(),
		verbose:     true, // Enable verbose logging by default
	}
}

// SetVerbose enables or disables verbose logging
func (sm *SystemManager) SetVerbose(verbose bool) {
	sm.verbose = verbose
	sm.resolver.verbose = verbose
}

// RegisterSystem registers a system with its metadata
func (sm *SystemManager) RegisterSystem(info *SystemInfo) error {
	// Check for duplicate registration
	if _, exists := sm.systems[info.Type]; exists {
		return fmt.Errorf("system %s is already registered", info.Type)
	}

	if sm.verbose {
		fmt.Printf("[SystemManager] Registering system: %s\n", info.Type)
		fmt.Printf("  - Dependencies: %v\n", info.Dependencies)
		fmt.Printf("  - Provides: %v\n", info.Provides)
		fmt.Printf("  - Requires: %v\n", info.Requires)
		fmt.Printf("  - Drawable: %v, Optional: %v\n", info.Drawable, info.Optional)
	}

	sm.systems[info.Type] = info
	return nil
}

// BuildExecutionOrder builds the correct execution order using robust dependency resolution
func (sm *SystemManager) BuildExecutionOrder() error {
	if sm.verbose {
		fmt.Println("\n[SystemManager] Building execution order...")
		fmt.Printf("Total systems registered: %d\n", len(sm.systems))
	}

	// Use the dependency resolver to build the execution order
	return sm.resolver.ResolveDependencies(sm.systems, &sm.updateOrder, &sm.drawOrder)
}

// ResolveDependencies implements robust dependency resolution based on libsolv practices
func (dr *DependencyResolver) ResolveDependencies(systems map[SystemType]*SystemInfo, updateOrder *[]SystemType, drawOrder *[]SystemType) error {
	if dr.verbose {
		fmt.Println("\n[DependencyResolver] Starting dependency resolution...")
	}

	// Clear existing orders
	*updateOrder = make([]SystemType, 0)
	*drawOrder = make([]SystemType, 0)

	// Build capability map and dependency graphs
	if err := dr.buildGraphs(systems); err != nil {
		return err
	}

	// Check for conflicts
	if err := dr.checkConflicts(systems); err != nil {
		return err
	}

	// Resolve capability dependencies
	if err := dr.resolveCapabilities(systems); err != nil {
		return err
	}

	// Perform topological sort
	order, err := dr.topologicalSort()
	if err != nil {
		return err
	}

	*updateOrder = order

	// Build draw order
	for _, systemType := range *updateOrder {
		if systems[systemType].Drawable {
			*drawOrder = append(*drawOrder, systemType)
		}
	}

	if dr.verbose {
		fmt.Printf("[DependencyResolver] Resolution completed successfully\n")
		fmt.Printf("Update order: %v\n", *updateOrder)
		fmt.Printf("Draw order: %v\n", *drawOrder)
	}

	return nil
}

// buildGraphs builds the dependency and capability graphs
func (dr *DependencyResolver) buildGraphs(systems map[SystemType]*SystemInfo) error {
	if dr.verbose {
		fmt.Println("\n[DependencyResolver] Building dependency and capability graphs...")
	}

	dr.systems = systems

	// Build capability map
	if dr.verbose {
		fmt.Println("Building capability map...")
	}
	for systemType, info := range systems {
		for _, capability := range info.Provides {
			dr.capabilityMap[capability] = append(dr.capabilityMap[capability], systemType)
			if dr.verbose {
				fmt.Printf("  %s provides capability: %s\n", systemType, capability)
			}
		}
	}

	// Print capability map
	if dr.verbose {
		fmt.Println("\nCapability Map:")
		for capability, providers := range dr.capabilityMap {
			fmt.Printf("  %s -> %v\n", capability, providers)
		}
	}

	// Build dependency graph
	if dr.verbose {
		fmt.Println("\nBuilding dependency graph...")
	}
	for systemType, info := range systems {
		deps := make([]SystemType, 0)

		// Add direct dependencies
		deps = append(deps, info.Dependencies...)
		if dr.verbose && len(info.Dependencies) > 0 {
			fmt.Printf("  %s direct dependencies: %v\n", systemType, info.Dependencies)
		}

		// Add capability-based dependencies
		for _, requiredCapability := range info.Requires {
			if providers, exists := dr.capabilityMap[requiredCapability]; exists {
				// For now, take the first provider. In a more sophisticated system,
				// we could implement provider selection based on preferences
				if len(providers) > 0 {
					deps = append(deps, providers[0])
					if dr.verbose {
						fmt.Printf("  %s requires capability '%s' -> depends on %s\n", systemType, requiredCapability, providers[0])
					}
				}
			} else {
				return fmt.Errorf("system %s requires capability '%s' but no system provides it", systemType, requiredCapability)
			}
		}

		dr.dependencyGraph[systemType] = deps

		// Build reverse graph for impact analysis
		for _, dep := range deps {
			dr.reverseGraph[dep] = append(dr.reverseGraph[dep], systemType)
		}
	}

	// Print dependency graph
	if dr.verbose {
		fmt.Println("\nDependency Graph:")
		for systemType, deps := range dr.dependencyGraph {
			fmt.Printf("  %s -> %v\n", systemType, deps)
		}

		fmt.Println("\nReverse Dependency Graph:")
		for systemType, dependents := range dr.reverseGraph {
			fmt.Printf("  %s <- %v\n", systemType, dependents)
		}
	}

	return nil
}

// checkConflicts verifies that no conflicting systems are present
func (dr *DependencyResolver) checkConflicts(systems map[SystemType]*SystemInfo) error {
	if dr.verbose {
		fmt.Println("\n[DependencyResolver] Checking for system conflicts...")
	}

	for systemType, info := range systems {
		for _, conflict := range info.Conflicts {
			if _, exists := systems[conflict]; exists {
				if dr.verbose {
					fmt.Printf("  CONFLICT: %s conflicts with %s\n", systemType, conflict)
				}
				return fmt.Errorf("conflict detected: system %s conflicts with system %s", systemType, conflict)
			}
		}
	}

	if dr.verbose {
		fmt.Println("  No conflicts detected")
	}
	return nil
}

// resolveCapabilities resolves capability-based dependencies
func (dr *DependencyResolver) resolveCapabilities(systems map[SystemType]*SystemInfo) error {
	if dr.verbose {
		fmt.Println("\n[DependencyResolver] Resolving capability dependencies...")
	}

	// This is a simplified implementation. A full resolver would:
	// 1. Handle multiple providers for the same capability
	// 2. Implement provider selection based on preferences
	// 3. Handle version constraints
	// 4. Implement backtracking for complex dependency resolution

	for systemType, info := range systems {
		for _, requiredCapability := range info.Requires {
			providers := dr.capabilityMap[requiredCapability]
			if len(providers) == 0 {
				if dr.verbose {
					fmt.Printf("  ERROR: %s requires capability '%s' but no system provides it\n", systemType, requiredCapability)
				}
				return fmt.Errorf("system %s requires capability '%s' but no system provides it", systemType, requiredCapability)
			}
			if dr.verbose {
				fmt.Printf("  %s requires '%s' -> provided by %v\n", systemType, requiredCapability, providers)
			}
			// For now, we just verify that at least one provider exists
			// In a more sophisticated system, we'd select the best provider
		}
	}

	if dr.verbose {
		fmt.Println("  All capability dependencies resolved")
	}
	return nil
}

// topologicalSort performs topological sorting of the dependency graph
func (dr *DependencyResolver) topologicalSort() ([]SystemType, error) {
	if dr.verbose {
		fmt.Println("\n[DependencyResolver] Performing topological sort...")
	}

	// Kahn's algorithm for topological sorting
	inDegree := make(map[SystemType]int)

	// Calculate in-degrees (number of dependencies each system has)
	for systemType := range dr.systems {
		inDegree[systemType] = 0
	}

	for systemType, deps := range dr.dependencyGraph {
		inDegree[systemType] = len(deps)
	}

	if dr.verbose {
		fmt.Println("Initial in-degrees:")
		for systemType, degree := range inDegree {
			fmt.Printf("  %s: %d\n", systemType, degree)
		}
	}

	// Find systems with no incoming edges
	queue := make([]SystemType, 0)
	for systemType, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, systemType)
		}
	}

	if dr.verbose {
		fmt.Printf("Initial queue (systems with no dependencies): %v\n", queue)
	}

	result := make([]SystemType, 0)

	// Process queue
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		if dr.verbose {
			fmt.Printf("Processing: %s (queue size: %d)\n", current, len(queue))
		}

		// Reduce in-degree for all neighbors
		for _, neighbor := range dr.reverseGraph[current] {
			inDegree[neighbor]--
			if dr.verbose {
				fmt.Printf("  Reducing in-degree for %s: %d -> %d\n", neighbor, inDegree[neighbor]+1, inDegree[neighbor])
			}
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
				if dr.verbose {
					fmt.Printf("  Adding %s to queue\n", neighbor)
				}
			}
		}
	}

	// Check for cycles
	if len(result) != len(dr.systems) {
		if dr.verbose {
			fmt.Printf("CYCLE DETECTED! Processed %d systems, but %d systems exist\n", len(result), len(dr.systems))
		}

		// Find systems that weren't processed (part of a cycle)
		processed := make(map[SystemType]bool)
		for _, systemType := range result {
			processed[systemType] = true
		}

		cycleSystems := make([]string, 0)
		for systemType := range dr.systems {
			if !processed[systemType] {
				cycleSystems = append(cycleSystems, string(systemType))
			}
		}

		if dr.verbose {
			fmt.Printf("Systems in cycle: %v\n", cycleSystems)
			fmt.Println("Dependency details for systems in cycle:")
			for _, systemName := range cycleSystems {
				systemType := SystemType(systemName)
				if info, exists := dr.systems[systemType]; exists {
					fmt.Printf("  %s:\n", systemType)
					fmt.Printf("    Dependencies: %v\n", info.Dependencies)
					fmt.Printf("    Requires: %v\n", info.Requires)
					fmt.Printf("    Provides: %v\n", info.Provides)
					if deps, exists := dr.dependencyGraph[systemType]; exists {
						fmt.Printf("    Resolved dependencies: %v\n", deps)
					}
				}
			}
		}

		return nil, fmt.Errorf("circular dependency detected among systems: %v", cycleSystems)
	}

	if dr.verbose {
		fmt.Printf("Topological sort completed successfully: %v\n", result)
	}

	return result, nil
}

// UpdateAll updates all systems in the correct order
func (sm *SystemManager) UpdateAll(deltaTime float64, entities []Entity, eventDispatcher *events.EventDispatcher) {
	for _, systemType := range sm.updateOrder {
		if info, exists := sm.systems[systemType]; exists {
			info.System.Update(deltaTime, entities, eventDispatcher)
		}
	}
}

// DrawAll draws all drawable systems in the correct order
func (sm *SystemManager) DrawAll(screen interface{}, entities []Entity) {
	for _, systemType := range sm.drawOrder {
		if info, exists := sm.systems[systemType]; exists {
			// Use type assertion to check if system implements DrawableSystem
			if drawable, ok := info.System.(DrawableSystem); ok {
				// Type assert screen to *ebiten.Image for type safety
				if ebitenScreen, ok := screen.(*ebiten.Image); ok {
					drawable.Draw(ebitenScreen, entities)
				}
			}
		}
	}
}

// GetSystem returns a system by type
func (sm *SystemManager) GetSystem(systemType SystemType) (System, bool) {
	if info, exists := sm.systems[systemType]; exists {
		return info.System, true
	}
	return nil, false
}

// GetSystemInfo returns system info by type
func (sm *SystemManager) GetSystemInfo(systemType SystemType) (*SystemInfo, bool) {
	if info, exists := sm.systems[systemType]; exists {
		return info, true
	}
	return nil, false
}

// GetUpdateOrder returns the current update order
func (sm *SystemManager) GetUpdateOrder() []SystemType {
	return sm.updateOrder
}

// GetDrawOrder returns the current draw order
func (sm *SystemManager) GetDrawOrder() []SystemType {
	return sm.drawOrder
}

// PrintExecutionOrder prints the current execution order for debugging
func (sm *SystemManager) PrintExecutionOrder() {
	fmt.Println("=== System Execution Order ===")
	fmt.Println("Update Order:")
	for i, systemType := range sm.updateOrder {
		fmt.Printf("  %d. %s\n", i+1, systemType)
	}
	fmt.Println("Draw Order:")
	for i, systemType := range sm.drawOrder {
		fmt.Printf("  %d. %s\n", i+1, systemType)
	}
	fmt.Println("==============================")
}

// PrintDependencyGraph prints the current dependency graph for debugging
func (sm *SystemManager) PrintDependencyGraph() {
	fmt.Println("=== Dependency Graph ===")
	for systemType, info := range sm.systems {
		fmt.Printf("%s:\n", systemType)
		fmt.Printf("  Dependencies: %v\n", info.Dependencies)
		fmt.Printf("  Requires: %v\n", info.Requires)
		fmt.Printf("  Provides: %v\n", info.Provides)
		if deps, exists := sm.resolver.dependencyGraph[systemType]; exists {
			fmt.Printf("  Resolved deps: %v\n", deps)
		}
		fmt.Println()
	}
	fmt.Println("==========================")
}
