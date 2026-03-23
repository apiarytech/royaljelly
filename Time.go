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

package royaljelly

import "time"

func (t *TIMESPEC) AFTER(u TIMESPEC) BOOL {
	return BOOL(time.Time(*t).After(time.Time(u)))
}

func (t *TIMESPEC) BEFORE(u TIMESPEC) BOOL {
	return BOOL(time.Time(*t).Before(time.Time(u)))
}

func (t *TIMESPEC) EQUAL(u TIMESPEC) BOOL {
	return BOOL(time.Time(*t).Equal(time.Time(u)))
}

func (t *TIMESPEC) MONTH_STRING() STRING {
	return STRING(time.Time(*t).Month().String())
}

func (t *TIMESPEC) WEEKDAY_STRING() STRING {
	return STRING(time.Time(*t).Weekday().String())
}

func (t *TIMESPEC) ISZERO() BOOL {
	return BOOL(time.Time(*t).IsZero())
}

func (t *TIMESPEC) TIME() (tm TIME) {
	// This conversion is conceptually tricky. IEC TIME is a duration.
	// A TIMESPEC is a point in time. A common interpretation is the duration since midnight.
	t_time := time.Time(*t)
	midnight := time.Date(t_time.Year(), t_time.Month(), t_time.Day(), 0, 0, 0, 0, t_time.Location())
	return TIME(t_time.Sub(midnight))
}

func (t *TIMESPEC) DATE() (d DATE) {
	t_time := time.Time(*t)
	return DATE(time.Date(t_time.Year(), t_time.Month(), t_time.Day(), 0, 0, 0, 0, t_time.Location()))
}

func (t *TIMESPEC) CLOCK() (tod TOD) {
	return TOD(time.Time(*t))
}

func (t *TIMESPEC) DATETIME() (dt DT) {
	return DT(time.Time(*t))
}

func (t *TIMESPEC) YEAR() LINT {
	return LINT(time.Time(*t).Year())
}

func (t *TIMESPEC) MONTH() LINT {
	return LINT(time.Time(*t).Month())
}

func (t *TIMESPEC) HOUR() LINT {
	return LINT(time.Time(*t).Hour())
}

func (t *TIMESPEC) MINUTE() LINT {
	return LINT(time.Time(*t).Minute())
}

func (t *TIMESPEC) MILLISECOND() LINT {
	return LINT(time.Time(*t).Nanosecond() / 1e6)
}

func (t *TIMESPEC) SECOND() LINT {
	return LINT(time.Time(*t).Second())
}

func (t *TIMESPEC) DAY() LINT {
	return LINT(time.Time(*t).Day())
}

func (t *TIMESPEC) ISOWEEK() (year, week LINT) {
	y, w := time.Time(*t).ISOWeek()
	return LINT(y), LINT(w)
}

func (t *TIMESPEC) WEEKDAY() LINT {
	return LINT(time.Time(*t).Weekday())
}

func NOW() TIMESPEC {
	return TIMESPEC(time.Now())
}

func UTC(t TIMESPEC) TIMESPEC {
	return TIMESPEC(time.Time(t).UTC())
}

func LOCAL(t TIMESPEC) TIMESPEC {
	return TIMESPEC(time.Time(t).Local())
}

/*****************************************************************/
/* Non-Standard but useful Time Conversion Functions             */
/*****************************************************************/

// STRING_TO_TIME converts a string representation into a TIME duration.
// It expects a format compatible with Go's time.ParseDuration (e.g., "1h30m15s").
func STRING_TO_TIME(in STRING) (TIME, error) {
	d, err := time.ParseDuration(string(in))
	if err != nil {
		return 0, err
	}
	return TIME(d), nil
}

// STRING_TO_DATE converts a string representation (e.g., "2026-03-22") into a DATE.
func STRING_TO_DATE(in STRING) (DATE, error) {
	t, err := time.Parse("2006-01-02", string(in))
	if err != nil {
		return DATE(time.Time{}), err
	}
	return DATE(t), nil
}

// STRING_TO_TOD converts a string representation (e.g., "15:04:05") into a TIME_OF_DAY.
func STRING_TO_TOD(in STRING) (TOD, error) {
	// We parse it against a known date, then the date part is ignored by the TOD type's usage.
	t, err := time.Parse("2006-01-02 15:04:05", "1970-01-01 "+string(in))
	if err != nil {
		return TOD(time.Time{}), err
	}
	return TOD(t), nil
}

// STRING_TO_DT converts a string representation (e.g., "2026-03-22-15:04:05") into a DATE_AND_TIME.
func STRING_TO_DT(in STRING) (DT, error) {
	t, err := time.Parse("2006-01-02-15:04:05", string(in))
	if err != nil {
		return DT(time.Time{}), err
	}
	return DT(t), nil
}

/*
TO_DT and other conversions
*/

// DT_TO_TM extracts the components of a DT into a TM struct.
func DT_TO_TM(in DT) TM {
	t := time.Time(in)
	return TM{
		d:  t.Day(),
		h:  t.Hour(),
		m:  t.Minute(),
		s:  t.Second(),
		ms: t.Nanosecond() / 1e6,
	}
}

// TM_TO_DT converts a TM struct into a DT (DATE_AND_TIME).
// It uses the current year and month, which is a common approach when only time components are provided.
func TM_TO_DT(in TM) DT {
	now := time.Now()
	return DT(time.Date(now.Year(), now.Month(), in.d, in.h, in.m, in.s, in.ms*1e6, now.Location()))
}

func TOD_TO_DT(in TOD) DT {
	return DT(time.Time(in))
}

func DATE_TO_DT(in DATE) DT {
	return DT(time.Time(in))
}
