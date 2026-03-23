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
// The result type is determined by the largest of the input types,
// following IEC 61131-3 type promotion rules.
func AND(inputs []interface{}) (interface{}, error) {
	if len(inputs) == 0 {
		// IEC 61131-3 does not define behavior for empty inputs. Returning 0 is a safe default.
		return LWORD(0), nil
	}

	if len(inputs) == 1 {
		return inputs[0], nil // With one input, the result is the input itself.
	}

	// Initialize accumulator with all bits set to 1.
	acc, err := anyToULINT(inputs[0])
	if err != nil {
		return nil, fmt.Errorf("AND: error converting first element %v (type %T) to ULINT: %w", inputs[0], inputs[0], err)
	}
	for i := 1; i < len(inputs); i++ {
		val, err := anyToULINT(inputs[i])
		if err != nil {
			return nil, fmt.Errorf("AND: error converting element %v (type %T) to ULINT: %w", inputs[i], inputs[i], err)
		}
		acc &= val
	}

	targetType := reflect.TypeOf(inputs[len(inputs)-1])
	// Determine the largest integer-like type present
	targetType := reflect.TypeOf(inputs[0])
	for _, num := range inputs {
		if getTypeRank(reflect.TypeOf(num)) > getTypeRank(targetType) {
			targetType = reflect.TypeOf(num)
		}
	}
	result, err := convertToTargetType(acc, targetType)
	if err != nil {
		return nil, fmt.Errorf("AND: error converting final result to target type %v: %w", targetType, err)
	}
	return result, nil
}

// NOT performs a bitwise NOT on a single ANY_BIT type.
func NOT(in interface{}) (interface{}, error) {
	val, err := anyToULINT(in)
	if err != nil {
		return nil, fmt.Errorf("NOT: error converting %v (type %T) to ULINT: %w", in, in, err)
	}

	acc := ^val

	targetType := reflect.TypeOf(in)
	result, err := convertToTargetType(acc, targetType)
	if err != nil {
		return nil, fmt.Errorf("NOT: error converting final result to target type %v: %w", targetType, err)
	}
	return result, nil
}

// OR performs a bitwise OR on a slice of ANY_BIT types.
// The result type is determined by the largest of the input types,
// following IEC 61131-3 type promotion rules.
func OR(inputs []interface{}) (interface{}, error) {
	if len(inputs) == 0 {
		return LWORD(0), nil
	}

	if len(inputs) == 1 {
		return inputs[0], nil
	}

	acc, err := anyToULINT(inputs[0])
	if err != nil {
		return nil, fmt.Errorf("OR: error converting first element %v (type %T) to ULINT: %w", inputs[0], inputs[0], err)
	}

	for i := 1; i < len(inputs); i++ {
		val, err := anyToULINT(inputs[i])
		if err != nil {
			return nil, fmt.Errorf("OR: error converting element %v (type %T) to ULINT: %w", inputs[i], inputs[i], err)
		}
		acc |= val
	}

	targetType := reflect.TypeOf(inputs[len(inputs)-1])
	// Determine the largest integer-like type present
	targetType := reflect.TypeOf(inputs[0])
	for _, num := range inputs {
		if getTypeRank(reflect.TypeOf(num)) > getTypeRank(targetType) {
			targetType = reflect.TypeOf(num)
		}
	}
	result, err := convertToTargetType(acc, targetType)
	if err != nil {
		return nil, fmt.Errorf("OR: error converting final result to target type %v: %w", targetType, err)
	}
	return result, nil
}

// XOR performs a bitwise XOR on a slice of ANY_BIT types.
// The result type is determined by the largest of the input types,
// following IEC 61131-3 type promotion rules.
func XOR(inputs []interface{}) (interface{}, error) {
	if len(inputs) == 0 {
		return LWORD(0), nil
	}

	if len(inputs) == 1 {
		return inputs[0], nil
	}

	acc, err := anyToULINT(inputs[0])
	if err != nil {
		return nil, fmt.Errorf("XOR: error converting first element %v (type %T) to ULINT: %w", inputs[0], inputs[0], err)
	}

	for i := 1; i < len(inputs); i++ {
		val, err := anyToULINT(inputs[i])
		if err != nil {
			return nil, fmt.Errorf("XOR: error converting element %v (type %T) to ULINT: %w", inputs[i], inputs[i], err)
		}
		acc ^= val
	}

	targetType := reflect.TypeOf(inputs[len(inputs)-1])
	// Determine the largest integer-like type present
	targetType := reflect.TypeOf(inputs[0])
	for _, num := range inputs {
		if getTypeRank(reflect.TypeOf(num)) > getTypeRank(targetType) {
			targetType = reflect.TypeOf(num)
		}
	}
	result, err := convertToTargetType(acc, targetType)
	if err != nil {
		return nil, fmt.Errorf("XOR: error converting final result to target type %v: %w", targetType, err)
	}
	return result, nil
}

// bitShift is a generic helper for all shift and rotate operations.
func bitShift(op string, in interface{}, n int) (interface{}, error) {
	if n < 0 {
		return nil, fmt.Errorf("%s: shift/rotate count 'n' cannot be negative, got %d", op, n)
	}

	val, err := anyToULINT(in)
	if err != nil {
		return nil, fmt.Errorf("%s: error converting %v (type %T) to ULINT: %w", op, in, in, err)
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
		return nil, fmt.Errorf("bitShift: unknown operation: %s", op)
	}

	result, err := convertToTargetType(acc, targetType)
	if err != nil {
		return nil, fmt.Errorf("%s: error converting final result to target type %v: %w", op, targetType, err)
	}
	return result, nil
}

// ROL performs a bitwise rotation to the left.
func ROL(in interface{}, n int) (interface{}, error) {
	return bitShift("ROL", in, n)
}

// ROR performs a bitwise rotation to the right.
func ROR(in interface{}, n int) (interface{}, error) {
	return bitShift("ROR", in, n)
}

// SHL performs a bitwise left shift, filling with zeros.
func SHL(in interface{}, n int) (interface{}, error) {
	return bitShift("SHL", in, n)
}

// SHR performs a bitwise right shift, filling with zeros.
func SHR(in interface{}, n int) (interface{}, error) {
	return bitShift("SHR", in, n)
}
