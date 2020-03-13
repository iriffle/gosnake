package game

import (
	"os"
	"strconv"
	"time"

	"github.com/google/logger"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/gdamore/tcell/views"
	"github.com/stjiub/gosnake/entity"
	"github.com/stjiub/gosnake/gamemap"
	"github.com/stjiub/gosnake/style"
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
	m *gamemap.GameMap

	// Number of random bits that should be present on map at a time
	numBits int = 5

	// Text to be displayed at bottom for controls
	controls        string = "w/s/a/d = up/down/left/right - esc = quit - f1 = restart - f12 = pause"
	mainOptions            = []string{"Play", "High Scores", "Settings"}
	playerOptions          = []string{"1 Player", "2 Player"}
	gameModeOptions        = []string{"Basic", "Advanced", "Battle"}
	PlayerRunes            = []rune{'█', '■', '◆', '࿖', 'ᚙ', '▚', 'ↀ', 'ↈ', 'ʘ', '֍', '߷', '⁂', 'O', 'o', '=', '#', '$', '+', '-', '!', '('}
	PlayerColors           = []string{"white", "black", "silver", "green", "lime", "blue", "navy", "aqua", "teal", "red", "purple", "fuschia"}
	BiteRunes              = []rune{BiteUpRune, BiteDownRune, BiteLeftRune, BiteRightRune, BiteAllRune, BiteExplodeRune}
)

// Game is the main game struct and is used to store and compute general game logic.
type Game struct {

	// Screen and views
	screen tcell.Screen    // Main Screen
	gview  *views.ViewPort // Game view port
	sview  *views.ViewPort // Controls view port
	sbar   *views.TextBar  // Controls text bar

	// Game structs
	players  []*entity.Player // All players in game
	entities []*entity.Entity // All entities currently in game
	bites    []*entity.Bite   // All bites currently  in game (triangles)
	bits     []*entity.Bit    // All bits currently in game (square dots)
	items    []*entity.Item
	gameMap  *gamemap.GameMap // Game map
	biteMap  *gamemap.GameMap // Bite map

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

	style.Style
}

func NewGame(numPlayers int, curProfiles []*Profile, scoreFile, proFile string) *Game {
	g := Game{
		numPlayers:  numPlayers,
		curProfiles: curProfiles,
		scoreFile:   scoreFile,
		proFile:     proFile,
	}

	return &g
}

// InitScreen initializes the tcell screen and sets views/bars and styles.
func (g *Game) InitScreen() error {

	// Set style
	g.SetDefaultStyle()

	encoding.Register()

	// Prepare screen
	if screen, err := tcell.NewScreen(); err != nil {
		logger.Errorf("Failed to create screen: %v", err)
		os.Exit(1)
	} else if err = screen.Init(); err != nil {
		logger.Errorf("Failed to initialize screen: %v", err)
		os.Exit(1)
	} else {
		screen.SetStyle(g.DefStyle)
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
	g.sbar.SetStyle(g.DefStyle)

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
	renderSnakeLogo(g, MapWidth/2, MapHeight/2)
	renderGoLogo(g, MapWidth/2, MapHeight/2)
	i := g.handleMenu(mainOptions)
	switch i {
	case ItemExit:
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
	renderSnakeLogo(g, MapWidth/2, MapHeight/2)
	renderGoLogo(g, MapWidth/2, MapHeight/2)
	i := g.handleMenu(playerOptions)
	switch i {
	case ItemExit:
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
	for cMenu == MenuProfile || cMenu == MenuEdit || cMenu == MenuRemove {
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
			if cMenu == MenuProfile {
				profileList = append(profileList, "New Profile", "Edit Profile ", "Remove Profile ")
			}

			// Draw the Select Profile text
			g.screen.Clear()
			renderSnakeLogo(g, MapWidth/2, MapHeight/2)
			renderGoLogo(g, MapWidth/2, MapHeight/2)
			if cMenu == MenuProfile {
				// If 1 Player mode don't show a number, if 2 player then
				// show which player number during profile select
				if g.numPlayers > 1 {
					pNum = strconv.Itoa(a + 1)
				} else {
					pNum = ""
				}
				renderCenterStr(g.gview, MapWidth, MapHeight-4, g.DefStyle, ("  Select Profile " + pNum + ":"))
			} else if cMenu == MenuEdit {
				renderCenterStr(g.gview, MapWidth, MapHeight-4, g.DefStyle, ("  Edit Profile:"))
			} else {
				renderCenterStr(g.gview, MapWidth, MapHeight-4, g.DefStyle, ("  Remove Profile:"))
			}
			g.screen.Show()

			// Draw and handle the player select menu. The list of menu items
			// is generated using the list of profiles read from file.
			if len(profileList) > 0 {
				i := g.handleMenu(profileList)

				// Drop back to MenuMain if Escape is pressed
				if i == ItemExit {
					return MenuMain
					// If any of the profiles are selected then add them to the current profile list
					// and either proceed to to InitGame or continue loop for second player
				} else if i < (len(profileList)-3) && cMenu == MenuProfile {
					g.curProfiles = append(g.curProfiles, g.profiles[i])
					g.state = Play
					if a == g.numPlayers-1 {
						cMenu = MenuMain
					}
					continue
					// If "New Profile" is selected then run getPlayerName to get a name and
					// create a profile from that name
				} else if i == (len(profileList)-3) && cMenu == MenuProfile {
					i := CreateProfile(g)
					if i == MenuMain {
						break
					}
				} else if i == (len(profileList)-2) && cMenu == MenuProfile {
					_ = g.MenuProfile(MenuEdit)
					return MenuProfile
				} else if i == (len(profileList)-1) && cMenu == MenuProfile {
					_ = g.MenuProfile(MenuRemove)
					return MenuProfile
				} else if cMenu == MenuEdit {
					_ = EditProfile(g, g.profiles[i])
				} else if cMenu == MenuRemove {
					_ = RemoveProfile(g, i)
				}
			} else {
				return MenuProfile
			}
		}
	}
	return cMenu
}

func (g *Game) MenuScore(cMenu int) int {
	for cMenu == MenuScore {
		g.screen.Clear()
		renderHighScoreScreen(g, g.DefStyle, MaxHighScores)

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
	m = &gamemap.GameMap{
		Width:  MapWidth,
		Height: MapHeight,
		X:      MapStartX,
		Y:      MapStartY,
	}
	g.gameMap = m
	m.InitMap()
	m.InitMapBoundary(WallRune, FloorRune, g.DefStyle)
	InitLevel1(g)
	logger.Info("Created game map and set to level 1.")

	biteMap := &gamemap.GameMap{
		Width:  m.Width,
		Height: m.Height,
	}
	biteMap.InitMap()
	biteMap.InitMapBoundary(WallRune, FloorRune, g.DefStyle)
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
		p := entity.NewPlayer(x, y, 0, (DirLeft - i), pChar, pName, pStyle)
		g.players = append(g.players, p)
	}
	g.players[0].SetScore(0)
	for i := 0; i < numBits; i++ {
		b := entity.NewRandomBit(m, 10, BitRune, g.BitStyle)
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
		p.InitChans()
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
			case bitPos := <-g.players[i].BitChan:
				g.removeBit(bitPos)
			case bitePos := <-g.players[i].BiteChan:
				g.removeBite(bitePos)
			default:
				continue
			}
		}

		// Render the game
		renderAll(g, g.DefStyle, m)

		// Keep track of FPS
		g.getFPS()
		g.frames++
	}

	// If game ends then kill the handlePlayer goroutines
	for _, p := range g.players {
		p.QuitChan <- true
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
	m := NewMainMenu(options, g.DefStyle, g.SelStyle, 0)
	m.SetSelected(0)
	m.ChangeSelected()
	for choice == 0 {
		renderMenu(g, m, g.DefStyle)
		choice = handleMenuInput(g, m)
	}
	if choice == 1 {
		choice = m.GetSelected()
	}
	return choice
}

// handlePlayer is the player loop and handles a player's
// state and  interaction with objects and the game map.
func (g *Game) handlePlayer(p *entity.Player) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Error in handlePlayer goroutine: %v", r)
		}
	}()
	// Continuously loop unless killed through p.ch channel
	for {
		select {
		default:
			scoreChange := false

			// Check which direction player should be moving
			dx, dy := p.CheckDirection()

			// Check if player is blocked at all
			if p.IsBlocked(m, g.biteMap, g.entities, g.players, dx, dy) {
				name := p.GetName()
				score := p.GetScore()
				// Read high scores from file, compare against current scores
				// and make changes if necessary
				g.getScores()
				g.scores2, scoreChange = UpdateScores(g.scores2, name, score, g.mode, MaxHighScores)
				if scoreChange {
					WriteScores(g.scores1, g.scores2, g.scoreFile)
				}

				// Run if in 2 player mode
				if g.numPlayers == 1 {
					// Kill player
					p.Kill(g.BiteExplodedStyle)
					logger.Infof("Player died: %v", name)
					// Wait a short period of time then restart the game
					time.Sleep(100 * time.Millisecond)
					g.Restart()
				} else {
					g.bits = p.DropBits(g.bits, BitRune, BitRandom, g.DefStyle)
					p.Reset(MapWidth/2, MapHeight/2, 3, g.BiteExplodedStyle)
					logger.Infof("Player died: %v", name)
				}

			} else {
				// Move player if not blocked
				p.Move(dx, dy)
			}
			// Check if player is on a bit or bite
			bitPos := g.IsOnBit(p)
			if bitPos != -1 {
				p.BitChan <- bitPos
			}
			bitePos := g.IsOnBite(p, m)
			if bitePos != -1 {
				p.BiteChan <- bitePos
			}
			g.IsOnItem(p)

			// Calculate player's speed based on their score.
			// Movement is done by causing the player goroutine
			// to sleep for a set amount of time.
			//p.speed += p.score / 200
			time.Sleep(g.moveInterval(0, p.GetDirection()))

		// Quit goroutine if signaled
		case <-p.QuitChan:
			return
		}
	}
}

// handleBits causes bits on map to move in a random direction in timed intervals.
func (g *Game) moveBits(m *gamemap.GameMap) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Error in handleBits goroutine: %v", r)
		}
	}()
	for {
		select {
		default:
			// Move bits in a random direction after a set amount of time
			for i := range g.bits {
				state := g.bits[i].GetState()
				switch state {
				case BitRandom:
					g.bits[i].MoveRandom(m)
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
func (g *Game) handleLevel(m *gamemap.GameMap) {
	for _, p := range g.players {
		score := p.GetScore()
		name := p.GetName()
		// Level 2
		if score >= Level2 {
			if g.level < 2 {
				InitLevel2(g)
				g.level = 2
				logger.Info(name + " reached level 2!")
			}
		}
		// Level 3
		if score >= Level3 {
			if g.level < 3 {
				InitLevel3(g)
				g.level = 3
				logger.Info(name + " reached level 3!")
			}
		}
		// Level 4
		if score >= Level4 {
			if g.level < 4 {
				InitLevel4(g)
				g.level = 4
				logger.Info(name + " reached level 4!")
			}
		}
		// Level 5
		if score >= Level5 {
			if g.level < 5 {
				InitLevel5(g)
				g.level = 5
				logger.Info(name + " reached level 5!")
			}
		}
		// Level 6
		if score >= Level6 {
			if g.level < 6 {
				InitLevel6(g)
				g.level = 6
				logger.Info(name + " reached level 6!")
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
				p.QuitChan <- true
				chQuit = true
			}
		}
		logger.Info("Pausing game...")

		// Render "PAUSED" to screen
		renderCenterStr(g.gview, MapWidth, MapHeight-4, g.BitStyle, "PAUSED")
		g.screen.Show()

		// If unpaused then restart player and bit goroutines
		if g.state == Play {
			for _, p := range g.players {
				go g.handlePlayer(p)
			}
			go g.moveBits(m)
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

// Check if player is on top of a bit
func (g *Game) IsOnBit(p *entity.Player) int {
	i := p.CheckBitPos(g.bits)
	if i != -1 {
		b := g.bits[i]
		points := b.GetPoints()
		char := p.GetChar(0)
		style := p.GetStyle(0)
		p.AddScore(points)
		p.AddSegment(1, char, style)
	}
	return i
}

// Determine if player is on a bite and if so trigger explosion
func (g *Game) IsOnBite(p *entity.Player, m *gamemap.GameMap) int {
	i := p.CheckBitePos(g.bites)
	if i != -1 {
		b := g.bites[i]
		char := p.GetChar(0)
		style := p.GetStyle(0)
		p.AddScore(50)
		p.AddSegment(4, char, style)
		go b.ExplodeBite(m, g.biteMap, BiteExplodeRune, g.BiteExplodedStyle, g.DefStyle)
		return i
	}
	return -1
}

func (g *Game) IsOnItem(p *entity.Player) {
	i := p.CheckItemPos(g.items)
	if i != -1 {
		p.AddItem(g.items[i])
		g.removeItem(i)
	}
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

func (g *Game) removeItem(i int) {
	g.items[i] = g.items[len(g.items)-1]
	g.items[len(g.items)-1] = nil
	g.items = g.items[:len(g.items)-1]
}

// getScores reads scores from the game's scoreFile and stores them in it's
// scores1 and scores2 variables.
func (g *Game) getScores() {
	byteData := ReadFile(g.scoreFile)
	g.scores1, g.scores2 = DecodeScores(byteData)
	logger.Infof("Loaded high scores from file: %v", g.scoreFile)
}

func (g *Game) GetState() int {
	return g.state
}

func (g *Game) GetNumPlayers() int {
	return g.numPlayers
}

func (g *Game) GetCurProfiles() []*Profile {
	return g.curProfiles
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
