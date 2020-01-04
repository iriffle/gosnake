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
	entities []*Entity
	player   *Entity
	gameMap  *GameMap
	dx, dy   int
	defStyle tcell.Style
)

type Game struct {
	screen tcell.Screen
	lview  *views.ViewPort
	sview  *views.ViewPort
	sbar   *views.TextBar
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

	player = &Entity{
		name: "player",
		char: "ʘ",
		style: tcell.StyleDefault.
			Background(bgColor).
			Foreground(tcell.ColorWhite),
		layer: 1,
		x:     MapWidth / 2,
		y:     MapHeight / 2,
		pos:   []Coord{Coord{MapWidth / 2, MapHeight / 2}},
	}

	entities = append(entities, player)

	return nil
}

func (g *Game) Run() error {

	renderAll(g, defStyle, gameMap, entities)

	for {
		g.screen.Show()

		handleInput(g.screen, player)

		renderAll(g, defStyle, gameMap, entities)
	}
	return nil
}
