package math

import (
	"encoding/json"
)

type Vec struct {
	X int
	Y int
}

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
