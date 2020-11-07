package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/blammo"
)

type menuList struct {
	*blammo.Log
	choices  []string
	selected int
}

func (m *menuList) handleEvent(ui *UI, e tcell.Event) bool {
	switch t := e.(type) {
	case *tcell.EventKey:
		key := t.Key()
		switch key {
		case tcell.KeyDown:
			m.selected = (m.selected + 1) % len(m.choices)
		case tcell.KeyUp:
			if m.selected == 0 {
				m.selected = len(m.choices) - 1
			} else {
				m.selected--
			}
		}
	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return true
}

func (m *menuList) draw(b *buffer) {
	for i, choice := range m.choices {
		if i == m.selected {
			b.writeString("▷ "+choice, math.Vec{0, i}, tcell.StyleDefault)
		} else {
			b.writeString("  "+choice, math.Vec{0, i}, tcell.StyleDefault)
		}
	}
}

func (m *menuList) setFocus(bool) {}
