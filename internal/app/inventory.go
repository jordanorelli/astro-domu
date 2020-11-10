package app

import (
	"fmt"

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
	if k, ok := e.(*tcell.EventKey); ok {
		if k.Key() == tcell.KeyESC {
			return changeFn(func(ui *UI) {
				if ui.root == inGameView {
					inGameView.focus(0)
				}
			})
		}
	}
	return nil
}

func (v *inventoryView) draw(img canvas, st *state) {
	writeString(img, "Inventory", math.Vec{0, 0}, tcell.StyleDefault)
	for i, item := range st.inventory.items {
		line := fmt.Sprintf("- %s", item.name)
		writeString(img, line, math.Vec{0, i + 2}, tcell.StyleDefault)
	}
}

type openInventory struct{}

func (openInventory) exec(ui *UI) {
	if ui.root == inGameView {
		ui.state.detail = &inventoryView{}
		inGameView.focus(1)
	}
}
