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

// Package readme holds the code snippets shown in README.md as runnable
// examples, so `go test ./...` fails if the README drifts from the API.
// Keep each example in sync with the matching README section.
package readme

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/apiarytech/royaljelly/config"
	"github.com/apiarytech/royaljelly/core"
	"github.com/apiarytech/royaljelly/fb/timers"
	"github.com/apiarytech/royaljelly/iec"
	"github.com/apiarytech/royaljelly/std/arithmetic"
	"github.com/apiarytech/royaljelly/std/selection"
	"github.com/apiarytech/royaljelly/std/strings"
	"github.com/apiarytech/royaljelly/vars"
)

// README: "Quick start: build a configuration in code".
func Example_quickStart() {
	// Program state lives in a struct, like the local variables of an IEC program.
	type Blinker struct {
		Lamp iec.BOOL
		Runs int
	}
	blinker := &Blinker{}

	task := core.NewTask("BlinkTask", core.CyclicTask, 1, 20*time.Millisecond).
		WithProgram(&core.Program{
			Name: "Blink",
			Logic: func(now time.Time) {
				blinker.Lamp = !blinker.Lamp
				blinker.Runs++
			},
		})
	resource := (&core.Resource{Name: "MainCPU", Cycle: 5 * time.Millisecond}).WithTask(task)
	cfg := (&core.Configuration{Name: "Demo"}).WithResource(resource)

	// Run validates the configuration, starts every resource, and stops them
	// all when the context ends.
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := cfg.Run(ctx); err != nil {
		fmt.Println("start failed:", err)
		return
	}
	fmt.Println("blinked:", blinker.Runs > 0)
	// Output: blinked: true
}

// README: "Load the structure from a config file".
func Example_configFile() {
	// A factory creates one program instance per "program:" line.
	config.RegisterProgramFactory("CounterProgram", func(params map[string]string) (func(time.Time), error) {
		var count iec.LINT
		if err := config.ParseLINT(params, "initial_value", &count); err != nil {
			return nil, err
		}
		return func(now time.Time) { count++ }, nil
	})

	cfg, err := config.LoadConfigurationFromString(`
name: EncapsulationExample
resource: MainCPU
  cycle: 50ms
  task: TaskA
    type: Cyclic
    priority: 1
    interval: 250ms
    program: CounterA CounterProgram
      param: initial_value 100
`)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, r := range cfg.Resources() {
		for _, t := range r.Tasks() {
			fmt.Printf("%s/%s every %v\n", r.Name, t.Name, t.Interval)
		}
	}
	// Output: MainCPU/TaskA every 250ms
}

// README: "Share data between resources".
func Example_sharedData() {
	var setpoint vars.Shared[iec.REAL]
	setpoint.Store(21.5) // e.g. written by an HMI goroutine

	var pi vars.ProcessImage
	pi.Write(func(img *vars.Image) { img.I.B[0] = true }) // e.g. an input driver sets %IX0

	var start iec.BOOL
	pi.Read(func(img *vars.Image) { start = img.I.B[0] })
	fmt.Println(setpoint.Load(), start)
	// Output: 21.5 true
}

// README: "Faults, watchdog and metrics".
func Example_faults() {
	task := core.NewTask("Control", core.CyclicTask, 1, 10*time.Millisecond)
	task.Watchdog = 5 * time.Millisecond // Report runs longer than 5ms.
	task.AddProgram(&core.Program{Name: "Faulty", Logic: func(time.Time) { panic("sensor fault") }})

	faults := make(chan core.Fault, 1)
	resource := &core.Resource{
		Name: "MainCPU",
		OnFault: func(f core.Fault) {
			select {
			case faults <- f:
			default:
			}
		},
	}
	resource.AddTask(task)
	if err := resource.Start(); err != nil {
		fmt.Println(err)
		return
	}
	f := <-faults
	resource.Stop()

	fmt.Println(f.Kind, f.Program, errors.Is(f, core.ErrProgramPanic))
	fmt.Println("runs recorded:", task.Stats().Runs > 0)
	// Output:
	// panic Faulty true
	// runs recorded: true
}

// README: "Function blocks".
func Example_functionBlock() {
	ton := timers.TON{PT: iec.TIME(5 * time.Second)}
	now := time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)
	for scan := 0; scan <= 6; scan++ {
		ton.IN = scan >= 1 // Input switches on at scan 1.
		ton.Execute(now)
		fmt.Printf("scan %d: Q=%v ET=%v\n", scan, ton.Q, time.Duration(ton.ET))
		now = now.Add(time.Second)
	}
	// Output:
	// scan 0: Q=false ET=0s
	// scan 1: Q=false ET=0s
	// scan 2: Q=false ET=1s
	// scan 3: Q=false ET=2s
	// scan 4: Q=false ET=3s
	// scan 5: Q=false ET=4s
	// scan 6: Q=true ET=5s
}

// README: "Standard functions".
func Example_standardFunctions() {
	sum := arithmetic.ADD(iec.REAL(10), iec.REAL(5.5))
	length := strings.LEN("Hello, RoyalJelly!")
	largest, err := selection.MAX(iec.LINT(100), iec.LINT(50), iec.LINT(120))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sum, length, largest)
	// Output: 15.5 18 120
}

// README: "Data types".
func Example_dataTypes() {
	var myBool iec.BOOL = true
	var myInt iec.INT = 123
	var myReal iec.REAL = 45.67
	var myTime iec.TIME = iec.TIME(10 * time.Second)
	var myText iec.WSTRING = "Grüße"
	fmt.Println(myBool, myInt, myReal, time.Duration(myTime), myText)
	// Output: true 123 45.67 10s Grüße
}
