// Package render provides terminal rendering capabilities for the CLI Dino Game.
//
// This package handles all terminal output, screen management, and provides
// a buffer-based rendering system for smooth game updates. It includes
// functionality for:
//
//   - Terminal raw mode setup and restoration
//   - Screen clearing and cursor positioning
//   - Buffer-based rendering for smooth updates
//   - Drawing primitives (characters, strings, boxes)
//   - Terminal size detection and handling
//
// The main type is Renderer, which manages the terminal state and provides
// methods for drawing to a screen buffer that can be flushed to the terminal.
//
// Example usage:
//
//	renderer, err := render.NewRenderer()
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer renderer.RestoreTerminal()
//
//	err = renderer.SetRawMode()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	renderer.Clear()
//	renderer.DrawString(10, 5, "Hello, World!")
//	renderer.Flush()
package render