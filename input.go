package main

import (
	"os"

	"github.com/gdamore/tcell"
)

func handleInput(s tcell.Screen, p *Player) {
	ev := s.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Rune() == 'w' {
			dx, dy = 0, -1
		}
		if ev.Rune() == 's' {
			dx, dy = 0, 1
		}
		if ev.Rune() == 'a' {
			dx, dy = -1, 0
		}
		if ev.Rune() == 'd' {
			dx, dy = 1, 0
		}
		// if ev.Rune() == 'r' {
		// 	add := p.pos{
		// 		x: (p.Object.x - 1),
		// 		y: (p.Object.y - 1),
		// 	}
		// 	p.pos = append(p.pos, p.Object.x)
		// }
		if ev.Key() == tcell.KeyF12 {
			s.Fini()
			os.Exit(0)
		}
		if !gameMap.IsBlocked(p.Object.x+dx, p.Object.y+dy) {
			p.MovePlayer(dx, dy)
		}
	}
}
