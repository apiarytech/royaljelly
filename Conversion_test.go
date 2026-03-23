package royaljelly

import (
	"reflect"
	"testing"
	"time"
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
	// Note: String to numeric is handled by anyToLINT/anyToLREAL, not direct conversion functions
	// in the same way as other types.
	val, err := anyToLINT(STRING("123"))
	if err != nil || val != 123 {
		t.Error("anyToLINT from STRING failed")
	}

	fVal, err := anyToLREAL(STRING("-123.45"))
	if err != nil || fVal != -123.45 {
		t.Error("anyToLREAL from STRING failed")
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
		if USINT_TO_BCD_BYTE(99) != 0x99 {
			t.Errorf("USINT_TO_BCD_BYTE failed, got 0x%X", USINT_TO_BCD_BYTE(99))
		}
		if BYTE_BCD_TO_USINT(0x99) != 99 {
			t.Errorf("BYTE_BCD_TO_USINT failed, got %d, want 99", BYTE_BCD_TO_USINT(0x99))
		}
		if WORD_BCD_TO_UINT(0x1234) != 1234 {
			t.Errorf("WORD_BCD_TO_UINT failed, got %d", WORD_BCD_TO_UINT(0x1234))
		}
	})

	t.Run("Invalid BCD Nibble", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("BYTE_BCD_TO_USINT with invalid nibble did not panic")
			}
		}()
		_ = BYTE_BCD_TO_USINT(0xA0)
	})

	t.Run("Value too large for BCD", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("USINT_TO_BCD_BYTE with value > 99 did not panic")
			}
		}()
		_ = USINT_TO_BCD_BYTE(101) // Use 101 to avoid ambiguity with hex 0x100
	})
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

		if resultVal.Kind() != reflect.Int64 {
			t.Fatalf("DATE.CONVERT() returned kind %v; want reflect.Int64", resultVal.Kind())
		}

		resultLINT, ok := resultVal.Interface().(LINT)
		if !ok {
			t.Fatalf("DATE.CONVERT() did not return a LINT value")
		}
		if resultLINT != expected {
			t.Errorf("DATE.CONVERT() = %d; want %d", resultLINT, expected)
		}
	})

	t.Run("DT.CONVERT()", func(t *testing.T) {
		// DT is a point in time, CONVERT should return milliseconds since Unix epoch.
		dt := DT(time.Unix(2000, 500*1e6).UTC()) // 2000 seconds and 500 ms since epoch
		expected := LINT(2000*1000 + 500)        // in milliseconds

		resultVal := dt.CONVERT()

		if resultVal.Kind() != reflect.Int64 {
			t.Fatalf("DT.CONVERT() returned kind %v; want reflect.Int64", resultVal.Kind())
		}

		resultLINT, ok := resultVal.Interface().(LINT)
		if !ok {
			t.Fatalf("DT.CONVERT() did not return a LINT value")
		}
		if resultLINT != expected {
			t.Errorf("DT.CONVERT() = %d; want %d", resultLINT, expected)
		}
	})

	t.Run("TOD.CONVERT()", func(t *testing.T) {
		// TOD is time since midnight.
		tod := TOD(time.Date(0, 0, 0, 1, 2, 3, 456*1e6, time.UTC))
		expected := LINT((1*time.Hour + 2*time.Minute + 3*time.Second + 456*time.Millisecond).Milliseconds())

		resultVal := tod.CONVERT()

		if resultVal.Kind() != reflect.Int64 {
			t.Fatalf("TOD.CONVERT() returned kind %v; want reflect.Int64", resultVal.Kind())
		}

		resultLINT, ok := resultVal.Interface().(LINT)
		if !ok {
			t.Fatalf("TOD.CONVERT() did not return a LINT value")
		}
		if resultLINT != expected {
			t.Errorf("TOD.CONVERT() = %d; want %d", resultLINT, expected)
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
	})

	t.Run("UINT Conversions", func(t *testing.T) {
		if UINT_TO_USINT(300) != 255 {
			t.Errorf("UINT_TO_USINT(300) should be clamped to 255, got %d", UINT_TO_USINT(300))
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
	})

	t.Run("DT Conversions", func(t *testing.T) {
		dt := DT(time.Date(2025, 1, 1, 10, 20, 30, 0, time.UTC))
		expectedDate := DATE(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))

		if DT_TO_DATE(dt) != expectedDate {
			t.Error("DT_TO_DATE conversion failed")
		}
		if DT_TO_TOD(dt) != TOD(time.Time(dt)) {
			t.Error("DT_TO_TOD conversion failed")
		}
	})
}
