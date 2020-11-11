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
		ID:         -1,
		Position:   math.Vec{5, 4},
		Glyph:      'o',
		solid:      true,
		pickupable: true,
		name:       "a rock",
		behavior:   doNothing{},
	})

	foyer.addEntity(&entity{
		ID:         -2,
		Position:   math.Vec{5, 5},
		Glyph:      'o',
		solid:      true,
		pickupable: true,
		name:       "another rock",
		behavior:   doNothing{},
	})

	foyer.addEntity(&entity{
		ID:         -3,
		Position:   math.Vec{5, 6},
		Glyph:      'o',
		solid:      true,
		pickupable: true,
		name:       "yet another rock (YAR)",
		behavior:   doNothing{},
	})

	foyer.addEntity(&entity{
		ID:       -4,
		Position: math.Vec{9, 5},
		Glyph:    '◇',
		name:     "Door to Hall",
		behavior: &door{
			Log:  log.Child("door"),
			to:   "hall",
			exit: -5,
		},
	})

	hall := room{
		Log:     log.Child("hall"),
		name:    "hall",
		Rect:    math.CreateRect(20, 4),
		tiles:   make([]tile, 80),
		players: make(map[string]*player),
	}

	hall.addEntity(&entity{
		ID:       -5,
		Position: math.Vec{0, 2},
		Glyph:    '◇',
		name:     "Door to Foyer",
		behavior: &door{
			Log:  log.Child("door"),
			to:   "foyer",
			exit: -4,
		},
	})

	kitchen := room{
		Log:     log.Child("kitchen"),
		name:    "kitchen",
		Rect:    math.CreateRect(6, 12),
		tiles:   make([]tile, 72),
		players: make(map[string]*player),
	}

	foyer.addEntity(&entity{
		ID:       -6,
		Position: math.Vec{9, 7},
		Glyph:    '◇',
		name:     "Door to Kitchen",
		behavior: &door{
			Log:  log.Child("door"),
			to:   "kitchen",
			exit: -7,
		},
	})

	kitchen.addEntity(&entity{
		ID:       -7,
		Position: math.Vec{0, 2},
		Glyph:    '◇',
		name:     "Door to Foyer",
		behavior: &door{
			Log:  log.Child("door"),
			to:   "foyer",
			exit: -6,
		},
	})

	log.Info("created foyer with bounds: %#v having width: %d height: %d area: %d", foyer.Rect, foyer.Width, foyer.Height, foyer.Area())
	return &world{
		Log: log,
		rooms: map[string]*room{
			"foyer":   &foyer,
			"hall":    &hall,
			"kitchen": &kitchen,
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

	r := w.rooms["foyer"]
	if hall, ok := w.rooms["hall"]; ok {
		if len(hall.players) < len(r.players) {
			r = hall
		}
	}

	if len(r.players) >= 100 {
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
			Seq:   90,
			Wants: &spawnPlayer{},
		},
	}
	r.players[c.login.Name] = &p
	w.players[c.login.Name] = &p

	w.Info("starting player...")
	p.start(w.inbox, c.conn, r)
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
			RoomName: r.name,
			RoomSize: r.Rect,
			Entities: r.allEntities(),
			Players:  r.playerAvatars(),
		}

		delta := r.lastFrame.Diff(frame)
		if delta != nil {
			w.Info("%s delta: %s", r.name, delta)
		}
		r.lastFrame = frame

		for _, p := range r.players {
			switch {
			case p.fullSync:
				p.send(wire.Response{Body: frame})
				p.fullSync = false
			case delta != nil:
				p.send(wire.Response{Body: delta})
			}
		}
	}
}
