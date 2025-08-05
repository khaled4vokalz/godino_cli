package entities

import (
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
	initialX := dino.X
	deltaTime := 0.1 // 100ms

	dino.Update(deltaTime)

	// Test that dinosaur stays in fixed position (no horizontal movement)
	if dino.X != initialX {
		t.Errorf("Expected X position to remain %f after update, got %f", initialX, dino.X)
	}
}

func TestDinosaurUpdate_PositionUnchanged(t *testing.T) {
	dino := NewDinosaur(15.0)
	dino.X = 25.0 // Set any position
	originalX := dino.X

	dino.Update(0.1)

	// Position should remain unchanged
	if dino.X != originalX {
		t.Errorf("Expected X position to remain %f, got %f", originalX, dino.X)
	}
}

func TestDinosaurUpdate_AnimationFrames(t *testing.T) {
	dino := NewDinosaur(15.0)
	dino.IsRunning = true
	dino.IsJumping = false

	// Set animation update time to past to trigger frame change
	dino.lastAnimUpdate = time.Now().Add(-time.Millisecond * 300)
	initialFrame := dino.AnimFrame

	dino.Update(0.1)

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
	dino.IsJumping = true
	dino.lastAnimUpdate = time.Now().Add(-time.Millisecond * 300)
	initialFrame := dino.AnimFrame

	dino.Update(0.1)

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
	initialX := dino.X
	deltaTime := 0.05 // 50ms per update
	numUpdates := 10

	for i := 0; i < numUpdates; i++ {
		dino.Update(deltaTime)
	}

	// Test that position remains static after multiple updates
	if dino.X != initialX {
		t.Errorf("Expected X position to remain %f after %d updates, got %f", initialX, numUpdates, dino.X)
	}
}
