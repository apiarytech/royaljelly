package main

import (
	"fmt"
	"time"

	. "github.com/apiarytech/royaljelly/core"
)

// This example simulates a 4-way traffic light controller using a state machine.
//
// Sequence:
// 1. North/South lights are Green, East/West are Red.
// 2. North/South lights turn Yellow, East/West remain Red.
// 3. North/South lights turn Red, East/West turn Green.
// 4. North/South remain Red, East/West turn Yellow.
// 5. The cycle repeats.

func main() {
	// --- Instantiate and Initialize the Program ---
	trafficLogic := &TrafficLightProgram{}
	trafficLogic.Init()

	// --- Configure and Assemble the PLC using a fluent API ---
	resource := (&Resource{Name: "MainCPU", Cycle: 100 * time.Millisecond}).
		WithTask(
			NewTask("TrafficLightTask", CyclicTask, 1, 1*time.Second).
				WithProgram(&Program{
					Name:  "TrafficLightLogic",
					Logic: trafficLogic.Logic,
				}),
		)

	// --- Start the PLC ---
	fmt.Println("--- Traffic Light Simulation ---")
	resource.Start()

	// Keep the simulation running for a while.
	time.Sleep(40 * time.Second)
	resource.Stop()
	fmt.Println("--- Simulation Complete ---")
}
