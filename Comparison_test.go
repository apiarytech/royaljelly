package royaljelly

import (
	"testing"
	"time"
)

func TestGT(t *testing.T) {
	testCases := []struct {
		name     string
		inputs   []interface{}
		expected BOOL
	}{
		{"LINTs true", []interface{}{LINT(100), LINT(50), LINT(10)}, true},
		{"LINTs false", []interface{}{LINT(100), LINT(100), LINT(10)}, false},
		{"REALs true", []interface{}{REAL(10.5), REAL(5.5)}, true},
		{"REALs false", []interface{}{REAL(10.5), REAL(10.6)}, false},
		{"Mixed true", []interface{}{LREAL(100.0), INT(50), REAL(10.5)}, true},
		{"Mixed false", []interface{}{LREAL(100.0), INT(100), REAL(10.5)}, false},
		{"Strings true", []interface{}{STRING("z"), STRING("m"), STRING("a")}, true},
		{"Strings false", []interface{}{STRING("a"), STRING("z")}, false},
		{"TIME true", []interface{}{TIME(time.Hour), TIME(time.Minute)}, true},
		{"TIME false", []interface{}{TIME(time.Minute), TIME(time.Hour)}, false},
		{"BOOLs false", []interface{}{BOOL(true), BOOL(false), BOOL(true)}, false},
		{"Less than 2 inputs", []interface{}{LINT(10)}, false},
		{"Incompatible types", []interface{}{LINT(10), STRING("abc")}, false},
		{"Incompatible TIME and INT", []interface{}{TIME(time.Second), INT(1000)}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GT(tc.inputs)
			if result != tc.expected {
				t.Errorf("GT(%v) = %v; want %v", tc.inputs, result, tc.expected)
			}
		})
	}
}

func TestGE(t *testing.T) {
	testCases := []struct {
		name     string
		inputs   []interface{}
		expected BOOL
	}{
		{"LINTs true", []interface{}{LINT(100), LINT(50), LINT(10)}, true},
		{"LINTs with equal true", []interface{}{LINT(100), LINT(100), LINT(10)}, true},
		{"LINTs false", []interface{}{LINT(100), LINT(99), LINT(100)}, false},
		{"REALs true", []interface{}{REAL(10.5), REAL(5.5)}, true},
		{"REALs with equal true", []interface{}{REAL(10.5), REAL(10.5)}, true},
		{"Mixed true", []interface{}{LREAL(100.0), INT(100), REAL(10.5)}, true},
		{"Strings true", []interface{}{STRING("z"), STRING("m"), STRING("a")}, true},
		{"Strings with equal true", []interface{}{STRING("z"), STRING("z"), STRING("a")}, true},
		{"TIME true", []interface{}{TIME(time.Hour), TIME(time.Hour)}, true},
		{"Less than 2 inputs", []interface{}{LINT(10)}, false},
		{"Incompatible types", []interface{}{LINT(10), STRING("abc")}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GE(tc.inputs)
			if result != tc.expected {
				t.Errorf("GE(%v) = %v; want %v", tc.inputs, result, tc.expected)
			}
		})
	}
}

func TestEQ(t *testing.T) {
	testCases := []struct {
		name     string
		inputs   []interface{}
		expected BOOL
	}{
		{"LINTs true", []interface{}{LINT(50), LINT(50), LINT(50)}, true},
		{"LINTs false", []interface{}{LINT(50), LINT(50), LINT(51)}, false},
		{"REALs true", []interface{}{REAL(5.5), REAL(5.5)}, true},
		{"REALs false", []interface{}{REAL(5.5), REAL(5.6)}, false},
		{"Mixed true", []interface{}{LREAL(50.0), INT(50), REAL(50.0), DINT(50)}, true},
		{"Mixed false", []interface{}{LREAL(50.0), INT(51)}, false},
		{"Strings true", []interface{}{STRING("hello"), STRING("hello")}, true},
		{"Strings false", []interface{}{STRING("hello"), STRING("world")}, false},
		{"TIME true", []interface{}{TIME(time.Hour), TIME(60 * time.Minute)}, true},
		{"BOOLs true", []interface{}{BOOL(true), BOOL(true)}, true},
		{"BOOLs false", []interface{}{BOOL(true), BOOL(false)}, false},
		{"Less than 2 inputs", []interface{}{LINT(10)}, false},
		{"Incompatible types", []interface{}{LINT(10), STRING("abc")}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := EQ(tc.inputs)
			if result != tc.expected {
				t.Errorf("EQ(%v) = %v; want %v", tc.inputs, result, tc.expected)
			}
		})
	}
}

func TestLE(t *testing.T) {
	testCases := []struct {
		name     string
		inputs   []interface{}
		expected BOOL
	}{
		{"LINTs true", []interface{}{LINT(10), LINT(50), LINT(100)}, true},
		{"LINTs with equal true", []interface{}{LINT(10), LINT(50), LINT(50)}, true},
		{"LINTs false", []interface{}{LINT(10), LINT(50), LINT(49)}, false},
		{"REALs true", []interface{}{REAL(5.5), REAL(10.5)}, true},
		{"REALs with equal true", []interface{}{REAL(5.5), REAL(5.5)}, true},
		{"Mixed true", []interface{}{INT(10), REAL(50.0), LREAL(100.0)}, true},
		{"Strings true", []interface{}{STRING("a"), STRING("m"), STRING("z")}, true},
		{"Strings with equal true", []interface{}{STRING("a"), STRING("m"), STRING("m")}, true},
		{"TIME true", []interface{}{TIME(time.Minute), TIME(time.Hour)}, true},
		{"Less than 2 inputs", []interface{}{LINT(10)}, false},
		{"Incompatible types", []interface{}{LINT(10), STRING("abc")}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := LE(tc.inputs)
			if result != tc.expected {
				t.Errorf("LE(%v) = %v; want %v", tc.inputs, result, tc.expected)
			}
		})
	}
}

func TestLT(t *testing.T) {
	testCases := []struct {
		name     string
		inputs   []interface{}
		expected BOOL
	}{
		{"LINTs true", []interface{}{LINT(10), LINT(50), LINT(100)}, true},
		{"LINTs false", []interface{}{LINT(10), LINT(50), LINT(50)}, false},
		{"REALs true", []interface{}{REAL(5.5), REAL(10.5)}, true},
		{"REALs false", []interface{}{REAL(10.5), REAL(10.5)}, false},
		{"Mixed true", []interface{}{INT(10), REAL(50.0), LREAL(100.0)}, true},
		{"Mixed false", []interface{}{INT(10), REAL(100.0), LREAL(50.0)}, false},
		{"Strings true", []interface{}{STRING("a"), STRING("m"), STRING("z")}, true},
		{"Strings false", []interface{}{STRING("z"), STRING("a")}, false},
		{"TIME true", []interface{}{TIME(time.Minute), TIME(time.Hour)}, true},
		{"TIME false", []interface{}{TIME(time.Hour), TIME(time.Minute)}, false},
		{"BOOLs true", []interface{}{BOOL(false), BOOL(true)}, true},
		{"Less than 2 inputs", []interface{}{LINT(10)}, false},
		{"Incompatible types", []interface{}{LINT(10), STRING("abc")}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := LT(tc.inputs)
			if result != tc.expected {
				t.Errorf("LT(%v) = %v; want %v", tc.inputs, result, tc.expected)
			}
		})
	}
}

func TestNE(t *testing.T) {
	testCases := []struct {
		name     string
		inputs   []interface{}
		expected BOOL
	}{
		{"LINTs true", []interface{}{LINT(10), LINT(50)}, true},
		{"LINTs false", []interface{}{LINT(50), LINT(50)}, false},
		{"REALs true", []interface{}{REAL(5.5), REAL(10.5)}, true},
		{"REALs false", []interface{}{REAL(5.5), REAL(5.5)}, false},
		{"Mixed true", []interface{}{INT(10), REAL(50.0)}, true},
		{"Mixed false", []interface{}{INT(50), LREAL(50.0)}, false},
		{"Strings true", []interface{}{STRING("a"), STRING("z")}, true},
		{"Strings false", []interface{}{STRING("a"), STRING("a")}, false},
		{"TIME true", []interface{}{TIME(time.Minute), TIME(time.Hour)}, true},
		{"TIME false", []interface{}{TIME(time.Minute), TIME(60 * time.Second)}, false},
		{"BOOLs true", []interface{}{BOOL(false), BOOL(true)}, true},
		{"BOOLs false", []interface{}{BOOL(true), BOOL(true)}, false},
		{"Less than 2 inputs", []interface{}{LINT(10)}, false},
		{"More than 2 inputs", []interface{}{LINT(10), LINT(20), LINT(30)}, false},
		{"Incompatible types", []interface{}{LINT(10), STRING("abc")}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NE(tc.inputs)
			if result != tc.expected {
				t.Errorf("NE(%v) = %v; want %v", tc.inputs, result, tc.expected)
			}
		})
	}
}
