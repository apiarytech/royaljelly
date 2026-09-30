package arithmetic

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/convert"
	"github.com/apiarytech/royaljelly/iec"
	"github.com/apiarytech/royaljelly/std/conversion"
)

func TestADD(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"LINTs", func(t *testing.T) {
			result := ADD(iec.LINT(10), iec.LINT(20), iec.LINT(30))
			expected := iec.LINT(60)
			if result != expected {
				t.Errorf("ADD() = %v; want %v", result, expected)
			}
		}},
		{"REALs", func(t *testing.T) {
			result := ADD(iec.REAL(1.5), iec.REAL(2.5))
			expected := iec.REAL(4.0)
			if result != expected {
				t.Errorf("ADD() = %v; want %v", result, expected)
			}
		}},
		{"Empty", func(t *testing.T) {
			result := ADD[iec.DINT]()
			expected := iec.DINT(0)
			if result != expected {
				t.Errorf("ADD() = %v; want %v", result, expected)
			}
		}},
		{"Single", func(t *testing.T) {
			result := ADD(iec.DINT(42))
			expected := iec.DINT(42)
			if result != expected {
				t.Errorf("ADD() = %v; want %v", result, expected)
			}
		}},
		{"TIME", func(t *testing.T) {
			result := ADD(iec.TIME(time.Second), iec.TIME(time.Minute))
			expected := iec.TIME(61 * time.Second)
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
			result := SUB(iec.LINT(100), iec.LINT(20), iec.LINT(30))
			expected := iec.LINT(50)
			if result != expected {
				t.Errorf("SUB() = %v; want %v", result, expected)
			}
		}},
		{"REALs", func(t *testing.T) {
			result := SUB(iec.REAL(10.5), iec.REAL(2.5))
			expected := iec.REAL(8.0)
			if result != expected {
				t.Errorf("SUB() = %v; want %v", result, expected)
			}
		}},
		{"Empty", func(t *testing.T) {
			result := SUB[iec.DINT]()
			expected := iec.DINT(0)
			if result != expected {
				t.Errorf("SUB() = %v; want %v", result, expected)
			}
		}},
		{"Single", func(t *testing.T) {
			result := SUB(iec.DINT(42))
			expected := iec.DINT(42)
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
			result := MUL(iec.LINT(2), iec.LINT(3), iec.LINT(4))
			expected := iec.LINT(24)
			if result != expected {
				t.Errorf("MUL() = %v; want %v", result, expected)
			}
		}},
		{"REALs", func(t *testing.T) {
			result := MUL(iec.REAL(1.5), iec.REAL(2.0))
			expected := iec.REAL(3.0)
			if result != expected {
				t.Errorf("MUL() = %v; want %v", result, expected)
			}
		}},
		{"With zero", func(t *testing.T) {
			result := MUL(iec.DINT(100), iec.DINT(0), iec.DINT(50))
			expected := iec.DINT(0)
			if result != expected {
				t.Errorf("MUL() = %v; want %v", result, expected)
			}
		}},
		{"Empty", func(t *testing.T) {
			result := MUL[iec.LINT]()
			expected := iec.LINT(1)
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
			result, err := DIV(iec.LINT(100), iec.LINT(10), iec.LINT(2))
			expected := iec.LINT(5)
			if err != nil || result != expected {
				t.Errorf("DIV() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false},
		{"REALs", func(t *testing.T) {
			result, err := DIV(iec.REAL(20.0), iec.REAL(4.0))
			expected := iec.REAL(5.0)
			if err != nil || result != expected {
				t.Errorf("DIV() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false},
		{"Empty", func(t *testing.T) {
			result, err := DIV[iec.DINT]()
			expected := iec.DINT(0)
			if err != nil || result != expected {
				t.Errorf("DIV() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false},
		{"Single", func(t *testing.T) {
			result, err := DIV(iec.DINT(42))
			expected := iec.DINT(42)
			if err != nil || result != expected {
				t.Errorf("DIV() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false},
		{"Integer Div by Zero", func(t *testing.T) {
			_, err := DIV(iec.LINT(100), iec.LINT(0))
			if err == nil {
				t.Errorf("DIV() did not return an error; expected error")
			}
		}, true},
		{"LREAL Div by Zero", func(t *testing.T) {
			_, err := DIV(iec.LREAL(100.0), iec.LREAL(0.0))
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
			result, err := MOD(iec.LINT(10), iec.LINT(3))
			expected := iec.LINT(1)
			if err != nil || result != expected {
				t.Errorf("MOD() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false, "LINTs"},
		{func(t *testing.T) {
			result, err := MOD(iec.LINT(25), iec.LINT(12), iec.LINT(2))
			expected := iec.LINT(1)
			if err != nil || result != expected {
				t.Errorf("MOD() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false, "Chain"},
		{func(t *testing.T) {
			result, err := MOD(iec.LINT(-10), iec.LINT(3))
			expected := iec.LINT(-1)
			if err != nil || result != expected {
				t.Errorf("MOD() = %v, err: %v; want %v, nil", result, err, expected)
			}
		}, false, "With negative"},
		{func(t *testing.T) {
			_, err := MOD(iec.LINT(10), iec.LINT(0))
			if err == nil {
				t.Errorf("MOD() did not return an error; expected error")
			}
		}, true, "Mod by zero"},
		{func(t *testing.T) {
			_, err := MOD(iec.LINT(10))
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
		{"Move LINT", func(t *testing.T) { MOVE(iec.LINT(123)) }},
		{"Move REAL", func(t *testing.T) { MOVE(iec.REAL(45.6)) }},
		{"Move STRING", func(t *testing.T) { MOVE(iec.STRING("hello")) }},
		{"Move BOOL", func(t *testing.T) { MOVE(iec.BOOL(true)) }},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.test)
	}
}

func TestTimeArithmeticFunctions(t *testing.T) {
	t.Run("ADD_TIME", func(t *testing.T) {
		in1 := iec.TIME(time.Hour)
		in2 := iec.TIME(30 * time.Minute)
		expected := iec.TIME(90 * time.Minute)
		result := ADD_TIME(in1, in2)
		if result != expected {
			t.Errorf("ADD_TIME(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("ADD_TOD", func(t *testing.T) {
		in1 := iec.TOD(time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC))
		in2 := iec.TIME(15 * time.Minute)
		expected := iec.TOD(time.Date(0, 0, 0, 10, 15, 0, 0, time.UTC))
		result := ADD_TOD(in1, in2)
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ADD_TOD(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("ADD_DT", func(t *testing.T) {
		in1 := iec.DT(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))
		in2 := iec.TIME(3 * time.Hour)
		expected := iec.DT(time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC))
		result := ADD_DT(in1, in2)
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ADD_DT(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("DIV_TIME by zero", func(t *testing.T) {
		_, err := DIV_TIME(iec.TIME(time.Minute), iec.INT(0))
		if err == nil {
			t.Error("DIV_TIME by zero should have returned an error")
		}
		_, err = DIV_TIME(iec.TIME(time.Minute), "not a number")
		if err == nil {
			t.Error("DIV_TIME with invalid divisor should have returned an error")
		}
	})

	t.Run("SUB_TIME", func(t *testing.T) {
		in1 := iec.TIME(time.Hour)
		in2 := iec.TIME(20 * time.Minute)
		expected := iec.TIME(40 * time.Minute)
		result := SUB_TIME(in1, in2)
		if result != expected {
			t.Errorf("SUB_TIME(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("SUB_DATE", func(t *testing.T) {
		in1 := iec.DATE(time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC))
		in2 := iec.DATE(time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC))
		expected := iec.TIME(10 * 24 * time.Hour)
		result := SUB_DATE(in1, in2)
		if result != expected {
			t.Errorf("SUB_DATE(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("SUB_TOD", func(t *testing.T) {
		tod1 := iec.TOD(time.Date(0, 0, 0, 12, 0, 0, 0, time.UTC))
		time1 := iec.TIME(time.Hour)
		tod2 := iec.TOD(time.Date(0, 0, 0, 11, 0, 0, 0, time.UTC))

		// TOD - TIME -> TOD
		expected1 := iec.TOD(time.Date(0, 0, 0, 11, 0, 0, 0, time.UTC))
		result1 := SUB_TOD_TIME(tod1, time1)
		if !time.Time(result1).Equal(time.Time(expected1)) {
			t.Errorf("SUB_TOD(TOD, TIME) = %v; want %v", result1, expected1)
		}

		// TOD - TOD -> TIME
		expected2 := iec.TIME(time.Hour)
		result2 := SUB_TOD_TOD(tod1, tod2)
		if result2 != expected2 {
			t.Errorf("SUB_TOD(TOD, TOD) = %v; want %v", result2, expected2)
		}
	})

	t.Run("SUB_DT", func(t *testing.T) {
		dt1 := iec.DT(time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC))
		time1 := iec.TIME(24 * time.Hour)
		dt2 := iec.DT(time.Date(2024, 3, 14, 12, 0, 0, 0, time.UTC))

		// DT - TIME -> DT
		expected1 := iec.DT(time.Date(2024, 3, 14, 12, 0, 0, 0, time.UTC))
		result1 := SUB_DT_TIME(dt1, time1)
		if !time.Time(result1).Equal(time.Time(expected1)) {
			t.Errorf("SUB_DT(DT, TIME) = %v; want %v", result1, expected1)
		}

		// DT - DT -> TIME
		expected2 := iec.TIME(24 * time.Hour)
		result2 := SUB_DT_DT(dt1, dt2)
		if result2 != expected2 {
			t.Errorf("SUB_DT(DT, DT) = %v; want %v", result2, expected2)
		}
	})

	t.Run("MUL_TIME", func(t *testing.T) {
		in1 := iec.TIME(10 * time.Second)
		in2 := iec.INT(6)
		expected := iec.TIME(time.Minute)
		result, err := MUL_TIME(in1, in2)
		if err != nil {
			t.Fatalf("MUL_TIME returned an unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("MUL_TIME(%v, %v) = %v; want %v", in1, in2, result, expected)
		}

		in3 := iec.REAL(2.5)
		expected2 := iec.TIME(25 * time.Second)
		result2, err := MUL_TIME(in1, in3)
		if err != nil {
			t.Fatalf("MUL_TIME(TIME, REAL) returned an unexpected error: %v", err)
		}
		if result2 != expected2 {
			t.Errorf("MUL_TIME(TIME, REAL) = %v; want %v", result2, expected2)
		}
	})

	t.Run("DIV_TIME", func(t *testing.T) {
		in1 := iec.TIME(time.Minute)
		in2 := iec.INT(4)
		expected := iec.TIME(15 * time.Second)
		result, err := DIV_TIME(in1, in2)
		if err != nil {
			t.Fatalf("DIV_TIME returned an unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("DIV_TIME(%v, %v) = %v; want %v", in1, in2, result, expected)
		}

		in3 := iec.REAL(2.5)
		expected2 := iec.TIME(24 * time.Second)
		result2, err := DIV_TIME(in1, in3)
		if err != nil {
			t.Fatalf("DIV_TIME(TIME, REAL) returned an unexpected error: %v", err)
		}
		if result2 != expected2 {
			t.Errorf("DIV_TIME(TIME, REAL) = %v; want %v", result2, expected2)
		}
	})

	t.Run("CONCAT_DATE_TOD", func(t *testing.T) {
		in1 := iec.DATE(time.Date(2025, 10, 21, 0, 0, 0, 0, time.UTC))
		in2 := iec.TOD(time.Date(0, 0, 0, 16, 30, 5, 0, time.UTC))
		expected := iec.DT(time.Date(2025, 10, 21, 16, 30, 5, 0, time.UTC))
		result := CONCAT_DATE_TOD(in1, in2)
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("CONCAT_DATE_TOD(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})
}

func TestConversionError(t *testing.T) {
	t.Run("STRING to LINT failure", func(t *testing.T) {
		invalidInput := iec.STRING("not-a-number")
		_, err := convert.AnyToLINT(invalidInput)

		if err == nil {
			t.Fatal("AnyToLINT with invalid string should have returned an error, but got nil")
		}

		var convErr *convert.ConversionError
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
		_, err := convert.AnyToLREAL(invalidInput)

		if err == nil {
			t.Fatal("AnyToLREAL with unsupported type should have returned an error, but got nil")
		}

		var convErr *convert.ConversionError
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
		expected iec.LREAL
		hasError bool
	}{
		{"SINT", iec.SINT(-10), iec.LREAL(-10), false},
		{"UINT", iec.UINT(100), iec.LREAL(100), false},
		{"REAL", iec.REAL(123.45), iec.LREAL(123.44999694824219), false},
		{"BOOL true", iec.BOOL(true), iec.LREAL(1.0), false},
		{"BOOL false", iec.BOOL(false), iec.LREAL(0.0), false},
		{"STRING int", iec.STRING("123"), iec.LREAL(123), false},
		{"STRING float", iec.STRING("123.45"), iec.LREAL(123.45), false},
		{"STRING invalid", iec.STRING("abc"), 0, true},
		{"Unsupported type", struct{}{}, 0, true},
		{"TIME", iec.TIME(2 * time.Second), iec.LREAL(2000), false},
		{"DINT", iec.DINT(12345), iec.LREAL(12345), false},
		{"LINT", iec.LINT(123456789), iec.LREAL(123456789), false},
		{"USINT", iec.USINT(255), iec.LREAL(255), false},
		{"UDINT", iec.UDINT(40000), iec.LREAL(40000), false},
		{"ULINT", iec.ULINT(9876543210), iec.LREAL(9876543210), false},
		{"LREAL", iec.LREAL(987.65), iec.LREAL(987.65), false},
		{"native float32", float32(1.23), iec.LREAL(1.230000019073486300), false},
		{"native float64", float64(4.56), iec.LREAL(4.56), false},
		{"BYTE", iec.BYTE(0xAB), iec.LREAL(171), false},
		{"WORD", iec.WORD(0xABCD), iec.LREAL(43981), false},
		{"DWORD", iec.DWORD(0xABCDEF), iec.LREAL(11259375), false},
		{"LWORD", iec.LWORD(0x1234567890ABCDEF), iec.LREAL(1311768467294899700), false},
		{"DATE", iec.DATE(time.UnixMilli(1678886400000)), iec.LREAL(1678886400000), false},        // 2023-03-15 12:00:00 UTC
		{"TOD", iec.TOD(time.Date(0, 1, 1, 14, 30, 15, 0, time.UTC)), iec.LREAL(52215000), false}, // 14h, 30m, 15s in ms
		{"DT", iec.DT(time.UnixMilli(1678886400000)), iec.LREAL(1678886400000), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Use a tolerance for float comparisons
			const tolerance = 1e-9
			result, err := convert.AnyToLREAL(tc.input)

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
				if diff := iec.LREAL(result - tc.expected); diff < -tolerance || diff > tolerance {
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
		expected iec.LINT
		hasError bool
	}{
		{"SINT", iec.SINT(-10), iec.LINT(-10), false},
		{"INT", iec.INT(-200), iec.LINT(-200), false},
		{"DINT", iec.DINT(30000), iec.LINT(30000), false},
		{"LINT", iec.LINT(1234567890), iec.LINT(1234567890), false},
		{"USINT", iec.USINT(250), iec.LINT(250), false},
		{"UINT", iec.UINT(100), iec.LINT(100), false},
		{"UDINT", iec.UDINT(65000), iec.LINT(65000), false},
		{"ULINT", iec.ULINT(123456789012345), iec.LINT(123456789012345), false},
		{"ULINT overflow", iec.ULINT(math.MaxInt64 + 1), iec.MAXLINT, true},
		{"REAL truncates", iec.REAL(123.75), iec.LINT(123), false},
		{"LREAL truncates", iec.LREAL(-456.99), iec.LINT(-456), false},
		{"BOOL true", iec.BOOL(true), iec.LINT(1), false},
		{"BOOL false", iec.BOOL(false), iec.LINT(0), false},
		{"BYTE", iec.BYTE(0xFE), iec.LINT(254), false},
		{"WORD", iec.WORD(0xFFFE), iec.LINT(65534), false},
		{"DWORD", iec.DWORD(0xFFFFFFFE), iec.LINT(4294967294), false},
		{"LWORD", iec.LWORD(0x100000000), iec.LINT(4294967296), false},
		{"LWORD overflow", iec.LWORD(math.MaxUint64), iec.LINT(-1), false},
		{"STRING int", iec.STRING("123"), iec.LINT(123), false},
		{"STRING invalid", iec.STRING("abc"), 0, true},
		{"STRING float", iec.STRING("123.45"), 0, true}, // ParseInt fails on floats
		{"TIME", iec.TIME(3 * time.Second), iec.LINT(3000), false},
		{"DATE", iec.DATE(time.UnixMilli(1678886400000)), iec.LINT(1678886400000), false},
		{"TOD", iec.TOD(time.Date(0, 0, 0, 1, 2, 3, 0, time.UTC)), iec.LINT(3723000), false},
		{"DT", iec.DT(time.UnixMilli(1678886400000)), iec.LINT(1678886400000), false},
		{"Unsupported type", struct{}{}, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convert.AnyToLINT(tc.input)

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
		expected iec.ULINT
		hasError bool
	}{
		{"SINT positive", iec.SINT(10), iec.ULINT(10), false},
		{"SINT negative", iec.SINT(-10), iec.ULINT(0xFFFFFFFFFFFFFFF6), false}, // two's complement
		{"INT", iec.INT(30000), iec.ULINT(30000), false},
		{"DINT", iec.DINT(-50000), iec.ULINT(0xFFFFFFFFFFFF3CB0), false},
		{"LINT", iec.LINT(123456789), iec.ULINT(123456789), false},
		{"USINT", iec.USINT(255), iec.ULINT(255), false},
		{"UINT", iec.UINT(100), iec.ULINT(100), false},
		{"UDINT", iec.UDINT(4000000000), iec.ULINT(4000000000), false},
		{"ULINT", iec.ULINT(999999999999), iec.ULINT(999999999999), false},
		{"REAL positive truncates", iec.REAL(123.75), iec.ULINT(123), false},
		{"LREAL negative returns error", iec.LREAL(-456.99), 0, true},
		{"BOOL true", iec.BOOL(true), iec.ULINT(1), false},
		{"BOOL false", iec.BOOL(false), iec.ULINT(0), false},
		{"BYTE", iec.BYTE(0xAB), iec.ULINT(171), false},
		{"WORD", iec.WORD(0xABCD), iec.ULINT(43981), false},
		{"DWORD", iec.DWORD(0xABCDEF01), iec.ULINT(2882400001), false},
		{"LWORD", iec.LWORD(0x1234567890ABCDEF), iec.ULINT(1311768467294899695), false},
		{"STRING int", iec.STRING("123"), iec.ULINT(123), false},
		{"STRING hex", iec.STRING("0xFF"), iec.ULINT(255), false},
		{"STRING invalid", iec.STRING("abc"), 0, true},
		{"TIME", iec.TIME(3 * time.Second), iec.ULINT(3000), false},
		{"DT", iec.DT(time.UnixMilli(1678886400000)), iec.ULINT(1678886400000), false},
		{"DATE", iec.DATE(time.UnixMilli(1678886400000)), iec.ULINT(1678886400000), false},
		{"TOD", iec.TOD(time.Date(0, 0, 0, 1, 2, 3, 0, time.UTC)), iec.ULINT(3723000), false},
		{"Unsupported type", struct{}{}, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convert.AnyToULINT(tc.input)

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
		result, err := convert.AnyToREAL(iec.LINT(123))
		if err != nil {
			t.Fatalf("AnyToREAL failed with error: %v", err)
		}
		if result != iec.REAL(123.0) {
			t.Errorf("AnyToREAL(LINT(123)) = %f; want 123.0", result)
		}
	})

	t.Run("Invalid conversion", func(t *testing.T) {
		_, err := convert.AnyToREAL(iec.STRING("abc"))
		if err == nil {
			t.Error("AnyToREAL(STRING(\"abc\")) should have returned an error")
		}
	})
}

func TestConvertTo(t *testing.T) {
	t.Run("LREAL_TO_DINT", func(t *testing.T) {
		input := iec.LREAL(123.7)
		expected := iec.DINT(123) // Note: LREAL_TO_DINT conversion truncates
		result, err := convert.ConvertTo[iec.DINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[DINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[DINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("SINT_TO_LREAL", func(t *testing.T) {
		input := iec.SINT(-50)
		expected := iec.LREAL(-50.0)
		result, err := convert.ConvertTo[iec.LREAL](input)
		if err != nil {
			t.Fatalf("ConvertTo[LREAL] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[LREAL](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("STRING_TO_LINT", func(t *testing.T) {
		input := iec.STRING("12345")
		expected := iec.LINT(12345)
		result, err := convert.ConvertTo[iec.LINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[LINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[LINT](%q) = %v; want %v", input, result, expected)
		}
	})

	t.Run("TIME_TO_STRING", func(t *testing.T) {
		input := iec.TIME(5 * time.Second) // Convert to seconds
		expected := iec.STRING("T#5s")
		result := conversion.TIME_TO_STRING(input)
		if result != expected {
			t.Errorf("ConvertTo[STRING](%v) = %q; want %q", input, result, expected)
		}
	})

	t.Run("UINT_TO_UDINT", func(t *testing.T) {
		input := iec.UINT(65000)
		expected := iec.UDINT(65000)
		result, err := convert.ConvertTo[iec.UDINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[UDINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[UDINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_BYTE", func(t *testing.T) {
		input := iec.LINT(255)
		expected := iec.BYTE(255)
		result, err := convert.ConvertTo[iec.BYTE](input)
		if err != nil {
			t.Fatalf("ConvertTo[BYTE] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[BYTE](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("INT_TO_TIME", func(t *testing.T) {
		input := iec.INT(5000)
		expected := iec.TIME(5 * time.Second)
		result, err := convert.ConvertTo[iec.TIME](input)
		if err != nil {
			t.Fatalf("ConvertTo[TIME] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[TIME](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_SINT", func(t *testing.T) {
		input := iec.LINT(120)
		expected := iec.SINT(120)
		result, err := convert.ConvertTo[iec.SINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[SINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[SINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("DINT_TO_UINT", func(t *testing.T) {
		input := iec.DINT(65000)
		expected := iec.UINT(65000)
		result, err := convert.ConvertTo[iec.UINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[UINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[UINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("REAL_TO_ULINT", func(t *testing.T) {
		input := iec.REAL(1234567.8)
		expected := iec.ULINT(1234567)
		result, err := convert.ConvertTo[iec.ULINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[ULINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[ULINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("INT_TO_BOOL", func(t *testing.T) {
		input := iec.INT(1)
		expected := iec.BOOL(true)
		result, err := convert.ConvertTo[iec.BOOL](input)
		if err != nil {
			t.Fatalf("ConvertTo[BOOL] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[BOOL](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("STRING_TO_BOOL", func(t *testing.T) {
		input := iec.STRING("true")
		expected := iec.BOOL(true)
		result, err := convert.ConvertTo[iec.BOOL](input)
		if err != nil {
			t.Fatalf("ConvertTo[BOOL] from STRING failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[BOOL](%q) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_WORD", func(t *testing.T) {
		input := iec.LINT(0xABCD)
		expected := iec.WORD(0xABCD)
		result, err := convert.ConvertTo[iec.WORD](input)
		if err != nil {
			t.Fatalf("ConvertTo[WORD] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[WORD](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_DATE", func(t *testing.T) {
		ms := int64(1678886400000)
		input := iec.LINT(ms)
		expected := iec.DATE(time.UnixMilli(ms))
		result, err := convert.ConvertTo[iec.DATE](input)
		if err != nil {
			t.Fatalf("ConvertTo[DATE] failed: %v", err)
		}
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ConvertTo[DATE](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LREAL_TO_INT", func(t *testing.T) {
		input := iec.LREAL(32767.8)
		expected := iec.INT(32767) // Clamped
		result, err := convert.ConvertTo[iec.INT](input)
		if err != nil {
			t.Fatalf("ConvertTo[INT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[INT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("DINT_TO_USINT", func(t *testing.T) {
		input := iec.DINT(250)
		expected := iec.USINT(250)
		result, err := convert.ConvertTo[iec.USINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[USINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[USINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_DWORD", func(t *testing.T) {
		input := iec.LINT(0xABCDEF)
		expected := iec.DWORD(0xABCDEF)
		result, err := convert.ConvertTo[iec.DWORD](input)
		if err != nil {
			t.Fatalf("ConvertTo[DWORD] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[DWORD](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("UINT_TO_LWORD", func(t *testing.T) {
		input := iec.UINT(0xFFFF)
		expected := iec.LWORD(0xFFFF)
		result, err := convert.ConvertTo[iec.LWORD](input)
		if err != nil {
			t.Fatalf("ConvertTo[LWORD] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[LWORD](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_TOD", func(t *testing.T) {
		ms := int64(3723000) // 1h, 2m, 3s
		input := iec.LINT(ms)
		expected := iec.TOD(time.Time{}.Add(time.Duration(ms) * time.Millisecond))
		result, err := convert.ConvertTo[iec.TOD](input)
		if err != nil {
			t.Fatalf("ConvertTo[TOD] failed: %v", err)
		}
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ConvertTo[TOD](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_DT", func(t *testing.T) {
		ms := int64(1678886400000)
		input := iec.LINT(ms)
		expected := iec.DT(time.UnixMilli(ms))
		result, err := convert.ConvertTo[iec.DT](input)
		if err != nil {
			t.Fatalf("ConvertTo[DT] failed: %v", err)
		}
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ConvertTo[DT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("DINT_TO_REAL", func(t *testing.T) {
		input := iec.DINT(12345)
		expected := iec.REAL(12345.0)
		result, err := convert.ConvertTo[iec.REAL](input)
		if err != nil {
			t.Fatalf("ConvertTo[REAL] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[REAL](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("Invalid STRING_TO_LINT", func(t *testing.T) {
		input := iec.STRING("abc")
		_, err := convert.ConvertTo[iec.LINT](input)
		if err == nil {
			t.Errorf("Expected an error for invalid string conversion, but got nil")
		}
	})

	t.Run("Unsupported Target Type", func(t *testing.T) {
		// Use a type that is not handled in the ConvertTo switch statement.
		input := iec.LINT(123)
		_, err := convert.ConvertTo[struct{}](input)
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
		if !convert.IsPlcFloat(iec.REAL(1.0)) {
			t.Error("IsPlcFloat(REAL) should be true")
		}
		if !convert.IsPlcFloat(iec.LREAL(1.0)) {
			t.Error("IsPlcFloat(LREAL) should be true")
		}
		if convert.IsPlcFloat(iec.DINT(1)) {
			t.Error("IsPlcFloat(DINT) should be false")
		}
	})

	t.Run("IsPlcInt", func(t *testing.T) {
		if !convert.IsPlcInt(iec.SINT(1)) {
			t.Error("IsPlcInt(SINT) should be true")
		}
		if !convert.IsPlcInt(iec.UINT(1)) {
			t.Error("IsPlcInt(UINT) should be true")
		}
		if !convert.IsPlcInt(iec.BOOL(true)) {
			t.Error("IsPlcInt(BOOL) should be true")
		}
		if convert.IsPlcInt(iec.REAL(1.0)) {
			t.Error("IsPlcInt(REAL) should be false")
		}
	})

	t.Run("IsPlcTimeType", func(t *testing.T) {
		if !convert.IsPlcTimeType(iec.TIME(0)) {
			t.Error("IsPlcTimeType(TIME) should be true")
		}
		if !convert.IsPlcTimeType(iec.DT{}) {
			t.Error("IsPlcTimeType(DT) should be true")
		}
		if convert.IsPlcTimeType(iec.LINT(0)) {
			t.Error("IsPlcTimeType(LINT) should be false")
		}
	})
}
