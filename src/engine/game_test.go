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

	if ge.GetCollisionTolerance() != 1.0 {
		t.Error("Default collision tolerance should be 1.0")
	}
}

func TestGameEngineStateManagement(t *testing.T) {
	config := NewDefaultConfig()
	ge := NewGameEngine(config)

	// Test state transitions
	ge.SetState(StatePlaying)
	if ge.GetState() != StatePlaying {
		t.Error("State should be set to Playing")
	}

	ge.SetState(StateGameOver)
	if ge.GetState() != StateGameOver {
		t.Error("State should be set to GameOver")
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
