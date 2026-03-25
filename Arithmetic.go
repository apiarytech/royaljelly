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
	"math"
	"reflect"
	"strconv"
	"time"
)

/*********************************/
/*  IEC Arithmatic definitions   */
/*********************************/

// --- Helper Functions ---

// isPlcFloatType checks if the given reflect.Type is one of the PLC float types.
func isPlcFloatType(rt reflect.Type) bool {
	return rt == reflect.TypeOf(REAL(0)) || rt == reflect.TypeOf(LREAL(0)) || rt == reflect.TypeOf(float32(0)) || rt == reflect.TypeOf(float64(0))
}

// isPlcIntType checks if the given reflect.Type is one of the PLC integer or bit-string types
// that can be treated as integers in arithmetic.
func isPlcIntType(rt reflect.Type) bool {
	switch rt {
	case reflect.TypeOf(SINT(0)), reflect.TypeOf(INT(0)), reflect.TypeOf(DINT(0)), reflect.TypeOf(LINT(0)),
		reflect.TypeOf(USINT(0)), reflect.TypeOf(UINT(0)), reflect.TypeOf(UDINT(0)), reflect.TypeOf(ULINT(0)),
		reflect.TypeOf(BOOL(false)), // BOOL is treated as 0 or 1
		reflect.TypeOf(BYTE(0)), reflect.TypeOf(WORD(0)), reflect.TypeOf(DWORD(0)), reflect.TypeOf(LWORD(0)):
		return true
	default:
		return false
	}
}

// isPlcTimeType checks if the given reflect.Type is one of the PLC time/date types.
func isPlcTimeType(rt reflect.Type) bool {
	switch rt {
	case reflect.TypeOf(TIME(0)), reflect.TypeOf(DATE{}), reflect.TypeOf(TOD{}), reflect.TypeOf(DT{}):
		return true
	default:
		return false
	}
}

// anytoREAL converts a supported PLC type to REAL (float32).
// It leverages anyToLREAL to avoid code duplication and then converts the result.
func anytoREAL(val interface{}) (REAL, error) {
	lrealVal, err := anyToLREAL(val)
	if err != nil {
		// Modify the error message to be specific to this function's context.
		return 0, fmt.Errorf("anytoREAL: failed during intermediate conversion to LREAL: %w", err)
	}
	return REAL(lrealVal), nil
}

// anyToLREAL converts a supported PLC type to LREAL.
func anyToLREAL(val interface{}) (LREAL, error) {
	switch v := val.(type) {
	case SINT:
		return LREAL(v), nil
	case INT:
		return LREAL(v), nil
	case DINT:
		return LREAL(v), nil
	case LINT:
		return LREAL(v), nil // Note: Max LINT > Max LREAL representable without precision loss
	case USINT:
		return LREAL(v), nil
	case UINT:
		return LREAL(v), nil
	case UDINT:
		return LREAL(v), nil
	case ULINT:
		return LREAL(v), nil // Note: Max ULINT > Max LREAL representable without precision loss
	case REAL:
		return LREAL(v), nil
	case LREAL:
		return v, nil
	case float32: // Handle raw Go float32
		return LREAL(v), nil
	case float64: // Handle raw Go float64
		return LREAL(v), nil
	case BOOL:
		if v {
			return 1.0, nil
		} else {
			return 0.0, nil
		}
	case BYTE:
		return LREAL(v), nil
	case WORD:
		return LREAL(v), nil
	case DWORD:
		return LREAL(v), nil
	case LWORD:
		return LREAL(v), nil // Note: Max LWORD (uint64) > Max LREAL
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
	default:
		return 0, fmt.Errorf("anyToLREAL: unsupported type %T for conversion to LREAL", val)
	}
}

// anyToLINT converts a supported PLC type to LINT.
func anyToLINT(val interface{}) (LINT, error) {
	switch v := val.(type) {
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
		return LINT(v), nil // Note: Max ULINT > Max LINT, potential overflow
	case REAL:
		return LINT(v), nil // Truncation
	case LREAL:
		return LINT(v), nil // Truncation
	case BOOL:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case BYTE:
		return LINT(v), nil
	case WORD:
		return LINT(v), nil
	case DWORD:
		return LINT(v), nil
	case LWORD:
		return LINT(v), nil // Note: Max LWORD (uint64) > Max LINT
	case TIME:
		return LINT(time.Duration(v).Milliseconds()), nil
	case DATE:
		return LINT(time.Time(v).UnixMilli()), nil
	case TOD:
		t_time := time.Time(v)
		midnight := time.Date(t_time.Year(), t_time.Month(), t_time.Day(), 0, 0, 0, 0, t_time.Location())
		return LINT(t_time.Sub(midnight).Milliseconds()), nil
	case DT:
		return LINT(time.Time(v).UnixMilli()), nil
	case STRING:
		i, err := strconv.Atoi(string(v))
		if err != nil {
			return 0, fmt.Errorf("anyToLINT: cannot parse STRING '%s' to LINT: %w", v, err)
		}
		return LINT(i), nil
	default:
		return 0, fmt.Errorf("anyToLINT: unsupported type %T for conversion to LINT", val)
	}
}

// anyToULINT converts a supported PLC bitwise type to ULINT.
func anyToULINT(val interface{}) (ULINT, error) {
	switch v := val.(type) {
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

// convertToTargetType converts an accumulated value (LREAL or LINT) to the target PLC type.
func convertToTargetType(accumulator interface{}, targetType reflect.Type) (interface{}, error) {
	sourceIsLREAL := false
	var sourceLREAL LREAL
	var sourceLINT LINT
	var sourceULINT ULINT
	var originalAccumulator interface{} // To hold the original accumulator for direct return if types match

	// If the accumulator is already of the target type, return it directly.
	if reflect.TypeOf(accumulator) == targetType {
		return accumulator, nil
	}

	if accLREAL, ok := accumulator.(LREAL); ok {
		sourceIsLREAL = true
		sourceLREAL = accLREAL
	} else if accTIME, ok := accumulator.(TIME); ok {
		// TIME is handled directly without promoting to LREAL/LINT
		sourceIsLREAL = true
		sourceLREAL = LREAL(accTIME)
	} else if accLINT, ok := accumulator.(LINT); ok {
		sourceLINT = accLINT
	} else if accULINT, ok := accumulator.(ULINT); ok {
		sourceULINT = accULINT
	} else {
		return nil, fmt.Errorf("convertToTargetType: accumulator is of unhandled type %T, expected LREAL or LINT", accumulator)
	}

	// Handle potential Inf/NaN when converting LREAL to integer types
	if sourceIsLREAL && (math.IsInf(float64(sourceLREAL), 0) || math.IsNaN(float64(sourceLREAL))) {
		if isPlcIntType(targetType) || targetType == reflect.TypeOf(TIME(0)) ||
			targetType == reflect.TypeOf(DATE{}) || targetType == reflect.TypeOf(TOD{}) || targetType == reflect.TypeOf(DT{}) {
			// IEC 61131-3 doesn't strictly define this; common behavior is implementation-dependent.
			// Options: return error, return 0, return max/min int.
			// For now, returning an error is safest.
			return nil, fmt.Errorf("convertToTargetType: cannot convert LREAL Inf/NaN to target type %v", targetType)
		}
	}

	switch targetType {
	case reflect.TypeOf(SINT(0)):
		if sourceIsLREAL {
			return SINT(sourceLREAL), nil
		} else if sourceULINT != 0 {
			return SINT(sourceULINT), nil
		} else { // sourceLINT
			return SINT(sourceLINT), nil
		}
	case reflect.TypeOf(INT(0)):
		if sourceIsLREAL {
			return INT(sourceLREAL), nil
		} else if sourceULINT != 0 {
			return INT(sourceULINT), nil
		} else { // sourceLINT
			return INT(sourceLINT), nil
		}
	case reflect.TypeOf(DINT(0)):
		if sourceIsLREAL {
			return DINT(sourceLREAL), nil
		} else if sourceULINT != 0 {
			return DINT(sourceULINT), nil
		} else { // sourceLINT
			return DINT(sourceLINT), nil
		}
	case reflect.TypeOf(LINT(0)):
		if sourceIsLREAL {
			return LINT(sourceLREAL), nil
		} else if sourceULINT != 0 {
			return LINT(sourceULINT), nil
		} else { // sourceLINT
			return LINT(sourceLINT), nil
		}
	case reflect.TypeOf(USINT(0)):
		if sourceIsLREAL {
			return USINT(sourceLREAL), nil
		} else if sourceULINT != 0 {
			return USINT(sourceULINT), nil
		} else { // sourceLINT
			return USINT(sourceLINT), nil
		}
	case reflect.TypeOf(UINT(0)):
		if sourceIsLREAL {
			return UINT(sourceLREAL), nil
		} else if sourceULINT != 0 {
			return UINT(sourceULINT), nil
		} else { // sourceLINT
			return UINT(sourceLINT), nil
		}
	case reflect.TypeOf(UDINT(0)):
		if sourceIsLREAL {
			return UDINT(sourceLREAL), nil
		} else if sourceULINT != 0 {
			return UDINT(sourceULINT), nil
		} else { // sourceLINT
			return UDINT(sourceLINT), nil
		}
	case reflect.TypeOf(ULINT(0)):
		if sourceIsLREAL {
			return ULINT(sourceLREAL), nil
		} else if sourceULINT != 0 {
			return ULINT(sourceULINT), nil
		} else { // sourceLINT
			return ULINT(sourceLINT), nil
		}
	case reflect.TypeOf(REAL(0)):
		if sourceIsLREAL {
			return REAL(sourceLREAL), nil
		} else if sourceULINT != 0 {
			return REAL(sourceULINT), nil
		} else { // sourceLINT
			return REAL(sourceLINT), nil
		}
	case reflect.TypeOf(LREAL(0)):
		if sourceIsLREAL {
			return LREAL(sourceLREAL), nil
		} else if sourceULINT != 0 {
			return LREAL(sourceULINT), nil
		} else { // sourceLINT
			return LREAL(sourceLINT), nil
		}
	case reflect.TypeOf(BOOL(false)):
		if sourceIsLREAL {
			return BOOL(sourceLREAL != 0), nil
		} else if sourceULINT != 0 {
			return BOOL(sourceULINT != 0), nil
		} else { // sourceLINT
			return BOOL(sourceLINT != 0), nil
		}
	case reflect.TypeOf(BYTE(0)):
		if sourceIsLREAL {
			return SubByte(REAL(sourceLREAL))
		} else if sourceULINT != 0 {
			return SubByte(ULINT(sourceULINT))
		} else { // sourceLINT
			return SubByte(LINT(sourceLINT))
		}
	case reflect.TypeOf(WORD(0)):
		if sourceIsLREAL {
			return SubWord(REAL(sourceLREAL))
		} else if sourceULINT != 0 {
			return SubWord(ULINT(sourceULINT))
		} else { // sourceLINT
			return SubWord(LINT(sourceLINT))
		}
	case reflect.TypeOf(DWORD(0)):
		if sourceIsLREAL {
			return SubDword(REAL(sourceLREAL))
		} else if sourceULINT != 0 {
			return SubDword(ULINT(sourceULINT))
		} else { // sourceLINT
			return SubDword(LINT(sourceLINT))
		}
	case reflect.TypeOf(LWORD(0)):
		if sourceIsLREAL {
			return SubLword(REAL(sourceLREAL))
		} else if sourceULINT != 0 {
			return SubLword(ULINT(sourceULINT))
		} else { // sourceLINT
			return SubLword(LINT(sourceLINT))
		}
	case reflect.TypeOf(TIME(0)):
		if sourceIsLREAL {
			return SubTime(LREAL(sourceLREAL))
		} else if sourceULINT != 0 {
			return SubTime(ULINT(sourceULINT))
		} else { // sourceLINT
			return SubTime(LINT(sourceLINT | LINT(sourceULINT)))
		}
	case reflect.TypeOf(DATE(time.Time{})):
		if _, ok := originalAccumulator.(DATE); ok { // If original accumulator was DATE, return it directly
			return originalAccumulator, nil
		}
		if sourceIsLREAL {
			return SubDate(LREAL(sourceLREAL))
		} else if sourceULINT != 0 {
			return SubDate(ULINT(sourceULINT))
		} else { // sourceLINT
			return SubDate(LINT(sourceLINT))
		}
	case reflect.TypeOf(TOD(time.Time{})):
		if _, ok := originalAccumulator.(TOD); ok { // If original accumulator was TOD, return it directly
			return originalAccumulator, nil
		}
		if sourceIsLREAL {
			return SubTod(LREAL(sourceLREAL))
		} else if sourceULINT != 0 { // Added this for consistency
			return SubTod(ULINT(sourceULINT))
		} else {
			return SubTod(LINT(sourceLINT))
		}
	case reflect.TypeOf(DT(time.Time{})):
		if _, ok := originalAccumulator.(DT); ok { // If original accumulator was DT, return it directly
			return originalAccumulator, nil
		}
		if sourceIsLREAL {
			return SubDt(LREAL(sourceLREAL))
		} else if sourceULINT != 0 {
			return SubDt(ULINT(sourceULINT))
		} else { // sourceLINT
			return SubDt(LINT(sourceLINT))
		}
	case reflect.TypeOf(STRING("")):
		if s, ok := originalAccumulator.(fmt.Stringer); ok { // Use Stringer interface for time types
			return STRING(s.String()), nil
		} else if sourceIsLREAL {
			return STRING(fmt.Sprintf("%g", sourceLREAL)), nil
		} else if sourceULINT != 0 { // Changed to decimal for general ULINT
			return STRING(fmt.Sprintf("%d", sourceULINT)), nil
		} else { // sourceLINT
			return STRING(fmt.Sprintf("%d", sourceLINT)), nil
		}
	default:
		return nil, fmt.Errorf("convertToTargetType: unsupported target type %v", targetType)
	}
}

// getTypeRank returns a numerical rank for a given PLC type to determine dominance.
// Higher numbers have higher precedence.
func getTypeRank(t reflect.Type) int {
	switch t {
	// Floats
	case reflect.TypeOf(LREAL(0)):
		return 20
	case reflect.TypeOf(REAL(0)):
		return 19

	// 64-bit Integers
	case reflect.TypeOf(LINT(0)):
		return 18
	case reflect.TypeOf(ULINT(0)), reflect.TypeOf(LWORD(0)):
		return 17

	// 32-bit Integers
	case reflect.TypeOf(DINT(0)):
		return 16
	case reflect.TypeOf(UDINT(0)), reflect.TypeOf(DWORD(0)):
		return 15

	// 16-bit Integers
	case reflect.TypeOf(INT(0)):
		return 14
	case reflect.TypeOf(UINT(0)), reflect.TypeOf(WORD(0)):
		return 13

	// 8-bit Integers
	case reflect.TypeOf(SINT(0)):
		return 12
	case reflect.TypeOf(USINT(0)), reflect.TypeOf(BYTE(0)):
		return 11

	// Other
	case reflect.TypeOf(BOOL(false)):
		return 1
	default:
		return 0 // Includes STRING, TIME, DATE, etc. which have special handling
	}
}

// --- End Helper Functions ---

// ADD performs addition on a slice of mixed PLC data types.
// The result type is determined by the type of the last element in the slice.
// Intermediate calculations are performed using LREAL if any float-like type is present, otherwise LINT. It returns the result and an error if one occurs.
func ADD(nums []interface{}) (interface{}, error) {
	if len(nums) == 0 {
		// IEC defines the additive identity as 0. Return a default integer type.
		return 0, nil
	}

	// --- Special Case: Time Arithmetic (Table 30) ---
	// Check if the first argument is a TIME, TOD or DT to handle time arithmetic.
	switch first := nums[0].(type) {
	case TIME:
		acc := first
		for i := 1; i < len(nums); i++ {
			if duration, ok := nums[i].(TIME); ok {
				acc += duration
			} else {
				// If we encounter a non-TIME type, fall back to the generic numeric addition logic.
				goto numeric_add
			}
		}
		return acc, nil
	case TOD:
		acc := time.Time(first)
		for i := 1; i < len(nums); i++ {
			if duration, ok := nums[i].(TIME); ok {
				acc = acc.Add(time.Duration(duration))
			} else {
				return nil, fmt.Errorf("ADD: invalid type for addition with TOD; expected TIME, got %T", nums[i])
			}
		}
		return TOD(acc), nil
	case DT:
		acc := time.Time(first)
		for i := 1; i < len(nums); i++ {
			if duration, ok := nums[i].(TIME); ok {
				acc = acc.Add(time.Duration(duration))
			} else {
				return nil, fmt.Errorf("ADD: invalid type for addition with DT; expected TIME, got %T", nums[i])
			}
		}
		return DT(acc), nil
	}

numeric_add:
	// --- Default Case: Numeric and Duration Arithmetic ---
	if len(nums) == 0 {
		return LINT(0), nil // IEC defines the additive identity as 0.
	}

	useLREAL := false
	for _, num := range nums {
		rt := reflect.TypeOf(num)
		if isPlcFloatType(rt) {
			useLREAL = true
			break
		}
		if s, ok := num.(STRING); ok {
			// Try parsing as float to see if LREAL accumulation is needed
			if _, err := strconv.ParseFloat(string(s), 64); err == nil {
				useLREAL = true // If any string looks like a float, use LREAL
				// break // Don't break, another element might be an actual REAL/LREAL
			}
		}
	}

	var finalAccumulator interface{}
	targetType := reflect.TypeOf(nums[0]) // Start with the first type

	if useLREAL {
		var accLREAL LREAL
		for _, num := range nums {
			val, err := anyToLREAL(num)
			if err != nil {
				return nil, fmt.Errorf("ADD: error converting %v (type %T) to LREAL for accumulation: %w", num, num, err)
			}
			accLREAL += val
		}
		finalAccumulator = accLREAL
		// Determine the largest float type present
		for _, num := range nums {
			numType := reflect.TypeOf(num)
			if getTypeRank(numType) > getTypeRank(targetType) {
				targetType = numType
			}
		}
	} else { // All integer-like (or bools, time types as int, strings as int)
		var accLINT LINT
		for _, num := range nums {
			val, err := anyToLINT(num)
			if err != nil {
				return nil, fmt.Errorf("ADD: error converting %v (type %T) to LINT for accumulation: %w", num, num, err)
			}
			accLINT += val
		}
		finalAccumulator = accLINT

		// Determine the largest integer-like type present
		for i := 1; i < len(nums); i++ {
			numType := reflect.TypeOf(nums[i])
			if getTypeRank(numType) > getTypeRank(targetType) {
				targetType = numType
			}
		}
	}

	result, err := convertToTargetType(finalAccumulator, targetType)
	if err != nil {
		return nil, fmt.Errorf("ADD: error converting final sum to target type %v: %w. Accumulator was: %v", targetType, err, finalAccumulator)
	}
	return result, nil
}

// SUB performs subtraction on a slice of mixed PLC data types.
// nums[0] - nums[1] - nums[2]...
// The result type is determined by the type of the last element.
func SUB(nums []interface{}) (interface{}, error) {
	if len(nums) == 0 {
		// IEC doesn't define SUB for zero inputs. Returning 0 is a safe default.
		return 0, nil
	}
	if len(nums) == 1 {
		return nums[0], nil
	}

	// --- Special Case: Time Arithmetic (Table 30) ---
	// Check for specific time-based subtraction patterns.
	switch first := nums[0].(type) {
	case TIME:
		// Handle TIME - TIME - TIME ...
		acc := first
		for i := 1; i < len(nums); i++ {
			if duration, ok := nums[i].(TIME); ok {
				acc -= duration
			} else {
				// If we encounter a non-TIME type, fall back to the generic numeric logic.
				goto numeric_sub
			}
		}
		return acc, nil
	case DATE:
		// Handle DATE - DATE -> TIME
		if len(nums) == 2 {
			if date2, ok := nums[1].(DATE); ok {
				return TIME(time.Time(first).Sub(time.Time(date2))), nil
			}
		}
	case TOD:
		// Handle TOD - TIME -> TOD  OR  TOD - TOD -> TIME
		if len(nums) == 2 {
			if duration, ok := nums[1].(TIME); ok {
				return TOD(time.Time(first).Add(-time.Duration(duration))), nil
			}
			if tod2, ok := nums[1].(TOD); ok {
				return TIME(time.Time(first).Sub(time.Time(tod2))), nil
			}
		}
	case DT:
		// Handle DT - TIME -> DT  OR  DT - DT -> TIME
		if len(nums) == 2 {
			if duration, ok := nums[1].(TIME); ok {
				return DT(time.Time(first).Add(-time.Duration(duration))), nil
			}
			if dt2, ok := nums[1].(DT); ok {
				return TIME(time.Time(first).Sub(time.Time(dt2))), nil
			}
		}
	}

numeric_sub:
	// --- Default Case: Numeric Subtraction ---
	useLREAL := false
	for _, num := range nums {
		rt := reflect.TypeOf(num)
		if isPlcFloatType(rt) {
			useLREAL = true
			break
		}
		if s, ok := num.(STRING); ok {
			if _, err := strconv.ParseFloat(string(s), 64); err == nil {
				useLREAL = true
			}
		}
	}

	var finalAccumulator interface{}
	targetType := reflect.TypeOf(nums[0])
	for _, num := range nums {
		numType := reflect.TypeOf(num)
		// Promote target type based on rank
		if getTypeRank(numType) > getTypeRank(targetType) {
			targetType = numType
		}
	}

	// --- Special Case: Time Arithmetic (Table 30) ---

	if useLREAL {
		accLREAL, err := anyToLREAL(nums[0])
		if err != nil {
			return nil, fmt.Errorf("SUB: error converting first element %v (type %T) to LREAL: %w", nums[0], nums[0], err)
		}
		for i := 1; i < len(nums); i++ {
			val, err := anyToLREAL(nums[i])
			if err != nil {
				return nil, fmt.Errorf("SUB: error converting element %v (type %T) to LREAL: %w", nums[i], nums[i], err)
			}
			accLREAL -= val
		}
		finalAccumulator = accLREAL
	} else { // All integer-like
		accLINT, err := anyToLINT(nums[0])
		if err != nil {
			return nil, fmt.Errorf("SUB: error converting first element %v (type %T) to LINT: %w", nums[0], nums[0], err)
		}
		for i := 1; i < len(nums); i++ {
			val, err := anyToLINT(nums[i])
			if err != nil {
				return nil, fmt.Errorf("SUB: error converting element %v (type %T) to LINT: %w", nums[i], nums[i], err)
			}
			accLINT -= val
		}
		finalAccumulator = accLINT
	}

	result, err := convertToTargetType(finalAccumulator, targetType)
	if err != nil {
		return nil, fmt.Errorf("SUB: error converting final result to target type %v: %w. Accumulator was: %v", targetType, err, finalAccumulator)
	}
	return result, nil
}

// MUL performs multiplication on a slice of mixed PLC data types.
// The result type is determined by the type of the last element.
func MUL(nums []interface{}) (interface{}, error) {

	if len(nums) == 0 {
		// IEC defines the multiplicative identity as 1. Return a default integer type.
		return LINT(1), nil
	}

	// --- Special Case: Time Arithmetic (Table 30) ---
	if first, ok := nums[0].(TIME); ok {
		// Convert TIME to LREAL (milliseconds) for calculation to maintain consistency.
		acc, err := anyToLREAL(first)
		if err != nil {
			return nil, fmt.Errorf("MUL: error converting initial TIME value %v to LREAL: %w", first, err)
		}
		for i := 1; i < len(nums); i++ {
			multiplier, err := anyToLREAL(nums[i])
			if err != nil {
				return nil, fmt.Errorf("MUL: invalid type for multiplication with TIME; expected ANY_NUM, got %T at index %d", nums[i], i)
			}
			acc *= multiplier
		}
		result, err := convertToTargetType(LREAL(acc), reflect.TypeOf(TIME(0)))
		if err != nil {
			return nil, fmt.Errorf("MUL: error converting final product to target type TIME: %w. Accumulator was: %v", err, acc)
		}
		return result, nil
	} else {
		var finalAccumulator interface{}
		targetType := reflect.TypeOf(nums[0])
		for _, num := range nums {
			numType := reflect.TypeOf(num)
			if getTypeRank(numType) > getTypeRank(targetType) {
				targetType = numType
			}
		}
		// --- Default Case: Numeric Multiplication ---
		useLREAL := false
		for _, num := range nums {
			rt := reflect.TypeOf(num)
			if isPlcFloatType(rt) {
				useLREAL = true
				break
			}
			if s, ok := num.(STRING); ok {
				if _, err := strconv.ParseFloat(string(s), 64); err == nil {
					useLREAL = true
				}
			}
		}

		if useLREAL {
			var accLREAL LREAL = 1.0
			for _, num := range nums {
				val, err := anyToLREAL(num)
				if err != nil {
					return nil, fmt.Errorf("MUL: error converting %v (type %T) to LREAL: %w", num, num, err)
				}
				accLREAL *= val
			}
			finalAccumulator = accLREAL
		} else { // All integer-like
			var accLINT LINT = 1
			for _, num := range nums {
				val, err := anyToLINT(num)
				if err != nil {
					return nil, fmt.Errorf("MUL: error converting %v (type %T) to LINT: %w", num, num, err)
				}
				accLINT *= val
			}
			finalAccumulator = accLINT
		}
		result, err := convertToTargetType(finalAccumulator, targetType)
		if err != nil {
			return nil, fmt.Errorf("MUL: error converting final product to target type %v: %w. Accumulator was: %v", targetType, err, finalAccumulator)
		}
		return result, nil
	}
}

// DIV performs division on a slice of mixed PLC data types.
// (...(nums[0] / nums[1]) / nums[2] ...)
// DIV(TIME, ANY_NUM): The result data type is TIME.
// DIV(TIME, TIME): The result data type is REAL (or LREAL in your implementation, given its higher precision).
// The result type is determined by the type of the last element.
func DIV(nums []interface{}) (interface{}, error) {

	if len(nums) < 1 {
		// IEC defines division by a single number as the number itself.
		// Returning a default integer type for empty slice.
		return LINT(0), nil
	} else if len(nums) == 1 {
		// Division of a single element is just the element itself.
		return nums[0], nil
	}

	// --- Special Case: Time Arithmetic (Table 30) ---
	if first, ok := nums[0].(TIME); ok {
		// Handle TIME / TIME -> LREAL
		if len(nums) == 2 { // This rule applies only for exactly two arguments
			if secondTime, ok := nums[1].(TIME); ok { // Check if second argument is TIME
				if time.Duration(secondTime) == 0 {
					return nil, fmt.Errorf("DIV: division by zero (TIME / TIME)")
				}
				finalAccumulator := LREAL(float64(first) / float64(secondTime))
				// Per standard, TIME/TIME -> REAL. We use LREAL for precision.
				result, err := convertToTargetType(finalAccumulator, reflect.TypeOf(LREAL(0)))
				if err != nil {
					return nil, fmt.Errorf("DIV: error converting final result to target type LREAL: %w. Accumulator was: %v", err, finalAccumulator)
				}
				return result, nil
			}
		}
		// Handle TIME / ANY_NUM -> TIME
		acc, err := anyToLREAL(first) // Convert to milliseconds as LREAL
		if err != nil {
			return nil, fmt.Errorf("DIV: error converting initial TIME value %v to LREAL: %w", first, err)
		}
		for i := 1; i < len(nums); i++ {
			divisor, err := anyToLREAL(nums[i])
			if err != nil {
				return nil, fmt.Errorf("DIV: invalid type for division with TIME; expected ANY_NUM, got %T", nums[i])
			}
			if divisor == 0.0 {
				return nil, fmt.Errorf("DIV: division by zero with TIME operand")
			}
			acc /= divisor
		}
		result, err := convertToTargetType(LREAL(acc), reflect.TypeOf(TIME(0)))
		if err != nil {
			return nil, fmt.Errorf("DIV: error converting final result to target type TIME: %w. Accumulator was: %v", err, acc)
		}
		return result, nil
	} else {
		// --- Default Case: Numeric Division ---
		//var finalAccumulator interface{}
		targetType := reflect.TypeOf(nums[0])
		for _, num := range nums {
			numType := reflect.TypeOf(num)
			if getTypeRank(numType) > getTypeRank(targetType) {
				targetType = numType
			}
		}
	}

	useLREAL := false
	for _, num := range nums {
		rt := reflect.TypeOf(num)
		if isPlcFloatType(rt) {
			useLREAL = true
			break
		}
		if s, ok := num.(STRING); ok {
			if _, err := strconv.ParseFloat(string(s), 64); err == nil {
				useLREAL = true
			}
		}
	}

	var finalAccumulator interface{}
	targetType := reflect.TypeOf(nums[0])
	if useLREAL {
		accLREAL, err := anyToLREAL(nums[0])
		if err != nil {
			return nil, fmt.Errorf("DIV: error converting first element %v (type %T) to LREAL: %w", nums[0], nums[0], err)
		}
		for i := 1; i < len(nums); i++ {
			divisor, err := anyToLREAL(nums[i])
			if err != nil {
				return nil, fmt.Errorf("DIV: error converting divisor %v (type %T) to LREAL: %w", nums[i], nums[i], err)
			}
			// LREAL division by zero yields Inf/NaN, handled by IEEE 754.
			accLREAL /= divisor
		}
		finalAccumulator = accLREAL
	} else { // All integer-like
		accLINT, err := anyToLINT(nums[0])
		if err != nil {
			return nil, fmt.Errorf("DIV: error converting first element %v (type %T) to LINT: %w", nums[0], nums[0], err)
		}
		for i := 1; i < len(nums); i++ {
			divisor, err := anyToLINT(nums[i])
			if err != nil {
				return nil, fmt.Errorf("DIV: error converting divisor %v (type %T) to LINT: %w", nums[i], nums[i], err)
			}
			if divisor == 0 {
				return nil, fmt.Errorf("DIV: division by zero (integer context, type %T)", accLINT)
			}
			accLINT /= divisor
		}
		finalAccumulator = accLINT
	}

	result, err := convertToTargetType(finalAccumulator, targetType)
	if err != nil {
		return nil, fmt.Errorf("DIV: error converting final result to target type %v: %w. Accumulator was: %v", targetType, err, finalAccumulator)
	}
	return result, nil

}

// MOD performs modulo on a slice of mixed integer-like PLC data types.
// (...(nums[0] % nums[1]) % nums[2] ...)
// The result type is determined by the type of the last element.
// All inputs and the target type must be integer-like types (actual integers, BOOL, BYTE, WORD, etc., or strings parsable to int).
func MOD(nums []interface{}) (interface{}, error) {

	if len(nums) < 2 {
		return nil, fmt.Errorf("MOD: input slice must have at least two elements")
	}

	targetType := reflect.TypeOf(nums[0])
	for _, num := range nums {
		numType := reflect.TypeOf(num)
		// Promote target type based on rank
		if getTypeRank(numType) > getTypeRank(targetType) {
			targetType = numType
		}
	}

	// Target for MOD must be integer-like or a time type that can be represented as an integer.
	if !isPlcIntType(targetType) && targetType != reflect.TypeOf(TIME(0)) &&
		targetType != reflect.TypeOf(DATE{}) && targetType != reflect.TypeOf(TOD{}) &&
		targetType != reflect.TypeOf(DT{}) { // Time types are stored as int millis
		return nil, fmt.Errorf("MOD: target type %v must be an integer-like or time type", targetType)
	}

	// All inputs must be convertible to LINT for MOD.
	// Floats or strings that parse to floats are not allowed for MOD's accumulation.
	for i, num := range nums {
		if _, ok := num.(REAL); ok {
			return nil, fmt.Errorf("MOD: input element at index %d (%v, type REAL) is not allowed for MOD operation", i, num)
		}
		if _, ok := num.(LREAL); ok {
			return nil, fmt.Errorf("MOD: input element at index %d (%v, type LREAL) is not allowed for MOD operation", i, num)
		}
		if s, ok := num.(STRING); ok {
			// Attempt to parse the string as a float.
			parsedFloat, errFloat := strconv.ParseFloat(string(s), 64)
			if errFloat == nil { // Successfully parsed as a float.
				// Now, check if it's also a clean integer representation.
				valInt, errInt := strconv.ParseInt(string(s), 10, 64)
				if errInt != nil || float64(LINT(valInt)) != parsedFloat {
					// If it parses as float but not cleanly as an int, it's problematic for MOD
					return nil, fmt.Errorf("MOD: input STRING element at index %d ('%s') parses as float but not as a clean integer, not suitable for integer MOD", i, s)
				}
			}
		}
	}

	accLINT, err := anyToLINT(nums[0])
	if err != nil {
		return nil, fmt.Errorf("MOD: error converting first element %v (type %T) to LINT: %w", nums[0], nums[0], err)
	}

	for i := 1; i < len(nums); i++ {
		divisor, err := anyToLINT(nums[i])
		if err != nil {
			return nil, fmt.Errorf("MOD: error converting divisor %v (type %T) to LINT: %w", nums[i], nums[i], err)
		}
		if divisor == 0 {
			return nil, fmt.Errorf("MOD: modulo by zero (integer context, type %T)", accLINT)
		}
		accLINT %= divisor
	}

	result, err := convertToTargetType(accLINT, targetType)
	if err != nil {
		return nil, fmt.Errorf("MOD: error converting final result to target type %v: %w. Accumulator was: %v", targetType, err, accLINT)
	}
	return result, nil
}

// MOVE performs an assignment of the input value.
// The standard defines this as a non-extensible function with one input and one output.
func MOVE(in interface{}) interface{} {
	// In this implementation, the function simply returns the input value.
	// The calling context would perform the assignment, e.g., `myVar := MOVE(otherVar)`.
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

// SUB_TOD subtracts a TIME from a TIME_OF_DAY, or two TIME_OF_DAYs.
// The return type depends on the inputs.
func SUB_TOD(in1, in2 interface{}) (interface{}, error) {
	t1, ok1 := in1.(TOD)
	if !ok1 {
		return nil, fmt.Errorf("SUB_TOD: first input must be of type TOD, got %T", in1)
	}

	// Case 1: TOD - TIME -> TOD
	if t2, ok2 := in2.(TIME); ok2 {
		return TOD(time.Time(t1).Add(-time.Duration(t2))), nil
	}

	// Case 2: TOD - TOD -> TIME
	if t2, ok2 := in2.(TOD); ok2 {
		return TIME(time.Time(t1).Sub(time.Time(t2))), nil
	}

	return nil, fmt.Errorf("SUB_TOD: second input must be of type TIME or TOD, got %T", in2)
}

// SUB_DT subtracts a TIME from a DATE_AND_TIME, or two DATE_AND_TIMEs.
// The return type depends on the inputs.
func SUB_DT(in1, in2 interface{}) (interface{}, error) {
	t1, ok1 := in1.(DT)
	if !ok1 {
		return nil, fmt.Errorf("SUB_DT: first input must be of type DT, got %T", in1)
	}

	// Case 1: DT - TIME -> DT
	if t2, ok2 := in2.(TIME); ok2 {
		return DT(time.Time(t1).Add(-time.Duration(t2))), nil
	}

	// Case 2: DT - DT -> TIME
	if t2, ok2 := in2.(DT); ok2 {
		return TIME(time.Time(t1).Sub(time.Time(t2))), nil
	}

	return nil, fmt.Errorf("SUB_DT: second input must be of type TIME or DT, got %T", in2)
}

// MUL_TIME multiplies a TIME duration by a numeric value.
func MUL_TIME(in1 TIME, in2 interface{}) (TIME, error) {
	val, err := anyToLREAL(in2)
	if err != nil {
		return 0, fmt.Errorf("MUL_TIME: error converting multiplier %v to LREAL: %w", in2, err)
	}
	return TIME(float64(in1) * float64(val)), nil
}

// DIV_TIME divides a TIME duration by a numeric value.
func DIV_TIME(in1 TIME, in2 interface{}) (TIME, error) {
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
