package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/belt-mud/internal/exit"
	"github.com/jordanorelli/blammo"
	"github.com/prometheus/common/log"
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
	log.Debug("sceen created")

	if err := screen.Init(); err != nil {
		exit.WithMessage(1, "unable to initialize screen: %v", err)
	}
	log.Debug("screen initialized")
	width, height := screen.Size()
	log.Debug("screen width: %d", width)
	log.Debug("screen height: %d", height)
	log.Debug("screen colors: %v", screen.Colors())
	log.Debug("screen has mouse: %v", screen.HasMouse())

	defer func() {
		log.Debug("finalizing screen")
		screen.Fini()
	}()

	ui.menu()
	log.Debug("clearing screen")
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
		ui.screen.Clear()
		ui.screen.SetContent(position.x, position.y, '+', nil, tcell.StyleDefault)
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
			log.Debug("screen saw key event: %v", v.Key())
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
			log.Debug("screen saw unhandled event of type %T", e)
		}
	}
}
