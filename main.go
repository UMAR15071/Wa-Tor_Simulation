package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Shark and Fish colors
var (
	sharkColor = color.RGBA{255, 0, 0, 255} // Red
	fishColor  = color.RGBA{0, 0, 0, 255}   // Black
)

// Cell size is constant
const cellSize = 7 // Each grid cell is 7x7 pixels

// Game implements ebiten.Game interface.
type Game struct {
	gridWidth    int
	gridHeight   int
	screenWidth  int
	screenHeight int
	sharks       []Position
	fishes       []Position
}

// Position represents the x, y coordinates of an entity
type Position struct {
	x int
	y int
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Future movement logic can go here
	return nil
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
	var gridWidth, gridHeight, sharksCount, fishesCount int

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

	// Calculate screen dimensions
	screenWidth := gridWidth * cellSize
	screenHeight := gridHeight * cellSize

	if gridWidth <= 0 || gridHeight <= 0 || sharksCount < 0 || fishesCount < 0 {
		log.Fatal("Invalid input: all values must be positive integers")
	}

	// Initialize game
	game := &Game{
		gridWidth:    gridWidth,
		gridHeight:   gridHeight,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		sharks:       generatePositions(sharksCount, gridWidth, gridHeight),
		fishes:       generatePositions(fishesCount, gridWidth, gridHeight),
	}

	// Set the window size to match the screen dimensions
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Wa-Tor Simulation (Grid-Based)")

	// Start the game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
