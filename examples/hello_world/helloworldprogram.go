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

package main

import (
	"time"

	"github.com/apiarytech/royaljelly/iec"
	"github.com/apiarytech/royaljelly/vars"
)

// HelloWorldProgram turns a lamp on while a button is pressed.
//
// The button is written by main (standing in for an input driver) and the lamp
// is read by main, while Logic runs on the scheduler goroutine. The two sides
// run concurrently, so both signals are vars.Shared values.
type HelloWorldProgram struct {
	MyButton vars.Shared[iec.BOOL] // (* Input from a switch *)
	MyLamp   vars.Shared[iec.BOOL] // (* Output to an indicator light *)
}

func (p *HelloWorldProgram) Init() {
	p.MyButton.Store(false)
	p.MyLamp.Store(false)
}

func (p *HelloWorldProgram) Logic(now time.Time) {
	p.MyLamp.Store(p.MyButton.Load())
}
