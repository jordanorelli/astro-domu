package wire

type Request struct {
	Seq  int         `json:"seq"`
	Type Tag         `json:"type"`
	Body interface{} `json:"body"`
}

func NewRequest(seq int, v Value) Request {
	return Request{
		Seq:  seq,
		Type: v.NetTag(),
		Body: v,
	}
}
