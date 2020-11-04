package sim

type tile struct {
	floor floor
	here  []*entity
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

func (t *tile) removeEntity(id int) {
	here := t.here[:0]
	for _, e := range t.here {
		if e.ID != id {
			here = append(here, e)
		}
	}
	t.here = here
}
