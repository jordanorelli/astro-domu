package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
)

type containerView struct {
	refWidth  int
	refHeight int
	children  []*node
	focussed  int
}

func (c *containerView) handleEvent(e tcell.Event) change {
	switch e.(type) {
	case *tcell.EventKey:
		return c.children[c.focussed].handleEvent(e)

	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}

	return nil
}

func (c *containerView) draw(img canvas, st *state) {
	bounds := img.bounds()

	for _, n := range c.children {
		n.draw(cut(img, c.scaleFrame(n.frame, bounds)), st)
	}
}

func (c *containerView) scaleFrame(frame math.Rect, bounds math.Rect) math.Rect {
	xscale := bounds.Width / c.refWidth
	yscale := bounds.Height / c.refHeight
	return math.Rect{
		Origin: math.Vec{xscale * frame.Origin.X, yscale * frame.Origin.Y},
		Width:  xscale * frame.Width,
		Height: yscale * frame.Height,
	}
}

type node struct {
	view
	frame math.Rect
}
