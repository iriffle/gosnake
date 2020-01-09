package main

import (
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/mattn/go-runewidth"
)

var scores string

func renderAll(g *Game, style tcell.Style, gameMap *GameMap, players []*Player, bits []*Bit) {
	g.gview.Clear()
	renderMap(g.gview, gameMap)
	level := "Level: " + strconv.Itoa(g.level)
	renderCenterStr(g.gview, gameMap.Width, gameMap.Height-2, style, level)
	for _, player := range g.players {
		scores = "Score: " + strconv.Itoa(player.score) + " "
	}
	renderCenterStr(g.gview, gameMap.Width, gameMap.Height, style, scores)
	renderBits(g.gview, bits)
	renderPlayers(g.gview, players)
	g.cbar.SetCenter(Controls, ControlStyle)
	g.screen.Show()

}

func renderMap(v *views.ViewPort, gameMap *GameMap) {
	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			if gameMap.Objects[x][y].blocked == true {
				renderRune(v, x, y, gameMap.Objects[x][y].style, gameMap.Objects[x][y].char)
			} else {
				renderRune(v, x, y, gameMap.Objects[x][y].style, gameMap.Objects[x][y].char)
			}
		}
	}
}

func renderMenu(v *views.ViewPort, w, h int, style tcell.Style) {
	renderCenterStr(v, w, h-4, style, "1 Player")
	renderCenterStr(v, w, h, style, "2 Player")
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
