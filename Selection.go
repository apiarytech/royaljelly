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
func MAX(inputs []interface{}) (interface{}, error) {
	if len(inputs) < 2 {
		return nil, fmt.Errorf("MAX: function requires 2 or more inputs")
	}

	maxVal := inputs[0]
	for i := 1; i < len(inputs); i++ {
		res, err := compare(maxVal, inputs[i])
		if err != nil {
			return nil, fmt.Errorf("MAX: error comparing values: %w", err)
		}
		// If maxVal is less than the current input, update maxVal.
		if res < 0 {
			maxVal = inputs[i]
		}
	}
	return maxVal, nil
}

// MIN returns the minimum value from a series of two or more inputs.
// All inputs must be of a comparable elementary type.
func MIN(inputs []interface{}) (interface{}, error) {
	if len(inputs) < 2 {
		return nil, fmt.Errorf("MIN: function requires 2 or more inputs")
	}

	minVal := inputs[0]
	for i := 1; i < len(inputs); i++ {
		res, err := compare(minVal, inputs[i])
		if err != nil {
			return nil, fmt.Errorf("MIN: error comparing values: %w", err)
		}
		// If minVal is greater than the current input, update minVal.
		if res > 0 {
			minVal = inputs[i]
		}
	}
	return minVal, nil
}

// LIMIT constrains a value to be within a specified minimum (MN) and maximum (MX) range.
// The output is: MN if IN < MN; MX if IN > MX; otherwise IN.
// This is equivalent to MIN(MAX(IN, MN), MX).
// All inputs must be of a comparable elementary type.
func LIMIT(MN, IN, MX interface{}) (interface{}, error) {
	// Check if IN is below the minimum limit
	resMN, errMN := compare(IN, MN)
	if errMN != nil {
		return nil, fmt.Errorf("LIMIT: error comparing IN and MN: %w", errMN)
	}
	if resMN < 0 {
		return MN, nil
	}

	// Check if IN is above the maximum limit
	resMX, errMX := compare(IN, MX)
	if errMX != nil {
		return nil, fmt.Errorf("LIMIT: error comparing IN and MX: %w", errMX)
	}
	if resMX > 0 {
		return MX, nil
	}

	// If within limits, return IN
	return IN, nil
}

// MUX selects one input from a list based on an integer selector K.
// The first input is the selector K (ANY_INT). The following inputs are the values to select from.
// The function returns inputs[K+1].
func MUX(inputs []interface{}) (interface{}, error) {
	if len(inputs) < 2 {
		// Must have at least a selector (K) and one option (IN0).
		return nil, fmt.Errorf("MUX: function requires at least 2 inputs (a selector K and one value)")
	}

	// First input is the selector K.
	kVal, err := anyToLINT(inputs[0])
	if err != nil {
		return nil, fmt.Errorf("MUX: selector K must be an integer type, got error: %w", err)
	}

	// The selectable values start from the second element in the slice.
	options := inputs[1:]
	k := int(kVal)

	// Check if K is within the valid range of indices for the options.
	if k < 0 || k >= len(options) {
		// Standard specifies this is an error. Returning nil and an error is a safe default.
		return nil, fmt.Errorf("MUX: selector K (%d) is out of bounds", k)
	}

	return options[k], nil
}
