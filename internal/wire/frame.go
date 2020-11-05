package wire

import "github.com/jordanorelli/astro-domu/internal/math"

type Frame struct {
	RoomName string         `json:"room_name"`
	RoomSize math.Rect      `json:"room_size"`
	Entities map[int]Entity `json:"entities"`
	Players  map[string]int `json:"players"`
}

func (Frame) NetTag() string { return "frame" }

func init() {
	Register(func() Value { return new(Frame) })
}
