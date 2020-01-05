package entities

import (
	"github.com/gdamore/tcell"
)

type Object struct {
	X, Y, OX, OY int
	Char         rune
	Style        tcell.Style
}

func NewObject(x, y int, char rune, style tcell.Style) Object {
	o := Object{
		x,
		y,
		x,
		y,
		char,
		style}
	return o
}

func (o *Object) MoveObject(dx, dy int) {
	o.X += dx
	o.Y += dy
}
