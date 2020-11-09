package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/sim"
	"github.com/jordanorelli/blammo"
)

type gameView struct {
	*blammo.Log
	inFocus bool
}

func (v *gameView) handleEvent(e tcell.Event) change {
	switch t := e.(type) {
	case *tcell.EventKey:
		key := t.Key()
		if key == tcell.KeyRune {
			switch t.Rune() {
			case 'w':
				return move{0, -1}
			case 'a':
				return move{-1, 0}
			case 's':
				return move{0, 1}
			case 'd':
				return move{1, 0}
			}
		}
	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return nil
}

func (v *gameView) draw(b *buffer, st *state) {
	v.drawHeader(b, st)

	// fill in background dots first
	for x := 0; x < st.room.Width; x++ {
		for y := 0; y < st.room.Height; y++ {
			b.set(x+1, y+2, tile{r: '·', style: tcell.StyleDefault})
		}
	}
	b.set(0, 1, tile{r: '┌'})
	b.set(st.room.Width+1, 1, tile{r: '┐'})
	b.set(0, st.room.Height+2, tile{r: '└'})
	b.set(st.room.Width+1, st.room.Height+2, tile{r: '┘'})
	for x := 0; x < st.room.Width; x++ {
		b.set(x+1, 1, tile{r: '─'})
		b.set(x+1, st.room.Height+2, tile{r: '─'})
	}
	for y := 0; y < st.room.Height; y++ {
		b.set(0, y+2, tile{r: '│'})
		b.set(st.room.Width+1, y+2, tile{r: '│'})
	}
	for _, e := range st.room.Entities {
		pos := e.Position
		b.set(pos.X+1, pos.Y+2, tile{r: e.Glyph, style: tcell.StyleDefault})
	}

}

func (v *gameView) drawHeader(b *buffer, st *state) {
	// the first row is the name of the room that we're currently in
	var style tcell.Style
	style = style.Background(tcell.NewRGBColor(64, 64, 128))
	style = style.Foreground(tcell.NewRGBColor(0, 0, 0))

	runes := []rune(st.room.Name)

	for i := 0; i < b.width; i++ {
		if i < len(runes) {
			b.set(i, 0, tile{r: runes[i], style: style})
		} else {
			b.set(i, 0, tile{r: ' ', style: style})
		}
	}
}

func (v *gameView) setFocus(yes bool) { v.inFocus = yes }

type move struct {
	x int
	y int
}

func (m move) exec(ui *UI) {
	ui.client.Send(sim.Move{m.x, m.y})
}
