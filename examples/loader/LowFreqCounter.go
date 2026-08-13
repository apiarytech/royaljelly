package main

import (
	"fmt"
	"time"

	. "github.com/apiarytech/royaljelly/iec"
)

// LowFreqCounter encapsulates the state and logic for the low-frequency program.
type LowFreqCounter struct {
	count LINT
}

// Logic is the method that will be executed by the scheduler.
func (p *LowFreqCounter) Logic(now time.Time) {
	p.count++
	fmt.Printf("[%s] Low Frequency Task running.  Count: %d\n", now.Format("15:04:05.000"), p.count)
}
