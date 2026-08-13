package main

import (
	"fmt"
	"time"

	"github.com/apiarytech/royaljelly/core"
)

var (
	lowFreqCount core.LINT = 0
)

// lowFreqLogicFunc is the logic for the program that runs less frequently.
func LowFreqLogicFunc(now time.Time) {
	lowFreqCount++
	fmt.Printf("[%s] Low Frequency Task running.  Count: %d\n", now.Format("15:04:05.000"), lowFreqCount)
}
