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
	"sync"
	"time"

	. "github.com/apiarytech/royaljelly/core"
	. "github.com/apiarytech/royaljelly/iec"
)

// RedundantProgram holds the logic and its output.
// Each instance will have its own state.
type RedundantProgram struct {
	Output      LINT
	InjectFault BOOL // Flag to control fault injection
	// internal state for intermittent faults
	faultyCount LINT
}

// FaultyRedundantProgram is a struct that embeds the original
// RedundantProgram and adds fault injection capabilities for this example.
type FaultyRedundantProgram struct {
	RedundantProgram // Embed the original program
	faultyCount      LINT
}

// Logic is the function that will be executed by the PLC task.
func (p *RedundantProgram) Logic(now time.Time) {
	p.Output++

	// If fault injection is enabled, occasionally produce a wrong result.
	if p.InjectFault {
		p.faultyCount++
		// On every 3rd execution of this task, add extra to the output
		// to force a mismatch with the other core.
		if p.faultyCount > 2 && p.faultyCount%3 == 0 {
			p.Output += 10 // This will cause a mismatch.
			fmt.Printf("      ⚡️ Fault injected! Maliciously changed output to: %d\n", p.Output)
		}
	}
}

// Logic overrides the embedded Logic method to introduce faults.
func (p *FaultyRedundantProgram) Logic(now time.Time) {
	// Call the original, non-faulty logic first.
	p.RedundantProgram.Logic(now)

	// Now, add the fault injection logic.
	p.faultyCount++
	// On every 3rd execution, add to the output to cause a mismatch.
	if p.faultyCount%3 == 0 {
		p.Output += 10 // This will cause a mismatch.
		fmt.Printf("      ⚡️ Fault injected! Maliciously changed output to: %d\n", p.Output)
	}
}

func main() {
	// --- 1. Instantiate program instances ---
	// One is the standard program, the other is our faulty version.
	progInstance1 := &RedundantProgram{}
	progInstance2 := &FaultyRedundantProgram{}

	// This variable will hold the final, verified result.
	var confirmedOutput LINT
	var mu sync.RWMutex // Mutex to protect confirmedOutput

	// --- 2. Create two "CPU Cores" (Resources) ---
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
	voterTask := NewTask("VoterTask", CyclicTask, 10, 1*time.Second)
	voterTask.WithProgram(&Program{
		Name: "ResultVoter",
		Logic: func(now time.Time) {
			mu.Lock()
			defer mu.Unlock()

			fmt.Printf("[%s] --- Voter Running ---\n", now.Format("15:04:05"))
			fmt.Printf("      Core 1 Output: %d\n", progInstance1.Output)
			fmt.Printf("      Core 2 Output: %d\n", progInstance2.Output)

			// Only update the confirmed output if both instances agree.
			if progInstance1.Output == progInstance2.Output {
				confirmedOutput = progInstance1.Output
				fmt.Printf("      ✅ Results match. Confirmed output is now: %d\n", confirmedOutput)
			} else {
				// If they don't match, the `confirmedOutput` is left unchanged.
				fmt.Printf("      ❌ Results DO NOT match. Confirmed output remains: %d\n", confirmedOutput)
			}
		},
	})
	cpuCore1.WithTask(voterTask)

	// --- 5. Start the PLC ---
	fmt.Println("Starting redundant PLC simulation with fault injection...")
	cpuCore1.Start()
	cpuCore2.Start()

	// Let the simulation run long enough to see faults.
	time.Sleep(8 * time.Second)

	// Stop the resources
	cpuCore1.Stop()
	cpuCore2.Stop()
	fmt.Println("\nSimulation complete.")
}
