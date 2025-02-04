package graphics

import "github.com/nsf/termbox-go"

const (
	LABEL    = "Hakulife"
	TOPLEFT  = 0x250F
	TOPRIGHT = 0x2513
	BOTLEFT  = 0x2517
	BOTRIGHT = 0x251B
	HLINE    = 0x2501
	VLINE    = 0x2503
	ALIVE    = 0x2588
)

var stateLabels [3]string = [3]string{"CLOSED", "RUNNING", "PAUSED"}

func Init() {
	termbox.Init()
}

func Close() {
	termbox.Close()
}

func update() {
	termbox.Flush()
}

func GetScreenSize() (int, int) {
	return termbox.Size()
}

func renderFrame(gameState int) {
	scrW, scrH := termbox.Size()
	termbox.SetCell(0, 0, TOPLEFT, termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(scrW-1, 0, TOPRIGHT, termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(0, scrH-1, BOTLEFT, termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(scrW-1, scrH-1, BOTRIGHT, termbox.ColorWhite, termbox.ColorBlack)

	completeLabel := " " + LABEL + " - " + stateLabels[gameState] + " "

	for col := 1; col < scrW-1; col++ {
		if col <= len(completeLabel) {

			termbox.SetCell(col, 0, rune(completeLabel[col-1]), termbox.ColorWhite, termbox.ColorBlack)
		} else {
			termbox.SetCell(col, 0, HLINE, termbox.ColorWhite, termbox.ColorBlack)
		}

		termbox.SetCell(col, scrH-1, HLINE, termbox.ColorWhite, termbox.ColorBlack)
	}

	for row := 1; row < scrH-1; row++ {
		termbox.SetCell(0, row, VLINE, termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(scrW-1, row, VLINE, termbox.ColorWhite, termbox.ColorBlack)
	}
}

func renderBoard(board [][]int) {
	rows := len(board)
	cols := len(board[0])

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			cell := ' '

			if board[row][col] == 1 {
				cell = ALIVE
			}

			termbox.SetCell((col*2)+1, row+1, cell, termbox.ColorWhite, termbox.ColorBlack)
			termbox.SetCell((col*2)+2, row+1, cell, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

func Render(board [][]int, gameState int) {
	renderFrame(gameState)
	renderBoard(board)
	update()
}
