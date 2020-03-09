package main

import (
	"github.com/gdamore/tcell"
)

// Entity struct
type Entity struct {
	pos       []*Object
	direction int
	speed     int
}

// Create a new Entity
func NewEntity(x, y, direction, speed int, char rune, style tcell.Style) *Entity {
	o := NewObject(x, y, char, style, true)
	e := Entity{
		direction: direction,
		speed:     speed,
	}
	e.pos = append(e.pos, o)
	return &e
}

func NewDisplayEntity(w, h, size, x, y int, char rune, style tcell.Style) *Entity {
	e := NewEntity(w, h, DirAll, 0, char, style)
	for i := 0; i < size; i++ {
		e.pos[i].ox += x
		e.pos[i].oy += y
		e.AddSegment(1, e.pos[0].char, e.pos[0].style)
	}
	return e
}

func NewColorEntity(w, h int, char rune, style tcell.Style) *Entity {
	e := NewEntity(w, h, DirAll, 0, char, style)
	for i := 0; i < len(PlayerColors); i++ {
		style := StringToStyle(PlayerColors[i], PlayerColors[1])
		e.pos[i].oy++
		e.pos[i].style = style
		if i < len(PlayerColors)-1 {
			e.AddSegment(1, e.pos[0].char, style)
		}
	}
	return e
}

// Move the entity's segments
func (e *Entity) Move(dx, dy int) {
	first := true
	e.pos[0].ox = e.pos[0].x
	e.pos[0].oy = e.pos[0].y
	e.pos[0].x += dx
	e.pos[0].y += dy

	for i := range e.pos {
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

// Get Entity's current direction and return dx, dy
// in order to change Entity's movement if change necessary
func (e *Entity) CheckDirection(g *Game) (int, int) {
	dx, dy := 0, 0
	switch e.direction {
	case DirUp:
		dy--
	case DirDown:
		dy++
	case DirLeft:
		dx--
	case DirRight:
		dx++
	}

	return dx, dy
}

func (e *Entity) GetDirection() int {
	return e.direction
}

// Add a segment to the entity
func (e *Entity) AddSegment(num int, char rune, style tcell.Style) {
	for i := 0; i < num; i++ {
		x := e.pos[len(e.pos)-1].ox
		y := e.pos[len(e.pos)-1].oy
		o := NewObject(x, y, char, style, true)
		e.pos = append(e.pos, o)
	}
}

func (e *Entity) RemoveSegment(num int) {
	for i := 0; i < num; i++ {
		e.pos[len(e.pos)-1] = nil
		e.pos = e.pos[:len(e.pos)-1]
	}
}

// Check if player is blocked by an object on the map
func (e *Entity) IsBlockedByMap(m *GameMap, dx, dy int) bool {
	if m.Objects[e.pos[0].x+dx][e.pos[0].y+dy].blocked {
		return true
	}
	return false
}

func (e *Entity) SetChar(char rune) {
	for i := range e.pos {
		e.pos[i].char = char
	}
}

func (e *Entity) SetStyle(style tcell.Style) {
	for i := range e.pos {
		e.pos[i].style = style
	}
}

func (e *Entity) RotateDisplay(entities []*Entity, rotation int) {
	e.SetChar(entities[rotation].pos[0].char)
	e.SetStyle(entities[rotation].pos[0].style)
}
