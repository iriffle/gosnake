package main

import (
	"os"

	"github.com/gdamore/tcell"
)

func handleInput(g *Game, p *Player) {
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
		switch ev.Rune() {
		case 'w':
			if !(p.direction == 2) {
				p.direction = 1
			}
		case 's':
			if !(p.direction == 1) {
				p.direction = 2
			}
		case 'a':
			if !(p.direction == 4) {
				p.direction = 3
			}
		case 'd':
			if !(p.direction == 3) {
				p.direction = 4
			}
		case 'r':
			p.AddSegment(p.pos[0].char, p.pos[0].style)
		}
	}
}
