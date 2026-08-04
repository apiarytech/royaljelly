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

package arithmetic

import (
	"fmt"
	"time"

	. "github.com/apiarytech/royaljelly/core"
)

/*********************************/
/*  IEC Arithmatic definitions   */
/*********************************/

// --- Helper Functions ---

// ADD performs addition on a slice of numeric types.
// It uses generics to operate on any numeric type defined in the ANY_NUM constraint.
func ADD[T ANY_NUM](nums ...T) T {
	var sum T
	for _, num := range nums {
		sum += num
	}
	return sum
}

// SUB performs subtraction on a slice of numeric types.
// It uses generics to operate on any numeric type defined in the ANY_NUM constraint.
// It calculates nums[0] - nums[1] - nums[2]...
func SUB[T ANY_NUM](nums ...T) T {
	if len(nums) == 0 {
		var zero T
		return zero
	}
	result := nums[0]
	for i := 1; i < len(nums); i++ {
		result -= nums[i]
	}
	return result
}

// MUL performs multiplication on a slice of numeric types.
// It uses generics to operate on any numeric type defined in the ANY_NUM constraint.
func MUL[T ANY_NUM](nums ...T) T {
	if len(nums) == 0 {
		var one T = 1
		return one
	}
	result := T(1)
	for _, num := range nums {
		result *= num
	}
	return result
}

// DIV performs division on a slice of numeric types.
// It uses generics to operate on any numeric type defined in the ANY_NUM or ANY_INT constraint.
// It calculates (...(nums[0] / nums[1]) / nums[2] ...).
func DIV[T ANY_NUM](nums ...T) (T, error) {
	if len(nums) == 0 {
		var zero T
		return zero, nil
	}
	result := nums[0]
	for i := 1; i < len(nums); i++ {
		divisor := nums[i]
		var zero T
		if divisor == zero {
			return zero, fmt.Errorf("DIV: division by zero")
		}
		result /= divisor
	}
	return result, nil
}

// MOD performs the modulo operation on a slice of integer types.
// It uses generics to operate on any integer type defined in the ANY_INT constraint.
// It calculates (...(nums[0] % nums[1]) % nums[2] ...).
func MOD[T ANY_INT](nums ...T) (T, error) {
	if len(nums) < 2 {
		return 0, fmt.Errorf("MOD: function requires at least 2 inputs")
	}
	result := nums[0]
	for i := 1; i < len(nums); i++ {
		divisor := nums[i]
		var zero T
		if divisor == 0 {
			return zero, fmt.Errorf("MOD: modulo by zero")
		}
		result %= divisor
	}
	return result, nil
}

// MOVE performs an assignment of the input value.
// The standard defines this as a non-extensible function with one input and one output.
// Using generics ensures type safety, as the output type is guaranteed to be the same as the input type.
func MOVE[T any](in T) T {
	return in
}

/*****************************************************************/
/* IEC 61131-3 Standard Functions of Time Data Types (Table 30)  */
/*****************************************************************/

// ADD_TIME adds two TIME durations.
func ADD_TIME(in1, in2 TIME) TIME {
	return in1 + in2
}

// ADD_TOD adds a TIME duration to a TIME_OF_DAY.
func ADD_TOD(in1 TOD, in2 TIME) TOD {
	return TOD(time.Time(in1).Add(time.Duration(in2)))
}

// ADD_DT adds a TIME duration to a DATE_AND_TIME.
func ADD_DT(in1 DT, in2 TIME) DT {
	return DT(time.Time(in1).Add(time.Duration(in2)))
}

// SUB_TIME subtracts two TIME durations.
func SUB_TIME(in1, in2 TIME) TIME {
	return in1 - in2
}

// SUB_DATE subtracts two DATEs, resulting in a TIME duration.
func SUB_DATE(in1, in2 DATE) TIME {
	return TIME(time.Time(in1).Sub(time.Time(in2)))
}

// SUB_TOD_TIME subtracts a TIME from a TIME_OF_DAY, resulting in a new TIME_OF_DAY.
func SUB_TOD_TIME(in1 TOD, in2 TIME) TOD {
	return TOD(time.Time(in1).Add(-time.Duration(in2)))
}

// SUB_TOD_TOD subtracts two TIME_OF_DAY values, resulting in a TIME duration.
func SUB_TOD_TOD(in1, in2 TOD) TIME {
	return TIME(time.Time(in1).Sub(time.Time(in2)))
}

// SUB_DT_TIME subtracts a TIME from a DATE_AND_TIME, resulting in a new DATE_AND_TIME.
func SUB_DT_TIME(in1 DT, in2 TIME) DT {
	return DT(time.Time(in1).Add(-time.Duration(in2)))
}

// SUB_DT_DT subtracts two DATE_AND_TIME values, resulting in a TIME duration.
func SUB_DT_DT(in1, in2 DT) TIME {
	return TIME(time.Time(in1).Sub(time.Time(in2)))
}

// MUL_TIME multiplies a TIME duration by a numeric value.
func MUL_TIME(in1 TIME, in2 any) (TIME, error) {
	val, err := AnyToLREAL(in2)
	if err != nil {
		return 0, &ConversionError{
			Value:    in2,
			FromType: GetTypeName(in2),
			ToType:   "LREAL (for multiplication with TIME)",
			Reason:   "invalid multiplier for TIME",
			Err:      err,
		}
	}
	return TIME(float64(in1) * float64(val)), nil
}

// DIV_TIME divides a TIME duration by a numeric value.
func DIV_TIME(in1 TIME, in2 any) (TIME, error) {
	val, err := AnyToLREAL(in2)
	if err != nil {
		return 0, &ConversionError{
			Value:    in2,
			FromType: GetTypeName(in2),
			ToType:   "LREAL (for division with TIME)",
			Reason:   "invalid divisor for TIME",
			Err:      err,
		}
	}
	if val == 0 {
		return 0, fmt.Errorf("DIV_TIME: division by zero")
	}
	return TIME(float64(in1) / float64(val)), nil
}

// CONCAT_DATE_TOD concatenates a DATE and a TIME_OF_DAY to create a DATE_AND_TIME.
func CONCAT_DATE_TOD(in1 DATE, in2 TOD) DT {
	d := time.Time(in1)
	t := time.Time(in2)
	return DT(time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), d.Location()))
}
