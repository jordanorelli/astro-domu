package sim

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

// world is the entire simulated world. A world consists of many rooms.
type world struct {
	*blammo.Log
	inbox   chan Request
	connect chan connect
	nextID  chan int

	done         chan bool
	lastEntityID int
	rooms        map[string]*room
	players      map[string]*player
}

type connect struct {
	conn   *websocket.Conn
	login  wire.Login
	failed chan error
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

	foyer.addEntity(&entity{
		ID:       -1,
		Position: math.Vec{5, 5},
		Glyph:    'd',
		solid:    true,
		behavior: doNothing{},
	})

	foyer.addEntity(&entity{
		ID:       -2,
		Position: math.Vec{9, 0},
		Glyph:    '+',
		behavior: doNothing{},
	})

	foyer.addEntity(&entity{
		ID:       -3,
		Position: math.Vec{9, 1},
		Glyph:    '-',
		behavior: doNothing{},
	})

	foyer.addEntity(&entity{
		ID:       -4,
		Position: math.Vec{9, 5},
		Glyph:    '◇',
		behavior: &door{
			Log: log.Child("door"),
		},
	})

	hall := room{
		Log:     log.Child("hall"),
		name:    "hall",
		Rect:    math.CreateRect(20, 4),
		tiles:   make([]tile, bounds.Area()),
		players: make(map[string]*player),
	}

	log.Info("created foyer with bounds: %#v having width: %d height: %d area: %d", foyer.Rect, foyer.Width, foyer.Height, foyer.Area())
	return &world{
		Log: log,
		rooms: map[string]*room{
			"foyer": &foyer,
			"hall":  &hall,
		},
		done:    make(chan bool),
		inbox:   make(chan Request),
		connect: make(chan connect),
		players: make(map[string]*player),
		nextID:  make(chan int),
	}
}

func (w *world) run(hz int) {
	defer w.Info("simulation has exited run loop")

	go func() {
		lastID := 1
		for {
			select {
			case <-w.done:
				return
			case w.nextID <- lastID:
				lastID++
			}
		}
	}()

	period := time.Second / time.Duration(hz)
	w.Info("starting world with a tick rate of %dhz, frame duration of %v", hz, period)
	ticker := time.NewTicker(period)
	lastTick := time.Now()

	for {
		select {
		case c := <-w.connect:
			w.register(c)
			w.Info("finished registration for: %v", c)

		case req := <-w.inbox:
			w.Info("read request off of inbox: %v", req)
			w.handleRequest(req)
			w.Info("finished handling request: %v", req)

		case <-ticker.C:
			w.tick(time.Since(lastTick))
			lastTick = time.Now()

		case <-w.done:
			for _, p := range w.players {
				p.stop <- true
			}
			return
		}
	}
}

func (w *world) handleRequest(req Request) {
	w.Info("read from inbox: %v", req)

	if req.From == "" {
		w.Error("request has no from: %v", req)
		return
	}

	p, ok := w.players[req.From]
	if !ok {
		w.Error("received non login request of type %T from unknown player %q", req.Wants, req.From)
		return
	}

	if p.pending == nil {
		p.pending = &req
	} else {
		p.outbox <- wire.ErrorResponse(req.Seq, "you already have a request for this frame")
	}
}

func (w *world) register(c connect) {
	w.Info("register: %#v", c.login)
	foyer := w.rooms["foyer"]
	if len(foyer.players) >= 100 {
		c.failed <- fmt.Errorf("room is full")
		close(c.failed)
		return
	}

	p := player{
		Log:    w.Log.Child("players").Child(c.login.Name),
		name:   c.login.Name,
		outbox: make(chan wire.Response, 8),
		pending: &Request{
			From:  c.login.Name,
			Seq:   1,
			Wants: &spawnPlayer{},
		},
	}
	foyer.players[c.login.Name] = &p
	w.players[c.login.Name] = &p

	w.Info("starting player...")
	p.start(w.inbox, c.conn, foyer)
}

func (w *world) stop() error {
	w.Info("stopping simulation")
	w.done <- true
	w.Info("simulation stopped")
	return nil
}

func (w *world) tick(d time.Duration) {
	// run all player effects
	for _, r := range w.rooms {
		for _, p := range r.players {
			if p.pending == nil {
				continue
			}
			req := p.pending
			p.pending = nil

			res := req.Wants.exec(w, r, p, req.Seq)
			if res.reply != nil {
				p.send(wire.Response{Re: req.Seq, Body: res.reply})
			} else {
				p.send(wire.Response{Re: req.Seq, Body: wire.OK{}})
			}
		}
	}

	// run all object updates
	for _, r := range w.rooms {
		for _, t := range r.tiles {
			t.update(d)
		}
	}

	// check all overlapping entities
	for _, r := range w.rooms {
		for _, t := range r.tiles {
			t.overlaps()
		}
	}

	// send frame data to all players
	for _, r := range w.rooms {
		frame := wire.Frame{
			Entities: r.allEntities(),
			Players:  r.playerAvatars(),
		}

		for _, p := range r.players {
			p.send(wire.Response{Body: frame})
		}
	}
}
