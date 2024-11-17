// Name:			Syed Muhammad Umar
//Student Number:	C00278724

package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// Shark and Fish colors
var (
	sharkColor = color.RGBA{255, 0, 0, 255} // Red
	fishColor  = color.RGBA{0, 0, 0, 255}   // Black
)

// Game implements ebiten.Game interface.
type Game struct {
	screenWidth  int
	screenHeight int
	sharks       int
	fishes       int
	sharkPos     []Position // Positions of sharks
	fishPos      []Position // Positions of fishes
}

// Position represents the x, y coordinates of an entity
type Position struct {
	x int
	y int
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Update logic can be added here (e.g., move sharks and fishes)
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	// Set the background color to sky blue
	skyBlue := color.RGBA{135, 206, 235, 255}
	screen.Fill(skyBlue)

	// Draw sharks
	for _, pos := range g.sharkPos {
		screen.Set(pos.x, pos.y, sharkColor)
	}

	// Draw fishes
	for _, pos := range g.fishPos {
		screen.Set(pos.x, pos.y, fishColor)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// Return fixed screen dimensions, ensuring they are positive.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

// generatePositions generates random positions for entities
func generatePositions(count, width, height int) []Position {
	positions := make([]Position, count)
	for i := 0; i < count; i++ {
		positions[i] = Position{
			x: rand.Intn(width),  // Random X position within screen width
			y: rand.Intn(height), // Random Y position within screen height
		}
	}
	return positions
}

func main() {
	var screenWidth, screenHeight, sharks, fishes int

	// Input dimensions and entity counts
	fmt.Println("Please enter screen dimensions")
	fmt.Print("Width:  ")
	fmt.Scan(&screenWidth)
	fmt.Print("Height: ")
	fmt.Scan(&screenHeight)

	fmt.Println("Please enter number of sharks and fishes")
	fmt.Print("Sharks: ")
	fmt.Scan(&sharks)
	fmt.Print("Fishes: ")
	fmt.Scan(&fishes)

	if screenWidth <= 0 || screenHeight <= 0 || sharks < 0 || fishes < 0 {
		log.Fatal("Invalid input: screen dimensions and entity counts must be positive")
	}

	// Initialize game
	game := &Game{
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		sharks:       sharks,
		fishes:       fishes,
		sharkPos:     generatePositions(sharks, screenWidth, screenHeight),
		fishPos:      generatePositions(fishes, screenWidth, screenHeight),
	}

	// Set the window size to match the screen dimensions
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Wa-Tor Simulation")

	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
