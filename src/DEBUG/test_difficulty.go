package main

import (
	"cli-dino-game/src/engine"
	"cli-dino-game/src/spawner"
	"fmt"
)

func testDifficultyProgression() {
	config := engine.NewDefaultConfig()
	config.ScreenWidth = 80
	config.ScreenHeight = 20

	spawnerInstance := spawner.NewObstacleSpawner(config, 80, 19)

	fmt.Printf("=== Difficulty Progression Test ===\n")
	fmt.Printf("Base spawn rate: %.2f obstacles/sec\n", config.SpawnRate)
	fmt.Printf("Max spawn rate: %.2f obstacles/sec\n\n", config.SpawnRate*2.0)

	timePoints := []float64{0, 10, 20, 30, 40, 60, 90, 120, 180}

	for _, gameTime := range timePoints {
		// Simulate game time
		spawnerInstance.SetDifficulty(config.SpawnRate, config.SpawnRate*2.0, 0.02)
		
		// Calculate current metrics using the same formulas as in spawner
		// Spawn rate calculation
		difficultyMultiplier := 1.0 + (gameTime * 0.02 / 30.0)
		currentSpawnRate := config.SpawnRate * difficultyMultiplier
		if currentSpawnRate > config.SpawnRate*2.0 {
			currentSpawnRate = config.SpawnRate*2.0
		}

		// Speed multiplier calculation
		speedIncrease := 1.0 + (gameTime * 0.02 / 10.0)
		if speedIncrease > 1.8 {
			speedIncrease = 1.8
		}

		// Gap calculation
		baseMinGap := 25.0
		difficultyReduction := gameTime * 0.1
		if difficultyReduction > 8.0 {
			difficultyReduction = 8.0
		}
		minGap := baseMinGap - difficultyReduction
		if minGap < 18.0 {
			minGap = 18.0
		}

		// Bird availability
		birdsAvailable := gameTime > 30.0
		birdPercentage := 0.0
		if birdsAvailable {
			birdMultiplier := (gameTime - 30.0) / 60.0
			if birdMultiplier > 1.0 {
				birdMultiplier = 1.0
			}
			birdPercentage = (0.05 + 0.03 + 0.02) * birdMultiplier * 100
		}

		fmt.Printf("Time: %3.0fs | Spawn: %.2f/s | Speed: %.2fx | MinGap: %.0f | Birds: %.1f%%\n",
			gameTime, currentSpawnRate, speedIncrease, minGap, birdPercentage)
	}
}

func main() {
	testDifficultyProgression()
}
