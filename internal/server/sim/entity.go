package sim

import (
	"time"

	"github.com/jordanorelli/astro-domu/internal/math"
)

type entity struct {
	ID       int      `json:"id"`
	Position math.Vec `json:"pos"`
	Glyph    rune     `json:"glyph"`
	behavior
}

type behavior interface {
	// update is the standard tick function
	update(time.Duration)
}

type doNothing struct{}

func (d doNothing) update(time.Duration) {}
