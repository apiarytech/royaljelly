//go:build !tinygo

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

	. "github.com/apiarytech/royaljelly/core"
	. "github.com/apiarytech/royaljelly/iec"
)

// RedundantProgram holds the logic and its output.
// Each instance will have its own state.
type RedundantProgram struct {
	Output LINT
}

// Logic is the function that will be executed by the PLC task.
func (p *RedundantProgram) Logic(now time.Time) {
	p.Output++
}

func main() {
	// --- 1. Instantiate two separate instances of the same program logic ---
	// Each instance has its own memory (`Output` field).
	progInstance1 := &RedundantProgram{}
	progInstance2 := &RedundantProgram{}

	// This variable will hold the final, verified result.
	var confirmedOutput LINT

	// --- 2. Create two "CPU Cores" (Resources) ---
	// Assign each resource to a specific OS-level CPU core.
	cpuCore1 := &Resource{Name: "CPUCore1", Cycle: 100 * time.Millisecond, Affinity: 1}
	cpuCore2 := &Resource{Name: "CPUCore2", Cycle: 100 * time.Millisecond, Affinity: 2}

	// --- 3. Create tasks and assign one program instance to each core ---
	task1 := NewTask("RedundantTask1", CyclicTask, 1, 500*time.Millisecond)
	task1.WithProgram(&Program{Name: "Logic1", Logic: progInstance1.Logic})
	cpuCore1.WithTask(task1)

	task2 := NewTask("RedundantTask2", CyclicTask, 1, 500*time.Millisecond)
	task2.WithProgram(&Program{Name: "Logic2", Logic: progInstance2.Logic})
	cpuCore2.WithTask(task2)

	// --- 4. Create a "Voter" program to compare results ---
	// This program runs on one of the cores (or could be on a third, slower one).
	// It runs at a slightly lower frequency to ensure the redundant tasks have run.
	voterTask := NewTask("VoterTask", CyclicTask, 10, 1*time.Second)
	voterTask.WithProgram(&Program{
		Name: "ResultVoter",
		Logic: func(now time.Time) {
			fmt.Printf("[%s] --- Voter Running ---\n", now.Format("15:04:05"))
			fmt.Printf("      Core 1 Output: %d\n", progInstance1.Output)
			fmt.Printf("      Core 2 Output: %d\n", progInstance2.Output)

			// The core of the verification logic:
			// Only update the confirmed output if both instances agree.
			if progInstance1.Output == progInstance2.Output {
				confirmedOutput = progInstance1.Output
				fmt.Printf("      ✅ Results match. Confirmed output is now: %d\n", confirmedOutput)
			} else {
				// If they don't match, the `confirmedOutput` is left unchanged.
				// You could also add error handling, logging, or fault logic here.
				fmt.Printf("      ❌ Results DO NOT match. Confirmed output remains: %d\n", confirmedOutput)
			}
		},
	})
	cpuCore1.WithTask(voterTask) // Add the voter to one of the cores.

	// --- 5. Add resources to the main configuration ---
	config := &Configuration{Name: "RedundantConfig"}
	config.WithResource(cpuCore1).WithResource(cpuCore2)

	// --- 6. Start the PLC ---
	fmt.Println("Starting redundant PLC simulation...")
	cpuCore1.Start()
	cpuCore2.Start()

	// Let the simulation run for a few seconds.
	time.Sleep(5 * time.Second)

	// Stop the resources
	cpuCore1.Stop()
	cpuCore2.Stop()
	fmt.Println("\nSimulation complete.")
}
