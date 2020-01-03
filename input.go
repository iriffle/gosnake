package main

import (
	"github.com/gdamore/tcell"
	"os"
)

func handleInput(s tcell.Screen, entity *Entity) {
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
		if ev.Key() == tcell.KeyF12 {
			s.Fini()
			os.Exit(0)
		}
		if !gameMap.IsBlocked(entity.x+dx, entity.y+dy) {
			entity.Move(dx, dy)
		}
	}
}
