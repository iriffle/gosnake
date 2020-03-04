package main

const (
	// Game states
	Play = iota
	Quit
	Pause
	Restart
	MainMenu
)

const (
	// Menu pages
	MenuMain = iota
	MenuPlayer
	MenuMode
	MenuScore
	MenuProfile
	MenuSettings
)

const (
	// Direction
	DirUp = iota
	DirDown
	DirLeft
	DirRight
	DirAll
)

const (
	// Game modes
	Player1 = iota
	Player2
	Battle
)

const (
	// Levels
	Level2 = 200
	Level3 = 400
	Level4 = 600
	Level5 = 800
)
