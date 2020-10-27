package sim

// tile is an individual cell within the world simulation. Everything happens
// on a tile.
type tile struct {
	// floor is the surface for the tile. All things sit atop the floor. The
	// floor informs clients how to draw the tile in the event that the tile's
	// contents are empty. The floor also determins whether or not the tile is
	// traversable.
	floor floor

	// contents is all of the entities on this tile. A given tile may have many
	// entities.
	contents map[int]entity
}

func (t *tile) addEntity(e entity) {
	if t.contents == nil {
		t.contents = make(map[int]entity)
	}
	t.contents[e.id()] = e
}

func (t *tile) removeEntity(id int) entity {
	if e, here := t.contents[id]; here {
		delete(t.contents, id)
		return e
	}
	return nil
}
