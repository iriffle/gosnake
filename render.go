package main

import (
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/mattn/go-runewidth"
)

func renderAll(g *Game, style tcell.Style, m *GameMap, players []*Player, bits []*Bit) {
	g.gview.Clear()
	renderMap(g.gview, m)
	if g.numPlayers == 1 {
		renderLevel(g.gview, g.level, m.Width, m.Height, style)
	}
	renderScore(g.gview, g.players, m.Width, m.Height, style)
	renderBits(g.gview, bits)
	renderPlayers(g.gview, players)
	g.sbar.SetCenter(Controls, ControlStyle)
	g.sbar.Draw()
	if g.debug {
		renderConsole(g, cviewWidth, cviewHeight, DebugStyle)
		g.cbar.Draw()
	}
	g.screen.Show()
}

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

func renderMenu(v *views.ViewPort, w, h int, style tcell.Style) {
	renderCenterStr(v, w, h-4, style, "1 Player")
	renderCenterStr(v, w, h, style, "2 Player")
}

func renderScore(v *views.ViewPort, players []*Player, w, h int, style tcell.Style) {
	scores := ""
	for i, player := range players {
		scores = player.name + ": " + strconv.Itoa(player.score) + " "
		renderCenterStr(v, w, h+i, style, scores)
	}
}

func renderLevel(v *views.ViewPort, l, w, h int, style tcell.Style) {
	level := "level: " + strconv.Itoa(l)
	renderCenterStr(v, w, h-2, style, level)
}

func renderEntity(v *views.ViewPort, p *Player) {
	for _, pos := range p.pos {
		var comb []rune
		comb = nil
		c := pos.char
		v.SetContent(pos.x, pos.y, c, comb, pos.style)
	}

}

func renderPlayers(v *views.ViewPort, players []*Player) {
	for _, player := range players {
		renderEntity(v, player)
	}
}

func renderBits(v *views.ViewPort, bits []*Bit) {
	for _, bit := range bits {
		renderRune(v, bit.x, bit.y, bit.style, bit.char)
	}
}

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

func renderCenterStr(v *views.ViewPort, w, h int, style tcell.Style, str string) {
	x := (w / 2) - (len(str) / 2)
	y := (h / 2)
	renderStr(v, x, y, style, str)
}

func renderRune(v *views.ViewPort, x, y int, style tcell.Style, char rune) {
	var comb []rune
	comb = nil
	v.SetContent(x, y, char, comb, style)
}

func renderConsole(g *Game, w, h int, style tcell.Style) {
	renderCenterStr(g.cview, w, h, style, strconv.Itoa(g.state))
}
