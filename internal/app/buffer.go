package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
)

// buffer is a rect of tiles
type buffer struct {
	width  int
	height int
	tiles  []tile
}

func newBuffer(width, height int) *buffer {
	b := &buffer{width: width, height: height}
	b.clear()
	return b
}

func (b *buffer) set(x, y int, t tile) bool {
	n := y*b.width + x
	if n >= len(b.tiles) {
		return false
	}
	b.tiles[n] = t
	return true
}

func (b *buffer) get(x, y int) (tile, bool) {
	n := y*b.width + x
	if n >= len(b.tiles) {
		return tile{}, false
	}
	return b.tiles[n], true
}

func (b *buffer) writeString(s string, start math.Vec, style tcell.Style) {
	for i, r := range []rune(s) {
		if !b.set(start.X+i, start.Y, tile{r: r, style: style}) {
			break
		}
	}
}

func (b *buffer) clear() { b.tiles = make([]tile, b.width*b.height) }

func (b *buffer) blit(s tcell.Screen, offset math.Vec) {
	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			t, ok := b.get(x, y)
			if ok {
				s.SetContent(x+offset.X, y+offset.Y, t.r, nil, t.style)
			}
		}
	}
}

func (b *buffer) fill(style tcell.Style) {
	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			b.set(x, y, tile{r: ' ', style: style})
		}
	}
}
