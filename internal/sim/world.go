package sim

import (
	"time"

	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

// world is the entire simulated world. A world consists of many rooms.
type world struct {
	*blammo.Log
	Inbox chan Request

	rooms        []room
	done         chan bool
	lastEntityID int
	players      map[string]*player
}

func newWorld(log *blammo.Log) *world {
	bounds := math.CreateRect(10, 10)
	foyer := room{
		Log:     log.Child("foyer"),
		name:    "foyer",
		Rect:    bounds,
		tiles:   make([]tile, bounds.Area()),
		players: make(map[string]*player),
	}
	foyer.tiles[55].here = &entity{
		ID:       777,
		Position: math.Vec{5, 5},
		Glyph:    'd',
		behavior: doNothing{},
	}
	log.Info("created foyer with bounds: %#v having width: %d height: %d area: %d", foyer.Rect, foyer.Width, foyer.Height, foyer.Area())
	return &world{
		Log:     log,
		rooms:   []room{foyer},
		done:    make(chan bool),
		Inbox:   make(chan Request),
		players: make(map[string]*player),
	}
}

func (w *world) run(hz int) {
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

func (w *world) stop() error {
	w.Info("stopping simulation")
	w.done <- true
	w.Info("simulation stopped")
	return nil
}

func (w *world) tick(d time.Duration) {
	for _, r := range w.rooms {
		r.update(d)
	}
}
