package main

import (
	"os"

	"github.com/gdamore/tcell"
)

var (
	dx, dy int
)

func handleInput(s tcell.Screen, p *player) {
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
		if ev.Rune() == 'r' {
			p.AddSegment('O', tcell.StyleDefault.
				Background(tcell.ColorDarkSlateBlue).
				Foreground(tcell.ColorWhite))
		}
		if ev.Key() == tcell.KeyF12 {
			s.Fini()
			os.Exit(0)
		}
		if !gameMap.IsBlocked(p.pos[0].x+dx, p.pos[0].y+dy) {
			p.MoveEntity(dx, dy)
		}
	}
}
