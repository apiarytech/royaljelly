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
	"time"
)

const (
	base = 255
)

type Addresses struct {
	B  [base]BOOL
	C  [base]BYTE
	D  [base]DWORD
	L  [base]LWORD
	R  [base]REAL
	LR [base]LREAL
	S  [base]STRING
	W  [base]WSTRING
	T  TIME
}

var (
	I        Addresses
	Q        Addresses
	M        Addresses
	IEC_TIME time.Time
)
