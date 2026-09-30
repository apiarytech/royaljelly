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

// pinThread is unsupported under TinyGo, which runs goroutines cooperatively on
// a single thread.
func pinThread(coreID int) error {
	return ErrAffinityUnsupported
}
