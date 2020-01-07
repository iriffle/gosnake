package main

import (
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/mattn/go-runewidth"
)

func renderAll(g *Game, style tcell.Style, gameMap *GameMap, players []*Player, bits []*Bit) {
	g.lview.Clear()
	renderMap(g, gameMap)
	renderBits(g, bits)
	renderPlayers(g, players)
	score := strconv.Itoa(players[0].score)
	renderStr(g.sview, 0, 0, style, ("Score: " + score))
}

func renderMap(g *Game, gameMap *GameMap) {
	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			if gameMap.Tiles[x][y].blocked == true {
				renderStr(g.lview, x, y, gameMap.Tiles[x][y].style, "▒")
			} else {
				renderStr(g.lview, x, y, gameMap.Tiles[x][y].style, string(' '))
			}
		}
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

func renderEntity(v *views.ViewPort, p *Player) {
	for _, pos := range p.pos {
		var comb []rune
		comb = nil
		c := pos.char
		v.SetContent(pos.x, pos.y, c, comb, pos.style)
	}

}

func renderPlayers(g *Game, players []*Player) {
	for _, player := range players {
		renderEntity(g.lview, player)
	}
}

func renderBits(g *Game, bits []*Bit) {
	for _, bit := range bits {
		renderRune(g.lview, bit.x, bit.y, bit.style, bit.char)
	}
}
