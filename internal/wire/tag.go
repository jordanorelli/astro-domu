package wire

import (
	"encoding/json"
	"fmt"
)

type Tag uint

const (
	T_None Tag = iota
	T_Error
	T_OK
	T_Client_Move
)

func (t Tag) String() string {
	switch t {
	case T_Error:
		return "error"
	case T_OK:
		return "ok"
	case T_Client_Move:
		return "self/move"
	default:
		panic("unknown type tag")
	}
}

func (t *Tag) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}

	switch name {
	case "error":
		*t = T_Error
		return nil
	case "ok":
		*t = T_OK
		return nil
	case "self/move":
		*t = T_Client_Move
		return nil
	default:
		return fmt.Errorf("unknown type tag: %q", name)
	}
}

func (t Tag) MarshalJSON() ([]byte, error) { return json.Marshal(t.String()) }
