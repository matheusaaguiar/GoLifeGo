package main

import (
	"fmt"
	"time"

	"encoding/json"
	"image/color"
	"math/rand/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// === Grid ===
type Grid struct {
	Lines int
	Cols int
	cells [] uint8
}

func (lifeGrid *Grid) at(i, j int) uint8 {
	return lifeGrid.cells[i * lifeGrid.Lines + j]
}

func (lifeGrid *Grid) alive(i, j int) {
	lifeGrid.cells[i * lifeGrid.Lines + j] = 1
}

func (lifeGrid *Grid) dead(i, j int) {
	lifeGrid.cells[i * lifeGrid.Lines + j] = 0
}

func (lifeGrid *Grid) isAlive(i, j int) bool {
	if (0 <= i && i <= lifeGrid.Lines && 0 <= j && j <= lifeGrid.Cols) {
		return lifeGrid.at(i, j) == 1
	}
	return false
}

func (lifeGrid *Grid) gridCenter() (int, int) {
	return lifeGrid.Lines / 2, lifeGrid.Cols / 2
}

func (lifeGrid *Grid) countNeighbours(i, j int) int {
	var neighbours int
	for dx := -1; dx < 2; dx++ {
		for dy := -1; dy < 2; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			x := i + dx
			y := j + dy
			if x >= 0 && x < lifeGrid.Lines && y >= 0 && y < lifeGrid.Cols {
				neighbours += int(lifeGrid.at(x, y))
			}
		}
	}
	return neighbours
}

func NewGrid(lines, cols int) *Grid {
	newGrid := &Grid{
		Lines: lines,
		Cols: cols,
		cells: make([] uint8, lines*cols),
	}
	return newGrid
}

func updateGrid(lifeGrid *Grid) *Grid {
	nextGrid := NewGrid(lifeGrid.Lines, lifeGrid.Cols)

	for i := range lifeGrid.Lines {
		for j := range lifeGrid.Cols {
			neighbours := lifeGrid.countNeighbours(i, j)
			if lifeGrid.isAlive(i, j) && 2 <= neighbours && neighbours <= 3 {
				nextGrid.alive(i, j)
			} else if !lifeGrid.isAlive(i, j) && neighbours == 3 {
				nextGrid.alive(i, j)
			}
		}
	}
	return nextGrid
}

func VerticalLinePattern(lifeGrid *Grid, i, j, size int) {
	lifeGrid.alive(i, j)

	counter := 1
	pos := 1
	for {
		if i - pos >= 0 {
			lifeGrid.alive(i - pos, j)
			counter++
		}
		if counter >= size {
			break
		}
		if i + pos < lifeGrid.Lines {
			lifeGrid.alive(i + pos, j)
			counter++
		}
		if counter >= size {
			break
		}
		pos++
	}
}

func blockPattern(lifeGrid *Grid, i, j int) {
	lifeGrid.alive(i, j)
	if i + 1 >= 0 {
		lifeGrid.alive(i + 1, j)
	}
	if j + 1 >= 0 {
		lifeGrid.alive(i, j + 1)
		if i + 1 >= 0 {
			lifeGrid.alive(i + 1, j + 1)
		}
	}
}

func gliderPattern(lifeGrid *Grid, i, j int) {
	glider := [3][3]int {
		{0, 1, 0},
		{0, 0, 1},
		{1, 1, 1},
	}
	for di := -1; di <= 1; di++ {
		for dj := -1; dj <= 1; dj++ {
			x := i + di
			y := j + dj
			if x >= 0 && x < lifeGrid.Lines && y >= 0 && y < lifeGrid.Cols {
				if glider[1 + di][1 + dj] == 1 {
					lifeGrid.alive(x, y)
				}
			}
		}
	}

}



// === Grid Input/Output ===

type jsonGrid struct {
	Lines int `json:"lines"`
	Cols int `json:"cols"`
	Cells []uint8 `json:"cells"`
}

func (grid Grid) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonGrid{
		Lines: grid.Lines,
		Cols: grid.Cols,
		Cells: grid.cells,
	})
}

func (grid *Grid) UnmarshalJSON(data []byte) error {
	var aux jsonGrid
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	grid.Lines = aux.Lines
	grid.Cols = aux.Cols
	grid.cells = aux.Cells
	return nil
}

func printGrid(lifeGrid *Grid) {
	// Clear screen and move cursor to top-left
	//fmt.Print("\033[2J\033[H")
	fmt.Print("\n\n\n\n\n")
	for i := range lifeGrid.Lines {
		for j := range lifeGrid.Cols {
			fmt.Print(lifeGrid.at(i, j))
		}
		fmt.Println()
	}
}

// === GUI ===
// TODO
type GuiGrid struct {

}

func newGuiGrid(width int) *fyne.Container {
	grid := container.NewGridWithColumns(width)
	for i := 0; i < width * width; i++ {
		square := canvas.NewRectangle(color.White)
		square.SetMinSize(fyne.NewSize(20, 20))
		// Add a thin border so white squares are visible against the background
		square.StrokeColor = color.Gray{Y: 128}
		square.StrokeWidth = 1

		grid.Add(square)
	}
	return grid
}

// randomizeGrid randomly sets each rectangle in the grid to black or white.
func randomizeGuiGrid(grid *fyne.Container) {
	for _, cell := range grid.Objects {
		if square, ok := cell.(*canvas.Rectangle); ok {
			if rand.IntN(2) == 1 {
				square.FillColor = color.Black
			} else {
				square.FillColor = color.White
			}
		}
	}
	grid.Refresh()
}

func updateGui(guiGrid *fyne.Container, lifeGrid *Grid) {
	n := 0
	for _, cell := range guiGrid.Objects {
		if square, ok := cell.(*canvas.Rectangle); ok {
			i, j := n / lifeGrid.Lines, n % lifeGrid.Cols
			n++
			if lifeGrid.isAlive(i, j) {
				square.FillColor = color.Black
			} else {
				square.FillColor = color.White
			}
		}
	}
	guiGrid.Refresh()
}

func main() {

	maxIterations := 100
	width := 19
	height := 19

	lifeGrid := NewGrid(width, height)
	ci, cj := lifeGrid.gridCenter()
	VerticalLinePattern(lifeGrid, ci, cj, 16)
	// blockPattern(lifeGrid, ci, cj)
	//gliderPattern(lifeGrid, ci, cj)

	// Setup GUI
	mainApp := app.New()
	mainWindow := mainApp.NewWindow("Game of Life")
	// TODO: refactor to struct
	guiGrid := newGuiGrid(width)
	updateGui(guiGrid, lifeGrid)

	// Header label above the grid
	label := widget.NewLabelWithStyle("Game of Life", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	iterations := 0
	counterLabel := widget.NewLabel(fmt.Sprintf("Iterations: %d", iterations))
	// Theme-aware divider between header and grid
	separator := widget.NewSeparator()
	// Stack the centered label and the separator vertically
	header := container.NewVBox(
		container.NewCenter(label),
		container.NewCenter(counterLabel),
		separator,
	)
	// Place header on top, grid fills the rest
	content := container.NewBorder(header, nil, nil, nil, guiGrid)

	mainWindow.SetContent(content)
	mainWindow.Resize(fyne.NewSize(float32(width*20 + width*5), float32(height*20 + height*5)))

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			fyne.DoAndWait(func() {
				lifeGrid = updateGrid(lifeGrid)
				// printGrid(lifeGrid)  // CLI
				updateGui(guiGrid, lifeGrid)
				iterations++
				counterLabel.SetText(fmt.Sprintf("Iterations: %d", iterations))
				if iterations >= maxIterations {
					mainApp.Quit()
				}
			})
		}
	}()

	mainWindow.ShowAndRun()
}
