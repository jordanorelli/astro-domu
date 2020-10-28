package sim

import (
	"time"

	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

// player represents a player character in the simulation
type player struct {
	*blammo.Log
	sessionID int
	outbox    chan wire.Response
	entityID  int
	pending   []Request
}

func (p *player) update(dt time.Duration) {}

func (p *player) id() int { return p.entityID }

type Move [2]int

func (Move) NetTag() string { return "move" }

// SpawnPlayer is a request to spawn a player
type SpawnPlayer struct {
	Outbox chan wire.Response
}

func (s SpawnPlayer) exec(w *World, from string) {
	w.Info("spawn player requested for: %s", from)
}

func (SpawnPlayer) NetTag() string { return "player/spawn" }

// PlayerSpawned is an announcement that a player has spawned
type PlayerSpawned struct {
	Name string
}

func init() {
	wire.Register(func() wire.Value { return new(Move) })
}
