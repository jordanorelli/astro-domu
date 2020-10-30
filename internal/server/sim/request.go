package sim

type Request struct {
	From  string
	Seq   int
	Wants Effect
}
