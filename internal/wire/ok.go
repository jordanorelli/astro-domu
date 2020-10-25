package wire

type OK struct {}
func (OK) NetTag() Tag { return T_OK }
