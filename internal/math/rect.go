package math

type Rect struct {
	Origin Vec `json:"origin"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func CreateRect(width, height int) Rect {
	return Rect{
		Origin: Vec{0, 0},
		Width:  width,
		Height: height,
	}
}

func (b Rect) Area() int { return b.Width * b.Height }

func (b Rect) Contains(v Vec) bool {
	return v.X >= b.Origin.X &&
		v.X < b.Origin.X+b.Width &&
		v.Y >= b.Origin.Y &&
		v.Y < b.Height
}
