package game

import (
	"github.com/gdamore/tcell"
	"github.com/stjiub/gosnake/entity"
	"github.com/stjiub/gosnake/gamemap"
	"github.com/stjiub/gosnake/style"
)

// Handle main game player input
func handleInput(g *Game) {
	// Make adjustments depending on if 1 or 2 player game
	var p2 *entity.Player
	p := g.players[0]
	if len(g.players) > 1 {
		p2 = g.players[1]
	} else {
		p2 = p
	}
	dir := p.GetDirection()
	dir2 := p2.GetDirection()
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
		if ev.Key() == tcell.KeyUp || ev.Rune() == 'w' {
			// Prevent player from turning into themselves
			if !(dir2 == DirDown) {
				p.SetDirection(DirUp)
			}
		}
		if ev.Key() == tcell.KeyUp {
			// Prevent player from turning into themselves
			if !(dir2 == DirDown) {
				p2.SetDirection(DirUp)
			}
		}
		if ev.Key() == tcell.KeyDown {
			if !(dir2 == DirUp) {
				p2.SetDirection(DirDown)
			}
		}
		if ev.Key() == tcell.KeyLeft {
			if !(dir2 == DirRight) {
				p2.SetDirection(DirLeft)
			}
		}
		if ev.Key() == tcell.KeyRight {
			//if !(dir == DirLeft) {
			p2.SetDirection(DirRight)
			//}
		}
		if ev.Rune() == 'w' {
			if !(dir == DirDown) {
				p.SetDirection(DirUp)
			}
		}
		if ev.Rune() == 's' {
			if !(dir == DirUp) {
				p.SetDirection(DirDown)
			}
		}
		if ev.Rune() == 'a' {
			if !(dir == DirRight) {
				p.SetDirection(DirLeft)
			}
		}
		if ev.Rune() == 'd' {
			if !(dir == DirLeft) {
				p.SetDirection(DirRight)
			}
		}
		if ev.Rune() == 'f' {
			p.ActivateItem()
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
			return ItemExit
		} else if ev.Key() == tcell.KeyUp {
			if s > 0 {
				m.items[s-1].selected = true
				m.items[s].selected = false
				m.ChangeSelected()
				return ItemNone
			}
		} else if ev.Key() == tcell.KeyDown {
			if s < (len(m.items) - 1) {
				m.items[s+1].selected = true
				m.items[s].selected = false
				m.ChangeSelected()
				return ItemNone
			}
		} else if ev.Key() == tcell.KeyEnter {
			return ItemEnter
		}
	}
	return ItemNone
}

// Handle profile input
func handleProfileInput(g *Game, entities []*entity.Entity, oColor, oChar *gamemap.Object, char, color *Menu, cColors []string, rotation int, fgMode bool) (int, int, []string) {
	var s, cs int
	var sty tcell.Style
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
			return ItemExit, rotation, cColors
		} else if ev.Key() == tcell.KeyLeft || ev.Rune() == 'a' {
			if s > 0 {
				char.SetSelectOnly(s - 1)
				for i := range entities {
					entities[i].SetChar(PlayerRunes[char.GetSelected()])
				}
				char.ChangeSelected()
				oChar.MoveCurPos(-2, 0)
				return ItemNone, rotation, cColors
			}
		} else if ev.Key() == tcell.KeyRight || ev.Rune() == 'd' {
			if s < (len(char.items) - 1) {
				char.SetSelectOnly(s + 1)
				for i := range entities {
					entities[i].SetChar(PlayerRunes[char.GetSelected()])
				}
				color.ChangeSelected()
				oChar.MoveCurPos(2, 0)
				return ItemNone, rotation, cColors
			}
		} else if ev.Key() == tcell.KeyUp || ev.Rune() == 'w' {
			if cs > 0 {

				color.SetSelectOnly(cs - 1)
				if fgMode {
					sty = style.StringToStyle(color.items[cs-1].str, cColors[1])
					cColors[0] = color.items[cs-1].str
				} else {
					sty = style.StringToStyle(cColors[0], color.items[cs-1].str)
					cColors[1] = color.items[cs-1].str
				}
				for i := range entities {
					entities[i].SetStyle(sty)
				}
				color.ChangeSelected()
				oColor.MoveCurPos(0, -1)
				return ItemNone, rotation, cColors
			}
		} else if ev.Key() == tcell.KeyDown || ev.Rune() == 's' {
			if cs < (len(color.items) - 1) {
				color.SetSelectOnly(cs + 1)
				if fgMode {
					sty = style.StringToStyle(color.items[cs+1].str, cColors[1])
					cColors[0] = color.items[cs+1].str
				} else {
					sty = style.StringToStyle(cColors[0], color.items[cs+1].str)
					cColors[1] = color.items[cs+1].str
				}
				for i := range entities {
					entities[i].SetStyle(sty)
				}
				char.ChangeSelected()
				oColor.MoveCurPos(0, 1)
				return ItemNone, rotation, cColors
			}
		} else if ev.Rune() == 'r' {
			switch rotation {
			case Horizontal:
				return ItemNone, DiagLeft, cColors
			case DiagLeft:
				return ItemNone, Vertical, cColors
			case Vertical:
				return ItemNone, DiagRight, cColors
			case DiagRight:
				return ItemNone, Horizontal, cColors
			}
		} else if ev.Rune() == 'c' {
			switch fgMode {
			case true:
				return BGMode, rotation, cColors
			case false:
				return FGMode, rotation, cColors
			}
		} else if ev.Key() == tcell.KeyEnter {
			return ItemEnter, rotation, cColors
		}
	}
	return ItemNone, rotation, cColors
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
