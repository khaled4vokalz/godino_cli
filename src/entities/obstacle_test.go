package entities

import (
	"cli-dino-game/src/engine"
	"testing"
)

func TestNewObstacle(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 15.0
	startX := 80.0

	tests := []struct {
		name      string
		obstType  ObstacleType
		expectedW float64
		expectedH float64
		expectedY float64
	}{
		{
			name:      "Small cactus",
			obstType:  CactusSmall,
			expectedW: 2.0,
			expectedH: 4.0,
			expectedY: 12.0, // groundLevel - height + 1
		},
		{
			name:      "Medium cactus",
			obstType:  CactusMedium,
			expectedW: 3.0,
			expectedH: 6.0,
			expectedY: 10.0, // groundLevel - height + 1
		},
		{
			name:      "Large cactus",
			obstType:  CactusLarge,
			expectedW: 4.0,
			expectedH: 8.0,
			expectedY: 8.0, // groundLevel - height + 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obstacle := NewObstacle(tt.obstType, startX, groundLevel, config)

			if obstacle.X != startX {
				t.Errorf("Expected X position %f, got %f", startX, obstacle.X)
			}
			if obstacle.Y != tt.expectedY {
				t.Errorf("Expected Y position %f, got %f", tt.expectedY, obstacle.Y)
			}
			if obstacle.Width != tt.expectedW {
				t.Errorf("Expected width %f, got %f", tt.expectedW, obstacle.Width)
			}
			if obstacle.Height != tt.expectedH {
				t.Errorf("Expected height %f, got %f", tt.expectedH, obstacle.Height)
			}
			if obstacle.Speed != config.ObstacleSpeed {
				t.Errorf("Expected speed %f, got %f", config.ObstacleSpeed, obstacle.Speed)
			}
			if obstacle.ObstType != tt.obstType {
				t.Errorf("Expected obstacle type %v, got %v", tt.obstType, obstacle.ObstType)
			}
			if !obstacle.Active {
				t.Error("Expected obstacle to be active")
			}
		})
	}
}

func TestObstacleUpdate(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 15.0
	startX := 80.0
	deltaTime := 1.0 / 30.0 // 30 FPS

	obstacle := NewObstacle(CactusSmall, startX, groundLevel, config)
	initialX := obstacle.X

	// Test normal movement
	obstacle.Update(deltaTime)
	expectedX := initialX - (config.ObstacleSpeed * deltaTime)
	if obstacle.X != expectedX {
		t.Errorf("Expected X position %f after update, got %f", expectedX, obstacle.X)
	}
	if !obstacle.Active {
		t.Error("Expected obstacle to remain active")
	}

	// Test obstacle going off-screen
	obstacle.X = -3.0 // Position where obstacle + width < 0 (width is 2.0, so -3.0 + 2.0 = -1.0 < 0)
	obstacle.Update(deltaTime)
	if obstacle.Active {
		t.Error("Expected obstacle to be inactive after going off-screen")
	}
}

func TestObstacleUpdateInactive(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 15.0
	startX := 80.0
	deltaTime := 1.0 / 30.0

	obstacle := NewObstacle(CactusSmall, startX, groundLevel, config)
	obstacle.Active = false
	initialX := obstacle.X

	obstacle.Update(deltaTime)

	// Inactive obstacles should not move
	if obstacle.X != initialX {
		t.Errorf("Expected inactive obstacle to not move, but X changed from %f to %f", initialX, obstacle.X)
	}
}

func TestObstacleGetBounds(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 15.0
	startX := 80.0

	obstacle := NewObstacle(CactusMedium, startX, groundLevel, config)
	bounds := obstacle.GetBounds()

	if bounds.X != obstacle.X {
		t.Errorf("Expected bounds X %f, got %f", obstacle.X, bounds.X)
	}
	if bounds.Y != obstacle.Y {
		t.Errorf("Expected bounds Y %f, got %f", obstacle.Y, bounds.Y)
	}
	if bounds.Width != obstacle.Width {
		t.Errorf("Expected bounds width %f, got %f", obstacle.Width, bounds.Width)
	}
	if bounds.Height != obstacle.Height {
		t.Errorf("Expected bounds height %f, got %f", obstacle.Height, bounds.Height)
	}
}

func TestObstacleGetASCIIArt(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 15.0
	startX := 80.0

	tests := []struct {
		name         string
		obstType     ObstacleType
		expectedRows int
	}{
		{
			name:         "Small cactus ASCII",
			obstType:     CactusSmall,
			expectedRows: 3,
		},
		{
			name:         "Medium cactus ASCII",
			obstType:     CactusMedium,
			expectedRows: 4,
		},
		{
			name:         "Large cactus ASCII",
			obstType:     CactusLarge,
			expectedRows: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obstacle := NewObstacle(tt.obstType, startX, groundLevel, config)
			art := obstacle.GetASCIIArt()

			if len(art) != tt.expectedRows {
				t.Errorf("Expected %d rows of ASCII art, got %d", tt.expectedRows, len(art))
			}

			// Check that all rows are non-empty
			for i, row := range art {
				if len(row) == 0 {
					t.Errorf("Row %d of ASCII art is empty", i)
				}
			}
		})
	}
}

func TestObstaclePositionMethods(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 15.0
	startX := 80.0

	obstacle := NewObstacle(CactusSmall, startX, groundLevel, config)

	// Test GetPosition
	x, y := obstacle.GetPosition()
	if x != startX {
		t.Errorf("Expected X position %f, got %f", startX, x)
	}
	if y != obstacle.Y {
		t.Errorf("Expected Y position %f, got %f", obstacle.Y, y)
	}

	// Test SetPosition
	newX, newY := 50.0, 10.0
	obstacle.SetPosition(newX, newY)
	x, y = obstacle.GetPosition()
	if x != newX {
		t.Errorf("Expected X position %f after SetPosition, got %f", newX, x)
	}
	if y != newY {
		t.Errorf("Expected Y position %f after SetPosition, got %f", newY, y)
	}
}

func TestObstacleActiveMethods(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 15.0
	startX := 80.0

	obstacle := NewObstacle(CactusSmall, startX, groundLevel, config)

	// Test initial active state
	if !obstacle.IsActive() {
		t.Error("Expected new obstacle to be active")
	}

	// Test deactivation
	obstacle.Deactivate()
	if obstacle.IsActive() {
		t.Error("Expected obstacle to be inactive after Deactivate()")
	}
}

func TestObstacleGetters(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 15.0
	startX := 80.0

	obstacle := NewObstacle(CactusMedium, startX, groundLevel, config)

	// Test GetType
	if obstacle.GetType() != CactusMedium {
		t.Errorf("Expected obstacle type %v, got %v", CactusMedium, obstacle.GetType())
	}

	// Test GetSpeed
	if obstacle.GetSpeed() != config.ObstacleSpeed {
		t.Errorf("Expected speed %f, got %f", config.ObstacleSpeed, obstacle.GetSpeed())
	}

	// Test SetSpeed
	newSpeed := 25.0
	obstacle.SetSpeed(newSpeed)
	if obstacle.GetSpeed() != newSpeed {
		t.Errorf("Expected speed %f after SetSpeed, got %f", newSpeed, obstacle.GetSpeed())
	}

	// Test GetDimensions
	width, height := obstacle.GetDimensions()
	if width != obstacle.Width {
		t.Errorf("Expected width %f, got %f", obstacle.Width, width)
	}
	if height != obstacle.Height {
		t.Errorf("Expected height %f, got %f", obstacle.Height, height)
	}
}

func TestObstacleIsOffScreen(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 15.0

	obstacle := NewObstacle(CactusSmall, 10.0, groundLevel, config)

	// Test on-screen obstacle
	if obstacle.IsOffScreen() {
		t.Error("Expected obstacle at X=10 to be on-screen")
	}

	// Test off-screen obstacle (X + Width < 0)
	obstacle.SetPosition(-3.0, groundLevel) // Width is 2.0, so -3.0 + 2.0 = -1.0 < 0
	if !obstacle.IsOffScreen() {
		t.Error("Expected obstacle at X=-3 to be off-screen")
	}

	// Test edge case (exactly at edge)
	obstacle.SetPosition(-2.0, groundLevel) // Width is 2.0, so -2.0 + 2.0 = 0.0
	if obstacle.IsOffScreen() {
		t.Error("Expected obstacle at X=-2 (edge) to not be off-screen")
	}
}

func TestObstacleTypeString(t *testing.T) {
	tests := []struct {
		obstType ObstacleType
		expected string
	}{
		{CactusSmall, "CactusSmall"},
		{CactusMedium, "CactusMedium"},
		{CactusLarge, "CactusLarge"},
		{ObstacleType(999), "Unknown"}, // Invalid type
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.obstType.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestObstacleMovementLifecycle(t *testing.T) {
	config := engine.NewDefaultConfig()
	groundLevel := 15.0
	startX := 80.0
	deltaTime := 1.0 / 30.0 // 30 FPS

	obstacle := NewObstacle(CactusSmall, startX, groundLevel, config)

	// Simulate obstacle moving across screen
	steps := 0
	for obstacle.IsActive() && steps < 1000 { // Safety limit
		obstacle.Update(deltaTime)
		steps++

		// Verify obstacle is moving left
		if obstacle.X > startX {
			t.Error("Obstacle should be moving left (decreasing X)")
			break
		}
	}

	// Verify obstacle became inactive due to going off-screen
	if obstacle.IsActive() {
		t.Error("Expected obstacle to become inactive after moving off-screen")
	}
	if !obstacle.IsOffScreen() {
		t.Error("Expected obstacle to be off-screen when inactive")
	}

	// Verify it took a reasonable number of steps
	if steps < 10 {
		t.Errorf("Expected obstacle to take more steps to cross screen, took %d", steps)
	}
}
