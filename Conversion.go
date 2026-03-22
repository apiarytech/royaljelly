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

func SubByte(in interface{}) (out BYTE) {
	val, err := anyToULINT(in)
	if err != nil {
		panic(fmt.Sprintf("SubByte: conversion error: %v", err))
	}
	return BYTE(val)
}

func SubWord(in interface{}) (out WORD) {
	val, err := anyToULINT(in)
	if err != nil {
		panic(fmt.Sprintf("SubWord: conversion error: %v", err))
	}
	return WORD(val)
}

func SubDword(in interface{}) (out DWORD) {
	val, err := anyToULINT(in)
	if err != nil {
		panic(fmt.Sprintf("SubDword: conversion error: %v", err))
	}
	return DWORD(val)
}

func SubLword(in interface{}) (out LWORD) {
	val, err := anyToULINT(in)
	if err != nil {
		panic(fmt.Sprintf("SubLword: conversion error: %v", err))
	}
	return LWORD(val)
}

func SubDt(in interface{}) (out DT) {
	val, err := anyToLINT(in)
	if err != nil {
		panic(fmt.Sprintf("SubDt: conversion error: %v", err))
	}
	// Assuming the integer value represents milliseconds since Unix epoch
	out = DT(time.UnixMilli(int64(val)))
	return out
}

func SubDate(in interface{}) (out DATE) {
	val, err := anyToLINT(in)
	if err != nil {
		panic(fmt.Sprintf("SubDate: conversion error: %v", err))
	}
	// Assuming the integer value represents milliseconds since Unix epoch
	out = DATE(time.UnixMilli(int64(val)))
	return out
}

func SubTod(in interface{}) (out TOD) {
	// Assuming input is milliseconds since midnight
	val, err := anyToLINT(in)
	if err != nil {
		panic(fmt.Sprintf("SubTod: conversion error: %v", err))
	}
	// A TOD is a duration since midnight on an arbitrary day.
	out = TOD(time.Time{}.Add(time.Duration(val) * time.Millisecond))
	return out
}

func SubTime(in interface{}) (out TIME) {
	val, err := anyToLINT(in)
	if err != nil {
		panic(fmt.Sprintf("SubTime: conversion error: %v", err))
	}
	out = TIME(time.Duration(val) * time.Millisecond)
	return out
}

func SubBool(in BOOL) (out int) {
	if in {
		out = 1
	} else {
		out = 0
	}
	return out
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
func BOOL_TO_BYTE(in BOOL) BYTE { return SubByte(in) }

// BOOL_TO_WORD conversion
func BOOL_TO_WORD(in BOOL) WORD { return SubWord(in) }

// BOOL_TO_DWORD conversion
func BOOL_TO_DWORD(in BOOL) DWORD { return SubDword(in) }

// BOOL_TO_LWORD conversion
func BOOL_TO_LWORD(in BOOL) LWORD { return SubLword(in) }

// BOOL_TO_INT conversion
func BOOL_TO_INT(in BOOL) (out INT) { return INT(SubBool(in)) }

// BOOL_TO_SINT conversion
func BOOL_TO_SINT(in BOOL) (out SINT) { return SINT(SubBool(in)) }

// BOOL_TO_UINT conversion
func BOOL_TO_UINT(in BOOL) (out UINT) { return UINT(SubBool(in)) }

// BOOL_TO_UDINT conversion
func BOOL_TO_UDINT(in BOOL) (out UDINT) { return UDINT(SubBool(in)) }

// BOOL_TO_DINT conversion
func BOOL_TO_DINT(in BOOL) (out DINT) { return DINT(SubBool(in)) }

// BOOL_TO_USINT conversion
func BOOL_TO_USINT(in BOOL) (out USINT) { return USINT(SubBool(in)) }

// BOOL_TO_ULINT conversion
func BOOL_TO_ULINT(in BOOL) (out ULINT) { return ULINT(SubBool(in)) }

// BOOL_TO_LINT conversion
func BOOL_TO_LINT(in BOOL) (out LINT) { return LINT(SubBool(in)) }

// BOOL_TO_STRING conversion
func BOOL_TO_STRING(in BOOL) (out STRING) { return STRING(fmt.Sprint(in)) }

// BOOL_TO_REAL conversion
func BOOL_TO_REAL(in BOOL) (out REAL) { return REAL(SubBool(in)) }

// BOOL_TO_LREAL conversion
func BOOL_TO_LREAL(in BOOL) (out LREAL) { return LREAL(SubBool(in)) }

// BOOL_TO_TIME conversion
func BOOL_TO_TIME(in BOOL) (out TIME) { return SubTime(in) }

// BOOL_TO_TOD conversion
func BOOL_TO_TOD(in BOOL) (out TOD) { return SubTod(in) }

// BOOL_TO_DATE conversion
func BOOL_TO_DATE(in BOOL) (out DATE) { return SubDate(in) }

// BOOL_TO_DT conversion
func BOOL_TO_DT(in BOOL) (out DT) { return SubDt(in) }

/*
BYTE_TO * Conversion Section
*/

// BYTE_TO_BOOL conversion
func BYTE_TO_BOOL(in BYTE) BOOL { return (in.Value() > 0) }

// BYTE_TO_SINT conversion
func BYTE_TO_SINT(in BYTE) SINT { return SINT(in.Value()) }

// BYTE_TO_INT conversion
func BYTE_TO_INT(in BYTE) INT { return INT(in.Value()) }

// BYTE_TO_DINT conversion
func BYTE_TO_DINT(in BYTE) DINT { return DINT(in.Value()) }

// BYTE_TO_LINT conversion
func BYTE_TO_LINT(in BYTE) LINT { return LINT(in.Value()) }

// BYTE_TO_ANYINT conversion
func BYTE_TO_ANYINT(in BYTE) ANYINT { return ANYINT(in.Value()) }

// BYTE_TO_USINT conversion
func BYTE_TO_USINT(in BYTE) USINT { return USINT(in.Value()) }

// BYTE_TO_UINT conversion
func BYTE_TO_UINT(in BYTE) UINT { return UINT(in.Value()) }

// BYTE_TO_Udint conversion
func BYTE_TO_UDINT(in BYTE) UDINT { return UDINT(in.Value()) }

// BYTE_TO_ULINT conversion
func BYTE_TO_ULINT(in BYTE) ULINT { return ULINT(in.Value()) }

// BYTE_TO_ANYUINT conversion
func BYTE_TO_ANYUINT(in BYTE) ANYUINT { return ANYUINT(in.Value()) }

// BYTE_TO_REAL conversion
func BYTE_TO_REAL(in BYTE) REAL { return REAL(in.Value()) }

// BYTE_TO_LREAL conversion
func BYTE_TO_LREAL(in BYTE) LREAL { return LREAL(in.Value()) }

// BYTE_TO_WORD conversion
func BYTE_TO_WORD(in BYTE) WORD { return SubWord(in.Value()) }

// BYTE_TO_DWORD conversion
func BYTE_TO_DWORD(in BYTE) DWORD { return SubDword(in.Value()) }

// BYTE_TO_LWORD conversion
func BYTE_TO_LWORD(in BYTE) LWORD { return SubLword(in.Value()) }

// BYTE_TO_STRING conversion
func BYTE_TO_STRING(in BYTE) STRING { return STRING(strconv.FormatUint(uint64(in), 10)) }

// BYTES_TO_STRING conversion
func BYTES_TO_STRING(in []BYTE) (out STRING) {
	byteSlice := make([]byte, len(in))
	for i, b := range in {
		byteSlice[i] = b.Value()
	}
	return STRING(byteSlice)
}

// BYTE_TO_DATE conversion
func BYTE_TO_DATE(in BYTE) DATE { return SubDate(in) }

// BYTE_TO_DT conversion
func BYTE_TO_DT(in BYTE) DT { return SubDt(in) }

// BYTE_TO_TOD conversion
func BYTE_TO_TOD(in BYTE) TOD { return SubTod(in) }

// BYTE_TO_TIME conversion
func BYTE_TO_TIME(in BYTE) TIME { return SubTime(in) }

/*
WORD_TO * Conversion Section
*/

// WORD_TO_BOOL covnersion
func WORD_TO_BOOL(in WORD) BOOL { return (in.Value() > 0) }

// WORD_TO_SINT conversion
func WORD_TO_SINT(in WORD) SINT { return SINT(in.Value()) }

// WORD_TO_INT conversion
func WORD_TO_INT(in WORD) INT { return INT(in.Value()) }

// WORD_TO_DINT conversion
func WORD_TO_DINT(in WORD) DINT { return DINT(in.Value()) }

// WORD_TO_DATE conversion
func WORD_TO_DATE(in WORD) DATE { return SubDate(in) }

// WORD_TO_Lint conversion
func WORD_TO_LINT(in WORD) LINT { return LINT(in.Value()) }

// WORD_TO_AnyInt conversion
func WORD_TO_AnyInt(in WORD) ANYINT { return ANYINT(in.Value()) }

// WORD_TO_USINT conversion
func WORD_TO_USINT(in WORD) USINT { return USINT(in.Value()) }

// WORD_TO_UINT conversion
func WORD_TO_UINT(in WORD) UINT { return UINT(in.Value()) }

// WORD_TO_UDINT conversion
func WORD_TO_UDINT(in WORD) UDINT { return UDINT(in.Value()) }

// WORD_TO_ULINT conversion
func WORD_TO_ULINT(in WORD) ULINT { return ULINT(in.Value()) }

// WORD_TO_ANYUINT conversion
func WORD_TO_ANYUINT(in WORD) ANYUINT { return ANYUINT(in.Value()) }

// WORD_TO_REAL conversion
func WORD_TO_REAL(in WORD) REAL { return REAL(in.Value()) }

// WORD_TO_LREAL conversion
func WORD_TO_LREAL(in WORD) LREAL { return LREAL(in.Value()) }

// WORD_TO_BYTE conversion
func WORD_TO_BYTE(in WORD) BYTE { return SubByte(in) }

// WORD_TO_DWORD conversion
func WORD_TO_DWORD(in WORD) DWORD { return SubDword(in) }

// WORD_TO_DT conversion
func WORD_TO_DT(in WORD) DT { return SubDt(in) }

// WORD_TO_TOD conversion
func WORD_TO_TOD(in WORD) TOD { return SubTod(in) }

// WORD_TO_LWORD conversion
func WORD_TO_LWORD(in WORD) LWORD { return SubLword(in) }

// WORD_TO_STRING conversion
func WORD_TO_STRING(in WORD) STRING { return STRING(strconv.FormatUint(uint64(in.Value()), 10)) }

// WORD_TO_TIM conversion
func WORD_TO_TIME(in WORD) TIME { return SubTime(in) }

/*
DWORD_TO * Conversion Section
*/

// DWORD_TO_BOOL covnersion
func DWORD_TO_BOOL(in DWORD) BOOL { return in.Value() > 0 }

// DWORD_TO_SINT conversion
func DWORD_TO_SINT(in DWORD) SINT { return SINT(in.Value()) }

// DWORD_TO_INT conversion
func DWORD_TO_INT(in DWORD) INT { return INT(in.Value()) }

// DWORD_TO_DINT conversion
func DWORD_TO_DINT(in DWORD) DINT { return DINT(in.Value()) }

// DWORD_TO_DATE conversion
func DWORD_TO_DATE(in DWORD) DATE { return SubDate(in) }

// DWORD_TO_DT conversion
func DWORD_TO_DT(in DWORD) DT { return SubDt(in) }

// DWORD_TO_TOD conversion
func DWORD_TO_TOD(in DWORD) TOD { return SubTod(in) }

// DWORD_TO_LINT conversion
func DWORD_TO_LINT(in DWORD) LINT { return LINT(in.Value()) }

// DWORD_TO_ANYINT conversion
func DWORD_TO_ANYINT(in DWORD) ANYINT { return ANYINT(in.Value()) }

// DWORD_TO_USINT conversion
func DWORD_TO_USINT(in DWORD) USINT { return USINT(in.Value()) }

// DWORD_TO_UINT conversion
func DWORD_TO_UINT(in DWORD) UINT { return UINT(in.Value()) }

// DWORD_TO_UDINT conversion
func DWORD_TO_UDINT(in DWORD) UDINT { return UDINT(in.Value()) }

// DWORD_TO_ULINT conversion
func DWORD_TO_ULINT(in DWORD) ULINT { return ULINT(in.Value()) }

// DWORD_TO_ANYUINT conversion
func DWORD_TO_ANYUINT(in DWORD) ANYUINT { return ANYUINT(in.Value()) }

// DWORD_TO_REAL conversion
func DWORD_TO_REAL(in DWORD) REAL { return REAL(in.Value()) }

// DWORD_TO_LREAL conversion
func DWORD_TO_LREAL(in DWORD) LREAL { return LREAL(in.Value()) }

// DWORD_TO_BYTE conversion
func DWORD_TO_BYTE(in DWORD) BYTE { return SubByte(in) }

// DWORD_TO_WORD conversion
func DWORD_TO_WORD(in DWORD) WORD { return SubWord(in) }

// DWORD_TO_LWORD conversion
func DWORD_TO_LWORD(in DWORD) LWORD { return SubLword(in) }

// DWORD_TO_STRING conversion
func DWORD_TO_STRING(in DWORD) STRING { return STRING(strconv.FormatUint(uint64(in.Value()), 10)) }

// DWORD_TO_TIME conversion
func DWORD_TO_TIME(in DWORD) TIME { return SubTime(in) }

/*
LWORD_TO * Conversion Section
*/

// LWORD_TO_BOOL covnersion
func LWORD_TO_BOOL(in LWORD) (out BOOL) { return in.Value() > 0 }

// LWORD_TO_SINT conversion
func LWORD_TO_SINT(in LWORD) SINT { return SINT(in.Value()) }

// LWORD_TO_INT conversion
func LWORD_TO_INT(in LWORD) INT { return INT(in.Value()) }

// LWORD_TO_DINT conversion
func LWORD_TO_DINT(in LWORD) DINT { return DINT(in.Value()) }

// LWORD_TO_LINT conversion
func LWORD_TO_LINT(in LWORD) LINT { return LINT(in.Value()) }

// LWORD_TO_ANYINT conversion
func LWORD_TO_ANYINT(in LWORD) ANYINT { return ANYINT(in.Value()) }

// LWORD_TO_USINT conversion
func LWORD_TO_USINT(in LWORD) USINT { return USINT(in.Value()) }

// LWORD_TO_UINT conversion
func LWORD_TO_UINT(in LWORD) UINT { return UINT(in.Value()) }

// LWORD_TO_UDINT conversion
func LWORD_TO_UDINT(in LWORD) UDINT { return UDINT(in.Value()) }

// LWORD_TO_ULINT conversion
func LWORD_TO_ULINT(in LWORD) ULINT { return ULINT(in.Value()) }

// LWORD_TO_ANYUINT conversion
func LWORD_TO_ANYUINT(in LWORD) ANYUINT { return ANYUINT(in.Value()) }

// LWORD_TO_REAL conversion
func LWORD_TO_REAL(in LWORD) REAL { return REAL(in.Value()) }

// LWORD_TO_LREAL conversion
func LWORD_TO_LREAL(in LWORD) LREAL { return LREAL(in.Value()) }

// LWORD_TO_BYTE conversion
func LWORD_TO_BYTE(in LWORD) BYTE { return SubByte(in) }

// LWORD_TO_WORD conversion
func LWORD_TO_WORD(in LWORD) WORD { return SubWord(in) }

// LWORD_TO_DWORD conversion
func LWORD_TO_DWORD(in LWORD) DWORD { return SubDword(in) }

// LWORD_TO_STRING conversion
func LWORD_TO_STRING(in LWORD) STRING { return STRING(strconv.FormatUint(in.Value(), 10)) }

// LWORD_TO_DATE conversion
func LWORD_TO_DATE(in LWORD) DATE { return SubDate(in) }

// LWORD_TO_DT conversion
func LWORD_TO_DT(in LWORD) DT { return SubDt(in) }

// LWORD_TO_TOD conversion
func LWORD_TO_TOD(in LWORD) TOD { return SubTod(in) }

// LWORD_TO_TIME converion
func LWORD_TO_TIME(in LWORD) TIME { return SubTime(in) }

/*
REAL_TO * Conversion Section
*/

// REAL_TO_SINT conversion
func REAL_TO_SINT(in REAL) SINT { return SINT(math.RoundToEven(float64(in))) }

// REAL_TO_LINT conversion
func REAL_TO_LINT(in REAL) LINT { return LINT(math.RoundToEven(float64(in))) }

// REAL_TO_DINT conversion
func REAL_TO_DINT(in REAL) DINT { return DINT(math.RoundToEven(float64(in))) }

// REAL_TO_DATE conversion
func REAL_TO_DATE(in REAL) (out DATE) { return SubDate(in) }

// REAL_TO_DWORD conversion
func REAL_TO_DWORD(in REAL) (out DWORD) { return SubDword(in) }

// REAL_TO_DT conversion
func REAL_TO_DT(in REAL) (out DT) { return SubDt(in) }

// REAL_TO_TOD conversion
func REAL_TO_TOD(in REAL) (out TOD) { return SubTod(in) }

// REAL_TO_UDINT conversion
func REAL_TO_UDINT(in REAL) (out UDINT) { return UDINT(math.RoundToEven(float64(in))) }

// REAL_TO_WORD conversion
func REAL_TO_WORD(in REAL) WORD { return SubWord(in) }

// REAL_TO_STRING conversion
func REAL_TO_STRING(in REAL) STRING { return STRING(strconv.FormatFloat(float64(in), 'g', -1, 32)) }

// REAL_TO_LWORD conversion
func REAL_TO_LWORD(in REAL) LWORD { return SubLword(in) }

// REAL_TO_UINT conversion
func REAL_TO_UINT(in REAL) UINT { return UINT(math.RoundToEven(float64(in))) }

// REAL_TO_LREAL conversion
func REAL_TO_LREAL(in REAL) LREAL { return LREAL(in) }

// REAL_TO_BYTE conversion
func REAL_TO_BYTE(in REAL) BYTE { return SubByte(in) }

// REAL_TO_USINT conversion
func REAL_TO_USINT(in REAL) USINT { return USINT(math.RoundToEven(float64(in))) }

// REAL_TO_ULINT conversion
func REAL_TO_ULINT(in REAL) ULINT { return ULINT(math.RoundToEven(float64(in))) }

// REAL_TO_BOOL conversion
func REAL_TO_BOOL(in REAL) BOOL { return in > 0 }

// REAL_TO_TIME conversion
func REAL_TO_TIME(in REAL) (out TIME) { return SubTime(in) }

// REAL_TO_INT conversion
func REAL_TO_INT(in REAL) INT { return INT(math.RoundToEven(float64(in))) }

/*
LREAL_TO * Conversion Section
*/
// LREAL_TO_REAL conversion
func LREAL_TO_REAL(in LREAL) (out REAL) { return REAL(in) }

// LREAL_TO_SINT conversion
func LREAL_TO_SINT(in LREAL) (out SINT) { return SINT(math.RoundToEven(float64(in))) }

// LREAL_TO_LINT conversion
func LREAL_TO_LINT(in LREAL) (out LINT) { return LINT(math.RoundToEven(float64(in))) }

// LREAL_TO_DINT conversion
func LREAL_TO_DINT(in LREAL) (out DINT) { return DINT(math.RoundToEven(float64(in))) }

// LREAL_TO_DATE conversion
func LREAL_TO_DATE(in LREAL) (out DATE) { return SubDate(in) }

// LREAL_TO_DWORD conversion
func LREAL_TO_DWORD(in LREAL) (out DWORD) { return SubDword(in) }

// LREAL_TO_DT conversion
func LREAL_TO_DT(in LREAL) (out DT) { return SubDt(in) }

// LREAL_TO_TOD conversion
func LREAL_TO_TOD(in LREAL) (out TOD) { return SubTod(in) }

// LREAL_TO_UDINT conversion
func LREAL_TO_UDINT(in LREAL) (out UDINT) { return UDINT(math.RoundToEven(float64(in))) }

// LREAL_TO_WORD conversion
func LREAL_TO_WORD(in LREAL) (out WORD) { return SubWord(in) }

// LREAL_TO_STRING conversion
func LREAL_TO_STRING(in LREAL) (out STRING) {
	return STRING(strconv.FormatFloat(float64(in), 'g', -1, 64))
}

// LREAL_TO_LWORD conversion
func LREAL_TO_LWORD(in LREAL) (out LWORD) { return SubLword(in) }

// LREAL_TO_UINT conversion
func LREAL_TO_UINT(in LREAL) (out UINT) { return UINT(math.RoundToEven(float64(in))) }

// LREAL_TO_BYTE conversion
func LREAL_TO_BYTE(in LREAL) (out BYTE) { return SubByte(in) }

// LREAL_TO_USINT conversion
func LREAL_TO_USINT(in LREAL) (out USINT) { return USINT(math.RoundToEven(float64(in))) }

// LREAL_TO_ULINT conversion
func LREAL_TO_ULINT(in LREAL) (out ULINT) { return ULINT(math.RoundToEven(float64(in))) }

// LREAL_TO_BOOL conversion
func LREAL_TO_BOOL(in LREAL) (out BOOL) { return in > 0 }

// LREAL_TO_TIME conversion
func LREAL_TO_TIME(in LREAL) (out TIME) { return SubTime(in) }

// LREAL_TO_INT conversion
func LREAL_TO_INT(in LREAL) (out INT) { return INT(math.RoundToEven(float64(in))) }

/*
SINT_TO * Conversion section
*/

// SINT_TO_REAL conversion
func SINT_TO_REAL(in SINT) REAL { return REAL(in) }

// SINT_TO_LINT conversion
func SINT_TO_LINT(in SINT) LINT { return LINT(in) }

// SINT_TO_DINT conversion
func SINT_TO_DINT(in SINT) DINT { return DINT(in) }

// SINT_TO_DATE conversion
func SINT_TO_DATE(in SINT) (out DATE) { return SubDate(in) }

// SINT_TO_DWORD conversion
func SINT_TO_DWORD(in SINT) DWORD { return SubDword(in) }

// SINT_TO_DT conversion
func SINT_TO_DT(in SINT) (out DT) { return SubDt(in) }

// SINT_TO_TOD conversion
func SINT_TO_TOD(in SINT) (out TOD) { return SubTod(in) }

// SINT_TO_UDINT conversion
func SINT_TO_UDINT(in SINT) UDINT { return UDINT(in) }

// SINT_TO_WORD conversion
func SINT_TO_WORD(in SINT) WORD { return SubWord(in) }

// SINT_TO_STRING conversion
func SINT_TO_STRING(in SINT) STRING { return STRING(strconv.FormatInt(int64(in), 10)) }

// SINT_TO_LWORD conversion
func SINT_TO_LWORD(in SINT) LWORD { return SubLword(in) }

// SINT_TO_UINT conversion
func SINT_TO_UINT(in SINT) UINT { return UINT(in) }

// SINT_TO_LREAL conversion
func SINT_TO_LREAL(in SINT) LREAL { return LREAL(in) }

// SINT_TO_BYTE conversion
func SINT_TO_BYTE(in SINT) BYTE { return SubByte(in) }

// SINT_TO_USINT conversion
func SINT_TO_USINT(in SINT) USINT { return USINT(in) }

// SINT_TO_ULINT conversion
func SINT_TO_ULINT(in SINT) ULINT { return ULINT(in) }

// SINT_TO_BOOL conversion
func SINT_TO_BOOL(in SINT) BOOL { return BOOL(in != 0) }

// SINT_TO_TIME conversion
func SINT_TO_TIME(in SINT) (out TIME) { return SubTime(in) }

// SINT_TO_INT conversion
func SINT_TO_INT(in SINT) INT { return INT(in) }

/*
INT_TO * Conversion section
*/

// INT_TO_REAL conversion
func INT_TO_REAL(in INT) REAL { return REAL(in) }

// INT_TO_SINT conversion
func INT_TO_SINT(in INT) SINT { return SINT(in) }

// INT_TO_LINT conversion
func INT_TO_LINT(in INT) LINT { return LINT(in) }

// INT_TO_DINT conversion
func INT_TO_DINT(in INT) DINT { return DINT(in) }

// INT_TO_DATE conversion seconds to DATE
func INT_TO_DATE(in INT) (out DATE) { return SubDate(in) }

// INT_TO_DWORD conversion
func INT_TO_DWORD(in INT) (out DWORD) { return SubDword(in) }

// INT_TO_DT conversion
func INT_TO_DT(in INT) (out DT) { return SubDt(in) }

// INT_TO_TOD conversion
func INT_TO_TOD(in INT) (out TOD) { return SubTod(in) }

// INT_TO_UDINT conversion
func INT_TO_UDINT(in INT) UDINT { return UDINT(in) }

// INT_TO_WORD conversion
func INT_TO_WORD(in INT) WORD { return SubWord(in) }

// INT_TO_STRING conversion
func INT_TO_STRING(in INT) STRING { return STRING(strconv.FormatInt(int64(in), 10)) }

// INT_TO_LWORD conversion
func INT_TO_LWORD(in INT) LWORD { return SubLword(in) }

// INT_TO_UINT conversion
func INT_TO_UINT(in INT) UINT { return UINT(in) }

// INT_TO_LREAL conversion
func INT_TO_LREAL(in INT) LREAL { return LREAL(in) }

// INT_TO_BYTE conversion
func INT_TO_BYTE(in INT) BYTE { return SubByte(in) }

// INT_TO_USINT conversion
func INT_TO_USINT(in INT) USINT { return USINT(in) }

// INT_TO_ULINT conversion
func INT_TO_ULINT(in INT) ULINT { return ULINT(in) }

// INT_TO_BOOL conversion
func INT_TO_BOOL(in INT) BOOL { return in > 0 }

// INT_TO_TIME conversion
func INT_TO_TIME(in INT) (out TIME) { return SubTime(in) }

/*
LINT_TO * Conversion section
*/

// LINT_TO_REAL conversion
func LINT_TO_REAL(in LINT) (out REAL) { return REAL(in) }

// LINT_TO_SINT conversion
func LINT_TO_SINT(in LINT) SINT { return SINT(in) }

// LINT_TO_DINT conversion
func LINT_TO_DINT(in LINT) DINT { return DINT(in) }

// LINT_TO_DWORD conversion
func LINT_TO_DWORD(in LINT) DWORD { return SubDword(in) }

// LINT_TO_DATE conversion
func LINT_TO_DATE(in LINT) (out DATE) { return SubDate(in) }

// LINT_TO_TOD conversion
func LINT_TO_TOD(in LINT) (out TOD) { return SubTod(in) }

// LINT_TO_UDINT conversion
func LINT_TO_UDINT(in LINT) UDINT { return UDINT(in) }

// LINT_TO_WORD conversion
func LINT_TO_WORD(in LINT) WORD { return SubWord(in) }

// LINT_TO_STRING conversion
func LINT_TO_STRING(in LINT) STRING { return STRING(strconv.FormatInt(int64(in), 10)) }

// LINT_TO_LWORD conversion
func LINT_TO_LWORD(in LINT) LWORD { return SubLword(in) }

// LINT_TO_UINT conversion
func LINT_TO_UINT(in LINT) UINT { return UINT(in) }

// LINT_TO_LREAL conversion
func LINT_TO_LREAL(in LINT) LREAL { return LREAL(in) }

// LINT_TO_BYTE conversion
func LINT_TO_BYTE(in LINT) BYTE { return SubByte(in) }

// LINT_TO_USINT conversion
func LINT_TO_USINT(in LINT) USINT { return USINT(in) }

// LINT_TO_ULINT conversion
func LINT_TO_ULINT(in LINT) (out ULINT) { return ULINT(in) }

// LINT_TO_BOOL conversion
func LINT_TO_BOOL(in LINT) BOOL { return in > 0 }

// LINT_TO_TIME conversion
func LINT_TO_TIME(in LINT) (out TIME) { return SubTime(in) }

// LINT_TO_INT conversion
func LINT_TO_INT(in LINT) INT { return INT(in) }

/*
// DINT_TO Conversions
*/

// DINT_TO_REAL conversion
func DINT_TO_REAL(in DINT) REAL { return REAL(in) }

// DINT_TO_SINT conversion
func DINT_TO_SINT(in DINT) SINT { return SINT(in) }

// DINT_TO_LINT conversion
func DINT_TO_LINT(in DINT) LINT { return LINT(in) }

// DINT_TO_DATE conversion
func DINT_TO_DATE(in DINT) (out DATE) { return SubDate(in) }

// DINT_TO_DWORD conversion
func DINT_TO_DWORD(in DINT) DWORD { return SubDword(in) }

// DINT_TO_DT conversion
func DINT_TO_DT(in DINT) (out DT) { return SubDt(in) }

// DINT_TO_TOD conversion
func DINT_TO_TOD(in DINT) (out TOD) { return SubTod(in) }

// DINT_TO_UDINT conversion
func DINT_TO_UDINT(in DINT) UDINT { return UDINT(in) }

// DINT_TO_WORD conversion
func DINT_TO_WORD(in DINT) WORD { return SubWord(in) }

// DINT_TO_STRING conversion
func DINT_TO_STRING(in DINT) STRING { return STRING(strconv.FormatInt(int64(in), 10)) }

// DINT_TO_LWORD conversion
func DINT_TO_LWORD(in DINT) LWORD { return SubLword(in) }

// DINT_TO_UINT conversion
func DINT_TO_UINT(in DINT) UINT { return UINT(in) }

// DINT_TO_LREAL conversion
func DINT_TO_LREAL(in DINT) LREAL { return LREAL(in) }

// DINT_TO_BYTE conversion
func DINT_TO_BYTE(in DINT) BYTE { return SubByte(in) }

// DINT_TO_USINT conversion
func DINT_TO_USINT(in DINT) (out USINT) { return USINT(in) }

// DINT_TO_ULINT conversion
func DINT_TO_ULINT(in DINT) (out ULINT) { return ULINT(in) }

// DINT_TO_BOOL conversion
func DINT_TO_BOOL(in DINT) BOOL { return in > 0 }

// DINT_TO_TIME conversion
func DINT_TO_TIME(in DINT) (out TIME) { return SubTime(in) }

// DINT_TO_INT conversion
func DINT_TO_INT(in DINT) INT { return INT(in) }

/*
USINT_TO * Conversion section
*/
// USINT_TO_REAL conversion
func USINT_TO_REAL(in USINT) (out REAL) { return REAL(in) }

// USINT_TO_SINT conversion
func USINT_TO_SINT(in USINT) (out SINT) { return SINT(in) }

// USINT_TO_LINT conversion
func USINT_TO_LINT(in USINT) (out LINT) { return LINT(in) }

// USINT_TO_DINT conversion
func USINT_TO_DINT(in USINT) (out DINT) { return DINT(in) }

// USINT_TO_DATE conversion
func USINT_TO_DATE(in USINT) (out DATE) { return SubDate(in) }

// USINT_TO_DWORD conversion
func USINT_TO_DWORD(in USINT) (out DWORD) { return SubDword(in) }

// USINT_TO_DT conversion
func USINT_TO_DT(in USINT) (out DT) { return SubDt(in) }

// USINT_TO_TOD conversion
func USINT_TO_TOD(in USINT) (out TOD) { return SubTod(in) }

// USINT_TO_UDINT conversion
func USINT_TO_UDINT(in USINT) (out UDINT) { return UDINT(in) }

// USINT_TO_WORD conversion
func USINT_TO_WORD(in USINT) (out WORD) { return SubWord(in) }

// USINT_TO_STRING conversion
func USINT_TO_STRING(in USINT) (out STRING) { return STRING(strconv.FormatUint(uint64(in), 10)) }

// USINT_TO_LWORD conversion
func USINT_TO_LWORD(in USINT) (out LWORD) { return SubLword(in) }

// USINT_TO_UINT conversion
func USINT_TO_UINT(in USINT) (out UINT) { return UINT(in) }

// USINT_TO_LREAL conversion
func USINT_TO_LREAL(in USINT) (out LREAL) { return LREAL(in) }

// USINT_TO_BYTE conversion
func USINT_TO_BYTE(in USINT) (out BYTE) { return SubByte(in) }

// USINT_TO_ULINT conversion
func USINT_TO_ULINT(in USINT) (out ULINT) { return ULINT(in) }

// USINT_TO_BOOL conversion
func USINT_TO_BOOL(in USINT) (out BOOL) { return in > 0 }

// USINT_TO_TIME conversion
func USINT_TO_TIME(in USINT) (out TIME) { return SubTime(in) }

// USINT_TO_INT conversion
func USINT_TO_INT(in USINT) (out INT) { return INT(in) }

/*
UINT_TO conversion
*/
// UINT_TO_REAL conversion
func UINT_TO_REAL(in UINT) (out REAL) { return REAL(in) }

// UINT_TO_SINT conversion
func UINT_TO_SINT(in UINT) (out SINT) { return SINT(in) }

// UINT_TO_LINT conversion
func UINT_TO_LINT(in UINT) (out LINT) { return LINT(in) }

// UINT_TO_DINT conversion
func UINT_TO_DINT(in UINT) (out DINT) { return DINT(in) }

// UINT_TO_DATE conversion
func UINT_TO_DATE(in UINT) (out DATE) { return SubDate(in) }

// UINT_TO_DWORD conversion
func UINT_TO_DWORD(in UINT) (out DWORD) { return SubDword(in) }

// UINT_TO_DT conversion
func UINT_TO_DT(in UINT) (out DT) { return SubDt(in) }

// UINT_TO_TOD conversion
func UINT_TO_TOD(in UINT) (out TOD) { return SubTod(in) }

// UINT_TO_UDINT conversion
func UINT_TO_UDINT(in UINT) (out UDINT) { return UDINT(in) }

// UINT_TO_WORD conversion
func UINT_TO_WORD(in UINT) (out WORD) { return SubWord(in) }

// UINT_TO_STRING conversion
func UINT_TO_STRING(in UINT) (out STRING) { return STRING(strconv.FormatUint(uint64(in), 10)) }

// UINT_TO_LWORD conversion
func UINT_TO_LWORD(in UINT) (out LWORD) { return SubLword(in) }

// UINT_TO_LREAL conversion
func UINT_TO_LREAL(in UINT) (out LREAL) { return LREAL(in) }

// UINT_TO_BYTE conversion
func UINT_TO_BYTE(in UINT) (out BYTE) { return SubByte(in) }

// UINT_TO_USINT conversion
func UINT_TO_USINT(in UINT) (out USINT) { return USINT(in) }

// UINT_TO_ULINT conversion
func UINT_TO_ULINT(in UINT) (out ULINT) { return ULINT(in) }

// UINT_TO_BOOL conversion
func UINT_TO_BOOL(in UINT) (out BOOL) { return in > 0 }

// UINT_TO_TIME conversion
func UINT_TO_TIME(in UINT) (out TIME) { return SubTime(in) }

// UINT_TO_INT conversion
func UINT_TO_INT(in UINT) (out INT) { return INT(in) }

/*
UDINT_TO * Conversion section
*/
// UDINT_TO_REAL conversion
func UDINT_TO_REAL(in UDINT) (out REAL) { return REAL(in) }

// UDINT_TO_SINT conversion
func UDINT_TO_SINT(in UDINT) (out SINT) { return SINT(in) }

// UDINT_TO_LINT conversion
func UDINT_TO_LINT(in UDINT) (out LINT) { return LINT(in) }

// UDINT_TO_DINT conversion
func UDINT_TO_DINT(in UDINT) (out DINT) { return DINT(in) }

// UDINT_TO_DATE conversion
func UDINT_TO_DATE(in UDINT) (out DATE) { return SubDate(in) }

// UDINT_TO_DWORD conversion
func UDINT_TO_DWORD(in UDINT) DWORD { return SubDword(in) }

// UDINT_TO_DT conversion
func UDINT_TO_DT(in UDINT) (out DT) { return SubDt(in) }

// UDINT_TO_TOD conversion
func UDINT_TO_TOD(in UDINT) (out TOD) { return SubTod(in) }

// UDINT_TO_WORD conversion
func UDINT_TO_WORD(in UDINT) WORD { return SubWord(in) }

// UDINT_TO_STRING conversion
func UDINT_TO_STRING(in UDINT) (out STRING) { return STRING(strconv.FormatUint(uint64(in), 10)) }

// UDINT_TO_LWORD conversion
func UDINT_TO_LWORD(in UDINT) (out LWORD) { return SubLword(in) }

// UDINT_TO_UINT conversion
func UDINT_TO_UINT(in UDINT) (out UINT) { return UINT(in) }

// UDINT_TO_LREAL conversion
func UDINT_TO_LREAL(in UDINT) (out LREAL) { return LREAL(in) }

// UDINT_TO_BYTE conversion
func UDINT_TO_BYTE(in UDINT) (out BYTE) { return SubByte(in) }

// UDINT_TO_USINT conversion
func UDINT_TO_USINT(in UDINT) (out USINT) { return USINT(in) }

// UDINT_TO_ULINT conversion
func UDINT_TO_ULINT(in UDINT) (out ULINT) { return ULINT(in) }

// UDINT_TO_BOOL conversion
func UDINT_TO_BOOL(in UDINT) (out BOOL) { return in > 0 }

// UDINT_TO_TIME conversion
func UDINT_TO_TIME(in UDINT) (out TIME) { return SubTime(in) }

// UDINT_TO_INT conversion
func UDINT_TO_INT(in UDINT) (out INT) { return INT(in) }

/*
ULINT_TO * Conversion section
*/
// ULINT_TO_REAL conversion
func ULINT_TO_REAL(in ULINT) (out REAL) { return REAL(in) }

// ULINT_TO_SINT conversion
func ULINT_TO_SINT(in ULINT) (out SINT) { return SINT(in) }

// ULINT_TO_LINT conversion
func ULINT_TO_LINT(in ULINT) (out LINT) { return LINT(in) }

// ULINT_TO_DINT conversion
func ULINT_TO_DINT(in ULINT) (out DINT) { return DINT(in) }

// ULINT_TO_DATE conversion
func ULINT_TO_DATE(in ULINT) (out DATE) { return SubDate(in) }

// ULINT_TO_DWORD conversion
func ULINT_TO_DWORD(in ULINT) (out DWORD) { return SubDword(in) }

// ULINT_TO_DT conversion
func ULINT_TO_DT(in ULINT) (out DT) { return SubDt(in) }

// ULINT_TO_TOD conversion
func ULINT_TO_TOD(in ULINT) (out TOD) { return SubTod(in) }

// ULINT_TO_UDINT conversion
func ULINT_TO_UDINT(in ULINT) (out UDINT) { return UDINT(in) }

// ULINT_TO_WORD conversion
func ULINT_TO_WORD(in ULINT) (out WORD) { return SubWord(in) }

// ULINT_TO_STRING conversion
func ULINT_TO_STRING(in ULINT) (out STRING) { return STRING(strconv.FormatUint(uint64(in), 10)) }

// ULINT_TO_LWORD conversion
func ULINT_TO_LWORD(in ULINT) (out LWORD) { return SubLword(in) }

// ULINT_TO_UINT conversion
func ULINT_TO_UINT(in ULINT) (out UINT) { return UINT(in) }

// ULINT_TO_LREAL conversion
func ULINT_TO_LREAL(in ULINT) (out LREAL) { return LREAL(in) }

// ULINT_TO_BYTE conversion
func ULINT_TO_BYTE(in ULINT) (out BYTE) { return SubByte(in) }

// ULINT_TO_USINT conversion
func ULINT_TO_USINT(in ULINT) (out USINT) { return USINT(in) }

// ULINT_TO_BOOL conversion
func ULINT_TO_BOOL(in ULINT) (out BOOL) { return in > 0 }

// ULINT_TO_TIME conversion
func ULINT_TO_TIME(in ULINT) (out TIME) { return SubTime(in) }

// ULINT_TO_INT conversion
func ULINT_TO_INT(in ULINT) (out INT) { return INT(in) }

/*
DATE_TO * Conversion section
*/

// DATE_TO_REAL conversion
func DATE_TO_REAL(in DATE) (out REAL) { return REAL(time.Time(in).UnixMilli()) }

// DATE_TO_SINT conversion
func DATE_TO_SINT(in DATE) (out SINT) { return SINT(time.Time(in).UnixMilli()) }

// DATE_TO_LINT conversion
func DATE_TO_LINT(in DATE) (out LINT) { return LINT(time.Time(in).UnixMilli()) }

// DATE_TO_DINT conversion
func DATE_TO_DINT(in DATE) (out DINT) { return DINT(time.Time(in).UnixMilli()) }

// DATE_TO_BYTE conversion
func DATE_TO_BYTE(in DATE) (out BYTE) { return SubByte(in) }

// DATE_TO_WORD conversion
func DATE_TO_WORD(in DATE) (out WORD) { return SubWord(in) }

// DATE_TO_DWORD conversion
func DATE_TO_DWORD(in DATE) (out DWORD) { return SubDword(in) }

// DATE_TO_LWORD conversion
func DATE_TO_LWORD(in DATE) (out LWORD) { return SubLword(in) }

// DATE_TO_UDINT conversion
func DATE_TO_UDINT(in DATE) (out UDINT) { return UDINT(time.Time(in).UnixMilli()) }

// DATE_TO_STRING conversion
func DATE_TO_STRING(in DATE) (out STRING) { return STRING(in.String()) }

// DATE_TO_UINT conversion
func DATE_TO_UINT(in DATE) (out UINT) { return UINT(time.Time(in).UnixMilli()) }

// DATE_TO_LREAL conversion
func DATE_TO_LREAL(in DATE) (out LREAL) { return LREAL(time.Time(in).UnixMilli()) }

// DATE_TO_USINT conversion
func DATE_TO_USINT(in DATE) (out USINT) { return USINT(time.Time(in).UnixMilli()) }

// DATE_TO_ULINT conversion
func DATE_TO_ULINT(in DATE) (out ULINT) { return ULINT(time.Time(in).UnixMilli()) }

// DATE_TO_INT conversion
func DATE_TO_INT(in DATE) (out INT) { return INT(time.Time(in).UnixMilli()) }

// DATE_TO_TIME
func DATE_TO_TIME(in DATE) (out TIME) {
	// Converts the DATE (a point in time) to a TIME (duration)
	// representing the milliseconds elapsed since the Unix epoch.
	return TIME(time.Time(in).UnixMilli() * int64(time.Millisecond))
}

/*
DT_TO conversion
*/

// DT_TO_REAL conversion
func DT_TO_REAL(in DT) (out REAL) { return REAL(time.Time(in).UnixMilli()) }

// DT_TO_SINT conversion
func DT_TO_SINT(in DT) (out SINT) { return SINT(time.Time(in).UnixMilli()) }

// DT_TO_LINT conversion
func DT_TO_LINT(in DT) (out LINT) { return LINT(time.Time(in).UnixMilli()) }

// DT_TO_DINT conversion
func DT_TO_DINT(in DT) (out DINT) { return DINT(time.Time(in).UnixMilli()) }

// DT_TO_DWORD conversion
func DT_TO_DWORD(in DT) (out DWORD) { return SubDword(in) }

// DT_TO_UDINT conversion
func DT_TO_UDINT(in DT) (out UDINT) { return UDINT(time.Time(in).UnixMilli()) }

// DT_TO_WORD conversion
func DT_TO_WORD(in DT) (out WORD) { return SubWord(in) }

// DT_TO_STRING conversion
func DT_TO_STRING(in DT) (out STRING) { return STRING(in.String()) }

// DT_TO_LWORD conversion
func DT_TO_LWORD(in DT) (out LWORD) { return SubLword(in) }

// DT_TO_UINT conversion
func DT_TO_UINT(in DT) (out UINT) { return UINT(time.Time(in).UnixMilli()) }

// DT_TO_LREAL conversion
func DT_TO_LREAL(in DT) (out LREAL) { return LREAL(time.Time(in).UnixMilli()) }

// DT_TO_BYTE conversion
func DT_TO_BYTE(in DT) (out BYTE) { return SubByte(in) }

// DT_TO_USINT conversion
func DT_TO_USINT(in DT) (out USINT) { return USINT(time.Time(in).UnixMilli()) }

// DT_TO_ULINT conversion
func DT_TO_ULINT(in DT) (out ULINT) { return ULINT(time.Time(in).UnixMilli()) }

// DT_TO_INT conversion
func DT_TO_INT(in DT) (out INT) { return INT(time.Time(in).UnixMilli()) }

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
func TOD_TO_REAL(in TOD) (out REAL) { return REAL(TOD_TO_LINT(in)) }

// TOD_TO_SINT conversion
func TOD_TO_SINT(in TOD) (out SINT) { return SINT(TOD_TO_LINT(in)) }

// TOD_TO_LINT conversion
func TOD_TO_LINT(in TOD) (out LINT) { return LINT(in.CONVERT().Interface().(LINT)) }

// TOD_TO_DINT conversion
func TOD_TO_DINT(in TOD) (out DINT) { return DINT(TOD_TO_LINT(in)) }

// TOD_TO_DWORD conversion
func TOD_TO_DWORD(in TOD) (out DWORD) { return SubDword(in) }

// TOD_TO_UDINT conversion
func TOD_TO_UDINT(in TOD) (out UDINT) { return UDINT(TOD_TO_LINT(in)) }

// TOD_TO_WORD conversion
func TOD_TO_WORD(in TOD) (out WORD) { return SubWord(in) }

// TOD_TO_STRING conversion
func TOD_TO_STRING(in TOD) (out STRING) { return STRING(in.String()) }

// TOD_TO_LWORD conversion
func TOD_TO_LWORD(in TOD) (out LWORD) { return SubLword(in) }

// TOD_TO_UINT conversion
func TOD_TO_UINT(in TOD) (out UINT) { return UINT(TOD_TO_LINT(in)) }

// TOD_TO_LREAL conversion
func TOD_TO_LREAL(in TOD) (out LREAL) { return LREAL(TOD_TO_LINT(in)) }

// TOD_TO_BYTE conversion
func TOD_TO_BYTE(in TOD) (out BYTE) { return SubByte(in) }

// TOD_TO_USINT conversion
func TOD_TO_USINT(in TOD) (out USINT) { return USINT(TOD_TO_LINT(in)) }

// TOD_TO_ULINT conversion
func TOD_TO_ULINT(in TOD) (out ULINT) { return ULINT(TOD_TO_LINT(in)) }

// TOD_TO_INT conversion
func TOD_TO_INT(in TOD) (out INT) { return INT(TOD_TO_LINT(in)) }

/*
TIME_TO conversion
*/

// TIME_TO_REAL conversion
func TIME_TO_REAL(in TIME) (out REAL) { return REAL(time.Duration(in).Milliseconds()) }

// TIME_TO_SINT conversion
func TIME_TO_SINT(in TIME) (out SINT) { return SINT(time.Duration(in).Milliseconds()) }

// TIME_TO_LINT conversion
func TIME_TO_LINT(in TIME) (out LINT) { return LINT(time.Duration(in).Milliseconds()) }

// TIME_TO_DINT conversion
func TIME_TO_DINT(in TIME) (out DINT) { return DINT(time.Duration(in).Milliseconds()) }

// TIME_TO_DWORD conversion
func TIME_TO_DWORD(in TIME) (out DWORD) { return SubDword(time.Duration(in).Milliseconds()) }

// TIME_TO_UDINT conversion
func TIME_TO_UDINT(in TIME) (out UDINT) { return UDINT(time.Duration(in).Milliseconds()) }

// TIME_TO_WORD conversion
func TIME_TO_WORD(in TIME) (out WORD) { return SubWord(time.Duration(in).Milliseconds()) }

// TIME_TO_STRING conversion
func TIME_TO_STRING(in TIME) (out STRING) { return STRING(in.String()) }

// TIME_TO_LWORD conversion
func TIME_TO_LWORD(in TIME) (out LWORD) { return SubLword(time.Duration(in).Milliseconds()) }

// TIME_TO_UINT conversion
func TIME_TO_UINT(in TIME) (out UINT) { return UINT(time.Duration(in).Milliseconds()) }

// TIME_TO_LREAL conversion
func TIME_TO_LREAL(in TIME) (out LREAL) { return LREAL(time.Duration(in).Milliseconds()) }

// TIME_TO_BYTE conversion
func TIME_TO_BYTE(in TIME) (out BYTE) { return SubByte(time.Duration(in).Milliseconds()) }

// TIME_TO_USINT conversion
func TIME_TO_USINT(in TIME) (out USINT) { return USINT(time.Duration(in).Milliseconds()) }

// TIME_TO_ULINT conversion
func TIME_TO_ULINT(in TIME) (out ULINT) { return ULINT(time.Duration(in).Milliseconds()) }

// TIME_TO_INT conversion
func TIME_TO_INT(in TIME) (out INT) { return INT(time.Duration(in).Milliseconds()) }

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
func uintToBCD(in uint64) uint64 {
	var res uint64
	var shift uint
	val := in
	if val == 0 {
		return 0
	}
	for val > 0 && shift < 64 {
		digit := val % 10
		res |= (digit << shift)
		val /= 10
		shift += 4
	}
	if val > 0 {
		panic("uintToBCD: input value too large for 64-bit BCD representation")
	}
	return res
}

// bcdToUint converts a BCD representation to a uint64.
func bcdToUint(in uint64) uint64 {
	var res uint64
	var factor uint64 = 1
	tempVal := in
	for i := 0; i < 16; i++ { // Process up to 16 nibbles (64 bits)
		nibble := (tempVal >> (i * 4)) & 0xF
		if nibble > 9 {
			panic(fmt.Sprintf("bcdToUint: invalid BCD nibble %d", nibble))
		}
		res += nibble * factor
		if tempVal>>((i+1)*4) == 0 {
			break // Stop if remaining bits are zero
		}
		factor *= 10
	}
	return res
}

func USINT_TO_BCD_BYTE(in USINT) BYTE   { return BYTE(uintToBCD(uint64(in))) }
func UINT_TO_BCD_WORD(in UINT) WORD     { return WORD(uintToBCD(uint64(in))) }
func UDINT_TO_BCD_DWORD(in UDINT) DWORD { return DWORD(uintToBCD(uint64(in))) }
func ULINT_TO_BCD_LWORD(in ULINT) LWORD { return LWORD(uintToBCD(uint64(in))) }

func BYTE_BCD_TO_USINT(in BYTE) USINT   { return USINT(bcdToUint(uint64(in))) }
func WORD_BCD_TO_UINT(in WORD) UINT     { return UINT(bcdToUint(uint64(in))) }
func DWORD_BCD_TO_UDINT(in DWORD) UDINT { return UDINT(bcdToUint(uint64(in))) }
func LWORD_BCD_TO_ULINT(in LWORD) ULINT { return ULINT(bcdToUint(uint64(in))) }
