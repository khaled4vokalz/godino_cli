package render

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

// Renderer handles all terminal output and screen management
type Renderer struct {
	width      int
	height     int
	buffer     [][]rune
	termWidth  int
	termHeight int
	oldState   *term.State
	isRawMode  bool
}

// NewRenderer creates a new renderer instance
func NewRenderer() (*Renderer, error) {
	r := &Renderer{}

	// Get terminal size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to get terminal size: %w", err)
	}

	r.termWidth = width
	r.termHeight = height
	r.width = width
	r.height = height

	// Initialize buffer
	r.initBuffer()

	return r, nil
}

// initBuffer initializes the screen buffer
func (r *Renderer) initBuffer() {
	r.buffer = make([][]rune, r.height)
	for i := range r.buffer {
		r.buffer[i] = make([]rune, r.width)
		for j := range r.buffer[i] {
			r.buffer[i][j] = ' '
		}
	}
}

// SetRawMode enables raw terminal mode for immediate input
func (r *Renderer) SetRawMode() error {
	if r.isRawMode {
		return nil
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}

	r.oldState = oldState
	r.isRawMode = true

	// Set up signal handler for graceful cleanup
	r.setupSignalHandler()

	return nil
}

// RestoreTerminal restores the terminal to its original state
func (r *Renderer) RestoreTerminal() error {
	if !r.isRawMode || r.oldState == nil {
		return nil
	}

	// Show cursor and clear screen
	fmt.Print("\033[?25h\033[2J\033[H")

	err := term.Restore(int(os.Stdin.Fd()), r.oldState)
	if err != nil {
		return fmt.Errorf("failed to restore terminal: %w", err)
	}

	r.isRawMode = false
	return nil
}

// setupSignalHandler sets up signal handling for graceful cleanup
func (r *Renderer) setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		r.RestoreTerminal()
		os.Exit(0)
	}()
}

// Clear clears the screen buffer
func (r *Renderer) Clear() {
	for i := range r.buffer {
		for j := range r.buffer[i] {
			r.buffer[i][j] = ' '
		}
	}
}

// SetCursor positions the cursor at the specified coordinates
func (r *Renderer) SetCursor(x, y int) {
	fmt.Printf("\033[%d;%dH", y+1, x+1)
}

// HideCursor hides the terminal cursor
func (r *Renderer) HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor shows the terminal cursor
func (r *Renderer) ShowCursor() {
	fmt.Print("\033[?25h")
}

// ClearScreen clears the entire terminal screen
func (r *Renderer) ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// DrawAt draws a character at the specified position in the buffer
func (r *Renderer) DrawAt(x, y int, char rune) {
	if x >= 0 && x < r.width && y >= 0 && y < r.height {
		r.buffer[y][x] = char
	}
}

// DrawString draws a string at the specified position in the buffer
func (r *Renderer) DrawString(x, y int, text string) {
	runes := []rune(text)
	for i, char := range runes {
		if x+i >= r.width {
			break
		}
		r.DrawAt(x+i, y, char)
	}
}

// DrawBox draws a rectangular box in the buffer
func (r *Renderer) DrawBox(x, y, width, height int, char rune) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			r.DrawAt(x+dx, y+dy, char)
		}
	}
}

// Flush outputs the buffer to the terminal
func (r *Renderer) Flush() {
	r.SetCursor(0, 0)

	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			fmt.Printf("%c", r.buffer[y][x])
		}
		if y < r.height-1 {
			fmt.Print("\n")
		}
	}
}

// GetSize returns the current terminal size
func (r *Renderer) GetSize() (int, int) {
	return r.width, r.height
}

// UpdateSize updates the renderer size (useful for handling terminal resize)
func (r *Renderer) UpdateSize() error {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return fmt.Errorf("failed to get terminal size: %w", err)
	}

	r.termWidth = width
	r.termHeight = height
	r.width = width
	r.height = height

	// Reinitialize buffer with new size
	r.initBuffer()

	return nil
}

// DrawScore renders the current score and high score in the top-right corner
func (r *Renderer) DrawScore(currentScore, highScore int) {
	scoreText := fmt.Sprintf("Score: %d", currentScore)
	highScoreText := fmt.Sprintf("High: %d", highScore)

	// Position score in top-right corner
	scoreX := r.width - len(scoreText) - 1
	highScoreX := r.width - len(highScoreText) - 1

	if scoreX >= 0 {
		r.DrawString(scoreX, 0, scoreText)
	}
	if highScoreX >= 0 {
		r.DrawString(highScoreX, 1, highScoreText)
	}
}

// DrawGameOverScreen renders the game over screen with final score
func (r *Renderer) DrawGameOverScreen(finalScore, highScore int, isNewHighScore bool) {
	// Clear the screen first
	r.Clear()

	// Calculate center positions
	centerX := r.width / 2
	centerY := r.height / 2

	// Game Over title
	gameOverText := "GAME OVER"
	titleX := centerX - len(gameOverText)/2
	if titleX >= 0 && titleX+len(gameOverText) < r.width {
		r.DrawString(titleX, centerY-3, gameOverText)
	}

	// Final score
	finalScoreText := fmt.Sprintf("Final Score: %d", finalScore)
	scoreX := centerX - len(finalScoreText)/2
	if scoreX >= 0 && scoreX+len(finalScoreText) < r.width {
		r.DrawString(scoreX, centerY-1, finalScoreText)
	}

	// High score or new high score message
	var highScoreText string
	if isNewHighScore {
		highScoreText = "NEW HIGH SCORE!"
	} else {
		highScoreText = fmt.Sprintf("High Score: %d", highScore)
	}
	highScoreX := centerX - len(highScoreText)/2
	if highScoreX >= 0 && highScoreX+len(highScoreText) < r.width {
		r.DrawString(highScoreX, centerY, highScoreText)
	}

	// Restart instruction
	restartText := "Press 'R' to restart or 'Q' to quit"
	restartX := centerX - len(restartText)/2
	if restartX >= 0 && restartX+len(restartText) < r.width {
		r.DrawString(restartX, centerY+2, restartText)
	}
}

// DrawStartScreen renders the start/menu screen with instructions
func (r *Renderer) DrawStartScreen() {
	// Clear the screen first
	r.Clear()

	// Calculate center positions
	centerX := r.width / 2
	centerY := r.height / 2

	// Game title
	titleLines := []string{
		"  ████████  ██  ██    ██  ████████",
		"  ██     ██ ██  ███   ██  ██     ██",
		"  ██     ██ ██  ████  ██  ██     ██",
		"  ██     ██ ██  ██ ██ ██  ██     ██",
		"  ████████  ██  ██  ████  ████████",
	}

	// Draw title (if it fits)
	startY := centerY - len(titleLines) - 3
	if startY >= 0 {
		for i, line := range titleLines {
			titleX := centerX - len(line)/2
			if titleX >= 0 && titleX+len(line) < r.width && startY+i < r.height {
				r.DrawString(titleX, startY+i, line)
			}
		}
	}

	// Game subtitle
	subtitleText := "CLI Dinosaur Game"
	subtitleX := centerX - len(subtitleText)/2
	subtitleY := startY + len(titleLines) + 1
	if subtitleX >= 0 && subtitleY < r.height {
		r.DrawString(subtitleX, subtitleY, subtitleText)
	}

	// Instructions
	instructions := []string{
		"Press SPACE or UP ARROW to jump",
		"Avoid the cacti!",
		"",
		"Press SPACE to start",
		"Press 'Q' to quit",
	}

	instructionStartY := centerY + 2
	for i, instruction := range instructions {
		if instruction == "" {
			continue // Skip empty lines
		}
		instrX := centerX - len(instruction)/2
		instrY := instructionStartY + i
		if instrX >= 0 && instrY < r.height && instrX+len(instruction) < r.width {
			r.DrawString(instrX, instrY, instruction)
		}
	}
}

// DrawControlInstructions renders control instructions during gameplay
func (r *Renderer) DrawControlInstructions() {
	// Draw controls in bottom-left corner
	controlText := "SPACE/UP: Jump | Q: Quit"
	if len(controlText) < r.width {
		r.DrawString(1, r.height-1, controlText)
	}
}

// DrawCenteredText draws text centered horizontally at the specified y position
func (r *Renderer) DrawCenteredText(y int, text string) {
	if y < 0 || y >= r.height {
		return
	}

	x := (r.width - len(text)) / 2
	if x >= 0 && x+len(text) <= r.width {
		r.DrawString(x, y, text)
	}
}

// DrawBorder draws a border around the screen
func (r *Renderer) DrawBorder() {
	// Top and bottom borders
	for x := 0; x < r.width; x++ {
		r.DrawAt(x, 0, '─')
		r.DrawAt(x, r.height-1, '─')
	}

	// Left and right borders
	for y := 0; y < r.height; y++ {
		r.DrawAt(0, y, '│')
		r.DrawAt(r.width-1, y, '│')
	}

	// Corners
	r.DrawAt(0, 0, '┌')
	r.DrawAt(r.width-1, 0, '┐')
	r.DrawAt(0, r.height-1, '└')
	r.DrawAt(r.width-1, r.height-1, '┘')
}
