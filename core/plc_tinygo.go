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
	"time"
)

// Start begins the resource's priority-based task scheduler.
// This version is for TinyGo and runs on a single thread without affinity.
func (r *Resource) Start() {
	r.mu.Lock()
	if r.running {
		r.mu.Unlock()
		return
	}

	SortTasks(r.Tasks)
	r.running = true
	r.stopChan = make(chan struct{})
	r.wg.Add(1)
	r.mu.Unlock()

	go func() {
		if r.Cycle == 0 {
			r.Cycle = time.Millisecond
		}
		ticker := time.NewTicker(r.Cycle)
		defer ticker.Stop()
		defer r.wg.Done()

		r.schedulerLoop(ticker)
	}()
}
