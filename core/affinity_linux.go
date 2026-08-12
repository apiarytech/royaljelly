//go:build linux && !tinygo && !windows

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
	"golang.org/x/sys/unix"
)

// setAffinity pins the current OS thread to a specific CPU core.
func setAffinity(coreID int) error {
	// Create a CPU set with a single core.
	var cpuSet unix.CPUSet
	cpuSet.Set(coreID)

	// Set the affinity for the current thread (pid=0).
	return unix.SchedSetaffinity(0, &cpuSet)
}
