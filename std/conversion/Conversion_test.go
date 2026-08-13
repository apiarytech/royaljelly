package conversion

import (
	"errors"
	"testing"
	"time"

	. "github.com/apiarytech/royaljelly/core"
	. "github.com/apiarytech/royaljelly/iec"
	. "github.com/apiarytech/royaljelly/std/time"
)

func TestBoolConversions(t *testing.T) {
	if BOOL_TO_INT(true) != 1 {
		t.Error("BOOL_TO_INT(true) should be 1")
	}
	if BOOL_TO_STRING(false) != "false" {
		t.Error("BOOL_TO_STRING(false) should be 'false'")
	}
}

func TestIntConversions(t *testing.T) {
	if INT_TO_REAL(123) != 123.0 {
		t.Error("INT_TO_REAL failed")
	}
	if DINT_TO_STRING(-456) != "-456" {
		t.Error("DINT_TO_STRING failed")
	}
	if LINT_TO_BOOL(0) != false {
		t.Error("LINT_TO_BOOL(0) failed")
	}
	if LINT_TO_BOOL(1) != true {
		t.Error("LINT_TO_BOOL(1) failed")
	}
	if LINT_TO_SINT(130) != 127 {
		t.Errorf("LINT_TO_SINT overflow failed, got %d, want 127", LINT_TO_SINT(130))
	}
	if LINT_TO_SINT(-130) != -128 {
		t.Errorf("LINT_TO_SINT underflow failed, got %d, want -128", LINT_TO_SINT(-130))
	}
	if DINT_TO_INT(40000) != 32767 {
		t.Errorf("DINT_TO_INT overflow failed, got %d, want 32767", DINT_TO_INT(40000))
	}
	if ULINT_TO_UINT(70000) != 65535 {
		t.Errorf("ULINT_TO_UINT overflow failed, got %d, want 65535", ULINT_TO_UINT(70000))
	}
}

func TestRealConversions(t *testing.T) {
	if REAL_TO_INT(123.7) != 124 {
		t.Errorf("REAL_TO_INT(123.7) was %d, want 124", REAL_TO_INT(123.7))
	}
	if LREAL_TO_DINT(-45.6) != -46 {
		t.Errorf("LREAL_TO_DINT(-45.6) was %d, want -46", LREAL_TO_DINT(-45.6))
	}
	if REAL_TO_BOOL(0.0) != false {
		t.Error("REAL_TO_BOOL(0.0) failed")
	}
	if LREAL_TO_SINT(200.0) != 127 {
		t.Errorf("LREAL_TO_SINT overflow failed, got %d, want 127", LREAL_TO_SINT(200.0))
	}
	if LREAL_TO_USINT(-10.0) != 0 {
		t.Errorf("LREAL_TO_USINT underflow failed, got %d, want 0", LREAL_TO_USINT(-10.0))
	}
}

func TestStringConversions(t *testing.T) {
	// Note: String to numeric is handled by AnyToLINT/AnyToLREAL, not direct conversion functions
	// in the same way as other types.
	val, err := AnyToLINT(STRING("123"))
	if err != nil || val != 123 {
		t.Error("AnyToLINT from STRING failed")
	}

	fVal, err := AnyToLREAL(STRING("-123.45"))
	if err != nil || fVal != -123.45 {
		t.Error("AnyToLREAL from STRING failed")
	}
}

func TestTimeConversions(t *testing.T) {
	d := DATE(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	dt := DT(time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC))
	tod := TOD(time.Date(0, 0, 0, 1, 2, 3, 0, time.UTC))
	tm := TIME(10 * time.Second)

	if DATE_TO_STRING(d) != "D#2024-01-01" {
		t.Errorf("DATE_TO_STRING failed, got %s", DATE_TO_STRING(d))
	}

	if DT_TO_DATE(dt) != DATE(time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)) {
		t.Error("DT_TO_DATE failed")
	}

	// TOD has date components from its creation, but they should be ignored in conversion to LINT
	expectedTodMs := LINT((1*time.Hour + 2*time.Minute + 3*time.Second).Milliseconds())
	if TOD_TO_LINT(tod) != expectedTodMs {
		t.Errorf("TOD_TO_LINT failed, got %d, want %d", TOD_TO_LINT(tod), expectedTodMs)
	}

	if TIME_TO_LREAL(tm) != 10000.0 {
		t.Error("TIME_TO_LREAL failed")
	}
}

func TestBitFloatConversions(t *testing.T) {
	var r REAL = -123.45
	bits := REAL_TO_BITS(r)
	r2 := BITS_TO_REAL(bits)
	if r != r2 {
		t.Errorf("REAL <-> BITS conversion failed. In: %f, Out: %f", r, r2)
	}

	var lr LREAL = 9876.54321
	lbits := LREAL_TO_BITS(lr)
	lr2 := BITS_TO_LREAL(lbits)
	if lr != lr2 {
		t.Errorf("LREAL <-> BITS conversion failed. In: %f, Out: %f", lr, lr2)
	}
}

func TestBCDConversions(t *testing.T) {
	t.Run("Valid BCD Conversions", func(t *testing.T) {
		val1, err1 := USINT_TO_BCD_BYTE(99)
		if err1 != nil || val1 != 0x99 {
			t.Errorf("USINT_TO_BCD_BYTE failed, got 0x%X, err: %v", val1, err1)
		}

		val2, err2 := BYTE_BCD_TO_USINT(0x99)
		if err2 != nil || val2 != 99 {
			t.Errorf("BYTE_BCD_TO_USINT failed, got %d, err: %v", val2, err2)
		}

		val3, err3 := WORD_BCD_TO_UINT(0x1234)
		if err3 != nil || val3 != 1234 {
			t.Errorf("WORD_BCD_TO_UINT failed, got %d, err: %v", val3, err3)
		}
	})

	t.Run("Invalid BCD Nibble", func(t *testing.T) {
		_, err := BYTE_BCD_TO_USINT(0xA0)
		if err == nil {
			t.Error("BYTE_BCD_TO_USINT with invalid nibble should have returned an error")
		}

		var convErr *ConversionError
		if !errors.As(err, &convErr) {
			t.Errorf("Expected a ConversionError, but got %T", err)
		}
	})

	t.Run("Value too large for BCD", func(t *testing.T) {
		_, err := USINT_TO_BCD_BYTE(101) // Use 101 to avoid ambiguity with hex 0x100
		if err == nil {
			t.Error("USINT_TO_BCD_BYTE with value > 99 should have returned an error")
		}

		var convErr *ConversionError
		if !errors.As(err, &convErr) {
			t.Errorf("Expected a ConversionError, but got %T", err)
		} else if convErr.Reason != "value out of range" {
			t.Errorf("Unexpected error reason: got %q, want 'value out of range'", convErr.Reason)
		}
	})
}

func TestBytesToString(t *testing.T) {
	input := []BYTE{0x48, 0x65, 0x6C, 0x6C, 0x6F} // "Hello"
	expected := STRING("Hello")
	result := BYTES_TO_STRING(input)
	if result != expected {
		t.Errorf("BYTES_TO_STRING() = %q; want %q", result, expected)
	}
}

func TestValueMethods(t *testing.T) {
	t.Run("BOOL.Value()", func(t *testing.T) {
		bTrue := BOOL(true)
		if bTrue.Value() != true {
			t.Errorf("BOOL(true).Value() = %v; want true", bTrue.Value())
		}
		bFalse := BOOL(false)
		if bFalse.Value() != false {
			t.Errorf("BOOL(false).Value() = %v; want false", bFalse.Value())
		}
	})

	t.Run("BYTE.Value()", func(t *testing.T) {
		b := BYTE(0xAB)
		if b.Value() != 0xAB {
			t.Errorf("BYTE(0xAB).Value() = 0x%X; want 0xAB", b.Value())
		}
		bZero := BYTE(0)
		if bZero.Value() != 0 {
			t.Errorf("BYTE(0).Value() = 0x%X; want 0x0", bZero.Value())
		}
		// Test conversion with clamping
		if SINT_TO_USINT(-1) != 0 {
			t.Errorf("SINT_TO_USINT(-1) should be clamped to 0, got %d", SINT_TO_USINT(-1))
		}
		if INT_TO_SINT(300) != 127 {
			t.Errorf("INT_TO_SINT(300) should be clamped to 127, got %d", INT_TO_SINT(300))
		}
	})

	t.Run("WORD.Value()", func(t *testing.T) {
		w := WORD(0xABCD)
		if w.Value() != 0xABCD {
			t.Errorf("WORD(0xABCD).Value() = 0x%X; want 0xABCD", w.Value())
		}
	})

	t.Run("DWORD.Value()", func(t *testing.T) {
		d := DWORD(0x12345678)
		if d.Value() != 0x12345678 {
			t.Errorf("DWORD(0x12345678).Value() = 0x%X; want 0x12345678", d.Value())
		}
	})

	t.Run("LWORD.Value()", func(t *testing.T) {
		l := LWORD(0x1234567890ABCDEF)
		if l.Value() != 0x1234567890ABCDEF {
			t.Errorf("LWORD(0x1234567890ABCDEF).Value() = 0x%X; want 0x1234567890ABCDEF", l.Value())
		}
	})

	t.Run("REAL.Value()", func(t *testing.T) {
		r := REAL(123.45)
		// Using a pointer to REAL for the Value() method
		if (&r).Value() != float32(123.45) {
			t.Errorf("REAL(123.45).Value() = %f; want %f", (&r).Value(), float32(123.45))
		}
		rZero := REAL(0.0)
		if (&rZero).Value() != float32(0.0) {
			t.Errorf("REAL(0.0).Value() = %f; want %f", (&rZero).Value(), float32(0.0))
		}
	})
}

func TestConvertMethods(t *testing.T) {
	t.Run("DATE.CONVERT()", func(t *testing.T) {
		// DATE is a point in time, CONVERT should return milliseconds since Unix epoch.
		d := DATE(time.Unix(1000, 0).UTC()) // 1000 seconds since epoch
		expected := LINT(1000 * 1000)       // 1000 seconds in milliseconds

		resultVal := d.CONVERT()

		if resultVal != expected {
			t.Errorf("DATE.CONVERT() = %d; want %d", resultVal, expected)
		}
	})

	t.Run("DT.CONVERT()", func(t *testing.T) {
		// DT is a point in time, CONVERT should return milliseconds since Unix epoch.
		dt := DT(time.Unix(2000, 500*1e6).UTC()) // 2000 seconds and 500 ms since epoch
		expected := LINT(2000*1000 + 500)        // in milliseconds

		resultVal := dt.CONVERT()

		if resultVal != expected {
			t.Errorf("DT.CONVERT() = %d; want %d", resultVal, expected)
		}
	})

	t.Run("TOD.CONVERT()", func(t *testing.T) {
		// TOD is time since midnight.
		tod := TOD(time.Date(0, 0, 0, 1, 2, 3, 456*1e6, time.UTC))
		expected := LINT((1*time.Hour + 2*time.Minute + 3*time.Second + 456*time.Millisecond).Milliseconds())

		resultVal := tod.CONVERT()

		if resultVal != expected {
			t.Errorf("TOD.CONVERT() = %d; want %d", resultVal, expected)
		}
	})
}

func TestSubConversionErrors(t *testing.T) {
	invalidInput := "not a number"

	t.Run("SubByte Error", func(t *testing.T) {
		_, err := SubByte(invalidInput)
		if err == nil {
			t.Error("SubByte with invalid input should have returned an error")
		}
	})
	t.Run("SubWord Error", func(t *testing.T) {
		_, err := SubWord(invalidInput)
		if err == nil {
			t.Error("SubWord with invalid input should have returned an error")
		}
	})
	t.Run("SubDword Error", func(t *testing.T) {
		_, err := SubDword(invalidInput)
		if err == nil {
			t.Error("SubDword with invalid input should have returned an error")
		}
	})
	t.Run("SubLword Error", func(t *testing.T) {
		_, err := SubLword(invalidInput)
		if err == nil {
			t.Error("SubLword with invalid input should have returned an error")
		}
	})
	t.Run("SubDt Error", func(t *testing.T) {
		_, err := SubDt(invalidInput)
		if err == nil {
			t.Error("SubDt with invalid input should have returned an error")
		}
	})
	t.Run("SubDate Error", func(t *testing.T) {
		_, err := SubDate(invalidInput)
		if err == nil {
			t.Error("SubDate with invalid input should have returned an error")
		}
	})

}

func TestAllConversions(t *testing.T) {
	t.Run("BOOL Conversions", func(t *testing.T) {
		if BYTE_TO_BOOL(1) != true || BYTE_TO_BOOL(0) != false {
			t.Error("BYTE_TO_BOOL conversion failed")
		}
		if WORD_TO_BOOL(1) != true || WORD_TO_BOOL(0) != false {
			t.Error("WORD_TO_BOOL conversion failed")
		}
		if DWORD_TO_BOOL(1) != true || DWORD_TO_BOOL(0) != false {
			t.Error("DWORD_TO_BOOL conversion failed")
		}
		if LWORD_TO_BOOL(1) != true || LWORD_TO_BOOL(0) != false {
			t.Error("LWORD_TO_BOOL conversion failed")
		}
		if BOOL_TO_BYTE(true) != 1 {
			t.Error("BOOL_TO_BYTE failed")
		}
		if BOOL_TO_WORD(true) != 1 {
			t.Error("BOOL_TO_WORD failed")
		}
		if BOOL_TO_DWORD(true) != 1 {
			t.Error("BOOL_TO_DWORD failed")
		}
		if BOOL_TO_LWORD(true) != 1 {
			t.Error("BOOL_TO_LWORD failed")
		}
		if BOOL_TO_LINT(true) != 1 {
			t.Error("BOOL_TO_LINT failed")
		}
		if BOOL_TO_SINT(true) != 1 || BOOL_TO_SINT(false) != 0 {
			t.Error("BOOL_TO_SINT failed")
		}
		if BOOL_TO_DINT(true) != 1 || BOOL_TO_DINT(false) != 0 {
			t.Error("BOOL_TO_DINT failed")
		}
		if BOOL_TO_USINT(true) != 1 || BOOL_TO_USINT(false) != 0 {
			t.Error("BOOL_TO_USINT failed")
		}
		if BOOL_TO_UINT(true) != 1 || BOOL_TO_UINT(false) != 0 {
			t.Error("BOOL_TO_UINT failed")
		}
		if BOOL_TO_UDINT(true) != 1 || BOOL_TO_UDINT(false) != 0 {
			t.Error("BOOL_TO_UDINT failed")
		}
		if BOOL_TO_ULINT(true) != 1 || BOOL_TO_ULINT(false) != 0 {
			t.Error("BOOL_TO_ULINT failed")
		}
		if BOOL_TO_REAL(true) != 1.0 || BOOL_TO_REAL(false) != 0.0 {
			t.Error("BOOL_TO_REAL failed")
		}
		if BOOL_TO_LREAL(true) != 1.0 || BOOL_TO_LREAL(false) != 0.0 {
			t.Error("BOOL_TO_LREAL failed")
		}
		if BOOL_TO_TIME(true) != TIME(1*time.Millisecond) {
			t.Error("BOOL_TO_TIME failed")
		}
		if BOOL_TO_DATE(true) != DATE(time.UnixMilli(1)) {
			t.Error("BOOL_TO_DATE failed")
		}
		if BOOL_TO_TOD(true) != TOD(time.Time{}.Add(1*time.Millisecond)) {
			t.Error("BOOL_TO_TOD failed")
		}
		if BOOL_TO_DT(true) != DT(time.UnixMilli(1)) {
			t.Error("BOOL_TO_DT failed")
		}
	})

	t.Run("BYTE Conversions", func(t *testing.T) {
		if BYTE_TO_SINT(100) != 100 || BYTE_TO_INT(100) != 100 || BYTE_TO_DINT(100) != 100 || BYTE_TO_LINT(100) != 100 {
			t.Error("BYTE_TO integer conversions failed")
		}
		if BYTE_TO_USINT(100) != 100 || BYTE_TO_UINT(100) != 100 || BYTE_TO_UDINT(100) != 100 || BYTE_TO_ULINT(100) != 100 {
			t.Error("BYTE_TO unsigned integer conversions failed")
		}
		if BYTE_TO_REAL(105) != 105 || BYTE_TO_LREAL(105) != 105 {
			t.Error("BYTE_TO real conversions failed")
		}
		if BYTE_TO_STRING(65) != "65" {
			t.Error("BYTE_TO_STRING conversion failed")
		}
		if BYTE_TO_WORD(10) != 10 {
			t.Error("BYTE_TO_WORD failed")
		}
		if BYTE_TO_DWORD(10) != 10 {
			t.Error("BYTE_TO_DWORD failed")
		}
		if BYTE_TO_LWORD(10) != 10 {
			t.Error("BYTE_TO_LWORD failed")
		}
		if BYTE_TO_TIME(100) != TIME(100*time.Millisecond) {
			t.Error("BYTE_TO_TIME failed")
		}
		if BYTE_TO_DATE(100) != DATE(time.UnixMilli(100)) {
			t.Error("BYTE_TO_DATE failed")
		}
		if BYTE_TO_DT(100) != DT(time.UnixMilli(100)) {
			t.Error("BYTE_TO_DT failed")
		}
		if BYTE_TO_TOD(100) != TOD(time.Time{}.Add(100*time.Millisecond)) {
			t.Error("BYTE_TO_TOD failed")
		}
	})

	t.Run("WORD Conversions", func(t *testing.T) {
		if WORD_TO_SINT(100) != 100 || WORD_TO_INT(100) != 100 || WORD_TO_DINT(100) != 100 || WORD_TO_LINT(100) != 100 {
			t.Error("WORD_TO integer conversions failed")
		}
		if WORD_TO_USINT(100) != 100 || WORD_TO_UINT(100) != 100 || WORD_TO_UDINT(100) != 100 || WORD_TO_ULINT(100) != 100 {
			t.Error("WORD_TO unsigned integer conversions failed")
		}
		if WORD_TO_REAL(105) != 105 || WORD_TO_LREAL(105) != 105 {
			t.Error("WORD_TO real conversions failed")
		}
		if WORD_TO_STRING(65) != "65" {
			t.Error("WORD_TO_STRING conversion failed")
		}
		if WORD_TO_BYTE(0x1FF) != 0xFF {
			t.Error("WORD_TO_BYTE failed")
		}
		if WORD_TO_DWORD(10) != 10 {
			t.Error("WORD_TO_DWORD failed")
		}
		if WORD_TO_LWORD(10) != 10 {
			t.Error("WORD_TO_LWORD failed")
		}
		if WORD_TO_TIME(5000) != TIME(5*time.Second) {
			t.Error("WORD_TO_TIME failed")
		}
		if WORD_TO_DATE(5000) != DATE(time.UnixMilli(5000)) {
			t.Error("WORD_TO_DATE failed")
		}
		if WORD_TO_DT(5000) != DT(time.UnixMilli(5000)) {
			t.Error("WORD_TO_DT failed")
		}
		if WORD_TO_TOD(5000) != TOD(time.Time{}.Add(5000*time.Millisecond)) {
			t.Error("WORD_TO_TOD failed")
		}
	})

	t.Run("DWORD Conversions", func(t *testing.T) {
		if DWORD_TO_SINT(100) != 100 || DWORD_TO_INT(100) != 100 || DWORD_TO_DINT(100) != 100 || DWORD_TO_LINT(100) != 100 {
			t.Error("DWORD_TO integer conversions failed")
		}
		if DWORD_TO_USINT(100) != 100 || DWORD_TO_UINT(100) != 100 || DWORD_TO_UDINT(100) != 100 || DWORD_TO_ULINT(100) != 100 {
			t.Error("DWORD_TO unsigned integer conversions failed")
		}
		if DWORD_TO_REAL(105) != 105 || DWORD_TO_LREAL(105) != 105 {
			t.Error("DWORD_TO real conversions failed")
		}
		if DWORD_TO_STRING(65) != "65" {
			t.Error("DWORD_TO_STRING conversion failed")
		}
		if DWORD_TO_BYTE(0x1FF) != 0xFF {
			t.Error("DWORD_TO_BYTE failed")
		}
		if DWORD_TO_WORD(0x1FFFF) != 0xFFFF {
			t.Error("DWORD_TO_WORD failed")
		}
		if DWORD_TO_LWORD(10) != 10 {
			t.Error("DWORD_TO_LWORD failed")
		}
		if DWORD_TO_TIME(5000) != TIME(5*time.Second) {
			t.Error("DWORD_TO_TIME failed")
		}
		if DWORD_TO_DATE(5000) != DATE(time.UnixMilli(5000)) {
			t.Error("DWORD_TO_DATE failed")
		}
		if DWORD_TO_DT(5000) != DT(time.UnixMilli(5000)) {
			t.Error("DWORD_TO_DT failed")
		}
		if DWORD_TO_TOD(5000) != TOD(time.Time{}.Add(5000*time.Millisecond)) {
			t.Error("DWORD_TO_TOD failed")
		}
	})

	t.Run("LWORD Conversions", func(t *testing.T) {
		if LWORD_TO_SINT(100) != 100 || LWORD_TO_INT(100) != 100 || LWORD_TO_DINT(100) != 100 || LWORD_TO_LINT(100) != 100 {
			t.Error("LWORD_TO integer conversions failed")
		}
		if LWORD_TO_USINT(100) != 100 || LWORD_TO_UINT(100) != 100 || LWORD_TO_UDINT(100) != 100 || LWORD_TO_ULINT(100) != 100 {
			t.Error("LWORD_TO unsigned integer conversions failed")
		}
		if LWORD_TO_REAL(105) != 105 || LWORD_TO_LREAL(105) != 105 {
			t.Error("LWORD_TO real conversions failed")
		}
		if LWORD_TO_STRING(65) != "65" {
			t.Error("LWORD_TO_STRING conversion failed")
		}
		if LWORD_TO_BYTE(0x1FF) != 0xFF {
			t.Error("LWORD_TO_BYTE failed")
		}
		if LWORD_TO_WORD(0x1FFFF) != 0xFFFF {
			t.Error("LWORD_TO_WORD failed")
		}
		if LWORD_TO_DWORD(0x1FFFFFFFF) != 0xFFFFFFFF {
			t.Error("LWORD_TO_DWORD failed")
		}
		if LWORD_TO_TIME(5000) != TIME(5*time.Second) {
			t.Error("LWORD_TO_TIME failed")
		}
		if LWORD_TO_DATE(5000) != DATE(time.UnixMilli(5000)) {
			t.Error("LWORD_TO_DATE failed")
		}
		if LWORD_TO_DT(5000) != DT(time.UnixMilli(5000)) {
			t.Error("LWORD_TO_DT failed")
		}
		if LWORD_TO_TOD(5000) != TOD(time.Time{}.Add(5000*time.Millisecond)) {
			t.Error("LWORD_TO_TOD failed")
		}
	})

	t.Run("REAL Conversions", func(t *testing.T) {
		if REAL_TO_SINT(100.0) != 100 || REAL_TO_INT(100.0) != 100 || REAL_TO_DINT(100.0) != 100 || REAL_TO_LINT(100.0) != 100 {
			t.Error("REAL_TO integer conversions failed")
		}
		if REAL_TO_USINT(100.0) != 100 || REAL_TO_UINT(100.0) != 100 || REAL_TO_UDINT(100.0) != 100 || REAL_TO_ULINT(100.0) != 100 {
			t.Error("REAL_TO unsigned integer conversions failed")
		}
		if REAL_TO_LREAL(10.5) != 10.5 {
			t.Error("REAL_TO real conversions failed")
		}
		if REAL_TO_STRING(65.0) != "65" {
			t.Error("REAL_TO_STRING conversion failed")
		}
		if REAL_TO_BYTE(256.0) != 255 {
			t.Error("REAL_TO_BYTE failed")
		}
		if REAL_TO_WORD(65536.0) != 65535 {
			t.Error("REAL_TO_WORD failed")
		}
		if REAL_TO_DWORD(4294967296.0) != 4294967295 {
			t.Error("REAL_TO_DWORD failed")
		}
		if REAL_TO_LWORD(0) != 0 {
			t.Error("REAL_TO_LWORD failed")
		}
		if REAL_TO_TIME(5000.5) != TIME(5*time.Second) {
			t.Error("REAL_TO_TIME failed")
		}
		if REAL_TO_DATE(5000.5) != DATE(time.UnixMilli(5000)) {
			t.Error("REAL_TO_DATE failed")
		}
		if REAL_TO_DT(5000.5) != DT(time.UnixMilli(5000)) {
			t.Error("REAL_TO_DT failed")
		}
		if REAL_TO_TOD(5000.5) != TOD(time.Time{}.Add(5000*time.Millisecond)) {
			t.Error("REAL_TO_TOD failed")
		}
	})

	t.Run("LREAL Conversions", func(t *testing.T) {
		if LREAL_TO_SINT(100.0) != 100 || LREAL_TO_INT(100.0) != 100 || LREAL_TO_DINT(100.0) != 100 || LREAL_TO_LINT(100.0) != 100 {
			t.Error("LREAL_TO integer conversions failed")
		}
		if LREAL_TO_USINT(100.0) != 100 || LREAL_TO_UINT(100.0) != 100 || LREAL_TO_UDINT(100.0) != 100 || LREAL_TO_ULINT(100.0) != 100 {
			t.Error("LREAL_TO unsigned integer conversions failed")
		}
		if LREAL_TO_STRING(65.0) != "65" {
			t.Error("LREAL_TO_STRING conversion failed")
		}
		if LREAL_TO_BYTE(256.0) != 255 {
			t.Error("LREAL_TO_BYTE failed")
		}
		if LREAL_TO_WORD(65536.0) != 65535 {
			t.Error("LREAL_TO_WORD failed")
		}
		if LREAL_TO_DWORD(4294967296.0) != 4294967295 {
			t.Error("LREAL_TO_DWORD failed")
		}
		if LREAL_TO_LWORD(0) != 0 {
			t.Error("LREAL_TO_LWORD failed")
		}
		if LREAL_TO_TIME(5000.5) != TIME(5*time.Second) {
			t.Error("LREAL_TO_TIME failed")
		}
		if LREAL_TO_DATE(5000.5) != DATE(time.UnixMilli(5000)) {
			t.Error("LREAL_TO_DATE failed")
		}
		if LREAL_TO_DT(5000.5) != DT(time.UnixMilli(5000)) {
			t.Error("LREAL_TO_DT failed")
		}
		if LREAL_TO_TOD(5000.5) != TOD(time.Time{}.Add(5000*time.Millisecond)) {
			t.Error("LREAL_TO_TOD failed")
		}
		if LREAL_TO_ULINT(100.5) != 100 {
			t.Error("LREAL_TO_ULINT failed for positive value")
		}
		if LREAL_TO_ULINT(-10.5) != 0 { // Clamped to 0
			t.Error("LREAL_TO_ULINT failed for negative value (clamping)")
		}
		if LREAL_TO_REAL(123.45) != 123.45 {
			t.Error("LREAL_TO_REAL failed")
		}
		if LREAL_TO_BOOL(1.0) != true || LREAL_TO_BOOL(0.0) != false {
			t.Error("LREAL_TO_BOOL failed")
		}
	})

	t.Run("SINT Conversions", func(t *testing.T) {
		if SINT_TO_REAL(100) != 100.0 || SINT_TO_LREAL(100) != 100.0 {
			t.Error("SINT_TO real conversions failed")
		}

		if SINT_TO_DINT(100) != 100 || SINT_TO_LINT(100) != 100 {
			t.Error("SINT_TO integer conversions failed")
		}

		if SINT_TO_UDINT(100) != 100 || SINT_TO_ULINT(100) != 100 {
			t.Error("SINT_TO unsigned integer conversions failed")
		}

		if SINT_TO_STRING(65) != "65" {
			t.Error("SINT_TO_STRING conversion failed")
		}
		if SINT_TO_BYTE(10) != 10 {
			t.Error("SINT_TO_BYTE failed")
		}
		if SINT_TO_WORD(10) != 10 {
			t.Error("SINT_TO_WORD failed")
		}
		if SINT_TO_DWORD(10) != 10 {
			t.Error("SINT_TO_DWORD failed")
		}
		if SINT_TO_LWORD(10) != 10 {
			t.Error("SINT_TO_LWORD failed")
		}
		if SINT_TO_TIME(100) != TIME(100*time.Millisecond) {
			t.Error("SINT_TO_TIME failed")
		}
		if SINT_TO_DATE(100) != DATE(time.UnixMilli(100)) {
			t.Error("SINT_TO_DATE failed")
		}
		if SINT_TO_DT(100) != DT(time.UnixMilli(100)) {
			t.Error("SINT_TO_DT failed")
		}
		if SINT_TO_TOD(100) != TOD(time.Time{}.Add(100*time.Millisecond)) {
			t.Error("SINT_TO_TOD failed")
		}
		if SINT_TO_UINT(-1) != 0 {
			t.Error("SINT_TO_UINT failed")
		}
		if SINT_TO_BOOL(1) != true || SINT_TO_BOOL(0) != false {
			t.Error("SINT_TO_BOOL failed")
		}
		if SINT_TO_INT(100) != 100 {
			t.Error("SINT_TO_INT failed")
		}
		if SINT_TO_ULINT(100) != 100 {
			t.Error("SINT_TO_ULINT failed")
		}
		if SINT_TO_REAL(100) != 100.0 {
			t.Error("SINT_TO_REAL failed")
		}
		if SINT_TO_LINT(100) != 100 {
			t.Error("SINT_TO_LINT failed")
		}
		if SINT_TO_DINT(100) != 100 {
			t.Error("SINT_TO_DINT failed")
		}
		if SINT_TO_STRING(-123) != "-123" {
			t.Error("SINT_TO_STRING failed")
		}
		if LINT_TO_LWORD(123) != 123 {
			t.Error("LINT_TO_LWORD failed")
		}
		if LINT_TO_DWORD(123) != 123 {
			t.Error("LINT_TO_DWORD failed")
		}
	})

	t.Run("INT Conversions", func(t *testing.T) {
		if INT_TO_REAL(100) != 100.0 || INT_TO_LREAL(100) != 100.0 {
			t.Error("INT_TO real conversions failed")
		}

		if INT_TO_DINT(100) != 100 || INT_TO_LINT(100) != 100 {
			t.Error("INT_TO integer conversions failed")
		}

		if INT_TO_UDINT(100) != 100 || INT_TO_ULINT(100) != 100 {
			t.Error("INT_TO unsigned integer conversions failed")
		}

		if INT_TO_STRING(65) != "65" {
			t.Error("INT_TO_STRING conversion failed")
		}
		if INT_TO_BYTE(300) != 255 {
			t.Error("INT_TO_BYTE failed")
		}
		if INT_TO_LWORD(10) != 10 {
			t.Error("INT_TO_LWORD failed")
		}
		if INT_TO_TIME(5000) != TIME(5*time.Second) {
			t.Error("INT_TO_TIME failed")
		}
		if INT_TO_DATE(5000) != DATE(time.UnixMilli(5000)) {
			t.Error("INT_TO_DATE failed")
		}
		if INT_TO_DT(5000) != DT(time.UnixMilli(5000)) {
			t.Error("INT_TO_DT failed")
		}
		if INT_TO_TOD(5000) != TOD(time.Time{}.Add(5000*time.Millisecond)) {
			t.Error("INT_TO_TOD failed")
		}
		if INT_TO_DWORD(10) != 10 {
			t.Error("INT_TO_DWORD failed")
		}
		if INT_TO_WORD(32767) != 32767 {
			t.Errorf("INT_TO_WORD failed, got %d", INT_TO_WORD(32767))
		}
		if INT_TO_UINT(-1) != 0 {
			t.Error("INT_TO_UINT failed")
		}
		if INT_TO_USINT(300) != 255 {
			t.Error("INT_TO_USINT failed")
		}
		if INT_TO_BOOL(1) != true || INT_TO_BOOL(0) != false {
			t.Error("INT_TO_BOOL failed")
		}
		if INT_TO_ULINT(100) != 100 {
			t.Error("INT_TO_ULINT failed")
		}
		if INT_TO_SINT(100) != 100 {
			t.Error("INT_TO_SINT failed")
		}
		if INT_TO_LINT(100) != 100 {
			t.Error("INT_TO_LINT failed")
		}
		if INT_TO_DINT(100) != 100 {
			t.Error("INT_TO_DINT failed")
		}
		if INT_TO_REAL(100) != 100.0 {
			t.Error("INT_TO_REAL failed")
		}
		if INT_TO_LREAL(100) != 100.0 {
			t.Error("INT_TO_LREAL failed")
		}
		if LINT_TO_LWORD(123) != 123 {
			t.Error("LINT_TO_LWORD failed")
		}
		if LINT_TO_DWORD(123) != 123 {
			t.Error("LINT_TO_DWORD failed")
		}
	})

	t.Run("LINT Conversions", func(t *testing.T) {
		if LINT_TO_SINT(100) != 100 || LINT_TO_INT(100) != 100 || LINT_TO_DINT(100) != 100 {
			t.Error("LINT_TO smaller signed integer conversions failed")
		}
		if LINT_TO_USINT(100) != 100 || LINT_TO_UINT(100) != 100 || LINT_TO_UDINT(100) != 100 {
			t.Error("LINT_TO smaller unsigned integer conversions failed")
		}
		if LINT_TO_REAL(105) != 105.0 || LINT_TO_LREAL(105) != 105.0 {
			t.Error("LINT_TO real conversions failed")
		}
		if LINT_TO_STRING(65) != "65" {
			t.Error("LINT_TO_STRING conversion failed")
		}
		if LINT_TO_BYTE(300) != 255 {
			t.Error("LINT_TO_BYTE failed")
		}
		if LINT_TO_WORD(70000) != 65535 {
			t.Error("LINT_TO_WORD failed")
		}
		if LINT_TO_TIME(5000) != TIME(5*time.Second) {
			t.Error("LINT_TO_TIME failed")
		}
		if LINT_TO_DATE(5000) != DATE(time.UnixMilli(5000)) {
			t.Error("LINT_TO_DATE failed")
		}
		if LINT_TO_TOD(5000) != TOD(time.Time{}.Add(5000*time.Millisecond)) {
			t.Error("LINT_TO_TOD failed")
		}
		if LINT_TO_ULINT(100) != 100 {
			t.Error("LINT_TO_ULINT failed")
		}
		if LINT_TO_LWORD(123) != 123 {
			t.Error("LINT_TO_LWORD failed")
		}
		if LINT_TO_DWORD(123) != 123 {
			t.Error("LINT_TO_DWORD failed")
		}
		if LINT_TO_INT(123) != 123 {
			t.Error("LINT_TO_INT failed")
		}
	})

	t.Run("DINT Conversions", func(t *testing.T) {
		if DINT_TO_SINT(100) != 100 || DINT_TO_INT(100) != 100 || DINT_TO_LINT(100) != 100 {
			t.Error("DINT_TO other signed integer conversions failed")
		}
		if DINT_TO_USINT(100) != 100 || DINT_TO_UINT(100) != 100 || DINT_TO_UDINT(100) != 100 || DINT_TO_ULINT(100) != 100 {
			t.Error("DINT_TO unsigned integer conversions failed")
		}
		if DINT_TO_REAL(105) != 105.0 || DINT_TO_LREAL(105) != 105.0 {
			t.Error("DINT_TO real conversions failed")
		}
		if DINT_TO_STRING(65) != "65" {
			t.Error("DINT_TO_STRING conversion failed")
		}
		if DINT_TO_BYTE(300) != 255 {
			t.Error("DINT_TO_BYTE failed")
		}
		if DINT_TO_WORD(70000) != 65535 {
			t.Error("DINT_TO_WORD failed")
		}
		if DINT_TO_DWORD(0) != 0 {
			t.Error("DINT_TO_DWORD failed")
		}
		if DINT_TO_TIME(5000) != TIME(5*time.Second) {
			t.Error("DINT_TO_TIME failed")
		}
		if DINT_TO_DATE(5000) != DATE(time.UnixMilli(5000)) {
			t.Error("DINT_TO_DATE failed")
		}
		if DINT_TO_DT(5000) != DT(time.UnixMilli(5000)) {
			t.Error("DINT_TO_DT failed")
		}
		if DINT_TO_TOD(5000) != TOD(time.Time{}.Add(5000*time.Millisecond)) {
			t.Error("DINT_TO_TOD failed")
		}
		if DINT_TO_ULINT(100) != 100 {
			t.Error("DINT_TO_ULINT failed")
		}
		if DINT_TO_LWORD(123) != 123 {
			t.Error("DINT_TO_LWORD failed")
		}
		if DINT_TO_LWORD(123) != 123 {
			t.Error("DINT_TO_LWORD failed")
		}
		if DINT_TO_LWORD(123) != 123 {
			t.Error("DINT_TO_LWORD failed")
		}
		if DINT_TO_BOOL(1) != true || DINT_TO_BOOL(0) != false {
			t.Error("DINT_TO_BOOL failed")
		}
		if DINT_TO_SINT(100) != 100 {
			t.Error("DINT_TO_SINT failed")
		}
		if DINT_TO_INT(100) != 100 {
			t.Error("DINT_TO_INT failed")
		}
		if DINT_TO_LINT(100) != 100 {
			t.Error("DINT_TO_LINT failed")
		}
	})

	t.Run("UINT Conversions", func(t *testing.T) {
		if UINT_TO_USINT(300) != 255 {
			t.Errorf("UINT_TO_USINT(300) should be clamped to 255, got %d", UINT_TO_USINT(300))
		}
		if UINT_TO_DATE(5000) != DATE(time.UnixMilli(5000)) {
			t.Error("UINT_TO_DATE failed")
		}
		if UINT_TO_DT(5000) != DT(time.UnixMilli(5000)) {
			t.Error("UINT_TO_DT failed")
		}
		if UINT_TO_ULINT(100) != 100 {
			t.Error("UINT_TO_ULINT failed")
		}
		if UINT_TO_SINT(100) != 100 {
			t.Error("UINT_TO_SINT failed")
		}
		if UINT_TO_LINT(100) != 100 {
			t.Error("UINT_TO_LINT failed")
		}
		if UINT_TO_DINT(100) != 100 {
			t.Error("UINT_TO_DINT failed")
		}
		if UINT_TO_DWORD(100) != 100 {
			t.Error("UINT_TO_DWORD failed")
		}
		if UINT_TO_LWORD(100) != 100 {
			t.Error("UINT_TO_LWORD failed")
		}
		if UINT_TO_BOOL(1) != true {
			t.Error("UINT_TO_BOOL failed")
		}
		if UINT_TO_SINT(100) != 100 {
			t.Error("UINT_TO_SINT failed")
		}
		if UINT_TO_LINT(100) != 100 {
			t.Error("UINT_TO_LINT failed")
		}
		if UINT_TO_DINT(100) != 100 {
			t.Error("UINT_TO_DINT failed")
		}
		if UINT_TO_DWORD(100) != 100 {
			t.Error("UINT_TO_DWORD failed")
		}
		if UINT_TO_LWORD(100) != 100 {
			t.Error("UINT_TO_LWORD failed")
		}
		if UINT_TO_BOOL(1) != true {
			t.Error("UINT_TO_BOOL failed")
		}
		if UINT_TO_REAL(100) != 100.0 {
			t.Error("UINT_TO_REAL failed")
		}
		if UINT_TO_TOD(5000) != TOD(time.Time{}.Add(5000*time.Millisecond)) {
			t.Error("UINT_TO_TOD failed")
		}
		if UINT_TO_UDINT(100) != 100 {
			t.Error("UINT_TO_UDINT failed")
		}
		if UINT_TO_WORD(100) != 100 {
			t.Error("UINT_TO_WORD failed")
		}
		if UINT_TO_LREAL(100) != 100.0 {
			t.Error("UINT_TO_LREAL failed")
		}
		if UINT_TO_BYTE(300) != 255 {
			t.Error("UINT_TO_BYTE failed")
		}
		if UINT_TO_TIME(5000) != TIME(5*time.Second) {
			t.Error("UINT_TO_TIME failed")
		}
		if UINT_TO_STRING(12345) != "12345" {
			t.Error("UINT_TO_STRING failed")
		}
		if UINT_TO_INT(100) != 100 {
			t.Error("UINT_TO_INT failed")
		}
	})

	t.Run("USINT Conversions", func(t *testing.T) {
		if USINT_TO_SINT(128) != 127 {
			t.Errorf("USINT_TO_SINT(128) should be clamped to 127, got %d", USINT_TO_SINT(128))
		}
		if USINT_TO_REAL(150) != 150.0 {
			t.Error("USINT_TO_REAL conversion failed")
		}
		if USINT_TO_STRING(200) != "200" {
			t.Error("USINT_TO_STRING conversion failed")
		}
		if USINT_TO_BYTE(10) != 10 {
			t.Error("USINT_TO_BYTE failed")
		}
		if USINT_TO_TIME(100) != TIME(100*time.Millisecond) {
			t.Error("USINT_TO_TIME failed")
		}
		if USINT_TO_DATE(100) != DATE(time.UnixMilli(100)) {
			t.Error("USINT_TO_DATE failed")
		}
		if USINT_TO_DT(100) != DT(time.UnixMilli(100)) {
			t.Error("USINT_TO_DT failed")
		}
		if USINT_TO_TOD(100) != TOD(time.Time{}.Add(100*time.Millisecond)) {
			t.Error("USINT_TO_TOD failed")
		}
		if USINT_TO_UDINT(100) != 100 {
			t.Error("USINT_TO_UDINT failed")
		}
		if USINT_TO_ULINT(100) != 100 {
			t.Error("USINT_TO_ULINT failed")
		}
		if USINT_TO_LINT(100) != 100 {
			t.Error("USINT_TO_LINT failed")
		}
		if USINT_TO_DINT(100) != 100 {
			t.Error("USINT_TO_DINT failed")
		}
		if USINT_TO_DWORD(100) != 100 {
			t.Error("USINT_TO_DWORD failed")
		}
		if USINT_TO_LWORD(100) != 100 {
			t.Error("USINT_TO_LWORD failed")
		}
		if USINT_TO_UINT(100) != 100 {
			t.Error("USINT_TO_UINT failed")
		}
		if USINT_TO_ULINT(100) != 100 {
			t.Error("USINT_TO_ULINT failed")
		}
		if USINT_TO_LINT(100) != 100 {
			t.Error("USINT_TO_LINT failed")
		}
		if USINT_TO_DINT(100) != 100 {
			t.Error("USINT_TO_DINT failed")
		}
		if USINT_TO_DWORD(100) != 100 {
			t.Error("USINT_TO_DWORD failed")
		}
		if USINT_TO_LWORD(100) != 100 {
			t.Error("USINT_TO_LWORD failed")
		}
		if USINT_TO_UINT(100) != 100 {
			t.Error("USINT_TO_UINT failed")
		}
		if USINT_TO_WORD(100) != 100 {
			t.Error("USINT_TO_WORD failed")
		}
		if USINT_TO_BOOL(1) != true {
			t.Error("USINT_TO_BOOL failed")
		}
		if USINT_TO_LINT(100) != 100 {
			t.Error("USINT_TO_LINT failed")
		}
		if USINT_TO_DINT(100) != 100 {
			t.Error("USINT_TO_DINT failed")
		}
		if USINT_TO_DWORD(100) != 100 {
			t.Error("USINT_TO_DWORD failed")
		}
		if USINT_TO_LWORD(100) != 100 {
			t.Error("USINT_TO_LWORD failed")
		}
		if USINT_TO_UINT(100) != 100 {
			t.Error("USINT_TO_UINT failed")
		}
		if USINT_TO_BOOL(1) != true {
			t.Error("USINT_TO_BOOL failed")
		}
		if USINT_TO_INT(100) != 100 {
			t.Error("USINT_TO_INT failed")
		}
		if USINT_TO_LREAL(150) != 150.0 {
			t.Error("USINT_TO_LREAL failed")
		}
		if USINT_TO_STRING(234) != "234" {
			t.Error("USINT_TO_STRING failed")
		}
	})

	t.Run("UDINT Conversions", func(t *testing.T) {
		if UDINT_TO_INT(40000) != 32767 {
			t.Errorf("UDINT_TO_INT(40000) should be clamped to 32767, got %d", UDINT_TO_INT(40000))
		}
		if UDINT_TO_REAL(12345) != 12345.0 {
			t.Error("UDINT_TO_REAL conversion failed")
		}
		if UDINT_TO_STRING(54321) != "54321" {
			t.Error("UDINT_TO_STRING conversion failed")
		}
		if UDINT_TO_WORD(70000) != 65535 {
			t.Error("UDINT_TO_WORD failed")
		}
		if UDINT_TO_TIME(5000) != TIME(5*time.Second) {
			t.Error("UDINT_TO_TIME failed")
		}
		if UDINT_TO_DATE(5000) != DATE(time.UnixMilli(5000)) {
			t.Error("UDINT_TO_DATE failed")
		}
		if UDINT_TO_DT(5000) != DT(time.UnixMilli(5000)) {
			t.Error("UDINT_TO_DT failed")
		}
		if UDINT_TO_TOD(5000) != TOD(time.Time{}.Add(5000*time.Millisecond)) {
			t.Error("UDINT_TO_TOD failed")
		}
		if UDINT_TO_ULINT(100) != 100 {
			t.Error("UDINT_TO_ULINT failed")
		}
		if UDINT_TO_SINT(100) != 100 {
			t.Error("UDINT_TO_SINT failed")
		}
		if UDINT_TO_LINT(100) != 100 {
			t.Error("UDINT_TO_LINT failed")
		}
		if UDINT_TO_DINT(100) != 100 {
			t.Error("UDINT_TO_DINT failed")
		}
		if UDINT_TO_DWORD(100) != 100 {
			t.Error("UDINT_TO_DWORD failed")
		}
		if UDINT_TO_LWORD(100) != 100 {
			t.Error("UDINT_TO_LWORD failed")
		}
		if UDINT_TO_BOOL(1) != true {
			t.Error("UDINT_TO_BOOL failed")
		}
		if UDINT_TO_LREAL(12345) != 12345.0 {
			t.Error("UDINT_TO_LREAL failed")
		}
		if UDINT_TO_BYTE(300) != 255 {
			t.Error("UDINT_TO_BYTE failed")
		}
		if UDINT_TO_USINT(300) != 255 {
			t.Error("UDINT_TO_USINT failed")
		}
		if UDINT_TO_UINT(70000) != 65535 {
			t.Error("UDINT_TO_UINT failed")
		}
	})

	t.Run("ULINT Conversions", func(t *testing.T) {
		if ULINT_TO_DATE(5000) != DATE(time.UnixMilli(5000)) {
			t.Error("ULINT_TO_DATE failed")
		}
		if ULINT_TO_DT(5000) != DT(time.UnixMilli(5000)) {
			t.Error("ULINT_TO_DT failed")
		}
		if ULINT_TO_TOD(5000) != TOD(time.Time{}.Add(5000*time.Millisecond)) {
			t.Error("ULINT_TO_TOD failed")
		}
		if ULINT_TO_SINT(100) != 100 {
			t.Error("ULINT_TO_SINT failed")
		}
		if ULINT_TO_LINT(100) != 100 {
			t.Error("ULINT_TO_LINT failed")
		}
		if ULINT_TO_DINT(100) != 100 {
			t.Error("ULINT_TO_DINT failed")
		}
		if ULINT_TO_DWORD(100) != 100 {
			t.Error("ULINT_TO_DWORD failed")
		}
		if ULINT_TO_LWORD(100) != 100 {
			t.Error("ULINT_TO_LWORD failed")
		}
		if ULINT_TO_BOOL(1) != true {
			t.Error("ULINT_TO_BOOL failed")
		}
		if ULINT_TO_REAL(100) != 100.0 {
			t.Error("ULINT_TO_REAL failed")
		}
		if ULINT_TO_LREAL(100) != 100.0 {
			t.Error("ULINT_TO_LREAL failed")
		}
		if ULINT_TO_TIME(5000) != TIME(5*time.Second) {
			t.Error("ULINT_TO_TIME failed")
		}
		if ULINT_TO_SINT(100) != 100 {
			t.Error("ULINT_TO_SINT failed")
		}
		if ULINT_TO_LINT(100) != 100 {
			t.Error("ULINT_TO_LINT failed")
		}
		if ULINT_TO_DINT(100) != 100 {
			t.Error("ULINT_TO_DINT failed")
		}
		if ULINT_TO_DWORD(100) != 100 {
			t.Error("ULINT_TO_DWORD failed")
		}
		if ULINT_TO_LWORD(100) != 100 {
			t.Error("ULINT_TO_LWORD failed")
		}
		if ULINT_TO_BOOL(1) != true {
			t.Error("ULINT_TO_BOOL failed")
		}
		if ULINT_TO_REAL(100) != 100.0 {
			t.Error("ULINT_TO_REAL failed")
		}
		if ULINT_TO_LREAL(100) != 100.0 {
			t.Error("ULINT_TO_LREAL failed")
		}
		if ULINT_TO_TIME(5000) != TIME(5*time.Second) {
			t.Error("ULINT_TO_TIME failed")
		}
		if ULINT_TO_SINT(100) != 100 {
			t.Error("ULINT_TO_SINT failed")
		}
		if ULINT_TO_LINT(100) != 100 {
			t.Error("ULINT_TO_LINT failed")
		}
		if ULINT_TO_DINT(100) != 100 {
			t.Error("ULINT_TO_DINT failed")
		}
		if ULINT_TO_DWORD(100) != 100 {
			t.Error("ULINT_TO_DWORD failed")
		}
		if ULINT_TO_LWORD(100) != 100 {
			t.Error("ULINT_TO_LWORD failed")
		}
		if ULINT_TO_BOOL(1) != true {
			t.Error("ULINT_TO_BOOL failed")
		}
		if ULINT_TO_REAL(100) != 100.0 {
			t.Error("ULINT_TO_REAL failed")
		}
		if ULINT_TO_LREAL(100) != 100.0 {
			t.Error("ULINT_TO_LREAL failed")
		}
		if ULINT_TO_TIME(5000) != TIME(5*time.Second) {
			t.Error("ULINT_TO_TIME failed")
		}
		if ULINT_TO_STRING(1234567890) != "1234567890" {
			t.Error("ULINT_TO_STRING failed")
		}
		if ULINT_TO_UDINT(100) != 100 {
			t.Error("ULINT_TO_UDINT failed")
		}
		if ULINT_TO_WORD(70000) != 65535 {
			t.Error("ULINT_TO_WORD failed")
		}
		if ULINT_TO_BYTE(300) != 255 {
			t.Error("ULINT_TO_BYTE failed")
		}
		if ULINT_TO_USINT(300) != 255 {
			t.Error("ULINT_TO_USINT failed")
		}
		if ULINT_TO_INT(100) != 100 {
			t.Error("ULINT_TO_INT failed")
		}
	})

	t.Run("TIME Conversions", func(t *testing.T) {
		tm := TIME(5 * time.Second)
		if TIME_TO_LINT(tm) != 5000 {
			t.Error("TIME_TO_LINT conversion failed")
		}
		if TIME_TO_LREAL(tm) != 5000.0 {
			t.Error("TIME_TO_LREAL conversion failed")
		}
		// The string representation is "T#5s"
		if TIME_TO_STRING(tm) != "T#5s" {
			t.Errorf("TIME_TO_STRING conversion failed, got %s", TIME_TO_STRING(tm))
		}
		if TIME_TO_SINT(tm) != 127 {
			t.Error("TIME_TO_SINT failed")
		}
		if TIME_TO_BYTE(tm) != 255 {
			t.Error("TIME_TO_BYTE failed")
		}
		if TIME_TO_WORD(tm) != 5000 {
			t.Error("TIME_TO_WORD failed")
		}
		if TIME_TO_DWORD(tm) != 5000 {
			t.Error("TIME_TO_DWORD failed")
		}
		if TIME_TO_LWORD(tm) != 5000 {
			t.Error("TIME_TO_LWORD failed")
		}
		if TIME_TO_DINT(tm) != 5000 {
			t.Error("TIME_TO_DINT failed")
		}
		if TIME_TO_UDINT(tm) != 5000 {
			t.Error("TIME_TO_UDINT failed")
		}
		if TIME_TO_ULINT(tm) != 5000 {
			t.Error("TIME_TO_ULINT failed")
		}
		if TIME_TO_UINT(tm) != 5000 {
			t.Error("TIME_TO_UINT failed")
		}
		if TIME_TO_USINT(tm) != 255 {
			t.Error("TIME_TO_USINT failed")
		}
		if TIME_TO_INT(tm) != 5000 {
			t.Error("TIME_TO_INT failed")
		}
	})

	t.Run("DATE Conversions", func(t *testing.T) {
		d := DATE(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
		expectedMillis := LINT(time.Time(d).UnixMilli())

		if DATE_TO_LINT(d) != expectedMillis {
			t.Error("DATE_TO_LINT conversion failed")
		}
		if DATE_TO_DT(d) != DT(time.Time(d)) {
			t.Error("DATE_TO_DT conversion failed")
		}
		if DATE_TO_SINT(d) != 127 {
			t.Error("DATE_TO_SINT failed")
		}
		if DATE_TO_BYTE(d) != 255 {
			t.Error("DATE_TO_BYTE failed")
		}
		if DATE_TO_WORD(d) != 65535 {
			t.Error("DATE_TO_WORD failed")
		}
		if DATE_TO_DWORD(d) != 4294967295 {
			t.Error("DATE_TO_DWORD failed")
		}
		if DATE_TO_LWORD(d) != 18446744073709551615 {
			t.Error("DATE_TO_LWORD failed")
		}
		if DATE_TO_TIME(d) != TIME(expectedMillis*1000000) {
			t.Error("DATE_TO_TIME failed")
		}
		if DATE_TO_REAL(d) != REAL(expectedMillis) {
			t.Error("DATE_TO_REAL failed")
		}
		if DATE_TO_DINT(d) != MAXDINT {
			t.Error("DATE_TO_DINT failed")
		}
		if DATE_TO_UDINT(d) != MAXUDINT {
			t.Error("DATE_TO_UDINT failed")
		}
		if DATE_TO_UINT(d) != MAXUINT {
			t.Error("DATE_TO_UINT failed")
		}
		if DATE_TO_LREAL(d) != LREAL(expectedMillis) {
			t.Error("DATE_TO_LREAL failed")
		}
		if DATE_TO_USINT(d) != MAXUSINT {
			t.Error("DATE_TO_USINT failed")
		}
		if DATE_TO_INT(d) != MAXINT {
			t.Error("DATE_TO_INT failed")
		}
		if DATE_TO_ULINT(d) != ULINT(expectedMillis) {
			t.Error("DATE_TO_ULINT failed")
		}
	})

	t.Run("TOD Conversions", func(t *testing.T) {
		tod := TOD(time.Date(0, 0, 0, 10, 20, 30, 0, time.UTC))
		expectedMillis := LINT((10*time.Hour + 20*time.Minute + 30*time.Second).Milliseconds())

		if TOD_TO_LINT(tod) != expectedMillis {
			t.Errorf("TOD_TO_LINT conversion failed, got %d, want %d", TOD_TO_LINT(tod), expectedMillis)
		}
		if TOD_TO_STRING(tod) != "TOD#10:20:30.000" {
			t.Errorf("TOD_TO_STRING conversion failed, got %s", TOD_TO_STRING(tod))
		}
		if TOD_TO_SINT(tod) != 127 {
			t.Error("TOD_TO_SINT failed")
		}
		if TOD_TO_BYTE(tod) != 255 {
			t.Error("TOD_TO_BYTE failed")
		}
		if TOD_TO_WORD(tod) != 65535 {
			t.Error("TOD_TO_WORD failed")
		}
		if TOD_TO_DWORD(tod) != 37230000 {
			t.Error("TOD_TO_DWORD failed")
		}
		if TOD_TO_LWORD(tod) != 37230000 {
			t.Error("TOD_TO_LWORD failed")
		}
		if TOD_TO_REAL(tod) != REAL(expectedMillis) {
			t.Error("TOD_TO_REAL failed")
		}
		if TOD_TO_UDINT(tod) != UDINT(expectedMillis) {
			t.Error("TOD_TO_UDINT failed")
		}
		if TOD_TO_UINT(tod) != UINT(expectedMillis) {
			t.Error("TOD_TO_UINT failed")
		}
		if TOD_TO_LREAL(tod) != LREAL(expectedMillis) {
			t.Error("TOD_TO_LREAL failed")
		}
		if TOD_TO_INT(tod) != INT(expectedMillis) {
			t.Error("TOD_TO_INT failed")
		}
		if TOD_TO_ULINT(tod) != ULINT(expectedMillis) {
			t.Error("TOD_TO_ULINT failed")
		}
		if TOD_TO_DINT(tod) != DINT(expectedMillis) {
			t.Error("TOD_TO_DINT failed")
		}
		if TOD_TO_USINT(tod) != USINT(expectedMillis) {
			t.Error("TOD_TO_USINT failed")
		}
	})

	t.Run("DT Conversions", func(t *testing.T) {
		dt := DT(time.Date(2025, 1, 1, 10, 20, 30, 0, time.UTC))
		expectedMillis := LINT(time.Time(dt).UnixMilli())
		expectedDate := DATE(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

		if DT_TO_DATE(dt) != expectedDate {
			t.Error("DT_TO_DATE conversion failed")
		}
		if DT_TO_TOD(dt) != TOD(time.Time(dt)) {
			t.Error("DT_TO_TOD conversion failed")
		}
		if DT_TO_SINT(dt) != 127 {
			t.Error("DT_TO_SINT failed")
		}
		if DT_TO_BYTE(dt) != 255 {
			t.Error("DT_TO_BYTE failed")
		}
		if DT_TO_WORD(dt) != 65535 {
			t.Error("DT_TO_WORD failed")
		}
		if DT_TO_DWORD(dt) != 4294967295 {
			t.Error("DT_TO_DWORD failed")
		}
		if DT_TO_LWORD(dt) != 18446744073709551615 {
			t.Error("DT_TO_LWORD failed")
		}
		if DT_TO_REAL(dt) != REAL(expectedMillis) {
			t.Error("DT_TO_REAL failed")
		}
		if DT_TO_LINT(dt) != LINT(expectedMillis) {
			t.Error("DT_TO_LINT failed")
		}
		if DT_TO_DINT(dt) != MAXDINT {
			t.Error("DT_TO_DINT failed")
		}
		if DT_TO_UDINT(dt) != MAXUDINT {
			t.Error("DT_TO_UDINT failed")
		}
		if DT_TO_UINT(dt) != MAXUINT {
			t.Error("DT_TO_UINT failed")
		}
		if DT_TO_LREAL(dt) != LREAL(expectedMillis) {
			t.Error("DT_TO_LREAL failed")
		}
		if DT_TO_INT(dt) != MAXINT {
			t.Error("DT_TO_INT failed")
		}
		if DT_TO_ULINT(dt) != ULINT(expectedMillis) {
			t.Error("DT_TO_ULINT failed")
		}
		if DT_TO_USINT(DT(time.UnixMilli(100))) != 100 {
			t.Error("DT_TO_USINT failed")
		}
		if DT_TO_STRING(dt) != "DT#2025-01-01-10:20:30" {
			t.Errorf("DT_TO_STRING failed, got %s", DT_TO_STRING(dt))
		}
	})
}

func TestBCDConversionsExtended(t *testing.T) {
	t.Run("UINT_TO_BCD_WORD", func(t *testing.T) {
		val, err := UINT_TO_BCD_WORD(1234)
		if err != nil || val != 0x1234 {
			t.Errorf("UINT_TO_BCD_WORD failed, got 0x%X, err: %v", val, err)
		}
	})

	t.Run("UDINT_TO_BCD_DWORD", func(t *testing.T) {
		val, err := UDINT_TO_BCD_DWORD(12345678)
		if err != nil || val != 0x12345678 {
			t.Errorf("UDINT_TO_BCD_DWORD failed, got 0x%X, err: %v", val, err)
		}
	})

	t.Run("ULINT_TO_BCD_LWORD", func(t *testing.T) {
		val, err := ULINT_TO_BCD_LWORD(1234567890123456)
		if err != nil || val != 0x1234567890123456 {
			t.Errorf("ULINT_TO_BCD_LWORD failed, got 0x%X, err: %v", val, err)
		}
	})

	t.Run("DWORD_BCD_TO_UDINT", func(t *testing.T) {
		val, err := DWORD_BCD_TO_UDINT(0x12345678)
		if err != nil || val != 12345678 {
			t.Errorf("DWORD_BCD_TO_UDINT failed, got %d, err: %v", val, err)
		}
	})

	t.Run("LWORD_BCD_TO_ULINT", func(t *testing.T) {
		val, err := LWORD_BCD_TO_ULINT(0x1234567890123456)
		if err != nil || val != 1234567890123456 {
			t.Errorf("LWORD_BCD_TO_ULINT failed, got %d, err: %v", val, err)
		}
	})
}
