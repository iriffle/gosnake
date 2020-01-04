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
		if ev.Rune() == 'r' {
			add := Coord{
				x: (p.snake.moving.object.x - 1),
				y: (p.snake.moving.object.y - 1),
			}
			p.snake.pos = append(p.snake.pos, add)
		}
		if ev.Key() == tcell.KeyF12 {
			s.Fini()
			os.Exit(0)
		}
		if !gameMap.IsBlocked(p.snake.moving.object.x+dx, p.snake.moving.object.y+dy) {
			p.snake.moving.Move(dx, dy)
		}
	}
}
