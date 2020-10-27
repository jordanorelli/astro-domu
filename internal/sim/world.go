package sim

import (
	"strconv"
	"time"

	"github.com/jordanorelli/blammo"
)

// World is the entire simulated world. A world consists of many rooms.
type World struct {
	*blammo.Log
	rooms        []room
	done         chan bool
	lastEntityID int
}

func NewWorld(log *blammo.Log) *World {
	foyer := room{
		Log:    log.Child("foyer"),
		name:   "foyer",
		origin: point{0, 0},
		width:  10,
		height: 10,
		tiles:  make([]tile, 100),
	}
	return &World{
		Log:   log,
		rooms: []room{foyer},
		done:  make(chan bool),
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

func (w *World) SpawnPlayer(id int) int {
	w.lastEntityID++
	r := w.rooms[0]
	w.Info("spawning player with id: %d into room %q", id, r.name)
	t := &r.tiles[0]
	p := player{
		Log:       w.Child("players").Child(strconv.Itoa(id)),
		sessionID: id,
		entityID:  w.lastEntityID,
	}
	t.addEntity(&p)
	return p.entityID
}

func (w *World) DespawnPlayer(id int) {
	w.Info("despawning player with id: %d", id)
	for _, r := range w.rooms {
		for _, t := range r.tiles {
			if e := t.removeEntity(id); e != nil {
				w.Info("player removed from room %q", r.name)
				return
			}
		}
	}
	w.Error("player was not found in any room")
}

func (w *World) tick(d time.Duration) {
	w.Info("tick. elapsed: %v", d)
	for _, r := range w.rooms {
		r.update(d)
	}
}
