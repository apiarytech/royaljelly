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

import (
	"fmt"
	"math"
	"reflect"
)

/*********************************/
/* IEC 61131-3 Standard Functions*/
/*********************************/

// ABS returns the absolute value of the input. Overloaded for ANY_NUM.
func ABS(in interface{}) interface{} {
	if isPlcFloatType(reflect.TypeOf(in)) {
		val, err := anyToLREAL(in)
		if err != nil {
			panic(fmt.Sprintf("ABS: error converting %v to LREAL: %v", in, err))
		}
		if val < 0 {
			val = -val
		}
		res, _ := convertToTargetType(val, reflect.TypeOf(in))
		return res
	}

	// Integer path - check if unsigned
	if _, ok := in.(ULINT); ok || reflect.TypeOf(in).Kind() == reflect.Uint64 { // Check for unsigned types
		// For unsigned types, ABS is a no-op, just return the input.
		return in
	} else {
		val, err := anyToLINT(in)
		if err != nil {
			panic(fmt.Sprintf("ABS: error converting %v to LINT: %v", in, err))
		}
		if val < 0 {
			val = -val
		}
		res, _ := convertToTargetType(val, reflect.TypeOf(in))
		return res
	}
}

// SQRT returns the square root of the input. Overloaded for ANY_REAL.
func SQRT(in interface{}) interface{} {
	val, err := anyToLREAL(in)
	if err != nil {
		panic(fmt.Sprintf("SQRT: error converting %v to LREAL: %v", in, err))
	}
	resVal := LREAL(math.Sqrt(float64(val)))
	res, _ := convertToTargetType(resVal, reflect.TypeOf(in))
	return res
}

// LN returns the natural logarithm of the input. Overloaded for ANY_REAL.
func LN(in interface{}) interface{} {
	val, err := anyToLREAL(in)
	if err != nil {
		panic(fmt.Sprintf("LN: error converting %v to LREAL: %v", in, err))
	}
	resVal := LREAL(math.Log(float64(val)))
	res, _ := convertToTargetType(resVal, reflect.TypeOf(in))
	return res
}

// LOG returns the base 10 logarithm of the input. Overloaded for ANY_REAL.
func LOG(in interface{}) interface{} {
	val, err := anyToLREAL(in)
	if err != nil {
		panic(fmt.Sprintf("LOG: error converting %v to LREAL: %v", in, err))
	}
	resVal := LREAL(math.Log10(float64(val)))
	res, _ := convertToTargetType(resVal, reflect.TypeOf(in))
	return res
}

// EXP returns e**IN. Overloaded for ANY_REAL.
func EXP(in interface{}) interface{} {
	val, err := anyToLREAL(in)
	if err != nil {
		panic(fmt.Sprintf("EXP: error converting %v to LREAL: %v", in, err))
	}
	resVal := LREAL(math.Exp(float64(val)))
	res, _ := convertToTargetType(resVal, reflect.TypeOf(in))
	return res
}

// SIN returns the sine of the input. Overloaded for ANY_REAL.
func SIN(in interface{}) interface{} {
	val, err := anyToLREAL(in)
	if err != nil {
		panic(fmt.Sprintf("SIN: error converting %v to LREAL: %v", in, err))
	}
	resVal := LREAL(math.Sin(float64(val)))
	res, _ := convertToTargetType(resVal, reflect.TypeOf(in))
	return res
}

// COS returns the cosine of the input. Overloaded for ANY_REAL.
func COS(in interface{}) interface{} {
	val, err := anyToLREAL(in)
	if err != nil {
		panic(fmt.Sprintf("COS: error converting %v to LREAL: %v", in, err))
	}
	resVal := LREAL(math.Cos(float64(val)))
	res, _ := convertToTargetType(resVal, reflect.TypeOf(in))
	return res
}

// TAN returns the tangent of the input. Overloaded for ANY_REAL.
func TAN(in interface{}) interface{} {
	val, err := anyToLREAL(in)
	if err != nil {
		panic(fmt.Sprintf("TAN: error converting %v to LREAL: %v", in, err))
	}
	resVal := LREAL(math.Tan(float64(val)))
	res, _ := convertToTargetType(resVal, reflect.TypeOf(in))
	return res
}

// ASIN returns the arcsine of the input. Overloaded for ANY_REAL.
func ASIN(in interface{}) interface{} {
	val, err := anyToLREAL(in)
	if err != nil {
		panic(fmt.Sprintf("ASIN: error converting %v to LREAL: %v", in, err))
	}
	resVal := LREAL(math.Asin(float64(val)))
	res, _ := convertToTargetType(resVal, reflect.TypeOf(in))
	return res
}

// ACOS returns the arccosine of the input. Overloaded for ANY_REAL.
func ACOS(in interface{}) interface{} {
	val, err := anyToLREAL(in)
	if err != nil {
		panic(fmt.Sprintf("ACOS: error converting %v to LREAL: %v", in, err))
	}
	resVal := LREAL(math.Acos(float64(val)))
	res, _ := convertToTargetType(resVal, reflect.TypeOf(in))
	return res
}

// ATAN returns the arctangent of the input. Overloaded for ANY_REAL.
func ATAN(in interface{}) interface{} {
	val, err := anyToLREAL(in)
	if err != nil {
		panic(fmt.Sprintf("ATAN: error converting %v to LREAL: %v", in, err))
	}
	resVal := LREAL(math.Atan(float64(val)))
	res, _ := convertToTargetType(resVal, reflect.TypeOf(in))
	return res
}

// EXPT performs exponentiation (IN1**IN2). Overloaded for ANY_REAL(IN1) and ANY_NUM(IN2).
func EXPT(in1, in2 interface{}) interface{} {
	base, err1 := anyToLREAL(in1)
	if err1 != nil {
		panic(fmt.Sprintf("EXPT: error converting base %v to LREAL: %v", in1, err1))
	}
	exponent, err2 := anyToLREAL(in2) // Promote exponent to LREAL for math.Pow
	if err2 != nil {
		panic(fmt.Sprintf("EXPT: error converting exponent %v to LREAL: %v", in2, err2))
	}
	resVal := LREAL(math.Pow(float64(base), float64(exponent)))
	res, _ := convertToTargetType(resVal, reflect.TypeOf(in1)) // Return type matches first input
	return res
}

// TRUNC truncates a real number to an integer. Overloaded for ANY_REAL to ANY_INT.
// The specific integer type is determined by the context, here we return LINT as a general form.
// For a version that returns a specific integer type, a typed function like TRUNC_TO_INT would be needed.
func TRUNC(in interface{}) LINT {
	if !isPlcFloatType(reflect.TypeOf(in)) {
		panic(fmt.Sprintf("TRUNC: input must be of type ANY_REAL, got %T", in))
	}
	val, err := anyToLREAL(in)
	if err != nil {
		panic(fmt.Sprintf("TRUNC: error converting %v to LREAL: %v", in, err))
	}
	return LINT(val) // LREAL to LINT performs truncation
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
