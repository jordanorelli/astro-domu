package sim

import (
	"time"

	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

// World is the entire simulated world. A world consists of many rooms.
type World struct {
	*blammo.Log
	Inbox chan Request

	rooms        []room
	done         chan bool
	lastEntityID int
	players      map[string]*player
}

func NewWorld(log *blammo.Log) *World {
	bounds := math.CreateBounds(10, 10)
	foyer := room{
		Log:     log.Child("foyer"),
		name:    "foyer",
		origin:  point{0, 0},
		bounds:  bounds,
		tiles:   make([]tile, bounds.Area()),
		players: make(map[string]*player),
	}
	foyer.tiles[55].here = &entity{
		ID:       777,
		Position: math.Vec{5, 5},
		Glyph:    'd',
		behavior: doNothing{},
	}
	return &World{
		Log:     log,
		rooms:   []room{foyer},
		done:    make(chan bool),
		Inbox:   make(chan Request),
		players: make(map[string]*player),
	}
}

func (w *World) Run(hz int) {
	defer w.Info("simulation has exited run loop")

	period := time.Second / time.Duration(hz)
	w.Info("starting world with a tick rate of %dhz, frame duration of %v", hz, period)
	ticker := time.NewTicker(period)
	lastTick := time.Now()

	for {
		select {
		case req := <-w.Inbox:
			w.Info("read from inbox: %v", req)

			if req.From == "" {
				w.Error("request has no from: %v", req)
				break
			}

			if spawn, ok := req.Wants.(*SpawnPlayer); ok {
				if _, ok := w.players[req.From]; ok {
					spawn.Outbox <- wire.ErrorResponse(req.Seq, "a player is already logged in as %q", req.From)
					break
				}
				spawn.exec(&w.rooms[0], nil, req.Seq)
				p := w.rooms[0].players[req.From]
				w.players[req.From] = p
				break
			}

			p, ok := w.players[req.From]
			if !ok {
				w.Error("received non login request of type %T from unknown player %q", req.Wants, req.From)
			}

			if p.pending == nil {
				p.pending = &req
			} else {
				p.outbox <- wire.ErrorResponse(req.Seq, "you already have a request for this frame")
			}

		case <-ticker.C:
			w.tick(time.Since(lastTick))
			lastTick = time.Now()

		case <-w.done:
			return
		}
	}
}

func (w *World) Stop() error {
	w.Info("stopping simulation")
	w.done <- true
	return nil
}

func (w *World) tick(d time.Duration) {
	for _, r := range w.rooms {
		r.update(d)
	}
}
