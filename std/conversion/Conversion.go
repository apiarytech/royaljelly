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

package conversion

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/apiarytech/royaljelly/convert"
	"github.com/apiarytech/royaljelly/iec"
)

/*
BOOL_TO conversion
*/

// BOOL_TO_BYTE conversion
func BOOL_TO_BYTE(in iec.BOOL) iec.BYTE { out, _ := convert.SubByte(in); return out }

// BOOL_TO_WORD conversion
func BOOL_TO_WORD(in iec.BOOL) iec.WORD { out, _ := convert.SubWord(in); return out }

// BOOL_TO_DWORD conversion
func BOOL_TO_DWORD(in iec.BOOL) iec.DWORD { out, _ := convert.SubDword(in); return out }

// BOOL_TO_LWORD conversion
func BOOL_TO_LWORD(in iec.BOOL) iec.LWORD { out, _ := convert.SubLword(in); return out }

// BOOL_TO_INT conversion
func BOOL_TO_INT(in iec.BOOL) iec.INT { return iec.INT(convert.SubBool(in)) }

// BOOL_TO_SINT conversion
func BOOL_TO_SINT(in iec.BOOL) iec.SINT {
	return iec.SINT(convert.ClampLINT(iec.LINT(convert.SubBool(in)), iec.MINSINT, iec.MAXSINT))
}

// BOOL_TO_UINT conversion
func BOOL_TO_UINT(in iec.BOOL) iec.UINT {
	return iec.UINT(convert.ClampULINT(iec.ULINT(convert.SubBool(in)), iec.MAXUINT))
}

// BOOL_TO_UDINT conversion
func BOOL_TO_UDINT(in iec.BOOL) iec.UDINT { return iec.UDINT(convert.SubBool(in)) }

// BOOL_TO_DINT conversion
func BOOL_TO_DINT(in iec.BOOL) iec.DINT { return iec.DINT(convert.SubBool(in)) }

// BOOL_TO_USINT conversion
func BOOL_TO_USINT(in iec.BOOL) iec.USINT {
	return iec.USINT(convert.ClampULINT(iec.ULINT(convert.SubBool(in)), iec.MAXUSINT))
}

// BOOL_TO_ULINT conversion
func BOOL_TO_ULINT(in iec.BOOL) iec.ULINT {
	return iec.ULINT(convert.ClampULINT(iec.ULINT(convert.SubBool(in)), iec.MAXULINT))
}

// BOOL_TO_LINT conversion
func BOOL_TO_LINT(in iec.BOOL) iec.LINT { return iec.LINT(convert.SubBool(in)) }

// BOOL_TO_STRING conversion
func BOOL_TO_STRING(in iec.BOOL) (out iec.STRING) { return iec.STRING(fmt.Sprint(in)) }

// BOOL_TO_REAL conversion
func BOOL_TO_REAL(in iec.BOOL) iec.REAL { return iec.REAL(convert.SubBool(in)) }

// BOOL_TO_LREAL conversion
func BOOL_TO_LREAL(in iec.BOOL) iec.LREAL { return iec.LREAL(convert.SubBool(in)) }

// BOOL_TO_TIME conversion
func BOOL_TO_TIME(in iec.BOOL) iec.TIME { out, _ := convert.SubTime(in); return out }

// BOOL_TO_TOD conversion
func BOOL_TO_TOD(in iec.BOOL) iec.TOD { out, _ := convert.SubTod(in); return out }

// BOOL_TO_DATE conversion
func BOOL_TO_DATE(in iec.BOOL) iec.DATE { out, _ := convert.SubDate(in); return out }

// BOOL_TO_DT conversion
func BOOL_TO_DT(in iec.BOOL) iec.DT { out, _ := convert.SubDt(in); return out }

/*
BYTE_TO * Conversion Section
*/

// BYTE_TO_BOOL conversion
func BYTE_TO_BOOL(in iec.BYTE) iec.BOOL { return (in.Value() > 0) }

// BYTE_TO_SINT conversion
func BYTE_TO_SINT(in iec.BYTE) iec.SINT {
	return iec.SINT(convert.ClampULINT(iec.ULINT(in), iec.MAXSINT))
}

// BYTE_TO_INT conversion
func BYTE_TO_INT(in iec.BYTE) iec.INT { return iec.INT(convert.ClampULINT(iec.ULINT(in), iec.MAXINT)) }

// BYTE_TO_DINT conversion
func BYTE_TO_DINT(in iec.BYTE) iec.DINT {
	return iec.DINT(convert.ClampULINT(iec.ULINT(in), iec.MAXDINT))
}

// BYTE_TO_LINT conversion
func BYTE_TO_LINT(in iec.BYTE) iec.LINT { return iec.LINT(in.Value()) }

// BYTE_TO_USINT conversion
func BYTE_TO_USINT(in iec.BYTE) iec.USINT { return iec.USINT(in.Value()) }

// BYTE_TO_UINT conversion
func BYTE_TO_UINT(in iec.BYTE) iec.UINT {
	return iec.UINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUINT))
}

// BYTE_TO_Udint conversion
func BYTE_TO_UDINT(in iec.BYTE) iec.UDINT {
	return iec.UDINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUDINT))
}

// BYTE_TO_ULINT conversion
func BYTE_TO_ULINT(in iec.BYTE) iec.ULINT { return iec.ULINT(in.Value()) }

// BYTE_TO_REAL conversion
func BYTE_TO_REAL(in iec.BYTE) iec.REAL { return iec.REAL(in.Value()) }

// BYTE_TO_LREAL conversion
func BYTE_TO_LREAL(in iec.BYTE) iec.LREAL { return iec.LREAL(in.Value()) }

// BYTE_TO_WORD conversion
func BYTE_TO_WORD(in iec.BYTE) iec.WORD { return iec.WORD(in) }

// BYTE_TO_DWORD conversion
func BYTE_TO_DWORD(in iec.BYTE) iec.DWORD { return iec.DWORD(in) }

// BYTE_TO_LWORD conversion
func BYTE_TO_LWORD(in iec.BYTE) iec.LWORD { return iec.LWORD(in) }

// BYTE_TO_STRING conversion
func BYTE_TO_STRING(in iec.BYTE) iec.STRING { return iec.STRING(strconv.FormatUint(uint64(in), 10)) }

// BYTES_TO_STRING conversion
func BYTES_TO_STRING(in []iec.BYTE) iec.STRING {
	byteSlice := make([]byte, len(in))
	for i, b := range in {
		byteSlice[i] = b.Value()
	}
	return iec.STRING(byteSlice)
}

// BYTE_TO_DATE conversion
func BYTE_TO_DATE(in iec.BYTE) iec.DATE { out, _ := convert.SubDate(in); return out }

// BYTE_TO_DT conversion
func BYTE_TO_DT(in iec.BYTE) iec.DT { out, _ := convert.SubDt(in); return out }

// BYTE_TO_TOD conversion
func BYTE_TO_TOD(in iec.BYTE) iec.TOD { out, _ := convert.SubTod(in); return out }

// BYTE_TO_TIME conversion
func BYTE_TO_TIME(in iec.BYTE) iec.TIME { out, _ := convert.SubTime(in); return out }

/*
WORD_TO * Conversion Section
*/

// WORD_TO_BOOL covnersion
func WORD_TO_BOOL(in iec.WORD) iec.BOOL { return (in.Value() > 0) }

// WORD_TO_SINT conversion
func WORD_TO_SINT(in iec.WORD) iec.SINT {
	return iec.SINT(convert.ClampULINT(iec.ULINT(in), iec.MAXSINT))
}

// WORD_TO_INT conversion
func WORD_TO_INT(in iec.WORD) iec.INT { return iec.INT(convert.ClampULINT(iec.ULINT(in), iec.MAXINT)) }

// WORD_TO_DINT conversion
func WORD_TO_DINT(in iec.WORD) iec.DINT {
	return iec.DINT(convert.ClampULINT(iec.ULINT(in), iec.MAXDINT))
}

// WORD_TO_DATE conversion
func WORD_TO_DATE(in iec.WORD) iec.DATE { out, _ := convert.SubDate(in); return out }

// WORD_TO_Lint conversion
func WORD_TO_LINT(in iec.WORD) iec.LINT {
	return iec.LINT(convert.ClampULINT(iec.ULINT(in), iec.MAXLINT))
}

// WORD_TO_USINT conversion
func WORD_TO_USINT(in iec.WORD) iec.USINT {
	return iec.USINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUSINT))
}

// WORD_TO_UINT conversion
func WORD_TO_UINT(in iec.WORD) iec.UINT { return iec.UINT(in.Value()) }

// WORD_TO_UDINT conversion
func WORD_TO_UDINT(in iec.WORD) iec.UDINT {
	return iec.UDINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUDINT))
}

// WORD_TO_ULINT conversion
func WORD_TO_ULINT(in iec.WORD) iec.ULINT { return iec.ULINT(in.Value()) }

// WORD_TO_REAL conversion
func WORD_TO_REAL(in iec.WORD) iec.REAL {
	return iec.REAL(convert.ClampLREAL(iec.LREAL(in), -iec.MAXREAL, iec.MAXREAL))
}

// WORD_TO_LREAL conversion
func WORD_TO_LREAL(in iec.WORD) iec.LREAL { return iec.LREAL(in.Value()) }

// WORD_TO_BYTE conversion
func WORD_TO_BYTE(in iec.WORD) iec.BYTE {
	return iec.BYTE(convert.ClampULINT(iec.ULINT(in), iec.MAXUSINT))
}

// WORD_TO_DWORD conversion
func WORD_TO_DWORD(in iec.WORD) iec.DWORD { return iec.DWORD(in) }

// WORD_TO_DT conversion
func WORD_TO_DT(in iec.WORD) iec.DT { out, _ := convert.SubDt(in); return out }

// WORD_TO_TOD conversion
func WORD_TO_TOD(in iec.WORD) iec.TOD { out, _ := convert.SubTod(in); return out }

// WORD_TO_LWORD conversion
func WORD_TO_LWORD(in iec.WORD) iec.LWORD { return iec.LWORD(in) }

// WORD_TO_STRING conversion
func WORD_TO_STRING(in iec.WORD) iec.STRING {
	return iec.STRING(strconv.FormatUint(uint64(in.Value()), 10))
}

// WORD_TO_TIM conversion
func WORD_TO_TIME(in iec.WORD) iec.TIME { out, _ := convert.SubTime(in); return out }

/*
DWORD_TO * Conversion Section
*/

// DWORD_TO_BOOL covnersion
func DWORD_TO_BOOL(in iec.DWORD) iec.BOOL { return in.Value() > 0 }

// DWORD_TO_SINT conversion
func DWORD_TO_SINT(in iec.DWORD) iec.SINT {
	return iec.SINT(convert.ClampULINT(iec.ULINT(in), iec.MAXSINT))
}

// DWORD_TO_INT conversion
func DWORD_TO_INT(in iec.DWORD) iec.INT {
	return iec.INT(convert.ClampULINT(iec.ULINT(in), iec.MAXINT))
}

// DWORD_TO_DINT conversion
func DWORD_TO_DINT(in iec.DWORD) iec.DINT {
	return iec.DINT(convert.ClampULINT(iec.ULINT(in), iec.MAXDINT))
}

// DWORD_TO_DATE conversion
func DWORD_TO_DATE(in iec.DWORD) iec.DATE { out, _ := convert.SubDate(in); return out }

// DWORD_TO_DT conversion
func DWORD_TO_DT(in iec.DWORD) iec.DT { out, _ := convert.SubDt(in); return out }

// DWORD_TO_TOD conversion
func DWORD_TO_TOD(in iec.DWORD) iec.TOD { out, _ := convert.SubTod(in); return out }

// DWORD_TO_LINT conversion
func DWORD_TO_LINT(in iec.DWORD) iec.LINT {
	return iec.LINT(convert.ClampULINT(iec.ULINT(in), iec.MAXLINT))
}

// DWORD_TO_USINT conversion
func DWORD_TO_USINT(in iec.DWORD) iec.USINT {
	return iec.USINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUSINT))
}

// DWORD_TO_UINT conversion
func DWORD_TO_UINT(in iec.DWORD) iec.UINT {
	return iec.UINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUINT))
}

// DWORD_TO_UDINT conversion
func DWORD_TO_UDINT(in iec.DWORD) iec.UDINT { return iec.UDINT(in.Value()) }

// DWORD_TO_ULINT conversion
func DWORD_TO_ULINT(in iec.DWORD) iec.ULINT {
	return iec.ULINT(convert.ClampULINT(iec.ULINT(in), iec.MAXULINT))
}

// DWORD_TO_REAL conversion
func DWORD_TO_REAL(in iec.DWORD) iec.REAL {
	return iec.REAL(convert.ClampLREAL(iec.LREAL(in), -iec.MAXREAL, iec.MAXREAL))
}

// DWORD_TO_LREAL conversion
func DWORD_TO_LREAL(in iec.DWORD) iec.LREAL { return iec.LREAL(in.Value()) }

// DWORD_TO_BYTE conversion
func DWORD_TO_BYTE(in iec.DWORD) iec.BYTE {
	return iec.BYTE(convert.ClampULINT(iec.ULINT(in), iec.MAXUSINT))
}

// DWORD_TO_WORD conversion
func DWORD_TO_WORD(in iec.DWORD) iec.WORD {
	return iec.WORD(convert.ClampULINT(iec.ULINT(in), iec.MAXUINT))
}

// DWORD_TO_LWORD conversion
func DWORD_TO_LWORD(in iec.DWORD) iec.LWORD { return iec.LWORD(in) }

// DWORD_TO_STRING conversion
func DWORD_TO_STRING(in iec.DWORD) iec.STRING {
	return iec.STRING(strconv.FormatUint(uint64(in.Value()), 10))
}

// DWORD_TO_TIME conversion
func DWORD_TO_TIME(in iec.DWORD) iec.TIME { out, _ := convert.SubTime(in); return out }

/*
LWORD_TO * Conversion Section
*/

// LWORD_TO_BOOL covnersion
func LWORD_TO_BOOL(in iec.LWORD) iec.BOOL { return in.Value() > 0 }

// LWORD_TO_SINT conversion
func LWORD_TO_SINT(in iec.LWORD) iec.SINT {
	return iec.SINT(convert.ClampULINT(iec.ULINT(in), iec.MAXSINT))
}

// LWORD_TO_INT conversion
func LWORD_TO_INT(in iec.LWORD) iec.INT {
	return iec.INT(convert.ClampULINT(iec.ULINT(in), iec.MAXINT))
}

// LWORD_TO_DINT conversion
func LWORD_TO_DINT(in iec.LWORD) iec.DINT {
	return iec.DINT(convert.ClampULINT(iec.ULINT(in), iec.MAXDINT))
}

// LWORD_TO_LINT conversion
func LWORD_TO_LINT(in iec.LWORD) iec.LINT {
	return iec.LINT(convert.ClampULINT(iec.ULINT(in), iec.MAXLINT))
}

// LWORD_TO_USINT conversion
func LWORD_TO_USINT(in iec.LWORD) iec.USINT {
	return iec.USINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUSINT))
}

// LWORD_TO_UINT conversion
func LWORD_TO_UINT(in iec.LWORD) iec.UINT {
	return iec.UINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUINT))
}

// LWORD_TO_UDINT conversion
func LWORD_TO_UDINT(in iec.LWORD) iec.UDINT {
	return iec.UDINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUDINT))
}

// LWORD_TO_ULINT conversion
func LWORD_TO_ULINT(in iec.LWORD) iec.ULINT { return iec.ULINT(in.Value()) }

// LWORD_TO_REAL conversion
func LWORD_TO_REAL(in iec.LWORD) iec.REAL {
	return iec.REAL(convert.ClampLREAL(iec.LREAL(in), -iec.MAXREAL, iec.MAXREAL))
}

// LWORD_TO_LREAL conversion
func LWORD_TO_LREAL(in iec.LWORD) iec.LREAL { return iec.LREAL(in.Value()) }

// LWORD_TO_BYTE conversion
func LWORD_TO_BYTE(in iec.LWORD) iec.BYTE {
	return iec.BYTE(convert.ClampULINT(iec.ULINT(in), iec.MAXUSINT))
}

// LWORD_TO_WORD conversion
func LWORD_TO_WORD(in iec.LWORD) iec.WORD {
	return iec.WORD(convert.ClampULINT(iec.ULINT(in), iec.MAXUINT))
}

// LWORD_TO_DWORD conversion
func LWORD_TO_DWORD(in iec.LWORD) iec.DWORD {
	return iec.DWORD(convert.ClampULINT(iec.ULINT(in), iec.MAXUDINT))
}

// LWORD_TO_STRING conversion
func LWORD_TO_STRING(in iec.LWORD) iec.STRING { return iec.STRING(strconv.FormatUint(in.Value(), 10)) }

// LWORD_TO_DATE conversion
func LWORD_TO_DATE(in iec.LWORD) iec.DATE { out, _ := convert.SubDate(in); return out }

// LWORD_TO_DT conversion
func LWORD_TO_DT(in iec.LWORD) iec.DT { out, _ := convert.SubDt(in); return out }

// LWORD_TO_TOD conversion
func LWORD_TO_TOD(in iec.LWORD) iec.TOD { out, _ := convert.SubTod(in); return out }

// LWORD_TO_TIME converion
func LWORD_TO_TIME(in iec.LWORD) iec.TIME { out, _ := convert.SubTime(in); return out }

/*
REAL_TO * Conversion Section
*/

// REAL_TO_SINT conversion
func REAL_TO_SINT(in iec.REAL) iec.SINT {
	return iec.SINT(convert.RoundAndClampLREAL(iec.LREAL(in), iec.MINSINT, iec.MAXSINT))
}

// REAL_TO_LINT conversion
func REAL_TO_LINT(in iec.REAL) iec.LINT {
	return iec.LINT(convert.RoundAndClampLREAL(iec.LREAL(in), iec.MINLINT, iec.MAXLINT))
}

// REAL_TO_DINT conversion
func REAL_TO_DINT(in iec.REAL) iec.DINT {
	return iec.DINT(convert.RoundAndClampLREAL(iec.LREAL(in), iec.MINDINT, iec.MAXDINT))
}

// REAL_TO_DATE conversion
func REAL_TO_DATE(in iec.REAL) iec.DATE { out, _ := convert.SubDate(in); return out }

// REAL_TO_DWORD conversion
func REAL_TO_DWORD(in iec.REAL) iec.DWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.DWORD(convert.ClampULINT(val, iec.MAXUDINT))
}

// REAL_TO_DT conversion
func REAL_TO_DT(in iec.REAL) iec.DT { out, _ := convert.SubDt(in); return out }

// REAL_TO_TOD conversion
func REAL_TO_TOD(in iec.REAL) iec.TOD { out, _ := convert.SubTod(in); return out }

// REAL_TO_UDINT conversion
func REAL_TO_UDINT(in iec.REAL) iec.UDINT {
	return iec.UDINT(convert.RoundAndClampLREAL(iec.LREAL(in), 0, iec.MAXUDINT))
}

// REAL_TO_WORD conversion
func REAL_TO_WORD(in iec.REAL) iec.WORD {
	val, _ := convert.AnyToULINT(in)
	return iec.WORD(convert.ClampULINT(val, iec.MAXUINT))
}

// REAL_TO_STRING conversion
func REAL_TO_STRING(in iec.REAL) iec.STRING {
	return iec.STRING(strconv.FormatFloat(float64(in), 'g', -1, 32))
}

// REAL_TO_LWORD conversion
func REAL_TO_LWORD(in iec.REAL) iec.LWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.LWORD(val)
}

// REAL_TO_UINT conversion
func REAL_TO_UINT(in iec.REAL) iec.UINT {
	return iec.UINT(convert.RoundAndClampLREAL(iec.LREAL(in), 0, iec.MAXUINT))
}

// REAL_TO_LREAL conversion
func REAL_TO_LREAL(in iec.REAL) iec.LREAL { return iec.LREAL(in) }

// REAL_TO_BYTE conversion
func REAL_TO_BYTE(in iec.REAL) iec.BYTE {
	val, _ := convert.AnyToULINT(in)
	return iec.BYTE(convert.ClampULINT(val, iec.MAXUSINT))
}

// REAL_TO_USINT conversion
func REAL_TO_USINT(in iec.REAL) iec.USINT {
	return iec.USINT(convert.RoundAndClampLREAL(iec.LREAL(in), 0, iec.MAXUSINT))
}

// REAL_TO_ULINT conversion
func REAL_TO_ULINT(in iec.REAL) iec.ULINT {
	return iec.ULINT(convert.RoundAndClampLREAL(iec.LREAL(in), 0, iec.MAXULINT))
}

// REAL_TO_BOOL conversion
func REAL_TO_BOOL(in iec.REAL) iec.BOOL { return in > 0 }

// REAL_TO_TIME conversion
func REAL_TO_TIME(in iec.REAL) iec.TIME { out, _ := convert.SubTime(in); return out }

// REAL_TO_INT conversion
func REAL_TO_INT(in iec.REAL) iec.INT {
	return iec.INT(convert.RoundAndClampLREAL(iec.LREAL(in), iec.MININT, iec.MAXINT))
}

/*
LREAL_TO * Conversion Section
*/
// LREAL_TO_REAL conversion
func LREAL_TO_REAL(in iec.LREAL) iec.REAL {
	return iec.REAL(convert.ClampLREAL(in, -iec.MAXREAL, iec.MAXREAL))
}

// LREAL_TO_SINT conversion
func LREAL_TO_SINT(in iec.LREAL) iec.SINT {
	return iec.SINT(convert.RoundAndClampLREAL(in, iec.MINSINT, iec.MAXSINT))
}

// LREAL_TO_LINT conversion
func LREAL_TO_LINT(in iec.LREAL) iec.LINT {
	return iec.LINT(convert.RoundAndClampLREAL(in, iec.MINLINT, iec.MAXLINT))
}

// LREAL_TO_DINT conversion
func LREAL_TO_DINT(in iec.LREAL) iec.DINT {
	return iec.DINT(convert.RoundAndClampLREAL(in, iec.MINDINT, iec.MAXDINT))
}

// LREAL_TO_DATE conversion
func LREAL_TO_DATE(in iec.LREAL) iec.DATE { out, _ := convert.SubDate(in); return out }

// LREAL_TO_DWORD conversion
func LREAL_TO_DWORD(in iec.LREAL) iec.DWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.DWORD(convert.ClampULINT(val, iec.MAXUDINT))
}

// LREAL_TO_DT conversion
func LREAL_TO_DT(in iec.LREAL) iec.DT { out, _ := convert.SubDt(in); return out }

// LREAL_TO_TOD conversion
func LREAL_TO_TOD(in iec.LREAL) iec.TOD { out, _ := convert.SubTod(in); return out }

// LREAL_TO_UDINT conversion
func LREAL_TO_UDINT(in iec.LREAL) iec.UDINT {
	return iec.UDINT(convert.RoundAndClampLREAL(in, 0, iec.MAXUDINT))
}

// LREAL_TO_WORD conversion
func LREAL_TO_WORD(in iec.LREAL) iec.WORD {
	val, _ := convert.AnyToULINT(in)
	return iec.WORD(convert.ClampULINT(val, iec.MAXUINT))
}

// LREAL_TO_STRING conversion
func LREAL_TO_STRING(in iec.LREAL) iec.STRING {
	return iec.STRING(strconv.FormatFloat(float64(in), 'g', -1, 64))
}

// LREAL_TO_LWORD conversion
func LREAL_TO_LWORD(in iec.LREAL) iec.LWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.LWORD(val)
}

// LREAL_TO_UINT conversion
func LREAL_TO_UINT(in iec.LREAL) iec.UINT {
	return iec.UINT(convert.RoundAndClampLREAL(in, 0, iec.MAXUINT))
}

// LREAL_TO_BYTE conversion
func LREAL_TO_BYTE(in iec.LREAL) iec.BYTE {
	val, _ := convert.AnyToULINT(in)
	return iec.BYTE(convert.ClampULINT(val, iec.MAXUSINT))
}

// LREAL_TO_USINT conversion
func LREAL_TO_USINT(in iec.LREAL) iec.USINT {
	return iec.USINT(convert.RoundAndClampLREAL(in, 0, iec.MAXUSINT))
}

// LREAL_TO_ULINT conversion
func LREAL_TO_ULINT(in iec.LREAL) iec.ULINT {
	return iec.ULINT(convert.RoundAndClampLREAL(in, 0, iec.MAXULINT))
}

// LREAL_TO_BOOL conversion
func LREAL_TO_BOOL(in iec.LREAL) iec.BOOL { return in > 0 }

// LREAL_TO_TIME conversion
func LREAL_TO_TIME(in iec.LREAL) iec.TIME { out, _ := convert.SubTime(in); return out }

// LREAL_TO_INT conversion
func LREAL_TO_INT(in iec.LREAL) iec.INT {
	return iec.INT(convert.RoundAndClampLREAL(in, iec.MININT, iec.MAXINT))
}

/*
SINT_TO * Conversion section
*/

// SINT_TO_REAL conversion
func SINT_TO_REAL(in iec.SINT) iec.REAL { return iec.REAL(in) }

// SINT_TO_LINT conversion
func SINT_TO_LINT(in iec.SINT) iec.LINT { return iec.LINT(in) }

// SINT_TO_DINT conversion
func SINT_TO_DINT(in iec.SINT) iec.DINT {
	return iec.DINT(convert.ClampLINT(iec.LINT(in), iec.MINDINT, iec.MAXDINT))
}

// SINT_TO_DATE conversion
func SINT_TO_DATE(in iec.SINT) iec.DATE { out, _ := convert.SubDate(in); return out }

// SINT_TO_DWORD conversion
func SINT_TO_DWORD(in iec.SINT) iec.DWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.DWORD(val)
}

// SINT_TO_DT conversion
func SINT_TO_DT(in iec.SINT) iec.DT { out, _ := convert.SubDt(in); return out }

// SINT_TO_TOD conversion
func SINT_TO_TOD(in iec.SINT) iec.TOD { out, _ := convert.SubTod(in); return out }

// SINT_TO_UDINT conversion
func SINT_TO_UDINT(in iec.SINT) iec.UDINT {
	return iec.UDINT(convert.ClampLINT(iec.LINT(in), 0, iec.MAXUDINT))
}

// SINT_TO_WORD conversion
func SINT_TO_WORD(in iec.SINT) iec.WORD {
	val, _ := convert.AnyToULINT(in)
	return iec.WORD(val)
}

// SINT_TO_STRING conversion
func SINT_TO_STRING(in iec.SINT) iec.STRING { return iec.STRING(strconv.FormatInt(int64(in), 10)) }

// SINT_TO_LWORD conversion
func SINT_TO_LWORD(in iec.SINT) iec.LWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.LWORD(val)
}

// SINT_TO_UINT conversion
func SINT_TO_UINT(in iec.SINT) iec.UINT {
	return iec.UINT(convert.ClampLINT(iec.LINT(in), 0, iec.MAXUINT))
}

// SINT_TO_LREAL conversion
func SINT_TO_LREAL(in iec.SINT) iec.LREAL { return iec.LREAL(in) }

// SINT_TO_BYTE conversion
func SINT_TO_BYTE(in iec.SINT) iec.BYTE {
	val, _ := convert.AnyToULINT(in)
	return iec.BYTE(val)
}

// SINT_TO_USINT conversion
func SINT_TO_USINT(in iec.SINT) iec.USINT {
	return iec.USINT(convert.ClampLINT(iec.LINT(in), 0, iec.MAXUSINT))
}

// SINT_TO_ULINT conversion
func SINT_TO_ULINT(in iec.SINT) iec.ULINT {
	return iec.ULINT(convert.ClampLINT(iec.LINT(in), 0, ULINT_TO_LINT(iec.MAXULINT)))
}

// SINT_TO_BOOL conversion
func SINT_TO_BOOL(in iec.SINT) iec.BOOL { return iec.BOOL(in != 0) }

// SINT_TO_TIME conversion
func SINT_TO_TIME(in iec.SINT) iec.TIME {
	val, _ := convert.AnyToLINT(in)
	return iec.TIME(val * iec.LINT(time.Millisecond))
}

// SINT_TO_INT conversion
func SINT_TO_INT(in iec.SINT) iec.INT { return iec.INT(in) }

/*
INT_TO * Conversion section
*/

// INT_TO_REAL conversion
func INT_TO_REAL(in iec.INT) iec.REAL { return iec.REAL(in) }

// INT_TO_SINT conversion
func INT_TO_SINT(in iec.INT) iec.SINT {
	return iec.SINT(convert.ClampLINT(iec.LINT(in), iec.MINSINT, iec.MAXSINT))
}

// INT_TO_LINT conversion
func INT_TO_LINT(in iec.INT) iec.LINT { return iec.LINT(in) }

// INT_TO_DINT conversion
func INT_TO_DINT(in iec.INT) iec.DINT {
	return iec.DINT(convert.ClampLINT(iec.LINT(in), iec.MINDINT, iec.MAXDINT))
}

// INT_TO_DATE conversion seconds to DATE
func INT_TO_DATE(in iec.INT) iec.DATE { out, _ := convert.SubDate(in); return out }

// INT_TO_DWORD conversion
func INT_TO_DWORD(in iec.INT) iec.DWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.DWORD(val)
}

// INT_TO_DT conversion
func INT_TO_DT(in iec.INT) iec.DT { out, _ := convert.SubDt(in); return out }

// INT_TO_TOD conversion
func INT_TO_TOD(in iec.INT) iec.TOD { out, _ := convert.SubTod(in); return out }

// INT_TO_UDINT conversion
func INT_TO_UDINT(in iec.INT) iec.UDINT {
	return iec.UDINT(convert.ClampLINT(iec.LINT(in), 0, iec.MAXUDINT))
}

// INT_TO_WORD conversion
func INT_TO_WORD(in iec.INT) iec.WORD {
	val, _ := convert.AnyToULINT(in)
	return iec.WORD(convert.ClampULINT(val, iec.MAXUINT))
}

// INT_TO_STRING conversion
func INT_TO_STRING(in iec.INT) iec.STRING { return iec.STRING(strconv.FormatInt(int64(in), 10)) }

// INT_TO_LWORD conversion
func INT_TO_LWORD(in iec.INT) iec.LWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.LWORD(val)
}

// INT_TO_UINT conversion
func INT_TO_UINT(in iec.INT) iec.UINT {
	return iec.UINT(convert.ClampLINT(iec.LINT(in), 0, iec.MAXUINT))
}

// INT_TO_LREAL conversion
func INT_TO_LREAL(in iec.INT) iec.LREAL { return iec.LREAL(in) }

// INT_TO_BYTE conversion
func INT_TO_BYTE(in iec.INT) iec.BYTE {
	val, _ := convert.AnyToULINT(in)
	return iec.BYTE(convert.ClampULINT(val, iec.MAXUSINT))
}

// INT_TO_USINT conversion
func INT_TO_USINT(in iec.INT) iec.USINT {
	return iec.USINT(convert.ClampLINT(iec.LINT(in), 0, iec.MAXUSINT))
}

// INT_TO_ULINT conversion
func INT_TO_ULINT(in iec.INT) iec.ULINT {
	return iec.ULINT(convert.ClampLINT(iec.LINT(in), 0, ULINT_TO_LINT(iec.MAXULINT)))
}

// INT_TO_BOOL conversion
func INT_TO_BOOL(in iec.INT) iec.BOOL { return in > 0 }

// INT_TO_TIME conversion
func INT_TO_TIME(in iec.INT) iec.TIME { out, _ := convert.SubTime(in); return out }

/*
LINT_TO * Conversion section
*/

// LINT_TO_REAL conversion
func LINT_TO_REAL(in iec.LINT) iec.REAL { return iec.REAL(in) }

// LINT_TO_SINT conversion
func LINT_TO_SINT(in iec.LINT) iec.SINT {
	return iec.SINT(convert.ClampLINT(in, iec.MINSINT, iec.MAXSINT))
}

// LINT_TO_DINT conversion
func LINT_TO_DINT(in iec.LINT) iec.DINT {
	return iec.DINT(convert.ClampLINT(in, iec.MINDINT, iec.MAXDINT))
}

// LINT_TO_DWORD conversion
func LINT_TO_DWORD(in iec.LINT) iec.DWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.DWORD(convert.ClampULINT(val, iec.MAXUDINT))
}

// LINT_TO_DATE conversion
func LINT_TO_DATE(in iec.LINT) iec.DATE { out, _ := convert.SubDate(in); return out }

// LINT_TO_TOD conversion
func LINT_TO_TOD(in iec.LINT) iec.TOD { out, _ := convert.SubTod(in); return out }

// LINT_TO_UDINT conversion
func LINT_TO_UDINT(in iec.LINT) iec.UDINT { return iec.UDINT(convert.ClampLINT(in, 0, iec.MAXUDINT)) }

// LINT_TO_UINT conversion
func LINT_TO_UINT(in iec.LINT) iec.UINT { return iec.UINT(convert.ClampLINT(in, 0, iec.MAXUINT)) }

// LINT_TO_WORD conversion
func LINT_TO_WORD(in iec.LINT) iec.WORD {
	val, _ := convert.AnyToULINT(in)
	return iec.WORD(convert.ClampULINT(val, iec.MAXUINT))
}

// LINT_TO_STRING conversion
func LINT_TO_STRING(in iec.LINT) iec.STRING { return iec.STRING(strconv.FormatInt(int64(in), 10)) }

// LINT_TO_LWORD conversion
func LINT_TO_LWORD(in iec.LINT) iec.LWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.LWORD(val)
}

// LINT_TO_ULINT conversion
func LINT_TO_ULINT(in iec.LINT) iec.ULINT {
	// If the input is negative, it should be Clamped to 0 for unsigned conversion.
	if in < 0 {
		return 0
	}
	return iec.ULINT(in)
}

// LINT_TO_LREAL conversion
func LINT_TO_LREAL(in iec.LINT) iec.LREAL { return iec.LREAL(in) }

// LINT_TO_BYTE conversion
func LINT_TO_BYTE(in iec.LINT) iec.BYTE {
	val, _ := convert.AnyToULINT(in)
	return iec.BYTE(convert.ClampULINT(val, iec.MAXUSINT))
}

// LINT_TO_USINT conversion
func LINT_TO_USINT(in iec.LINT) iec.USINT { return iec.USINT(convert.ClampLINT(in, 0, iec.MAXUSINT)) }

// LINT_TO_BOOL conversion
func LINT_TO_BOOL(in iec.LINT) iec.BOOL { return in > 0 }

// LINT_TO_TIME conversion
func LINT_TO_TIME(in iec.LINT) iec.TIME { out, _ := convert.SubTime(in); return out }

// LINT_TO_INT conversion
func LINT_TO_INT(in iec.LINT) iec.INT { return iec.INT(convert.ClampLINT(in, iec.MININT, iec.MAXINT)) }

/*
// DINT_TO Conversions
*/

// DINT_TO_REAL conversion
func DINT_TO_REAL(in iec.DINT) iec.REAL { return iec.REAL(in) }

// DINT_TO_SINT conversion
func DINT_TO_SINT(in iec.DINT) iec.SINT {
	return iec.SINT(convert.ClampLINT(iec.LINT(in), iec.MINSINT, iec.MAXSINT))
}

// DINT_TO_LINT conversion
func DINT_TO_LINT(in iec.DINT) iec.LINT { return iec.LINT(in) }

// DINT_TO_DATE conversion
func DINT_TO_DATE(in iec.DINT) iec.DATE { out, _ := convert.SubDate(in); return out }

// DINT_TO_DWORD conversion
func DINT_TO_DWORD(in iec.DINT) iec.DWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.DWORD(convert.ClampULINT(val, iec.MAXUDINT))
}

// DINT_TO_DT conversion
func DINT_TO_DT(in iec.DINT) iec.DT { out, _ := convert.SubDt(in); return out }

// DINT_TO_TOD conversion
func DINT_TO_TOD(in iec.DINT) iec.TOD { out, _ := convert.SubTod(in); return out }

// DINT_TO_UDINT conversion
func DINT_TO_UDINT(in iec.DINT) iec.UDINT {
	return iec.UDINT(convert.ClampLINT(iec.LINT(in), 0, iec.MAXUDINT))
}

// DINT_TO_WORD conversion
func DINT_TO_WORD(in iec.DINT) iec.WORD {
	val, _ := convert.AnyToULINT(in)
	return iec.WORD(convert.ClampULINT(val, iec.MAXUINT))
}

// DINT_TO_STRING conversion
func DINT_TO_STRING(in iec.DINT) iec.STRING { return iec.STRING(strconv.FormatInt(int64(in), 10)) }

// DINT_TO_LWORD conversion
func DINT_TO_LWORD(in iec.DINT) iec.LWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.LWORD(val)
}

// DINT_TO_UINT conversion
func DINT_TO_UINT(in iec.DINT) iec.UINT {
	return iec.UINT(convert.ClampLINT(iec.LINT(in), 0, iec.MAXUINT))
}

// DINT_TO_LREAL conversion
func DINT_TO_LREAL(in iec.DINT) iec.LREAL { return iec.LREAL(in) }

// DINT_TO_BYTE conversion
func DINT_TO_BYTE(in iec.DINT) iec.BYTE {
	val, _ := convert.AnyToULINT(in)
	return iec.BYTE(convert.ClampULINT(val, iec.MAXUSINT))
}

// DINT_TO_USINT conversion
func DINT_TO_USINT(in iec.DINT) iec.USINT {
	return iec.USINT(convert.ClampLINT(iec.LINT(in), 0, iec.MAXUSINT))
}

// DINT_TO_ULINT conversion
func DINT_TO_ULINT(in iec.DINT) iec.ULINT {
	return iec.ULINT(convert.ClampLINT(iec.LINT(in), 0, ULINT_TO_LINT(iec.MAXULINT)))
}

// DINT_TO_BOOL conversion
func DINT_TO_BOOL(in iec.DINT) iec.BOOL { return in > 0 }

// DINT_TO_TIME conversion
func DINT_TO_TIME(in iec.DINT) iec.TIME { out, _ := convert.SubTime(in); return out }

// DINT_TO_INT conversion
func DINT_TO_INT(in iec.DINT) iec.INT {
	return iec.INT(convert.ClampLINT(iec.LINT(in), iec.MININT, iec.MAXINT))
}

/*
USINT_TO * Conversion section
*/
// USINT_TO_REAL conversion
func USINT_TO_REAL(in iec.USINT) iec.REAL { return iec.REAL(in) }

// USINT_TO_SINT conversion
func USINT_TO_SINT(in iec.USINT) iec.SINT {
	return iec.SINT(convert.ClampULINT(iec.ULINT(in), iec.MAXSINT))
}

// USINT_TO_LINT conversion
func USINT_TO_LINT(in iec.USINT) iec.LINT {
	return iec.LINT(convert.ClampULINT(iec.ULINT(in), iec.MAXLINT))
}

// USINT_TO_DINT conversion
func USINT_TO_DINT(in iec.USINT) iec.DINT {
	return iec.DINT(convert.ClampULINT(iec.ULINT(in), iec.MAXDINT))
}

// USINT_TO_DATE conversion
func USINT_TO_DATE(in iec.USINT) iec.DATE { val, _ := convert.SubDate(iec.LINT(in)); return val }

// USINT_TO_DWORD conversion
func USINT_TO_DWORD(in iec.USINT) iec.DWORD { return iec.DWORD(in) }

// USINT_TO_DT conversion
func USINT_TO_DT(in iec.USINT) iec.DT { val, _ := convert.SubDt(iec.LINT(in)); return val }

// USINT_TO_TOD conversion
func USINT_TO_TOD(in iec.USINT) iec.TOD { val, _ := convert.SubTod(iec.LINT(in)); return val }

// USINT_TO_UDINT conversion
func USINT_TO_UDINT(in iec.USINT) iec.UDINT {
	return iec.UDINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUDINT))
}

// USINT_TO_WORD conversion
func USINT_TO_WORD(in iec.USINT) iec.WORD { return iec.WORD(in) }

// USINT_TO_STRING conversion
func USINT_TO_STRING(in iec.USINT) iec.STRING { return iec.STRING(strconv.FormatUint(uint64(in), 10)) }

// USINT_TO_LWORD conversion
func USINT_TO_LWORD(in iec.USINT) iec.LWORD { return iec.LWORD(in) }

// USINT_TO_UINT conversion
func USINT_TO_UINT(in iec.USINT) iec.UINT {
	return iec.UINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUINT))
}

// USINT_TO_LREAL conversion
func USINT_TO_LREAL(in iec.USINT) iec.LREAL { return iec.LREAL(in) }

// USINT_TO_BYTE conversion
func USINT_TO_BYTE(in iec.USINT) iec.BYTE { return iec.BYTE(in) }

// USINT_TO_ULINT conversion
func USINT_TO_ULINT(in iec.USINT) iec.ULINT {
	return iec.ULINT(convert.ClampULINT(iec.ULINT(in), iec.MAXULINT))
}

// USINT_TO_BOOL conversion
func USINT_TO_BOOL(in iec.USINT) iec.BOOL { return in > 0 }

// USINT_TO_TIME conversion
func USINT_TO_TIME(in iec.USINT) iec.TIME { val, _ := convert.SubTime(iec.LINT(in)); return val }

// USINT_TO_INT conversion
func USINT_TO_INT(in iec.USINT) iec.INT {
	return iec.INT(convert.ClampULINT(iec.ULINT(in), iec.MAXINT))
}

/*
UINT_TO conversion
*/
// UINT_TO_REAL conversion
func UINT_TO_REAL(in iec.UINT) iec.REAL { return iec.REAL(in) }

// UINT_TO_SINT conversion
func UINT_TO_SINT(in iec.UINT) iec.SINT {
	return iec.SINT(convert.ClampULINT(iec.ULINT(in), iec.MAXSINT))
}

// UINT_TO_LINT conversion
func UINT_TO_LINT(in iec.UINT) iec.LINT {
	return iec.LINT(convert.ClampULINT(iec.ULINT(in), iec.MAXLINT))
}

// UINT_TO_DINT conversion
func UINT_TO_DINT(in iec.UINT) iec.DINT {
	return iec.DINT(convert.ClampULINT(iec.ULINT(in), iec.MAXDINT))
}

// UINT_TO_DATE conversion
func UINT_TO_DATE(in iec.UINT) iec.DATE { val, _ := convert.SubDate(iec.LINT(in)); return val }

// UINT_TO_DWORD conversion
func UINT_TO_DWORD(in iec.UINT) iec.DWORD { return iec.DWORD(in) }

// UINT_TO_DT conversion
func UINT_TO_DT(in iec.UINT) iec.DT { val, _ := convert.SubDt(iec.LINT(in)); return val }

// UINT_TO_TOD conversion
func UINT_TO_TOD(in iec.UINT) iec.TOD { val, _ := convert.SubTod(iec.LINT(in)); return val }

// UINT_TO_UDINT conversion
func UINT_TO_UDINT(in iec.UINT) iec.UDINT {
	return iec.UDINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUDINT))
}

// UINT_TO_WORD conversion
func UINT_TO_WORD(in iec.UINT) iec.WORD { return iec.WORD(in) }

// UINT_TO_STRING conversion
func UINT_TO_STRING(in iec.UINT) iec.STRING { return iec.STRING(strconv.FormatUint(uint64(in), 10)) }

// UINT_TO_LWORD conversion
func UINT_TO_LWORD(in iec.UINT) iec.LWORD { return iec.LWORD(in) }

// UINT_TO_LREAL conversion
func UINT_TO_LREAL(in iec.UINT) iec.LREAL { return iec.LREAL(in) }

// UINT_TO_BYTE conversion
func UINT_TO_BYTE(in iec.UINT) iec.BYTE {
	return iec.BYTE(convert.ClampULINT(iec.ULINT(in), iec.MAXUSINT))
}

// UINT_TO_USINT conversion
func UINT_TO_USINT(in iec.UINT) iec.USINT {
	return iec.USINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUSINT))
}

// UINT_TO_ULINT conversion
func UINT_TO_ULINT(in iec.UINT) iec.ULINT {
	return iec.ULINT(convert.ClampULINT(iec.ULINT(in), iec.MAXULINT))
}

// UINT_TO_BOOL conversion
func UINT_TO_BOOL(in iec.UINT) iec.BOOL { return in > 0 }

// UINT_TO_TIME conversion
func UINT_TO_TIME(in iec.UINT) iec.TIME { val, _ := convert.SubTime(iec.LINT(in)); return val }

// UINT_TO_INT conversion
func UINT_TO_INT(in iec.UINT) iec.INT { return iec.INT(convert.ClampULINT(iec.ULINT(in), iec.MAXINT)) }

/*
UDINT_TO * Conversion section
*/
// UDINT_TO_REAL conversion
func UDINT_TO_REAL(in iec.UDINT) iec.REAL { return iec.REAL(in) }

// UDINT_TO_SINT conversion
func UDINT_TO_SINT(in iec.UDINT) iec.SINT {
	return iec.SINT(convert.ClampULINT(iec.ULINT(in), iec.MAXSINT))
}

// UDINT_TO_LINT conversion
func UDINT_TO_LINT(in iec.UDINT) iec.LINT {
	return iec.LINT(convert.ClampULINT(iec.ULINT(in), iec.MAXLINT))
}

// UDINT_TO_DINT conversion
func UDINT_TO_DINT(in iec.UDINT) iec.DINT {
	return iec.DINT(convert.ClampULINT(iec.ULINT(in), iec.MAXDINT))
}

// UDINT_TO_DATE conversion
func UDINT_TO_DATE(in iec.UDINT) iec.DATE { val, _ := convert.SubDate(iec.LINT(in)); return val }

// UDINT_TO_DWORD conversion
func UDINT_TO_DWORD(in iec.UDINT) iec.DWORD {
	return iec.DWORD(in)
}

// UDINT_TO_DT conversion
func UDINT_TO_DT(in iec.UDINT) iec.DT { val, _ := convert.SubDt(iec.LINT(in)); return val }

// UDINT_TO_TOD conversion
func UDINT_TO_TOD(in iec.UDINT) iec.TOD { val, _ := convert.SubTod(iec.LINT(in)); return val }

// UDINT_TO_WORD conversion
func UDINT_TO_WORD(in iec.UDINT) iec.WORD {
	val, _ := convert.AnyToULINT(in)
	return iec.WORD(convert.ClampULINT(val, iec.MAXUINT))
}

// UDINT_TO_STRING conversion
func UDINT_TO_STRING(in iec.UDINT) iec.STRING { return iec.STRING(strconv.FormatUint(uint64(in), 10)) }

// UDINT_TO_LWORD conversion
func UDINT_TO_LWORD(in iec.UDINT) iec.LWORD { return iec.LWORD(in) }

// UDINT_TO_UINT conversion
func UDINT_TO_UINT(in iec.UDINT) iec.UINT {
	return iec.UINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUINT))
}

// UDINT_TO_LREAL conversion
func UDINT_TO_LREAL(in iec.UDINT) iec.LREAL { return iec.LREAL(in) }

// UDINT_TO_BYTE conversion
func UDINT_TO_BYTE(in iec.UDINT) iec.BYTE {
	return iec.BYTE(convert.ClampULINT(iec.ULINT(in), iec.MAXUSINT))
}

// UDINT_TO_USINT conversion
func UDINT_TO_USINT(in iec.UDINT) iec.USINT {
	return iec.USINT(convert.ClampULINT(iec.ULINT(in), iec.MAXUSINT))
}

// UDINT_TO_ULINT conversion
func UDINT_TO_ULINT(in iec.UDINT) iec.ULINT {
	return iec.ULINT(convert.ClampULINT(iec.ULINT(in), iec.MAXULINT))
}

// UDINT_TO_BOOL conversion
func UDINT_TO_BOOL(in iec.UDINT) iec.BOOL { return in > 0 }

// UDINT_TO_TIME conversion
func UDINT_TO_TIME(in iec.UDINT) iec.TIME { val, _ := convert.SubTime(iec.LINT(in)); return val }

// UDINT_TO_INT conversion
func UDINT_TO_INT(in iec.UDINT) iec.INT {
	return iec.INT(convert.ClampULINT(iec.ULINT(in), iec.MAXINT))
}

/*
ULINT_TO * Conversion section
*/
// ULINT_TO_REAL conversion
func ULINT_TO_REAL(in iec.ULINT) iec.REAL { return iec.REAL(in) }

// ULINT_TO_SINT conversion
func ULINT_TO_SINT(in iec.ULINT) iec.SINT { return iec.SINT(convert.ClampULINT(in, iec.MAXSINT)) }

// ULINT_TO_LINT conversion
func ULINT_TO_LINT(in iec.ULINT) iec.LINT { return iec.LINT(convert.ClampULINT(in, iec.MAXLINT)) }

// ULINT_TO_DINT conversion
func ULINT_TO_DINT(in iec.ULINT) iec.DINT { return iec.DINT(convert.ClampULINT(in, iec.MAXDINT)) }

// ULINT_TO_DATE conversion
func ULINT_TO_DATE(in iec.ULINT) iec.DATE { val, _ := convert.SubDate(iec.LINT(in)); return val }

// ULINT_TO_DWORD conversion
func ULINT_TO_DWORD(in iec.ULINT) iec.DWORD { return iec.DWORD(convert.ClampULINT(in, iec.MAXUDINT)) }

// ULINT_TO_DT conversion
func ULINT_TO_DT(in iec.ULINT) iec.DT { val, _ := convert.SubDt(iec.LINT(in)); return val }

// ULINT_TO_TOD conversion
func ULINT_TO_TOD(in iec.ULINT) iec.TOD { val, _ := convert.SubTod(iec.LINT(in)); return val }

// ULINT_TO_UDINT conversion
func ULINT_TO_UDINT(in iec.ULINT) iec.UDINT { return iec.UDINT(in) }

// ULINT_TO_WORD conversion
func ULINT_TO_WORD(in iec.ULINT) iec.WORD { return iec.WORD(convert.ClampULINT(in, iec.MAXUINT)) }

// ULINT_TO_STRING conversion
func ULINT_TO_STRING(in iec.ULINT) iec.STRING { return iec.STRING(strconv.FormatUint(uint64(in), 10)) }

// ULINT_TO_LWORD conversion
func ULINT_TO_LWORD(in iec.ULINT) iec.LWORD { return iec.LWORD(in) }

// ULINT_TO_UINT conversion
func ULINT_TO_UINT(in iec.ULINT) iec.UINT { return iec.UINT(convert.ClampULINT(in, iec.MAXUINT)) }

// ULINT_TO_LREAL conversion
func ULINT_TO_LREAL(in iec.ULINT) iec.LREAL { return iec.LREAL(in) }

// ULINT_TO_BYTE conversion
func ULINT_TO_BYTE(in iec.ULINT) iec.BYTE { return iec.BYTE(convert.ClampULINT(in, iec.MAXUSINT)) }

// ULINT_TO_USINT conversion
func ULINT_TO_USINT(in iec.ULINT) iec.USINT { return iec.USINT(convert.ClampULINT(in, iec.MAXUSINT)) }

// ULINT_TO_BOOL conversion
func ULINT_TO_BOOL(in iec.ULINT) iec.BOOL { return in > 0 }

// ULINT_TO_TIME conversion
func ULINT_TO_TIME(in iec.ULINT) iec.TIME { val, _ := convert.SubTime(iec.LINT(in)); return val }

// ULINT_TO_INT conversion
func ULINT_TO_INT(in iec.ULINT) iec.INT { return iec.INT(convert.ClampULINT(in, iec.MAXINT)) }

/*
DATE_TO * Conversion section
*/

// DATE_TO_REAL conversion
func DATE_TO_REAL(in iec.DATE) iec.REAL { return iec.REAL(time.Time(in).UnixMilli()) }

// DATE_TO_SINT conversion
func DATE_TO_SINT(in iec.DATE) iec.SINT {
	return iec.SINT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), iec.MINSINT, iec.MAXSINT))
}

// DATE_TO_LINT conversion
func DATE_TO_LINT(in iec.DATE) iec.LINT { return iec.LINT(time.Time(in).UnixMilli()) }

// DATE_TO_DINT conversion
func DATE_TO_DINT(in iec.DATE) iec.DINT {
	return iec.DINT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), iec.MINDINT, iec.MAXDINT))
}

// DATE_TO_BYTE conversion
func DATE_TO_BYTE(in iec.DATE) iec.BYTE {
	val, _ := convert.AnyToULINT(in)
	return iec.BYTE(convert.ClampULINT(val, iec.MAXUSINT))
}

// DATE_TO_WORD conversion
func DATE_TO_WORD(in iec.DATE) iec.WORD {
	val, _ := convert.AnyToULINT(in)
	return iec.WORD(convert.ClampULINT(val, iec.MAXUINT))
}

// DATE_TO_DWORD conversion
func DATE_TO_DWORD(in iec.DATE) iec.DWORD {
	return iec.DWORD(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), 0, iec.MAXUDINT))
}

// DATE_TO_LWORD conversion
func DATE_TO_LWORD(in iec.DATE) iec.LWORD {
	return iec.LWORD(convert.ClampLINT(in.CONVERT(), 0, -1)) // -1 for max ULINT
}

// DATE_TO_UDINT conversion
func DATE_TO_UDINT(in iec.DATE) iec.UDINT {
	return iec.UDINT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), 0, iec.MAXUDINT))
}

// DATE_TO_STRING conversion
func DATE_TO_STRING(in iec.DATE) iec.STRING { return iec.STRING(in.String()) }

// DATE_TO_UINT conversion
func DATE_TO_UINT(in iec.DATE) iec.UINT {
	return iec.UINT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), 0, iec.MAXUINT))
}

// DATE_TO_LREAL conversion
func DATE_TO_LREAL(in iec.DATE) iec.LREAL { return iec.LREAL(time.Time(in).UnixMilli()) }

// DATE_TO_USINT conversion
func DATE_TO_USINT(in iec.DATE) iec.USINT {
	return iec.USINT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), 0, iec.MAXUSINT))
}

// DATE_TO_ULINT conversion
func DATE_TO_ULINT(in iec.DATE) iec.ULINT {
	val := iec.LINT(time.Time(in).UnixMilli())
	if val < 0 {
		return 0
	}
	return iec.ULINT(val)
}

// DATE_TO_INT conversion
func DATE_TO_INT(in iec.DATE) iec.INT {
	return iec.INT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), iec.MININT, iec.MAXINT))
}

// DATE_TO_TIME
func DATE_TO_TIME(in iec.DATE) iec.TIME {
	// Converts the DATE (a point in time) to a TIME (duration)
	// representing the milliseconds elapsed since the Unix epoch.
	return iec.TIME(time.Time(in).UnixMilli() * int64(time.Millisecond))
}

/*
DT_TO conversion
*/

// DT_TO_REAL conversion
func DT_TO_REAL(in iec.DT) iec.REAL { return iec.REAL(time.Time(in).UnixMilli()) }

// DT_TO_SINT conversion
func DT_TO_SINT(in iec.DT) iec.SINT {
	return iec.SINT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), iec.MINSINT, iec.MAXSINT))
}

// DT_TO_LINT conversion
func DT_TO_LINT(in iec.DT) iec.LINT { return iec.LINT(time.Time(in).UnixMilli()) }

// DT_TO_DINT conversion
func DT_TO_DINT(in iec.DT) iec.DINT {
	return iec.DINT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), iec.MINDINT, iec.MAXDINT))
}

// DT_TO_DWORD conversion
func DT_TO_DWORD(in iec.DT) iec.DWORD {
	return iec.DWORD(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), 0, iec.MAXUDINT))
}

// DT_TO_USINT conversion
func DT_TO_USINT(in iec.DT) iec.USINT {
	return iec.USINT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), 0, iec.MAXUSINT))
}

// DT_TO_UDINT conversion
func DT_TO_UDINT(in iec.DT) iec.UDINT {
	return iec.UDINT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), 0, iec.MAXUDINT))
}

// DT_TO_WORD conversion
func DT_TO_WORD(in iec.DT) iec.WORD {
	return iec.WORD(convert.ClampLINT(in.CONVERT(), 0, iec.MAXUINT))
}

// DT_TO_STRING conversion
func DT_TO_STRING(in iec.DT) iec.STRING { return iec.STRING(in.String()) }

// DT_TO_LWORD conversion
func DT_TO_LWORD(in iec.DT) iec.LWORD {
	return iec.LWORD(convert.ClampLINT(in.CONVERT(), 0, -1)) // -1 for max ULINT
}

// DT_TO_UINT conversion
func DT_TO_UINT(in iec.DT) iec.UINT {
	return iec.UINT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), 0, iec.MAXUINT))
}

// DT_TO_LREAL conversion
func DT_TO_LREAL(in iec.DT) iec.LREAL { return iec.LREAL(time.Time(in).UnixMilli()) }

// DT_TO_BYTE conversion
func DT_TO_BYTE(in iec.DT) iec.BYTE {
	val, _ := convert.AnyToULINT(in)
	return iec.BYTE(convert.ClampULINT(val, iec.MAXUSINT))
}

// DT_TO_ULINT conversion
func DT_TO_ULINT(in iec.DT) iec.ULINT {
	val := iec.LINT(time.Time(in).UnixMilli())
	if val < 0 {
		return 0
	}
	return iec.ULINT(val)
}

// DT_TO_INT conversion
func DT_TO_INT(in iec.DT) iec.INT {
	return iec.INT(convert.ClampLINT(iec.LINT(time.Time(in).UnixMilli()), iec.MININT, iec.MAXINT))
}

// DT_TO_DATE extracts the DATE part from a DATE_AND_TIME value.
func DT_TO_DATE(in iec.DT) iec.DATE {
	t := time.Time(in)
	// Returns a new DATE with the time part zeroed out, preserving the location.
	return iec.DATE(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()))
}

// DT_TO_TOD extracts the TIME_OF_DAY part from a DATE_AND_TIME value.
func DT_TO_TOD(in iec.DT) iec.TOD {
	// The standard implies the date part is zeroed out. A simple cast to TOD is sufficient
	// as the interpretation of a TOD value focuses only on the time part.
	return iec.TOD(in)
}

/*
TOD_TO conversion
*/

// TOD_TO_REAL conversion
func TOD_TO_REAL(in iec.TOD) iec.REAL { return iec.REAL(in.CONVERT()) }

// TOD_TO_SINT conversion
func TOD_TO_SINT(in iec.TOD) iec.SINT {
	return iec.SINT(convert.ClampLINT(in.CONVERT(), iec.MINSINT, iec.MAXSINT))
}

// TOD_TO_LINT conversion
func TOD_TO_LINT(in iec.TOD) iec.LINT { return in.CONVERT() }

// TOD_TO_DINT conversion
func TOD_TO_DINT(in iec.TOD) iec.DINT { return iec.DINT(TOD_TO_LINT(in)) }

// TOD_TO_DWORD conversion
func TOD_TO_DWORD(in iec.TOD) iec.DWORD {
	val, _ := convert.AnyToULINT(in)
	return iec.DWORD(val)
}

// TOD_TO_UDINT conversion
func TOD_TO_UDINT(in iec.TOD) iec.UDINT { return iec.UDINT(TOD_TO_LINT(in)) }

// TOD_TO_WORD conversion
func TOD_TO_WORD(in iec.TOD) iec.WORD {
	val, _ := convert.AnyToULINT(in)
	return iec.WORD(convert.ClampULINT(val, iec.MAXUINT))
}

// TOD_TO_STRING conversion
func TOD_TO_STRING(in iec.TOD) iec.STRING { return iec.STRING(in.String()) }

// TOD_TO_LWORD conversion
func TOD_TO_LWORD(in iec.TOD) iec.LWORD { out, _ := convert.SubLword(in); return out }

// TOD_TO_UINT conversion
func TOD_TO_UINT(in iec.TOD) iec.UINT { return iec.UINT(TOD_TO_LINT(in)) }

// TOD_TO_LREAL conversion
func TOD_TO_LREAL(in iec.TOD) iec.LREAL { return iec.LREAL(in.CONVERT()) }

// TOD_TO_BYTE conversion
func TOD_TO_BYTE(in iec.TOD) iec.BYTE {
	val, _ := convert.AnyToULINT(in)
	return iec.BYTE(convert.ClampULINT(val, iec.MAXUSINT))
}

// TOD_TO_USINT conversion
func TOD_TO_USINT(in iec.TOD) iec.USINT { return iec.USINT(TOD_TO_LINT(in)) }

// TOD_TO_ULINT conversion
func TOD_TO_ULINT(in iec.TOD) iec.ULINT { return iec.ULINT(TOD_TO_LINT(in)) }

// TOD_TO_INT conversion
func TOD_TO_INT(in iec.TOD) iec.INT { return iec.INT(TOD_TO_LINT(in)) }

/*
TIME_TO conversion
*/

// TIME_TO_REAL conversion
func TIME_TO_REAL(in iec.TIME) iec.REAL { return iec.REAL(time.Duration(in).Milliseconds()) }

// TIME_TO_SINT conversion
func TIME_TO_SINT(in iec.TIME) iec.SINT {
	val, _ := convert.AnyToLINT(in)
	return iec.SINT(convert.ClampLINT(val, iec.MINSINT, iec.MAXSINT))
}

// TIME_TO_LINT conversion
func TIME_TO_LINT(in iec.TIME) iec.LINT { return iec.LINT(time.Duration(in).Milliseconds()) }

// TIME_TO_DINT conversion
func TIME_TO_DINT(in iec.TIME) iec.DINT { return iec.DINT(time.Duration(in).Milliseconds()) }

// TIME_TO_DWORD conversion
func TIME_TO_DWORD(in iec.TIME) iec.DWORD {
	val, _ := convert.AnyToULINT(iec.LINT(time.Duration(in).Milliseconds()))
	return iec.DWORD(val)
}

// TIME_TO_UDINT conversion
func TIME_TO_UDINT(in iec.TIME) iec.UDINT { return iec.UDINT(time.Duration(in).Milliseconds()) }

// TIME_TO_WORD conversion
func TIME_TO_WORD(in iec.TIME) iec.WORD {
	val, _ := convert.AnyToULINT(iec.LINT(time.Duration(in).Milliseconds()))
	return iec.WORD(convert.ClampULINT(val, iec.MAXUINT))
}

// TIME_TO_STRING conversion
func TIME_TO_STRING(in iec.TIME) iec.STRING { return iec.STRING(in.String()) }

// TIME_TO_LWORD conversion
func TIME_TO_LWORD(in iec.TIME) iec.LWORD {
	val, _ := convert.AnyToULINT(iec.LINT(time.Duration(in).Milliseconds()))
	return iec.LWORD(val)
}

// TIME_TO_UINT conversion
func TIME_TO_UINT(in iec.TIME) iec.UINT { return iec.UINT(time.Duration(in).Milliseconds()) }

// TIME_TO_LREAL conversion
func TIME_TO_LREAL(in iec.TIME) iec.LREAL { return iec.LREAL(time.Duration(in).Milliseconds()) }

// TIME_TO_BYTE conversion
func TIME_TO_BYTE(in iec.TIME) iec.BYTE {
	val, _ := convert.AnyToULINT(iec.LINT(time.Duration(in).Milliseconds()))
	return iec.BYTE(convert.ClampULINT(val, iec.MAXUSINT))
}

// TIME_TO_USINT conversion
func TIME_TO_USINT(in iec.TIME) iec.USINT {
	return iec.USINT(convert.ClampLINT(iec.LINT(time.Duration(in).Milliseconds()), 0, iec.MAXUSINT))
}

// TIME_TO_ULINT conversion
func TIME_TO_ULINT(in iec.TIME) iec.ULINT { return iec.ULINT(time.Duration(in).Milliseconds()) }

// TIME_TO_INT conversion
func TIME_TO_INT(in iec.TIME) iec.INT { return iec.INT(time.Duration(in).Milliseconds()) }

/*
Math conversion of float to bits and vice versa
*/

// REAL_TO_BITS conversion: floats to uint32 as bits
func REAL_TO_BITS(in iec.REAL) iec.UDINT { return iec.UDINT(math.Float32bits(float32(in))) }

// BITS_TO_REAL converstion: uint32 bits to float
func BITS_TO_REAL(in iec.UDINT) iec.REAL { return iec.REAL(math.Float32frombits(uint32(in))) }

// LREAL_TO_BITS conversion: floats to uint32 as bits
func LREAL_TO_BITS(in iec.LREAL) iec.ULINT { return iec.ULINT(math.Float64bits(float64(in))) }

// BITS_TO_REAL converstion: uint32 bits to float
func BITS_TO_LREAL(in iec.ULINT) iec.LREAL { return iec.LREAL(math.Float64frombits(uint64(in))) }

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
		return 0, &convert.ConversionError{
			Value:    in,
			FromType: "uint64",
			ToType:   "BCD",
			Reason:   "input value too large for 64-bit BCD representation",
		}
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
			return 0, &convert.ConversionError{
				Value:    in,
				FromType: "BCD",
				ToType:   "uint64",
				Reason:   fmt.Sprintf("invalid BCD nibble %d", nibble),
			}
		}
		res += nibble * factor
		if tempVal>>((i+1)*4) == 0 {
			break // Stop if remaining bits are zero
		}
		factor *= 10
	}
	return res, nil
}

func USINT_TO_BCD_BYTE(in iec.USINT) (iec.BYTE, error) {
	out, err := uintToBCD(uint64(in))
	if err != nil || out > iec.MAXUSINT {
		return 0, &convert.ConversionError{Value: in, FromType: "USINT", ToType: "BCD_BYTE", Reason: "value out of range", Err: err}
	}
	return iec.BYTE(out), nil
}
func UINT_TO_BCD_WORD(in iec.UINT) (iec.WORD, error) {
	out, err := uintToBCD(uint64(in))
	return iec.WORD(out), err
}
func UDINT_TO_BCD_DWORD(in iec.UDINT) (iec.DWORD, error) {
	out, err := uintToBCD(uint64(in))
	return iec.DWORD(out), err
}
func ULINT_TO_BCD_LWORD(in iec.ULINT) (iec.LWORD, error) {
	out, err := uintToBCD(uint64(in))
	return iec.LWORD(out), err
}

func BYTE_BCD_TO_USINT(in iec.BYTE) (iec.USINT, error) {
	out, err := bcdToUint(uint64(in))
	if err != nil {
		return 0, err
	}
	return iec.USINT(convert.ClampULINT(iec.ULINT(out), iec.MAXUSINT)), nil
}
func WORD_BCD_TO_UINT(in iec.WORD) (iec.UINT, error) {
	out, err := bcdToUint(uint64(in))
	return iec.UINT(out), err
}
func DWORD_BCD_TO_UDINT(in iec.DWORD) (iec.UDINT, error) {
	out, err := bcdToUint(uint64(in))
	return iec.UDINT(out), err
}
func LWORD_BCD_TO_ULINT(in iec.LWORD) (iec.ULINT, error) {
	out, err := bcdToUint(uint64(in))
	return iec.ULINT(out), err
}
