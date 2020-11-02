package sim

import (
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type player struct {
	*blammo.Log
	*room
	name    string
	outbox  chan wire.Response
	pending *Request
	avatar  *entity
}

type Move math.Vec

func (Move) NetTag() string { return "move" }

func (m *Move) exec(r *room, p *player, seq int) result {
	pos := p.avatar.Position
	target := pos.Add(math.Vec(*m))
	p.Info("running move for player %s from %v to %v", p.name, *m, target)
	if !p.room.Contains(target) {
		return result{reply: wire.Errorf("target cell (%d, %d) is out of bounds", target.X, target.Y)}
	}

	currentTile := r.getTile(pos)
	nextTile := r.getTile(target)
	if nextTile.here != nil {
		return result{reply: wire.Errorf("target cell (%d, %d) is occupied", target.X, target.Y)}
	}

	currentTile.here, nextTile.here = nil, p.avatar
	p.avatar.Position = target
	return result{reply: wire.OK{}}
}

// SpawnPlayer is a request to spawn a player
type SpawnPlayer struct {
	Outbox chan wire.Response
	Name   string
	queued bool
}

var lastEntityID = 0

func (s *SpawnPlayer) exec(r *room, p *player, seq int) result {
	if !s.queued {
		r.Info("spawn player requested for: %s", s.Name)

		if _, ok := r.players[s.Name]; ok {
			s.Outbox <- wire.ErrorResponse(seq, "a player is already logged in as %q", s.Name)
			return result{}
		}

		lastEntityID++
		avatar := &entity{
			ID:       lastEntityID,
			Position: math.Vec{0, 0},
			Glyph:    '@',
			behavior: doNothing{},
		}
		p := &player{
			Log:    r.Log.Child("players").Child(s.Name),
			room:   r,
			name:   s.Name,
			outbox: s.Outbox,
			avatar: avatar,
		}
		p.pending = &Request{Seq: seq, From: s.Name, Wants: s}
		r.players[s.Name] = p
		r.tiles[0].here = p.avatar
		s.queued = true
		return result{}
	}

	welcome := wire.Welcome{
		Rooms:   make(map[string]wire.Room),
		Players: make(map[string]wire.Player),
	}
	ents := make(map[int]wire.Entity)
	for id, e := range r.allEntities() {
		ents[id] = wire.Entity{
			ID:       id,
			Position: e.Position,
			Glyph:    e.Glyph,
		}
	}
	welcome.Rooms[r.name] = wire.Room{
		Name:     r.name,
		Rect:     r.Rect,
		Entities: ents,
	}
	for _, p := range r.players {
		welcome.Players[p.name] = wire.Player{
			Name:   p.name,
			Avatar: p.avatar.ID,
			Room:   r.name,
		}
	}
	return result{
		reply: welcome,
	}
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
}
