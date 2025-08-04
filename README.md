# LBaaS Packet Catcher 🎮

A 2D game parody of "Eggs Electronica Wolf and Hare" where you play as a **Load Balancer as a Service (LBaaS)** catching falling network packets!

## 🎯 Game Concept

Instead of a wolf catching eggs, you control a **load balancer** that must catch falling **network packets** before they hit the ground. Each packet type has its own color and represents different network protocols:

- 🔴 **HTTP** - Red packets
- 🟢 **HTTPS** - Green packets  
- 🔵 **TCP** - Blue packets
- 🟡 **UDP** - Yellow packets
- 🟣 **WebSocket** - Magenta packets

## 🎮 Controls

- **A** or **Left Arrow** - Move load balancer left
- **D** or **Right Arrow** - Move load balancer right
- **Mouse (Left Click)** - Move load balancer to mouse position
- **Ctrl+X** - Exit game
- **R** - Restart game (when game over)
- **UP/DOWN** - Select game mode in menu
- **ENTER** - Start game

## 🏆 Scoring

- **SLA-based scoring**: Maintain Service Level Agreement targets
- **Dynamic SLA display**: Updates in real-time
- **SLA checks**: Every 100 packets against target
- **Game modes**: Different SLA targets and error budgets
- **Maximum 10,000 packets**: Game ends when limit reached

## 🚀 How to Run

1. **Install Go** (if not already installed)
   ```bash
   # On Fedora/RHEL
   sudo dnf install golang
   ```

2. **Download dependencies**
   ```bash
   go mod tidy
   ```

3. **Build and run the game**
   ```bash
   ./build.sh
   ./lbbaspack
   ```

## 🎨 Features

### Core Gameplay
- **Progressive difficulty**: Packet spawn rate increases over time
- **Multiple packet types**: Different network protocols with unique colors
- **SLA System**: Service Level Agreement targets with different game modes
- **Dynamic SLA Display**: Real-time SLA percentage updates
- **Mouse Control**: Click to move load balancer to mouse position
- **Fullscreen Mode**: Immersive gaming experience
- **Ctrl+X Exit**: Quick exit with keyboard shortcut

### Power-ups & Special Abilities
- **Speed Boost** (Yellow): Doubles load balancer movement speed
- **Wide Catch** (Cyan): Increases catch area by 50%
- **Multi-Catch** (Magenta): Can catch multiple packets simultaneously
- **Time Slow** (Blue): Slows down falling packets
- **Shield** (Green): Protects against missed packets
- **Auto-Balancer** (Orange): Automatically distributes packets to least-loaded backend

### Backend Visualization
- **Backend Visualization**: See packets flow to backend servers
- **Round-robin Load Balancing**: Packets distributed evenly across backends
- **Smart Load Balancing**: Auto-balancer finds least-loaded backend
- **Packet Counters**: Real-time packet counts per backend

### Visual Effects
- **Particle effects**: Visual feedback when catching packets
- **Power-up indicators**: Visual feedback for active abilities
- **Gradient background**: More polished visual appearance
- **Level progression**: Levels based on packets caught

## 🎯 Game Modes

1. **Mission Critical** - 99.95% SLA target (3 error budget)
2. **Business Critical** - 99.5% SLA target (10 error budget)
3. **Business Operational** - 99% SLA target (25 error budget)
4. **Office Productivity** - 95% SLA target (50 error budget)
5. **Best Effort** - 90% SLA target (100 error budget)

## 🏗️ Current Project Structure

The game is built with a modular architecture following SOLID/DRY principles:

```
lbbaspack/
├── main.go           # Entry point and game loop
├── game.go           # Game state and core logic
├── packet.go         # Packet spawning and management
├── backend.go        # Backend management and load balancing
├── powerup.go        # Power-up system and effects
├── particle.go       # Particle effects system
├── ui.go            # UI drawing functions
├── constants.go      # Game constants and configuration
├── build.sh          # Build script
└── README.md         # This file
```

### Module Responsibilities

- **main.go**: Ebiten setup, game loop, and high-level coordination
- **game.go**: Game state management and core game logic
- **packet.go**: Packet creation, spawning, and movement logic
- **backend.go**: Backend management and load balancing algorithms
- **powerup.go**: Power-up spawning, activation, and effect management
- **particle.go**: Particle system for visual effects
- **ui.go**: User interface drawing and display logic
- **constants.go**: Centralized game constants and configuration

## 🔄 Planned Advanced Refactoring

The current structure is still quite monolithic. We're planning a more sophisticated architecture following game development best practices:

### 🎯 Target Architecture: Entity-Component-System (ECS)

```
lbbaspack/
├── main.go                    # Entry point
├── engine/                    # Core game engine
│   ├── game.go               # Game state management
│   ├── systems/              # ECS systems
│   │   ├── input.go          # Input handling system
│   │   ├── movement.go       # Movement system
│   │   ├── collision.go      # Collision detection system
│   │   ├── spawning.go       # Entity spawning system
│   │   ├── rendering.go      # Rendering system
│   │   └── ui.go            # UI rendering system
│   ├── components/           # ECS components
│   │   ├── transform.go      # Position, rotation, scale
│   │   ├── sprite.go         # Visual representation
│   │   ├── physics.go        # Velocity, acceleration
│   │   ├── health.go         # Health and damage
│   │   ├── powerup.go        # Power-up effects
│   │   └── ai.go            # AI behavior
│   ├── entities/             # Entity definitions
│   │   ├── loadbalancer.go   # Load balancer entity
│   │   ├── packet.go         # Packet entity
│   │   ├── powerup.go        # Power-up entity
│   │   ├── backend.go        # Backend entity
│   │   └── particle.go       # Particle entity
│   └── events/               # Event system
│       ├── events.go         # Event definitions
│       └── dispatcher.go     # Event handling
├── systems/                   # Game-specific systems
│   ├── sla.go               # SLA calculation system
│   ├── loadbalancing.go     # Load balancing logic
│   ├── powerup_manager.go   # Power-up management
│   └── game_state.go        # Game state transitions
├── ui/                       # User interface
│   ├── menu.go              # Menu system
│   ├── hud.go               # Heads-up display
│   └── game_over.go         # Game over screen
├── config/                   # Configuration
│   ├── game_config.go       # Game settings
│   ├── powerup_config.go    # Power-up definitions
│   └── packet_config.go     # Packet type definitions
└── utils/                    # Utilities
    ├── math.go              # Math utilities
    ├── colors.go            # Color definitions
    └── constants.go         # Game constants
```

### 🎮 Benefits of ECS Architecture

1. **Separation of Concerns**: Each system handles one aspect (movement, rendering, etc.)
2. **Composition over Inheritance**: Entities are composed of components
3. **Easy Extension**: Add new components/systems without modifying existing code
4. **Performance**: Systems can be optimized independently
5. **Testability**: Each system can be unit tested in isolation
6. **Parallel Processing**: Systems can run in parallel where possible

### 🔧 Key Improvements

- **Event-Driven Architecture**: Loose coupling between systems
- **State Management**: Clear game state transitions
- **Resource Management**: Proper asset loading and cleanup
- **Input Abstraction**: Platform-independent input handling
- **Rendering Pipeline**: Efficient rendering with batching
- **Physics System**: Proper collision detection and response
- **Audio System**: Sound effects and music management

## 🛠️ Technical Details

- Built with **Ebitengine v2** (Go 2D game engine)
- Written in **Go 1.24.4**
- **Modular architecture** following SOLID/DRY principles
- Window size: 800x600 pixels (fullscreen)
- Cross-platform support (Linux, Windows, macOS)

## 🎯 Game Mechanics

- **SLA-based gameplay**: Must maintain SLA above target
- **Dynamic difficulty**: Packet speed increases with lost packets
- **Backend visualization**: See packets flow to backend servers
- **Mouse control**: Precise positioning with mouse
- **Fullscreen experience**: Immersive gaming
- **Strategic depth**: Balance between catching packets and maintaining SLA
- **Power-up system**: Collect special abilities for enhanced gameplay

Enjoy catching those network packets! 🌐📦 