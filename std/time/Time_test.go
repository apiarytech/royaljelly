package time

import (
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

func TestTimePackageFunctions(t *testing.T) {
	now := time.Now()
	ts1 := iec.TIMESPEC(now)

	t.Run("Global Functions", func(t *testing.T) {
		// NOW()
		nowResult := NOW()
		// Check if it's approximately now (within a small delta)
		if time.Since(time.Time(nowResult)) > 50*time.Millisecond {
			t.Errorf("NOW() result %v is not approximately current time %v", nowResult, time.Now())
		}

		// UTC()
		utcTime := UTC(ts1)
		if time.Time(utcTime).Location() != time.UTC {
			t.Errorf("UTC() did not convert to UTC. Got %v, want UTC", time.Time(utcTime).Location())
		}

		// LOCAL()
		localTime := LOCAL(utcTime) // Convert UTC back to local
		if time.Time(localTime).Location() != time.Local {
			t.Errorf("LOCAL() did not convert to Local. Got %v, want Local", time.Time(localTime).Location())
		}
		// To make this test robust against DST changes during execution,
		// we build the expected local time from the UTC time directly.
		expectedLocal := time.Time(utcTime).In(time.Local)
		if !time.Time(localTime).Equal(expectedLocal) {
			t.Errorf("LOCAL() did not produce the correct local time. Got %v, want %v", time.Time(localTime), expectedLocal)
		}
	})

	t.Run("DT_TO_TM", func(t *testing.T) {
		dt := iec.DT(time.Date(2024, 3, 15, 10, 30, 45, 123456789, time.UTC))
		tm := DT_TO_TM(dt)
		if tm.D != 15 || tm.H != 10 || tm.M != 30 || tm.S != 45 || tm.Ms != 123 {
			t.Errorf("DT_TO_TM(%v) = %+v; want {d:15, h:10, m:30, s:45, ms:123}", dt, tm)
		}
	})

	t.Run("TM_TO_DT", func(t *testing.T) {
		// TM_TO_DT uses current year/month, so we need to account for that.
		tm := iec.TM{D: 1, H: 11, M: 22, S: 33, Ms: 444}
		resultTime := time.Time(TM_TO_DT(tm))
		// Build the expected time using the year and month from the *actual* result
		// to avoid race conditions if the test runs across a date boundary (e.g., midnight).
		expectedDT := time.Date(resultTime.Year(), resultTime.Month(), tm.D, tm.H, tm.M, tm.S, tm.Ms*1e6, resultTime.Location())
		if !resultTime.Equal(expectedDT) {
			t.Errorf("TM_TO_DT(%+v) = %v; want %v", tm, resultTime, expectedDT)
		}
	})

	t.Run("TOD_TO_DT", func(t *testing.T) {
		todTime := time.Date(1970, 1, 1, 14, 0, 0, 0, time.UTC)
		tod := iec.TOD(todTime)
		dt := TOD_TO_DT(tod)
		// TOD_TO_DT simply casts TOD to DT, preserving all underlying time.Time fields.
		if !time.Time(dt).Equal(time.Time(tod)) {
			t.Errorf("TOD_TO_DT(%v) = %v; want %v", tod, time.Time(dt), time.Time(tod))
		}
	})

	t.Run("DATE_TO_DT", func(t *testing.T) {
		dateTime := time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC)
		date := iec.DATE(dateTime)
		dt := DATE_TO_DT(date)
		// DATE_TO_DT simply casts DATE to DT, preserving all underlying time.Time fields.
		if !time.Time(dt).Equal(time.Time(date)) {
			t.Errorf("DATE_TO_DT(%v) = %v; want %v", date, time.Time(dt), time.Time(date))
		}
	})
}

func TestStringConversionFunctions(t *testing.T) {
	t.Run("STRING_TO_TIME", func(t *testing.T) {
		// Valid case
		result, err := STRING_TO_TIME("1m30s")
		if err != nil {
			t.Fatalf("STRING_TO_TIME valid case failed with error: %v", err)
		}
		expected := iec.TIME(90 * time.Second)
		if result != expected {
			t.Errorf("STRING_TO_TIME(%q) = %v; want %v", "1m30s", result, expected)
		}

		// Invalid case
		_, err = STRING_TO_TIME("invalid-duration")
		if err == nil {
			t.Error("STRING_TO_TIME with invalid input should have returned an error")
		}
	})

	t.Run("STRING_TO_DATE", func(t *testing.T) {
		// Valid case
		result, err := STRING_TO_DATE("2025-07-14")
		if err != nil {
			t.Fatalf("STRING_TO_DATE valid case failed with error: %v", err)
		}
		expected := iec.DATE(time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC))
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("STRING_TO_DATE(%q) = %v; want %v", "2025-07-14", result, expected)
		}

		// Invalid case
		_, err = STRING_TO_DATE("2025/07/14")
		if err == nil {
			t.Error("STRING_TO_DATE with invalid format should have returned an error")
		}
	})

	t.Run("STRING_TO_TOD", func(t *testing.T) {
		// Valid case
		result, err := STRING_TO_TOD("18:45:10")
		if err != nil {
			t.Fatalf("STRING_TO_TOD valid case failed with error: %v", err)
		}
		// The result is on the zero date, like INITTOD and decoded TODs (issue #29).
		expected := iec.NewTOD(18*time.Hour + 45*time.Minute + 10*time.Second)
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("STRING_TO_TOD(%q) = %v; want %v", "18:45:10", result, expected)
		}

		// Invalid case
		_, err = STRING_TO_TOD("6:45:10 PM")
		if err == nil {
			t.Error("STRING_TO_TOD with invalid format should have returned an error")
		}
	})

	t.Run("STRING_TO_DT", func(t *testing.T) {
		// Valid case
		result, err := STRING_TO_DT("2026-01-02-03:04:05")
		if err != nil {
			t.Fatalf("STRING_TO_DT valid case failed with error: %v", err)
		}
		expected := iec.DT(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC))
		if !time.Time(result).Equal(time.Time(expected)) {
			t.Errorf("STRING_TO_DT(%q) = %v; want %v", "2026-01-02-03:04:05", result, expected)
		}

		// Invalid case
		_, err = STRING_TO_DT("2026/01/02 03:04:05")
		if err == nil {
			t.Error("STRING_TO_DT with invalid format should have returned an error")
		}
	})
}

// TestStringRoundTripIssue30 covers issue #30: the *_TO_STRING output parses back.
func TestStringRoundTripIssue30(t *testing.T) {
	d := iec.TIME(90 * time.Minute)
	if got, err := STRING_TO_TIME(iec.STRING(d.String())); err != nil || got != d {
		t.Errorf("STRING_TO_TIME(%q) = %v, %v", d.String(), got, err)
	}
	for _, s := range []iec.STRING{"T#1H30M", "TIME#1h30m", "t#1h_30m", "1h30m"} {
		if got, err := STRING_TO_TIME(s); err != nil || got != d {
			t.Errorf("STRING_TO_TIME(%q) = %v, %v", s, got, err)
		}
	}
	date := iec.DATE(time.Date(2026, 10, 22, 0, 0, 0, 0, time.UTC))
	if got, err := STRING_TO_DATE(iec.STRING(date.String())); err != nil || !time.Time(got).Equal(time.Time(date)) {
		t.Errorf("STRING_TO_DATE(%q) = %v, %v", date.String(), time.Time(got), err)
	}
	if _, err := STRING_TO_DATE("DATE#2026-10-22"); err != nil {
		t.Errorf("STRING_TO_DATE(DATE#...) failed: %v", err)
	}
	dt := iec.DT(time.Date(2026, 10, 22, 14, 30, 5, 0, time.UTC))
	if got, err := STRING_TO_DT(iec.STRING(dt.String())); err != nil || !time.Time(got).Equal(time.Time(dt)) {
		t.Errorf("STRING_TO_DT(%q) = %v, %v", dt.String(), time.Time(got), err)
	}
	if _, err := STRING_TO_DT("DATE_AND_TIME#2026-10-22-14:30:05"); err != nil {
		t.Errorf("STRING_TO_DT(DATE_AND_TIME#...) failed: %v", err)
	}
	tod := iec.NewTOD(14*time.Hour + 30*time.Minute + 5*time.Second)
	if got, err := STRING_TO_TOD(iec.STRING(tod.String())); err != nil || got != tod {
		t.Errorf("STRING_TO_TOD(%q) = %v, %v", tod.String(), time.Time(got), err)
	}
}
