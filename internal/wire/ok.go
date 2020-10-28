package wire

type OK struct{}

func (OK) NetTag() string { return "ok" }

func init() {
	Register(func() Value { return new(OK) })
}
