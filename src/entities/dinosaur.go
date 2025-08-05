package entities

import (
	"cli-dino-game/src/engine"
	"time"
)

// Dinosaur represents the player-controlled dinosaur character
type Dinosaur struct {
	// Position and movement
	X         float64 // Horizontal position
	Y         float64 // Vertical position
	VelocityY float64 // Vertical velocity for jumping

	// State management
	IsJumping   bool    // Whether the dinosaur is currently jumping
	IsRunning   bool    // Whether the dinosaur is in running state
	AnimFrame   int     // Current animation frame for running
	GroundLevel float64 // Y position of the ground

	// Animation timing
	lastAnimUpdate time.Time
	animSpeed      time.Duration

	// Dimensions for collision detection
	Width  float64
	Height float64
}

// NewDinosaur creates a new dinosaur with default values
func NewDinosaur(groundLevel float64) *Dinosaur {
	return &Dinosaur{
		X:              15.0, // Fixed position on screen
		Y:              groundLevel,
		VelocityY:      0.0,
		IsJumping:      false,
		IsRunning:      true,
		AnimFrame:      0,
		GroundLevel:    groundLevel,
		lastAnimUpdate: time.Now(),
		animSpeed:      time.Millisecond * 200, // Animation frame duration
		Width:          8.0,                    // Width of dinosaur sprite
		Height:         6.0,                    // Height of dinosaur sprite
	}
}

// Update updates the dinosaur's state and position
func (d *Dinosaur) Update(deltaTime float64) {
	// Dinosaur stays in a fixed horizontal position
	// The world/obstacles will scroll past the dinosaur instead
	// No horizontal movement needed - X position remains constant

	// Update running animation if on ground
	if !d.IsJumping && d.IsRunning {
		now := time.Now()
		if now.Sub(d.lastAnimUpdate) >= d.animSpeed {
			d.AnimFrame = (d.AnimFrame + 1) % 2 // Alternate between 2 frames
			d.lastAnimUpdate = now
		}
	}

	// Future: Vertical movement (jumping) will be handled here
	// when jump mechanics are implemented in later tasks
}

// GetBounds returns the collision rectangle for the dinosaur
func (d *Dinosaur) GetBounds() engine.Rectangle {
	return engine.Rectangle{
		X:      d.X,
		Y:      d.Y,
		Width:  d.Width,
		Height: d.Height,
	}
}

// GetASCIIArt returns the ASCII art representation of the dinosaur
func (d *Dinosaur) GetASCIIArt() []string {
	if d.IsJumping {
		// Jumping dinosaur sprite
		return []string{
			"        ████████",
			"        ██    ██",
			"        ████████",
			"        ██████  ",
			"██      ██      ",
			"████    ██      ",
		}
	}

	// Running animation - alternate between two frames
	if d.AnimFrame == 0 {
		// Running frame 1
		return []string{
			"        ████████",
			"        ██    ██",
			"        ████████",
			"        ██████  ",
			"██      ██      ",
			"████    ██  ██  ",
		}
	} else {
		// Running frame 2
		return []string{
			"        ████████",
			"        ██    ██",
			"        ████████",
			"        ██████  ",
			"██      ██      ",
			"████    ██    ██",
		}
	}
}

// GetPosition returns the current position of the dinosaur
func (d *Dinosaur) GetPosition() (float64, float64) {
	return d.X, d.Y
}

// SetPosition sets the dinosaur's position
func (d *Dinosaur) SetPosition(x, y float64) {
	d.X = x
	d.Y = y
}

// IsOnGround returns true if the dinosaur is on the ground
func (d *Dinosaur) IsOnGround() bool {
	return d.Y >= d.GroundLevel && !d.IsJumping
}
