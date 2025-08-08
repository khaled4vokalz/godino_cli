package main

import (
	"cli-dino-game/src/background"
	"cli-dino-game/src/engine"
	"cli-dino-game/src/entities"
	"cli-dino-game/src/input"
	"cli-dino-game/src/render"
	"cli-dino-game/src/spawner"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Game represents the main game application
type Game struct {
	engine       *engine.GameEngine
	renderer     *render.Renderer
	inputHandler *input.InputHandler
	dinosaur     *entities.Dinosaur
	spawner      *spawner.ObstacleSpawner
	background   *background.BackgroundManager
	config       *engine.Config

	// Game loop control
	running bool
	ticker  *time.Ticker

	// Graceful shutdown
	shutdownChan chan os.Signal
}

// NewGame creates a new game instance
func NewGame() (*Game, error) {
	// Create default configuration
	config := engine.NewDefaultConfig()
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create renderer
	renderer, err := render.NewRenderer()
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}

	// Update config with actual terminal size
	termWidth, termHeight := renderer.GetSize()
	config.ScreenWidth = termWidth
	config.ScreenHeight = termHeight

	// Create game engine
	gameEngine := engine.NewGameEngine(config)

	// Create input handler
	inputHandler := input.NewInputHandler()

	// Create dinosaur
	groundLevel := float64(config.ScreenHeight - 5) // Leave space for dinosaur sprite
	dinosaur := entities.NewDinosaur(groundLevel)

	// Calculate the actual ground line position (where obstacles should sit)
	actualGroundY := groundLevel + dinosaur.Height

	// Create obstacle spawner
	obstacleSpawner := spawner.NewObstacleSpawner(config, float64(config.ScreenWidth), actualGroundY)

	// Create background manager
	backgroundManager := background.NewBackgroundManager(float64(config.ScreenWidth), float64(config.ScreenHeight), actualGroundY)

	// Setup graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	game := &Game{
		engine:       gameEngine,
		renderer:     renderer,
		inputHandler: inputHandler,
		dinosaur:     dinosaur,
		spawner:      obstacleSpawner,
		background:   backgroundManager,
		config:       config,
		running:      false,
		shutdownChan: shutdownChan,
	}

	return game, nil
}

// Run starts the main game loop
func (g *Game) Run() error {
	// Termbox is already initialized by the renderer
	defer g.renderer.Close()

	// Start input handler
	if err := g.inputHandler.Start(); err != nil {
		return fmt.Errorf("failed to start input handler: %w", err)
	}
	defer g.inputHandler.Stop()

	// Setup game loop timing
	frameDuration := time.Second / time.Duration(g.config.TargetFPS)
	g.ticker = time.NewTicker(frameDuration)
	defer g.ticker.Stop()

	// Initialize game state
	g.running = true
	g.engine.SetState(engine.StateMenu)

	// Main game loop
	for g.running {
		select {
		case <-g.ticker.C:
			// Update game state
			g.update()
			// Render frame
			g.render()

		case inputEvent := <-g.inputHandler.GetInputChannel():
			// Handle input
			g.handleInput(inputEvent)

		case <-g.shutdownChan:
			// Graceful shutdown
			g.shutdown()
			return nil
		}
	}

	return nil
}

// update handles all game logic updates
func (g *Game) update() {
	// Update game engine timing
	g.engine.Update()
	deltaTime := g.engine.GetDeltaTime()

	switch g.engine.GetState() {
	case engine.StatePlaying:
		// Update dinosaur
		g.dinosaur.Update(deltaTime, g.config)

		// Update obstacle spawner
		g.spawner.Update(deltaTime)

		// Update background elements
		g.background.Update(deltaTime)

		// Check collisions
		g.checkCollisions()

	case engine.StateGameOver:
		// Game over state - no updates needed

	case engine.StateMenu:
		// Menu state - minimal updates
	}
}

// render handles all rendering
func (g *Game) render() {
	// Clear screen buffer
	g.renderer.Clear()

	switch g.engine.GetState() {
	case engine.StateMenu:
		g.renderMenu()

	case engine.StatePlaying:
		g.renderGame()

	case engine.StateGameOver:
		g.renderGameOver()
	}

	// Flush buffer to screen
	g.renderer.Flush()
}

// renderMenu renders the main menu
func (g *Game) renderMenu() {
	// Use the new start screen renderer
	g.renderer.DrawStartScreen()
}

// renderGame renders the main gameplay
func (g *Game) renderGame() {
	// Render ground line
	width, _ := g.renderer.GetSize()
	groundY := int(g.dinosaur.GroundLevel) + int(g.dinosaur.Height)
	groundChar := '-'
	if g.config.UseUnicode {
		groundChar = '▔'
	}
	for x := 0; x < width; x++ {
		g.renderer.DrawAt(x, groundY, groundChar)
	}

	// Render background elements (behind everything else)
	g.renderBackground()

	// Render dinosaur
	g.renderDinosaur()

	// Render obstacles
	g.renderObstacles()

	// Render UI
	g.renderUI()
}

// renderDinosaur renders the dinosaur sprite
func (g *Game) renderDinosaur() {
	art := g.dinosaur.GetASCIIArtWithConfig(g.config.UseUnicode)
	x := int(g.dinosaur.X)
	y := int(g.dinosaur.Y)

	for i, line := range art {
		g.renderer.DrawString(x, y+i, line)
	}
}

// renderObstacles renders all active obstacles
func (g *Game) renderObstacles() {
	obstacles := g.spawner.GetObstacles()
	for _, obstacle := range obstacles {
		if obstacle.IsActive() {
			art := obstacle.GetASCIIArtWithConfig(g.config.UseUnicode)
			x := int(obstacle.X)
			y := int(obstacle.Y)

			for i, line := range art {
				g.renderer.DrawString(x, y+i, line)
			}
		}
	}
}

// renderBackground renders background elements (clouds, hills)
func (g *Game) renderBackground() {
	elements := g.background.GetElements()
	for _, element := range elements {
		sprite := element.GetSprite(g.config.UseUnicode)
		x := int(element.X)
		y := int(element.Y)

		// Use different colors for different elements
		color := "ash" // Default for clouds
		if element.Type == background.Hill {
			color = "dark" // Darker color for hills
		}

		for i, line := range sprite {
			g.renderer.DrawStringWithColor(x, y+i, line, color)
		}
	}
}

// renderUI renders the game UI (score, etc.)
func (g *Game) renderUI() {
	// Use the new score display renderer
	g.renderer.DrawScore(g.engine.GetCurrentScore(), g.engine.GetHighScore())

	// Draw control instructions at the bottom
	g.renderer.DrawControlInstructions()
}

// renderGameOver renders the game over screen
func (g *Game) renderGameOver() {
	// Use the new game over screen renderer
	g.renderer.DrawGameOverScreen(
		g.engine.GetCurrentScore(),
		g.engine.GetHighScore(),
		g.engine.IsNewHighScore(),
	)
}

// handleInput processes input events
func (g *Game) handleInput(event input.InputEvent) {
	switch event.Key {
	case input.KeyCtrlC, input.KeyQ:
		g.shutdown()

	case input.KeySpace, input.KeyUp:
		switch g.engine.GetState() {
		case engine.StateMenu:
			g.startGame()
		case engine.StatePlaying:
			g.dinosaur.Jump(g.config)
		}

	case input.KeyR:
		if g.engine.GetState() == engine.StateGameOver {
			g.restartGame()
		}
	}
}

// startGame starts a new game
func (g *Game) startGame() {
	g.engine.Start()
	g.spawner.Reset()
	g.background.Reset()
}

// restartGame restarts the game from game over state
func (g *Game) restartGame() {
	g.engine.Restart()
	g.spawner.Reset()
	g.background.Reset()
}

// checkCollisions checks for collisions between dinosaur and obstacles
func (g *Game) checkCollisions() {
	dinosaurBounds := g.dinosaur.GetBounds()
	obstacles := g.spawner.GetObstacles()

	for _, obstacle := range obstacles {
		if obstacle.IsActive() {
			obstacleBounds := obstacle.GetBounds()
			if g.engine.CheckCollision(dinosaurBounds, obstacleBounds) {
				g.engine.TriggerGameOver()
				return
			}
		}
	}

	// Award points for obstacles that have passed the dinosaur
	for _, obstacle := range obstacles {
		if obstacle.IsActive() && obstacle.X+obstacle.Width < g.dinosaur.X {
			g.engine.AddObstacleBonus()
			obstacle.Deactivate() // Prevent multiple bonuses for same obstacle
		}
	}
}

// shutdown gracefully shuts down the game
func (g *Game) shutdown() {
	g.running = false
}

// Cleanup performs cleanup operations
func (g *Game) Cleanup() {
	if g.ticker != nil {
		g.ticker.Stop()
	}
	g.engine.Cleanup()
}

func main() {
	// Parse command line flags
	useUnicode := flag.Bool("unicode", true, "Use Unicode characters for rendering (default: true for better visuals)")
	asciiMode := flag.Bool("ascii", false, "Use ASCII characters instead of Unicode (for terminals with poor Unicode support)")
	flag.Parse()

	// Create game instance
	game, err := NewGame()
	if err != nil {
		log.Fatalf("Failed to create game: %v", err)
	}
	defer game.Cleanup()

	// Set Unicode preference
	if *asciiMode {
		game.config.UseUnicode = false
	} else {
		game.config.UseUnicode = *useUnicode
	}

	// Run the game
	if err := game.Run(); err != nil {
		log.Fatalf("Game error: %v", err)
	}

	fmt.Println("Thanks for playing CLI Dino Game!")
}
