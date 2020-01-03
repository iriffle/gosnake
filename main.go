package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

const (
	WindowSizeX = 100
	WindowSizeY = 35
	MapWidth    = WindowSizeX
	MapHeight   = WindowSizeY
	Title       = "TestGame"
)

var (
	entities []*Entity
	gameMap  *GameMap
	dx, dy   int
	defStyle tcell.Style
)

func main() {
	fmt.Println("test")
	encoding.Register()
	s, err := tcell.NewScreen()

	// Error caching
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Set default colors and style
	bgColor := tcell.ColorDarkSlateBlue
	fgColor := tcell.ColorBurlyWood
	defStyle = tcell.StyleDefault.Background(bgColor).Foreground(fgColor)

	// Prepare screen
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.Clear()

	gameMap = &GameMap{
		Width:  MapWidth,
		Height: MapHeight,
	}

	gameMap.InitializeMap()

	player := &Entity{
		name:  "player",
		char:  "@",
		style: tcell.StyleDefault.Background(bgColor).Foreground(tcell.ColorDarkSlateGray),
		layer: 1,
		x:     1,
		y:     1,
	}

	npc := &Entity{
		name:  "npc",
		char:  "N",
		style: tcell.StyleDefault.Background(bgColor).Foreground(tcell.ColorDarkOrange),
		layer: 1,
		x:     10,
		y:     10,
	}

	entities = append(entities, player, npc)

	renderAll(s, defStyle, gameMap, entities)

	for {
		s.Show()

		handleInput(s, player)

		renderAll(s, defStyle, gameMap, entities)
	}
}
