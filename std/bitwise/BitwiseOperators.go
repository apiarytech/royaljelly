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

package bitwise

import (
	"unsafe"

	. "github.com/apiarytech/royaljelly/iec"
)

// SetBit sets a bit at pos in the integer n. Returns the modified value.
func SetBit[T ANY_INT](n T, pos uint) T {
	// Cast 1 to the generic type T before shifting. This ensures the bitmask
	// is created with the correct type and width, preventing overflow issues.
	return n | (T(1) << pos)
}

// ClearBit clears a bit at pos in n. Returns the modified value.
func ClearBit[T ANY_INT](n T, pos uint) T {
	mask := ^(T(1) << pos)
	return n & mask
}

// HasBit checks if a bit at pos in n is set.
func HasBit[T ANY_INT](n T, pos uint) BOOL {
	val := n & (T(1) << pos)
	return val != 0
}

// AND performs a bitwise AND on a slice of unsigned integer types.
// All inputs must be of the same type.
func AND[T ANY_UINTS](inputs ...T) T {
	if len(inputs) == 0 {
		var zero T
		return zero
	}
	acc := inputs[0]
	for i := 1; i < len(inputs); i++ {
		acc &= inputs[i]
	}
	return acc
}

// Overload for BOOL type to perform logical AND.
func AND_BOOL[T ANY_BOOL](inputs ...T) T {
	if len(inputs) == 0 {
		var zero T
		return zero
	}
	acc := inputs[0]
	for i := 1; i < len(inputs); i++ {
		acc = acc && inputs[i]
	}
	return acc
}

// NOT performs a bitwise NOT on a single ANY_BIT type.
func NOT[T ANY_INT](inputs T) T {
	return ^inputs
}

// NOT_BOOL performs a logical NOT on a single BOOL.
func NOT_BOOL[T ANY_BOOL](inputs T) T {
	return !inputs
}

// OR performs a bitwise OR on a slice of ANY_BIT types.
// All inputs must be of the same type.
func OR[T ANY_INT](inputs ...T) T {
	if len(inputs) == 0 {
		var zero T
		return zero
	}
	acc := inputs[0]
	for i := 1; i < len(inputs); i++ {
		acc |= inputs[i]
	}
	return acc
}

// Overload for BOOL type to perform logical OR.
func OR_BOOL[T ANY_BOOL](inputs ...T) T {
	if len(inputs) == 0 {
		var zero T
		return zero
	}
	acc := inputs[0]
	for i := 1; i < len(inputs); i++ {
		acc = acc || inputs[i]
	}
	return acc
}

func XOR[T ANY_UINTS](inputs ...T) T {
	if len(inputs) == 0 {
		var zero T
		return zero
	}
	acc := inputs[0]
	for i := 1; i < len(inputs); i++ {
		acc ^= inputs[i]
	}
	return acc
}

func XOR_BOOL[T ANY_BOOL](inputs ...T) T {
	if len(inputs) == 0 {
		var zero T
		return zero
	}
	acc := inputs[0]
	for i := 1; i < len(inputs); i++ {
		acc = acc != inputs[i]
	}
	return acc
}

// SHL performs a bitwise left shift, filling with zeros.
func SHL[T ANY_INT](in T, n uint) T {
	return in << n
}

// SHR performs a bitwise right shift, filling with zeros.
func SHR[T ANY_INT](in T, n uint) T {
	return in >> n
}

// ROL performs a bitwise rotation to the left.
func ROL[T ANY_INT](in T, n int) T {
	bitSize := unsafe.Sizeof(in) * 8
	return (in << uint(n)) | (in >> (uint(bitSize) - uint(n)))
}

// ROR performs a bitwise rotation to the right.
func ROR[T ANY_INT](in T, n int) T {
	bitSize := unsafe.Sizeof(in) * 8
	return (in >> uint(n)) | (in << (uint(bitSize) - uint(n)))
}
