package main

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/mattn/go-runewidth"
)

// Render all of the game in main game loop
func renderAll(g *Game, style tcell.Style, m *GameMap) {

	// Clear screen for redraw
	g.gview.Clear()

	// Draw game map
	renderMap(g.gview, m)

	// Draw the Bite explosion map
	renderMap(g.gview, g.biteMap)

	if g.numPlayers == 1 {
		renderLevel(g.gview, g.level, m.Width, m.Height, g.style.SelStyle)
	}
	renderScore(g.gview, g.players, m.Width, m.Height, g.style.SelStyle)
	renderBits(g.gview, g.bits)
	renderBites(g.gview, g.bites)
	renderEntities(g.gview, g.entities)
	renderPlayers(g.gview, g.players)
	g.sbar.SetCenter(controls, g.style.DefStyle)
	g.sbar.Draw()
	g.screen.Show()
}

// Render the game map
func renderMap(v *views.ViewPort, m *GameMap) {
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			renderRune(v, x, y, m.Objects[x][y].style, m.Objects[x][y].char)
		}
	}
}

// Render a menu
func renderMenu(g *Game, m *Menu, style tcell.Style) {
	for i := range m.items {
		renderStr(g.gview, m.items[i].x, m.items[i].y, m.items[i].style, m.items[i].str)
	}
	g.screen.Show()
}

func renderProfile(g *Game, m *Menu, w, h int, style tcell.Style) {
	renderChars(g.gview, m, w, h)
	g.screen.Show()
}

// Render the name selection screen
func renderNameSelect(g *Game, w, h int, hStr, charStr string) {
	g.gview.Clear()
	renderCenterStr(g.gview, w, h, g.style.DefStyle, hStr)
	renderCenterStr(g.gview, w, h+2, g.style.SelStyle, charStr+"|")
	g.screen.Show()
}

// Render the High Score screen
func renderHighScoreScreen(g *Game, style tcell.Style, max int) {
	g.gview.Clear()
	g.gview.Fill(' ', style)
	renderCenterStr(g.gview, MapWidth, 4, style, "High Scores")
	renderCenterStr(g.gview, MapWidth, 6, style, strings.Repeat("=", MapWidth-10))
	renderCenterStr(g.gview, MapWidth, 10, style, "1 Player:")
	renderHighScores(g, Player1, 14)
	renderCenterStr(g.gview, MapWidth, (16 + max*2), style, "2 Player:")
	renderHighScores(g, Player2, (20 + max*2))

	g.screen.Show()
}

// Render a list of scores for the high score screen
func renderHighScores(g *Game, mode, lastScorePos int) {
	var scores []*Score
	if mode == Player1 {
		scores = g.scores1
	} else if mode == Player2 {
		scores = g.scores2
	}
	for i := range scores {
		if i < len(scores) {
			renderCenterStr(g.gview, MapWidth, lastScorePos+i, g.style.SelStyle, (scores[i].Name + " - " + strconv.Itoa(scores[i].Score)))
		} else {
			renderCenterStr(g.gview, MapWidth, lastScorePos+i, g.style.DefStyle, "----")
		}
		lastScorePos++
	}
}

// Render the player scores in middle of screen
func renderScore(v *views.ViewPort, players []*Player, w, h int, style tcell.Style) {
	scores := ""
	for i := range players {
		scores = players[i].name + ": " + strconv.Itoa(players[i].score) + " "
		renderCenterStr(v, w, h+i, style, scores)
	}
}

// Render the current level in middle of screen
func renderLevel(v *views.ViewPort, l, w, h int, style tcell.Style) {
	level := "level: " + strconv.Itoa(l)
	renderCenterStr(v, w, h-2, style, level)
}

// Render a Player
func renderPlayer(v *views.ViewPort, p *Player) {
	for i := range p.pos {
		var comb []rune
		comb = nil
		c := p.pos[i].char
		v.SetContent(p.pos[i].x, p.pos[i].y, c, comb, p.pos[i].style)
	}
}

// Render an Entity
func renderEntity(v *views.ViewPort, e *Entity) {
	for i := range e.pos {
		var comb []rune
		comb = nil
		c := e.pos[i].char
		v.SetContent(e.pos[i].x, e.pos[i].y, c, comb, e.pos[i].style)
	}

}

// Render all Entities
func renderEntities(v *views.ViewPort, entities []*Entity) {
	for i := range entities {
		renderEntity(v, entities[i])
	}
}

// Render all Players
func renderPlayers(v *views.ViewPort, players []*Player) {
	for i := range players {
		renderPlayer(v, players[i])
	}
}

// Render all objects
func renderObjects(v *views.ViewPort, objects []*Object) {
	for i := range objects {
		renderRune(v, objects[i].x, objects[i].y, objects[i].style, objects[i].char)
	}
}

// Render all Bits
func renderBits(v *views.ViewPort, bits []*Bit) {
	for i := range bits {
		renderRune(v, bits[i].x, bits[i].y, bits[i].style, bits[i].char)
	}
}

// Render all Bites
func renderBites(v *views.ViewPort, bites []*Bite) {
	for i := range bites {
		renderRune(v, bites[i].x, bites[i].y, bites[i].style, bites[i].char)
	}
}

// Render a string at given position
func renderStr(v *views.ViewPort, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		v.SetContent(x, y, c, comb, style)
		x += w
	}
}

// Render a string in the center of the screen
func renderCenterStr(v *views.ViewPort, w, h int, style tcell.Style, str string) {
	x := (w / 2) - (len(str) / 2)
	y := (h / 2)
	renderStr(v, x, y, style, str)
}

// Render a string at given position
func renderChars(v *views.ViewPort, m *Menu, x, y int) {
	for _, i := range m.items {
		for _, c := range i.str {
			var comb []rune
			w := runewidth.RuneWidth(c)
			if w == 0 {
				comb = []rune{c}
				c = ' '
				w = 1
			}
			v.SetContent(x, y, c, comb, i.style)
			x += w + 1
		}
	}
}

// Render a single rune to the screen
func renderRune(v *views.ViewPort, x, y int, style tcell.Style, char rune) {
	var comb []rune
	comb = nil
	v.SetContent(x, y, char, comb, style)
}
