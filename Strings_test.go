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
	res, err := LEFT("ABCDEFG", 3)
	if err != nil || res != "ABC" {
		t.Error("LEFT failed")
	}
	res, err = LEFT("ABC", 5)
	if err != nil || res != "ABC" {
		t.Error("LEFT with L > len(s) failed")
	}
	res, err = LEFT("ABC", 0)
	if err != nil || res != "" {
		t.Error("LEFT with L=0 failed")
	}
	_, err = LEFT("ABC", -1)
	if err == nil {
		t.Error("LEFT with negative L should have returned an error")
	}
}

func TestRIGHT(t *testing.T) {
	res, err := RIGHT("ABCDEFG", 3)
	if err != nil || res != "EFG" {
		t.Error("RIGHT failed")
	}
	res, err = RIGHT("ABC", 5)
	if err != nil || res != "ABC" {
		t.Error("RIGHT with L > len(s) failed")
	}
	res, err = RIGHT("ABC", 0)
	if err != nil || res != "" {
		t.Error("RIGHT with L=0 failed")
	}
	_, err = RIGHT("ABC", -1)
	if err == nil {
		t.Error("RIGHT with negative L should have returned an error")
	}
}

func TestMID(t *testing.T) {
	res, err := MID("ABCDEFG", 3, 3)
	if err != nil || res != "CDE" {
		t.Errorf("MID failed: %v", err)
	}
	res, err = MID("ABCDEFG", 10, 5)
	if err != nil || res != "EFG" {
		t.Errorf("MID with L > remaining length failed: %v", err)
	}
	res, err = MID("ABCDEFG", 3, 8)
	if err != nil || res != "" {
		t.Errorf("MID with P > len(s) failed: %v", err)
	}
	_, err = MID("ABC", 1, 0)
	if err == nil {
		t.Error("MID with P < 1 should have returned an error")
	}
	_, err = MID("ABC", -1, 1)
	if err == nil {
		t.Error("MID with negative L should have returned an error")
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
		hasError bool
	}{
		{"Insert in middle", "ABCD", "XYZ", 2, "ABXYZCD", false},
		{"Insert at start (P=0)", "ABCD", "XYZ", 0, "XYZABCD", false},
		{"Insert at end (P=len)", "ABCD", "XYZ", 4, "ABCDXYZ", false},
		{"Insert past end (P>len)", "ABCD", "XYZ", 5, "ABCDXYZ", false},
		{"Insert with negative P", "ABCD", "XYZ", -1, "", true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := INSERT(tc.in1, tc.in2, tc.p)
			if tc.hasError {
				if err == nil {
					t.Errorf("INSERT expected an error but got none")
				}
			} else if err != nil {
				t.Errorf("INSERT returned an unexpected error: %v", err)
			} else if result != tc.expected {
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

func TestREPEAT(t *testing.T) {
	t.Run("Valid repeat", func(t *testing.T) {
		res, err := REPEAT("A", 5)
		if err != nil || res != "AAAAA" {
			t.Errorf("REPEAT('A', 5) failed. Got %q, err: %v", res, err)
		}
	})

	t.Run("Repeat zero times", func(t *testing.T) {
		res, err := REPEAT("ABC", 0)
		if err != nil || res != "" {
			t.Errorf("REPEAT('ABC', 0) failed. Got %q, err: %v", res, err)
		}
	})

	t.Run("Repeat with negative count", func(t *testing.T) {
		_, err := REPEAT("A", -1)
		if err == nil {
			t.Error("REPEAT with negative count should have returned an error")
		}
	})
}

func TestREPLACE(t *testing.T) {
	testCases := []struct {
		name     string
		in1      STRING
		in2      STRING
		l        LINT
		p        LINT
		expected STRING
		hasError bool
	}{
		{"Standard replace in middle", "ABCDEFG", "XYZ", 3, 3, "ABXYZFG", false},
		{"Replace at start (P=1)", "ABCDEFG", "XYZ", 2, 1, "XYZCDEFG", false},
		{"Replace at start (P=0)", "ABCDEFG", "XYZ", 2, 0, "XYZCDEFG", false}, // P=0 should act like P=1
		{"Replace at end", "ABCDEFG", "XYZ", 2, 6, "ABCDEXYZ", false},
		{"Insert behavior (L=0)", "ABCD", "XYZ", 0, 2, "ABXYZCD", false},      // Should behave like INSERT
		{"Insert at start (L=0, P=0)", "ABCD", "XYZ", 0, 0, "XYZABCD", false}, // Should behave like INSERT
		{"Delete behavior (IN2 empty)", "ABCDEFG", "", 3, 3, "ABFG", false},
		{"L exceeds string length", "ABCDEFG", "XYZ", 10, 5, "ABCDXYZ", false},
		{"P exceeds string length", "ABC", "XYZ", 1, 4, "ABCXYZ", false},
		{"Negative L", "abc", "d", -1, 1, "", true},
		{"Negative P", "abc", "d", 1, -1, "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := REPLACE(tc.in1, tc.in2, tc.l, tc.p)
			if tc.hasError {
				if err == nil {
					t.Errorf("REPLACE expected an error but got none")
				}
			} else if err != nil {
				t.Errorf("REPLACE returned an unexpected error: %v", err)
			} else if result != tc.expected {
				t.Errorf("REPLACE(%q, %q, %d, %d) = %q; want %q", tc.in1, tc.in2, tc.l, tc.p, result, tc.expected)
			}
		})
	}
}
