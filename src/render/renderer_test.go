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