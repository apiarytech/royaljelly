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

// Variable is one of a program's variables as its host sees it, so the host
// can bind the variable to its own data: a tag database, retain storage, an
// HMI. A program's Variables method lists them; the beedance transpiler
// generates it.
//
// Get and Set read and write the variable without reflection. Set reports
// false, and changes nothing, when v is not of the variable's Go type (for an
// INT, iec.INT). A host calls them only between scans: a program's state is
// its own while it runs.
type Variable struct {
	// Name is the variable's name as declared.
	Name string
	// Block is its declaration block: VAR_INPUT, VAR_OUTPUT, VAR_IN_OUT or
	// VAR.
	Block string
	// Type is its declared IEC 61131-3 type, e.g. "REAL".
	Type string
	// Address is a located variable's address, e.g. "%QX0.1"; empty
	// otherwise.
	Address  string
	Retain   bool
	Constant bool
	Get      func() any
	Set      func(v any) bool
}
