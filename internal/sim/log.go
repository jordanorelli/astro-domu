package sim

import (
	"os"

	"github.com/jordanorelli/blammo"
)

func defaultLog() *blammo.Log {
	stdout := blammo.NewLineWriter(os.Stdout)
	stderr := blammo.NewLineWriter(os.Stderr)

	options := []blammo.Option{
		blammo.DebugWriter(stdout),
		blammo.InfoWriter(stdout),
		blammo.ErrorWriter(stderr),
	}

	return blammo.NewLog("sim", options...)
}
