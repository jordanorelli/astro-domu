package main

import (
	"os"
	"time"

	"github.com/jordanorelli/belt-mud/internal/exit"
	"github.com/jordanorelli/blammo"
)

func newLog(path string) *blammo.Log {
	f, err := os.OpenFile("./belt.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		exit.WithMessage(1, "unable to open log file %q for writing: %v", err)
	}

	w := blammo.NewLineWriter(f)

	options := []blammo.Option{
		blammo.DebugWriter(w),
		blammo.InfoWriter(w),
		blammo.ErrorWriter(w),
	}

	return blammo.NewLog("belt", options...)
}

func main() {
	log := newLog("./belt.log")

	start := time.Now()
	log.Info("starting at: %v", start)
	defer func() {
		finished := time.Now()
		log.Info("finished at: %v", finished)
		log.Info("total play time: %v", finished.Sub(start))
	}()

	ui := ui{
		Log: log.Child("ui"),
	}
	ui.run()
}
