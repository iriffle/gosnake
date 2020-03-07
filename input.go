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
			g.Return()
			return

			// Quit game and close window if terminal window exited
		} else if ev.Key() == tcell.KeyExit {
			g.Quit()
			return

			// Restart game if F1 pressed
		} else if ev.Key() == tcell.KeyF1 {
			g.Restart()
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
	for i := range m.items {
		if m.items[i].selected {
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

// Handle profile input
func handleProfileInput(g *Game, e *Entity, oColor, oChar *Object, char, color *Menu, curColors []string, rotation int, cMode bool) (int, []string) {
	var s, cs int
	var style tcell.Style
	for i := range char.items {
		if char.items[i].selected {
			s = i
		}
	}
	for i := range color.items {
		if color.items[i].selected {
			cs = i
		}
	}
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyExit {
			g.state = Quit
			return -1, curColors
		} else if ev.Key() == tcell.KeyLeft || ev.Rune() == 'a' {
			if s > 0 {
				char.items[s-1].selected = true
				char.items[s].selected = false
				e.SetChar(PlayerRunes[char.GetSelected()])
				char.ChangeSelected()
				oChar.x -= 2
				return 0, curColors
			}
		} else if ev.Key() == tcell.KeyRight || ev.Rune() == 'd' {
			if s < (len(char.items) - 1) {
				char.items[s+1].selected = true
				char.items[s].selected = false
				e.SetChar(PlayerRunes[char.GetSelected()])
				char.ChangeSelected()
				oChar.x += 2
				return 0, curColors
			}
		} else if ev.Key() == tcell.KeyUp || ev.Rune() == 'w' {
			if cs > 0 {

				color.items[cs-1].selected = true
				color.items[cs].selected = false
				if cMode {
					style = StringToStyle(color.items[cs-1].str, curColors[1])
					curColors[0] = color.items[cs-1].str
				} else {
					style = StringToStyle(curColors[0], color.items[cs-1].str)
					curColors[1] = color.items[cs-1].str
				}
				e.SetStyle(style)
				char.ChangeSelected()
				oColor.y--
				return 0, curColors
			}
		} else if ev.Key() == tcell.KeyDown || ev.Rune() == 's' {
			if cs < (len(color.items) - 1) {
				color.items[cs+1].selected = true
				color.items[cs].selected = false
				if cMode {
					style = StringToStyle(color.items[cs+1].str, curColors[1])
					curColors[0] = color.items[cs+1].str
				} else {
					style = StringToStyle(curColors[0], color.items[cs+1].str)
					curColors[1] = color.items[cs+1].str
				}
				e.SetStyle(style)
				char.ChangeSelected()
				oColor.y++
				return 0, curColors
			}
		} else if ev.Rune() == 'r' {
			switch rotation {
			case 0:
				return 2, curColors
			case 1:
				return 3, curColors
			case 2:
				return 4, curColors
			case 3:
				return 5, curColors
			}
		} else if ev.Rune() == 'c' {
			switch cMode {
			case true:
				return 6, curColors
			case false:
				return 7, curColors
			}
		} else if ev.Key() == tcell.KeyEnter {
			return 1, curColors
		}
	}
	return 0, curColors
}

func handleStringInput(g *Game) rune {
	ev := g.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyExit {
			return '\v'
		} else if ev.Key() == tcell.KeyEnter {
			return '\r'
		} else if ev.Key() == tcell.KeyBackspace {
			return '\t'
		} else {
			return ev.Rune()
		}
	}
	return '\n'
}
