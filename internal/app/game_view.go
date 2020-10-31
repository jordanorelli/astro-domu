package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/server/sim"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type gameView struct {
	*blammo.Log
	width  int
	height int
	// entities map[int]wire.Entity
}

func (v *gameView) handleEvent(ui *UI, e tcell.Event) bool {
	switch t := e.(type) {
	case *tcell.EventKey:
		key := t.Key()
		if key == tcell.KeyRune {
			switch t.Rune() {
			case 'w':
				v.move(ui, 0, -1)
			case 'a':
				v.move(ui, -1, 0)
			case 's':
				v.move(ui, 0, 1)
			case 'd':
				v.move(ui, 1, 0)
			}
		}
	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return true
}

func (v *gameView) notify(wv wire.Value) {
	v.Error("ignoring notifications at the moment: %v", wv)
	// if e, ok := v.(*wire.Entity); ok {
	// 	m.entities[e.ID] = *e
	// }
}

func (v *gameView) move(ui *UI, dx, dy int) {
	_, err := ui.client.Send(sim.Move{dx, dy})
	if err != nil {
		return
	}
	// e := reply.Body.(*wire.Entity)
	// v.entities[e.ID] = *e
}

func (v *gameView) draw(ui *UI) {
	offset := point{1, 1}

	// fill in background dots first
	for x := 0; x < v.width; x++ {
		for y := 0; y < v.height; y++ {
			ui.screen.SetContent(x+offset.x, y+offset.y, '·', nil, tcell.StyleDefault)
		}
	}

	// frame it
	ui.screen.SetContent(offset.x-1, offset.y-1, '┌', nil, tcell.StyleDefault)
	ui.screen.SetContent(offset.x+v.width, offset.y-1, '┐', nil, tcell.StyleDefault)
	ui.screen.SetContent(offset.x-1, offset.y+v.height, '└', nil, tcell.StyleDefault)
	ui.screen.SetContent(offset.x+v.width, offset.y+v.height, '┘', nil, tcell.StyleDefault)
	for x := 0; x < v.width; x++ {
		ui.screen.SetContent(x+offset.x, offset.y-1, '─', nil, tcell.StyleDefault)
		ui.screen.SetContent(x+offset.x, offset.y+v.height, '─', nil, tcell.StyleDefault)
	}
	for y := 0; y < v.height; y++ {
		ui.screen.SetContent(offset.x-1, y+offset.y, '│', nil, tcell.StyleDefault)
		ui.screen.SetContent(offset.x+v.width, y+offset.y, '│', nil, tcell.StyleDefault)
	}

	// for _, e := range v.entities {
	// 	ui.screen.SetContent(e.Position[0]+offset.x, e.Position[1]+offset.y, e.Glyph, nil, tcell.StyleDefault)
	// }
}
