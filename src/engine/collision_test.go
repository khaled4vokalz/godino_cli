package engine

import (
	"testing"
)

func TestNewCollisionDetector(t *testing.T) {
	cd := NewCollisionDetector()
	if cd == nil {
		t.Fatal("NewCollisionDetector should not return nil")
	}
	if cd.debugMode {
		t.Error("New collision detector should have debug mode disabled by default")
	}
}

func TestSetDebugMode(t *testing.T) {
	cd := NewCollisionDetector()

	// Test enabling debug mode
	cd.SetDebugMode(true)
	if !cd.debugMode {
		t.Error("Debug mode should be enabled")
	}

	// Test disabling debug mode
	cd.SetDebugMode(false)
	if cd.debugMode {
		t.Error("Debug mode should be disabled")
	}
}

func TestCheckCollision_NoCollision(t *testing.T) {
	cd := NewCollisionDetector()

	// Test rectangles that don't overlap
	rect1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	rect2 := Rectangle{X: 15, Y: 0, Width: 10, Height: 10}

	if cd.CheckCollision(rect1, rect2) {
		t.Error("Rectangles should not collide")
	}
}

func TestCheckCollision_WithCollision(t *testing.T) {
	cd := NewCollisionDetector()

	// Test rectangles that overlap
	rect1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	rect2 := Rectangle{X: 5, Y: 5, Width: 10, Height: 10}

	if !cd.CheckCollision(rect1, rect2) {
		t.Error("Rectangles should collide")
	}
}

func TestCheckCollision_EdgeTouching(t *testing.T) {
	cd := NewCollisionDetector()

	// Test rectangles that touch at edges (should not collide)
	rect1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	rect2 := Rectangle{X: 10, Y: 0, Width: 10, Height: 10}

	if cd.CheckCollision(rect1, rect2) {
		t.Error("Rectangles touching at edge should not collide")
	}
}

func TestCheckCollision_CompleteOverlap(t *testing.T) {
	cd := NewCollisionDetector()

	// Test one rectangle completely inside another
	rect1 := Rectangle{X: 0, Y: 0, Width: 20, Height: 20}
	rect2 := Rectangle{X: 5, Y: 5, Width: 10, Height: 10}

	if !cd.CheckCollision(rect1, rect2) {
		t.Error("Completely overlapping rectangles should collide")
	}
}

func TestCheckCollision_SameRectangle(t *testing.T) {
	cd := NewCollisionDetector()

	// Test identical rectangles
	rect := Rectangle{X: 10, Y: 10, Width: 5, Height: 5}

	if !cd.CheckCollision(rect, rect) {
		t.Error("Identical rectangles should collide")
	}
}

func TestCheckCollisionWithTolerance_NoCollisionWithTolerance(t *testing.T) {
	cd := NewCollisionDetector()

	// Test rectangles that would collide without tolerance but not with tolerance
	rect1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	rect2 := Rectangle{X: 8, Y: 8, Width: 10, Height: 10}
	tolerance := 3.0

	// Should collide without tolerance
	if !cd.CheckCollision(rect1, rect2) {
		t.Error("Rectangles should collide without tolerance")
	}

	// Should not collide with tolerance
	if cd.CheckCollisionWithTolerance(rect1, rect2, tolerance) {
		t.Error("Rectangles should not collide with tolerance")
	}
}

func TestCheckCollisionWithTolerance_CollisionWithTolerance(t *testing.T) {
	cd := NewCollisionDetector()

	// Test rectangles that collide even with tolerance
	rect1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	rect2 := Rectangle{X: 2, Y: 2, Width: 10, Height: 10}
	tolerance := 1.0

	if !cd.CheckCollisionWithTolerance(rect1, rect2, tolerance) {
		t.Error("Rectangles should collide even with tolerance")
	}
}

func TestCheckCollisionWithTolerance_ZeroDimensions(t *testing.T) {
	cd := NewCollisionDetector()

	// Test with tolerance that makes rectangles have zero or negative dimensions
	rect1 := Rectangle{X: 0, Y: 0, Width: 4, Height: 4}
	rect2 := Rectangle{X: 1, Y: 1, Width: 4, Height: 4}
	tolerance := 3.0 // This will make adjusted rectangles have negative dimensions

	if cd.CheckCollisionWithTolerance(rect1, rect2, tolerance) {
		t.Error("Should not collide when tolerance makes rectangles have zero/negative dimensions")
	}
}

func TestGetCollisionInfo_NoCollision(t *testing.T) {
	cd := NewCollisionDetector()

	rect1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	rect2 := Rectangle{X: 15, Y: 0, Width: 10, Height: 10}

	info := cd.GetCollisionInfo(rect1, rect2)

	if info.HasCollision {
		t.Error("Should not have collision")
	}
	if info.OverlapX != 0 || info.OverlapY != 0 || info.OverlapArea != 0 {
		t.Error("Overlap values should be zero for non-colliding rectangles")
	}
}

func TestGetCollisionInfo_WithCollision(t *testing.T) {
	cd := NewCollisionDetector()

	rect1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	rect2 := Rectangle{X: 5, Y: 5, Width: 10, Height: 10}

	info := cd.GetCollisionInfo(rect1, rect2)

	if !info.HasCollision {
		t.Error("Should have collision")
	}

	expectedOverlapX := 5.0 // overlap from x=5 to x=10
	expectedOverlapY := 5.0 // overlap from y=5 to y=10
	expectedArea := 25.0    // 5 * 5

	if info.OverlapX != expectedOverlapX {
		t.Errorf("Expected OverlapX %.1f, got %.1f", expectedOverlapX, info.OverlapX)
	}
	if info.OverlapY != expectedOverlapY {
		t.Errorf("Expected OverlapY %.1f, got %.1f", expectedOverlapY, info.OverlapY)
	}
	if info.OverlapArea != expectedArea {
		t.Errorf("Expected OverlapArea %.1f, got %.1f", expectedArea, info.OverlapArea)
	}
}

func TestGetCollisionInfo_PartialOverlap(t *testing.T) {
	cd := NewCollisionDetector()

	rect1 := Rectangle{X: 0, Y: 0, Width: 15, Height: 10}
	rect2 := Rectangle{X: 10, Y: 5, Width: 10, Height: 10}

	info := cd.GetCollisionInfo(rect1, rect2)

	if !info.HasCollision {
		t.Error("Should have collision")
	}

	expectedOverlapX := 5.0 // overlap from x=10 to x=15
	expectedOverlapY := 5.0 // overlap from y=5 to y=10
	expectedArea := 25.0    // 5 * 5

	if info.OverlapX != expectedOverlapX {
		t.Errorf("Expected OverlapX %.1f, got %.1f", expectedOverlapX, info.OverlapX)
	}
	if info.OverlapY != expectedOverlapY {
		t.Errorf("Expected OverlapY %.1f, got %.1f", expectedOverlapY, info.OverlapY)
	}
	if info.OverlapArea != expectedArea {
		t.Errorf("Expected OverlapArea %.1f, got %.1f", expectedArea, info.OverlapArea)
	}
}

// Test collision scenarios specific to the dino game
func TestDinosaurObstacleCollision_SmallCactus(t *testing.T) {
	cd := NewCollisionDetector()

	// Simulate dinosaur bounds (typical size and position)
	dinosaur := Rectangle{X: 15, Y: 14, Width: 8, Height: 6} // On ground at y=20, height=6 means y=14-20

	// Small cactus bounds
	smallCactus := Rectangle{X: 20, Y: 16, Width: 2, Height: 4} // Ground at y=20, height=4 means y=16-20

	// Should collide when dinosaur runs into cactus
	if !cd.CheckCollision(dinosaur, smallCactus) {
		t.Error("Dinosaur should collide with small cactus")
	}
}

func TestDinosaurObstacleCollision_JumpOverSmallCactus(t *testing.T) {
	cd := NewCollisionDetector()

	// Simulate jumping dinosaur bounds (higher Y position)
	dinosaur := Rectangle{X: 15, Y: 8, Width: 8, Height: 6} // Jumping high

	// Small cactus bounds
	smallCactus := Rectangle{X: 20, Y: 16, Width: 2, Height: 4}

	// Should not collide when dinosaur jumps over small cactus
	if cd.CheckCollision(dinosaur, smallCactus) {
		t.Error("Jumping dinosaur should not collide with small cactus")
	}
}

func TestDinosaurObstacleCollision_LargeCactus(t *testing.T) {
	cd := NewCollisionDetector()

	// Simulate dinosaur bounds
	dinosaur := Rectangle{X: 15, Y: 14, Width: 8, Height: 6}

	// Large cactus bounds (taller, harder to jump over)
	largeCactus := Rectangle{X: 20, Y: 12, Width: 4, Height: 8}

	// Should collide when dinosaur runs into large cactus
	if !cd.CheckCollision(dinosaur, largeCactus) {
		t.Error("Dinosaur should collide with large cactus")
	}
}

func TestDinosaurObstacleCollision_JumpOverLargeCactus(t *testing.T) {
	cd := NewCollisionDetector()

	// Simulate high jumping dinosaur bounds
	dinosaur := Rectangle{X: 15, Y: 6, Width: 8, Height: 6} // Very high jump

	// Large cactus bounds
	largeCactus := Rectangle{X: 20, Y: 12, Width: 4, Height: 8}

	// Should not collide when dinosaur jumps high enough over large cactus
	if cd.CheckCollision(dinosaur, largeCactus) {
		t.Error("High jumping dinosaur should not collide with large cactus")
	}
}

func TestDinosaurObstacleCollision_WithTolerance(t *testing.T) {
	cd := NewCollisionDetector()

	// Test collision with tolerance for more forgiving gameplay
	dinosaur := Rectangle{X: 15, Y: 14, Width: 8, Height: 6}
	smallCactus := Rectangle{X: 22, Y: 16, Width: 2, Height: 4} // Just barely touching

	tolerance := 1.0

	// Should collide without tolerance
	if !cd.CheckCollision(dinosaur, smallCactus) {
		t.Error("Should collide without tolerance")
	}

	// Should not collide with tolerance (more forgiving)
	if cd.CheckCollisionWithTolerance(dinosaur, smallCactus, tolerance) {
		t.Error("Should not collide with tolerance for more forgiving gameplay")
	}
}

func TestMinMaxHelpers(t *testing.T) {
	// Test min function
	if min(5.0, 3.0) != 3.0 {
		t.Error("min(5.0, 3.0) should return 3.0")
	}
	if min(2.0, 7.0) != 2.0 {
		t.Error("min(2.0, 7.0) should return 2.0")
	}
	if min(4.0, 4.0) != 4.0 {
		t.Error("min(4.0, 4.0) should return 4.0")
	}

	// Test max function
	if max(5.0, 3.0) != 5.0 {
		t.Error("max(5.0, 3.0) should return 5.0")
	}
	if max(2.0, 7.0) != 7.0 {
		t.Error("max(2.0, 7.0) should return 7.0")
	}
	if max(4.0, 4.0) != 4.0 {
		t.Error("max(4.0, 4.0) should return 4.0")
	}
}
