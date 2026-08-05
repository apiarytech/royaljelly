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

	// --- Robust waiting logic ---
	// Instead of a fixed sleep, poll until the desired state is reached or a timeout occurs.
	// This makes the test resilient to scheduler jitter on different environments like GitHub Actions.
	timeout := time.After(200 * time.Millisecond)
	for runCount.Load() < 2 {
		select {
		case <-timeout:
			t.Fatalf("Timeout: Expected 2 runs, but got %d", runCount.Load())
		case <-time.After(5 * time.Millisecond):
			// Poll every 5ms
		}
	}

	if runCount.Load() != 2 {
		t.Errorf("Expected 2 runs, but got %d", runCount.Load())
	}

	task.Disable()
	// Wait a bit to ensure the disabled task doesn't run again.
	time.Sleep(55 * time.Millisecond)
	resource.Stop()

	if runCount.Load() > 2 {
		t.Errorf("Task ran after being disabled. Expected 2 total runs, but got %d", runCount.Load())
	}
}

func TestRemoveFunctions(t *testing.T) {
	t.Run("TestRemoveProgram", func(t *testing.T) {
		task := NewTask("TestTask", CyclicTask, 1, 1*time.Second)
		prog1 := &Program{Name: "Prog1"}
		prog2 := &Program{Name: "Prog2"}
		task.AddProgram(prog1)
		task.AddProgram(prog2)

		if len(task.Programs) != 2 {
			t.Fatalf("Expected 2 programs initially, got %d", len(task.Programs))
		}

		// Test removing an existing program
		removed := task.RemoveProgram("Prog1")
		if !removed {
			t.Error("Expected RemoveProgram to return true for an existing program")
		}
		if len(task.Programs) != 1 {
			t.Errorf("Expected 1 program after removal, got %d", len(task.Programs))
		}
		if task.Programs[0].Name != "Prog2" {
			t.Errorf("Incorrect program remained. Expected 'Prog2', got '%s'", task.Programs[0].Name)
		}

		// Test removing a non-existent program
		removed = task.RemoveProgram("NonExistentProg")
		if removed {
			t.Error("Expected RemoveProgram to return false for a non-existent program")
		}
		if len(task.Programs) != 1 {
			t.Errorf("Program count should not change when removing a non-existent program. Got %d", len(task.Programs))
		}
	})

	t.Run("TestRemoveTask", func(t *testing.T) {
		resource := &Resource{Name: "TestResource"}
		task1 := NewTask("Task1", CyclicTask, 1, 1*time.Second)
		task2 := NewTask("Task2", CyclicTask, 2, 1*time.Second)
		resource.AddTask(task1)
		resource.AddTask(task2)

		if len(resource.Tasks) != 2 {
			t.Fatalf("Expected 2 tasks initially, got %d", len(resource.Tasks))
		}

		// Test removing an existing task
		removed := resource.RemoveTask("Task1")
		if !removed {
			t.Error("Expected RemoveTask to return true for an existing task")
		}
		if len(resource.Tasks) != 1 {
			t.Errorf("Expected 1 task after removal, got %d", len(resource.Tasks))
		}
		if resource.Tasks[0].Name != "Task2" {
			t.Errorf("Incorrect task remained. Expected 'Task2', got '%s'", resource.Tasks[0].Name)
		}

		// Test removing a non-existent task
		removed = resource.RemoveTask("NonExistentTask")
		if removed {
			t.Error("Expected RemoveTask to return false for a non-existent task")
		}
		if len(resource.Tasks) != 1 {
			t.Errorf("Task count should not change when removing a non-existent task. Got %d", len(resource.Tasks))
		}
	})

	t.Run("TestRemoveResource", func(t *testing.T) {
		config := &Configuration{Name: "TestConfig"}
		res1 := &Resource{Name: "Resource1", Cycle: 10 * time.Millisecond}
		res2 := &Resource{Name: "Resource2", Cycle: 10 * time.Millisecond}
		config.WithResource(res1).WithResource(res2)

		// Start one of the resources to test the stop functionality
		res1.Start()
		time.Sleep(20 * time.Millisecond) // Let it run for a moment

		if len(config.Resources) != 2 {
			t.Fatalf("Expected 2 resources initially, got %d", len(config.Resources))
		}

		// Test removing an existing (and running) resource
		removed := config.RemoveResource("Resource1")
		if !removed {
			t.Error("Expected RemoveResource to return true for an existing resource")
		}
		if len(config.Resources) != 1 {
			t.Errorf("Expected 1 resource after removal, got %d", len(config.Resources))
		}
		if config.Resources[0].Name != "Resource2" {
			t.Errorf("Incorrect resource remained. Expected 'Resource2', got '%s'", config.Resources[0].Name)
		}

		// Verify the removed resource is no longer running
		if res1.running {
			t.Error("Removed resource should have been stopped, but its 'running' flag is still true")
		} else {
			// As an extra check, try to stop it again to ensure it doesn't panic or deadlock.
			// This tests the idempotency of the Stop() method.
			res1.Stop()
		}
	})
}

func TestFindFunctions(t *testing.T) {
	// Setup a hierarchy: Config -> Resource -> Task -> Program
	config := &Configuration{Name: "TestConfig"}
	res1 := &Resource{Name: "Resource1", Cycle: 10 * time.Millisecond}
	res2 := &Resource{Name: "Resource2", Cycle: 10 * time.Millisecond}
	config.WithResource(res1).WithResource(res2)

	task1 := NewTask("Task1", CyclicTask, 1, 1*time.Second)
	task2 := NewTask("Task2", CyclicTask, 2, 1*time.Second)
	res1.WithTask(task1).WithTask(task2)

	prog1 := &Program{Name: "Prog1"}
	prog2 := &Program{Name: "Prog2"}
	task1.WithProgram(prog1).WithProgram(prog2)

	t.Run("TestFindResource", func(t *testing.T) {
		// Test finding an existing resource
		foundRes := config.FindResource("Resource1")
		if foundRes == nil {
			t.Error("Expected to find Resource1, but got nil")
		}
		if foundRes != res1 {
			t.Errorf("Found incorrect resource. Expected %p, got %p", res1, foundRes)
		}

		// Test finding a non-existent resource
		foundRes = config.FindResource("NonExistentResource")
		if foundRes != nil {
			t.Errorf("Expected to find nil for non-existent resource, but got %p", foundRes)
		}
	})

	t.Run("TestFindTask", func(t *testing.T) {
		// Test finding an existing task
		foundTask := res1.FindTask("Task1")
		if foundTask == nil {
			t.Error("Expected to find Task1 in Resource1, but got nil")
		}
		if foundTask != task1 {
			t.Errorf("Found incorrect task. Expected %p, got %p", task1, foundTask)
		}

		// Test finding a non-existent task in res1
		foundTask = res1.FindTask("NonExistentTask")
		if foundTask != nil {
			t.Errorf("Expected to find nil for non-existent task, but got %p", foundTask)
		}

		// Test finding a task in the wrong resource
		foundTask = res2.FindTask("Task1")
		if foundTask != nil {
			t.Errorf("Expected to find nil for Task1 in Resource2, but got %p", foundTask)
		}
	})

	t.Run("TestFindProgram", func(t *testing.T) {
		// Test finding an existing program
		foundProg := task1.FindProgram("Prog1")
		if foundProg == nil {
			t.Error("Expected to find Prog1 in Task1, but got nil")
		}
		if foundProg != prog1 {
			t.Errorf("Found incorrect program. Expected %p, got %p", prog1, foundProg)
		}

		// Test finding a non-existent program in task1
		foundProg = task1.FindProgram("NonExistentProg")
		if foundProg != nil {
			t.Errorf("Expected to find nil for non-existent program, but got %p", foundProg)
		}

		// Test finding a program in the wrong task (assuming task2 has no programs initially)
		foundProg = task2.FindProgram("Prog1")
		if foundProg != nil {
			t.Errorf("Expected to find nil for Prog1 in Task2, but got %p", foundProg)
		}
	})
}
