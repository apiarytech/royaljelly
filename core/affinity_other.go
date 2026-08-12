//go:build !linux && !windows

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

// setAffinity is a no-op on unsupported operating systems.
func setAffinity(coreID int) error {
	// CPU affinity is not implemented for this OS.
	return nil
}
