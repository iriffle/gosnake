package main

import (
	"github.com/gdamore/tcell"
)

const (

	// Preset colors
	Black   = tcell.ColorBlack
	Maroon  = tcell.ColorMaroon
	Green   = tcell.ColorGreen
	Navy    = tcell.ColorNavy
	Olive   = tcell.ColorOlive
	Purple  = tcell.ColorPurple
	Teal    = tcell.ColorTeal
	Silver  = tcell.ColorSilver
	Gray    = tcell.ColorGray
	Red     = tcell.ColorRed
	Blue    = tcell.ColorBlue
	Lime    = tcell.ColorLime
	Yellow  = tcell.ColorYellow
	Fuchsia = tcell.ColorFuchsia
	Aqua    = tcell.ColorAqua
	White   = tcell.ColorWhite

	DefBGStyle = Black
	DefFGStyle = Silver
	SelFGStyle = Aqua
)

type Style struct {
	DefStyle          tcell.Style
	SelStyle          tcell.Style
	SelStyleBG        tcell.Style
	BitStyle          tcell.Style
	BiteStyle         tcell.Style
	BiteExplodedStyle tcell.Style
	DefBGColor        tcell.Color
	DefFGColor        tcell.Color
	DefSelColor       tcell.Color
}

func (s *Style) SetDefaultStyle() {
	s.DefStyle = GetStyle(DefBGStyle, DefFGStyle)
	s.SelStyle = GetStyle(DefBGStyle, SelFGStyle)
	s.SelStyleBG = GetStyle(Aqua, DefFGStyle)
	s.BitStyle = GetStyle(Black, White)
	s.BiteStyle = GetStyle(Black, Fuchsia)
	s.BiteExplodedStyle = GetStyle(Black, Red)
	s.DefBGColor = Black
	s.DefFGColor = Silver
	s.DefSelColor = Aqua
}

func (s *Style) SetNewStyle(defStyle, selStyle, selStyleBG, bitStyle, biteStyle, biteExplodedStyle tcell.Style, defBGColor, defFGColor, defSelColor tcell.Color, playerColors []tcell.Style) {
	s.DefStyle = defStyle
	s.SelStyle = selStyle
	s.SelStyleBG = selStyleBG
	s.BitStyle = bitStyle
	s.BiteStyle = biteStyle
	s.BiteExplodedStyle = biteExplodedStyle
}

// Generate a tcell style using a provided background and foreground color
func GetStyle(bg tcell.Color, fg tcell.Color) tcell.Style {
	style := tcell.StyleDefault.
		Background(bg).
		Foreground(fg)

	return style
}

func StringToStyle(fg, bg string) tcell.Style {
	fgColor := tcell.GetColor(fg)
	bgColor := tcell.GetColor(bg)
	style := GetStyle(bgColor, fgColor)
	return style
}
