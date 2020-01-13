package main

import (
	"github.com/gdamore/tcell"
)

// Selectable menu item
type MenuItem struct {
	x, y, w  int
	str      string
	style    tcell.Style
	selected bool
}

// Menu of MenuItems
type Menu struct {
	items    []*MenuItem
	defStyle tcell.Style
	selStyle tcell.Style
}

// Create new MenuItem
func NewMenuItem(w, h int, str string, style tcell.Style) MenuItem {
	l := len(str)
	x := (w / 2) - (l / 2)
	y := h / 2
	i := MenuItem{
		x,
		y,
		l,
		str,
		style,
		false,
	}
	return i
}

// Create NewMenu with slice of MenuItems
func NewMenu(items []*MenuItem, defStyle, selStyle tcell.Style) Menu {
	m := Menu{
		items,
		defStyle,
		selStyle,
	}
	return m
}

// Create player menu screen
func NewPlayerMenu(menuOptions [3]string, defStyle, selStyle tcell.Style) Menu {
	var items []*MenuItem
	for _, option := range menuOptions {
		p := NewMenuItem(GameWidth, GameHeight, option, defStyle)
		items = append(items, &p)
	}
	m := NewMenu(items, defStyle, selStyle)
	m.AdjustItemPos()
	return m
}

// Sets current highlighted MenuItem
func (m *Menu) SetSelected(i int) {
	m.items[i].selected = true
}

// Gets current highlighted MenuItem
func (m *Menu) GetSelected() int {
	for i, item := range m.items {
		if item.selected {
			return i
		}
	}
	return 0
}

// Changes current highlighted MenuItem
func (m *Menu) ChangeSelected() {
	for _, item := range m.items {
		if item.selected {
			item.style = m.selStyle
		} else {
			item.style = m.defStyle
		}
	}
}

// Change MenuItem position based on number of items in Menu
func (m *Menu) AdjustItemPos() {
	for i, item := range m.items {
		item.y = item.y + i
	}
}
