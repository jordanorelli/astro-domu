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
	for _, p := range r.players {
		for _, req := range p.pending {
			res := req.Wants.exec(r, p, req.Seq)
			p.outbox <- wire.Response{Re: req.Seq, Body: res.reply}
			for _, p2 := range r.players {
				if p2 == p {
					continue
				}
				p2.outbox <- wire.Response{Body: res.reply}
			}
		}
		p.pending = p.pending[0:0]
	}

	for _, t := range r.tiles {
		if t.here != nil {
			t.here.update(dt)
		}
	}
}

func (r *room) allEntities() []entity {
	all := make([]entity, 0, 4)
	for _, t := range r.tiles {
		if t.here != nil {
			all = append(all, *t.here)
		}
	}
	return all
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
