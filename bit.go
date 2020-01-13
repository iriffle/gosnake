package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

// Bit struct
type Bit struct {
	Object
	points int
	moving bool
}

// Create new Bit
func NewBit(x, y, points int, char rune, moving bool, style tcell.Style) Bit {
	o := NewObject(x, y, char, style, false)
	b := Bit{
		o,
		points,
		moving,
	}
	return b
}

func NewBitLineH(g *Game, x, y, points, numBits int, char rune, style tcell.Style) {
	for i := 0; i < numBits; i++ {
		x += 2
		b := NewBit(x, y, points, char, false, style)
		g.bits = append(g.bits, &b)
	}
}

func NewBitLineV(g *Game, x, y, points, numBits int, char rune, style tcell.Style) {
	for i := 0; i < numBits; i++ {
		y += 1
		b := NewBit(x, y, points, char, false, style)
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
			b = NewBit(randX, randY, points, char, true, style)
			break
		}
	}
	return b
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

func NewRandomBite(m *GameMap, char rune, style tcell.Style) Bit {
	var bite Bit
	for {
		randX := rand.Intn(m.Width)
		randY := rand.Intn(m.Height)
		if randX < m.Width-1 && randX > 1 && randY < m.Height-1 && randY > 1 {
			bite = NewBit(randX, randY, 50, char, false, style)
			break
		}
	}
	return bite

}

func (b *Bit) ExplodeBite(g *Game, m *GameMap) {
	b.style = BiteExplodedStyle
	biteMap := &GameMap{
		Width:  m.Width,
		Height: m.Height,
	}
	biteMap.InitMap()
	biteMap.InitMapBoundary(wallRune, floorRune, DefStyle)
	g.maps = append(g.maps, biteMap)
	time.Sleep(500 * time.Millisecond)
	go b.ExplodeXRight(biteMap, m)
	go b.ExplodeXLeft(biteMap, m)
	go b.ExplodeYUp(biteMap, m)
	go b.ExplodeYDown(biteMap, m)
	time.Sleep(10 * time.Second)
	emptyMap := &GameMap{}
	i := len(g.maps) - 1
	g.maps[i] = emptyMap
	g.maps = g.maps[:len(g.maps)-1]
	b = &Bit{}
	i = len(g.bites) - 1
	g.bites[i] = b
	g.bites = g.bites[:len(g.bites)-1]

}

func (b *Bit) ExplodeXRight(biteMap, m *GameMap) {
	for x := b.x + 1; x < m.Width-1; x++ {
		time.Sleep(30 * time.Millisecond)
		biteMap.Objects[x][b.y].char = b.char
		biteMap.Objects[x][b.y].style = b.style
		biteMap.Objects[x][b.y].blocked = true
	}
}
func (b *Bit) ExplodeXLeft(biteMap, m *GameMap) {
	for x := b.x - 1; x > 1; x-- {
		time.Sleep(30 * time.Millisecond)
		biteMap.Objects[x][b.y].char = b.char
		biteMap.Objects[x][b.y].style = b.style
		biteMap.Objects[x][b.y].blocked = true
	}
}

func (b *Bit) ExplodeYDown(biteMap, m *GameMap) {
	for y := b.y + 1; y < m.Height-1; y++ {
		time.Sleep(30 * time.Millisecond)
		biteMap.Objects[b.x][y].char = b.char
		biteMap.Objects[b.x][y].style = b.style
		biteMap.Objects[b.x][y].blocked = true
	}
}

func (b *Bit) ExplodeYUp(biteMap, m *GameMap) {
	for y := b.y - 1; y > 0; y-- {
		time.Sleep(30 * time.Millisecond)
		biteMap.Objects[b.x][y].char = b.char
		biteMap.Objects[b.x][y].style = b.style
		biteMap.Objects[b.x][y].blocked = true
	}
}

// Generate random boolean output
func randBool() bool {
	return rand.Uint64()&(1<<63) == 0
}
