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

/********************************************************************************/
/* This file implements the standard edge detection function blocks R_TRIG      */
/* and F_TRIG as per IEC 61131-3, Table 35.                                     */
/********************************************************************************/

// R_TRIG is a rising edge detection function block.
type R_TRIG struct {
	// Input
	CLK BOOL // Clock input

	// Output
	Q BOOL // Output pulse

	// Internal memory
	mem BOOL // Memory of the previous state of CLK
}

// F_TRIG is a falling edge detection function block.
type F_TRIG struct {
	// Input
	CLK BOOL // Clock input

	// Output
	Q BOOL // Output pulse

	// Internal memory
	mem BOOL // Memory of the previous state of CLK
}

// INIT initializes the R_TRIG function block.
func (fb *R_TRIG) INIT() {
	fb.CLK = false
	fb.Q = false
	fb.mem = false
}

// R_TRIG executes the rising edge detection logic.
func (fb *R_TRIG) R_TRIG() {
	fb.Q = fb.CLK && !fb.mem
	fb.mem = fb.CLK
}

// INIT initializes the F_TRIG function block.
func (fb *F_TRIG) INIT() {
	fb.CLK = false
	fb.Q = false
	fb.mem = false
}

// F_TRIG executes the falling edge detection logic.
func (fb *F_TRIG) F_TRIG() {
	fb.Q = !fb.CLK && fb.mem
	fb.mem = fb.CLK
}
