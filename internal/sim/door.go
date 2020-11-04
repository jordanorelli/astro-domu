package sim

import (
	"time"

	"github.com/jordanorelli/blammo"
)

type door struct {
	*blammo.Log
}

func (d *door) update(time.Duration) {}

func (d *door) onStartOverlap(e *entity) {
	d.Info("start overlap: %v", e)
}

func (d *door) onOverlap(e *entity) {
	d.Info("overlap: %v", e)
}

func (d *door) onStopOverlap(e *entity) {
	d.Info("stop overlap: %v", e)
}

/*


┌──────────┐
│··········│
│··········│
│··········│        ┌──────┐
│··········│        │      │
│··········│        │      │
│·····d····◇░░░░░░░░◇      │
│··········│        │      │
│··········│        │      │
│··········│        └──────┘
│··········◇
└──────────┘

*/
