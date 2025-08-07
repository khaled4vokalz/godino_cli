package spawner

import (
	"cli-dino-game/src/engine"
	"cli-dino-game/src/entities"
	"math/rand"
	"time"
)

// ObstacleSpawner manages the spawning of obstacles with random intervals and difficulty progression
type ObstacleSpawner struct {
	config         *engine.Config
	obstacles      []*entities.Obstacle
	lastSpawnTime  time.Time
	nextSpawnDelay time.Duration
	gameTime       float64
	screenWidth    float64
	groundLevel    float64
	rng            *rand.Rand

	// Difficulty progression parameters
	baseSpawnRate    float64 // Base spawn rate (obstacles per second)
	maxSpawnRate     float64 // Maximum spawn rate
	difficultyRamp   float64 // How quickly difficulty increases
	minSpawnInterval time.Duration
	maxSpawnInterval time.Duration

	// Obstacle type distribution
	typeWeights map[entities.ObstacleType]float64
}

// NewObstacleSpawner creates a new obstacle spawner
func NewObstacleSpawner(config *engine.Config, screenWidth, groundLevel float64) *ObstacleSpawner {
	spawner := &ObstacleSpawner{
		config:           config,
		obstacles:        make([]*entities.Obstacle, 0, 10), // Pre-allocate for efficiency
		screenWidth:      screenWidth,
		groundLevel:      groundLevel,
		rng:              rand.New(rand.NewSource(time.Now().UnixNano())),
		baseSpawnRate:    config.SpawnRate,
		maxSpawnRate:     config.SpawnRate * 3.0,  // Max 3x base rate
		difficultyRamp:   0.1,                     // Difficulty increases by 10% every 10 seconds
		minSpawnInterval: time.Millisecond * 300,  // Minimum 0.3 seconds between spawns
		maxSpawnInterval: time.Millisecond * 2500, // Maximum 2.5 seconds between spawns
		typeWeights: map[entities.ObstacleType]float64{
			entities.CactusSmall:  0.35, // 35% chance
			entities.CactusMedium: 0.25, // 25% chance
			entities.CactusLarge:  0.15, // 15% chance
			entities.BirdLow:      0.10, // 10% chance
			entities.BirdMid:      0.10, // 10% chance
			entities.BirdHigh:     0.05, // 5% chance
		},
	}

	// Initialize first spawn delay
	spawner.scheduleNextSpawn()
	return spawner
}

// Update updates the spawner and manages obstacle spawning
func (s *ObstacleSpawner) Update(deltaTime float64) {
	s.gameTime += deltaTime

	// Check if it's time to spawn a new obstacle
	if time.Since(s.lastSpawnTime) >= s.nextSpawnDelay {
		s.spawnObstacle()
		s.scheduleNextSpawn()
	}

	// Update all active obstacles
	for i := len(s.obstacles) - 1; i >= 0; i-- {
		obstacle := s.obstacles[i]
		obstacle.Update(deltaTime)

		// Remove inactive obstacles for memory efficiency
		if !obstacle.IsActive() {
			s.removeObstacle(i)
		}
	}
}

// spawnObstacle creates and spawns a new obstacle
func (s *ObstacleSpawner) spawnObstacle() {
	// Choose obstacle type based on weighted distribution
	obstType := s.selectObstacleType()

	// Calculate spawn position with proper spacing
	spawnX := s.calculateSpawnPosition()

	// Create new obstacle
	obstacle := entities.NewObstacle(obstType, spawnX, s.groundLevel, s.config)

	// Apply current difficulty speed multiplier
	speedMultiplier := s.getDifficultySpeedMultiplier()
	obstacle.SetSpeed(obstacle.GetSpeed() * speedMultiplier)

	// Add to obstacle list
	s.obstacles = append(s.obstacles, obstacle)
	s.lastSpawnTime = time.Now()
}

// scheduleNextSpawn calculates the delay until the next obstacle spawn
func (s *ObstacleSpawner) scheduleNextSpawn() {
	// Calculate current spawn rate based on difficulty progression
	currentSpawnRate := s.getCurrentSpawnRate()

	// Convert spawn rate to interval (seconds between spawns)
	baseInterval := 1.0 / currentSpawnRate

	// Add randomness to the interval (±50% variation)
	randomFactor := 0.5 + s.rng.Float64()*1.0 // Range: 0.5 to 1.5
	interval := time.Duration(baseInterval*randomFactor*1000) * time.Millisecond

	// Clamp to min/max intervals
	if interval < s.minSpawnInterval {
		interval = s.minSpawnInterval
	}
	if interval > s.maxSpawnInterval {
		interval = s.maxSpawnInterval
	}

	s.nextSpawnDelay = interval
}

// calculateSpawnPosition determines where to spawn the next obstacle with random spacing
func (s *ObstacleSpawner) calculateSpawnPosition() float64 {
	// Base spawn position (just off-screen)
	baseSpawnX := s.screenWidth + 2.0

	// Find the rightmost active obstacle that's still relevant for spacing
	rightmostX := baseSpawnX
	for _, obstacle := range s.obstacles {
		if obstacle.IsActive() && obstacle.X > s.screenWidth-20.0 { // Only consider obstacles near or off-screen
			obstacleRight := obstacle.X + obstacle.Width
			if obstacleRight > rightmostX {
				rightmostX = obstacleRight
			}
		}
	}

	// Define minimum and maximum gaps between obstacles
	minGap := 15.0 // Minimum distance for jumpability
	maxGap := 45.0 // Maximum distance to keep game challenging

	// Generate random gap within the range
	randomGap := minGap + s.rng.Float64()*(maxGap-minGap)

	// Calculate final spawn position
	spawnX := rightmostX + randomGap

	// Ensure we don't spawn too close to screen edge
	if spawnX < baseSpawnX {
		spawnX = baseSpawnX + randomGap
	}

	return spawnX
}

// selectObstacleType chooses an obstacle type based on weighted distribution and game time
func (s *ObstacleSpawner) selectObstacleType() entities.ObstacleType {
	// Create dynamic weights based on game time
	weights := make(map[entities.ObstacleType]float64)

	// Always include cacti
	weights[entities.CactusSmall] = 0.5
	weights[entities.CactusMedium] = 0.3
	weights[entities.CactusLarge] = 0.2

	// Only include birds after 5 seconds of gameplay (even earlier for testing)
	if s.gameTime > 5.0 {
		// Much more aggressive bird introduction
		birdMultiplier := (s.gameTime - 5.0) / 10.0 // Reach full strength in just 10 seconds
		if birdMultiplier > 1.0 {
			birdMultiplier = 1.0
		}

		// Much higher bird weights - make them very visible
		weights[entities.BirdLow] = 0.3 * birdMultiplier  // 30% at full strength
		weights[entities.BirdMid] = 0.2 * birdMultiplier  // 20% at full strength
		weights[entities.BirdHigh] = 0.1 * birdMultiplier // 10% at full strength

		// Significantly reduce cactus weights to make room for birds
		totalBirdWeight := weights[entities.BirdLow] + weights[entities.BirdMid] + weights[entities.BirdHigh]
		weights[entities.CactusSmall] = 0.5 - (totalBirdWeight * 0.4)
		weights[entities.CactusMedium] = 0.3 - (totalBirdWeight * 0.3)
		weights[entities.CactusLarge] = 0.2 - (totalBirdWeight * 0.3)
	}

	// Calculate total weight
	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight
	}

	// Generate random value
	randomValue := s.rng.Float64() * totalWeight

	// Select type based on cumulative weights
	cumulative := 0.0
	for obstType, weight := range weights {
		cumulative += weight
		if randomValue <= cumulative {
			return obstType
		}
	}

	// Fallback to small cactus
	return entities.CactusSmall
}

// getCurrentSpawnRate calculates the current spawn rate based on difficulty progression
func (s *ObstacleSpawner) getCurrentSpawnRate() float64 {
	// Increase spawn rate over time
	difficultyMultiplier := 1.0 + (s.gameTime * s.difficultyRamp / 10.0)
	currentRate := s.baseSpawnRate * difficultyMultiplier

	// Cap at maximum spawn rate
	if currentRate > s.maxSpawnRate {
		currentRate = s.maxSpawnRate
	}

	return currentRate
}

// getDifficultySpeedMultiplier calculates speed multiplier based on game time
func (s *ObstacleSpawner) getDifficultySpeedMultiplier() float64 {
	// Gradually increase obstacle speed over time - more aggressive progression
	speedIncrease := 1.0 + (s.gameTime * 0.1 / 5.0) // 10% increase every 5 seconds
	maxSpeedMultiplier := 2.5                       // Cap at 2.5x speed

	if speedIncrease > maxSpeedMultiplier {
		speedIncrease = maxSpeedMultiplier
	}

	return speedIncrease
}

// removeObstacle removes an obstacle at the specified index
func (s *ObstacleSpawner) removeObstacle(index int) {
	// Efficient removal by swapping with last element
	lastIndex := len(s.obstacles) - 1
	if index != lastIndex {
		s.obstacles[index] = s.obstacles[lastIndex]
	}
	s.obstacles = s.obstacles[:lastIndex]
}

// GetObstacles returns all active obstacles
func (s *ObstacleSpawner) GetObstacles() []*entities.Obstacle {
	return s.obstacles
}

// GetActiveObstacleCount returns the number of active obstacles
func (s *ObstacleSpawner) GetActiveObstacleCount() int {
	return len(s.obstacles)
}

// Reset resets the spawner state for a new game
func (s *ObstacleSpawner) Reset() {
	s.obstacles = s.obstacles[:0] // Clear slice but keep capacity
	s.gameTime = 0.0
	s.lastSpawnTime = time.Now()
	s.scheduleNextSpawn()
}

// SetDifficulty allows manual adjustment of difficulty parameters
func (s *ObstacleSpawner) SetDifficulty(baseRate, maxRate, ramp float64) {
	s.baseSpawnRate = baseRate
	s.maxSpawnRate = maxRate
	s.difficultyRamp = ramp
}

// SetObstacleTypeWeights allows customization of obstacle type distribution
func (s *ObstacleSpawner) SetObstacleTypeWeights(weights map[entities.ObstacleType]float64) {
	s.typeWeights = make(map[entities.ObstacleType]float64)
	for obstType, weight := range weights {
		if weight > 0 {
			s.typeWeights[obstType] = weight
		}
	}
}

// GetGameTime returns the current game time
func (s *ObstacleSpawner) GetGameTime() float64 {
	return s.gameTime
}

// GetCurrentSpawnRate returns the current spawn rate for debugging/display
func (s *ObstacleSpawner) GetCurrentSpawnRate() float64 {
	return s.getCurrentSpawnRate()
}

// GetNextSpawnDelay returns the time until next spawn for debugging/display
func (s *ObstacleSpawner) GetNextSpawnDelay() time.Duration {
	elapsed := time.Since(s.lastSpawnTime)
	if elapsed >= s.nextSpawnDelay {
		return 0
	}
	return s.nextSpawnDelay - elapsed
}
