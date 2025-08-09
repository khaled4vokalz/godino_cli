# CLI Dino Game 🦕

A terminal-based dinosaur runner game inspired by the Chrome T-Rex game, built in Go as a **learning project for agent-based development**.

> **Note**: This project was developed through iterative collaboration with AI agents, demonstrating modern software development workflows using AI assistance for coding, debugging, and feature implementation.

## Demo

https://github.com/user-attachments/assets/19bfa6de-c849-47d2-9939-a43b364fd4a4

## Features

- **Jump over obstacles** with `Space` or `↑`
- **Progressive difficulty** - speed and obstacles increase over time
- **Multiple obstacle types** - cacti and birds (birds appear after 15s)
- **Beautiful graphics** - Unicode characters with ASCII fallback
- **Smooth animations** - running, jumping, and background scrolling
- **Continuous hill backgrounds** - generated using sine waves
- **Collision detection** - precise AABB with configurable tolerance

## Quick Start

### Prerequisites

- Go 1.16 or higher
- Terminal with Unicode support (recommended for best visuals)

```bash
# Build and run
go build
./cli-dino-game

# ASCII mode for compatibility
./cli-dino-game -ascii
```

## Controls

- **Start/Jump**: `Space` or `↑`
- **Restart**: `R` (after game over)
- **Quit**: `Q` or `Ctrl+C`

## Agent-Based Development Journey

This project showcases how modern development can leverage AI agents for:

- **Feature Implementation**: Jump mechanics, collision detection, obstacle spawning
- **Bug Fixing**: Bird collision issues, difficulty balancing, collision tolerance
- **Visual Improvements**: Background rendering, animation systems, Unicode graphics
- **Code Architecture**: Clean separation of concerns across multiple packages
- **Testing & Debugging**: Systematic problem-solving and iterative improvements

### Key Development Iterations

1. **Core Game Loop** - Basic dino, jumping, and obstacle mechanics
2. **Collision System** - Implemented AABB detection with tolerance tuning
3. **Difficulty Balancing** - Refined spawn rates and progression curves
4. **Visual Polish** - Added continuous hills, improved sprites, smooth animations
5. **Bug Fixes** - Resolved bird collision detection and timing issues

## Project Structure

```text
cli-dino-game/
├── main.go                 # Game loop and coordination
├── src/
│   ├── background/         # Hills and cloud generation
│   ├── engine/            # Game state and collision detection
│   ├── entities/          # Dinosaur and obstacles
│   ├── input/             # Keyboard handling
│   ├── render/            # Terminal graphics
│   ├── score/             # Scoring system
│   └── spawner/           # Obstacle generation
└── go.mod
```

## Technical Highlights

- **Real-time terminal rendering** using [termbox-go](https://github.com/nsf/termbox-go)
- **Physics simulation** with gravity and velocity
- **Procedural background generation** using mathematical functions
- **Configurable game parameters** for fine-tuning
- **Cross-platform compatibility** (macOS, Linux, Windows)

## Development Notes

Built entirely through agent-assisted development, this project demonstrates:

- Effective human-AI collaboration patterns
- Iterative problem-solving approaches
- Code quality through AI-guided refactoring
- Systematic debugging methodologies

---

**A fun game and a learning experience in AI-assisted development!** 🎮🤖
