package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type inventory struct {
	items []item
}

func (inv *inventory) removeItem(id int) bool {
	for i, item := range inv.items {
		if item.ID == id {
			inv.items = append(inv.items[:i], inv.items[i+1:]...)
			return true
		}
	}
	return false
}

type item struct {
	ID   int
	name string
}

type inventoryView struct {
	*blammo.Log
	highlight int
	avatar    *wire.Entity // probably but not necessarily the player avatar
	*inventory
	keyHandler func(*tcell.EventKey) change
}

func (v *inventoryView) handleEvent(e tcell.Event) change {
	if v.keyHandler == nil {
		v.keyHandler = v.selectHandler
	}

	switch t := e.(type) {
	case *tcell.EventKey:
		return v.keyHandler(t)
	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return nil
}

func (v *inventoryView) selectHandler(e *tcell.EventKey) change {
	key := e.Key()

	if key == tcell.KeyRune {
		switch e.Rune() {
		case 'p':
			v.keyHandler = v.putdownHandler
			return nil
		}
	}

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

	return nil
}

func (v *inventoryView) putdownHandler(e *tcell.EventKey) change {
	key := e.Key()

	if key == tcell.KeyRune {

		v.Info("key %c pressed while holding %d at %v", e.Rune(), v.items[v.highlight].ID, v.avatar.Position)

		switch e.Rune() {
		case 'w':
			v.keyHandler = v.selectHandler
			return putdown{
				ID:       v.items[v.highlight].ID,
				Location: v.avatar.Position.Add(math.Up),
			}
		case 'a':
			v.keyHandler = v.selectHandler
			return putdown{
				ID:       v.items[v.highlight].ID,
				Location: v.avatar.Position.Add(math.Left),
			}
		case 's':
			v.keyHandler = v.selectHandler
			return putdown{
				ID:       v.items[v.highlight].ID,
				Location: v.avatar.Position.Add(math.Down),
			}
		case 'd':
			v.keyHandler = v.selectHandler
			return putdown{
				ID:       v.items[v.highlight].ID,
				Location: v.avatar.Position.Add(math.Right),
			}
		}
	}

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

	return nil
}

func (v *inventoryView) draw(img canvas, st *state) {
	v.avatar = st.avatar
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
		ui.state.detail = &inventoryView{
			Log:       ui.Log.Child("inventory"),
			avatar:    ui.state.avatar,
			inventory: &ui.state.inventory,
		}
		inGameView.focus(1)
	}
}
