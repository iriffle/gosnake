package main

import (
	"flag"
	"github.com/google/logger"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	logFile   = "log.txt"
	proFile   = "profiles.json"
	scoreFile = "hs.json"
)

var (
	verbose = flag.Bool("verbose", false, "print info level logs to stdout")

	// Keep track of previous game values
	lastGameState  int = Play
	lastNumPlayers int
	curProfiles    []*Profile
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
		g := &Game{numPlayers: lastNumPlayers, curProfiles: curProfiles, scoreFile: scoreFile, proFile: proFile}

		// Initialize screen
		g.InitScreen()

		// Open main menu
		if lastGameState == Play || lastGameState == MainMenu {
			g.MainMenu()
		}
		// Setup a game
		g.InitGame()

		// Run the game
		if g.state == Play {
			g.RunGame()
		}

		// Quit game if signaled
		if g.state == Quit {
			g.QuitGame()
		}

		// Save game values
		lastGameState = g.state
		lastNumPlayers = g.numPlayers
		curProfiles = g.curProfiles
	}
}
