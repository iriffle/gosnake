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
	GameWidth  = 100
	GameHeight = 35
	GameStartX = 0
	GameStartY = 0

	// Control bar values
	SViewWidth  = GameWidth
	SViewHeight = 1
	SViewStartX = 0
	SViewStartY = GameHeight + 1

	// Console values
	CViewWidth  = GameWidth
	CViewHeight = 1
	CViewStartX = 0
	CViewStartY = SViewStartY + 2

	// Preset colors
	DefBGColor         tcell.Color = tcell.ColorBlack
	DefFGColor         tcell.Color = tcell.ColorSteelBlue
	ScreenBGColor      tcell.Color = tcell.ColorBlack
	ScreenFGColor      tcell.Color = tcell.ColorWhite
	BitFGColor         tcell.Color = tcell.ColorWhite
	BiteFGColor        tcell.Color = tcell.ColorAqua
	BiteFGExplodeColor tcell.Color = tcell.ColorRed
)

var (
	// Current game map
	m *GameMap

	// Text to be displayed at bottom for controls
	controls     string = "w/s/a/d = up/down/left/right - esc = quit - f1 = restart - f12 = pause"
	menuOptions         = [3]string{"1 Player", "2 Player", "High Scores"}
	playerColors        = []tcell.Color{tcell.ColorGreen, tcell.ColorRed, tcell.ColorSilver, tcell.ColorAqua}

	// Runes to be used on map
	playerRune rune = '█'
	bitRune    rune = '■'
	wallRune   rune = '▒'
	floorRune  rune = ' '
	BiteRune   rune = '▲'

	// Number of bits that should be present on map
	numBits int = 5

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
	BiteStyle tcell.Style = tcell.StyleDefault.
			Background(DefBGColor).
			Foreground(BiteFGColor)
	BiteExplodedStyle tcell.Style = tcell.StyleDefault.
				Background(DefBGColor).
				Foreground(BiteFGExplodeColor)
	ControlStyle tcell.Style = tcell.StyleDefault.
			Background(ScreenBGColor).
			Foreground(ScreenFGColor)
	ConsoleStyle tcell.Style = tcell.StyleDefault.
			Background(ScreenBGColor).
			Foreground(ScreenFGColor)
)

// Main game struct
type Game struct {
	screen     tcell.Screen
	gview      *views.ViewPort
	sview      *views.ViewPort
	sbar       *views.TextBar
	cview      *views.ViewPort
	cbar       *views.TextBar
	players    []*Player
	entities   []*Entity
	bites      []*Bit
	bits       []*Bit
	maps       []*GameMap
	colors     []tcell.Color
	state      int
	mode       int
	level      int
	numPlayers int
	fps        int
	frames     int
	debug      bool
	bitQuit    chan bool
}

// Initialize the screen and set views/bars and styles
func (g *Game) InitScreen() error {
	// Prepare screen
	encoding.Register()
	if screen, err := tcell.NewScreen(); err != nil {
		return err
	} else if err = screen.Init(); err != nil {
		return err
	} else {
		screen.SetStyle(ScreenStyle)
		g.screen = screen
	}

	if g.screen.HasMouse() {
		g.screen.EnableMouse()
	}
	g.screen.ShowCursor(CViewStartX, CViewStartY)
	g.gview = views.NewViewPort(g.screen, GameStartX, GameStartY, GameWidth, GameHeight)
	g.sview = views.NewViewPort(g.screen, SViewStartX, SViewStartY, SViewWidth, SViewHeight)
	g.sbar = views.NewTextBar()
	g.sbar.SetView(g.sview)
	g.sbar.SetStyle(ControlStyle)

	if g.debug {
		g.cview = views.NewViewPort(g.screen, CViewStartX, CViewStartY, CViewWidth, CViewHeight)
		g.cbar = views.NewTextBar()
		g.cbar.SetView(g.cview)
		g.cbar.SetStyle(ConsoleStyle)
	}

	return nil
}

// Launch main menu screen
func (g *Game) MainMenu() {
	choice := false
	m := NewPlayerMenu(menuOptions, DefStyle, ScreenStyle)
	m.SetSelected(0)
	m.ChangeSelected()
	for !choice {
		renderMenu(g, &m, DefStyle)
		choice = handleMenu(g, &m)
	}
	i := m.GetSelected()
	switch i {
	case 0:
		g.numPlayers = 1
	case 1:
		g.numPlayers = 2
	case 2:
		g.screen.Fini()
		os.Exit(0)
	}

}

// Initialize game
func (g *Game) InitGame() {
	g.state = Play
	g.level = 1

	m = &GameMap{}
	g.maps = append(g.maps, m)
	m.InitLevel1(g, wallRune, floorRune, bitRune, DefStyle)
	g.colors = playerColors

	x := GameWidth / 2

	for i := 0; i < g.numPlayers; i++ {
		y := (GameHeight / 2) + (i * 2)

		pName := "player"
		pName = pName + strconv.Itoa(i+1)

		pStyle := tcell.StyleDefault.
			Background(DefBGColor).
			Foreground(g.colors[i])
		p := NewPlayer(x, y, 0, 3+i, playerRune, pName, pStyle)
		g.players = append(g.players, &p)
	}
	g.players[0].score = 0
	for i := 0; i < numBits; i++ {
		b := NewRandomBit(m, 10, bitRune, BitStyle)
		g.bits = append(g.bits, &b)
	}
}

// Run main game loop
func (g *Game) Run() {
	for _, p := range g.players {
		p.ch = make(chan bool)
		go g.handlePlayer(p)
	}

	for g.state == Play || g.state == Pause {
		g.handleLevel(m)
		g.getFPS()
		go handleInput(g)
		g.handlePause()
		renderAll(g, DefStyle, m)
		g.frames++
	}
	for _, p := range g.players {
		p.ch <- true
	}
	//g.bitQuit <- true
}

// Display high score screen
func (g *Game) HighScores() {
	// TODO
}

// Quit game
func (g *Game) Quit() {
	g.screen.Fini()
	os.Exit(0)
}

// Pause game until unpaused
func (g *Game) handlePause() {
	chQuit := false
	for g.state == Pause {
		if !chQuit {
			for _, p := range g.players {
				p.ch <- true
				chQuit = true
			}
			g.bitQuit <- true
		}
		renderCenterStr(g.gview, GameWidth, GameHeight-4, BitStyle, "PAUSED")
		g.screen.Show()

		if g.state == Play {
			for _, p := range g.players {
				//p.ch = make(chan bool)
				go g.handlePlayer(p)
			}
			//bitQuit := make(chan bool)
			go g.handleBits(m)
		}
	}
}

// Player movevement loop
func (g *Game) handlePlayer(p *Player) {
	for {
		select {
		default:
			dx, dy := p.CheckDirection(g)
			if p.IsBlocked(m, g.maps, g.players) {
				if g.numPlayers > 1 {
					if p.IsBlockedByPlayer(g.players) {
						for _, i := range p.pos {
							b := NewBit(i.ox, i.oy, 10, bitRune, true, BitStyle)
							g.bits = append(g.bits, &b)
						}
					}
					p.Reset(GameWidth/2, GameHeight/2, 3)
				} else {
					p.Kill()
					time.Sleep(100 * time.Millisecond)
					g.screen.Fini()
					g.state = Restart
				}
			} else {
				p.Move(dx, dy)
			}
			p.IsOnBit(g)
			p.IsOnBite(g, m)
			p.speed += p.score / 200
			time.Sleep(g.moveInterval(p.speed, p.direction))
		case <-p.ch:
			return
		}
	}
}

// Bit movement loop
func (g *Game) handleBits(m *GameMap) {
	for {
		select {
		default:
			for _, bit := range g.bits {
				if bit.moving {
					bit.Move(m)
				}
			}
			time.Sleep(500 * time.Millisecond)
		case <-g.bitQuit:
			return
		}
	}
}

// func (g *Game) handleEntities(m *GameMap) {
// 	for {
// 		select {
// 		default:
// 			for _, e := range g.entities {

// 			}
// 		}
// 	}
// }

func (g *Game) handleLevel(m *GameMap) {
	for _, p := range g.players {
		switch p.score {
		case 200:
			if g.level < 2 {
				m.InitLevel2(g, wallRune, floorRune, DefStyle)
				g.level = 2
			}
		case 400:
			if g.level < 3 {
				m.InitLevel3(g, wallRune, floorRune, BiteExplodedStyle)
				g.level = 3
			}
		}
	}
}

// Calculate FPS
func (g *Game) getFPS() {
	time.AfterFunc(1*time.Second, func() {
		g.fps = g.frames
		g.frames = 0
	})
}

// Calculate speed of player
func (g *Game) moveInterval(speed, direction int) time.Duration {
	ms := 150
	switch direction {
	case 1, 2:
		ms = 200
	}
	ms -= (speed / 100)
	return time.Duration(ms) * time.Millisecond
}

// Remove bit from map
func (g *Game) removeBit(i int) {
	b := &Bit{}
	g.bits[i] = g.bits[len(g.bits)-1]
	g.bits[len(g.bits)-1] = b
	g.bits = g.bits[:len(g.bits)-1]
}
