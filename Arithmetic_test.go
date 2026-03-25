package royaljelly

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestADD(t *testing.T) {
	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectError bool
	}{
		{"Add LINTs", []interface{}{LINT(10), LINT(20), LINT(30)}, LINT(60), false},
		{"Add REALs", []interface{}{REAL(1.5), REAL(2.5)}, REAL(4.0), false},
		{"Add Mixed Int/Float", []interface{}{INT(5), REAL(2.5), LINT(10)}, LREAL(17.5), false},
		{"Add with BOOL", []interface{}{BOOL(true), INT(5), BOOL(false)}, INT(6), false},
		{"Add with STRING", []interface{}{STRING("10"), INT(5)}, INT(15), false},
		{"Add with STRING float", []interface{}{STRING("10.5"), REAL(5)}, REAL(15.5), false},
		{"Add with invalid STRING", []interface{}{STRING("abc"), INT(5)}, nil, true},
		{"Add empty", []interface{}{}, 0, false},
		{"Add single", []interface{}{DINT(42)}, DINT(42), false},
		{"Add TIME", []interface{}{TIME(time.Second), TIME(time.Minute)}, TIME(61 * time.Second), false},
		{"Add TOD+TIME", []interface{}{TOD(time.Date(0, 0, 0, 10, 30, 0, 0, time.UTC)), TIME(time.Hour)}, TOD(time.Date(0, 0, 0, 11, 30, 0, 0, time.UTC)), false},
		{"Add DT+TIME", []interface{}{DT(time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)), TIME(24 * time.Hour)}, DT(time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC)), false},
		{"Add multiple TIME", []interface{}{TIME(time.Second), TIME(time.Minute), TIME(time.Hour)}, TIME(3661 * time.Second), false},
		{"Add TOD with multiple TIME", []interface{}{TOD(time.Date(0, 0, 0, 8, 0, 0, 0, time.UTC)), TIME(15 * time.Minute), TIME(30 * time.Second)}, TOD(time.Date(0, 0, 0, 8, 15, 30, 0, time.UTC)), false},
		{"Add DT with multiple TIME", []interface{}{DT(time.Date(2024, 5, 20, 10, 0, 0, 0, time.UTC)), TIME(2 * time.Hour), TIME(45 * time.Minute)}, DT(time.Date(2024, 5, 20, 12, 45, 0, 0, time.UTC)), false},
		{"Add TOD with invalid type", []interface{}{TOD(time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC)), INT(5)}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ADD(tc.inputs)

			if tc.expectError {
				if err == nil {
					t.Errorf("ADD(%v) did not return an error; expected error", tc.inputs)
				}
			} else {
				if err != nil {
					t.Errorf("ADD(%v) returned an unexpected error: %v", tc.inputs, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("ADD(%v) = %v; want %v", tc.inputs, result, tc.expected)
				}
			}
		})
	}
}

func TestSUB(t *testing.T) {
	// Base times for testing
	date1 := DATE(time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC))
	date2 := DATE(time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC))
	tod1 := TOD(time.Date(0, 0, 0, 14, 30, 0, 0, time.UTC))
	tod2 := TOD(time.Date(0, 0, 0, 10, 15, 0, 0, time.UTC))
	dt1 := DT(time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC))
	dt2 := DT(time.Date(2024, 3, 15, 10, 15, 0, 0, time.UTC))

	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectError bool
	}{
		{"Sub LINTs", []interface{}{LINT(100), LINT(20), LINT(30)}, LINT(50), false},
		{"Sub REALs", []interface{}{REAL(10.5), REAL(2.5)}, REAL(8.0), false},
		{"Sub Mixed Int/Float", []interface{}{LINT(20), REAL(2.5)}, LREAL(17.5), false},
		{"Sub with invalid STRING", []interface{}{STRING("abc"), INT(5)}, nil, true},
		{"Sub empty", []interface{}{}, 0, false},
		{"Sub single", []interface{}{DINT(42)}, DINT(42), false},
		{"TIME - TIME", []interface{}{TIME(time.Hour), TIME(15 * time.Minute)}, TIME(45 * time.Minute), false},
		{"DATE - DATE", []interface{}{date1, date2}, TIME(5 * 24 * time.Hour), false},
		{"TOD - TIME", []interface{}{tod1, TIME(time.Hour)}, TOD(time.Date(0, 0, 0, 13, 30, 0, 0, time.UTC)), false},
		{"TOD - TOD", []interface{}{tod1, tod2}, TIME(4*time.Hour + 15*time.Minute), false},
		{"DT - TIME", []interface{}{dt1, TIME(time.Hour)}, DT(time.Date(2024, 3, 15, 13, 30, 0, 0, time.UTC)), false},
		{"DT - DT", []interface{}{dt1, dt2}, TIME(4*time.Hour + 15*time.Minute), false},
		//{"TOD - invalid type", []interface{}{tod1, INT(5)}, nil, true,},
		//{"DT - invalid type", []interface{}{dt1, BOOL(true)}, nil, true,},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SUB(tc.inputs)

			if tc.expectError {
				if err == nil {
					t.Errorf("SUB(%v) did not return an error; expected error", tc.inputs)
				}
			} else {
				if err != nil {
					t.Errorf("SUB(%v) returned an unexpected error: %v", tc.inputs, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("SUB(%v) = %v (type %T); want %v (type %T)", tc.inputs, result, result, tc.expected, tc.expected)
				}
			}
		})
	}
}

func TestMUL(t *testing.T) {
	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectError bool
	}{
		{"Mul LINTs", []interface{}{LINT(2), LINT(3), LINT(4)}, LINT(24), false},
		{"Mul REALs", []interface{}{REAL(1.5), REAL(2.0)}, REAL(3.0), false},
		{"Mul Mixed Int/Float", []interface{}{INT(5), REAL(2.5)}, LREAL(12.5), false},
		{"Mul with zero", []interface{}{DINT(100), INT(0), LINT(50)}, LINT(0), false},
		{"Mul empty", []interface{}{}, LINT(1), false},
		{"Mul with invalid STRING", []interface{}{STRING("abc"), INT(5)}, nil, true},
		{"TIME * ANY_NUM", []interface{}{TIME(time.Second * 10), INT(3)}, TIME(time.Second * 30), false},
		{"TIME * REAL", []interface{}{TIME(60 * time.Second), REAL(1.5)}, TIME(90 * time.Second), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MUL(tc.inputs)

			if tc.expectError {
				if err == nil {
					t.Errorf("MUL(%v) did not return an error; expected error", tc.inputs)
				}
			} else {
				if err != nil {
					t.Errorf("MUL(%v) returned an unexpected error: %v", tc.inputs, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("MUL(%v) = %v; want %v", tc.inputs, result, tc.expected)
				}
			}
		})
	}
}

func TestDIV(t *testing.T) {
	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectError bool
	}{
		{"Div LINTs", []interface{}{LINT(100), LINT(10), LINT(2)}, LINT(5), false},
		{"Div REALs", []interface{}{REAL(20.0), REAL(4.0)}, REAL(5.0), false},
		{"Div Mixed Int/Float", []interface{}{LINT(25), REAL(2.5)}, LREAL(10.0), false},
		{"Div empty", []interface{}{}, LINT(0), false},
		{"Div single", []interface{}{DINT(42)}, DINT(42), false},
		{"TIME / ANY_NUM", []interface{}{TIME(time.Minute), INT(2)}, TIME(30 * time.Second), false},
		{"TIME / TIME", []interface{}{TIME(time.Minute), TIME(30 * time.Second)}, LREAL(2.0), false},
		{"TIME / TIME by zero", []interface{}{TIME(time.Minute), TIME(0)}, nil, true},
		{"TIME / REAL", []interface{}{TIME(time.Minute), REAL(1.5)}, TIME(40 * time.Second), false},
		{"Integer Div by Zero", []interface{}{LINT(100), LINT(0)}, nil, true},
		{"TIME Div by Zero", []interface{}{TIME(time.Minute), INT(0)}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := DIV(tc.inputs)

			if tc.expectError {
				if err == nil {
					t.Errorf("DIV(%v) did not return an error; expected error", tc.inputs)
				}
			} else {
				if err != nil {
					t.Errorf("DIV(%v) returned an unexpected error: %v", tc.inputs, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("DIV(%v) = %v; want %v", tc.inputs, result, tc.expected)
				}
			}
		})
	}
}

func TestMOD(t *testing.T) {
	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectError bool
	}{
		{"Mod LINTs", []interface{}{LINT(10), LINT(3)}, LINT(1), false},
		{"Mod chain", []interface{}{LINT(25), LINT(12), LINT(2)}, LINT(1), false},
		{"Mod with negative", []interface{}{LINT(-10), LINT(3)}, LINT(-1), false},
		{"Mod by zero", []interface{}{LINT(10), LINT(0)}, nil, true},
		{"Mod with float", []interface{}{LINT(10), REAL(3.0)}, nil, true},
		{"Mod not enough args", []interface{}{LINT(10)}, nil, true},
		{"Mod with invalid string", []interface{}{STRING("abc"), LINT(10)}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MOD(tc.inputs)

			if tc.expectError {
				if err == nil {
					t.Errorf("MOD(%v) did not return an error; expected error", tc.inputs)
				}
			} else {
				if err != nil {
					t.Errorf("MOD(%v) returned an unexpected error: %v", tc.inputs, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("MOD(%v) = %v; want %v", tc.inputs, result, tc.expected)
				}
			}
		})
	}
}

func TestMOVE(t *testing.T) {
	testCases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{"Move LINT", LINT(123), LINT(123)},
		{"Move REAL", REAL(45.6), REAL(45.6)},
		{"Move STRING", STRING("hello"), STRING("hello")},
		{"Move BOOL", BOOL(true), BOOL(true)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MOVE(tc.input)
			if result != tc.expected {
				t.Errorf("MOVE(%v) = %v; want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestTimeArithmeticFunctions(t *testing.T) {
	t.Run("ADD_TIME", func(t *testing.T) {
		in1 := TIME(time.Hour)
		in2 := TIME(30 * time.Minute)
		expected := TIME(90 * time.Minute)
		result, err := ADD([]interface{}{in1, in2})
		if err != nil {
			t.Fatalf("ADD_TIME returned an unexpected error: %v", err)
		}
		if result != expected {
			t.Errorf("ADD_TIME(%v, %v) = %v; want %v", in1, in2, result, expected)
		}
	})

	t.Run("ADD_TOD", func(t *testing.T) {
		in1 := TOD(time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC))
		in2 := TIME(15 * time.Minute)
		expected := TOD(time.Date(0, 0, 0, 10, 15, 0, 0, time.UTC))
		result, err := ADD([]interface{}{in1, in2})
		if err != nil {
			t.Fatalf("ADD(TOD, TIME) returned an unexpected error: %v", err)
		}
		if !time.Time(result.(TOD)).Equal(time.Time(expected)) {
			t.Errorf("ADD(TOD, TIME) = %v; want %v", result, expected)
		}
	})

	t.Run("ADD_DT", func(t *testing.T) {
		in1 := DT(time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC))
		in2 := TIME(3 * time.Hour)
		expected := DT(time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC))
		result, err := ADD([]interface{}{in1, in2})
		if err != nil {
			t.Fatalf("ADD(DT, TIME) returned an unexpected error: %v", err)
		}
		if !time.Time(result.(DT)).Equal(time.Time(expected)) {
			t.Errorf("ADD(DT, TIME) = %v; want %v", result, expected)
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
		result1, err1 := SUB([]interface{}{tod1, time1})
		if err1 != nil {
			t.Fatalf("SUB(TOD, TIME) returned an unexpected error: %v", err1)
		}
		if !time.Time(result1.(TOD)).Equal(time.Time(expected1)) {
			t.Errorf("SUB(TOD, TIME) = %v; want %v", result1, expected1)
		}
		result1a, err1a := SUB_TOD(tod1, time1)
		if err1a != nil {
			t.Fatalf("SUB_TOD(TOD, TIME) returned an unexpected error: %v", err1a)
		}
		if !time.Time(result1a.(TOD)).Equal(time.Time(expected1)) {
			t.Errorf("SUB_TOD(TOD, TIME) = %v; want %v", result1a, expected1)
		}

		// TOD - TOD -> TIME
		expected2 := TIME(time.Hour)
		result2, err2 := SUB([]interface{}{tod1, tod2})
		if err2 != nil {
			t.Fatalf("SUB(TOD, TOD) returned an unexpected error: %v", err2)
		}
		if result2.(TIME) != expected2 {
			t.Errorf("SUB(TOD, TOD) = %v; want %v", result2, expected2)
		}

		result2a, err2a := SUB_TOD(tod1, tod2)
		if err2a != nil {
			t.Fatalf("SUB_TOD(TOD, TOD) returned an unexpected error: %v", err2a)
		}
		if result2.(TIME) != expected2 {
			t.Errorf("SUB_TOD(TOD, TOD) = %v; want %v", result2a, expected2)
		}
	})

	t.Run("SUB_DT", func(t *testing.T) {
		dt1 := DT(time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC))
		time1 := TIME(24 * time.Hour)
		dt2 := DT(time.Date(2024, 3, 14, 12, 0, 0, 0, time.UTC))

		// DT - TIME -> DT
		expected1 := DT(time.Date(2024, 3, 14, 12, 0, 0, 0, time.UTC))
		result1, err1 := SUB([]interface{}{dt1, time1})
		if err1 != nil {
			t.Fatalf("SUB(DT, TIME) returned an unexpected error: %v", err1)
		}
		if !time.Time(result1.(DT)).Equal(time.Time(expected1)) {
			t.Errorf("SUB_DT(DT, TIME) = %v; want %v", result1, expected1)
		}

		// DT - DT -> TIME
		expected2 := TIME(24 * time.Hour)
		result2, err2 := SUB([]interface{}{dt1, dt2})
		if err2 != nil {
			t.Fatalf("SUB(DT, DT) returned an unexpected error: %v", err2)
		}
		if result2.(TIME) != expected2 {
			t.Errorf("SUB_DT(DT, DT) = %v; want %v", result2, expected2)
		}

		// Test error on invalid second argument
		t.Run("SUB_DT with invalid type", func(t *testing.T) {
			_, err := SUB_DT(dt1, LINT(123))
			if err == nil {
				t.Error("SUB_DT with invalid type did not return an error")
			}
		})
	})

	t.Run("SUB_DT with invalid first argument", func(t *testing.T) {
		_, err := SUB_DT(LINT(123), TIME(time.Second))
		if err == nil {
			t.Error("SUB_DT with invalid first argument did not return an error")
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
			result, err := anyToLREAL(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("anyToLREAL(%v) expected an error, but got none", tc.input)
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
			result, err := anyToLINT(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("anyToLINT(%v) expected an error, but got none", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("anyToLINT(%v) returned an error: %v", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("anyToLINT(%v) = %d; want %d", tc.input, result, tc.expected)
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
			result, err := anyToULINT(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("anyToULINT(%v) expected an error, but got none", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("anyToULINT(%v) returned an error: %v", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("anyToULINT(%v) = %d; want %d", tc.input, result, tc.expected)
				}
			}
		})
	}
}

func TestAnyToREAL(t *testing.T) {
	t.Run("Valid conversion", func(t *testing.T) {
		result, err := anytoREAL(LINT(123))
		if err != nil {
			t.Fatalf("anytoREAL failed with error: %v", err)
		}
		if result != REAL(123.0) {
			t.Errorf("anytoREAL(LINT(123)) = %f; want 123.0", result)
		}
	})

	t.Run("Invalid conversion", func(t *testing.T) {
		_, err := anytoREAL(STRING("abc"))
		if err == nil {
			t.Error("anytoREAL(STRING(\"abc\")) should have returned an error")
		}
	})
}

func TestConvertToTargetType(t *testing.T) {
	t.Run("Successful Conversions", func(t *testing.T) {
		testCases := []struct {
			name        string
			accumulator interface{}
			targetType  reflect.Type
			expected    interface{}
		}{
			// LINT to other types
			{"LINT to SINT", LINT(127), reflect.TypeOf(SINT(0)), SINT(127)},
			{"LINT to INT", LINT(32767), reflect.TypeOf(INT(0)), INT(32767)},
			{"LINT to DINT", LINT(123456), reflect.TypeOf(DINT(0)), DINT(123456)},
			{"LINT to LINT", LINT(98765), reflect.TypeOf(LINT(0)), LINT(98765)},
			{"LINT to USINT", LINT(255), reflect.TypeOf(USINT(0)), USINT(255)},
			{"LINT to UINT", LINT(65535), reflect.TypeOf(UINT(0)), UINT(65535)},
			{"LINT to UDINT", LINT(123456), reflect.TypeOf(UDINT(0)), UDINT(123456)},
			{"LINT to ULINT", LINT(98765), reflect.TypeOf(ULINT(0)), ULINT(98765)},
			{"LINT to REAL", LINT(123), reflect.TypeOf(REAL(0)), REAL(123.0)},
			{"LINT to LREAL", LINT(456), reflect.TypeOf(LREAL(0)), LREAL(456.0)},
			{"LINT to BOOL true", LINT(1), reflect.TypeOf(BOOL(false)), BOOL(true)},
			{"LINT to BOOL false", LINT(0), reflect.TypeOf(BOOL(false)), BOOL(false)},
			{"LINT to BYTE", LINT(0xAB), reflect.TypeOf(BYTE(0)), BYTE(0xAB)},
			{"LINT to WORD", LINT(0xABCD), reflect.TypeOf(WORD(0)), WORD(0xABCD)},
			{"LINT to DWORD", LINT(0xABCDEF), reflect.TypeOf(DWORD(0)), DWORD(0xABCDEF)},
			{"LINT to LWORD", LINT(0x12345678), reflect.TypeOf(LWORD(0)), LWORD(0x12345678)},
			{"LINT to TIME", LINT(5000), reflect.TypeOf(TIME(0)), TIME(5 * time.Second)},
			{"LINT to DATE", LINT(1678886400000), reflect.TypeOf(DATE{}), DATE(time.UnixMilli(1678886400000))},
			{"LINT to TOD", LINT(3723000), reflect.TypeOf(TOD{}), TOD(time.Time{}.Add(3723000 * time.Millisecond))},
			{"LINT to DT", LINT(1678886400000), reflect.TypeOf(DT{}), DT(time.UnixMilli(1678886400000))},
			{"LINT to STRING", LINT(12345), reflect.TypeOf(STRING("")), STRING("12345")},

			// LREAL to other types
			{"LREAL to SINT", LREAL(12.7), reflect.TypeOf(SINT(0)), SINT(12)},
			{"LREAL to INT", LREAL(327.6), reflect.TypeOf(INT(0)), INT(327)},
			{"LREAL to DINT", LREAL(12345.6), reflect.TypeOf(DINT(0)), DINT(12345)},
			{"LREAL to LINT", LREAL(9876.5), reflect.TypeOf(LINT(0)), LINT(9876)},
			{"LREAL to USINT", LREAL(25.5), reflect.TypeOf(USINT(0)), USINT(25)},
			{"LREAL to UINT", LREAL(655.3), reflect.TypeOf(UINT(0)), UINT(655)},
			{"LREAL to UDINT", LREAL(12345.6), reflect.TypeOf(UDINT(0)), UDINT(12345)},
			{"LREAL to ULINT", LREAL(9876.5), reflect.TypeOf(ULINT(0)), ULINT(9876)},
			{"LREAL to REAL", LREAL(123.456), reflect.TypeOf(REAL(0)), REAL(123.456)},
			{"LREAL to LREAL", LREAL(456.789), reflect.TypeOf(LREAL(0)), LREAL(456.789)},
			{"LREAL to BOOL true", LREAL(0.1), reflect.TypeOf(BOOL(false)), BOOL(true)},
			{"LREAL to BOOL false", LREAL(0.0), reflect.TypeOf(BOOL(false)), BOOL(false)},
			{"LREAL to BYTE", LREAL(171.2), reflect.TypeOf(BYTE(0)), BYTE(171)},
			{"LREAL to WORD", LREAL(43981.9), reflect.TypeOf(WORD(0)), WORD(43981)},
			{"LREAL to DWORD", LREAL(11259375.1), reflect.TypeOf(DWORD(0)), DWORD(11259375)},
			{"LREAL to LWORD", LREAL(1.3e18), reflect.TypeOf(LWORD(0)), LWORD(1300000013958512640)},
			{"LREAL to TIME", LREAL(5000.5), reflect.TypeOf(TIME(0)), TIME(5 * time.Second)},
			{"LREAL to STRING", LREAL(123.45), reflect.TypeOf(STRING("")), STRING("123.45")},

			// ULINT to other types
			{"ULINT to SINT", ULINT(127), reflect.TypeOf(SINT(0)), SINT(127)},
			{"ULINT to INT", ULINT(32767), reflect.TypeOf(INT(0)), INT(32767)},
			{"ULINT to DINT", ULINT(123456), reflect.TypeOf(DINT(0)), DINT(123456)},
			{"ULINT to LINT", ULINT(98765), reflect.TypeOf(LINT(0)), LINT(98765)},
			{"ULINT to USINT", ULINT(255), reflect.TypeOf(USINT(0)), USINT(255)},
			{"ULINT to UINT", ULINT(65535), reflect.TypeOf(UINT(0)), UINT(65535)},
			{"ULINT to UDINT", ULINT(123456), reflect.TypeOf(UDINT(0)), UDINT(123456)},
			{"ULINT to ULINT", ULINT(98765), reflect.TypeOf(ULINT(0)), ULINT(98765)},
			{"ULINT to REAL", ULINT(123), reflect.TypeOf(REAL(0)), REAL(123.0)},
			{"ULINT to LREAL", ULINT(456), reflect.TypeOf(LREAL(0)), LREAL(456.0)},
			{"ULINT to BOOL true", ULINT(1), reflect.TypeOf(BOOL(false)), BOOL(true)},
			{"ULINT to BOOL false", ULINT(0), reflect.TypeOf(BOOL(false)), BOOL(false)},
			{"ULINT to BYTE", ULINT(0xAB), reflect.TypeOf(BYTE(0)), BYTE(0xAB)},
			{"ULINT to WORD", ULINT(0xABCD), reflect.TypeOf(WORD(0)), WORD(0xABCD)},
			{"ULINT to DWORD", ULINT(0xABCDEF), reflect.TypeOf(DWORD(0)), DWORD(0xABCDEF)},
			{"ULINT to LWORD", ULINT(0x12345678), reflect.TypeOf(LWORD(0)), LWORD(0x12345678)},
			{"ULINT to TIME", ULINT(5000), reflect.TypeOf(TIME(0)), TIME(5 * time.Second)},
			{"ULINT to STRING", ULINT(12345), reflect.TypeOf(STRING("")), STRING("12345")},

			// Special TIME accumulator
			{"TIME to TIME", TIME(time.Second), reflect.TypeOf(TIME(0)), TIME(time.Second)},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, err := convertToTargetType(tc.accumulator, tc.targetType)
				if err != nil {
					t.Fatalf("convertToTargetType failed with error: %v", err)
				}

				// Using Sprintf for consistent comparison, especially for time types
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("convertToTargetType(%v, %v) = %v (type %T); want %v (type %T)",
						tc.accumulator, tc.targetType.Name(), result, result, tc.expected, tc.expected)
				}
			})
		}
	})

	t.Run("LREAL to Inf to INT should panic", func(t *testing.T) {
		inf := LREAL(math.Inf(1))
		// Directly call convertToTargetType to test the error path
		_, err := convertToTargetType(inf, reflect.TypeOf(INT(0)))
		if err == nil {
			t.Fatalf("Expected an error when converting Inf to INT, but got nil")
		}
	})

	t.Run("Unhandled accumulator type", func(t *testing.T) {
		// defer func() {
		// 	if r := recover(); r == nil {
		// 		t.Error("Expected panic for unhandled accumulator type")
		// 	}
		// }()
		// We can't call convertToTargetType directly, but we can trigger the error path.
		// This is a bit tricky as the public functions wrap it.
		// We'll test the error path inside convertToTargetType itself.
		_, err := convertToTargetType("a string accumulator", reflect.TypeOf(INT(0)))
		if err == nil {
			t.Error("Expected error for unhandled accumulator type")
		}
	})

	t.Run("Unhandled target type", func(t *testing.T) {
		_, err := convertToTargetType(LINT(123), reflect.TypeOf(struct{}{}))
		if err == nil {
			t.Error("Expected error for unhandled target type")
		}
	})
}

func TestTimeSpecificFunctions(t *testing.T) {
	t.Run("SUB_TOD with invalid type", func(t *testing.T) {
		_, err := SUB_TOD(TOD(time.Now()), LINT(5))
		if err == nil {
			t.Error("SUB_TOD with invalid type did not return an error")
		}
	})
	t.Run("SUB_DT with invalid type", func(t *testing.T) {
		_, err := SUB_DT(DT(time.Now()), LINT(5))
		if err == nil {
			t.Error("SUB_DT with invalid type did not return an error")
		}
	})
}

func TestADD_T(t *testing.T) {
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
}
