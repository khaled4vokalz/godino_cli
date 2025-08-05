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