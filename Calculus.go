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

type DERIVATIVE struct {
	//Input Variables
	EN  BOOL //enable
	ENO BOOL //enable output

	RUN   BOOL
	XIN   REAL
	CYCLE TIME
	//Output Variables
	XOUT REAL
	//Internal Variables
	X1, X2, X3 REAL
}

type INTEGRAL struct {
	//Input Variables
	EN  BOOL //enable
	ENO BOOL //enable output

	RUN   BOOL
	R1    BOOL
	XIN   REAL
	X0    REAL
	CYCLE TIME
	//Output Variables
	XOUT REAL
	Q    BOOL
}

type HYSTERESIS struct {
	//Input Variables
	EN  BOOL //enable
	ENO BOOL //enable output

	//Internal Variables
	XIN1, XIN2, EPS REAL
	Q               BOOL
}

type PID struct {
	//Input Variables
	EN   BOOL //enable
	ENO  BOOL //enable output
	AUTO BOOL

	//Setpoints
	PV, SP, X0 REAL

	//Parameters
	KP, TR, TD     REAL
	CYCLE          TIME
	CONTROL_ACTION BOOL //0 - DIRECT (SP-PV); 1- REVERSE(PV-SP)

	//Internal Variables
	XOUT0 REAL
	XOUT1 REAL
	XOUT2 REAL
	//
	ERROR REAL
	ITERM INTEGRAL
	DTERM DERIVATIVE
}

// INIT initializes the INTEGRAL function block to its default state.
func (Int *INTEGRAL) INIT() {
	Int.EN = true
	Int.ENO = true
	Int.RUN = false
	Int.R1 = false
	Int.XIN = 0
	Int.CYCLE = INITTIME
	Int.XOUT = 0
	Int.X0 = 0
	Int.Q = false
}

// INIT initializes the DERIVATIVE function block to its default state.
func (Dev *DERIVATIVE) INIT() {
	Dev.EN = true
	Dev.ENO = true
	Dev.RUN = false
	Dev.XIN = 0
	Dev.CYCLE = INITTIME
	Dev.XOUT = 0
	Dev.X1 = 0
	Dev.X2 = 0
	Dev.X3 = 0
}

// INIT initializes the PID function block and its internal components to their default states.
func (Pid *PID) INIT() {
	Pid.EN = true
	Pid.ENO = true
	Pid.AUTO = false
	Pid.PV = 0
	Pid.SP = 0
	Pid.X0 = 0
	Pid.KP = 0
	Pid.TR = 0
	Pid.TD = 0
	Pid.CYCLE = INITTIME
	Pid.XOUT0 = 0
	Pid.XOUT1 = 0
	Pid.XOUT2 = 0
	Pid.ERROR = 0
	Pid.ITERM.INIT()
	Pid.DTERM.INIT()
}

// DERIVATIVE calculates the derivative of the input signal XIN over time.
// When RUN is true, it continuously calculates the rate of change. When RUN is false, the output is reset to 0.
func (Dt *DERIVATIVE) DERIVATIVE() {
	if !Dt.EN {
		Dt.ENO = false
		return
	} else {
		Dt.ENO = true
	}

	if Dt.RUN {
		Dt.XOUT = (3.0*(Dt.XIN-Dt.X3) + Dt.X1 - Dt.X2) / (10.0 * TIME_TO_REAL(Dt.CYCLE) / 1000)
		Dt.X3 = Dt.X2
		Dt.X2 = Dt.X1
		Dt.X1 = Dt.XIN
	} else {
		Dt.XOUT = 0.0
		Dt.X1 = Dt.XIN
		Dt.X2 = Dt.XIN
		Dt.X3 = Dt.XIN
	}

}

// INTEGRAL calculates the time integral of the input signal XIN.
// When RUN is true, it accumulates the input value over the CYCLE time.
// The integration can be reset to an initial value X0 by setting R1 to true.
func (Intg *INTEGRAL) INTEGRAL() {
	if !Intg.EN {
		Intg.ENO = false
		return
	} else {
		Intg.ENO = true
	}

	Intg.Q = !Intg.R1

	if Intg.R1 {
		// When R1 is true (e.g., PID in manual), the integral output should be held at X0.
		// This is key for bumpless transfer when switching to AUTO.
		Intg.XOUT = Intg.X0 // Set the integral output directly to the manual value.
	} else if Intg.RUN {
		Intg.XOUT = Intg.XOUT + Intg.XIN*TIME_TO_REAL(Intg.CYCLE)/1000
	}
}

// HYSTERESIS provides a boolean output with a deadband, determined by EPS.
// The output Q becomes true when XIN1 exceeds XIN2 + EPS and becomes false when XIN1 drops below XIN2 - EPS.
func (Hyst *HYSTERESIS) HYSTERESIS() {
	if !Hyst.EN {
		Hyst.ENO = false
		return
	} else {
		Hyst.ENO = true
	}

	if Hyst.Q {
		if Hyst.XIN1 < (Hyst.XIN2 - Hyst.EPS) {
			Hyst.Q = false
		}
	} else if Hyst.XIN1 > (Hyst.XIN2 + Hyst.EPS) {
		Hyst.Q = true
	}
}

// PID implements a Proportional-Integral-Derivative controller.
// It calculates an output value based on the error between a setpoint (SP) and a process variable (PV).
func (PID *PID) PID() {
	if !PID.EN {
		PID.ENO = false
		return
	} else {
		PID.ENO = true
	}
	// Calculate the error
	if PID.CONTROL_ACTION {
		PID.ERROR = PID.PV - PID.SP
	} else {
		PID.ERROR = PID.SP - PID.PV
	}

	PID.ITERM.EN = PID.EN
	PID.ITERM.RUN = PID.AUTO
	PID.ITERM.R1 = !PID.AUTO
	// The integral term's input should be the error scaled by KP/TR for the standard form.
	if PID.TR != 0 {
		PID.ITERM.XIN = (PID.KP / PID.TR) * PID.ERROR
	}
	// For bumpless transfer, when switching from Manual (AUTO=false) to Auto,
	// the integral term should start at a value that makes the initial auto output
	// match the manual output X0.
	// X0 = pTerm + iTerm => iTerm = X0 - pTerm
	PID.ITERM.X0 = PID.X0 - (PID.KP * PID.ERROR)
	PID.ITERM.CYCLE = PID.CYCLE
	PID.ITERM.INTEGRAL()

	PID.DTERM.EN = PID.EN
	PID.DTERM.RUN = PID.AUTO
	PID.DTERM.XIN = PID.ERROR
	PID.DTERM.CYCLE = PID.CYCLE
	PID.DTERM.DERIVATIVE()

	// Proportional Term
	pTerm := PID.KP * PID.ERROR
	// Integral Term (already calculated by INTEGRAL block)
	iTerm := PID.ITERM.XOUT
	// Derivative Term
	dTerm := PID.KP * PID.TD * PID.DTERM.XOUT

	// Standard Parallel PID: Y = Kp*e + Ki*Integral(e) + Kd*Derivative(e)
	PID.XOUT0 = pTerm + iTerm + dTerm
	// Non-Interactive (Series & Parallel) PID: Y = Kp * (e + 1/Tr * Integral(e) + Td * Derivative(e))
	PID.XOUT1 = PID.KP * (PID.ERROR + (1/PID.TR)*PID.ITERM.XOUT + PID.TD*PID.DTERM.XOUT)
	// Interactive (Series) PID: Y = Kp * (e + 1/Tr * Integral(e)) * (1 + Td*Derivative(e))
	PID.XOUT2 = PID.KP * (PID.ERROR + (1/PID.TR)*PID.ITERM.XOUT) * (1 + PID.TD*PID.DTERM.XOUT)
}
