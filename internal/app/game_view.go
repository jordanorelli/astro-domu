package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/sim"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type gameView struct {
	*blammo.Log
	room *wire.Room
	me   *wire.Entity
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

func (v *gameView) move(ui *UI, dx, dy int) {
	// fuck lol
	go ui.client.Send(sim.Move{dx, dy})
}

func (v *gameView) draw(b *buffer) {
	v.drawHeader(b)

	// fill in background dots first
	for x := 0; x < v.room.Width; x++ {
		for y := 0; y < v.room.Height; y++ {
			b.set(x+1, y+2, tile{r: '·', style: tcell.StyleDefault})
		}
	}
	b.set(0, 1, tile{r: '┌'})
	b.set(v.room.Width+1, 1, tile{r: '┐'})
	b.set(0, v.room.Height+2, tile{r: '└'})
	b.set(v.room.Width+1, v.room.Height+2, tile{r: '┘'})
	for x := 0; x < v.room.Width; x++ {
		b.set(x+1, 1, tile{r: '─'})
		b.set(x+1, v.room.Height+2, tile{r: '─'})
	}
	for y := 0; y < v.room.Height; y++ {
		b.set(0, y+2, tile{r: '│'})
		b.set(v.room.Width+1, y+2, tile{r: '│'})
	}
	for _, e := range v.room.Entities {
		pos := e.Position
		b.set(pos.X+1, pos.Y+2, tile{r: e.Glyph, style: tcell.StyleDefault})
	}

}

func (v *gameView) drawHeader(b *buffer) {
	// the first row is the name of the room that we're currently in
	var style tcell.Style
	style = style.Background(tcell.NewRGBColor(64, 64, 128))
	style = style.Foreground(tcell.NewRGBColor(0, 0, 0))

	runes := []rune(v.room.Name)

	for i := 0; i < b.width; i++ {
		if i < len(runes) {
			b.set(i, 0, tile{r: runes[i], style: style})
		} else {
			b.set(i, 0, tile{r: ' ', style: style})
		}
	}
}
