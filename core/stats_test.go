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

package core

import (
	"errors"
	"testing"
	"time"
)

// TestExecutionTimeStats checks the shortest, longest and mean run times.
func TestExecutionTimeStats(t *testing.T) {
	sleeps := []time.Duration{2 * time.Millisecond, 20 * time.Millisecond, 8 * time.Millisecond}
	run := 0
	task := NewTask("T", CyclicTask, 1, time.Second)
	task.AddProgram(&Program{Name: "P", Logic: func(time.Time) { time.Sleep(sleeps[run]); run++ }})
	r := (&Resource{Name: "R"}).WithTask(task)
	if s := task.Stats(); s.MinExecutionTime != 0 || s.MaxExecutionTime != 0 || s.AvgExecutionTime != 0 {
		t.Fatalf("before any run: %+v", s)
	}
	scanAt(r, 0, time.Second, 2*time.Second)

	s := task.Stats()
	if s.Runs != 3 {
		t.Fatalf("Runs = %d, want 3", s.Runs)
	}
	if s.MinExecutionTime < 2*time.Millisecond || s.MinExecutionTime >= 8*time.Millisecond {
		t.Errorf("MinExecutionTime = %v, want the 2 ms run", s.MinExecutionTime)
	}
	if s.MaxExecutionTime < 20*time.Millisecond {
		t.Errorf("MaxExecutionTime = %v, want at least the 20 ms run", s.MaxExecutionTime)
	}
	if s.AvgExecutionTime <= s.MinExecutionTime || s.AvgExecutionTime >= s.MaxExecutionTime {
		t.Errorf("AvgExecutionTime = %v, want between %v and %v", s.AvgExecutionTime, s.MinExecutionTime, s.MaxExecutionTime)
	}
	if s.ExecutionTime < 8*time.Millisecond || s.ExecutionTime >= 20*time.Millisecond {
		t.Errorf("ExecutionTime = %v, want the last (8 ms) run", s.ExecutionTime)
	}

	task.ResetStats()
	if s := task.Stats(); s != (TaskStats{}) {
		t.Errorf("after ResetStats: %+v", s)
	}
}

// TestFaultStats checks that each kind of fault is counted and kept as the
// last fault with its time.
func TestFaultStats(t *testing.T) {
	task := NewTask("T", CyclicTask, 1, time.Second)
	task.AddProgram(&Program{Name: "Bad", Logic: func(time.Time) { panic("boom") }})
	r := (&Resource{Name: "R", OnFault: func(Fault) {}}).WithTask(task)
	if s := task.Stats(); s.Faults != 0 || s.LastFault != nil || !s.LastFaultTime.IsZero() {
		t.Fatalf("before any fault: %+v", s)
	}

	scanAt(r, 0)
	s := task.Stats()
	if s.Faults != 1 || s.LastFault == nil || s.LastFault.Kind != FaultPanic || !errors.Is(s.LastFault, ErrProgramPanic) {
		t.Fatalf("after a panic: %+v", s)
	}
	if s.LastFault.Program != "Bad" || s.LastFault.Resource != "R" || !s.LastFaultTime.Equal(t0) {
		t.Errorf("last fault %+v at %v, want program Bad on R at the scan time", s.LastFault, s.LastFaultTime)
	}

	// Scheduled runs missed: the overrun is reported before the run, whose
	// panic is then the last fault.
	scanAt(r, 3*time.Second+time.Millisecond)
	s = task.Stats()
	if s.Faults != 3 || s.Overruns == 0 || s.LastFault.Kind != FaultPanic || !s.LastFaultTime.Equal(t0.Add(3*time.Second+time.Millisecond)) {
		t.Errorf("after an overrun and a second panic: %+v", s)
	}

	// The returned fault is a copy: changing it does not change the task's.
	s.LastFault.Task = "changed"
	if task.Stats().LastFault.Task != "T" {
		t.Error("Stats returned the task's own fault, not a copy")
	}

	task.ResetStats()
	if s := task.Stats(); s.Faults != 0 || s.LastFault != nil || !s.LastFaultTime.IsZero() {
		t.Errorf("after ResetStats: %+v", s)
	}
}

// TestWatchdogFaultStats checks that a watchdog expiration is the last fault.
func TestWatchdogFaultStats(t *testing.T) {
	task := NewTask("Slow", CyclicTask, 1, time.Second)
	task.Watchdog = 5 * time.Millisecond
	task.AddProgram(&Program{Name: "P", Logic: func(time.Time) { time.Sleep(50 * time.Millisecond) }})
	done := make(chan struct{}, 1)
	r := (&Resource{Name: "R", OnFault: func(Fault) { done <- struct{}{} }}).WithTask(task)
	before := time.Now()
	scanAt(r, 0)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("watchdog did not fire")
	}
	s := task.Stats()
	if s.Faults != 1 || s.LastFault == nil || s.LastFault.Kind != FaultWatchdog || s.LastFaultTime.Before(before) {
		t.Errorf("after a watchdog expiration: %+v", s)
	}
}
