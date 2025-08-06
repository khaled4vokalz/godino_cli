package main

import (
	"cli-dino-game/src/engine"
	"cli-dino-game/src/entities"
	"cli-dino-game/src/input"
	"cli-dino-game/src/render"
	"cli-dino-game/src/spawner"
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
	groundLevel := float64(config.ScreenHeight - 7) // Leave space for dinosaur sprite
	dinosaur := entities.NewDinosaur(groundLevel)

	// Create obstacle spawner
	obstacleSpawner := spawner.NewObstacleSpawner(config, float64(config.ScreenWidth), groundLevel)

	// Setup graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	game := &Game{
		engine:       gameEngine,
		renderer:     renderer,
		inputHandler: inputHandler,
		dinosaur:     dinosaur,
		spawner:      obstacleSpawner,
		config:       config,
		running:      false,
		shutdownChan: shutdownChan,
	}

	return game, nil
}

// Run starts the main game loop
func (g *Game) Run() error {
	// Initialize terminal
	if err := g.renderer.SetRawMode(); err != nil {
		return fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer g.renderer.RestoreTerminal()

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

	// Clear screen and hide cursor
	g.renderer.ClearScreen()
	g.renderer.HideCursor()

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
	width, height := g.renderer.GetSize()

	// Title
	title := "CLI DINO GAME"
	titleX := (width - len(title)) / 2
	titleY := height/2 - 3
	g.renderer.DrawString(titleX, titleY, title)

	// Instructions
	instructions := []string{
		"Press SPACE or UP to jump",
		"Press Q to quit",
		"Press SPACE to start",
	}

	for i, instruction := range instructions {
		instrX := (width - len(instruction)) / 2
		instrY := titleY + 2 + i
		g.renderer.DrawString(instrX, instrY, instruction)
	}

	// High score
	highScore := fmt.Sprintf("High Score: %d", g.engine.GetHighScore())
	scoreX := (width - len(highScore)) / 2
	scoreY := height - 3
	g.renderer.DrawString(scoreX, scoreY, highScore)
}

// renderGame renders the main gameplay
func (g *Game) renderGame() {
	// Render ground line
	width, _ := g.renderer.GetSize()
	groundY := int(g.dinosaur.GroundLevel) + int(g.dinosaur.Height)
	for x := 0; x < width; x++ {
		g.renderer.DrawAt(x, groundY, '▔')
	}

	// Render dinosaur
	g.renderDinosaur()

	// Render obstacles
	g.renderObstacles()

	// Render UI
	g.renderUI()
}

// renderDinosaur renders the dinosaur sprite
func (g *Game) renderDinosaur() {
	art := g.dinosaur.GetASCIIArt()
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
			art := obstacle.GetASCIIArt()
			x := int(obstacle.X)
			y := int(obstacle.Y)

			for i, line := range art {
				g.renderer.DrawString(x, y+i, line)
			}
		}
	}
}

// renderUI renders the game UI (score, etc.)
func (g *Game) renderUI() {
	// Current score
	score := fmt.Sprintf("Score: %d", g.engine.GetCurrentScore())
	g.renderer.DrawString(2, 1, score)

	// High score
	highScore := fmt.Sprintf("High: %d", g.engine.GetHighScore())
	g.renderer.DrawString(2, 2, highScore)
}

// renderGameOver renders the game over screen
func (g *Game) renderGameOver() {
	width, height := g.renderer.GetSize()

	// Game Over title
	gameOver := "GAME OVER"
	gameOverX := (width - len(gameOver)) / 2
	gameOverY := height/2 - 2
	g.renderer.DrawString(gameOverX, gameOverY, gameOver)

	// Final score
	finalScore := fmt.Sprintf("Final Score: %d", g.engine.GetCurrentScore())
	scoreX := (width - len(finalScore)) / 2
	scoreY := gameOverY + 2
	g.renderer.DrawString(scoreX, scoreY, finalScore)

	// High score indicator
	if g.engine.IsNewHighScore() {
		newHigh := "NEW HIGH SCORE!"
		newHighX := (width - len(newHigh)) / 2
		newHighY := scoreY + 1
		g.renderer.DrawString(newHighX, newHighY, newHigh)
	}

	// Instructions
	restart := "Press R to restart or Q to quit"
	restartX := (width - len(restart)) / 2
	restartY := height/2 + 4
	g.renderer.DrawString(restartX, restartY, restart)
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
}

// restartGame restarts the game from game over state
func (g *Game) restartGame() {
	g.engine.Restart()
	g.spawner.Reset()
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
	// Create game instance
	game, err := NewGame()
	if err != nil {
		log.Fatalf("Failed to create game: %v", err)
	}
	defer game.Cleanup()

	// Run the game
	if err := game.Run(); err != nil {
		log.Fatalf("Game error: %v", err)
	}

	fmt.Println("Thanks for playing CLI Dino Game!")
}
