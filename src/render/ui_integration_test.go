package render

import (
	"testing"
)

// TestUIIntegration tests the integration of all UI components
func TestUIIntegration(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}

	// Test complete UI workflow - ensure methods don't panic
	t.Run("StartScreen", func(t *testing.T) {
		renderer.Clear()
		renderer.DrawStartScreen()
		// Test passes if no panic occurs
	})

	t.Run("GameplayUI", func(t *testing.T) {
		renderer.Clear()

		// Draw game UI elements - should not panic
		renderer.DrawScore(1500, 2000)
		renderer.DrawControlInstructions()
		// Test passes if no panic occurs
	})

	t.Run("GameOverScreen", func(t *testing.T) {
		renderer.Clear()

		// Test regular game over - should not panic
		renderer.DrawGameOverScreen(1500, 2000, false)

		// Test new high score - should not panic
		renderer.Clear()
		renderer.DrawGameOverScreen(2500, 2000, true)
		// Test passes if no panic occurs
	})

	t.Run("UIStateTransitions", func(t *testing.T) {
		// Test transitioning between different UI states - should not panic

		// Start screen
		renderer.Clear()
		renderer.DrawStartScreen()

		// Gameplay UI
		renderer.Clear()
		renderer.DrawScore(100, 500)
		renderer.DrawControlInstructions()

		// Game over screen
		renderer.Clear()
		renderer.DrawGameOverScreen(100, 500, false)
		// Test passes if no panic occurs
	})
}

// TestUIResponsiveness tests UI behavior with different screen sizes
func TestUIResponsiveness(t *testing.T) {
	testSizes := []struct {
		name   string
		width  int
		height int
	}{
		{"Large", 120, 30},
		{"Medium", 80, 24},
		{"Small", 60, 20},
		{"Minimal", 40, 10},
	}

	for _, size := range testSizes {
		t.Run(size.name, func(t *testing.T) {
			renderer := &Renderer{
				width:  size.width,
				height: size.height,
			}

			// Test that UI elements don't panic with different screen sizes
			renderer.DrawStartScreen()
			renderer.Clear()

			renderer.DrawScore(12345, 67890)
			renderer.DrawControlInstructions()
			renderer.Clear()

			renderer.DrawGameOverScreen(12345, 67890, true)
			renderer.Clear()

			// Test centered text with various sizes
			renderer.DrawCenteredText(size.height/2, "Test Text")

			// Test border drawing
			renderer.DrawBorder()
			// Test passes if no panic occurs
		})
	}
}

// TestUIContentAccuracy tests that UI methods work with various inputs
func TestUIContentAccuracy(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}

	testCases := []struct {
		name         string
		currentScore int
		highScore    int
		isNewHigh    bool
	}{
		{"LowScore", 100, 1000, false},
		{"HighScore", 5000, 1000, true},
		{"EqualScore", 1000, 1000, false},
		{"ZeroScore", 0, 500, false},
		{"LargeNumbers", 999999, 888888, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test that methods don't panic with various score values
			renderer.Clear()
			renderer.DrawScore(tc.currentScore, tc.highScore)

			// Test game over screen with various inputs
			renderer.Clear()
			renderer.DrawGameOverScreen(tc.currentScore, tc.highScore, tc.isNewHigh)
			// Test passes if no panic occurs
		})
	}
}
