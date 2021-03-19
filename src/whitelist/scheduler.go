package whitelist

import (
	"time"
)

type jobTicker struct {
	timer *time.Timer
}

func (t *jobTicker) updateTimer(diff time.Duration) {
	if t.timer == nil {
		t.timer = time.NewTimer(diff)
	} else {
		t.timer.Reset(diff)
	}
}
