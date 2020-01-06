package main

import (
	"math/rand"
	"time"

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

func NewRandomBit(MapStartX, MapStartY, MapWidth, MapHeight int, points int, char rune, style tcell.Style) Bit {
	rand.Seed(time.Now().UnixNano())
	minX := MapStartX + 2
	maxX := minX + MapWidth - 2
	minY := MapStartY + 2
	maxY := minY + MapHeight - 2
	randX := rand.Intn(maxX - minX)
	randY := rand.Intn(maxY - minY)
	b := NewBit(randX, randY, points, char, style)
	return b
}
