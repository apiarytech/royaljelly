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

package triggers_test

import (
	"fmt"

	"github.com/apiarytech/royaljelly/fb/triggers"
	"github.com/apiarytech/royaljelly/iec"
)

func Example_rTRIG() {
	// Q is true for exactly one call after CLK goes from false to true.
	var rt triggers.R_TRIG
	for _, clk := range []iec.BOOL{false, true, true, false, true} {
		rt.CLK = clk
		rt.R_TRIG()
		fmt.Print(rt.Q, " ")
	}
	fmt.Println()
	// Output: false true false false true
}

func Example_fTRIG() {
	// Q is true for exactly one call after CLK goes from true to false.
	var ft triggers.F_TRIG
	for _, clk := range []iec.BOOL{true, false, false, true, false} {
		ft.CLK = clk
		ft.F_TRIG()
		fmt.Print(ft.Q, " ")
	}
	fmt.Println()
	// Output: false true false false true
}

func Example_sR() {
	// Set-dominant latch: when S1 and R are both true, Q1 is set.
	var sr triggers.SR_FB
	sr.INIT()
	sr.S1 = true
	sr.SR()
	fmt.Println("set:", sr.Q1)
	sr.S1 = false
	sr.SR()
	fmt.Println("held:", sr.Q1)
	sr.S1, sr.R = true, true
	sr.SR()
	fmt.Println("both:", sr.Q1)
	sr.S1 = false
	sr.SR()
	fmt.Println("reset:", sr.Q1)
	// Output:
	// set: true
	// held: true
	// both: true
	// reset: false
}

func Example_rS() {
	// Reset-dominant latch: when S and R1 are both true, Q1 is reset.
	var rs triggers.RS_FB
	rs.INIT()
	rs.S = true
	rs.RS()
	fmt.Println("set:", rs.Q1)
	rs.R1 = true
	rs.RS()
	fmt.Println("both:", rs.Q1)
	// Output:
	// set: true
	// both: false
}
