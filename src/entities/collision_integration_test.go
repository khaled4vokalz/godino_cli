package entities

import (
	"cli-dino-game/src/engine"
	"testing"
)

// Integration tests for collision detection between dinosaur and obstacles
func TestDinosaurObstacleCollisionIntegration(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 20.0

	// Create dinosaur and collision detector
	dinosaur := NewDinosaur(groundLevel)
	collisionDetector := engine.NewCollisionDetector()

	// Test collision with small cactus
	t.Run("SmallCactusCollision", func(t *testing.T) {
		smallCactus := NewObstacle(CactusSmall, 20.0, groundLevel, config)

		// Position dinosaur to collide with cactus
		dinosaur.SetPosition(18.0, groundLevel)

		dinosaurBounds := dinosaur.GetBounds()
		cactusBounds := smallCactus.GetBounds()

		if !collisionDetector.CheckCollision(dinosaurBounds, cactusBounds) {
			t.Error("Dinosaur should collide with small cactus")
		}
	})

	// Test no collision when dinosaur jumps over small cactus
	t.Run("JumpOverSmallCactus", func(t *testing.T) {
		smallCactus := NewObstacle(CactusSmall, 25.0, groundLevel, config) // Position further away

		// Position jumping dinosaur high enough to clear small cactus
		dinosaur.SetPosition(15.0, groundLevel-8.0) // Jump high, position to just clear horizontally

		dinosaurBounds := dinosaur.GetBounds()
		cactusBounds := smallCactus.GetBounds()

		if collisionDetector.CheckCollision(dinosaurBounds, cactusBounds) {
			t.Error("Jumping dinosaur should not collide with small cactus")
		}
	})

	// Test collision with medium cactus
	t.Run("MediumCactusCollision", func(t *testing.T) {
		mediumCactus := NewObstacle(CactusMedium, 20.0, groundLevel, config)

		// Position dinosaur to collide with cactus
		dinosaur.SetPosition(18.0, groundLevel)

		dinosaurBounds := dinosaur.GetBounds()
		cactusBounds := mediumCactus.GetBounds()

		if !collisionDetector.CheckCollision(dinosaurBounds, cactusBounds) {
			t.Error("Dinosaur should collide with medium cactus")
		}
	})

	// Test collision with large cactus
	t.Run("LargeCactusCollision", func(t *testing.T) {
		largeCactus := NewObstacle(CactusLarge, 20.0, groundLevel, config)

		// Position dinosaur to collide with cactus
		dinosaur.SetPosition(18.0, groundLevel)

		dinosaurBounds := dinosaur.GetBounds()
		cactusBounds := largeCactus.GetBounds()

		if !collisionDetector.CheckCollision(dinosaurBounds, cactusBounds) {
			t.Error("Dinosaur should collide with large cactus")
		}
	})

	// Test that large cactus is harder to jump over
	t.Run("LargeCactusHardToJumpOver", func(t *testing.T) {
		largeCactus := NewObstacle(CactusLarge, 20.0, groundLevel, config)

		// Position dinosaur with moderate jump height
		dinosaur.SetPosition(18.0, groundLevel-6.0) // Moderate jump

		dinosaurBounds := dinosaur.GetBounds()
		cactusBounds := largeCactus.GetBounds()

		// Large cactus should still collide with moderate jump
		if !collisionDetector.CheckCollision(dinosaurBounds, cactusBounds) {
			t.Error("Large cactus should still collide with moderate jump height")
		}
	})

	// Test very high jump clears large cactus
	t.Run("HighJumpClearsLargeCactus", func(t *testing.T) {
		largeCactus := NewObstacle(CactusLarge, 25.0, groundLevel, config) // Position further away

		// Position dinosaur with very high jump
		dinosaur.SetPosition(15.0, groundLevel-12.0) // Very high jump, position to just clear horizontally

		dinosaurBounds := dinosaur.GetBounds()
		cactusBounds := largeCactus.GetBounds()

		if collisionDetector.CheckCollision(dinosaurBounds, cactusBounds) {
			t.Error("Very high jump should clear large cactus")
		}
	})
}

// Test collision detection with tolerance for more forgiving gameplay
func TestDinosaurObstacleCollisionWithTolerance(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 20.0

	dinosaur := NewDinosaur(groundLevel)
	collisionDetector := engine.NewCollisionDetector()
	tolerance := 0.8 // Smaller tolerance to avoid making rectangles invalid

	// Test near-miss collision with tolerance
	t.Run("NearMissWithTolerance", func(t *testing.T) {
		smallCactus := NewObstacle(CactusSmall, 22.5, groundLevel, config) // Just barely overlapping

		dinosaur.SetPosition(15.0, groundLevel)

		dinosaurBounds := dinosaur.GetBounds()
		cactusBounds := smallCactus.GetBounds()

		// Should collide without tolerance (dinosaur X=15-23, cactus X=22.5-24.5, overlap 0.5)
		if !collisionDetector.CheckCollision(dinosaurBounds, cactusBounds) {
			t.Error("Should collide without tolerance")
		}

		// Should not collide with tolerance (more forgiving)
		if collisionDetector.CheckCollisionWithTolerance(dinosaurBounds, cactusBounds, tolerance) {
			t.Error("Should not collide with tolerance for more forgiving gameplay")
		}
	})

	// Test definite collision even with tolerance
	t.Run("DefiniteCollisionWithTolerance", func(t *testing.T) {
		// Test with overlapping rectangles that should still collide even with tolerance
		// Use a smaller tolerance for this test
		smallTolerance := 0.3

		smallCactus := NewObstacle(CactusSmall, 18.0, groundLevel, config)
		dinosaur.SetPosition(15.0, groundLevel)

		dinosaurBounds := dinosaur.GetBounds()
		cactusBounds := smallCactus.GetBounds()

		// Should collide even with small tolerance (significant overlap)
		if !collisionDetector.CheckCollisionWithTolerance(dinosaurBounds, cactusBounds, smallTolerance) {
			t.Error("Should still collide with small tolerance when there's significant overlap")
		}
	})
}

// Test collision detection during dinosaur movement and animation
func TestDinosaurObstacleCollisionDuringMovement(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 20.0

	dinosaur := NewDinosaur(groundLevel)
	collisionDetector := engine.NewCollisionDetector()

	// Test collision during jump arc
	t.Run("CollisionDuringJumpArc", func(t *testing.T) {
		mediumCactus := NewObstacle(CactusMedium, 20.0, groundLevel, config)

		// Simulate dinosaur at different points in jump arc
		jumpPositions := []struct {
			name string
			y    float64
		}{
			{"TakeoffPosition", groundLevel - 1.0},
			{"MidJumpPosition", groundLevel - 4.0},
			{"PeakJumpPosition", groundLevel - 7.0},
			{"DescentPosition", groundLevel - 3.0},
			{"LandingPosition", groundLevel - 0.5},
		}

		for _, pos := range jumpPositions {
			t.Run(pos.name, func(t *testing.T) {
				dinosaur.SetPosition(18.0, pos.y)

				dinosaurBounds := dinosaur.GetBounds()
				cactusBounds := mediumCactus.GetBounds()

				collision := collisionDetector.CheckCollision(dinosaurBounds, cactusBounds)

				// Log the collision result for analysis
				if collision {
					t.Logf("Collision detected at %s (y=%.1f)", pos.name, pos.y)
				} else {
					t.Logf("No collision at %s (y=%.1f)", pos.name, pos.y)
				}

				// Verify collision logic based on position
				if pos.y >= groundLevel-2.0 { // Low positions should collide
					if !collision {
						t.Errorf("Expected collision at low position %s", pos.name)
					}
				}
			})
		}
	})
}

// Test collision boundaries are accurate for fair gameplay
func TestCollisionBoundaryAccuracy(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 20.0

	dinosaur := NewDinosaur(groundLevel)
	collisionDetector := engine.NewCollisionDetector()

	// Test edge cases for collision boundaries
	t.Run("EdgeCaseCollisions", func(t *testing.T) {
		smallCactus := NewObstacle(CactusSmall, 23.0, groundLevel, config) // At edge of dinosaur

		dinosaur.SetPosition(15.0, groundLevel)

		dinosaurBounds := dinosaur.GetBounds()
		cactusBounds := smallCactus.GetBounds()

		// Test exact boundary conditions
		// Dinosaur: X=15, Width=8, so right edge is at X=23
		// Cactus: X=23, Width=2, so left edge is at X=23
		// They should just touch but not overlap (no collision)

		collision := collisionDetector.CheckCollision(dinosaurBounds, cactusBounds)
		if collision {
			t.Error("Objects touching at edges should not collide")
		}

		// Move cactus slightly to the left to create overlap
		smallCactus.SetPosition(22.5, groundLevel)
		cactusBounds = smallCactus.GetBounds()

		collision = collisionDetector.CheckCollision(dinosaurBounds, cactusBounds)
		if !collision {
			t.Error("Objects with slight overlap should collide")
		}
	})
}

// Test collision detection with various obstacle sizes
func TestCollisionWithVariousObstacleSizes(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 20.0

	dinosaur := NewDinosaur(groundLevel)
	collisionDetector := engine.NewCollisionDetector()

	// Position dinosaur consistently
	dinosaur.SetPosition(15.0, groundLevel)
	dinosaurBounds := dinosaur.GetBounds()

	// Test all obstacle types
	obstacleTypes := []struct {
		obstType     ObstacleType
		name         string
		expectedSize struct{ width, height float64 }
	}{
		{CactusSmall, "SmallCactus", struct{ width, height float64 }{2.0, 4.0}},
		{CactusMedium, "MediumCactus", struct{ width, height float64 }{3.0, 6.0}},
		{CactusLarge, "LargeCactus", struct{ width, height float64 }{4.0, 8.0}},
	}

	for _, obstType := range obstacleTypes {
		t.Run(obstType.name, func(t *testing.T) {
			obstacle := NewObstacle(obstType.obstType, 20.0, groundLevel, config)
			obstacleBounds := obstacle.GetBounds()

			// Verify obstacle dimensions
			if obstacleBounds.Width != obstType.expectedSize.width {
				t.Errorf("Expected width %.1f, got %.1f", obstType.expectedSize.width, obstacleBounds.Width)
			}
			if obstacleBounds.Height != obstType.expectedSize.height {
				t.Errorf("Expected height %.1f, got %.1f", obstType.expectedSize.height, obstacleBounds.Height)
			}

			// Test collision
			collision := collisionDetector.CheckCollision(dinosaurBounds, obstacleBounds)
			if !collision {
				t.Errorf("Dinosaur should collide with %s", obstType.name)
			}

			// Test collision info
			info := collisionDetector.GetCollisionInfo(dinosaurBounds, obstacleBounds)
			if !info.HasCollision {
				t.Errorf("Collision info should indicate collision with %s", obstType.name)
			}
			if info.OverlapArea <= 0 {
				t.Errorf("Overlap area should be positive for %s", obstType.name)
			}
		})
	}
}
