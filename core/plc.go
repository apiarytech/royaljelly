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

// Program represents a collection of logic (POUs) that is scheduled by a Task.
// In a real-world scenario, the `Logic` function would contain the user's
// control code written in Go, instantiating and calling Function Blocks.
type Program struct {
	Name     string
	Logic    func(now time.Time) // The user-defined logic for the program.
	InitFunc func()              // Optional: Logic to reset the program's state.
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
	// Cyclic tasks run at a fixed interval.
	CyclicTask TaskType = iota
	// EventDriven tasks run when triggered by an event.
	EventDrivenTask
)

// Task controls the execution of one or more Programs.
type Task struct {
	Name     string
	Type     TaskType
	Priority int           // Lower number means higher priority.
	Interval time.Duration // For Cyclic tasks.
	Programs []*Program
	Enabled  bool // If false, the task will not be scheduled.

	// --- Runtime Metrics ---
	executionTime time.Duration // Last execution time of the task's programs.
	cycleTime     time.Duration // Time between the last two executions (delta).
	drift         time.Duration // For cyclic tasks, the deviation from the scheduled interval.

	// Internal state for the scheduler
	runSignal chan struct{} // For event-driven tasks.
	lastRun   time.Time     // For cyclic tasks.
	mu        sync.RWMutex  // Protects all mutable fields in the task, including metrics.
}

// NewTask creates and initializes a new task.
func NewTask(name string, taskType TaskType, priority int, interval time.Duration) *Task {
	return &Task{
		Name:      name,
		Type:      taskType,
		Priority:  priority,
		Interval:  interval,
		Programs:  []*Program{},
		Enabled:   true,                    // Tasks are enabled by default.
		runSignal: make(chan struct{}, 10), // Buffered channel for events
	}
}

// AddProgram associates a program with the task.
func (t *Task) AddProgram(p *Program) {
	t.Programs = append(t.Programs, p)
}

// WithProgram adds a program to the task and returns the task for chaining.
func (t *Task) WithProgram(p *Program) *Task {
	t.AddProgram(p)
	return t
}

// Enable allows the resource scheduler to execute this task.
func (t *Task) Enable() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Enabled = true
}

// Disable prevents the resource scheduler from executing this task.
func (t *Task) Disable() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Enabled = false
}

// Reset calls the Init() function on all programs within the task.
// This allows for re-initializing the state of the logic controlled by the task.
func (t *Task) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	// Also reset the lastRun time to allow immediate execution if enabled.
	t.lastRun = time.Time{}
	for _, p := range t.Programs {
		p.Init()
	}
}

// Trigger executes an event-driven task.
func (t *Task) Trigger() error {
	if t.Type != EventDrivenTask {
		return fmt.Errorf("cannot trigger non-event-driven task '%s'", t.Name)
	}
	select {
	case t.runSignal <- struct{}{}:
		// Signal sent
	default:
		// Signal already pending, do nothing
	}
	return nil
}

// RemoveProgram removes a program from the task by its name.
// It returns true if the program was found and removed, false otherwise.
func (t *Task) RemoveProgram(name string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, p := range t.Programs {
		if p.Name == name {
			// Remove the element at index i from t.Programs.
			t.Programs = append(t.Programs[:i], t.Programs[i+1:]...)
			return true
		}
	}

	// Program not found.
	return false
}

// FindProgram finds a program within the task by its name.
// It returns the program pointer or nil if not found.
func (t *Task) FindProgram(name string) *Program {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, p := range t.Programs {
		if p.Name == name {
			return p
		}
	}
	return nil
}

// ExecutionTime returns the duration of the last execution of the task's programs.
// This method is safe for concurrent use.
func (t *Task) ExecutionTime() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.executionTime
}

// CycleTime returns the time between the last two executions of the task (delta).
// This method is safe for concurrent use.
func (t *Task) CycleTime() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.cycleTime
}

// Drift returns the deviation from the scheduled interval for cyclic tasks.
// This method is safe for concurrent use.
func (t *Task) Drift() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.drift
}

// taskSorter implements the sort.Interface for a slice of *Task pointers,
// allowing for reflection-free sorting.
type taskSorter []*Task

func (s taskSorter) Len() int           { return len(s) }
func (s taskSorter) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s taskSorter) Less(i, j int) bool { return s[i].Priority < s[j].Priority }

// SortTasks sorts a slice of tasks by priority.
func SortTasks(tasks []*Task) {
	sort.Sort(taskSorter(tasks))
}

// Resource represents a processing unit within the configuration, like a CPU.
type Resource struct {
	Name     string
	Cycle    time.Duration // The base scan cycle of the resource scheduler.
	Affinity int           // Optional: OS CPU core to pin to. -1 or 0 means no affinity.
	Tasks    []*Task

	stopChan chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

// AddTask adds a task to the resource.
func (r *Resource) AddTask(t *Task) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Tasks = append(r.Tasks, t) // Add the new task
	SortTasks(r.Tasks)           // Re-sort the entire slice
}

// WithTask adds a task to the resource and returns the resource for chaining.
func (r *Resource) WithTask(t *Task) *Resource {
	r.AddTask(t)
	return r
}

// RemoveTask removes a task from the resource by its name.
// It returns true if the task was found and removed, false otherwise.
// This operation is safe to call even when the resource is running.
func (r *Resource) RemoveTask(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, t := range r.Tasks {
		if t.Name == name {
			r.Tasks = append(r.Tasks[:i], r.Tasks[i+1:]...)
			SortTasks(r.Tasks) // Re-sort after removal to maintain priority order.
			return true
		}
	}
	return false
}

// FindTask finds a task within the resource by its name.
// It returns the task pointer or nil if not found.
func (r *Resource) FindTask(name string) *Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.Tasks {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (r *Resource) WithResource(t *Task) *Resource {
	return r
}

// shouldRun checks if a task should execute in the current cycle and updates its
// state accordingly. This method must be called within a task's write lock (t.mu.Lock).
func (t *Task) shouldRun(now time.Time) bool {
	if !t.Enabled {
		return false
	}

	switch t.Type {
	case CyclicTask:
		if now.Sub(t.lastRun) >= t.Interval {
			return true
		}
	case EventDrivenTask:
		// Use a non-blocking receive to atomically check for and consume a signal.
		select {
		case <-t.runSignal:
			return true // Signal was present and consumed.
		default:
			return false // No signal.
		}
	}
	return false
}

// updateMetrics calculates and sets the task's runtime metrics after an execution.
// This method must be called within a task's write lock (t.mu.Lock).
func (t *Task) updateMetrics(runTime, executionStartTime time.Time) {
	// Calculate CycleTime (delta) if it's not the first run.
	if !t.lastRun.IsZero() {
		t.cycleTime = runTime.Sub(t.lastRun)
	}

	// Calculate Drift for cyclic tasks.
	if t.Type == CyclicTask {
		// Drift is the deviation from the *scheduled* run time.
		scheduledRunTime := t.lastRun.Add(t.Interval)
		t.drift = runTime.Sub(scheduledRunTime)
	} else {
		// Drift is not applicable to event-driven tasks.
		t.drift = 0
	}

	// Record the time of this run for the next cycle's calculation.
	t.lastRun = runTime

	// Record the total execution time for this run.
	t.executionTime = time.Since(executionStartTime)
}

// schedulerLoop is the core execution loop for a resource. It's called by the platform-specific Start methods.
func (r *Resource) schedulerLoop(ticker *time.Ticker) {
	for {
		select {
		case now := <-ticker.C:
			// For standard Go, we copy the task slice to minimize lock contention,
			// as another goroutine could be trying to modify the task list.
			// For TinyGo, its cooperative scheduler makes this copy unnecessary,
			// so we can iterate directly over the locked slice to save memory.
			var tasksToRun []*Task
			r.mu.Lock()
			tasksToRun = make([]*Task, len(r.Tasks))
			copy(tasksToRun, r.Tasks)
			r.mu.Unlock()

			// In TinyGo, we can hold the read lock for the duration of the loop,
			// as there's no preemption.
			if tasksToRun == nil { // This block will be optimized for TinyGo
				r.mu.Lock()
				tasksToRun = r.Tasks
				defer r.mu.Unlock()
			}

			for _, task := range tasksToRun {
				// The decision to run and the update of task state should be atomic.
				// We can determine if a run is needed inside the task's write lock.
				var shouldRun bool
				task.mu.Lock()
				shouldRun = task.shouldRun(now)
				task.mu.Unlock()

				if shouldRun {
					executionStartTime := time.Now()

					for _, p := range task.Programs {
						func(program *Program) {
							defer func(progName string) {
								if rec := recover(); rec != nil {
									fmt.Printf("PANIC RECOVERED: Task '%s', Program '%s': %v\n", task.Name, progName, rec)
								}
							}(program.Name)
							program.Execute(now)
						}(p)
					}

					// Now, lock the task once to update all its runtime metrics.
					task.mu.Lock()
					task.updateMetrics(now, executionStartTime)
					task.mu.Unlock()
				}
			}
		case <-r.stopChan:
			return
		}
	}
}

// Stop terminates the resource's task scheduler.
func (r *Resource) Stop() {
	r.mu.Lock()
	if !r.running {
		r.mu.Unlock()
		return
	}
	close(r.stopChan)
	r.running = false
	r.mu.Unlock()

	r.wg.Wait() // Wait for the scheduler goroutine to exit.
}

// Pause disables all tasks within the resource, effectively pausing its logical execution
// without stopping the underlying scheduler goroutine. This is safe to call on a running resource.
func (r *Resource) Pause() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, task := range r.Tasks {
		task.Disable()
	}
}

// Resume enables all tasks within the resource, resuming its logical execution.
// This is safe to call on a running resource.
func (r *Resource) Resume() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, task := range r.Tasks {
		task.Enable()
	}
}

// Configuration is the top-level element, representing the entire PLC system.
type Configuration struct {
	Name      string
	Resources []*Resource
}

// WithResource adds a resource to the configuration and returns the configuration for chaining.
func (c *Configuration) WithResource(r *Resource) *Configuration {
	c.Resources = append(c.Resources, r)
	return c
}

// RemoveResource removes a resource from the configuration by its name.
// It returns true if the resource was found and removed, false otherwise.
func (c *Configuration) RemoveResource(name string) bool {
	for i, r := range c.Resources {
		if r.Name == name {
			// Stop the resource before removing it to ensure clean shutdown.
			if r.running {
				r.Stop()
			}
			c.Resources = append(c.Resources[:i], c.Resources[i+1:]...)
			return true
		}
	}
	return false
}

// FindResource finds a resource within the configuration by its name.
// It returns the resource pointer or nil if not found.
func (c *Configuration) FindResource(name string) *Resource {
	// The list of resources is typically static or managed at a higher level,
	// so locking is not implemented here, matching the pattern of RemoveResource.
	for _, r := range c.Resources {
		if r.Name == name {
			return r
		}
	}
	return nil
}
