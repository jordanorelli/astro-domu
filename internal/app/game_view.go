package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/sim"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type gameView struct {
	*blammo.Log
	inFocus    bool
	keyHandler func(*tcell.EventKey) change
	statusLine string
}

func (v *gameView) handleEvent(e tcell.Event) change {
	if v.keyHandler == nil {
		v.keyHandler = v.walkHandler
		v.statusLine = "(walk)"
	}

	switch t := e.(type) {
	case *tcell.EventKey:
		return v.keyHandler(t)
	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return nil
}

func (v *gameView) walkHandler(e *tcell.EventKey) change {
	if e.Key() == tcell.KeyRune {
		switch e.Rune() {
		case 'w':
			return move{0, -1}
		case 'a':
			return move{-1, 0}
		case 's':
			return move{0, 1}
		case 'd':
			return move{1, 0}
		case 'i':
			return openInventory{}
		case 'l':
			v.keyHandler = v.lookHandler
			v.statusLine = "(look)"
		case 'p':
			v.keyHandler = v.pickupHandler
			v.statusLine = "(pickup)"
		case 'P':
			v.keyHandler = v.putdownHandler
			v.statusLine = "(put down)"
		}
	}
	return nil
}

func (v *gameView) lookHandler(e *tcell.EventKey) change {
	if e.Key() == tcell.KeyESC {
		v.keyHandler = v.walkHandler
		v.statusLine = "(walk)"
		return nil
	}

	if e.Key() == tcell.KeyRune {
		switch e.Rune() {
		case 'w':
			v.keyHandler = v.walkHandler
			v.statusLine = "(walk)"
			return &lookAt{x: 0, y: -1}
		case 'a':
			v.keyHandler = v.walkHandler
			v.statusLine = "(walk)"
			return &lookAt{x: -1, y: 0}
		case 's':
			v.keyHandler = v.walkHandler
			v.statusLine = "(walk)"
			return &lookAt{x: 0, y: 1}
		case 'd':
			v.keyHandler = v.walkHandler
			v.statusLine = "(walk)"
			return &lookAt{x: 1, y: 0}
		}
	}
	return nil
}

func (v *gameView) pickupHandler(e *tcell.EventKey) change {
	if e.Key() == tcell.KeyESC {
		v.keyHandler = v.walkHandler
		v.statusLine = "(walk)"
		return nil
	}

	if e.Key() == tcell.KeyRune {
		switch e.Rune() {
		case 'w':
			v.keyHandler = v.walkHandler
			v.statusLine = "(walk)"
			return &pickup{x: 0, y: -1}
		case 'a':
			v.keyHandler = v.walkHandler
			v.statusLine = "(walk)"
			return &pickup{x: -1, y: 0}
		case 's':
			v.keyHandler = v.walkHandler
			v.statusLine = "(walk)"
			return &pickup{x: 0, y: 1}
		case 'd':
			v.keyHandler = v.walkHandler
			v.statusLine = "(walk)"
			return &pickup{x: 1, y: 0}
		}
	}
	return nil
}

func (v *gameView) putdownHandler(e *tcell.EventKey) change {
	v.keyHandler = v.walkHandler
	v.statusLine = "(walk)"
	return nil
}

func (v *gameView) draw(img canvas, st *state) {
	fill(img, tcell.StyleDefault.Background(tcell.NewRGBColor(0, 0, 12)))
	v.drawHeader(img, st)

	// fill in background dots first
	dim := tcell.StyleDefault.Foreground(tcell.NewRGBColor(76, 72, 83))
	for x := 0; x < st.room.Width; x++ {
		for y := 0; y < st.room.Height; y++ {
			img.setTile(x+1, y+2, tile{r: '·', style: dim})
		}
	}
	img.setTile(0, 1, tile{r: '┌'})
	img.setTile(st.room.Width+1, 1, tile{r: '┐'})
	img.setTile(0, st.room.Height+2, tile{r: '└'})
	img.setTile(st.room.Width+1, st.room.Height+2, tile{r: '┘'})
	for x := 0; x < st.room.Width; x++ {
		img.setTile(x+1, 1, tile{r: '─'})
		img.setTile(x+1, st.room.Height+2, tile{r: '─'})
	}
	for y := 0; y < st.room.Height; y++ {
		img.setTile(0, y+2, tile{r: '│'})
		img.setTile(st.room.Width+1, y+2, tile{r: '│'})
	}
	for _, e := range st.room.Entities {
		pos := e.Position
		img.setTile(pos.X+1, pos.Y+2, tile{r: e.Glyph, style: tcell.StyleDefault})
	}
	writeString(img, v.statusLine, math.Vec{0, img.bounds().Height - 1}, tcell.StyleDefault)
}

func (v *gameView) drawHeader(img canvas, st *state) {
	// the first row is the name of the room that we're currently in
	var style tcell.Style
	style = style.Background(tcell.NewRGBColor(64, 64, 128))
	style = style.Foreground(tcell.NewRGBColor(0, 0, 0))

	runes := []rune(st.room.Name)

	bounds := img.bounds()
	for i := 0; i < bounds.Width; i++ {
		if i < len(runes) {
			img.setTile(i, 0, tile{r: runes[i], style: style})
		} else {
			img.setTile(i, 0, tile{r: ' ', style: style})
		}
	}
}

func (v *gameView) setFocus(yes bool) { v.inFocus = yes }

type move struct {
	x int
	y int
}

func (m move) exec(ui *UI) {
	go ui.client.Send(sim.Move{m.x, m.y})
}

type lookAt struct {
	x       int
	y       int
	results *sim.Look
}

func (l *lookAt) exec(ui *UI) {
	go func() {
		res, err := ui.client.Send(sim.LookAt{l.x, l.y})
		if err != nil {
			ui.Error("look error: %v", err)
			return
		}

		look, ok := res.Body.(*sim.Look)
		if !ok {
			ui.Error("look response is not look: %v", res.Body)
			return
		}

		l.results = look
		ui.misc <- func(ui *UI) {
			ui.Info("setting detail view to %T", l)
			ui.state.detail = l
		}
	}()
}

func (l *lookAt) handleEvent(e tcell.Event) change {
	return nil
}

func (l *lookAt) draw(img canvas, st *state) {
	if len(l.results.Here) == 0 {
		writeString(img, "there's nothing here...", math.Vec{0, 0}, tcell.StyleDefault)
		return
	}

	for i, item := range l.results.Here {
		writeString(img, item.Name, math.Vec{0, i}, tcell.StyleDefault)
	}
}

type pickup struct {
	x, y int
}

func (p pickup) exec(ui *UI) {
	go func() {
		res, err := ui.client.Send(sim.Pickup{p.x, p.y})
		if err != nil {
			ui.state.detail = textView(err.Error())
			return
		}

		switch v := res.Body.(type) {
		case *sim.Pickedup:
			ui.misc <- func(ui *UI) {
				ui.state.detail = textView(fmt.Sprintf("you picked up: %s", v.Name))
				ui.state.inventory.items = append(ui.state.inventory.items, item{name: v.Name})
			}
		case wire.Error:
			ui.misc <- func(ui *UI) {
				ui.state.detail = textView(v.Error())
			}
		default:
			ui.misc <- func(ui *UI) {
				ui.state.detail = textView(fmt.Sprintf("unexpected pickup response type: %T", res.Body))
			}
		}
	}()
}
