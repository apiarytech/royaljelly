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

	if executionCount.Load() < 2 {
		t.Fatalf("Task was expected to run at least 2 times, but ran %d times", executionCount.Load())
	}

	// 1. Test ExecutionTime
	execTime := cyclicTask.ExecutionTime()
	if execTime <= programSleep {
		t.Errorf("ExecutionTime (%v) should be greater than the program's sleep time (%v)", execTime, programSleep)
	}
	// It should also be reasonably close, not excessively long.
	if execTime > programSleep*3 {
		t.Logf("Warning: ExecutionTime (%v) is much longer than program sleep time (%v)", execTime, programSleep)
	}
	t.Logf("Metric - ExecutionTime: %v", execTime)

	// 2. Test CycleTime (delta between last two runs)
	cycleTime := cyclicTask.CycleTime()
	if cycleTime <= 0 {
		t.Errorf("CycleTime should be a positive duration, but got %v", cycleTime)
	}
	// The cycle time should be close to the task's interval.
	// We allow for some deviation due to scheduler timing.
	expectedCycleTime := taskInterval
	minCycle := expectedCycleTime - resourceCycle*2
	maxCycle := expectedCycleTime + resourceCycle*2
	if cycleTime < minCycle || cycleTime > maxCycle {
		t.Errorf("CycleTime (%v) is outside the expected range [%v, %v]", cycleTime, minCycle, maxCycle)
	}
	t.Logf("Metric - CycleTime: %v", cycleTime)

	// 3. Test Drift
	drift := cyclicTask.Drift()
	// Drift is the difference between the actual run time and the scheduled run time.
	// It can be positive or negative but should be small.
	maxDrift := resourceCycle * 2 // Should not drift more than a couple of resource cycles.
	if drift > maxDrift || drift < -maxDrift {
		t.Errorf("Drift (%v) is larger than the expected maximum (%v)", drift, maxDrift)
	}
	t.Logf("Metric - Drift: %v", drift)
}

func TestTaskEnableDisable(t *testing.T) {
	resource := &Resource{Name: "EnableDisableCPU", Cycle: 10 * time.Millisecond}
	taskInterval := 30 * time.Millisecond
	var runCount atomic.Int32

	task := NewTask("SwitchableTask", CyclicTask, 1, taskInterval)
	task.AddProgram(&Program{
		Name: "Counter",
		Logic: func(now time.Time) {
			runCount.Add(1)
		},
	})
	resource.AddTask(task)
	resource.Start()
	defer resource.Stop()

	// 1. Wait for the task to run (it starts enabled by default).
	// Use robust polling to avoid flaky tests.
	waitTimeout := time.After(200 * time.Millisecond)
	initialRun := false
	for !initialRun {
		select {
		case <-waitTimeout:
			t.Fatalf("Timeout: Task did not run initially")
		case <-time.After(5 * time.Millisecond):
			if runCount.Load() > 0 {
				initialRun = true
			}
		}
	}
	t.Logf("Task ran %d time(s) initially.", runCount.Load())

	// 2. Disable the task and verify it stops running.
	task.Disable()
	currentRuns := runCount.Load()
	time.Sleep(taskInterval * 2) // Wait for a couple of cycles where it should have run.

	if runCount.Load() != currentRuns {
		t.Errorf("Task ran after being disabled. Expected %d runs, but got %d", currentRuns, runCount.Load())
	}
	t.Log("Task correctly stopped after being disabled.")

	// 3. Enable the task again and verify it resumes.
	task.Enable()
	waitTimeout = time.After(200 * time.Millisecond)
	resumedRun := false
	for !resumedRun {
		select {
		case <-waitTimeout:
			t.Fatalf("Timeout: Task did not resume after being enabled")
		case <-time.After(5 * time.Millisecond):
			if runCount.Load() > currentRuns {
				resumedRun = true
			}
		}
	}
	if runCount.Load() <= currentRuns {
		t.Errorf("Task did not resume after being enabled. Expected more than %d runs, got %d", currentRuns, runCount.Load())
	}
	t.Logf("Task correctly resumed after being enabled, running %d time(s).", runCount.Load())
}

func TestEventDrivenTaskTrigger(t *testing.T) {
	resource := &Resource{Name: "EventCPU", Cycle: 5 * time.Millisecond}
	var runCount atomic.Int32

	eventTask := NewTask("EventTask", EventDrivenTask, 1, 0)
	eventTask.AddProgram(&Program{
		Name: "EventCounter",
		Logic: func(now time.Time) {
			runCount.Add(1)
		},
	})

	resource.AddTask(eventTask)
	resource.Start()
	defer resource.Stop()

	// 1. Verify it doesn't run without a trigger.
	time.Sleep(50 * time.Millisecond)
	if runCount.Load() != 0 {
		t.Fatalf("Event task ran without a trigger. Count: %d", runCount.Load())
	}

	// 2. Trigger the task and verify it runs once.
	t.Log("Triggering event task...")
	if err := eventTask.Trigger(); err != nil {
		t.Fatalf("Trigger() failed: %v", err)
	}

	// Wait for the run count to become 1
	waitTimeout := time.After(100 * time.Millisecond)
	select {
	case <-waitTimeout:
		t.Fatalf("Timeout: Event task did not run after trigger. Count: %d", runCount.Load())
	case <-time.After(50 * time.Millisecond): // Give it ample time to execute
		if runCount.Load() != 1 {
			t.Errorf("Expected run count of 1 after trigger, but got %d", runCount.Load())
		}
	}
	t.Log("Event task ran successfully after trigger.")

	// 3. Test that triggering a non-event-driven task returns an error.
	cyclicTask := NewTask("CyclicForTriggerTest", CyclicTask, 1, 1*time.Second)
	err := cyclicTask.Trigger()
	if err == nil {
		t.Error("Expected an error when triggering a cyclic task, but got nil")
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

func TestBuilderMethods(t *testing.T) {
	t.Run("TestWithResource", func(t *testing.T) {
		config := &Configuration{Name: "TestConfig"}
		res1 := &Resource{Name: "R1"}
		res2 := &Resource{Name: "R2"}

		// Chain the calls
		config.WithResource(res1).WithResource(res2)

		if len(config.Resources) != 2 {
			t.Fatalf("Expected 2 resources after chaining WithResource, got %d", len(config.Resources))
		}
		if config.Resources[0] != res1 || config.Resources[1] != res2 {
			t.Error("WithResource did not add resources in the correct order")
		}
	})

	t.Run("TestWithTask", func(t *testing.T) {
		resource := &Resource{Name: "TestResource"}
		task1 := NewTask("T1", CyclicTask, 1, 1*time.Second)
		task2 := NewTask("T2", CyclicTask, 2, 1*time.Second)

		resource.WithTask(task1).WithTask(task2)

		if len(resource.Tasks) != 2 {
			t.Fatalf("Expected 2 tasks after chaining WithTask, got %d", len(resource.Tasks))
		}
		if resource.Tasks[0] != task1 || resource.Tasks[1] != task2 {
			t.Error("WithTask did not add tasks in the correct order")
		}
	})
}

func TestResourceStartIdempotency(t *testing.T) {
	resource := &Resource{Name: "IdempotentCPU", Cycle: 10 * time.Millisecond}
	var runCount atomic.Int32

	task := NewTask("IdempotentTask", CyclicTask, 1, 20*time.Millisecond)
	task.AddProgram(&Program{
		Name: "Counter",
		Logic: func(now time.Time) {
			runCount.Add(1)
		},
	})
	resource.AddTask(task)

	// Start the resource
	resource.Start()

	// Immediately try to start it again. This should be a no-op.
	resource.Start()

	time.Sleep(50 * time.Millisecond) // Let it run for a few cycles

	// Stop the resource. If Start() was not idempotent, this might hang or panic.
	resource.Stop()

	if runCount.Load() == 0 {
		t.Error("Task did not run, indicating the scheduler might not have started correctly.")
	}
	t.Logf("Task ran %d times. Calling Start() multiple times did not cause issues.", runCount.Load())
}
