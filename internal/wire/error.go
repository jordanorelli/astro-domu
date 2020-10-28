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

func (e *Error) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*e = Errorf(s)
	return nil
}

func Errorf(t string, args ...interface{}) Error {
	return Error{val: fmt.Errorf(t, args...)}
}

func init() {
	Register(func() Value { return new(Error) })
}
