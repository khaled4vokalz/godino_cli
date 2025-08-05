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

	// Demonstrate key parsing
	handler2 := NewInputHandler()
	spaceKey := handler2.parseKey([]byte{' '})
	upKey := handler2.parseKey([]byte{27, 91, 65})

	fmt.Printf("Space key parsed: %s\n", spaceKey.String())
	fmt.Printf("Up arrow parsed: %s\n", upKey.String())

	// Channel is available for communication
	fmt.Printf("Input channel ready: %t\n", inputChan != nil)

	// Output:
	// Key: Space
	// Space key parsed: Space
	// Up arrow parsed: Up
	// Input channel ready: true
}
