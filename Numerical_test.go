package royaljelly

import (
	"math"
	"testing"
)

const float64EqualityThreshold = 1e-9

func almostEqual(a, b LREAL) bool {
	return math.Abs(float64(a-b)) <= float64EqualityThreshold
}

func TestSUMLINT(t *testing.T) {
	// This test remains valid for the deprecated function.
	// It can be removed when the function is fully removed.
	t.Run("Basic Sum", func(t *testing.T) {
		m := map[STRING]LINT{
			"a": 10,
			"b": 20,
			"c": -5,
		}
		expected := LINT(25)
		result := SUMLINT(m)
		if result != expected {
			t.Errorf("SUMLINT() = %d; want %d", result, expected)
		}
	})

	t.Run("Empty Map", func(t *testing.T) {
		m := make(map[STRING]LINT)
		expected := LINT(0)
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
		m := map[STRING]REAL{
			"a": 10.5,
			"b": 20.25,
			"c": -5.0,
		}
		expected := REAL(25.75)
		result := SUMREAL(m)
		if result != expected {
			t.Errorf("SUMREAL() = %f; want %f", result, expected)
		}
	})

	t.Run("Empty Map", func(t *testing.T) {
		m := make(map[STRING]REAL)
		expected := REAL(0)
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
		m := map[STRING]LREAL{
			"a": 100.125,
			"b": 200.250,
			"c": -50.0,
		}
		expected := LREAL(250.375)
		result := SUMLREAL(m)
		if result != expected {
			t.Errorf("SUMLREAL() = %f; want %f", result, expected)
		}
	})

	t.Run("Empty Map", func(t *testing.T) {
		m := make(map[STRING]LREAL)
		expected := LREAL(0)
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
		m := map[int]LINT{
			1: 100,
			2: 200,
			3: 300,
		}
		expected := LINT(600)
		result := SUMLINTorLREAL(m)
		if result != expected {
			t.Errorf("SUMLINTorLREAL() with LINT = %d; want %d", result, expected)
		}
	})

	t.Run("Sum LREAL with string key", func(t *testing.T) {
		m := map[string]LREAL{
			"x": 1.1,
			"y": 2.2,
			"z": 3.3,
		}
		// Use a tolerance for float comparison
		expected := LREAL(6.6)
		result := SUMLINTorLREAL(m)
		if result < expected-1e-9 || result > expected+1e-9 {
			t.Errorf("SUMLINTorLREAL() with LREAL = %f; want %f", result, expected)
		}
	})
}

func TestSUM(t *testing.T) {
	// This test remains valid for the non-standard generic SUM function.
	t.Run("Sum INT", func(t *testing.T) {
		m := map[string]INT{
			"one": 1,
			"two": 2,
		}
		expected := INT(3)
		result := SUM(m)
		if result != expected {
			t.Errorf("SUM() with INT = %d; want %d", result, expected)
		}
	})

	t.Run("Sum UINT", func(t *testing.T) {
		m := map[int]UINT{
			1: 1000,
			2: 2000,
		}
		expected := UINT(3000)
		result := SUM(m)
		if result != expected {
			t.Errorf("SUM() with UINT = %d; want %d", result, expected)
		}
	})
}

func TestABS(t *testing.T) {
	t.Run("Negative LINT", func(t *testing.T) {
		if res := ABS(LINT(-100)); res != 100 {
			t.Errorf("ABS(-100) = %v; want 100", res)
		}
	})
	t.Run("Positive LINT", func(t *testing.T) {
		if res := ABS(LINT(50)); res != 50 {
			t.Errorf("ABS(50) = %v; want 50", res)
		}
	})
	t.Run("Negative REAL", func(t *testing.T) {
		if res := ABS(REAL(-123.45)); res != 123.45 {
			t.Errorf("ABS(-123.45) = %v; want 123.45", res)
		}
	})
}

func TestSQRT(t *testing.T) {
	t.Run("Perfect square REAL", func(t *testing.T) {
		if res := SQRT(REAL(25.0)); res != 5.0 {
			t.Errorf("SQRT(25.0) = %v; want 5.0", res)
		}
	})
	t.Run("Non-perfect square LREAL", func(t *testing.T) {
		if res := SQRT(LREAL(2.0)); !almostEqual(res, 1.414213562) {
			t.Errorf("SQRT(2.0) = %v; want ~1.414", res)
		}
	})
	t.Run("Negative REAL", func(t *testing.T) {
		if res := SQRT(REAL(-4.0)); !math.IsNaN(float64(res)) {
			t.Errorf("SQRT(-4.0) = %v; want NaN", res)
		}
	})
}

func TestLogarithms(t *testing.T) {
	t.Run("LN", func(t *testing.T) {
		result := LN(LREAL(math.E))
		if !almostEqual(result, 1.0) {
			t.Errorf("LN(e) = %v; want 1.0", result)
		}
	})

	t.Run("LOG", func(t *testing.T) {
		result := LOG(LREAL(100.0))
		if !almostEqual(result, 2.0) {
			t.Errorf("LOG(100) = %v; want 2.0", result)
		}
	})

	t.Run("LN of zero", func(t *testing.T) {
		result := LN(REAL(0))
		if !math.IsInf(float64(result), -1) {
			t.Errorf("LN(0) = %v; want -Inf", result)
		}
	})
}

func TestEXP(t *testing.T) {
	t.Run("EXP of 1", func(t *testing.T) {
		result := EXP(LREAL(1.0))
		if !almostEqual(result, LREAL(math.E)) {
			t.Errorf("EXP(1.0) = %v; want %v", result, math.E)
		}
	})

	t.Run("EXP of 0", func(t *testing.T) {
		result := EXP(LREAL(0.0))
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
		expected    LREAL
		expectError bool
	}{
		{"Integer base and exp", LREAL(2), INT(8), 256.0, false},
		{"Real base, integer exp", REAL(2.5), DINT(2), 6.25, false},
		{"Integer base, real exp", REAL(4), REAL(0.5), 2.0, false},
		{"Negative exponent", LREAL(10.0), SINT(-2), 0.01, false},
		{"Zero exponent", REAL(123.45), INT(0), 1.0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result LREAL

			// Use a type switch to call the generic function with the correct concrete types.
			// The exponent must be converted to a real type to match the generic constraint.
			switch base := tc.base.(type) {
			case REAL:
				switch exponent := tc.exponent.(type) {
				case INT:
					result = LREAL(EXPT(base, REAL(exponent)))
				case DINT:
					result = LREAL(EXPT(base, REAL(exponent)))
				case REAL:
					result = LREAL(EXPT(base, exponent))
				default:
					t.Fatalf("unhandled exponent type for REAL base in test: %T", tc.exponent)
				}
			case LREAL:
				switch exponent := tc.exponent.(type) {
				case SINT:
					result = EXPT(base, LREAL(exponent))
				case INT:
					result = EXPT(base, LREAL(exponent))
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
	pi := LREAL(math.Pi)
	testCases := []struct {
		name     string
		fn       func(LREAL) LREAL
		input    LREAL
		expected LREAL
	}{
		{"SIN(0)", SIN[LREAL], 0, 0},
		{"SIN(pi/2)", SIN[LREAL], pi / 2, 1},
		{"COS(0)", COS[LREAL], 0, 1},
		{"COS(pi)", COS[LREAL], pi, -1},
		{"TAN(0)", TAN[LREAL], 0, 0},
		{"TAN(pi/4)", TAN[LREAL], pi / 4, 1},
		{"ASIN(1)", ASIN[LREAL], 1, pi / 2},
		{"ACOS(1)", ACOS[LREAL], 1, 0},
		{"ATAN(1)", ATAN[LREAL], 1, pi / 4},
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
		expected    DINT
		expectPanic bool
		expectError bool
	}{
		{"Positive REAL", REAL(123.75), DINT(123), false, false},
		{"Negative LREAL", LREAL(-45.9), DINT(-45), false, false},
		{"Zero REAL", REAL(0.0), DINT(0), false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result DINT
			switch v := tc.input.(type) {
			case REAL:
				result = TRUNC(v)
			case LREAL:
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
