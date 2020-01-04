package main

import (
	"github.com/gdamore/tcell"
)

type Coord struct {
	x, y int
}

type Entity struct {
	name  string
	char  string
	layer int
	style tcell.Style
	x     int
	y     int
	pos   []Coord
}

func (e *Entity) Move(dx int, dy int) {
	// Move the Entity by the amount (dx, dy)
	e.x += dx
	e.y += dy
}
