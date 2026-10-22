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

package counters

import (
	"github.com/apiarytech/royaljelly/fb/triggers"
	"github.com/apiarytech/royaljelly/iec"
)

// CTU is the IEC 61131-3 count-up counter. CV saturates at the INT maximum.
//
// The block runs only while EN is TRUE. EN is FALSE in a zero value, so set
// EN to TRUE before calling Execute; a zero-value CTU{} otherwise does nothing.
// The same applies to CTD and CTUD.
type CTU struct {
	//ENABLES
	EN  iec.BOOL //enable
	ENO iec.BOOL //enable output
	//INPUTS
	CU iec.BOOL
	R  iec.BOOL
	PV iec.INT
	//OUTPUTS
	Q  iec.BOOL
	CV iec.INT
	//INTERNAL
	re          triggers.R_TRIG
	initialized iec.BOOL
}

// INIT CTU Initialization of Counters
func (CTU *CTU) INIT() {
	CTU.Q = false
	CTU.CV = 0
	CTU.re.INIT()
	CTU.initialized = true
}

// CTU method to update CTU parameters
func (CTU *CTU) Execute() {
	if !CTU.EN {
		CTU.ENO = false
		return
	} else {
		CTU.ENO = true
	}
	if !CTU.initialized {
		CTU.INIT()
	} else {
		CTU.ENO = true
	}

	CTU.re.CLK = CTU.CU
	CTU.re.R_TRIG()

	if CTU.R {
		CTU.CV = 0
	} else if CTU.re.Q && CTU.CV < iec.MAXINT { // CV may pass PV but saturates at PVmax
		CTU.CV += 1
	}

	CTU.Q = CTU.CV >= CTU.PV

}

// CTD is the IEC 61131-3 count-down counter. CV saturates at the INT minimum.
// Set EN to TRUE before calling Execute, as for CTU.
type CTD struct {
	//ENABLES
	EN  iec.BOOL //enable
	ENO iec.BOOL //enable output
	//INPUTS
	CD iec.BOOL
	LD iec.BOOL
	PV iec.INT
	//OUTPUTS
	Q  iec.BOOL
	CV iec.INT
	//INTERNAL
	re          triggers.R_TRIG
	initialized iec.BOOL
}

// INIT CTD Countdown Timer Initialization of Counters
func (CTD *CTD) INIT() {
	CTD.Q = false
	CTD.CV = 0
	CTD.re.INIT()
	CTD.initialized = true
}

// CTD method to update CTD parameters
func (CTD *CTD) Execute() {
	if !CTD.EN {
		CTD.ENO = false
		return
	} else {
		CTD.ENO = true
	}
	if !CTD.initialized {
		CTD.INIT()
	} else {
		CTD.ENO = true
	}

	CTD.re.CLK = CTD.CD
	CTD.re.R_TRIG()

	if CTD.LD {
		CTD.CV = CTD.PV
	} else if CTD.re.Q && CTD.CV > iec.MININT { // CV may pass 0 but saturates at PVmin
		CTD.CV -= 1
	}

	CTD.Q = CTD.CV <= 0
}

// CTUD is the IEC 61131-3 up-down counter. CV saturates at the INT limits.
// Set EN to TRUE before calling Execute, as for CTU.
type CTUD struct {
	//ENABLES
	EN  iec.BOOL //enable
	ENO iec.BOOL //enable output
	//INPUTS
	CU iec.BOOL
	CD iec.BOOL
	R  iec.BOOL
	LD iec.BOOL
	PV iec.INT
	//OUTPUS
	QU iec.BOOL
	QD iec.BOOL
	CV iec.INT
	//INTERNAL
	reUP        triggers.R_TRIG
	reDOWN      triggers.R_TRIG
	initialized iec.BOOL
}

// INIT CTD Countdown Timer Initialization of Counters
func (CTUD *CTUD) INIT() {
	CTUD.QD = false
	CTUD.QU = false
	CTUD.CV = 0
	CTUD.reUP.INIT()
	CTUD.reDOWN.INIT()
	CTUD.initialized = true
}

// CTD method to update CTD parameters
func (CTUD *CTUD) Execute() {
	if !CTUD.EN {
		CTUD.ENO = false
		return
	} else {
		CTUD.ENO = true
	}
	if !CTUD.initialized {
		CTUD.INIT()
	} else {
		CTUD.ENO = true
	}

	CTUD.reUP.CLK = CTUD.CU
	CTUD.reUP.R_TRIG()
	CTUD.reDOWN.CLK = CTUD.CD
	CTUD.reDOWN.R_TRIG()

	if CTUD.R {
		CTUD.CV = 0
	} else if CTUD.LD {
		CTUD.CV = CTUD.PV
	} else {
		// As per standard, if both CU and CD have a rising edge in the same scan, no action is taken.
		if CTUD.reUP.Q && !CTUD.reDOWN.Q && CTUD.CV < iec.MAXINT { // CV may pass PV but saturates at PVmax
			// Count up only on a rising edge of CU
			CTUD.CV += 1
		} else if CTUD.reDOWN.Q && !CTUD.reUP.Q && CTUD.CV > iec.MININT { // CV may pass 0 but saturates at PVmin
			// Count down only on a rising edge of CD
			CTUD.CV -= 1
		}
	}
	CTUD.QU = CTUD.CV >= CTUD.PV
	CTUD.QD = CTUD.CV <= 0
}
