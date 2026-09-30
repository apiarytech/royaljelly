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

package calculus_test

import (
	"fmt"
	"time"

	"github.com/apiarytech/royaljelly/fb/calculus"
	"github.com/apiarytech/royaljelly/iec"
)

func ExampleINTEGRAL() {
	// Integrate a constant input of 2.0 per second, called every 100ms.
	var in calculus.INTEGRAL
	in.INIT()
	in.RUN = true
	in.XIN = 2.0
	in.CYCLE = iec.TIME(100 * time.Millisecond)
	for i := 0; i < 10; i++ {
		if err := in.INTEGRAL(); err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Printf("after 1s: %.2f\n", in.XOUT)
	// Output: after 1s: 2.00
}

func ExampleDERIVATIVE() {
	// The derivative of a ramp rising 0.5 per 100ms settles at 5.0 per second.
	var d calculus.DERIVATIVE
	d.INIT()
	d.CYCLE = iec.TIME(100 * time.Millisecond)
	d.DERIVATIVE() // RUN is false: primes the history with XIN.
	d.RUN = true
	for i := 1; i <= 5; i++ {
		d.XIN = iec.REAL(i) * 0.5
		if err := d.DERIVATIVE(); err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Printf("slope: %.2f\n", d.XOUT)
	// Output: slope: 5.00
}

func ExampleHYSTERESIS() {
	// Q switches on above XIN2+EPS and off below XIN2-EPS, ignoring noise in between.
	h := calculus.HYSTERESIS{EN: true, XIN2: 50, EPS: 2}
	for _, temp := range []iec.REAL{49, 51, 53, 51, 49, 47} {
		h.XIN1 = temp
		h.HYSTERESIS()
		fmt.Printf("%.0f:%v ", temp, h.Q)
	}
	fmt.Println()
	// Output: 49:false 51:false 53:true 51:true 49:true 47:false
}

func ExamplePID() {
	// A proportional-integral controller in automatic mode drives its output up
	// while the process variable is below the setpoint.
	var pid calculus.PID
	pid.INIT()
	pid.AUTO = true
	pid.DIRECT_ACTION = true // Error = SP - PV.
	pid.KP, pid.TR, pid.TD = 2, 1, 0
	pid.CYCLE = iec.TIME(100 * time.Millisecond)
	pid.SP, pid.PV = 10, 8
	for i := 0; i < 3; i++ {
		if err := pid.PID(); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("error=%.1f output=%.2f\n", pid.ERROR, pid.XOUT)
	}
	// Output:
	// error=2.0 output=4.40
	// error=2.0 output=4.80
	// error=2.0 output=5.20
}
