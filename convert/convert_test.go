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

package convert

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

func TestGetTypeName(t *testing.T) {
	testCases := []struct {
		input    any
		expected string
	}{
		{iec.BOOL(true), "BOOL"},
		{iec.SINT(1), "SINT"},
		{iec.INT(1), "INT"},
		{iec.DINT(1), "DINT"},
		{iec.LINT(1), "LINT"},
		{iec.USINT(1), "USINT"},
		{iec.UINT(1), "UINT"},
		{iec.UDINT(1), "UDINT"},
		{iec.ULINT(1), "ULINT"},
		{iec.REAL(1.0), "REAL"},
		{iec.LREAL(1.0), "LREAL"},
		{iec.STRING("s"), "STRING"},
		{iec.TIME(0), "Time/Date Type"},
		{iec.DATE{}, "Time/Date Type"},
		{iec.TOD{}, "Time/Date Type"},
		{iec.DT{}, "Time/Date Type"},
		{iec.BYTE(1), "Bit-String Type"},
		{iec.WORD(1), "Bit-String Type"},
		{iec.DWORD(1), "Bit-String Type"},
		{iec.LWORD(1), "Bit-String Type"},
		{struct{}{}, "unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			if name := GetTypeName(tc.input); name != tc.expected {
				t.Errorf("GetTypeName(%T) = %q; want %q", tc.input, name, tc.expected)
			}
		})
	}
}

func TestTypeCheckFunctions(t *testing.T) {
	t.Run("IsPlcFloat", func(t *testing.T) {
		if !IsPlcFloat(iec.REAL(1.0)) {
			t.Error("IsPlcFloat(REAL) should be true")
		}
		if !IsPlcFloat(iec.LREAL(1.0)) {
			t.Error("IsPlcFloat(LREAL) should be true")
		}
		if IsPlcFloat(iec.DINT(1)) {
			t.Error("IsPlcFloat(DINT) should be false")
		}
	})

	t.Run("IsPlcInt", func(t *testing.T) {
		if !IsPlcInt(iec.SINT(1)) {
			t.Error("IsPlcInt(SINT) should be true")
		}
		if !IsPlcInt(iec.UINT(1)) {
			t.Error("IsPlcInt(UINT) should be true")
		}
		if !IsPlcInt(iec.BOOL(true)) {
			t.Error("IsPlcInt(BOOL) should be true")
		}
		if IsPlcInt(iec.REAL(1.0)) {
			t.Error("IsPlcInt(REAL) should be false")
		}
	})

	t.Run("IsPlcTimeType", func(t *testing.T) {
		if !IsPlcTimeType(iec.TIME(0)) {
			t.Error("IsPlcTimeType(TIME) should be true")
		}
		if !IsPlcTimeType(iec.DT{}) {
			t.Error("IsPlcTimeType(DT) should be true")
		}
		if IsPlcTimeType(iec.LINT(0)) {
			t.Error("IsPlcTimeType(LINT) should be false")
		}
	})
}

func TestAnyToLREAL(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected iec.LREAL
		hasError bool
	}{
		{"SINT", iec.SINT(-10), iec.LREAL(-10), false},
		{"UINT", iec.UINT(100), iec.LREAL(100), false},
		{"REAL", iec.REAL(123.45), iec.LREAL(123.44999694824219), false},
		{"BOOL true", iec.BOOL(true), iec.LREAL(1.0), false},
		{"STRING float", iec.STRING("123.45"), iec.LREAL(123.45), false},
		{"STRING invalid", iec.STRING("abc"), 0, true},
		{"Unsupported type", struct{}{}, 0, true},
		{"TIME", iec.TIME(2 * time.Second), iec.LREAL(2000), false}, // TIME is in nanoseconds, converted to milliseconds
		{"LREAL", iec.LREAL(987.65), iec.LREAL(987.65), false},
		{"BYTE", iec.BYTE(0xAB), iec.LREAL(171), false},
		{"DATE", iec.DATE(time.UnixMilli(1678886400000)), iec.LREAL(1678886400000), false},
		{"TOD", iec.TOD(time.Date(0, 1, 1, 14, 30, 15, 0, time.UTC)), iec.LREAL(52215000), false},
		{"DT", iec.DT(time.UnixMilli(1678886400000)), iec.LREAL(1678886400000), false},
		{"float32", float32(1.23), iec.LREAL(1.2300000190734863), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			const tolerance = 1e-9
			result, err := AnyToLREAL(tc.input)

			if tc.hasError {
				if err == nil {
					t.Errorf("AnyToLREAL(%v) expected an error, but got none", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("AnyToLREAL(%v) returned an error: %v", tc.input, err)
				}
				if diff := iec.LREAL(result - tc.expected); diff < -tolerance || diff > tolerance {
					t.Errorf("AnyToLREAL(%v) = %g; want %g", tc.input, result, tc.expected)
				}
			}
		})
	}
}

func TestAnyToLINT(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected iec.LINT
		hasError bool
	}{
		{"SINT", iec.SINT(-10), iec.LINT(-10), false},
		{"UINT", iec.UINT(100), iec.LINT(100), false},
		{"ULINT overflow", iec.ULINT(math.MaxInt64 + 1), iec.MAXLINT, true},
		{"REAL truncates", iec.REAL(123.75), iec.LINT(123), false},
		{"BOOL true", iec.BOOL(true), iec.LINT(1), false},
		{"BYTE", iec.BYTE(0xFE), iec.LINT(254), false},
		{"STRING int", iec.STRING("123"), iec.LINT(123), false},
		{"STRING invalid", iec.STRING("abc"), 0, true},
		{"TIME", iec.TIME(3 * time.Second), iec.LINT(3000), false},
		{"DATE", iec.DATE(time.UnixMilli(1678886400000)), iec.LINT(1678886400000), false},
		{"TOD", iec.TOD(time.Date(0, 0, 0, 1, 2, 3, 0, time.UTC)), iec.LINT(3723000), false},
		{"DT", iec.DT(time.UnixMilli(1678886400000)), iec.LINT(1678886400000), false},
		{"UDINT valid", iec.UDINT(iec.MAXUDINT), iec.LINT(iec.MAXUDINT), false},
		{"REAL overflow >", iec.REAL(float32(iec.MAXLINT) * 2), iec.MAXLINT, true},
		{"REAL overflow <", iec.REAL(float32(iec.MINLINT) * 2), iec.MINLINT, true},
		{"LWORD", iec.LWORD(123), iec.LINT(123), false},
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
		input    any
		expected iec.ULINT
		hasError bool
	}{
		{"SINT positive", iec.SINT(10), iec.ULINT(10), false},
		{"SINT negative", iec.SINT(-10), iec.ULINT(0xFFFFFFFFFFFFFFF6), false}, // two's complement
		{"LREAL negative returns error", iec.LREAL(-456.99), 0, true},
		{"BOOL true", iec.BOOL(true), iec.ULINT(1), false},
		{"BYTE", iec.BYTE(0xAB), iec.ULINT(171), false},
		{"STRING hex", iec.STRING("0xFF"), iec.ULINT(255), false},
		{"STRING invalid", iec.STRING("abc"), 0, true},
		{"TIME", iec.TIME(3 * time.Second), iec.ULINT(3000), false},
		{"DATE", iec.DATE(time.UnixMilli(1678886400000)), iec.ULINT(1678886400000), false},
		{"TOD", iec.TOD(time.Date(0, 0, 0, 1, 2, 3, 0, time.UTC)), iec.ULINT(3723000), false},
		{"DT", iec.DT(time.UnixMilli(1678886400000)), iec.ULINT(1678886400000), false},
		{"REAL positive", iec.REAL(123.75), iec.ULINT(123), false},
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

func TestConversionError(t *testing.T) {
	t.Run("Error formatting", func(t *testing.T) {
		baseErr := fmt.Errorf("base error")
		err := &ConversionError{
			Value:    iec.LINT(123),
			FromType: "LINT",
			ToType:   "SINT",
			Reason:   "overflow",
			Err:      baseErr,
		}
		_ = err.Error() // Exercise the Error() method
		if !errors.Is(err, err.Unwrap()) {
			t.Error("Unwrap should return the inner error")
		}
	})
}

func TestAnyToREAL(t *testing.T) {
	t.Run("Valid LINT conversion", func(t *testing.T) {
		result, err := AnyToREAL(iec.LINT(123))
		if err != nil {
			t.Fatalf("AnyToREAL failed with error: %v", err)
		}
		if result != iec.REAL(123.0) {
			t.Errorf("AnyToREAL(LINT(123)) = %f; want 123.0", result)
		}
	})

	t.Run("Overflow from LREAL", func(t *testing.T) {
		_, err := AnyToREAL(iec.LREAL(math.MaxFloat64))
		if err == nil {
			t.Error("AnyToREAL from huge LREAL should have returned an error")
		}
	})

	t.Run("Underflow from LREAL", func(t *testing.T) {
		_, err := AnyToREAL(iec.LREAL(-math.MaxFloat64))
		if err == nil {
			t.Error("AnyToREAL from huge negative LREAL should have returned an error")
		}
	})

	t.Run("Invalid STRING conversion", func(t *testing.T) {
		_, err := AnyToREAL(iec.STRING("abc"))
		if err == nil {
			t.Error("AnyToREAL(STRING(\"abc\")) should have returned an error")
		}
	})
}

func TestAnyToBOOL(t *testing.T) {
	testCases := []struct {
		name     string
		input    any
		expected iec.BOOL
		hasError bool
	}{
		{"LINT non-zero", iec.LINT(123), true, false},
		{"LINT zero", iec.LINT(0), false, false},
		{"REAL non-zero", iec.REAL(0.1), true, false},
		{"REAL zero", iec.REAL(0.0), false, false},
		{"STRING true", iec.STRING("true"), true, false},
		{"STRING t", iec.STRING("t"), true, false},
		{"STRING 1", iec.STRING("1"), true, false},
		{"STRING false", iec.STRING("false"), false, true}, // AnyToLINT fails, and string is not "true"
		{"STRING empty", iec.STRING(""), false, true},
		{"Unsupported", struct{}{}, false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := AnyToBOOL(tc.input)
			if tc.hasError {
				if err == nil {
					t.Errorf("AnyToBOOL(%v) expected an error, but got none", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("AnyToBOOL(%v) returned an error: %v", tc.input, err)
				}
				if result != tc.expected {
					t.Errorf("AnyToBOOL(%v) = %v; want %v", tc.input, result, tc.expected)
				}
			}
		})
	}
}

func TestConvertTo(t *testing.T) {
	t.Run("LREAL_TO_DINT", func(t *testing.T) {
		input := iec.LREAL(123.7)
		expected := iec.DINT(123) // Note: LREAL_TO_DINT conversion truncates
		result, err := ConvertTo[iec.DINT](input)
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
		result, err := ConvertTo[iec.LREAL](input)
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
		result, err := ConvertTo[iec.LINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[LINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[LINT](%q) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LREAL_TO_INT_Overflow", func(t *testing.T) {
		input := iec.LREAL(iec.MAXINT + 1)
		_, err := ConvertTo[iec.INT](input)
		if err == nil {
			t.Error("Expected overflow error for LREAL to INT conversion, but got nil")
		}
	})

	t.Run("LREAL_TO_INT_Underflow", func(t *testing.T) {
		input := iec.LREAL(iec.MININT - 1)
		_, err := ConvertTo[iec.INT](input)
		if err == nil {
			t.Error("Expected underflow error for LREAL to INT conversion, but got nil")
		}
	})

	t.Run("LINT_TO_SINT_Overflow", func(t *testing.T) {
		input := iec.LINT(300) // MAXSINT is 127
		_, err := ConvertTo[iec.SINT](input)
		if err == nil {
			t.Error("Expected overflow error for LINT to SINT conversion, but got nil")
		}
		var convErr *ConversionError
		if !errors.As(err, &convErr) {
			t.Fatalf("Expected error of type *ConversionError, but got %T", err)
		}
		if convErr.Reason != "overflow" {
			t.Errorf("Expected reason 'overflow', but got %q", convErr.Reason)
		}
	})

	t.Run("LINT_TO_USINT_Overflow", func(t *testing.T) {
		input := iec.LINT(300) // MAXUSINT is 255
		_, err := ConvertTo[iec.USINT](input)
		if err == nil {
			t.Error("Expected overflow error for LINT to USINT conversion, but got nil")
		}
	})

	t.Run("LINT_TO_USINT_Negative", func(t *testing.T) {
		input := iec.LINT(-1)
		_, err := ConvertTo[iec.USINT](input)
		if err == nil {
			t.Error("Expected error for negative LINT to USINT conversion, but got nil")
		}
	})

	t.Run("UINT_TO_UDINT", func(t *testing.T) {
		input := iec.UINT(65000)
		expected := iec.UDINT(65000)
		result, err := ConvertTo[iec.UDINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[UDINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[UDINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("DINT_TO_UINT_Overflow", func(t *testing.T) {
		input := iec.DINT(iec.MAXUINT + 1)
		_, err := ConvertTo[iec.UINT](input)
		if err == nil {
			t.Error("Expected overflow error for DINT to UINT conversion, but got nil")
		}
	})

	t.Run("REAL_TO_ULINT", func(t *testing.T) {
		input := iec.REAL(1234567.8)
		expected := iec.ULINT(1234567)
		result, err := ConvertTo[iec.ULINT](input)
		if err != nil {
			t.Fatalf("ConvertTo[ULINT] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[ULINT](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_BYTE", func(t *testing.T) {
		input := iec.LINT(255)
		expected := iec.BYTE(255)
		result, err := ConvertTo[iec.BYTE](input)
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
		result, err := ConvertTo[iec.TIME](input)
		if err != nil {
			t.Fatalf("ConvertTo[TIME] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[TIME](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("INT_TO_BOOL", func(t *testing.T) {
		input := iec.INT(1)
		expected := iec.BOOL(true)
		result, err := ConvertTo[iec.BOOL](input)
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
		result, err := ConvertTo[iec.BOOL](input)
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
		result, err := ConvertTo[iec.WORD](input)
		if err != nil {
			t.Fatalf("ConvertTo[WORD] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[WORD](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_DWORD", func(t *testing.T) {
		input := iec.LINT(0xABCDEF)
		expected := iec.DWORD(0xABCDEF)
		result, err := ConvertTo[iec.DWORD](input)
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
		result, err := ConvertTo[iec.LWORD](input)
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
		result, err := ConvertTo[iec.TOD](input)
		if err != nil {
			t.Fatalf("ConvertTo[TOD] failed: %v", err)
		}
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ConvertTo[TOD](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_DATE", func(t *testing.T) {
		ms := int64(1678886400000)
		input := iec.LINT(ms)
		expected := iec.DATE(time.UnixMilli(ms))
		result, err := ConvertTo[iec.DATE](input)
		if err != nil {
			t.Fatalf("ConvertTo[DATE] failed: %v", err)
		}
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("ConvertTo[DATE](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("LINT_TO_DT", func(t *testing.T) {
		ms := int64(1678886400000)
		input := iec.LINT(ms)
		expected := iec.DT(time.UnixMilli(ms))
		result, err := ConvertTo[iec.DT](input)
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
		result, err := ConvertTo[iec.REAL](input)
		if err != nil {
			t.Fatalf("ConvertTo[REAL] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[REAL](%v) = %v; want %v", input, result, expected)
		}
	})

	t.Run("ANY_TO_STRING", func(t *testing.T) {
		input := iec.DINT(123)
		expected := iec.STRING("123")
		result, err := ConvertTo[iec.STRING](input)
		if err != nil {
			t.Fatalf("ConvertTo[STRING] failed: %v", err)
		}
		if result != expected {
			t.Errorf("ConvertTo[STRING](%v) = %q; want %q", input, result, expected)
		}
	})

	t.Run("Unsupported Target Type", func(t *testing.T) {
		input := iec.LINT(123)
		_, err := ConvertTo[struct{}](input)
		if err == nil {
			t.Error("Expected an error for unsupported target type, but got nil")
		}
	})
}

func TestSubFunctions(t *testing.T) {
	t.Run("SubByte", func(t *testing.T) {
		res, err := SubByte(iec.LINT(123))
		if err != nil || res != 123 {
			t.Errorf("SubByte failed, got %v, %v", res, err)
		}
		_, err = SubByte(iec.LINT(300))
		if err == nil {
			t.Error("SubByte should have failed on overflow")
		}
	})
	t.Run("SubWord", func(t *testing.T) {
		res, err := SubWord(iec.LINT(123))
		if err != nil || res != 123 {
			t.Errorf("SubWord failed, got %v, %v", res, err)
		}
	})
	t.Run("SubDword", func(t *testing.T) {
		res, err := SubDword(iec.LINT(123))
		if err != nil || res != 123 {
			t.Errorf("SubDword failed, got %v, %v", res, err)
		}
	})
	t.Run("SubLword", func(t *testing.T) {
		res, err := SubLword(iec.LINT(123))
		if err != nil || res != 123 {
			t.Errorf("SubLword failed, got %v, %v", res, err)
		}
	})
	t.Run("SubDt", func(t *testing.T) {
		_, err := SubDt(iec.STRING("abc"))
		if err == nil {
			t.Error("SubDt should have failed on invalid input")
		}
	})
	t.Run("SubDate", func(t *testing.T) {
		_, err := SubDate(iec.STRING("abc"))
		if err == nil {
			t.Error("SubDate should have failed on invalid input")
		}
	})
	t.Run("SubTod", func(t *testing.T) {
		_, err := SubTod(iec.STRING("abc"))
		if err == nil {
			t.Error("SubTod should have failed on invalid input")
		}
	})
	t.Run("SubBool", func(t *testing.T) {
		if SubBool(true) != 1 || SubBool(false) != 0 {
			t.Error("SubBool conversion is incorrect")
		}
	})
}

func TestClampFunctions(t *testing.T) {
	t.Run("ClampLINT", func(t *testing.T) {
		if ClampLINT(150, 0, 100) != 100 {
			t.Error("ClampLINT max failed")
		}
		if ClampLINT(-50, 0, 100) != 0 {
			t.Error("ClampLINT min failed")
		}
		if ClampLINT(50, 0, 100) != 50 {
			t.Error("ClampLINT within range failed")
		}
	})

	t.Run("ClampULINT", func(t *testing.T) {
		if ClampULINT(150, 100) != 100 {
			t.Error("ClampULINT max failed")
		}
		if ClampULINT(50, 100) != 50 {
			t.Error("ClampULINT within range failed")
		}
	})

	t.Run("ClampLREAL", func(t *testing.T) {
		if ClampLREAL(150.5, 0.0, 100.0) != 100.0 {
			t.Error("ClampLREAL max failed")
		}
		if ClampLREAL(-50.5, 0.0, 100.0) != 0.0 {
			t.Error("ClampLREAL min failed")
		}
		if ClampLREAL(50.5, 0.0, 100.0) != 50.5 {
			t.Error("ClampLREAL within range failed")
		}
	})
}

func TestAlmostEqual(t *testing.T) {
	if !AlmostEqual(1.0, 1.0+1e-10) {
		t.Error("AlmostEqual should be true for very close numbers")
	}
	if AlmostEqual(1.0, 1.0+1e-8) {
		t.Error("AlmostEqual should be false for numbers outside the threshold")
	}
}

func TestRoundAndClampLREAL(t *testing.T) {
	if RoundAndClampLREAL(50.7, 0, 100) != 51.0 {
		t.Error("RoundAndClampLREAL rounding up failed")
	}
	if RoundAndClampLREAL(50.4, 0, 100) != 50.0 {
		t.Error("RoundAndClampLREAL rounding down failed")
	}
	if RoundAndClampLREAL(150.2, 0, 100) != 100.0 {
		t.Error("RoundAndClampLREAL clamping max failed")
	}
}
