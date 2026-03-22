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

// SR_FB is a Set-dominant bistable function block (Flip-Flop).
type SR_FB struct {
	//INPUT
	EN  BOOL //enable
	ENO BOOL //enable output

	S1 BOOL
	R  BOOL

	//OUTPUT
	Q1 BOOL
}

// RS_FB is a Reset-dominant bistable function block (Flip-Flop).
type RS_FB struct {
	//INPUT
	EN  BOOL //enable
	ENO BOOL //enable output

	S  BOOL
	R1 BOOL

	//OUTPUT
	Q1 BOOL
}

// Init initializes the SR function block.
func (SR *SR_FB) INIT() {
	SR.EN = true
	SR.ENO = true
	SR.S1 = false
	SR.R = false
	SR.Q1 = false
}

// SR executes the Set-dominant logic.
func (SR *SR_FB) SR() {
	if !SR.EN {
		SR.ENO = false
		return
	}
	SR.ENO = true
	SR.Q1 = SR.S1 || (!SR.R && SR.Q1)
}

// Init initializes the RS function block.
func (RS *RS_FB) INIT() {
	RS.EN = true
	RS.ENO = true
	RS.S = false
	RS.R1 = false
	RS.Q1 = false
}

// RS executes the Reset-dominant logic.
func (RS *RS_FB) RS() {
	if !RS.EN {
		RS.ENO = false
		return
	}
	RS.ENO = true
	// Corrected logic: Q1 is set if S is true, but only if R1 is not true.
	RS.Q1 = !RS.R1 && (RS.S || RS.Q1)
}
