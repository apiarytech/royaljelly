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

// This example extends the redundancy example: one of the two redundant copies
// periodically corrupts its result, and the voter detects the mismatch. It also
// shows a program panic being reported through the fault handler while the
// other programs keep running.
//
// CPU affinity is supported on Linux and Windows. On other platforms, set
// Affinity to 0 on both resources.
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/apiarytech/royaljelly/core"
	"github.com/apiarytech/royaljelly/iec"
	"github.com/apiarytech/royaljelly/vars"
)

// RedundantProgram holds one copy of the logic. Output is shared with the voter,
// which runs on another resource.
type RedundantProgram struct {
	count  iec.LINT
	Output vars.Shared[iec.LINT]

	// InjectFault makes every third run publish a corrupted value.
	InjectFault bool
}

// Logic is the function executed by the PLC task.
func (p *RedundantProgram) Logic(now time.Time) {
	p.count++
	out := p.count
	if p.InjectFault && p.count%3 == 0 {
		out += 10 // Corrupt the published value to force a mismatch.
		fmt.Printf("      ⚡️ Fault injected! Published output %d instead of %d\n", out, p.count)
	}
	if p.InjectFault && p.count == 7 {
		panic("simulated hardware exception")
	}
	p.Output.Store(out)
}

func main() {
	healthy := &RedundantProgram{}
	faulty := &RedundantProgram{InjectFault: true}
	var confirmedOutput iec.LINT // Only the voter's resource touches this.

	cpuCore1 := &core.Resource{Name: "CPUCore1", Cycle: 100 * time.Millisecond, Affinity: 1}
	cpuCore2 := &core.Resource{Name: "CPUCore2", Cycle: 100 * time.Millisecond, Affinity: 2}

	cpuCore1.WithTask(core.NewTask("RedundantTask1", core.CyclicTask, 1, 500*time.Millisecond).
		WithProgram(&core.Program{Name: "Logic1", Logic: healthy.Logic}))
	cpuCore2.WithTask(core.NewTask("RedundantTask2", core.CyclicTask, 1, 500*time.Millisecond).
		WithProgram(&core.Program{Name: "Logic2", Logic: faulty.Logic}))

	cpuCore1.WithTask(core.NewTask("VoterTask", core.CyclicTask, 10, 1*time.Second).
		WithProgram(&core.Program{
			Name: "ResultVoter",
			Logic: func(now time.Time) {
				out1, out2 := healthy.Output.Load(), faulty.Output.Load()
				fmt.Printf("[%s] --- Voter: core 1 = %d, core 2 = %d\n", now.Format("15:04:05"), out1, out2)
				if out1 == out2 {
					confirmedOutput = out1
					fmt.Printf("      ✅ Results match. Confirmed output is now: %d\n", confirmedOutput)
				} else {
					fmt.Printf("      ❌ Results DO NOT match. Confirmed output remains: %d\n", confirmedOutput)
				}
			},
		}))

	cfg := &core.Configuration{
		Name: "FaultInjectionConfig",
		// The fault handler receives recovered panics, overruns and watchdog trips
		// from every resource that has no handler of its own.
		OnFault: func(f core.Fault) {
			if errors.Is(f, core.ErrProgramPanic) {
				fmt.Printf("      🛑 %s program %q panicked: %v (other programs keep running)\n", f.Resource, f.Program, f.Value)
				return
			}
			fmt.Println("      ⚠️ fault:", f)
		},
	}
	cfg.WithResource(cpuCore1).WithResource(cpuCore2)

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	fmt.Println("Starting redundant PLC simulation with fault injection...")
	if err := cfg.Run(ctx); err != nil {
		panic(err)
	}
	fmt.Println("\nSimulation complete.")
}
