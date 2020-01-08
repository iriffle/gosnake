package main

import (
	"github.com/gdamore/tcell"
)

type GameMap struct {
	Width   int
	Height  int
	Objects [][]*Object
}

func (m *GameMap) InitializeMap(defStyle tcell.Style) {
	// Set up a map where all the border (edge) Tiles are walls (block movement, and sight)
	// This is just a test method, we will build maps more dynamically in the future.
	m.Objects = make([][]*Object, m.Width)
	for i := range m.Objects {
		m.Objects[i] = make([]*Object, m.Height)
	}

	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if x == 0 || x == m.Width-1 || y == 0 || y == m.Height-1 {
				m.Objects[x][y] = &Object{x, y, x, y, '▒', defStyle, true}
			} else {
				m.Objects[x][y] = &Object{x, y, x, y, ' ', defStyle, false}
			}
		}
	}
}

func (m *GameMap) IsBlocked(x int, y int) bool {
	// Check to see if the provided coordinates contain a blocked tile
	if m.Objects[x][y].blocked {
		return true
	} else {
		return false
	}
}
