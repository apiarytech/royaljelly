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

	INITBOOL    BOOL    = false
	INITBYTE    BYTE    = 0
	INITWORD    WORD    = 0
	INITDWORD   DWORD   = 0
	INITLWORD   LWORD   = 0
	INITSINT    SINT    = 0
	INITINT     INT     = 0
	INITDINT    DINT    = 0
	INITLINT    LINT    = 0
	INITUSINT   USINT   = 0
	INITUINT    UINT    = 0
	INITUDINT   UDINT   = 0
	INITULINT   ULINT   = 0
	INITREAL    REAL    = 0.0
	INITLREAL   LREAL   = 0.0
	INITSTRING  STRING  = ""
	INITWSTRING WSTRING = ' '
	INITTIME    TIME    = TIME(0)
	INITDATE    DATE    = DATE(time.Time{})
	INITTOD     TOD     = TOD(time.Time{})
	INITDT      DT      = DT(time.Time{})

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
