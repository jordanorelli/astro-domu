package sim

import (
	"time"

	"github.com/jordanorelli/astro-domu/internal/wire"
)

type Entity struct {
	ID       int    `json:"id"`
	Position [2]int `json:"pos"`
	Glyph    rune   `json:"glyph"`
	behavior
}

func (Entity) NetTag() string { return "entity" }

func init() {
	wire.Register(func() wire.Value { return new(Entity) })
}

type behavior interface {
	// update is the standard tick function
	update(time.Duration)
}

type doNothing struct{}

func (d doNothing) update(time.Duration) {}
