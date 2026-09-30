/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - GPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package timers_test

import (
	"fmt"
	"time"

	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
)

// start is a fixed scan time so the examples are deterministic. In a program,
// pass the scan time that the scheduler gives to Program.Logic.
var start = time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)

func at(ms int) time.Time { return start.Add(time.Duration(ms) * time.Millisecond) }

func ExampleTON() {
	// Q turns on once IN has been true for PT.
	ton := timers.TON{PT: iec.TIME(100 * time.Millisecond)}
	ton.IN = true
	for _, ms := range []int{0, 50, 100} {
		ton.Execute(at(ms))
		fmt.Printf("t=%dms IN=%v Q=%v ET=%v\n", ms, ton.IN, ton.Q, time.Duration(ton.ET))
	}
	ton.IN = false
	ton.Execute(at(150))
	fmt.Printf("t=150ms IN=%v Q=%v ET=%v\n", ton.IN, ton.Q, time.Duration(ton.ET))
	// Output:
	// t=0ms IN=true Q=false ET=0s
	// t=50ms IN=true Q=false ET=50ms
	// t=100ms IN=true Q=true ET=100ms
	// t=150ms IN=false Q=false ET=0s
}

func ExampleTOF() {
	// Q follows IN on, and stays on for PT after IN turns off.
	tof := timers.TOF{PT: iec.TIME(100 * time.Millisecond)}
	steps := []struct {
		ms int
		in iec.BOOL
	}{{0, true}, {10, false}, {60, false}, {110, false}}
	for _, s := range steps {
		tof.IN = s.in
		tof.Execute(at(s.ms))
		fmt.Printf("t=%dms IN=%v Q=%v\n", s.ms, tof.IN, tof.Q)
	}
	// Output:
	// t=0ms IN=true Q=true
	// t=10ms IN=false Q=true
	// t=60ms IN=false Q=true
	// t=110ms IN=false Q=false
}

func ExampleTP() {
	// A rising edge on IN produces a pulse of exactly PT, however long IN stays on.
	tp := timers.TP{PT: iec.TIME(100 * time.Millisecond)}
	steps := []struct {
		ms int
		in iec.BOOL
	}{{0, true}, {50, true}, {100, true}, {150, false}}
	for _, s := range steps {
		tp.IN = s.in
		tp.Execute(at(s.ms))
		fmt.Printf("t=%dms IN=%v Q=%v\n", s.ms, tp.IN, tp.Q)
	}
	// Output:
	// t=0ms IN=true Q=true
	// t=50ms IN=true Q=true
	// t=100ms IN=true Q=false
	// t=150ms IN=false Q=false
}
