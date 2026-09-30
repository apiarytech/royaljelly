package bitwise

import (
	"testing"

	"github.com/apiarytech/royaljelly/iec"
)

func TestHasBit(t *testing.T) {
	if !HasBit(iec.BYTE(0b1000), 3) {
		t.Error("HasBit(0b1000, 3) should be true")
	}
	if HasBit(iec.BYTE(0b1000), 2) {
		t.Error("HasBit(0b1000, 2) should be false")
	}
	if !HasBit(iec.DINT(-1), 31) {
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
		{"SINT", iec.SINT(0b1010), uint(0), iec.SINT(0b1011)},
		{"INT", iec.INT(0), uint(14), iec.INT(1 << 14)},
		{"DINT", iec.DINT(0x12345670), uint(3), iec.DINT(0x12345678)},
		{"LINT", iec.LINT(0x0123456789ABCDEF), uint(60), iec.LINT(0x1123456789ABCDEF)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Type assertion to the expected type for the operation
			// This ensures that tc.n and tc.expected are passed as their concrete ANY_INT type.
			switch v := tc.n.(type) {
			case iec.SINT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(iec.SINT))
			case iec.INT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(iec.INT))
			case iec.DINT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(iec.DINT))
			case iec.LINT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(iec.LINT))
			case iec.USINT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(iec.USINT))
			case iec.UINT:
				genericBitwiseTest(t, "SetBit", SetBit, v, tc.pos, tc.expected.(iec.UINT))
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
		{"SINT", iec.SINT(0b1011), uint(0), iec.SINT(0b1010)},
		{"INT", iec.USINT(0b10000000), uint(8), iec.USINT(1 << 7)},
		{"DINT", iec.DINT(0x12345678), uint(3), iec.DINT(0x12345670)},
		{"LINT", iec.LINT(0x1123456789ABCDEF), uint(60), iec.LINT(0x0123456789ABCDEF)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Type assertion to the expected type for the operation
			// This ensures that tc.n and tc.expected are passed as their concrete ANY_INT type.
			switch v := tc.n.(type) {
			case iec.SINT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(iec.SINT))
			case iec.INT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(iec.INT))
			case iec.DINT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(iec.DINT))
			case iec.LINT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(iec.LINT))
			case iec.USINT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(iec.USINT))
			case iec.UINT:
				genericBitwiseTest(t, "ClearBit", ClearBit, v, tc.pos, tc.expected.(iec.UINT))
			default:
				t.Fatalf("unhandled type for ClearBit test case: %T", v)
			}
		})
	}
}

// genericBitwiseTest is a helper function to test generic bitwise operations like SetBit and ClearBit.
func genericBitwiseTest[T iec.ANY_INT](t *testing.T, opName string, op func(T, uint) T, n T, pos uint, expected T) {
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
			result := AND(iec.BYTE(0b1100), iec.BYTE(0b1010))
			expected := iec.BYTE(0b1000)
			if result != expected {
				t.Errorf("AND() = %v; want %v", result, expected)
			}
		}},
		{"WORDs", func(t *testing.T) {
			result := AND(iec.WORD(0xFF00), iec.WORD(0x00FF), iec.WORD(0xFFFF))
			expected := iec.WORD(0x0000)
			if result != expected {
				t.Errorf("AND() = %v; want %v", result, expected)
			}
		}},
		{"Empty", func(t *testing.T) {
			result := AND[iec.UINT]()
			expected := iec.UINT(0)
			if result != expected {
				t.Errorf("AND() = %v; want %v", result, expected)
			}
		}},
		{"BOOLs", func(t *testing.T) {
			result := AND_BOOL(iec.BOOL(true), iec.BOOL(true), iec.BOOL(false))
			expected := iec.BOOL(false)
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
			result := OR(iec.BYTE(0b1100), iec.BYTE(0b1010))
			expected := iec.BYTE(0b1110)
			if result != expected {
				t.Errorf("OR() = %v; want %v", result, expected)
			}
		}},
		{"WORDs", func(t *testing.T) {
			result := OR(iec.WORD(0xFF00), iec.WORD(0x00FF))
			expected := iec.WORD(0xFFFF)
			if result != expected {
				t.Errorf("OR() = %v; want %v", result, expected)
			}
		}},
		{"BOOLs", func(t *testing.T) {
			result := OR_BOOL(iec.BOOL(true), iec.BOOL(false), iec.BOOL(false))
			expected := iec.BOOL(true)
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
			result := XOR(iec.BYTE(0b1100), iec.BYTE(0b1010))
			expected := iec.BYTE(0b0110)
			if result != expected {
				t.Errorf("XOR() = %v; want %v", result, expected)
			}
		}},
		{"WORDs", func(t *testing.T) {
			result := XOR(iec.WORD(0xFF00), iec.WORD(0xFFFF))
			expected := iec.WORD(0x00FF)
			if result != expected {
				t.Errorf("XOR() = %v; want %v", result, expected)
			}
		}},
		{"BOOLs", func(t *testing.T) {
			result := XOR_BOOL(iec.BOOL(true), iec.BOOL(true), iec.BOOL(false))
			expected := iec.BOOL(false)
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
			result := NOT(iec.BYTE(0b11110000))
			expected := iec.BYTE(0b00001111)
			if result != expected {
				t.Errorf("NOT() = %v; want %v", result, expected)
			}
		}},
		{"DINT", func(t *testing.T) {
			result := NOT(iec.DINT(0))
			expected := iec.DINT(-1)
			if result != expected {
				t.Errorf("NOT() = %v; want %v", result, expected)
			}
		}},
		{"BOOL", func(t *testing.T) {
			result := NOT_BOOL(iec.BOOL(true))
			expected := iec.BOOL(false)
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
		{"BYTE", iec.BYTE(0b00001111), 4, iec.BYTE(0b11110000)},
		{"WORD", iec.WORD(1), 15, iec.WORD(32768)},
		{"UDINT", iec.UDINT(0x0FFFFFFF), 4, iec.UDINT(0xFFFFFFF0)},
		{"LINT", iec.LINT(1), 62, iec.LINT(0x4000000000000000)},
		{"Shift by 0", iec.INT(123), 0, iec.INT(123)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch v := tc.in.(type) {
			case iec.BYTE:
				if SHL(v, tc.n) != tc.expected.(iec.BYTE) {
					t.Errorf("SHL failed")
				}
			case iec.WORD:
				if SHL(v, tc.n) != tc.expected.(iec.WORD) {
					t.Errorf("SHL failed")
				}
			case iec.UDINT:
				if SHL(v, tc.n) != tc.expected.(iec.UDINT) {
					t.Errorf("SHL failed")
				}
			case iec.LINT:
				if SHL(v, tc.n) != tc.expected.(iec.LINT) {
					t.Errorf("SHL failed")
				}
			case iec.INT:
				if SHL(v, tc.n) != tc.expected.(iec.INT) {
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
		{"BYTE", iec.BYTE(0b11110000), 4, iec.BYTE(0b00001111)},
		{"WORD", iec.WORD(32768), 15, iec.WORD(1)},
		{"UDINT", iec.UDINT(0xFFFFFFFF), 1, iec.UDINT(0x7FFFFFFF)}, // Logical shift for signed (IEC 61131-3 compliant)
		{"LINT", iec.LINT(0x4000000000000000), 4, iec.LINT(0x0400000000000000)},
		{"Shift by 0", iec.INT(123), 0, iec.INT(123)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch v := tc.in.(type) {
			case iec.BYTE:
				if SHR(v, tc.n) != tc.expected.(iec.BYTE) {
					t.Errorf("SHR failed")
				}
			case iec.WORD:
				if SHR(v, tc.n) != tc.expected.(iec.WORD) {
					t.Errorf("SHR failed")
				}
			case iec.UDINT:
				if SHR(v, tc.n) != tc.expected.(iec.UDINT) {
					t.Errorf("SHR failed")
				}
			case iec.LINT:
				if SHR(v, tc.n) != tc.expected.(iec.LINT) {
					t.Errorf("SHR failed")
				}
			case iec.INT:
				if SHR(v, tc.n) != tc.expected.(iec.INT) {
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
		{"BYTE", iec.BYTE(0b11000001), 1, iec.BYTE(0b10000011)},
		{"WORD", iec.WORD(0x8001), 1, iec.WORD(0x0003)},
		{"DINT", iec.UDINT(0xC0000000), 2, iec.UDINT(0x00000003)},
		{"LINT", iec.LINT(1), 64, iec.LINT(1)}, // Rotate by full width
		{"Rotate by 0", iec.LINT(123), 0, iec.LINT(123)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch v := tc.in.(type) {
			case iec.BYTE:
				if ROL(v, tc.n) != tc.expected.(iec.BYTE) {
					t.Errorf("ROL failed")
				}
			case iec.WORD:
				if ROL(v, tc.n) != tc.expected.(iec.WORD) {
					t.Errorf("ROL failed")
				}
			case iec.UDINT:
				if ROL(v, tc.n) != tc.expected.(iec.UDINT) {
					t.Errorf("ROL failed")
				}
			case iec.LINT:
				if ROL(v, tc.n) != tc.expected.(iec.LINT) {
					t.Errorf("ROL failed")
				}
			case iec.INT:
				if ROL(v, tc.n) != tc.expected.(iec.INT) {
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
		{"BYTE", iec.BYTE(0b11000001), 1, iec.BYTE(0b11100000)},
		{"WORD", iec.WORD(0x0003), 1, iec.WORD(0x8001)},
		{"DINT", iec.UDINT(0x00000003), 2, iec.UDINT(0xC0000000)},
		{"LINT", iec.LINT(1), 64, iec.LINT(1)}, // Rotate by full width
		{"Rotate by 0", iec.INT(123), 0, iec.INT(123)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch v := tc.in.(type) {
			case iec.BYTE:
				if ROR(v, tc.n) != tc.expected.(iec.BYTE) {
					t.Errorf("ROR failed")
				}
			case iec.WORD:
				if ROR(v, tc.n) != tc.expected.(iec.WORD) {
					t.Errorf("ROR failed")
				}
			case iec.UDINT:
				if ROR(v, tc.n) != tc.expected.(iec.UDINT) {
					t.Errorf("ROR failed")
				}
			case iec.LINT:
				if ROR(v, tc.n) != tc.expected.(iec.LINT) {
					t.Errorf("ROR failed")
				}
			case iec.INT:
				if ROR(v, tc.n) != tc.expected.(iec.INT) {
					t.Errorf("ROR failed")
				}
			default:
				t.Fatalf("unhandled type for ROR test case: %T", v)
			}
		})
	}
}
