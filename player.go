package main

import (
	"github.com/gdamore/tcell"
)

type Player struct {
	Entity
	name  string
	score int
}

func NewPlayer(x, y, score, direction int, char rune, name string, style tcell.Style) Player {
	e := NewEntity(x, y, direction, char, style)
	p := Player{
		e,
		name,
		score}
	return p
}
