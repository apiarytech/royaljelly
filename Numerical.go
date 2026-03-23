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
)

/*********************************/
/* IEC 61131-3 Standard Functions*/
/*********************************/

// ABS returns the absolute value of the input. Overloaded for ANY_NUM.
func ABS(in interface{}) (interface{}, error) {
	if _, ok := in.(STRING); ok {
		return nil, fmt.Errorf("ABS: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_INT or STRING_TO_REAL")
	}

	if isPlcFloatType(reflect.TypeOf(in)) {
		val, err := anyToLREAL(in)
		if err != nil {
			return nil, fmt.Errorf("ABS: error converting %v to LREAL: %w", in, err)
		}
		if val < 0 {
			val = -val
		}
		return convertToTargetType(val, reflect.TypeOf(in))
	}

	// Integer path - check if unsigned
	if _, ok := in.(ULINT); ok || reflect.TypeOf(in).Kind() == reflect.Uint64 { // Check for unsigned types
		// For unsigned types, ABS is a no-op, just return the input.
		return in, nil
	} else {
		val, err := anyToLINT(in)
		if err != nil {
			return nil, fmt.Errorf("ABS: error converting %v to LINT: %w", in, err)
		}
		if val < 0 {
			val = -val
		}
		return convertToTargetType(val, reflect.TypeOf(in))
	}
}

// SQRT returns the square root of the input. Overloaded for ANY_REAL.
func SQRT(in interface{}) (interface{}, error) {
	if _, ok := in.(STRING); ok {
		return nil, fmt.Errorf("SQRT: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_REAL")
	}

	val, err := anyToLREAL(in)
	if err != nil {
		return nil, fmt.Errorf("SQRT: error converting %v to LREAL: %w", in, err)
	}
	resVal := LREAL(math.Sqrt(float64(val)))
	return convertToTargetType(resVal, reflect.TypeOf(in))
}

// LN returns the natural logarithm of the input. Overloaded for ANY_REAL.
func LN(in interface{}) (interface{}, error) {
	if _, ok := in.(STRING); ok {
		return nil, fmt.Errorf("LN: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_REAL")
	}

	val, err := anyToLREAL(in)
	if err != nil {
		return nil, fmt.Errorf("LN: error converting %v to LREAL: %w", in, err)
	}
	resVal := LREAL(math.Log(float64(val)))
	return convertToTargetType(resVal, reflect.TypeOf(in))
}

// LOG returns the base 10 logarithm of the input. Overloaded for ANY_REAL.
func LOG(in interface{}) (interface{}, error) {
	if _, ok := in.(STRING); ok {
		return nil, fmt.Errorf("LOG: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_REAL")
	}

	val, err := anyToLREAL(in)
	if err != nil {
		return nil, fmt.Errorf("LOG: error converting %v to LREAL: %w", in, err)
	}
	resVal := LREAL(math.Log10(float64(val)))
	return convertToTargetType(resVal, reflect.TypeOf(in))
}

// EXP returns e**IN. Overloaded for ANY_REAL.
func EXP(in interface{}) (interface{}, error) {
	if _, ok := in.(STRING); ok {
		return nil, fmt.Errorf("EXP: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_REAL")
	}

	val, err := anyToLREAL(in)
	if err != nil {
		return nil, fmt.Errorf("EXP: error converting %v to LREAL: %w", in, err)
	}
	resVal := LREAL(math.Exp(float64(val)))
	return convertToTargetType(resVal, reflect.TypeOf(in))
}

// SIN returns the sine of the input. Overloaded for ANY_REAL.
func SIN(in interface{}) (interface{}, error) {
	if _, ok := in.(STRING); ok {
		return nil, fmt.Errorf("SIN: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_REAL")
	}

	val, err := anyToLREAL(in)
	if err != nil {
		return nil, fmt.Errorf("SIN: error converting %v to LREAL: %w", in, err)
	}
	resVal := LREAL(math.Sin(float64(val)))
	return convertToTargetType(resVal, reflect.TypeOf(in))
}

// COS returns the cosine of the input. Overloaded for ANY_REAL.
func COS(in interface{}) (interface{}, error) {
	if _, ok := in.(STRING); ok {
		return nil, fmt.Errorf("COS: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_REAL")
	}

	val, err := anyToLREAL(in)
	if err != nil {
		return nil, fmt.Errorf("COS: error converting %v to LREAL: %w", in, err)
	}
	resVal := LREAL(math.Cos(float64(val)))
	return convertToTargetType(resVal, reflect.TypeOf(in))
}

// TAN returns the tangent of the input. Overloaded for ANY_REAL.
func TAN(in interface{}) (interface{}, error) {
	if _, ok := in.(STRING); ok {
		return nil, fmt.Errorf("TAN: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_REAL")
	}

	val, err := anyToLREAL(in)
	if err != nil {
		return nil, fmt.Errorf("TAN: error converting %v to LREAL: %w", in, err)
	}
	resVal := LREAL(math.Tan(float64(val)))
	return convertToTargetType(resVal, reflect.TypeOf(in))
}

// ASIN returns the arcsine of the input. Overloaded for ANY_REAL.
func ASIN(in interface{}) (interface{}, error) {
	if _, ok := in.(STRING); ok {
		return nil, fmt.Errorf("ASIN: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_REAL")
	}

	val, err := anyToLREAL(in)
	if err != nil {
		return nil, fmt.Errorf("ASIN: error converting %v to LREAL: %w", in, err)
	}
	resVal := LREAL(math.Asin(float64(val)))
	return convertToTargetType(resVal, reflect.TypeOf(in))
}

// ACOS returns the arccosine of the input. Overloaded for ANY_REAL.
func ACOS(in interface{}) (interface{}, error) {
	if _, ok := in.(STRING); ok {
		return nil, fmt.Errorf("ACOS: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_REAL")
	}

	val, err := anyToLREAL(in)
	if err != nil {
		return nil, fmt.Errorf("ACOS: error converting %v to LREAL: %w", in, err)
	}
	resVal := LREAL(math.Acos(float64(val)))
	return convertToTargetType(resVal, reflect.TypeOf(in))
}

// ATAN returns the arctangent of the input. Overloaded for ANY_REAL.
func ATAN(in interface{}) (interface{}, error) {
	if _, ok := in.(STRING); ok {
		return nil, fmt.Errorf("ATAN: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_REAL")
	}

	val, err := anyToLREAL(in)
	if err != nil {
		return nil, fmt.Errorf("ATAN: error converting %v to LREAL: %w", in, err)
	}
	resVal := LREAL(math.Atan(float64(val)))
	return convertToTargetType(resVal, reflect.TypeOf(in))
}

// EXPT performs exponentiation (IN1**IN2). Overloaded for ANY_REAL(IN1) and ANY_NUM(IN2).
func EXPT(in1, in2 interface{}) (interface{}, error) {
	if _, ok := in1.(STRING); ok {
		return nil, fmt.Errorf("EXPT: implicit string conversion is not supported for base. Use explicit conversion functions")
	}
	if _, ok := in2.(STRING); ok {
		return nil, fmt.Errorf("EXPT: implicit string conversion is not supported for exponent. Use explicit conversion functions")
	}

	base, err1 := anyToLREAL(in1)
	if err1 != nil {
		return nil, fmt.Errorf("EXPT: error converting base %v to LREAL: %w", in1, err1)
	}
	exponent, err2 := anyToLREAL(in2) // Promote exponent to LREAL for math.Pow
	if err2 != nil {
		return nil, fmt.Errorf("EXPT: error converting exponent %v to LREAL: %w", in2, err2)
	}
	resVal := LREAL(math.Pow(float64(base), float64(exponent)))
	return convertToTargetType(resVal, reflect.TypeOf(in1)) // Return type matches first input
}

// TRUNC truncates a real number to an integer. Overloaded for ANY_REAL to ANY_INT.
// As per IEC 61131-3, the result of truncating a REAL is a DINT.
func TRUNC(in interface{}) (DINT, error) {
	if _, ok := in.(STRING); ok {
		return 0, fmt.Errorf("TRUNC: implicit string conversion is not supported. Use explicit conversion functions like STRING_TO_REAL")
	}
	if !isPlcFloatType(reflect.TypeOf(in)) {
		return 0, fmt.Errorf("TRUNC: input must be of type ANY_REAL, got %T", in)
	}
	val, err := anyToLREAL(in)
	if err != nil {
		return 0, fmt.Errorf("TRUNC: error converting %v to LREAL: %w", in, err)
	}
	return DINT(val), nil // LREAL to DINT performs truncation
}

/*********************************/
/* Non-Standard Helper Functions */
/*********************************/

// SUMLINT adds together the values of m.
// Deprecated: Use the generic SUM function instead.
func SUMLINT(m map[STRING]LINT) LINT {
	var s LINT
	for _, v := range m {
		s += v
	}
	return s
}

// SUMREAL adds together the values of m.
// Deprecated: Use the generic SUM function instead.
func SUMREAL(m map[STRING]REAL) REAL {
	var s REAL
	for _, v := range m {
		s += v
	}
	return s
}

// SUMLREAL adds together the values of m.
// Deprecated: Use the generic SUM function instead.
func SUMLREAL(m map[STRING]LREAL) LREAL {
	var s LREAL
	for _, v := range m {
		s += v
	}
	return s
}

// SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// as types for map values.
// Deprecated: Use the generic SUM function instead.
func SUMLINTorLREAL[K comparable, V LINT | LREAL](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// SumNumbers sums the values of map m. It supports both integers
// and floats as map values.
func SUM[K comparable, V ANY_NUM](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
