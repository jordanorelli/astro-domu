package sim

import (
	"time"
)

type Entity struct {
	ID       int    `json:"id"`
	Position [2]int `json:"pos"`
	Glyph    rune   `json:"glyph"`
	behavior
}

type behavior interface {
	// update is the standard tick function
	update(time.Duration)
}
