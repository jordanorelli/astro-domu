package ui

// buffer is a rect of tiles
type buffer struct {
	width  int
	height int
	tiles  []tile
}

func newBuffer(width, height int) *buffer {
	return &buffer{
		width:  width,
		height: height,
		tiles:  make([]tile, width*height),
	}
}

func (b *buffer) set(x, y int, t tile) { b.tiles[y*b.width+x] = t }
func (b *buffer) get(x, y int) tile    { return b.tiles[y*b.width+x] }
