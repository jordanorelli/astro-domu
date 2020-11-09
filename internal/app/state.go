package app

import "github.com/jordanorelli/astro-domu/internal/wire"

type state struct {
	playerName string
	room       *wire.Room
}
