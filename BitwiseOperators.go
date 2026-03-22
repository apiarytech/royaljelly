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
	"fmt"
	"math/bits"
	"reflect"
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

// AND performs a bitwise AND on a slice of ANY_BIT types.
// The result type is determined by the type of the last element.
func AND(inputs []interface{}) interface{} {
	if len(inputs) == 0 {
		// IEC 61131-3 does not define behavior for empty inputs. Returning 0 is a safe default.
		return LWORD(0)
	}

	if len(inputs) == 1 {
		return inputs[0] // With one input, the result is the input itself.
	}

	// Initialize accumulator with all bits set to 1.
	acc, err := anyToULINT(inputs[0])
	if err != nil {
		panic(fmt.Sprintf("AND: error converting first element %v (type %T) to ULINT: %v", inputs[0], inputs[0], err))
	}
	for i := 1; i < len(inputs); i++ {
		val, err := anyToULINT(inputs[i])
		if err != nil {
			panic(fmt.Sprintf("AND: error converting element %v (type %T) to ULINT: %v", inputs[i], inputs[i], err))
		}
		acc &= val
	}

	targetType := reflect.TypeOf(inputs[len(inputs)-1])
	result, err := convertToTargetType(acc, targetType)
	if err != nil {
		panic(fmt.Sprintf("AND: error converting final result to target type %v: %v", targetType, err))
	}
	return result
}

// NOT performs a bitwise NOT on a single ANY_BIT type.
func NOT(in interface{}) interface{} {
	val, err := anyToULINT(in)
	if err != nil {
		panic(fmt.Sprintf("NOT: error converting %v (type %T) to ULINT: %v", in, in, err))
	}

	acc := ^val

	targetType := reflect.TypeOf(in)
	result, err := convertToTargetType(acc, targetType)
	if err != nil {
		panic(fmt.Sprintf("NOT: error converting final result to target type %v: %v", targetType, err))
	}
	return result
}

// OR performs a bitwise OR on a slice of ANY_BIT types.
func OR(inputs []interface{}) interface{} {
	if len(inputs) == 0 {
		return LWORD(0)
	}

	if len(inputs) == 1 {
		return inputs[0]
	}

	acc, err := anyToULINT(inputs[0])
	if err != nil {
		panic(fmt.Sprintf("OR: error converting first element %v (type %T) to ULINT: %v", inputs[0], inputs[0], err))
	}

	for i := 1; i < len(inputs); i++ {
		val, err := anyToULINT(inputs[i])
		if err != nil {
			panic(fmt.Sprintf("OR: error converting element %v (type %T) to ULINT: %v", inputs[i], inputs[i], err))
		}
		acc |= val
	}

	targetType := reflect.TypeOf(inputs[len(inputs)-1])
	result, err := convertToTargetType(acc, targetType)
	if err != nil {
		panic(fmt.Sprintf("OR: error converting final result to target type %v: %v", targetType, err))
	}
	return result
}

// XOR performs a bitwise XOR on a slice of ANY_BIT types.
func XOR(inputs []interface{}) interface{} {
	if len(inputs) == 0 {
		return LWORD(0)
	}

	if len(inputs) == 1 {
		return inputs[0]
	}

	acc, err := anyToULINT(inputs[0])
	if err != nil {
		panic(fmt.Sprintf("XOR: error converting first element %v (type %T) to ULINT: %v", inputs[0], inputs[0], err))
	}

	for i := 1; i < len(inputs); i++ {
		val, err := anyToULINT(inputs[i])
		if err != nil {
			panic(fmt.Sprintf("XOR: error converting element %v (type %T) to ULINT: %v", inputs[i], inputs[i], err))
		}
		acc ^= val
	}

	targetType := reflect.TypeOf(inputs[len(inputs)-1])
	result, err := convertToTargetType(acc, targetType)
	if err != nil {
		panic(fmt.Sprintf("XOR: error converting final result to target type %v: %v", targetType, err))
	}
	return result
}

// bitShift is a generic helper for all shift and rotate operations.
func bitShift(op string, in interface{}, n int) interface{} {
	if n < 0 {
		panic(fmt.Sprintf("%s: shift/rotate count 'n' cannot be negative, got %d", op, n))
	}

	val, err := anyToULINT(in)
	if err != nil {
		panic(fmt.Sprintf("%s: error converting %v (type %T) to ULINT: %v", op, in, in, err))
	}

	targetType := reflect.TypeOf(in)
	bitSize := targetType.Bits()

	var acc ULINT

	switch op {
	case "SHL":
		acc = val << uint(n)
	case "SHR":
		acc = val >> uint(n)
	case "ROL":
		switch bitSize {
		case 8:
			acc = ULINT(bits.RotateLeft8(uint8(val), n))
		case 16:
			acc = ULINT(bits.RotateLeft16(uint16(val), n))
		case 32:
			acc = ULINT(bits.RotateLeft32(uint32(val), n))
		default:
			acc = ULINT(bits.RotateLeft64(uint64(val), n))
		}
	case "ROR":
		switch bitSize {
		case 8:
			acc = ULINT(bits.RotateLeft8(uint8(val), -n))
		case 16:
			acc = ULINT(bits.RotateLeft16(uint16(val), -n))
		case 32:
			acc = ULINT(bits.RotateLeft32(uint32(val), -n))
		default:
			acc = ULINT(bits.RotateLeft64(uint64(val), -n))
		}
	default:
		panic("bitShift: unknown operation: " + op)
	}

	result, err := convertToTargetType(acc, targetType)
	if err != nil {
		panic(fmt.Sprintf("%s: error converting final result to target type %v: %v", op, targetType, err))
	}
	return result
}

// ROL performs a bitwise rotation to the left.
func ROL(in interface{}, n int) interface{} {
	return bitShift("ROL", in, n)
}

// ROR performs a bitwise rotation to the right.
func ROR(in interface{}, n int) interface{} {
	return bitShift("ROR", in, n)
}

// SHL performs a bitwise left shift, filling with zeros.
func SHL(in interface{}, n int) interface{} {
	return bitShift("SHL", in, n)
}

// SHR performs a bitwise right shift, filling with zeros.
func SHR(in interface{}, n int) interface{} {
	return bitShift("SHR", in, n)
}
