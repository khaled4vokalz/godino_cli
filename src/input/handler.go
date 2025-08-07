package input

import (
	"time"

	"github.com/nsf/termbox-go"
)

// InputHandler manages keyboard input using termbox-go
type InputHandler struct {
	inputChan chan InputEvent
	done      chan bool
}

// NewInputHandler creates a new termbox-based InputHandler instance
func NewInputHandler() *InputHandler {
	return &InputHandler{
		inputChan: make(chan InputEvent, 10), // Buffered channel to prevent blocking
		done:      make(chan bool),
	}
}

// Start begins the input processing loop in a separate goroutine
func (h *InputHandler) Start() error {
	// Start input processing goroutine
	go h.processInput()
	return nil
}

// Stop stops the input processing
func (h *InputHandler) Stop() error {
	close(h.done)
	return nil
}

// GetInputChannel returns the channel for receiving input events
func (h *InputHandler) GetInputChannel() <-chan InputEvent {
	return h.inputChan
}

// processInput runs in a separate goroutine to handle keyboard input using termbox
func (h *InputHandler) processInput() {
	for {
		select {
		case <-h.done:
			return
		default:
			// Poll for events with timeout
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				key := h.parseTermboxKey(ev)
				if key != KeyUnknown {
					event := InputEvent{
						Key:  key,
						Time: time.Now(),
					}

					// Non-blocking send to channel
					select {
					case h.inputChan <- event:
					default:
						// Channel full, drop event
					}
				}
			case termbox.EventResize:
				// Handle resize events if needed
				continue
			}
		}
	}
}

// parseTermboxKey converts termbox key events to our Key type
func (h *InputHandler) parseTermboxKey(ev termbox.Event) Key {
	switch {
	case ev.Key == termbox.KeySpace:
		return KeySpace
	case ev.Key == termbox.KeyArrowUp:
		return KeyUp
	case ev.Key == termbox.KeyCtrlC:
		return KeyCtrlC
	case ev.Ch != 0:
		// Handle character keys
		switch ev.Ch {
		case 'q', 'Q':
			return KeyQ
		case 'r', 'R':
			return KeyR
		default:
			return KeyUnknown
		}
	default:
		return KeyUnknown
	}
}
