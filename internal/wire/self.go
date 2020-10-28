package wire

type Self_Move struct {
	Delta bool `json:"delta"`
	X     int  `json:"x"`
	Y     int  `json:"y"`
}

func (Self_Move) NetTag() string { return "self/move" }

func init() {
	Register(func() Value { return new(Self_Move) })
}
