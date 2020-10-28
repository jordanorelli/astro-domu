package wire

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	val error
}

func (e Error) Error() string  { return e.val.Error() }
func (e Error) NetTag() string { return "error" }
func (e Error) Unwrap() error  { return e.val }

func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Error())
}

func Errorf(t string, args ...interface{}) Error {
	return Error{val: fmt.Errorf(t, args...)}
}

func init() {
	Register(func() Value { return new(Error) })
}
