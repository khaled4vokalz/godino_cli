package main

import (
	"cli-dino-game/src/engine"
	"cli-dino-game/src/entities"
	"fmt"
)

func testCollisionTolerance() {
	// Create a config
	config := engine.NewDefaultConfig()
	config.ScreenWidth = 80
	config.ScreenHeight = 20

	// Simulate ground level calculation from main.go
	groundLevel := float64(config.ScreenHeight - 5) // Leave space for dinosaur sprite = 15
	dinosaur := entities.NewDinosaur(groundLevel)    // Dinosaur at y=15
	
	// Calculate the actual ground line position (where obstacles should sit)
	actualGroundY := groundLevel + dinosaur.Height // 15 + 4 = 19

	fmt.Printf("=== Collision Tolerance Test (0.8) ===\n")
	fmt.Printf("Dinosaur bounds: %s\n", dinosaur.GetBounds().String())

	// Create different types of obstacles at the same X position as dinosaur
	smallCactus := entities.NewObstacle(entities.CactusSmall, dinosaur.X, actualGroundY, config)
	mediumCactus := entities.NewObstacle(entities.CactusMedium, dinosaur.X, actualGroundY, config)
	largeCactus := entities.NewObstacle(entities.CactusLarge, dinosaur.X, actualGroundY, config)
	birdLow := entities.NewObstacle(entities.BirdLow, dinosaur.X, actualGroundY, config)
	birdMid := entities.NewObstacle(entities.BirdMid, dinosaur.X, actualGroundY, config)
	birdHigh := entities.NewObstacle(entities.BirdHigh, dinosaur.X, actualGroundY, config)

	fmt.Printf("Small Cactus bounds: %s\n", smallCactus.GetBounds().String())
	fmt.Printf("Medium Cactus bounds: %s\n", mediumCactus.GetBounds().String())
	fmt.Printf("Large Cactus bounds: %s\n", largeCactus.GetBounds().String())
	fmt.Printf("Bird Low bounds: %s\n", birdLow.GetBounds().String())
	fmt.Printf("Bird Mid bounds: %s\n", birdMid.GetBounds().String())
	fmt.Printf("Bird High bounds: %s\n", birdHigh.GetBounds().String())

	// Create game engine to test collision
	gameEngine := engine.NewGameEngine(config)
	fmt.Printf("Default collision tolerance: %.1f\n\n", gameEngine.GetCollisionTolerance())

	// Test collisions with default tolerance (0.8)
	fmt.Printf("Collision tests with tolerance %.1f:\n", gameEngine.GetCollisionTolerance())
	fmt.Printf("Dinosaur vs Small Cactus: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), smallCactus.GetBounds()))
	fmt.Printf("Dinosaur vs Medium Cactus: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), mediumCactus.GetBounds()))
	fmt.Printf("Dinosaur vs Large Cactus: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), largeCactus.GetBounds()))
	fmt.Printf("Dinosaur vs Bird Low: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdLow.GetBounds()))
	fmt.Printf("Dinosaur vs Bird Mid: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdMid.GetBounds()))
	fmt.Printf("Dinosaur vs Bird High: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdHigh.GetBounds()))

	// Test with no tolerance to see raw collision
	gameEngine.SetCollisionTolerance(0)
	fmt.Printf("\nWith zero tolerance:\n")
	fmt.Printf("Dinosaur vs Small Cactus: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), smallCactus.GetBounds()))
	fmt.Printf("Dinosaur vs Medium Cactus: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), mediumCactus.GetBounds()))
	fmt.Printf("Dinosaur vs Large Cactus: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), largeCactus.GetBounds()))
	fmt.Printf("Dinosaur vs Bird Low: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdLow.GetBounds()))
	fmt.Printf("Dinosaur vs Bird Mid: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdMid.GetBounds()))
	fmt.Printf("Dinosaur vs Bird High: %t\n", gameEngine.CheckCollision(dinosaur.GetBounds(), birdHigh.GetBounds()))

	// Test with jumping dinosaur vs large cactus to see if we can clear it
	fmt.Printf("\n=== Jump Test ===\n")
	jumpingDino := entities.NewDinosaur(groundLevel)
	
	// Test different jump heights
	jumpHeights := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	for _, height := range jumpHeights {
		jumpingDino.Y = groundLevel - height
		fmt.Printf("Jump height %.1f (Y=%.1f): ", height, jumpingDino.Y)
		
		gameEngine.SetCollisionTolerance(0.8)
		collision := gameEngine.CheckCollision(jumpingDino.GetBounds(), largeCactus.GetBounds())
		fmt.Printf("vs Large Cactus = %t\n", collision)
	}
}

func main() {
	testCollisionTolerance()
}
