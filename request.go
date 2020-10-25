package main

import "encoding/json"

type request struct {
	// client-side sequence number. Clients are expected to send requests with
	// a monotonically-incrementing sequence number.
	Seq int `json:"seq"`

	// command represents the command to be executed on the server. Command
	// names are globally unique.
	Command string `json:"cmd"`

	// args are the arguments passed to the command for parameterized commands
	// executed on the server.
	Args map[string]json.RawMessage `json:"args"`
}
