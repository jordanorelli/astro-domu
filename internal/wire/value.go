package wire

type Value interface {
	NetTag() string
}

var registry = make(map[string]func() Value)

func Register(f func() Value) {
	v := f()
	t := v.NetTag()
	if _, exists := registry[t]; exists {
		panic("cannot register type: a type already exists with tag " + t)
	}
	registry[t] = f
}

func New(name string) Value {
	f, ok := registry[name]
	if !ok {
		return nil
	}
	return f()
}
