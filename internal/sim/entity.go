package sim

import (
	"time"

	"github.com/jordanorelli/astro-domu/internal/math"
)

type entity struct {
	ID          int      `json:"id"`
	Position    math.Vec `json:"pos"`
	Glyph       rune     `json:"glyph"`
	name        string
	description string
	solid       bool `json:"-"`
	overlapped  map[int]*entity
	pickupable  bool
	behavior
}

func (e *entity) overlapping(others ...*entity) {
	idx := make(map[int]*entity, len(others))
	for i, _ := range others {
		e2 := others[i]
		if e2 == e {
			continue
		}
		idx[e2.ID] = e2
	}

	for _, e2 := range e.overlapped {
		if idx[e2.ID] == nil {
			e.stopOverlap(e2)
		}
	}

	for _, e2 := range idx {
		e.overlap(e2)
	}

	e.overlapped = idx
}

func (e *entity) overlap(e2 *entity) {
	if e.overlapped == nil {
		e.overlapped = make(map[int]*entity, 4)
	}

	type overlapStarter interface {
		onStartOverlap(*entity)
	}

	if e.behavior != nil && e.overlapped[e2.ID] == nil {
		if b, ok := e.behavior.(overlapStarter); ok {
			b.onStartOverlap(e2)
		}
	}

	type overlapper interface {
		onOverlap(*entity)
	}
	if e.behavior != nil {
		if b, ok := e.behavior.(overlapper); ok {
			b.onOverlap(e2)
		}
	}
}

func (e *entity) stopOverlap(e2 *entity) {
	type overlapStopper interface {
		onStopOverlap(*entity)
	}
	if e.behavior != nil {
		if b, ok := e.behavior.(overlapStopper); ok {
			b.onStopOverlap(e2)
		}
	}
}

type behavior interface {
	// update is the standard tick function
	update(time.Duration)
}

type doNothing struct{}

func (d doNothing) update(time.Duration) {}
