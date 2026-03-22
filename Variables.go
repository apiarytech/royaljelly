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
