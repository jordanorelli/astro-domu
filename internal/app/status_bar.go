package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
)

type statusBar struct {
	inFocus bool
	ui      *UI
}

func (s *statusBar) handleEvent(ui *UI, e tcell.Event) bool { return false }

func (s *statusBar) draw(b *buffer) {
	if s.ui == nil {
		return
	}
	line := fmt.Sprintf("shown: %08d cleared: %08d messages: %08d", s.ui.showCount, s.ui.clearCount, s.ui.msgCount)
	b.writeString(line, math.Vec{0, 0}, tcell.StyleDefault)
}

func (s *statusBar) setFocus(enabled bool) { s.inFocus = enabled }
