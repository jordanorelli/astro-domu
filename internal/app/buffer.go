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
	b := &buffer{
		width:  width,
		height: height,
		tiles:  make([]tile, width*height),
	}
	return b
}

func (b *buffer) setTile(x, y int, t tile) {
	n := y*b.width + x
	if n >= len(b.tiles) {
		return
	}
	b.tiles[n] = t
}

func (b *buffer) getTile(x, y int) tile {
	n := y*b.width + x
	if n >= len(b.tiles) {
		return tile{}
	}
	return b.tiles[n]
}

func (b *buffer) clear() {
	for i, _ := range b.tiles {
		b.tiles[i] = tile{}
	}
}

func (b *buffer) blit(s tcell.Screen, offset math.Vec) {
	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			t := b.getTile(x, y)
			s.SetContent(x+offset.X, y+offset.Y, t.r, nil, t.style)
		}
	}
}

func (b *buffer) fill(style tcell.Style) {
	for x := 0; x < b.width; x++ {
		for y := 0; y < b.height; y++ {
			b.setTile(x, y, tile{r: ' ', style: style})
		}
	}
}

func (b *buffer) bounds() math.Rect { return math.Rect{Width: b.width, Height: b.height} }
