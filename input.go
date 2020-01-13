package main

import (
	"os"

	"github.com/gdamore/tcell"
)

// Handle main game input
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
			g.state = MainMenu
			return
		} else if ev.Key() == tcell.KeyExit {
			g.state = Quit
			return
		} else if ev.Key() == tcell.KeyF1 {
			g.screen.Fini()
			g.state = Restart
			return
		} else if ev.Key() == tcell.KeyF12 {
			if g.state == Play {
				g.state = Pause
				return
			} else if g.state == Pause {
				g.state = Play
				return
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

// Handle paused game input
func handlePause(g *Game) {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			g.state = Quit
			return
		case tcell.KeyExit:
			g.state = Quit
			return
		case tcell.KeyF1:
			g.screen.Fini()
			g.state = Restart
			return
		case tcell.KeyF12:
			g.state = Play
			return
		}
	}
}

// Handle main menu input
func handleMenu(g *Game, m *Menu) bool {
	var s int
	for i, item := range m.items {
		if item.selected {
			s = i
		}
	}
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyExit {
			g.screen.Fini()
			os.Exit(0)
		} else if ev.Key() == tcell.KeyUp {
			if s > 0 {
				m.items[s-1].selected = true
				m.items[s].selected = false
				m.ChangeSelected()
				return false
			}
		} else if ev.Key() == tcell.KeyDown {
			if s < (len(m.items) - 1) {
				m.items[s+1].selected = true
				m.items[s].selected = false
				m.ChangeSelected()
				return false
			}
		} else if ev.Key() == tcell.KeyEnter {
			return true
		}
	}
	return false
}
