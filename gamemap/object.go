package gamemap

import (
	"github.com/gdamore/tcell"
)

// Object struct
type Object struct {
	x, y, ox, oy int
	char         rune
	style        tcell.Style
	blocked      bool
}

// Create new Object
func NewObject(x, y int, char rune, style tcell.Style, blocked bool) *Object {
	o := Object{
		x,
		y,
		x,
		y,
		char,
		style,
		blocked}
	return &o
}

func (o *Object) GetCurPos() (int, int) {
	return o.x, o.y
}

func (o *Object) GetLastPos() (int, int) {
	return o.ox, o.oy
}

// Move Object
func (o *Object) Move(dx, dy int) {
	o.ox = o.x
	o.oy = o.y
	o.x += dx
	o.y += dy
}

func (o *Object) MoveCurPos(dx, dy int) {
	o.x += dx
	o.y += dy
}

func (o *Object) MoveLastPos(dx, dy int) {
	o.ox += dx
	o.oy += dy
}

func (o *Object) SetPos(x, y int) {
	o.x = x
	o.y = y
	o.ox = x
	o.oy = y
}

func (o *Object) GetChar() rune {
	return o.char
}

func (o *Object) SetChar(char rune) {
	o.char = char
}

func (o *Object) GetStyle() tcell.Style {
	return o.style
}

func (o *Object) SetStyle(style tcell.Style) {
	o.style = style
}

// Check if Object is blocked
func (o *Object) IsBlocked() bool {
	// Check to see if the provided coordinates contain a blocked tile
	if o.blocked {
		return true
	}
	return false
}

func (o *Object) Block() {
	o.blocked = true
}

func (o *Object) Unblock() {
	o.blocked = false
}
