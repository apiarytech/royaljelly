package royaljelly

import (
	"fmt"
	"testing"
	"time"
)

func TestSEL(t *testing.T) {
	testCases := []struct {
		name     string
		g        BOOL
		in0      interface{}
		in1      interface{}
		expected interface{}
	}{
		{"Select IN0 (false)", false, LINT(10), LINT(20), LINT(10)},
		{"Select IN1 (true)", true, LINT(10), LINT(20), LINT(20)},
		{"Select REAL", true, REAL(1.5), REAL(2.5), REAL(2.5)},
		{"Select STRING", false, STRING("a"), STRING("b"), STRING("a")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SEL(tc.g, tc.in0, tc.in1)
			if result != tc.expected {
				t.Errorf("SEL(%v, %v, %v) = %v; want %v", tc.g, tc.in0, tc.in1, result, tc.expected)
			}
		})
	}
}

func TestMAX(t *testing.T) {
	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectError bool
	}{
		{"LINTs", []interface{}{LINT(10), LINT(50), LINT(20)}, LINT(50), false},
		{"REALs", []interface{}{REAL(10.5), REAL(10.6), REAL(10.1)}, REAL(10.6), false},
		{"Mixed Int/Float", []interface{}{INT(100), REAL(100.1)}, REAL(100.1), false},
		{"Strings", []interface{}{STRING("apple"), STRING("orange"), STRING("banana")}, STRING("orange"), false},
		{"TIME", []interface{}{TIME(time.Second), TIME(time.Minute)}, TIME(time.Minute), false},
		{"Not enough inputs", []interface{}{LINT(10)}, nil, true},
		{"Incompatible types", []interface{}{LINT(10), STRING("apple")}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MAX(tc.inputs)
			if tc.expectError {
				if err == nil {
					t.Errorf("MAX(%v) did not return an error; expected error", tc.inputs)
				}
			} else {
				if err != nil {
					t.Errorf("MAX(%v) returned an unexpected error: %v", tc.inputs, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("MAX(%v) = %v; want %v", tc.inputs, result, tc.expected)
				}
			}
		})
	}
}

func TestMIN(t *testing.T) {
	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectError bool
	}{
		{"LINTs", []interface{}{LINT(10), LINT(50), LINT(20)}, LINT(10), false},
		{"REALs", []interface{}{REAL(10.5), REAL(10.6), REAL(10.1)}, REAL(10.1), false},
		{"Mixed Int/Float", []interface{}{INT(100), REAL(99.9)}, REAL(99.9), false},
		{"Strings", []interface{}{STRING("apple"), STRING("orange"), STRING("banana")}, STRING("apple"), false},
		{"TIME", []interface{}{TIME(time.Second), TIME(time.Minute)}, TIME(time.Second), false},
		{"Not enough inputs", []interface{}{LINT(10)}, nil, true},
		{"Incompatible types", []interface{}{LINT(10), STRING("apple")}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MIN(tc.inputs)
			if tc.expectError {
				if err == nil {
					t.Errorf("MIN(%v) did not return an error; expected error", tc.inputs)
				}
			} else {
				if err != nil {
					t.Errorf("MIN(%v) returned an unexpected error: %v", tc.inputs, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("MIN(%v) = %v; want %v", tc.inputs, result, tc.expected)
				}
			}
		})
	}
}

func TestLIMIT(t *testing.T) {
	testCases := []struct {
		name        string
		mn          interface{}
		in          interface{}
		mx          interface{}
		expected    interface{}
		expectError bool
	}{
		{"Within limits", LINT(10), LINT(50), LINT(100), LINT(50), false},
		{"Below minimum", LINT(10), LINT(5), LINT(100), LINT(10), false},
		{"Above maximum", LINT(10), LINT(150), LINT(100), LINT(100), false},
		{"REALs within limits", REAL(10.0), REAL(50.5), REAL(100.0), REAL(50.5), false},
		{"Strings within limits", STRING("a"), STRING("b"), STRING("c"), STRING("b"), false},
		{"compatible IN and MN", LINT(10), STRING("50"), LINT(100), STRING("50"), false},
		{"compatible IN and MX", LINT(10), LINT(50), STRING("100"), LINT(50), false},
		{"compatible MN and MX", LINT(10), LINT(50), REAL(100.0), LINT(50), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := LIMIT(tc.mn, tc.in, tc.mx)
			if tc.expectError {
				if err == nil {
					t.Errorf("LIMIT(%v, %v, %v) did not return an error; expected error", tc.mn, tc.in, tc.mx)
				}
			} else {
				if err != nil {
					t.Errorf("LIMIT(%v, %v, %v) returned an unexpected error: %v", tc.mn, tc.in, tc.mx, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("LIMIT(%v, %v, %v) = %v; want %v", tc.mn, tc.in, tc.mx, result, tc.expected)
				}
			}
		})
	}
}

func TestMUX(t *testing.T) {
	options := []interface{}{STRING("a"), STRING("b"), STRING("c"), STRING("d")}

	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectError bool
	}{
		{"Select 0", []interface{}{INT(0), options[0], options[1], options[2]}, options[0], false},
		{"Select 2", []interface{}{INT(2), options[0], options[1], options[2], options[3]}, options[2], false},
		{"Select with LINT", []interface{}{LINT(1), REAL(10.0), REAL(20.0)}, REAL(20.0), false},
		{"Selector out of bounds (negative)", []interface{}{INT(-1), options[0]}, nil, true},
		{"Selector out of bounds (too high)", []interface{}{INT(2), options[0], options[1]}, nil, true},
		{"Not enough inputs", []interface{}{INT(0)}, nil, true},
		{"Invalid selector type", []interface{}{STRING("abc"), options[0]}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MUX(tc.inputs)
			if tc.expectError {
				if err == nil {
					t.Errorf("MUX(%v) did not return an error; expected error", tc.inputs)
				}
			} else {
				if err != nil {
					t.Errorf("MUX(%v) returned an unexpected error: %v", tc.inputs, err)
				}
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("MUX(%v) = %v; want %v", tc.inputs, result, tc.expected)
				}
			}
		})
	}
}
