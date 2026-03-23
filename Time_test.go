package royaljelly

import (
	"testing"
	"time"
)

func TestTIMESPEC_Methods(t *testing.T) {
	// Create two TIMESPEC instances for comparison
	now := time.Now()
	ts1 := TIMESPEC(now)
	ts2 := TIMESPEC(now.Add(time.Second))

	t.Run("Comparison", func(t *testing.T) {
		if !ts1.BEFORE(ts2) {
			t.Error("ts1.BEFORE(ts2) should be true")
		}
		if ts1.AFTER(ts2) {
			t.Error("ts1.AFTER(ts2) should be false")
		}
		if !ts1.EQUAL(TIMESPEC(now)) {
			t.Error("ts1.EQUAL(ts1) should be true")
		}
	})

	t.Run("MONTH_STRING", func(t *testing.T) {
		expected := STRING(now.Month().String())
		if ts1.MONTH_STRING() != expected {
			t.Errorf("MONTH_STRING() = %q; want %q", ts1.MONTH_STRING(), expected)
		}
	})

	t.Run("WEEKDAY_STRING", func(t *testing.T) {
		expected := STRING(now.Weekday().String())
		if ts1.WEEKDAY_STRING() != expected {
			t.Errorf("WEEKDAY_STRING() = %q; want %q", ts1.WEEKDAY_STRING(), expected)
		}
	})

	t.Run("Extraction", func(t *testing.T) {
		if ts1.YEAR() != LINT(now.Year()) {
			t.Errorf("YEAR() mismatch. Got %d, want %d", ts1.YEAR(), now.Year())
		}
		if ts1.MONTH() != LINT(now.Month()) {
			t.Errorf("MONTH() mismatch. Got %d, want %d", ts1.MONTH(), now.Month())
		}
		if ts1.DAY() != LINT(now.Day()) {
			t.Errorf("DAY() mismatch. Got %d, want %d", ts1.DAY(), now.Day())
		}
		if ts1.WEEKDAY() != LINT(now.Weekday()) {
			t.Errorf("WEEKDAY() mismatch. Got %d, want %d", ts1.WEEKDAY(), now.Weekday())
		}
		if ts1.HOUR() != LINT(now.Hour()) {
			t.Errorf("HOUR() mismatch. Got %d, want %d", ts1.HOUR(), now.Hour())
		}
		if ts1.MINUTE() != LINT(now.Minute()) {
			t.Errorf("MINUTE() mismatch. Got %d, want %d", ts1.MINUTE(), now.Minute())
		}
		// Millisecond extraction from Nanosecond might have precision issues, check within a range
		expectedMs := LINT(now.Nanosecond() / 1e6)
		if ts1.MILLISECOND() != expectedMs {
			t.Errorf("MILLISECOND() mismatch. Got %d, want %d", ts1.MILLISECOND(), expectedMs)
		}
		if ts1.SECOND() != LINT(now.Second()) {
			t.Errorf("SECOND() mismatch. Got %d, want %d", ts1.SECOND(), now.Second())
		}
		year, week := ts1.ISOWEEK()
		expectedYear, expectedWeek := now.ISOWeek()
		if year != LINT(expectedYear) || week != LINT(expectedWeek) {
			t.Errorf("ISOWEEK() mismatch. Got year %d, week %d; want year %d, week %d", year, week, expectedYear, expectedWeek)
		}
	})

	t.Run("Conversion", func(t *testing.T) {
		// Test ISZERO
		zeroTime := TIMESPEC(time.Time{})
		if !zeroTime.ISZERO() {
			t.Error("ISZERO() should be true for zero time")
		}
		if ts1.ISZERO() {
			t.Error("ISZERO() should be false for non-zero time")
		}

		// Test TIME() - duration since midnight
		expectedTime := TIME(now.Sub(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())))
		if ts1.TIME() != expectedTime {
			t.Errorf("TIME() conversion failed. Got %v, want %v", ts1.TIME(), expectedTime)
		}

		// Test DATE()
		expectedDate := DATE(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()))
		if ts1.DATE() != expectedDate {
			t.Errorf("DATE() conversion failed. Got %v, want %v", ts1.DATE(), expectedDate)
		}

		if ts1.DATETIME() != DT(now) {
			t.Error("DATETIME() conversion failed")
		}
		if ts1.CLOCK() != TOD(now) {
			t.Error("CLOCK() conversion failed")
		}
	})

	t.Run("Global Functions", func(t *testing.T) {
		// NOW()
		nowResult := NOW()
		// Check if it's approximately now (within a small delta)
		if time.Since(time.Time(nowResult)) > 10*time.Millisecond || time.Since(time.Time(nowResult)) < -10*time.Millisecond {
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
		// Check if the time value is preserved (ignoring location for comparison)
		if time.Time(localTime).Year() != time.Time(ts1).Year() || time.Time(localTime).Month() != time.Time(ts1).Month() {
			t.Errorf("LOCAL() did not preserve time value. Got %v, want %v", time.Time(localTime), time.Time(ts1))
		}
	})

	t.Run("DT_TO_TM", func(t *testing.T) {
		dt := DT(time.Date(2024, 3, 15, 10, 30, 45, 123456789, time.UTC))
		tm := DT_TO_TM(dt)
		if tm.d != 15 || tm.h != 10 || tm.m != 30 || tm.s != 45 || tm.ms != 123 {
			t.Errorf("DT_TO_TM(%v) = %+v; want {d:15, h:10, m:30, s:45, ms:123}", dt, tm)
		}
	})

	t.Run("TM_TO_DT", func(t *testing.T) {
		// TM_TO_DT uses current year/month, so we need to account for that.
		now := time.Now()
		tm := TM{d: 1, h: 11, m: 22, s: 33, ms: 444}
		dt := TM_TO_DT(tm)
		expectedDT := time.Date(now.Year(), now.Month(), tm.d, tm.h, tm.m, tm.s, tm.ms*1e6, now.Location())
		if !time.Time(dt).Equal(expectedDT) {
			t.Errorf("TM_TO_DT(%+v) = %v; want %v", tm, time.Time(dt), expectedDT)
		}
	})

	t.Run("TOD_TO_DT", func(t *testing.T) {
		tod := TOD(time.Date(0, 0, 0, 14, 0, 0, 0, time.UTC))
		dt := TOD_TO_DT(tod)
		// TOD_TO_DT simply casts TOD to DT, preserving all underlying time.Time fields.
		// The date components of TOD are usually ignored in its interpretation, but present in the underlying time.Time.
		// So, the DT should be identical to the TOD's underlying time.Time.
		if !time.Time(dt).Equal(time.Time(tod)) {
			t.Errorf("TOD_TO_DT(%v) = %v; want %v", tod, time.Time(dt), time.Time(tod))
		}
	})

	t.Run("DATE_TO_DT", func(t *testing.T) {
		date := DATE(time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC))
		dt := DATE_TO_DT(date)
		// DATE_TO_DT simply casts DATE to DT, preserving all underlying time.Time fields.
		// The time components of DATE are usually zeroed out, but present in the underlying time.Time.
		// So, the DT should be identical to the DATE's underlying time.Time.
		if !time.Time(dt).Equal(time.Time(date)) {
			t.Errorf("DATE_TO_DT(%v) = %v; want %v", date, time.Time(dt), time.Time(date))
		}
	})
}
