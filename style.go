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
	BitStyle          tcell.Style
	BiteStyle         tcell.Style
	BiteExplodedStyle tcell.Style
	PlayerColors      []tcell.Style
}

func SetDefaultStyle() *Style {

	s := Style{
		DefStyle:          getStyle(DefBGStyle, DefFGStyle),
		SelStyle:          getStyle(DefBGStyle, SelFGStyle),
		BitStyle:          getStyle(Black, White),
		BiteStyle:         getStyle(Black, Fuchsia),
		BiteExplodedStyle: getStyle(Black, Red),

		PlayerColors: []tcell.Style{getStyle(DefBGStyle, tcell.ColorGreen), getStyle(DefBGStyle, tcell.ColorRed), getStyle(DefBGStyle, tcell.ColorSilver), getStyle(DefBGStyle, tcell.ColorAqua)},
	}

	return &s
}

func SetNewStyle(defStyle, selStyle, bitStyle, biteStyle, biteExplodedStyle tcell.Style, playerColors []tcell.Style) *Style {
	s := Style{
		DefStyle:          defStyle,
		SelStyle:          selStyle,
		BitStyle:          bitStyle,
		BiteStyle:         biteStyle,
		BiteExplodedStyle: biteExplodedStyle,
		PlayerColors:      playerColors,
	}

	return &s
}

// Generate a tcell style using a provided background and foreground color
func getStyle(bg tcell.Color, fg tcell.Color) tcell.Style {
	style := tcell.StyleDefault.
		Background(bg).
		Foreground(fg)

	return style
}
