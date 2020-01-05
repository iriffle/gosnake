package entities

import (
	"github.com/gdamore/tcell"
)

type Player struct {
	Entity
	Name  string
	Score int
}

func NewPlayer(x, y, score int, char rune, name string, style tcell.Style) Player {
	e := NewEntity(x, y, char, style)
	p := Player{
		e,
		name,
		score}
	return p
}
