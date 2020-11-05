package sim

import "github.com/jordanorelli/astro-domu/internal/wire"

type Effect interface {
	exec(*world, *room, *player, int) result
}

type effect func(*world, *room, *player, int) result

func (f effect) exec(w *world, r *room, p *player, seq int) result { return f(w, r, p, seq) }

type result struct {
	reply wire.Value
}
