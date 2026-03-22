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
		{"Select IN0 with LINT", false, LINT(10), LINT(20), LINT(10)},
		{"Select IN1 with LINT", true, LINT(10), LINT(20), LINT(20)},
		{"Select IN0 with REAL", false, REAL(1.1), REAL(2.2), REAL(1.1)},
		{"Select IN1 with REAL", true, REAL(1.1), REAL(2.2), REAL(2.2)},
		{"Select IN0 with STRING", false, STRING("A"), STRING("B"), STRING("A")},
		{"Select IN1 with STRING", true, STRING("A"), STRING("B"), STRING("B")},
		{"Select IN0 with different types", false, LINT(100), REAL(20.5), LINT(100)},
		{"Select IN1 with different types", true, LINT(100), REAL(20.5), REAL(20.5)},
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
		expectPanic bool
	}{
		{"LINTs", []interface{}{LINT(10), LINT(50), LINT(20)}, LINT(50), false},
		{"REALs", []interface{}{REAL(10.5), REAL(50.1), REAL(20.2)}, REAL(50.1), false},
		{"Mixed Int/Float", []interface{}{REAL(10.5), LINT(60), INT(20)}, LREAL(60.0), false},
		{"Strings", []interface{}{STRING("apple"), STRING("orange"), STRING("banana")}, STRING("orange"), false},
		{"TIME", []interface{}{TIME(time.Second), TIME(time.Minute), TIME(time.Millisecond)}, TIME(time.Minute), false},
		{"Single element", []interface{}{DINT(42)}, DINT(42), false},
		{"Empty slice", []interface{}{}, nil, false},
		{"Incompatible types", []interface{}{LINT(10), STRING("abc")}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.expectPanic {
					if r == nil {
						t.Errorf("MAX(%v) did not panic; expected panic", tc.inputs)
					}
				} else if r != nil {
					t.Errorf("MAX(%v) panicked; got %v", tc.inputs, r)
				}
			}()

			result := MAX(tc.inputs)
			if !tc.expectPanic {
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
		expectPanic bool
	}{
		{"LINTs", []interface{}{LINT(10), LINT(50), LINT(20)}, LINT(10), false},
		{"REALs", []interface{}{REAL(10.5), REAL(50.1), REAL(20.2)}, REAL(10.5), false},
		{"Mixed Int/Float", []interface{}{REAL(10.5), LINT(60), INT(5)}, LREAL(5.0), false},
		{"Strings", []interface{}{STRING("apple"), STRING("orange"), STRING("banana")}, STRING("apple"), false},
		{"TIME", []interface{}{TIME(time.Second), TIME(time.Minute), TIME(time.Millisecond)}, TIME(time.Millisecond), false},
		{"Single element", []interface{}{DINT(42)}, DINT(42), false},
		{"Empty slice", []interface{}{}, nil, false},
		{"Incompatible types", []interface{}{LINT(10), STRING("abc")}, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.expectPanic {
					if r == nil {
						t.Errorf("MIN(%v) did not panic; expected panic", tc.inputs)
					}
				} else if r != nil {
					t.Errorf("MIN(%v) panicked; got %v", tc.inputs, r)
				}
			}()

			result := MIN(tc.inputs)
			if !tc.expectPanic {
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("MIN(%v) = %v; want %v", tc.inputs, result, tc.expected)
				}
			}
		})
	}
}

func TestLIMIT(t *testing.T) {
	testCases := []struct {
		name     string
		mn       interface{}
		in       interface{}
		mx       interface{}
		expected interface{}
	}{
		{"LINT below min", LINT(10), LINT(5), LINT(20), LINT(10)},
		{"LINT above max", LINT(10), LINT(25), LINT(20), LINT(20)},
		{"LINT within range", LINT(10), LINT(15), LINT(20), LINT(15)},
		{"REAL below min", REAL(10.0), REAL(5.5), REAL(20.0), REAL(10.0)},
		{"REAL above max", REAL(10.0), REAL(25.5), REAL(20.0), REAL(20.0)},
		{"REAL within range", REAL(10.0), REAL(15.5), REAL(20.0), REAL(15.5)},
		{"Mixed types within range", LINT(10), REAL(15.0), LINT(20), REAL(15.0)},
		{"Mixed types below min", LINT(10), INT(5), REAL(20.0), LINT(10)},
		{"String below min", STRING("b"), STRING("a"), STRING("d"), STRING("b")},
		{"String above max", STRING("b"), STRING("e"), STRING("d"), STRING("d")},
		{"String within range", STRING("b"), STRING("c"), STRING("d"), STRING("c")},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := LIMIT(tc.mn, tc.in, tc.mx)
			if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
				t.Errorf("LIMIT(%v, %v, %v) = %v; want %v", tc.mn, tc.in, tc.mx, result, tc.expected)
			}
		})
	}
}

func TestMUX(t *testing.T) {
	options := []interface{}{
		STRING("apple"),
		STRING("banana"),
		STRING("cherry"),
		STRING("date"),
	}

	testCases := []struct {
		name        string
		inputs      []interface{}
		expected    interface{}
		expectPanic bool
	}{
		{"Select first", append([]interface{}{INT(0)}, options...), STRING("apple"), false},
		{"Select middle", append([]interface{}{LINT(2)}, options...), STRING("cherry"), false},
		{"Select last", append([]interface{}{DINT(3)}, options...), STRING("date"), false},
		{"K out of bounds (negative)", append([]interface{}{INT(-1)}, options...), nil, false},
		{"K out of bounds (too large)", append([]interface{}{INT(4)}, options...), nil, false},
		{"Empty inputs", []interface{}{}, nil, false},
		{"Only K provided", []interface{}{INT(0)}, nil, false},
		{"K is not integer", append([]interface{}{REAL(1.0)}, options...), STRING("banana"), false},
		{"K is string", append([]interface{}{STRING("1")}, options...), STRING("banana"), false},
		{"K is invalid string", append([]interface{}{STRING("abc")}, options...), nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.expectPanic {
					if r == nil {
						t.Errorf("MUX(%v) did not panic; expected panic", tc.inputs)
					}
				} else if r != nil {
					t.Errorf("MUX(%v) panicked; got %v", tc.inputs, r)
				}
			}()

			result := MUX(tc.inputs)
			if !tc.expectPanic {
				if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tc.expected) {
					t.Errorf("MUX(%v) = %v; want %v", tc.inputs, result, tc.expected)
				}
			}
		})
	}
}
