package app

import (
	"github.com/jordanorelli/astro-domu/internal/sim"
	"github.com/jordanorelli/astro-domu/internal/wire"
)

type state struct {
	playerName string
	room       *wire.Room
	history    []sim.ChatMessage
}
