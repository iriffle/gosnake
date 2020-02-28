package main

import "github.com/gdamore/tcell"

type Profile struct {
	name   string
	player *Player
	style  tcell.Style
}

func NewProfile(name string, style tcell.Style) *Profile {
	p := Profile{
		name:  name,
		style: style,
	}

	return &p
}

func (p *Profile) SetName(name string) {
	p.name = name
}

func (p *Profile) SetPlayer(player *Player) {
	p.player = player
	p.player.name = p.name
}

func (p *Profile) SetStyle(style tcell.Style) {
	p.style = style
}
