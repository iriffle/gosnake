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
	MenuEdit
	MenuRemove
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
	Level2 = 20
	Level3 = 40
	Level4 = 60
	Level5 = 80
	Level6 = 100
)

const (
	ItemExit  = -1
	ItemNone  = 0
	ItemEnter = 1
	BGMode    = 2
	FGMode    = 3
)

const (
	Horizontal = iota
	DiagLeft
	Vertical
	DiagRight
)
