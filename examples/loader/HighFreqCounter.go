package main

import (
	"fmt"
	"time"

	"github.com/apiarytech/royaljelly/core"
)

// HighFreqCounter encapsulates the state and logic for the high-frequency program.
type HighFreqCounter struct {
	count core.LINT
}

// Logic is the method that will be executed by the scheduler.
func (p *HighFreqCounter) Logic(now time.Time) {
	p.count++
	fmt.Printf("[%s] High Frequency Task running. Count: %d\n", now.Format("15:04:05.000"), p.count)
}
