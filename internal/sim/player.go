package sim

import (
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
	entity  *Entity
}

type Move [2]int

func (Move) NetTag() string { return "move" }

func (m *Move) exec(r *room, p *player, seq int) result {
	pos := p.entity.Position
	target := [2]int{pos[0] + m[0], pos[1] + m[1]}
	p.Info("running move for player %s from %v to %v", p.name, *m, target)
	if target[0] >= r.width || target[0] < 0 {
		return result{reply: wire.Errorf("target cell (%d, %d) is out of bounds", target[0], target[1])}
	}
	if target[1] >= r.height || target[1] < 0 {
		return result{reply: wire.Errorf("target cell (%d, %d) is out of bounds", target[0], target[1])}
	}
	n := target[1]*r.width + target[0]
	if r.tiles[n].here != nil {
		return result{reply: wire.Errorf("target cell (%d, %d) is occupied", target[0], target[1])}
	}
	p.entity.Position = target
	return result{reply: p.entity, announce: p.entity}
}

// SpawnPlayer is a request to spawn a player
type SpawnPlayer struct {
	Outbox chan wire.Response
	Name   string
	queued bool
}

func (s *SpawnPlayer) exec(r *room, _ *player, seq int) result {
	if !s.queued {
		r.Info("spawn player requested for: %s", s.Name)

		if _, ok := r.players[s.Name]; ok {
			s.Outbox <- wire.ErrorResponse(seq, "a player is already logged in as %q", s.Name)
			return result{}
		}

		p := &player{
			Log:     r.Log.Child("players").Child(s.Name),
			room:    r,
			name:    s.Name,
			outbox:  s.Outbox,
			pending: make([]Request, 0, 32),
			entity: &Entity{
				ID:       999,
				Position: [2]int{0, 0},
				Glyph:    '@',
			},
		}
		p.pending = append(p.pending, Request{Seq: seq, From: s.Name, Wants: s})
		r.players[s.Name] = p
		s.queued = true
		return result{}
	}

	return result{
		reply: Welcome{
			Room:     r.name,
			Size:     [2]int{r.width, r.height},
			Contents: r.allEntities(),
		},
	}
}

func (SpawnPlayer) NetTag() string { return "player/spawn" }

// PlayerSpawned is an announcement that a player has spawned
type PlayerSpawned struct {
	Name     string `json:"name"`
	Position [2]int `json:"position"`
}

func (PlayerSpawned) NetTag() string { return "player/spawned" }

type Welcome struct {
	Room     string   `json:"room"`
	Size     [2]int   `json:"size"`
	Contents []Entity `json:"contents"`
}

/*

{
	"name": "foyer",
	"width": 10,
	"height": 10,
	"contents": [
	  [5, 3, 10],
	],
	"entities": [
	  [3, "pawn", {"name": "bones"}],
	  [10, "pawn", {"name": "steve"}]
	]
}

*/

func (Welcome) NetTag() string { return "player/welcome" }

func init() {
	wire.Register(func() wire.Value { return new(Move) })
	wire.Register(func() wire.Value { return new(Welcome) })
	// wire.Register(func() wire.Value { return new(pawn) })
}
