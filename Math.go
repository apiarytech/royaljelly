/* Copyright (C) 2026 Franklin D. Amador
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
	"math/rand"
	"time"
)

// globalRand is a package-level random number generator.
// It is seeded once during package initialization to ensure different
// execution runs produce different random sequences.
var globalRand *rand.Rand

func init() {
	// Initialize the global random generator with a seed based on the current time.
	// Using rand.NewSource is preferred over the global rand.Seed for creating
	// a generator that is not shared with other packages.
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

/*********************************/
/* IEC 61131-3 Standard Functions*/
/*********************************/

// NAN not a number
var NAN float64 = math.NaN()

// INF infinity
var POSINF float64 = math.Inf(1)
var NEGINF float64 = math.Inf(-1)

/*********************************/
/* Go 'math' Package Wrappers    */
/* (Non-standard extensions)     */
/*********************************/

// ACOS returns the arccosine, in radians, of in1.
func ACOS_LREAL(in1 LREAL) LREAL {
	return LREAL(math.Acos(float64(in1)))
}

// ACOSH returns the arccosine, in radians, of in1.
func ACOSH(in1 LREAL) LREAL {
	return LREAL(math.Acosh(float64(in1)))
}

// ASIN returns the arcsine, in radians, of in1.
func ASIN_LREAL(in1 LREAL) LREAL {
	return LREAL(math.Asin(float64(in1)))
}

// ASIN returns the arcsine, in radians, of in1.
func ASINH(in1 LREAL) LREAL {
	return LREAL(math.Asinh(float64(in1)))
}

// ATAN returns the arctangent, in radians, of in1.
func ATAN_LREAL(in1 LREAL) LREAL {
	return LREAL(math.Atan(float64(in1)))
}

// ATAN returns the arctangent, in radians, of in1.
func ATAN2(in1, in2 LREAL) LREAL {
	return LREAL(math.Atan2(float64(in1), float64(in2)))
}

// CBRT returns the cube root of in1.
func CBRT(in1 LREAL) LREAL {
	return LREAL(math.Cbrt(float64(in1)))
}

// CEIL returns the least integer value greater than or equal to in1.
func CEIL(in1 LREAL) LREAL {
	return LREAL(math.Ceil(float64(in1)))
}

// COPYSIGN returns a value with the magnitude of in1 and the sign of sign.
func COPYSIGN(in1, sign LREAL) LREAL {
	return LREAL(math.Copysign(float64(in1), float64(sign)))
}

// COS returns the arccosine, in radians, of in1.
func COS_LREAL(in1 LREAL) LREAL {
	return LREAL(math.Cos(float64(in1)))
}

// COSH returns the arccosine, in radians, of in1.
func COSH(in1 LREAL) LREAL {
	return LREAL(math.Cosh(float64(in1)))
}

// DIM returns the maximum of in1-in2 or 0.
func DIM(in1, in2 LREAL) LREAL {
	return LREAL(math.Dim(float64(in1), float64(in2)))
}

// ERF returns the erro fucntion of in1
func ERF(in1 LREAL) LREAL {
	return LREAL(math.Erf(float64(in1)))
}

// ERFC returns the erro fucntion of in1
func ERFC(in1 LREAL) LREAL {
	return LREAL(math.Erfc(float64(in1)))
}

// ERFCINV returns the inverse of Erfc(x).
func ERFCINV(in1 LREAL) LREAL {
	return LREAL(math.Erfcinv(float64(in1)))
}

// ERFCINV returns the inverse error function of in1.
func ERFINV(in1 LREAL) LREAL {
	return LREAL(math.Erfinv(float64(in1)))
}

// EXPT returns e**x, the base-e exponential of in1.
func EXP_LREAL(in1 LREAL) LREAL {
	return LREAL(math.Exp(float64(in1)))
}

// EXP returns 2**x, the base-2 exponential of in1.
func EXP2(in1 LREAL) LREAL {
	return LREAL(math.Exp2(float64(in1)))
}

// EXPM1 returns e**in1 - 1, the base-e exponential of in1 minus 1. It is more accurate than Exp(in1) - 1 when in1 is near zero.
func EXPM1(in1 LREAL) LREAL {
	return LREAL(math.Expm1(float64(in1)))
}

// FMA returns in1 * in2 + in3, computed with only one rounding. (That is, FMA returns the fused multiply-add of in1, in2, and in3.)
func FMA(in1, in2, in3 LREAL) LREAL {
	return LREAL(math.FMA(float64(in1), float64(in2), float64(in3)))
}

// FLOOR returns the greatest integer value less than or equal to x.
func FLOOR(in1 LREAL) LREAL {
	return LREAL(math.Floor(float64(in1)))
}

// FREXP breaks f into a normalized fraction and an integral power of two. It returns frac and exp satisfying f == frac × 2**exp, with the absolute value of frac in the interval [½, 1].
func FREXP(in1 LREAL) (LREAL, LINT) {
	frac, exp := math.Frexp(float64(in1))
	return LREAL(frac), LINT(exp)
}

// GAMMA returns the Gamma function of in1.
func GAMMA(in1 LREAL) LREAL {
	return LREAL(math.Gamma(float64(in1)))
}

// Hypot returns Sqrt(in1*in1 + in2*in2), taking care to avoid unnecessary overflow and underflow.
func HYPOT(in1, in2 LREAL) LREAL {
	return LREAL(math.Hypot(float64(in1), float64(in2)))
}

// Ilogb returns the binary exponent of x as an integer.
func ILOGB(in1 LREAL) LINT {
	return LINT(math.Ilogb(float64(in1)))
}

// INF returns positive infinity if sign >= 0, negative infinity if sign < 0.
func INF(sign LINT) LREAL {
	return LREAL(math.Inf(int(sign)))
}

// ISINF reports whether in1 is an infinity, according to sign. If sign > 0, IsInf reports whether in1 is positive infinity. If sign < 0, IsInf reports whether in1 is negative infinity. If sign == 0, IsInf reports whether in1 is either infinity.
func ISINF(in1 LREAL, sign LINT) BOOL {
	return BOOL(math.IsInf(float64(in1), int(sign)))
}

// ISNAN reports whether in1 is an IEEE 754 “not-a-number” value.
func ISNAN(in1 LREAL) BOOL {
	return BOOL(math.IsNaN(float64(in1)))
}

// J0 returns the order-zero Bessel function of the first kind.
func J0(in1 LREAL) LREAL {
	return LREAL(math.J0(float64(in1)))
}

// J1 returns the order-one Bessel function of the first kind.
func J1(in1 LREAL) LREAL {
	return LREAL(math.J1(float64(in1)))
}

// Jn returns the order-n Bessel function of the first kind.
func JN(in1 LINT, in2 LREAL) LREAL {
	return LREAL(math.Jn(int(in1), float64(in2)))
}

// LDEXP is the inverse of Frexp. It returns frac × 2**exp.
func LDEXP(frac LREAL, exp LINT) LREAL {
	return LREAL(math.Ldexp(float64(frac), int(exp)))
}

// LGAMMA returns the natural logarithm and sign (-1 or +1) of Gamma(x).
func LGAMMA(in1 LREAL) (LREAL, LINT) {
	x, y := math.Lgamma(float64(in1))
	return LREAL(x), LINT(y)
}

// LOG (aka LN) returns the natural logarithm of in1.
func LOG_LREAL(in1 LREAL) LREAL {
	return LREAL(math.Log(float64(in1)))
}

// LOG10 returns the decimal logarithm of in1. The special cases are the same as for Log.
func LOG10(in1 LREAL) LREAL {
	return LREAL(math.Log10(float64(in1)))
}

// LOG1P returns the natural logarithm of 1 plus its argument x. It is more accurate than Log(1 + in1) when in1 is near zero.
func LOG1P(in1 LREAL) LREAL {
	return LREAL(math.Log1p(float64(in1)))
}

// LOG2 returns the binary logarithm of in1. The special cases are the same as for Log.
func LOG2(in1 LREAL) LREAL {
	return LREAL(math.Log2(float64(in1)))
}

// LOGB returns the binary exponent of in1.
func LOGB(in1 LREAL) LREAL {
	return LREAL(math.Logb(float64(in1)))
}

// MOD returns the floating-point remainder of x/y. The magnitude of the result is less than y and its sign agrees with that of x.
func MODL(in1, in2 LREAL) LREAL {
	return LREAL(math.Mod(float64(in1), float64(in2)))
}

// MODF returns integer and fractional floating-point numbers that sum to f. Both values have the same sign as f.
func MODF(in1 LREAL) (LREAL, LREAL) {
	x, y := math.Modf(float64(in1))
	return LREAL(x), LREAL(y)
}

// NEXTAFTER returns the next representable float64 value after x towards y.
func NEXTAFTER(x, y LREAL) LREAL {
	return LREAL(math.Nextafter(float64(x), float64(y)))
}

// NEXTAFTER32 returns the next representable float32 value after x towards y.
func NEXTAFTER32(x, y REAL) REAL {
	return REAL(math.Nextafter32(float32(x), float32(y)))
}

// POW returns in1**in2, the base-in1 exponential of in2.
func POW(in1 LREAL, in2 LREAL) LREAL {
	return LREAL(math.Pow(float64(in1), float64(in2)))
}

// POW10 returns in1**in2, the base-in1 exponential of in2.
func POW10(in1 LINT) (LREAL, error) {
	n := int(in1)
	// Prevent panics from math.Pow10 for out-of-range inputs.
	if n > 308 || n < -324 {
		return 0, fmt.Errorf("POW10: input %d is out of range", n)
	}
	return LREAL(math.Pow10(n)), nil
}

// REMAINDER returns the IEEE 754 floating-point remainder of x/y.
func REMAINDER(x, y LREAL) LREAL {
	return LREAL(math.Remainder(float64(x), float64(y)))
}

// ROUND returns the nearest integer, rounding half away from zero.
func ROUND(in1 LREAL) LREAL {
	return LREAL(math.Round(float64(in1)))
}

// ROUNTOEVEN returns the nearest integer, rounding ties to even.
func ROUNTOEVEN(in1 LREAL) LREAL {
	return LREAL(math.RoundToEven(float64(in1)))
}

// SINGBIT reports whether x is negative or negative zero.
func SIGNBIT(in1 LREAL) BOOL {
	return BOOL(math.Signbit(float64(in1)))
}

// SIN returns the sine of the radian argument in1.
func SIN_LREAL(in1 LREAL) LREAL {
	return LREAL(math.Sin(float64(in1)))
}

// SINCOS  returns Sin(in1), Cos(in1).
func SINCOS(in1 LREAL) (LREAL, LREAL) {
	x, y := math.Sincos(float64(in1))
	return LREAL(x), LREAL(y)
}

// SINH returns the hyperbolic sine of x.
func SINH(in1 LREAL) LREAL {
	return LREAL(math.Sinh(float64(in1)))
}

// SQRT returns the square root of in1.
func SQRT_LREAL(in1 LREAL) LREAL {
	return LREAL(math.Sqrt(float64(in1)))
}

// TAN returns the tangent of the radian argument in1.
func TAN_LREAL(in1 LREAL) LREAL {
	return LREAL(math.Tan(float64(in1)))
}

// TANH returns the hyperbolic tangent of in1.
func TANH(in1 LREAL) LREAL {
	return LREAL(math.Tanh(float64(in1)))
}

// TRUNC returns the integer value of in1.
func TRUNC_LREAL(in1 LREAL) LREAL {
	return LREAL(math.Trunc(float64(in1)))
}

// Y0 returns the order-zero Bessel function of the second kind.
func Y0(x LREAL) LREAL {
	return LREAL(math.Y0(float64(x)))
}

// Y1 returns the order-one Bessel function of the second kind.
func Y1(x LREAL) LREAL {
	return LREAL(math.Y1(float64(x)))
}

// YN returns the order-n Bessel function of the second kind.
func YN(n LINT, x LREAL) LREAL {
	return LREAL(math.Yn(int(n), float64(x)))
}

// RAND generates a random value of a specified IEC 61131-3 data type.
// The function takes a selector `dataType` to determine the output type.
// It returns the generated value as an interface{} and an error if the
// selector is invalid. This is a non-standard extension.
func RAND(datatype interface{}) (interface{}, error) {
	switch targetType := datatype.(type) {
	case BOOL:
		return BOOL(globalRand.Intn(2) == 1), nil
	case SINT:
		return SINT(int8(globalRand.Int63())), nil
	case INT:
		return INT(int16(globalRand.Int63())), nil
	case DINT:
		return DINT(globalRand.Int31()), nil
	case LINT:
		return LINT(globalRand.Int63()), nil
	case USINT:
		return USINT(uint8(globalRand.Uint32())), nil
	case UINT:
		return UINT(uint16(globalRand.Uint32())), nil
	case UDINT:
		return UDINT(globalRand.Uint32()), nil
	case ULINT:
		return ULINT(globalRand.Uint64()), nil
	case BYTE:
		return BYTE(uint8(globalRand.Uint32())), nil
	case WORD:
		return WORD(uint16(globalRand.Uint32())), nil
	case DWORD:
		return DWORD(globalRand.Uint32()), nil
	case LWORD:
		return LWORD(globalRand.Uint64()), nil
	case REAL:
		return REAL(globalRand.Float32()), nil
	case LREAL:
		return LREAL(globalRand.Float64()), nil
	case TIME:
		// Generate a random duration, e.g., up to 24 hours
		return TIME(time.Duration(globalRand.Int63n(int64(24 * time.Hour)))), nil
	case DATE:
		// Generate a random date within a reasonable range, e.g., past 50 years
		randomSeconds := globalRand.Int63n(50 * 365 * 24 * 60 * 60)
		return DATE(time.Now().Add(-time.Second * time.Duration(randomSeconds))), nil
	case TOD:
		// Generate a random time of day (duration from midnight)
		return TOD(time.Time{}.Add(time.Duration(globalRand.Int63n(int64(24 * time.Hour))))), nil
	case DT:
		// Generate a random date and time
		randomSeconds := globalRand.Int63n(50 * 365 * 24 * 60 * 60)
		return DT(time.Now().Add(-time.Second * time.Duration(randomSeconds))), nil
	case STRING:
		// Generate a random string of a variable length (e.g., 5-15 chars)
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		length := globalRand.Intn(11) + 5
		b := make([]byte, length)
		for i := range b {
			b[i] = charset[globalRand.Intn(len(charset))]
		}
		return STRING(b), nil
	default:
		return nil, fmt.Errorf("RAND: invalid data type selector: %v", targetType)
	}
}
