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

package numerical

import (
	"math"

	"github.com/apiarytech/royaljelly/iec"
)

/*********************************/
/* IEC 61131-3 Standard Functions*/
/*********************************/

// ABS returns the absolute value of the input. Overloaded for signed numeric types.
func ABS[T iec.ANY_SINTS | iec.ANY_REAL](in T) T {
	if in < 0 {
		return -in
	}
	return in
}

// SQRT returns the square root of the input. Overloaded for ANY_REAL.
func SQRT[T iec.ANY_REAL](in T) T {
	return T(math.Sqrt(float64(in)))
}

// LN returns the natural logarithm of the input. Overloaded for ANY_REAL.
func LN[T iec.ANY_REAL](in T) T {
	return T(math.Log(float64(in)))
}

// LOG returns the base 10 logarithm of the input. Overloaded for ANY_REAL.
func LOG[T iec.ANY_REAL](in T) T {
	return T(math.Log10(float64(in)))
}

// EXP returns e**IN. Overloaded for ANY_REAL.
func EXP[T iec.ANY_REAL](in T) T {
	return T(math.Exp(float64(in)))
}

// SIN returns the sine of the input. Overloaded for ANY_REAL.
func SIN[T iec.ANY_REAL](in T) T {
	return T(math.Sin(float64(in)))
}

// COS returns the cosine of the input. Overloaded for ANY_REAL.
func COS[T iec.ANY_REAL](in T) T {
	return T(math.Cos(float64(in)))
}

// TAN returns the tangent of the input. Overloaded for ANY_REAL.
func TAN[T iec.ANY_REAL](in T) T {
	return T(math.Tan(float64(in)))
}

// ASIN returns the arcsine of the input. Overloaded for ANY_REAL.
func ASIN[T iec.ANY_REAL](in T) T {
	return T(math.Asin(float64(in)))
}

// ACOS returns the arccosine of the input. Overloaded for ANY_REAL.
func ACOS[T iec.ANY_REAL](in T) T {
	return T(math.Acos(float64(in)))
}

// ATAN returns the arctangent of the input. Overloaded for ANY_REAL.
func ATAN[T iec.ANY_REAL](in T) T {
	return T(math.Atan(float64(in)))
}

// EXPT performs exponentiation (IN1**IN2). Overloaded for ANY_REAL(IN1) and ANY_NUM(IN2).
func EXPT[T1, T2 iec.ANY_REAL](base T1, exponent T2) T1 {
	// The result type is the same as the base type (ANY_REAL).
	return T1(math.Pow(float64(base), float64(exponent)))
}

// TRUNC truncates a real number to an integer. Overloaded for ANY_REAL to ANY_INT.
// As per IEC 61131-3, the result of truncating a REAL is a DINT.
func TRUNC[T iec.ANY_REAL](in T) iec.DINT {
	return iec.DINT(in)
}

/*********************************/
/* Non-Standard Helper Functions */
/*********************************/

// SUMLINT adds together the values of m.
// Deprecated: Use the generic SUM function instead.
func SUMLINT(m map[iec.STRING]iec.LINT) iec.LINT {
	var s iec.LINT
	for _, v := range m {
		s += v
	}
	return s
}

// SUMREAL adds together the values of m.
// Deprecated: Use the generic SUM function instead.
func SUMREAL(m map[iec.STRING]iec.REAL) iec.REAL {
	var s iec.REAL
	for _, v := range m {
		s += v
	}
	return s
}

// SUMLREAL adds together the values of m.
// Deprecated: Use the generic SUM function instead.
func SUMLREAL(m map[iec.STRING]iec.LREAL) iec.LREAL {
	var s iec.LREAL
	for _, v := range m {
		s += v
	}
	return s
}

// SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// as types for map values.
// Deprecated: Use the generic SUM function instead.
func SUMLINTorLREAL[K comparable, V iec.LINT | iec.LREAL](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// SumNumbers sums the values of map m. It supports both integers
// and floats as map values.
func SUM[K comparable, V iec.ANY_NUM](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
