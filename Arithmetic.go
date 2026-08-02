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
	"strconv"
	"time"
)

/*********************************/
/*  IEC Arithmatic definitions   */
/*********************************/

// --- Helper Functions ---

// isPlcFloatType checks if the given Type is one of the PLC float types.
func IsPlcFloat[T any](val T) bool {
	switch any(val).(type) {
	case REAL, LREAL, float32, float64:
		return true
	default:
		return false
	}
}

// isPlcIntType checks if the given Type is one of the PLC integer or bit-string types
// that can be treated as integers in arithmetic.
func IsPlcInt[T any](val T) bool {
	switch any(val).(type) {
	case SINT, INT, DINT, LINT,
		USINT, UINT, UDINT, ULINT,
		BOOL, BYTE, WORD, DWORD, LWORD:
		return true
	default:
		return false
	}
}

// isPlcTimeType checks if the given Type is one of the PLC time/date types.
func IsPlcTimeType[T any](val T) bool {
	switch any(val).(type) {
	case TIME, DATE, TOD, DT:
		return true
	default:
		return false
	}
}

// anytoREAL converts a supported PLC type to REAL (float32).
// It leverages anyToLREAL to avoid code duplication and then converts the result.
func anytoREAL[T any](val T) (REAL, error) {
	lrealVal, err := anyToLREAL(val)
	if err != nil {
		// Modify the error message to be specific to this function's context.
		return 0, fmt.Errorf("anytoREAL: failed during intermediate conversion to LREAL: %w", err)
	}
	return REAL(lrealVal), nil
}

// anyToLREAL converts a supported PLC type to LREAL. It uses a type switch on the generic parameter T.
func anyToLREAL[T any](val T) (LREAL, error) {
	switch v := any(val).(type) {
	case BOOL:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	case TIME:
		return LREAL(time.Duration(v).Milliseconds()), nil
	case DATE:
		return LREAL(time.Time(v).UnixMilli()), nil
	case TOD:
		t_time := time.Time(v)
		midnight := time.Date(t_time.Year(), t_time.Month(), t_time.Day(), 0, 0, 0, 0, t_time.Location())
		return LREAL(t_time.Sub(midnight).Milliseconds()), nil
	case DT:
		return LREAL(time.Time(v).UnixMilli()), nil
	case STRING:
		f, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			return 0, fmt.Errorf("anyToLREAL: cannot parse STRING '%s' to LREAL: %w", v, err)
		}
		return LREAL(f), nil
	// Handle numeric types explicitly to avoid reflection.
	case SINT:
		return LREAL(v), nil
	case INT:
		return LREAL(v), nil
	case DINT:
		return LREAL(v), nil
	case LINT:
		return LREAL(v), nil
	case USINT:
		return LREAL(v), nil
	case UINT:
		return LREAL(v), nil
	case UDINT:
		return LREAL(v), nil
	case ULINT:
		return LREAL(v), nil
	case REAL:
		return LREAL(v), nil
	case LREAL:
		return v, nil
	case float32:
		return LREAL(v), nil
	case float64:
		return LREAL(v), nil
	case BYTE:
		return LREAL(v), nil
	case WORD:
		return LREAL(v), nil
	case DWORD:
		return LREAL(v), nil
	case LWORD:
		return LREAL(v), nil
	default:
		return 0, fmt.Errorf("anyToLREAL: unsupported type %T for conversion", v)
	}
}

// anyToLINT converts a supported PLC type to LINT.
func anyToLINT[T any](val T) (LINT, error) {
	switch v := any(val).(type) {
	case BOOL:
		if v {
			return 1, nil
		}
		return 0, nil
	case TIME:
		return LINT(time.Duration(v).Milliseconds()), nil
	case DATE:
		return LINT(time.Time(v).UnixMilli()), nil
	case TOD:
		t := time.Time(v)
		midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		return LINT(t.Sub(midnight).Milliseconds()), nil
	case DT:
		return LINT(time.Time(v).UnixMilli()), nil
	case STRING:
		i, err := strconv.Atoi(string(v))
		if err != nil {
			return 0, fmt.Errorf("anyToLINT: cannot parse STRING '%s' to LINT: %w", v, err)
		}
		return LINT(i), nil
	// Handle numeric types explicitly.
	case SINT:
		return LINT(v), nil
	case INT:
		return LINT(v), nil
	case DINT:
		return LINT(v), nil
	case LINT:
		return v, nil
	case USINT:
		return LINT(v), nil
	case UINT:
		return LINT(v), nil
	case UDINT:
		return LINT(v), nil
	case ULINT:
		return LINT(v), nil // Note: Potential for overflow
	case REAL:
		return LINT(v), nil // Note: Truncation will occur
	case LREAL:
		return LINT(v), nil // Note: Truncation will occur
	case BYTE, WORD, DWORD, LWORD:
		// Convert bit-string types through ULINT to handle them as unsigned integers
		uVal, _ := anyToULINT(v)
		return LINT(uVal), nil
	default:
		return 0, fmt.Errorf("anyToLINT: unsupported type %T for conversion", v)
	}
}

// anyToULINT converts a supported PLC bitwise type to ULINT.
func anyToULINT[T any](val T) (ULINT, error) {
	switch v := any(val).(type) {
	case SINT:
		return ULINT(v), nil
	case INT:
		return ULINT(v), nil
	case DINT:
		return ULINT(v), nil
	case LINT:
		return ULINT(v), nil
	case USINT:
		return ULINT(v), nil
	case UINT:
		return ULINT(v), nil
	case UDINT:
		return ULINT(v), nil
	case ULINT:
		return v, nil
	case REAL:
		if v < 0 {
			return REAL_TO_ULINT(v), nil // Handle negative floats via LINT
		}
		return ULINT(v), nil // Truncation for positive floats
	case LREAL:
		if v < 0 {
			return LREAL_TO_ULINT(v), nil // Handle negative floats via LINT
		}
		return ULINT(v), nil // Truncation for positive floats
	case BOOL:
		if v {
			return 1, nil
		}
		return 0, nil
	case BYTE:
		return ULINT(v), nil
	case WORD:
		return ULINT(v), nil
	case DWORD:
		return ULINT(v), nil
	case LWORD:
		return ULINT(v), nil
	case TIME:
		return ULINT(time.Duration(v).Milliseconds()), nil
	case DATE:
		return ULINT(time.Time(v).UnixMilli()), nil
	case TOD:
		t_time := time.Time(v)
		midnight := time.Date(t_time.Year(), t_time.Month(), t_time.Day(), 0, 0, 0, 0, t_time.Location())
		return ULINT(t_time.Sub(midnight).Milliseconds()), nil
	case DT:
		return ULINT(time.Time(v).UnixMilli()), nil
	case STRING:
		// Bitwise operations on strings are not standard, but can be interpreted as parsing to an integer.
		i, err := strconv.ParseUint(string(v), 0, 64) // Use ParseUint with base 0 for auto-detection (e.g., "0xFF")
		if err != nil {
			return 0, fmt.Errorf("anyToULINT: cannot parse STRING '%s' to ULINT: %w", v, err)
		}
		return ULINT(i), nil
	default:
		// Return an error for types that are not bitwise-compatible.
		return 0, fmt.Errorf("anyToULINT: unsupported type %T for conversion to ULINT", val)
	}
}

// ConvertTo provides a generic, type-safe way to convert any supported PLC type to a specific target type `T`.
// It replaces the need for the reflection-based `convertToTargetType` function.
func ConvertTo[T any](in any) (T, error) {
	var zero T
	var result any
	var err error

	// Use a type switch on the *target* type `T`.
	// We create a zero value of T and switch on its type.
	switch any(zero).(type) {
	case SINT:
		val, err := anyToLINT(in)
		if err != nil {
			return zero, err
		}
		result = SINT(val)
	case DINT:
		val, err := anyToLINT(in)
		if err != nil {
			return zero, err
		}
		result = DINT(val)
	case LINT:
		result, err = anyToLINT(in)
	case INT:
		val, err := anyToLINT(in)
		if err != nil {
			return zero, err
		}
		result = INT(val)
	case USINT:
		val, err := anyToULINT(in)
		if err != nil {
			return zero, err
		}
		result = USINT(val)
	case UINT:
		val, err := anyToULINT(in)
		if err != nil {
			return zero, err
		}
		result = UINT(val)
	case UDINT:
		val, err := anyToULINT(in)
		if err != nil {
			return zero, err
		}
		result = UDINT(val)
	case ULINT:
		result, err = anyToULINT(in)
	case BOOL:
		val, err := anyToLINT(in)
		if err != nil {
			// Fallback for non-numeric types like STRING "true"
			if s, ok := in.(STRING); ok && (s == "true" || s == "TRUE") {
				val = 1
			} else {
				return zero, err
			}
		}
		result = BOOL(val != 0)
	case BYTE:
		val, err := anyToULINT(in)
		if err != nil {
			return zero, err
		}
		result = BYTE(val)
	case WORD:
		val, err := anyToULINT(in)
		if err != nil {
			return zero, err
		}
		result = WORD(val)
	case DWORD:
		val, err := anyToULINT(in)
		if err != nil {
			return zero, err
		}
		result = DWORD(val)
	case LWORD:
		val, err := anyToULINT(in)
		if err != nil {
			return zero, err
		}
		result = LWORD(val)
	case TIME:
		result, err = SubTime(in)
	case DATE:
		result, err = SubDate(in)
	case TOD:
		result, err = SubTod(in)
	case DT:
		result, err = SubDt(in)
	case REAL:
		result, err = anytoREAL(in)
	case LREAL:
		result, err = anyToLREAL(in)
	case STRING:
		// Special handling for string conversion
		if s, ok := in.(fmt.Stringer); ok {
			result = STRING(s.String())
		} else {
			result = STRING(fmt.Sprintf("%v", in))
		}
	default:
		// For any other type, we can attempt a promotion to LREAL or LINT as an intermediate.
		// This part shows the complexity of a truly universal converter.
		// For this example, we'll primarily rely on the explicit cases.
		return zero, fmt.Errorf("ConvertTo: unsupported target type %T", zero)
	}

	if err != nil {
		return zero, err
	}

	// Final type assertion to the requested generic type T.
	if finalResult, ok := result.(T); ok {
		return finalResult, nil
	}

	return zero, fmt.Errorf("ConvertTo: internal error, failed to cast result to %T", zero)
}

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
	val, err := anyToLREAL(in2)
	if err != nil {
		return 0, fmt.Errorf("MUL_TIME: error converting multiplier %v to LREAL: %w", in2, err)
	}
	return TIME(float64(in1) * float64(val)), nil
}

// DIV_TIME divides a TIME duration by a numeric value.
func DIV_TIME(in1 TIME, in2 any) (TIME, error) {
	val, err := anyToLREAL(in2)
	if err != nil {
		return 0, fmt.Errorf("DIV_TIME: error converting divisor %v to LREAL: %w", in2, err)
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
