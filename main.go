package main

import (
	"log"
	"time"

	"github.com/nsf/termbox-go"
)

const animationSpeed = 10 * time.Millisecond

var eventQueue = make(chan termbox.Event)

func main() {
	err := termbox.Init()
	if err != nil {
		log.Print(err)
		return
	}
	defer termbox.Close()

	g := NewGame()
	go pollEvents()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Ch == 'q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyCtrlD {
					return
				} else {
					handleKeyPress(g, ev)
					g.play()
				}
			}
		case <-g.turnTimer.C:
			g.play()
		default:
			render(g)
			time.Sleep(animationSpeed)
		}
	}
}

func pollEvents() {
	for {
		eventQueue <- termbox.PollEvent()
	}
}

func handleKeyPress(g *Game, ev termbox.Event) {
	switch {
	case ev.Key == termbox.KeyArrowLeft:
		g.moveLeft()
	case ev.Key == termbox.KeyArrowRight:
		g.moveRight()
	case ev.Key == termbox.KeyArrowUp:
		g.moveUp()
	case ev.Key == termbox.KeyArrowDown:
		g.moveDown()
	case ev.Ch == 's':
		g.start()
	case ev.Ch == 'p':
		g.pause()
	}
}
