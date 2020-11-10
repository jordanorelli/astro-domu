package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
)

type inventory struct {
	items []item
}

type item struct {
	name string
}

type inventoryView struct {
}

func (v *inventoryView) handleEvent(e tcell.Event) change {
	return nil
}

func (v *inventoryView) draw(img canvas, st *state) {
	writeString(img, "Inventory", math.Vec{0, 0}, tcell.StyleDefault)
}

type openInventory struct{}

func (openInventory) exec(ui *UI) {
	if ui.root == inGameView {
		ui.state.detail = &inventoryView{}
		inGameView.focus(1)
	}
}
