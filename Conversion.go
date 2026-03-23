/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.
 */package royaljelly

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

func lpad(s string, pad string, plength int) string {
	for i := len(s); i < plength; i++ {
		s = pad + s
	}
	return s
}

// Value of BOOL
func (in *BOOL) Value() bool {
	return *in == true
}

/* // Value of TIMESPEC
func (in *TIMESPEC) Value() (out STRING) {
} */
// String returns the IEC 61131-3 string representation of a TIME value (e.g., "T#5s").
func (in TIME) String() string {
	return "T#" + time.Duration(in).String()
}

// String returns the IEC 61131-3 string representation of a TIME_OF_DAY value (e.g., "TOD#15:04:05.000").
func (in TOD) String() string {
	return "TOD#" + time.Time(in).Format("15:04:05.000")
}

// String returns the IEC 61131-3 string representation of a DATE value (e.g., "D#2006-01-02").
func (in DATE) String() string {
	return "D#" + time.Time(in).Format("2006-01-02")
}

// String returns the IEC 61131-3 string representation of a DATE_AND_TIME value (e.g., "DT#2006-01-02-15:04:05").
func (in DT) String() string {
	return "DT#" + time.Time(in).Format("2006-01-02-15:04:05")
}

// Value of Byte
func (in BYTE) Value() uint8 { return uint8(in) }

// Value of WORD
func (in WORD) Value() uint16 { return uint16(in) }

// Value of DWORD
func (in DWORD) Value() uint32 { return uint32(in) }

// Value of LWORD
func (in LWORD) Value() uint64 { return uint64(in) }

// Value of REAL
func (in *REAL) Value() float32 {
	return float32(*in)
}

/*
simplified converstions
*/

func SubByte(in interface{}) (BYTE, error) {
	val, err := anyToULINT(in)
	if err != nil {
		return 0, fmt.Errorf("SubByte: conversion error: %w", err)
	}
	return BYTE(val), nil
}

func SubWord(in interface{}) (WORD, error) {
	val, err := anyToULINT(in)
	if err != nil {
		return 0, fmt.Errorf("SubWord: conversion error: %w", err)
	}
	return WORD(val), nil
}

func SubDword(in interface{}) (DWORD, error) {
	val, err := anyToULINT(in)
	if err != nil {
		return 0, fmt.Errorf("SubDword: conversion error: %w", err)
	}
	return DWORD(val), nil
}

func SubLword(in interface{}) (LWORD, error) {
	val, err := anyToULINT(in)
	if err != nil {
		return 0, fmt.Errorf("SubLword: conversion error: %w", err)
	}
	return LWORD(val), nil
}

func SubDt(in interface{}) (DT, error) {
	val, err := anyToLINT(in)
	if err != nil {
		return DT{}, fmt.Errorf("SubDt: conversion error: %w", err)
	}
	// Assuming the integer value represents milliseconds since Unix epoch
	return DT(time.UnixMilli(int64(val))), nil
}

func SubDate(in interface{}) (DATE, error) {
	val, err := anyToLINT(in)
	if err != nil {
		return DATE{}, fmt.Errorf("SubDate: conversion error: %w", err)
	}
	// Assuming the integer value represents milliseconds since Unix epoch
	return DATE(time.UnixMilli(int64(val))), nil
}

func SubTod(in interface{}) (TOD, error) {
	// Assuming input is milliseconds since midnight
	val, err := anyToLINT(in)
	if err != nil {
		return INITTOD, fmt.Errorf("SubTod: conversion error: %v", err)
	}
	// A TOD is a duration since midnight on an arbitrary day.
	return TOD(time.Time{}.Add(time.Duration(val) * time.Millisecond)), nil
}

func SubTime(in interface{}) (TIME, error) {
	val, err := anyToLINT(in)
	if err != nil {
		return 0, fmt.Errorf("SubTime: conversion error: %w", err)
	}
	return TIME(time.Duration(val) * time.Millisecond), nil
}

func SubBool(in BOOL) (out INT) {
	if in {
		out = 1
	} else {
		out = 0
	}
	return out
}

func clampLINT(val LINT, min, max LINT) LINT {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func clampULINT(val ULINT, max ULINT) ULINT {
	if val > max {
		return max
	}
	return val
}

func clampLREAL(val LREAL, min, max LREAL) LREAL {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func roundAndClampLREAL(val LREAL, min, max LREAL) LREAL {
	rounded := math.RoundToEven(float64(val))
	return clampLREAL(LREAL(rounded), min, max)
}

func SubFromTo(in1 interface{}, in2 interface{}) interface{} {
	out := reflect.TypeOf(in2)
	return reflect.ValueOf(in1).Convert(out).Interface()
}

func (in *DATE) CONVERT() (out reflect.Value) {
	out = reflect.ValueOf(LINT(time.Time(*in).UnixMilli()))
	return out
}

func (in *DT) CONVERT() (out reflect.Value) {
	out = reflect.ValueOf(LINT(time.Time(*in).UnixMilli()))
	return out
}

func (in *TOD) CONVERT() (out reflect.Value) {
	t := time.Time(*in)
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	out = reflect.ValueOf(LINT(t.Sub(midnight).Milliseconds()))
	return out
}

/*
BOOL_TO conversion
*/

// BOOL_TO_BYTE conversion
func BOOL_TO_BYTE(in BOOL) BYTE { out, _ := SubByte(in); return out }

// BOOL_TO_WORD conversion
func BOOL_TO_WORD(in BOOL) WORD { out, _ := SubWord(in); return out }

// BOOL_TO_DWORD conversion
func BOOL_TO_DWORD(in BOOL) DWORD { out, _ := SubDword(in); return out }

// BOOL_TO_LWORD conversion
func BOOL_TO_LWORD(in BOOL) LWORD { out, _ := SubLword(in); return out }

// BOOL_TO_INT conversion
func BOOL_TO_INT(in BOOL) INT { return INT(SubBool(in)) }

// BOOL_TO_SINT conversion
func BOOL_TO_SINT(in BOOL) SINT { return SINT(clampLINT(LINT(SubBool(in)), MINSINT, MAXSINT)) }

// BOOL_TO_UINT conversion
func BOOL_TO_UINT(in BOOL) UINT { return UINT(clampULINT(ULINT(SubBool(in)), MAXUINT)) }

// BOOL_TO_UDINT conversion
func BOOL_TO_UDINT(in BOOL) UDINT { return UDINT(SubBool(in)) }

// BOOL_TO_DINT conversion
func BOOL_TO_DINT(in BOOL) DINT { return DINT(SubBool(in)) }

// BOOL_TO_USINT conversion
func BOOL_TO_USINT(in BOOL) USINT { return USINT(clampULINT(ULINT(SubBool(in)), MAXUSINT)) }

// BOOL_TO_ULINT conversion
func BOOL_TO_ULINT(in BOOL) ULINT { return ULINT(clampULINT(ULINT(SubBool(in)), MAXULINT)) }

// BOOL_TO_LINT conversion
func BOOL_TO_LINT(in BOOL) LINT { return LINT(SubBool(in)) }

// BOOL_TO_STRING conversion
func BOOL_TO_STRING(in BOOL) (out STRING) { return STRING(fmt.Sprint(in)) }

// BOOL_TO_REAL conversion
func BOOL_TO_REAL(in BOOL) REAL { return REAL(SubBool(in)) }

// BOOL_TO_LREAL conversion
func BOOL_TO_LREAL(in BOOL) LREAL { return LREAL(SubBool(in)) }

// BOOL_TO_TIME conversion
func BOOL_TO_TIME(in BOOL) TIME { out, _ := SubTime(in); return out }

// BOOL_TO_TOD conversion
func BOOL_TO_TOD(in BOOL) TOD { out, _ := SubTod(in); return out }

// BOOL_TO_DATE conversion
func BOOL_TO_DATE(in BOOL) DATE { out, _ := SubDate(in); return out }

// BOOL_TO_DT conversion
func BOOL_TO_DT(in BOOL) DT { out, _ := SubDt(in); return out }

/*
BYTE_TO * Conversion Section
*/

// BYTE_TO_BOOL conversion
func BYTE_TO_BOOL(in BYTE) BOOL { return (in.Value() > 0) }

// BYTE_TO_SINT conversion
func BYTE_TO_SINT(in BYTE) SINT { return SINT(clampULINT(ULINT(in), MAXSINT)) }

// BYTE_TO_INT conversion
func BYTE_TO_INT(in BYTE) INT { return INT(clampULINT(ULINT(in), MAXINT)) }

// BYTE_TO_DINT conversion
func BYTE_TO_DINT(in BYTE) DINT { return DINT(clampULINT(ULINT(in), MAXDINT)) }

// BYTE_TO_LINT conversion
func BYTE_TO_LINT(in BYTE) LINT { return LINT(in.Value()) }

// BYTE_TO_ANYINT conversion
func BYTE_TO_ANYINT(in BYTE) ANYINT { return ANYINT(clampULINT(ULINT(in), MAXANYINT)) }

// BYTE_TO_USINT conversion
func BYTE_TO_USINT(in BYTE) USINT { return USINT(in.Value()) }

// BYTE_TO_UINT conversion
func BYTE_TO_UINT(in BYTE) UINT { return UINT(clampULINT(ULINT(in), MAXUINT)) }

// BYTE_TO_Udint conversion
func BYTE_TO_UDINT(in BYTE) UDINT { return UDINT(clampULINT(ULINT(in), MAXUDINT)) }

// BYTE_TO_ULINT conversion
func BYTE_TO_ULINT(in BYTE) ULINT { return ULINT(in.Value()) }

// BYTE_TO_ANYUINT conversion
func BYTE_TO_ANYUINT(in BYTE) ANYUINT { return ANYUINT(clampULINT(ULINT(in), MAXANYUINT)) }

// BYTE_TO_REAL conversion
func BYTE_TO_REAL(in BYTE) REAL { return REAL(in.Value()) }

// BYTE_TO_LREAL conversion
func BYTE_TO_LREAL(in BYTE) LREAL { return LREAL(in.Value()) }

// BYTE_TO_WORD conversion
func BYTE_TO_WORD(in BYTE) WORD { out, _ := SubWord(in.Value()); return out }

// BYTE_TO_DWORD conversion
func BYTE_TO_DWORD(in BYTE) DWORD { out, _ := SubDword(in.Value()); return out }

// BYTE_TO_LWORD conversion
func BYTE_TO_LWORD(in BYTE) LWORD { out, _ := SubLword(in.Value()); return out }

// BYTE_TO_STRING conversion
func BYTE_TO_STRING(in BYTE) STRING { return STRING(strconv.FormatUint(uint64(in), 10)) }

// BYTES_TO_STRING conversion
func BYTES_TO_STRING(in []BYTE) STRING {
	byteSlice := make([]byte, len(in))
	for i, b := range in {
		byteSlice[i] = b.Value()
	}
	return STRING(byteSlice)
}

// BYTE_TO_DATE conversion
func BYTE_TO_DATE(in BYTE) DATE { out, _ := SubDate(in); return out }

// BYTE_TO_DT conversion
func BYTE_TO_DT(in BYTE) DT { out, _ := SubDt(in); return out }

// BYTE_TO_TOD conversion
func BYTE_TO_TOD(in BYTE) TOD { out, _ := SubTod(in); return out }

// BYTE_TO_TIME conversion
func BYTE_TO_TIME(in BYTE) TIME { out, _ := SubTime(in); return out }

/*
WORD_TO * Conversion Section
*/

// WORD_TO_BOOL covnersion
func WORD_TO_BOOL(in WORD) BOOL { return (in.Value() > 0) }

// WORD_TO_SINT conversion
func WORD_TO_SINT(in WORD) SINT { return SINT(clampULINT(ULINT(in), MAXSINT)) }

// WORD_TO_INT conversion
func WORD_TO_INT(in WORD) INT { return INT(clampULINT(ULINT(in), MAXINT)) }

// WORD_TO_DINT conversion
func WORD_TO_DINT(in WORD) DINT { return DINT(clampULINT(ULINT(in), MAXDINT)) }

// WORD_TO_DATE conversion
func WORD_TO_DATE(in WORD) DATE { out, _ := SubDate(in); return out }

// WORD_TO_Lint conversion
func WORD_TO_LINT(in WORD) LINT { return LINT(clampULINT(ULINT(in), MAXLINT)) }

// WORD_TO_AnyInt conversion
func WORD_TO_AnyInt(in WORD) ANYINT { return ANYINT(clampULINT(ULINT(in), MAXANYINT)) }

// WORD_TO_USINT conversion
func WORD_TO_USINT(in WORD) USINT { return USINT(clampULINT(ULINT(in), MAXUSINT)) }

// WORD_TO_UINT conversion
func WORD_TO_UINT(in WORD) UINT { return UINT(in.Value()) }

// WORD_TO_UDINT conversion
func WORD_TO_UDINT(in WORD) UDINT { return UDINT(clampULINT(ULINT(in), MAXUDINT)) }

// WORD_TO_ULINT conversion
func WORD_TO_ULINT(in WORD) ULINT { return ULINT(in.Value()) }

// WORD_TO_ANYUINT conversion
func WORD_TO_ANYUINT(in WORD) ANYUINT { return ANYUINT(clampULINT(ULINT(in), MAXANYUINT)) }

// WORD_TO_REAL conversion
func WORD_TO_REAL(in WORD) REAL { return REAL(clampLREAL(LREAL(in), -MAXREAL, MAXREAL)) }

// WORD_TO_LREAL conversion
func WORD_TO_LREAL(in WORD) LREAL { return LREAL(in.Value()) }

// WORD_TO_BYTE conversion
func WORD_TO_BYTE(in WORD) BYTE { out, _ := SubByte(in); return out }

// WORD_TO_DWORD conversion
func WORD_TO_DWORD(in WORD) DWORD { out, _ := SubDword(in); return out }

// WORD_TO_DT conversion
func WORD_TO_DT(in WORD) DT { out, _ := SubDt(in); return out }

// WORD_TO_TOD conversion
func WORD_TO_TOD(in WORD) TOD { out, _ := SubTod(in); return out }

// WORD_TO_LWORD conversion
func WORD_TO_LWORD(in WORD) LWORD { out, _ := SubLword(in); return out }

// WORD_TO_STRING conversion
func WORD_TO_STRING(in WORD) STRING { return STRING(strconv.FormatUint(uint64(in.Value()), 10)) }

// WORD_TO_TIM conversion
func WORD_TO_TIME(in WORD) TIME { out, _ := SubTime(in); return out }

/*
DWORD_TO * Conversion Section
*/

// DWORD_TO_BOOL covnersion
func DWORD_TO_BOOL(in DWORD) BOOL { return in.Value() > 0 }

// DWORD_TO_SINT conversion
func DWORD_TO_SINT(in DWORD) SINT { return SINT(clampULINT(ULINT(in), MAXSINT)) }

// DWORD_TO_INT conversion
func DWORD_TO_INT(in DWORD) INT { return INT(clampULINT(ULINT(in), MAXINT)) }

// DWORD_TO_DINT conversion
func DWORD_TO_DINT(in DWORD) DINT { return DINT(clampULINT(ULINT(in), MAXDINT)) }

// DWORD_TO_DATE conversion
func DWORD_TO_DATE(in DWORD) DATE { out, _ := SubDate(in); return out }

// DWORD_TO_DT conversion
func DWORD_TO_DT(in DWORD) DT { out, _ := SubDt(in); return out }

// DWORD_TO_TOD conversion
func DWORD_TO_TOD(in DWORD) TOD { out, _ := SubTod(in); return out }

// DWORD_TO_LINT conversion
func DWORD_TO_LINT(in DWORD) LINT { return LINT(clampULINT(ULINT(in), MAXLINT)) }

// DWORD_TO_ANYINT conversion
func DWORD_TO_ANYINT(in DWORD) ANYINT { return ANYINT(clampULINT(ULINT(in), MAXANYINT)) }

// DWORD_TO_USINT conversion
func DWORD_TO_USINT(in DWORD) USINT { return USINT(clampULINT(ULINT(in), MAXUSINT)) }

// DWORD_TO_UINT conversion
func DWORD_TO_UINT(in DWORD) UINT { return UINT(clampULINT(ULINT(in), MAXUINT)) }

// DWORD_TO_UDINT conversion
func DWORD_TO_UDINT(in DWORD) UDINT { return UDINT(in.Value()) }

// DWORD_TO_ULINT conversion
func DWORD_TO_ULINT(in DWORD) ULINT { return ULINT(clampULINT(ULINT(in), MAXULINT)) }

// DWORD_TO_ANYUINT conversion
func DWORD_TO_ANYUINT(in DWORD) ANYUINT { return ANYUINT(clampULINT(ULINT(in), MAXANYUINT)) }

// DWORD_TO_REAL conversion
func DWORD_TO_REAL(in DWORD) REAL { return REAL(clampLREAL(LREAL(in), -MAXREAL, MAXREAL)) }

// DWORD_TO_LREAL conversion
func DWORD_TO_LREAL(in DWORD) LREAL { return LREAL(in.Value()) }

// DWORD_TO_BYTE conversion
func DWORD_TO_BYTE(in DWORD) BYTE { out, _ := SubByte(in); return out }

// DWORD_TO_WORD conversion
func DWORD_TO_WORD(in DWORD) WORD { out, _ := SubWord(in); return out }

// DWORD_TO_LWORD conversion
func DWORD_TO_LWORD(in DWORD) LWORD { out, _ := SubLword(in); return out }

// DWORD_TO_STRING conversion
func DWORD_TO_STRING(in DWORD) STRING { return STRING(strconv.FormatUint(uint64(in.Value()), 10)) }

// DWORD_TO_TIME conversion
func DWORD_TO_TIME(in DWORD) TIME { out, _ := SubTime(in); return out }

/*
LWORD_TO * Conversion Section
*/

// LWORD_TO_BOOL covnersion
func LWORD_TO_BOOL(in LWORD) BOOL { return in.Value() > 0 }

// LWORD_TO_SINT conversion
func LWORD_TO_SINT(in LWORD) SINT { return SINT(clampULINT(ULINT(in), MAXSINT)) }

// LWORD_TO_INT conversion
func LWORD_TO_INT(in LWORD) INT { return INT(clampULINT(ULINT(in), MAXINT)) }

// LWORD_TO_DINT conversion
func LWORD_TO_DINT(in LWORD) DINT { return DINT(clampULINT(ULINT(in), MAXDINT)) }

// LWORD_TO_LINT conversion
func LWORD_TO_LINT(in LWORD) LINT { return LINT(clampULINT(ULINT(in), MAXLINT)) }

// LWORD_TO_ANYINT conversion
func LWORD_TO_ANYINT(in LWORD) ANYINT { return ANYINT(clampULINT(ULINT(in), MAXANYINT)) }

// LWORD_TO_USINT conversion
func LWORD_TO_USINT(in LWORD) USINT { return USINT(clampULINT(ULINT(in), MAXUSINT)) }

// LWORD_TO_UINT conversion
func LWORD_TO_UINT(in LWORD) UINT { return UINT(clampULINT(ULINT(in), MAXUINT)) }

// LWORD_TO_UDINT conversion
func LWORD_TO_UDINT(in LWORD) UDINT { return UDINT(clampULINT(ULINT(in), MAXUDINT)) }

// LWORD_TO_ULINT conversion
func LWORD_TO_ULINT(in LWORD) ULINT { return ULINT(in.Value()) }

// LWORD_TO_ANYUINT conversion
func LWORD_TO_ANYUINT(in LWORD) ANYUINT { return ANYUINT(clampULINT(ULINT(in), MAXANYUINT)) }

// LWORD_TO_REAL conversion
func LWORD_TO_REAL(in LWORD) REAL { return REAL(clampLREAL(LREAL(in), -MAXREAL, MAXREAL)) }

// LWORD_TO_LREAL conversion
func LWORD_TO_LREAL(in LWORD) LREAL { return LREAL(in.Value()) }

// LWORD_TO_BYTE conversion
func LWORD_TO_BYTE(in LWORD) BYTE { out, _ := SubByte(in); return out }

// LWORD_TO_WORD conversion
func LWORD_TO_WORD(in LWORD) WORD { out, _ := SubWord(in); return out }

// LWORD_TO_DWORD conversion
func LWORD_TO_DWORD(in LWORD) DWORD { out, _ := SubDword(in); return out }

// LWORD_TO_STRING conversion
func LWORD_TO_STRING(in LWORD) STRING { return STRING(strconv.FormatUint(in.Value(), 10)) }

// LWORD_TO_DATE conversion
func LWORD_TO_DATE(in LWORD) DATE { out, _ := SubDate(in); return out }

// LWORD_TO_DT conversion
func LWORD_TO_DT(in LWORD) DT { out, _ := SubDt(in); return out }

// LWORD_TO_TOD conversion
func LWORD_TO_TOD(in LWORD) TOD { out, _ := SubTod(in); return out }

// LWORD_TO_TIME converion
func LWORD_TO_TIME(in LWORD) TIME { out, _ := SubTime(in); return out }

/*
REAL_TO * Conversion Section
*/

// REAL_TO_SINT conversion
func REAL_TO_SINT(in REAL) SINT { return SINT(roundAndClampLREAL(LREAL(in), MINSINT, MAXSINT)) }

// REAL_TO_LINT conversion
func REAL_TO_LINT(in REAL) LINT { return LINT(roundAndClampLREAL(LREAL(in), MINLINT, MAXLINT)) }

// REAL_TO_DINT conversion
func REAL_TO_DINT(in REAL) DINT { return DINT(roundAndClampLREAL(LREAL(in), MINDINT, MAXDINT)) }

// REAL_TO_DATE conversion
func REAL_TO_DATE(in REAL) DATE { out, _ := SubDate(in); return out }

// REAL_TO_DWORD conversion
func REAL_TO_DWORD(in REAL) DWORD { out, _ := SubDword(in); return out }

// REAL_TO_DT conversion
func REAL_TO_DT(in REAL) DT { out, _ := SubDt(in); return out }

// REAL_TO_TOD conversion
func REAL_TO_TOD(in REAL) TOD { out, _ := SubTod(in); return out }

// REAL_TO_UDINT conversion
func REAL_TO_UDINT(in REAL) UDINT { return UDINT(roundAndClampLREAL(LREAL(in), 0, MAXUDINT)) }

// REAL_TO_WORD conversion
func REAL_TO_WORD(in REAL) WORD { out, _ := SubWord(in); return out }

// REAL_TO_STRING conversion
func REAL_TO_STRING(in REAL) STRING { return STRING(strconv.FormatFloat(float64(in), 'g', -1, 32)) }

// REAL_TO_LWORD conversion
func REAL_TO_LWORD(in REAL) LWORD { out, _ := SubLword(in); return out }

// REAL_TO_UINT conversion
func REAL_TO_UINT(in REAL) UINT { return UINT(roundAndClampLREAL(LREAL(in), 0, MAXUINT)) }

// REAL_TO_LREAL conversion
func REAL_TO_LREAL(in REAL) LREAL { return LREAL(in) }

// REAL_TO_BYTE conversion
func REAL_TO_BYTE(in REAL) BYTE { out, _ := SubByte(in); return out }

// REAL_TO_USINT conversion
func REAL_TO_USINT(in REAL) USINT { return USINT(roundAndClampLREAL(LREAL(in), 0, MAXUSINT)) }

// REAL_TO_ULINT conversion
func REAL_TO_ULINT(in REAL) ULINT { return ULINT(roundAndClampLREAL(LREAL(in), 0, MAXULINT)) }

// REAL_TO_BOOL conversion
func REAL_TO_BOOL(in REAL) BOOL { return in > 0 }

// REAL_TO_TIME conversion
func REAL_TO_TIME(in REAL) TIME { out, _ := SubTime(in); return out }

// REAL_TO_INT conversion
func REAL_TO_INT(in REAL) INT { return INT(roundAndClampLREAL(LREAL(in), MININT, MAXINT)) }

/*
LREAL_TO * Conversion Section
*/
// LREAL_TO_REAL conversion
func LREAL_TO_REAL(in LREAL) REAL { return REAL(clampLREAL(in, -MAXREAL, MAXREAL)) }

// LREAL_TO_SINT conversion
func LREAL_TO_SINT(in LREAL) SINT { return SINT(roundAndClampLREAL(in, MINSINT, MAXSINT)) }

// LREAL_TO_LINT conversion
func LREAL_TO_LINT(in LREAL) LINT { return LINT(roundAndClampLREAL(in, MINLINT, MAXLINT)) }

// LREAL_TO_DINT conversion
func LREAL_TO_DINT(in LREAL) DINT { return DINT(roundAndClampLREAL(in, MINDINT, MAXDINT)) }

// LREAL_TO_DATE conversion
func LREAL_TO_DATE(in LREAL) DATE { out, _ := SubDate(in); return out }

// LREAL_TO_DWORD conversion
func LREAL_TO_DWORD(in LREAL) DWORD { out, _ := SubDword(in); return out }

// LREAL_TO_DT conversion
func LREAL_TO_DT(in LREAL) DT { out, _ := SubDt(in); return out }

// LREAL_TO_TOD conversion
func LREAL_TO_TOD(in LREAL) TOD { out, _ := SubTod(in); return out }

// LREAL_TO_UDINT conversion
func LREAL_TO_UDINT(in LREAL) UDINT { return UDINT(roundAndClampLREAL(in, 0, MAXUDINT)) }

// LREAL_TO_WORD conversion
func LREAL_TO_WORD(in LREAL) WORD { out, _ := SubWord(in); return out }

// LREAL_TO_STRING conversion
func LREAL_TO_STRING(in LREAL) STRING {
	return STRING(strconv.FormatFloat(float64(in), 'g', -1, 64))
}

// LREAL_TO_LWORD conversion
func LREAL_TO_LWORD(in LREAL) LWORD { out, _ := SubLword(in); return out }

// LREAL_TO_UINT conversion
func LREAL_TO_UINT(in LREAL) UINT { return UINT(roundAndClampLREAL(in, 0, MAXUINT)) }

// LREAL_TO_BYTE conversion
func LREAL_TO_BYTE(in LREAL) BYTE { out, _ := SubByte(in); return out }

// LREAL_TO_USINT conversion
func LREAL_TO_USINT(in LREAL) USINT { return USINT(roundAndClampLREAL(in, 0, MAXUSINT)) }

// LREAL_TO_ULINT conversion
func LREAL_TO_ULINT(in LREAL) ULINT { return ULINT(roundAndClampLREAL(in, 0, MAXULINT)) }

// LREAL_TO_BOOL conversion
func LREAL_TO_BOOL(in LREAL) BOOL { return in > 0 }

// LREAL_TO_TIME conversion
func LREAL_TO_TIME(in LREAL) TIME { out, _ := SubTime(in); return out }

// LREAL_TO_INT conversion
func LREAL_TO_INT(in LREAL) INT { return INT(roundAndClampLREAL(in, MININT, MAXINT)) }

/*
SINT_TO * Conversion section
*/

// SINT_TO_REAL conversion
func SINT_TO_REAL(in SINT) REAL { return REAL(in) }

// SINT_TO_LINT conversion
func SINT_TO_LINT(in SINT) LINT { return LINT(in) }

// SINT_TO_DINT conversion
func SINT_TO_DINT(in SINT) DINT { return DINT(clampLINT(LINT(in), MINDINT, MAXDINT)) }

// SINT_TO_DATE conversion
func SINT_TO_DATE(in SINT) DATE { out, _ := SubDate(in); return out }

// SINT_TO_DWORD conversion
func SINT_TO_DWORD(in SINT) DWORD { out, _ := SubDword(in); return out }

// SINT_TO_DT conversion
func SINT_TO_DT(in SINT) DT { out, _ := SubDt(in); return out }

// SINT_TO_TOD conversion
func SINT_TO_TOD(in SINT) TOD { out, _ := SubTod(in); return out }

// SINT_TO_UDINT conversion
func SINT_TO_UDINT(in SINT) UDINT { return UDINT(clampLINT(LINT(in), 0, MAXUDINT)) }

// SINT_TO_WORD conversion
func SINT_TO_WORD(in SINT) WORD { out, _ := SubWord(in); return out }

// SINT_TO_STRING conversion
func SINT_TO_STRING(in SINT) STRING { return STRING(strconv.FormatInt(int64(in), 10)) }

// SINT_TO_LWORD conversion
func SINT_TO_LWORD(in SINT) LWORD { out, _ := SubLword(in); return out }

// SINT_TO_UINT conversion
func SINT_TO_UINT(in SINT) UINT { return UINT(clampLINT(LINT(in), 0, MAXUINT)) }

// SINT_TO_LREAL conversion
func SINT_TO_LREAL(in SINT) LREAL { return LREAL(in) }

// SINT_TO_BYTE conversion
func SINT_TO_BYTE(in SINT) BYTE { out, _ := SubByte(in); return out }

// SINT_TO_USINT conversion
func SINT_TO_USINT(in SINT) USINT { return USINT(clampLINT(LINT(in), 0, MAXUSINT)) }

// SINT_TO_ULINT conversion
func SINT_TO_ULINT(in SINT) ULINT { return ULINT(clampLINT(LINT(in), 0, ULINT_TO_LINT(MAXULINT))) }

// SINT_TO_BOOL conversion
func SINT_TO_BOOL(in SINT) BOOL { return BOOL(in != 0) }

// SINT_TO_TIME conversion
func SINT_TO_TIME(in SINT) TIME { out, _ := SubTime(in); return out }

// SINT_TO_INT conversion
func SINT_TO_INT(in SINT) INT { return INT(in) }

/*
INT_TO * Conversion section
*/

// INT_TO_REAL conversion
func INT_TO_REAL(in INT) REAL { return REAL(in) }

// INT_TO_SINT conversion
func INT_TO_SINT(in INT) SINT { return SINT(clampLINT(LINT(in), MINSINT, MAXSINT)) }

// INT_TO_LINT conversion
func INT_TO_LINT(in INT) LINT { return LINT(in) }

// INT_TO_DINT conversion
func INT_TO_DINT(in INT) DINT { return DINT(clampLINT(LINT(in), MINDINT, MAXDINT)) }

// INT_TO_DATE conversion seconds to DATE
func INT_TO_DATE(in INT) DATE { out, _ := SubDate(in); return out }

// INT_TO_DWORD conversion
func INT_TO_DWORD(in INT) DWORD { out, _ := SubDword(in); return out }

// INT_TO_DT conversion
func INT_TO_DT(in INT) DT { out, _ := SubDt(in); return out }

// INT_TO_TOD conversion
func INT_TO_TOD(in INT) TOD { out, _ := SubTod(in); return out }

// INT_TO_UDINT conversion
func INT_TO_UDINT(in INT) UDINT { return UDINT(clampLINT(LINT(in), 0, MAXUDINT)) }

// INT_TO_WORD conversion
func INT_TO_WORD(in INT) WORD { out, _ := SubWord(in); return out }

// INT_TO_STRING conversion
func INT_TO_STRING(in INT) STRING { return STRING(strconv.FormatInt(int64(in), 10)) }

// INT_TO_LWORD conversion
func INT_TO_LWORD(in INT) LWORD { out, _ := SubLword(in); return out }

// INT_TO_UINT conversion
func INT_TO_UINT(in INT) UINT { return UINT(clampLINT(LINT(in), 0, MAXUINT)) }

// INT_TO_LREAL conversion
func INT_TO_LREAL(in INT) LREAL { return LREAL(in) }

// INT_TO_BYTE conversion
func INT_TO_BYTE(in INT) BYTE { out, _ := SubByte(in); return out }

// INT_TO_USINT conversion
func INT_TO_USINT(in INT) USINT { return USINT(clampLINT(LINT(in), 0, MAXUSINT)) }

// INT_TO_ULINT conversion
func INT_TO_ULINT(in INT) ULINT { return ULINT(clampLINT(LINT(in), 0, ULINT_TO_LINT(MAXULINT))) }

// INT_TO_BOOL conversion
func INT_TO_BOOL(in INT) BOOL { return in > 0 }

// INT_TO_TIME conversion
func INT_TO_TIME(in INT) TIME { out, _ := SubTime(in); return out }

/*
LINT_TO * Conversion section
*/

// LINT_TO_REAL conversion
func LINT_TO_REAL(in LINT) REAL { return REAL(in) }

// LINT_TO_SINT conversion
func LINT_TO_SINT(in LINT) SINT { return SINT(clampLINT(in, MINSINT, MAXSINT)) }

// LINT_TO_DINT conversion
func LINT_TO_DINT(in LINT) DINT { return DINT(clampLINT(in, MINDINT, MAXDINT)) }

// LINT_TO_DWORD conversion
func LINT_TO_DWORD(in LINT) DWORD { out, _ := SubDword(in); return out }

// LINT_TO_DATE conversion
func LINT_TO_DATE(in LINT) DATE { out, _ := SubDate(in); return out }

// LINT_TO_TOD conversion
func LINT_TO_TOD(in LINT) TOD { out, _ := SubTod(in); return out }

// LINT_TO_UDINT conversion
func LINT_TO_UDINT(in LINT) UDINT { return UDINT(clampLINT(in, 0, MAXUDINT)) }

// LINT_TO_WORD conversion
func LINT_TO_WORD(in LINT) WORD { out, _ := SubWord(in); return out }

// LINT_TO_STRING conversion
func LINT_TO_STRING(in LINT) STRING { return STRING(strconv.FormatInt(int64(in), 10)) }

// LINT_TO_LWORD conversion
func LINT_TO_LWORD(in LINT) LWORD { out, _ := SubLword(in); return out }

// LINT_TO_UINT conversion
func LINT_TO_UINT(in LINT) UINT { return UINT(clampLINT(in, 0, MAXUINT)) }

// LINT_TO_LREAL conversion
func LINT_TO_LREAL(in LINT) LREAL { return LREAL(in) }

// LINT_TO_BYTE conversion
func LINT_TO_BYTE(in LINT) BYTE { out, _ := SubByte(in); return out }

// LINT_TO_USINT conversion
func LINT_TO_USINT(in LINT) USINT { return USINT(clampLINT(in, 0, MAXUSINT)) }

// LINT_TO_ULINT conversion
func LINT_TO_ULINT(in LINT) ULINT { return ULINT(clampLINT(in, 0, -1)) }

// LINT_TO_BOOL conversion
func LINT_TO_BOOL(in LINT) BOOL { return in > 0 }

// LINT_TO_TIME conversion
func LINT_TO_TIME(in LINT) TIME { out, _ := SubTime(in); return out }

// LINT_TO_INT conversion
func LINT_TO_INT(in LINT) INT { return INT(clampLINT(in, MININT, MAXINT)) }

/*
// DINT_TO Conversions
*/

// DINT_TO_REAL conversion
func DINT_TO_REAL(in DINT) REAL { return REAL(in) }

// DINT_TO_SINT conversion
func DINT_TO_SINT(in DINT) SINT { return SINT(clampLINT(LINT(in), MINSINT, MAXSINT)) }

// DINT_TO_LINT conversion
func DINT_TO_LINT(in DINT) LINT { return LINT(in) }

// DINT_TO_DATE conversion
func DINT_TO_DATE(in DINT) DATE { out, _ := SubDate(in); return out }

// DINT_TO_DWORD conversion
func DINT_TO_DWORD(in DINT) DWORD { out, _ := SubDword(in); return out }

// DINT_TO_DT conversion
func DINT_TO_DT(in DINT) DT { out, _ := SubDt(in); return out }

// DINT_TO_TOD conversion
func DINT_TO_TOD(in DINT) TOD { out, _ := SubTod(in); return out }

// DINT_TO_UDINT conversion
func DINT_TO_UDINT(in DINT) UDINT { return UDINT(clampLINT(LINT(in), 0, MAXUDINT)) }

// DINT_TO_WORD conversion
func DINT_TO_WORD(in DINT) WORD { out, _ := SubWord(in); return out }

// DINT_TO_STRING conversion
func DINT_TO_STRING(in DINT) STRING { return STRING(strconv.FormatInt(int64(in), 10)) }

// DINT_TO_LWORD conversion
func DINT_TO_LWORD(in DINT) LWORD { out, _ := SubLword(in); return out }

// DINT_TO_UINT conversion
func DINT_TO_UINT(in DINT) UINT { return UINT(clampLINT(LINT(in), 0, MAXUINT)) }

// DINT_TO_LREAL conversion
func DINT_TO_LREAL(in DINT) LREAL { return LREAL(in) }

// DINT_TO_BYTE conversion
func DINT_TO_BYTE(in DINT) BYTE { out, _ := SubByte(in); return out }

// DINT_TO_USINT conversion
func DINT_TO_USINT(in DINT) USINT { return USINT(clampLINT(LINT(in), 0, MAXUSINT)) }

// DINT_TO_ULINT conversion
func DINT_TO_ULINT(in DINT) ULINT { return ULINT(clampLINT(LINT(in), 0, ULINT_TO_LINT(MAXULINT))) }

// DINT_TO_BOOL conversion
func DINT_TO_BOOL(in DINT) BOOL { return in > 0 }

// DINT_TO_TIME conversion
func DINT_TO_TIME(in DINT) TIME { out, _ := SubTime(in); return out }

// DINT_TO_INT conversion
func DINT_TO_INT(in DINT) INT { return INT(clampLINT(LINT(in), MININT, MAXINT)) }

/*
USINT_TO * Conversion section
*/
// USINT_TO_REAL conversion
func USINT_TO_REAL(in USINT) REAL { return REAL(in) }

// USINT_TO_SINT conversion
func USINT_TO_SINT(in USINT) SINT { return SINT(clampULINT(ULINT(in), MAXSINT)) }

// USINT_TO_LINT conversion
func USINT_TO_LINT(in USINT) LINT { return LINT(clampULINT(ULINT(in), MAXLINT)) }

// USINT_TO_DINT conversion
func USINT_TO_DINT(in USINT) DINT { return DINT(clampULINT(ULINT(in), MAXDINT)) }

// USINT_TO_DATE conversion
func USINT_TO_DATE(in USINT) DATE { out, _ := SubDate(in); return out }

// USINT_TO_DWORD conversion
func USINT_TO_DWORD(in USINT) DWORD { out, _ := SubDword(in); return out }

// USINT_TO_DT conversion
func USINT_TO_DT(in USINT) DT { out, _ := SubDt(in); return out }

// USINT_TO_TOD conversion
func USINT_TO_TOD(in USINT) TOD { out, _ := SubTod(in); return out }

// USINT_TO_UDINT conversion
func USINT_TO_UDINT(in USINT) UDINT { return UDINT(clampULINT(ULINT(in), MAXUDINT)) }

// USINT_TO_WORD conversion
func USINT_TO_WORD(in USINT) WORD { out, _ := SubWord(in); return out }

// USINT_TO_STRING conversion
func USINT_TO_STRING(in USINT) STRING { return STRING(strconv.FormatUint(uint64(in), 10)) }

// USINT_TO_LWORD conversion
func USINT_TO_LWORD(in USINT) LWORD { out, _ := SubLword(in); return out }

// USINT_TO_UINT conversion
func USINT_TO_UINT(in USINT) UINT { return UINT(clampULINT(ULINT(in), MAXUINT)) }

// USINT_TO_LREAL conversion
func USINT_TO_LREAL(in USINT) LREAL { return LREAL(in) }

// USINT_TO_BYTE conversion
func USINT_TO_BYTE(in USINT) BYTE { out, _ := SubByte(in); return out }

// USINT_TO_ULINT conversion
func USINT_TO_ULINT(in USINT) ULINT { return ULINT(clampULINT(ULINT(in), MAXULINT)) }

// USINT_TO_BOOL conversion
func USINT_TO_BOOL(in USINT) BOOL { return in > 0 }

// USINT_TO_TIME conversion
func USINT_TO_TIME(in USINT) TIME { out, _ := SubTime(in); return out }

// USINT_TO_INT conversion
func USINT_TO_INT(in USINT) INT { return INT(clampULINT(ULINT(in), MAXINT)) }

/*
UINT_TO conversion
*/
// UINT_TO_REAL conversion
func UINT_TO_REAL(in UINT) REAL { return REAL(in) }

// UINT_TO_SINT conversion
func UINT_TO_SINT(in UINT) SINT { return SINT(clampULINT(ULINT(in), MAXSINT)) }

// UINT_TO_LINT conversion
func UINT_TO_LINT(in UINT) LINT { return LINT(clampULINT(ULINT(in), MAXLINT)) }

// UINT_TO_DINT conversion
func UINT_TO_DINT(in UINT) DINT { return DINT(clampULINT(ULINT(in), MAXDINT)) }

// UINT_TO_DATE conversion
func UINT_TO_DATE(in UINT) DATE { out, _ := SubDate(in); return out }

// UINT_TO_DWORD conversion
func UINT_TO_DWORD(in UINT) DWORD { out, _ := SubDword(in); return out }

// UINT_TO_DT conversion
func UINT_TO_DT(in UINT) DT { out, _ := SubDt(in); return out }

// UINT_TO_TOD conversion
func UINT_TO_TOD(in UINT) TOD { out, _ := SubTod(in); return out }

// UINT_TO_UDINT conversion
func UINT_TO_UDINT(in UINT) UDINT { return UDINT(clampULINT(ULINT(in), MAXUDINT)) }

// UINT_TO_WORD conversion
func UINT_TO_WORD(in UINT) WORD { out, _ := SubWord(in); return out }

// UINT_TO_STRING conversion
func UINT_TO_STRING(in UINT) STRING { return STRING(strconv.FormatUint(uint64(in), 10)) }

// UINT_TO_LWORD conversion
func UINT_TO_LWORD(in UINT) LWORD { out, _ := SubLword(in); return out }

// UINT_TO_LREAL conversion
func UINT_TO_LREAL(in UINT) LREAL { return LREAL(in) }

// UINT_TO_BYTE conversion
func UINT_TO_BYTE(in UINT) BYTE { out, _ := SubByte(in); return out }

// UINT_TO_USINT conversion
func UINT_TO_USINT(in UINT) USINT { return USINT(clampULINT(ULINT(in), MAXUSINT)) }

// UINT_TO_ULINT conversion
func UINT_TO_ULINT(in UINT) ULINT { return ULINT(clampULINT(ULINT(in), MAXULINT)) }

// UINT_TO_BOOL conversion
func UINT_TO_BOOL(in UINT) BOOL { return in > 0 }

// UINT_TO_TIME conversion
func UINT_TO_TIME(in UINT) TIME { out, _ := SubTime(in); return out }

// UINT_TO_INT conversion
func UINT_TO_INT(in UINT) INT { return INT(clampULINT(ULINT(in), MAXINT)) }

/*
UDINT_TO * Conversion section
*/
// UDINT_TO_REAL conversion
func UDINT_TO_REAL(in UDINT) REAL { return REAL(in) }

// UDINT_TO_SINT conversion
func UDINT_TO_SINT(in UDINT) SINT { return SINT(clampULINT(ULINT(in), MAXSINT)) }

// UDINT_TO_LINT conversion
func UDINT_TO_LINT(in UDINT) LINT { return LINT(clampULINT(ULINT(in), MAXLINT)) }

// UDINT_TO_DINT conversion
func UDINT_TO_DINT(in UDINT) DINT { return DINT(clampULINT(ULINT(in), MAXDINT)) }

// UDINT_TO_DATE conversion
func UDINT_TO_DATE(in UDINT) DATE { out, _ := SubDate(in); return out }

// UDINT_TO_DWORD conversion
func UDINT_TO_DWORD(in UDINT) DWORD { out, _ := SubDword(in); return out }

// UDINT_TO_DT conversion
func UDINT_TO_DT(in UDINT) DT { out, _ := SubDt(in); return out }

// UDINT_TO_TOD conversion
func UDINT_TO_TOD(in UDINT) TOD { out, _ := SubTod(in); return out }

// UDINT_TO_WORD conversion
func UDINT_TO_WORD(in UDINT) WORD { out, _ := SubWord(in); return out }

// UDINT_TO_STRING conversion
func UDINT_TO_STRING(in UDINT) STRING { return STRING(strconv.FormatUint(uint64(in), 10)) }

// UDINT_TO_LWORD conversion
func UDINT_TO_LWORD(in UDINT) LWORD { out, _ := SubLword(in); return out }

// UDINT_TO_UINT conversion
func UDINT_TO_UINT(in UDINT) UINT { return UINT(clampULINT(ULINT(in), MAXUINT)) }

// UDINT_TO_LREAL conversion
func UDINT_TO_LREAL(in UDINT) LREAL { return LREAL(in) }

// UDINT_TO_BYTE conversion
func UDINT_TO_BYTE(in UDINT) BYTE { out, _ := SubByte(in); return out }

// UDINT_TO_USINT conversion
func UDINT_TO_USINT(in UDINT) USINT { return USINT(clampULINT(ULINT(in), MAXUSINT)) }

// UDINT_TO_ULINT conversion
func UDINT_TO_ULINT(in UDINT) ULINT { return ULINT(clampULINT(ULINT(in), MAXULINT)) }

// UDINT_TO_BOOL conversion
func UDINT_TO_BOOL(in UDINT) BOOL { return in > 0 }

// UDINT_TO_TIME conversion
func UDINT_TO_TIME(in UDINT) TIME { out, _ := SubTime(in); return out }

// UDINT_TO_INT conversion
func UDINT_TO_INT(in UDINT) INT { return INT(clampULINT(ULINT(in), MAXINT)) }

/*
ULINT_TO * Conversion section
*/
// ULINT_TO_REAL conversion
func ULINT_TO_REAL(in ULINT) REAL { return REAL(in) }

// ULINT_TO_SINT conversion
func ULINT_TO_SINT(in ULINT) SINT { return SINT(clampULINT(in, MAXSINT)) }

// ULINT_TO_LINT conversion
func ULINT_TO_LINT(in ULINT) LINT { return LINT(clampULINT(in, MAXLINT)) }

// ULINT_TO_DINT conversion
func ULINT_TO_DINT(in ULINT) DINT { return DINT(clampULINT(in, MAXDINT)) }

// ULINT_TO_DATE conversion
func ULINT_TO_DATE(in ULINT) DATE { out, _ := SubDate(in); return out }

// ULINT_TO_DWORD conversion
func ULINT_TO_DWORD(in ULINT) DWORD { out, _ := SubDword(in); return out }

// ULINT_TO_DT conversion
func ULINT_TO_DT(in ULINT) DT { out, _ := SubDt(in); return out }

// ULINT_TO_TOD conversion
func ULINT_TO_TOD(in ULINT) TOD { out, _ := SubTod(in); return out }

// ULINT_TO_UDINT conversion
func ULINT_TO_UDINT(in ULINT) UDINT { return UDINT(in) }

// ULINT_TO_WORD conversion
func ULINT_TO_WORD(in ULINT) WORD { out, _ := SubWord(in); return out }

// ULINT_TO_STRING conversion
func ULINT_TO_STRING(in ULINT) STRING { return STRING(strconv.FormatUint(uint64(in), 10)) }

// ULINT_TO_LWORD conversion
func ULINT_TO_LWORD(in ULINT) LWORD { out, _ := SubLword(in); return out }

// ULINT_TO_UINT conversion
func ULINT_TO_UINT(in ULINT) UINT { return UINT(clampULINT(in, MAXUINT)) }

// ULINT_TO_LREAL conversion
func ULINT_TO_LREAL(in ULINT) LREAL { return LREAL(in) }

// ULINT_TO_BYTE conversion
func ULINT_TO_BYTE(in ULINT) BYTE { out, _ := SubByte(in); return out }

// ULINT_TO_USINT conversion
func ULINT_TO_USINT(in ULINT) USINT { return USINT(clampULINT(in, MAXUSINT)) }

// ULINT_TO_BOOL conversion
func ULINT_TO_BOOL(in ULINT) BOOL { return in > 0 }

// ULINT_TO_TIME conversion
func ULINT_TO_TIME(in ULINT) TIME { out, _ := SubTime(in); return out }

// ULINT_TO_INT conversion
func ULINT_TO_INT(in ULINT) INT { return INT(clampULINT(in, MAXINT)) }

/*
DATE_TO * Conversion section
*/

// DATE_TO_REAL conversion
func DATE_TO_REAL(in DATE) REAL { return REAL(time.Time(in).UnixMilli()) }

// DATE_TO_SINT conversion
func DATE_TO_SINT(in DATE) SINT {
	return SINT(clampLINT(LINT(time.Time(in).UnixMilli()), MINSINT, MAXSINT))
}

// DATE_TO_LINT conversion
func DATE_TO_LINT(in DATE) LINT { return LINT(time.Time(in).UnixMilli()) }

// DATE_TO_DINT conversion
func DATE_TO_DINT(in DATE) DINT {
	return DINT(clampLINT(LINT(time.Time(in).UnixMilli()), MINDINT, MAXDINT))
}

// DATE_TO_BYTE conversion
func DATE_TO_BYTE(in DATE) BYTE { out, _ := SubByte(in); return out }

// DATE_TO_WORD conversion
func DATE_TO_WORD(in DATE) WORD { out, _ := SubWord(in); return out }

// DATE_TO_DWORD conversion
func DATE_TO_DWORD(in DATE) DWORD { out, _ := SubDword(in); return out }

// DATE_TO_LWORD conversion
func DATE_TO_LWORD(in DATE) LWORD { out, _ := SubLword(in); return out }

// DATE_TO_UDINT conversion
func DATE_TO_UDINT(in DATE) UDINT {
	return UDINT(clampLINT(LINT(time.Time(in).UnixMilli()), 0, MAXUDINT))
}

// DATE_TO_STRING conversion
func DATE_TO_STRING(in DATE) STRING { return STRING(in.String()) }

// DATE_TO_UINT conversion
func DATE_TO_UINT(in DATE) UINT { return UINT(clampLINT(LINT(time.Time(in).UnixMilli()), 0, MAXUINT)) }

// DATE_TO_LREAL conversion
func DATE_TO_LREAL(in DATE) LREAL { return LREAL(time.Time(in).UnixMilli()) }

// DATE_TO_USINT conversion
func DATE_TO_USINT(in DATE) USINT {
	return USINT(clampLINT(LINT(time.Time(in).UnixMilli()), 0, MAXUSINT))
}

// DATE_TO_ULINT conversion
func DATE_TO_ULINT(in DATE) ULINT { return ULINT(clampLINT(LINT(time.Time(in).UnixMilli()), 0, -1)) }

// DATE_TO_INT conversion
func DATE_TO_INT(in DATE) INT { return INT(clampLINT(LINT(time.Time(in).UnixMilli()), MININT, MAXINT)) }

// DATE_TO_TIME
func DATE_TO_TIME(in DATE) TIME {
	// Converts the DATE (a point in time) to a TIME (duration)
	// representing the milliseconds elapsed since the Unix epoch.
	return TIME(time.Time(in).UnixMilli() * int64(time.Millisecond))
}

/*
DT_TO conversion
*/

// DT_TO_REAL conversion
func DT_TO_REAL(in DT) REAL { return REAL(time.Time(in).UnixMilli()) }

// DT_TO_SINT conversion
func DT_TO_SINT(in DT) SINT {
	return SINT(clampLINT(LINT(time.Time(in).UnixMilli()), MINSINT, MAXSINT))
}

// DT_TO_LINT conversion
func DT_TO_LINT(in DT) LINT { return LINT(time.Time(in).UnixMilli()) }

// DT_TO_DINT conversion
func DT_TO_DINT(in DT) DINT {
	return DINT(clampLINT(LINT(time.Time(in).UnixMilli()), MINDINT, MAXDINT))
}

// DT_TO_DWORD conversion
func DT_TO_DWORD(in DT) DWORD { out, _ := SubDword(in); return out }

// DT_TO_UDINT conversion
func DT_TO_UDINT(in DT) UDINT { return UDINT(clampLINT(LINT(time.Time(in).UnixMilli()), 0, MAXUDINT)) }

// DT_TO_WORD conversion
func DT_TO_WORD(in DT) WORD { out, _ := SubWord(in); return out }

// DT_TO_STRING conversion
func DT_TO_STRING(in DT) STRING { return STRING(in.String()) }

// DT_TO_LWORD conversion
func DT_TO_LWORD(in DT) LWORD { out, _ := SubLword(in); return out }

// DT_TO_UINT conversion
func DT_TO_UINT(in DT) UINT { return UINT(clampLINT(LINT(time.Time(in).UnixMilli()), 0, MAXUINT)) }

// DT_TO_LREAL conversion
func DT_TO_LREAL(in DT) LREAL { return LREAL(time.Time(in).UnixMilli()) }

// DT_TO_BYTE conversion
func DT_TO_BYTE(in DT) BYTE { out, _ := SubByte(in); return out }

// DT_TO_USINT conversion
func DT_TO_USINT(in DT) USINT { return USINT(time.Time(in).UnixMilli()) }

// DT_TO_ULINT conversion
func DT_TO_ULINT(in DT) ULINT { return ULINT(clampLINT(LINT(time.Time(in).UnixMilli()), 0, -1)) }

// DT_TO_INT conversion
func DT_TO_INT(in DT) INT { return INT(clampLINT(LINT(time.Time(in).UnixMilli()), MININT, MAXINT)) }

// DT_TO_DATE extracts the DATE part from a DATE_AND_TIME value.
func DT_TO_DATE(in DT) DATE {
	t := time.Time(in)
	// Returns a new DATE with the time part zeroed out, preserving the location.
	return DATE(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()))
}

// DT_TO_TOD extracts the TIME_OF_DAY part from a DATE_AND_TIME value.
func DT_TO_TOD(in DT) TOD {
	// The standard implies the date part is zeroed out. A simple cast to TOD is sufficient
	// as the interpretation of a TOD value focuses only on the time part.
	return TOD(in)
}

/*
TOD_TO conversion
*/

// TOD_TO_REAL conversion
func TOD_TO_REAL(in TOD) REAL { return REAL(TOD_TO_LINT(in)) }

// TOD_TO_SINT conversion
func TOD_TO_SINT(in TOD) SINT { return SINT(TOD_TO_LINT(in)) }

// TOD_TO_LINT conversion
func TOD_TO_LINT(in TOD) LINT { return LINT(in.CONVERT().Interface().(LINT)) }

// TOD_TO_DINT conversion
func TOD_TO_DINT(in TOD) DINT { return DINT(TOD_TO_LINT(in)) }

// TOD_TO_DWORD conversion
func TOD_TO_DWORD(in TOD) DWORD { out, _ := SubDword(in); return out }

// TOD_TO_UDINT conversion
func TOD_TO_UDINT(in TOD) UDINT { return UDINT(TOD_TO_LINT(in)) }

// TOD_TO_WORD conversion
func TOD_TO_WORD(in TOD) WORD { out, _ := SubWord(in); return out }

// TOD_TO_STRING conversion
func TOD_TO_STRING(in TOD) STRING { return STRING(in.String()) }

// TOD_TO_LWORD conversion
func TOD_TO_LWORD(in TOD) LWORD { out, _ := SubLword(in); return out }

// TOD_TO_UINT conversion
func TOD_TO_UINT(in TOD) UINT { return UINT(TOD_TO_LINT(in)) }

// TOD_TO_LREAL conversion
func TOD_TO_LREAL(in TOD) LREAL { return LREAL(TOD_TO_LINT(in)) }

// TOD_TO_BYTE conversion
func TOD_TO_BYTE(in TOD) BYTE { out, _ := SubByte(in); return out }

// TOD_TO_USINT conversion
func TOD_TO_USINT(in TOD) USINT { return USINT(TOD_TO_LINT(in)) }

// TOD_TO_ULINT conversion
func TOD_TO_ULINT(in TOD) ULINT { return ULINT(TOD_TO_LINT(in)) }

// TOD_TO_INT conversion
func TOD_TO_INT(in TOD) INT { return INT(TOD_TO_LINT(in)) }

/*
TIME_TO conversion
*/

// TIME_TO_REAL conversion
func TIME_TO_REAL(in TIME) REAL { return REAL(time.Duration(in).Milliseconds()) }

// TIME_TO_SINT conversion
func TIME_TO_SINT(in TIME) SINT { return SINT(time.Duration(in).Milliseconds()) }

// TIME_TO_LINT conversion
func TIME_TO_LINT(in TIME) LINT { return LINT(time.Duration(in).Milliseconds()) }

// TIME_TO_DINT conversion
func TIME_TO_DINT(in TIME) DINT { return DINT(time.Duration(in).Milliseconds()) }

// TIME_TO_DWORD conversion
func TIME_TO_DWORD(in TIME) DWORD { out, _ := SubDword(time.Duration(in).Milliseconds()); return out }

// TIME_TO_UDINT conversion
func TIME_TO_UDINT(in TIME) UDINT { return UDINT(time.Duration(in).Milliseconds()) }

// TIME_TO_WORD conversion
func TIME_TO_WORD(in TIME) WORD { out, _ := SubWord(time.Duration(in).Milliseconds()); return out }

// TIME_TO_STRING conversion
func TIME_TO_STRING(in TIME) STRING { return STRING(in.String()) }

// TIME_TO_LWORD conversion
func TIME_TO_LWORD(in TIME) LWORD { out, _ := SubLword(time.Duration(in).Milliseconds()); return out }

// TIME_TO_UINT conversion
func TIME_TO_UINT(in TIME) UINT { return UINT(time.Duration(in).Milliseconds()) }

// TIME_TO_LREAL conversion
func TIME_TO_LREAL(in TIME) LREAL { return LREAL(time.Duration(in).Milliseconds()) }

// TIME_TO_BYTE conversion
func TIME_TO_BYTE(in TIME) BYTE { out, _ := SubByte(time.Duration(in).Milliseconds()); return out }

// TIME_TO_USINT conversion
func TIME_TO_USINT(in TIME) USINT { return USINT(time.Duration(in).Milliseconds()) }

// TIME_TO_ULINT conversion
func TIME_TO_ULINT(in TIME) ULINT { return ULINT(time.Duration(in).Milliseconds()) }

// TIME_TO_INT conversion
func TIME_TO_INT(in TIME) INT { return INT(time.Duration(in).Milliseconds()) }

/*
Math conversion of float to bits and vice versa
*/

// REAL_TO_BITS conversion: floats to uint32 as bits
func REAL_TO_BITS(in REAL) UDINT { return UDINT(math.Float32bits(float32(in))) }

// BITS_TO_REAL converstion: uint32 bits to float
func BITS_TO_REAL(in UDINT) REAL { return REAL(math.Float32frombits(uint32(in))) }

// LREAL_TO_BITS conversion: floats to uint32 as bits
func LREAL_TO_BITS(in LREAL) ULINT { return ULINT(math.Float64bits(float64(in))) }

// BITS_TO_REAL converstion: uint32 bits to float
func BITS_TO_LREAL(in ULINT) LREAL { return LREAL(math.Float64frombits(uint64(in))) }

/*
BCD_TO and TO_BCD conversions
*/

// uintToBCD converts a uint64 to its BCD representation.
func uintToBCD(in uint64) (uint64, error) {
	var res uint64
	var shift uint
	val := in
	if val == 0 {
		return 0, nil
	}
	for val > 0 && shift < 64 {
		digit := val % 10
		res |= (digit << shift)
		val /= 10
		shift += 4
	}
	if val > 0 {
		return 0, fmt.Errorf("uintToBCD: input value %d too large for 64-bit BCD representation", in)
	}
	return res, nil
}

// bcdToUint converts a BCD representation to a uint64.
func bcdToUint(in uint64) (uint64, error) {
	var res uint64
	var factor uint64 = 1
	tempVal := in
	for i := 0; i < 16; i++ { // Process up to 16 nibbles (64 bits)
		nibble := (tempVal >> (i * 4)) & 0xF
		if nibble > 9 {
			return 0, fmt.Errorf("bcdToUint: invalid BCD nibble %d in value 0x%X", nibble, in)
		}
		res += nibble * factor
		if tempVal>>((i+1)*4) == 0 {
			break // Stop if remaining bits are zero
		}
		factor *= 10
	}
	return res, nil
}

func USINT_TO_BCD_BYTE(in USINT) BYTE {
	out, err := uintToBCD(uint64(in))
	if err != nil || out > MAXUSINT {
		// According to some interpretations, invalid BCD should result in 0.
		// Or it could be an error state. Panicking is an option for unrecoverable states.
		panic(fmt.Sprintf("USINT_TO_BCD_BYTE: value %d is invalid for BCD conversion to BYTE: %v", in, err))
	}
	return BYTE(out)
}
func UINT_TO_BCD_WORD(in UINT) WORD     { out, _ := uintToBCD(uint64(in)); return WORD(out) }
func UDINT_TO_BCD_DWORD(in UDINT) DWORD { out, _ := uintToBCD(uint64(in)); return DWORD(out) }
func ULINT_TO_BCD_LWORD(in ULINT) LWORD { out, _ := uintToBCD(uint64(in)); return LWORD(out) }

func BYTE_BCD_TO_USINT(in BYTE) USINT {
	out, err := bcdToUint(uint64(in))
	if err != nil {
		panic(fmt.Sprintf("BYTE_BCD_TO_USINT: invalid BCD value 0x%X: %v", in, err))
	}
	return USINT(clampULINT(ULINT(out), MAXUSINT))
}
func WORD_BCD_TO_UINT(in WORD) UINT     { out, _ := bcdToUint(uint64(in)); return UINT(out) }
func DWORD_BCD_TO_UDINT(in DWORD) UDINT { out, _ := bcdToUint(uint64(in)); return UDINT(out) }
func LWORD_BCD_TO_ULINT(in LWORD) ULINT { out, _ := bcdToUint(uint64(in)); return ULINT(out) }
