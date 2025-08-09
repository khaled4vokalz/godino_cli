package background

import (
	"math"
	"math/rand"
	"time"
)

// BackgroundElementType represents different types of background elements
type BackgroundElementType int

const (
	Cloud BackgroundElementType = iota
	Hill
	Mountain
)

// BackgroundElement represents a decorative background element
type BackgroundElement struct {
	Type    BackgroundElementType
	X       float64 // Horizontal position
	Y       float64 // Vertical position
	Width   float64 // Width for positioning
	Height  float64 // Height for positioning
	Speed   float64 // Scroll speed (slower than obstacles for parallax effect)
	Active  bool    // Whether the element is active
	Variant int     // Different variants of the same type
}

// HillProfile represents a continuous hill landscape
type HillProfile struct {
	heights []float64 // Height at each X position
	width   int        // Total width of the profile
	offset  float64    // Current scroll offset
	speed   float64    // Scroll speed
}

// BackgroundManager manages all background elements
type BackgroundManager struct {
	elements       []*BackgroundElement
	hillProfile    *HillProfile
	screenWidth    float64
	screenHeight   float64
	groundLevel    float64
	rng            *rand.Rand
	lastCloudSpawn time.Time
}

// NewBackgroundManager creates a new background manager
func NewBackgroundManager(screenWidth, screenHeight, groundLevel float64) *BackgroundManager {
	bm := &BackgroundManager{
		elements:     make([]*BackgroundElement, 0, 20),
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		groundLevel:  groundLevel,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	
	// Create continuous hill profile
	bm.hillProfile = bm.generateHillProfile()
	
	return bm
}

// generateHillProfile creates a continuous hill landscape using sine waves
func (bm *BackgroundManager) generateHillProfile() *HillProfile {
	// Create a wide profile that extends beyond screen for smooth scrolling
	profileWidth := int(bm.screenWidth * 4) // 4x screen width for seamless looping
	heights := make([]float64, profileWidth)
	
	// Generate continuous hills using multiple sine waves for natural variation
	for x := 0; x < profileWidth; x++ {
		// Combine multiple sine waves for natural-looking hills
		normalizedX := float64(x) / float64(profileWidth) * 4 * math.Pi
		
		// Large rolling hills (base layer) - restored larger heights
		baseHeight := 6.0 + 4.0*math.Sin(normalizedX*0.5)
		
		// Medium hills (detail layer) - restored
		mediumHeight := 2.5 * math.Sin(normalizedX*1.2+1.5)
		
		// Small hills (fine detail layer) - restored
		smallHeight := 1.2 * math.Sin(normalizedX*2.3+0.7)
		
		// Combine all layers with some randomness
		totalHeight := baseHeight + mediumHeight + smallHeight + (bm.rng.Float64()-0.5)*1.0
		
		// Ensure minimum height and reasonable maximum
		if totalHeight < 2.0 {
			totalHeight = 2.0
		}
		if totalHeight > 15.0 { // Good height for dramatic but not overwhelming hills
			totalHeight = 15.0
		}
		
		heights[x] = totalHeight
	}
	
	return &HillProfile{
		heights: heights,
		width:   profileWidth,
		offset:  0,
		speed:   8.0, // Hills scroll at medium speed
	}
}

// Update updates all background elements
func (bm *BackgroundManager) Update(deltaTime float64) {
	// Update hill profile scrolling
	bm.hillProfile.offset += bm.hillProfile.speed * deltaTime
	// Loop the hills when we've scrolled through one cycle
	if bm.hillProfile.offset >= float64(bm.hillProfile.width)/2 {
		bm.hillProfile.offset = 0
	}
	
	// Spawn clouds periodically
	bm.spawnElements()

	// Update existing cloud elements
	for i := len(bm.elements) - 1; i >= 0; i-- {
		element := bm.elements[i]
		element.X -= element.Speed * deltaTime

		// Remove elements that have moved off-screen
		if element.X+element.Width < -20 {
			bm.removeElement(i)
		}
	}
}

// spawnElements creates new background elements when needed
func (bm *BackgroundManager) spawnElements() {
	now := time.Now()

	// Spawn clouds every 15-30 seconds
	if now.Sub(bm.lastCloudSpawn) > time.Duration(15000+bm.rng.Intn(15000))*time.Millisecond {
		bm.spawnCloud()
		bm.lastCloudSpawn = now
	}
}

// spawnCloud creates a new cloud element
func (bm *BackgroundManager) spawnCloud() {
	// Position clouds in the upper area of the screen
	upperArea := bm.screenHeight / 2.5
	cloudY := 1 + bm.rng.Float64()*upperArea

	cloud := &BackgroundElement{
		Type:    Cloud,
		X:       bm.screenWidth + 10,
		Y:       cloudY,
		Width:   12 + bm.rng.Float64()*8, // Variable width clouds
		Height:  2 + bm.rng.Float64()*1,  // Variable height clouds
		Speed:   3 + bm.rng.Float64()*2,  // Slow parallax movement
		Active:  true,
		Variant: bm.rng.Intn(3), // 3 different cloud shapes
	}
	bm.elements = append(bm.elements, cloud)
}

// GetHillHeightAt returns the hill height at a specific screen X coordinate
func (bm *BackgroundManager) GetHillHeightAt(screenX float64) float64 {
	// Calculate the position in the hill profile
	profileX := int(bm.hillProfile.offset + screenX)
	
	// Loop the profile seamlessly
	profileX = profileX % bm.hillProfile.width
	if profileX < 0 {
		profileX += bm.hillProfile.width
	}
	
	return bm.hillProfile.heights[profileX]
}

// IsActive returns whether the element is active (helper method)
func (be *BackgroundElement) IsActive() bool {
	return be.Active
}

// removeElement removes a background element at the specified index
func (bm *BackgroundManager) removeElement(index int) {
	lastIndex := len(bm.elements) - 1
	if index != lastIndex {
		bm.elements[index] = bm.elements[lastIndex]
	}
	bm.elements = bm.elements[:lastIndex]
}

// GetElements returns all active background elements
func (bm *BackgroundManager) GetElements() []*BackgroundElement {
	return bm.elements
}

// Reset clears all background elements
func (bm *BackgroundManager) Reset() {
	bm.elements = bm.elements[:0]
	bm.lastCloudSpawn = time.Now()
	// Regenerate hill profile for variety
	bm.hillProfile = bm.generateHillProfile()
}

// GetSprite returns the sprite for a background element
func (be *BackgroundElement) GetSprite(useUnicode bool) []string {
	if useUnicode {
		switch be.Type {
		case Cloud:
			switch be.Variant {
			case 0:
				return []string{
					"  ☁☁☁☁☁☁☁☁  ",
					"☁☁☁☁☁☁☁☁☁☁☁☁",
				}
			case 1:
				return []string{
					"   ☁☁☁☁☁☁   ",
					"☁☁☁☁☁☁☁☁☁☁☁☁",
				}
			case 2:
				return []string{
					"    ☁☁☁☁    ",
					"  ☁☁☁☁☁☁☁☁  ",
					"☁☁☁☁☁☁☁☁☁☁☁☁",
				}
			}
		}
	} else {
		// ASCII versions
		switch be.Type {
		case Cloud:
			switch be.Variant {
			case 0:
				return []string{
					"    .-~~~-.    ",
					"  .-~       ~-.",
				}
			case 1:
				return []string{
					"  .-~~~-.  ",
					" (       ) ",
				}
			case 2:
				return []string{
					"    .---.    ",
					"  .-~     ~-.",
					" (           )",
				}
			}
		}
	}

	// Fallback
	return []string{"?"}
}
