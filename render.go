package main

import (
	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

func renderStr(s tcell.Screen, x int, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func renderMap(s tcell.Screen, gameMap *GameMap) {
	// Render the game map. If a tile is blocked and blocks sight, draw a '#', if it is not blocked, and does not block
	// sight, draw a '.'
	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			fgColor := tcell.ColorBurlyWood //bgStyle := tcell.GetColor(gameMap.Tiles[x][y].BGColor)
			//fgStyle := tcell.GetColor(gameMap.Tiles[x][y].FGColor)
			tileStyle := tcell.StyleDefault.Foreground(fgColor)
			if gameMap.Tiles[x][y].Blocked == true {
				renderStr(s, x, y, tileStyle, "#")
			} else {
				renderStr(s, x, y, tileStyle, string(' '))
			}
		}
	}
}

func renderEntities(s tcell.Screen, entities []*Entity) {
	// Draw every Entity present in the game. This gets called on each iteration of the game loop.
	for _, e := range entities {
		renderStr(s, e.x, e.y, e.style, e.char)
	}
}

func renderAll(s tcell.Screen, style tcell.Style, gameMap *GameMap, entities []*Entity) {
	// Convenience function to render all entities, followed by rendering the game map
	renderMap(s, gameMap)
	renderEntities(s, entities)
}
