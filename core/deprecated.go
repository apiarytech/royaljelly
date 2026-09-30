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

package core

import (
	"github.com/apiarytech/royaljelly/convert"
	"github.com/apiarytech/royaljelly/iec"
)

// This file keeps the type-conversion helpers that used to live in package core
// available for existing callers. They moved to package convert so that the
// standard-function packages no longer depend on the scheduler.
// New code should import github.com/apiarytech/royaljelly/convert directly.

// ConversionError is an alias of convert.ConversionError.
//
// Deprecated: use convert.ConversionError.
type ConversionError = convert.ConversionError

// GetTypeName is deprecated: use convert.GetTypeName.
func GetTypeName(v any) string { return convert.GetTypeName(v) }

// IsPlcFloat is deprecated: use convert.IsPlcFloat.
func IsPlcFloat[T any](val T) bool { return convert.IsPlcFloat(val) }

// IsPlcInt is deprecated: use convert.IsPlcInt.
func IsPlcInt[T any](val T) bool { return convert.IsPlcInt(val) }

// IsPlcTimeType is deprecated: use convert.IsPlcTimeType.
func IsPlcTimeType[T any](val T) bool { return convert.IsPlcTimeType(val) }

// AnyToREAL is deprecated: use convert.AnyToREAL.
func AnyToREAL[T any](val T) (iec.REAL, error) { return convert.AnyToREAL(val) }

// AnyToLREAL is deprecated: use convert.AnyToLREAL.
func AnyToLREAL[T any](val T) (iec.LREAL, error) { return convert.AnyToLREAL(val) }

// AnyToLINT is deprecated: use convert.AnyToLINT.
func AnyToLINT[T any](val T) (iec.LINT, error) { return convert.AnyToLINT(val) }

// AnyToULINT is deprecated: use convert.AnyToULINT.
func AnyToULINT[T any](val T) (iec.ULINT, error) { return convert.AnyToULINT(val) }

// AnyToBOOL is deprecated: use convert.AnyToBOOL.
func AnyToBOOL[T any](val T) (iec.BOOL, error) { return convert.AnyToBOOL(val) }

// ConvertTo is deprecated: use convert.ConvertTo.
func ConvertTo[T any](in any) (T, error) { return convert.ConvertTo[T](in) }

// SubByte is deprecated: use convert.SubByte.
func SubByte(in any) (iec.BYTE, error) { return convert.SubByte(in) }

// SubWord is deprecated: use convert.SubWord.
func SubWord(in any) (iec.WORD, error) { return convert.SubWord(in) }

// SubDword is deprecated: use convert.SubDword.
func SubDword(in any) (iec.DWORD, error) { return convert.SubDword(in) }

// SubLword is deprecated: use convert.SubLword.
func SubLword(in any) (iec.LWORD, error) { return convert.SubLword(in) }

// SubDt is deprecated: use convert.SubDt.
func SubDt(in any) (iec.DT, error) { return convert.SubDt(in) }

// SubDate is deprecated: use convert.SubDate.
func SubDate(in any) (iec.DATE, error) { return convert.SubDate(in) }

// SubTod is deprecated: use convert.SubTod.
func SubTod(in any) (iec.TOD, error) { return convert.SubTod(in) }

// SubTime is deprecated: use convert.SubTime.
func SubTime(in any) (iec.TIME, error) { return convert.SubTime(in) }

// SubBool is deprecated: use convert.SubBool.
func SubBool(in iec.BOOL) iec.INT { return convert.SubBool(in) }

// ClampLINT is deprecated: use convert.ClampLINT.
func ClampLINT(val, min, max iec.LINT) iec.LINT { return convert.ClampLINT(val, min, max) }

// ClampULINT is deprecated: use convert.ClampULINT.
func ClampULINT(val, max iec.ULINT) iec.ULINT { return convert.ClampULINT(val, max) }

// ClampLREAL is deprecated: use convert.ClampLREAL.
func ClampLREAL(val, min, max iec.LREAL) iec.LREAL { return convert.ClampLREAL(val, min, max) }

// RoundAndClampLREAL is deprecated: use convert.RoundAndClampLREAL.
func RoundAndClampLREAL(val, min, max iec.LREAL) iec.LREAL {
	return convert.RoundAndClampLREAL(val, min, max)
}

// AlmostEqual is deprecated: use convert.AlmostEqual.
func AlmostEqual(a, b iec.LREAL) bool { return convert.AlmostEqual(a, b) }
