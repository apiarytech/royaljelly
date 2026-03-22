/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software

 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.
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
