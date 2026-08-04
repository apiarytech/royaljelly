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
	running  bool
	mu       sync.Mutex
}

// AddTask adds a task to the resource.
func (r *Resource) AddTask(t *Task) {
	r.Tasks = append(r.Tasks, t)
}

// WithTask adds a task to the resource and returns the resource for chaining.
func (r *Resource) WithTask(t *Task) *Resource {
	r.AddTask(t)
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
	r.mu.Unlock()

	go func() {
		// Use a ticker for the resource's main execution loop.
		// A reasonable default could be 1ms for high-speed control.
		if r.Cycle == 0 {
			r.Cycle = time.Millisecond
		}
		ticker := time.NewTicker(r.Cycle)
		defer ticker.Stop()

		for {
			select {
			case now := <-ticker.C:
				// Iterate through tasks in priority order.
				for _, task := range r.Tasks {
					task.mu.Lock() // Lock the task for the entire check-and-run operation.

					shouldRun := false
					if task.Enabled {
						switch task.Type {
						case CyclicTask:
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

					if shouldRun {
						executionStartTime := time.Now()
						// Calculate CycleTime (delta) if it's not the first run.
						if !task.lastRun.IsZero() {
							task.CycleTime = now.Sub(task.lastRun)
						}

						// Calculate Drift for cyclic tasks.
						if task.Type == CyclicTask && !task.lastRun.IsZero() {
							scheduledRunTime := task.lastRun.Add(task.Interval)
							task.Drift = now.Sub(scheduledRunTime)
						} else {
							task.Drift = 0 // Drift is not applicable here.
						}
						task.lastRun = now

						for _, p := range task.Programs {
							func() {
								defer func() {
									if rec := recover(); rec != nil {
										// Use a type switch for robust, reflection-free error logging.
										switch v := rec.(type) {
										case error:
											fmt.Printf("PANIC RECOVERED: Task '%s', Program '%s': %s\n", task.Name, p.Name, v.Error())
										default:
											fmt.Printf("PANIC RECOVERED: Task '%s', Program '%s': %v\n", task.Name, p.Name, v)
										}
									}
								}()
								p.Execute(now)
							}()
						}

						// Record the total execution time for this run.
						task.ExecutionTime = time.Since(executionStartTime)
					}

					task.mu.Unlock()
				}
			case <-r.stopChan:
				return
			}
		}
	}()
}

// Stop terminates the resource's task scheduler.
func (r *Resource) Stop() {
	close(r.stopChan)
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
