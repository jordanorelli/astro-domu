package sim

import "github.com/jordanorelli/astro-domu/internal/wire"

type Effect interface {
	wire.Value
	exec(*World, string)
}
