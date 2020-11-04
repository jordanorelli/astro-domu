package app

import (
	"github.com/gdamore/tcell/v2"
)

type view interface {
	handleEvent(*UI, tcell.Event) bool
	draw(*buffer)
	setFocus(bool)
}
