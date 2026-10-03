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

package iec

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// DATE, TIME_OF_DAY (TOD), DT and TIMESPEC are defined on time.Time, so they do
// not inherit its encoding methods. Without the methods below, encoding/json
// writes them as {} and encoding/gob rejects them. This file gives them:
//
//   - Text and JSON encoding in ISO 8601 form:
//     DT and TIMESPEC  "2026-10-01T14:30:05Z" (RFC 3339, fractional seconds only when non-zero)
//     DATE             "2026-10-01"
//     TOD              "14:30:05" (fractional seconds only when non-zero)
//     Decoding also accepts the IEC 61131-3 literal forms, such as
//     "DT#2026-10-01-14:30:05", "D#2026-10-01" and "TOD#14:30:05.5".
//   - Binary encoding (used by encoding/gob), which round-trips exactly,
//     including the time zone offset.
//
// TIME is a duration based on int64 and is encoded as a number of nanoseconds.

const (
	dateLayout = "2006-01-02"
	todLayout  = "15:04:05.999999999"
	// iecDTLayout is the date-time part of an IEC literal, e.g. DT#2026-10-01-14:30:05.5.
	iecDTLayout = "2006-01-02-15:04:05.999999999"
	// localDTLayout is an ISO date-time without a zone, read as UTC.
	localDTLayout = "2006-01-02T15:04:05.999999999"
)

// unmarshalJSONText decodes a JSON value for a time type whose text form is
// handled by fromText. It accepts:
//   - a JSON string, decoded by fromText;
//   - null, which leaves the value unchanged, as encoding/json does for other types;
//   - an empty object {}, also leaving the value unchanged. Releases up to and
//     including v0.1.0-beta1 had no encoding for these types and wrote them as {},
//     so data saved by those versions (for example honeycomb's persisted tags)
//     keeps loading instead of failing.
func unmarshalJSONText(typeName string, data []byte, fromText func([]byte) error) error {
	s := strings.TrimSpace(string(data))
	if s == "null" {
		return nil
	}
	if len(s) >= 2 && s[0] == '{' && strings.TrimSpace(s[1:len(s)-1]) == "" && s[len(s)-1] == '}' {
		return nil
	}
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return fmt.Errorf("iec: invalid %s %s: want a JSON string", typeName, s)
	}
	return fromText([]byte(text))
}

// trimPrefix removes the first matching IEC type prefix, ignoring case.
func trimPrefix(s string, prefixes ...string) (string, bool) {
	for _, p := range prefixes {
		if len(s) >= len(p) && strings.EqualFold(s[:len(p)], p) {
			return s[len(p):], true
		}
	}
	return s, false
}

// parseDateTime parses an RFC 3339 date-time, an ISO date-time without a zone
// (as UTC), or an IEC DT literal.
func parseDateTime(typeName string, text []byte) (time.Time, error) {
	s := strings.TrimSpace(string(text))
	if rest, ok := trimPrefix(s, "DATE_AND_TIME#", "DT#"); ok {
		t, err := time.Parse(iecDTLayout, rest)
		if err != nil {
			return time.Time{}, fmt.Errorf("iec: invalid %s literal %q: %w", typeName, s, err)
		}
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse(localDTLayout, s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("iec: invalid %s %q: want RFC 3339 such as 2026-10-01T14:30:05Z or DT#2026-10-01-14:30:05", typeName, s)
}

// --- DT ---

// MarshalText encodes the date and time in RFC 3339 form.
func (d DT) MarshalText() ([]byte, error) {
	return []byte(time.Time(d).Format(time.RFC3339Nano)), nil
}

// UnmarshalText decodes RFC 3339, an ISO date-time without a zone (as UTC), or
// an IEC literal such as DT#2026-10-01-14:30:05.
func (d *DT) UnmarshalText(text []byte) error {
	t, err := parseDateTime("DT", text)
	if err != nil {
		return err
	}
	*d = DT(t)
	return nil
}

// UnmarshalJSON decodes a JSON string as UnmarshalText does. It also accepts
// null and the legacy {} form, leaving the value unchanged.
func (d *DT) UnmarshalJSON(data []byte) error {
	return unmarshalJSONText("DT", data, d.UnmarshalText)
}

// MarshalBinary encodes the value exactly, including its time zone offset.
func (d DT) MarshalBinary() ([]byte, error) { return time.Time(d).MarshalBinary() }

// UnmarshalBinary decodes a value written by MarshalBinary.
func (d *DT) UnmarshalBinary(data []byte) error { return (*time.Time)(d).UnmarshalBinary(data) }

// --- TIMESPEC ---

// MarshalText encodes the instant in RFC 3339 form.
func (t TIMESPEC) MarshalText() ([]byte, error) {
	return []byte(time.Time(t).Format(time.RFC3339Nano)), nil
}

// UnmarshalText decodes the same forms as DT.UnmarshalText.
func (t *TIMESPEC) UnmarshalText(text []byte) error {
	v, err := parseDateTime("TIMESPEC", text)
	if err != nil {
		return err
	}
	*t = TIMESPEC(v)
	return nil
}

// UnmarshalJSON decodes a JSON string as UnmarshalText does. It also accepts
// null and the legacy {} form, leaving the value unchanged.
func (t *TIMESPEC) UnmarshalJSON(data []byte) error {
	return unmarshalJSONText("TIMESPEC", data, t.UnmarshalText)
}

// MarshalBinary encodes the value exactly, including its time zone offset.
func (t TIMESPEC) MarshalBinary() ([]byte, error) { return time.Time(t).MarshalBinary() }

// UnmarshalBinary decodes a value written by MarshalBinary.
func (t *TIMESPEC) UnmarshalBinary(data []byte) error {
	return (*time.Time)(t).UnmarshalBinary(data)
}

// --- DATE ---

// MarshalText encodes the calendar date as YYYY-MM-DD.
func (d DATE) MarshalText() ([]byte, error) {
	return []byte(time.Time(d).Format(dateLayout)), nil
}

// UnmarshalText decodes YYYY-MM-DD or an IEC literal such as D#2026-10-01. The
// result is midnight UTC on that date.
func (d *DATE) UnmarshalText(text []byte) error {
	s := strings.TrimSpace(string(text))
	rest, _ := trimPrefix(s, "DATE#", "D#")
	t, err := time.Parse(dateLayout, rest)
	if err != nil {
		return fmt.Errorf("iec: invalid DATE %q: want 2026-10-01 or D#2026-10-01", s)
	}
	*d = DATE(t)
	return nil
}

// UnmarshalJSON decodes a JSON string as UnmarshalText does. It also accepts
// null and the legacy {} form, leaving the value unchanged.
func (d *DATE) UnmarshalJSON(data []byte) error {
	return unmarshalJSONText("DATE", data, d.UnmarshalText)
}

// MarshalBinary encodes the value exactly, including its time zone offset.
func (d DATE) MarshalBinary() ([]byte, error) { return time.Time(d).MarshalBinary() }

// UnmarshalBinary decodes a value written by MarshalBinary.
func (d *DATE) UnmarshalBinary(data []byte) error { return (*time.Time)(d).UnmarshalBinary(data) }

// --- TIME_OF_DAY (TOD) ---

// MarshalText encodes the clock time as hh:mm:ss, with fractional seconds only
// when they are non-zero.
func (t TIME_OF_DAY) MarshalText() ([]byte, error) {
	return []byte(time.Time(t).Format(todLayout)), nil
}

// UnmarshalText decodes hh:mm:ss[.fraction] or an IEC literal such as
// TOD#14:30:05.5. The date part is the zero date (0001-01-01 UTC), matching INITTOD.
func (t *TIME_OF_DAY) UnmarshalText(text []byte) error {
	s := strings.TrimSpace(string(text))
	rest, _ := trimPrefix(s, "TIME_OF_DAY#", "TOD#")
	clock, err := time.Parse(todLayout, rest)
	if err != nil {
		return fmt.Errorf("iec: invalid TOD %q: want 14:30:05 or TOD#14:30:05", s)
	}
	// time.Parse puts a bare clock on year 0; move it to the zero date.
	*t = TIME_OF_DAY(time.Time{}.Add(clock.Sub(time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC))))
	return nil
}

// UnmarshalJSON decodes a JSON string as UnmarshalText does. It also accepts
// null and the legacy {} form, leaving the value unchanged.
func (t *TIME_OF_DAY) UnmarshalJSON(data []byte) error {
	return unmarshalJSONText("TOD", data, t.UnmarshalText)
}

// MarshalBinary encodes the value exactly, including its date and time zone offset.
func (t TIME_OF_DAY) MarshalBinary() ([]byte, error) { return time.Time(t).MarshalBinary() }

// UnmarshalBinary decodes a value written by MarshalBinary.
func (t *TIME_OF_DAY) UnmarshalBinary(data []byte) error {
	return (*time.Time)(t).UnmarshalBinary(data)
}
