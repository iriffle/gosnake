package main

import (
	"os"

	"github.com/gdamore/tcell"
)

// var (
// 	dx, dy int
// )

func handleInput(s tcell.Screen, p *Player) {
	ev := s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyF12 {
			s.Fini()
			os.Exit(0)
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
