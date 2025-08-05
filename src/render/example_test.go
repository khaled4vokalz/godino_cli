package render

import (
	"fmt"
)

// ExampleRenderer demonstrates basic renderer usage
func ExampleRenderer() {
	// This example shows how to use the renderer
	// Note: This won't actually run in automated tests due to terminal requirements
	
	renderer := &Renderer{
		width:  20,
		height: 10,
	}
	renderer.initBuffer()
	
	// Draw some content
	renderer.DrawString(2, 2, "Hello, World!")
	renderer.DrawBox(1, 1, 18, 8, '*')
	renderer.DrawString(5, 5, "CLI Dino Game")
	
	fmt.Println("Renderer example completed")
	// Output: Renderer example completed
}

// ExampleRenderer_basicOperations demonstrates basic renderer operations
func ExampleRenderer_basicOperations() {
	renderer := &Renderer{
		width:  10,
		height: 5,
	}
	renderer.initBuffer()
	
	// Clear the buffer
	renderer.Clear()
	
	// Draw at specific position
	renderer.DrawAt(5, 2, 'X')
	
	// Draw a string
	renderer.DrawString(0, 0, "Test")
	
	// Get size
	width, height := renderer.GetSize()
	fmt.Printf("Size: %dx%d", width, height)
	
	// Output: Size: 10x5
}