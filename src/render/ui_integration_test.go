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
	renderer.initBuffer()

	// Test complete UI workflow
	t.Run("StartScreen", func(t *testing.T) {
		renderer.Clear()
		renderer.DrawStartScreen()

		// Verify the screen is not empty
		hasContent := false
		for _, row := range renderer.buffer {
			for _, char := range row {
				if char != ' ' {
					hasContent = true
					break
				}
			}
			if hasContent {
				break
			}
		}
		if !hasContent {
			t.Error("Start screen should have content")
		}
	})

	t.Run("GameplayUI", func(t *testing.T) {
		renderer.Clear()

		// Draw game UI elements
		renderer.DrawScore(1500, 2000)
		renderer.DrawControlInstructions()

		// Verify score is displayed
		content := bufferToString(renderer.buffer)
		if !contains(content, "Score: 1500") {
			t.Error("Score should be displayed during gameplay")
		}
		if !contains(content, "High: 2000") {
			t.Error("High score should be displayed during gameplay")
		}
		if !contains(content, "SPACE/UP: Jump") {
			t.Error("Control instructions should be displayed during gameplay")
		}
	})

	t.Run("GameOverScreen", func(t *testing.T) {
		renderer.Clear()

		// Test regular game over
		renderer.DrawGameOverScreen(1500, 2000, false)
		content := bufferToString(renderer.buffer)

		if !contains(content, "GAME OVER") {
			t.Error("Game over screen should display 'GAME OVER'")
		}
		if !contains(content, "Final Score: 1500") {
			t.Error("Game over screen should display final score")
		}
		if !contains(content, "High Score: 2000") {
			t.Error("Game over screen should display high score when not new")
		}

		// Test new high score
		renderer.Clear()
		renderer.DrawGameOverScreen(2500, 2000, true)
		content = bufferToString(renderer.buffer)

		if !contains(content, "NEW HIGH SCORE!") {
			t.Error("Game over screen should display new high score message")
		}
	})

	t.Run("UIStateTransitions", func(t *testing.T) {
		// Test transitioning between different UI states

		// Start screen
		renderer.Clear()
		renderer.DrawStartScreen()
		startContent := bufferToString(renderer.buffer)

		// Gameplay UI
		renderer.Clear()
		renderer.DrawScore(100, 500)
		renderer.DrawControlInstructions()
		gameContent := bufferToString(renderer.buffer)

		// Game over screen
		renderer.Clear()
		renderer.DrawGameOverScreen(100, 500, false)
		gameOverContent := bufferToString(renderer.buffer)

		// Verify each state has unique content
		if startContent == gameContent {
			t.Error("Start screen and gameplay UI should be different")
		}
		if gameContent == gameOverContent {
			t.Error("Gameplay UI and game over screen should be different")
		}
		if startContent == gameOverContent {
			t.Error("Start screen and game over screen should be different")
		}
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
			renderer.initBuffer()

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
		})
	}
}

// TestUIContentAccuracy tests that UI displays correct information
func TestUIContentAccuracy(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}
	renderer.initBuffer()

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
			// Test score display accuracy
			renderer.Clear()
			renderer.DrawScore(tc.currentScore, tc.highScore)
			content := bufferToString(renderer.buffer)

			expectedScore := "Score: " + intToString(tc.currentScore)
			expectedHigh := "High: " + intToString(tc.highScore)

			if !contains(content, expectedScore) {
				t.Errorf("Expected to find '%s' in score display", expectedScore)
			}
			if !contains(content, expectedHigh) {
				t.Errorf("Expected to find '%s' in score display", expectedHigh)
			}

			// Test game over screen accuracy
			renderer.Clear()
			renderer.DrawGameOverScreen(tc.currentScore, tc.highScore, tc.isNewHigh)
			content = bufferToString(renderer.buffer)

			expectedFinal := "Final Score: " + intToString(tc.currentScore)
			if !contains(content, expectedFinal) {
				t.Errorf("Expected to find '%s' in game over screen", expectedFinal)
			}

			if tc.isNewHigh {
				if !contains(content, "NEW HIGH SCORE!") {
					t.Error("Expected to find 'NEW HIGH SCORE!' message")
				}
			} else {
				expectedHighScore := "High Score: " + intToString(tc.highScore)
				if !contains(content, expectedHighScore) {
					t.Errorf("Expected to find '%s' in game over screen", expectedHighScore)
				}
			}
		})
	}
}

// Helper function to convert int to string (simple implementation)
func intToString(n int) string {
	if n == 0 {
		return "0"
	}

	var result string
	negative := n < 0
	if negative {
		n = -n
	}

	for n > 0 {
		digit := n % 10
		result = string(rune('0'+digit)) + result
		n /= 10
	}

	if negative {
		result = "-" + result
	}

	return result
}
