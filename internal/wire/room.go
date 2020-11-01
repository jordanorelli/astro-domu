package wire

import (
	"github.com/jordanorelli/astro-domu/internal/math"
)

// Room represents a 2-dimensional coordinate space.
type Room struct {
	Name     string          `json:"name"`
	Bounds   math.Bounds     `json:"bounds"`
	Entities map[int]*Entity `json:"entities"`
}

func (r Room) Width() int  { return r.Bounds.Width() }
func (r Room) Height() int { return r.Bounds.Height() }
