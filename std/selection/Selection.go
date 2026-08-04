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

package selection

import (
	"fmt"

	. "github.com/apiarytech/royaljelly/core"
)

// SEL selects one of two inputs based on a boolean selector.
// If G is FALSE, it returns IN0. If G is TRUE, it returns IN1.
// The types of IN0, IN1, and the output must be the same.
func SEL[T any](G BOOL, IN0, IN1 T) T {
	if G {
		return IN1
	}
	return IN0
}

// MAX returns the maximum value from a series of two or more inputs of an ordered type.
func MAX[T ANY_MAGNITUDE](inputs ...T) (T, error) {
	if len(inputs) < 2 {
		var zero T
		return zero, fmt.Errorf("MAX: function requires 2 or more inputs")
	}

	maxVal := inputs[0]
	for i := 1; i < len(inputs); i++ {
		if maxVal < inputs[i] {
			maxVal = inputs[i]
		}
	}
	return maxVal, nil
}

// MIN returns the minimum value from a series of two or more inputs of an ordered type.
func MIN[T ANY_MAGNITUDE](inputs ...T) (T, error) {
	if len(inputs) < 2 {
		var zero T
		return zero, fmt.Errorf("MIN: function requires 2 or more inputs")
	}

	minVal := inputs[0]
	for i := 1; i < len(inputs); i++ {
		if minVal > inputs[i] {
			minVal = inputs[i]
		}
	}
	return minVal, nil
}

// LIMIT constrains a value to be within a specified minimum (MN) and maximum (MX) range.
// The output is: MN if IN < MN; MX if IN > MX; otherwise IN. All inputs must be of the same ordered type.
func LIMIT[T ANY_MAGNITUDE](MN, IN, MX T) T {
	if IN < MN {
		return MN
	}
	if IN > MX {
		return MX
	}
	return IN
}

// MUX selects one input from a list based on an integer selector K.
// The function returns one of the `options` based on the index `K`.
func MUX[K ANY_INT, T any](selector K, options ...T) (T, error) {
	if len(options) == 0 {
		var zero T
		return zero, fmt.Errorf("MUX: function requires at least one option to select from")
	}

	// Convert the generic integer selector to a standard int for indexing.
	k := int(selector)

	// Check if K is within the valid range of indices for the options.
	if k < 0 || k >= len(options) {
		var zero T
		return zero, fmt.Errorf("MUX: selector K (%d) is out of bounds for %d options", k, len(options))
	}

	return options[k], nil
}
