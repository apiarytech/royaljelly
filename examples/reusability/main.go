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
	"fmt"
	"math/rand"
	"time"

	"github.com/apiarytech/royaljelly/config"
	"github.com/apiarytech/royaljelly/core"
)

// PIDController encapsulates the state and logic for a PID control loop.
// This struct acts as a reusable Function Block type.
type PIDController struct {
	Name            string
	Setpoint        core.REAL // Desired value
	Kp              core.REAL // Proportional gain
	Ki              core.REAL // Integral gain
	Kd              core.REAL // Derivative gain
	ProcessVariable core.REAL // Measured value from a sensor
	Output          core.REAL // Calculated output to an actuator

	// Internal state for PID calculation
	integral  core.REAL
	lastError core.REAL
}

// Logic simulates a full PID controller.
func (p *PIDController) Logic(now time.Time) {
	// Simulate reading a sensor value that fluctuates around the setpoint.
	p.ProcessVariable = p.Setpoint + core.REAL(rand.Float32()*2-1) // Fluctuates by +/- 1.0

	// Calculate error
	err := p.Setpoint - p.ProcessVariable

	// Integral term
	p.integral += err

	// Derivative term
	derivative := err - p.lastError

	// Calculate the PID output
	p.Output = (p.Kp * err) + (p.Ki * p.integral) + (p.Kd * derivative)

	// Update state for the next cycle
	p.lastError = err

	fmt.Printf("[%s] PID '%s': Setpoint=%.2f, PV=%.2f, Output=%.2f\n",
		now.Format("15:04:05.000"), p.Name, p.Setpoint, p.ProcessVariable, p.Output)
}

// NewPIDControllerFactory is the factory function that the config loader uses
// to create and configure instances of the PIDController.
func NewPIDControllerFactory(params map[string]string) (func(time.Time), error) {
	// Create a new instance with default values.
	instance := &PIDController{Kp: 1.0, Ki: 0.0, Kd: 0.0}

	// Use the parameter helpers to parse values from the config file.
	if err := config.ParseREAL(params, "setpoint", &instance.Setpoint); err != nil {
		return nil, err
	}
	if err := config.ParseREAL(params, "kp", &instance.Kp); err != nil {
		return nil, err
	}
	if err := config.ParseREAL(params, "ki", &instance.Ki); err != nil {
		return nil, err
	}
	if err := config.ParseREAL(params, "kd", &instance.Kd); err != nil {
		return nil, err
	}

	// The loader automatically provides the instance name.
	instance.Name = params["name"]

	// Return the 'Logic' method as the function to be executed.
	return instance.Logic, nil
}

func main() {
	// 1. Register the "PID_Controller" type by providing its factory.
	config.RegisterProgramFactory("PID_Controller", NewPIDControllerFactory)

	fmt.Println("Loading configuration from 'config.txt'...")
	cfg, err := config.LoadConfigurationFromFile("config.txt")
	if err != nil {
		panic(err)
	}

	// 2. Start all configured resources.
	for _, res := range cfg.Resources {
		res.Start()
	}

	fmt.Println("\nSimulation running for 3 seconds...")
	time.Sleep(3 * time.Second)
	fmt.Println("\nSimulation complete.")
}
