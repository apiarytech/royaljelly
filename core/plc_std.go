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
	"fmt"
	"runtime"
	"time"
)

func init() {
	// For standard Go builds, automatically set GOMAXPROCS to use all available CPU cores.
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// Start begins the resource's priority-based task scheduler.
// This version is for standard Go and supports CPU affinity.
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
	affinity := r.Affinity // Capture affinity before unlocking
	r.mu.Unlock()

	go func() {
		// Pin this goroutine to a specific OS thread if requested.
		// This is a prerequisite for influencing CPU core placement.
		if affinity > 0 {
			runtime.LockOSThread()

			// Now, use gopsutil to set the CPU affinity for the locked thread.
			// The affinity value is 1-based, so we subtract 1 for the 0-based core index.
			coreID := affinity - 1 // Convert 1-based affinity to 0-based core ID
			if err := setAffinity(coreID); err != nil {
				panic(fmt.Sprintf("CRITICAL: Failed to set CPU affinity for resource '%s' to core %d: %v", r.Name, affinity, err))
			}
		}

		// Use a ticker for the resource's main execution loop.
		if r.Cycle == 0 {
			r.Cycle = time.Millisecond
		}
		ticker := time.NewTicker(r.Cycle)
		defer ticker.Stop()
		defer r.wg.Done()

		// The core scheduling loop.
		r.schedulerLoop(ticker)
	}()
}
