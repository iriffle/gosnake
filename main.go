package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	// Used if log flag provided
	logFile string

	// Keep track of previous game values
	lastGameState  int = Play
	lastNumPlayers int

	files = []string{"1.dat", "2.dat", "profiles.dat"}
)

func main() {

	var logfile string

	// Check if log flag provided
	flag.StringVar(&logfile, "log", logfile, "Log file for debugging log")
	flag.Parse()

	// Set rand seed
	rand.Seed(time.Now().UnixNano())

	// Set logging
	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			//if f, e := os.Create(logfile); e == nil {
			log.Fatalf("Error opening log file")
			//}
		}
		defer f.Close()
	} else {
		log.SetOutput(ioutil.Discard)
	}

	// Game loop
	for {
		// Create game
		g := &Game{numPlayers: lastNumPlayers, files: files}

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
	}
}
