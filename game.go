package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/gdamore/tcell/views"
)

const (
	// Game states
	Play     = iota
	Quit     = iota
	Pause    = iota
	Restart  = iota
	MainMenu = iota
)

const (
	// Menu pages
	MenuMain     = iota
	MenuPlayer   = iota
	MenuMode     = iota
	MenuScore    = iota
	MenuSettings = iota
)

const (
	// Direction
	DirUp    = iota
	DirDown  = iota
	DirLeft  = iota
	DirRight = iota
	DirAll   = iota
)

const (
	// Game modes
	Basic    = iota
	Advanced = iota
	Battle   = iota

	// Map values
	MapWidth  = 100
	MapHeight = 35
	MapStartX = 0
	MapStartY = 0

	// Control bar values
	SViewWidth  = MapWidth
	SViewHeight = 1
	SViewStartX = 0
	SViewStartY = MapHeight + 1

	// Console values
	CViewWidth  = MapWidth
	CViewHeight = 1
	CViewStartX = 0
	CViewStartY = SViewStartY + 2

	// High Score count
	MaxHighScores = 5

	// Preset colors
	Black   = tcell.ColorBlack
	Maroon  = tcell.ColorMaroon
	Green   = tcell.ColorGreen
	Navy    = tcell.ColorNavy
	Olive   = tcell.ColorOlive
	Purple  = tcell.ColorPurple
	Teal    = tcell.ColorTeal
	Silver  = tcell.ColorSilver
	Gray    = tcell.ColorGray
	Red     = tcell.ColorRed
	Blue    = tcell.ColorBlue
	Lime    = tcell.ColorLime
	Yellow  = tcell.ColorYellow
	Fuchsia = tcell.ColorFuchsia
	Aqua    = tcell.ColorAqua
	White   = tcell.ColorWhite

	DefBGStyle = Black
	DefFGStyle = Silver
	SelFGStyle = Aqua
)

var (
	// Current game map
	m *GameMap

	// Text to be displayed at bottom for controls
	controls        string = "w/s/a/d = up/down/left/right - esc = quit - f1 = restart - f12 = pause"
	mainOptions            = []string{"Play", "High Scores", "Settings"}
	playerOptions          = []string{"1 Player", "2 Player"}
	gameModeOptions        = []string{"Basic", "Advanced", "Battle"}
	playerColors           = []tcell.Color{tcell.ColorGreen, tcell.ColorRed, tcell.ColorSilver, tcell.ColorAqua}

	// Runes to be used on map
	PlayerRune      rune = '█'
	BitRune         rune = '■'
	WallRune        rune = '▒'
	FloorRune       rune = ' '
	BiteUpRune      rune = '▲'
	BiteDownRune    rune = '▼'
	BiteLeftRune    rune = '◄'
	BiteRightRune   rune = '►'
	BiteAllRune     rune = '◆'
	BiteExplodeRune rune = '░'

	// Number of bits that should be present on map
	numBits int = 5

	// Preset styles
	DefStyle          = GetStyle(DefBGStyle, DefFGStyle)
	SelStyle          = GetStyle(DefBGStyle, SelFGStyle)
	BitStyle          = GetStyle(Black, White)
	BiteStyle         = GetStyle(Black, Fuchsia)
	BiteExplodedStyle = GetStyle(Black, Red)
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
	bites      []*Bite
	bits       []*Bit
	maps       []*GameMap
	colors     []tcell.Color
	state      int
	mode       int
	level      int
	numPlayers int
	fps        int
	frames     int
	bitQuit    chan bool
	scores     [][]string
	scoreFile  string
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
		screen.SetStyle(DefStyle)
		g.screen = screen
	}

	if g.screen.HasMouse() {
		g.screen.EnableMouse()
	}
	g.screen.ShowCursor(CViewStartX, CViewStartY)
	g.gview = views.NewViewPort(g.screen, MapStartX, MapStartY, MapWidth, MapHeight)
	g.sview = views.NewViewPort(g.screen, SViewStartX, SViewStartY, SViewWidth, SViewHeight)
	g.sbar = views.NewTextBar()
	g.sbar.SetView(g.sview)
	g.sbar.SetStyle(DefStyle)

	return nil
}

// Launch main menu screen
func (g *Game) MainMenu() {
	g.state = MainMenu
	cMenu := MenuMain
	g.readScores()
	for g.state != Play {
		// Main menu
		if cMenu == MenuMain {
			i := g.handleMenu(mainOptions)
			switch i {
			case -1:
				g.screen.Fini()
				os.Exit(0)
			case 0:
				cMenu = MenuPlayer
			case 1:
				cMenu = MenuScore
			}
		}
		// Player number menu
		if cMenu == MenuPlayer {
			i := g.handleMenu(playerOptions)
			switch i {
			case -1:
				cMenu = 0
			case 0:
				g.numPlayers = 1
				g.state = Play
				break
			case 1:
				g.numPlayers = 2
				g.state = Play
				break
			}
		}
		// High score screen
		for cMenu == MenuScore {
			err := renderHighScores(g, DefStyle)
			if err != nil {
				log.Println(err)
			}
			ev := g.screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					cMenu = MenuMain
					break
				}
			}
		}
	}
}

// Initialize game
func (g *Game) InitGame() {
	g.state = Play
	g.level = 1

	m = &GameMap{
		Width:  MapWidth,
		Height: MapHeight,
		X:      MapStartX,
		Y:      MapStartY,
	}
	g.maps = append(g.maps, m)
	m.InitMap()
	m.InitMapBoundary(WallRune, FloorRune, DefStyle)
	m.InitLevel1(g)
	g.colors = playerColors

	x := MapWidth / 2

	// Create a player for selected number of players
	for i := 0; i < g.numPlayers; i++ {
		y := (MapHeight / 2) + (i * 2)

		pName := "player"
		pName = pName + strconv.Itoa(i+1)

		pStyle := tcell.StyleDefault.
			Background(DefBGStyle).
			Foreground(g.colors[i])
		p := NewPlayer(x, y, 0, (DirLeft - i), PlayerRune, pName, pStyle)
		g.players = append(g.players, &p)
	}
	g.players[0].score = 0
	for i := 0; i < numBits; i++ {
		b := NewRandomBit(m, 10, BitRune, BitStyle)
		g.bits = append(g.bits, &b)
	}
	log.Println("Initialized game with ", strconv.Itoa(g.numPlayers), " players.")
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

func (g *Game) handleMenu(options []string) int {
	choice := 0
	m := NewMainMenu(options, DefStyle, SelStyle, 0)
	m.SetSelected(0)
	m.ChangeSelected()
	for choice == 0 {
		renderMenu(g, &m, DefStyle)
		choice = handleMenuInput(g, &m)
	}
	if choice == 1 {
		choice = m.GetSelected()
	}
	return choice
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
		}
		renderCenterStr(g.gview, MapWidth, MapHeight-4, BitStyle, "PAUSED")
		g.screen.Show()

		if g.state == Play {
			for _, p := range g.players {
				go g.handlePlayer(p)
			}
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
			if p.IsBlocked(m, g.maps, g.entities, g.players, dx, dy) {
				if g.numPlayers > 1 {
					//if p.IsBlockedByPlayer(g.players) {
					for _, i := range p.pos {
						b := NewBit(i.ox, i.oy, 10, BitRune, BitRandom, BitStyle)
						g.bits = append(g.bits, &b)
					}
					//}
					g.readScores()
					scoreChange := g.checkScores()
					if scoreChange {
						g.writeScores()
					}
					p.Reset(MapWidth/2, MapHeight/2, 3)
				} else {
					p.Kill()
					g.readScores()
					scoreChange := g.checkScores()
					if scoreChange {
						g.writeScores()
					}
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
				switch bit.state {
				case BitRandom:
					bit.Move(m)
				}
			}
			time.Sleep(500 * time.Millisecond)
		case <-g.bitQuit:
			return
		}
	}
}

// Change level based on player score
func (g *Game) handleLevel(m *GameMap) {
	for _, p := range g.players {
		switch p.score {
		case 200:
			if g.level < 2 {
				m.InitLevel2(g)
				g.level = 2
				log.Println(p.name + " reached level 2!")
			}
		case 400:
			if g.level < 3 {
				m.InitLevel3(g)
				g.level = 3
				log.Println(p.name + " reached level 3!")
			}
		case 600:
			if g.level < 4 {
				m.InitLevel4(g)
				g.level = 4
				log.Println(p.name + " reached level 4!")
			}
		}
	}
}

func (g *Game) readScores() {
	g.scores = nil
	f, err := os.Open(g.scoreFile)
	if err != nil {
		f, err := os.Create(g.scoreFile)
		if err != nil {
			log.Println(err, f)
		}
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Println(err)
		}
	}()

	s := bufio.NewScanner(f)
	for s.Scan() {
		score := strings.Split(s.Text(), ":")
		g.scores = append(g.scores, score)
	}
	err = s.Err()
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
	}
}

func (g *Game) writeScores() {
	f, err := os.OpenFile(g.scoreFile, os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			//log.Fatal(err)
			fmt.Println(err)
		}
	}()
	for _, v := range g.scores {
		_, err := fmt.Fprintln(f, strings.Join(v[:], ":"))
		if err != nil {
			log.Println(err)
		}
	}
}

func (g *Game) checkScores() bool {
	scoreChange := false
	var newScores [][]string
	if g.scores != nil {
		for _, p := range g.players {
			for i, s := range g.scores {
				scoreStr, err := strconv.Atoi(s[2])
				if err != nil {
					log.Println(err)
				}
				numPlayers, err := strconv.Atoi(s[0])
				if err != nil {
					log.Println(err)
				}
				if p.score > scoreStr && numPlayers == g.numPlayers {
					var newScore []string
					scoreChange = true
					newScore = append(newScore, strconv.Itoa(g.numPlayers), p.name, strconv.Itoa(p.score))
					for a := 0; a < i; a++ {
						newScores = append(newScores, g.scores[a])
					}
					newScores = append(newScores, newScore)
					if i <= len(g.scores)-1 {
						for a := i; a < len(g.scores); a++ {
							newScores = append(newScores, g.scores[a])
						}
					}
					break
				} else if len(g.scores) < MaxHighScores && numPlayers == g.numPlayers {
					var newScore []string
					scoreChange = true
					newScore = append(newScore, strconv.Itoa(g.numPlayers), p.name, strconv.Itoa(p.score))
					newScores = append(g.scores, newScore)
					break
				}
			}
		}
		if scoreChange {
			g.scores = nil
			if len(newScores) > MaxHighScores {
				for i := 0; i < MaxHighScores; i++ {
					g.scores = append(g.scores, newScores[i])
				}
			} else {
				for i := 0; i < len(newScores); i++ {
					g.scores = append(g.scores, newScores[i])
				}
			}
		}
	} else {
		for _, p := range g.players {
			var newScore []string
			scoreChange = true
			newScore = append(newScore, strconv.Itoa(g.numPlayers), p.name, strconv.Itoa(p.score))
			g.scores = append(g.scores, newScore)
		}
	}
	return scoreChange
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
	ms := 80 //120
	switch direction {
	case DirUp, DirDown:
		ms = 140 //180
	}
	//ms -= (speed / 100)
	return time.Duration(ms) * time.Millisecond
}

// Remove bit from map
func (g *Game) removeBit(i int) {
	b := &Bit{}
	g.bits[i] = g.bits[len(g.bits)-1]
	g.bits[len(g.bits)-1] = b
	g.bits = g.bits[:len(g.bits)-1]
}

func GetStyle(bg tcell.Color, fg tcell.Color) tcell.Style {
	style := tcell.StyleDefault.
		Background(bg).
		Foreground(fg)

	return style
}
