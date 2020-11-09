package app

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jordanorelli/astro-domu/internal/exit"
	"github.com/jordanorelli/astro-domu/internal/math"
	"github.com/jordanorelli/astro-domu/internal/sim"
	"github.com/jordanorelli/astro-domu/internal/wire"
	"github.com/jordanorelli/blammo"
)

type UI struct {
	*blammo.Log
	screen        tcell.Screen
	client        *wire.Client
	notifications <-chan wire.Response

	state state
	root  *node
}

func (ui *UI) Run() {
	ui.setupTerminal()
	defer ui.clearTerminal()

	ui.root = mainMenu

	input := make(chan tcell.Event)
	go ui.pollInput(input)

	tick := time.Tick(time.Second / time.Duration(30))

	width, height := ui.screen.Size()
	b := newBuffer(width, height)

	// b, prev := newBuffer(width, height), newBuffer(width, height)

	for {
		notify := ui.notifications

		select {
		case e := <-input:
			switch v := e.(type) {
			case *tcell.EventKey:
				key := v.Key()
				if key == tcell.KeyCtrlC {
					ui.Info("saw ctrl+c keyboard input, shutting down")
					return
				}
			}
			ui.Info("input event: %v", e)
			wants := ui.root.handleEvent(e)
			if wants != nil {
				wants.exec(ui)
			}
		case <-tick:
			b.clear()
			ui.root.draw(b, &ui.state)
			b.blit(ui.screen, math.Vec{0, 0})
			ui.screen.Show()
		case n := <-notify:
			ui.Info("notification: %v", n)
			ui.handleNotification(n.Body)
		}
	}
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

func (ui *UI) handleNotification(v wire.Value) bool {
	switch n := v.(type) {

	case *wire.Entity:
		ui.state.room.Entities[n.ID] = *n
		return true

	case *wire.Frame:
		if ui.state.room == nil {
			ui.state.room = new(wire.Room)
		}

		ui.state.room.Name = n.RoomName
		ui.state.room.Rect = n.RoomSize
		ui.state.room.Entities = n.Entities
		return true

	case *wire.Delta:
		if n.RoomSize != nil {
			ui.state.room.Rect = *n.RoomSize
		}

		if len(n.Entities) > 0 {
			for id, e := range n.Entities {
				if e != nil {
					ui.state.room.Entities[id] = *e
				} else {
					delete(ui.state.room.Entities, id)
				}
			}
		}
		return true

	case *sim.ChatMessage:
		// ui.chatView.history = append(ui.chatView.history, *n)
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

func (ui *UI) pollInput(c chan tcell.Event) {
	defer close(c)

	for {
		e := ui.screen.PollEvent()
		if e == nil {
			ui.Info("run loop sees nil event, breaking out")
			return
		}
		c <- e
	}
}

var mainMenu = &node{
	view: &menuList{
		choices: []menuItem{
			menuItem{
				name: "join",
				onSelect: changeFn(func(ui *UI) {
					ui.root = joinForm
				}),
			},
			menuItem{
				name: "exit",
				onSelect: changeFn(func(ui *UI) {
					panic("this is bad programming")
				}),
			},
		},
	},
}

var joinForm = &node{
	view: &form{
		fields: []textField{
			textField{
				label: "What is your name?",
				textInput: textInput{
					prompt: "> ",
				},
			},
		},
	},
}
