package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/logger"
	"github.com/stjiub/gosnake/game"
)

const (
	logFile   = "log.txt"
	proFile   = "profiles.json"
	scoreFile = "hs.json"
)

var (
	verbose = flag.Bool("verbose", false, "print info level logs to stdout")

	// Keep track of previous game values
	lastGameState  int = game.Play
	lastNumPlayers int
	curProfiles    []*game.Profile
)

func main() {
	flag.Parse()

	// Set rand seed
	rand.Seed(time.Now().UnixNano())

	// Set logging
	lf, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		//if f, e := os.Create(logfile); e == nil {
		logger.Fatalf("Error opening log file: %v", err)
		//}
	}
	defer lf.Close()

	defer logger.Init("Error log", *verbose, true, lf).Close()
	logger.SetFlags(log.LstdFlags)

	// Game loop
	for {
		// Create game
		g := game.NewGame(lastNumPlayers, curProfiles, scoreFile, proFile)

		// Initialize screen
		err := g.InitScreen()
		if err != nil {
			logger.Fatalf("Error initializing screen: %v", err)
		}

		// Open main menu
		if lastGameState == game.Play || lastGameState == game.MainMenu {
			err := g.MainMenu()
			if err != nil {
				logger.Fatalf("Error running MainMenu: %v", err)
			}
		}

		if g.GetState() == game.Play {
			// Setup a game
			err := g.InitMap()
			if err != nil {
				logger.Fatalf("Error initializing map: %v", err)
			}
			err = g.InitPlayers()
			if err != nil {
				logger.Fatalf("Error initializing players: %v", err)
			}

			// Run the game
			err = g.Run()
			if err != nil {
				logger.Errorf("Error during main game loop: %v", err)
			}
		}

		// Quit game if signaled
		if g.GetState() == game.Quit {
			g.Quit()
		}

		// Save game values
		lastGameState = g.GetState()
		lastNumPlayers = g.GetNumPlayers()
		curProfiles = g.GetCurProfiles()
	}
}
