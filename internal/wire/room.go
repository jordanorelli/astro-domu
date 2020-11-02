package wire

import (
	"github.com/jordanorelli/astro-domu/internal/math"
)

// Room represents a 2-dimensional coordinate space.
type Room struct {
	Name      string `json:"name"`
	math.Rect `json:"bounds"`
	Entities  map[int]Entity `json:"entities"`
}
