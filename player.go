package main

import (
	"time"

	"github.com/gdamore/tcell"
)

// The player struct
type Player struct {
	Entity
	name  string
	score int
	count int
	ch    chan bool
}

// Make a new player
func NewPlayer(x, y, score, direction int, char rune, name string, style tcell.Style) Player {
	e := NewEntity(x, y, direction, 1, char, style)
	p := Player{
		Entity: e,
		name:   name,
		score:  score,
	}
	return p
}

// Reset player's score and set back to middle of screen
func (p *Player) Reset(x, y, direction int) {
	style := p.pos[0].style
	p.Kill()
	p.score = 0
	p.Entity = NewEntity(x, y, direction, 1, p.pos[0].char, p.pos[0].style)
	p.pos[0].style = style
}

func (p *Player) Kill() {
	for i, _ := range p.pos {
		p.pos[i].style = BiteExplodedStyle
		time.Sleep(20 * time.Millisecond)
	}
}

// Check if player is on top of a bit
func (p *Player) IsOnBit(g *Game) {
	onBit, i := p.CheckBitPos(g.bits)
	if onBit {
		b := g.bits[i]
		p.score += b.points
		p.AddSegment(p.pos[0].char, p.pos[0].style)
		g.removeBit(i)
	}
}

// Check the position of bits in relation to the player
// and see if there is a match
func (p *Player) CheckBitPos(bits []*Bit) (bool, int) {
	i := 0
	for i, bit := range bits {
		if p.pos[0].x == bit.x && p.pos[0].y == bit.y {
			return true, i
		}
	}
	return false, i
}

func (p *Player) CheckBitePos(bites []*Bite) (bool, int) {
	i := 0
	for i, bite := range bites {
		if p.pos[0].x == bite.x && p.pos[0].y == bite.y {
			return true, i
		}
	}
	return false, i
}

func (p *Player) IsOnBite(g *Game, m *GameMap) {
	onBite, i := p.CheckBitePos(g.bites)
	if onBite {
		b := g.bites[i]
		p.score += 50
		for i := 0; i < 4; i++ {
			p.AddSegment(p.pos[0].char, p.pos[0].style)
		}
		go b.ExplodeBite(g, m)
	}
}

// Check if player is blocked
func (p *Player) IsBlocked(m *GameMap, bites []*GameMap, entities []*Entity, players []*Player, dx, dy int) bool {
	if p.IsBlockedByPlayer(players, dx, dy) {
		return true
	}
	if p.IsBlockedBySelf(dx, dy) {
		return true
	}
	if p.IsBlockedByMap(m, dx, dy) {
		return true
	}
	if p.IsBlockedByEntity(entities, players, dx, dy) {
		return true
	}
	for _, bite := range bites {
		if p.IsBlockedByMap(bite, dx, dy) {
			return true
		}
	}
	return false
}

// Check if player is blocked by its own body
func (p *Player) IsBlockedBySelf(dx, dy int) bool {
	for a, i := range p.pos {
		if p.pos[0].x+dx == i.x && p.pos[0].y+dy == i.y && !(a == 0) {
			return true
		}
	}
	return false
}

// Check if player is blocked by another player
func (p *Player) IsBlockedByPlayer(players []*Player, dx, dy int) bool {
	for _, player := range players {
		for _, i := range player.pos {
			if p.pos[0].x+dx == i.x && p.pos[0].y+dy == i.y && !(p.name == player.name) {
				return true
			}
		}
	}
	return false
}

func (p *Player) IsBlockedByEntity(entities []*Entity, players []*Player, dx, dy int) bool {
	for _, p := range players {
		for _, entity := range entities {
			for _, i := range entity.pos {
				if p.pos[0].x+dx == i.x && p.pos[0].y+dy == i.y && i.blocked {
					return true
				}
			}
		}
	}
	return false
}
