package engine

import (
	"testing"
)

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()

	if config == nil {
		t.Fatal("NewDefaultConfig returned nil")
	}

	// Test default values
	if config.ScreenWidth != 80 {
		t.Errorf("Expected ScreenWidth 80, got %d", config.ScreenWidth)
	}
	if config.ScreenHeight != 20 {
		t.Errorf("Expected ScreenHeight 20, got %d", config.ScreenHeight)
	}
	if config.TargetFPS != 30 {
		t.Errorf("Expected TargetFPS 30, got %d", config.TargetFPS)
	}
	if config.JumpVelocity != 15.0 {
		t.Errorf("Expected JumpVelocity 15.0, got %f", config.JumpVelocity)
	}
	if config.Gravity != 50.0 {
		t.Errorf("Expected Gravity 50.0, got %f", config.Gravity)
	}
	if config.ObstacleSpeed != 20.0 {
		t.Errorf("Expected ObstacleSpeed 20.0, got %f", config.ObstacleSpeed)
	}
	if config.SpawnRate != 2.0 {
		t.Errorf("Expected SpawnRate 2.0, got %f", config.SpawnRate)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid default config",
			config:      NewDefaultConfig(),
			expectError: false,
		},
		{
			name: "zero screen width",
			config: &Config{
				ScreenWidth:   0,
				ScreenHeight:  20,
				TargetFPS:     30,
				JumpVelocity:  15.0,
				Gravity:       50.0,
				ObstacleSpeed: 20.0,
				SpawnRate:     2.0,
			},
			expectError: true,
			errorMsg:    "screen width must be positive",
		},
		{
			name: "negative screen height",
			config: &Config{
				ScreenWidth:   80,
				ScreenHeight:  -5,
				TargetFPS:     30,
				JumpVelocity:  15.0,
				Gravity:       50.0,
				ObstacleSpeed: 20.0,
				SpawnRate:     2.0,
			},
			expectError: true,
			errorMsg:    "screen height must be positive",
		},
		{
			name: "zero target FPS",
			config: &Config{
				ScreenWidth:   80,
				ScreenHeight:  20,
				TargetFPS:     0,
				JumpVelocity:  15.0,
				Gravity:       50.0,
				ObstacleSpeed: 20.0,
				SpawnRate:     2.0,
			},
			expectError: true,
			errorMsg:    "target FPS must be positive",
		},
		{
			name: "screen width too small",
			config: &Config{
				ScreenWidth:   30,
				ScreenHeight:  20,
				TargetFPS:     30,
				JumpVelocity:  15.0,
				Gravity:       50.0,
				ObstacleSpeed: 20.0,
				SpawnRate:     2.0,
			},
			expectError: true,
			errorMsg:    "screen width too small (minimum 40)",
		},
		{
			name: "screen height too small",
			config: &Config{
				ScreenWidth:   80,
				ScreenHeight:  5,
				TargetFPS:     30,
				JumpVelocity:  15.0,
				Gravity:       50.0,
				ObstacleSpeed: 20.0,
				SpawnRate:     2.0,
			},
			expectError: true,
			errorMsg:    "screen height too small (minimum 10)",
		},
		{
			name: "FPS too high",
			config: &Config{
				ScreenWidth:   80,
				ScreenHeight:  20,
				TargetFPS:     150,
				JumpVelocity:  15.0,
				Gravity:       50.0,
				ObstacleSpeed: 20.0,
				SpawnRate:     2.0,
			},
			expectError: true,
			errorMsg:    "target FPS too high (maximum 120)",
		},
		{
			name: "negative jump velocity",
			config: &Config{
				ScreenWidth:   80,
				ScreenHeight:  20,
				TargetFPS:     30,
				JumpVelocity:  -5.0,
				Gravity:       50.0,
				ObstacleSpeed: 20.0,
				SpawnRate:     2.0,
			},
			expectError: true,
			errorMsg:    "jump velocity must be positive",
		},
		{
			name: "zero gravity",
			config: &Config{
				ScreenWidth:   80,
				ScreenHeight:  20,
				TargetFPS:     30,
				JumpVelocity:  15.0,
				Gravity:       0.0,
				ObstacleSpeed: 20.0,
				SpawnRate:     2.0,
			},
			expectError: true,
			errorMsg:    "gravity must be positive",
		},
		{
			name: "negative obstacle speed",
			config: &Config{
				ScreenWidth:   80,
				ScreenHeight:  20,
				TargetFPS:     30,
				JumpVelocity:  15.0,
				Gravity:       50.0,
				ObstacleSpeed: -10.0,
				SpawnRate:     2.0,
			},
			expectError: true,
			errorMsg:    "obstacle speed must be positive",
		},
		{
			name: "zero spawn rate",
			config: &Config{
				ScreenWidth:   80,
				ScreenHeight:  20,
				TargetFPS:     30,
				JumpVelocity:  15.0,
				Gravity:       50.0,
				ObstacleSpeed: 20.0,
				SpawnRate:     0.0,
			},
			expectError: true,
			errorMsg:    "spawn rate must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestGameStateString(t *testing.T) {
	tests := []struct {
		state    GameState
		expected string
	}{
		{StateMenu, "Menu"},
		{StatePlaying, "Playing"},
		{StateGameOver, "GameOver"},
		{GameState(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.state.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestConfigString(t *testing.T) {
	config := NewDefaultConfig()
	result := config.String()
	expected := "Config{Screen: 80x20, FPS: 30, Jump: 15.0, Gravity: 50.0, Speed: 20.0, Spawn: 2.0}"

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestRectangleIntersects(t *testing.T) {
	tests := []struct {
		name     string
		rect1    Rectangle
		rect2    Rectangle
		expected bool
	}{
		{
			name:     "overlapping rectangles",
			rect1:    Rectangle{X: 0, Y: 0, Width: 10, Height: 10},
			rect2:    Rectangle{X: 5, Y: 5, Width: 10, Height: 10},
			expected: true,
		},
		{
			name:     "non-overlapping rectangles",
			rect1:    Rectangle{X: 0, Y: 0, Width: 10, Height: 10},
			rect2:    Rectangle{X: 20, Y: 20, Width: 10, Height: 10},
			expected: false,
		},
		{
			name:     "touching rectangles (edge case)",
			rect1:    Rectangle{X: 0, Y: 0, Width: 10, Height: 10},
			rect2:    Rectangle{X: 10, Y: 0, Width: 10, Height: 10},
			expected: false,
		},
		{
			name:     "one rectangle inside another",
			rect1:    Rectangle{X: 0, Y: 0, Width: 20, Height: 20},
			rect2:    Rectangle{X: 5, Y: 5, Width: 5, Height: 5},
			expected: true,
		},
		{
			name:     "identical rectangles",
			rect1:    Rectangle{X: 10, Y: 10, Width: 15, Height: 15},
			rect2:    Rectangle{X: 10, Y: 10, Width: 15, Height: 15},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect1.Intersects(tt.rect2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
			// Test symmetry
			result2 := tt.rect2.Intersects(tt.rect1)
			if result2 != tt.expected {
				t.Errorf("Intersection should be symmetric. Expected %v, got %v", tt.expected, result2)
			}
		})
	}
}

func TestRectangleContains(t *testing.T) {
	tests := []struct {
		name     string
		rect1    Rectangle
		rect2    Rectangle
		expected bool
	}{
		{
			name:     "larger rectangle contains smaller",
			rect1:    Rectangle{X: 0, Y: 0, Width: 20, Height: 20},
			rect2:    Rectangle{X: 5, Y: 5, Width: 5, Height: 5},
			expected: true,
		},
		{
			name:     "smaller rectangle does not contain larger",
			rect1:    Rectangle{X: 5, Y: 5, Width: 5, Height: 5},
			rect2:    Rectangle{X: 0, Y: 0, Width: 20, Height: 20},
			expected: false,
		},
		{
			name:     "identical rectangles contain each other",
			rect1:    Rectangle{X: 10, Y: 10, Width: 15, Height: 15},
			rect2:    Rectangle{X: 10, Y: 10, Width: 15, Height: 15},
			expected: true,
		},
		{
			name:     "overlapping but not containing",
			rect1:    Rectangle{X: 0, Y: 0, Width: 10, Height: 10},
			rect2:    Rectangle{X: 5, Y: 5, Width: 10, Height: 10},
			expected: false,
		},
		{
			name:     "non-overlapping rectangles",
			rect1:    Rectangle{X: 0, Y: 0, Width: 10, Height: 10},
			rect2:    Rectangle{X: 20, Y: 20, Width: 10, Height: 10},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rect1.Contains(tt.rect2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRectangleCenter(t *testing.T) {
	tests := []struct {
		name      string
		rect      Rectangle
		expectedX float64
		expectedY float64
	}{
		{
			name:      "rectangle at origin",
			rect:      Rectangle{X: 0, Y: 0, Width: 10, Height: 10},
			expectedX: 5.0,
			expectedY: 5.0,
		},
		{
			name:      "rectangle with offset",
			rect:      Rectangle{X: 10, Y: 20, Width: 6, Height: 8},
			expectedX: 13.0,
			expectedY: 24.0,
		},
		{
			name:      "rectangle with decimal values",
			rect:      Rectangle{X: 1.5, Y: 2.5, Width: 3.0, Height: 4.0},
			expectedX: 3.0,
			expectedY: 4.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := tt.rect.Center()
			if x != tt.expectedX {
				t.Errorf("Expected center X %f, got %f", tt.expectedX, x)
			}
			if y != tt.expectedY {
				t.Errorf("Expected center Y %f, got %f", tt.expectedY, y)
			}
		})
	}
}

func TestRectangleString(t *testing.T) {
	rect := Rectangle{X: 1.5, Y: 2.5, Width: 10.0, Height: 15.0}
	result := rect.String()
	expected := "Rectangle{X: 1.5, Y: 2.5, W: 10.0, H: 15.0}"

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
