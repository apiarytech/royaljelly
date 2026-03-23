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

// CTU Counter struct for all Count-Up/Count-Down/Count-UpDown Counters
type CTU struct {
	//ENABLES
	EN  BOOL //enable
	ENO BOOL //enable output
	//INPUTS
	CU BOOL
	R  BOOL
	PV INT
	//OUTPUTS
	Q  BOOL
	CV INT
	//INTERNAL
	re R_TRIG
}

// INIT CTU Initialization of Counters
func (CT *CTU) INIT() {
	CT.EN = true
	CT.ENO = true

	CT.CU = false
	CT.re.INIT()
	CT.R = false
	CT.PV = 0

	CT.Q = false
	CT.CV = 0
}

// CTU method to update CTU parameters
func (CTU *CTU) CTU() {
	if !CTU.EN {
		CTU.ENO = false
		return
	} else {
		CTU.ENO = true
	}

	CTU.re.CLK = CTU.CU
	CTU.re.R_TRIG()

	if CTU.R {
		CTU.CV = 0
	} else if CTU.re.Q { // Standard allows CV to increment past PV
		CTU.CV += 1
	}

	CTU.Q = CTU.CV >= CTU.PV

}

// CT Counter struct for all Count-Up/Count-Down/Count-UpDown Counters
type CTD struct {
	//ENABLES
	EN  BOOL //enable
	ENO BOOL //enable output
	//INPUTS
	CD BOOL
	LD BOOL
	PV INT
	//OUTPUTS
	Q  BOOL
	CV INT
	//INTERNAL
	re R_TRIG
}

// INIT CTD Countdown Timer Initialization of Counters
func (CT *CTD) INIT() {
	CT.EN = true
	CT.ENO = true

	CT.CD = false
	CT.re.INIT()
	CT.LD = false
	CT.PV = 0

	CT.Q = false
	CT.CV = 0
}

// CTD method to update CTD parameters
func (CTD *CTD) CTD() {
	if !CTD.EN {
		CTD.ENO = false
		return
	} else {
		CTD.ENO = true
	}

	CTD.re.CLK = CTD.CD
	CTD.re.R_TRIG()

	if CTD.LD {
		CTD.CV = CTD.PV
	} else if CTD.re.Q { // Standard allows CV to decrement past 0
		CTD.CV -= 1
	}

	CTD.Q = CTD.CV <= 0
}

// CTUD Counter struct for all Count-Up/Count-Down/Count-UpDown Counters
type CTUD struct {
	//ENABLES
	EN  BOOL //enable
	ENO BOOL //enable output
	//INPUTS
	CU BOOL
	CD BOOL
	R  BOOL
	LD BOOL
	PV INT
	//OUTPUS
	QU BOOL
	QD BOOL
	CV INT
	//INTERNAL
	reUP   R_TRIG
	reDOWN R_TRIG
}

// INIT CTD Countdown Timer Initialization of Counters
func (CT *CTUD) INIT() {
	//
	CT.EN = true
	CT.ENO = true

	CT.CU = false
	CT.CD = false
	CT.R = false
	CT.LD = false
	CT.PV = 0

	CT.QD = false
	CT.QU = false
	CT.CV = 0

	CT.reUP.INIT()
	CT.reDOWN.INIT()
}

// CTD method to update CTD parameters
func (CTUD *CTUD) CTUD() {
	if !CTUD.EN {
		CTUD.ENO = false
		return
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
		if CTUD.reUP.Q && !CTUD.reDOWN.Q { // Standard allows CV to increment past PV
			// Count up only on a rising edge of CU
			CTUD.CV += 1
		} else if CTUD.reDOWN.Q && !CTUD.reUP.Q { // Standard allows CV to decrement past 0
			// Count down only on a rising edge of CD
			CTUD.CV -= 1
		}
	}
	CTUD.QU = CTUD.CV >= CTUD.PV
	CTUD.QD = CTUD.CV <= 0
}
