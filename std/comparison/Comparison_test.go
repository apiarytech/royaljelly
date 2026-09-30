package comparison

import (
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

func TestGT(t *testing.T) {
	testCases := []struct {
		name     string
		result   iec.BOOL
		expected iec.BOOL
	}{
		{"LINTs true", GT(iec.LINT(100), iec.LINT(50), iec.LINT(10)), true},
		{"LINTs false", GT(iec.LINT(100), iec.LINT(100), iec.LINT(10)), false},
		{"REALs true", GT(iec.REAL(10.5), iec.REAL(5.5)), true},
		{"REALs false", GT(iec.REAL(10.5), iec.REAL(10.6)), false},
		{"LREALs true", GT(iec.LREAL(100.0), iec.LREAL(50), iec.LREAL(10.5)), true},
		{"LREALs false", GT(iec.LREAL(100.0), iec.LREAL(100), iec.LREAL(10.5)), false},
		{"Strings true", GT(iec.STRING("z"), iec.STRING("m"), iec.STRING("a")), true},
		{"Strings false", GT(iec.STRING("a"), iec.STRING("z")), false},
		{"TIME true", GT(iec.TIME(time.Hour), iec.TIME(time.Minute)), true},
		{"TIME false", GT(iec.TIME(time.Minute), iec.TIME(time.Hour)), false},
		{"Less than 2 inputs", GT(iec.LINT(10)), false},
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
		result   iec.BOOL
		expected iec.BOOL
	}{
		{"LINTs true", GE(iec.LINT(100), iec.LINT(50), iec.LINT(10)), true},
		{"LINTs with equal true", GE(iec.LINT(100), iec.LINT(100), iec.LINT(10)), true},
		{"LINTs false", GE(iec.LINT(100), iec.LINT(99), iec.LINT(100)), false},
		{"REALs true", GE(iec.REAL(10.5), iec.REAL(5.5)), true},
		{"REALs with equal true", GE(iec.REAL(10.5), iec.REAL(10.5)), true},
		{"LREALs true", GE(iec.LREAL(100.0), iec.LREAL(100), iec.LREAL(10.5)), true},
		{"Strings true", GE(iec.STRING("z"), iec.STRING("m"), iec.STRING("a")), true},
		{"Strings with equal true", GE(iec.STRING("z"), iec.STRING("z"), iec.STRING("a")), true},
		{"TIME true", GE(iec.TIME(time.Hour), iec.TIME(time.Hour)), true},
		{"Less than 2 inputs", GE(iec.LINT(10)), false},
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
		result   iec.BOOL
		expected iec.BOOL
	}{
		{"LINTs true", EQ(iec.LINT(50), iec.LINT(50), iec.LINT(50)), true},
		{"LINTs false", EQ(iec.LINT(50), iec.LINT(50), iec.LINT(51)), false},
		{"REALs true", EQ(iec.REAL(5.5), iec.REAL(5.5)), true},
		{"REALs false", EQ(iec.REAL(5.5), iec.REAL(5.6)), false},
		{"LREALs true", EQ(iec.LREAL(50.0), iec.LREAL(50.0)), true},
		{"Strings true", EQ(iec.STRING("hello"), iec.STRING("hello")), true},
		{"Strings false", EQ(iec.STRING("hello"), iec.STRING("world")), false},
		{"TIME true", EQ(iec.TIME(time.Hour), iec.TIME(60*time.Minute)), true},
		{"BOOLs true", EQ(iec.BOOL(true), iec.BOOL(true)), true},
		{"BOOLs false", EQ(iec.BOOL(true), iec.BOOL(false)), false},
		{"Less than 2 inputs", EQ(iec.LINT(10)), false},
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
		result   iec.BOOL
		expected iec.BOOL
	}{
		{"LINTs true", LE(iec.LINT(10), iec.LINT(50), iec.LINT(100)), true},
		{"LINTs with equal true", LE(iec.LINT(10), iec.LINT(50), iec.LINT(50)), true},
		{"LINTs false", LE(iec.LINT(10), iec.LINT(50), iec.LINT(49)), false},
		{"REALs true", LE(iec.REAL(5.5), iec.REAL(10.5)), true},
		{"REALs with equal true", LE(iec.REAL(5.5), iec.REAL(5.5)), true},
		{"LREALs true", LE(iec.LREAL(10), iec.LREAL(50.0), iec.LREAL(100.0)), true},
		{"Strings true", LE(iec.STRING("a"), iec.STRING("m"), iec.STRING("z")), true},
		{"Strings with equal true", LE(iec.STRING("a"), iec.STRING("m"), iec.STRING("m")), true},
		{"TIME true", LE(iec.TIME(time.Minute), iec.TIME(time.Hour)), true},
		{"Less than 2 inputs", LE(iec.LINT(10)), false},
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
		result   iec.BOOL
		expected iec.BOOL
	}{
		{"LINTs true", LT(iec.LINT(10), iec.LINT(50), iec.LINT(100)), true},
		{"LINTs false", LT(iec.LINT(10), iec.LINT(50), iec.LINT(50)), false},
		{"REALs true", LT(iec.REAL(5.5), iec.REAL(10.5)), true},
		{"REALs false", LT(iec.REAL(10.5), iec.REAL(10.5)), false},
		{"LREALs true", LT(iec.LREAL(10), iec.LREAL(50.0), iec.LREAL(100.0)), true},
		{"LREALs false", LT(iec.LREAL(10), iec.LREAL(100.0), iec.LREAL(50.0)), false},
		{"Strings true", LT(iec.STRING("a"), iec.STRING("m"), iec.STRING("z")), true},
		{"Strings false", LT(iec.STRING("z"), iec.STRING("a")), false},
		{"TIME true", LT(iec.TIME(time.Minute), iec.TIME(time.Hour)), true},
		{"TIME false", LT(iec.TIME(time.Hour), iec.TIME(time.Minute)), false},
		{"Less than 2 inputs", LT(iec.LINT(10)), false},
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
		result   iec.BOOL
		expected iec.BOOL
	}{
		{"LINTs true", NE(iec.LINT(10), iec.LINT(50)), true},
		{"LINTs false", NE(iec.LINT(50), iec.LINT(50)), false},
		{"REALs true", NE(iec.REAL(5.5), iec.REAL(10.5)), true},
		{"REALs false", NE(iec.REAL(5.5), iec.REAL(5.5)), false},
		{"LREALs false", NE(iec.LREAL(50.0), iec.LREAL(50.0)), false},
		{"Strings true", NE(iec.STRING("a"), iec.STRING("z")), true},
		{"Strings false", NE(iec.STRING("a"), iec.STRING("a")), false},
		{"TIME true", NE(iec.TIME(time.Minute), iec.TIME(time.Hour)), true},
		{"TIME false", NE(iec.TIME(time.Minute), iec.TIME(60*time.Second)), false},
		{"BOOLs true", NE(iec.BOOL(false), iec.BOOL(true)), true},
		{"BOOLs false", NE(iec.BOOL(true), iec.BOOL(true)), false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.result != tc.expected {
				t.Errorf("NE() = %v; want %v", tc.result, tc.expected)
			}
		})
	}
}
