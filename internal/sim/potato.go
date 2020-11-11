package sim

import (
	"time"
)

type stageOfGrowth int

const (
	unborn stageOfGrowth = iota
	planted
	sapling
	unripe
	ripe
	overripe
	rotting
	dead
)

type potato struct {
	planted time.Duration
	stage   stageOfGrowth
}

func (p *potato) update(e *entity, dt time.Duration) {
	if p.stage > unborn {
		p.planted += dt
	}
}

type percent int

func (p *potato) progress() percent {
	if p.stage <= unborn {
		return 0
	}
	return 100
}
