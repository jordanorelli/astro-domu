package sim

import (
	"time"

	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type room struct {
	*blammo.Log
	name    string
	origin  point
	width   int
	height  int
	tiles   []tile
	players map[string]*player
}

func (r *room) update(dt time.Duration) {
	// announcements := make([]result, 0, 8)

	for _, p := range r.players {
		for _, req := range p.pending {
			res := req.Wants.exec(r, p.name, req.Seq)
			p.outbox <- wire.Response{Re: req.Seq, Body: res.reply}
		}
		p.pending = p.pending[0:0]
	}

	for _, t := range r.tiles {
		for _, e := range t.contents {
			e.update(dt)
		}
	}
}

func (r *room) addPlayer(p *player) {
	r.players[p.name] = p
}

func (r *room) removePlayer(name string) bool {
	if _, ok := r.players[name]; ok {
		delete(r.players, name)
		return true
	}
	return false
}
