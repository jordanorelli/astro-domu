package wire

type OK struct{}

func (OK) NetTag() string { return "ok" }

func (OK) MarshalJSON() ([]byte, error) { return []byte(`"ok"`), nil }

func init() {
	Register(func() Value { return new(OK) })
}
