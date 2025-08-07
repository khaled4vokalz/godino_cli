package render

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

// Renderer handles all terminal output and screen management using termbox-go
type Renderer struct {
	width  int
	height int
}

// NewRenderer creates a new renderer instance using termbox-go
func NewRenderer() (*Renderer, error) {
	// Initialize termbox
	err := termbox.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize termbox: %w", err)
	}

	// Set input mode for better key handling
	termbox.SetInputMode(termbox.InputEsc)

	// Get terminal size
	width, height := termbox.Size()

	return &Renderer{
		width:  width,
		height: height,
	}, nil
}

// Close closes the termbox and restores terminal
func (r *Renderer) Close() {
	termbox.Close()
}

// Clear clears the screen buffer
func (r *Renderer) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

// DrawAt draws a character at the specified position
func (r *Renderer) DrawAt(x, y int, char rune) {
	if x >= 0 && x < r.width && y >= 0 && y < r.height {
		termbox.SetCell(x, y, char, termbox.ColorDefault, termbox.ColorDefault)
	}
}

// DrawString draws a string at the specified position
func (r *Renderer) DrawString(x, y int, text string) {
	charPos := 0
	for _, char := range text {
		if x+charPos >= r.width {
			break
		}
		r.DrawAt(x+charPos, y, char)
		charPos++
	}
}

// DrawBox draws a rectangular box
func (r *Renderer) DrawBox(x, y, width, height int, char rune) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			r.DrawAt(x+dx, y+dy, char)
		}
	}
}

// Flush renders the buffer to the terminal (termbox handles double buffering)
func (r *Renderer) Flush() {
	termbox.Flush()
}

// GetSize returns the current terminal size
func (r *Renderer) GetSize() (int, int) {
	return r.width, r.height
}

// UpdateSize updates the renderer size (useful for handling terminal resize)
func (r *Renderer) UpdateSize() error {
	width, height := termbox.Size()
	r.width = width
	r.height = height
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

	// Simple game title that fits in most terminals
	titleText := "CLI DINO GAME"
	titleX := centerX - len(titleText)/2
	titleY := centerY - 4
	if titleX >= 0 && titleY >= 0 && titleX+len(titleText) < r.width {
		r.DrawString(titleX, titleY, titleText)
	}

	// Simple dinosaur sprite for the menu (using Unicode for better visuals)
	dinoSprite := []string{
		"  ████",
		"  █  █",
		"  ████",
		"█ ██ █",
	}

	dinoX := centerX - 3
	dinoY := centerY - 1
	if dinoX >= 0 && dinoY >= 0 {
		for i, line := range dinoSprite {
			if dinoY+i < r.height && dinoX+len(line) < r.width {
				r.DrawString(dinoX, dinoY+i, line)
			}
		}
	}

	// Instructions
	instructions := []string{
		"SPACE/UP: Jump | Q: Quit",
		"",
		"Press SPACE to start",
	}

	instructionStartY := centerY + 4
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
