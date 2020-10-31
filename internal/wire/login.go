package wire

type Login struct {
	Name string `json:"name"`
}

func (Login) NetTag() string { return "login" }

func init() {
	Register(func() Value { return new(Login) })
}
