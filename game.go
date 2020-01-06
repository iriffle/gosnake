package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/gdamore/tcell/views"
	"time"
)

const (
	MapWidth    = 100
	MapHeight   = 35
	MapStartX   = 0
	MapStartY   = 1
	SViewStartX = 0
	SViewStartY = 0
)

var (
	gameMap    *GameMap
	defBgColor tcell.Color = tcell.ColorDarkSlateBlue
	defFgColor tcell.Color = tcell.ColorWhite

	screenBgColor tcell.Color = tcell.ColorBlack
	screenFgColor tcell.Color = tcell.ColorWhite

	player1FgColor tcell.Color = tcell.ColorGreen

	defStyle tcell.Style = tcell.StyleDefault.
			Background(defBgColor).
			Foreground(defFgColor)
	screenStyle tcell.Style = tcell.StyleDefault.
			Background(screenBgColor).
			Foreground(screenFgColor)
	player1Style tcell.Style = tcell.StyleDefault.
			Background(defBgColor).
			Foreground(player1FgColor)

	dx, dy int

	bitRune rune = '*'
)

type Game struct {
	screen  tcell.Screen
	lview   *views.ViewPort
	sview   *views.ViewPort
	sbar    *views.TextBar
	players []*Player
	bits    []*Bit
}

func (g *Game) moveInterval() time.Duration {
	ms := 70
	return time.Duration(ms) * time.Millisecond
}

func (g *Game) Init() error {
	encoding.Register()

	if screen, err := tcell.NewScreen(); err != nil {
		return err
	} else if err = screen.Init(); err != nil {
		return err
	} else {
		screen.SetStyle(screenStyle)
		g.screen = screen
	}

	// Prepare screen
	g.screen.EnableMouse()
	g.screen.Clear()
	g.lview = views.NewViewPort(g.screen, MapStartX, MapStartY, MapWidth, MapHeight)
	g.sview = views.NewViewPort(g.screen, SViewStartX, SViewStartY, -1, 1)
	g.sbar = views.NewTextBar()
	g.sbar.SetView(g.sview)

	gameMap = &GameMap{
		Width:  MapWidth,
		Height: MapHeight,
	}

	gameMap.InitializeMap(defStyle)

	x := MapWidth / 2
	y := MapHeight / 2
	p1 := NewPlayer(x, y, 0, 3, '█', "Player1", player1Style)
	g.players = append(g.players, &p1)

	// b := NewBit(10, 10, 10, '*', pStyle)
	// g.bits = append(g.bits, &b)
	b := NewRandomBit(MapStartX, MapStartY, MapWidth, MapHeight, 10, bitRune, defStyle)
	g.bits = append(g.bits, &b)
	return nil
}

func (g *Game) Run() error {

	var b Bit
	renderAll(g, defStyle, gameMap, g.players, g.bits)

mainloop:
	for {
		g.screen.Show()

		for _, p := range g.players {
			dx, dy = 0, 0
			go handleInput(g.screen, p)

			switch p.direction {
			case 1:
				dy--
			case 2:
				dy++
			case 3:
				dx--
			case 4:
				dx++
			}
			if !gameMap.IsBlocked(p.pos[0].x+dx, p.pos[0].y+dy) {
				p.MoveEntity(dx, dy)
			} else {
				g.screen.Clear()
				renderStr(g.sview, SViewStartX, SViewStartY, defStyle, "Game Over")
				break mainloop
			}

			for i, bit := range g.bits {
				if p.pos[0].x == bit.x && p.pos[0].y == bit.y {
					p.score += bit.points
					p.AddSegment(p.pos[0].char, p.pos[0].style)
					g.bits = append(g.bits[:i], g.bits[i+1:]...)
					b = NewRandomBit(MapStartX, MapStartY, MapWidth, MapHeight, 10, bitRune, defStyle)
				}
			}
		}
		g.bits = append(g.bits, &b)
		renderAll(g, defStyle, gameMap, g.players, g.bits)
		time.Sleep(g.moveInterval())
	}
	return nil
}
