package wire

type Welcome struct {
	Rooms     map[string]Room   `json:"rooms"`
	Players   map[string]Player `json:"players"`
	Inventory []Entity          `json:"inventory"`
}

func (Welcome) NetTag() string { return "welcome" }

func init() {
	Register(func() Value { return new(Welcome) })
}
