package main

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

// Move Object
func (o *Object) Move(dx, dy int) {
	o.x += dx
	o.y += dy
}

// Check if Object is blocked
func (o *Object) IsBlocked(x int, y int) bool {
	// Check to see if the provided coordinates contain a blocked tile
	if o.x == x && o.y == y && o.blocked {
		return true
	}
	return false
}
