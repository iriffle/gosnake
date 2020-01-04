package main

import (
	"github.com/gdamore/tcell"
)

type Coord struct {
	x, y int
}

type Object struct {
	char  string
	layer int
	style tcell.Style
	x     int
	y     int
}

type Moving struct {
	object    *Object
	direction string
	speed     int
}

type Snake struct {
	moving *Moving
	pos    []Coord
}

type Player struct {
	snake *Snake
	name  string
	score int
}

func (m *Moving) Move(dx int, dy int) {
	// Move the Entity by the amount (dx, dy)
	m.object.x += dx
	m.object.y += dy
}
