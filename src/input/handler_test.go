package input

import (
	"testing"
	"time"
)

func TestNewInputHandler(t *testing.T) {
	handler := NewInputHandler()

	if handler == nil {
		t.Fatal("NewInputHandler() returned nil")
	}

	if handler.inputChan == nil {
		t.Error("InputHandler.inputChan is nil")
	}

	if handler.done == nil {
		t.Error("InputHandler.done is nil")
	}
}

func TestGetInputChannel(t *testing.T) {
	handler := NewInputHandler()
	channel := handler.GetInputChannel()

	if channel == nil {
		t.Error("GetInputChannel() returned nil")
	}

	// Verify it's the same channel
	if channel != handler.inputChan {
		t.Error("GetInputChannel() returned different channel than internal inputChan")
	}
}

func TestInputHandlerChannelBuffering(t *testing.T) {
	handler := NewInputHandler()

	// Test that we can send multiple events without blocking
	// (since channel is buffered with size 10)
	for i := 0; i < 5; i++ {
		select {
		case handler.inputChan <- InputEvent{Key: KeySpace, Time: time.Now()}:
			// Success - channel accepted the event
		default:
			t.Errorf("Channel blocked on event %d, expected buffered channel to accept", i)
		}
	}

	// Verify we can read the events back
	for i := 0; i < 5; i++ {
		select {
		case event := <-handler.inputChan:
			if event.Key != KeySpace {
				t.Errorf("Event %d: expected KeySpace, got %v", i, event.Key)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Timeout waiting for event %d", i)
		}
	}
}

func TestInputHandlerStop(t *testing.T) {
	handler := NewInputHandler()

	// Test stopping without starting (should not panic)
	err := handler.Stop()
	if err != nil {
		t.Errorf("Stop() without Start() returned error: %v", err)
	}
}

// TestInputEventTiming verifies that InputEvent captures timing correctly
func TestInputEventTiming(t *testing.T) {
	before := time.Now()
	event := InputEvent{
		Key:  KeySpace,
		Time: time.Now(),
	}
	after := time.Now()

	if event.Time.Before(before) || event.Time.After(after) {
		t.Errorf("InputEvent time %v is not between %v and %v",
			event.Time, before, after)
	}
}
