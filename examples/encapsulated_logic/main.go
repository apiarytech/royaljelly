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
	"github.com/apiarytech/royaljelly/core"
)

var (
	// globalCycleCount acts as a VAR_GLOBAL, accessible by all programs in this package.
	globalCycleCount core.LINT = 0
)

// CounterProgram encapsulates the state (local tags) and logic for a counter.
// This is analogous to a Function Block in IEC 61131-3.
type CounterProgram struct {
	// This field 'Name' helps identify the instance.
	Name string
	// This field 'Output' is a "local tag" or instance variable.
	Output core.LINT
}

// Logic is the method that will be executed by the scheduler.
// It operates on the instance's own data, not on global variables.
func (p *CounterProgram) Logic(now time.Time) {
	// Increment the instance's local variable (VAR)
	p.Output++
	// Increment the shared global variable (VAR_GLOBAL)
	globalCycleCount++

	fmt.Printf("[%s] Counter '%s' running. Local Output: %d, Global Count: %d\n", now.Format("15:04:05.000"), p.Name, p.Output, globalCycleCount)
}

func main() {
	// 1. Create two independent instances of our CounterProgram.
	// Each has its own 'Output' variable.
	counterA := &CounterProgram{Name: "A"}
	counterB := &CounterProgram{Name: "B"}

	// 2. Register the 'Logic' method of each instance with the config loader.
	config.RegisterProgram("CounterA", counterA.Logic)
	config.RegisterProgram("CounterB", counterB.Logic)

	fmt.Println("Loading configuration from 'config.txt'...")
	cfg, err := config.LoadConfigurationFromFile("config.txt")
	if err != nil {
		panic(err)
	}

	// 3. Start all configured resources.
	for _, res := range cfg.Resources {
		res.Start()
	}

	fmt.Println("\nSimulation running for 2 seconds...")
	time.Sleep(2 * time.Second)
	fmt.Println("\nSimulation complete.")
	fmt.Printf("Final Global Cycle Count: %d\n", globalCycleCount)
}
