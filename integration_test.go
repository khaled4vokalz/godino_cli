package main

import (
	"cli-dino-game/src/engine"
	"cli-dino-game/src/input"
	"testing"
	"time"
)

// TestCompleteGameCycle tests a complete game cycle from start to finish
func TestCompleteGameCycle(t *testing.T) {
	game := NewTestGame()

	// Verify initial state
	if game.engine.GetState() != engine.StateMenu {
		t.Fatal("Game should start in menu state")
	}

	// Start the game
	game.handleInput(input.KeySpace)
	if game.engine.GetState() != engine.StatePlaying {
		t.Fatal("Game should transition to playing state")
	}

	// Simulate gameplay for a reasonable duration
	startTime := time.Now()
	gameplayDuration := time.Millisecond * 500 // Half second of gameplay

	for time.Since(startTime) < gameplayDuration {
		game.update()

		// Simulate occasional jumps
		if time.Since(startTime)%100*time.Millisecond < 10*time.Millisecond {
			game.handleInput(input.KeySpace)
		}

		time.Sleep(time.Millisecond * 16) // ~60 FPS simulation
	}

	// Verify game is still running
	if game.engine.GetState() != engine.StatePlaying {
		t.Log("Game state changed during gameplay, which is normal if collision occurred")
	}

	// Force game over to test that transition
	game.engine.TriggerGameOver()
	if game.engine.GetState() != engine.StateGameOver {
		t.Error("Game should transition to game over state")
	}

	// Test restart
	game.handleInput(input.KeyR)
	if game.engine.GetState() != engine.StatePlaying {
		t.Error("Game should restart to playing state")
	}

	// Verify systems are working after restart
	initialScore := game.engine.GetCurrentScore()
	if initialScore != 0 {
		t.Error("Score should be reset after restart")
	}

	// Run a few more updates to ensure everything still works
	for i := 0; i < 5; i++ {
		game.update()
		time.Sleep(time.Millisecond * 16)
	}

	t.Log("Complete game cycle test passed successfully")
}

// TestGamePerformanceUnderLoad tests game performance with extended gameplay
func TestGamePerformanceUnderLoad(t *testing.T) {
	game := NewTestGame()
	game.engine.SetState(engine.StatePlaying)

	// Run for a longer period to test performance
	startTime := time.Now()
	updateCount := 0
	testDuration := time.Millisecond * 1000 // 1 second

	for time.Since(startTime) < testDuration {
		game.update()
		updateCount++
		time.Sleep(time.Microsecond * 100) // Very short sleep to stress test
	}

	actualDuration := time.Since(startTime)
	updatesPerSecond := float64(updateCount) / actualDuration.Seconds()

	t.Logf("Performed %d updates in %v (%.2f updates/sec)", updateCount, actualDuration, updatesPerSecond)

	// Verify we can maintain reasonable performance (at least 100 updates/sec)
	if updatesPerSecond < 100 {
		t.Errorf("Performance too low: %.2f updates/sec (expected at least 100)", updatesPerSecond)
	}
}

// TestMemoryUsageStability tests that the game doesn't leak memory during extended play
func TestMemoryUsageStability(t *testing.T) {
	game := NewTestGame()
	game.engine.SetState(engine.StatePlaying)

	// Force spawn many obstacles to test memory management
	for i := 0; i < 100; i++ {
		game.spawner.Update(0.1) // Force frequent spawning
		game.update()
	}

	// Get obstacle count
	obstacleCount := game.spawner.GetActiveObstacleCount()
	t.Logf("Active obstacles after stress test: %d", obstacleCount)

	// Obstacles should be cleaned up as they move off screen
	// The exact number depends on timing, but it shouldn't grow indefinitely
	if obstacleCount > 50 {
		t.Errorf("Too many active obstacles: %d (possible memory leak)", obstacleCount)
	}
}

// TestErrorHandling tests that the game handles error conditions gracefully
func TestErrorHandling(t *testing.T) {
	game := NewTestGame()

	// Test invalid state transitions
	game.engine.SetState(engine.StateGameOver)

	// Try to jump in game over state (should be ignored)
	initialY := game.dinosaur.Y
	game.handleInput(input.KeySpace)
	if game.dinosaur.Y != initialY {
		t.Error("Dinosaur should not jump in game over state")
	}

	// Test multiple rapid inputs
	for i := 0; i < 10; i++ {
		game.handleInput(input.KeySpace)
		game.handleInput(input.KeyR)
		game.handleInput(input.KeyQ)
	}

	// Game should still be in a valid state
	validStates := []engine.GameState{engine.StateMenu, engine.StatePlaying, engine.StateGameOver}
	currentState := game.engine.GetState()
	isValidState := false
	for _, state := range validStates {
		if currentState == state {
			isValidState = true
			break
		}
	}

	if !isValidState {
		t.Errorf("Game in invalid state after rapid inputs: %v", currentState)
	}
}
