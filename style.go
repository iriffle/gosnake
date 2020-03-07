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
}

func SetDefaultStyle() *Style {

	s := Style{
		DefStyle:          GetStyle(DefBGStyle, DefFGStyle),
		SelStyle:          GetStyle(DefBGStyle, SelFGStyle),
		SelStyleBG:        GetStyle(Aqua, DefFGStyle),
		BitStyle:          GetStyle(Black, White),
		BiteStyle:         GetStyle(Black, Fuchsia),
		BiteExplodedStyle: GetStyle(Black, Red),
	}

	return &s
}

func SetNewStyle(defStyle, selStyle, selStyleBG, bitStyle, biteStyle, biteExplodedStyle tcell.Style, playerColors []tcell.Style) *Style {
	s := Style{
		DefStyle:          defStyle,
		SelStyle:          selStyle,
		SelStyleBG:        selStyleBG,
		BitStyle:          bitStyle,
		BiteStyle:         biteStyle,
		BiteExplodedStyle: biteExplodedStyle,
	}

	return &s
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
