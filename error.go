package main

func errorResponse(re int, err error) response {
	var body struct {
		Message string `json:"message"`
	}
	body.Message = err.Error()
	return response{
		Re:   re,
		Type: "error",
		Body: body,
	}
}
