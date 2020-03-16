package entity

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/stjiub/gosnake/gamemap"
)

// The player struct
type Player struct {
	name     string
	score    int
	count    int
	items    []*Item
	style    tcell.Style
	QuitChan chan bool
	BitChan  chan int
	BiteChan chan int
	*Entity
}

// Make a new player
func NewPlayer(x, y, score, direction int, char rune, name string, sty tcell.Style) *Player {
	e := NewEntity(x, y, direction, 1, char, sty)
	p := Player{
		Entity: e,
		name:   name,
		score:  score,
	}
	return &p
}

func (p *Player) InitChans() {
	p.QuitChan = make(chan bool)
	p.BitChan = make(chan int, 1)
	p.BiteChan = make(chan int, 1)
}

// Reset player's score and set back to middle of screen
func (p *Player) Reset(x, y, direction int, biteExplodedStyle tcell.Style) {
	style := p.pos[0].GetStyle()
	char := p.pos[0].GetChar()
	p.Kill(biteExplodedStyle)
	p.score = 0
	p.Entity = NewEntity(x, y, direction, 1, char, style)
	p.pos[0].SetStyle(style)
}

func (p *Player) Kill(biteExplodedStyle tcell.Style) {
	for i := range p.pos {
		p.pos[i].SetStyle(biteExplodedStyle)
		time.Sleep(20 * time.Millisecond)
	}
}

// Check the position of bits in relation to the player
// and see if there is a match
func (p *Player) CheckBitPos(bits []*Bit) int {
	px, py := p.pos[0].GetCurPos()
	for i := range bits {
		bx, by := bits[i].GetCurPos()
		if px == bx && py == by {
			return i
		}
	}
	return -1
}

// Check the position of bites in relation to the player
// and see if there is a match
func (p *Player) CheckBitePos(bites []*Bit) int {
	px, py := p.pos[0].GetCurPos()
	for i := range bites {
		bx, by := bites[i].GetCurPos()
		if px == bx && py == by {
			return i
		}
	}
	return -1
}

// Check if player is blocked
func (p *Player) IsBlocked(m *gamemap.GameMap, biteMap *gamemap.GameMap, entities []*Entity, players []*Player, dx, dy int) bool {
	if p.IsBlockedByMap(m, dx, dy) {
		return true
	}
	for i := range p.items {
		if p.items[i].effect == WallPass && p.items[i].activated {
			return false
		}
	}
	if p.IsBlockedByPlayer(players, dx, dy) {
		return true
	}
	if p.IsBlockedBySelf(dx, dy) {
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
	px, py := p.pos[0].GetCurPos()
	for i := range p.pos {
		ix, iy := p.pos[i].GetCurPos()
		if px+dx == ix && py+dy == iy && !(i == 0) && !(i == 1) {
			return true
		}
	}
	return false
}

// Check if player is blocked by another player
func (p *Player) IsBlockedByPlayer(players []*Player, dx, dy int) bool {
	px, py := p.pos[0].GetCurPos()
	for e := range players {
		for i := range players[e].pos {
			ix, iy := players[e].pos[i].GetCurPos()
			if px+dx == ix && py+dy == iy && !(p.name == players[e].name) {
				return true
			}
		}
	}
	return false
}

// Check if player is blocked by an entity
func (p *Player) IsBlockedByEntity(entities []*Entity, players []*Player, dx, dy int) bool {
	px, py := p.pos[0].GetCurPos()
	for e := range entities {
		for i := range entities[e].pos {
			ix, iy := entities[e].pos[i].GetCurPos()
			if px+dx == ix && py+dy == iy && entities[e].pos[i].IsBlocked() {
				return true
			}
		}
	}
	return false
}

// Generate bits where player's body was during collision
func (p *Player) DropBits(bits []*Bit, char rune, random int, sty tcell.Style) []*Bit {
	for i := range p.pos {
		ox, oy := p.pos[i].GetLastPos()
		b := NewBit(ox, oy, 10, char, random, DirNone, sty)
		bits = append(bits, b)
	}
	return bits
}

func (p *Player) AddItem(item *Item) {
	p.items = append(p.items, item)
	p.AdjustItemPos()
}

func (p *Player) RemoveItem(i int) {
	p.items[i] = p.items[len(p.items)-1]
	p.items[len(p.items)-1] = nil
	p.items = p.items[:len(p.items)-1]
	p.AdjustItemPos()
}

func (p *Player) AdjustItemPos() {
	for i := range p.items {
		p.items[i].pos = i
	}
}

func (p *Player) CheckItemPos(items []*Item) int {
	px, py := p.pos[0].GetCurPos()
	for i := range items {
		ix, iy := items[i].GetCurPos()
		if px == ix && py == iy {
			return i
		}
	}
	return -1
}

func (p *Player) AddScore(score int) {
	p.score += score
}

func (p *Player) SetScore(score int) {
	p.score = score
}

func (p *Player) GetScore() int {
	return p.score
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) SetName(name string) {
	p.name = name
}

func (p *Player) ActivateItem() {
	p.items[0].Activate(p)
}
