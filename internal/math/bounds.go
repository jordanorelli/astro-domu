package math

import "encoding/json"

type Bounds struct {
	Min Vec `json:"min"`
	Max Vec `json:"max"`
}

func CreateBounds(width, height int) Bounds {
	return Bounds{
		Min: Vec{0, 0},
		Max: Vec{width - 1, height - 1},
	}
}

func (b Bounds) Width() int  { return Abs(b.Max.X - b.Min.X) }
func (b Bounds) Height() int { return Abs(b.Max.Y - b.Min.Y) }
func (b Bounds) Area() int   { return b.Width() * b.Height() }

func (b Bounds) Contains(v Vec) bool {
	return v.X >= b.Min.X && v.X <= b.Max.X && v.Y >= b.Min.Y && v.Y <= b.Max.Y
}

func (b Bounds) MarshalJSON() ([]byte, error) { return json.Marshal([2]Vec{b.Min, b.Max}) }

func (b *Bounds) UnmarshalJSON(buf []byte) error {
	var raw [2]Vec
	if err := json.Unmarshal(buf, &raw); err != nil {
		return err
	}
	b.Min = raw[0]
	b.Max = raw[1]
	return nil
}
