package app

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/exit"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type UI struct {
	*blammo.Log
	PlayerName string
	screen     tcell.Screen
	view       view
	room       *wire.Room
	client     *wire.Client
}

func (ui *UI) Run() {
	ui.setupTerminal()
	defer ui.clearTerminal()

	if err := ui.connect(); err != nil {
		return
	}

	res, err := ui.client.Send(wire.Login{Name: ui.PlayerName})
	if err != nil {
		ui.Error("login error: %v", err)
		return
	}

	welcome, ok := res.Body.(*wire.Welcome)
	if !ok {
		ui.Error("unexpected initial message of type %t", res.Body)
		return
	}
	ui.Info("welcome: %v", welcome)
	meta := welcome.Players[ui.PlayerName]
	room := welcome.Rooms[meta.Room]
	ui.room = &room
	ui.view = &gameView{
		Log:  ui.Child("game-view"),
		room: &room,
		me:   room.Entities[meta.Avatar],
	}

	ui.Info("running ui")
	if ui.handleUserInput() {
		ui.Info("user requested close")
		ui.Info("closing client")
		ui.client.Close()
		ui.Info("client closed")
		ui.Info("finalizing screen")
	}
	ui.Info("run loop done, shutting down")
}

func (ui *UI) connect() error {
	ui.client = &wire.Client{
		Log:  ui.Child("client"),
		Host: "127.0.0.1",
		Port: 12805,
	}

	c, err := ui.client.Dial()
	if err != nil {
		return fmt.Errorf("unable to dial server: %v", err)
	}
	go ui.handleNotifications(c)
	return nil
}

func (ui *UI) setupTerminal() {
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
}

func (ui *UI) clearTerminal() {
	ui.screen.Clear()
	ui.screen.Fini()
}

func (ui *UI) handleNotifications(c <-chan wire.Response) {
	for n := range c {
		if ui.handleNotification(n.Body) {
			ui.view.draw(ui)
			ui.screen.Show()
		}
	}
	ui.Info("notifications channel is closed so we must be done")
	ui.Info("clearing and finalizing screen from notifications goroutine")
	ui.screen.Clear()
	ui.screen.Fini()
	ui.Info("screen finalized")
}

func (ui *UI) handleNotification(v wire.Value) bool {
	switch n := v.(type) {

	case *wire.Entity:
		ui.room.Entities[n.ID] = n
		return true

	default:
		ui.Info("ignoring notification: %v", n)
		return false
	}
}

// writeString writes a string in the given style from left to right beginning
// at the location (x, y). Writing of the screen just fails silently so don't
// do that.
func (ui *UI) writeString(x, y int, s string, style tcell.Style) {
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

func (ui *UI) handleUserInput() bool {
	ui.screen.Clear()
	ui.view.draw(ui)

	for {
		e := ui.screen.PollEvent()
		if e == nil {
			ui.Info("run loop sees nil event, breaking out")
			// someone else shut us down, so return false
			return false
		}

		switch v := e.(type) {
		case *tcell.EventKey:
			key := v.Key()
			if key == tcell.KeyCtrlC {
				ui.Info("saw ctrl+c keyboard input, shutting down")
				// we want to shut things down
				return true
			}
		}

		ui.view.handleEvent(ui, e)
		ui.screen.Clear()
		ui.view.draw(ui)
		ui.screen.Show()
	}
}
