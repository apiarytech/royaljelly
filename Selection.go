/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software

 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.
 */

package royaljelly

import "fmt"

// SEL selects one of two inputs based on a boolean selector.
// If G is FALSE, it returns IN0. If G is TRUE, it returns IN1.
// The types of IN0, IN1, and the output must be the same.
func SEL(G BOOL, IN0, IN1 interface{}) interface{} {
	// Basic type check for safety, though Go's type system helps.
	// A more rigorous implementation might use reflection to ensure IN0 and IN1 are of compatible types.
	if G {
		return IN1
	}
	return IN0
}

// MAX returns the maximum value from a series of two or more inputs.
// All inputs must be of a comparable elementary type.
func MAX(inputs []interface{}) interface{} {
	if len(inputs) < 1 {
		// Returning nil for empty or single-element slices as the maximum is undefined.
		// The standard implies extensible functions take 2 or more inputs.
		return nil
	}

	maxVal := inputs[0]
	for i := 1; i < len(inputs); i++ {
		res, err := compare(maxVal, inputs[i])
		if err != nil {
			panic(fmt.Sprintf("MAX: error comparing values: %v", err))
		}
		// If maxVal is less than the current input, update maxVal.
		if res < 0 {
			maxVal = inputs[i]
		}
	}
	return maxVal
}

// MIN returns the minimum value from a series of two or more inputs.
// All inputs must be of a comparable elementary type.
func MIN(inputs []interface{}) interface{} {
	if len(inputs) < 1 {
		// Returning nil for empty or single-element slices as the minimum is undefined.
		return nil
	}

	minVal := inputs[0]
	for i := 1; i < len(inputs); i++ {
		res, err := compare(minVal, inputs[i])
		if err != nil {
			panic(fmt.Sprintf("MIN: error comparing values: %v", err))
		}
		// If minVal is greater than the current input, update minVal.
		if res > 0 {
			minVal = inputs[i]
		}
	}
	return minVal
}

// LIMIT constrains a value to be within a specified minimum (MN) and maximum (MX) range.
// The output is: MN if IN < MN; MX if IN > MX; otherwise IN.
// This is equivalent to MIN(MAX(IN, MN), MX).
// All inputs must be of a comparable elementary type.
func LIMIT(MN, IN, MX interface{}) interface{} {
	// Check if IN is below the minimum limit
	resMN, errMN := compare(IN, MN)
	if errMN != nil {
		panic(fmt.Sprintf("LIMIT: error comparing IN and MN: %v", errMN))
	}
	if resMN < 0 {
		return MN
	}

	// Check if IN is above the maximum limit
	resMX, errMX := compare(IN, MX)
	if errMX != nil {
		panic(fmt.Sprintf("LIMIT: error comparing IN and MX: %v", errMX))
	}
	if resMX > 0 {
		return MX
	}

	// If within limits, return IN
	return IN
}

// MUX selects one input from a list based on an integer selector K.
// The first input is the selector K (ANY_INT). The following inputs are the values to select from.
// The function returns inputs[K+1].
func MUX(inputs []interface{}) interface{} {
	if len(inputs) < 2 {
		// Must have at least a selector (K) and one option (IN0).
		return nil
	}

	// First input is the selector K.
	kVal, err := anyToLINT(inputs[0])
	if err != nil {
		panic(fmt.Sprintf("MUX: selector K must be an integer type, got error: %v", err))
	}

	// The selectable values start from the second element in the slice.
	options := inputs[1:]
	k := int(kVal)

	// Check if K is within the valid range of indices for the options.
	if k < 0 || k >= len(options) {
		// Standard specifies this is an error. Returning nil is a safe default.
		return nil
	}

	return options[k]
}
