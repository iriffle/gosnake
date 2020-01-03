package main

import ()

type Tile struct {
	Blocked     bool
	BlocksSight bool
	BGColor     string
	FGColor     string
}

type GameMap struct {
	Width  int
	Height int
	Tiles  [][]*Tile
}

func (m *GameMap) InitializeMap() {
	// Set up a map where all the border (edge) Tiles are walls (block movement, and sight)
	// This is just a test method, we will build maps more dynamically in the future.
	m.Tiles = make([][]*Tile, m.Width)
	for i := range m.Tiles {
		m.Tiles[i] = make([]*Tile, m.Height)
	}

	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if x == 0 || x == m.Width-1 || y == 0 || y == m.Height-1 {
				m.Tiles[x][y] = &Tile{true, true, "ColorBrown", "ColorDarkGreen"}
			} else {
				m.Tiles[x][y] = &Tile{false, false, "ColorBrown", "ColorDarkGreen"}
			}
		}
	}
}

func (m *GameMap) IsBlocked(x int, y int) bool {
	// Check to see if the provided coordinates contain a blocked tile
	if m.Tiles[x][y].Blocked {
		return true
	} else {
		return false
	}
}
