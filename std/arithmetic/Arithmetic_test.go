package arithmetic

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	. "github.com/apiarytech/royaljelly/core"
	. "github.com/apiarytech/royaljelly/std/conversion"
)

func TestADD(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"LINTs", func(t *testing.T) {
			result := ADD(LINT(10), LINT(20), LINT(30))
			expected := LINT(60)
			if result != expected {
				t.Errorf("ADD() = %v; want %v", result, expected)
			}
		}},
		{"REALs", func(t *testing.T) {
			result := ADD(REAL(1.5), REAL(2.5))
			expected := REAL(4.0)
			if result != expected {
				t.Errorf("ADD() = %v; want %v", result, expected)
			}
		}},
		{"Empty", func(t *testing.T) {
			result := ADD[DINT]()
			expected := DINT(0)
			if result != expected {
				t.Errorf("ADD() = %v; want %v", result, expected)
			}
		}},
		{"Single", func(t *testing.T) {
			result := ADD(DINT(42))
			expected := DINT(42)
			if result != expected {
				t.Errorf("ADD() = %v; want %v", result, expected)
			}
		}},
		{"TIME", func(t *testing.T) {
			result := ADD(TIME(time.Second), TIME(time.Minute))
			expected := TIME(61 * time.Second)
			if result != expected {
				t.Errorf("ADD() = %v; want %v", result, expected)
			}
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}

func TestSUB(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"LINTs", func(t *testing.T) {
			result := SUB(LINT(100), LINT(20), LINT(30))
			expected := LINT(50)
			if result != expected {
				t.Errorf("SUB() = %v; want %v", result, expected)
			}
		}},
		{"REALs", func(t *testing.T) {
			result := SUB(REAL(10.5), REAL(2.5))
			expected := REAL(8.0)
			if result != expected {
				t.Errorf("SUB() = %v; want %v", result, expected)
			}
		}},
		{"Empty", func(t *testing.T) {
			result := SUB[DINT]()
			expected := DINT(0)
			if result != expected {
				t.Errorf("SUB() = %v; want %v", result, expected)
			}
		}},
		{"Single", func(t *testing.T) {
			result := SUB(DINT(42))
			expected := DINT(42)
			if result != expected {
				t.Errorf("SUB() = %v; want %v", result, expected)
			}
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}

func TestMUL(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"LINTs", func(t *testing.T) {
			result := MUL(LINT(2), LINT(3), LINT(4))
			expected := LINT(24)
			if result != expected {
				t.Errorf("MUL() = %v; want %v", result, expected)
			}
		}},
		{"REALs", func(t *testing.T) {
			result := MUL(REAL(1.5), REAL(2.0))
			expected := REAL(3.0)
			if result != expected {
				t.Errorf("MUL() = %v; want %v", result, expected)
			}
		}},
		{"With zero", func(t *testing.T) {
			result := MUL(DINT(100), DINT(0), DINT(50))
			expected := DINT(0)
			if result != expected {
				t.Errorf("MUL() = %v; want %v", result, expected)
			}
		}},
		{"Empty", func(t *testing.T) {
			result := MUL[LINT]()
			expected := LINT(1)
			if result != expected {
				t.Errorf("MUL() = %v; want %v", result, expected)
			}
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}

func TestDIV(t *testing.T) {
	testCases := []struct {
		name        string
		testFunc    func(*testing.T)
		expectError bool
	}{
		{"LINTs", func(t *testing.T) {
			result, err := DIV(LINT(100), LINT(10), LINT(2))
			expected := LINT(5)
			if err != nil || result != expected {
				t.Errorf("DIV() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false},
		{"REALs", func(t *testing.T) {
			result, err := DIV(REAL(20.0), REAL(4.0))
			expected := REAL(5.0)
			if err != nil || result != expected {
				t.Errorf("DIV() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false},
		{"Empty", func(t *testing.T) {
			result, err := DIV[DINT]()
			expected := DINT(0)
			if err != nil || result != expected {
				t.Errorf("DIV() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false},
		{"Single", func(t *testing.T) {
			result, err := DIV(DINT(42))
			expected := DINT(42)
			if err != nil || result != expected {
				t.Errorf("DIV() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false},
		{"Integer Div by Zero", func(t *testing.T) {
			_, err := DIV(LINT(100), LINT(0))
			if err == nil {
				t.Errorf("DIV() did not return an error; expected error")
			}
		}, true},
		{"LREAL Div by Zero", func(t *testing.T) {
			_, err := DIV(LREAL(100.0), LREAL(0.0))
			if err == nil {
				t.Errorf("DIV() with LREAL did not return an error for division by zero; expected error")
			}
		}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This wrapper allows tc.expectError to be available inside the testFunc
			tc.testFunc(t)
		})
	}
}

func TestMOD(t *testing.T) {
	testCases := []struct {
		testFunc    func(*testing.T)
		expectError bool
		name        string
	}{
		{func(t *testing.T) {
			result, err := MOD(LINT(10), LINT(3))
			expected := LINT(1)
			if err != nil || result != expected {
				t.Errorf("MOD() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false, "LINTs"},
		{func(t *testing.T) {
			result, err := MOD(LINT(25), LINT(12), LINT(2))
			expected := LINT(1)
			if err != nil || result != expected {
				t.Errorf("MOD() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false, "Chain"},
		{func(t *testing.T) {
			result, err := MOD(LINT(-10), LINT(3))
			expected := LINT(-1)
			if err != nil || result != expected {
				t.Errorf("MOD() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false, "With negative"},
		{func(t *testing.T) {
			_, err := MOD(LINT(10), LINT(0))
			if err == nil {
				t.Errorf("MOD() did not return an error; expected error")
			}
		}, true, "Mod by zero"},
		{func(t *testing.T) {
			_, err := MOD(LINT(10))
			if err == nil {
				t.Errorf("MOD() did not return an error; expected error")
			}
		}, true, "Not enough args"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.testFunc)
	}
}

func TestMOVE(t *testing.T) {
	testCases := []struct {
		name string
		test func(*testing.T)
	}{
		{"Move LINT", func(t *testing.T) { MOVE(LINT(123)) }},
		{"Move REAL", func(t *testing.T) { MOVE(REAL(45.6)) }},
		{"Move STRING", func(t *testing.T) { MOVE(STRING("hello")) }},
		{"Move BOOL", func(t *testing.T) { MOVE(BOOL(true)) }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.test)
	}
}

func TestTimeArithmeticFunctions(t *testing.T) {
	t.Run("ADD_TIME", func(t *testing.T) {
		in1 := TIME(time.Hour)
		in2 := TIME(30 * time.Minute)
		expected := TIME(90 * time.Minute)
		result := ADD_TIME(in1, in2)
		if result != expected {
			t.Errorf("ADD_TIME(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("ADD_TOD", func(t *testing.T) {
		in1 := TOD(time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC))
		in2 := TIME(15 * time.Minute)
		expected := TOD(time.Date(0, 0, 0, 10, 15, 0, 0, time.UTC))
		result := ADD_TOD(in1, in2)
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ADD_TOD(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("ADD_DT", func(t *testing.T) {
		in1 := DT(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))
		in2 := TIME(3 * time.Hour)
		expected := DT(time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC))
		result := ADD_DT(in1, in2)
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ADD_DT(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("DIV_TIME by zero", func(t *testing.T) {
		_, err := DIV_TIME(TIME(time.Minute), INT(0))
		if err == nil {
			t.Error("DIV_TIME by zero should have returned an error")
		}
		_, err = DIV_TIME(TIME(time.Minute), "not a number")
		if err == nil {
			t.Error("DIV_TIME with invalid divisor should have returned an error")
		}
	})

	t.Run("SUB_TIME", func(t *testing.T) {
		in1 := TIME(time.Hour)
		in2 := TIME(20 * time.Minute)
		expected := TIME(40 * time.Minute)
		result := SUB_TIME(in1, in2)
		if result != expected {
			t.Errorf("SUB_TIME(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("SUB_DATE", func(t *testing.T) {
		in1 := DATE(time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC))
		in2 := DATE(time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC))
		expected := TIME(10 * 24 * time.Hour)
		result := SUB_DATE(in1, in2)
		if result != expected {
			t.Errorf("SUB_DATE(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("SUB_TOD", func(t *testing.T) {
		tod1 := TOD(time.Date(0, 0, 0, 12, 0, 0, 0, time.UTC))
		time1 := TIME(time.Hour)
		tod2 := TOD(time.Date(0, 0, 0, 11, 0, 0, 0, time.UTC))

		// TOD - TIME -> TOD
		expected1 := TOD(time.Date(0, 0, 0, 11, 0, 0, 0, time.UTC))
		result1 := SUB_TOD_TIME(tod1, time1)
		if !time.Time(result1).Equal(time.Time(expected1)) {
			t.Errorf("SUB_TOD(TOD, TIME) = %v; want %v", result1, expected1)
		}

		// TOD - TOD -> TIME
		expected2 := TIME(time.Hour)
		result2 := SUB_TOD_TOD(tod1, tod2)
		if result2 != expected2 {
			t.Errorf("SUB_TOD(TOD, TOD) = %v; want %v", result2, expected2)
		}
	})

	t.Run("SUB_DT", func(t *testing.T) {
		dt1 := DT(time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC))
		time1 := TIME(24 * time.Hour)
		dt2 := DT(time.Date(2024, 3, 14, 12, 0, 0, 0, time.UTC))

		// DT - TIME -> DT
		expected1 := DT(time.Date(2024, 3, 14, 12, 0, 0, 0, time.UTC))
		result1 := SUB_DT_TIME(dt1, time1)
		if !time.Time(result1).Equal(time.Time(expected1)) {
			t.Errorf("SUB_DT(DT, TIME) = %v; want %v", result1, expected1)
		}

		// DT - DT -> TIME
		expected2 := TIME(24 * time.Hour)
		result2 := SUB_DT_DT(dt1, dt2)
		if result2 != expected2 {
			t.Errorf("SUB_DT(DT, DT) = %v; want %v", result2, expected2)
		}
	})

	t.Run("MUL_TIME", func(t *testing.T) {
		in1 := TIME(10 * time.Second)
		in2 := INT(6)
		expected := TIME(time.Minute)
		result, err := MUL_TIME(in1, in2)
		if err != nil {
			t.Fatalf("MUL_TIME returned an unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("MUL_TIME(%v, %v) = %v; want %v", in1, in2, result, expected)
		}

		in3 := REAL(2.5)
		expected2 := TIME(25 * time.Second)
		result2, err := MUL_TIME(in1, in3)
		if err != nil {
			t.Fatalf("MUL_TIME(TIME, REAL) returned an unexpected error: %v", err)
		}
		if result2 != expected2 {
			t.Errorf("MUL_TIME(TIME, REAL) = %v; want %v", result2, expected2)
		}
	})

	t.Run("DIV_TIME", func(t *testing.T) {
		in1 := TIME(time.Minute)
		in2 := INT(4)
		expected := TIME(15 * time.Second)
		result, err := DIV_TIME(in1, in2)
		if err != nil {
			t.Fatalf("DIV_TIME returned an unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("DIV_TIME(%v, %v) = %v; want %v", in1, in2, result, expected)
		}

		in3 := REAL(2.5)
		expected2 := TIME(24 * time.Second)
		result2, err := DIV_TIME(in1, in3)
		if err != nil {
			t.Fatalf("DIV_TIME(TIME, REAL) returned an unexpected error: %v", err)
		}
		if result2 != expected2 {
			t.Errorf("DIV_TIME(TIME, REAL) = %v; want %v", result2, expected2)
		}
	})

	t.Run("CONCAT_DATE_TOD", func(t *testing.T) {
		in1 := DATE(time.Date(2025, 10, 21, 0, 0, 0, 0, time.UTC))
		in2 := TOD(time.Date(0, 0, 0, 16, 30, 5, 0, time.UTC))
		expected := DT(time.Date(2025, 10, 21, 16, 30, 5, 0, time.UTC))
		result := CONCAT_DATE_TOD(in1, in2)
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("CONCAT_DATE_TOD(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})
}

func TestConversionError(t *testing.T) {
	t.Run("STRING to LINT failure", func(t *testing.T) {
		invalidInput := STRING("not-a-number")
		_, err := AnyToLINT(invalidInput)

		if err == nil {
			t.Fatal("AnyToLINT with invalid string should have returned an error, but got nil")
		}

		var convErr *ConversionError
		if !errors.As(err, &convErr) {
			t.Fatalf("Expected error of type *ConversionError, but got %T", err)
		}

		if convErr.Value != invalidInput {
			t.Errorf("ConversionError.Value = %v; want %v", convErr.Value, invalidInput)
		}
		if convErr.FromType != "STRING" {
			t.Errorf("ConversionError.FromType = %q; want 'STRING'", convErr.FromType)
		}
		if convErr.ToType != "LINT" {
			t.Errorf("ConversionError.ToType = %q; want 'LINT'", convErr.ToType)
		}
		if convErr.Reason != "string could not be parsed as an integer" {
			t.Errorf("ConversionError.Reason = %q; want 'string could not be parsed as an integer'", convErr.Reason)
		}
		if convErr.Err == nil {
			t.Error("ConversionError.Err should not be nil for a parse error")
		}
	})

	t.Run("Unsupported type to LREAL failure", func(t *testing.T) {
		invalidInput := struct{}{} // An unsupported type
		_, err := AnyToLREAL(invalidInput)

		if err == nil {
			t.Fatal("AnyToLREAL with unsupported type should have returned an error, but got nil")
		}

		var convErr *ConversionError
		if !errors.As(err, &convErr) {
			t.Fatalf("Expected error of type *ConversionError, but got %T", err)
		}

		if convErr.ToType != "LREAL" {
			t.Errorf("ConversionError.ToType = %q; want 'LREAL'", convErr.ToType)
		}
	})
}

func TestAnyToLREAL(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected LREAL
		hasError bool
	}{
		{"SINT", SINT(-10), LREAL(-10), false},
		{"UINT", UINT(100), LREAL(100), false},
		{"REAL", REAL(123.45), LREAL(123.44999694824219), false},
		{"BOOL true", BOOL(true), LREAL(1.0), false},
		{"BOOL false", BOOL(false), LREAL(0.0), false},
		{"STRING int", STRING("123"), LREAL(123), false},
		{"STRING float", STRING("123.45"), LREAL(123.45), false},
		{"STRING invalid", STRING("abc"), 0, true},
		{"Unsupported type", struct{}{}, 0, true},
		{"TIME", TIME(2 * time.Second), LREAL(2000), false},
		{"DINT", DINT(12345), LREAL(12345), false},
		{"LINT", LINT(123456789), LREAL(123456789), false},
		{"USINT", USINT(255), LREAL(255), false},
		{"UDINT", UDINT(40000), LREAL(40000), false},
		{"ULINT", ULINT(9876543210), LREAL(9876543210), false},
		{"LREAL", LREAL(987.65), LREAL(987.65), false},
		{"native float32", float32(1.23), LREAL(1.230000019073486300), false},
		{"native float64", float64(4.56), LREAL(4.56), false},
		{"BYTE", BYTE(0xAB), LREAL(171), false},
		{"WORD", WORD(0xABCD), LREAL(43981), false},
		{"DWORD", DWORD(0xABCDEF), LREAL(11259375), false},
		{"LWORD", LWORD(0x1234567890ABCDEF), LREAL(1311768467294899700), false},
		{"DATE", DATE(time.UnixMilli(1678886400000)), LREAL(1678886400000), false},        // 2023-03-15 12:00:00 UTC
		{"TOD", TOD(time.Date(0, 1, 1, 14, 30, 15, 0, time.UTC)), LREAL(52215000), false}, // 14h, 30m, 15s in ms
		{"DT", DT(time.UnixMilli(1678886400000)), LREAL(1678886400000), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a tolerance for float comparisons
			const tolerance = 1e-9
			result, err := AnyToLREAL(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("AnyToLREAL(%v) expected an error, but got none", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("anyToLREAL(%v) returned an error: %v", tc.input, err)
				}
				// Comparing floats can be tricky due to precision.
				// For this test, direct comparison is fine for most cases, but a tolerance check is more robust.
				if diff := LREAL(result - tc.expected); diff < -tolerance || diff > tolerance {
					// The conversion from REAL to LREAL might introduce tiny precision differences.
					// We re-check by converting the input to string and comparing.
					if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
						t.Errorf("anyToLREAL(%v) = %g; want %g", tc.input, result, tc.expected)
					}
				}
			}
		})
	}
}

func TestAnyToLINT(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected LINT
		hasError bool
	}{
		{"SINT", SINT(-10), LINT(-10), false},
		{"INT", INT(-200), LINT(-200), false},
		{"DINT", DINT(30000), LINT(30000), false},
		{"LINT", LINT(1234567890), LINT(1234567890), false},
		{"USINT", USINT(250), LINT(250), false},
		{"UINT", UINT(100), LINT(100), false},
		{"UDINT", UDINT(65000), LINT(65000), false},
		{"ULINT", ULINT(123456789012345), LINT(123456789012345), false},
		{"ULINT overflow", ULINT(math.MaxInt64 + 1), LINT(math.MinInt64), false},
		{"REAL truncates", REAL(123.75), LINT(123), false},
		{"LREAL truncates", LREAL(-456.99), LINT(-456), false},
		{"BOOL true", BOOL(true), LINT(1), false},
		{"BOOL false", BOOL(false), LINT(0), false},
		{"BYTE", BYTE(0xFE), LINT(254), false},
		{"WORD", WORD(0xFFFE), LINT(65534), false},
		{"DWORD", DWORD(0xFFFFFFFE), LINT(4294967294), false},
		{"LWORD", LWORD(0x100000000), LINT(4294967296), false},
		{"LWORD overflow", LWORD(math.MaxUint64), LINT(-1), false},
		{"STRING int", STRING("123"), LINT(123), false},
		{"STRING invalid", STRING("abc"), 0, true},
		{"STRING float", STRING("123.45"), 0, true}, // ParseInt fails on floats
		{"TIME", TIME(3 * time.Second), LINT(3000), false},
		{"DATE", DATE(time.UnixMilli(1678886400000)), LINT(1678886400000), false},
		{"TOD", TOD(time.Date(0, 0, 0, 1, 2, 3, 0, time.UTC)), LINT(3723000), false},
		{"DT", DT(time.UnixMilli(1678886400000)), LINT(1678886400000), false},
		{"Unsupported type", struct{}{}, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := AnyToLINT(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("AnyToLINT(%v) expected an error, but got none", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("AnyToLINT(%v) returned an error: %v", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("AnyToLINT(%v) = %d; want %d", tc.input, result, tc.expected)
				}
			}
		})
	}
}

func TestAnyToULINT(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected ULINT
		hasError bool
	}{
		{"SINT positive", SINT(10), ULINT(10), false},
		{"SINT negative", SINT(-10), ULINT(0xFFFFFFFFFFFFFFF6), false}, // two's complement
		{"INT", INT(30000), ULINT(30000), false},
		{"DINT", DINT(-50000), ULINT(0xFFFFFFFFFFFF3CB0), false},
		{"LINT", LINT(123456789), ULINT(123456789), false},
		{"USINT", USINT(255), ULINT(255), false},
		{"UINT", UINT(100), ULINT(100), false},
		{"UDINT", UDINT(4000000000), ULINT(4000000000), false},
		{"ULINT", ULINT(999999999999), ULINT(999999999999), false},
		{"REAL truncates", REAL(123.75), ULINT(123), false},
		{"LREAL truncates negative", LREAL(-456.99), LREAL_TO_ULINT(-456), false},
		{"BOOL true", BOOL(true), ULINT(1), false},
		{"BOOL false", BOOL(false), ULINT(0), false},
		{"BYTE", BYTE(0xAB), ULINT(171), false},
		{"WORD", WORD(0xABCD), ULINT(43981), false},
		{"DWORD", DWORD(0xABCDEF01), ULINT(2882400001), false},
		{"LWORD", LWORD(0x1234567890ABCDEF), ULINT(1311768467294899695), false},
		{"STRING int", STRING("123"), ULINT(123), false},
		{"STRING hex", STRING("0xFF"), ULINT(255), false},
		{"STRING invalid", STRING("abc"), 0, true},
		{"TIME", TIME(3 * time.Second), ULINT(3000), false},
		{"DT", DT(time.UnixMilli(1678886400000)), ULINT(1678886400000), false},
		{"DATE", DATE(time.UnixMilli(1678886400000)), ULINT(1678886400000), false},
		{"TOD", TOD(time.Date(0, 0, 0, 1, 2, 3, 0, time.UTC)), ULINT(3723000), false},
		{"Unsupported type", struct{}{}, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := AnyToULINT(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("AnyToULINT(%v) expected an error, but got none", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("AnyToULINT(%v) returned an error: %v", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("AnyToULINT(%v) = %d; want %d", tc.input, result, tc.expected)
				}
			}
		})
	}
}

func TestAnyToREAL(t *testing.T) {
	t.Run("Valid conversion", func(t *testing.T) {
		result, err := AnyToREAL(LINT(123))
		if err != nil {
			t.Fatalf("AnyToREAL failed with error: %v", err)
		}
		if result != REAL(123.0) {
			t.Errorf("AnyToREAL(LINT(123)) = %f; want 123.0", result)
		}
	})

	t.Run("Invalid conversion", func(t *testing.T) {
		_, err := AnyToREAL(STRING("abc"))
		if err == nil {
			t.Error("AnyToREAL(STRING(\"abc\")) should have returned an error")
		}
	})
}

func TestConvertTo(t *testing.T) {
	t.Run("LREAL_TO_DINT", func(t *testing.T) {
		input := LREAL(123.7)
		expected := DINT(123) // Note: LREAL_TO_DINT conversion truncates
		result, err := ConvertTo[DINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[DINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[DINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("SINT_TO_LREAL", func(t *testing.T) {
		input := SINT(-50)
		expected := LREAL(-50.0)
		result, err := ConvertTo[LREAL](input)
		if err != nil {
			t.Fatalf("ConvertTo[LREAL] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[LREAL](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("STRING_TO_LINT", func(t *testing.T) {
		input := STRING("12345")
		expected := LINT(12345)
		result, err := ConvertTo[LINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[LINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[LINT](%q) = %v; want %v", input, result, expected)
		}
	})

	t.Run("TIME_TO_STRING", func(t *testing.T) {
		input := TIME(5 * time.Second)
		expected := STRING("T#5s")
		result, err := ConvertTo[STRING](input)
		if err != nil {
			t.Fatalf("ConvertTo[STRING] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[STRING](%v) = %q; want %q", input, result, expected)
		}
	})

	t.Run("UINT_TO_UDINT", func(t *testing.T) {
		input := UINT(65000)
		expected := UDINT(65000)
		result, err := ConvertTo[UDINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[UDINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[UDINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_BYTE", func(t *testing.T) {
		input := LINT(255)
		expected := BYTE(255)
		result, err := ConvertTo[BYTE](input)
		if err != nil {
			t.Fatalf("ConvertTo[BYTE] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[BYTE](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("INT_TO_TIME", func(t *testing.T) {
		input := INT(5000)
		expected := TIME(5 * time.Second)
		result, err := ConvertTo[TIME](input)
		if err != nil {
			t.Fatalf("ConvertTo[TIME] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[TIME](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_SINT", func(t *testing.T) {
		input := LINT(120)
		expected := SINT(120)
		result, err := ConvertTo[SINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[SINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[SINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("DINT_TO_UINT", func(t *testing.T) {
		input := DINT(65000)
		expected := UINT(65000)
		result, err := ConvertTo[UINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[UINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[UINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("REAL_TO_ULINT", func(t *testing.T) {
		input := REAL(1234567.8)
		expected := ULINT(1234567)
		result, err := ConvertTo[ULINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[ULINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[ULINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("INT_TO_BOOL", func(t *testing.T) {
		input := INT(1)
		expected := BOOL(true)
		result, err := ConvertTo[BOOL](input)
		if err != nil {
			t.Fatalf("ConvertTo[BOOL] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[BOOL](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("STRING_TO_BOOL", func(t *testing.T) {
		input := STRING("true")
		expected := BOOL(true)
		result, err := ConvertTo[BOOL](input)
		if err != nil {
			t.Fatalf("ConvertTo[BOOL] from STRING failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[BOOL](%q) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_WORD", func(t *testing.T) {
		input := LINT(0xABCD)
		expected := WORD(0xABCD)
		result, err := ConvertTo[WORD](input)
		if err != nil {
			t.Fatalf("ConvertTo[WORD] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[WORD](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_DATE", func(t *testing.T) {
		ms := int64(1678886400000)
		input := LINT(ms)
		expected := DATE(time.UnixMilli(ms))
		result, err := ConvertTo[DATE](input)
		if err != nil {
			t.Fatalf("ConvertTo[DATE] failed: %v", err)
		}
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ConvertTo[DATE](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LREAL_TO_INT", func(t *testing.T) {
		input := LREAL(32767.8)
		expected := INT(32767) // Clamped
		result, err := ConvertTo[INT](input)
		if err != nil {
			t.Fatalf("ConvertTo[INT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[INT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("DINT_TO_USINT", func(t *testing.T) {
		input := DINT(250)
		expected := USINT(250)
		result, err := ConvertTo[USINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[USINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[USINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_DWORD", func(t *testing.T) {
		input := LINT(0xABCDEF)
		expected := DWORD(0xABCDEF)
		result, err := ConvertTo[DWORD](input)
		if err != nil {
			t.Fatalf("ConvertTo[DWORD] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[DWORD](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("UINT_TO_LWORD", func(t *testing.T) {
		input := UINT(0xFFFF)
		expected := LWORD(0xFFFF)
		result, err := ConvertTo[LWORD](input)
		if err != nil {
			t.Fatalf("ConvertTo[LWORD] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[LWORD](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_TOD", func(t *testing.T) {
		ms := int64(3723000) // 1h, 2m, 3s
		input := LINT(ms)
		expected := TOD(time.Time{}.Add(time.Duration(ms) * time.Millisecond))
		result, err := ConvertTo[TOD](input)
		if err != nil {
			t.Fatalf("ConvertTo[TOD] failed: %v", err)
		}
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ConvertTo[TOD](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_DT", func(t *testing.T) {
		ms := int64(1678886400000)
		input := LINT(ms)
		expected := DT(time.UnixMilli(ms))
		result, err := ConvertTo[DT](input)
		if err != nil {
			t.Fatalf("ConvertTo[DT] failed: %v", err)
		}
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ConvertTo[DT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("DINT_TO_REAL", func(t *testing.T) {
		input := DINT(12345)
		expected := REAL(12345.0)
		result, err := ConvertTo[REAL](input)
		if err != nil {
			t.Fatalf("ConvertTo[REAL] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[REAL](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("Invalid STRING_TO_LINT", func(t *testing.T) {
		input := STRING("abc")
		_, err := ConvertTo[LINT](input)
		if err == nil {
			t.Errorf("Expected an error for invalid string conversion, but got nil")
		}
	})

	t.Run("Unsupported Target Type", func(t *testing.T) {
		// Use a type that is not handled in the ConvertTo switch statement.
		input := LINT(123)
		_, err := ConvertTo[struct{}](input)
		if err == nil {
			t.Error("Expected an error for unsupported target type, but got nil")
		}
	})
}

func TestTimeSpecificFunctions(t *testing.T) {
	// This function can be used for more specific time-related tests if needed.
	// Currently, the main time arithmetic is covered in TestTimeArithmeticFunctions.
}

func TestTypeCheckFunctions(t *testing.T) {
	t.Run("IsPlcFloat", func(t *testing.T) {
		if !IsPlcFloat(REAL(1.0)) {
			t.Error("IsPlcFloat(REAL) should be true")
		}
		if !IsPlcFloat(LREAL(1.0)) {
			t.Error("IsPlcFloat(LREAL) should be true")
		}
		if IsPlcFloat(DINT(1)) {
			t.Error("IsPlcFloat(DINT) should be false")
		}
	})

	t.Run("IsPlcInt", func(t *testing.T) {
		if !IsPlcInt(SINT(1)) {
			t.Error("IsPlcInt(SINT) should be true")
		}
		if !IsPlcInt(UINT(1)) {
			t.Error("IsPlcInt(UINT) should be true")
		}
		if !IsPlcInt(BOOL(true)) {
			t.Error("IsPlcInt(BOOL) should be true")
		}
		if IsPlcInt(REAL(1.0)) {
			t.Error("IsPlcInt(REAL) should be false")
		}
	})

	t.Run("IsPlcTimeType", func(t *testing.T) {
		if !IsPlcTimeType(TIME(0)) {
			t.Error("IsPlcTimeType(TIME) should be true")
		}
		if !IsPlcTimeType(DT{}) {
			t.Error("IsPlcTimeType(DT) should be true")
		}
		if IsPlcTimeType(LINT(0)) {
			t.Error("IsPlcTimeType(LINT) should be false")
		}
	})
}
