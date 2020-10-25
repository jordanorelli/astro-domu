package main

type response struct {
	Re   int         `json:"re"`
	Type string      `json:"type"`
	Body interface{} `json:"body,omitempty"`
}

func ok(re int) response {
	return response{
		Re:   re,
		Type: "ok",
	}
}
