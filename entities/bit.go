package entities

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

type Bit struct {
	Object
	Points int
}

func NewBit(x, y, points int, char rune, style tcell.Style) Bit {
	o := NewObject(x, y, char, style)
	b := Bit{
		o,
		points,
	}
	return b
}

func SetBit(MapStartX, MapStartY, MapWidth, MapHeight int, points int, char rune, style tcell.Style) Bit {
	rand.Seed(time.Now().UnixNano())
	minX := MapStartX
	maxX := minX + MapWidth
	minY := MapStartY
	maxY := minY + MapHeight
	randX := rand.Intn(maxX-minX+1) + minX
	randY := rand.Intn(maxY-minY+1) + minY
	b := NewBit(randX, randY, points, char, style)
	return b
}
