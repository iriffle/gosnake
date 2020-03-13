package main

import (
	"github.com/gdamore/tcell"
	"time"
)

type Item struct {
	effect    int
	activated bool
	pos       int
	duration  time.Duration
	ch        chan bool
	Object
}

func NewItem(x, y, effect int, duration time.Duration, char rune, style tcell.Style) *Item {
	i := Item{
		effect:    effect,
		activated: false,
		duration:  duration,
	}
	i.x = x
	i.y = y
	i.ox = x
	i.oy = y
	i.char = char
	i.style = style
	i.blocked = false

	return &i
}

func (i *Item) Activate(p *Player) {
	switch i.effect {
	case WallPass:
		go i.WallPass(p)
	}
}

func (i *Item) WallPass(p *Player) {
	i.activated = true
	for end := time.Now().Add(i.duration); ; {
		if time.Now().After(end) {
			i.activated = false
			p.SetChar('*')
			p.RemoveItem(i.pos)
			break
		}
	}
}

// func (i *Item) SlowSpeed(p *Player) {
// 	for end := time.Now().Add()
// }
