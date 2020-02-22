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

	lastGameState  int = Play
	lastNumPlayers int
)

func main() {

	var logfile string

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
		game := &Game{numPlayers: lastNumPlayers, scoreFile: scoreFile}
		// Initialize screen
		if err := game.InitScreen(); err != nil {
			log.Println("Failed to initialize game: %v\n", err)
			os.Exit(1)
		}
		// Open main menu
		if lastGameState == Play || lastGameState == MainMenu {
			game.MainMenu()
		}
		// Setup a game
		game.InitGame()
		// Run the game
		game.Run()
		// Quit game if signaled
		if game.state == Quit {
			game.Quit()
		}
		lastGameState = game.state
		lastNumPlayers = game.numPlayers
	}
}
