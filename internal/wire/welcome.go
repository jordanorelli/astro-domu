package wire

import (
	"github.com/jordanorelli/astro-domu/internal/math"
)

type Welcome struct {
	Room struct {
		Origin math.Vec `json:"origin"`
		Width  int      `json:"width"`
		Height int      `json:"height"`
	} `json:"room"`
	Entities map[int]Entity `json:"entities"`
}

func (Welcome) NetTag() string { return "welcome" }

func init() {
	Register(func() Value { return new(Welcome) })
}
