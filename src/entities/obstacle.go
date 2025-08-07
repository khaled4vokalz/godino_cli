package entities

import (
	"cli-dino-game/src/engine"
)

// ObstacleType represents different types of obstacles
type ObstacleType int

const (
	CactusSmall ObstacleType = iota
	CactusMedium
	CactusLarge
)

// String returns the string representation of ObstacleType
func (ot ObstacleType) String() string {
	switch ot {
	case CactusSmall:
		return "CactusSmall"
	case CactusMedium:
		return "CactusMedium"
	case CactusLarge:
		return "CactusLarge"
	default:
		return "Unknown"
	}
}

// Obstacle represents an obstacle that the dinosaur must avoid
type Obstacle struct {
	// Position and movement
	X     float64 // Horizontal position
	Y     float64 // Vertical position (ground level)
	Speed float64 // Movement speed from right to left

	// Obstacle properties
	ObstType ObstacleType // Type of obstacle
	Width    float64      // Width for collision detection
	Height   float64      // Height for collision detection

	// State
	Active bool // Whether the obstacle is active (on screen)
}

// NewObstacle creates a new obstacle of the specified type
func NewObstacle(obstType ObstacleType, x, groundLevel float64, config *engine.Config) *Obstacle {
	obstacle := &Obstacle{
		X:        x,
		Y:        groundLevel,
		Speed:    config.ObstacleSpeed,
		ObstType: obstType,
		Active:   true,
	}

	// Set dimensions based on obstacle type
	switch obstType {
	case CactusSmall:
		obstacle.Width = 3.0
		obstacle.Height = 3.0
		obstacle.Y = groundLevel - obstacle.Height // Adjust Y to sit on ground
	case CactusMedium:
		obstacle.Width = 3.0
		obstacle.Height = 4.0
		obstacle.Y = groundLevel - obstacle.Height
	case CactusLarge:
		obstacle.Width = 5.0
		obstacle.Height = 5.0
		obstacle.Y = groundLevel - obstacle.Height
	}

	return obstacle
}

// Update updates the obstacle's position and state
func (o *Obstacle) Update(deltaTime float64) {
	if !o.Active {
		return
	}

	// Move obstacle from right to left
	o.X -= o.Speed * deltaTime

	// Deactivate obstacle if it moves off-screen (left edge)
	if o.X+o.Width < 0 {
		o.Active = false
	}
}

// GetBounds returns the collision rectangle for the obstacle
func (o *Obstacle) GetBounds() engine.Rectangle {
	return engine.Rectangle{
		X:      o.X,
		Y:      o.Y,
		Width:  o.Width,
		Height: o.Height,
	}
}

// GetASCIIArt returns the ASCII art representation of the obstacle
func (o *Obstacle) GetASCIIArt() []string {
	return o.GetASCIIArtWithConfig(false) // Default to ASCII
}

// GetASCIIArtWithConfig returns the ASCII art with Unicode/ASCII choice
func (o *Obstacle) GetASCIIArtWithConfig(useUnicode bool) []string {
	if useUnicode {
		switch o.ObstType {
		case CactusSmall:
			return []string{
				" ╷",
				" │",
				"═══",
			}
		case CactusMedium:
			return []string{
				" ╷ ",
				"═╪═",
				" │ ",
				"═══",
			}
		case CactusLarge:
			return []string{
				"  ╷  ",
				"══╪══",
				"  │  ",
				"  │  ",
				"═════",
			}
		default:
			return []string{
				" ╷",
				" │",
				"═══",
			}
		}
	} else {
		switch o.ObstType {
		case CactusSmall:
			return []string{
				" #",
				" #",
				"###",
			}
		case CactusMedium:
			return []string{
				" # ",
				"###",
				" # ",
				"###",
			}
		case CactusLarge:
			return []string{
				"  #  ",
				"#####",
				"  #  ",
				"  #  ",
				"#####",
			}
		default:
			return []string{
				" #",
				" #",
				"###",
			}
		}
	}
}

// GetPosition returns the current position of the obstacle
func (o *Obstacle) GetPosition() (float64, float64) {
	return o.X, o.Y
}

// SetPosition sets the obstacle's position
func (o *Obstacle) SetPosition(x, y float64) {
	o.X = x
	o.Y = y
}

// IsActive returns whether the obstacle is currently active
func (o *Obstacle) IsActive() bool {
	return o.Active
}

// Deactivate marks the obstacle as inactive
func (o *Obstacle) Deactivate() {
	o.Active = false
}

// GetType returns the obstacle type
func (o *Obstacle) GetType() ObstacleType {
	return o.ObstType
}

// GetSpeed returns the obstacle's movement speed
func (o *Obstacle) GetSpeed() float64 {
	return o.Speed
}

// SetSpeed sets the obstacle's movement speed
func (o *Obstacle) SetSpeed(speed float64) {
	o.Speed = speed
}

// IsOffScreen returns true if the obstacle has moved completely off the left side of the screen
func (o *Obstacle) IsOffScreen() bool {
	return o.X+o.Width < 0
}

// GetDimensions returns the width and height of the obstacle
func (o *Obstacle) GetDimensions() (float64, float64) {
	return o.Width, o.Height
}
