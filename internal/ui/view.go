package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/sim"
	"github.com/jordanorelli/astro-domu/internal/wire"
)

type view interface {
	handleEvent(*UI, tcell.Event) bool
	notify(wire.Value)
	draw(*UI)
}

type roomDisplay struct {
	width    int
	height   int
	entities map[int]sim.Entity
}

func (m *roomDisplay) handleEvent(ui *UI, e tcell.Event) bool {
	switch v := e.(type) {
	case *tcell.EventKey:
		key := v.Key()
		if key == tcell.KeyRune {
			switch v.Rune() {
			case 'w':
				m.move(ui, 0, -1)
			case 'a':
				m.move(ui, -1, 0)
			case 's':
				m.move(ui, 0, 1)
			case 'd':
				m.move(ui, 1, 0)
			}
		}
	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return true
}

func (m *roomDisplay) notify(v wire.Value) {
	if e, ok := v.(*sim.Entity); ok {
		m.entities[e.ID] = *e
	}
}

func (m *roomDisplay) move(ui *UI, dx, dy int) {
	reply, err := ui.client.Send(sim.Move{dx, dy})
	if err != nil {
		return
	}
	e := reply.Body.(*sim.Entity)
	m.entities[e.ID] = *e
}

func (m *roomDisplay) draw(ui *UI) {
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

	for _, e := range m.entities {
		ui.screen.SetContent(e.Position[0]+offset.x, e.Position[1]+offset.y, e.Glyph, nil, tcell.StyleDefault)
	}
}
