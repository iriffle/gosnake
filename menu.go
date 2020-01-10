package main

import (
	"github.com/gdamore/tcell"
)

type MenuItem struct {
	x, y, w  int
	str      string
	style    tcell.Style
	selected bool
}

type Menu struct {
	items    []*MenuItem
	defStyle tcell.Style
	selStyle tcell.Style
}

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

func NewMenu(items []*MenuItem, defStyle, selStyle tcell.Style) Menu {
	m := Menu{
		items,
		defStyle,
		selStyle,
	}
	return m
}

func NewPlayerMenu(menuOptions [3]string, defStyle, selStyle tcell.Style) Menu {
	var items []*MenuItem
	for _, option := range menuOptions {
		p := NewMenuItem(MapWidth, MapHeight, option, DefStyle)
		items = append(items, &p)
	}
	m := NewMenu(items, defStyle, selStyle)
	m.AdjustItemPos()
	return m
}

func (m *Menu) SetSelected(i int) {
	m.items[i].selected = true
}

func (m *Menu) GetSelected() int {
	for i, item := range m.items {
		if item.selected {
			return i
		}
	}
	return 0
}

func (m *Menu) ChangeSelected() {
	for _, item := range m.items {
		if item.selected {
			item.style = m.selStyle
		} else {
			item.style = m.defStyle
		}
	}
}

func (m *Menu) AdjustItemPos() {
	for i, item := range m.items {
		item.y = item.y + i
	}
}
