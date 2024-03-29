package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/jordanorelli/astro-domu/internal/app"
	"github.com/jordanorelli/astro-domu/internal/exit"
	"github.com/jordanorelli/astro-domu/internal/sim"
	"github.com/jordanorelli/blammo"
)

func newLog(path string) *blammo.Log {
	f, err := os.OpenFile("./astro.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		exit.WithMessage(1, "unable to open log file %q for writing: %v", err)
	}

	w := blammo.NewLineWriter(f)

	options := []blammo.Option{
		blammo.DebugWriter(w),
		blammo.InfoWriter(w),
		blammo.ErrorWriter(w),
	}

	return blammo.NewLog("astro", options...)
}

func main() {
	if len(os.Args) < 2 {
		runClient()
	}

	switch os.Args[1] {
	case "client":
		runClient()
	case "server":
		s := sim.Server{}
		if err := s.Start("cdm.jordanorelli.com", 12805); err != nil {
			exit.WithMessage(1, "unable to start server: %v", err)
		}
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		s.Shutdown()

	default:
		exit.WithMessage(1, "supported options are [client|server]")
	}
}

func runClient() {
	log := newLog("./astro.log").Child("client")

	start := time.Now()
	log.Info("starting at: %v", start)
	defer func() {
		finished := time.Now()
		log.Info("finished at: %v", finished)
		log.Info("total play time: %v", finished.Sub(start))
	}()

	app := app.UI{
		Log: log.Child("ui"),
	}
	app.Run()
}
