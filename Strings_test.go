package royaljelly

import "testing"

func TestLEN(t *testing.T) {
	if LEN("ABCDE") != 5 {
		t.Error("LEN(\"ABCDE\") should be 5")
	}
	if LEN("") != 0 {
		t.Error("LEN(\"\") should be 0")
	}
}

func TestLEFT(t *testing.T) {
	if LEFT("ABCDEFG", 3) != "ABC" {
		t.Error("LEFT failed")
	}
	if LEFT("ABC", 5) != "ABC" {
		t.Error("LEFT with L > len(s) failed")
	}
	if LEFT("ABC", 0) != "" {
		t.Error("LEFT with L=0 failed")
	}
}

func TestRIGHT(t *testing.T) {
	if RIGHT("ABCDEFG", 3) != "EFG" {
		t.Error("RIGHT failed")
	}
	if RIGHT("ABC", 5) != "ABC" {
		t.Error("RIGHT with L > len(s) failed")
	}
	if RIGHT("ABC", 0) != "" {
		t.Error("RIGHT with L=0 failed")
	}
}

func TestMID(t *testing.T) {
	if MID("ABCDEFG", 3, 3) != "CDE" {
		t.Error("MID failed")
	}
	if MID("ABCDEFG", 10, 5) != "EFG" {
		t.Error("MID with L > remaining length failed")
	}
	if MID("ABCDEFG", 3, 8) != "" {
		t.Error("MID with P > len(s) failed")
	}
}

func TestCONCAT(t *testing.T) {
	if CONCAT("A", "B", "C") != "ABC" {
		t.Error("CONCAT failed")
	}
	if CONCAT() != "" {
		t.Error("CONCAT with no inputs failed")
	}
	if CONCAT("Hello") != "Hello" {
		t.Error("CONCAT with one input failed")
	}
}

func TestINSERT(t *testing.T) {
	testCases := []struct {
		name     string
		in1      STRING
		in2      STRING
		p        LINT
		expected STRING
	}{
		{"Insert in middle", "ABCD", "XYZ", 2, "ABXYZCD"},
		{"Insert at start (P=0)", "ABCD", "XYZ", 0, "XYZABCD"},
		{"Insert at end (P=len)", "ABCD", "XYZ", 4, "ABCDXYZ"},
		{"Insert past end (P>len)", "ABCD", "XYZ", 5, "ABCDXYZ"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := INSERT(tc.in1, tc.in2, tc.p)
			if result != tc.expected {
				t.Errorf("INSERT(%q, %q, %d) = %q; want %q", tc.in1, tc.in2, tc.p, result, tc.expected)
			}
		})
	}
}

func TestDELETE(t *testing.T) {
	if DELETE("ABCDEFG", 3, 3) != "ABFG" {
		t.Error("DELETE failed")
	}
	if DELETE("ABCDEFG", 10, 3) != "AB" {
		t.Error("DELETE with L > remaining length failed")
	}
	if DELETE("ABCDEFG", 2, 8) != "ABCDEFG" {
		t.Error("DELETE with P > len(s) should be no-op")
	}
}

func TestFIND(t *testing.T) {
	if FIND("ABCABC", "BCA") != 2 {
		t.Error("FIND failed to find substring")
	}
	if FIND("ABCABC", "XYZ") != 0 {
		t.Error("FIND should return 0 for not found")
	}
	if FIND("ABC", "A") != 1 {
		t.Error("FIND at start of string failed")
	}
	if FIND("", "") != 1 {
		t.Error("FIND of empty in empty should be 1")
	}
}

func TestREPLACE(t *testing.T) {
	testCases := []struct {
		name     string
		in1      STRING
		in2      STRING
		l        LINT
		p        LINT
		expected STRING
	}{
		{"Standard replace", "ABCDEFG", "XYZ", 3, 3, "ABXYZFG"},
		{"Replace at start", "ABCDEFG", "XYZ", 2, 1, "XYZCDEFG"},
		{"Replace at end", "ABCDEFG", "XYZ", 2, 6, "ABCDEXYZ"},
		{"Insert (L=0)", "ABCDEFG", "XYZ", 0, 3, "ABCXYZDEFG"},
		{"Delete (IN2 empty)", "ABCDEFG", "", 3, 3, "ABFG"},
		{"L exceeds string length", "ABCDEFG", "XYZ", 10, 5, "ABCDXYZ"},
		{"P exceeds string length", "ABC", "XYZ", 1, 4, "ABCXYZ"},
		{"P is 1", "ABCDEFG", "X", 2, 1, "XCDEFG"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := REPLACE(tc.in1, tc.in2, tc.l, tc.p)
			if result != tc.expected {
				t.Errorf("REPLACE(%q, %q, %d, %d) = %q; want %q", tc.in1, tc.in2, tc.l, tc.p, result, tc.expected)
			}
		})
	}

	t.Run("Panic on negative P", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for P < 0, but did not get one")
			}
		}()
		REPLACE("abc", "d", 1, -1)
	})
}
