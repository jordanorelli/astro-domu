package wire

type Response struct {
	Re   int         `json:"re,omitempty"`
	Type Tag         `json:"type"`
	Body interface{} `json:"body"`
}

func NewResponse(re int, v Value) Response {
	return Response{
		Re:   re,
		Type: v.NetTag(),
		Body: v,
	}
}

func ErrorResponse(re int, t string, args ...interface{}) Response {
	return NewResponse(re, Errorf(t, args...))
}
