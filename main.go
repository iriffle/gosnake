package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	lastGameState  int = Play
	lastNumPlayers int
)

func main() {

	var logfile string

	flag.StringVar(&logfile, "log", logfile, "Log file for debugging log")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	if logfile != "" {
		if f, e := os.Create(logfile); e == nil {
			log.SetOutput(f)
		}
	} else {
		log.SetOutput(ioutil.Discard)
	}

	for {
		game := &Game{debug: true,
			numPlayers: lastNumPlayers}
		if err := game.InitScreen(); err != nil {
			fmt.Printf("Failed to initialize game: %v\n", err)
			os.Exit(1)
		}
		if lastGameState == Play || lastGameState == MainMenu {
			game.MainMenu()
		}
		game.InitGame()
		game.Run()
		if game.state == Quit {
			game.Quit()
		}
		lastGameState = game.state
		lastNumPlayers = game.numPlayers
	}
}
