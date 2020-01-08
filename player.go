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

func (p *Player) IsPlayerBlocked(m *GameMap, players []*Player) bool {
	if p.IsPlayerBlockedByPlayer(players) {
		return true
	}
	if p.IsPlayerBlockedBySelf() {
		return true
	}
	if p.IsPlayerBlockedByMap(m) {
		return true
	}
	return false
}

func (p *Player) IsPlayerBlockedBySelf() bool {
	for a, i := range p.pos {
		if p.pos[0].x == i.x && p.pos[0].y == i.y && !(a == 0) {
			return true
		}
	}
	return false
}

func (p *Player) IsPlayerBlockedByPlayer(players []*Player) bool {
	for _, player := range players {
		for _, i := range player.pos {
			if p.pos[0].x == i.x && p.pos[0].y == i.y && !(p.name == player.name) {
				return true
			}
		}
	}
	return false
}

func (p *Player) IsPlayerBlockedByMap(m *GameMap) bool {
	if m.Objects[p.pos[0].x][p.pos[0].y].blocked {
		return true
	} else {
		return false
	}
}
