package game

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/views"
	"github.com/mattn/go-runewidth"
	"github.com/stjiub/gosnake/entity"
	"github.com/stjiub/gosnake/gamemap"
	"github.com/stjiub/gosnake/style"
)

// Render all of the game in main game loop
func renderAll(g *Game, style tcell.Style, m *gamemap.GameMap) {

	// Clear screen for redraw
	g.gview.Clear()
	g.screen.ShowCursor(0, MapHeight)

	// Draw game map
	renderMap(g.gview, m)

	// Draw the Bite explosion map
	renderMap(g.gview, g.biteMap)

	if g.numPlayers == 1 {
		renderLevel(g.gview, g.level, m.Width, m.Height, g.SelStyle)
	}
	renderScore(g.gview, g.players, m.Width, m.Height, g.SelStyle)
	renderBits(g.gview, g.bits)
	renderBites(g.gview, g.bites)
	renderItems(g.gview, g.items)
	renderEntities(g.gview, g.entities)
	renderPlayers(g.gview, g.players)
	g.sbar.SetCenter(controls, g.DefStyle)
	g.sbar.Draw()
	g.screen.Show()
}

// Render the game map
func renderMap(v *views.ViewPort, m *gamemap.GameMap) {
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			renderRune(v, x, y, m.Objects[x][y].GetStyle(), m.Objects[x][y].GetChar())
		}
	}
}

// Render a menu
func renderMenu(g *Game, m *Menu, style tcell.Style) {
	for i := range m.items {
		renderStr(g.gview, m.items[i].x, m.items[i].y, m.items[i].style, m.items[i].str)
	}
	g.screen.Show()
}

func renderProfile(g *Game, m *Menu, w, h int, style tcell.Style) {
	renderChars(g.gview, m, w, h)
	g.screen.Show()
}

// Render the name selection screen
func renderNameSelect(g *Game, w, h int, hStr, charStr string) {
	g.gview.Clear()
	renderCenterStr(g.gview, w, h, g.DefStyle, hStr)
	renderCenterStr(g.gview, w, h+2, g.SelStyle, charStr+"|")
	g.screen.Show()
}

// Render the High Score screen
func renderHighScoreScreen(g *Game, style tcell.Style, max int) {
	g.gview.Clear()
	g.gview.Fill(' ', style)
	renderCenterStr(g.gview, MapWidth, 4, style, "High Scores")
	renderCenterStr(g.gview, MapWidth, 6, style, strings.Repeat("=", MapWidth-10))
	renderCenterStr(g.gview, MapWidth, 10, style, "1 Player:")
	renderHighScores(g, Player1, 14)
	renderCenterStr(g.gview, MapWidth, (16 + max*2), style, "2 Player:")
	renderHighScores(g, Player2, (20 + max*2))

	g.screen.Show()
}

// Render a list of scores for the high score screen
func renderHighScores(g *Game, mode, lastScorePos int) {
	var scores []*Score
	if mode == Player1 {
		scores = g.scores1
	} else if mode == Player2 {
		scores = g.scores2
	}
	for i := range scores {
		if i < len(scores) {
			renderCenterStr(g.gview, MapWidth, lastScorePos+i, g.SelStyle, (scores[i].Name + " - " + strconv.Itoa(scores[i].Score)))
		} else {
			renderCenterStr(g.gview, MapWidth, lastScorePos+i, g.DefStyle, "----")
		}
		lastScorePos++
	}
}

// Render the player scores in middle of screen
func renderScore(v *views.ViewPort, players []*entity.Player, w, h int, style tcell.Style) {
	scores := ""
	for i := range players {
		scores = players[i].GetName() + ": " + strconv.Itoa(players[i].GetScore()) + " "
		renderCenterStr(v, w, h+i, style, scores)
	}
}

// Render the current level in middle of screen
func renderLevel(v *views.ViewPort, l, w, h int, style tcell.Style) {
	level := "level: " + strconv.Itoa(l)
	renderCenterStr(v, w, h-2, style, level)
}

// Render a Player
func renderPlayer(v *views.ViewPort, p *entity.Player) {
	for i := 0; i < p.GetLength(); i++ {
		var comb []rune
		comb = nil
		c := p.GetChar(i)
		x, y := p.GetCurPos(i)
		sty := p.GetStyle(i)
		v.SetContent(x, y, c, comb, sty)
	}
}

// Render an Entity
func renderEntity(v *views.ViewPort, e *entity.Entity) {
	for i := 0; i < e.GetLength(); i++ {
		var comb []rune
		comb = nil
		c := e.GetChar(i)
		x, y := e.GetCurPos(i)
		sty := e.GetStyle(i)
		v.SetContent(x, y, c, comb, sty)
	}

}

// Render all Entities
func renderEntities(v *views.ViewPort, entities []*entity.Entity) {
	for i := range entities {
		renderEntity(v, entities[i])
	}
}

// Render all Players
func renderPlayers(v *views.ViewPort, players []*entity.Player) {
	for i := range players {
		renderPlayer(v, players[i])
	}
}

// Render all objects
func renderObjects(v *views.ViewPort, objects []*gamemap.Object) {
	for i := range objects {
		x, y := objects[i].GetCurPos()
		char := objects[i].GetChar()
		sty := objects[i].GetStyle()
		renderRune(v, x, y, sty, char)
	}
}

// Render all Bits
func renderBits(v *views.ViewPort, bits []*entity.Bit) {
	for i := range bits {
		x, y := bits[i].GetCurPos()
		char := bits[i].GetChar()
		sty := bits[i].GetStyle()
		renderRune(v, x, y, sty, char)
	}
}

// Render all Bites
func renderBites(v *views.ViewPort, bites []*entity.Bite) {
	for i := range bites {
		x, y := bites[i].GetCurPos()
		char := bites[i].GetChar()
		sty := bites[i].GetStyle()
		renderRune(v, x, y, sty, char)
	}
}

func renderItems(v *views.ViewPort, items []*entity.Item) {
	for i := range items {
		x, y := items[i].GetCurPos()
		char := items[i].GetChar()
		sty := items[i].GetStyle()
		renderRune(v, x, y, sty, char)
	}
}

// Render a string at given position
func renderStr(v *views.ViewPort, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		v.SetContent(x, y, c, comb, style)
		x += w
	}
}

// Render a string in the center of the screen
func renderCenterStr(v *views.ViewPort, w, h int, style tcell.Style, str string) {
	x := (w / 2) - (len(str) / 2)
	y := (h / 2)
	renderStr(v, x, y, style, str)
}

// Render a string at given position
func renderChars(v *views.ViewPort, m *Menu, x, y int) {
	for _, i := range m.items {
		for _, c := range i.str {
			var comb []rune
			w := runewidth.RuneWidth(c)
			if w == 0 {
				comb = []rune{c}
				c = ' '
				w = 1
			}
			v.SetContent(x, y, c, comb, i.style)
			x += w + 1
		}
	}
}

// Render a single rune to the screen
func renderRune(v *views.ViewPort, x, y int, style tcell.Style, char rune) {
	var comb []rune
	comb = nil
	v.SetContent(x, y, char, comb, style)
}

func renderGoLogo(g *Game, w, h int) {
	w = w - 10
	h = h - 12
	sty := style.GetStyle(g.DefBGColor, style.White)
	renderStr(g.gview, w+6, h, sty, "______")
	renderStr(g.gview, w+5, h+1, sty, "//   ) )")
	renderStr(g.gview, w+4, h+2, sty, "//         ___")
	renderStr(g.gview, w+3, h+3, sty, "//  ____  //   ) )")
	renderStr(g.gview, w+2, h+4, sty, "//    / / //   / /")
	renderStr(g.gview, w+1, h+5, sty, "((____/ / ((___/ /")
	g.screen.Show()
}

func renderSnakeLogo(g *Game, w, h int) {
	w = w - 33
	h = h - 17
	sty := style.GetStyle(g.DefBGColor, style.Green)
	redSty := style.GetStyle(g.DefBGColor, style.Red)
	renderStr(g.gview, w, h, sty, "           /^\\/^\\")
	renderStr(g.gview, w, h+1, sty, "         _|__|  O|")
	renderStr(g.gview, w+8, h+2, sty, "/~     \\_/ \\")
	renderStr(g.gview, w, h+2, redSty, "\\/")
	renderStr(g.gview, w+6, h+3, sty, "|__________/  \\")
	renderStr(g.gview, w, h+3, redSty, " \\____")
	renderStr(g.gview, w, h+4, sty, "       \\_______      \\")
	renderStr(g.gview, w, h+5, sty, "                `\\     \\                   \\")
	renderStr(g.gview, w, h+6, sty, "                  |     |                   \\")
	renderStr(g.gview, w, h+7, sty, "                 /      /                    \\")
	renderStr(g.gview, w, h+8, sty, "                /     /                       \\\\")
	renderStr(g.gview, w, h+9, sty, "              /      /                         \\ \\")
	renderStr(g.gview, w, h+10, sty, "             /     /                            \\  \\")
	renderStr(g.gview, w, h+11, sty, "           /     /             _----_           \\   \\")
	renderStr(g.gview, w, h+12, sty, "          /     /           _-~      ~-_         |   |")
	renderStr(g.gview, w, h+13, sty, "         (      (        _-~    _--_    ~-_     _/   |")
	renderStr(g.gview, w, h+14, sty, "          \\      ~-____-~    _-~    ~-_    ~-_-~    /")
	renderStr(g.gview, w, h+15, sty, "            ~-_           _-~          ~-_       _-~")
	renderStr(g.gview, w, h+16, sty, "               ~--______-~                ~-___-~")
	g.screen.Show()
}
