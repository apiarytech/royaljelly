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
 */

package royaljelly

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

// ANYINT signed int
type ANYINT int

// USINT signed int 8 bit
type USINT uint8

// UINT signed int 16 bit
type UINT uint16

// UDINT signed int 32 bit
type UDINT uint32

// ULINT signed int 64 bit
type ULINT uint64

// ANYUINT signed int
type ANYUINT uint

// REAL float 32 bit
type REAL float32

// LREAL float64 bit
type LREAL float64

// COMPLEX Real & Imaginary 64 bit
type COMPLEX complex64

// LCOMPLEX Real & Imaginary 128 bit
type LCOMPLEX complex128

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
	d  int
	h  int
	m  int
	s  int
	ms int
}

// ANY_BIT
type ANY_BIT interface {
	~bool |
		~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int8 | ~int16 | ~int32 | ~int64
}

// ANY_DATE
type ANY_DATE interface {
	time.Time
}

// ANY_STRING
type ANY_STRING interface {
	~string | ~rune
}

// ANY_INT
type ANY_INT interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int8 | ~int16 | ~int32 | ~int64
}

// ANY_REAL
type ANY_REAL interface {
	~float32 | ~float64
}

// ANY_NUM
type ANY_NUM interface {
	ANY_INT | ANY_REAL
}

// ANY_MAGNITUDE
type ANY_MAGNITUDE interface {
	ANY_NUM | TIME
}

// ANY
type ANY interface {
	ANY_BIT | ANY_MAGNITUDE | ANY_DATE | ANY_STRING
}

// ANY_COMPLEX
type ANY_COMPLEX interface {
	~complex64 | ~complex128
}
