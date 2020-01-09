package main

import (
	"os"

	"github.com/gdamore/tcell"
)

func handleInput(g *Game) {
	var p2 *Player
	p := g.players[0]
	if len(g.players) > 1 {
		p2 = g.players[1]
	} else {
		p2 = p
	}
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape {
			g.screen.Fini()
			os.Exit(0)
		}
		if ev.Key() == tcell.KeyF1 {
			g.screen.Fini()
			g.state = 1
		}
		if ev.Key() == tcell.KeyF12 {
			if g.state == 2 {
				g.state = 0
			} else {
				g.state = 2
				g.Pause(p)
			}
		}
		if ev.Key() == tcell.KeyUp {
			if !(p2.direction == 2) {
				p2.direction = 1
			}
		}
		if ev.Key() == tcell.KeyDown {
			if !(p2.direction == 1) {
				p2.direction = 2
			}
		}
		if ev.Key() == tcell.KeyLeft {
			if !(p2.direction == 4) {
				p2.direction = 3
			}
		}
		if ev.Key() == tcell.KeyRight {
			if !(p2.direction == 3) {
				p2.direction = 4
			}
		}
		if ev.Rune() == 'w' {
			if !(p.direction == 2) {
				p.direction = 1
			}
		}
		if ev.Rune() == 's' {
			if !(p.direction == 1) {
				p.direction = 2
			}
		}
		if ev.Rune() == 'a' {
			if !(p.direction == 4) {
				p.direction = 3
			}
		}
		if ev.Rune() == 'd' {
			if !(p.direction == 3) {
				p.direction = 4
			}
		}
	}
}

func handleMenu(g *Game, choice int) int {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyUp {
			return 1
		} else if ev.Key() == tcell.KeyDown {
			return 2
		} else if ev.Key() == tcell.KeyEnter {
			return 3
		}
	}
	return choice
}
