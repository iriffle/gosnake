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
	Basic = iota
	Advanced
	Battle
)
