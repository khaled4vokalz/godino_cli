package engine

import (
	"time"
)

// GameEngine manages the overall game state and coordinates all game systems
type GameEngine struct {
	// Game state
	state     GameState
	running   bool
	gameOver  bool
	startTime time.Time

	// Configuration
	config *Config

	// Collision detection
	collisionDetector  *CollisionDetector
	collisionTolerance float64 // For more forgiving gameplay

	// Game timing
	lastUpdate time.Time
	deltaTime  float64
}

// NewGameEngine creates a new game engine with the specified configuration
func NewGameEngine(config *Config) *GameEngine {
	return &GameEngine{
		state:              StateMenu,
		running:            false,
		gameOver:           false,
		config:             config,
		collisionDetector:  NewCollisionDetector(),
		collisionTolerance: 1.0, // Default tolerance for fair gameplay
		lastUpdate:         time.Now(),
	}
}

// GetState returns the current game state
func (ge *GameEngine) GetState() GameState {
	return ge.state
}

// SetState sets the game state
func (ge *GameEngine) SetState(state GameState) {
	ge.state = state
}

// IsRunning returns whether the game is currently running
func (ge *GameEngine) IsRunning() bool {
	return ge.running
}

// IsGameOver returns whether the game is in game over state
func (ge *GameEngine) IsGameOver() bool {
	return ge.gameOver
}

// Start starts the game engine
func (ge *GameEngine) Start() {
	ge.running = true
	ge.gameOver = false
	ge.state = StatePlaying
	ge.startTime = time.Now()
	ge.lastUpdate = time.Now()
}

// Stop stops the game engine
func (ge *GameEngine) Stop() {
	ge.running = false
}

// Reset resets the game engine to initial state
func (ge *GameEngine) Reset() {
	ge.state = StateMenu
	ge.running = false
	ge.gameOver = false
	ge.startTime = time.Time{}
	ge.lastUpdate = time.Now()
}

// Update updates the game engine timing
func (ge *GameEngine) Update() {
	now := time.Now()
	ge.deltaTime = now.Sub(ge.lastUpdate).Seconds()
	ge.lastUpdate = now
}

// GetDeltaTime returns the time elapsed since the last update
func (ge *GameEngine) GetDeltaTime() float64 {
	return ge.deltaTime
}

// GetConfig returns the game configuration
func (ge *GameEngine) GetConfig() *Config {
	return ge.config
}

// SetCollisionTolerance sets the collision tolerance for more forgiving gameplay
func (ge *GameEngine) SetCollisionTolerance(tolerance float64) {
	ge.collisionTolerance = tolerance
}

// GetCollisionTolerance returns the current collision tolerance
func (ge *GameEngine) GetCollisionTolerance() float64 {
	return ge.collisionTolerance
}

// EnableCollisionDebug enables debug mode for collision detection
func (ge *GameEngine) EnableCollisionDebug(enabled bool) {
	ge.collisionDetector.SetDebugMode(enabled)
}

// CheckCollision checks for collision between two rectangles
func (ge *GameEngine) CheckCollision(rect1, rect2 Rectangle) bool {
	if ge.collisionTolerance > 0 {
		return ge.collisionDetector.CheckCollisionWithTolerance(rect1, rect2, ge.collisionTolerance)
	}
	return ge.collisionDetector.CheckCollision(rect1, rect2)
}

// GetCollisionInfo returns detailed collision information
func (ge *GameEngine) GetCollisionInfo(rect1, rect2 Rectangle) CollisionInfo {
	return ge.collisionDetector.GetCollisionInfo(rect1, rect2)
}

// TriggerGameOver triggers the game over state
func (ge *GameEngine) TriggerGameOver() {
	ge.gameOver = true
	ge.state = StateGameOver
	ge.running = false
}

// GetGameDuration returns how long the current game has been running
func (ge *GameEngine) GetGameDuration() time.Duration {
	if ge.startTime.IsZero() {
		return 0
	}
	return time.Since(ge.startTime)
}

// Restart restarts the game from game over state
func (ge *GameEngine) Restart() {
	ge.Reset()
	ge.Start()
}
