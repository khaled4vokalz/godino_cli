package render

import (
	"os"
	"testing"
)

func TestNewRenderer(t *testing.T) {
	// Skip if not in a terminal environment
	if !isTerminal() {
		t.Skip("Skipping terminal test in non-terminal environment")
	}

	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}

	if renderer == nil {
		t.Fatal("NewRenderer() returned nil renderer")
	}

	if renderer.width <= 0 || renderer.height <= 0 {
		t.Errorf("Invalid terminal size: width=%d, height=%d", renderer.width, renderer.height)
	}

	if renderer.buffer == nil {
		t.Error("Buffer not initialized")
	}

	if len(renderer.buffer) != renderer.height {
		t.Errorf("Buffer height mismatch: expected %d, got %d", renderer.height, len(renderer.buffer))
	}

	if len(renderer.buffer[0]) != renderer.width {
		t.Errorf("Buffer width mismatch: expected %d, got %d", renderer.width, len(renderer.buffer[0]))
	}
}

func TestRendererInitBuffer(t *testing.T) {
	renderer := &Renderer{
		width:  10,
		height: 5,
	}

	renderer.initBuffer()

	if len(renderer.buffer) != renderer.height {
		t.Errorf("Buffer height mismatch: expected %d, got %d", renderer.height, len(renderer.buffer))
	}

	for i, row := range renderer.buffer {
		if len(row) != renderer.width {
			t.Errorf("Buffer row %d width mismatch: expected %d, got %d", i, renderer.width, len(row))
		}

		for j, char := range row {
			if char != ' ' {
				t.Errorf("Buffer[%d][%d] not initialized to space: got %c", i, j, char)
			}
		}
	}
}

func TestRendererClear(t *testing.T) {
	renderer := &Renderer{
		width:  5,
		height: 3,
	}
	renderer.initBuffer()

	// Fill buffer with test data
	for i := range renderer.buffer {
		for j := range renderer.buffer[i] {
			renderer.buffer[i][j] = 'X'
		}
	}

	// Clear the buffer
	renderer.Clear()

	// Verify all positions are cleared
	for i, row := range renderer.buffer {
		for j, char := range row {
			if char != ' ' {
				t.Errorf("Buffer[%d][%d] not cleared: expected ' ', got %c", i, j, char)
			}
		}
	}
}

func TestRendererDrawAt(t *testing.T) {
	renderer := &Renderer{
		width:  10,
		height: 5,
	}
	renderer.initBuffer()

	// Test valid coordinates
	renderer.DrawAt(2, 1, 'A')
	if renderer.buffer[1][2] != 'A' {
		t.Errorf("DrawAt(2, 1, 'A') failed: expected 'A', got %c", renderer.buffer[1][2])
	}

	// Test boundary coordinates
	renderer.DrawAt(0, 0, 'B')
	if renderer.buffer[0][0] != 'B' {
		t.Errorf("DrawAt(0, 0, 'B') failed: expected 'B', got %c", renderer.buffer[0][0])
	}

	renderer.DrawAt(9, 4, 'C')
	if renderer.buffer[4][9] != 'C' {
		t.Errorf("DrawAt(9, 4, 'C') failed: expected 'C', got %c", renderer.buffer[4][9])
	}

	// Test out-of-bounds coordinates (should not panic)
	renderer.DrawAt(-1, 0, 'D')
	renderer.DrawAt(0, -1, 'E')
	renderer.DrawAt(10, 0, 'F')
	renderer.DrawAt(0, 5, 'G')
}

func TestRendererDrawString(t *testing.T) {
	renderer := &Renderer{
		width:  10,
		height: 5,
	}
	renderer.initBuffer()

	// Test normal string drawing
	text := "Hello"
	renderer.DrawString(2, 1, text)

	expected := []rune(text)
	for i, expectedChar := range expected {
		if renderer.buffer[1][2+i] != expectedChar {
			t.Errorf("DrawString failed at position %d: expected %c, got %c",
				i, expectedChar, renderer.buffer[1][2+i])
		}
	}

	// Test string that exceeds width
	longText := "This is a very long string"
	renderer.DrawString(5, 2, longText)

	// Should only draw up to the width boundary
	for i := 5; i < renderer.width; i++ {
		expectedChar := rune(longText[i-5])
		if renderer.buffer[2][i] != expectedChar {
			t.Errorf("DrawString overflow test failed at position %d: expected %c, got %c",
				i, expectedChar, renderer.buffer[2][i])
		}
	}
}

func TestRendererDrawBox(t *testing.T) {
	renderer := &Renderer{
		width:  10,
		height: 5,
	}
	renderer.initBuffer()

	// Draw a 3x2 box at position (2, 1)
	renderer.DrawBox(2, 1, 3, 2, '#')

	// Verify the box is drawn correctly
	for y := 1; y < 3; y++ {
		for x := 2; x < 5; x++ {
			if renderer.buffer[y][x] != '#' {
				t.Errorf("DrawBox failed at position (%d, %d): expected '#', got %c",
					x, y, renderer.buffer[y][x])
			}
		}
	}

	// Verify areas outside the box are not affected
	if renderer.buffer[0][2] != ' ' {
		t.Errorf("DrawBox affected area outside box: expected ' ', got %c", renderer.buffer[0][2])
	}
	if renderer.buffer[1][1] != ' ' {
		t.Errorf("DrawBox affected area outside box: expected ' ', got %c", renderer.buffer[1][1])
	}
}

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

func TestRendererUpdateSize(t *testing.T) {
	// Skip if not in a terminal environment
	if !isTerminal() {
		t.Skip("Skipping terminal test in non-terminal environment")
	}

	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}

	originalWidth, originalHeight := renderer.GetSize()

	err = renderer.UpdateSize()
	if err != nil {
		t.Errorf("UpdateSize() failed: %v", err)
	}

	newWidth, newHeight := renderer.GetSize()

	// Size should remain the same (or be updated if terminal was resized)
	if newWidth <= 0 || newHeight <= 0 {
		t.Errorf("UpdateSize() resulted in invalid size: (%d, %d)", newWidth, newHeight)
	}

	// Buffer should be reinitialized
	if len(renderer.buffer) != newHeight {
		t.Errorf("Buffer height not updated: expected %d, got %d", newHeight, len(renderer.buffer))
	}

	if len(renderer.buffer[0]) != newWidth {
		t.Errorf("Buffer width not updated: expected %d, got %d", newWidth, len(renderer.buffer[0]))
	}

	t.Logf("Original size: (%d, %d), New size: (%d, %d)",
		originalWidth, originalHeight, newWidth, newHeight)
}

func TestRendererRawModeOperations(t *testing.T) {
	// Skip if not in a terminal environment
	if !isTerminal() {
		t.Skip("Skipping terminal test in non-terminal environment")
	}

	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() failed: %v", err)
	}

	// Test setting raw mode
	err = renderer.SetRawMode()
	if err != nil {
		t.Errorf("SetRawMode() failed: %v", err)
	}

	if !renderer.isRawMode {
		t.Error("isRawMode flag not set after SetRawMode()")
	}

	if renderer.oldState == nil {
		t.Error("oldState not saved after SetRawMode()")
	}

	// Test setting raw mode again (should not error)
	err = renderer.SetRawMode()
	if err != nil {
		t.Errorf("Second SetRawMode() call failed: %v", err)
	}

	// Test restoring terminal
	err = renderer.RestoreTerminal()
	if err != nil {
		t.Errorf("RestoreTerminal() failed: %v", err)
	}

	if renderer.isRawMode {
		t.Error("isRawMode flag still set after RestoreTerminal()")
	}

	// Test restoring again (should not error)
	err = renderer.RestoreTerminal()
	if err != nil {
		t.Errorf("Second RestoreTerminal() call failed: %v", err)
	}
}

// isTerminal checks if we're running in a terminal environment
func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func TestDrawScore(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}
	renderer.initBuffer()

	// Test drawing score
	currentScore := 1500
	highScore := 2000
	renderer.DrawScore(currentScore, highScore)

	// Check that score text is drawn in the top-right area
	scoreText := "Score: 1500"
	highScoreText := "High: 2000"

	// Find the score text in the buffer
	found := false
	for y := 0; y < 3; y++ { // Check first few rows
		line := string(renderer.buffer[y])
		if contains(line, scoreText) {
			found = true
			break
		}
	}
	if !found {
		t.Error("Score text not found in renderer buffer")
	}

	// Find the high score text
	found = false
	for y := 0; y < 3; y++ {
		line := string(renderer.buffer[y])
		if contains(line, highScoreText) {
			found = true
			break
		}
	}
	if !found {
		t.Error("High score text not found in renderer buffer")
	}
}

func TestDrawGameOverScreen(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}
	renderer.initBuffer()

	// Test game over screen with new high score
	renderer.DrawGameOverScreen(2500, 2000, true)

	// Convert buffer to string for easier searching
	content := bufferToString(renderer.buffer)

	// Check for expected text
	expectedTexts := []string{
		"GAME OVER",
		"Final Score: 2500",
		"NEW HIGH SCORE!",
		"Press 'R' to restart",
	}

	for _, text := range expectedTexts {
		if !contains(content, text) {
			t.Errorf("Expected text '%s' not found in game over screen", text)
		}
	}

	// Test game over screen without new high score
	renderer.Clear()
	renderer.DrawGameOverScreen(1500, 2000, false)

	content = bufferToString(renderer.buffer)

	if contains(content, "NEW HIGH SCORE!") {
		t.Error("Should not show 'NEW HIGH SCORE!' when it's not a new high score")
	}
	if !contains(content, "High Score: 2000") {
		t.Error("Should show regular high score when it's not a new high score")
	}
}

func TestDrawStartScreen(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}
	renderer.initBuffer()

	renderer.DrawStartScreen()

	// Convert buffer to string for easier searching
	content := bufferToString(renderer.buffer)

	// Check for expected text
	expectedTexts := []string{
		"CLI Dinosaur Game",
		"Press SPACE or UP ARROW to jump",
		"Avoid the cacti!",
		"Press SPACE to start",
		"Press 'Q' to quit",
	}

	for _, text := range expectedTexts {
		if !contains(content, text) {
			t.Errorf("Expected text '%s' not found in start screen", text)
		}
	}
}

func TestDrawControlInstructions(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}
	renderer.initBuffer()

	renderer.DrawControlInstructions()

	// Check that control instructions are drawn in the bottom area
	controlText := "SPACE/UP: Jump | Q: Quit"

	// Check the last few rows for the control text
	found := false
	for y := renderer.height - 3; y < renderer.height; y++ {
		if y >= 0 {
			line := string(renderer.buffer[y])
			if contains(line, controlText) {
				found = true
				break
			}
		}
	}
	if !found {
		t.Error("Control instructions not found in renderer buffer")
	}
}

func TestDrawCenteredText(t *testing.T) {
	renderer := &Renderer{
		width:  80,
		height: 24,
	}
	renderer.initBuffer()

	testText := "Centered Text"
	testY := 10
	renderer.DrawCenteredText(testY, testText)

	// Check that text is centered
	line := string(renderer.buffer[testY])
	if !contains(line, testText) {
		t.Error("Centered text not found in buffer")
	}

	// Find the position of the text
	textPos := findInString(line, testText)
	expectedPos := (renderer.width - len(testText)) / 2

	// Allow some tolerance for positioning
	if textPos < expectedPos-1 || textPos > expectedPos+1 {
		t.Errorf("Text not properly centered. Expected around position %d, found at %d", expectedPos, textPos)
	}

	// Test edge cases
	renderer.Clear()
	renderer.DrawCenteredText(-1, "Invalid Y")  // Should not panic
	renderer.DrawCenteredText(100, "Invalid Y") // Should not panic
}

func TestDrawBorder(t *testing.T) {
	renderer := &Renderer{
		width:  10,
		height: 5,
	}
	renderer.initBuffer()

	renderer.DrawBorder()

	// Check corners
	if renderer.buffer[0][0] != '┌' {
		t.Errorf("Expected '┌' at top-left corner, got '%c'", renderer.buffer[0][0])
	}
	if renderer.buffer[0][renderer.width-1] != '┐' {
		t.Errorf("Expected '┐' at top-right corner, got '%c'", renderer.buffer[0][renderer.width-1])
	}
	if renderer.buffer[renderer.height-1][0] != '└' {
		t.Errorf("Expected '└' at bottom-left corner, got '%c'", renderer.buffer[renderer.height-1][0])
	}
	if renderer.buffer[renderer.height-1][renderer.width-1] != '┘' {
		t.Errorf("Expected '┘' at bottom-right corner, got '%c'", renderer.buffer[renderer.height-1][renderer.width-1])
	}

	// Check top and bottom borders
	for x := 1; x < renderer.width-1; x++ {
		if renderer.buffer[0][x] != '─' {
			t.Errorf("Expected '─' at top border position (%d,0), got '%c'", x, renderer.buffer[0][x])
		}
		if renderer.buffer[renderer.height-1][x] != '─' {
			t.Errorf("Expected '─' at bottom border position (%d,%d), got '%c'", x, renderer.height-1, renderer.buffer[renderer.height-1][x])
		}
	}

	// Check left and right borders
	for y := 1; y < renderer.height-1; y++ {
		if renderer.buffer[y][0] != '│' {
			t.Errorf("Expected '│' at left border position (0,%d), got '%c'", y, renderer.buffer[y][0])
		}
		if renderer.buffer[y][renderer.width-1] != '│' {
			t.Errorf("Expected '│' at right border position (%d,%d), got '%c'", renderer.width-1, y, renderer.buffer[y][renderer.width-1])
		}
	}
}

func TestUIRenderingEdgeCases(t *testing.T) {
	// Test with very small terminal (edge case)
	renderer := &Renderer{
		width:  10,
		height: 5,
	}
	renderer.initBuffer()

	// These should not panic even with small screen
	renderer.DrawScore(999999, 888888) // Very long numbers
	renderer.DrawGameOverScreen(12345, 54321, true)
	renderer.DrawStartScreen()
	renderer.DrawControlInstructions()

	// Test drawing outside bounds
	renderer.DrawCenteredText(-1, "Invalid Y")
	renderer.DrawCenteredText(100, "Invalid Y")
}

// Helper functions for testing
func bufferToString(buffer [][]rune) string {
	var result string
	for _, row := range buffer {
		result += string(row) + "\n"
	}
	return result
}

func contains(s, substr string) bool {
	return findInString(s, substr) != -1
}

func findInString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
