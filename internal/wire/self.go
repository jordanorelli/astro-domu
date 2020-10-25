package wire

type Self_Move struct {
	Delta bool `json:"delta"`
	X     int  `json:"x"`
	Y     int  `json:"y"`
}

func (Self_Move) NetTag() Tag { return T_Client_Move }
