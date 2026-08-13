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
	"github.com/apiarytech/royaljelly/iec"
)

// embeddedConfigString holds the content of config.txt as a Go string literal.
// This allows the application to be self-contained without external config files.
const embeddedConfigString = `
name: EmbeddedConfigExample

resource: MainCPU
  cycle: 50ms
  affinity: 1

  task: HighFrequencyTask
    type: Cyclic
    priority: 1
    interval: 250ms
    program: HighFreqInstance HighFreqProgram

  task: LowFrequencyTask
    type: Cyclic
    priority: 10
    interval: 1s
    program: LowFreqInstance LowFreqProgram
`

// HighFreqCounter and LowFreqCounter logic (can be in separate files or here)
type HighFreqCounter struct{ count iec.LINT }

func (p *HighFreqCounter) Logic(now time.Time) {
	p.count++
	fmt.Printf("[%s] High Frequency Task running. Count: %d\n", now.Format("15:04:05.000"), p.count)
}

type LowFreqCounter struct{ count iec.LINT }

func (p *LowFreqCounter) Logic(now time.Time) {
	p.count++
	fmt.Printf("[%s] Low Frequency Task running.  Count: %d\n", now.Format("15:04:05.000"), p.count)
}

func main() {
	hfc := &HighFreqCounter{}
	lfc := &LowFreqCounter{}

	config.RegisterProgramFactory("HighFreqProgram", func(params map[string]string) (func(time.Time), error) { return hfc.Logic, nil })
	config.RegisterProgramFactory("LowFreqProgram", func(params map[string]string) (func(time.Time), error) { return lfc.Logic, nil })

	fmt.Println("Loading configuration from embedded string...")
	cfg, err := config.LoadConfigurationFromString(embeddedConfigString)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Configuration '%s' loaded successfully with %d resource(s).\n", cfg.Name, len(cfg.Resources))

	for _, res := range cfg.Resources {
		res.Start()
	}
	fmt.Println("\nSimulation running for 5 seconds...")
	time.Sleep(5 * time.Second)
	for _, res := range cfg.Resources {
		res.Stop()
	}
	fmt.Println("\nSimulation complete.")
}
