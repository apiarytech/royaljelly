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
)

// pinThread locks the scheduler goroutine to its OS thread and pins that thread
// to the given 0-based CPU core.
func pinThread(coreID int) error {
	if n := runtime.NumCPU(); coreID < 0 || coreID >= n {
		return fmt.Errorf("core %d out of range (this machine has %d cores)", coreID+1, n)
	}
	runtime.LockOSThread()
	return setAffinity(coreID)
}
