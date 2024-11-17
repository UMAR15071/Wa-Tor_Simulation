package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Shark and Fish colors
var (
	sharkColor = color.RGBA{255, 0, 0, 255} // Red
	fishColor  = color.RGBA{0, 0, 0, 255}   // Black
)

// Cell size is constant
const cellSize = 7

// Chronon duration (in seconds)
const chrononDuration = 0.1

// Directions for movement
var directions = []Position{
	{x: 0, y: -1}, // Up
	{x: 0, y: 1},  // Down
	{x: -1, y: 0}, // Left
	{x: 1, y: 0},  // Right
}

// Position represents the x, y coordinates of an entity
type Position struct {
	x int
	y int
}

// Game implements ebiten.Game interface.
type Game struct {
	gridWidth     int
	gridHeight    int
	screenWidth   int
	screenHeight  int
	sharks        []Position
	fishes        []Position
	lastUpdate    time.Time
	routinesCount int
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	now := time.Now()
	if now.Sub(g.lastUpdate).Seconds() >= chrononDuration {
		// Update sharks and fishes using multiple goroutines
		g.sharks = updateEntitiesParallel(g.sharks, g.gridWidth, g.gridHeight, g.routinesCount)
		g.fishes = updateEntitiesParallel(g.fishes, g.gridWidth, g.gridHeight, g.routinesCount)

		// Reset the last update time
		g.lastUpdate = now
	}
	return nil
}

// updateEntitiesParallel updates entities in parallel using the specified number of routines
func updateEntitiesParallel(entities []Position, gridWidth, gridHeight, routinesCount int) []Position {
	chunkSize := (len(entities) + routinesCount - 1) / routinesCount // Divide work into chunks
	results := make(chan []Position, routinesCount)

	var wg sync.WaitGroup

	for i := 0; i < routinesCount; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(entities) {
			end = len(entities)
		}

		wg.Add(1)
		go func(subEntities []Position) {
			defer wg.Done()
			results <- moveEntities(subEntities, gridWidth, gridHeight)
		}(entities[start:end])
	}

	wg.Wait()
	close(results)

	// Collect results from all routines
	updatedEntities := make([]Position, 0, len(entities))
	for res := range results {
		updatedEntities = append(updatedEntities, res...)
	}

	return updatedEntities
}

// moveEntities moves each entity randomly within the grid
func moveEntities(entities []Position, gridWidth, gridHeight int) []Position {
	for i, pos := range entities {
		// Choose a random direction
		dir := directions[rand.Intn(len(directions))]

		// Calculate new position
		newX := pos.x + dir.x
		newY := pos.y + dir.y

		// Check boundaries and update position if valid
		if newX >= 0 && newX < gridWidth && newY >= 0 && newY < gridHeight {
			entities[i] = Position{x: newX, y: newY}
		}
	}
	return entities
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Set the background color to sky blue
	skyBlue := color.RGBA{135, 206, 235, 255}
	screen.Fill(skyBlue)

	// Draw grid lines
	for i := 0; i <= g.gridWidth; i++ {
		x := i * cellSize
		ebitenutil.DrawLine(screen, float64(x), 0, float64(x), float64(g.screenHeight), color.Black)
	}
	for i := 0; i <= g.gridHeight; i++ {
		y := i * cellSize
		ebitenutil.DrawLine(screen, 0, float64(y), float64(g.screenWidth), float64(y), color.Black)
	}

	// Draw sharks as red filled cells
	for _, pos := range g.sharks {
		ebitenutil.DrawRect(screen, float64(pos.x*cellSize), float64(pos.y*cellSize), float64(cellSize), float64(cellSize), sharkColor)
	}

	// Draw fishes as black filled cells
	for _, pos := range g.fishes {
		ebitenutil.DrawRect(screen, float64(pos.x*cellSize), float64(pos.y*cellSize), float64(cellSize), float64(cellSize), fishColor)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

// generatePositions generates random positions for entities within the grid
func generatePositions(count, gridWidth, gridHeight int) []Position {
	positions := make([]Position, count)
	for i := 0; i < count; i++ {
		positions[i] = Position{
			x: rand.Intn(gridWidth),  // Random X position within grid width
			y: rand.Intn(gridHeight), // Random Y position within grid height
		}
	}
	return positions
}

func main() {
	var gridWidth, gridHeight, sharksCount, fishesCount, routinesCount int

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Input grid dimensions and entity counts
	fmt.Println("Enter the number of grid cells (width and height):")
	fmt.Print("Grid Width:  ")
	fmt.Scan(&gridWidth)
	fmt.Print("Grid Height: ")
	fmt.Scan(&gridHeight)

	fmt.Println("Enter the number of sharks and fishes:")
	fmt.Print("Sharks: ")
	fmt.Scan(&sharksCount)
	fmt.Print("Fishes: ")
	fmt.Scan(&fishesCount)

	fmt.Println("Enter the number of routines:")
	fmt.Print("Routines: ")
	fmt.Scan(&routinesCount)

	// Calculate screen dimensions
	screenWidth := gridWidth * cellSize
	screenHeight := gridHeight * cellSize

	if gridWidth <= 0 || gridHeight <= 0 || sharksCount < 0 || fishesCount < 0 || routinesCount <= 0 {
		log.Fatal("Invalid input: all values must be positive integers")
	}

	// Initialize game
	game := &Game{
		gridWidth:     gridWidth,
		gridHeight:    gridHeight,
		screenWidth:   screenWidth,
		screenHeight:  screenHeight,
		sharks:        generatePositions(sharksCount, gridWidth, gridHeight),
		fishes:        generatePositions(fishesCount, gridWidth, gridHeight),
		lastUpdate:    time.Now(),
		routinesCount: routinesCount,
	}

	// Set the window size to match the screen dimensions
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Wa-Tor Simulation (Multi-threaded)")

	// Start the game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
