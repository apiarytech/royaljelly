//go:build tinygo

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

// waitForCondition is a helper for TinyGo's cooperative scheduler.
// It yields control by sleeping for a very short duration, allowing other
// goroutines to run.
func waitForCondition(t *testing.T, timeout time.Duration, condition func() bool) {
	t.Helper()
	start := time.Now()
	for time.Since(start) < timeout {
		if condition() {
			return
		}
		// Yield to the scheduler to allow other goroutines to run.
		// A small sleep is a common way to do this in TinyGo.
		time.Sleep(1 * time.Millisecond)
	}
	t.Fatalf("Timeout waiting for condition")
}

func TestTaskEnableDisableTinyGo(t *testing.T) {
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
	waitForCondition(t, 200*time.Millisecond, func() bool {
		return runCount.Load() > 0
	})
	t.Logf("Task ran %d time(s) initially.", runCount.Load())

	// 2. Disable the task and verify it stops running.
	task.Disable()
	currentRuns := runCount.Load()
	// Wait for a couple of cycles where it should have run.
	// In a cooperative test, we can't just sleep; we must yield.
	time.Sleep(taskInterval * 2)

	if runCount.Load() != currentRuns {
		t.Errorf("Task ran after being disabled. Expected %d runs, but got %d", currentRuns, runCount.Load())
	}
	t.Log("Task correctly stopped after being disabled.")

	// 3. Enable the task again and verify it resumes.
	task.Enable()
	waitForCondition(t, 200*time.Millisecond, func() bool {
		return runCount.Load() > currentRuns
	})

	if runCount.Load() <= currentRuns {
		t.Errorf("Task did not resume after being enabled. Expected more than %d runs, got %d", currentRuns, runCount.Load())
	}
	t.Logf("Task correctly resumed after being enabled, running %d time(s).", runCount.Load())
}

func TestEventDrivenTaskTriggerTinyGo(t *testing.T) {
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

	// The existing test logic for event-driven tasks is mostly cooperative
	// and should work, but we place it here for consistency.
	time.Sleep(50 * time.Millisecond) // Yield to ensure it doesn't run
	if runCount.Load() != 0 {
		t.Fatalf("Event task ran without a trigger. Count: %d", runCount.Load())
	}

	if err := eventTask.Trigger(); err != nil {
		t.Fatalf("Trigger() failed: %v", err)
	}

	waitForCondition(t, 100*time.Millisecond, func() bool {
		return runCount.Load() == 1
	})

	if runCount.Load() != 1 {
		t.Errorf("Expected run count of 1 after trigger, but got %d", runCount.Load())
	}
}
