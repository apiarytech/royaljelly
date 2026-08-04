package comparison

import (
	"testing"
	"time"

	. "github.com/apiarytech/royaljelly/core"
)

func TestGT(t *testing.T) {
	testCases := []struct {
		name     string
		result   BOOL
		expected BOOL
	}{
		{"LINTs true", GT(LINT(100), LINT(50), LINT(10)), true},
		{"LINTs false", GT(LINT(100), LINT(100), LINT(10)), false},
		{"REALs true", GT(REAL(10.5), REAL(5.5)), true},
		{"REALs false", GT(REAL(10.5), REAL(10.6)), false},
		{"LREALs true", GT(LREAL(100.0), LREAL(50), LREAL(10.5)), true},
		{"LREALs false", GT(LREAL(100.0), LREAL(100), LREAL(10.5)), false},
		{"Strings true", GT(STRING("z"), STRING("m"), STRING("a")), true},
		{"Strings false", GT(STRING("a"), STRING("z")), false},
		{"TIME true", GT(TIME(time.Hour), TIME(time.Minute)), true},
		{"TIME false", GT(TIME(time.Minute), TIME(time.Hour)), false},
		{"Less than 2 inputs", GT(LINT(10)), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.result != tc.expected {
				t.Errorf("GT() = %v; want %v", tc.result, tc.expected)
			}
		})
	}
}

func TestGE(t *testing.T) {
	testCases := []struct {
		name     string
		result   BOOL
		expected BOOL
	}{
		{"LINTs true", GE(LINT(100), LINT(50), LINT(10)), true},
		{"LINTs with equal true", GE(LINT(100), LINT(100), LINT(10)), true},
		{"LINTs false", GE(LINT(100), LINT(99), LINT(100)), false},
		{"REALs true", GE(REAL(10.5), REAL(5.5)), true},
		{"REALs with equal true", GE(REAL(10.5), REAL(10.5)), true},
		{"LREALs true", GE(LREAL(100.0), LREAL(100), LREAL(10.5)), true},
		{"Strings true", GE(STRING("z"), STRING("m"), STRING("a")), true},
		{"Strings with equal true", GE(STRING("z"), STRING("z"), STRING("a")), true},
		{"TIME true", GE(TIME(time.Hour), TIME(time.Hour)), true},
		{"Less than 2 inputs", GE(LINT(10)), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.result != tc.expected {
				t.Errorf("GE() = %v; want %v", tc.result, tc.expected)
			}
		})
	}
}

func TestEQ(t *testing.T) {
	testCases := []struct {
		name     string
		result   BOOL
		expected BOOL
	}{
		{"LINTs true", EQ(LINT(50), LINT(50), LINT(50)), true},
		{"LINTs false", EQ(LINT(50), LINT(50), LINT(51)), false},
		{"REALs true", EQ(REAL(5.5), REAL(5.5)), true},
		{"REALs false", EQ(REAL(5.5), REAL(5.6)), false},
		{"LREALs true", EQ(LREAL(50.0), LREAL(50.0)), true},
		{"Strings true", EQ(STRING("hello"), STRING("hello")), true},
		{"Strings false", EQ(STRING("hello"), STRING("world")), false},
		{"TIME true", EQ(TIME(time.Hour), TIME(60*time.Minute)), true},
		{"BOOLs true", EQ(BOOL(true), BOOL(true)), true},
		{"BOOLs false", EQ(BOOL(true), BOOL(false)), false},
		{"Less than 2 inputs", EQ(LINT(10)), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.result != tc.expected {
				t.Errorf("EQ() = %v; want %v", tc.result, tc.expected)
			}
		})
	}
}

func TestLE(t *testing.T) {
	testCases := []struct {
		name     string
		result   BOOL
		expected BOOL
	}{
		{"LINTs true", LE(LINT(10), LINT(50), LINT(100)), true},
		{"LINTs with equal true", LE(LINT(10), LINT(50), LINT(50)), true},
		{"LINTs false", LE(LINT(10), LINT(50), LINT(49)), false},
		{"REALs true", LE(REAL(5.5), REAL(10.5)), true},
		{"REALs with equal true", LE(REAL(5.5), REAL(5.5)), true},
		{"LREALs true", LE(LREAL(10), LREAL(50.0), LREAL(100.0)), true},
		{"Strings true", LE(STRING("a"), STRING("m"), STRING("z")), true},
		{"Strings with equal true", LE(STRING("a"), STRING("m"), STRING("m")), true},
		{"TIME true", LE(TIME(time.Minute), TIME(time.Hour)), true},
		{"Less than 2 inputs", LE(LINT(10)), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.result != tc.expected {
				t.Errorf("LE() = %v; want %v", tc.result, tc.expected)
			}
		})
	}
}

func TestLT(t *testing.T) {
	testCases := []struct {
		name     string
		result   BOOL
		expected BOOL
	}{
		{"LINTs true", LT(LINT(10), LINT(50), LINT(100)), true},
		{"LINTs false", LT(LINT(10), LINT(50), LINT(50)), false},
		{"REALs true", LT(REAL(5.5), REAL(10.5)), true},
		{"REALs false", LT(REAL(10.5), REAL(10.5)), false},
		{"LREALs true", LT(LREAL(10), LREAL(50.0), LREAL(100.0)), true},
		{"LREALs false", LT(LREAL(10), LREAL(100.0), LREAL(50.0)), false},
		{"Strings true", LT(STRING("a"), STRING("m"), STRING("z")), true},
		{"Strings false", LT(STRING("z"), STRING("a")), false},
		{"TIME true", LT(TIME(time.Minute), TIME(time.Hour)), true},
		{"TIME false", LT(TIME(time.Hour), TIME(time.Minute)), false},
		{"Less than 2 inputs", LT(LINT(10)), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.result != tc.expected {
				t.Errorf("LT() = %v; want %v", tc.result, tc.expected)
			}
		})
	}
}

func TestNE(t *testing.T) {
	testCases := []struct {
		name     string
		result   BOOL
		expected BOOL
	}{
		{"LINTs true", NE(LINT(10), LINT(50)), true},
		{"LINTs false", NE(LINT(50), LINT(50)), false},
		{"REALs true", NE(REAL(5.5), REAL(10.5)), true},
		{"REALs false", NE(REAL(5.5), REAL(5.5)), false},
		{"LREALs false", NE(LREAL(50.0), LREAL(50.0)), false},
		{"Strings true", NE(STRING("a"), STRING("z")), true},
		{"Strings false", NE(STRING("a"), STRING("a")), false},
		{"TIME true", NE(TIME(time.Minute), TIME(time.Hour)), true},
		{"TIME false", NE(TIME(time.Minute), TIME(60*time.Second)), false},
		{"BOOLs true", NE(BOOL(false), BOOL(true)), true},
		{"BOOLs false", NE(BOOL(true), BOOL(true)), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.result != tc.expected {
				t.Errorf("NE() = %v; want %v", tc.result, tc.expected)
			}
		})
	}
}
