package entities

import (
	"cli-dino-game/src/engine"
	"testing"
	"time"
)

func TestNewDinosaur(t *testing.T) {
	groundLevel := 15.0
	dino := NewDinosaur(groundLevel)

	// Test initial position and state
	if dino.X != 15.0 {
		t.Errorf("Expected initial X position to be 15.0, got %f", dino.X)
	}
	if dino.Y != groundLevel {
		t.Errorf("Expected initial Y position to be %f, got %f", groundLevel, dino.Y)
	}
	if dino.GroundLevel != groundLevel {
		t.Errorf("Expected ground level to be %f, got %f", groundLevel, dino.GroundLevel)
	}
	if dino.VelocityY != 0.0 {
		t.Errorf("Expected initial velocity to be 0.0, got %f", dino.VelocityY)
	}
	if dino.IsJumping {
		t.Error("Expected dinosaur to not be jumping initially")
	}
	if !dino.IsRunning {
		t.Error("Expected dinosaur to be running initially")
	}
	if dino.AnimFrame != 0 {
		t.Errorf("Expected initial animation frame to be 0, got %d", dino.AnimFrame)
	}
	if dino.Width != 8.0 {
		t.Errorf("Expected width to be 8.0, got %f", dino.Width)
	}
	if dino.Height != 6.0 {
		t.Errorf("Expected height to be 6.0, got %f", dino.Height)
	}
}

func TestDinosaurUpdate_StaticPosition(t *testing.T) {
	dino := NewDinosaur(15.0)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}
	initialX := dino.X
	deltaTime := 0.1 // 100ms

	dino.Update(deltaTime, config)

	// Test that dinosaur stays in fixed position (no horizontal movement)
	if dino.X != initialX {
		t.Errorf("Expected X position to remain %f after update, got %f", initialX, dino.X)
	}
}

func TestDinosaurUpdate_PositionUnchanged(t *testing.T) {
	dino := NewDinosaur(15.0)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}
	dino.X = 25.0 // Set any position
	originalX := dino.X

	dino.Update(0.1, config)

	// Position should remain unchanged
	if dino.X != originalX {
		t.Errorf("Expected X position to remain %f, got %f", originalX, dino.X)
	}
}

func TestDinosaurUpdate_AnimationFrames(t *testing.T) {
	dino := NewDinosaur(15.0)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}
	dino.IsRunning = true
	dino.IsJumping = false

	// Set animation update time to past to trigger frame change
	dino.lastAnimUpdate = time.Now().Add(-time.Millisecond * 300)
	initialFrame := dino.AnimFrame

	dino.Update(0.1, config)

	// Animation frame should have changed
	if dino.AnimFrame == initialFrame {
		t.Error("Expected animation frame to change after sufficient time")
	}

	// Frame should be 0 or 1 (alternating)
	if dino.AnimFrame < 0 || dino.AnimFrame > 1 {
		t.Errorf("Expected animation frame to be 0 or 1, got %d", dino.AnimFrame)
	}
}

func TestDinosaurUpdate_NoAnimationWhenJumping(t *testing.T) {
	dino := NewDinosaur(15.0)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}
	dino.IsJumping = true
	dino.lastAnimUpdate = time.Now().Add(-time.Millisecond * 300)
	initialFrame := dino.AnimFrame

	dino.Update(0.1, config)

	// Animation frame should not change when jumping
	if dino.AnimFrame != initialFrame {
		t.Error("Expected animation frame to not change when jumping")
	}
}

func TestDinosaurGetBounds(t *testing.T) {
	dino := NewDinosaur(15.0)
	dino.X = 25.0
	dino.Y = 10.0

	bounds := dino.GetBounds()

	if bounds.X != dino.X {
		t.Errorf("Expected bounds X to be %f, got %f", dino.X, bounds.X)
	}
	if bounds.Y != dino.Y {
		t.Errorf("Expected bounds Y to be %f, got %f", dino.Y, bounds.Y)
	}
	if bounds.Width != dino.Width {
		t.Errorf("Expected bounds width to be %f, got %f", dino.Width, bounds.Width)
	}
	if bounds.Height != dino.Height {
		t.Errorf("Expected bounds height to be %f, got %f", dino.Height, bounds.Height)
	}
}

func TestDinosaurGetASCIIArt_Running(t *testing.T) {
	dino := NewDinosaur(15.0)
	dino.IsJumping = false
	dino.IsRunning = true

	// Test frame 0
	dino.AnimFrame = 0
	art := dino.GetASCIIArt()
	if len(art) == 0 {
		t.Error("Expected ASCII art to have content")
	}
	if len(art) != 6 {
		t.Errorf("Expected ASCII art to have 6 lines, got %d", len(art))
	}

	// Test frame 1
	dino.AnimFrame = 1
	art2 := dino.GetASCIIArt()
	if len(art2) != 6 {
		t.Errorf("Expected ASCII art to have 6 lines, got %d", len(art2))
	}

	// The two frames should be different
	if art[5] == art2[5] { // Last line should be different between frames
		t.Error("Expected different ASCII art between animation frames")
	}
}

func TestDinosaurGetASCIIArt_Jumping(t *testing.T) {
	dino := NewDinosaur(15.0)
	dino.IsJumping = true

	art := dino.GetASCIIArt()
	if len(art) == 0 {
		t.Error("Expected ASCII art to have content")
	}
	if len(art) != 6 {
		t.Errorf("Expected ASCII art to have 6 lines, got %d", len(art))
	}
}

func TestDinosaurGetPosition(t *testing.T) {
	dino := NewDinosaur(15.0)
	dino.X = 25.0
	dino.Y = 10.0

	x, y := dino.GetPosition()
	if x != dino.X {
		t.Errorf("Expected X position to be %f, got %f", dino.X, x)
	}
	if y != dino.Y {
		t.Errorf("Expected Y position to be %f, got %f", dino.Y, y)
	}
}

func TestDinosaurSetPosition(t *testing.T) {
	dino := NewDinosaur(15.0)
	newX, newY := 30.0, 5.0

	dino.SetPosition(newX, newY)

	if dino.X != newX {
		t.Errorf("Expected X position to be %f, got %f", newX, dino.X)
	}
	if dino.Y != newY {
		t.Errorf("Expected Y position to be %f, got %f", newY, dino.Y)
	}
}

func TestDinosaurIsOnGround(t *testing.T) {
	groundLevel := 15.0
	dino := NewDinosaur(groundLevel)

	// Test on ground
	dino.Y = groundLevel
	dino.IsJumping = false
	if !dino.IsOnGround() {
		t.Error("Expected dinosaur to be on ground")
	}

	// Test above ground
	dino.Y = groundLevel - 5.0
	if dino.IsOnGround() {
		t.Error("Expected dinosaur to not be on ground when above ground level")
	}

	// Test jumping state
	dino.Y = groundLevel
	dino.IsJumping = true
	if dino.IsOnGround() {
		t.Error("Expected dinosaur to not be on ground when jumping")
	}
}

func TestDinosaurMultipleUpdates(t *testing.T) {
	dino := NewDinosaur(15.0)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}
	initialX := dino.X
	deltaTime := 0.05 // 50ms per update
	numUpdates := 10

	for i := 0; i < numUpdates; i++ {
		dino.Update(deltaTime, config)
	}

	// Test that position remains static after multiple updates
	if dino.X != initialX {
		t.Errorf("Expected X position to remain %f after %d updates, got %f", initialX, numUpdates, dino.X)
	}
}

// Jump mechanics tests

func TestDinosaurJump_FromGround(t *testing.T) {
	groundLevel := 15.0
	dino := NewDinosaur(groundLevel)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}

	// Ensure dinosaur is on ground initially
	if !dino.IsOnGround() {
		t.Error("Expected dinosaur to be on ground initially")
	}

	dino.Jump(config)

	// After jumping, dinosaur should be in jumping state
	if !dino.IsJumping {
		t.Error("Expected dinosaur to be jumping after Jump() call")
	}
	if dino.VelocityY != -config.JumpVelocity {
		t.Errorf("Expected velocity to be -%f, got %f", config.JumpVelocity, dino.VelocityY)
	}
	if dino.IsRunning {
		t.Error("Expected dinosaur to stop running when jumping")
	}
}

func TestDinosaurJump_WhenAlreadyJumping(t *testing.T) {
	groundLevel := 15.0
	dino := NewDinosaur(groundLevel)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}

	// Set dinosaur to jumping state
	dino.IsJumping = true
	dino.VelocityY = -10.0
	initialVelocity := dino.VelocityY

	dino.Jump(config)

	// Jump should be ignored when already jumping
	if dino.VelocityY != initialVelocity {
		t.Errorf("Expected velocity to remain %f when already jumping, got %f", initialVelocity, dino.VelocityY)
	}
}

func TestDinosaurJump_WhenAboveGround(t *testing.T) {
	groundLevel := 15.0
	dino := NewDinosaur(groundLevel)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}

	// Position dinosaur above ground but not jumping
	dino.Y = groundLevel - 5.0
	dino.IsJumping = false
	initialVelocity := dino.VelocityY

	dino.Jump(config)

	// Jump should be ignored when not on ground
	if dino.VelocityY != initialVelocity {
		t.Errorf("Expected velocity to remain %f when above ground, got %f", initialVelocity, dino.VelocityY)
	}
	if dino.IsJumping {
		t.Error("Expected dinosaur to not be jumping when above ground")
	}
}

func TestDinosaurUpdate_JumpPhysics(t *testing.T) {
	groundLevel := 15.0
	dino := NewDinosaur(groundLevel)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}
	deltaTime := 0.1 // 100ms

	// Start jump
	dino.Jump(config)
	initialY := dino.Y
	initialVelocity := dino.VelocityY

	// Update once
	dino.Update(deltaTime, config)

	// Check that position was updated using initial velocity (before gravity was applied)
	expectedY := initialY + initialVelocity*deltaTime
	if dino.Y != expectedY {
		t.Errorf("Expected Y position to be %f after velocity application, got %f", expectedY, dino.Y)
	}

	// Check that gravity was applied to velocity (for next frame)
	expectedVelocity := initialVelocity + config.Gravity*deltaTime
	if dino.VelocityY != expectedVelocity {
		t.Errorf("Expected velocity to be %f after gravity application, got %f", expectedVelocity, dino.VelocityY)
	}

	// Should still be jumping
	if !dino.IsJumping {
		t.Error("Expected dinosaur to still be jumping")
	}
}

func TestDinosaurUpdate_Landing(t *testing.T) {
	groundLevel := 15.0
	dino := NewDinosaur(groundLevel)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}

	// Set dinosaur to falling state (positive velocity, near ground)
	dino.IsJumping = true
	dino.VelocityY = 10.0      // Falling down
	dino.Y = groundLevel - 0.5 // Just above ground
	dino.IsRunning = false

	dino.Update(0.1, config)

	// Should have landed
	if dino.IsJumping {
		t.Error("Expected dinosaur to have landed")
	}
	if dino.Y != groundLevel {
		t.Errorf("Expected Y position to be %f after landing, got %f", groundLevel, dino.Y)
	}
	if dino.VelocityY != 0.0 {
		t.Errorf("Expected velocity to be 0 after landing, got %f", dino.VelocityY)
	}
	if !dino.IsRunning {
		t.Error("Expected dinosaur to resume running after landing")
	}
}

func TestDinosaurUpdate_JumpArc(t *testing.T) {
	groundLevel := 15.0
	dino := NewDinosaur(groundLevel)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}
	deltaTime := 0.05 // 50ms updates

	// Start jump
	dino.Jump(config)

	// Track positions during jump
	positions := []float64{dino.Y}
	velocities := []float64{dino.VelocityY}

	// Simulate jump for more frames to ensure landing
	for i := 0; i < 20 && dino.IsJumping; i++ {
		dino.Update(deltaTime, config)
		positions = append(positions, dino.Y)
		velocities = append(velocities, dino.VelocityY)
	}

	// Should have gone up first (Y decreases)
	if positions[1] >= positions[0] {
		t.Error("Expected dinosaur to move up initially")
	}

	// Velocity should increase (become more positive) due to gravity
	if velocities[len(velocities)-1] <= velocities[0] {
		t.Error("Expected velocity to increase due to gravity")
	}

	// Should eventually land back on ground
	if dino.IsJumping {
		t.Error("Expected dinosaur to land within simulation time")
	}
	if dino.Y != groundLevel {
		t.Errorf("Expected final Y position to be %f, got %f", groundLevel, dino.Y)
	}
}

func TestDinosaurGetJumpHeight(t *testing.T) {
	groundLevel := 15.0
	dino := NewDinosaur(groundLevel)

	// On ground
	height := dino.GetJumpHeight()
	if height != 0.0 {
		t.Errorf("Expected jump height to be 0 on ground, got %f", height)
	}

	// Above ground
	dino.Y = groundLevel - 5.0
	height = dino.GetJumpHeight()
	if height != 5.0 {
		t.Errorf("Expected jump height to be 5.0, got %f", height)
	}

	// Below ground (shouldn't happen in normal gameplay)
	dino.Y = groundLevel + 2.0
	height = dino.GetJumpHeight()
	if height != 0.0 {
		t.Errorf("Expected jump height to be 0 when below ground, got %f", height)
	}
}

func TestDinosaurJumpStateTransitions(t *testing.T) {
	groundLevel := 15.0
	dino := NewDinosaur(groundLevel)
	config := &engine.Config{Gravity: 50.0, JumpVelocity: 15.0}

	// Initial state: on ground, running, not jumping
	if !dino.IsRunning || dino.IsJumping || !dino.IsOnGround() {
		t.Error("Expected initial state: running, not jumping, on ground")
	}

	// Jump: should transition to jumping, not running
	dino.Jump(config)
	if dino.IsRunning || !dino.IsJumping {
		t.Error("Expected state after jump: not running, jumping")
	}

	// Simulate complete jump cycle
	for i := 0; i < 20 && dino.IsJumping; i++ {
		dino.Update(0.05, config)
	}

	// After landing: should be running, not jumping, on ground
	if !dino.IsRunning || dino.IsJumping || !dino.IsOnGround() {
		t.Error("Expected state after landing: running, not jumping, on ground")
	}
}

func TestDinosaurJumpWithDifferentConfigs(t *testing.T) {
	groundLevel := 15.0

	// Test with high jump velocity
	dino1 := NewDinosaur(groundLevel)
	config1 := &engine.Config{Gravity: 50.0, JumpVelocity: 20.0}
	dino1.Jump(config1)

	// Test with low jump velocity
	dino2 := NewDinosaur(groundLevel)
	config2 := &engine.Config{Gravity: 50.0, JumpVelocity: 10.0}
	dino2.Jump(config2)

	if dino1.VelocityY >= dino2.VelocityY {
		t.Error("Expected higher jump velocity to result in more negative initial velocity")
	}

	// Simulate one update to see effect of different gravities
	dino3 := NewDinosaur(groundLevel)
	config3 := &engine.Config{Gravity: 100.0, JumpVelocity: 15.0}
	dino3.Jump(config3)

	dino4 := NewDinosaur(groundLevel)
	config4 := &engine.Config{Gravity: 25.0, JumpVelocity: 15.0}
	dino4.Jump(config4)

	dino3.Update(0.1, config3)
	dino4.Update(0.1, config4)

	// Higher gravity should result in higher velocity after update
	if dino3.VelocityY <= dino4.VelocityY {
		t.Error("Expected higher gravity to result in higher velocity after update")
	}
}
