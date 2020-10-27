package sim

import "time"

// entity is any entity that can be simulated.
type entity interface {
	id() int

	// update is the standard tick function
	update(time.Duration)
}
