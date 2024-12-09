// Name:			Syed Muhammad Umar
//Student number: 	C00278724

package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	FISH = iota
	SHARK
)

const (
	NORTH = iota
	SOUTH
	EAST
	WEST
)

type coordinate struct {
	x, y int
}

var (
	sharkColor = color.RGBA{255, 0, 0, 255} // Red
	fishColor  = color.RGBA{0, 0, 0, 255}   // Black
	skyBlue    = color.RGBA{135, 206, 235, 255}
)

type Game struct {
	gridWidth    int
	gridHeight   int
	screenWidth  int
	screenHeight int
	fishesCount  int
	sharksCount  int
	fbreed       int
	sBreed       int
	starve       int
	routines     int
}

type creature struct {
	age, health, species int
	asset                color.RGBA
	chronon              int
}

var cellSize = 3
var tick = 0
var wm [][]*creature

func (g *Game) Update() error {
	tick++
	Chronon(tick, g)
	return nil
}
func Chronon(c int, g *Game) {
	for y := 0; y < g.gridHeight; y++ {
		for x := 0; x < g.gridWidth; x++ {
			if wm[x][y] == nil || wm[x][y].chronon == c {
				continue
			}

			wm[x][y].chronon = c // Mark creature as processed
			wm[x][y].age++       // Increment age

			if wm[x][y].species == FISH {
				if wm[x][y].age%g.fbreed == 0 { // Fish reproduces
					if spawnFish(x, y, c, g) { // Only increment count if reproduction succeeds
						g.fishesCount++
					}
				} else {
					moveCreature(wm[x][y], x, y, g, false)
				}
			} else if wm[x][y].species == SHARK {
				wm[x][y].health-- // Shark gets hungrier

				if wm[x][y].health <= 0 {
					wm[x][y] = nil  // Shark dies
					g.sharksCount-- // Decrease shark count
				} else if wm[x][y].age%g.sBreed == 0 { // Shark reproduces
					spawnShark(x, y, c, g)
					g.sharksCount++
				} else {
					if moveCreature(wm[x][y], x, y, g, true) {
						g.fishesCount-- // Only decrement fish count if a fish is eaten
					}
				}
			}
		}
	}
}

func moveCreature(c *creature, x, y int, g *Game, isShark bool) bool {
	start := coordinate{x, y}
	directions := shuffledDirections()

	// Sharks prioritize finding fish
	if isShark {
		n, s, e, w := adjacent(x, y, g)
		neighbors := []coordinate{n, s, e, w}

		for _, dir := range directions {
			target := neighbors[dir]
			if wm[target.x][target.y] != nil && wm[target.x][target.y].species == FISH {
				// Shark eats the fish
				wm[target.x][target.y] = c
				wm[x][y] = nil
				c.health = g.starve // Reset shark's health
				return true         // Indicates a fish was eaten
			}
		}
	}

	// Move to an empty space
	if target := findAvailableSpace(start, directions, g); target != nil {
		wm[target.x][target.y] = c
		wm[x][y] = nil
		return false // No fish was eaten, just movement
	}

	return false // No valid move found
}

func spawnFish(x, y, c int, g *Game) bool {
	start := coordinate{x, y}
	directions := shuffledDirections()

	if target := findAvailableSpace(start, directions, g); target != nil {
		wm[target.x][target.y] = &creature{
			age:     0,
			health:  0,
			species: FISH,
			asset:   fishColor,
			chronon: c,
		}
		return true
	}
	return false
}

func spawnShark(x, y, c int, g *Game) {
	start := coordinate{x, y}
	directions := shuffledDirections()

	if target := findAvailableSpace(start, directions, g); target != nil {
		wm[target.x][target.y] = &creature{
			age:     0,
			health:  g.starve,
			species: SHARK,
			asset:   sharkColor,
			chronon: c,
		}
	}
}

func shuffledDirections() []int {
	directions := []int{NORTH, SOUTH, EAST, WEST}
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
	return directions
}

func adjacent(x, y int, g *Game) (coordinate, coordinate, coordinate, coordinate) {
	height, width := g.gridHeight, g.gridWidth

	n := coordinate{x, (y + height - 1) % height}
	s := coordinate{x, (y + 1) % height}
	e := coordinate{(x + 1) % width, y}
	w := coordinate{(x + width - 1) % width, y}

	return n, s, e, w
}
func findAvailableSpace(start coordinate, directions []int, g *Game) *coordinate {
	n, s, e, w := adjacent(start.x, start.y, g)
	neighbors := []coordinate{n, s, e, w}

	for _, dir := range directions {
		target := neighbors[dir]
		if wm[target.x][target.y] == nil { // Check if space is empty
			return &target
		}
	}
	return nil
}

func initWator(game *Game) {
	wm = make([][]*creature, game.gridWidth)
	for i := range wm {
		wm[i] = make([]*creature, game.gridHeight)
	}
	pop := 0

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < game.fishesCount; i++ {
		for {
			if pop == game.gridHeight*game.gridWidth {
				break
			}
			x := r.Intn(game.gridWidth)
			y := r.Intn(game.gridHeight)
			if wm[x][y] == nil {
				wm[x][y] = &creature{
					age:     r.Intn(game.fbreed),
					species: FISH,
					asset:   fishColor,
				}
				pop++
				break
			}
		}
	}

	for i := 0; i < game.sharksCount; i++ {
		for {
			if pop == game.gridHeight*game.gridWidth {
				break
			}
			x := r.Intn(game.gridWidth)
			y := r.Intn(game.gridHeight)
			if wm[x][y] == nil {
				wm[x][y] = &creature{
					age:     r.Intn(game.sBreed),
					species: SHARK,
					health:  game.starve,
					asset:   sharkColor,
				}
				pop++
				break
			}
		}
	}
}
func (g *Game) Draw(screen *ebiten.Image) {
	// Set background
	screen.Fill(skyBlue)

	// Draw creatures
	for x := 0; x < g.gridWidth; x++ {
		for y := 0; y < g.gridHeight; y++ {
			if wm[x][y] != nil {
				xPos := x * cellSize
				yPos := y * cellSize
				ebitenutil.DrawRect(screen, float64(xPos), float64(yPos), float64(cellSize), float64(cellSize), wm[x][y].asset)
			}
		}
	}

	// Display simulation stats as text
	info := "Tick: " + strconv.Itoa(tick) + "\n" +
		"Sharks: " + strconv.Itoa(g.sharksCount) + "\n" +
		"Fish: " + strconv.Itoa(g.fishesCount)
	ebitenutil.DebugPrint(screen, info)

	// Draw the graph panel
	panelWidth := 200
	panelHeight := g.screenHeight
	panelOffsetX := g.screenWidth - panelWidth
	ebitenutil.DrawRect(screen, float64(panelOffsetX), 0, float64(panelWidth), float64(panelHeight), color.RGBA{R: 200, G: 200, B: 200, A: 255})

	// Draw the bar graph within the panel
	barWidth := 30
	barMaxHeight := panelHeight - 100
	graphOffsetX := panelOffsetX + 50
	graphOffsetY := 50

	// Normalize the heights of the bars
	totalCreatures := g.fishesCount + g.sharksCount
	if totalCreatures == 0 {
		totalCreatures = 1 // Prevent division by zero
	}

	// Calculate the heights of the bars
	fishBarHeight := int(float64(g.fishesCount) / float64(totalCreatures) * float64(barMaxHeight))
	sharkBarHeight := int(float64(g.sharksCount) / float64(totalCreatures) * float64(barMaxHeight))

	// Draw the bars
	ebitenutil.DrawRect(screen, float64(graphOffsetX), float64(graphOffsetY+barMaxHeight-fishBarHeight), float64(barWidth), float64(fishBarHeight), fishColor)
	ebitenutil.DrawRect(screen, float64(graphOffsetX+barWidth+10), float64(graphOffsetY+barMaxHeight-sharkBarHeight), float64(barWidth), float64(sharkBarHeight), sharkColor)

	// Add labels for the bars
	ebitenutil.DebugPrintAt(screen, "Fish", graphOffsetX, graphOffsetY+barMaxHeight+10)
	ebitenutil.DebugPrintAt(screen, "Sharks", graphOffsetX+barWidth+10, graphOffsetY+barMaxHeight+10)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

func main() {
	var gridWidth, gridHeight, sharksCount, fishesCount, routines int

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
	fmt.Scan(&routines)

	// Calculate screen dimensions
	screenWidth := gridWidth * cellSize
	screenHeight := gridHeight * cellSize

	// Initialize game
	game := &Game{
		gridWidth:    gridWidth,
		gridHeight:   gridHeight,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		fishesCount:  fishesCount,
		sharksCount:  sharksCount,
		fbreed:       150,
		sBreed:       150,
		starve:       150,
		routines:     routines,
	}

	initWator(game)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Wa-Tor Simulation (Multi-threaded)")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
