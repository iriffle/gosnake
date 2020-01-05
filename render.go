package main

import (
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/mattn/go-runewidth"
)

func renderAll(g *Game, style tcell.Style, gameMap *GameMap, players []*player, bits []*bit) {
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
			fgColor := tcell.ColorBurlyWood
			tileStyle := tcell.StyleDefault.Foreground(fgColor)
			if gameMap.Tiles[x][y].Blocked == true {
				renderStr(g.lview, x, y, tileStyle, "▒")
			} else {
				renderStr(g.lview, x, y, tileStyle, string(' '))
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

func renderRune(v *views.ViewPort, x, y int, style tcell.Style, char rune) {
	var comb []rune
	comb = nil
	v.SetContent(x, y, char, comb, style)
}

func renderEntity(v *views.ViewPort, p *player) {
	for _, pos := range p.pos {
		var comb []rune
		comb = nil
		c := pos.char
		v.SetContent(pos.x, pos.y, c, comb, pos.style)
	}

}

func renderPlayers(g *Game, players []*player) {
	for _, player := range players {
		renderEntity(g.lview, player)
	}
}

func renderBits(g *Game, bits []*bit) {
	for _, bit := range bits {
		renderRune(g.lview, bit.x, bit.y, bit.style, bit.char)
	}
}
