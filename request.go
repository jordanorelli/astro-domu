package main

import "encoding/json"

// requestBody is a container structure that contains the data of a user
// request. This is what clients serialize directly.
type requestBody struct {
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

type request struct {
	body requestBody
}
