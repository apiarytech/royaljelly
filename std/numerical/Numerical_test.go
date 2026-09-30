package numerical

import (
	"math"
	"testing"

	"github.com/apiarytech/royaljelly/iec"
)

const float64EqualityThreshold = 1e-9

func almostEqual(a, b iec.LREAL) bool {
	return math.Abs(float64(a-b)) <= float64EqualityThreshold
}

func TestSUMLINT(t *testing.T) {
	// This test remains valid for the deprecated function.
	// It can be removed when the function is fully removed.
	t.Run("Basic Sum", func(t *testing.T) {
		m := map[iec.STRING]iec.LINT{
			"a": 10,
			"b": 20,
			"c": -5,
		}
		expected := iec.LINT(25)
		result := SUMLINT(m)
		if result != expected {
			t.Errorf("SUMLINT() = %d; want %d", result, expected)
		}
	})

	t.Run("Empty Map", func(t *testing.T) {
		m := make(map[iec.STRING]iec.LINT)
		expected := iec.LINT(0)
		result := SUMLINT(m)
		if result != expected {
			t.Errorf("SUMLINT() on empty map = %d; want %d", result, expected)
		}
	})
}

func TestSUMREAL(t *testing.T) {
	// This test remains valid for the deprecated function.
	// It can be removed when the function is fully removed.
	t.Run("Basic Sum", func(t *testing.T) {
		m := map[iec.STRING]iec.REAL{
			"a": 10.5,
			"b": 20.25,
			"c": -5.0,
		}
		expected := iec.REAL(25.75)
		result := SUMREAL(m)
		if result != expected {
			t.Errorf("SUMREAL() = %f; want %f", result, expected)
		}
	})

	t.Run("Empty Map", func(t *testing.T) {
		m := make(map[iec.STRING]iec.REAL)
		expected := iec.REAL(0)
		result := SUMREAL(m)
		if result != expected {
			t.Errorf("SUMREAL() on empty map = %f; want %f", result, expected)
		}
	})
}

func TestSUMLREAL(t *testing.T) {
	// This test remains valid for the deprecated function.
	// It can be removed when the function is fully removed.
	t.Run("Basic Sum", func(t *testing.T) {
		m := map[iec.STRING]iec.LREAL{
			"a": 100.125,
			"b": 200.250,
			"c": -50.0,
		}
		expected := iec.LREAL(250.375)
		result := SUMLREAL(m)
		if result != expected {
			t.Errorf("SUMLREAL() = %f; want %f", result, expected)
		}
	})

	t.Run("Empty Map", func(t *testing.T) {
		m := make(map[iec.STRING]iec.LREAL)
		expected := iec.LREAL(0)
		result := SUMLREAL(m)
		if result != expected {
			t.Errorf("SUMLREAL() on empty map = %f; want %f", result, expected)
		}
	})
}

func TestSUMLINTorLREAL(t *testing.T) {
	// This test remains valid for the deprecated function.
	// It can be removed when the function is fully removed.
	t.Run("Sum LINT with int key", func(t *testing.T) {
		m := map[int]iec.LINT{
			1: 100,
			2: 200,
			3: 300,
		}
		expected := iec.LINT(600)
		result := SUMLINTorLREAL(m)
		if result != expected {
			t.Errorf("SUMLINTorLREAL() with LINT = %d; want %d", result, expected)
		}
	})

	t.Run("Sum LREAL with string key", func(t *testing.T) {
		m := map[string]iec.LREAL{
			"x": 1.1,
			"y": 2.2,
			"z": 3.3,
		}
		// Use a tolerance for float comparison
		expected := iec.LREAL(6.6)
		result := SUMLINTorLREAL(m)
		if result < expected-1e-9 || result > expected+1e-9 {
			t.Errorf("SUMLINTorLREAL() with LREAL = %f; want %f", result, expected)
		}
	})
}

func TestSUM(t *testing.T) {
	// This test remains valid for the non-standard generic SUM function.
	t.Run("Sum INT", func(t *testing.T) {
		m := map[string]iec.INT{
			"one": 1,
			"two": 2,
		}
		expected := iec.INT(3)
		result := SUM(m)
		if result != expected {
			t.Errorf("SUM() with INT = %d; want %d", result, expected)
		}
	})

	t.Run("Sum UINT", func(t *testing.T) {
		m := map[int]iec.UINT{
			1: 1000,
			2: 2000,
		}
		expected := iec.UINT(3000)
		result := SUM(m)
		if result != expected {
			t.Errorf("SUM() with UINT = %d; want %d", result, expected)
		}
	})
}

func TestABS(t *testing.T) {
	t.Run("Negative LINT", func(t *testing.T) {
		if res := ABS(iec.LINT(-100)); res != 100 {
			t.Errorf("ABS(-100) = %v; want 100", res)
		}
	})
	t.Run("Positive LINT", func(t *testing.T) {
		if res := ABS(iec.LINT(50)); res != 50 {
			t.Errorf("ABS(50) = %v; want 50", res)
		}
	})
	t.Run("Negative REAL", func(t *testing.T) {
		if res := ABS(iec.REAL(-123.45)); res != 123.45 {
			t.Errorf("ABS(-123.45) = %v; want 123.45", res)
		}
	})
}

func TestSQRT(t *testing.T) {
	t.Run("Perfect square REAL", func(t *testing.T) {
		if res := SQRT(iec.REAL(25.0)); res != 5.0 {
			t.Errorf("SQRT(25.0) = %v; want 5.0", res)
		}
	})
	t.Run("Non-perfect square LREAL", func(t *testing.T) {
		if res := SQRT(iec.LREAL(2.0)); !almostEqual(res, 1.414213562) {
			t.Errorf("SQRT(2.0) = %v; want ~1.414", res)
		}
	})
	t.Run("Negative REAL", func(t *testing.T) {
		if res := SQRT(iec.REAL(-4.0)); !math.IsNaN(float64(res)) {
			t.Errorf("SQRT(-4.0) = %v; want NaN", res)
		}
	})
}

func TestLogarithms(t *testing.T) {
	t.Run("LN", func(t *testing.T) {
		result := LN(iec.LREAL(math.E))
		if !almostEqual(result, 1.0) {
			t.Errorf("LN(e) = %v; want 1.0", result)
		}
	})

	t.Run("LOG", func(t *testing.T) {
		result := LOG(iec.LREAL(100.0))
		if !almostEqual(result, 2.0) {
			t.Errorf("LOG(100) = %v; want 2.0", result)
		}
	})

	t.Run("LN of zero", func(t *testing.T) {
		result := LN(iec.REAL(0))
		if !math.IsInf(float64(result), -1) {
			t.Errorf("LN(0) = %v; want -Inf", result)
		}
	})
}

func TestEXP(t *testing.T) {
	t.Run("EXP of 1", func(t *testing.T) {
		result := EXP(iec.LREAL(1.0))
		if !almostEqual(result, iec.LREAL(math.E)) {
			t.Errorf("EXP(1.0) = %v; want %v", result, math.E)
		}
	})

	t.Run("EXP of 0", func(t *testing.T) {
		result := EXP(iec.LREAL(0.0))
		if !almostEqual(result, 1.0) {
			t.Errorf("EXP(0.0) = %v; want 1.0", result)
		}
	})
}

func TestEXPT(t *testing.T) {
	testCases := []struct {
		name        string
		base        interface{}
		exponent    interface{}
		expected    iec.LREAL
		expectError bool
	}{
		{"Integer base and exp", iec.LREAL(2), iec.INT(8), 256.0, false},
		{"Real base, integer exp", iec.REAL(2.5), iec.DINT(2), 6.25, false},
		{"Integer base, real exp", iec.REAL(4), iec.REAL(0.5), 2.0, false},
		{"Negative exponent", iec.LREAL(10.0), iec.SINT(-2), 0.01, false},
		{"Zero exponent", iec.REAL(123.45), iec.INT(0), 1.0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result iec.LREAL

			// Use a type switch to call the generic function with the correct concrete types.
			// The exponent must be converted to a real type to match the generic constraint.
			switch base := tc.base.(type) {
			case iec.REAL:
				switch exponent := tc.exponent.(type) {
				case iec.INT:
					result = iec.LREAL(EXPT(base, iec.REAL(exponent)))
				case iec.DINT:
					result = iec.LREAL(EXPT(base, iec.REAL(exponent)))
				case iec.REAL:
					result = iec.LREAL(EXPT(base, exponent))
				default:
					t.Fatalf("unhandled exponent type for REAL base in test: %T", tc.exponent)
				}
			case iec.LREAL:
				switch exponent := tc.exponent.(type) {
				case iec.SINT:
					result = EXPT(base, iec.LREAL(exponent))
				case iec.INT:
					result = EXPT(base, iec.LREAL(exponent))
				default:
					t.Fatalf("unhandled exponent type for LREAL base in test: %T", tc.exponent)
				}
			default:
				if !tc.expectError {
					t.Fatalf("unhandled base type in test: %T", tc.base)
				}
				return // Exit test for expected error cases.
			}

			if !almostEqual(result, tc.expected) {
				t.Errorf("EXPT(%v, %v) = %v; want %v", tc.base, tc.exponent, result, tc.expected)
			}
		})
	}
}

func TestTrigonometric(t *testing.T) {
	pi := iec.LREAL(math.Pi)
	testCases := []struct {
		name     string
		fn       func(iec.LREAL) iec.LREAL
		input    iec.LREAL
		expected iec.LREAL
	}{
		{"SIN(0)", SIN[iec.LREAL], 0, 0},
		{"SIN(pi/2)", SIN[iec.LREAL], pi / 2, 1},
		{"COS(0)", COS[iec.LREAL], 0, 1},
		{"COS(pi)", COS[iec.LREAL], pi, -1},
		{"TAN(0)", TAN[iec.LREAL], 0, 0},
		{"TAN(pi/4)", TAN[iec.LREAL], pi / 4, 1},
		{"ASIN(1)", ASIN[iec.LREAL], 1, pi / 2},
		{"ACOS(1)", ACOS[iec.LREAL], 1, 0},
		{"ATAN(1)", ATAN[iec.LREAL], 1, pi / 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.fn(tc.input)
			if !almostEqual(result, tc.expected) {
				t.Errorf("%s = %v; want %v", tc.name, result, tc.expected)
			}
		})
	}
}

func TestTRUNC(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		expected    iec.DINT
		expectPanic bool
		expectError bool
	}{
		{"Positive REAL", iec.REAL(123.75), iec.DINT(123), false, false},
		{"Negative LREAL", iec.LREAL(-45.9), iec.DINT(-45), false, false},
		{"Zero REAL", iec.REAL(0.0), iec.DINT(0), false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result iec.DINT
			switch v := tc.input.(type) {
			case iec.REAL:
				result = TRUNC(v)
			case iec.LREAL:
				result = TRUNC(v)
			default:
				t.Fatalf("unhandled type for TRUNC test: %T", v)
			}
			if result != tc.expected {
				t.Errorf("TRUNC(%v) = %v; want %v", tc.input, result, tc.expected)
			}
		})
	}
}
