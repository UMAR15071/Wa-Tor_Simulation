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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var xcoord, ycoord int

	for y := 0; y < g.gridHeight; y++ {
		for x := 0; x < g.gridWidth; x++ {
			north, south, east, west := adjacent(x, y, g)

			if wm[x][y] == nil {
				continue
			}

			if wm[x][y].chronon == c {
				continue
			}
			wm[x][y].age += 1
			d := r.Intn(3)
			switch wm[x][y].species {
			case FISH:
				foundspace := false
				wm[x][y].chronon = c
				for i := 0; i < 4; i++ {
					d += i

					switch d % 4 {
					case NORTH:
						xcoord, ycoord = north.x, north.y
					case SOUTH:
						xcoord, ycoord = south.x, south.y
					case EAST:
						xcoord, ycoord = east.x, east.y
					case WEST:
						xcoord, ycoord = west.x, west.y
					}

					// Check bounds
					if wm[xcoord][ycoord] == nil {
						wm[xcoord][ycoord] = wm[x][y]
						wm[xcoord][ycoord].age = 0 // New fish
						foundspace = true
						break
					}
				}
				if !foundspace {
					wm[x][y] = nil
				}

			case SHARK:
				foundfish := false
				wm[x][y].chronon = c

				// Sharks get hungrier each turn
				wm[x][y].health--
				if wm[x][y].health <= 0 {
					wm[x][y] = nil // Shark starves to death
					break
				}

				for i := 0; i < 4; i++ {
					d += i
					switch d % 4 {
					case NORTH:
						xcoord, ycoord = north.x, north.y
					case SOUTH:
						xcoord, ycoord = south.x, south.y
					case EAST:
						xcoord, ycoord = east.x, east.y
					case WEST:
						xcoord, ycoord = west.x, west.y
					}

					// Check bounds
					if wm[xcoord][ycoord] == nil {
						continue
					}

					// Found fish to eat
					if wm[xcoord][ycoord].species == FISH {
						foundfish = true
						wm[xcoord][ycoord] = wm[x][y]
						wm[xcoord][ycoord].health = g.starve
						wm[x][y] = nil
						break
					}
				}

				if !foundfish {
					// No fish found, move to empty space
					for i := 0; i < 4; i++ {
						d += i
						switch d % 4 {
						case NORTH:
							xcoord, ycoord = north.x, north.y
						case SOUTH:
							xcoord, ycoord = south.x, south.y
						case EAST:
							xcoord, ycoord = east.x, east.y
						case WEST:
							xcoord, ycoord = west.x, west.y
						}

						if wm[xcoord][ycoord] == nil {
							wm[xcoord][ycoord] = wm[x][y]
							wm[x][y] = nil

							// Check if shark can reproduce
							if wm[xcoord][ycoord].age > 0 && wm[xcoord][ycoord].age%g.sBreed == 0 {
								wm[x][y] = &creature{
									age:     0,
									health:  g.starve,
									species: SHARK,
									asset:   sharkColor,
									chronon: c,
								}
							}
							break
						}
					}
				}
			}
		}
	}
}

func adjacent(x, y int, g *Game) (coordinate, coordinate, coordinate, coordinate) {

	var n, s, e, w coordinate
	if y == 0 {
		n.y = *&g.gridHeight - 1
	} else {
		n.y = y - 1
	}
	n.x = x
	if y == *&g.gridHeight-1 {
		s.y = 0
	} else {
		s.y = y + 1
	}
	s.x = x
	if x == *&g.gridWidth-1 {
		e.x = 0
	} else {
		e.x = x + 1
	}
	e.y = y
	if x == 0 {
		w.x = *&g.gridWidth - 1
	} else {
		w.x = x - 1
	}
	w.y = y

	return n, s, e, w
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

	// Debug display
	ebitenutil.DebugPrint(screen, "Tick: "+strconv.Itoa(tick))
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
		fbreed:       50,
		sBreed:       200,
		starve:       100,
		routines:     routines,
	}

	initWator(game)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Wa-Tor Simulation (Multi-threaded)")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
