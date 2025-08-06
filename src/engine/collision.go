package engine

import (
	"fmt"
)

// CollisionDetector handles collision detection between game entities
type CollisionDetector struct {
	// Debug mode for collision detection
	debugMode bool
}

// NewCollisionDetector creates a new collision detector
func NewCollisionDetector() *CollisionDetector {
	return &CollisionDetector{
		debugMode: false,
	}
}

// SetDebugMode enables or disables debug mode for collision detection
func (cd *CollisionDetector) SetDebugMode(enabled bool) {
	cd.debugMode = enabled
}

// CheckCollision performs AABB collision detection between two rectangles
func (cd *CollisionDetector) CheckCollision(rect1, rect2 Rectangle) bool {
	collision := rect1.Intersects(rect2)

	if cd.debugMode && collision {
		fmt.Printf("Collision detected: %s intersects %s\n", rect1.String(), rect2.String())
	}

	return collision
}

// CheckCollisionWithTolerance performs collision detection with a tolerance margin
// This can be used to make the game more forgiving by reducing the effective collision area
func (cd *CollisionDetector) CheckCollisionWithTolerance(rect1, rect2 Rectangle, tolerance float64) bool {
	// Reduce the collision rectangles by the tolerance amount
	adjustedRect1 := Rectangle{
		X:      rect1.X + tolerance,
		Y:      rect1.Y + tolerance,
		Width:  rect1.Width - (2 * tolerance),
		Height: rect1.Height - (2 * tolerance),
	}

	adjustedRect2 := Rectangle{
		X:      rect2.X + tolerance,
		Y:      rect2.Y + tolerance,
		Width:  rect2.Width - (2 * tolerance),
		Height: rect2.Height - (2 * tolerance),
	}

	// Ensure adjusted rectangles have positive dimensions
	if adjustedRect1.Width <= 0 || adjustedRect1.Height <= 0 ||
		adjustedRect2.Width <= 0 || adjustedRect2.Height <= 0 {
		return false
	}

	return cd.CheckCollision(adjustedRect1, adjustedRect2)
}

// GetCollisionInfo returns detailed information about a collision
type CollisionInfo struct {
	HasCollision bool
	OverlapX     float64
	OverlapY     float64
	OverlapArea  float64
}

// GetCollisionInfo returns detailed collision information between two rectangles
func (cd *CollisionDetector) GetCollisionInfo(rect1, rect2 Rectangle) CollisionInfo {
	info := CollisionInfo{
		HasCollision: cd.CheckCollision(rect1, rect2),
	}

	if !info.HasCollision {
		return info
	}

	// Calculate overlap dimensions
	overlapLeft := max(rect1.X, rect2.X)
	overlapRight := min(rect1.X+rect1.Width, rect2.X+rect2.Width)
	overlapTop := max(rect1.Y, rect2.Y)
	overlapBottom := min(rect1.Y+rect1.Height, rect2.Y+rect2.Height)

	info.OverlapX = overlapRight - overlapLeft
	info.OverlapY = overlapBottom - overlapTop
	info.OverlapArea = info.OverlapX * info.OverlapY

	return info
}

// Helper functions for min/max
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
