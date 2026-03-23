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

import "fmt"

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
	re   R_TRIG // Added for edge detection on R1
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
	KP, TR, TD    REAL
	CYCLE         TIME
	DIRECT_ACTION BOOL // true: SP-PV (heating), false: PV-SP (cooling)

	//Internal Variables
	XOUT   REAL // Simplified to a single output
	ERROR  REAL
	ITERM  INTEGRAL
	DTERM  DERIVATIVE
	autoRE R_TRIG // Rising edge trigger for AUTO
	autoFE F_TRIG // Falling edge trigger for AUTO
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
	Int.re.INIT()
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
	Pid.XOUT = 0 // Initialize the single output
	Pid.ERROR = 0
	Pid.ITERM.INIT()
	Pid.DTERM.INIT()
	Pid.autoRE.INIT()
	Pid.autoFE.INIT()
}

// DERIVATIVE calculates the derivative of the input signal XIN over time.
// When RUN is true, it continuously calculates the rate of change. When RUN is false, the output is reset to 0.
func (Dt *DERIVATIVE) DERIVATIVE() error {
	if !Dt.EN {
		Dt.ENO = false
		return nil
	} else {
		Dt.ENO = true
	}

	if Dt.RUN {
		cycleInSeconds := TIME_TO_REAL(Dt.CYCLE) / 1000
		if cycleInSeconds == 0 {
			return fmt.Errorf("DERIVATIVE error: CYCLE time cannot be zero")
		}
		Dt.XOUT = (3.0*(Dt.XIN-Dt.X3) + Dt.X1 - Dt.X2) / (10.0 * cycleInSeconds)
		Dt.X3 = Dt.X2
		Dt.X2 = Dt.X1
		Dt.X1 = Dt.XIN
	} else {
		Dt.XOUT = 0.0
		Dt.X1 = Dt.XIN // Reset history on not running
		Dt.X2 = Dt.XIN
		Dt.X3 = Dt.XIN
	}
	return nil
}

// INTEGRAL calculates the time integral of the input signal XIN.
// When RUN is true, it accumulates the input value over the CYCLE time.
// The integration can be reset to an initial value X0 by setting R1 to true.
func (Intg *INTEGRAL) INTEGRAL() error {
	if !Intg.EN {
		Intg.ENO = false
		return nil
	} else {
		Intg.ENO = true
	}

	Intg.re.CLK = Intg.R1
	Intg.re.R_TRIG()

	Intg.Q = !Intg.R1

	if Intg.re.Q {
		// On the rising edge of R1 (transition to manual), set the integral output to the initial condition.
		// This is key for bumpless transfer.
		Intg.XOUT = Intg.X0 // Set the integral output directly to the manual value.
	} else if Intg.RUN {
		cycleInSeconds := TIME_TO_REAL(Intg.CYCLE) / 1000
		if cycleInSeconds <= 0 {
			return fmt.Errorf("INTEGRAL error: CYCLE time must be positive")
		}
		Intg.XOUT = Intg.XOUT + Intg.XIN*cycleInSeconds
	}
	return nil
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
func (PID *PID) PID() error {
	if !PID.EN {
		PID.ENO = false
		return nil
	} else {
		PID.ENO = true
	}

	// 1. Detect mode changes first
	PID.autoRE.CLK = PID.AUTO
	PID.autoRE.R_TRIG()
	PID.autoFE.CLK = PID.AUTO
	PID.autoFE.F_TRIG()

	// 2. Handle Auto-to-Manual bumpless transfer
	if PID.autoFE.Q { // On falling edge of AUTO (Auto -> Manual)
		PID.X0 = PID.XOUT // Set manual setpoint to last auto output
	}

	// Calculate the error
	if !PID.DIRECT_ACTION { // Reverse action (cooling)
		PID.ERROR = PID.PV - PID.SP
	} else {
		PID.ERROR = PID.SP - PID.PV
	}

	// 3. Handle Manual-to-Auto bumpless transfer
	if PID.autoRE.Q { // On rising edge of AUTO (Manual -> Auto)
		// Calculate the initial integral value needed to make the first auto output match the manual output X0.
		// X0 = pTerm + iTerm => iTerm = X0 - pTerm
		PID.ITERM.X0 = PID.X0 - (PID.KP * PID.ERROR)
	}

	// 4. Configure and execute sub-blocks
	PID.ITERM.EN = PID.EN
	PID.ITERM.RUN = PID.AUTO
	PID.ITERM.R1 = !PID.AUTO
	if PID.TR == 0 {
		return fmt.Errorf("PID error: TR (integral time) cannot be zero")
	} else {
		PID.ITERM.XIN = (PID.KP / PID.TR) * PID.ERROR
	}
	PID.ITERM.CYCLE = PID.CYCLE
	if err := PID.ITERM.INTEGRAL(); err != nil {
		return fmt.Errorf("PID error in ITERM: %w", err)
	}

	PID.DTERM.EN = PID.EN
	PID.DTERM.RUN = PID.AUTO
	PID.DTERM.XIN = PID.ERROR
	PID.DTERM.CYCLE = PID.CYCLE
	if err := PID.DTERM.DERIVATIVE(); err != nil {
		return fmt.Errorf("PID error in DTERM: %w", err)
	}

	// 5. Calculate final output
	if PID.AUTO {
		pTerm := PID.KP * PID.ERROR
		iTerm := PID.ITERM.XOUT
		dTerm := PID.KP * PID.TD * PID.DTERM.XOUT
		PID.XOUT = pTerm + iTerm + dTerm
	} else {
		PID.XOUT = PID.X0 // In manual mode, output is simply X0
	}
	return nil
}
