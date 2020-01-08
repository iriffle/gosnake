package main

import (
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/gdamore/tcell/views"
)

const (
	MapWidth    = 100
	MapHeight   = 35
	MapStartX   = 0
	MapStartY   = 0
	SViewWidth  = 40
	SViewHeight = 20
	SViewStartX = (MapWidth / 2) - (SViewWidth / 2)
	SViewStartY = (MapHeight / 2) - (SViewHeight / 2)

	DefBGColor     tcell.Color = tcell.ColorDarkSlateBlue
	DefFGColor     tcell.Color = tcell.ColorSteelBlue
	ScreenBGColor  tcell.Color = tcell.ColorBlack
	ScreenFGColor  tcell.Color = tcell.ColorWhite
	Player1FGColor tcell.Color = tcell.ColorGreen
	BitFGColor     tcell.Color = tcell.ColorWhite
)

var (
	gameMap *GameMap

	DefStyle tcell.Style = tcell.StyleDefault.
			Background(DefBGColor).
			Foreground(DefFGColor)
	ScreenStyle tcell.Style = tcell.StyleDefault.
			Background(ScreenBGColor).
			Foreground(ScreenFGColor)
	Player1Style tcell.Style = tcell.StyleDefault.
			Background(DefBGColor).
			Foreground(Player1FGColor)
	BitStyle tcell.Style = tcell.StyleDefault.
			Background(DefBGColor).
			Foreground(BitFGColor)

	dx, dy int

	bitRune   rune = '*'
	wallRune  rune = '▒'
	floorRune rune = ' '
	numBits   int  = 5
)

type Game struct {
	screen  tcell.Screen
	lview   *views.ViewPort
	sview   *views.ViewPort
	sbar    *views.TextBar
	players []*Player
	bits    []*Bit
	state   int
}

func (g *Game) moveInterval(score int) time.Duration {
	ms := 80 - (score / 10)
	return time.Duration(ms) * time.Millisecond
}

func (g *Game) Init() error {
	encoding.Register()

	if screen, err := tcell.NewScreen(); err != nil {
		return err
	} else if err = screen.Init(); err != nil {
		return err
	} else {
		screen.SetStyle(ScreenStyle)
		g.screen = screen
	}

	// Prepare screen
	g.screen.EnableMouse()
	g.screen.Clear()
	g.lview = views.NewViewPort(g.screen, MapStartX, MapStartY, MapWidth, MapHeight)
	//g.sview = views.NewViewPort(g.screen, SViewStartX, SViewStartY, SViewWidth, SViewHeight)
	//g.sbar = views.NewTextBar()
	//g.sbar.SetView(g.sview)

	g.state = 0

	gameMap = &GameMap{
		Width:  MapWidth,
		Height: MapHeight,
	}

	gameMap.InitializeMap(wallRune, floorRune, DefStyle)

	x := MapWidth / 2
	y := MapHeight / 2
	p1 := NewPlayer(x, y, 0, 3, '█', "Player1", Player1Style)
	g.players = append(g.players, &p1)

	for i := 0; i < numBits; i++ {
		b := NewRandomBit(MapStartX, MapStartY, MapWidth, MapHeight, 10, bitRune, BitStyle)
		g.bits = append(g.bits, &b)
	}
	return nil
}

func (g *Game) Run() error {

	var b Bit
	renderAll(g, DefStyle, gameMap, g.players, g.bits)

	for !(g.state == 1) {
		g.screen.Show()

		go func() {
			for _, p := range g.players {
				dx, dy = 0, 0
				go handleInput(g, p)

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
				go func() {
					if p.IsPlayerBlocked(gameMap, g.players) {
						g.screen.Fini()
						g.state = 1
					} else {
						p.MoveEntity(dx, dy)
					}
				}()

				for i, bit := range g.bits {
					if p.pos[0].x == bit.x && p.pos[0].y == bit.y {
						p.score += bit.points
						p.AddSegment(p.pos[0].char, p.pos[0].style)
						g.bits = append(g.bits[:i], g.bits[i+1:]...)
						b = NewRandomBit(MapStartX, MapStartY, MapWidth, MapHeight, 10, bitRune, BitStyle)
					}
				}
			}
		}()
		g.bits = append(g.bits, &b)
		renderAll(g, DefStyle, gameMap, g.players, g.bits)
		time.Sleep(g.moveInterval(g.players[0].score))
	}
	return nil
}

func (g *Game) Pause() {
	for g.state == 2 {
		go handleInput(g, g.players[0])
	}
}
