package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
)

const (
	WindowSizeX = 100
	WindowSizeY = 35
	MapWidth    = WindowSizeX
	MapHeight   = WindowSizeY
)

var (
	gameMap  *GameMap
	dx, dy   int
	defStyle tcell.Style
)

type Game struct {
	screen  tcell.Screen
	lview   *views.ViewPort
	sview   *views.ViewPort
	sbar    *views.TextBar
	players []*Player
}

func (g *Game) Init() error {

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
	g.lview = views.NewViewPort(g.screen, 0, 1, -1, -1)
	g.sview = views.NewViewPort(g.screen, 0, 0, -1, 1)
	g.sbar = views.NewTextBar()
	g.sbar.SetView(g.sview)

	gameMap = &GameMap{
		Width:  MapWidth,
		Height: MapHeight,
	}

	gameMap.InitializeMap()

	p1Object := &Object{
		char:  "O",
		layer: 1,
		style: tcell.StyleDefault.
			Background(bgColor).
			Foreground(tcell.ColorWhite),
		x: MapWidth / 2,
		y: MapHeight / 2,
	}

	p1Moving := &Moving{
		object:    p1Object,
		direction: "left",
		speed:     1,
	}

	p1Snake := &Snake{
		moving: p1Moving,
		pos:    []Coord{Coord{MapWidth / 2, MapHeight / 2}},
	}

	p1 := &Player{
		snake: p1Snake,
		name:  "Player1",
		score: 0,
	}

	g.players = append(g.players, p1)

	return nil
}

func (g *Game) Run() error {

	renderAll(g, defStyle, gameMap, g.players)

	for {
		g.screen.Show()

		handleInput(g.screen, g.players[0])

		renderAll(g, defStyle, gameMap, g.players)
	}
	return nil
}
