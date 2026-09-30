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

// Package core implements the IEC 61131-3 software model: a Configuration holds
// Resources, a Resource runs a scheduler over Tasks, and a Task executes Programs.
//
// # Timing model
//
// Each Resource runs one scheduler goroutine driven by a ticker with period
// Resource.Cycle. On every tick ("scan") the scheduler walks its tasks in
// priority order (lower number first) and runs every task that is due:
//
//   - A cyclic task is due when the scan time reaches its next scheduled time.
//     The next scheduled time advances by exactly Task.Interval, so jitter does
//     not accumulate. The effective interval is quantized to the resource cycle,
//     so choose a Cycle that divides the task intervals evenly.
//   - An event-driven task is due when Task.Trigger has been called since its
//     last run. Up to 10 pending triggers are queued; each scan consumes one.
//
// When a cyclic task falls behind by one or more whole intervals, the missed
// runs are skipped (not replayed), counted in Task.Overruns and reported to the
// fault handler as FaultOverrun. Task.Drift reports how late the last run
// started relative to its scheduled time.
//
// # Priorities are an ordering, not preemption
//
// Tasks on one Resource run sequentially on its single scheduler goroutine.
// Priority decides the order in which due tasks run within a scan; a running
// task is never interrupted by a higher-priority one. A slow program therefore
// delays every other task on the same Resource. To isolate timing-critical
// logic, put it on its own Resource (optionally pinned to a CPU core with
// Resource.Affinity) and exchange data through thread-safe variables such as
// vars.Shared.
//
// # Faults
//
// Panics in programs are recovered per program, so one faulty program does not
// stop the others. Panics, overruns, watchdog expirations and affinity failures
// are reported as a Fault to Resource.OnFault (or Configuration.OnFault). When
// no handler is set, panics and watchdog expirations are written to standard
// error and overruns are only counted.
//
// A task watchdog (Task.Watchdog) detects programs that run too long. Go cannot
// interrupt a running function, so the watchdog reports the fault while the
// program is still running; it does not abort it.
//
// # Concurrency
//
// The exported configuration fields (names, Cycle, Affinity, Interval, Priority,
// Watchdog, OnFault) must be set before the owning Resource is started and not
// changed afterwards. All methods are safe for concurrent use, including adding
// and removing tasks and programs on a running Resource. Program logic runs on
// the scheduler goroutine: do not call Resource.Stop or Configuration.Stop from
// inside a program, because Stop waits for the scheduler to exit.
package core
