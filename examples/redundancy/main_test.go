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
	"sync"
	"testing"
	"time"

	. "github.com/apiarytech/royaljelly/core"
)

// TestRedundancyWithFaultInjection verifies that the voter logic correctly
// handles matching and mismatching results from two redundant programs.
func TestRedundancyWithFaultInjection(t *testing.T) {
	// --- 1. Instantiate two separate instances of the same program logic ---
	// progInstance2 is configured to inject faults intermittently.
	progInstance1 := &RedundantProgram{}
	progInstance2 := &RedundantProgram{InjectFault: true}

	// This variable will hold the final, verified result.
	// We use a mutex to protect it during concurrent access from the test and voter task.
	var confirmedOutput LINT
	var mu sync.RWMutex

	// --- 2. Create two "CPU Cores" (Resources) ---
	// Use shorter cycle times for faster test execution.
	cpuCore1 := &Resource{Name: "TestCore1", Cycle: 10 * time.Millisecond}
	cpuCore2 := &Resource{Name: "TestCore2", Cycle: 10 * time.Millisecond}

	// --- 3. Create tasks and assign one program instance to each core ---
	// Use shorter intervals for the test.
	redundantTaskInterval := 50 * time.Millisecond
	task1 := NewTask("RedundantTask1", CyclicTask, 1, redundantTaskInterval)
	task1.WithProgram(&Program{Name: "Logic1", Logic: progInstance1.Logic})
	cpuCore1.WithTask(task1)

	task2 := NewTask("RedundantTask2", CyclicTask, 1, redundantTaskInterval)
	task2.WithProgram(&Program{Name: "Logic2", Logic: progInstance2.Logic})
	cpuCore2.WithTask(task2)

	// --- 4. Create a "Voter" program to compare results ---
	voterTask := NewTask("VoterTask", CyclicTask, 10, 100*time.Millisecond)
	voterTask.WithProgram(&Program{
		Name: "ResultVoter",
		Logic: func(now time.Time) {
			mu.Lock()
			defer mu.Unlock()
			// Only update the confirmed output if both instances agree.
			if progInstance1.Output == progInstance2.Output {
				confirmedOutput = progInstance1.Output
			}
		},
	})
	cpuCore1.WithTask(voterTask)

	// --- 5. Start the PLC and run the simulation ---
	cpuCore1.Start()
	cpuCore2.Start()

	// Let the simulation run long enough for a few cycles.
	// The fault is injected on the 3rd run of the logic.
	// Redundant tasks run at 50ms, so 3 runs take ~150ms.
	// The voter runs at 100ms.
	// Let's check the state after ~120ms (before the fault) and ~220ms (after the fault).

	// --- 6. Assertions ---

	// Check 1: Before the fault is injected.
	// After ~120ms, the redundant tasks should have run twice (at 50ms, 100ms).
	// The voter should have run once (at 100ms) and confirmed the matching result.
	time.Sleep(120 * time.Millisecond)

	mu.RLock()
	outputBeforeFault := confirmedOutput
	mu.RUnlock()

	if outputBeforeFault == 0 {
		t.Errorf("Expected confirmedOutput to be updated before fault, but it was 0")
	}
	t.Logf("State at 120ms: Confirmed output is %d (as expected)", outputBeforeFault)

	// Check 2: After the fault is injected.
	// At 150ms, progInstance2 injects a fault (Output becomes 3 + 10 = 13, while prog1 is 3).
	// At 200ms, the voter runs, sees the mismatch, and should NOT update confirmedOutput.
	time.Sleep(100 * time.Millisecond) // Sleep an additional 100ms to pass the 200ms mark.

	mu.RLock()
	outputAfterFault := confirmedOutput
	prog1Out := progInstance1.Output
	prog2Out := progInstance2.Output
	mu.RUnlock()

	if prog1Out == prog2Out {
		t.Fatalf("Fault was not correctly injected; outputs still match: %d == %d", prog1Out, prog2Out)
	}

	if outputAfterFault != outputBeforeFault {
		t.Errorf("FAIL: confirmedOutput was changed during a fault condition. Expected %d, got %d", outputBeforeFault, outputAfterFault)
	} else {
		t.Logf("PASS: confirmedOutput remained at %d during fault condition (Core1=%d, Core2=%d)", outputAfterFault, prog1Out, prog2Out)
	}

	// --- 7. Cleanup ---
	cpuCore1.Stop()
	cpuCore2.Stop()
}
