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

func (c *chatView) handleEvent(ui *UI, e tcell.Event) bool {
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
			// ugh lol
			go ui.client.Send(sim.SendChatMessage{Text: c.composing})
			c.composing = ""
			break
		}

		if key == tcell.KeyRune {
			c.composing = fmt.Sprintf("%s%c", c.composing, t.Rune())
			c.Info("composing: %v", c.composing)
		}

	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return false
}

func (c *chatView) draw(b *buffer) {
	chatHeight := b.height - 1
	for i := 0; i < math.Min(chatHeight, len(c.history)); i++ {
		msg := c.history[len(c.history)-1-i]
		s := fmt.Sprintf("%12s: %s", msg.From, msg.Text)
		b.writeString(s, math.Vec{0, b.height - 2 - i}, tcell.StyleDefault)
	}

	b.writeString(c.composing, math.Vec{0, b.height - 1}, tcell.StyleDefault)

	if c.inFocus {
		cursor := tile{r: tcell.RuneBlock, style: tcell.StyleDefault.Blink(true)}
		b.set(len([]rune(c.composing)), b.height-1, cursor)
	}
}

func (c *chatView) setFocus(yes bool) { c.inFocus = yes }
