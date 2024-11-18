package main

import (
	"fmt"
	"image/color"
	"log"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/exp/rand"
)

// Shark and Fish colors
var (
	sharkColor = color.RGBA{255, 0, 0, 255} // Red
	fishColor  = color.RGBA{0, 0, 0, 255}   // Black
)

// Cell size is constant
const cellSize = 7

// Add a mutex
var gridMutex sync.Mutex

const chrononDuration = 0.1

// Position represents the x, y coordinates of an entity
type Position struct {
	x int
	y int
}

// Directions for movement
var directions = []Position{
	{x: 0, y: -1}, // Up
	{x: 0, y: 1},  // Down
	{x: -1, y: 0}, // Left
	{x: 1, y: 0},  // Right
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

// Update the game state, updating sharks and fishes
func (g *Game) Update() error {
	now := time.Now()
	if now.Sub(g.lastUpdate).Seconds() >= chrononDuration {
		// Update sharks and fishes using the new logic
		g.sharks = updateEntitiesParallel(g.sharks, g.fishes, g.gridWidth, g.gridHeight, g.routinesCount, "Shark")
		g.fishes = updateEntitiesParallel(g.fishes, g.sharks, g.gridWidth, g.gridHeight, g.routinesCount, "Fish")

		// Reset the last update time
		g.lastUpdate = now
	}
	return nil
}

// updateEntitiesParallel updates entities in parallel using the specified number of routines
func updateEntitiesParallel(entities, otherEntities []Position, gridWidth, gridHeight, routinesCount int, entityType string) []Position {
	chunkSize := (len(entities) + routinesCount - 1) / routinesCount
	results := make(chan []Position, routinesCount)

	// Create a grid to track occupied cells
	grid := make(map[Position]string)

	// Populate the grid with all current entities
	for _, e := range entities {
		grid[e] = entityType
	}
	for _, e := range otherEntities {
		if entityType == "Shark" {
			grid[e] = "Fish"
		} else {
			grid[e] = "Shark"
		}
	}

	var wg sync.WaitGroup

	// Split the entities into chunks and update in parallel
	for i := 0; i < routinesCount; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(entities) {
			end = len(entities)
		}

		wg.Add(1)
		go func(subEntities []Position) {
			defer wg.Done()
			results <- moveEntities(subEntities, grid, gridWidth, gridHeight)
		}(entities[start:end])
	}

	wg.Wait()
	close(results)

	// Gather updated entities from all results
	updatedEntities := make([]Position, 0, len(entities))
	for res := range results {
		updatedEntities = append(updatedEntities, res...)
	}

	return updatedEntities
}

func moveEntities(entities []Position, grid map[Position]string, gridWidth, gridHeight int) []Position {
	// New slice to store entities after movement
	var newEntities []Position

	// Iterate over each entity
	for i, pos := range entities {
		// Shuffle directions for random movement
		rand.Shuffle(len(directions), func(i, j int) { directions[i], directions[j] = directions[j], directions[i] })

		// Try moving in shuffled directions
		for _, dir := range directions {
			// Calculate new position with wrap-around logic
			newX := (pos.x + dir.x + gridWidth) % gridWidth   // Wrap around horizontally
			newY := (pos.y + dir.y + gridHeight) % gridHeight // Wrap around vertically
			newPos := Position{x: newX, y: newY}

			// Lock grid access to ensure thread-safety
			gridMutex.Lock()

			// Assume the cell is empty until proven otherwise
			isEmpty := true

			// Check if the new position is occupied by a shark or fish
			if grid[newPos] == "Shark" || grid[newPos] == "Fish" {
				isEmpty = false
			}

			// If the new position is empty, move the entity
			if isEmpty {
				// Mark current position as empty
				delete(grid, pos)

				// Mark new position as occupied (depending on the entity type)
				if grid[pos] == "Shark" {
					grid[newPos] = "Shark"
				} else if grid[pos] == "Fish" {
					grid[newPos] = "Fish"
				}

				// Reproduce: Add the previous position as a new entity if needed
				newEntities = append(newEntities, pos)

				// Move the entity
				entities[i] = newPos
			}

			// Unlock grid after updating
			gridMutex.Unlock()

			// If the new position is empty, break out of the loop as no further movement needed
			if isEmpty {
				break
			}
		}
	}

	// Return the updated list of entities
	return entities
}

// Method to generate initial position of sharks and fishes
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

// Draw method is necessary because Ebiten interface requires it to be implemented.
// Draw draws the game screen.
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

func main() {
	var gridWidth, gridHeight, sharksCount, fishesCount, routinesCount int

	// Input grid dimensions and entity counts
	fmt.Print("Enter grid width: ")
	fmt.Scan(&gridWidth)
	fmt.Print("Enter grid height: ")
	fmt.Scan(&gridHeight)
	fmt.Print("Enter number of sharks: ")
	fmt.Scan(&sharksCount)
	fmt.Print("Enter number of fishes: ")
	fmt.Scan(&fishesCount)
	fmt.Print("Enter number of goroutines: ")
	fmt.Scan(&routinesCount)

	// Calculate screen dimensions
	screenWidth := gridWidth * cellSize
	screenHeight := gridHeight * cellSize

	if gridWidth <= 0 || gridHeight <= 0 || sharksCount < 0 || fishesCount < 0 || routinesCount <= 0 {
		log.Fatal("Invalid input: all values must be positive integers")
	}

	// Create a new game
	game := &Game{
		gridWidth:     gridWidth,
		gridHeight:    gridHeight,
		screenWidth:   screenWidth,
		screenHeight:  screenHeight,
		sharks:        generatePositions(sharksCount, gridWidth, gridHeight),
		fishes:        generatePositions(fishesCount, gridWidth, gridHeight),
		routinesCount: routinesCount,
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Wa-Tor Simulation (Multi-threaded)")

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
