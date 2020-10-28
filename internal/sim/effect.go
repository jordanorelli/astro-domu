package sim

import "github.com/jordanorelli/astro-domu/internal/wire"

type Effect interface {
	wire.Value
	exec(*room, string, int) result
}

type result struct {
	reply    wire.Value
	announce wire.Value
}
