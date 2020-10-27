package sim

import (
	"time"

	"github.com/jordanorelli/blammo"
)

type room struct {
	*blammo.Log
	name   string
	origin point
	width  int
	height int
	tiles  []tile
}

func (r *room) update(dt time.Duration) {
	r.Info("updating room")
	for _, t := range r.tiles {
		for _, e := range t.contents {
			e.update(dt)
		}
	}
}
