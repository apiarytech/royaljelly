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

package core

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestResourceSchedulerWithMetrics(t *testing.T) {
	resourceCycle := 10 * time.Millisecond
	taskInterval := 70 * time.Millisecond
	programSleep := 5 * time.Millisecond

	resource := &Resource{Name: "TestCPU", Cycle: resourceCycle}

	var executionCount atomic.Int32

	cyclicTask := NewTask("CyclicMetricsTask", CyclicTask, 10, taskInterval)
	cyclicTask.AddProgram(&Program{
		Name: "Workload",
		Logic: func(now time.Time) {
			time.Sleep(programSleep)
			executionCount.Add(1)
		},
	})

	resource.AddTask(cyclicTask)

	// Start the resource scheduler
	resource.Start()

	// Let the scheduler run long enough for the task to execute at least twice
	// to ensure CycleTime and Drift are calculated.
	time.Sleep(taskInterval*2 + resourceCycle*5) // e.g., 50ms * 2 + 10ms * 5 = 150ms

	// Stop the scheduler to safely inspect the metrics
	resource.Stop()

	// --- Assertions ---

	// Lock the task to safely read its metrics
	cyclicTask.mu.RLock()
	defer cyclicTask.mu.RUnlock()

	if executionCount.Load() < 2 {
		t.Fatalf("Task was expected to run at least 2 times, but ran %d times", executionCount.Load())
	}

	// 1. Test ExecutionTime
	if cyclicTask.ExecutionTime <= programSleep {
		t.Errorf("ExecutionTime (%v) should be greater than the program's sleep time (%v)", cyclicTask.ExecutionTime, programSleep)
	}
	// It should also be reasonably close, not excessively long.
	if cyclicTask.ExecutionTime > programSleep*3 {
		t.Logf("Warning: ExecutionTime (%v) is much longer than program sleep time (%v)", cyclicTask.ExecutionTime, programSleep)
	}
	t.Logf("Metric - ExecutionTime: %v", cyclicTask.ExecutionTime)

	// 2. Test CycleTime (delta between last two runs)
	if cyclicTask.CycleTime <= 0 {
		t.Errorf("CycleTime should be a positive duration, but got %v", cyclicTask.CycleTime)
	}
	// The cycle time should be close to the task's interval.
	// We allow for some deviation due to scheduler timing.
	expectedCycleTime := taskInterval
	minCycle := expectedCycleTime - resourceCycle*2
	maxCycle := expectedCycleTime + resourceCycle*2
	if cyclicTask.CycleTime < minCycle || cyclicTask.CycleTime > maxCycle {
		t.Errorf("CycleTime (%v) is outside the expected range [%v, %v]", cyclicTask.CycleTime, minCycle, maxCycle)
	}
	t.Logf("Metric - CycleTime: %v", cyclicTask.CycleTime)

	// 3. Test Drift
	// Drift is the difference between the actual run time and the scheduled run time.
	// It can be positive or negative but should be small.
	maxDrift := resourceCycle * 2 // Should not drift more than a couple of resource cycles.
	if cyclicTask.Drift > maxDrift || cyclicTask.Drift < -maxDrift {
		t.Errorf("Drift (%v) is larger than the expected maximum (%v)", cyclicTask.Drift, maxDrift)
	}
	t.Logf("Metric - Drift: %v", cyclicTask.Drift)
}

func TestTaskEnableDisable(t *testing.T) {
	resource := &Resource{Name: "EnableDisableCPU", Cycle: 10 * time.Millisecond}
	var runCount atomic.Int32

	task := NewTask("SwitchableTask", CyclicTask, 1, 20*time.Millisecond)
	task.AddProgram(&Program{
		Name: "Counter",
		Logic: func(now time.Time) {
			runCount.Add(1)
		},
	})

	resource.AddTask(task)
	resource.Start()

	time.Sleep(45 * time.Millisecond) // Should run 2 times (at ~20ms, ~40ms). Check before the 3rd run at ~60ms.
	if runCount.Load() != 2 {
		t.Errorf("Expected 2 runs, but got %d", runCount.Load())
	}

	task.Disable()
	fmt.Println("Task disabled")
	time.Sleep(55 * time.Millisecond) // Should not run anymore
	resource.Stop()

	if runCount.Load() != 2 {
		t.Errorf("Task ran after being disabled. Expected 2 total runs, but got %d", runCount.Load())
	}
}
