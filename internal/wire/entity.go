package wire

import (
	"github.com/jordanorelli/astro-domu/internal/math"
)

type Entity struct {
	ID       int      `json:"id"`
	Position math.Vec `json:"position"`
	Glyph    rune     `json:"glyph"`
}

func (Entity) NetTag() string { return "entity" }

type UpdateEntity struct {
	Room     string   `json:"room"`
	ID       int      `json:"id"`
	Position math.Vec `json:"position"`
	Glyph    rune     `json:"glyph"`
}

func (UpdateEntity) NetTag() string { return "entity/updated" }

func init() {
	Register(func() Value { return new(Entity) })
	Register(func() Value { return new(UpdateEntity) })
}
