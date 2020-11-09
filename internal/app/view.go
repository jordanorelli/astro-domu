package app

import (
	"github.com/gdamore/tcell/v2"
)

type view interface {
	handleEvent(tcell.Event) change
	draw(canvas, *state)
}

type focusable interface {
	setFocus(bool)
}

type change interface {
	exec(*UI)
}

type changeFn func(*UI)

func (f changeFn) exec(ui *UI) { f(ui) }
