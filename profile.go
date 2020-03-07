package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gdamore/tcell"
	"github.com/google/logger"
)

type Profile struct {
	Name  string
	Color string
	Char  rune
}

// NewProfile creates a new player profile with a given name
// and color.
func NewProfile(name string, color string, char rune) *Profile {
	p := Profile{
		Name:  name,
		Color: color,
		Char:  char,
	}
	return &p
}

// GetStyle returns a tcell Style based on the profile's color.
func (p *Profile) GetStyle() tcell.Style {
	color := tcell.GetColor(p.Color)
	style := GetStyle(DefBGStyle, color)
	return style
}

// AssignToPlayer assigns the current profile color to a player.
func (p *Profile) AssignToPlayer(player *Player) {
	player.name = p.Name
	color := tcell.GetColor(p.Color)
	style := GetStyle(DefBGStyle, color)
	player.pos[0].style = style
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
			charList := getCharList(PlayerRunes)
			char, color := EditProfile(g, charList, PlayerColors)
			if char == -1 {
				return MenuMain
			}
			p := NewProfile(name, PlayerColors[color], PlayerRunes[char])
			g.profiles = append(g.profiles, p)
			WriteProfiles(g.profiles, g.proFile)
			return MenuProfile
		}
		return MenuMain
	}
	return MenuMain
}

func EditProfile(g *Game, chars, colors []string) (int, int) {
	var entities []*Entity
	var objects []*Object
	var curColors []string
	char, color := 0, 0

	// Create char and color select menus
	charMenu := NewMainMenu(chars, g.style.DefStyle, g.style.SelStyle, 0)
	colorMenu := NewMainMenu(colors, g.style.DefStyle, g.style.DefStyle, 0)

	// Set general positioning for screen elements
	l := len(charMenu.items)
	w := (MapWidth / 2) - l
	h := (MapHeight / 2) - 8

	// Create Color Entity to display color selection
	eColor := NewEntity(w-2, h+2, DirAll, 0, PlayerRune, g.style.DefStyle)
	for i := 0; i < len(PlayerColors); i++ {
		style := StringToStyle(PlayerColors[i], PlayerColors[1])
		eColor.pos[i].oy++
		eColor.pos[i].style = style
		if i < len(PlayerColors)-1 {
			eColor.AddSegment(1, eColor.pos[0].char, style)
		}
	}

	// Create Display Entity to show current selected attributes on
	eDisplayH := NewEntity(w+8, h+8, DirAll, 0, PlayerRune, g.style.DefStyle)
	for i := 0; i < 20; i++ {
		eDisplayH.pos[i].ox++
		eDisplayH.AddSegment(1, eDisplayH.pos[0].char, eDisplayH.pos[0].style)
	}
	style := StringToStyle(PlayerColors[0], PlayerColors[1])
	eDisplayH.SetStyle(style)

	// Create Display Entity to show current selected attributes on
	eDisplayV := NewEntity(MapWidth/2, MapHeight/2-6, DirAll, 0, PlayerRune, g.style.DefStyle)
	for i := 0; i < 11; i++ {
		eDisplayV.pos[i].oy++
		eDisplayV.AddSegment(1, eDisplayV.pos[0].char, eDisplayV.pos[0].style)
	}
	eDisplayV.SetStyle(style)

	// Create Display Entity to show current selected attributes on
	eDisplayDL := NewEntity(MapWidth/2-6, MapHeight/2-6, DirAll, 0, PlayerRune, g.style.DefStyle)
	for i := 0; i < 11; i++ {
		eDisplayDL.pos[i].oy++
		eDisplayDL.pos[i].ox++
		eDisplayDL.AddSegment(1, eDisplayDL.pos[0].char, eDisplayDL.pos[0].style)
	}
	eDisplayDL.SetStyle(style)

	// Create Display Entity to show current selected attributes on
	eDisplayDR := NewEntity(MapWidth/2-6, MapHeight/2+5, DirAll, 0, PlayerRune, g.style.DefStyle)
	for i := 0; i < 11; i++ {
		eDisplayDR.pos[i].oy--
		eDisplayDR.pos[i].ox++
		eDisplayDR.AddSegment(1, eDisplayDR.pos[0].char, eDisplayDR.pos[0].style)
	}
	eDisplayDR.SetStyle(style)

	// Create dots to show what attributes are currently selected
	oChar := NewObject(w, h-1, '■', g.style.SelStyle, false)
	oColor := NewObject(w-4, h+2, '■', g.style.SelStyle, false)

	entities = append(entities, eDisplayH, eDisplayDL, eDisplayV, eDisplayDR)
	objects = append(objects, oColor, oChar)
	curColors = append(curColors, PlayerColors[0], PlayerColors[1])

	// Display editor and handle input
	rotation := 0
	cMode := true
	for char == 0 {
		g.gview.Clear()
		renderObjects(g.gview, objects)
		renderEntity(g.gview, eColor)
		renderEntity(g.gview, entities[rotation])
		renderProfile(g, charMenu, w, h, g.style.DefStyle)
		char, curColors = handleProfileInput(g, entities[rotation], oColor, oChar, charMenu, colorMenu, curColors, rotation, cMode)
		switch char {
		case 2:
			eDisplayDL.SetChar(entities[rotation].pos[0].char)
			eDisplayDL.SetStyle(entities[rotation].pos[0].style)
			char = 0
			rotation = 1
		case 3:
			eDisplayV.SetChar(entities[rotation].pos[0].char)
			eDisplayV.SetStyle(entities[rotation].pos[0].style)
			char = 0
			rotation = 2
		case 4:
			eDisplayDR.SetChar(entities[rotation].pos[0].char)
			eDisplayDR.SetStyle(entities[rotation].pos[0].style)
			char = 0
			rotation = 3
		case 5:
			eDisplayH.SetChar(entities[rotation].pos[0].char)
			eDisplayH.SetStyle(entities[rotation].pos[0].style)
			char = 0
			rotation = 0
		case 6:
			cMode = false
			char = 0
		case 7:
			cMode = true
			char = 0
		}
	}
	// Get selected attributes after enter is pressed
	if char == 1 {
		char = charMenu.GetSelected()
		color = colorMenu.GetSelected()
	}
	return char, color
}

func RotateDisplay(entities []*Entity) {

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
