// Deprecated: This file tests non-standard extensions from Math.go. The standard IEC 61131-3 numerical functions
// are tested in Numerical_test.go. This file may be removed in a future version.
package royaljelly

import (
	"math"
	"testing"
)

func TestMathWrappers(t *testing.T) {
	t.Run("CBRT", func(t *testing.T) {
		in := LREAL(27.0)
		expected := LREAL(3.0)
		result := CBRT(in)
		if !almostEqual(result, expected) {
			t.Errorf("CBRT(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("CEIL", func(t *testing.T) {
		in := LREAL(2.3)
		expected := LREAL(3.0)
		result := CEIL(in)
		if result != expected {
			t.Errorf("CEIL(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("FLOOR", func(t *testing.T) {
		in := LREAL(2.7)
		expected := LREAL(2.0)
		result := FLOOR(in)
		if result != expected {
			t.Errorf("FLOOR(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("ROUND", func(t *testing.T) {
		in1 := LREAL(2.7)
		expected1 := LREAL(3.0)
		result1 := ROUND(in1)
		if result1 != expected1 {
			t.Errorf("ROUND(%v) = %v; want %v", in1, result1, expected1)
		}

		in2 := LREAL(2.3)
		expected2 := LREAL(2.0)
		result2 := ROUND(in2)
		if result2 != expected2 {
			t.Errorf("ROUND(%v) = %v; want %v", in2, result2, expected2)
		}
	})

	t.Run("POW", func(t *testing.T) {
		base, exp := LREAL(3.0), LREAL(4.0)
		expected := LREAL(81.0)
		result := POW(base, exp)
		if !almostEqual(result, expected) {
			t.Errorf("POW(%v, %v) = %v; want %v", base, exp, result, expected)
		}
	})

	t.Run("ISNAN", func(t *testing.T) {
		nan := LREAL(math.NaN())
		num := LREAL(123.0)

		if !ISNAN(nan) {
			t.Errorf("ISNAN(NaN) = false; want true")
		}
		if ISNAN(num) {
			t.Errorf("ISNAN(123.0) = true; want false")
		}
	})

	t.Run("ISINF", func(t *testing.T) {
		posInf := LREAL(math.Inf(1))
		negInf := LREAL(math.Inf(-1))
		num := LREAL(123.0)

		if !ISINF(posInf, 1) {
			t.Errorf("ISINF(+Inf, 1) = false; want true")
		}
		if !ISINF(negInf, -1) {
			t.Errorf("ISINF(-Inf, -1) = false; want true")
		}
		if !ISINF(posInf, 0) {
			t.Errorf("ISINF(+Inf, 0) = false; want true")
		}
		if ISINF(num, 0) {
			t.Errorf("ISINF(123.0, 0) = true; want false")
		}
	})

	t.Run("ACOS_LREAL", func(t *testing.T) {
		in := LREAL(0)
		expected := LREAL(math.Pi / 2)
		result := ACOS_LREAL(in)
		if !almostEqual(result, expected) {
			t.Errorf("ACOS_LREAL(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("ASIN_LREAL", func(t *testing.T) {
		in := LREAL(0)
		expected := LREAL(0)
		result := ASIN_LREAL(in)
		if !almostEqual(result, expected) {
			t.Errorf("ASIN_LREAL(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("ATAN_LREAL", func(t *testing.T) {
		in := LREAL(0)
		expected := LREAL(0)
		result := ATAN_LREAL(in)
		if !almostEqual(result, expected) {
			t.Errorf("ATAN_LREAL(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("COS_LREAL", func(t *testing.T) {
		in := LREAL(0)
		expected := LREAL(1)
		result := COS_LREAL(in)
		if !almostEqual(result, expected) {
			t.Errorf("COS_LREAL(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("EXP_LREAL", func(t *testing.T) {
		in := LREAL(1)
		expected := LREAL(math.E)
		result := EXP_LREAL(in)
		if !almostEqual(result, expected) {
			t.Errorf("EXP_LREAL(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("INF", func(t *testing.T) {
		if INF(1) != LREAL(math.Inf(1)) {
			t.Errorf("INF(1) = %v; want %v", INF(1), math.Inf(1))
		}
		if INF(-1) != LREAL(math.Inf(-1)) {
			t.Errorf("INF(-1) = %v; want %v", INF(-1), math.Inf(-1))
		}
	})

	t.Run("LOG_LREAL", func(t *testing.T) {
		in := LREAL(math.E)
		expected := LREAL(1)
		result := LOG_LREAL(in)
		if !almostEqual(result, expected) {
			t.Errorf("LOG_LREAL(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("MODL", func(t *testing.T) {
		in1, in2 := LREAL(10.5), LREAL(3.0)
		expected := LREAL(1.5)
		result := MODL(in1, in2)
		if !almostEqual(result, expected) {
			t.Errorf("MODL(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("LOG10", func(t *testing.T) {
		in := LREAL(1000)
		expected := LREAL(3.0)
		result := LOG10(in)
		if !almostEqual(result, expected) {
			t.Errorf("LOG10(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("LOG2", func(t *testing.T) {
		in := LREAL(256)
		expected := LREAL(8.0)
		result := LOG2(in)
		if !almostEqual(result, expected) {
			t.Errorf("LOG2(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("ACOSH", func(t *testing.T) {
		in := LREAL(1)
		expected := LREAL(0)
		result := ACOSH(in)
		if !almostEqual(result, expected) {
			t.Errorf("ACOSH(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("ASINH", func(t *testing.T) {
		in := LREAL(0)
		expected := LREAL(0)
		result := ASINH(in)
		if !almostEqual(result, expected) {
			t.Errorf("ASINH(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("NEXTAFTER", func(t *testing.T) {
		x, y := LREAL(1.0), LREAL(2.0)
		expected := LREAL(math.Nextafter(float64(x), float64(y)))
		result := NEXTAFTER(x, y)
		if result != expected {
			t.Errorf("NEXTAFTER(%v, %v) = %v; want %v", x, y, result, expected)
		}
	})

	t.Run("NEXTAFTER32", func(t *testing.T) {
		x, y := REAL(1.0), REAL(2.0)
		expected := REAL(math.Nextafter32(float32(x), float32(y)))
		result := NEXTAFTER32(x, y)
		// Due to potential floating point inaccuracies, direct comparison might fail.
		// Use almostEqual for REAL if a threshold is defined for it, otherwise check string representation.
		// For Nextafter32, direct comparison is usually fine as it's about exact representation.
		if result != expected {
			t.Errorf("NEXTAFTER32(%v, %v) = %v; want %v", x, y, result, expected)

		}
	})

	t.Run("ATAN2", func(t *testing.T) {
		y, x := LREAL(1), LREAL(1)
		expected := LREAL(math.Pi / 4)
		result := ATAN2(y, x)
		if !almostEqual(result, expected) {
			t.Errorf("ATAN2(%v, %v) = %v; want %v", y, x, result, expected)
		}
	})

	t.Run("COPYSIGN", func(t *testing.T) {
		val, sign := LREAL(10.0), LREAL(-1.0)
		expected := LREAL(-10.0)
		result := COPYSIGN(val, sign)
		if result != expected {
			t.Errorf("COPYSIGN(%v, %v) = %v; want %v", val, sign, result, expected)
		}
	})

	t.Run("COSH", func(t *testing.T) {
		in := LREAL(0)
		expected := LREAL(1)
		result := COSH(in)
		if !almostEqual(result, expected) {
			t.Errorf("COSH(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("DIM", func(t *testing.T) {
		x, y := LREAL(5), LREAL(2)
		expected := LREAL(3)
		result := DIM(x, y)
		if result != expected {
			t.Errorf("DIM(%v, %v) = %v; want %v", x, y, result, expected)
		}
		x2, y2 := LREAL(2), LREAL(5)
		expected2 := LREAL(0)
		result2 := DIM(x2, y2)
		if result2 != expected2 {
			t.Errorf("DIM(%v, %v) = %v; want %v", x2, y2, result2, expected2)
		}
	})

	t.Run("ERF/ERFC/ERFINV/ERFCINV", func(t *testing.T) {
		if ERF(0) != 0 {
			t.Error("ERF(0) should be 0")
		}
		if ERFC(0) != 1 {
			t.Error("ERFC(0) should be 1")
		}
		if ERFINV(0) != 0 {
			t.Error("ERFINV(0) should be 0")
		}
		if ERFCINV(1) != 0 {
			t.Error("ERFCINV(1) should be 0")
		}
	})

	t.Run("EXP2", func(t *testing.T) {
		in := LREAL(3)
		expected := LREAL(8)
		result := EXP2(in)
		if result != expected {
			t.Errorf("EXP2(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("EXPM1", func(t *testing.T) {
		in := LREAL(1)
		expected := LREAL(math.E - 1)
		result := EXPM1(in)
		if !almostEqual(result, expected) {
			t.Errorf("EXPM1(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("FMA", func(t *testing.T) {
		in1, in2, in3 := LREAL(2), LREAL(3), LREAL(4)
		expected := LREAL(10)
		result := FMA(in1, in2, in3)
		if result != expected {
			t.Errorf("FMA(%v, %v, %v) = %v; want %v", in1, in2, in3, result, expected)
		}
	})

	t.Run("FREXP/LDEXP", func(t *testing.T) {
		in := LREAL(12.5)
		frac, exp := FREXP(in)
		if frac != 0.78125 || exp != 4 {
			t.Errorf("FREXP(%v) = %v, %v; want 0.78125, 4", in, frac, exp)
		}
		result := LDEXP(frac, exp)
		if result != in {
			t.Errorf("LDEXP(%v, %v) = %v; want %v", frac, exp, result, in)
		}
	})

	t.Run("GAMMA", func(t *testing.T) {
		in := LREAL(4)
		expected := LREAL(6) // (4-1)!
		result := GAMMA(in)
		if !almostEqual(result, expected) {
			t.Errorf("GAMMA(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("HYPOT", func(t *testing.T) {
		p, q := LREAL(3), LREAL(4)
		expected := LREAL(5)
		result := HYPOT(p, q)
		if result != expected {
			t.Errorf("HYPOT(%v, %v) = %v; want %v", p, q, result, expected)
		}
	})

	t.Run("ILOGB/LOGB", func(t *testing.T) {
		in := LREAL(128)
		if ILOGB(in) != 7 {
			t.Errorf("ILOGB(%v) = %v; want 7", in, ILOGB(in))
		}
		if LOGB(in) != 7.0 {
			t.Errorf("LOGB(%v) = %v; want 7.0", in, LOGB(in))
		}
	})

	t.Run("J0/J1/JN", func(t *testing.T) {
		if J0(0) != 1.0 {
			t.Error("J0(0) should be 1.0")
		}
		if J1(0) != 0.0 {
			t.Error("J1(0) should be 0.0")
		}
		if JN(2, 0) != 0.0 {
			t.Error("JN(2, 0) should be 0.0")
		}
	})

	t.Run("LGAMMA", func(t *testing.T) {
		lgamma, sign := LGAMMA(LREAL(4))
		expectedLgamma := LREAL(math.Log(6))
		if !almostEqual(lgamma, expectedLgamma) || sign != 1 {
			t.Errorf("LGAMMA(4) = %v, %v; want %v, 1", lgamma, sign, expectedLgamma)
		}
	})

	t.Run("LOG1P", func(t *testing.T) {
		result := LOG1P(0)
		if result != 0 {
			t.Errorf("LOG1P(0) = %v; want 0", result)
		}
	})

	t.Run("MODF", func(t *testing.T) {
		i, f := MODF(LREAL(3.14))
		if i != 3.0 || !almostEqual(f, 0.14) {
			t.Errorf("MODF(3.14) = %v, %v; want 3.0, 0.14", i, f)
		}
	})

	t.Run("POW10", func(t *testing.T) {
		res, err := POW10(3)
		if err != nil || res != 1000.0 {
			t.Errorf("POW10(3) failed. Got %v, err: %v", res, err)
		}
		_, err = POW10(400)
		if err == nil {
			t.Error("POW10(400) should have returned an error for overflow")
		}
		_, err = POW10(-400)
		if err == nil {
			t.Error("POW10(-400) should have returned an error for underflow")
		}
	})

	t.Run("REMAINDER", func(t *testing.T) {
		result := REMAINDER(LREAL(10.5), LREAL(3.0))
		if !almostEqual(result, -1.5) {
			t.Errorf("REMAINDER(10.5, 3.0) = %v; want -1.5", result)
		}
	})

	t.Run("ROUNTOEVEN", func(t *testing.T) {
		if ROUNTOEVEN(2.5) != 2.0 || ROUNTOEVEN(3.5) != 4.0 {
			t.Errorf("ROUNTOEVEN failed. 2.5->%v, 3.5->%v", ROUNTOEVEN(2.5), ROUNTOEVEN(3.5))
		}
	})

	t.Run("SIGNBIT", func(t *testing.T) {
		if !SIGNBIT(-1.0) || SIGNBIT(1.0) {
			t.Error("SIGNBIT failed")
		}
	})

	t.Run("SINCOS", func(t *testing.T) {
		s, c := SINCOS(0)
		if s != 0 || c != 1 {
			t.Errorf("SINCOS(0) = %v, %v; want 0, 1", s, c)
		}
	})

	t.Run("SIN_LREAL", func(t *testing.T) {
		in := LREAL(0)
		expected := LREAL(0)
		result := SIN_LREAL(in)
		if !almostEqual(result, expected) {
			t.Errorf("SIN_LREAL(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("SINH", func(t *testing.T) {
		in := LREAL(0)
		expected := LREAL(0)
		result := SINH(in)
		if !almostEqual(result, expected) {
			t.Errorf("SINH(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("SQRT_LREAL", func(t *testing.T) {
		in := LREAL(25.0)
		expected := LREAL(5.0)
		result := SQRT_LREAL(in)
		if !almostEqual(result, expected) {
			t.Errorf("SQRT_LREAL(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("TAN_LREAL", func(t *testing.T) {
		in := LREAL(0)
		expected := LREAL(0)
		result := TAN_LREAL(in)
		if !almostEqual(result, expected) {
			t.Errorf("TAN_LREAL(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("TANH", func(t *testing.T) {
		in := LREAL(0)
		expected := LREAL(0)
		result := TANH(in)
		if !almostEqual(result, expected) {
			t.Errorf("TANH(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("TRUNC_LREAL", func(t *testing.T) {
		in := LREAL(123.75)
		expected := LREAL(123.0)
		result := TRUNC_LREAL(in)
		if result != expected {
			t.Errorf("TRUNC_LREAL(%v) = %v; want %v", in, result, expected)
		}

		in2 := LREAL(-45.9)
		expected2 := LREAL(-45.0)
		result2 := TRUNC_LREAL(in2)
		if result2 != expected2 {
			t.Errorf("TRUNC_LREAL(%v) = %v; want %v", in2, result2, expected2)
		}

	})

	t.Run("Y0", func(t *testing.T) {
		in := LREAL(1.0)
		expected := LREAL(math.Y0(float64(in)))
		result := Y0(in)
		if !almostEqual(result, expected) {
			t.Errorf("Y0(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("Y1", func(t *testing.T) {
		in := LREAL(1.0)
		expected := LREAL(math.Y1(float64(in)))
		result := Y1(in)
		if !almostEqual(result, expected) {
			t.Errorf("Y1(%v) = %v; want %v", in, result, expected)
		}
	})

	t.Run("YN", func(t *testing.T) {
		n := ANYINT(0)
		x := LREAL(1.0)
		expected := LREAL(math.Yn(int(n), float64(x)))
		result := YN(n, x)
		if !almostEqual(result, expected) {
			t.Errorf("YN(%v, %v) = %v; want %v", n, x, result, expected)
		}
	})
}

func TestConstants(t *testing.T) {
	if !math.IsNaN(NAN) {
		t.Errorf("Constant NAN does not match math.NaN()")
	}
	if POSINF != math.Inf(1) {
		t.Errorf("Constant POSINF does not match math.Inf(1)")
	}
	if NEGINF != math.Inf(-1) {
		t.Errorf("Constant NEGINF does not match math.Inf(-1)")
	}
}
