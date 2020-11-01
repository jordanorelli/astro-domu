package sim

import (
	"time"

	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type room struct {
	*blammo.Log
	name string
	math.Bounds
	tiles   []tile
	players map[string]*player
}

func (r *room) update(dt time.Duration) {
	for _, p := range r.players {
		if p.pending == nil {
			continue
		}
		req := p.pending
		p.pending = nil

		res := req.Wants.exec(r, p, req.Seq)
		p.outbox <- wire.Response{Re: req.Seq, Body: res.reply}
		if res.announce != nil {
			for _, p2 := range r.players {
				if p2 == p {
					continue
				}
				p2.outbox <- wire.Response{Body: res.announce}
			}
		}
	}

	for _, t := range r.tiles {
		if t.here != nil {
			t.here.update(dt)
		}
	}
}

func (r *room) allEntities() map[int]*entity {
	all := make(map[int]*entity, 4)
	for _, t := range r.tiles {
		if t.here != nil {
			e := t.here
			all[e.ID] = e
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

func (r *room) getTile(pos math.Vec) *tile {
	if !r.Contains(pos) {
		return nil
	}
	n := pos.X*r.Width + pos.Y
	return &r.tiles[n]
}
