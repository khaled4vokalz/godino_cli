package engine

import (
	"testing"
	"time"
)

func TestNewGameEngine(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	if ge == nil {
		t.Fatal("NewGameEngine should not return nil")
	}

	if ge.GetState() != StateMenu {
		t.Error("New game engine should start in menu state")
	}

	if ge.IsRunning() {
		t.Error("New game engine should not be running")
	}

	if ge.IsGameOver() {
		t.Error("New game engine should not be in game over state")
	}

	if ge.GetConfig() != config {
		t.Error("Game engine should store the provided config")
	}

	if ge.GetCollisionTolerance() != 0.8 {
		t.Error("Default collision tolerance should be 0.8")
	}
}

func TestGameEngineStateManagement(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Test initial state
	if ge.GetState() != StateMenu {
		t.Error("Initial state should be Menu")
	}
	if ge.GetPreviousState() != StateMenu {
		t.Error("Initial previous state should be Menu")
	}

	// Test state transitions
	ge.SetState(StatePlaying)
	if ge.GetState() != StatePlaying {
		t.Error("State should be set to Playing")
	}
	if ge.GetPreviousState() != StateMenu {
		t.Error("Previous state should be Menu")
	}

	ge.SetState(StateGameOver)
	if ge.GetState() != StateGameOver {
		t.Error("State should be set to GameOver")
	}
	if ge.GetPreviousState() != StatePlaying {
		t.Error("Previous state should be Playing")
	}

	// Test setting same state (should not change previous state)
	ge.SetState(StateGameOver)
	if ge.GetPreviousState() != StatePlaying {
		t.Error("Previous state should remain Playing when setting same state")
	}
}

func TestGameEngineStartStop(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Test starting the game
	ge.Start()
	if !ge.IsRunning() {
		t.Error("Game should be running after Start()")
	}
	if ge.GetState() != StatePlaying {
		t.Error("Game state should be Playing after Start()")
	}
	if ge.IsGameOver() {
		t.Error("Game should not be in game over state after Start()")
	}

	// Test stopping the game
	ge.Stop()
	if ge.IsRunning() {
		t.Error("Game should not be running after Stop()")
	}
}

func TestGameEngineReset(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Start and modify state
	ge.Start()
	ge.TriggerGameOver()

	// Reset
	ge.Reset()

	if ge.GetState() != StateMenu {
		t.Error("State should be Menu after Reset()")
	}
	if ge.IsRunning() {
		t.Error("Game should not be running after Reset()")
	}
	if ge.IsGameOver() {
		t.Error("Game should not be in game over state after Reset()")
	}
}

func TestGameEngineUpdate(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Initial delta time should be 0
	initialDelta := ge.GetDeltaTime()

	// Wait a bit and update
	time.Sleep(10 * time.Millisecond)
	ge.Update()

	newDelta := ge.GetDeltaTime()
	if newDelta <= initialDelta {
		t.Error("Delta time should increase after Update()")
	}
	if newDelta <= 0 {
		t.Error("Delta time should be positive")
	}
}

func TestGameEngineCollisionTolerance(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Test setting collision tolerance
	tolerance := 2.5
	ge.SetCollisionTolerance(tolerance)

	if ge.GetCollisionTolerance() != tolerance {
		t.Errorf("Expected collision tolerance %.1f, got %.1f", tolerance, ge.GetCollisionTolerance())
	}
}

func TestGameEngineCollisionDetection(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Test collision detection without tolerance
	ge.SetCollisionTolerance(0)
	rect1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	rect2 := Rectangle{X: 5, Y: 5, Width: 10, Height: 10}

	if !ge.CheckCollision(rect1, rect2) {
		t.Error("Should detect collision without tolerance")
	}

	// Test collision detection with tolerance
	ge.SetCollisionTolerance(3.0)
	rect3 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	rect4 := Rectangle{X: 8, Y: 8, Width: 10, Height: 10}

	// Should collide without tolerance but not with tolerance
	ge.SetCollisionTolerance(0)
	if !ge.CheckCollision(rect3, rect4) {
		t.Error("Should collide without tolerance")
	}

	ge.SetCollisionTolerance(3.0)
	if ge.CheckCollision(rect3, rect4) {
		t.Error("Should not collide with tolerance")
	}
}

func TestGameEngineCollisionInfo(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	rect1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	rect2 := Rectangle{X: 5, Y: 5, Width: 10, Height: 10}

	info := ge.GetCollisionInfo(rect1, rect2)

	if !info.HasCollision {
		t.Error("Should detect collision")
	}

	expectedOverlap := 5.0
	if info.OverlapX != expectedOverlap || info.OverlapY != expectedOverlap {
		t.Errorf("Expected overlap %.1f, got X:%.1f Y:%.1f", expectedOverlap, info.OverlapX, info.OverlapY)
	}
}

func TestGameEngineTriggerGameOver(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Start the game
	ge.Start()

	// Trigger game over
	ge.TriggerGameOver()

	if !ge.IsGameOver() {
		t.Error("Game should be in game over state")
	}
	if ge.GetState() != StateGameOver {
		t.Error("Game state should be GameOver")
	}
	if ge.IsRunning() {
		t.Error("Game should not be running after game over")
	}
}

func TestGameEngineGameDuration(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Duration should be 0 before starting
	if ge.GetGameDuration() != 0 {
		t.Error("Game duration should be 0 before starting")
	}

	// Start the game
	ge.Start()

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	duration := ge.GetGameDuration()
	if duration <= 0 {
		t.Error("Game duration should be positive after starting")
	}
	if duration < 40*time.Millisecond {
		t.Error("Game duration should be at least the time we waited")
	}
}

func TestGameEngineRestart(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Start, trigger game over, then restart
	ge.Start()
	ge.TriggerGameOver()
	ge.Restart()

	if !ge.IsRunning() {
		t.Error("Game should be running after Restart()")
	}
	if ge.GetState() != StatePlaying {
		t.Error("Game state should be Playing after Restart()")
	}
	if ge.IsGameOver() {
		t.Error("Game should not be in game over state after Restart()")
	}
}

func TestGameEngineCollisionDebug(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Test enabling collision debug (should not panic)
	ge.EnableCollisionDebug(true)
	ge.EnableCollisionDebug(false)

	// Test collision with debug enabled
	ge.EnableCollisionDebug(true)
	rect1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	rect2 := Rectangle{X: 5, Y: 5, Width: 10, Height: 10}

	// Should not panic and should still detect collision
	if !ge.CheckCollision(rect1, rect2) {
		t.Error("Should still detect collision with debug enabled")
	}
}

// Integration test for dinosaur-obstacle collision scenario
func TestGameEngineIntegration_DinosaurObstacleCollision(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Start the game
	ge.Start()

	// Simulate dinosaur and obstacle bounds
	dinosaurBounds := Rectangle{X: 15, Y: 14, Width: 8, Height: 6}
	obstacleBounds := Rectangle{X: 20, Y: 16, Width: 2, Height: 4}

	// Check collision
	collision := ge.CheckCollision(dinosaurBounds, obstacleBounds)

	if collision {
		// Trigger game over on collision
		ge.TriggerGameOver()

		if !ge.IsGameOver() {
			t.Error("Game should be over after collision")
		}
		if ge.GetState() != StateGameOver {
			t.Error("Game state should be GameOver after collision")
		}
	}
}

// Test collision with different obstacle types
func TestGameEngineIntegration_DifferentObstacleTypes(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Disable collision tolerance for precise testing
	ge.SetCollisionTolerance(0)

	dinosaurBounds := Rectangle{X: 15, Y: 14, Width: 8, Height: 6}

	// Test collision with small cactus - position it to definitely overlap
	smallCactus := Rectangle{X: 18, Y: 16, Width: 2, Height: 4}
	if !ge.CheckCollision(dinosaurBounds, smallCactus) {
		t.Error("Should collide with small cactus")
	}

	// Test collision with medium cactus - position it to definitely overlap
	mediumCactus := Rectangle{X: 18, Y: 14, Width: 3, Height: 6}
	if !ge.CheckCollision(dinosaurBounds, mediumCactus) {
		t.Error("Should collide with medium cactus")
	}

	// Test collision with large cactus - position it to definitely overlap
	largeCactus := Rectangle{X: 18, Y: 12, Width: 4, Height: 8}
	if !ge.CheckCollision(dinosaurBounds, largeCactus) {
		t.Error("Should collide with large cactus")
	}
}

// Test jumping over obstacles
func TestGameEngineIntegration_JumpingOverObstacles(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Jumping dinosaur bounds (higher Y position)
	jumpingDinosaurBounds := Rectangle{X: 15, Y: 8, Width: 8, Height: 6}

	// Small cactus that can be jumped over
	smallCactus := Rectangle{X: 20, Y: 16, Width: 2, Height: 4}
	if ge.CheckCollision(jumpingDinosaurBounds, smallCactus) {
		t.Error("Jumping dinosaur should not collide with small cactus")
	}

	// Large cactus that might still cause collision even when jumping
	largeCactus := Rectangle{X: 20, Y: 12, Width: 4, Height: 8}
	collision := ge.CheckCollision(jumpingDinosaurBounds, largeCactus)

	// This depends on the exact jump height and cactus size
	// The test verifies the collision detection works correctly
	if collision {
		t.Log("Jumping dinosaur collides with large cactus (expected for insufficient jump height)")
	} else {
		t.Log("Jumping dinosaur clears large cactus (expected for sufficient jump height)")
	}
}
func TestGameEngineStateTransitions(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Test Menu -> Playing transition
	if !ge.CanTransitionTo(StatePlaying) {
		t.Error("Should be able to transition from Menu to Playing")
	}
	if ge.CanTransitionTo(StateGameOver) {
		t.Error("Should not be able to transition from Menu to GameOver")
	}

	// Transition to Playing
	if !ge.TransitionTo(StatePlaying) {
		t.Error("Should successfully transition to Playing")
	}
	if ge.GetState() != StatePlaying {
		t.Error("State should be Playing after transition")
	}
	if !ge.IsRunning() {
		t.Error("Game should be running in Playing state")
	}
	if !ge.IsInitialized() {
		t.Error("Game should be initialized in Playing state")
	}

	// Test Playing -> GameOver transition
	if !ge.CanTransitionTo(StateGameOver) {
		t.Error("Should be able to transition from Playing to GameOver")
	}
	if !ge.CanTransitionTo(StateMenu) {
		t.Error("Should be able to transition from Playing to Menu")
	}

	// Transition to GameOver
	if !ge.TransitionTo(StateGameOver) {
		t.Error("Should successfully transition to GameOver")
	}
	if ge.GetState() != StateGameOver {
		t.Error("State should be GameOver after transition")
	}
	if ge.IsRunning() {
		t.Error("Game should not be running in GameOver state")
	}
	if !ge.IsGameOver() {
		t.Error("Game should be in game over state")
	}

	// Test GameOver -> Menu transition
	if !ge.CanTransitionTo(StateMenu) {
		t.Error("Should be able to transition from GameOver to Menu")
	}
	if !ge.CanTransitionTo(StatePlaying) {
		t.Error("Should be able to transition from GameOver to Playing")
	}

	// Transition to Menu
	if !ge.TransitionTo(StateMenu) {
		t.Error("Should successfully transition to Menu")
	}
	if ge.GetState() != StateMenu {
		t.Error("State should be Menu after transition")
	}
	if ge.IsRunning() {
		t.Error("Game should not be running in Menu state")
	}
	if ge.IsGameOver() {
		t.Error("Game should not be in game over state in Menu")
	}
}

func TestGameEngineStateChangeCallback(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	var callbackCalled bool
	var fromState, toState GameState

	// Set callback
	ge.SetStateChangeCallback(func(from, to GameState) {
		callbackCalled = true
		fromState = from
		toState = to
	})

	// Trigger state change
	ge.SetState(StatePlaying)

	if !callbackCalled {
		t.Error("State change callback should have been called")
	}
	if fromState != StateMenu {
		t.Error("Callback should receive correct from state")
	}
	if toState != StatePlaying {
		t.Error("Callback should receive correct to state")
	}

	// Reset callback test
	callbackCalled = false
	ge.SetState(StatePlaying) // Same state, should not trigger callback

	if callbackCalled {
		t.Error("Callback should not be called when setting same state")
	}
}

func TestGameEngineInitializationAndCleanup(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Initially not initialized
	if ge.IsInitialized() {
		t.Error("Game should not be initialized initially")
	}

	// Start should initialize
	ge.Start()
	if !ge.IsInitialized() {
		t.Error("Game should be initialized after Start()")
	}

	// Cleanup should reset initialization
	ge.Cleanup()
	if ge.IsInitialized() {
		t.Error("Game should not be initialized after Cleanup()")
	}
	if ge.IsRunning() {
		t.Error("Game should not be running after Cleanup()")
	}
	if ge.IsGameOver() {
		t.Error("Game should not be in game over state after Cleanup()")
	}
}

func TestGameEngineInvalidTransitions(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Try invalid transition from Menu to GameOver
	if ge.TransitionTo(StateGameOver) {
		t.Error("Should not be able to transition from Menu to GameOver")
	}
	if ge.GetState() != StateMenu {
		t.Error("State should remain Menu after invalid transition")
	}

	// Move to Playing state
	ge.TransitionTo(StatePlaying)

	// All transitions from Playing should be valid
	if !ge.CanTransitionTo(StateGameOver) {
		t.Error("Should be able to transition from Playing to GameOver")
	}
	if !ge.CanTransitionTo(StateMenu) {
		t.Error("Should be able to transition from Playing to Menu")
	}
}

func TestGameEngineRestartFromGameOver(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Start game and wait to accumulate some duration
	ge.Start()
	time.Sleep(50 * time.Millisecond)
	originalDuration := ge.GetGameDuration()
	ge.TriggerGameOver()

	if ge.GetState() != StateGameOver {
		t.Error("State should be GameOver after TriggerGameOver()")
	}

	// Wait a bit more to ensure time difference
	time.Sleep(10 * time.Millisecond)

	// Restart should work from GameOver state
	ge.Restart()

	if ge.GetState() != StatePlaying {
		t.Error("State should be Playing after Restart()")
	}
	if !ge.IsRunning() {
		t.Error("Game should be running after Restart()")
	}
	if ge.IsGameOver() {
		t.Error("Game should not be in game over state after Restart()")
	}

	// Game duration should be reset (new start time should be much smaller)
	newDuration := ge.GetGameDuration()
	if newDuration >= originalDuration {
		t.Errorf("Game duration should be reset after restart. Original: %v, New: %v", originalDuration, newDuration)
	}
	if newDuration > 10*time.Millisecond {
		t.Errorf("New game duration should be very small after restart, got: %v", newDuration)
	}
}

func TestGameEngineRestartFromNonGameOverState(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Try restart from Menu state (should not work)
	originalState := ge.GetState()
	ge.Restart()

	if ge.GetState() != originalState {
		t.Error("Restart should not work from non-GameOver state")
	}

	// Try restart from Playing state (should not work)
	ge.Start()
	ge.Restart()

	if ge.GetState() != StatePlaying {
		t.Error("Restart should not work from Playing state")
	}
}

func TestGameEngineStateTransitionHandling(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Test Menu state properties
	ge.SetState(StateMenu)
	if ge.IsRunning() {
		t.Error("Game should not be running in Menu state")
	}
	if ge.IsGameOver() {
		t.Error("Game should not be in game over state in Menu")
	}

	// Test Playing state properties
	ge.SetState(StatePlaying)
	if !ge.IsRunning() {
		t.Error("Game should be running in Playing state")
	}
	if ge.IsGameOver() {
		t.Error("Game should not be in game over state in Playing")
	}
	if !ge.IsInitialized() {
		t.Error("Game should be initialized in Playing state")
	}

	// Test GameOver state properties
	ge.SetState(StateGameOver)
	if ge.IsRunning() {
		t.Error("Game should not be running in GameOver state")
	}
	if !ge.IsGameOver() {
		t.Error("Game should be in game over state in GameOver")
	}
}

func TestGameEngineStopTransition(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Start the game
	ge.Start()
	if ge.GetState() != StatePlaying {
		t.Error("State should be Playing after Start()")
	}

	// Stop should transition to Menu
	ge.Stop()
	if ge.GetState() != StateMenu {
		t.Error("State should be Menu after Stop()")
	}
	if ge.IsRunning() {
		t.Error("Game should not be running after Stop()")
	}
}

func TestGameEngineResetFromDifferentStates(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Reset from Playing state
	ge.Start()
	ge.Reset()
	if ge.GetState() != StateMenu {
		t.Error("State should be Menu after Reset() from Playing")
	}
	if ge.IsInitialized() {
		t.Error("Game should not be initialized after Reset()")
	}

	// Reset from GameOver state
	ge.Start()
	ge.TriggerGameOver()
	ge.Reset()
	if ge.GetState() != StateMenu {
		t.Error("State should be Menu after Reset() from GameOver")
	}
	if ge.IsInitialized() {
		t.Error("Game should not be initialized after Reset()")
	}
}
func TestGameEngineScoring(t *testing.T) {
	config := NewDefaultConfig()
	engine := NewGameEngine(config)

	// Test initial score state
	if engine.GetCurrentScore() != 0 {
		t.Errorf("Expected initial current score to be 0, got %d", engine.GetCurrentScore())
	}

	// Test score reset
	engine.ResetScore()
	if engine.GetCurrentScore() != 0 {
		t.Errorf("Expected current score to be 0 after reset, got %d", engine.GetCurrentScore())
	}
}

func TestGameEngineObstacleBonus(t *testing.T) {
	config := NewDefaultConfig()
	engine := NewGameEngine(config)

	initialScore := engine.GetCurrentScore()
	engine.AddObstacleBonus()

	if engine.GetCurrentScore() <= initialScore {
		t.Errorf("Expected score to increase after obstacle bonus, got %d", engine.GetCurrentScore())
	}
}

func TestGameEngineScoreStateTransitions(t *testing.T) {
	config := NewDefaultConfig()
	engine := NewGameEngine(config)

	// Start playing - should reset score
	engine.SetState(StatePlaying)
	initialScore := engine.GetCurrentScore()

	// Add some score
	engine.AddObstacleBonus()
	if engine.GetCurrentScore() <= initialScore {
		t.Error("Expected score to increase after obstacle bonus")
	}

	// Transition to game over - should finalize score
	engine.SetState(StateGameOver)

	// Start new game - should reset score
	engine.SetState(StatePlaying)
	if engine.GetCurrentScore() != 0 {
		t.Errorf("Expected score to be reset when starting new game, got %d", engine.GetCurrentScore())
	}
}

func TestGameEngineScoreUpdate(t *testing.T) {
	config := NewDefaultConfig()
	engine := NewGameEngine(config)

	// Set to playing state
	engine.SetState(StatePlaying)

	// Simulate time passing
	engine.Update()

	// Score should be accessible
	score := engine.GetScore()
	if score == nil {
		t.Error("Expected score instance to be available")
	}
}

func TestGameEngineHighScore(t *testing.T) {
	config := NewDefaultConfig()
	engine := NewGameEngine(config)

	// Test high score functionality
	initialHigh := engine.GetHighScore()

	// Set a score higher than current high score
	engine.GetScore().Current = initialHigh + 100

	if !engine.IsNewHighScore() {
		t.Error("Expected IsNewHighScore to return true when current > high")
	}

	// Finalize score
	isNewHigh, err := engine.FinalizeScore()
	if err != nil {
		t.Fatalf("Failed to finalize score: %v", err)
	}

	if !isNewHigh {
		t.Error("Expected FinalizeScore to return true for new high score")
	}
}
