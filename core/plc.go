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
	Name  string
	Logic func(now time.Time) // The user-defined logic for the program.
}

// Execute runs the program's defined logic.
func (p *Program) Execute(now time.Time) {
	if p.Logic != nil {
		p.Logic(now)
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
	Enabled  BOOL // If false, the task will not be scheduled.

	// --- Runtime Metrics ---
	ExecutionTime time.Duration // Last execution time of the task's programs.
	CycleTime     time.Duration // Time between the last two executions (delta).
	Drift         time.Duration // For cyclic tasks, the deviation from the scheduled interval.

	// Internal state for the scheduler
	runSignal chan struct{} // For event-driven tasks.
	lastRun   time.Time     // For cyclic tasks.
	mu        sync.RWMutex  // Protects all mutable fields in the task.
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
	Name  string
	Cycle time.Duration // The base scan cycle of the resource scheduler.
	Tasks []*Task

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

// Start begins the resource's priority-based task scheduler.
func (r *Resource) Start() {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return
	}

	// Sort tasks by priority (lower number is higher priority)
	SortTasks(r.Tasks)

	r.running = true
	r.stopChan = make(chan struct{})
	r.wg.Add(1)
	r.mu.Unlock()

	go func() {
		// Use a ticker for the resource's main execution loop.
		// A reasonable default could be 1ms for high-speed control.
		// This ensures the scheduler runs at least at the specified cycle.
		if r.Cycle == 0 {
			r.Cycle = time.Millisecond
		}
		ticker := time.NewTicker(r.Cycle)
		defer ticker.Stop()
		defer r.wg.Done()

		for {
			select {
			case now := <-ticker.C:
				// Acquire lock to get a consistent snapshot of tasks for this cycle.
				// This prevents race conditions if r.Tasks is modified (added/removed/sorted)
				// by another goroutine during this iteration.
				r.mu.Lock()
				currentTasks := make([]*Task, len(r.Tasks))
				copy(currentTasks, r.Tasks)
				r.mu.Unlock() // Release the lock as soon as the copy is made.

				for _, task := range currentTasks {
					var shouldRun bool
					task.mu.RLock() // Use a read lock to check if the task should run.
					if task.Enabled {
						switch task.Type {
						case CyclicTask:
							// Check if the interval has passed.
							if now.Sub(task.lastRun) >= task.Interval {
								shouldRun = true
							}
						case EventDrivenTask:
							// Check for and consume a signal atomically.
							if len(task.runSignal) > 0 {
								<-task.runSignal // Consume one signal
								shouldRun = true
							}
						}
					}
					task.mu.RUnlock()

					if shouldRun {
						executionStartTime := time.Now()

						// Lock the task with a write lock only when modifying its state.
						task.mu.Lock()
						// Calculate CycleTime (delta) if it's not the first run.
						if !task.lastRun.IsZero() {
							task.CycleTime = now.Sub(task.lastRun)
						}

						// Calculate Drift for cyclic tasks.
						if task.Type == CyclicTask {
							scheduledRunTime := task.lastRun.Add(task.Interval)
							task.Drift = now.Sub(scheduledRunTime)
						} else {
							task.Drift = 0 // Drift is not applicable here.
						}
						task.lastRun = now
						task.mu.Unlock() // *** Release the lock before executing user code ***

						for _, p := range task.Programs {
							// The panic recovery can be simplified slightly.
							// It's good practice to keep it to protect the scheduler loop.
							func(program *Program) {
								defer func(progName string) {
									if rec := recover(); rec != nil {
										fmt.Printf("PANIC RECOVERED: Task '%s', Program '%s': %v\n", task.Name, progName, rec)
									}
								}(program.Name)
								program.Execute(now)
							}(p)
						}

						// Re-acquire the lock briefly to update the final metric.
						task.mu.Lock()
						// Record the total execution time for this run.
						task.ExecutionTime = time.Since(executionStartTime)
						task.mu.Unlock()
					}
				}
			case <-r.stopChan:
				return
			}
		}
	}()
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
