package main

import (
	"time"

	"github.com/gdamore/tcell"
)

// The game map struct
type GameMap struct {
	Width    int
	Height   int
	X        int
	Y        int
	Objects  [][]*Object
	BitChan  chan bool
	BiteChan []chan bool
	WallChan []chan bool
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
	m.BitChan = make(chan bool, 2)
	go g.handleBits(m)
}

func (m *GameMap) InitLevel3(g *Game) {
	bChan := make(chan bool, 2)
	m.BiteChan = append(m.BiteChan, bChan)
	go m.RandomBites(g, 1, 3, (20 * time.Second), false, m.BiteChan[0])
}

func (m *GameMap) InitLevel4(g *Game) {
	m.MakeWallChan(2)
	go m.MovingWall(g, 1+15, m.Height/4, DirLeft, 2, 15, WallRune, g.style.DefStyle, m.WallChan[0])
	go m.MovingWall(g, m.Width-15, (m.Height - m.Height/4), DirRight, 2, 15, WallRune, g.style.DefStyle, m.WallChan[1])
}

func (m *GameMap) InitLevel5(g *Game) {
	bChan := make(chan bool, 2)
	m.BiteChan = append(m.BiteChan, bChan)
	go m.RandomBites(g, 1, 3, (20 * time.Second), true, m.BiteChan[1])
}

func (m *GameMap) InitLevel6(g *Game) {
	m.BiteChan[0] <- true
	m.MakeWallChan(4)
	go m.MovingWall(g, m.Width/4, 6, DirUp, 1, 7, WallRune, g.style.DefStyle, m.WallChan[2])
	go m.MovingWall(g, (m.Width/4 + 1), 6, DirUp, 1, 7, WallRune, g.style.DefStyle, m.WallChan[3])
	go m.MovingWall(g, (m.Width - m.Width/4), m.Height-6, DirDown, 1, 7, WallRune, g.style.DefStyle, m.WallChan[4])
	go m.MovingWall(g, ((m.Width - m.Width/4) - 1), m.Height-6, DirDown, 1, 7, WallRune, g.style.DefStyle, m.WallChan[5])
}

func (m *GameMap) InitLevel7(g *Game) {
	for i := range m.WallChan {
		m.WallChan[i] <- true
	}
	m.BitChan <- true
}

func (m *GameMap) RandomLines(g *Game, numTimes int) {
	//for i := 0; i < numTimes; i++ {
	for {
		NewRandomBitLine(g, m, 10, BitRune, g.style.BitStyle)
		time.Sleep(15 * time.Second)
	}
}

func (m *GameMap) RandomBits(g *Game, bitsGen, bitsMax int, dur time.Duration) {
	for {
		for i := 0; i < bitsGen; i++ {
			if len(g.bits)-bitsGen < bitsMax {
				newB := NewRandomBit(m, 10, BitRune, g.style.BitStyle)
				g.bits = append(g.bits, newB)
			}
		}
		time.Sleep(dur)
	}
}

func (m *GameMap) RandomBites(g *Game, bitesGen, bitesMax int, dur time.Duration, random bool, biteChan chan bool) {
	for {
		select {
		default:
			for i := 0; i < bitesGen; i++ {
				if len(g.bites)-bitesGen < bitesMax {
					newB := NewRandomBite(m, g.style.BiteStyle, random)
					g.bites = append(g.bites, newB)
				}
			}
			time.Sleep(dur)
		case <-biteChan:
			return
		}
	}
}

func (m *GameMap) MovingWall(g *Game, x, y, direction, speed, segments int, char rune, style tcell.Style, quit chan bool) {
	e := NewEntity(x, y, direction, speed, char, style)
	e.AddSegment(segments, char, style)
	g.entities = append(g.entities, e)
	for {
		select {
		default:
			dx, dy := e.CheckDirection(g)
			if e.IsBlockedByMap(m, dx, dy) {
				var newPos []*Object
				for i := range e.pos {
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
		case <-quit:
			return
		}
	}
}

func (m *GameMap) MakeWallChan(num int) {
	for i := 0; i < num; i++ {
		c := make(chan bool, 2)
		m.WallChan = append(m.WallChan, c)
	}
}
