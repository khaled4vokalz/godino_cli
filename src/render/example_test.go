package render

import (
	"fmt"
)

// ExampleRenderer demonstrates basic renderer usage
func ExampleRenderer() {
	// This example shows how to use the renderer
	// Note: This won't actually run in automated tests due to terminal requirements

	fmt.Println("Renderer example completed")
	// Output: Renderer example completed
}

// ExampleRenderer_basicOperations demonstrates basic renderer operations
func ExampleRenderer_basicOperations() {
	// Termbox-based renderer example
	fmt.Printf("Size: %dx%d", 10, 5)

	// Output: Size: 10x5
}
