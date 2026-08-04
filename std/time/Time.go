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

package time

import (
	"time"

	. "github.com/apiarytech/royaljelly/core"
)

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
		return 0, &ConversionError{
			Value:    in,
			FromType: "STRING",
			ToType:   "TIME",
			Reason:   "string could not be parsed as a duration (e.g., '1h30m15s')",
			Err:      err,
		}
	}
	return TIME(d), nil
}

// STRING_TO_DATE converts a string representation (e.g., "2026-03-22") into a DATE.
func STRING_TO_DATE(in STRING) (DATE, error) {
	t, err := time.Parse("2006-01-02", string(in))
	if err != nil {
		return DATE(time.Time{}), &ConversionError{
			Value:    in,
			FromType: "STRING",
			ToType:   "DATE",
			Reason:   "string could not be parsed as a date (format 'YYYY-MM-DD')",
			Err:      err,
		}
	}
	return DATE(t), nil
}

// STRING_TO_TOD converts a string representation (e.g., "15:04:05") into a TIME_OF_DAY.
func STRING_TO_TOD(in STRING) (TOD, error) {
	// We parse it against a known date, then the date part is ignored by the TOD type's usage.
	t, err := time.Parse("2006-01-02 15:04:05", "1970-01-01 "+string(in))
	if err != nil {
		return TOD(time.Time{}), &ConversionError{
			Value:    in,
			FromType: "STRING",
			ToType:   "TOD",
			Reason:   "string could not be parsed as a time of day (format 'HH:MM:SS')",
			Err:      err,
		}
	}
	return TOD(t), nil
}

// STRING_TO_DT converts a string representation (e.g., "2026-03-22-15:04:05") into a DATE_AND_TIME.
func STRING_TO_DT(in STRING) (DT, error) {
	t, err := time.Parse("2006-01-02-15:04:05", string(in))
	if err != nil {
		return DT(time.Time{}), &ConversionError{
			Value:    in,
			FromType: "STRING",
			ToType:   "DT",
			Reason:   "string could not be parsed as a date and time (format 'YYYY-MM-DD-HH:MM:SS')",
			Err:      err,
		}
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
		D:  t.Day(),
		H:  t.Hour(),
		M:  t.Minute(),
		S:  t.Second(),
		Ms: t.Nanosecond() / 1e6,
	}
}

// TM_TO_DT converts a TM struct into a DT (DATE_AND_TIME).
// It uses the current year and month, which is a common approach when only time components are provided.
func TM_TO_DT(in TM) DT {
	now := time.Now()
	return DT(time.Date(now.Year(), now.Month(), in.D, in.H, in.M, in.S, in.Ms*1e6, now.Location()))
}

func TOD_TO_DT(in TOD) DT {
	return DT(time.Time(in))
}

func DATE_TO_DT(in DATE) DT {
	return DT(time.Time(in))
}
