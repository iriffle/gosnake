package main

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
)

type object struct {
	x, y, ox, oy int
	char         rune
	style        tcell.Style
}

func NewObject(x, y int, char rune, style tcell.Style) object {
	o := object{
		x,
		y,
		x,
		y,
		char,
		style}
	return o
}

func (o *object) MoveObject(dx, dy int) {
	o.x += dx
	o.y += dy
}

type bit struct {
	object
	points int
}

func NewBit(x, y, points int, char rune, style tcell.Style) bit {
	o := NewObject(x, y, char, style)
	b := bit{
		o,
		points,
	}
	return b
}

func SetBit(g *Game, points int, char rune, style tcell.Style) {
	rand.Seed(time.Now().UnixNano())
	minX := MapStartX
	maxX := minX + MapWidth
	minY := MapStartY
	maxY := minY + MapHeight
	randX := rand.Intn(maxX-minX+1) + minX
	randY := rand.Intn(maxY-minY+1) + minY
	b := NewBit(randX, randY, points, char, style)
	g.bits = append(g.bits, &b)
}

type entity struct {
	pos []object
}

func NewEntity(x, y int, char rune, style tcell.Style) entity {
	o := NewObject(x, y, char, style)
	e := entity{}
	e.pos = append(e.pos, o)
	return e
}

func (e *entity) MoveEntity(dx, dy int) {
	first := true
	e.pos[0].ox = e.pos[0].x
	e.pos[0].oy = e.pos[0].y
	e.pos[0].x += dx
	e.pos[0].y += dy

	for i, _ := range e.pos {
		if !first {
			e.pos[i].ox = e.pos[i].x
			e.pos[i].oy = e.pos[i].y
			e.pos[i].x = e.pos[i-1].ox
			e.pos[i].y = e.pos[i-1].oy
		} else {
			first = false
		}
	}
}

func (e *entity) AddSegment(char rune, style tcell.Style) {
	x := e.pos[len(e.pos)-1].ox
	y := e.pos[len(e.pos)-1].oy
	o := NewObject(x, y, char, style)
	e.pos = append(e.pos, o)
}

type player struct {
	entity
	name  string
	score int
}

func NewPlayer(x, y, score int, char rune, name string, style tcell.Style) player {
	e := NewEntity(x, y, char, style)
	p := player{
		e,
		name,
		score}
	return p
}
