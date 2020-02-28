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
func (e *Entity) AddSegment(char rune, style tcell.Style) {
	x := e.pos[len(e.pos)-1].ox
	y := e.pos[len(e.pos)-1].oy
	o := NewObject(x, y, char, style, true)
	e.pos = append(e.pos, o)
}

// Check if player is blocked by an object on the map
func (e *Entity) IsBlockedByMap(m *GameMap, dx, dy int) bool {
	if m.Objects[e.pos[0].x+dx][e.pos[0].y+dy].blocked {
		return true
	} else {
		return false
	}
}
