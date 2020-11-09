package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
)

type canvas interface {
	getTile(x, y int) tile
	setTile(int, int, tile)
	bounds() math.Rect
}

func writeString(img canvas, s string, start math.Vec, style tcell.Style) {
	for i, r := range []rune(s) {
		img.setTile(start.X+i, start.Y, tile{r: r, style: style})
	}
}

func cut(img canvas, bounds math.Rect) canvas {
	return &section{img: img, frame: bounds}
}

func fill(img canvas, style tcell.Style) {
	bounds := img.bounds()
	for x := 0; x < bounds.Width; x++ {
		for y := 0; y < bounds.Height; y++ {
			img.setTile(x, y, tile{r: ' ', style: style})
		}
	}
}

type section struct {
	img   canvas
	frame math.Rect
}

func (s *section) getTile(x, y int) tile {
	if x < 0 || x >= s.frame.Width || y < 0 || y >= s.frame.Height {
		return tile{}
	}
	return s.img.getTile(x+s.frame.Origin.X, y+s.frame.Origin.Y)
}

func (s *section) setTile(x, y int, t tile) {
	if x < 0 || x >= s.frame.Width || y < 0 || y >= s.frame.Height {
		return
	}
	s.img.setTile(x+s.frame.Origin.X, y+s.frame.Origin.Y, t)
}

func (s *section) bounds() math.Rect { return s.frame }
