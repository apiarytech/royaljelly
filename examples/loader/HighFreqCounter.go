package main

import (
	"fmt"
	"time"

	"github.com/apiarytech/royaljelly/core"
)

var (
	highFreqCount core.LINT = 0
)

// highFreqLogicFunc is the logic for the program that runs more frequently.
func HighFreqLogicFunc(now time.Time) {
	highFreqCount++
	fmt.Printf("[%s] High Frequency Task running. Count: %d\n", now.Format("15:04:05.000"), highFreqCount)
}
