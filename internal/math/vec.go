package math

import (
	"encoding/json"
	"fmt"
)

type Vec struct {
	X int
	Y int
}

func (v Vec) String() string { return fmt.Sprintf("(%d, %d)", v.X, v.Y) }

func (v Vec) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]int{v.X, v.Y})
}

func (v *Vec) UnmarshalJSON(b []byte) error {
	var raw [2]int
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	v.X = raw[0]
	v.Y = raw[1]
	return nil
}

func (v Vec) Add(v2 Vec) Vec { return Vec{v.X + v2.X, v.Y + v2.Y} }

func (v Vec) Unit() Vec {
	var out Vec
	switch {
	case v.X < 0:
		out.X = -1
	case v.X > 0:
		out.X = 1
	}

	switch {
	case v.Y < 0:
		out.Y = -1
	case v.Y > 0:
		out.Y = 1
	}
	return out
}

// MDist calculates the manhattan distance between two vectors.
func (v Vec) MDist(v2 Vec) int { return Abs(v.X-v2.X) + Abs(v.Y-v2.Y) }
