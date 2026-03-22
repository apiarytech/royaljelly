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
	"math"
	"time"
)

// Math package variables
const (
	E   = math.E
	PI  = math.Pi
	PHI = math.Phi

	SQRT2   = math.Sqrt2
	SQRTE   = math.SqrtE
	SQRTPI  = math.SqrtPi
	SQRTPHI = math.SqrtPhi

	LN2    = math.Ln2
	LOG2E  = math.Log2E
	LN10   = math.Ln10
	LOG10E = math.Log10E

	MAXREAL              = math.MaxFloat32
	SMALLESTNONZEROREAL  = math.SmallestNonzeroFloat32
	MAXLREAL             = math.MaxFloat64
	SMALLESTNONZEROLREAL = math.SmallestNonzeroFloat64

	MAXANYINT = math.MaxInt
	MAXSINT   = math.MaxInt8
	MAXINT    = math.MaxInt16
	MAXDINT   = math.MaxInt32
	MAXLINT   = math.MaxInt64
	MINANYINT = math.MinInt
	MINSINT   = math.MinInt8
	MININT    = math.MinInt16
	MINDINT   = math.MinInt32
	MINLINT   = math.MinInt64

	MAXANYUINT = math.MaxUint
	MAXUSINT   = math.MaxUint8
	MAXUINT    = math.MaxUint16
	MAXUDINT   = math.MaxUint32
	MAXULINT   = math.MaxUint64

	INITBOOL  BOOL  = false
	INITBYTE  BYTE  = 0
	INITWORD  WORD  = 0
	INITDWORD DWORD = 0
	INITLWORD LWORD = 0
	INITSINT  SINT  = 0
	INITINT   INT   = 0
	INITDINT  DINT  = 0
	INITLINT  LINT  = 0
	INITUSINT USINT = 0
	INITUINT  UINT  = 0
	INITUDINT UDINT = 0
	INITULINT ULINT = 0
	INITREAL  REAL  = 0.0
	INITLREAL LREAL = 0.0

	// Special character constants for use in strings, as per IEC 61131-3 Table 6.
	DOLLAR       = STRING("$$")   // Dollar sign ($)
	SINGLE_QUOTE = STRING("$'")   // Single quote (')
	DOUBLE_QUOTE = STRING("$\"")  // Double quote (")
	CR           = STRING("$R")   // Carriage return
	LF           = STRING("$L")   // Line feed
	CRLF         = STRING("$R$L") // Carriage return + Line feed
	NEWLINE      = STRING("$N")   // Newline
	FORM_FEED    = STRING("$P")   // Form feed
	TAB          = STRING("$T")   // Tab
)

// Default initial values for non-constant types
var (
	// TIME is a duration, so TIME(0) is a valid zero duration (T#0s).
	INITTIME = TIME(0)
	// DATE, TOD, and DT are based on time.Time. Their zero value is time.Time{}.
	// This corresponds to the IEC default of 0001-01-01 00:00:00.
	INITDATE    = DATE(time.Time{})
	INITTOD     = TOD(time.Time{})
	INITDT      = DT(time.Time{})
	INITSTRING  = STRING("")
	INITWSTRING = WSTRING(' ')
)
