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
	"fmt"
	"sync"
	"time"
)

// DefaultCycle is the scan cycle used when Resource.Cycle is zero.
const DefaultCycle = time.Millisecond

// Resource represents a processing unit within the configuration, like a CPU.
// It runs one scheduler goroutine that executes its tasks in priority order.
//
// The exported fields are configuration. Set them before calling Start.
type Resource struct {
	Name     string
	Cycle    time.Duration // Scan cycle of the scheduler. Zero means DefaultCycle.
	Affinity int           // Optional: 1-based CPU core to pin the scheduler to. 0 or less means none.
	OnFault  FaultHandler  // Optional: receives faults. Nil uses Configuration.OnFault or standard error.

	mu      sync.Mutex
	tasks   []*Task // Sorted by priority. Copy-on-write: never modified in place once published.
	running bool
	stop    chan struct{} // Closed to stop the current scheduler goroutine.
	done    chan struct{} // Closed when the current scheduler goroutine has exited.
}

// AddTask adds a task to the resource and keeps the task list sorted by priority.
// It is safe to call on a running resource; the task is scheduled from the next scan.
func (r *Resource) AddTask(t *Task) {
	t.setOwner(r)
	r.mu.Lock()
	defer r.mu.Unlock()
	next := make([]*Task, len(r.tasks), len(r.tasks)+1)
	copy(next, r.tasks)
	next = append(next, t)
	SortTasks(next)
	r.tasks = next
}

// WithTask adds a task to the resource and returns the resource for chaining.
func (r *Resource) WithTask(t *Task) *Resource {
	r.AddTask(t)
	return r
}

// RemoveTask removes the first task with the given name.
// It returns true if the task was found and removed. It is safe to call on a
// running resource; a run of that task already in progress completes normally.
func (r *Resource) RemoveTask(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, t := range r.tasks {
		if t.Name == name {
			next := make([]*Task, 0, len(r.tasks)-1)
			next = append(next, r.tasks[:i]...)
			r.tasks = append(next, r.tasks[i+1:]...)
			return true
		}
	}
	return false
}

// FindTask returns the first task with the given name, or nil.
func (r *Resource) FindTask(name string) *Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.tasks {
		if t.Name == name {
			return t
		}
	}
	return nil
}

// Tasks returns a copy of the resource's tasks in priority order.
func (r *Resource) Tasks() []*Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*Task, len(r.tasks))
	copy(out, r.tasks)
	return out
}

// IsRunning reports whether the scheduler is running.
func (r *Resource) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.running
}

// Validate checks the resource and its tasks for configuration errors. It is
// called by Start, and returns all problems found joined into one error.
func (r *Resource) Validate() error {
	r.mu.Lock()
	cycle := r.Cycle
	tasks := r.tasks
	r.mu.Unlock()
	if cycle == 0 {
		cycle = DefaultCycle
	}

	var errs []error
	fail := func(format string, args ...any) {
		errs = append(errs, fmt.Errorf("resource '%s': "+format, append([]any{r.Name}, args...)...))
	}
	if cycle < 0 {
		fail("cycle must not be negative, got %v", cycle)
	}
	names := make(map[string]bool, len(tasks))
	priorities := make(map[int]string, len(tasks))
	for _, t := range tasks {
		if names[t.Name] {
			fail("duplicate task name '%s'", t.Name)
		}
		names[t.Name] = true
		if other, ok := priorities[t.Priority]; ok {
			fail("duplicate task priority %d (tasks '%s' and '%s')", t.Priority, other, t.Name)
		} else {
			priorities[t.Priority] = t.Name
		}
		switch t.Type {
		case CyclicTask:
			if t.Interval <= 0 {
				fail("cyclic task '%s' needs a positive interval, got %v", t.Name, t.Interval)
			} else if cycle >= t.Interval {
				fail("cycle time (%v) must be faster than task '%s' interval (%v)", cycle, t.Name, t.Interval)
			}
		case EventDrivenTask:
		default:
			fail("task '%s' has unknown type %v", t.Name, t.Type)
		}
		if t.Watchdog < 0 {
			fail("task '%s' watchdog must not be negative, got %v", t.Name, t.Watchdog)
		}
		progNames := make(map[string]bool)
		for _, p := range t.Programs() {
			if p == nil {
				fail("task '%s' contains a nil program", t.Name)
				continue
			}
			if progNames[p.Name] {
				fail("task '%s' has duplicate program name '%s'", t.Name, p.Name)
			}
			progNames[p.Name] = true
		}
	}
	return errors.Join(errs...)
}

// Start validates the resource and starts its scheduler goroutine. It returns an
// error if validation fails or the requested CPU affinity cannot be applied.
// Calling Start on a running resource does nothing and returns nil.
func (r *Resource) Start() error {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return nil
	}
	if r.Cycle == 0 {
		r.Cycle = DefaultCycle
	}
	r.mu.Unlock()

	if err := r.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return nil
	}
	stop, done := make(chan struct{}), make(chan struct{})
	r.stop, r.done, r.running = stop, done, true
	cycle, affinity := r.Cycle, r.Affinity
	r.mu.Unlock()

	ready := make(chan error, 1)
	go r.run(cycle, affinity, stop, done, ready)
	if err := <-ready; err != nil {
		<-done
		r.mu.Lock()
		if r.done == done {
			r.running, r.stop, r.done = false, nil, nil
		}
		r.mu.Unlock()
		return err
	}
	return nil
}

// run is the scheduler goroutine.
func (r *Resource) run(cycle time.Duration, affinity int, stop, done chan struct{}, ready chan<- error) {
	defer close(done)
	if affinity > 0 {
		if err := pinThread(affinity - 1); err != nil {
			ready <- fmt.Errorf("resource '%s': cannot pin to CPU %d: %w", r.Name, affinity, err)
			return
		}
	}
	ready <- nil

	ticker := time.NewTicker(cycle)
	defer ticker.Stop()
	for {
		select {
		case now := <-ticker.C:
			r.scan(now)
		case <-stop:
			return
		}
	}
}

// scan runs every due task once, in priority order. It does not allocate in the
// steady state: task and program lists are copy-on-write, so the scan only
// reads the current slice headers under the locks.
func (r *Resource) scan(now time.Time) {
	r.mu.Lock()
	tasks := r.tasks
	r.mu.Unlock()

	for _, t := range tasks {
		t.mu.Lock()
		run, missed := t.due(now)
		programs := t.programs
		t.mu.Unlock()

		if missed > 0 {
			r.taskFault(t, Fault{
				Kind: FaultOverrun,
				Task: t.Name,
				Time: now,
				Err:  fmt.Errorf("%w: %d scheduled run(s) skipped", ErrOverrun, missed),
			})
		}
		if !run {
			continue
		}

		start := time.Now()
		t.armWatchdog()
		for _, p := range programs {
			r.execute(t, p, now)
		}
		t.disarmWatchdog()
		end := time.Now()

		t.mu.Lock()
		t.finish(now, start, end)
		t.mu.Unlock()
	}
}

// execute runs one program and converts a panic into a fault.
func (r *Resource) execute(t *Task, p *Program, now time.Time) {
	defer func() {
		if rec := recover(); rec != nil {
			r.taskFault(t, Fault{
				Kind:    FaultPanic,
				Task:    t.Name,
				Program: p.Name,
				Time:    now,
				Value:   rec,
				Err:     fmt.Errorf("%w: %v", ErrProgramPanic, rec),
			})
		}
	}()
	p.Execute(now)
}

// taskFault records a fault of task t in its stats and delivers it.
func (r *Resource) taskFault(t *Task, f Fault) {
	f.Resource = r.Name
	t.recordFault(f)
	r.fault(f)
}

// fault delivers a fault to the configured handler.
func (r *Resource) fault(f Fault) {
	f.Resource = r.Name
	if h := r.OnFault; h != nil {
		h(f)
		return
	}
	defaultFaultHandler(f)
}

// Stop terminates the scheduler and waits for it to exit. It is safe to call
// more than once. Do not call Stop from inside a program running on this resource.
func (r *Resource) Stop() {
	r.mu.Lock()
	if !r.running {
		r.mu.Unlock()
		return
	}
	stop, done := r.stop, r.done
	r.running, r.stop, r.done = false, nil, nil
	r.mu.Unlock()

	close(stop)
	<-done
}

// Pause disables all tasks of the resource without stopping the scheduler.
func (r *Resource) Pause() {
	for _, t := range r.Tasks() {
		t.Disable()
	}
}

// Resume enables all tasks of the resource.
func (r *Resource) Resume() {
	for _, t := range r.Tasks() {
		t.Enable()
	}
}
