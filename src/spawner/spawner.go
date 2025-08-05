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
		minSpawnInterval: time.Millisecond * 800,  // Minimum 0.8 seconds between spawns
		maxSpawnInterval: time.Millisecond * 3000, // Maximum 3 seconds between spawns
		typeWeights: map[entities.ObstacleType]float64{
			entities.CactusSmall:  0.5, // 50% chance
			entities.CactusMedium: 0.3, // 30% chance
			entities.CactusLarge:  0.2, // 20% chance
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

	// Spawn obstacle off-screen to the right
	spawnX := s.screenWidth + 5.0 // Spawn 5 units off-screen

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

	// Add randomness to the interval (±30% variation)
	randomFactor := 0.7 + s.rng.Float64()*0.6 // Range: 0.7 to 1.3
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

// selectObstacleType chooses an obstacle type based on weighted distribution
func (s *ObstacleSpawner) selectObstacleType() entities.ObstacleType {
	// Calculate total weight
	totalWeight := 0.0
	for _, weight := range s.typeWeights {
		totalWeight += weight
	}

	// Generate random value
	randomValue := s.rng.Float64() * totalWeight

	// Select type based on cumulative weights
	cumulative := 0.0
	for obstType, weight := range s.typeWeights {
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
	// Gradually increase obstacle speed over time
	speedIncrease := 1.0 + (s.gameTime * 0.05 / 10.0) // 5% increase every 10 seconds
	maxSpeedMultiplier := 2.0                         // Cap at 2x speed

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
