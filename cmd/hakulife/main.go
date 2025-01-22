package main

import (
	"flag"

	"github.com/hakuthedeadboi/hakulife/pkg/gameoflife"
	"github.com/hakuthedeadboi/hakulife/pkg/graphics"
)

func main() {
	msecs := flag.Int("speed", 1000, "Time for one cycle in milliseconds.")
	flag.Parse()

	graphics.Init()
	defer graphics.Close()

	scrW, scrH := graphics.GetScreenSize()
	boardCols := ((scrW - 2) - (scrW % 2)) / 2
	boardRows := ((scrH - 2) - (scrH % 2))

	game := gameoflife.NewGame(boardRows, boardCols, *msecs)
	game.SetDrawingFunc(graphics.Render)
	game.Start()
}
