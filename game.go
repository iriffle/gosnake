package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/gdamore/tcell/views"
)

const (
	WindowSizeX = 100
	WindowSizeY = 35
	MapWidth    = WindowSizeX
	MapHeight   = WindowSizeY
	MapStartX   = 0
	MapStartY   = 1
	SViewStartX = 0
	SViewStartY = 0
)

var (
	gameMap  *GameMap
	defStyle tcell.Style
)

type Game struct {
	screen  tcell.Screen
	lview   *views.ViewPort
	sview   *views.ViewPort
	sbar    *views.TextBar
	players []*player
	bits    []*bit
}

func (g *Game) Init() error {
	encoding.Register()

	if screen, err := tcell.NewScreen(); err != nil {
		return err
	} else if err = screen.Init(); err != nil {
		return err
	} else {
		screen.SetStyle(tcell.StyleDefault.
			Background(tcell.ColorBlack).
			Foreground(tcell.ColorWhite))
		g.screen = screen
	}

	// Set default colors and style
	bgColor := tcell.ColorDarkSlateBlue

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

	gameMap.InitializeMap()

	x := MapWidth / 2
	y := MapHeight / 2
	pStyle := tcell.StyleDefault.
		Background(bgColor).
		Foreground(tcell.ColorWhite)
	p1 := NewPlayer(x, y, 0, 'O', "Player1", pStyle)
	g.players = append(g.players, &p1)

	// b := NewBit(10, 10, 10, '*', pStyle)
	// g.bits = append(g.bits, &b)
	SetBit(g, 10, '*', tcell.StyleDefault.
		Background(tcell.ColorDarkSlateBlue).
		Foreground(tcell.ColorWhite))

	return nil
}

func (g *Game) Run() error {

	renderAll(g, defStyle, gameMap, g.players, g.bits)

	for {
		g.screen.Show()

		handleInput(g.screen, g.players[0])

		for a, p := range g.players {
			for i, bit := range g.bits {
				if p.pos[0].x == bit.x && p.pos[0].y == bit.y {
					p.score += bit.points
					g.players[a].AddSegment('O', tcell.StyleDefault.
						Background(tcell.ColorDarkSlateBlue).
						Foreground(tcell.ColorWhite))
					g.bits = append(g.bits[:i], g.bits[i+1:]...)
					SetBit(g, 10, '*', tcell.StyleDefault.
						Background(tcell.ColorDarkSlateBlue).
						Foreground(tcell.ColorWhite))
				}
			}
		}

		renderAll(g, defStyle, gameMap, g.players, g.bits)
	}
	return nil
}
