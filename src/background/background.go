package background

import (
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

// BackgroundManager manages all background elements
type BackgroundManager struct {
	elements       []*BackgroundElement
	screenWidth    float64
	screenHeight   float64
	groundLevel    float64
	rng            *rand.Rand
	lastCloudSpawn time.Time
	lastHillSpawn  time.Time
}

// NewBackgroundManager creates a new background manager
func NewBackgroundManager(screenWidth, screenHeight, groundLevel float64) *BackgroundManager {
	return &BackgroundManager{
		elements:     make([]*BackgroundElement, 0, 20),
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		groundLevel:  groundLevel,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Update updates all background elements
func (bm *BackgroundManager) Update(deltaTime float64) {
	// Spawn new elements periodically
	bm.spawnElements()

	// Update existing elements
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

	// Spawn clouds every 20-40 seconds (very infrequent for larger, detailed clouds)
	if now.Sub(bm.lastCloudSpawn) > time.Duration(20000+bm.rng.Intn(20000))*time.Millisecond {
		bm.spawnCloud()
		bm.lastCloudSpawn = now
	}

	// Spawn hills much more frequently to ensure continuous coverage
	if now.Sub(bm.lastHillSpawn) > time.Duration(2000+bm.rng.Intn(3000))*time.Millisecond {
		bm.spawnHill()
		bm.lastHillSpawn = now
	}
}

// spawnCloud creates a new cloud element
func (bm *BackgroundManager) spawnCloud() {
	// Position clouds in the middle area of the screen (around 1/3 from top)
	middleArea := bm.screenHeight / 3.0
	cloudY := 2 + bm.rng.Float64()*middleArea // Clouds in middle area

	cloud := &BackgroundElement{
		Type:    Cloud,
		X:       bm.screenWidth + 15,    // Start further off-screen for larger clouds
		Y:       cloudY,                 // Position in middle area
		Width:   15,                     // Wider for proper cloud shapes
		Height:  2,                      // Exactly 2 lines tall
		Speed:   2 + bm.rng.Float64()*1, // Slower movement for larger clouds
		Active:  true,
		Variant: bm.rng.Intn(2), // 2 different cloud shapes
	}
	bm.elements = append(bm.elements, cloud)
}

// spawnHill creates a new hill element
func (bm *BackgroundManager) spawnHill() {
	// Create larger, more varied hills
	hillHeight := 3 + bm.rng.Float64()*6 // Larger hills (3-9 units tall)

	// Find the rightmost hill to connect to it
	rightmostX := bm.screenWidth
	for _, element := range bm.elements {
		if element.Type == Hill && element.IsActive() {
			elementRight := element.X + element.Width
			if elementRight > rightmostX {
				rightmostX = elementRight
			}
		}
	}

	hill := &BackgroundElement{
		Type:    Hill,
		X:       rightmostX - 2, // Overlap slightly to connect hills
		Y:       bm.groundLevel - hillHeight,
		Width:   15 + bm.rng.Float64()*10, // Much wider hills (15-25 units)
		Height:  hillHeight,
		Speed:   6 + bm.rng.Float64()*2, // Faster movement for hills
		Active:  true,
		Variant: bm.rng.Intn(2), // 2 different hill shapes
	}
	bm.elements = append(bm.elements, hill)
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
	bm.lastHillSpawn = time.Now()
}

// GetSprite returns the sprite for a background element
func (be *BackgroundElement) GetSprite(useUnicode bool) []string {
	if useUnicode {
		switch be.Type {
		case Cloud:
			switch be.Variant {
			case 0:
				return []string{
					"  ▁▁▁▁▁▁▁▁▁  ",
					"▁▁▁▁▁▁▁▁▁▁▁▁▁",
				}
			case 1:
				return []string{
					"   ▁▁▁▁▁▁   ",
					"▁▁▁▁▁▁▁▁▁▁▁▁",
				}
			}
		case Hill:
			switch be.Variant {
			case 0:
				return []string{
					"      ▓▓▓▓▓▓      ",
					"    ▓▓▓▓▓▓▓▓▓▓    ",
					"  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ",
					"▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓",
				}
			case 1:
				return []string{
					"    ▓▓▓▓▓▓▓▓    ",
					"  ▓▓▓▓▓▓▓▓▓▓▓▓  ",
					"▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓",
				}
			}
		}
	} else {
		// ASCII versions - inspired by asciiart.eu cloud designs
		switch be.Type {
		case Cloud:
			switch be.Variant {
			case 0:
				return []string{
					"    .-~~~-.    ",
					"  .-~       ~-.",
					" (             )",
				}
			case 1:
				return []string{
					"  .-~~~-.  ",
					" (       ) ",
				}
			}
		case Hill:
			switch be.Variant {
			case 0:
				return []string{
					"      ######      ",
					"    ############  ",
					"  ################",
					"##################",
				}
			case 1:
				return []string{
					"    ########    ",
					"  ############  ",
					"################",
				}
			}
		}
	}

	// Fallback
	return []string{"?"}
}
