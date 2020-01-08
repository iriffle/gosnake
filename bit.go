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

func NewRandomBit(mapStartX, mapStartY, mapWidth, mapHeight int, points int, char rune, style tcell.Style) Bit {
	randX := rand.Intn(mapWidth)
	randY := rand.Intn(mapHeight)
	if randX == mapWidth {
		randX = randX - 1
	}
	if randY == mapHeight {
		randY = randY - 1
	} else if randY == 1 {
		randY = randY + 1
	}
	b := NewBit(randX, randY, points, char, style)
	return b
}
