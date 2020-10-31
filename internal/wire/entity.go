package wire

import (
	"github.com/jordanorelli/astro-domu/internal/math"
)

type Entity struct {
	Position math.Vec `json:"position"`
	Glyph    rune     `json:"glyph"`
}

func (Entity) NetTag() string { return "entity" }

func init() {
	Register(func() Value { return new(Entity) })
}
