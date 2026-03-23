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

package royaljelly

import (
	"fmt"
	"reflect"
)

// compare is a helper function to perform comparison between two values.
// It promotes integer-like types to LINT and float-like types to LREAL before comparing.
// It returns:
// -1 if a < b
//
//	0 if a == b
//
// +1 if a > b
// An error if types are incompatible for comparison.
func compare(a, b interface{}) (int, error) {
	// Determine if we need to compare as LREAL or LINT
	typeA := reflect.TypeOf(a)
	typeB := reflect.TypeOf(b)

	isFloatA, isFloatB := isPlcFloatType(typeA), isPlcFloatType(typeB)
	isStringA, isStringB := typeA == reflect.TypeOf(STRING("")), typeB == reflect.TypeOf(STRING(""))
	isTimeA, isTimeB := isPlcTimeType(typeA), isPlcTimeType(typeB)

	// If one is a time type and the other is not (and not a string, which is handled next)
	if isTimeA && !isTimeB && !isStringB {
		return 0, fmt.Errorf("incompatible types for comparison: cannot compare time type %T with numeric type %T", a, b)
	}
	if isTimeB && !isTimeA && !isStringA {
		return 0, fmt.Errorf("incompatible types for comparison: cannot compare time type %T with numeric type %T", b, a)
	}

	if isFloatA || isFloatB { // If either is a float, promote both to LREAL
		valA, errA := anyToLREAL(a)
		if errA != nil {
			return 0, fmt.Errorf("could not convert %v to LREAL: %w", a, errA)
		}
		valB, errB := anyToLREAL(b)
		if errB != nil {
			return 0, fmt.Errorf("could not convert %v to LREAL: %w", b, errB)
		}
		if valA < valB {
			return -1, nil
		}
		if valA > valB {
			return 1, nil
		}
		return 0, nil
	}

	if isStringA && isStringB { // If both are strings, compare them lexicographically
		valA := a.(STRING)
		valB := b.(STRING)
		if valA < valB {
			return -1, nil
		}
		if valA > valB {
			return 1, nil
		}
		return 0, nil
	}

	// Default to LINT comparison for all other compatible types (integers, bools, time types)
	valA, errA := anyToLINT(a)
	if errA != nil {
		return 0, fmt.Errorf("could not convert %v to LINT: %w", a, errA)
	}
	valB, errB := anyToLINT(b)
	if errB != nil {
		return 0, fmt.Errorf("could not convert %v to LINT: %w", b, errB)
	}

	if valA < valB {
		return -1, nil
	}
	if valA > valB {
		return 1, nil
	}
	return 0, nil
}

// GT (Greater Than) checks if IN1 > IN2 > IN3 ...
func GT(inputs []interface{}) BOOL {
	if len(inputs) < 2 {
		// According to the standard, comparison functions are extensible (2 or more inputs).
		// Behavior for less than 2 inputs is undefined, returning false is a safe default.
		return false
	}
	for i := 0; i < len(inputs)-1; i++ {
		res, err := compare(inputs[i], inputs[i+1])
		if err != nil || res <= 0 { // If a > b is not true (error, a <= b)
			return false
		}
	}
	return true
}

// GE (Greater than or Equal) checks if IN1 >= IN2 >= IN3 ...
func GE(inputs []interface{}) BOOL {
	if len(inputs) < 2 {
		return false
	}
	for i := 0; i < len(inputs)-1; i++ {
		res, err := compare(inputs[i], inputs[i+1])
		if err != nil || res < 0 { // If a >= b is not true (error, a < b)
			return false
		}
	}
	return true
}

// EQ (Equal) checks if IN1 == IN2 == IN3 ...
func EQ(inputs []interface{}) BOOL {
	if len(inputs) < 2 {
		return false
	}
	for i := 0; i < len(inputs)-1; i++ {
		res, err := compare(inputs[i], inputs[i+1])
		if err != nil || res != 0 { // If a == b is not true
			return false
		}
	}
	return true
}

// LE (Less than or Equal) checks if IN1 <= IN2 <= IN3 ...
func LE(inputs []interface{}) BOOL {
	if len(inputs) < 2 {
		return false
	}
	for i := 0; i < len(inputs)-1; i++ {
		res, err := compare(inputs[i], inputs[i+1])
		if err != nil || res > 0 { // If a <= b is not true (error, a > b)
			return false
		}
	}
	return true
}

// LT (Less Than) checks if IN1 < IN2 < IN3 ...
func LT(inputs []interface{}) BOOL {
	if len(inputs) < 2 {
		return false
	}
	for i := 0; i < len(inputs)-1; i++ {
		res, err := compare(inputs[i], inputs[i+1])
		if err != nil || res >= 0 { // If a < b is not true (error, a >= b)
			return false
		}
	}
	return true
}

// NE (Not Equal) checks if IN1 != IN2. This function is not extensible.
func NE(inputs []interface{}) BOOL {
	if len(inputs) != 2 {
		// The standard specifies NE is not extensible and takes exactly two inputs.
		// Returning false for incorrect usage is a safe default.
		return false
	}
	res, err := compare(inputs[0], inputs[1])
	if err != nil {
		// Per IEC standard, if types are not comparable, the result is implementation-dependent.
		// Returning false indicates a programming error, which is a safer default than returning true.
		return false
	}
	return res != 0
}
