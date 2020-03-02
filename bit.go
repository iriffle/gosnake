package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

const (

	// Bit states
	BitStatic = iota
	BitMoving
	BitRandom
)

// Bit struct
type Bit struct {
	*Object
	points int
	state  int
}

type Bite struct {
	*Bit
	dir int
}

type Bits interface {
	CheckPos()
}

// Create new Bit
func NewBit(x, y, points int, char rune, state int, style tcell.Style) *Bit {
	o := NewObject(x, y, char, style, false)
	b := Bit{
		o,
		points,
		state,
	}
	return &b
}

// Create a horizontal line of bits of a given length
func NewBitLineH(g *Game, x, y, points, numBits int, char rune, style tcell.Style) {
	for i := 0; i < numBits; i++ {
		x += 2
		b := NewBit(x, y, points, char, 0, style)
		g.bits = append(g.bits, b)
	}
}

// Create a vertical line of bits of a given length
func NewBitLineV(g *Game, x, y, points, numBits int, char rune, style tcell.Style) {
	for i := 0; i < numBits; i++ {
		y += 1
		b := NewBit(x, y, points, char, 0, style)
		g.bits = append(g.bits, b)
	}
}

// Generate random coordinates for a Bit
func NewRandomBit(m *GameMap, points int, char rune, style tcell.Style) *Bit {
	var b *Bit
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

// Generate random coordinates for a Bit line
func NewRandomBitLine(g *Game, m *GameMap, points int, char rune, style tcell.Style) {
	for {
		randNum := rand.Intn(6) + 2
		randDir := randBool()
		randX := rand.Intn(m.Width)
		randY := rand.Intn(m.Height)
		if randDir {
			if randX < ((m.Width-1)-(randNum*2)) && randX > 1 && randY < m.Height-1 && randY > 1 {
				NewBitLineH(g, randX, randY, points, randNum, char, style)
				break
			}
		} else {
			if randX < m.Width-1 && randX > 1 && randY < ((m.Height-1)-randNum) && randY > 1 {
				NewBitLineV(g, randX, randY, points, randNum, char, style)
				break
			}
		}
	}
}

// Move a Bit in random direction
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

// NewBite creates a new Bite object.
func NewBite(m *GameMap, x, y, points, dir, state int, char rune, style tcell.Style) *Bite {
	bit := NewBit(x, y, points, char, state, style)
	bite := Bite{
		bit,
		dir,
	}
	return &bite
}

// NewRandomBite generates random coordinates and random explosion directions for a new Bite.
func NewRandomBite(m *GameMap, style tcell.Style, random bool) *Bite {
	var (
		bite *Bite
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

// ExplodeBite triggers a bite explosion based on bite direction type.
func (b *Bite) ExplodeBite(g *Game, m, biteMap *GameMap) {
	b.style = g.style.BiteExplodedStyle
	time.Sleep(500 * time.Millisecond)
	b.ExplodeDir(biteMap, m, BiteExplodeRune, g.style.BiteExplodedStyle, true, (30 * time.Millisecond))
	time.Sleep(10 * time.Second)
	b.ExplodeDir(biteMap, m, ' ', g.style.DefStyle, false, 0)

	for i, _ := range g.bites {
		if b.x == g.bites[i].x && b.y == g.bites[i].y {
			g.bites[i] = g.bites[len(g.bites)-1]
			g.bites[len(g.bites)-1] = nil
			g.bites = g.bites[:len(g.bites)-1]
			return
		}
	}
}

func (b *Bite) ExplodeDir(biteMap, m *GameMap, char rune, style tcell.Style, blocked bool, t time.Duration) {
	if b.dir == DirUp || b.dir == DirAll {
		go b.SetUp(biteMap, m, char, style, blocked, t)
	}
	if b.dir == DirDown || b.dir == DirAll {
		go b.SetDown(biteMap, m, char, style, blocked, t)
	}
	if b.dir == DirLeft || b.dir == DirAll {
		go b.SetLeft(biteMap, m, char, style, blocked, t)
	}
	if b.dir == DirRight || b.dir == DirAll {
		go b.SetRight(biteMap, m, char, style, blocked, t)
	}
}

// SetRight sets or clears an explosion to the right.
func (b *Bite) SetRight(biteMap, m *GameMap, char rune, style tcell.Style, blocked bool, t time.Duration) {
	for x := b.x + 1; x < m.Width-1; x++ {
		time.Sleep(t)
		SetObject(biteMap, x, b.y, char, style, blocked)
	}
}

// SetLeft sets or clears an explosion to the left.
func (b *Bite) SetLeft(biteMap, m *GameMap, char rune, style tcell.Style, blocked bool, t time.Duration) {
	for x := b.x - 1; x > 1; x-- {
		time.Sleep(t)
		SetObject(biteMap, x, b.y, char, style, blocked)
	}
}

// SetDown sets or clears an explosion down.
func (b *Bite) SetDown(biteMap, m *GameMap, char rune, style tcell.Style, blocked bool, t time.Duration) {
	for y := b.y + 1; y < m.Height-1; y++ {
		time.Sleep(t)
		SetObject(biteMap, b.x, y, char, style, blocked)
	}
}

// SetUp sets or clears and explosion up.
func (b *Bite) SetUp(biteMap, m *GameMap, char rune, style tcell.Style, blocked bool, t time.Duration) {
	for y := b.y - 1; y > 0; y-- {
		time.Sleep(t)
		SetObject(biteMap, b.x, y, char, style, blocked)
	}
}

// SetObject changes the state of an object on the biteMap.
func SetObject(biteMap *GameMap, x, y int, char rune, style tcell.Style, blocked bool) {
	biteMap.Objects[x][y].char = char
	biteMap.Objects[x][y].style = style
	biteMap.Objects[x][y].blocked = blocked
}

// randBool generates a random boolean output.
func randBool() bool {
	return rand.Uint64()&(1<<63) == 0
}
