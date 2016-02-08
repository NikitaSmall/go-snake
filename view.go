package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/nsf/termbox-go"
)

const (
	bgColor      = termbox.ColorBlue
	boardColor   = termbox.ColorBlack
	commentColor = termbox.ColorYellow

	defaultMarginWidth  = 2
	defaultMarginHeight = 1

	titleStartX = defaultMarginWidth
	titleStartY = defaultMarginHeight

	titleHeight = 1
	titleEndY   = titleStartY + titleHeight

	boardStartX = defaultMarginWidth
	boardStartY = defaultMarginHeight + titleEndY

	boardWidth  = 25
	boardHeight = 25

	cellWidth = 2

	boardEndX = boardStartX + boardWidth*cellWidth
	boardEndY = boardStartY + boardHeight

	instructionsStartX = boardEndX + defaultMarginWidth
	instructionsStartY = boardStartY

	title = "Snake in Golang"
)

var pieceColors = []termbox.Attribute{
	termbox.ColorBlack,
	termbox.ColorYellow,
	termbox.ColorRed,
	termbox.ColorWhite,
}

var instructions = []string{
	"Goal: Eat all the frogs!",
	"",
	"left   Left",
	"right  Right",
	"up     Up",
	"down   Down",
	"s      Start",
	"p      Pause",
	"esc,q  Exit",
	"",
	"Level: %v",
	"",
	"GAME OVER!",
}

func render(g *Game) {
	termbox.Clear(bgColor, bgColor)
	stringPrint(titleStartX, titleStartY, commentColor, bgColor, title)

	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			cellValue := g.board[y][x]
			absCellValue := int(math.Abs(float64(cellValue)))

			cellColor := pieceColors[absCellValue]
			for i := 0; i < cellWidth; i++ {
				termbox.SetCell(boardStartX+cellWidth*x+i, boardStartY+y, ' ', cellColor, cellColor)
			}
		}
	}

	for y, instruction := range instructions {
		switch {
		case strings.HasPrefix(instruction, "Level:"):
			instruction = fmt.Sprintf(instruction, g.level)
		case strings.HasPrefix(instruction, "GAME OVER"):
			if g.state != gameOver {
				instruction = ""
			}
		default:
		}

		stringPrint(instructionsStartX, instructionsStartY+y, commentColor, bgColor, instruction)
	}

	termbox.Flush()
}

func stringPrint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, char := range msg {
		termbox.SetCell(x, y, char, fg, bg)
		x++
	}
}
