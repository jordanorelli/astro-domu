package wire

import (
	"errors"
	"fmt"
)

type Error struct {
	Message string `json:"message"`
	parent  error
}

func (e Error) Error() string  { return e.Message }
func (e Error) NetTag() string { return "error" }
func (e Error) Unwrap() error  { return e.parent }

func Errorf(t string, args ...interface{}) Error {
	err := fmt.Errorf(t, args...)
	return Error{
		Message: err.Error(),
		parent:  errors.Unwrap(err),
	}
}

func init() {
	Register(func() Value { return new(Error) })
}
