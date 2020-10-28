package sim

import (
	"time"

	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

// player represents a player character in the simulation
type player struct {
	*blammo.Log
	*room
	name      string
	outbox    chan wire.Response
	sessionID int
	entityID  int
	pending   []Request
}

func (p *player) update(dt time.Duration) {}

func (p *player) id() int { return p.entityID }

type Move [2]int

func (Move) NetTag() string { return "move" }

func (m *Move) exec(r *room, from *player, seq int) {

}

// SpawnPlayer is a request to spawn a player
type SpawnPlayer struct {
	Outbox chan wire.Response
	queued bool
}

func (s *SpawnPlayer) exec(r *room, from string, seq int) result {
	if !s.queued {
		r.Info("spawn player requested for: %s", from)

		if _, ok := r.players[from]; ok {
			s.Outbox <- wire.ErrorResponse(seq, "a player is already logged in as %q", from)
			return result{}
		}

		p := &player{
			Log:     r.Log.Child("players").Child(from),
			room:    r,
			name:    from,
			outbox:  s.Outbox,
			pending: make([]Request, 0, 32),
		}
		p.pending = append(p.pending, Request{Seq: seq, From: from, Wants: s})
		r.players[from] = p
		s.queued = true
		return result{}
	}

	return result{
		reply: Welcome{Room: r.name},
	}
}

func (SpawnPlayer) NetTag() string { return "player/spawn" }

// PlayerSpawned is an announcement that a player has spawned
type PlayerSpawned struct {
	Name string
}

type Welcome struct {
	Room string `json:"room"`
}

func (Welcome) NetTag() string { return "player/welcome" }

func init() {
	wire.Register(func() wire.Value { return new(Move) })
	wire.Register(func() wire.Value { return new(Welcome) })
}
