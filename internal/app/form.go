package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
)

type form struct {
	fields []textField
	active int
}

func (f *form) handleEvent(e tcell.Event) change {
	switch t := e.(type) {
	case *tcell.EventKey:
		key := t.Key()
		switch key {
		case tcell.KeyEnter:
			return login{name: f.fields[0].textInput.entered}
		}
	}
	return f.fields[0].handleEvent(e)
}

func (f *form) draw(b *buffer, _ *state) {
	for i, field := range f.fields {
		b.writeString(field.label, math.Vec{0, i * 2}, tcell.StyleDefault)
		b.writeString(field.prompt+field.entered, math.Vec{0, i*2 + 1}, tcell.StyleDefault)
	}
}

type textField struct {
	label string
	textInput
}
