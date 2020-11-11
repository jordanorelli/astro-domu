package sim

import (
	"time"

	"github.com/jordanorelli/blammo"
)

type door struct {
	*blammo.Log
	to      string
	exit    int
	arrived int
}

func (d *door) update(*entity, time.Duration) {}

func (d *door) exec(w *world, r *room, p *player, seq int) result {
	p.Info("executing door to %q for player %s", d.to, p.name)

	dest, ok := w.rooms[d.to]
	if !ok {
		p.Error("door destination %q does not exist", d.to)
		return result{}
	}
	p.Info("found destination room %q", d.to)

	exit := dest.findEntity(d.exit)
	if exit == nil {
		p.Error("door exit %d does not exist", d.exit)
		return result{}
	}

	exitDoor, ok := exit.behavior.(*door)
	if !ok {
		p.Error("exit entity %d is not a door", d.exit)
		return result{}
	}
	p.Info("found exit door %v", exitDoor)

	t := dest.getTile(exit.Position)
	p.Info("exit tile: %v", t)
	if t.isOccupied() {
		p.Error("destination tile %v is occupied", t)
		return result{}
	}

	p.Info("removing player from room %s", r.name)
	r.removePlayer(p.name)
	p.Info("adding player to room %s", dest.name)
	if t.addEntity(p.avatar) {
		p.Info("added player avatar to tile %v", t)
		exitDoor.arrived = p.avatar.ID
		p.avatar.Position = exit.Position
	} else {
		p.Error("failed to add player avatar to tile %v", t)
	}
	dest.addPlayer(p)
	p.fullSync = true
	return result{}
}

func (d *door) onStartOverlap(e *entity) {
	if e.ID == d.arrived {
		return
	}
	if p, ok := e.behavior.(*player); ok {
		d.Info("player %s start overlap on door to %s", p.name, d.to)
		if p.pending != nil {
			d.Info("player %s starting overlap on door to %s has a pending request", p.name, d.to)
		} else {
			d.Info("player %s starting overlap on door to %s has NO pending request", p.name, d.to)
			p.pending = &Request{Wants: d}
		}
	}
}

func (d *door) onOverlap(e *entity) {
	if e.ID == d.arrived {
		return
	}
	// d.Info("overlap: %v", e)
	if p, ok := e.behavior.(*player); ok {
		d.Info("player %s is continuing overlap on door to %s", p.name, d.to)
		if p.pending != nil {
			d.Info("player %s continuing to overlap door to %s has a pending request", p.name, d.to)
		} else {
			d.Info("player %s continuing to overlap door to %s has NO pending request", p.name, d.to)
			p.pending = &Request{Wants: d}
		}
	}
}

func (d *door) onStopOverlap(e *entity) {
	if e.ID == d.arrived {
		d.arrived = 0
	}
	// d.Info("stop overlap: %v", e)
	if p, ok := e.behavior.(*player); ok {
		d.Info("player %s stepped off of door to %s", p.name, d.to)
	}
}

/*

   ┌──────────┐
   │··········│
   │·····@····│
   │··········│
   │··········│        ┌────────────────────┐
   │··········│        │····················│
   │·····d···◇│--------│◇···················│
   │··········│        │····················│
   │··········│        └────────────────────┘
   │··········│
   │·········◇│
   └──────────┘


   ┌──────────┐
   │··········│
   │·····@····│
   │··········│
   │··········│        ┌────────────────────┐
   │··········│        │····················│
   │·····d····◇--------◇····················│
   │··········│        │····················│
   │··········│        └────────────────────┘
   │··········│
   │·········◇│
   └──────────┘


   ┌──────────┐
   │··········│
   │·····@····│
   │··········│
   │··········│────────────────────┐
   │··········│····················│
   │·····d····◇····················│
   │··········│····················│
   │··········│────────────────────┘
   │··········│
   │·········◇│
   └──────────┘

*/
