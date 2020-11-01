package wire

type Frame struct {
	Entities map[int]Entity `json:"entities"`
	Players  map[string]int `json:"players"`
}

func (Frame) NetTag() string { return "frame" }

func init() {
	Register(func() Value { return new(Frame) })
}
