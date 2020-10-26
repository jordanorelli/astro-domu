package sim

import (
	"time"

	"github.com/jordanorelli/blammo"
)

// World is the entire simulated world. A world consists of many rooms.
type World struct {
	*blammo.Log
	rooms []room
}

func NewWorld(log *blammo.Log) *World {
	return &World{
		Log: log,
		rooms: []room{
			room{
				origin: point{0, 0},
				width:  10,
				height: 10,
			},
		},
	}
}

func (w *World) Run(hz int) {
	period := time.Second / time.Duration(hz)
	w.Info("starting world with a tick rate of %dhz, frame duration of %v", hz, period)
	ticker := time.NewTicker(period)
	lastTick := time.Now()
	for {
		select {
		case <-ticker.C:
			w.tick(time.Since(lastTick))
			lastTick = time.Now()
		}
	}
}

func (w *World) tick(d time.Duration) {
	w.Info("tick. elapsed: %v", d)
}
