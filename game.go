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

	// Text to be displayed at bottom for controls
	controls        string = "w/s/a/d = up/down/left/right - esc = quit - f1 = restart - f12 = pause"
	mainOptions            = []string{"Play", "High Scores", "Settings"}
	playerOptions          = []string{"1 Player", "2 Player"}
	gameModeOptions        = []string{"Basic", "Advanced", "Battle"}

	// Number of bits that should be present on map
	numBits int = 5
)

// Main game struct
type Game struct {

	// Screen and views
	screen tcell.Screen    // Main Screen
	gview  *views.ViewPort // Game view port
	sview  *views.ViewPort // Controls view port
	sbar   *views.TextBar  // Controls text bar
	//cview  *views.ViewPort // Console view port
	//cbar   *views.TextBar  // Console text bar

	// Game objects
	players  []*Player  // All players in game
	entities []*Entity  // All entities currently in game
	bites    []*Bite    // All bites currently  in game (triangles)
	bits     []*Bit     // All bits currently in game (square dots)
	maps     []*GameMap // All current game maps, including maps used for bites
	style    *Style     // Colors usable for players

	// Score tracking
	scores    [][]string // 1 player scores
	scores2   [][]string // 2 player scores
	scoreFile string     // File that scores are saved to

	// Misc variables
	state      int // Game state
	mode       int // Game mode
	level      int // Current game level
	numPlayers int // Chosen number of players for game
	fps        int // Game FPS
	frames     int
	bitQuit    chan bool // Used to close handlebits goroutine
}

// Initialize the screen and set views/bars and styles
func (g *Game) InitScreen() {

	// Set style
	s := SetDefaultStyle()
	g.style = s

	// Prepare screen
	encoding.Register()

	if screen, err := tcell.NewScreen(); err != nil {
		log.Println("Failed to create screen: ", err)
		os.Exit(1)
	} else if err = screen.Init(); err != nil {
		log.Println("Failed to initialize screen: ", err)
		os.Exit(1)
	} else {
		screen.SetStyle(g.style.DefStyle)
		g.screen = screen
	}

	// Enable mouse support. Not currently used
	// if g.screen.HasMouse() {
	// 	g.screen.EnableMouse()
	// }

	// Display cursor at bottom of screen. Seems to be an issue with
	// Windows Terminal and hiding the cursor completely
	g.screen.ShowCursor(CViewStartX, CViewStartY)

	// Create the main game viewport
	g.gview = views.NewViewPort(g.screen, MapStartX, MapStartY, MapWidth, MapHeight)

	// Create the secondary view port and text bars for the controls display
	g.sview = views.NewViewPort(g.screen, SViewStartX, SViewStartY, SViewWidth, SViewHeight)
	g.sbar = views.NewTextBar()
	g.sbar.SetView(g.sview)
	g.sbar.SetStyle(g.style.DefStyle)
}

// Launch main menu screen
func (g *Game) MainMenu() {

	// Setup main menu
	g.state = MainMenu
	cMenu := MenuMain

	// Read high scores from scoreFile
	g.readScores()

	// Run main menu until play or quit
	for g.state != Play {

		// Display the "Main Menu" menu
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

		// Display the Player number choice menu to decide
		// how many players will be playing
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

		// Display the high score screen
		for cMenu == MenuScore {
			renderHighScoreScreen(g, g.style.DefStyle)

			// Wait for Escape key to be pressed to return to Main Menu
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

	// Initialize game states
	g.state = Play
	g.level = 1

	// Create a game map
	m = &GameMap{
		Width:  MapWidth,
		Height: MapHeight,
		X:      MapStartX,
		Y:      MapStartY,
	}
	g.maps = append(g.maps, m)
	m.InitMap()
	m.InitMapBoundary(WallRune, FloorRune, g.style.DefStyle)
	m.InitLevel1(g)

	// Set player starting x value to middle of map
	x := MapWidth / 2

	// Create a player for selected number of players
	for i := 0; i < g.numPlayers; i++ {
		pName := ""
		y := (MapHeight / 2) + (i * 2)

		for pName == "" {
			pName = g.getPlayerName(i+1, m.Width, m.Height, g.style.DefStyle, g.style.SelStyle)
		}
		if pName == "-quit-" {
			g.QuitToMenu()
			return
		}

		pStyle := g.style.PlayerColors[i]
		p := NewPlayer(x, y, 0, (DirLeft - i), PlayerRune, pName, pStyle)
		g.players = append(g.players, &p)
	}
	g.players[0].score = 0
	for i := 0; i < numBits; i++ {
		b := NewRandomBit(m, 10, BitRune, g.style.BitStyle)
		g.bits = append(g.bits, &b)
	}
	log.Println("Initialized game with ", strconv.Itoa(g.numPlayers), " players.")
}

// Run the actual game
func (g *Game) RunGame() {

	// Run a goroutine for each player to handle their own loop
	// separately from each other and the main game loop
	for _, p := range g.players {
		p.ch = make(chan bool)
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

		// Render the game
		renderAll(g, g.style.DefStyle, m)

		// Keep track of FPS
		g.getFPS()
		g.frames++
	}

	// If game ends then kill the handlePlayer goroutines
	for _, p := range g.players {
		p.ch <- true
	}

	//g.bitQuit <- true
}

// Completely quit game back to terminal
func (g *Game) QuitGame() {
	g.state = Quit
	g.screen.Fini()
	os.Exit(0)
}

func (g *Game) QuitToMenu() {
	g.state = MainMenu
	g.screen.Fini()
	g.state = MainMenu
}

// Renders the menu screens and keeps track of which
// menu item is currently selected and which to move to
// based on input.
func (g *Game) handleMenu(options []string) int {
	choice := 0
	m := NewMainMenu(options, g.style.DefStyle, g.style.SelStyle, 0)
	m.SetSelected(0)
	m.ChangeSelected()
	for choice == 0 {
		renderMenu(g, &m, g.style.DefStyle)
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

	// If pause is called kill the player goroutines
	for g.state == Pause {
		if !chQuit {
			for _, p := range g.players {
				p.ch <- true
				chQuit = true
			}
		}

		// Render "PAUSED" to screen
		renderCenterStr(g.gview, MapWidth, MapHeight-4, g.style.BitStyle, "PAUSED")
		g.screen.Show()

		// If unpaused then restart player and bit goroutines
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
	var scoreChange bool

	// Continuously loop unless killed through p.ch channel
	for {
		select {
		default:

			// Check which direction player should be moving
			dx, dy := p.CheckDirection(g)

			// Check if player is blocked at all
			if p.IsBlocked(m, g.maps, g.entities, g.players, dx, dy) {

				// Run if in 2 player mode
				if g.numPlayers > 1 {

					// Generate bits where player's body was during collision
					for _, i := range p.pos {
						b := NewBit(i.ox, i.oy, 10, BitRune, BitRandom, g.style.BitStyle)
						g.bits = append(g.bits, &b)
					}

					// Read high scores from file, compare against current scores
					// and make changes if necessary
					g.readScores()
					g.scores2, scoreChange = g.checkScores()
					if scoreChange {
						g.writeScores()
					}

					// Reset the player
					p.Reset(MapWidth/2, MapHeight/2, 3, g.style.BiteExplodedStyle)

					// Run if in 1 player mode
				} else {

					// Kill player
					p.Kill(g.style.BiteExplodedStyle)

					// Read high scores from file, compare against current scores
					// and make changes if necessary
					g.readScores()
					g.scores, scoreChange = g.checkScores()
					if scoreChange {
						g.writeScores()
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
			p.IsOnBit(g)
			p.IsOnBite(g, m)

			// Calculate player's speed based on their score.
			// Movement is done by causing the player goroutine
			// to sleep for a set amount of time.
			p.speed += p.score / 200
			time.Sleep(g.moveInterval(p.speed, p.direction))

		// Quit goroutine if signaled
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
			// Move bits in a random direction after a set amount of time
			for _, bit := range g.bits {
				switch bit.state {
				case BitRandom:
					bit.Move(m)
				}
			}
			// Wait a set amount of time
			time.Sleep(500 * time.Millisecond)

		// Quit goroutine if signaled
		case <-g.bitQuit:
			return
		}
	}
}

// Change level based on player score
func (g *Game) handleLevel(m *GameMap) {
	for _, p := range g.players {
		switch p.score {

		// Level 2
		case 200:
			if g.level < 2 {
				m.InitLevel2(g)
				g.level = 2
				log.Println(p.name + " reached level 2!")
			}

		// Level 3
		case 400:
			if g.level < 3 {
				m.InitLevel3(g)
				g.level = 3
				log.Println(p.name + " reached level 3!")
			}

		// Level 4
		case 600:
			if g.level < 4 {
				m.InitLevel4(g)
				g.level = 4
				log.Println(p.name + " reached level 4!")
			}
		}
	}
}

// Read high scores from scoreFile
func (g *Game) readScores() {

	// Clear currently cached scores
	g.scores = nil
	g.scores2 = nil

	twoPlayer := false

	// Open the scoreFile
	f, err := os.Open(g.scoreFile)
	if err != nil {
		log.Println(err, f)
	}

	// Close the scoreFile on exit
	defer func() {
		if err = f.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Read scoreFile one line at a time
	s := bufio.NewScanner(f)
	for s.Scan() {

		// If hash symbol read then start appending scores
		// to the 2 player scores variable (g.scores2)
		if s.Text() == "#" {
			twoPlayer = true
			continue
		}

		// Split the scores into 3 values, split by colons
		score := strings.Split(s.Text(), ":")

		// Append the scores based on number of players
		if !(twoPlayer) {
			g.scores = append(g.scores, score)
		} else {
			g.scores2 = append(g.scores2, score)
		}
	}
	err = s.Err()
	if err != nil {
		log.Println(err)
	}
}

// Write high scores to scoreFile
func (g *Game) writeScores() {

	// Open score file overwriting any previous data
	f, err := os.OpenFile(g.scoreFile, os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
	}

	// Close scoreFile on exit
	defer func() {
		if err = f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Write the scores for 1 player and 2 player separated
	// by a hash symbol
	g.writeScore(f, g.scores)
	_, err = fmt.Fprintln(f, "#")
	if err != nil {
		log.Println(err)
	}
	g.writeScore(f, g.scores2)
}

// Loop through a slice of scores and write to file. Used by writeScores
func (g *Game) writeScore(f *os.File, scores [][]string) {
	for _, v := range scores {
		_, err := fmt.Fprintln(f, strings.Join(v[:], ":"))
		if err != nil {
			log.Println(err)
		}
	}
}

// Compare player's score against high score list to see if a new
// high score has been reached.
func (g *Game) checkScores() ([][]string, bool) {

	var (
		scores     [][]string
		newScores  [][]string
		numPlayers string
	)

	scoreChange := false

	// Determine which score slice to use based on number of players
	if g.numPlayers == 1 {
		scores = g.scores

		// numPlayers is used over g.numPlayers later when appending
		// to newScore slice as a bug seems to pop up occassionally if
		// strconv.Itoa is used more than once being passed to append.
		numPlayers = "1"
	} else {
		scores = g.scores2
		numPlayers = "2"
	}

	// If there are previous high scores then compare them to
	// player's current score
	if scores != nil {

		// Run for both players if more than one exists
		for _, p := range g.players {

			// Run through all scores in the list
			for i, s := range scores {

				// Score is saved as a string so it must be converted to
				// integer to compare
				scoreStr, err := strconv.Atoi(s[2])
				if err != nil {
					log.Println(err)
				}

				// Check if player's score is higher than current score from list
				if p.score > scoreStr {
					var newScore []string
					scoreChange = true

					// Create a formatted score of "number of players:player name:score"
					newScore = append(newScore, numPlayers, p.name, strconv.Itoa(p.score))

					// Append the previous high scores to the new high score list up until
					// where the newest high score should be inserted
					for a := 0; a < i; a++ {
						newScores = append(newScores, scores[a])
					}

					// Append the newest high score into the new high score list
					newScores = append(newScores, newScore)

					// Continue appending the rest of the previous high scores after the
					// newest high score until there are no scores left
					if i <= len(g.scores)-1 {
						for a := i; a < len(g.scores); a++ {
							newScores = append(newScores, scores[a])
						}
					}
					break

					// If the player's score is less than any of the previous high scores
					// but the number of previous high scores is less than the maximum
					// number of high scores saved, then add the score to the end of the list.
				} else if len(scores) < MaxHighScores {
					var newScore []string
					scoreChange = true
					newScore = append(newScore, numPlayers, p.name, strconv.Itoa(p.score))
					newScores = append(scores, newScore)
					break
				}
			}
		}

		// Check for changes in high score list
		if scoreChange {

			// Reset scores list
			scores = nil

			// If the number of high scores saved is longer than the maximum, then only
			// add scores up to the maximum back to the scores list
			if len(newScores) > MaxHighScores {
				for i := 0; i < MaxHighScores; i++ {
					scores = append(scores, newScores[i])
				}

				// If its not higher then add all of them
			} else {
				for i := 0; i < len(newScores); i++ {
					scores = append(scores, newScores[i])
				}
			}
		}

		// If no previous high scores present then add all player scores
		// to high score list
	} else {
		for _, p := range g.players {
			var newScore []string
			scoreChange = true
			newScore = append(newScore, strconv.Itoa(g.numPlayers), p.name, strconv.Itoa(p.score))
			scores = append(scores, newScore)
		}
	}
	return scores, scoreChange
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

// Remove a bit from game
func (g *Game) removeBit(i int) {
	b := &Bit{}
	g.bits[i] = g.bits[len(g.bits)-1]
	g.bits[len(g.bits)-1] = b
	g.bits = g.bits[:len(g.bits)-1]
}

func (g *Game) getPlayerName(playerNum, w, h int, headStyle, typeStyle tcell.Style) string {
	var (
		char       rune
		chars      []rune
		newChars   []rune
		charString string
	)

	for {
		// Reset screen and values
		g.gview.Clear()
		newChars = nil

		// Render to screen
		renderCenterStr(g.gview, w, h, headStyle, "Name of Player "+strconv.Itoa(playerNum)+":")
		renderCenterStr(g.gview, w, h+2, typeStyle, charString+"|")
		g.screen.Show()

		// Get input
		char = handleStringInput(g)

		// Evaluate input
		if char == '\r' {
			return charString
		} else if char == '\n' {
			continue
		} else if char == '\v' {
			return "-quit-"
		} else if char == '\t' {
			for i := 0; i < len(chars)-1; i++ {
				newChars = append(newChars, chars[i])
			}
			chars = newChars
			charString = string(chars)
		} else {
			chars = append(chars, char)
			charString = string(chars)
		}
	}
}
