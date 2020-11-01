package wire

type Player struct {
	Name   string `json:"name"`
	Room   string `json:"room"`
	Avatar int    `json:"avatar"`
}
