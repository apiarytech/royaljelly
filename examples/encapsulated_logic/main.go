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
	"time"

	"github.com/apiarytech/royaljelly/config"
	. "github.com/apiarytech/royaljelly/iec"
)

var (
	// globalCycleCount acts as a VAR_GLOBAL, accessible by all programs in this package.
	globalCycleCount LINT = 0
)

// CounterProgram encapsulates the state (local tags) and logic for a counter.
// This is analogous to a Function Block in IEC 61131-3.
type CounterProgram struct {
	// This field 'Name' helps identify the instance.
	Name string
	// This field 'Output' is a "local tag" or instance variable.
	Output   LINT
	StepSize LINT
}

// Logic is the method that will be executed by the scheduler.
// It operates on the instance's own data, not on global variables.
func (p *CounterProgram) Logic(now time.Time) {
	// Increment the instance's local variable (VAR)
	p.Output += p.StepSize
	// Increment the shared global variable (VAR_GLOBAL)
	globalCycleCount++

	fmt.Printf("[%s] Counter '%s' running. Local Output: %d, Global Count: %d\n", now.Format("15:04:05.000"), p.Name, p.Output, globalCycleCount)
}

// NewCounterProgramFactory is a factory function that creates counter programs.
// It knows how to parse parameters from the config file to initialize the program state.
func NewCounterProgramFactory(params map[string]string) (func(time.Time), error) {
	// Create a new instance of our program struct.
	instance := &CounterProgram{StepSize: 1} // Default step size to 1

	// Apply parameters from the config file.
	if err := config.ParseLINT(params, "initial_value", &instance.Output); err != nil {
		return nil, err
	}
	if err := config.ParseLINT(params, "step_size", &instance.StepSize); err != nil {
		return nil, err
	}

	if name, ok := params["name"]; ok {
		instance.Name = name
	}

	// Return the 'Logic' method as the function to be executed by the scheduler.
	return instance.Logic, nil
}

func main() {
	// 1. Register our program "type" with the config loader by providing its factory function.
	// The loader will call this factory for each instance defined in the config file.
	config.RegisterProgramFactory("CounterProgram", NewCounterProgramFactory)

	fmt.Println("Loading configuration from 'config.txt'...")
	cfg, err := config.LoadConfigurationFromFile("config.txt")
	if err != nil {
		panic(err)
	}

	// 2. Start all configured resources.
	for _, res := range cfg.Resources {
		res.Start()
	}

	fmt.Println("\nSimulation running for 2 seconds...")
	time.Sleep(2 * time.Second)
	fmt.Println("\nSimulation complete.")
	fmt.Printf("Final Global Cycle Count: %d\n", globalCycleCount)
}
