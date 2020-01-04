package main

import (
	"github.com/gdamore/tcell"
)

type Object struct {
	x, y, speed int
	char        string
	style       tcell.Style
	newDir      string
	oldDir      string
}

type Player struct {
	*Object
	pos   []Object
	name  string
	score int
}

func NewObject(char, newDir, oldDir string, style tcell.Style, x, y, speed int) *Object {
	o := &Object{
		char:   char,
		style:  style,
		x:      x,
		y:      y,
		newDir: newDir,
		oldDir: oldDir,
		speed:  speed,
	}
	return o
}

func NewPlayer(style tcell.Style, x, y, speed int, char, newDir, oldDir, name string) *Player {
	pObject := NewObject(char, newDir, oldDir, style, x, y, speed)
	p := &Player{
		Object: pObject,
		name:   name,
		score:  0,
	}
	p.pos = append(p.pos, *pObject)
	return p
}

func (o *Object) MoveObject(dx, dy int) {
	// Move the Entity by the amount (dx, dy)
	o.x += dx
	o.y += dy
}

func (p *Player) MovePlayer(dx, dy int) {
	var n int
	first := true
	for i, _ := range p.pos {
		if first {
			p.pos[0].x += dx
			p.pos[0].y += dy
			first = false
		} else {
			n = i - 1
			p.pos[i].x = p.pos[n].x
			p.pos[i].y = p.pos[n].y
		}
		// switch p.pos[i].newDir {
		// case "right":
		// 	p.pos[i].x++
		// case "left":
		// 	p.pos[i].x--
		// case "down":
		// 	p.pos[i].y++
		// case "up":
		// 	p.pos[i].y--
		// }
		// p.pos[i].oldDir = p.pos[i].newDir
	}
}
