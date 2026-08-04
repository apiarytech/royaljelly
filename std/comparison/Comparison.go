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

package comparison

import (
	. "github.com/apiarytech/royaljelly/core"
)

// GT (Greater Than) checks if IN1 > IN2 > IN3 ...
func GT[T ANY_MAGNITUDE](inputs ...T) BOOL {
	if len(inputs) < 2 {
		return false
	}
	for i := 0; i < len(inputs)-1; i++ {
		if inputs[i] <= inputs[i+1] {
			return false
		}
	}
	return true
}

// GE (Greater than or Equal) checks if IN1 >= IN2 >= IN3 ...
func GE[T ANY_MAGNITUDE](inputs ...T) BOOL {
	if len(inputs) < 2 {
		return false
	}
	for i := 0; i < len(inputs)-1; i++ {
		if inputs[i] < inputs[i+1] {
			return false
		}
	}
	return true
}

// EQ (Equal) checks if IN1 == IN2 == IN3 ... It can be used with any comparable type.
func EQ[T ANY_ELEMENTARY](inputs ...T) BOOL {
	if len(inputs) < 2 {
		return false
	}
	for i := 0; i < len(inputs)-1; i++ {
		if inputs[i] != inputs[i+1] {
			return false
		}
	}
	return true
}

// LE (Less than or Equal) checks if IN1 <= IN2 <= IN3 ...
func LE[T ANY_MAGNITUDE](inputs ...T) BOOL {
	if len(inputs) < 2 {
		return false
	}
	for i := 0; i < len(inputs)-1; i++ {
		if inputs[i] > inputs[i+1] {
			return false
		}
	}
	return true
}

// LT (Less Than) checks if IN1 < IN2 < IN3 ...
func LT[T ANY_MAGNITUDE](inputs ...T) BOOL {
	if len(inputs) < 2 {
		return false
	}
	for i := 0; i < len(inputs)-1; i++ {
		if inputs[i] >= inputs[i+1] {
			return false
		}
	}
	return true
}

// NE (Not Equal) checks if IN1 != IN2. This function is not extensible and can be used with any comparable type.
func NE[T ANY_ELEMENTARY](in1, in2 T) BOOL {
	return in1 != in2
}
