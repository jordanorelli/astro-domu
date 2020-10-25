package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type uiMode interface {
	HandleEvent(tcell.Event) bool
	draw(*ui)
}

type boxWalker struct {
	width    int
	height   int
	position point
}

func (m *boxWalker) HandleEvent(e tcell.Event) bool {
	switch v := e.(type) {
	case *tcell.EventKey:
		key := v.Key()
		if key == tcell.KeyRune {
			switch v.Rune() {
			case 'w':
				m.position.y--
			case 'a':
				m.position.x--
			case 's':
				m.position.y++
			case 'd':
				m.position.x++
			}
		}
	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return true
}

func (m *boxWalker) draw(ui *ui) {
	if m.position.x < 2 {
		m.position.x = 2
	}
	if m.position.x > 11 {
		m.position.x = 11
	}
	if m.position.y < 2 {
		m.position.y = 2
	}
	if m.position.y > 10 {
		m.position.y = 10
	}
	ui.writeString(1, 1, `┌──────────┐`, tcell.StyleDefault)
	ui.writeString(1, 2, `│··········│`, tcell.StyleDefault)
	ui.writeString(1, 3, `│··········│`, tcell.StyleDefault)
	ui.writeString(1, 4, `│··········│`, tcell.StyleDefault)
	ui.writeString(1, 5, `│··········│`, tcell.StyleDefault)
	ui.writeString(1, 6, `│··········│`, tcell.StyleDefault)
	ui.writeString(1, 7, `│··········│`, tcell.StyleDefault)
	ui.writeString(1, 8, `│··········│`, tcell.StyleDefault)
	ui.writeString(1, 9, `│··········│`, tcell.StyleDefault)
	ui.writeString(1, 10, `│··········│`, tcell.StyleDefault)
	ui.writeString(1, 11, `└──────────┘`, tcell.StyleDefault)
	ui.screen.SetContent(m.position.x, m.position.y, '@', nil, tcell.StyleDefault)
	ui.writeString(0, 12, fmt.Sprintf(" (%02d, %02d)", m.position.x, m.position.y), tcell.StyleDefault)
}
