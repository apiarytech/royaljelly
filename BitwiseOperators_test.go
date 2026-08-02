package royaljelly

import (
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
		name     string
		testFunc func(*testing.T)
	}{
		{"BYTEs", func(t *testing.T) {
			result := AND(BYTE(0b1100), BYTE(0b1010))
			expected := BYTE(0b1000)
			if result != expected {
				t.Errorf("AND() = %v; want %v", result, expected)
			}
		}},
		{"WORDs", func(t *testing.T) {
			result := AND(WORD(0xFF00), WORD(0x00FF), WORD(0xFFFF))
			expected := WORD(0x0000)
			if result != expected {
				t.Errorf("AND() = %v; want %v", result, expected)
			}
		}},
		{"Empty", func(t *testing.T) {
			result := AND[UINT]()
			expected := UINT(0)
			if result != expected {
				t.Errorf("AND() = %v; want %v", result, expected)
			}
		}},
		{"BOOLs", func(t *testing.T) {
			result := AND_BOOL(BOOL(true), BOOL(true), BOOL(false))
			expected := BOOL(false)
			if result != expected {
				t.Errorf("AND() with BOOLs = %v; want %v", result, expected)
			}
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}

func TestOR(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"BYTEs", func(t *testing.T) {
			result := OR(BYTE(0b1100), BYTE(0b1010))
			expected := BYTE(0b1110)
			if result != expected {
				t.Errorf("OR() = %v; want %v", result, expected)
			}
		}},
		{"WORDs", func(t *testing.T) {
			result := OR(WORD(0xFF00), WORD(0x00FF))
			expected := WORD(0xFFFF)
			if result != expected {
				t.Errorf("OR() = %v; want %v", result, expected)
			}
		}},
		{"BOOLs", func(t *testing.T) {
			result := OR_BOOL(BOOL(true), BOOL(false), BOOL(false))
			expected := BOOL(true)
			if result != expected {
				t.Errorf("OR() with BOOLs = %v; want %v", result, expected)
			}
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}

func TestXOR(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"BYTEs", func(t *testing.T) {
			result := XOR(BYTE(0b1100), BYTE(0b1010))
			expected := BYTE(0b0110)
			if result != expected {
				t.Errorf("XOR() = %v; want %v", result, expected)
			}
		}},
		{"WORDs", func(t *testing.T) {
			result := XOR(WORD(0xFF00), WORD(0xFFFF))
			expected := WORD(0x00FF)
			if result != expected {
				t.Errorf("XOR() = %v; want %v", result, expected)
			}
		}},
		{"BOOLs", func(t *testing.T) {
			result := XOR_BOOL(BOOL(true), BOOL(true), BOOL(false))
			expected := BOOL(false)
			if result != expected {
				t.Errorf("XOR() with BOOLs = %v; want %v", result, expected)
			}
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}

func TestNOT(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"BYTE", func(t *testing.T) {
			result := NOT(BYTE(0b11110000))
			expected := BYTE(0b00001111)
			if result != expected {
				t.Errorf("NOT() = %v; want %v", result, expected)
			}
		}},
		{"DINT", func(t *testing.T) {
			result := NOT(DINT(0))
			expected := DINT(-1)
			if result != expected {
				t.Errorf("NOT() = %v; want %v", result, expected)
			}
		}},
		{"BOOL", func(t *testing.T) {
			result := NOT_BOOL(BOOL(true))
			expected := BOOL(false)
			if result != expected {
				t.Errorf("NOT() with BOOL = %v; want %v", result, expected)
			}
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.testFunc(t)
		})
	}
}

func TestSHL(t *testing.T) {
	testCases := []struct {
		name     string
		in       interface{}
		n        uint
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
			switch v := tc.in.(type) {
			case BYTE:
				if SHL(v, tc.n) != tc.expected.(BYTE) {
					t.Errorf("SHL failed")
				}
			case WORD:
				if SHL(v, tc.n) != tc.expected.(WORD) {
					t.Errorf("SHL failed")
				}
			case UDINT:
				if SHL(v, tc.n) != tc.expected.(UDINT) {
					t.Errorf("SHL failed")
				}
			case LINT:
				if SHL(v, tc.n) != tc.expected.(LINT) {
					t.Errorf("SHL failed")
				}
			case INT:
				if SHL(v, tc.n) != tc.expected.(INT) {
					t.Errorf("SHL failed")
				}
			default:
				t.Fatalf("unhandled type for SHL test case: %T", v)
			}
		})
	}
}

func TestSHR(t *testing.T) {
	testCases := []struct {
		name     string
		in       interface{}
		n        uint
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
			switch v := tc.in.(type) {
			case BYTE:
				if SHR(v, tc.n) != tc.expected.(BYTE) {
					t.Errorf("SHR failed")
				}
			case WORD:
				if SHR(v, tc.n) != tc.expected.(WORD) {
					t.Errorf("SHR failed")
				}
			case UDINT:
				if SHR(v, tc.n) != tc.expected.(UDINT) {
					t.Errorf("SHR failed")
				}
			case LINT:
				if SHR(v, tc.n) != tc.expected.(LINT) {
					t.Errorf("SHR failed")
				}
			case INT:
				if SHR(v, tc.n) != tc.expected.(INT) {
					t.Errorf("SHR failed")
				}
			default:
				t.Fatalf("unhandled type for SHR test case: %T", v)
			}
		})
	}
}

func TestROL(t *testing.T) {
	testCases := []struct {
		name     string
		in       interface{}
		n        int
		expected interface{}
	}{
		{"BYTE", BYTE(0b11000001), 1, BYTE(0b10000011)},
		{"WORD", WORD(0x8001), 1, WORD(0x0003)},
		{"DINT", UDINT(0xC0000000), 2, UDINT(0x00000003)},
		{"LINT", LINT(1), 64, LINT(1)}, // Rotate by full width
		{"Rotate by 0", LINT(123), 0, LINT(123)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch v := tc.in.(type) {
			case BYTE:
				if ROL(v, tc.n) != tc.expected.(BYTE) {
					t.Errorf("ROL failed")
				}
			case WORD:
				if ROL(v, tc.n) != tc.expected.(WORD) {
					t.Errorf("ROL failed")
				}
			case UDINT:
				if ROL(v, tc.n) != tc.expected.(UDINT) {
					t.Errorf("ROL failed")
				}
			case LINT:
				if ROL(v, tc.n) != tc.expected.(LINT) {
					t.Errorf("ROL failed")
				}
			case INT:
				if ROL(v, tc.n) != tc.expected.(INT) {
					t.Errorf("ROL failed")
				}
			default:
				t.Fatalf("unhandled type for ROL test case: %T", v)
			}
		})
	}
}

func TestROR(t *testing.T) {
	testCases := []struct {
		name     string
		in       interface{}
		n        int
		expected interface{}
	}{
		{"BYTE", BYTE(0b11000001), 1, BYTE(0b11100000)},
		{"WORD", WORD(0x0003), 1, WORD(0x8001)},
		{"DINT", UDINT(0x00000003), 2, UDINT(0xC0000000)},
		{"LINT", LINT(1), 64, LINT(1)}, // Rotate by full width
		{"Rotate by 0", INT(123), 0, INT(123)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch v := tc.in.(type) {
			case BYTE:
				if ROR(v, tc.n) != tc.expected.(BYTE) {
					t.Errorf("ROR failed")
				}
			case WORD:
				if ROR(v, tc.n) != tc.expected.(WORD) {
					t.Errorf("ROR failed")
				}
			case UDINT:
				if ROR(v, tc.n) != tc.expected.(UDINT) {
					t.Errorf("ROR failed")
				}
			case LINT:
				if ROR(v, tc.n) != tc.expected.(LINT) {
					t.Errorf("ROR failed")
				}
			case INT:
				if ROR(v, tc.n) != tc.expected.(INT) {
					t.Errorf("ROR failed")
				}
			default:
				t.Fatalf("unhandled type for ROR test case: %T", v)
			}
		})
	}
}
