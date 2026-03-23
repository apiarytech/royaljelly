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

package royaljelly

import (
	"time"
)

/********************************************************************************/
/* This file implements the standard timer function blocks as per IEC 61131-3.  */
/* Each timer (TP, TON, TOF) is a separate function block with its own state.   */
/********************************************************************************/

// TP is a Pulse timer function block.
// It generates a pulse of a specified duration.
type TP struct {
	// Inputs
	IN BOOL // Trigger
	PT TIME // Preset Time (pulse duration)

	// Outputs
	Q  BOOL // Output pulse
	ET TIME // Elapsed Time

	// Internal state
	startTime time.Time
	timing    BOOL
	re        R_TRIG // Rising-edge detector
}

// INIT initializes the TP function block to its default state.
func (t *TP) INIT() {
	t.IN = false
	t.PT = 0
	t.Q = false
	t.ET = 0
	t.startTime = time.Time{}
	t.timing = false
	t.re.INIT()
}

// Execute runs the logic for the TP function block for the current scan cycle.
// The `now` parameter represents the current time of the scan.
func (t *TP) Execute(now time.Time) {
	t.re.CLK = t.IN
	t.re.R_TRIG()

	// On a rising edge of IN, start timing.
	if t.re.Q {
		t.timing = true
		t.startTime = now
	}

	if t.timing {
		elapsed := now.Sub(t.startTime)
		if elapsed < time.Duration(t.PT) {
			t.Q = true
			t.ET = TIME(elapsed)
		} else {
			// Pulse duration is over.
			t.Q = false
			t.ET = t.PT // Clamp ET to PT
			// The pulse is finished. Stop timing.
			t.timing = false
		}
	} else if !t.IN {
		// Not timing, reset outputs.
		t.Q = false
		t.ET = 0
	}
}

// TON is an On-Delay timer function block.
// It delays the setting of its output Q after its input IN becomes true.
type TON struct {
	// Inputs
	IN BOOL // Input
	PT TIME // Preset Time (delay duration)

	// Outputs
	Q  BOOL // Output
	ET TIME // Elapsed Time

	// Internal state
	startTime time.Time
}

// INIT initializes the TON function block to its default state.
func (t *TON) INIT() {
	t.IN = false
	t.PT = 0
	t.Q = false
	t.ET = 0
	t.startTime = time.Time{}
}

// Execute runs the logic for the TON function block for the current scan cycle.
// The `now` parameter represents the current time of the scan.
func (t *TON) Execute(now time.Time) {
	if !t.IN {
		// If input is false, reset everything.
		t.Q = false
		t.ET = 0
		t.startTime = time.Time{} // A zero time indicates the timer is not running.
		return
	}

	// If input is true and we haven't started timing yet.
	if t.startTime.IsZero() {
		t.startTime = now
	}

	elapsed := now.Sub(t.startTime)
	if elapsed < time.Duration(t.PT) {
		t.ET = TIME(elapsed)
		t.Q = false
	} else {
		// Timer has reached its preset time.
		t.ET = t.PT
		t.Q = true
	}
}

// TOF is an Off-Delay timer function block.
// It delays the resetting of its output Q after its input IN becomes false.
type TOF struct {
	// Inputs
	IN BOOL // Input
	PT TIME // Preset Time (delay duration)

	// Outputs
	Q  BOOL // Output
	ET TIME // Elapsed Time

	// Internal state
	startTime time.Time
	mem       BOOL // Internal memory of the previous state of IN
}

// INIT initializes the TOF function block to its default state.
func (t *TOF) INIT() {
	t.IN = false
	t.PT = 0
	t.Q = false
	t.ET = 0
	t.startTime = time.Time{}
	t.mem = false
}

// Execute runs the logic for the TOF function block for the current scan cycle.
// The `now` parameter represents the current time of the scan.
func (t *TOF) Execute(now time.Time) {
	// Detect the falling edge of IN
	fallingEdge := !t.IN && t.mem

	if t.IN {
		// If input is true, output is true and timer is reset.
		t.Q = true
		t.ET = 0
		t.startTime = time.Time{}
	} else if fallingEdge {
		// On the falling edge, start the timer.
		t.startTime = now
		t.Q = true // Q remains true as the off-delay starts
		t.ET = 0
	} else if !t.startTime.IsZero() {
		// This block executes when IN is false and the timer is running.
		elapsed := now.Sub(t.startTime)
		if elapsed < time.Duration(t.PT) {
			t.Q = true
			t.ET = TIME(elapsed)
		} else {
			t.Q = false
			t.ET = t.PT // Clamp ET to PT
		}
	} else {
		// If not IN, not a falling edge, and timer not running, then Q must be false.
		t.Q = false
		t.ET = 0
	}

	// Store the current state of IN for the next execution cycle.
	t.mem = t.IN
}
