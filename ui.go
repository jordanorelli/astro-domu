package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/belt-mud/internal/exit"
	"github.com/jordanorelli/blammo"
)

// ui represents our terminal-based user interface
type ui struct {
	*blammo.Log
	screen tcell.Screen
}

func (ui *ui) run() {
	screen, err := tcell.NewScreen()
	if err != nil {
		exit.WithMessage(1, "unable to create a screen: %v", err)
	}
	ui.screen = screen
	ui.Debug("sceen created")

	if err := screen.Init(); err != nil {
		exit.WithMessage(1, "unable to initialize screen: %v", err)
	}
	ui.Debug("screen initialized")
	width, height := screen.Size()
	ui.Debug("screen width: %d", width)
	ui.Debug("screen height: %d", height)
	ui.Debug("screen colors: %v", screen.Colors())
	ui.Debug("screen has mouse: %v", screen.HasMouse())

	defer func() {
		ui.Debug("finalizing screen")
		screen.Fini()
	}()

	ui.menu()
	ui.Debug("clearing screen")
	screen.Clear()

	time.Sleep(1 * time.Second)

}

// writeString writes a string in the given style from left to right beginning
// at the location (x, y). Writing of the screen just fails silently so don't
// do that.
func (ui *ui) writeString(x, y int, s string, style tcell.Style) {
	width, height := ui.screen.Size()
	if y > height {
		return
	}
	for i, r := range []rune(s) {
		x := x + i
		if x > width {
			return
		}
		ui.screen.SetContent(x, y, r, nil, style)
	}
}

func (ui *ui) menu() {
	ui.screen.Clear()
	_, height := ui.screen.Size()
	ui.writeString(0, height-1, "fart", tcell.StyleDefault)
	ui.screen.Sync()

	type point struct{ x, y int }
	position := point{10, 10}

	redraw := func() {
		if position.x < 2 {
			position.x = 2
		}
		if position.x > 11 {
			position.x = 11
		}
		if position.y < 2 {
			position.y = 2
		}
		if position.y > 10 {
			position.y = 10
		}
		ui.screen.Clear()
		ui.writeString(1, 1, `┌──────────┐`, tcell.StyleDefault)
		ui.writeString(1, 2, `│··········│`, tcell.StyleDefault)
		ui.writeString(1, 3, `│··········│`, tcell.StyleDefault)
		ui.writeString(1, 4, `│··········│`, tcell.StyleDefault)
		ui.writeString(1, 5, `│··········│`, tcell.StyleDefault)
		ui.writeString(1, 6, `│··········│`, tcell.StyleDefault)
		ui.writeString(1, 7, `│··········│`, tcell.StyleDefault)
		ui.writeString(1, 8, `│··········│`, tcell.StyleDefault)
		ui.writeString(1, 9, `│··········│`, tcell.StyleDefault)
		ui.writeString(1, 10, `│··········│`, tcell.StyleDefault)
		ui.writeString(1, 11, `└──────────┘`, tcell.StyleDefault)
		ui.screen.SetContent(position.x, position.y, '+', nil, tcell.StyleDefault)
		ui.writeString(0, 12, fmt.Sprintf(" (%02d, %02d)", position.x, position.y), tcell.StyleDefault)
		ui.screen.Show()
	}
	redraw()

	for {
		e := ui.screen.PollEvent()
		if e == nil {
			break
		}
		switch v := e.(type) {
		case *tcell.EventKey:
			ui.Debug("screen saw key event: %v", v.Key())
			key := v.Key()
			if key == tcell.KeyCtrlC {
				return
			}
			if key == tcell.KeyRune {
				switch v.Rune() {
				case 'w':
					position.y--
					redraw()
				case 'a':
					position.x--
					redraw()
				case 's':
					position.y++
					redraw()
				case 'd':
					position.x++
					redraw()
				}
			}
		default:
			ui.Debug("screen saw unhandled event of type %T", e)
		}
	}
}
