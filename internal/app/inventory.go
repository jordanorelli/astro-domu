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
	highlight int
	*inventory
}

func (v *inventoryView) handleEvent(e tcell.Event) change {
	switch t := e.(type) {
	case *tcell.EventKey:
		key := t.Key()
		switch key {
		case tcell.KeyEnter:
		case tcell.KeyDown:
			if len(v.items) > 0 {
				v.highlight = (v.highlight + 1) % len(v.items)
			}

		case tcell.KeyUp:
			if len(v.items) > 0 {
				v.highlight = (v.highlight - 1 + len(v.items)) % len(v.items)
			}

		case tcell.KeyESC:
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
	v.inventory = &st.inventory
	writeString(img, "Inventory", math.Vec{0, 0}, tcell.StyleDefault)
	for i, item := range st.inventory.items {
		line := fmt.Sprintf("- %s", item.name)
		if i == v.highlight {
			line = fmt.Sprintf("+ %s", item.name)
		}
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
