package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
)

type borderedView struct {
	inner    view
	hasFocus bool
}

func (b borderedView) handleEvent(e tcell.Event) change {
	return b.inner.handleEvent(e)
}

func (b borderedView) draw(img canvas, st *state) {
	s := tcell.StyleDefault
	if !b.hasFocus {
		s = s.Foreground(tcell.NewRGBColor(128, 128, 128))
	}

	bounds := img.bounds()
	img.setTile(0, 0, tile{r: '┌', style: s})
	img.setTile(bounds.Width-1, 0, tile{r: '┐', style: s})
	img.setTile(0, bounds.Height-1, tile{r: '└', style: s})
	img.setTile(bounds.Width-1, bounds.Height-1, tile{r: '┘', style: s})
	for x := 1; x < bounds.Width-1; x++ {
		img.setTile(x, 0, tile{r: '─', style: s})
		img.setTile(x, bounds.Height-1, tile{r: '─', style: s})
	}
	for y := 1; y < bounds.Height-1; y++ {
		img.setTile(0, y, tile{r: '│', style: s})
		img.setTile(bounds.Width-1, y, tile{r: '│', style: s})
	}

	b.inner.draw(cut(img, math.Rect{math.Vec{1, 1}, bounds.Width - 2, bounds.Height - 2}), st)
}

func (b *borderedView) setFocus(enabled bool) {
	b.hasFocus = enabled
	if v, ok := b.inner.(focusable); ok {
		v.setFocus(enabled)
	}
}
