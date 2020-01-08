package main

import (
	"github.com/gdamore/tcell"
)

type GameMap struct {
	Width   int
	Height  int
	Objects [][]*Object
}

func (m *GameMap) InitMap() {
	m.Objects = make([][]*Object, m.Width)
	for i := range m.Objects {
		m.Objects[i] = make([]*Object, m.Height)
	}
}

func (m *GameMap) InitMapBoundary(wallrune, floorRune rune, style tcell.Style) {

	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if x == 0 || x == m.Width-1 || y == 0 || y == m.Height-1 {
				m.Objects[x][y] = &Object{x, y, x, y, wallRune, style, true}
			} else {
				m.Objects[x][y] = &Object{x, y, x, y, floorRune, style, false}
			}
		}
	}
}

func (m *GameMap) InitLevel1(wallRune, floorRune rune, style tcell.Style) {
	m.InitMap()
	m.InitMapBoundary(wallRune, floorRune, style)
}
