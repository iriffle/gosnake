package entities

import (
	"github.com/gdamore/tcell"
)

type Entity struct {
	Pos []Object
}

func NewEntity(x, y int, char rune, style tcell.Style) Entity {
	o := NewObject(x, y, char, style)
	e := Entity{}
	e.Pos = append(e.Pos, o)
	return e
}

func (e *Entity) MoveEntity(dx, dy int) {
	first := true
	e.Pos[0].OX = e.Pos[0].X
	e.Pos[0].OY = e.Pos[0].Y
	e.Pos[0].X += dx
	e.Pos[0].Y += dy

	for i, _ := range e.Pos {
		if !first {
			e.Pos[i].OX = e.Pos[i].X
			e.Pos[i].OY = e.Pos[i].Y
			e.Pos[i].X = e.Pos[i-1].OX
			e.Pos[i].Y = e.Pos[i-1].OY
		} else {
			first = false
		}
	}
}

func (e *Entity) AddSegment(char rune, style tcell.Style) {
	x := e.Pos[len(e.Pos)-1].OX
	y := e.Pos[len(e.Pos)-1].OY
	o := NewObject(x, y, char, style)
	e.Pos = append(e.Pos, o)
}
