package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/sim"
	"github.com/jordanorelli/astro-domu/internal/wire"
)

type state struct {
	playerName string
	room       *wire.Room
	history    []sim.ChatMessage
	detail     view
}

// detailView is a proxy view. Wherever you put a detailView, you get the
// current view assigned to the detail in the room state
type detailView struct {
	showing view
}

func (v *detailView) handleEvent(e tcell.Event) change {
	if v.showing == nil {
		return nil
	}
	return v.showing.handleEvent(e)
}

func (v *detailView) draw(img canvas, st *state) {
	v.showing = st.detail
	if v.showing == nil {
		return
	}
	v.showing.draw(img, st)
}
