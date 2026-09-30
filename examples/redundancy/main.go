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

// This example runs the same logic on two resources pinned to different CPU
// cores and uses a voter to accept a result only when both copies agree.
//
// The two resources run concurrently, so their results are exchanged through
// vars.Shared values instead of plain fields; reading another resource's plain
// variables would be a data race.
//
// CPU affinity is supported on Linux and Windows. On other platforms, set
// Affinity to 0 on both resources.
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apiarytech/royaljelly/core"
	"github.com/apiarytech/royaljelly/iec"
	"github.com/apiarytech/royaljelly/vars"
)

// RedundantProgram holds one copy of the logic. Its internal state is private to
// the resource that runs it; only Output is shared.
type RedundantProgram struct {
	count  iec.LINT
	Output vars.Shared[iec.LINT]
}

// Logic is executed by the PLC task on its own resource.
func (p *RedundantProgram) Logic(now time.Time) {
	p.count++
	p.Output.Store(p.count)
}

func main() {
	// --- 1. Two independent instances of the same program logic ---
	progInstance1 := &RedundantProgram{}
	progInstance2 := &RedundantProgram{}

	// The verified result, readable from any goroutine.
	var confirmedOutput vars.Shared[iec.LINT]

	// --- 2. Two "CPU cores" (resources), each pinned to its own OS core ---
	cpuCore1 := &core.Resource{Name: "CPUCore1", Cycle: 100 * time.Millisecond, Affinity: 1}
	cpuCore2 := &core.Resource{Name: "CPUCore2", Cycle: 100 * time.Millisecond, Affinity: 2}

	// --- 3. One program instance per core ---
	cpuCore1.WithTask(core.NewTask("RedundantTask1", core.CyclicTask, 1, 500*time.Millisecond).
		WithProgram(&core.Program{Name: "Logic1", Logic: progInstance1.Logic}))
	cpuCore2.WithTask(core.NewTask("RedundantTask2", core.CyclicTask, 1, 500*time.Millisecond).
		WithProgram(&core.Program{Name: "Logic2", Logic: progInstance2.Logic}))

	// --- 4. A voter that compares the two results ---
	// It runs at a lower frequency on core 1 so both copies have produced output.
	cpuCore1.WithTask(core.NewTask("VoterTask", core.CyclicTask, 10, 1*time.Second).
		WithProgram(&core.Program{
			Name: "ResultVoter",
			Logic: func(now time.Time) {
				out1, out2 := progInstance1.Output.Load(), progInstance2.Output.Load()
				fmt.Printf("[%s] --- Voter Running ---\n", now.Format("15:04:05"))
				fmt.Printf("      Core 1 Output: %d\n", out1)
				fmt.Printf("      Core 2 Output: %d\n", out2)
				if out1 == out2 {
					confirmedOutput.Store(out1)
					fmt.Printf("      ✅ Results match. Confirmed output is now: %d\n", out1)
				} else {
					// The copies can differ by one for a moment because the cores are not
					// synchronized; a real voter would compare values from the same cycle.
					fmt.Printf("      ❌ Results DO NOT match. Confirmed output remains: %d\n", confirmedOutput.Load())
				}
			},
		}))

	// --- 5. Run the configuration for five seconds ---
	cfg := (&core.Configuration{Name: "RedundantConfig"}).WithResource(cpuCore1).WithResource(cpuCore2)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("Starting redundant PLC simulation...")
	if err := cfg.Run(ctx); err != nil {
		panic(err)
	}
	fmt.Printf("\nSimulation complete. Final confirmed output: %d\n", confirmedOutput.Load())
}
