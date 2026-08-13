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
	"time"
)

/*********************/
/*  IEC Types defs   */
/*********************/

// TRUE bool = 1
const TRUE bool = true

// FALSE bool = 0
const FALSE bool = false

// BOOL bool definition
type BOOL bool

// BYTE bit strings
type BYTE uint8

// WORD unsigned int 16 bit
type WORD uint16

// DWORD unsigned int 32 bit
type DWORD uint32

// LWORD unsgined int 64 bit
type LWORD uint64

// SINT signed int 8 bit
type SINT int8

// INT signed int 16 bit
type INT int16

// DINT signed int 32 bit
type DINT int32

// LINT signed int 64 bit
type LINT int64

// USINT signed int 8 bit
type USINT uint8

// UINT signed int 16 bit
type UINT uint16

// UDINT signed int 32 bit
type UDINT uint32

// ULINT signed int 64 bit
type ULINT uint64

// REAL float 32 bit
type REAL float32

// LREAL float64 bit
type LREAL float64

// COMPLEX Real & Imaginary 64 bit
type COMPLEX complex64

// LCOMPLEX Real & Imaginary 128 bit
type LCOMPLEX complex128

// CHAR definition
type CHAR byte

// WCHAR definition
type WCHAR rune

// STRING definition
type STRING string

// STRINGS definition
type STRINGS []string

// WSTRING definition
type WSTRING rune

// WSTRINGS definition
type WSTRINGS []rune

// TIME represents a duration as defined by IEC 61131-3. It is based on Go's time.Duration for easier manipulation.
type TIME time.Duration

// DATE represents a date as defined by IEC 61131-3.
type DATE time.Time

// TIME_OF_DAY (TOD) represents a time of day as defined by IEC 61131-3.
type TIME_OF_DAY time.Time

// DT (DATE_AND_TIME) represents a specific date and time as defined by IEC 61131-3.
type DT time.Time

// TOD is an alias for TIME_OF_DAY
type TOD = TIME_OF_DAY

// TIMESPEC is a generic time type, useful for internal representations.
type TIMESPEC time.Time

// STEP X current value, prevState previous Value, T time elapsed
type STEP struct {
	X         BOOL
	prevState BOOL
	T         TIMESPEC
}

// TM is a helper struct for timer constants, not a standard IEC type.
type TM struct {
	D  int
	H  int
	M  int
	S  int
	Ms int
}

// Generic ANY types for function overloading (conceptual)
type ANY interface{}

// ANY_BOOL
type ANY_BOOL interface {
	~bool
}

// ANY BYTES
type ANY_UINTS interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint // Corresponds to BYTE, WORD, DWORD, LWORD
}

type ANY_SINTS interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int // Corresponds to SINT, INT, DINT, LINT, USINT, UINT, UDINT, ULINT
}

// ANY_BIT_STRING represents the group of bit-string types as per IEC 61131-3.
type ANY_BIT interface {
	~bool | ANY_UINTS
}

// ANY_INT
type ANY_INT interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~int8 | ~int16 | ~int32 | ~int64 | ~int
}

// ANY_DATE
type ANY_DATE interface {
	DATE | TOD | DT
}

// ANY_STRING
type ANY_STRING interface {
	~string | ~rune
}

// ANY_REAL
type ANY_REAL interface {
	~float32 | ~float64
}

// ANY_NUM
type ANY_NUM interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~int8 | ~int16 | ~int32 | ~int64 | ~int | ~float32 | ~float64
}

// ANY_DURATION
type ANY_DURATION interface {
	TIME
}

// ANY_MAGNITUDE
type ANY_MAGNITUDE interface {
	ANY_NUM | ANY_DURATION | ANY_STRING
}

// ANY
type ANY_ELEMENTARY interface {
	ANY_BIT | ANY_MAGNITUDE
}

// ANY_COMPLEX
type ANY_COMPLEX interface {
	COMPLEX | LCOMPLEX
}

// Value of BOOL
func (in *BOOL) Value() bool {
	return *in == true
}

/* // Value of TIMESPEC
func (in *TIMESPEC) Value() (out STRING) {
} */
// String returns the IEC 61131-3 string representation of a TIME value (e.g., "T#5s").
func (in *TIME) String() string {
	return "T#" + time.Duration(*in).String()
}

// String returns the IEC 61131-3 string representation of a TIME_OF_DAY value (e.g., "TOD#15:04:05.000").
func (in *TOD) String() string {
	return "TOD#" + time.Time(*in).Format("15:04:05.000")
}

// String returns the IEC 61131-3 string representation of a DATE value (e.g., "D#2006-01-02").
func (in *DATE) String() string {
	return "D#" + time.Time(*in).Format("2006-01-02")
}

// String returns the IEC 61131-3 string representation of a DATE_AND_TIME value (e.g., "DT#2006-01-02-15:04:05").
func (in *DT) String() string {
	return "DT#" + time.Time(*in).Format("2006-01-02-15:04:05")
}

// Value of Byte
func (in *BYTE) Value() uint8 { return uint8(*in) }

// Value of WORD
func (in *WORD) Value() uint16 { return uint16(*in) }

// Value of DWORD
func (in *DWORD) Value() uint32 { return uint32(*in) }

// Value of LWORD
func (in *LWORD) Value() uint64 { return uint64(*in) }

// Value of REAL
func (in *REAL) Value() float32 {
	return float32(*in)
}

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

func (in *DATE) CONVERT() LINT {
	return LINT(time.Time(*in).UnixMilli())
}

func (in *DT) CONVERT() LINT {
	return LINT(time.Time(*in).UnixMilli())
}

func (in *TOD) CONVERT() LINT {
	t := time.Time(*in)
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return LINT(t.Sub(midnight).Milliseconds())
}
