package main

import (
	"github.com/gdamore/tcell/v2"
)

type uiMode interface {
	handleEvent(*ui, tcell.Event) bool
	draw(*ui)
}

type boxWalker struct {
	width    int
	height   int
	position point
}

func (m *boxWalker) handleEvent(ui *ui, e tcell.Event) bool {
	switch v := e.(type) {
	case *tcell.EventKey:
		key := v.Key()
		if key == tcell.KeyRune {
			switch v.Rune() {
			case 'w':
				ui.client.send("self/move", map[string]interface{}{"delta": []int{0, -1}})
				m.move(0, -1)
			case 'a':
				ui.client.send("self/move", map[string]interface{}{"delta": []int{-1, 0}})
				m.move(-1, 0)
			case 's':
				ui.client.send("self/move", map[string]interface{}{"delta": []int{0, 1}})
				m.move(0, 1)
			case 'd':
				ui.client.send("self/move", map[string]interface{}{"delta": []int{1, 0}})
				m.move(1, 0)
			}
		}
	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return true
}

func (m *boxWalker) move(dx, dy int) {
	m.position.x = clamp(m.position.x+dx, 0, m.width-1)
	m.position.y = clamp(m.position.y+dy, 0, m.height-1)
}

func (m *boxWalker) draw(ui *ui) {
	offset := point{1, 1}

	// fill in background dots first
	for x := 0; x < m.width; x++ {
		for y := 0; y < m.height; y++ {
			ui.screen.SetContent(x+offset.x, y+offset.y, '·', nil, tcell.StyleDefault)
		}
	}

	// frame it
	ui.screen.SetContent(offset.x-1, offset.y-1, '┌', nil, tcell.StyleDefault)
	ui.screen.SetContent(offset.x+m.width, offset.y-1, '┐', nil, tcell.StyleDefault)
	ui.screen.SetContent(offset.x-1, offset.y+m.height, '└', nil, tcell.StyleDefault)
	ui.screen.SetContent(offset.x+m.width, offset.y+m.height, '┘', nil, tcell.StyleDefault)
	for x := 0; x < m.width; x++ {
		ui.screen.SetContent(x+offset.x, offset.y-1, '─', nil, tcell.StyleDefault)
		ui.screen.SetContent(x+offset.x, offset.y+m.height, '─', nil, tcell.StyleDefault)
	}
	for y := 0; y < m.height; y++ {
		ui.screen.SetContent(offset.x-1, y+offset.y, '│', nil, tcell.StyleDefault)
		ui.screen.SetContent(offset.x+m.width, y+offset.y, '│', nil, tcell.StyleDefault)
	}

	// add all characters
	ui.screen.SetContent(m.position.x+offset.x, m.position.y+offset.y, '@', nil, tcell.StyleDefault)
}
