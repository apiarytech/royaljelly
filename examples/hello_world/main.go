package main

import (
	"fmt"
	"time"

	"github.com/apiarytech/royaljelly/core"
)

func main() {
	// --- Instantiate and Initialize the Programs ---
	helloLogic := &HelloWorldProgram{}
	helloLogic.Init()

	// A second, simple program to demonstrate adding multiple programs to one task.
	var scanCount int64
	loggingLogic := func(now time.Time) {
		scanCount++
		fmt.Printf("... Scan cycle %d ...\n", scanCount)
	}

	// --- Configure and Assemble the PLC ---
	// 1. Create a task and hold it in a variable.
	helloWorldTask := core.NewTask("HelloWorldTask", core.CyclicTask, 1, 250*time.Millisecond)

	// 2. Use the fluent helper to add multiple programs to the *same* task.
	helloWorldTask.
		WithProgram(&core.Program{Name: "HelloWorldLogic", Logic: helloLogic.Logic}).
		WithProgram(&core.Program{Name: "LoggingProgram", Logic: loggingLogic})

	// 3. Add the fully configured task to the resource.
	resource := (&core.Resource{Name: "MainCPU", Cycle: 100 * time.Millisecond}).WithTask(helloWorldTask)

	// --- Start the PLC ---
	fmt.Println("--- Hello World Simulation ---")
	if err := resource.Start(); err != nil {
		panic(err)
	}

	// --- Simulate Interaction ---
	fmt.Printf("Initial state: Button=%v, Lamp=%v\n", helloLogic.MyButton.Load(), helloLogic.MyLamp.Load())

	// Press the button
	time.Sleep(1 * time.Second)
	fmt.Println(">>> Pressing button...")
	helloLogic.MyButton.Store(true)
	time.Sleep(1 * time.Second) // Hold it for a second
	fmt.Printf("Lamp state with button pressed: %v\n", helloLogic.MyLamp.Load())

	resource.Stop()
	fmt.Println("--- Simulation Complete ---")
}
