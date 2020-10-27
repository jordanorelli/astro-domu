package sim

import (
	"time"

	"github.com/jordanorelli/blammo"
)

// player represents a player character in the simulation
type player struct {
	*blammo.Log
	entityID int
}

func (p *player) update(dt time.Duration) {
}

func (p *player) id() int { return p.entityID }
