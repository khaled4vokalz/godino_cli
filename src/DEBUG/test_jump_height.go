package main

import (
	"cli-dino-game/src/engine"
	"fmt"
	"math"
)

func simulateJump() {
	config := engine.NewDefaultConfig()
	
	// Jump physics calculation
	// At peak, velocity = 0, so: 0 = jumpVelocity - gravity * time
	// time_to_peak = jumpVelocity / gravity
	timeToPeak := config.JumpVelocity / config.Gravity
	
	// Height = initial_velocity * time - 0.5 * gravity * time^2
	// At peak: height = jumpVelocity * timeToPeak - 0.5 * gravity * timeToPeak^2
	maxHeight := config.JumpVelocity * timeToPeak - 0.5 * config.Gravity * math.Pow(timeToPeak, 2)
	
	fmt.Printf("=== Jump Physics Simulation ===\n")
	fmt.Printf("Jump Velocity: %.1f\n", config.JumpVelocity)
	fmt.Printf("Gravity: %.1f\n", config.Gravity)
	fmt.Printf("Time to peak: %.3f seconds\n", timeToPeak)
	fmt.Printf("Maximum jump height: %.1f units\n", maxHeight)
	
	// Check if this can clear obstacles
	fmt.Printf("\nCan clear obstacles requiring:\n")
	fmt.Printf("- 3 units height: %t\n", maxHeight >= 3)
	fmt.Printf("- 4 units height: %t\n", maxHeight >= 4)
	fmt.Printf("- 5 units height: %t\n", maxHeight >= 5)
}

func main() {
	simulateJump()
}
