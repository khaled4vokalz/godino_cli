package score

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Score manages the current game score and high score tracking
type Score struct {
	Current    int       `json:"current"`
	High       int       `json:"high"`
	Distance   float64   `json:"distance"`
	StartTime  time.Time `json:"start_time"`
	LastUpdate time.Time `json:"last_update"`

	// Scoring configuration
	TimeMultiplier     int     `json:"time_multiplier"`     // Points per second
	ObstacleBonus      int     `json:"obstacle_bonus"`      // Bonus points per obstacle
	DistanceMultiplier float64 `json:"distance_multiplier"` // Points per distance unit

	// Internal tracking
	obstaclesPassed int
	gameStartTime   time.Time
	lastScoreTime   time.Time
}

// ScoreData represents the persistent score data
type ScoreData struct {
	HighScore int `json:"high_score"`
}

// NewScore creates a new Score instance with default configuration
func NewScore() *Score {
	return &Score{
		Current:            0,
		High:               0,
		Distance:           0,
		StartTime:          time.Now(),
		LastUpdate:         time.Now(),
		TimeMultiplier:     10,  // 10 points per second
		ObstacleBonus:      100, // 100 points per obstacle
		DistanceMultiplier: 1.0, // 1 point per distance unit
		obstaclesPassed:    0,
		gameStartTime:      time.Now(),
		lastScoreTime:      time.Now(),
	}
}

// NewScoreWithConfig creates a new Score instance with custom configuration
func NewScoreWithConfig(timeMultiplier, obstacleBonus int, distanceMultiplier float64) *Score {
	score := NewScore()
	score.TimeMultiplier = timeMultiplier
	score.ObstacleBonus = obstacleBonus
	score.DistanceMultiplier = distanceMultiplier
	return score
}

// Reset resets the current score for a new game
func (s *Score) Reset() {
	s.Current = 0
	s.Distance = 0
	s.obstaclesPassed = 0
	s.gameStartTime = time.Now()
	s.lastScoreTime = time.Now()
	s.StartTime = time.Now()
	s.LastUpdate = time.Now()
}

// Update updates the score based on time elapsed
func (s *Score) Update(deltaTime float64) {
	now := time.Now()

	// Update distance (assuming constant movement)
	s.Distance += deltaTime * 10.0 // Arbitrary distance units per second

	// Calculate time-based score
	timeSinceLastScore := now.Sub(s.lastScoreTime).Seconds()
	if timeSinceLastScore >= 1.0 { // Update score every second
		timePoints := int(timeSinceLastScore) * s.TimeMultiplier
		s.Current += timePoints
		s.lastScoreTime = now
	}

	// Add distance-based points
	distancePoints := int(s.Distance * s.DistanceMultiplier)
	if distancePoints > s.Current {
		s.Current = distancePoints + (s.obstaclesPassed * s.ObstacleBonus)
	}

	s.LastUpdate = now
}

// AddObstacleBonus adds bonus points for successfully passing an obstacle
func (s *Score) AddObstacleBonus() {
	s.obstaclesPassed++
	s.Current += s.ObstacleBonus
	s.LastUpdate = time.Now()
}

// GetCurrent returns the current score
func (s *Score) GetCurrent() int {
	return s.Current
}

// GetHigh returns the high score
func (s *Score) GetHigh() int {
	return s.High
}

// GetDistance returns the distance traveled
func (s *Score) GetDistance() float64 {
	return s.Distance
}

// GetObstaclesPassed returns the number of obstacles passed
func (s *Score) GetObstaclesPassed() int {
	return s.obstaclesPassed
}

// GetGameDuration returns how long the current game has been running
func (s *Score) GetGameDuration() time.Duration {
	return time.Since(s.gameStartTime)
}

// IsNewHighScore checks if the current score is a new high score
func (s *Score) IsNewHighScore() bool {
	return s.Current > s.High
}

// UpdateHighScore updates the high score if current score is higher
func (s *Score) UpdateHighScore() bool {
	if s.IsNewHighScore() {
		s.High = s.Current
		return true
	}
	return false
}

// GetScoreBreakdown returns a breakdown of how the score was calculated
func (s *Score) GetScoreBreakdown() map[string]int {
	timeScore := int(s.GetGameDuration().Seconds()) * s.TimeMultiplier
	obstacleScore := s.obstaclesPassed * s.ObstacleBonus
	distanceScore := int(s.Distance * s.DistanceMultiplier)

	return map[string]int{
		"time":      timeScore,
		"obstacles": obstacleScore,
		"distance":  distanceScore,
		"total":     s.Current,
	}
}

// String returns a formatted string representation of the score
func (s *Score) String() string {
	return fmt.Sprintf("Score: %d (High: %d) | Distance: %.1f | Obstacles: %d | Time: %v",
		s.Current, s.High, s.Distance, s.obstaclesPassed, s.GetGameDuration().Truncate(time.Second))
}

// getScoreFilePath returns the path to the score file
func getScoreFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	scoreDir := filepath.Join(homeDir, ".cli-dino-game")
	if err := os.MkdirAll(scoreDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create score directory: %w", err)
	}

	return filepath.Join(scoreDir, "scores.json"), nil
}

// LoadHighScore loads the high score from persistent storage
func LoadHighScore() (int, error) {
	filePath, err := getScoreFilePath()
	if err != nil {
		return 0, err
	}

	// If file doesn't exist, return 0 (no high score yet)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return 0, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read score file: %w", err)
	}

	var scoreData ScoreData
	if err := json.Unmarshal(data, &scoreData); err != nil {
		return 0, fmt.Errorf("failed to parse score file: %w", err)
	}

	return scoreData.HighScore, nil
}

// SaveHighScore saves the high score to persistent storage
func SaveHighScore(highScore int) error {
	filePath, err := getScoreFilePath()
	if err != nil {
		return err
	}

	scoreData := ScoreData{
		HighScore: highScore,
	}

	data, err := json.Marshal(scoreData)
	if err != nil {
		return fmt.Errorf("failed to marshal score data: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write score file: %w", err)
	}

	return nil
}

// LoadHighScoreInto loads the high score from persistent storage into the Score instance
func (s *Score) LoadHighScoreInto() error {
	highScore, err := LoadHighScore()
	if err != nil {
		return err
	}
	s.High = highScore
	return nil
}

// SaveHighScoreFrom saves the high score from the Score instance to persistent storage
func (s *Score) SaveHighScoreFrom() error {
	return SaveHighScore(s.High)
}

// FinalizeScore finalizes the score at game end, updating high score if necessary
func (s *Score) FinalizeScore() (bool, error) {
	isNewHigh := s.UpdateHighScore()
	if isNewHigh {
		if err := s.SaveHighScoreFrom(); err != nil {
			return isNewHigh, fmt.Errorf("failed to save new high score: %w", err)
		}
	}
	return isNewHigh, nil
}
