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
	"sort"
	"sync"
	"time"
)

// POU (Program Organization Unit) defines the interface for executable logic blocks.
// Programs and Function Blocks are considered POUs.
type POU interface {
	Execute(now time.Time)
}

// Program represents a unit of control logic that is scheduled by a Task.
// Logic typically instantiates and calls Function Blocks.
type Program struct {
	Name     string
	Logic    func(now time.Time) // The user-defined logic for the program.
	InitFunc func()              // Optional: logic that resets the program's state.
}

// Execute runs the program's defined logic.
func (p *Program) Execute(now time.Time) {
	if p.Logic != nil {
		p.Logic(now)
	}
}

// Init runs the program's initialization logic, if defined.
func (p *Program) Init() {
	if p.InitFunc != nil {
		p.InitFunc()
	}
}

// TaskType defines the scheduling mechanism for a task.
type TaskType int

const (
	// CyclicTask runs at a fixed interval.
	CyclicTask TaskType = iota
	// EventDrivenTask runs when triggered by Task.Trigger.
	EventDrivenTask
)

// String returns the name of the task type.
func (t TaskType) String() string {
	switch t {
	case CyclicTask:
		return "Cyclic"
	case EventDrivenTask:
		return "EventDriven"
	default:
		return fmt.Sprintf("TaskType(%d)", int(t))
	}
}

// maxPendingEvents bounds how many triggers an event-driven task queues.
const maxPendingEvents = 10

// Task controls the execution of one or more Programs.
//
// The exported fields are configuration. Set them before the task's resource
// is started; the scheduler reads them without locking. A zero Task is enabled.
type Task struct {
	Name     string
	Type     TaskType
	Priority int           // Lower number means higher priority (runs first in a scan).
	Interval time.Duration // Period of a cyclic task. Must exceed the resource cycle.
	Watchdog time.Duration // Optional: report FaultWatchdog if one run takes longer than this.

	mu       sync.Mutex
	disabled bool
	programs []*Program // Copy-on-write: never modified in place once published.
	pending  int        // Queued triggers for an event-driven task.
	nextRun  time.Time  // Next scheduled run of a cyclic task; zero means "run at next scan".
	lastRun  time.Time  // Scan time of the previous run.
	owner    *Resource  // Resource the task was added to, for fault reporting.

	// Runtime metrics, protected by mu.
	executionTime time.Duration
	cycleTime     time.Duration
	drift         time.Duration
	runs          uint64
	overruns      uint64
	watchdogTrips uint64

	// Watchdog timer; only touched by the scheduler goroutine.
	wdTimer *time.Timer
}

// TaskStats is a consistent snapshot of a task's runtime metrics.
type TaskStats struct {
	Runs          uint64        // Completed runs.
	Overruns      uint64        // Scheduled runs skipped because the task fell behind.
	WatchdogTrips uint64        // Runs that exceeded Task.Watchdog.
	ExecutionTime time.Duration // Duration of the last run.
	CycleTime     time.Duration // Time between the start of the last two runs.
	Drift         time.Duration // How late the last run started relative to its schedule.
}

// NewTask creates an enabled task.
func NewTask(name string, taskType TaskType, priority int, interval time.Duration) *Task {
	return &Task{Name: name, Type: taskType, Priority: priority, Interval: interval}
}

// AddProgram appends a program to the task. It is safe to call on a running task;
// the program runs from the next scan on.
func (t *Task) AddProgram(p *Program) {
	t.mu.Lock()
	defer t.mu.Unlock()
	// The full slice expression forces a new backing array, so a scan that is
	// iterating the previous slice is never affected.
	t.programs = append(t.programs[:len(t.programs):len(t.programs)], p)
}

// WithProgram adds a program to the task and returns the task for chaining.
func (t *Task) WithProgram(p *Program) *Task {
	t.AddProgram(p)
	return t
}

// RemoveProgram removes the first program with the given name.
// It returns true if the program was found and removed.
func (t *Task) RemoveProgram(name string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	for i, p := range t.programs {
		if p.Name == name {
			next := make([]*Program, 0, len(t.programs)-1)
			next = append(next, t.programs[:i]...)
			t.programs = append(next, t.programs[i+1:]...)
			return true
		}
	}
	return false
}

// FindProgram returns the first program with the given name, or nil.
func (t *Task) FindProgram(name string) *Program {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, p := range t.programs {
		if p.Name == name {
			return p
		}
	}
	return nil
}

// Programs returns a copy of the task's program list.
func (t *Task) Programs() []*Program {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]*Program, len(t.programs))
	copy(out, t.programs)
	return out
}

// Enable allows the scheduler to execute this task. A cyclic task that is
// re-enabled runs at the next scan and then resumes its interval.
func (t *Task) Enable() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.disabled {
		t.disabled = false
		t.nextRun = time.Time{}
	}
}

// Disable prevents the scheduler from executing this task.
func (t *Task) Disable() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.disabled = true
}

// IsEnabled reports whether the task may be scheduled.
func (t *Task) IsEnabled() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return !t.disabled
}

// Reset calls Init on every program of the task and restarts its schedule, so a
// cyclic task runs at the next scan. Init runs on the caller's goroutine; if the
// resource is running, the program's state must tolerate concurrent access.
func (t *Task) Reset() {
	t.mu.Lock()
	t.nextRun = time.Time{}
	programs := t.programs
	t.mu.Unlock()
	for _, p := range programs {
		p.Init()
	}
}

// Trigger requests one run of an event-driven task. Triggers beyond the queue
// limit are dropped. It returns an error for a cyclic task.
func (t *Task) Trigger() error {
	if t.Type != EventDrivenTask {
		return fmt.Errorf("cannot trigger non-event-driven task '%s'", t.Name)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.pending < maxPendingEvents {
		t.pending++
	}
	return nil
}

// ExecutionTime returns the duration of the last run of the task's programs.
func (t *Task) ExecutionTime() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.executionTime
}

// CycleTime returns the time between the starts of the last two runs.
func (t *Task) CycleTime() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.cycleTime
}

// Drift returns how late the last run of a cyclic task started relative to its
// scheduled time. It is zero for the first run and for event-driven tasks.
func (t *Task) Drift() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.drift
}

// Overruns returns how many scheduled runs were skipped because the task fell behind.
func (t *Task) Overruns() uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.overruns
}

// Stats returns a consistent snapshot of all runtime metrics.
func (t *Task) Stats() TaskStats {
	t.mu.Lock()
	defer t.mu.Unlock()
	return TaskStats{
		Runs:          t.runs,
		Overruns:      t.overruns,
		WatchdogTrips: t.watchdogTrips,
		ExecutionTime: t.executionTime,
		CycleTime:     t.cycleTime,
		Drift:         t.drift,
	}
}

// due decides whether the task runs in the scan at time now and advances its
// schedule. It returns the number of scheduled runs that were skipped because
// the task fell behind. It must be called with t.mu held.
func (t *Task) due(now time.Time) (run bool, missed uint64) {
	if t.disabled {
		return false, 0
	}
	switch t.Type {
	case CyclicTask:
		if t.Interval <= 0 {
			// Invalid interval (rejected by Validate); run every scan rather than divide by zero.
			t.drift = 0
			return true, 0
		}
		if t.nextRun.IsZero() {
			t.drift = 0
			t.nextRun = now.Add(t.Interval)
			return true, 0
		}
		if now.Before(t.nextRun) {
			return false, 0
		}
		t.drift = now.Sub(t.nextRun)
		t.nextRun = t.nextRun.Add(t.Interval)
		if !t.nextRun.After(now) {
			missed = uint64(now.Sub(t.nextRun)/t.Interval) + 1
			t.nextRun = t.nextRun.Add(time.Duration(missed) * t.Interval)
			t.overruns += missed
		}
		return true, missed
	case EventDrivenTask:
		if t.pending > 0 {
			t.pending--
			t.drift = 0
			return true, 0
		}
	}
	return false, 0
}

// finish records metrics for a completed run. It must be called with t.mu held.
func (t *Task) finish(scan, start, end time.Time) {
	if !t.lastRun.IsZero() {
		t.cycleTime = scan.Sub(t.lastRun)
	}
	t.lastRun = scan
	t.executionTime = end.Sub(start)
	t.runs++
}

// armWatchdog starts the watchdog timer for one run. Scheduler goroutine only.
func (t *Task) armWatchdog() {
	if t.Watchdog <= 0 {
		return
	}
	if t.wdTimer == nil {
		t.wdTimer = time.AfterFunc(t.Watchdog, t.watchdogExpired)
		return
	}
	t.wdTimer.Reset(t.Watchdog)
}

// disarmWatchdog stops the watchdog timer after a run. Scheduler goroutine only.
func (t *Task) disarmWatchdog() {
	if t.wdTimer != nil {
		t.wdTimer.Stop()
	}
}

// watchdogExpired runs on the timer goroutine when a run exceeds Task.Watchdog.
func (t *Task) watchdogExpired() {
	t.mu.Lock()
	t.watchdogTrips++
	owner := t.owner
	t.mu.Unlock()
	if owner != nil {
		owner.fault(Fault{
			Kind: FaultWatchdog,
			Task: t.Name,
			Time: time.Now(),
			Err:  fmt.Errorf("%w: run exceeded %v", ErrWatchdog, t.Watchdog),
		})
	}
}

func (t *Task) setOwner(r *Resource) {
	t.mu.Lock()
	t.owner = r
	t.mu.Unlock()
}

// taskSorter implements sort.Interface for a slice of tasks without reflection.
type taskSorter []*Task

func (s taskSorter) Len() int           { return len(s) }
func (s taskSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s taskSorter) Less(i, j int) bool { return s[i].Priority < s[j].Priority }

// SortTasks sorts tasks by priority, keeping insertion order for equal priorities.
func SortTasks(tasks []*Task) {
	sort.Stable(taskSorter(tasks))
}
