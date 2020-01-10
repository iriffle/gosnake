package main

import (
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/gdamore/tcell/views"
)

const (
	// Game states
	Play     = 0
	Quit     = 1
	Pause    = 2
	Restart  = 3
	MainMenu = 4

	// Game modes
	Basic    = 0
	Advanced = 1
	Battle   = 2

	// Map values
	MapWidth  = 100
	MapHeight = 35
	MapStartX = 0
	MapStartY = 0

	// Control bar values
	sviewWidth  = MapWidth
	sviewHeight = 1
	sviewStartX = 0
	sviewStartY = MapHeight + 1

	cviewWidth  = MapWidth
	cviewHeight = 1
	cviewStartX = 0
	cviewStartY = sviewStartY + 2

	// Preset colors
	DefBGColor     tcell.Color = tcell.ColorBlack
	DefFGColor     tcell.Color = tcell.ColorSteelBlue
	ScreenBGColor  tcell.Color = tcell.ColorBlack
	ScreenFGColor  tcell.Color = tcell.ColorWhite
	Player1FGColor tcell.Color = tcell.ColorGreen
	BitFGColor     tcell.Color = tcell.ColorWhite
)

var (
	// Current game map
	m *GameMap

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
	DebugStyle tcell.Style = tcell.StyleDefault.
			Background(ScreenBGColor).
			Foreground(ScreenFGColor)

	// Text to be displayed at bottom for controls
	Controls    string = "w/s/a/d = up/down/left/right - esc = quit - f1 = restart - f12 = pause"
	menuOptions        = [3]string{"1 Player", "2 Player", "High Scores"}

	// Runes to be used on map
	playerRune rune = '█'
	bitRune    rune = '*'
	wallRune   rune = '▒'
	floorRune  rune = ' '

	// Number of bits that should be present on map
	numBits int = 5

	// Used for player movement
	dx, dy int

	// Menu variables
	p1Str    string = "1 Player"
	p2Str    string = "2 Player"
	p1Width  int    = MapWidth
	p1Height int    = MapHeight - 4
	p2Width  int    = MapWidth
	p2Height int    = MapHeight
)

type Game struct {
	screen     tcell.Screen
	gview      *views.ViewPort
	sview      *views.ViewPort
	sbar       *views.TextBar
	cview      *views.ViewPort
	cbar       *views.TextBar
	players    []*Player
	bits       []*Bit
	colors     []tcell.Color
	state      int
	mode       int
	level      int
	numPlayers int
	debug      bool
}

func (g *Game) InitScreen() error {
	encoding.Register()

	if screen, err := tcell.NewConsoleScreen(); err != nil {
		return err
	} else if err = screen.Init(); err != nil {
		return err
	} else {
		screen.SetStyle(ScreenStyle)
		g.screen = screen
	}

	// Prepare screen
	if g.screen.HasMouse() {
		g.screen.EnableMouse()
	}
	//g.screen.ShowCursor(cviewStartX, cviewStartY)
	g.screen.Clear()
	g.gview = views.NewViewPort(g.screen, MapStartX, MapStartY, MapWidth, MapHeight)
	g.sview = views.NewViewPort(g.screen, sviewStartX, sviewStartY, sviewWidth, sviewHeight)
	g.sbar = views.NewTextBar()
	g.sbar.SetView(g.sview)
	g.sbar.SetStyle(ControlStyle)

	if g.debug {
		g.cview = views.NewViewPort(g.screen, cviewStartX, cviewStartY, cviewWidth, cviewHeight)
		g.cbar = views.NewTextBar()
		g.cbar.SetView(g.cview)
		g.cbar.SetStyle(DebugStyle)
	}

	return nil
}

func (g *Game) MainMenu() {
	choice := false
	m := NewPlayerMenu(menuOptions, DefStyle, ScreenStyle)
	m.SetSelected(0)
	m.ChangeSelected()
	for !choice {
		renderMenu(g, &m, DefStyle)
		choice = handleMenu(g, &m)
	}

}

func (g *Game) InitGame() {
	g.state = Play
	g.level = 1

	m = &GameMap{
		Width:  MapWidth,
		Height: MapHeight,
	}
	m.InitLevel1(wallRune, floorRune, DefStyle)

	g.colors = append(g.colors, tcell.ColorGreen, tcell.ColorRed)

	x := MapWidth / 2

	for i := 0; i < g.numPlayers; i++ {
		y := (MapHeight / 2) + (i * 2)

		pName := "player"
		pName = pName + strconv.Itoa(i+1)

		pStyle := tcell.StyleDefault.
			Background(DefBGColor).
			Foreground(g.colors[i])
		p := NewPlayer(x, y, 0, 3+i, playerRune, pName, pStyle)
		g.players = append(g.players, &p)
	}

	for i := 0; i < numBits; i++ {
		b := NewRandomBit(m, 10, bitRune, BitStyle)
		g.bits = append(g.bits, &b)
	}
}

func (g *Game) Run() {
	renderAll(g, DefStyle, m)

	for g.state == Play || g.state == Pause {
		go handleInput(g)

		for _, p := range g.players {

			for g.state == Pause {
				renderCenterStr(g.gview, MapWidth, MapHeight-4, BitStyle, "PAUSED")
				g.screen.Show()
				continue
			}

			dx, dy = 0, 0
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

			if p.IsPlayerBlocked(m, g.players) {
				if g.numPlayers == 1 {
					g.screen.Fini()
					g.state = Restart
				} else {
					if p.IsPlayerBlockedByPlayer(g.players) {
						for _, i := range p.pos {
							b := NewBit(i.ox, i.oy, 10, bitRune, BitStyle)
							g.bits = append(g.bits, &b)
						}
					}
					p.ResetPlayer(MapWidth/2, MapHeight/2, 3)
				}
			} else {
				p.MoveEntity(dx, dy)
			}

			g.isOnBit(p)
		}

		renderAll(g, DefStyle, m)
		if g.state == Play {
			time.Sleep(g.moveInterval(g.players[0].score))
		}
	}
}

func (g *Game) Quit() {
	g.screen.Fini()
	os.Exit(0)
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
		newB := NewRandomBit(m, 10, bitRune, BitStyle)
		g.bits = append(g.bits, &newB)

	}
}
