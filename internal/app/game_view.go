package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/server/sim"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type gameView struct {
	*blammo.Log
	roomName    string
	width       int
	height      int
	me          wire.Entity
	allRooms    map[string]wire.Room
	allEntities map[int]wire.Entity
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
	switch z := wv.(type) {
	case *wire.UpdateEntity:
		if z.Room == v.roomName {
			v.Info("we want to read this one: %v", z)
		}
	}
}

func (v *gameView) move(ui *UI, dx, dy int) {
	reply, err := ui.client.Send(sim.Move{dx, dy})
	if err != nil {
		return
	}

	e := reply.Body.(*wire.UpdateEntity)
	// ughhhhhh
	v.me = wire.Entity{
		ID:       e.ID,
		Position: e.Position,
		Glyph:    e.Glyph,
	}
	v.allEntities[e.ID] = v.me
	// jfc this sucks
	v.allRooms[v.roomName].Entities[e.ID] = v.me
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

	for _, e := range v.allRooms[v.roomName].Entities {
		pos := e.Position.Add(math.Vec{1, 1})
		ui.screen.SetContent(pos.X, pos.Y, e.Glyph, nil, tcell.StyleDefault)
	}
}
