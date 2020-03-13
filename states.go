package main

// Game states
const (
	Play = iota
	Quit
	Pause
	Restart
	MainMenu
)

// Menu pages
const (
	MenuMain = iota
	MenuPlayer
	MenuMode
	MenuScore
	MenuProfile
	MenuEdit
	MenuRemove
	MenuSettings
)

// Direction
const (
	DirUp = iota
	DirDown
	DirLeft
	DirRight
	DirAll
)

// Game modes
const (
	Player1 = iota
	Player2
	Battle
)

// Levels
const (
	Level2 = 20
	Level3 = 40
	Level4 = 60
	Level5 = 80
	Level6 = 100
)

// Menu item input values
const (
	ItemExit  = -1
	ItemNone  = 0
	ItemEnter = 1
	BGMode    = 2
	FGMode    = 3
)

// Snake editor rotations
const (
	Horizontal = iota
	DiagLeft
	Vertical
	DiagRight
)

// Player item states
const (
	WallPass = iota
)
