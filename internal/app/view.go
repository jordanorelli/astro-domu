package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/wire"
)

type view interface {
	handleEvent(*UI, tcell.Event) bool
	notify(wire.Value)
	draw(*UI)
}
