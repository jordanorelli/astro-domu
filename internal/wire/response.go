package wire

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Re   int
	Body Value
}

func (r Response) MarshalJSON() ([]byte, error) {
	return json.Marshal([3]interface{}{r.Re, r.Body.NetTag(), r.Body})
}

func (r *Response) UnmarshalJSON(b []byte) error {
	var parts [3]json.RawMessage
	if err := json.Unmarshal(b, &parts); err != nil {
		return err
	}
	if err := json.Unmarshal(parts[0], &r.Re); err != nil {
		return err
	}

	var tag string
	if err := json.Unmarshal(parts[1], &tag); err != nil {
		return err
	}

	f, ok := registry[tag]
	if !ok {
		return fmt.Errorf("unknown tag: %q", tag)
	}
	v := f()
	if err := json.Unmarshal(parts[2], v); err != nil {
		return err
	}
	r.Body = v
	return nil
}

func ErrorResponse(re int, t string, args ...interface{}) Response {
	return Response{re, Errorf(t, args...)}
}
