package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

const (
	BitStatic = iota
	BitMoving = iota
	BitRandom = iota
)

// Bit struct
type Bit struct {
	Object
	points int
	state  int
}

type Bite struct {
	Bit
	dir int
}

type Bits interface {
	CheckPos()
}

// Create new Bit
func NewBit(x, y, points int, char rune, state int, style tcell.Style) Bit {
	o := NewObject(x, y, char, style, false)
	b := Bit{
		o,
		points,
		state,
	}
	return b
}

func NewBitLineH(g *Game, x, y, points, numBits int, char rune, style tcell.Style) {
	for i := 0; i < numBits; i++ {
		x += 2
		b := NewBit(x, y, points, char, 0, style)
		g.bits = append(g.bits, &b)
	}
}

func NewBitLineV(g *Game, x, y, points, numBits int, char rune, style tcell.Style) {
	for i := 0; i < numBits; i++ {
		y += 1
		b := NewBit(x, y, points, char, 0, style)
		g.bits = append(g.bits, &b)
	}
}

// Generate random coordinates for a Bit
func NewRandomBit(m *GameMap, points int, char rune, style tcell.Style) Bit {
	var b Bit
	for {
		randX := rand.Intn(m.Width)
		randY := rand.Intn(m.Height)
		if randX < m.Width-1 && randX > 1 && randY < m.Height-1 && randY > 1 {
			b = NewBit(randX, randY, points, char, 2, style)
			break
		}
	}
	return b
}

// Generate random coordinates for a Bit
func NewRandomBitLine(g *Game, m *GameMap, points int, char rune, style tcell.Style) {
	for {
		randNum := rand.Intn(6)
		randDir := randBool()
		randX := rand.Intn(m.Width)
		randY := rand.Intn(m.Height)
		if randDir {
			if randX < m.Width-1-randNum && randX > 1 && randY < m.Height-1 && randY > 1 {
				NewBitLineH(g, randX, randY, points, randNum, char, style)
				break
			}
		} else {
			if randX < m.Width-1 && randX > 1 && randY < m.Height-1-randNum && randY > 1 {
				NewBitLineH(g, randX, randY, points, randNum, char, style)
				break
			}
		}
	}
}

// Move Bit in random direction
func (b *Bit) Move(m *GameMap) {
	r := [2]int{0, 0}
	for i, _ := range r {
		random := randBool()
		if random {
			r[i] = 1
		}
		random = randBool()
		if random {
			r[i] -= (r[i] * 2)
		}
	}
	if !m.Objects[r[0]+b.x][r[1]+b.y].blocked {
		b.x += r[0]
		b.y += r[1]
	}
}

func NewBite(m *GameMap, x, y, points, dir, state int, char rune, style tcell.Style) Bite {
	bit := NewBit(x, y, points, char, state, style)
	bite := Bite{
		bit,
		dir,
	}
	return bite
}

func NewRandomBite(m *GameMap, style tcell.Style, random bool) Bite {
	var (
		bite Bite
		dir  int
		char rune
	)

	for {
		if random {
			randDir := rand.Intn(4)
			switch randDir {
			case DirUp:
				dir = DirUp
				char = BiteUpRune
			case DirDown:
				dir = DirDown
				char = BiteDownRune
			case DirLeft:
				dir = DirLeft
				char = BiteLeftRune
			case DirRight:
				dir = DirRight
				char = BiteRightRune
			case DirAll:
				dir = DirAll
				char = BiteAllRune
			}
		} else {
			dir = DirAll
			char = BiteAllRune
		}
		randX := rand.Intn(m.Width)
		randY := rand.Intn(m.Height)
		if randX < m.Width-1 && randX > 1 && randY < m.Height-1 && randY > 1 {
			bite = NewBite(m, randX, randY, 50, dir, BitStatic, char, style)
			break
		}
	}
	return bite

}

func (b *Bite) ExplodeBite(g *Game, m *GameMap) {
	b.style = BiteExplodedStyle
	biteMap := &GameMap{
		Width:  m.Width,
		Height: m.Height,
	}
	biteMap.InitMap()
	biteMap.InitMapBoundary(WallRune, FloorRune, DefStyle)
	g.maps = append(g.maps, biteMap)
	time.Sleep(500 * time.Millisecond)
	if b.dir == DirUp || b.dir == DirAll {
		go b.ExplodeYUp(biteMap, m)
	}
	if b.dir == DirDown || b.dir == DirAll {
		go b.ExplodeYDown(biteMap, m)
	}
	if b.dir == DirLeft || b.dir == DirAll {
		go b.ExplodeXLeft(biteMap, m)
	}
	if b.dir == DirRight || b.dir == DirAll {
		go b.ExplodeXRight(biteMap, m)
	}

	time.Sleep(10 * time.Second)
	emptyMap := &GameMap{}
	i := len(g.maps) - 1
	g.maps[i] = emptyMap
	g.maps = g.maps[:len(g.maps)-1]
	b = &Bite{}
	i = len(g.bites) - 1
	g.bites[i] = b
	g.bites = g.bites[:len(g.bites)-1]

}

func (b *Bite) ExplodeXRight(biteMap, m *GameMap) {
	for x := b.x + 1; x < m.Width-1; x++ {
		time.Sleep(30 * time.Millisecond)
		biteMap.Objects[x][b.y].char = BiteExplodeRune
		biteMap.Objects[x][b.y].style = BiteExplodedStyle
		biteMap.Objects[x][b.y].blocked = true
	}
}
func (b *Bite) ExplodeXLeft(biteMap, m *GameMap) {
	for x := b.x - 1; x > 1; x-- {
		time.Sleep(30 * time.Millisecond)
		biteMap.Objects[x][b.y].char = BiteExplodeRune
		biteMap.Objects[x][b.y].style = BiteExplodedStyle
		biteMap.Objects[x][b.y].blocked = true
	}
}

func (b *Bite) ExplodeYDown(biteMap, m *GameMap) {
	for y := b.y + 1; y < m.Height-1; y++ {
		time.Sleep(30 * time.Millisecond)
		biteMap.Objects[b.x][y].char = BiteExplodeRune
		biteMap.Objects[b.x][y].style = BiteExplodedStyle
		biteMap.Objects[b.x][y].blocked = true
	}
}

func (b *Bite) ExplodeYUp(biteMap, m *GameMap) {
	for y := b.y - 1; y > 0; y-- {
		time.Sleep(30 * time.Millisecond)
		biteMap.Objects[b.x][y].char = BiteExplodeRune
		biteMap.Objects[b.x][y].style = BiteExplodedStyle
		biteMap.Objects[b.x][y].blocked = true
	}
}

// Generate random boolean output
func randBool() bool {
	return rand.Uint64()&(1<<63) == 0
}
