package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
)

type textInput struct {
	prompt  string
	entered string
}

func (t *textInput) handleEvent(e tcell.Event) change {
	switch v := e.(type) {
	case *tcell.EventKey:
		key := v.Key()

		if key == tcell.KeyBackspace || key == tcell.KeyBackspace2 {
			line := []rune(t.entered)
			if len(line) > 0 {
				line = line[:len(line)-1]
			}
			t.entered = string(line)
			break
		}

		if key == tcell.KeyRune {
			t.entered = fmt.Sprintf("%s%c", t.entered, v.Rune())
		}

	default:
		// ui.Debug("screen saw unhandled event of type %T", e)
	}
	return nil
}

func (t *textInput) draw(img canvas, _ *state) {
	writeString(img, t.prompt+t.entered, math.Vec{0, 0}, tcell.StyleDefault)
}

type textView string

func (textView) handleEvent(tcell.Event) change { return nil }

func (t textView) draw(img canvas, _ *state) {
	writeString(img, string(t), math.Vec{0, 0}, tcell.StyleDefault)
}
