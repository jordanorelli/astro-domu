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
	tiles     []tile
	players   map[string]*player
	lastFrame wire.Frame
}

func (r *room) allEntities() map[int]wire.Entity {
	all := make(map[int]wire.Entity, 4)
	for _, t := range r.tiles {
		for _, e := range t.here {
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

func (r *room) findEntity(id int) *entity {
	for _, t := range r.tiles {
		for _, e := range t.here {
			if e.ID == id {
				return e
			}
		}
	}
	return nil
}

func (r *room) addEntity(e *entity) bool {
	t := r.getTile(e.Position)
	if t == nil {
		return false
	}
	return t.addEntity(e)
}

func (r *room) addPlayer(p *player) {
	r.players[p.name] = p
}

func (r *room) removePlayer(name string) bool {
	if p, ok := r.players[name]; ok {
		delete(r.players, name)
		t := r.getTile(p.avatar.Position)
		if t != nil {
			t.removeEntity(p.avatar.ID)
		}
		return true
	}
	return false
}

func (r *room) getTile(pos math.Vec) *tile {
	if !r.Contains(pos) {
		return nil
	}
	n := r.Width*pos.Y + pos.X
	return &r.tiles[n]
}
