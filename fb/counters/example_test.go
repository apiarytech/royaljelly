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

package counters_test

import (
	"fmt"

	"github.com/apiarytech/royaljelly/fb/counters"
)

func ExampleCTU() {
	// Count rising edges of CU; Q turns on when CV reaches PV.
	ctu := counters.CTU{EN: true, PV: 3}
	for i := 0; i < 3; i++ {
		ctu.CU = true // Rising edge: counts.
		ctu.Execute()
		ctu.CU = false
		ctu.Execute()
		fmt.Printf("CV=%d Q=%v\n", ctu.CV, ctu.Q)
	}
	ctu.R = true
	ctu.Execute()
	fmt.Printf("after reset CV=%d Q=%v\n", ctu.CV, ctu.Q)
	// Output:
	// CV=1 Q=false
	// CV=2 Q=false
	// CV=3 Q=true
	// after reset CV=0 Q=false
}

func ExampleCTD() {
	// LD loads PV; each rising edge of CD counts down; Q turns on at zero.
	ctd := counters.CTD{EN: true, PV: 2, LD: true}
	ctd.Execute()
	ctd.LD = false
	fmt.Printf("loaded CV=%d Q=%v\n", ctd.CV, ctd.Q)
	for i := 0; i < 2; i++ {
		ctd.CD = true
		ctd.Execute()
		ctd.CD = false
		ctd.Execute()
		fmt.Printf("CV=%d Q=%v\n", ctd.CV, ctd.Q)
	}
	// Output:
	// loaded CV=2 Q=false
	// CV=1 Q=false
	// CV=0 Q=true
}

func ExampleCTUD() {
	// Count up on CU and down on CD. QU reports CV >= PV and QD reports CV <= 0.
	ctud := counters.CTUD{EN: true, PV: 2}
	pulse := func(up bool) {
		if up {
			ctud.CU = true
		} else {
			ctud.CD = true
		}
		ctud.Execute()
		ctud.CU, ctud.CD = false, false
		ctud.Execute()
		fmt.Printf("CV=%d QU=%v QD=%v\n", ctud.CV, ctud.QU, ctud.QD)
	}
	pulse(true)
	pulse(true)
	pulse(false)
	pulse(false)
	// Output:
	// CV=1 QU=false QD=false
	// CV=2 QU=true QD=false
	// CV=1 QU=false QD=false
	// CV=0 QU=false QD=true
}
