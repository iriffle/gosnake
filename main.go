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
	// Saved high score file
	scoreFile = "hs.dat"

	// Used if log flag provided
	logFile string

	// Keep track of previous game values
	lastGameState  int = Play
	lastNumPlayers int

	initScreen = false
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

	// Check if high score file exists. If not then create it
	_, err := os.Stat(scoreFile)
	if os.IsNotExist(err) {
		f, err := os.Create(scoreFile)
		if err != nil {
			log.Println(err, f)
		}
	}

	// Game loop
	for {
		// Create game
		g := &Game{numPlayers: lastNumPlayers, scoreFile: scoreFile}

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
