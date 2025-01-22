/*
	export Game object
	call game.Start()

	game object:
		- 2D matrix: board
		- 2D matrix: buffer
		- rules: array of two items
			- where every item is an array of 9 bools

	algorithm description:
	- init
		- fill with random values
	- run
		- iterate through every cell
			- count cell neighbors
			- if rules[cell][neighbors] true
				- buffer the same position true
			- else
				- buffer the same position false
		- swap buffer and board pointers
		- draw
			-> import from graphics.go and
			-> give it the board object
*/

package gameoflife

import (
	"math/rand"
	"time"
)

type Game struct {
	rows          int
	cols          int
	board         [][]int
	buffer        [][]int
	rules         [2][9]int
	cycleDuration time.Duration
	drawingFunc   func([][]int)
}

func NewGame(rows int, cols int, durationMsecs int) *Game {
	board := make([][]int, rows)
	for i := 0; i < rows; i++ {
		board[i] = make([]int, cols)
	}

	buffer := make([][]int, rows)
	for i := 0; i < rows; i++ {
		buffer[i] = make([]int, cols)
	}

	game := &Game{
		rows:   rows,
		cols:   cols,
		board:  board,
		buffer: buffer,
		rules: [2][9]int{
			{0, 0, 0, 1, 0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0, 0, 0, 0, 0},
		},
		cycleDuration: time.Millisecond * time.Duration(durationMsecs),
	}

	return game
}

func (g *Game) fillWithRandom() {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			g.board[i][j] = rand.Intn(2)
		}
	}
}

func (g *Game) setRule(cell int, neighbors int, newValue int) {
	if !(cell == 0 || cell == 1) {
		panic("trying to set a rule for invalid cell state")
	}

	if !(neighbors >= 0 && neighbors < 9) {
		panic("trying to set a rule for an invalid count of neighbors")
	}

	if !(newValue == 0 || newValue == 1) {
		panic("trying to set a rule to an invalid value")
	}

	g.rules[cell][neighbors] = newValue
}

// dr stands for row delta
// dc stands for column delta
// nb stands for neighbor
func (g *Game) countNeighbors(row int, col int) int {
	count := 0

	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			nbRow := row + dr
			nbCol := col + dc

			// row out of game board
			if nbRow < 0 || nbRow >= g.rows {
				continue
			}

			// column out of game board
			if nbCol < 0 || nbCol >= g.cols {
				continue
			}

			// prevent from counting itself
			if dr == 0 && dc == 0 {
				continue
			}

			count += g.board[nbRow][nbCol]
		}
	}

	return count
}

func (g *Game) getNewState(row int, col int, neighbors int) int {
	return g.rules[g.board[row][col]][neighbors]
}

func (g *Game) generateNewGeneration() {
	for i := 0; i < g.rows; i++ {
		for j := 0; j < g.cols; j++ {
			g.buffer[i][j] = g.getNewState(i, j, g.countNeighbors(i, j))
		}
	}
}

func (g *Game) swapMatrices() {
	temp := g.board
	g.board = g.buffer
	g.buffer = temp
}

func (g *Game) wait(startTime time.Time) {
	elapsedTime := time.Since(startTime)
	remainingTime := g.cycleDuration - elapsedTime

	time.Sleep(remainingTime)
}

func (g *Game) SetDrawingFunc(drawingFunc func([][]int)) {
	g.drawingFunc = drawingFunc
}

func (g *Game) Start() {
	g.fillWithRandom()

	for {
		startTime := time.Now()
		g.generateNewGeneration()
		g.swapMatrices()
		g.drawingFunc(g.board)
		g.wait(startTime)
	}
}
