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
		animSpeed:      time.Millisecond * 150, // Animation frame duration for smoother 4-frame animation
		Width:          8.0,                    // Width of dinosaur sprite
		Height:         6.0,                    // Height of dinosaur sprite
	}
}

// Jump initiates a jump if the dinosaur is on the ground
func (d *Dinosaur) Jump(config *engine.Config) {
	// Only allow jumping if dinosaur is on the ground
	if d.IsOnGround() {
		d.IsJumping = true
		d.VelocityY = -config.JumpVelocity // Negative because Y increases downward
		d.IsRunning = false                // Stop running animation while jumping
	}
}

// Update updates the dinosaur's state and position
func (d *Dinosaur) Update(deltaTime float64, config *engine.Config) {
	// Dinosaur stays in a fixed horizontal position
	// The world/obstacles will scroll past the dinosaur instead
	// No horizontal movement needed - X position remains constant

	// Handle jumping physics
	if d.IsJumping {
		// Update vertical position based on current velocity (before applying gravity)
		d.Y += d.VelocityY * deltaTime

		// Apply gravity to velocity (for next frame)
		d.VelocityY += config.Gravity * deltaTime

		// Check for landing
		if d.Y >= d.GroundLevel {
			// Land on ground
			d.Y = d.GroundLevel
			d.VelocityY = 0.0
			d.IsJumping = false
			d.IsRunning = true // Resume running animation
		}
	} else {
		// Update running animation if on ground
		if d.IsRunning {
			now := time.Now()
			if now.Sub(d.lastAnimUpdate) >= d.animSpeed {
				d.AnimFrame = (d.AnimFrame + 1) % 4 // Cycle through 4 frames
				d.lastAnimUpdate = now
			}
		}
	}
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

	// Running animation with 4 frames for smoother animation
	switch d.AnimFrame {
	case 0:
		// Running frame 1 - left leg forward
		return []string{
			"        ████████",
			"        ██    ██",
			"        ████████",
			"        ██████  ",
			"██      ██      ",
			"████    ██  ██  ",
		}
	case 1:
		// Running frame 2 - both legs center
		return []string{
			"        ████████",
			"        ██    ██",
			"        ████████",
			"        ██████  ",
			"██      ██      ",
			"████    ████    ",
		}
	case 2:
		// Running frame 3 - right leg forward
		return []string{
			"        ████████",
			"        ██    ██",
			"        ████████",
			"        ██████  ",
			"██      ██      ",
			"████    ██    ██",
		}
	case 3:
		// Running frame 4 - both legs center (slight variation)
		return []string{
			"        ████████",
			"        ██    ██",
			"        ████████",
			"        ██████  ",
			"██      ██      ",
			"████      ██    ",
		}
	default:
		// Fallback to frame 0
		return []string{
			"        ████████",
			"        ██    ██",
			"        ████████",
			"        ██████  ",
			"██      ██      ",
			"████    ██  ██  ",
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

// GetJumpHeight returns the current height above ground level
func (d *Dinosaur) GetJumpHeight() float64 {
	if d.Y < d.GroundLevel {
		return d.GroundLevel - d.Y
	}
	return 0.0
}

// GetAnimationFrame returns the current animation frame
func (d *Dinosaur) GetAnimationFrame() int {
	return d.AnimFrame
}

// SetAnimationSpeed sets the speed of the running animation
func (d *Dinosaur) SetAnimationSpeed(speed time.Duration) {
	d.animSpeed = speed
}

// GetAnimationSpeed returns the current animation speed
func (d *Dinosaur) GetAnimationSpeed() time.Duration {
	return d.animSpeed
}

// ResetAnimation resets the animation to frame 0 and updates the timer
func (d *Dinosaur) ResetAnimation() {
	d.AnimFrame = 0
	d.lastAnimUpdate = time.Now()
}

// IsAnimating returns true if the dinosaur is currently animating (running)
func (d *Dinosaur) IsAnimating() bool {
	return d.IsRunning && !d.IsJumping
}
