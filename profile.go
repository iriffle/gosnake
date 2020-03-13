package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gdamore/tcell"
	"github.com/google/logger"
)

type Profile struct {
	Name    string
	FGColor string
	BGColor string
	Char    rune
}

var (
	profileControls string = "r = rotate - c = background/foreground"
)

// NewProfile creates a new player profile with a given name
// and color.
func NewProfile(name string, fgColor, bgColor string, char rune) *Profile {
	p := Profile{
		Name:    name,
		FGColor: fgColor,
		BGColor: bgColor,
		Char:    char,
	}
	return &p
}

// GetStyle returns a tcell Style based on the profile's color.
func (p *Profile) GetStyle() tcell.Style {
	fgColor := tcell.GetColor(p.FGColor)
	bgColor := tcell.GetColor(p.BGColor)
	style := GetStyle(bgColor, fgColor)
	return style
}

// AssignToPlayer assigns the current profile color to a player.
func (p *Profile) AssignToPlayer(player *Player) {
	player.name = p.Name
	fgColor := tcell.GetColor(p.FGColor)
	bgColor := tcell.GetColor(p.BGColor)
	style := GetStyle(bgColor, fgColor)
	player.style = style
}

// DecodeProfiles takes a JSON byte slice and converts it into
// a slice of Profiles.
func DecodeProfiles(byteValue []byte) []*Profile {
	var profiles []*Profile
	json.Unmarshal(byteValue, &profiles)
	return profiles
}

// EncodeProfiles takes a slice of Profiles and converts it to
// a JSON byte slice.
func EncodeProfiles(profiles []*Profile) []byte {
	file, _ := json.MarshalIndent(profiles, "", " ")
	return file
}

// WriteProfiles writes a slice of Profiles to a JSON file.
func WriteProfiles(profiles []*Profile, file string) {
	f, err := os.OpenFile(file, os.O_CREATE, 0660)
	if err != nil {
		logger.Errorf("Error creating new file: %v", err)
	}

	if err = f.Close(); err != nil {
		logger.Errorf("Error closing file: %v", err)
	}

	data := EncodeProfiles(profiles)
	_ = ioutil.WriteFile(file, data, 0644)
}

func CreateProfile(g *Game) int {
	var name string = ""
	for name == "" {
		name = GetProfileName(g, MapWidth, MapHeight)
		if name != "-quit-" {
			g.gview.Clear()
			p := NewProfile(name, PlayerColors[0], PlayerColors[1], PlayerRune)
			g.profiles = append(g.profiles, p)
			charList := getCharList(PlayerRunes)
			char := p.Edit(g, charList, PlayerColors)
			if char == ItemExit {
				return MenuProfile
			}
			WriteProfiles(g.profiles, g.proFile)
			return MenuProfile
		}
		return MenuProfile
	}
	return MenuProfile
}

func EditProfile(g *Game, p *Profile) int {
	g.gview.Clear()
	charList := getCharList(PlayerRunes)
	char := p.Edit(g, charList, PlayerColors)
	if char == ItemExit {
		return MenuProfile
	}
	WriteProfiles(g.profiles, g.proFile)
	return MenuProfile
}

func RemoveProfile(g *Game, i int) int {
	g.profiles[i] = g.profiles[len(g.profiles)-1]
	g.profiles[len(g.profiles)-1] = nil
	g.profiles = g.profiles[:len(g.profiles)-1]
	WriteProfiles(g.profiles, g.proFile)
	return MenuProfile
}

func (p *Profile) Edit(g *Game, chars, colors []string) int {
	// Create char and color select menus
	charMenu := NewMainMenu(chars, g.DefStyle, g.SelStyle, 0)
	colorMenu := NewMainMenu(colors, g.DefStyle, g.DefStyle, 0)

	// Set general positioning for screen elements
	l := len(charMenu.items)
	w := (MapWidth / 2) - l
	h := (MapHeight / 2) - 8

	// Create color entity to display color selection bar
	eColor := NewColorEntity(w-2, h+2, BitRune, g.DefStyle)

	// Create display entities for each rotation to show current selected attributes on
	style := StringToStyle(p.FGColor, p.BGColor)
	eDisplayH := NewDisplayEntity(w+12, h+8, 20, 1, 0, p.Char, style)
	eDisplayV := NewDisplayEntity(w+l, h+2, 11, 0, 1, p.Char, style)
	eDisplayDL := NewDisplayEntity(w+l-6, h+2, 11, 1, 1, p.Char, style)
	eDisplayDR := NewDisplayEntity(w+l-5, h+13, 11, 1, -1, p.Char, style)

	// Create dots to show what attributes are currently selected
	oChar := NewObject(w, h-1, BitRune, g.SelStyle, false)
	oColor := NewObject(w-4, h+2, BitRune, g.SelStyle, false)

	entities := []*Entity{eDisplayH, eDisplayDL, eDisplayV, eDisplayDR}
	objects := []*Object{oColor, oChar}
	cColors := []string{PlayerColors[0], PlayerColors[1]}

	rotation := Horizontal
	fgMode := true
	char := ItemNone

	// Display editor and handle input
	for char == ItemNone {
		g.gview.Clear()
		renderCenterStr(g.gview, MapWidth, MapHeight/4, g.SelStyle, "Edit Profile")
		//renderCenterStr(g.sview, MapWidth, 0, g.DefStyle, profileControls)
		renderObjects(g.gview, objects)
		renderEntity(g.gview, eColor)
		renderEntity(g.gview, entities[rotation])
		renderProfile(g, charMenu, w, h, g.DefStyle)
		char, rotation, cColors = handleProfileInput(g, entities, oColor, oChar, charMenu, colorMenu, cColors, rotation, fgMode)
		switch char {
		case BGMode:
			eColor.SetChar(PlayerRune)
			fgMode = false
			char = 0
		case FGMode:
			eColor.SetChar(BitRune)
			fgMode = true
			char = 0
		}
	}
	// Get selected attributes after enter is pressed
	if char == ItemEnter {
		p.Char = PlayerRunes[charMenu.GetSelected()]
		p.FGColor = cColors[0]
		p.BGColor = cColors[1]
	}

	return char
}

// GetProfileName allows a player to input their name.
func GetProfileName(g *Game, w, h int) string {
	var (
		char     rune
		chars    []rune
		newChars []rune
		hStr     string = "Name of Profile:"
		charStr  string
	)

	for {
		newChars = nil
		// Render the player select screen
		renderNameSelect(g, w, h, hStr, charStr)

		// Get input
		char = handleStringInput(g)

		// Evaluate input
		if char == '\r' {
			isPlayer := CheckProfileName(g, charStr)
			if isPlayer {
				hStr = "That profile name is taken. Provide a new name:"
				charStr = ""
				newChars = nil
				chars = nil
			} else {
				return charStr
			}
		} else if char == '\n' {
			continue
		} else if char == '\v' {
			return "-quit-"
		} else if char == '\t' {
			for i := 0; i < len(chars)-1; i++ {
				newChars = append(newChars, chars[i])
			}
			chars = newChars
			charStr = string(chars)
		} else {
			chars = append(chars, char)
			charStr = string(chars)
		}
	}
}

func CheckProfileName(g *Game, name string) bool {
	for i := range g.profiles {
		if name == g.profiles[i].Name {
			return true
		}
	}
	return false
}
