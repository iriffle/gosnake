package main

import (
	"github.com/gdamore/tcell"
)

// Entity struct
type Entity struct {
	pos       []Object
	direction int
}

// Create a new Entity
func NewEntity(x, y, direction int, char rune, style tcell.Style) Entity {
	o := NewObject(x, y, char, style, true)
	e := Entity{
		direction: direction,
	}
	e.pos = append(e.pos, o)
	return e
}

// Move the entity's segments
func (e *Entity) Move(dx, dy int) {
	first := true
	e.pos[0].ox = e.pos[0].x
	e.pos[0].oy = e.pos[0].y
	e.pos[0].x += dx
	e.pos[0].y += dy

	for i, _ := range e.pos {
		if !first {
			e.pos[i].ox = e.pos[i].x
			e.pos[i].oy = e.pos[i].y
			e.pos[i].x = e.pos[i-1].ox
			e.pos[i].y = e.pos[i-1].oy
		} else {
			first = false
		}
	}
}

// Add a segment to the entity
func (e *Entity) AddSegment(char rune, style tcell.Style) {
	x := e.pos[len(e.pos)-1].ox
	y := e.pos[len(e.pos)-1].oy
	o := NewObject(x, y, char, style, true)
	e.pos = append(e.pos, o)
}
