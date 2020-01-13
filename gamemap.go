package main

import (
	"time"

	"github.com/gdamore/tcell"
)

// The game map struct
type GameMap struct {
	Width   int
	Height  int
	Objects [][]*Object
}

// Generate an empty map
func (m *GameMap) InitMap() {
	m.Objects = make([][]*Object, m.Width)
	for i := range m.Objects {
		m.Objects[i] = make([]*Object, m.Height)
	}
}

// Generate walls around perimeter of map
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

func (m *GameMap) InitLevel(wallrune, floorRune rune, style tcell.Style) {
	m.InitMap()
	m.InitMapBoundary(wallRune, floorRune, style)
}

// Generate level 1 map which is just an open map with walls around perimeter
func (m *GameMap) InitLevel1(g *Game, wallRune, floorRune, bitRune rune, style tcell.Style) {
	m.Width = 80
	m.Height = 20
	m.InitLevel(wallRune, floorRune, style)
	go m.RandomBits(g, 2, 10, 3*time.Second)
	go m.RandomLines(g, 2)
}

func (m *GameMap) InitLevel2(g *Game, wallRune, floorRune rune, style tcell.Style) {
	m.Width = 90
	m.Height = 30
	m.InitLevel(wallRune, floorRune, style)
	g.bitQuit = make(chan bool)
	go g.handleBits(m)
}

func (m *GameMap) InitLevel3(g *Game, wallRune, floorRune rune, style tcell.Style) {
	m.Width = 100
	m.Height = 35
	m.InitLevel(wallRune, floorRune, style)
	go m.RandomBites(g, 1, 3, (20 * time.Second))
}

func (m *GameMap) RandomLines(g *Game, numTimes int) {
	for i := 0; i < numTimes; i++ {
		NewBitLineH(g, (m.Width/2)-6, m.Height-3, 10, 6, bitRune, BitStyle)
		NewBitLineH(g, (m.Width/2)-6, GameStartY+1, 10, 6, bitRune, BitStyle)
		time.Sleep(10 * time.Second)
		NewBitLineV(g, GameStartX+1, m.Height/2-6, 10, 6, bitRune, BitStyle)
		NewBitLineV(g, m.Width-2, m.Height/2-6, 10, 6, bitRune, BitStyle)
		time.Sleep(10 * time.Second)
	}
}

func (m *GameMap) RandomBits(g *Game, bitsGen, bitsMax int, dur time.Duration) {
	for {
		for i := 0; i < bitsGen; i++ {
			if len(g.bits)-bitsGen < bitsMax {
				newB := NewRandomBit(m, 10, bitRune, BitStyle)
				g.bits = append(g.bits, &newB)
			}
		}
		time.Sleep(dur)

	}
}

func (m *GameMap) RandomBites(g *Game, bitesGen, bitesMax int, dur time.Duration) {
	for {
		for i := 0; i < bitesGen; i++ {
			if len(g.bites)-bitesGen < bitesMax {
				newB := NewRandomBite(m, BiteRune, BiteStyle)
				g.bites = append(g.bites, &newB)
			}
		}
		time.Sleep(dur)
	}
}
