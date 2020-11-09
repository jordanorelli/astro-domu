package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/sim"
	"github.com/jordanorelli/blammo"
)

type chatView struct {
	*blammo.Log
	composing string
	inFocus   bool
	history   []sim.ChatMessage
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
			c.composing = ""
			return changeFn(func(ui *UI) {
				// ugh lol
				go ui.client.Send(sim.SendChatMessage{Text: c.composing})
			})
		}

		if key == tcell.KeyRune {
			c.composing = fmt.Sprintf("%s%c", c.composing, t.Rune())
			c.Info("composing: %v", c.composing)
		}

	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return nil
}

func (c *chatView) draw(img canvas, st *state) {
	bounds := img.bounds()
	chatHeight := bounds.Height - 1
	for i := 0; i < math.Min(chatHeight, len(c.history)); i++ {
		msg := c.history[len(c.history)-1-i]
		s := fmt.Sprintf("%12s: %s", msg.From, msg.Text)
		writeString(img, s, math.Vec{0, bounds.Height - 2 - i}, tcell.StyleDefault)
	}

	writeString(img, c.composing, math.Vec{0, bounds.Height - 1}, tcell.StyleDefault)

	if c.inFocus {
		cursor := tile{r: tcell.RuneBlock, style: tcell.StyleDefault.Blink(true)}
		img.setTile(len([]rune(c.composing)), bounds.Height-1, cursor)
	}
}

func (c *chatView) setFocus(yes bool) { c.inFocus = yes }
