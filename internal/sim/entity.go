package sim

// entity is any entity that can be simulated.
type entity interface {
	// update is the standard tick function
	update(time.Duration)
}
