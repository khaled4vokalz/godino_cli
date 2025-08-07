package main

import (
	"cli-dino-game/src/engine"
	"cli-dino-game/src/entities"
	"fmt"
)

func debugBirdCollision() {
	// Create a config
	config := engine.NewDefaultConfig()
	config.ScreenWidth = 80
	config.ScreenHeight = 20

	// Simulate ground level calculation from main.go
	groundLevel := float64(config.ScreenHeight - 5) // Leave space for dinosaur sprite = 15
	dinosaur := entities.NewDinosaur(groundLevel)    // Dinosaur at y=15
	
	// Calculate the actual ground line position (where obstacles should sit)
	actualGroundY := groundLevel + dinosaur.Height // 15 + 4 = 19

	fmt.Printf("Screen height: %d\n", config.ScreenHeight)
	fmt.Printf("Ground level (dinosaur Y): %.1f\n", groundLevel)
	fmt.Printf("Dinosaur height: %.1f\n", dinosaur.Height)
	fmt.Printf("Actual ground Y: %.1f\n", actualGroundY)
	fmt.Printf("Dinosaur bounds: %s\n", dinosaur.GetBounds().String())

	// Create different types of birds at the same X position as dinosaur
	birdLow := entities.NewObstacle(entities.BirdLow, dinosaur.X, actualGroundY, config)
	birdMid := entities.NewObstacle(entities.BirdMid, dinosaur.X, actualGroundY, config)
	birdHigh := entities.NewObstacle(entities.BirdHigh, dinosaur.X, actualGroundY, config)

	fmt.Printf("\nBird Low bounds: %s\n", birdLow.GetBounds().String())
	fmt.Printf("Bird Mid bounds: %s\n", birdMid.GetBounds().String())
	fmt.Printf("Bird High bounds: %s\n", birdHigh.GetBounds().String())

	// Create game engine to test collision
	gameEngine := engine.NewGameEngine(config)
	fmt.Printf("Default collision tolerance: %.1f\n", gameEngine.GetCollisionTolerance())

	// Test collisions
	fmt.Printf("\nCollision tests:\n")
	fmt.Printf("Dinosaur vs Bird Low: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdLow.GetBounds()))
	fmt.Printf("Dinosaur vs Bird Mid: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdMid.GetBounds()))
	fmt.Printf("Dinosaur vs Bird High: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdHigh.GetBounds()))

	// Test with no tolerance
	gameEngine.SetCollisionTolerance(0)
	fmt.Printf("\nWith zero tolerance:\n")
	fmt.Printf("Dinosaur vs Bird Low: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdLow.GetBounds()))
	fmt.Printf("Dinosaur vs Bird Mid: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdMid.GetBounds()))
	fmt.Printf("Dinosaur vs Bird High: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdHigh.GetBounds()))

	// Print ASCII art to visualize
	fmt.Printf("\nDinosaur ASCII art:\n")
	dinoArt := dinosaur.GetASCIIArt()
	for i, line := range dinoArt {
		fmt.Printf("Y=%.1f: %s\n", dinosaur.Y+float64(i), line)
	}

	fmt.Printf("\nBird Low ASCII art:\n")
	birdArt := birdLow.GetASCIIArt()
	for i, line := range birdArt {
		fmt.Printf("Y=%.1f: %s\n", birdLow.Y+float64(i), line)
	}
}

func main() {
	debugBirdCollision()
}
