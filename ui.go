package main

import (
	"context"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/belt-mud/internal/exit"
	"github.com/jordanorelli/blammo"
)

// ui represents our terminal-based user interface
type ui struct {
	*blammo.Log
	screen tcell.Screen
	mode   uiMode
	client *client
}

func (ui *ui) run() {
	ui.client = &client{
		Log:  ui.Child("client"),
		host: "127.0.0.1",
		port: 12805,
	}
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		ui.Debug("canceling client context")
		cancel()
		time.Sleep(time.Second)
	}()

	go ui.client.run(ctx)

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

	ui.mode = &boxWalker{
		width:  10,
		height: 6,
	}
	ui.menu()
	ui.Debug("clearing screen")
	screen.Clear()
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

	for {
		e := ui.screen.PollEvent()
		if e == nil {
			break
		}

		switch v := e.(type) {
		case *tcell.EventKey:
			key := v.Key()
			if key == tcell.KeyCtrlC {
				return
			}
		}

		ui.mode.handleEvent(ui, e)
		ui.screen.Clear()
		ui.mode.draw(ui)
		ui.screen.Show()
	}
}
