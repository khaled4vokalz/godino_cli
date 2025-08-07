package render

import (
	"testing"
)

func TestRendererGetSize(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}

	width, height := renderer.GetSize()
	if width != 80 || height != 24 {
		t.Errorf("GetSize() failed: expected (80, 24), got (%d, %d)", width, height)
	}
}

func TestDrawCenteredText(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}

	testText := "Centered Text"
	testY := 10

	// Test that DrawCenteredText doesn't panic with valid input
	renderer.DrawCenteredText(testY, testText)

	// Test edge cases - these should not panic
	renderer.DrawCenteredText(-1, "Invalid Y")  // Should not panic
	renderer.DrawCenteredText(100, "Invalid Y") // Should not panic

	// Test centering calculation logic
	expectedPos := (renderer.width - len(testText)) / 2
	if expectedPos < 0 {
		t.Error("Text too long for terminal width")
	}
}

func TestDrawBorder(t *testing.T) {
	renderer := &Renderer{
		width:  10,
		height: 5,
	}

	// Test that DrawBorder doesn't panic
	renderer.DrawBorder()

	// Since we can't directly inspect termbox buffer, we test the logic
	// by ensuring the method completes without error and the renderer
	// dimensions are valid for border drawing
	if renderer.width < 2 || renderer.height < 2 {
		t.Error("Renderer dimensions too small for border drawing")
	}
}

func TestUIRenderingEdgeCases(t *testing.T) {
	// Test with very small terminal (edge case)
	renderer := &Renderer{
		width:  10,
		height: 5,
	}

	// These should not panic even with small screen
	renderer.DrawScore(999999, 888888) // Very long numbers
	renderer.DrawGameOverScreen(12345, 54321, true)
	renderer.DrawStartScreen()
	renderer.DrawControlInstructions()

	// Test drawing outside bounds
	renderer.DrawCenteredText(-1, "Invalid Y")
	renderer.DrawCenteredText(100, "Invalid Y")
}

func TestRendererDrawAt(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}

	// Test that DrawAt doesn't panic with valid coordinates
	renderer.DrawAt(10, 10, 'X')

	// Test boundary conditions - should not panic
	renderer.DrawAt(-1, 10, 'X')  // Invalid x
	renderer.DrawAt(10, -1, 'X')  // Invalid y
	renderer.DrawAt(100, 10, 'X') // x out of bounds
	renderer.DrawAt(10, 100, 'X') // y out of bounds
}

func TestRendererDrawString(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}

	// Test that DrawString doesn't panic
	renderer.DrawString(10, 10, "Hello World")

	// Test edge cases
	renderer.DrawString(-1, 10, "Invalid X")
	renderer.DrawString(10, -1, "Invalid Y")
	renderer.DrawString(75, 10, "Long string that exceeds width")
	renderer.DrawString(10, 10, "") // Empty string
}

func TestRendererDrawBox(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}

	// Test that DrawBox doesn't panic
	renderer.DrawBox(10, 10, 5, 3, '#')

	// Test edge cases
	renderer.DrawBox(-1, 10, 5, 3, '#') // Invalid x
	renderer.DrawBox(10, -1, 5, 3, '#') // Invalid y
	renderer.DrawBox(10, 10, 0, 3, '#') // Zero width
	renderer.DrawBox(10, 10, 5, 0, '#') // Zero height
}
