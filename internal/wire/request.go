package wire

type Request struct {
	Seq  int         `json:"seq"`
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}
