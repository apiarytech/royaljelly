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
	testCases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{"Negative LINT", LINT(-100), LINT(100)},
		{"Positive LINT", LINT(50), LINT(50)},
		{"Zero DINT", DINT(0), DINT(0)},
		{"Negative REAL", REAL(-123.45), REAL(123.45)},
		{"Positive LREAL", LREAL(567.89), LREAL(567.89)},
		{"Unsigned UINT", UINT(200), UINT(200)},
		{"Unsigned ULINT", ULINT(999), ULINT(999)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ABS(tc.input)
			if result != tc.expected {
				t.Errorf("ABS(%v) = %v; want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestSQRT(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected LREAL
	}{
		{"Perfect square REAL", REAL(25.0), 5.0},
		{"Non-perfect square LREAL", LREAL(2.0), 1.414213562},
		{"Zero INT", INT(0), 0.0},
		{"Negative REAL", REAL(-4.0), LREAL(math.NaN())},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SQRT(tc.input)
			resLREAL, _ := anyToLREAL(result)

			if math.IsNaN(float64(tc.expected)) {
				if !math.IsNaN(float64(resLREAL)) {
					t.Errorf("SQRT(%v) = %v; want NaN", tc.input, resLREAL)
				}
			} else if !almostEqual(resLREAL, tc.expected) {
				t.Errorf("SQRT(%v) = %v; want %v", tc.input, resLREAL, tc.expected)
			}
		})
	}
}

func TestLogarithms(t *testing.T) {
	t.Run("LN", func(t *testing.T) {
		result := LN(LREAL(math.E))
		if !almostEqual(result.(LREAL), 1.0) {
			t.Errorf("LN(e) = %v; want 1.0", result)
		}
	})

	t.Run("LOG", func(t *testing.T) {
		result := LOG(LREAL(100.0))
		if !almostEqual(result.(LREAL), 2.0) {
			t.Errorf("LOG(100) = %v; want 2.0", result)
		}
	})

	t.Run("LN of zero", func(t *testing.T) {
		result := LN(REAL(0))
		if !math.IsInf(float64(result.(REAL)), -1) {
			t.Errorf("LN(0) = %v; want -Inf", result)
		}
	})
}

func TestEXP(t *testing.T) {
	t.Run("EXP of 1", func(t *testing.T) {
		result := EXP(LREAL(1.0))
		if !almostEqual(result.(LREAL), LREAL(math.E)) {
			t.Errorf("EXP(1.0) = %v; want %v", result, math.E)
		}
	})

	t.Run("EXP of 0", func(t *testing.T) {
		result := EXP(LREAL(0.0))
		if !almostEqual(result.(LREAL), 1.0) {
			t.Errorf("EXP(0.0) = %v; want 1.0", result)
		}
	})
}

func TestEXPT(t *testing.T) {
	testCases := []struct {
		name     string
		base     interface{}
		exp      interface{}
		expected LREAL
	}{
		{"Integer base and exp", LINT(2), INT(8), 256.0},
		{"Real base, integer exp", REAL(2.5), DINT(2), 6.25},
		{"Integer base, real exp", INT(4), REAL(0.5), 2.0},
		{"Negative exponent", LREAL(10.0), SINT(-2), 0.01},
		{"Zero exponent", REAL(123.45), INT(0), 1.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := EXPT(tc.base, tc.exp)
			resLREAL, _ := anyToLREAL(result)
			if !almostEqual(resLREAL, tc.expected) {
				t.Errorf("EXPT(%v, %v) = %v; want %v", tc.base, tc.exp, resLREAL, tc.expected)
			}
		})
	}
}

func TestTrigonometric(t *testing.T) {
	pi := LREAL(math.Pi)
	testCases := []struct {
		name     string
		fn       func(interface{}) interface{}
		input    LREAL
		expected LREAL
	}{
		{"SIN(0)", SIN, 0, 0},
		{"SIN(pi/2)", SIN, pi / 2, 1},
		{"COS(0)", COS, 0, 1},
		{"COS(pi)", COS, pi, -1},
		{"TAN(0)", TAN, 0, 0},
		{"TAN(pi/4)", TAN, pi / 4, 1},
		{"ASIN(1)", ASIN, 1, pi / 2},
		{"ACOS(1)", ACOS, 1, 0},
		{"ATAN(1)", ATAN, 1, pi / 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.fn(tc.input)
			if !almostEqual(result.(LREAL), tc.expected) {
				t.Errorf("%s = %v; want %v", tc.name, result, tc.expected)
			}
		})
	}
}

func TestTRUNC(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		expected    LINT
		expectPanic bool
	}{
		{"Positive REAL", REAL(123.75), 123, false},
		{"Negative LREAL", LREAL(-45.9), -45, false},
		{"Zero REAL", REAL(0.0), 0, false},
		{"Integer input", LINT(100), 0, true},
		{"String input", STRING("12.3"), 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.expectPanic {
					if r == nil {
						t.Errorf("TRUNC(%v) did not panic; expected panic", tc.input)
					}
				} else if r != nil {
					t.Errorf("TRUNC(%v) panicked; got %v", tc.input, r)
				}
			}()

			result := TRUNC(tc.input)
			if !tc.expectPanic {
				if result != tc.expected {
					t.Errorf("TRUNC(%v) = %v; want %v", tc.input, result, tc.expected)
				}
			}
		})
	}
}
