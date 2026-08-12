//go:build windows && !tinygo && !linux

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
	"golang.org/x/sys/windows"
)

var (
	modkernel32               = windows.NewLazySystemDLL("kernel32.dll")
	procSetThreadAffinityMask = modkernel32.NewProc("SetThreadAffinityMask")
)

// setAffinity pins the current OS thread to a specific CPU core.
func setAffinity(coreID int) error {
	// On Windows, the affinity is set with a bitmask.
	// 1 << coreID creates a mask with the bit for the desired core set to 1.
	affinityMask := uintptr(1) << coreID

	// Get the current thread handle and set its affinity mask.
	r1, _, err := procSetThreadAffinityMask.Call(uintptr(windows.CurrentThread()), affinityMask)
	if r1 == 0 {
		return err
	}
	return nil
}
