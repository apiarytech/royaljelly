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

	"github.com/apiarytech/royaljelly/convert"
	"github.com/apiarytech/royaljelly/iec"
)

func NOW() iec.TIMESPEC {
	return iec.TIMESPEC(time.Now())
}

func UTC(t iec.TIMESPEC) iec.TIMESPEC {
	return iec.TIMESPEC(time.Time(t).UTC())
}

func LOCAL(t iec.TIMESPEC) iec.TIMESPEC {
	return iec.TIMESPEC(time.Time(t).Local())
}

/*****************************************************************/
/* Non-Standard but useful Time Conversion Functions             */
/*****************************************************************/

// STRING_TO_TIME converts a string representation into a TIME duration.
// It expects a format compatible with Go's time.ParseDuration (e.g., "1h30m15s").
func STRING_TO_TIME(in iec.STRING) (iec.TIME, error) {
	d, err := time.ParseDuration(string(in))
	if err != nil {
		return 0, &convert.ConversionError{
			Value:    in,
			FromType: "STRING",
			ToType:   "TIME",
			Reason:   "string could not be parsed as a duration (e.g., '1h30m15s')",
			Err:      err,
		}
	}
	return iec.TIME(d), nil
}

// STRING_TO_DATE converts a string representation (e.g., "2026-03-22") into a DATE.
func STRING_TO_DATE(in iec.STRING) (iec.DATE, error) {
	t, err := time.Parse("2006-01-02", string(in))
	if err != nil {
		return iec.DATE(time.Time{}), &convert.ConversionError{
			Value:    in,
			FromType: "STRING",
			ToType:   "DATE",
			Reason:   "string could not be parsed as a date (format 'YYYY-MM-DD')",
			Err:      err,
		}
	}
	return iec.DATE(t), nil
}

// STRING_TO_TOD converts a string representation (e.g., "15:04:05") into a TIME_OF_DAY.
func STRING_TO_TOD(in iec.STRING) (iec.TOD, error) {
	// We parse it against a known date, then the date part is ignored by the TOD type's usage.
	t, err := time.Parse("2006-01-02 15:04:05", "1970-01-01 "+string(in))
	if err != nil {
		return iec.TOD(time.Time{}), &convert.ConversionError{
			Value:    in,
			FromType: "STRING",
			ToType:   "TOD",
			Reason:   "string could not be parsed as a time of day (format 'HH:MM:SS')",
			Err:      err,
		}
	}
	return iec.TOD(t), nil
}

// STRING_TO_DT converts a string representation (e.g., "2026-03-22-15:04:05") into a DATE_AND_TIME.
func STRING_TO_DT(in iec.STRING) (iec.DT, error) {
	t, err := time.Parse("2006-01-02-15:04:05", string(in))
	if err != nil {
		return iec.DT(time.Time{}), &convert.ConversionError{
			Value:    in,
			FromType: "STRING",
			ToType:   "DT",
			Reason:   "string could not be parsed as a date and time (format 'YYYY-MM-DD-HH:MM:SS')",
			Err:      err,
		}
	}
	return iec.DT(t), nil
}

/*
TO_DT and other conversions
*/

// DT_TO_TM extracts the components of a DT into a TM struct.
func DT_TO_TM(in iec.DT) iec.TM {
	t := time.Time(in)
	return iec.TM{
		D:  t.Day(),
		H:  t.Hour(),
		M:  t.Minute(),
		S:  t.Second(),
		Ms: t.Nanosecond() / 1e6,
	}
}

// TM_TO_DT converts a TM struct into a DT (DATE_AND_TIME).
// It uses the current year and month, which is a common approach when only time components are provided.
func TM_TO_DT(in iec.TM) iec.DT {
	now := time.Now()
	return iec.DT(time.Date(now.Year(), now.Month(), in.D, in.H, in.M, in.S, in.Ms*1e6, now.Location()))
}

func TOD_TO_DT(in iec.TOD) iec.DT {
	return iec.DT(time.Time(in))
}

func DATE_TO_DT(in iec.DATE) iec.DT {
	return iec.DT(time.Time(in))
}
