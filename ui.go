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
}

func (ui *ui) run() {
	screen, err := tcell.NewScreen()
	if err != nil {
		exit.WithMessage(1, "unable to create a screen: %v", err)
	}
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

	log.Debug("clearing screen")
	for {
		e := screen.PollEvent()
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
		default:
			log.Debug("screen saw unhandled event of type %T", e)
		}
	}
	screen.Clear()

	time.Sleep(1 * time.Second)

}
