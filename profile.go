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
}

// NewProfile creates a new player profile with a given name
// and color.
func NewProfile(name string, color string) *Profile {
	p := Profile{
		Name:  name,
		Color: color,
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
