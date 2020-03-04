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
	for _, item := range m.items {
		renderStr(g.gview, item.x, item.y, item.style, item.str)
	}
	g.screen.Show()
}

// Render the name selection screen
func renderNameSelect(g *Game, w, h int, charString string) {
	g.gview.Clear()
	renderCenterStr(g.gview, w, h, g.style.DefStyle, "Name of Player:")
	renderCenterStr(g.gview, w, h+2, g.style.SelStyle, charString+"|")
	g.screen.Show()
}

// Render the High Score screen
func renderHighScoreScreen(g *Game, style tcell.Style, max int) {
	g.gview.Clear()
	g.gview.Fill(' ', style)
	renderCenterStr(g.gview, MapWidth, 4, style, "High Scores")
	renderCenterStr(g.gview, MapWidth, 6, style, strings.Repeat("=", MapWidth-10))
	renderCenterStr(g.gview, MapWidth, 10, style, "1 Player:")
	renderHighScores(g, style, Player1, 14)
	renderCenterStr(g.gview, MapWidth, (16 + max*2), style, "2 Player:")
	renderHighScores(g, style, Player2, (20 + max*2))

	g.screen.Show()
}

// Render a list of scores for the high score screen
func renderHighScores(g *Game, style tcell.Style, mode, lastScorePos int) {
	var scores []*Score
	if mode == Player1 {
		scores = g.scores1
	} else if mode == Player2 {
		scores = g.scores2
	}
	for i, score := range scores {
		if i < len(scores) {
			renderCenterStr(g.gview, MapWidth, lastScorePos+i, style, (score.Name + " - " + strconv.Itoa(score.Score)))
		} else {
			renderCenterStr(g.gview, MapWidth, lastScorePos+i, style, "----")
		}
		lastScorePos++
	}
}

// Render the player scores in middle of screen
func renderScore(v *views.ViewPort, players []*Player, w, h int, style tcell.Style) {
	scores := ""
	for i, player := range players {
		scores = player.name + ": " + strconv.Itoa(player.score) + " "
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
	for _, pos := range p.pos {
		var comb []rune
		comb = nil
		c := pos.char
		v.SetContent(pos.x, pos.y, c, comb, pos.style)
	}
}

// Render an Entity
func renderEntity(v *views.ViewPort, e *Entity) {
	for _, pos := range e.pos {
		var comb []rune
		comb = nil
		c := pos.char
		v.SetContent(pos.x, pos.y, c, comb, pos.style)
	}

}

// Render all Entities
func renderEntities(v *views.ViewPort, entities []*Entity) {
	for _, entity := range entities {
		renderEntity(v, entity)
	}
}

// Render all Players
func renderPlayers(v *views.ViewPort, players []*Player) {
	for _, player := range players {
		renderPlayer(v, player)
	}
}

// Render all Bits
func renderBits(v *views.ViewPort, bits []*Bit) {
	for _, bit := range bits {
		renderRune(v, bit.x, bit.y, bit.style, bit.char)
	}
}

// Render all Bites
func renderBites(v *views.ViewPort, bites []*Bite) {
	for _, bite := range bites {
		renderRune(v, bite.x, bite.y, bite.style, bite.char)
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

// Render a single rune to the screen
func renderRune(v *views.ViewPort, x, y int, style tcell.Style, char rune) {
	var comb []rune
	comb = nil
	v.SetContent(x, y, char, comb, style)
}
