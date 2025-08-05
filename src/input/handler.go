package input

import (
	"os"
	"time"

	"golang.org/x/term"
)

// InputHandler manages keyboard input in a non-blocking manner
type InputHandler struct {
	inputChan chan InputEvent
	done      chan bool
	oldState  *term.State
}

// NewInputHandler creates a new InputHandler instance
func NewInputHandler() *InputHandler {
	return &InputHandler{
		inputChan: make(chan InputEvent, 10), // Buffered channel to prevent blocking
		done:      make(chan bool),
	}
}

// Start begins the input processing loop in a separate goroutine
func (h *InputHandler) Start() error {
	// Set terminal to raw mode for immediate key detection
	var err error
	h.oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	// Start input processing goroutine
	go h.processInput()
	return nil
}

// Stop stops the input processing and restores terminal state
func (h *InputHandler) Stop() error {
	close(h.done)

	// Restore terminal state
	if h.oldState != nil {
		return term.Restore(int(os.Stdin.Fd()), h.oldState)
	}
	return nil
}

// GetInputChannel returns the channel for receiving input events
func (h *InputHandler) GetInputChannel() <-chan InputEvent {
	return h.inputChan
}

// processInput runs in a separate goroutine to handle keyboard input
func (h *InputHandler) processInput() {
	buffer := make([]byte, 3) // Buffer for reading key sequences

	for {
		select {
		case <-h.done:
			return
		default:
			// Set a short read timeout to make this non-blocking
			n, err := os.Stdin.Read(buffer)
			if err != nil {
				continue
			}

			if n > 0 {
				key := h.parseKey(buffer[:n])
				if key != KeyUnknown {
					event := InputEvent{
						Key:  key,
						Time: time.Now(),
					}

					// Non-blocking send to channel
					select {
					case h.inputChan <- event:
					default:
						// Channel is full, drop the event
					}
				}
			}
		}
	}
}

// parseKey converts raw bytes to Key type
func (h *InputHandler) parseKey(data []byte) Key {
	if len(data) == 0 {
		return KeyUnknown
	}

	switch {
	case len(data) == 1:
		switch data[0] {
		case ' ':
			return KeySpace
		case 'q', 'Q':
			return KeyQ
		case 'r', 'R':
			return KeyR
		case 3: // Ctrl+C
			return KeyCtrlC
		}
	case len(data) == 3 && data[0] == 27 && data[1] == 91:
		// Arrow keys (ESC [ X)
		switch data[2] {
		case 65: // Up arrow
			return KeyUp
		}
	}

	return KeyUnknown
}
