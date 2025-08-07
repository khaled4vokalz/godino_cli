package input

import (
	"fmt"
	"time"
)

// ExampleInputHandler demonstrates basic usage of the InputHandler
func ExampleInputHandler() {
	handler := NewInputHandler()

	// Get the input channel
	inputChan := handler.GetInputChannel()

	// Simulate receiving an input event
	event := InputEvent{
		Key:  KeySpace,
		Time: time.Now(),
	}

	fmt.Printf("Key: %s\n", event.Key.String())

	// Demonstrate key types
	fmt.Printf("Space key: %s\n", KeySpace.String())
	fmt.Printf("Up arrow: %s\n", KeyUp.String())

	// Channel is available for communication
	fmt.Printf("Input channel ready: %t\n", inputChan != nil)

	// Output:
	// Key: Space
	// Space key: Space
	// Up arrow: Up
	// Input channel ready: true
}
