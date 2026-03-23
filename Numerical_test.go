package royaljelly

import (
	"math"
	"reflect"
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
		name        string
		input       interface{}
		expected    interface{}
		expectError bool
	}{
		{"Negative LINT", LINT(-100), LINT(100), false},
		{"Positive LINT", LINT(50), LINT(50), false},
		{"Zero DINT", DINT(0), DINT(0), false},
		{"Negative REAL", REAL(-123.45), REAL(123.45), false},
		{"Positive LREAL", LREAL(567.89), LREAL(567.89), false},
		{"Unsigned UINT", UINT(200), UINT(200), false},
		{"Unsigned ULINT", ULINT(999), ULINT(999), false},
		{"String input error", STRING("123"), nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ABS(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("ABS(%v) did not return an error; expected error", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("ABS(%v) returned an unexpected error: %v", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("ABS(%v) = %v; want %v", tc.input, result, tc.expected)
				}
			}
		})
	}
}

func TestSQRT(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		expected    LREAL
		expectError bool
	}{
		{"Perfect square REAL", REAL(25.0), 5.0, false},
		{"Non-perfect square LREAL", LREAL(2.0), 1.414213562, false},
		{"Zero INT", INT(0), 0.0, false},
		{"Negative REAL", REAL(-4.0), LREAL(math.NaN()), false},
		{"String input error", STRING("25.0"), 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SQRT(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("SQRT(%v) did not return an error; expected error", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("SQRT(%v) returned an unexpected error: %v", tc.input, err)
				}
				resLREAL, _ := anyToLREAL(result)

				if math.IsNaN(float64(tc.expected)) {
					if !math.IsNaN(float64(resLREAL)) {
						t.Errorf("SQRT(%v) = %v; want NaN", tc.input, resLREAL)
					}
				} else if !almostEqual(resLREAL, tc.expected) {
					t.Errorf("SQRT(%v) = %v; want %v", tc.input, resLREAL, tc.expected)
				}
			}
		})
	}
}

func TestLogarithms(t *testing.T) {
	t.Run("LN", func(t *testing.T) {
		result, err := LN(LREAL(math.E))
		if err != nil {
			t.Fatalf("LN returned an unexpected error: %v", err)
		}
		if !almostEqual(result.(LREAL), 1.0) {
			t.Errorf("LN(e) = %v; want 1.0", result)
		}
	})

	t.Run("LOG", func(t *testing.T) {
		result, err := LOG(LREAL(100.0))
		if err != nil {
			t.Fatalf("LOG returned an unexpected error: %v", err)
		}
		if !almostEqual(result.(LREAL), 2.0) {
			t.Errorf("LOG(100) = %v; want 2.0", result)
		}
	})

	t.Run("LN of zero", func(t *testing.T) {
		result, err := LN(REAL(0))
		if err != nil {
			t.Fatalf("LN(0) returned an unexpected error: %v", err)
		}
		if !math.IsInf(float64(result.(REAL)), -1) {
			t.Errorf("LN(0) = %v; want -Inf", result)
		}
	})

	t.Run("LN with string error", func(t *testing.T) {
		_, err := LN(STRING("1.0"))
		if err == nil {
			t.Error("LN(STRING) did not return an error")
		}
	})

	t.Run("LOG with string error", func(t *testing.T) {
		_, err := LOG(STRING("100.0"))
		if err == nil {
			t.Error("LOG(STRING) did not return an error")
		}
	})
}

func TestEXP(t *testing.T) {
	t.Run("EXP of 1", func(t *testing.T) {
		result, err := EXP(LREAL(1.0))
		if err != nil {
			t.Fatalf("EXP returned an unexpected error: %v", err)
		}
		if !almostEqual(result.(LREAL), LREAL(math.E)) {
			t.Errorf("EXP(1.0) = %v; want %v", result, math.E)
		}
	})

	t.Run("EXP of 0", func(t *testing.T) {
		result, err := EXP(LREAL(0.0))
		if err != nil {
			t.Fatalf("EXP returned an unexpected error: %v", err)
		}
		if !almostEqual(result.(LREAL), 1.0) {
			t.Errorf("EXP(0.0) = %v; want 1.0", result)
		}
	})

	t.Run("EXP with string error", func(t *testing.T) {
		_, err := EXP(STRING("1.0"))
		if err == nil {
			t.Error("EXP(STRING) did not return an error")
		}
	})
}

func TestEXPT(t *testing.T) {
	testCases := []struct {
		name        string
		base        interface{}
		exp         interface{}
		expected    LREAL
		expectError bool
	}{
		{"Integer base and exp", LINT(2), INT(8), 256.0, false},
		{"Real base, integer exp", REAL(2.5), DINT(2), 6.25, false},
		{"Integer base, real exp", INT(4), REAL(0.5), 2.0, false},
		{"Negative exponent", LREAL(10.0), SINT(-2), 0.01, false},
		{"Zero exponent", REAL(123.45), INT(0), 1.0, false},
		{"String base error", STRING("2.0"), INT(2), 0, true},
		{"String exponent error", LREAL(2.0), STRING("2"), 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := EXPT(tc.base, tc.exp)

			if tc.expectError {
				if err == nil {
					t.Errorf("EXPT(%v, %v) did not return an error; expected error", tc.base, tc.exp)
				}
			} else {
				if err != nil {
					t.Errorf("EXPT(%v, %v) returned an unexpected error: %v", tc.base, tc.exp, err)
				}
				resLREAL, _ := anyToLREAL(result)
				if !almostEqual(resLREAL, tc.expected) {
					t.Errorf("EXPT(%v, %v) = %v; want %v", tc.base, tc.exp, resLREAL, tc.expected)
				}
			}
		})
	}
}

func TestTrigonometric(t *testing.T) {
	pi := LREAL(math.Pi)
	testCases := []struct {
		name     string
		fn       func(interface{}) (interface{}, error)
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
			result, err := tc.fn(tc.input)
			if err != nil {
				t.Fatalf("%s returned an unexpected error: %v", tc.name, err)
			}
			if !almostEqual(result.(LREAL), tc.expected) {
				t.Errorf("%s = %v; want %v", tc.name, result, tc.expected)
			}
		})
	}

	t.Run("Trig with string panic", func(t *testing.T) {
		trigFuncs := map[string]func(interface{}) (interface{}, error){
			"SIN": ASIN, "COS": COS, "TAN": TAN, "ASIN": ASIN, "ACOS": ACOS, "ATAN": ATAN,
		}
		for name, fn := range trigFuncs {
			t.Run(name, func(t *testing.T) {
				_, err := fn(STRING("0.5"))
				if err == nil {
					t.Errorf("%s(STRING) did not return an error", name)
				}
			})
		}
	})
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
		{"Integer input", LINT(100), 0, false, true},
		{"String input", STRING("12.3"), 0, false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := TRUNC(tc.input)
			if tc.expectError {
				if err == nil {
					t.Errorf("TRUNC(%v) did not return an error; expected error", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("TRUNC(%v) returned an unexpected error: %v", tc.input, err)
				}
				// Check the value
				if result != tc.expected {
					t.Errorf("TRUNC(%v) = %v; want %v", tc.input, result, tc.expected)
				}
				// Check the type
				if reflect.TypeOf(result) != reflect.TypeOf(DINT(0)) {
					t.Errorf("TRUNC(%v) returned type %T; want %T", tc.input, result, DINT(0))
				}
			}
		})
	}
}
