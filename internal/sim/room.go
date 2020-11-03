package sim

import (
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type room struct {
	*blammo.Log
	name string
	math.Rect
	tiles   []tile
	players map[string]*player
}

func (r *room) allEntities() map[int]wire.Entity {
	all := make(map[int]wire.Entity, 4)
	for _, t := range r.tiles {
		if t.here != nil {
			e := t.here
			all[e.ID] = wire.Entity{
				ID:       e.ID,
				Position: e.Position,
				Glyph:    e.Glyph,
			}
		}
	}
	return all
}

func (r *room) playerAvatars() map[string]int {
	all := make(map[string]int, len(r.players))
	for nick, p := range r.players {
		all[nick] = p.avatar.ID
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
