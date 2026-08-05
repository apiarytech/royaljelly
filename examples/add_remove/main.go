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
	"sync/atomic"
	"time"

	"github.com/apiarytech/royaljelly/core"
)

// IsRunning checks the private 'running' field of a resource.
// This is a helper for demonstration purposes.
func IsRunning(r *core.Resource) bool {
	type runningChecker interface {
		IsRunning() bool
	}
	if checker, ok := interface{}(r).(runningChecker); ok {
		return checker.IsRunning()
	}
	// Fallback or error if the method doesn't exist.
	// For this example, we assume it might not be public.
	return false
}

func main() {
	// --- 1. Setup the PLC structure ---
	var highFreqCounter, lowFreqCounter atomic.Int32

	// Create a high-frequency task that runs every 250ms
	highFreqTask := core.NewTask("HighFrequencyTask", core.CyclicTask, 1, 250*time.Millisecond)
	highFreqTask.WithProgram(&core.Program{
		Name: "HighFreqCounter",
		Logic: func(now time.Time) {
			highFreqCounter.Add(1)
			fmt.Printf("[%s] HighFrequencyTask running (count: %d)\n", now.Format("15:04:05.000"), highFreqCounter.Load())
		},
	})
	// `ConvertTo` function in `helper.go` is quite large. Can you refactor it for better readability and maintenance?
	// Create a low-frequency task that runs every 1 second
	lowFreqTask := core.NewTask("LowFrequencyTask", core.CyclicTask, 10, 1*time.Second)
	lowFreqTask.WithProgram(&core.Program{
		Name: "LowFreqCounter",
		Logic: func(now time.Time) {
			lowFreqCounter.Add(1)
			fmt.Printf("[%s] ---- LowFrequencyTask running (count: %d)\n", now.Format("15:04:05.000"), lowFreqCounter.Load())
		},
	})

	// Create a resource and add the tasks to it
	resource := &core.Resource{Name: "MainCPU", Cycle: 50 * time.Millisecond}
	resource.WithTask(highFreqTask).WithTask(lowFreqTask)

	// Add the resource to the main configuration
	config := &core.Configuration{Name: "MainConfig"}
	config.WithResource(resource)

	// --- 2. Start the PLC ---
	fmt.Println("Starting resource. Both tasks should be running.")
	resource.Start()
	time.Sleep(3200 * time.Millisecond)

	// --- 3. Dynamically Remove a Task ---
	fmt.Println("\n>>> Removing 'LowFrequencyTask' from the running resource...")
	if resource.RemoveTask("LowFrequencyTask") {
		fmt.Println(">>> 'LowFrequencyTask' successfully removed.")
	}
	fmt.Println()
	time.Sleep(3200 * time.Millisecond)

	// --- 4. Dynamically Remove the Resource ---
	fmt.Println("\n>>> Removing 'MainCPU' resource from the configuration...")
	config.RemoveResource("MainCPU")
	fmt.Println(">>> 'MainCPU' resource successfully removed and stopped.")
	fmt.Println("\nDynamic reconfiguration complete. Application will now exit.")
}
