// Package input provides non-blocking keyboard input handling for the CLI dino game.
//
// The input package implements a concurrent input handling system that captures
// keyboard events in raw terminal mode and delivers them through Go channels.
// This allows the game to respond immediately to user input without blocking
// the main game loop.
//
// Key Features:
//   - Non-blocking keyboard input detection
//   - Support for game controls (Space, Up Arrow, Q, R, Ctrl+C)
//   - Channel-based communication for concurrent processing
//   - Proper terminal state management and restoration
//   - Input event timestamping for precise timing
//
// Basic Usage:
//
//	handler := input.NewInputHandler()
//	err := handler.Start()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer handler.Stop()
//
//	inputChan := handler.GetInputChannel()
//	for {
//		select {
//		case event := <-inputChan:
//			fmt.Printf("Key pressed: %s at %v\n", event.Key.String(), event.Time)
//		default:
//			// Continue with other game logic
//		}
//	}
//
// The package is designed to work with the game engine's main loop,
// providing responsive controls for the dinosaur jumping mechanics
// and game state management.
package input
