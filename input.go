package main

import (
	"github.com/gdamore/tcell"
)

// Handle main game player input
func handleInput(g *Game) {

	// Make adjustments depending on if 1 or 2 player game
	var p2 *Player
	p := g.players[0]
	if len(g.players) > 1 {
		p2 = g.players[1]
	} else {
		p2 = p
	}

	// Wait for input events and process
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:

		// Quit game and return to Main Menu if Escape key pressed
		if ev.Key() == tcell.KeyEscape {
			g.screen.Fini()
			g.state = MainMenu
			return

			// Quit game and close window if terminal window exited
		} else if ev.Key() == tcell.KeyExit {
			g.state = Quit
			return

			// Restart game if F1 pressed
		} else if ev.Key() == tcell.KeyF1 {
			g.screen.Fini()
			g.state = Restart
			return

			// Pause game if F12 pressed
		} else if ev.Key() == tcell.KeyF12 {
			if g.state == Play {
				g.state = Pause
				return
			} else if g.state == Pause {
				g.state = Play
				return
			}
		}

		// Handle player direction. Can use WSAD or Arrow keys
		// for 1 player. 2 player splits these up with WSAD for
		// player1 and Arrow keys for player2
		if ev.Key() == tcell.KeyUp {
			// Prevent player from turning into themselves
			if !(p2.direction == DirDown) {
				p2.direction = DirUp
			}
		}
		if ev.Key() == tcell.KeyDown {
			if !(p2.direction == DirUp) {
				p2.direction = DirDown
			}
		}
		if ev.Key() == tcell.KeyLeft {
			if !(p2.direction == DirRight) {
				p2.direction = DirLeft
			}
		}
		if ev.Key() == tcell.KeyRight {
			if !(p2.direction == DirLeft) {
				p2.direction = DirRight
			}
		}
		if ev.Rune() == 'w' {
			if !(p.direction == DirDown) {
				p.direction = DirUp
			}
		}
		if ev.Rune() == 's' {
			if !(p.direction == DirUp) {
				p.direction = DirDown
			}
		}
		if ev.Rune() == 'a' {
			if !(p.direction == DirRight) {
				p.direction = DirLeft
			}
		}
		if ev.Rune() == 'd' {
			if !(p.direction == DirLeft) {
				p.direction = DirRight
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
func handleMenuInput(g *Game, m *Menu) int {
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
			g.state = Quit
			return -1
		} else if ev.Key() == tcell.KeyUp {
			if s > 0 {
				m.items[s-1].selected = true
				m.items[s].selected = false
				m.ChangeSelected()
				return 0
			}
		} else if ev.Key() == tcell.KeyDown {
			if s < (len(m.items) - 1) {
				m.items[s+1].selected = true
				m.items[s].selected = false
				m.ChangeSelected()
				return 0
			}
		} else if ev.Key() == tcell.KeyEnter {
			return 1
		}
	}
	return 0
}
