package main

import (
	"os"
	"strconv"
	"time"

	"github.com/google/logger"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/gdamore/tcell/views"
)

const (
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

	// High Score count
	MaxHighScores = 5

	// Game runes
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
)

var (
	// Current game map
	m *GameMap

	// Number of random bits that should be present on map at a time
	numBits int = 5

	// Text to be displayed at bottom for controls
	controls        string = "w/s/a/d = up/down/left/right - esc = quit - f1 = restart - f12 = pause"
	mainOptions            = []string{"Play", "High Scores", "Settings"}
	playerOptions          = []string{"1 Player", "2 Player"}
	gameModeOptions        = []string{"Basic", "Advanced", "Battle"}
	PlayerRunes            = []rune{'█', '■', '◆', '࿖', 'ᚙ', '▚', 'ↀ', 'ↈ', 'ʘ', '֍', '߷', '⁂', 'O', 'o', '=', '#', '$'}
	PlayerColors           = []string{"green", "black", "navy", "silver", "purple", "teal", "red", "blue", "lime", "yellow", "fuchsia", "aqua", "white"}
)

// Game is the main game struct and is used to store and compute general game logic.
type Game struct {

	// Screen and views
	screen tcell.Screen    // Main Screen
	gview  *views.ViewPort // Game view port
	sview  *views.ViewPort // Controls view port
	sbar   *views.TextBar  // Controls text bar

	// Game structs
	players  []*Player // All players in game
	entities []*Entity // All entities currently in game
	bites    []*Bite   // All bites currently  in game (triangles)
	bits     []*Bit    // All bits currently in game (square dots)
	gameMap  *GameMap  // Game map
	biteMap  *GameMap  // Bite map
	style    *Style    // The game's current color styles

	// Score and profile tracking
	scores1     []*Score   // 1 player scores
	scores2     []*Score   // 2 player scores
	scoreFile   string     // File that stores the scores
	profiles    []*Profile // Player profiles
	curProfiles []*Profile // Currently selected profiles
	proFile     string     // File that stores the profiles

	// Misc variables
	state      int       // Game state
	mode       int       // Game mode
	level      int       // Current game level
	numPlayers int       // Chosen number of players for game
	fps        int       // Game FPS
	frames     int       // Used to track game FPS
	bitQuit    chan bool // Used to close handlebits goroutine
}

type Data interface {
	Encode() []byte
	Decode()
}

// InitScreen initializes the tcell screen and sets views/bars and styles.
func (g *Game) InitScreen() error {

	// Set style
	s := SetDefaultStyle()
	g.style = s

	encoding.Register()

	// Prepare screen
	if screen, err := tcell.NewScreen(); err != nil {
		logger.Errorf("Failed to create screen: %v", err)
		os.Exit(1)
	} else if err = screen.Init(); err != nil {
		logger.Errorf("Failed to initialize screen: %v", err)
		os.Exit(1)
	} else {
		screen.SetStyle(g.style.DefStyle)
		g.screen = screen
		logger.Info("Intialized screen...")
	}

	// Display cursor at bottom of screen. Seems to be an issue with
	// Windows Terminal and hiding the cursor completely
	g.screen.ShowCursor(0, MapHeight+3)

	// Create the main game viewport
	g.gview = views.NewViewPort(g.screen, MapStartX, MapStartY, MapWidth, MapHeight)

	// Create the secondary view port and text bars for the controls display
	g.sview = views.NewViewPort(g.screen, SViewStartX, SViewStartY, SViewWidth, SViewHeight)
	g.sbar = views.NewTextBar()
	g.sbar.SetView(g.sview)
	g.sbar.SetStyle(g.style.DefStyle)

	return nil
}

// MainMenu displays and handles input for the Main Menu.
func (g *Game) MainMenu() error {

	// Setup main menu
	g.state = MainMenu
	cMenu := MenuMain

	// Read high scores from scoreFile
	g.getScores()

	// Run main menu until play or quit
	for g.state != Play {
		// Display the "Main Menu" menu
		if cMenu == MenuMain {
			cMenu = g.MenuMain()
			logger.Infof("%v", cMenu)
		}
		// Display the Player number choice menu to decide
		// how many players will be playing
		if cMenu == MenuPlayer {
			cMenu = g.MenuPlayer()
			logger.Infof("%v", cMenu)
		}
		// Display the player profile menu and let players pick their
		// profile or create a new one.
		if cMenu == MenuProfile {
			cMenu = g.MenuProfile(cMenu)
		}
		// Display the high score screen
		if cMenu == MenuScore {
			cMenu = g.MenuScore(cMenu)
		}
	}
	return nil
}

func (g *Game) MenuMain() int {
	var cMenu int
	g.gview.Clear()
	i := g.handleMenu(mainOptions)
	switch i {
	case -1:
		g.screen.Fini()
		os.Exit(0)
	case 0:
		return MenuPlayer
	case 1:
		return MenuScore
	}
	return cMenu
}

func (g *Game) MenuPlayer() int {
	var cMenu int
	g.gview.Clear()
	i := g.handleMenu(playerOptions)
	switch i {
	case -1:
		return MenuMain
	case 0:
		g.numPlayers = 1
		g.mode = Player1
		return MenuProfile
	case 1:
		g.numPlayers = 2
		g.mode = Player2
		return MenuProfile
	}

	return cMenu
}

func (g *Game) MenuProfile(cMenu int) int {
	for cMenu == MenuProfile {
		g.gview.Clear()
		g.curProfiles = nil

		for a := 0; a < g.numPlayers; a++ {
			var profileList []string
			var pNum string

			// Read profiles from file and add to list to create menu items
			byteData := ReadFile(g.proFile)
			g.profiles = DecodeProfiles(byteData)
			for i := range g.profiles {
				profileList = append(profileList, g.profiles[i].Name)
			}

			// Add an entry for creating a new profile
			profileList = append(profileList, "New Profile")

			// If 1 Player mode don't show a number, if 2 player then
			// show which player number during profile select
			if g.numPlayers > 1 {
				pNum = strconv.Itoa(a + 1)
			} else {
				pNum = ""
			}
			// Draw the Select Profile text
			g.screen.Clear()
			renderCenterStr(g.gview, MapWidth, MapHeight-4, g.style.DefStyle, ("  Select Profile " + pNum + ":"))
			g.screen.Show()

			// Draw and handle the player select menu. The list of menu items
			// is generated using the list of profiles read from file.
			i := g.handleMenu(profileList)

			// Drop back to MenuMain if Escape is pressed
			if i == -1 {
				return MenuMain
				// If any of the profiles are selected then add them to the current profile list
				// and either proceed to to InitGame or continue loop for second player
			} else if i < (len(profileList) - 1) {
				g.curProfiles = append(g.curProfiles, g.profiles[i])
				g.state = Play
				if a == g.numPlayers-1 {
					cMenu = MenuMain
				}
				continue
				// If "New Profile" is selected then run getPlayerName to get a name and
				// create a profile from that name
			} else {
				i := CreateProfile(g)
				if i == MenuMain {
					break
				}
			}
		}
	}
	return cMenu
}

func (g *Game) MenuScore(cMenu int) int {
	for cMenu == MenuScore {
		g.screen.Clear()
		renderHighScoreScreen(g, g.style.DefStyle, MaxHighScores)

		// Wait for Escape key to be pressed to return to Main Menu
		ev := g.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				return MenuMain
			}
		}
	}
	return cMenu
}

// InitMap generates new maps for the game.
func (g *Game) InitMap() error {

	// Initialize game states
	g.level = 1

	// Create a game map
	m = &GameMap{
		Width:  MapWidth,
		Height: MapHeight,
		X:      MapStartX,
		Y:      MapStartY,
	}
	g.gameMap = m
	m.InitMap()
	m.InitMapBoundary(WallRune, FloorRune, g.style.DefStyle)
	m.InitLevel1(g)
	logger.Info("Created game map and set to level 1.")

	biteMap := &GameMap{
		Width:  m.Width,
		Height: m.Height,
	}
	biteMap.InitMap()
	biteMap.InitMapBoundary(WallRune, FloorRune, g.style.DefStyle)
	g.biteMap = biteMap

	return nil
}

// InitPlayers creates player objects for the game.
func (g *Game) InitPlayers() error {
	// Set player starting x value to middle of map
	x := MapWidth / 2

	// Create a player for selected number of players
	for i := 0; i < g.numPlayers; i++ {
		y := (MapHeight / 2) + (i * 2)

		// Get player vars from loaded profile
		pName := g.curProfiles[i].Name
		pStyle := g.curProfiles[i].GetStyle()
		pChar := g.curProfiles[i].Char

		// Create player and
		p := NewPlayer(x, y, 0, (DirLeft - i), pChar, pName, pStyle)
		g.players = append(g.players, p)
	}
	g.players[0].score = 0
	for i := 0; i < numBits; i++ {
		b := NewRandomBit(m, 10, BitRune, g.style.BitStyle)
		g.bits = append(g.bits, b)
	}
	logger.Info("Initialized game with ", strconv.Itoa(g.numPlayers), " players.")

	return nil
}

// Run is the main game loop.
func (g *Game) Run() error {

	// Run a goroutine for each player to handle their own loop
	// separately from each other and the main game loop
	for _, p := range g.players {
		p.quitChan = make(chan bool)
		p.bitChan = make(chan int, 1)
		p.biteChan = make(chan int, 1)
		go g.handlePlayer(p)
	}

	// The gameplay loop
	for g.state == Play || g.state == Pause {

		// Handle entities and objects on level
		g.handleLevel(m)

		// Run goroutine for player's input
		go handleInput(g)

		// Handle game if pause button is pressed
		g.handlePause()

		for i := range g.players {
			select {
			case bitPos := <-g.players[i].bitChan:
				g.removeBit(bitPos)
			case bitePos := <-g.players[i].biteChan:
				g.removeBite(bitePos)
			default:
				continue
			}
		}

		// Render the game
		renderAll(g, g.style.DefStyle, m)

		// Keep track of FPS
		g.getFPS()
		g.frames++
	}

	// If game ends then kill the handlePlayer goroutines
	for _, p := range g.players {
		p.quitChan <- true
	}

	return nil
}

// Quit completely exits the game back to terminal.
func (g *Game) Quit() {
	g.state = Quit
	g.screen.Fini()
	logger.Info("Quitting the game...")
	os.Exit(0)
}

// Return quits the current game and returns to the Main Menu.
func (g *Game) Return() {
	g.state = MainMenu
	g.screen.Fini()
}

// Restart resets the game in the same game mode with same players.
func (g *Game) Restart() {
	g.state = Restart
	logger.Info("Restarting the game...")
	g.screen.Fini()
}

// handleMenu renders the menu screens and keeps track of which
// menu item is currently selected and which to move to
// based on input.
func (g *Game) handleMenu(options []string) int {
	choice := 0
	m := NewMainMenu(options, g.style.DefStyle, g.style.SelStyle, 0)
	m.SetSelected(0)
	m.ChangeSelected()
	for choice == 0 {
		renderMenu(g, m, g.style.DefStyle)
		choice = handleMenuInput(g, m)
	}
	if choice == 1 {
		choice = m.GetSelected()
	}
	return choice
}

// handlePlayer is the player loop and handles a player's
// state and  interaction with objects and the game map.
func (g *Game) handlePlayer(p *Player) {

	// Continuously loop unless killed through p.ch channel
	for {
		select {
		default:
			scoreChange := false

			// Check which direction player should be moving
			dx, dy := p.CheckDirection(g)

			// Check if player is blocked at all
			if p.IsBlocked(m, g.biteMap, g.entities, g.players, dx, dy) {

				// Run if in 2 player mode
				if g.numPlayers > 1 {

					// Generate bits where player's body was during collision
					for _, i := range p.pos {
						b := NewBit(i.ox, i.oy, 10, BitRune, BitRandom, g.style.BitStyle)
						g.bits = append(g.bits, b)
					}

					// Read high scores from file, compare against current scores
					// and make changes if necessary
					g.getScores()
					g.scores2, scoreChange = UpdateScores(g.scores2, p.name, p.score, g.mode, MaxHighScores)
					if scoreChange {
						WriteScores(g.scores1, g.scores2, g.scoreFile)
					}

					// Reset the player
					p.Reset(MapWidth/2, MapHeight/2, 3, g.style.BiteExplodedStyle)
					logger.Infof("Player died: %v", p.name)

					// Run if in 1 player mode
				} else {

					// Kill player
					p.Kill(g.style.BiteExplodedStyle)
					logger.Infof("Player died: %v", p.name)

					// Read high scores from file, compare against current scores
					// and make changes if necessary
					g.getScores()
					g.scores1, scoreChange = UpdateScores(g.scores1, p.name, p.score, g.mode, MaxHighScores)
					if scoreChange {
						WriteScores(g.scores1, g.scores2, g.scoreFile)
					}

					// Wait a short period of time then restart the game
					time.Sleep(100 * time.Millisecond)
					g.screen.Fini()
					g.state = Restart
				}
			} else {
				// Move player if not blocked
				p.Move(dx, dy)
			}
			// Check if player is on a bit or bite
			bitPos := p.IsOnBit(g)
			if bitPos != -1 {
				p.bitChan <- bitPos
			}
			bitePos := p.IsOnBite(g, m)
			if bitePos != -1 {
				p.biteChan <- bitePos
			}

			// Calculate player's speed based on their score.
			// Movement is done by causing the player goroutine
			// to sleep for a set amount of time.
			//p.speed += p.score / 200
			time.Sleep(g.moveInterval(0, p.GetDirection()))

		// Quit goroutine if signaled
		case <-p.quitChan:
			return
		}
	}
}

// handleBits causes bits on map to move in a random direction in timed intervals.
func (g *Game) handleBits(m *GameMap) {
	for {
		select {
		default:
			// Move bits in a random direction after a set amount of time
			for i := range g.bits {
				switch g.bits[i].state {
				case BitRandom:
					g.bits[i].Move(m)
				}
			}
			// Wait a set amount of time
			time.Sleep(500 * time.Millisecond)

		// Quit goroutine if signaled
		case <-m.BitChan:
			return
		}
	}
}

// handleLevel checks the current score against the current level and
// changes the level if a certain score is reached.
func (g *Game) handleLevel(m *GameMap) {
	for _, p := range g.players {
		// Level 2
		if p.score >= Level2 {
			if g.level < 2 {
				m.InitLevel2(g)
				g.level = 2
				logger.Info(p.name + " reached level 2!")
			}
		}
		// Level 3
		if p.score >= Level3 {
			if g.level < 3 {
				m.InitLevel3(g)
				g.level = 3
				logger.Info(p.name + " reached level 3!")
			}
		}
		// Level 4
		if p.score >= Level4 {
			if g.level < 4 {
				m.InitLevel4(g)
				g.level = 4
				logger.Info(p.name + " reached level 4!")
			}
		}
		// Level 5
		if p.score >= Level5 {
			if g.level < 5 {
				m.InitLevel5(g)
				g.level = 5
				logger.Info(p.name + " reached level 5!")
			}
		}
		// Level 6
		if p.score >= Level6 {
			if g.level < 6 {
				m.InitLevel6(g)
				g.level = 6
				logger.Info(p.name + " reached level 6!")
			}
		}
	}
}

// handlePause controls the game and input during the "paused" state.
func (g *Game) handlePause() {
	chQuit := false

	// If pause is called kill the player goroutines
	for g.state == Pause {
		if !chQuit {
			for _, p := range g.players {
				p.quitChan <- true
				chQuit = true
			}
		}
		logger.Info("Pausing game...")

		// Render "PAUSED" to screen
		renderCenterStr(g.gview, MapWidth, MapHeight-4, g.style.BitStyle, "PAUSED")
		g.screen.Show()

		// If unpaused then restart player and bit goroutines
		if g.state == Play {
			for _, p := range g.players {
				go g.handlePlayer(p)
			}
			go g.handleBits(m)
			logger.Info("Resuming game...")
		}
	}
}

// getFPS tracks variables used to calculate the FPS of the game.
func (g *Game) getFPS() {
	time.AfterFunc(1*time.Second, func() {
		g.fps = g.frames
		g.frames = 0
	})
}

// moveInterval calculates a player's movement speed. Player speed is created by
// sleeping the player's loop for a set amount of time. The up and down direction
// have a decrease in speed in an attempt to even out the direction movement speeds.
// because of the way the terminal is designed vertical movement is normally much
// faster than horizontal.
func (g *Game) moveInterval(speed, direction int) time.Duration {
	ms := 80 //120
	switch direction {
	case DirUp, DirDown:
		ms = 140 //180
	}
	//ms -= (speed / 100)
	return time.Duration(ms) * time.Millisecond
}

// removeBit removes a particular bit from the game's bit slice in order to remove
// that bit from the game.
func (g *Game) removeBit(i int) {
	g.bits[i] = g.bits[len(g.bits)-1]
	g.bits[len(g.bits)-1] = nil
	g.bits = g.bits[:len(g.bits)-1]
}

// removeBit removes a particular bit from the game's bit slice in order to remove
// that bit from the game.
func (g *Game) removeBite(i int) {
	g.bites[i] = g.bites[len(g.bites)-1]
	g.bites[len(g.bites)-1] = nil
	g.bites = g.bites[:len(g.bites)-1]
}

// getScores reads scores from the game's scoreFile and stores them in it's
// scores1 and scores2 variables.
func (g *Game) getScores() {
	byteData := ReadFile(g.scoreFile)
	g.scores1, g.scores2 = DecodeScores(byteData)
	logger.Infof("Loaded high scores from file: %v", g.scoreFile)
}

func getCharList(list []rune) []string {
	var charList []string
	for i := range list {
		char, err := strconv.Unquote(strconv.QuoteRune(PlayerRunes[i]))
		if err != nil {
			logger.Errorf("Error removing quotes: %v", err)
		}
		charList = append(charList, char)
	}
	return charList
}
