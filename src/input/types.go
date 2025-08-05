package input

import "time"

// Key represents different keyboard inputs
type Key int

const (
	KeySpace Key = iota
	KeyUp
	KeyQ
	KeyR
	KeyCtrlC
	KeyUnknown
)

// InputEvent represents a keyboard input event
type InputEvent struct {
	Key  Key
	Time time.Time
}

// String returns a string representation of the Key
func (k Key) String() string {
	switch k {
	case KeySpace:
		return "Space"
	case KeyUp:
		return "Up"
	case KeyQ:
		return "Q"
	case KeyR:
		return "R"
	case KeyCtrlC:
		return "Ctrl+C"
	default:
		return "Unknown"
	}
}
