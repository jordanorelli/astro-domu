package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
)

type containerView struct {
	children []*node
	focussed int
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
	for _, n := range c.children {
		n.draw(cut(img, n.frame), st)
	}
}

type node struct {
	view
	frame math.Rect
}
