package engine

import (
	"cli-dino-game/src/score"
	"time"
)

// GameEngine manages the overall game state and coordinates all game systems
type GameEngine struct {
	// Game state
	state         GameState
	previousState GameState
	running       bool
	gameOver      bool
	startTime     time.Time
	initialized   bool

	// Configuration
	config *Config

	// Game scoring
	gameScore *score.Score

	// Collision detection
	collisionDetector  *CollisionDetector
	collisionTolerance float64 // For more forgiving gameplay

	// Game timing
	lastUpdate time.Time
	deltaTime  float64

	// State transition callbacks
	onStateChange func(from, to GameState)
}

// NewGameEngine creates a new game engine with the specified configuration
func NewGameEngine(config *Config) *GameEngine {
	gameScore := score.NewScore()
	// Load high score from persistent storage
	gameScore.LoadHighScoreInto()

	return &GameEngine{
		state:              StateMenu,
		previousState:      StateMenu,
		running:            false,
		gameOver:           false,
		initialized:        false,
		config:             config,
		gameScore:          gameScore,
		collisionDetector:  NewCollisionDetector(),
		collisionTolerance: 0.5, // Default tolerance for fair but not overly forgiving gameplay
		lastUpdate:         time.Now(),
	}
}

// GetState returns the current game state
func (ge *GameEngine) GetState() GameState {
	return ge.state
}

// SetState sets the game state and handles state transitions
func (ge *GameEngine) SetState(state GameState) {
	if ge.state == state {
		return // No change needed
	}

	previousState := ge.state
	ge.previousState = previousState
	ge.state = state

	// Handle state-specific logic
	ge.handleStateTransition(previousState, state)

	// Call state change callback if set
	if ge.onStateChange != nil {
		ge.onStateChange(previousState, state)
	}
}

// GetPreviousState returns the previous game state
func (ge *GameEngine) GetPreviousState() GameState {
	return ge.previousState
}

// SetStateChangeCallback sets a callback function to be called when state changes
func (ge *GameEngine) SetStateChangeCallback(callback func(from, to GameState)) {
	ge.onStateChange = callback
}

// handleStateTransition handles logic when transitioning between states
func (ge *GameEngine) handleStateTransition(from, to GameState) {
	switch to {
	case StateMenu:
		ge.running = false
		ge.gameOver = false
	case StatePlaying:
		if !ge.initialized {
			ge.initialize()
		}
		ge.running = true
		ge.gameOver = false
		if from != StatePlaying {
			ge.startTime = time.Now()
			ge.ResetScore() // Reset score when starting a new game
		}
	case StateGameOver:
		ge.running = false
		ge.gameOver = true
		// Finalize score when game ends
		if ge.gameScore != nil {
			ge.FinalizeScore()
		}
	}
}

// IsRunning returns whether the game is currently running
func (ge *GameEngine) IsRunning() bool {
	return ge.running
}

// IsGameOver returns whether the game is in game over state
func (ge *GameEngine) IsGameOver() bool {
	return ge.gameOver
}

// Initialize initializes the game engine for gameplay
func (ge *GameEngine) initialize() {
	ge.initialized = true
	ge.lastUpdate = time.Now()
	// Additional initialization logic can be added here
}

// IsInitialized returns whether the game engine has been initialized
func (ge *GameEngine) IsInitialized() bool {
	return ge.initialized
}

// Start starts the game engine
func (ge *GameEngine) Start() {
	ge.SetState(StatePlaying)
}

// Stop stops the game engine and transitions to menu
func (ge *GameEngine) Stop() {
	ge.SetState(StateMenu)
}

// Cleanup performs cleanup operations when shutting down the game
func (ge *GameEngine) Cleanup() {
	ge.running = false
	ge.gameOver = false
	ge.initialized = false
	ge.onStateChange = nil
	// Additional cleanup logic can be added here
}

// Reset resets the game engine to initial state
func (ge *GameEngine) Reset() {
	ge.SetState(StateMenu)
	ge.startTime = time.Time{}
	ge.lastUpdate = time.Now()
	ge.initialized = false
}

// Update updates the game engine timing and score
func (ge *GameEngine) Update() {
	now := time.Now()
	ge.deltaTime = now.Sub(ge.lastUpdate).Seconds()
	ge.lastUpdate = now

	// Update score if game is playing
	ge.UpdateScore()
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
	ge.SetState(StateGameOver)
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
	if ge.state == StateGameOver {
		ge.Reset()
		ge.Start()
	}
}

// CanTransitionTo checks if a state transition is valid
func (ge *GameEngine) CanTransitionTo(newState GameState) bool {
	switch ge.state {
	case StateMenu:
		return newState == StatePlaying
	case StatePlaying:
		return newState == StateGameOver || newState == StateMenu
	case StateGameOver:
		return newState == StateMenu || newState == StatePlaying
	default:
		return false
	}
}

// TransitionTo attempts to transition to a new state if valid
func (ge *GameEngine) TransitionTo(newState GameState) bool {
	if ge.CanTransitionTo(newState) {
		ge.SetState(newState)
		return true
	}
	return false
}

// GetScore returns the game score instance
func (ge *GameEngine) GetScore() *score.Score {
	return ge.gameScore
}

// UpdateScore updates the game score based on elapsed time
func (ge *GameEngine) UpdateScore() {
	if ge.state == StatePlaying && ge.gameScore != nil {
		ge.gameScore.Update(ge.deltaTime)
	}
}

// AddObstacleBonus adds bonus points for passing an obstacle
func (ge *GameEngine) AddObstacleBonus() {
	if ge.gameScore != nil {
		ge.gameScore.AddObstacleBonus()
	}
}

// ResetScore resets the score for a new game
func (ge *GameEngine) ResetScore() {
	if ge.gameScore != nil {
		ge.gameScore.Reset()
	}
}

// FinalizeScore finalizes the score at game end and handles high score persistence
func (ge *GameEngine) FinalizeScore() (bool, error) {
	if ge.gameScore != nil {
		return ge.gameScore.FinalizeScore()
	}
	return false, nil
}

// GetCurrentScore returns the current score value
func (ge *GameEngine) GetCurrentScore() int {
	if ge.gameScore != nil {
		return ge.gameScore.GetCurrent()
	}
	return 0
}

// GetHighScore returns the high score value
func (ge *GameEngine) GetHighScore() int {
	if ge.gameScore != nil {
		return ge.gameScore.GetHigh()
	}
	return 0
}

// IsNewHighScore checks if the current score is a new high score
func (ge *GameEngine) IsNewHighScore() bool {
	if ge.gameScore != nil {
		return ge.gameScore.IsNewHighScore()
	}
	return false
}
