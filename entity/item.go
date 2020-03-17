package entity

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/stjiub/gosnake/gamemap"
)

// Player item states
const (
	WallPass = iota
)

type Item struct {
	effect    int
	activated bool
	pos       int
	duration  time.Duration
	ch        chan bool
	gamemap.Object
}

func NewItem(x, y, effect int, duration time.Duration, char rune, style tcell.Style) *Item {
	i := Item{
		effect:    effect,
		activated: false,
		duration:  duration,
	}
	i.SetPos(x, y)
	i.SetChar(char)
	i.SetStyle(style)
	i.Unblock()

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
