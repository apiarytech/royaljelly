package royaljelly

import (
	"errors"
	"fmt"
	"testing"
)

func TestHasBit(t *testing.T) {
	if !HasBit(BYTE(0b1000), 3) {
		t.Error("HasBit(0b1000, 3) should be true")
	}
	if HasBit(BYTE(0b1000), 2) {
		t.Error("HasBit(0b1000, 2) should be false")
	}
	if !HasBit(DINT(-1), 31) {
		t.Error("HasBit(-1, 31) should be true for DINT")
	}
}

func TestSetBit(t *testing.T) {
	testCases := []struct {
		name     string
		n        interface{}
		pos      uint
		expected interface{}
	}{
		{"SINT", SINT(0b1010), uint(0), SINT(0b1011)},
		{"INT", INT(0), uint(14), INT(1 << 14)},
		{"DINT", DINT(0x12345670), uint(3), DINT(0x12345678)},
		{"LINT", LINT(0x0123456789ABCDEF), uint(60), LINT(0x1123456789ABCDEF)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Type assertion to the expected type for the operation
			// This ensures that tc.n and tc.expected are passed as their concrete ANY_INT type.
			switch v := tc.n.(type) {
			case SINT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(SINT))
			case INT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(INT))
			case DINT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(DINT))
			case LINT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(LINT))
			case USINT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(USINT))
			case UINT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(UINT))
			default:
				t.Fatalf("unhandled type for SetBit test case: %T", v)
			}
		})
	}
}

func TestClearBit(t *testing.T) {
	testCases := []struct {
		name     string
		n        interface{}
		pos      uint
		expected interface{}
	}{
		{"SINT", SINT(0b1011), uint(0), SINT(0b1010)},
		{"INT", USINT(0b10000000), uint(8), USINT(1 << 7)},
		{"DINT", DINT(0x12345678), uint(3), DINT(0x12345670)},
		{"LINT", LINT(0x1123456789ABCDEF), uint(60), LINT(0x0123456789ABCDEF)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Type assertion to the expected type for the operation
			// This ensures that tc.n and tc.expected are passed as their concrete ANY_INT type.
			switch v := tc.n.(type) {
			case SINT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(SINT))
			case INT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(INT))
			case DINT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(DINT))
			case LINT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(LINT))
			case USINT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(USINT))
			case UINT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(UINT))
			default:
				t.Fatalf("unhandled type for ClearBit test case: %T", v)
			}
		})
	}
}

// genericBitwiseTest is a helper function to test generic bitwise operations like SetBit and ClearBit.
func genericBitwiseTest[T ANY_INT](t *testing.T, opName string, op func(T, uint) T, n T, pos uint, expected T) {
	t.Helper()
	result := op(n, pos)
	if result != expected {
		t.Errorf("%s(%v, %d) = %v (0x%X); want %v (0x%X)", opName, n, pos, result, result, expected, expected)
	}
}

func TestAND(t *testing.T) {
	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectError bool
		panicValue  error
	}{
		{"BYTEs", []interface{}{BYTE(0b1100), BYTE(0b1010)}, BYTE(0b1000), false, nil},
		{"WORDs", []interface{}{WORD(0xFF00), WORD(0x00FF), WORD(0xFFFF)}, WORD(0x0000), false, nil},
		{"Mixed Types", []interface{}{LWORD(0xFF), USINT(0x0F)}, USINT(0x0F), false, nil},
		{"Empty", []interface{}{}, LWORD(0), false, nil},
		{"Single", []interface{}{DINT(123)}, DINT(123), false, nil},
		{"Invalid String", []interface{}{STRING("not-a-number"), INT(5)}, nil, true, errors.New("anyToULINT: cannot parse STRING 'not-a-number' to ULINT")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := AND(tc.inputs)

			if tc.expectError {
				if err == nil {
					t.Errorf("AND(%v) did not return an error; expected error", tc.inputs)
				}
			} else {
				if err != nil {
					t.Errorf("AND(%v) returned an unexpected error: %v", tc.inputs, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("AND(%v) = %v (type %T); want %v (type %T)", tc.inputs, result, result, tc.expected, tc.expected)
				}
			}
		})
	}
}

func TestOR(t *testing.T) {
	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectError bool
	}{
		{"BYTEs", []interface{}{BYTE(0b1100), BYTE(0b1010)}, BYTE(0b1110), false},
		{"WORDs", []interface{}{WORD(0xFF00), WORD(0x00FF)}, WORD(0xFFFF), false},
		{"Mixed Types", []interface{}{LWORD(0xF0), USINT(0x0F)}, USINT(0xFF), false},
		{"Empty", []interface{}{}, LWORD(0), false},
		{"Single", []interface{}{DINT(123)}, DINT(123), false},
		{"Invalid String", []interface{}{STRING("abc"), INT(5)}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := OR(tc.inputs)
			if tc.expectError {
				if err == nil {
					t.Errorf("OR(%v) did not return an error; expected error", tc.inputs)
				}
			} else {
				if err != nil {
					t.Errorf("OR(%v) returned an unexpected error: %v", tc.inputs, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("OR(%v) = %v; want %v", tc.inputs, result, tc.expected)
				}
			}
		})
	}
}

func TestXOR(t *testing.T) {
	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectError bool
	}{
		{"BYTEs", []interface{}{BYTE(0b1100), BYTE(0b1010)}, BYTE(0b0110), false},
		{"WORDs", []interface{}{WORD(0xFF00), WORD(0xFFFF)}, WORD(0x00FF), false},
		{"Mixed Types", []interface{}{LWORD(0xF0), USINT(0xFF)}, USINT(0x0F), false},
		{"Empty", []interface{}{}, LWORD(0), false},
		{"Single", []interface{}{DINT(123)}, DINT(123), false},
		{"Invalid String", []interface{}{STRING("abc"), INT(5)}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := XOR(tc.inputs)
			if tc.expectError {
				if err == nil {
					t.Errorf("XOR(%v) did not return an error; expected error", tc.inputs)
				}
			} else {
				if err != nil {
					t.Errorf("XOR(%v) returned an unexpected error: %v", tc.inputs, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("XOR(%v) = %v; want %v", tc.inputs, result, tc.expected)
				}
			}
		})
	}
}

func TestNOT(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		expected    interface{}
		expectError bool
	}{
		{"BYTE", BYTE(0b11110000), BYTE(0b00001111), false},
		{"WORD", WORD(0x00FF), WORD(0xFF00), false},
		{"DINT", DINT(0), DINT(-1), false},
		{"Invalid String", STRING("abc"), nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NOT(tc.input)
			if tc.expectError {
				if err == nil {
					t.Errorf("NOT(%v) did not return an error; expected error", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("NOT(%v) returned an unexpected error: %v", tc.input, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("NOT(%v) = %v; want %v", tc.input, result, tc.expected)
				}
			}
		})
	}
}

func TestSHL(t *testing.T) {
	testCases := []struct {
		name     string
		in       interface{}
		n        int
		expected interface{}
	}{
		{"BYTE", BYTE(0b00001111), 4, BYTE(0b11110000)},
		{"WORD", WORD(1), 15, WORD(32768)},
		{"UDINT", UDINT(0x0FFFFFFF), 4, UDINT(0xFFFFFFF0)},
		{"LINT", LINT(1), 62, LINT(0x4000000000000000)},
		{"Shift by 0", INT(123), 0, INT(123)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SHL(tc.in, tc.n)
			if err != nil {
				t.Fatalf("SHL(%v, %d) returned an unexpected error: %v", tc.in, tc.n, err)
			}
			if result != tc.expected {
				t.Errorf("SHL(%v, %d) = %v; want %v", tc.in, tc.n, result, tc.expected)
			}
		})
	}
}

func TestSHR(t *testing.T) {
	testCases := []struct {
		name     string
		in       interface{}
		n        int
		expected interface{}
	}{
		{"BYTE", BYTE(0b11110000), 4, BYTE(0b00001111)},
		{"WORD", WORD(32768), 15, WORD(1)},
		{"UDINT", UDINT(0xFFFFFFFF), 1, UDINT(0x7FFFFFFF)}, // Logical shift for signed (IEC 61131-3 compliant)
		{"LINT", LINT(0x4000000000000000), 4, LINT(0x0400000000000000)},
		{"Shift by 0", INT(123), 0, INT(123)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := SHR(tc.in, tc.n)
			if err != nil {
				t.Fatalf("SHR(%v, %d) returned an unexpected error: %v", tc.in, tc.n, err)
			}
			if result != tc.expected {
				t.Errorf("SHR(%v, %d) = %v; want %v", tc.in, tc.n, result, tc.expected)
			}
		})
	}
}

func TestROL(t *testing.T) {
	testCases := []struct {
		name        string
		in          interface{}
		n           int
		expected    interface{}
		expectError bool
	}{
		{"BYTE", BYTE(0b11000001), 1, BYTE(0b10000011), false},
		{"WORD", WORD(0x8001), 1, WORD(0x0003), false},
		{"DINT", UDINT(0xC0000000), 2, UDINT(0x00000003), false},
		{"LINT", LINT(1), 64, LINT(1), false}, // Rotate by full width
		{"Rotate by 0", LINT(123), 0, LINT(123), false},
		{"Negative shift", BYTE(1), -1, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ROL(tc.in, tc.n)

			if tc.expectError {
				if err == nil {
					t.Errorf("ROL(%v, %d) did not return an error; expected error", tc.in, tc.n)
				}
			} else {
				if err != nil {
					t.Errorf("ROL(%v, %d) returned an unexpected error: %v", tc.in, tc.n, err)
				}
				if result != tc.expected {
					t.Errorf("ROL(%v, %d) = %v; want %v", tc.in, tc.n, result, tc.expected)
				}
			}
		})
	}
}

func TestROR(t *testing.T) {
	testCases := []struct {
		name        string
		in          interface{}
		n           int
		expected    interface{}
		expectError bool
	}{
		{"BYTE", BYTE(0b11000001), 1, BYTE(0b11100000), false},
		{"WORD", WORD(0x0003), 1, WORD(0x8001), false},
		{"DINT", UDINT(0x00000003), 2, UDINT(0xC0000000), false},
		{"LINT", LINT(1), 64, LINT(1), false}, // Rotate by full width
		{"Rotate by 0", INT(123), 0, INT(123), false},
		{"Negative shift", BYTE(1), -1, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ROR(tc.in, tc.n)

			if tc.expectError {
				if err == nil {
					t.Errorf("ROR(%v, %d) did not return an error; expected error", tc.in, tc.n)
				}
			} else {
				if err != nil {
					t.Errorf("ROR(%v, %d) returned an unexpected error: %v", tc.in, tc.n, err)
				}
				if result != tc.expected {
					t.Errorf("ROR(%v, %d) = %v; want %v", tc.in, tc.n, result, tc.expected)
				}
			}
		})
	}
}
