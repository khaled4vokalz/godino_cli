package score

import (
	"os"
	"testing"
)

func TestNewScore(t *testing.T) {
	score := NewScore()

	if score.Current != 0 {
		t.Errorf("Expected current score to be 0, got %d", score.Current)
	}

	if score.High != 0 {
		t.Errorf("Expected high score to be 0, got %d", score.High)
	}

	if score.Distance != 0 {
		t.Errorf("Expected distance to be 0, got %f", score.Distance)
	}

	if score.TimeMultiplier != 10 {
		t.Errorf("Expected time multiplier to be 10, got %d", score.TimeMultiplier)
	}

	if score.ObstacleBonus != 100 {
		t.Errorf("Expected obstacle bonus to be 100, got %d", score.ObstacleBonus)
	}
}

func TestNewScoreWithConfig(t *testing.T) {
	timeMultiplier := 20
	obstacleBonus := 200
	distanceMultiplier := 2.0

	score := NewScoreWithConfig(timeMultiplier, obstacleBonus, distanceMultiplier)

	if score.TimeMultiplier != timeMultiplier {
		t.Errorf("Expected time multiplier to be %d, got %d", timeMultiplier, score.TimeMultiplier)
	}

	if score.ObstacleBonus != obstacleBonus {
		t.Errorf("Expected obstacle bonus to be %d, got %d", obstacleBonus, score.ObstacleBonus)
	}

	if score.DistanceMultiplier != distanceMultiplier {
		t.Errorf("Expected distance multiplier to be %f, got %f", distanceMultiplier, score.DistanceMultiplier)
	}
}

func TestScoreReset(t *testing.T) {
	score := NewScore()
	score.Current = 500
	score.Distance = 100.0
	score.obstaclesPassed = 5

	score.Reset()

	if score.Current != 0 {
		t.Errorf("Expected current score to be 0 after reset, got %d", score.Current)
	}

	if score.Distance != 0 {
		t.Errorf("Expected distance to be 0 after reset, got %f", score.Distance)
	}

	if score.obstaclesPassed != 0 {
		t.Errorf("Expected obstacles passed to be 0 after reset, got %d", score.obstaclesPassed)
	}
}
func TestScoreUpdate(t *testing.T) {
	score := NewScore()
	score.Reset()

	// Simulate 1 second passing
	deltaTime := 1.0
	score.Update(deltaTime)

	// Check that distance increased
	if score.Distance <= 0 {
		t.Errorf("Expected distance to increase after update, got %f", score.Distance)
	}

	// Check that score increased based on time
	if score.Current <= 0 {
		t.Errorf("Expected current score to increase after update, got %d", score.Current)
	}
}

func TestAddObstacleBonus(t *testing.T) {
	score := NewScore()
	initialScore := score.Current

	score.AddObstacleBonus()

	expectedScore := initialScore + score.ObstacleBonus
	if score.Current != expectedScore {
		t.Errorf("Expected score to be %d after obstacle bonus, got %d", expectedScore, score.Current)
	}

	if score.GetObstaclesPassed() != 1 {
		t.Errorf("Expected obstacles passed to be 1, got %d", score.GetObstaclesPassed())
	}
}

func TestIsNewHighScore(t *testing.T) {
	score := NewScore()
	score.High = 100
	score.Current = 50

	if score.IsNewHighScore() {
		t.Error("Expected IsNewHighScore to return false when current < high")
	}

	score.Current = 150
	if !score.IsNewHighScore() {
		t.Error("Expected IsNewHighScore to return true when current > high")
	}
}

func TestUpdateHighScore(t *testing.T) {
	score := NewScore()
	score.High = 100
	score.Current = 150

	updated := score.UpdateHighScore()

	if !updated {
		t.Error("Expected UpdateHighScore to return true for new high score")
	}

	if score.High != 150 {
		t.Errorf("Expected high score to be updated to 150, got %d", score.High)
	}

	// Test no update when current score is lower
	score.Current = 120
	updated = score.UpdateHighScore()

	if updated {
		t.Error("Expected UpdateHighScore to return false when no new high score")
	}

	if score.High != 150 {
		t.Errorf("Expected high score to remain 150, got %d", score.High)
	}
}
func TestGetScoreBreakdown(t *testing.T) {
	score := NewScore()
	score.Current = 500
	score.Distance = 50.0
	score.obstaclesPassed = 3

	breakdown := score.GetScoreBreakdown()

	if breakdown["obstacles"] != 300 { // 3 * 100
		t.Errorf("Expected obstacle score to be 300, got %d", breakdown["obstacles"])
	}

	if breakdown["distance"] != 50 { // 50.0 * 1.0
		t.Errorf("Expected distance score to be 50, got %d", breakdown["distance"])
	}

	if breakdown["total"] != 500 {
		t.Errorf("Expected total score to be 500, got %d", breakdown["total"])
	}
}

func TestScorePersistence(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test saving high score
	testHighScore := 1500
	err := SaveHighScore(testHighScore)
	if err != nil {
		t.Fatalf("Failed to save high score: %v", err)
	}

	// Test loading high score
	loadedScore, err := LoadHighScore()
	if err != nil {
		t.Fatalf("Failed to load high score: %v", err)
	}

	if loadedScore != testHighScore {
		t.Errorf("Expected loaded score to be %d, got %d", testHighScore, loadedScore)
	}
}

func TestScorePersistenceWithScore(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	score := NewScore()
	score.High = 2000

	// Test saving through Score instance
	err := score.SaveHighScoreFrom()
	if err != nil {
		t.Fatalf("Failed to save high score from Score instance: %v", err)
	}

	// Create new score instance and load
	newScore := NewScore()
	err = newScore.LoadHighScoreInto()
	if err != nil {
		t.Fatalf("Failed to load high score into Score instance: %v", err)
	}

	if newScore.High != 2000 {
		t.Errorf("Expected loaded high score to be 2000, got %d", newScore.High)
	}
}
func TestFinalizeScore(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	score := NewScore()
	score.High = 100
	score.Current = 200

	isNewHigh, err := score.FinalizeScore()
	if err != nil {
		t.Fatalf("Failed to finalize score: %v", err)
	}

	if !isNewHigh {
		t.Error("Expected FinalizeScore to return true for new high score")
	}

	if score.High != 200 {
		t.Errorf("Expected high score to be updated to 200, got %d", score.High)
	}

	// Verify persistence
	loadedScore, err := LoadHighScore()
	if err != nil {
		t.Fatalf("Failed to load high score after finalize: %v", err)
	}

	if loadedScore != 200 {
		t.Errorf("Expected persisted high score to be 200, got %d", loadedScore)
	}
}

func TestLoadHighScoreNonExistentFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Override the home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Try to load from non-existent file
	score, err := LoadHighScore()
	if err != nil {
		t.Fatalf("Expected no error when loading non-existent file, got: %v", err)
	}

	if score != 0 {
		t.Errorf("Expected score to be 0 when file doesn't exist, got %d", score)
	}
}

func TestScoreString(t *testing.T) {
	score := NewScore()
	score.Current = 500
	score.High = 1000
	score.Distance = 75.5
	score.obstaclesPassed = 5

	str := score.String()

	// Check that the string contains expected values
	if str == "" {
		t.Error("Expected non-empty string representation")
	}

	// The string should contain the score values
	// We won't check exact format since it might change, but ensure it's not empty
	t.Logf("Score string representation: %s", str)
}

func TestScoreGetters(t *testing.T) {
	score := NewScore()
	score.Current = 300
	score.High = 500
	score.Distance = 45.5
	score.obstaclesPassed = 3

	if score.GetCurrent() != 300 {
		t.Errorf("Expected GetCurrent() to return 300, got %d", score.GetCurrent())
	}

	if score.GetHigh() != 500 {
		t.Errorf("Expected GetHigh() to return 500, got %d", score.GetHigh())
	}

	if score.GetDistance() != 45.5 {
		t.Errorf("Expected GetDistance() to return 45.5, got %f", score.GetDistance())
	}

	if score.GetObstaclesPassed() != 3 {
		t.Errorf("Expected GetObstaclesPassed() to return 3, got %d", score.GetObstaclesPassed())
	}

	duration := score.GetGameDuration()
	if duration < 0 {
		t.Errorf("Expected GetGameDuration() to return positive duration, got %v", duration)
	}
}
