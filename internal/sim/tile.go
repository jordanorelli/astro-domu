package sim

import (
	"fmt"
	"time"
)

type tile struct {
	floor floor
	here  []*entity
}

func (t tile) String() string {
	ids := make([]int, len(t.here))
	for i, e := range t.here {
		ids[i] = e.ID
	}
	return fmt.Sprintf("{%c %v}", t.floor, ids)
}

func (t *tile) addEntity(e *entity) bool {
	if e.solid {
		for _, other := range t.here {
			if other.solid {
				return false
			}
		}
	}
	t.here = append(t.here, e)
	return true
}

func (t *tile) removeEntity(id int) bool {
	start := len(t.here)
	here := t.here[:0]
	for _, e := range t.here {
		if e.ID != id {
			here = append(here, e)
		}
	}
	t.here = here
	return len(t.here) != start
}

func (t *tile) hasEntity(id int) bool {
	for _, e := range t.here {
		if e.ID == id {
			return true
		}
	}
	return false
}

func (t *tile) isOccupied() bool {
	for _, e := range t.here {
		if e.solid {
			return true
		}
	}
	return false
}

func (t *tile) update(d time.Duration) {
	for _, e := range t.here {
		e.update(e, d)
	}
}

func (t *tile) overlaps() {
	switch len(t.here) {
	case 0:
		return
	case 1:
		t.here[0].overlapping()
		return
	default:
	}

	for _, e := range t.here {
		e.overlapping(t.here...)
	}
}
