package sim

import (
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

// player is a player session in the simulation, eek
type player struct {
	*blammo.Log
	*room
	name    string
	outbox  chan wire.Response
	pending []Request
	entity  *entity
}

type Move math.Vec

func (Move) NetTag() string { return "move" }

func (m *Move) exec(r *room, p *player, seq int) result {
	pos := p.entity.Position
	target := pos.Add(math.Vec(*m))
	p.Info("running move for player %s from %v to %v", p.name, *m, target)
	if target.X >= r.width || target.X < 0 {
		return result{reply: wire.Errorf("target cell (%d, %d) is out of bounds", target.X, target.Y)}
	}
	if target.Y >= r.height || target.Y < 0 {
		return result{reply: wire.Errorf("target cell (%d, %d) is out of bounds", target.X, target.Y)}
	}
	n := target.X*r.width + target.Y
	if r.tiles[n].here != nil {
		return result{reply: wire.Errorf("target cell (%d, %d) is occupied", target.X, target.Y)}
	}
	r.tiles[p.entity.Position.X*r.width+p.entity.Position.Y].here = nil
	p.entity.Position = target
	r.tiles[n].here = p.entity
	e := wire.Entity{
		Position: p.entity.Position,
		Glyph:    '@',
	}
	return result{reply: e, announce: e}
}

// SpawnPlayer is a request to spawn a player
type SpawnPlayer struct {
	Outbox chan wire.Response
	Name   string
	queued bool
}

var lastEntityID = 0

func (s *SpawnPlayer) exec(r *room, _ *player, seq int) result {
	if !s.queued {
		r.Info("spawn player requested for: %s", s.Name)

		if _, ok := r.players[s.Name]; ok {
			s.Outbox <- wire.ErrorResponse(seq, "a player is already logged in as %q", s.Name)
			return result{}
		}

		lastEntityID++
		p := &player{
			Log:     r.Log.Child("players").Child(s.Name),
			room:    r,
			name:    s.Name,
			outbox:  s.Outbox,
			pending: make([]Request, 0, 32),
			entity: &entity{
				ID:       lastEntityID,
				Position: math.Vec{0, 0},
				Glyph:    '@',
				behavior: doNothing{},
			},
		}
		p.pending = append(p.pending, Request{Seq: seq, From: s.Name, Wants: s})
		r.players[s.Name] = p
		r.tiles[0].here = p.entity
		s.queued = true
		return result{}
	}

	var welcome wire.Welcome
	welcome.Room.Width = r.width
	welcome.Room.Height = r.height
	welcome.Room.Origin = math.Vec{0, 0}
	return result{reply: welcome}
}

func (SpawnPlayer) NetTag() string { return "player/spawn" }

// PlayerSpawned is an announcement that a player has spawned
type PlayerSpawned struct {
	Name     string `json:"name"`
	Position [2]int `json:"position"`
}

func (PlayerSpawned) NetTag() string { return "player/spawned" }

func init() {
	wire.Register(func() wire.Value { return new(Move) })
	// wire.Register(func() wire.Value { return new(pawn) })
}
