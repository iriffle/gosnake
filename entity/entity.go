package entity

import (
	"github.com/gdamore/tcell"
	"github.com/stjiub/gosnake/gamemap"
	"github.com/stjiub/gosnake/style"
)

// Entity struct
type Entity struct {
	pos       []*gamemap.Object
	direction int
	speed     int
}

// Create a new Entity
func NewEntity(x, y, direction, speed int, char rune, sty tcell.Style) *Entity {
	o := gamemap.NewObject(x, y, char, sty, true)
	e := Entity{
		direction: direction,
		speed:     speed,
	}
	e.pos = append(e.pos, o)
	return &e
}

func NewDisplayEntity(w, h, size, x, y int, char rune, sty tcell.Style) *Entity {
	e := NewEntity(w, h, DirAll, 0, char, sty)
	for i := 0; i < size; i++ {
		e.pos[i].MoveLastPos(x, y)
		e.AddSegment(1, char, sty)
	}
	return e
}

func NewColorEntity(w, h int, char rune, colors []string, sty tcell.Style) *Entity {
	e := NewEntity(w, h, DirAll, 0, char, sty)
	for i := 0; i < len(colors); i++ {
		sty := style.StringToStyle(colors[i], colors[1])
		e.pos[i].MoveLastPos(0, 1)
		e.pos[i].SetStyle(sty)
		if i < len(colors)-1 {
			e.AddSegment(1, char, sty)
		}
	}
	return e
}

// Move the entity's segments
func (e *Entity) Move(dx, dy int) {
	first := true
	e.pos[0].Move(dx, dy)

	for i := range e.pos {
		if !first {
			x, y := e.pos[i-1].GetLastPos()
			e.pos[i].SetPos(x, y)
		} else {
			first = false
		}
	}
}

// Get Entity's current direction and return dx, dy
// in order to change Entity's movement if change necessary
func (e *Entity) CheckDirection() (int, int) {
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

// Add a segment to the entity
func (e *Entity) AddSegment(num int, char rune, sty tcell.Style) {
	for i := 0; i < num; i++ {
		x, y := e.pos[len(e.pos)-1].GetLastPos()
		o := gamemap.NewObject(x, y, char, sty, true)
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
func (e *Entity) IsBlockedByMap(m *gamemap.GameMap, dx, dy int) bool {
	x, y := e.pos[0].GetCurPos()
	if m.Objects[x+dx][y+dy].IsBlocked() {
		return true
	}
	return false
}

func (e *Entity) SetChar(char rune) {
	for i := range e.pos {
		e.pos[i].SetChar(char)
	}
}

func (e *Entity) SetStyle(style tcell.Style) {
	for i := range e.pos {
		e.pos[i].SetStyle(style)
	}
}

func (e *Entity) RotateDisplay(entities []*Entity, rotation int) {
	char := entities[rotation].pos[0].GetChar()
	style := entities[rotation].pos[0].GetStyle()
	e.SetChar(char)
	e.SetStyle(style)
}

func (e *Entity) GetDirection() int {
	return e.direction
}

func (e *Entity) SetDirection(dir int) {
	e.direction = dir
}

func (e *Entity) GetSpeed() int {
	return e.speed
}

func (e *Entity) GetLength() int {
	return len(e.pos)
}

func (e *Entity) GetSegment(i int) *gamemap.Object {
	return e.pos[i]
}

func (e *Entity) GetCurPos(i int) (int, int) {
	x, y := e.pos[i].GetCurPos()
	return x, y
}

func (e *Entity) NewPos(newPos []*gamemap.Object) {
	e.pos = newPos
}

func (e *Entity) GetChar(i int) rune {
	return e.pos[i].GetChar()
}

func (e *Entity) GetStyle(i int) tcell.Style {
	return e.pos[i].GetStyle()
}
