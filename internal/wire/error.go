package wire

import (
	"fmt"
)

type Error struct {
	val error
}

func (e Error) Error() string  { return e.val.Error() }
func (e Error) NetTag() string { return "error" }
func (e Error) Unwrap() error  { return e.val }

func Errorf(t string, args ...interface{}) Error {
	return Error{val: fmt.Errorf(t, args...)}
}

func init() {
	Register(func() Value { return new(Error) })
}
