package main

import (
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/gdamore/tcell/views"
)

const (
	// Map values
	GameWidth  = 90
	GameHeight = 30
	MapStartX  = 0
	MapStartY  = 0

	// Control bar values
	CViewWidth  = GameWidth
	CViewHeight = 1
	CViewStartX = 0
	CViewStartY = GameHeight + 1

	// Preset colors
	DefBGColor     tcell.Color = tcell.ColorDarkSlateBlue
	DefFGColor     tcell.Color = tcell.ColorSteelBlue
	ScreenBGColor  tcell.Color = tcell.ColorBlack
	ScreenFGColor  tcell.Color = tcell.ColorWhite
	Player1FGColor tcell.Color = tcell.ColorGreen
	BitFGColor     tcell.Color = tcell.ColorWhite
)

var (
	// Current game map
	gameMap *GameMap

	// Preset styles
	DefStyle tcell.Style = tcell.StyleDefault.
			Background(DefBGColor).
			Foreground(DefFGColor)
	ScreenStyle tcell.Style = tcell.StyleDefault.
			Background(ScreenBGColor).
			Foreground(ScreenFGColor)
	BitStyle tcell.Style = tcell.StyleDefault.
			Background(DefBGColor).
			Foreground(BitFGColor)
	ControlStyle tcell.Style = tcell.StyleDefault.
			Background(ScreenBGColor).
			Foreground(ScreenFGColor)

	// Text to be displayed at bottom for controls
	Controls string = "w/s/a/d = up/down/left/right - esc = quit - f1 = restart - f12 = pause"

	// Runes to be used on map
	playerRune rune = '█'
	bitRune    rune = '*'
	wallRune   rune = '▒'
	floorRune  rune = ' '

	// Number of bits that should be present on map
	numBits int = 5

	// Used for player movement
	dx, dy int
)

type Game struct {
	screen     tcell.Screen
	gview      *views.ViewPort
	cview      *views.ViewPort
	cbar       *views.TextBar
	players    []*Player
	bits       []*Bit
	state      int
	level      int
	numPlayers int
	colors     []tcell.Color
}

func (g *Game) InitScreen() error {
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
	g.gview = views.NewViewPort(g.screen, MapStartX, MapStartY, GameWidth, GameHeight)
	g.cview = views.NewViewPort(g.screen, CViewStartX, CViewStartY, CViewWidth, CViewHeight)
	g.cbar = views.NewTextBar()
	g.cbar.SetView(g.cview)
	g.cbar.SetStyle(ControlStyle)

	return nil
}

func (g *Game) MainMenu() {
	g.gview.Fill(' ', DefStyle)
	lastChoice, choice := 1, 1
	for !(choice == 3) {
		switch choice {
		case 1:
			renderCenterStr(g.gview, GameWidth, GameHeight-4, BitStyle, "1 Player")
			renderCenterStr(g.gview, GameWidth, GameHeight, DefStyle, "2 Player")
		case 2:
			renderCenterStr(g.gview, GameWidth, GameHeight-4, DefStyle, "1 Player")
			renderCenterStr(g.gview, GameWidth, GameHeight, BitStyle, "2 Player")
		}
		g.screen.Show()
		lastChoice = choice
		choice = handleMenu(g, choice)
	}
	switch lastChoice {
	case 1:
		g.numPlayers = 1
	case 2:
		g.numPlayers = 2
	}
}

func (g *Game) InitGame() {
	g.state = 0
	g.level = 1

	gameMap = &GameMap{
		Width:  GameWidth,
		Height: GameHeight,
	}
	gameMap.InitLevel1(wallRune, floorRune, DefStyle)

	g.colors = append(g.colors, tcell.ColorGreen, tcell.ColorRed)

	x := GameWidth / 2

	for i := 0; i < g.numPlayers; i++ {
		y := (GameHeight / 2) + (i * 2)

		pName := "player"
		pName = pName + strconv.Itoa(i)

		pStyle := tcell.StyleDefault.
			Background(DefBGColor).
			Foreground(g.colors[i])
		p := NewPlayer(x, y, 0, 3, playerRune, pName, pStyle)
		g.players = append(g.players, &p)
	}

	for i := 0; i < numBits; i++ {
		b := NewRandomBit(gameMap, 10, bitRune, BitStyle)
		g.bits = append(g.bits, &b)
	}
}

func (g *Game) Run() {
	renderAll(g, DefStyle, gameMap, g.players, g.bits)

	for !(g.state == 1) {

		for _, p := range g.players {
			dx, dy = 0, 0
			go handleInput(g)

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

			if p.IsPlayerBlocked(gameMap, g.players) {
				g.screen.Fini()
				g.state = 1
			} else {
				p.MoveEntity(dx, dy)
			}

			g.isOnBit(p)
		}

		renderAll(g, DefStyle, gameMap, g.players, g.bits)
		time.Sleep(g.moveInterval(g.players[0].score))
	}
}

func (g *Game) Pause(p *Player) {
	for g.state == 2 {
		go handleInput(g)
	}
}

func (g *Game) moveInterval(score int) time.Duration {
	ms := 80 - (score / 10)
	return time.Duration(ms) * time.Millisecond
}

func (g *Game) removeBit(i int) {
	b := &Bit{}
	g.bits[i] = g.bits[len(g.bits)-1]
	g.bits[len(g.bits)-1] = b
	g.bits = g.bits[:len(g.bits)-1]
}

func (g *Game) isOnBit(p *Player) {
	onBit, i := p.CheckPlayerOnBit(g.bits)
	if onBit {
		b := g.bits[i]
		p.score += b.points
		p.AddSegment(p.pos[0].char, p.pos[0].style)
		g.removeBit(i)
		newB := NewRandomBit(gameMap, 10, bitRune, BitStyle)
		g.bits = append(g.bits, &newB)

	}
}
