package entities

import (
	"cli-dino-game/src/engine"
	"time"
)

// ObstacleType represents different types of obstacles
type ObstacleType int

const (
	CactusSmall ObstacleType = iota
	CactusMedium
	CactusLarge
	BirdLow
	BirdMid
	BirdHigh
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
	case BirdLow:
		return "BirdLow"
	case BirdMid:
		return "BirdMid"
	case BirdHigh:
		return "BirdHigh"
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

	// Animation (for birds)
	AnimFrame      int           // Current animation frame
	lastAnimUpdate time.Time     // Last animation update time
	animSpeed      time.Duration // Animation frame duration

	// State
	Active bool // Whether the obstacle is active (on screen)
}

// NewObstacle creates a new obstacle of the specified type
func NewObstacle(obstType ObstacleType, x, groundLevel float64, config *engine.Config) *Obstacle {
	obstacle := &Obstacle{
		X:              x,
		Y:              groundLevel,
		Speed:          config.ObstacleSpeed,
		ObstType:       obstType,
		Active:         true,
		AnimFrame:      0,
		lastAnimUpdate: time.Now(),
		animSpeed:      time.Millisecond * 200, // Wing flapping speed
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
	case BirdLow:
		obstacle.Width = 4.0  // Use full sprite width for more realistic collision
		obstacle.Height = 2.0 // Use full sprite height for more realistic collision
		// groundLevel passed here is actualGroundY (where ground line is drawn)
		// Dinosaur is at (groundLevel - dinosaur.Height), which is (groundLevel - 4)
		// So dinosaur occupies Y from (groundLevel - 4) to groundLevel
		// BirdLow should be at dinosaur's lower body level
		obstacle.Y = groundLevel - 3.0 // Bird at dinosaur lower body level
	case BirdMid:
		obstacle.Width = 4.0           // Use full sprite width
		obstacle.Height = 2.0          // Use full sprite height
		obstacle.Y = groundLevel - 4.0 // Bird at dinosaur middle body level
	case BirdHigh:
		obstacle.Width = 4.0           // Use full sprite width
		obstacle.Height = 2.0          // Use full sprite height
		obstacle.Y = groundLevel - 5.0 // Bird at dinosaur head level
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

	// Update animation for birds
	if o.isBird() {
		now := time.Now()
		if now.Sub(o.lastAnimUpdate) >= o.animSpeed {
			o.AnimFrame = (o.AnimFrame + 1) % 2 // Birds have 2 animation frames
			o.lastAnimUpdate = now
		}
	}

	// Deactivate obstacle if it moves off-screen (left edge)
	if o.X+o.Width < 0 {
		o.Active = false
	}
}

// isBird returns true if this obstacle is a bird type
func (o *Obstacle) isBird() bool {
	return o.ObstType == BirdLow || o.ObstType == BirdMid || o.ObstType == BirdHigh
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
		case BirdLow, BirdMid, BirdHigh:
			if o.AnimFrame == 0 {
				return []string{
					"◦▲◦▲", // Wings up
					"▼ ▼ ",
				}
			} else {
				return []string{
					"◦▼◦▼", // Wings down
					"▲ ▲ ",
				}
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
		case BirdLow, BirdMid, BirdHigh:
			if o.AnimFrame == 0 {
				return []string{
					"^o^o", // Wings up
					" v v",
				}
			} else {
				return []string{
					"vo vo", // Wings down
					" ^ ^",
				}
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
