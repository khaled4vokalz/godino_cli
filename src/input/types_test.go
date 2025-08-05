package input

import (
	"testing"
	"time"
)

func TestKeyString(t *testing.T) {
	tests := []struct {
		key      Key
		expected string
	}{
		{KeySpace, "Space"},
		{KeyUp, "Up"},
		{KeyQ, "Q"},
		{KeyR, "R"},
		{KeyCtrlC, "Ctrl+C"},
		{KeyUnknown, "Unknown"},
	}

	for _, test := range tests {
		result := test.key.String()
		if result != test.expected {
			t.Errorf("Key.String() for %v: expected %s, got %s", test.key, test.expected, result)
		}
	}
}

func TestInputEvent(t *testing.T) {
	now := time.Now()
	event := InputEvent{
		Key:  KeySpace,
		Time: now,
	}

	if event.Key != KeySpace {
		t.Errorf("InputEvent.Key: expected %v, got %v", KeySpace, event.Key)
	}

	if event.Time != now {
		t.Errorf("InputEvent.Time: expected %v, got %v", now, event.Time)
	}
}
