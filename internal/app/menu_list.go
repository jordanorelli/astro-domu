package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
)

type menuList struct {
	choices   []menuItem
	highlight int
}

func (m *menuList) handleEvent(e tcell.Event) change {
	switch t := e.(type) {
	case *tcell.EventKey:
		key := t.Key()
		switch key {
		case tcell.KeyEnter:
			return m.choices[m.highlight].onSelect
		case tcell.KeyDown:
			m.highlight = (m.highlight + 1) % len(m.choices)
			return nil
		case tcell.KeyUp:
			if m.highlight == 0 {
				m.highlight = len(m.choices) - 1
			} else {
				m.highlight--
			}
			return nil
		}
	}
	return nil
}

func (m *menuList) draw(img canvas, _ *state) {
	for i, choice := range m.choices {
		if i == m.highlight {
			writeString(img, "▷ "+choice.name, math.Vec{0, i}, tcell.StyleDefault)
		} else {
			writeString(img, "  "+choice.name, math.Vec{0, i}, tcell.StyleDefault)
		}
	}
}

func (m *menuList) setFocus(bool) {}

type menuItem struct {
	name     string
	onSelect change
}
