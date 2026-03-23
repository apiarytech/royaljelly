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
