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
	"os"
	"time"
)

// Sentinel errors carried in Fault.Err. Test for them with errors.Is.
var (
	// ErrProgramPanic reports a panic recovered from a program.
	ErrProgramPanic = errors.New("program panicked")
	// ErrOverrun reports that a cyclic task missed one or more scheduled runs.
	ErrOverrun = errors.New("task overrun")
	// ErrWatchdog reports that a task ran longer than its watchdog time.
	ErrWatchdog = errors.New("task watchdog expired")
	// ErrAffinityUnsupported reports that CPU affinity is not available on this platform.
	ErrAffinityUnsupported = errors.New("CPU affinity is not supported on this platform")
)

// FaultKind classifies a Fault.
type FaultKind int

const (
	// FaultPanic means a program panicked. Fault.Value holds the panic value.
	FaultPanic FaultKind = iota
	// FaultOverrun means a cyclic task missed scheduled runs.
	FaultOverrun
	// FaultWatchdog means a task exceeded its watchdog time.
	FaultWatchdog
)

// String returns the name of the fault kind.
func (k FaultKind) String() string {
	switch k {
	case FaultPanic:
		return "panic"
	case FaultOverrun:
		return "overrun"
	case FaultWatchdog:
		return "watchdog"
	default:
		return fmt.Sprintf("FaultKind(%d)", int(k))
	}
}

// Fault describes a runtime problem detected by a Resource scheduler.
type Fault struct {
	Kind     FaultKind
	Resource string    // Name of the resource that detected the fault.
	Task     string    // Name of the affected task.
	Program  string    // Name of the affected program, when known.
	Time     time.Time // Scan time at which the fault was detected.
	Value    any       // For FaultPanic, the recovered panic value.
	Err      error     // Wraps one of the Err* sentinel errors.
}

// Error implements the error interface so a Fault can be logged or returned directly.
func (f Fault) Error() string {
	where := fmt.Sprintf("resource %q task %q", f.Resource, f.Task)
	if f.Program != "" {
		where += fmt.Sprintf(" program %q", f.Program)
	}
	return fmt.Sprintf("%s fault in %s: %v", f.Kind, where, f.Err)
}

// Unwrap returns the underlying error.
func (f Fault) Unwrap() error { return f.Err }

// FaultHandler receives faults from a scheduler. It may be called from the
// scheduler goroutine or from a watchdog timer goroutine, so it must be safe
// for concurrent use and should return quickly.
type FaultHandler func(Fault)

// defaultFaultHandler writes panics and watchdog faults to standard error.
// Overruns are only counted, to avoid flooding the log under sustained load.
func defaultFaultHandler(f Fault) {
	if f.Kind == FaultOverrun {
		return
	}
	fmt.Fprintln(os.Stderr, "royaljelly:", f.Error())
}
