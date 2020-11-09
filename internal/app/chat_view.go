package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/sim"
)

type chatView struct {
	composing string
	inFocus   bool
}

func (c *chatView) handleEvent(e tcell.Event) change {
	switch t := e.(type) {
	case *tcell.EventKey:
		key := t.Key()

		if key == tcell.KeyBackspace || key == tcell.KeyBackspace2 {
			line := []rune(c.composing)
			if len(line) > 0 {
				line = line[:len(line)-1]
			}
			c.composing = string(line)
			break
		}

		if key == tcell.KeyEnter {
			msg := c.composing
			c.composing = ""
			return changeFn(func(ui *UI) {
				// ugh lol
				go ui.client.Send(sim.SendChatMessage{Text: msg})
			})
		}

		if key == tcell.KeyRune {
			c.composing = fmt.Sprintf("%s%c", c.composing, t.Rune())
		}

	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return nil
}

func (c *chatView) draw(img canvas, st *state) {
	bounds := img.bounds()
	style := tcell.StyleDefault.Background(tcell.NewRGBColor(32, 32, 32))
	fill(img, style)
	chatHeight := bounds.Height - 1

	for i := 0; i < math.Min(chatHeight, len(st.history)); i++ {
		msg := st.history[len(st.history)-1-i]
		nameText := fmt.Sprintf("%12s: ", msg.From)
		if msg.From == st.playerName {
			style := style.Foreground(tcell.NewRGBColor(255, 32, 32))
			writeString(img, nameText, math.Vec{0, bounds.Height - 2 - i}, style)
		} else {
			style := style.Foreground(tcell.NewRGBColor(32, 32, 255))
			writeString(img, nameText, math.Vec{0, bounds.Height - 2 - i}, style)
		}
		writeString(img, msg.Text, math.Vec{14, bounds.Height - 2 - i}, style)
	}

	writeString(img, c.composing, math.Vec{0, bounds.Height - 1}, style)

	if c.inFocus {
		cursor := tile{r: tcell.RuneBlock, style: style.Blink(true)}
		img.setTile(len([]rune(c.composing)), bounds.Height-1, cursor)
	}
}

func (c *chatView) setFocus(yes bool) { c.inFocus = yes }
