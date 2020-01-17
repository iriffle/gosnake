package main

import (
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/mattn/go-runewidth"
)

// Render all of the game in main game loop
func renderAll(g *Game, style tcell.Style, m *GameMap) {
	g.gview.Clear()
	renderMap(g.gview, m)
	for _, b := range g.maps {
		renderMap(g.gview, b)
	}
	if g.numPlayers == 1 {
		renderLevel(g.gview, g.level, m.Width, m.Height, SelStyle)
	}
	renderScore(g.gview, g.players, m.Width, m.Height, SelStyle)
	renderBits(g.gview, g.bits)
	renderBites(g.gview, g.bites)
	renderEntities(g.gview, g.entities)
	renderPlayers(g.gview, g.players)
	g.sbar.SetCenter(controls, DefStyle)
	g.sbar.Draw()
	g.screen.Show()
}

// Render the game map
func renderMap(v *views.ViewPort, m *GameMap) {
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if m.Objects[x][y].blocked == true {
				renderRune(v, x, y, m.Objects[x][y].style, m.Objects[x][y].char)
			} else {
				renderRune(v, x, y, m.Objects[x][y].style, m.Objects[x][y].char)
			}
		}
	}
}

// Render the main menu
func renderMenu(g *Game, m *Menu, style tcell.Style) {
	g.gview.Fill(' ', style)
	for _, item := range m.items {
		renderStr(g.gview, item.x, item.y, item.style, item.str)
	}
	g.screen.Show()
}

// Render the player scores
func renderScore(v *views.ViewPort, players []*Player, w, h int, style tcell.Style) {
	scores := ""
	for i, player := range players {
		scores = player.name + ": " + strconv.Itoa(player.score) + " "
		renderCenterStr(v, w, h+i, style, scores)
	}
}

// Render the current level
func renderLevel(v *views.ViewPort, l, w, h int, style tcell.Style) {
	level := "level: " + strconv.Itoa(l)
	renderCenterStr(v, w, h-2, style, level)
}

// Render an entity
func renderPlayer(v *views.ViewPort, p *Player) {
	for _, pos := range p.pos {
		var comb []rune
		comb = nil
		c := pos.char
		v.SetContent(pos.x, pos.y, c, comb, pos.style)
	}

}

// Render an entity
func renderEntity(v *views.ViewPort, e *Entity) {
	for _, pos := range e.pos {
		var comb []rune
		comb = nil
		c := pos.char
		v.SetContent(pos.x, pos.y, c, comb, pos.style)
	}

}

func renderEntities(v *views.ViewPort, entities []*Entity) {
	for _, entity := range entities {
		renderEntity(v, entity)
	}
}

// Render all players
func renderPlayers(v *views.ViewPort, players []*Player) {
	for _, player := range players {
		renderPlayer(v, player)
	}
}

// Render all bits
func renderBits(v *views.ViewPort, bits []*Bit) {
	for _, bit := range bits {
		renderRune(v, bit.x, bit.y, bit.style, bit.char)
	}
}

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
