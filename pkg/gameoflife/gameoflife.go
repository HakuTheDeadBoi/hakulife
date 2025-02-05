package gameoflife

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	CLOSED             = 0
	RUNNING            = 1
	PAUSED             = 2
	MSECONDSINPUTDELAY = 10
)

type Game struct {
	rows          int
	cols          int
	board         [][]int
	buffer        [][]int
	rules         [2][9]int
	cycleDuration time.Duration
	drawingFunc   func([][]int, int)
	state         int
	keyChan       chan termbox.Event
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

	keyChan := make(chan termbox.Event, 50)

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
		state:         RUNNING,
		keyChan:       keyChan,
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

	if remainingTime > time.Duration(0) {
		time.Sleep(remainingTime)
	}
}

func (g *Game) readInput() {
	termbox.SetInputMode(termbox.InputEsc)
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			switch ev.Ch {
			case 'Q':
				g.state = CLOSED
			case 'q':
				g.state = CLOSED
			case 'R':
				g.state = PAUSED
				g.fillWithRandom()
				g.state = RUNNING
			case 'r':
				g.state = PAUSED
				g.fillWithRandom()
				g.state = RUNNING
			case 'P':
				if g.state == PAUSED {
					g.state = RUNNING
				} else {
					g.state = PAUSED
				}
			case 'p':
				if g.state == PAUSED {
					g.state = RUNNING
				} else {
					g.state = PAUSED
				}
			}
		}
		time.Sleep(MSECONDSINPUTDELAY * time.Millisecond)
	}
}

func (g *Game) SetDrawingFunc(drawingFunc func([][]int, int)) {
	g.drawingFunc = drawingFunc
}

func (g *Game) Start() {
	g.fillWithRandom()
	go g.readInput()

	for {
		startTime := time.Now()
		if g.state == CLOSED {
			return
		}

		if g.state == RUNNING {
			g.generateNewGeneration()
			g.swapMatrices()
		}

		g.drawingFunc(g.board, g.state)
		g.wait(startTime)
	}
}
