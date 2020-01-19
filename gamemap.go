package main

import (
	"time"

	"github.com/gdamore/tcell"
)

// The game map struct
type GameMap struct {
	Width   int
	Height  int
	X       int
	Y       int
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
func (m *GameMap) InitMapBoundary(wallRune, floorRune rune, style tcell.Style) {

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

// Generate level 1 map which is just an open map with walls around perimeter
func (m *GameMap) InitLevel1(g *Game) {
	go m.RandomBits(g, 2, 10, 3*time.Second)
	go m.RandomLines(g, 2)
}

func (m *GameMap) InitLevel2(g *Game) {
	g.bitQuit = make(chan bool)
	go g.handleBits(m)
}

func (m *GameMap) InitLevel3(g *Game) {
	go m.RandomBites(g, 1, 3, (20 * time.Second), false)
}

func (m *GameMap) InitLevel4(g *Game) {
	go m.MovingWall(g, m.Width/2, m.Height/4, 2, 1, 10, WallRune, DefStyle)
	go m.MovingWall(g, m.Width/2, (m.Height - m.Height/4), 3, 1, 10, WallRune, DefStyle)
}

func (m *GameMap) RandomLines(g *Game, numTimes int) {
	//for i := 0; i < numTimes; i++ {
	for {
		NewRandomBitLine(g, m, 10, BitRune, BitStyle)
		time.Sleep(15 * time.Second)
	}
}

func (m *GameMap) RandomBits(g *Game, bitsGen, bitsMax int, dur time.Duration) {
	for {
		for i := 0; i < bitsGen; i++ {
			if len(g.bits)-bitsGen < bitsMax {
				newB := NewRandomBit(m, 10, BitRune, BitStyle)
				g.bits = append(g.bits, &newB)
			}
		}
		time.Sleep(dur)

	}
}

func (m *GameMap) RandomBites(g *Game, bitesGen, bitesMax int, dur time.Duration, random bool) {
	for {
		for i := 0; i < bitesGen; i++ {
			if len(g.bites)-bitesGen < bitesMax {
				newB := NewRandomBite(m, BiteStyle, random)
				g.bites = append(g.bites, &newB)
			}
		}
		time.Sleep(dur)
	}
}

func (m *GameMap) MovingWall(g *Game, x, y, direction, speed, segments int, char rune, style tcell.Style) {
	e := NewEntity(x, y, direction, speed, char, style)
	for i := 0; i < segments; i++ {
		e.AddSegment(char, style)
	}
	g.entities = append(g.entities, &e)
	for {
		dx, dy := e.CheckDirection(g)
		if e.IsBlockedByMap(m, dx, dy) {
			var newPos []Object
			for i, _ := range e.pos {
				newPos = append(newPos, e.pos[len(e.pos)-1-i])
			}
			e.pos = newPos
			switch e.direction {
			case DirUp:
				e.direction = DirDown
			case DirDown:
				e.direction = DirUp
			case DirLeft:
				e.direction = DirRight
			case DirRight:
				e.direction = DirLeft
			}
		} else {
			e.Move(dx, dy)
		}
		time.Sleep(g.moveInterval(e.speed, e.direction))
	}
}
