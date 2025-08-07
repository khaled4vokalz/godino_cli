package engine

import (
	"errors"
	"fmt"
)

// Config holds all game configuration parameters
type Config struct {
	// Screen dimensions
	ScreenWidth  int `json:"screen_width"`
	ScreenHeight int `json:"screen_height"`

	// Game timing
	TargetFPS int `json:"target_fps"`

	// Physics constants
	JumpVelocity  float64 `json:"jump_velocity"`
	Gravity       float64 `json:"gravity"`
	ObstacleSpeed float64 `json:"obstacle_speed"`

	// Gameplay parameters
	SpawnRate float64 `json:"spawn_rate"`

	// Rendering options
	UseUnicode bool `json:"use_unicode"`
}

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StateGameOver
)

// String returns the string representation of GameState
func (gs GameState) String() string {
	switch gs {
	case StateMenu:
		return "Menu"
	case StatePlaying:
		return "Playing"
	case StateGameOver:
		return "GameOver"
	default:
		return "Unknown"
	}
}

// Rectangle represents a rectangular collision boundary
type Rectangle struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// NewDefaultConfig creates a configuration with sensible default values
func NewDefaultConfig() *Config {
	return &Config{
		ScreenWidth:   80,
		ScreenHeight:  20,
		TargetFPS:     15,
		JumpVelocity:  25.0,
		Gravity:       60.0,
		ObstacleSpeed: 18.0,
		SpawnRate:     2.0,
		UseUnicode:    true, // Default to Unicode for better visuals
	}
}

// Validate checks if the configuration values are valid
func (c *Config) Validate() error {
	if c.ScreenWidth <= 0 {
		return errors.New("screen width must be positive")
	}
	if c.ScreenHeight <= 0 {
		return errors.New("screen height must be positive")
	}
	if c.TargetFPS <= 0 {
		return errors.New("target FPS must be positive")
	}
	if c.JumpVelocity <= 0 {
		return errors.New("jump velocity must be positive")
	}
	if c.Gravity <= 0 {
		return errors.New("gravity must be positive")
	}
	if c.ObstacleSpeed <= 0 {
		return errors.New("obstacle speed must be positive")
	}
	if c.SpawnRate <= 0 {
		return errors.New("spawn rate must be positive")
	}

	// Additional validation for reasonable ranges
	if c.ScreenWidth < 40 {
		return errors.New("screen width too small (minimum 40)")
	}
	if c.ScreenHeight < 10 {
		return errors.New("screen height too small (minimum 10)")
	}
	if c.TargetFPS > 120 {
		return errors.New("target FPS too high (maximum 120)")
	}

	return nil
}

// String returns a formatted string representation of the config
func (c *Config) String() string {
	return fmt.Sprintf("Config{Screen: %dx%d, FPS: %d, Jump: %.1f, Gravity: %.1f, Speed: %.1f, Spawn: %.1f}",
		c.ScreenWidth, c.ScreenHeight, c.TargetFPS, c.JumpVelocity, c.Gravity, c.ObstacleSpeed, c.SpawnRate)
}

// Intersects checks if this rectangle intersects with another rectangle
func (r Rectangle) Intersects(other Rectangle) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}

// Contains checks if this rectangle completely contains another rectangle
func (r Rectangle) Contains(other Rectangle) bool {
	return r.X <= other.X &&
		r.Y <= other.Y &&
		r.X+r.Width >= other.X+other.Width &&
		r.Y+r.Height >= other.Y+other.Height
}

// Center returns the center point of the rectangle
func (r Rectangle) Center() (float64, float64) {
	return r.X + r.Width/2, r.Y + r.Height/2
}

// String returns a formatted string representation of the rectangle
func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle{X: %.1f, Y: %.1f, W: %.1f, H: %.1f}", r.X, r.Y, r.Width, r.Height)
}
