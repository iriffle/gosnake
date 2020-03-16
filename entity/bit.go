package entity

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
	"github.com/stjiub/gosnake/gamemap"
)

const (

	// Bit states
	BitStatic = iota
	BitMoving
	BitRandom
	Bite
)

const (
	DirUp = iota
	DirDown
	DirLeft
	DirRight
	DirAll
	DirNone
)

// Bit struct
type Bit struct {
	*gamemap.Object
	points int
	state  int
	dir    int
}

type Bits interface {
	CheckPos()
}

// Create new Bit
func NewBit(x, y, points int, char rune, state, dir int, style tcell.Style) *Bit {
	o := gamemap.NewObject(x, y, char, style, false)
	b := Bit{
		o,
		points,
		state,
		dir,
	}
	return &b
}

// Create a horizontal line of bits of a given length
func NewBitLineH(bits []*Bit, x, y, points, numBits int, char rune, style tcell.Style) []*Bit {
	for i := 0; i < numBits; i++ {
		x += 2
		b := NewBit(x, y, points, char, 0, DirNone, style)
		bits = append(bits, b)
	}
	return bits
}

// Create a vertical line of bits of a given length
func NewBitLineV(bits []*Bit, x, y, points, numBits int, char rune, style tcell.Style) []*Bit {
	for i := 0; i < numBits; i++ {
		y++
		b := NewBit(x, y, points, char, 0, DirNone, style)
		bits = append(bits, b)
	}
	return bits
}

// Generate random coordinates for a Bit
func NewRandomBit(m *gamemap.GameMap, points int, char rune, style tcell.Style) *Bit {
	var b *Bit
	for {
		randX := rand.Intn(m.Width)
		randY := rand.Intn(m.Height)
		if randX < m.Width-1 && randX > 1 && randY < m.Height-1 && randY > 1 {
			b = NewBit(randX, randY, points, char, 2, DirNone, style)
			break
		}
	}
	return b
}

// Generate random coordinates for a Bit line
func NewRandomBitLine(bits []*Bit, m *gamemap.GameMap, points int, char rune, style tcell.Style) []*Bit {
	for {
		randNum := rand.Intn(6) + 2
		randDir := randBool()
		randX := rand.Intn(m.Width)
		randY := rand.Intn(m.Height)
		if randDir {
			if randX < ((m.Width-1)-(randNum*2)) && randX > 1 && randY < m.Height-1 && randY > 1 {
				bits = NewBitLineH(bits, randX, randY, points, randNum, char, style)
				return bits
			}
		} else {
			if randX < m.Width-1 && randX > 1 && randY < ((m.Height-1)-randNum) && randY > 1 {
				bits = NewBitLineV(bits, randX, randY, points, randNum, char, style)
				return bits
			}
		}
		return bits
	}
}

func (b *Bit) GetState() int {
	return b.state
}

func (b *Bit) GetPoints() int {
	return b.points
}

// Move a Bit in random direction
func (b *Bit) MoveRandom(m *gamemap.GameMap) {
	r := [2]int{0, 0}
	for i := range r {
		random := randBool()
		if random {
			r[i] = 1
		}
		random = randBool()
		if random {
			r[i] -= (r[i] * 2)
		}
	}
	bx, by := b.GetCurPos()
	if !m.Objects[r[0]+bx][r[1]+by].IsBlocked() {
		b.Move(r[0], r[1])
	}
}

// NewRandomBite generates random coordinates and random explosion directions for a new Bite.
func NewRandomBite(m *gamemap.GameMap, runes []rune, style tcell.Style, random bool) *Bit {
	var (
		bite *Bit
		dir  int
		char rune
	)

	for {
		if random {
			randDir := rand.Intn(4)
			switch randDir {
			case DirUp:
				dir = DirUp
				char = runes[DirUp]
			case DirDown:
				dir = DirDown
				char = runes[DirDown]
			case DirLeft:
				dir = DirLeft
				char = runes[DirLeft]
			case DirRight:
				dir = DirRight
				char = runes[DirRight]
			case DirAll:
				dir = DirAll
				char = runes[DirAll]
			}
		} else {
			dir = DirAll
			char = runes[DirAll]
		}
		randX := rand.Intn(m.Width)
		randY := rand.Intn(m.Height)
		if randX < m.Width-1 && randX > 1 && randY < m.Height-1 && randY > 1 {
			bite = NewBit(randX, randY, 50, char, BitStatic, dir, style)
			break
		}
	}
	return bite

}

// ExplodeBite triggers a bite explosion based on bite direction type.
func (b *Bit) ExplodeBite(m, biteMap *gamemap.GameMap, biteExplodeRune rune, explodedStyle, defStyle tcell.Style) {
	b.SetStyle(explodedStyle)
	time.Sleep(500 * time.Millisecond)
	b.ExplodeDir(biteMap, m, biteExplodeRune, explodedStyle, true, (30 * time.Millisecond))
	time.Sleep(10 * time.Second)
	b.ExplodeDir(biteMap, m, ' ', defStyle, false, 0)
}

func (b *Bit) ExplodeDir(biteMap, m *gamemap.GameMap, char rune, style tcell.Style, blocked bool, t time.Duration) {
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
func (b *Bit) SetRight(biteMap, m *gamemap.GameMap, char rune, style tcell.Style, blocked bool, t time.Duration) {
	bx, by := b.GetCurPos()
	for x := bx + 1; x < m.Width-1; x++ {
		time.Sleep(t)
		SetObject(biteMap, x, by, char, style, blocked)
	}
}

// SetLeft sets or clears an explosion to the left.
func (b *Bit) SetLeft(biteMap, m *gamemap.GameMap, char rune, style tcell.Style, blocked bool, t time.Duration) {
	bx, by := b.GetCurPos()
	for x := bx - 1; x > 1; x-- {
		time.Sleep(t)
		SetObject(biteMap, x, by, char, style, blocked)
	}
}

// SetDown sets or clears an explosion down.
func (b *Bit) SetDown(biteMap, m *gamemap.GameMap, char rune, style tcell.Style, blocked bool, t time.Duration) {
	bx, by := b.GetCurPos()
	for y := by + 1; y < m.Height-1; y++ {
		time.Sleep(t)
		SetObject(biteMap, bx, y, char, style, blocked)
	}
}

// SetUp sets or clears and explosion up.
func (b *Bit) SetUp(biteMap, m *gamemap.GameMap, char rune, style tcell.Style, blocked bool, t time.Duration) {
	bx, by := b.GetCurPos()
	for y := by - 1; y > 0; y-- {
		time.Sleep(t)
		SetObject(biteMap, bx, y, char, style, blocked)
	}
}

// SetObject changes the state of an object on the biteMap.
func SetObject(biteMap *gamemap.GameMap, x, y int, char rune, style tcell.Style, blocked bool) {
	biteMap.Objects[x][y].SetChar(char)
	biteMap.Objects[x][y].SetStyle(style)
	biteMap.Objects[x][y].Block()
}

// randBool generates a random boolean output.
func randBool() bool {
	return rand.Uint64()&(1<<63) == 0
}
