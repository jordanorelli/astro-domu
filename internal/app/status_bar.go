package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/math"
)

type statusBar struct {
	inFocus    bool
	msgCount   int
	showCount  int
	clearCount int
}

func (s *statusBar) handleEvent(ui *UI, e tcell.Event) bool { return false }

func (s *statusBar) draw(img canvas, st *state) {
	line := fmt.Sprintf("shown: %08d cleared: %08d messages: %08d", s.showCount, s.clearCount, s.msgCount)
	writeString(img, line, math.Vec{0, 0}, tcell.StyleDefault)
}

func (s *statusBar) setFocus(enabled bool) { s.inFocus = enabled }
