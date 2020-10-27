package wire

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	Seq  int
	Body Value
}

func (r Request) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{r.Seq, r.Body.NetTag(), r.Body})
}

func (r *Request) UnmarshalJSON(b []byte) error {
	var parts [3]json.RawMessage
	if err := json.Unmarshal(b, &parts); err != nil {
		return err
	}
	if err := json.Unmarshal(parts[0], &r.Seq); err != nil {
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
