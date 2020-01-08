package main

import (
	"math/rand"

	"github.com/gdamore/tcell"
)

type Bit struct {
	Object
	points int
}

func NewBit(x, y, points int, char rune, style tcell.Style) Bit {
	o := NewObject(x, y, char, style, false)
	b := Bit{
		o,
		points,
	}
	return b
}

func NewRandomBit(m *GameMap, points int, char rune, style tcell.Style) Bit {
	var b Bit
	for {
		randX := rand.Intn(m.Width)
		randY := rand.Intn(m.Height)
		if randX < m.Width-1 && randX > 1 && randY < m.Height-1 && randY > 1 {
			b = NewBit(randX, randY, points, char, style)
			break
		}
	}
	return b
}
