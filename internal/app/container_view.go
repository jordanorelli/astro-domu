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
	switch v := e.(type) {
	case *tcell.EventKey:
		if v.Key() == tcell.KeyTab {
			c.nextFocus()
			return nil
		}
		return c.children[c.focussed].handleEvent(e)

	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}

	return nil
}

func (c *containerView) nextFocus() {
	setFocus := func(i int, enabled bool) bool {
		n := c.children[i]
		if v, ok := n.view.(focusable); ok {
			v.setFocus(enabled)
			return true
		}
		return false
	}

	for start, next := c.focussed, c.focussed+1; next != start; next++ {
		if next >= len(c.children) {
			next = 0
		}
		if setFocus(next, true) {
			c.focussed = next
			setFocus(start, false)
			return
		}
	}
}

func (c *containerView) draw(img canvas, st *state) {
	bounds := img.bounds()

	for _, n := range c.children {
		n.draw(cut(img, c.scaleFrame(n.frame, bounds)), st)
	}
}

func (c *containerView) scaleFrame(frame math.Rect, bounds math.Rect) math.Rect {
	var (
		xscale = bounds.Width / c.refWidth
		yscale = bounds.Height / c.refHeight
		xrem   = bounds.Width % c.refWidth
		yrem   = bounds.Height % c.refHeight
	)

	xremTaken := math.Min(frame.Origin.X, xrem)
	yremTaken := math.Min(frame.Origin.Y, yrem)

	xinflate := math.Max(xrem-xremTaken, 0)
	yinflate := math.Max(yrem-yremTaken, 0)

	return math.Rect{
		Origin: math.Vec{
			xscale*frame.Origin.X + xremTaken,
			yscale*frame.Origin.Y + yremTaken,
		},
		Width:  xscale*frame.Width + xinflate,
		Height: yscale*frame.Height + yinflate,
	}
}

type node struct {
	view
	frame math.Rect
}
