package main

import (
	"time"

	"github.com/gdamore/tcell"
)

// The player struct
type Player struct {
	*Entity
	name     string
	score    int
	count    int
	quitChan chan bool
	bitChan  chan int
	biteChan chan int
}

// Make a new player
func NewPlayer(x, y, score, direction int, char rune, name string, style tcell.Style) *Player {
	e := NewEntity(x, y, direction, 1, char, style)
	p := Player{
		Entity: e,
		name:   name,
		score:  score,
	}
	return &p
}

// Reset player's score and set back to middle of screen
func (p *Player) Reset(x, y, direction int, biteExplodedStyle tcell.Style) {
	style := p.pos[0].style
	p.Kill(biteExplodedStyle)
	p.score = 0
	p.Entity = NewEntity(x, y, direction, 1, p.pos[0].char, p.pos[0].style)
	p.pos[0].style = style
}

func (p *Player) Kill(biteExplodedStyle tcell.Style) {
	for i, _ := range p.pos {
		p.pos[i].style = biteExplodedStyle
		time.Sleep(20 * time.Millisecond)
	}
}

// Check if player is on top of a bit
func (p *Player) IsOnBit(g *Game) int {
	i := p.CheckBitPos(g.bits)
	if i != -1 {
		b := g.bits[i]
		p.score += b.points
		p.AddSegment(p.pos[0].char, p.pos[0].style)
	}
	return i
}

// Check the position of bits in relation to the player
// and see if there is a match
func (p *Player) CheckBitPos(bits []*Bit) int {
	for i, _ := range bits {
		if p.pos[0].x == bits[i].x && p.pos[0].y == bits[i].y {
			return i
		}
	}
	return -1
}

// Check the position of bites in relation to the player
// and see if there is a match
func (p *Player) CheckBitePos(bites []*Bite) int {
	for i, bite := range bites {
		if p.pos[0].x == bite.x && p.pos[0].y == bite.y {
			return i
		}
	}
	return -1
}

// Determine if player is on a bite and if so trigger explosion
func (p *Player) IsOnBite(g *Game, m *GameMap) int {
	i := p.CheckBitePos(g.bites)
	if i != -1 {
		b := g.bites[i]
		p.score += 50
		for i := 0; i < 4; i++ {
			p.AddSegment(p.pos[0].char, p.pos[0].style)
		}
		go b.ExplodeBite(m, g.biteMap, g.style.BiteExplodedStyle, g.style.DefStyle)
		return i
	}
	return -1
}

// Check if player is blocked
func (p *Player) IsBlocked(m *GameMap, biteMap *GameMap, entities []*Entity, players []*Player, dx, dy int) bool {
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
	if p.IsBlockedByMap(biteMap, dx, dy) {
		return true
	}

	return false
}

// Check if player is blocked by its own body
func (p *Player) IsBlockedBySelf(dx, dy int) bool {
	for i := range p.pos {
		if p.pos[0].x+dx == p.pos[i].x && p.pos[0].y+dy == p.pos[i].y && !(i == 0) && !(i == 1) {
			return true
		}
	}
	return false
}

// Check if player is blocked by another player
func (p *Player) IsBlockedByPlayer(players []*Player, dx, dy int) bool {
	for e := range players {
		for i := range players[e].pos {
			if p.pos[0].x+dx == players[e].pos[i].x && p.pos[0].y+dy == players[e].pos[i].y && !(p.name == players[e].name) {
				return true
			}
		}
	}
	return false
}

// Check if player is blocked by an entity
func (p *Player) IsBlockedByEntity(entities []*Entity, players []*Player, dx, dy int) bool {
	for e := range entities {
		for i := range entities[e].pos {
			if p.pos[0].x+dx == entities[e].pos[i].x && p.pos[0].y+dy == entities[e].pos[i].y && entities[e].pos[i].blocked {
				return true
			}
		}
	}
	return false
}
