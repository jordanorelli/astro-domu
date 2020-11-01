package math

type Bounds struct {
	Origin Vec `json:"origin"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func CreateBounds(width, height int) Bounds {
	return Bounds{
		Origin: Vec{0, 0},
		Width:  width,
		Height: height,
	}
}

func (b Bounds) Area() int { return b.Width * b.Height }

func (b Bounds) Contains(v Vec) bool {
	return v.X >= b.Origin.X &&
		v.X < b.Origin.X+b.Width &&
		v.Y >= b.Origin.Y &&
		v.Y < b.Height
}
