package royaljelly

import "testing"

func TestR_TRIG(t *testing.T) {
	testCases := []struct {
		name        string
		initialMem  BOOL
		clk         BOOL
		expectedQ   BOOL
		expectedMem BOOL
	}{
		{"Rising edge", false, true, true, true},
		{"No edge (high)", true, true, false, true},
		{"No edge (low)", false, false, false, false},
		{"Falling edge", true, false, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize the function block
			fb := R_TRIG{}
			fb.INIT()
			fb.mem = tc.initialMem
			fb.CLK = tc.clk

			// Execute the logic
			fb.R_TRIG()

			// Check the output Q
			if fb.Q != tc.expectedQ {
				t.Errorf("Q mismatch: got %v, want %v", fb.Q, tc.expectedQ)
			}

			// Check the internal memory
			if fb.mem != tc.expectedMem {
				t.Errorf("mem mismatch: got %v, want %v", fb.mem, tc.expectedMem)
			}
		})
	}
}

func TestF_TRIG(t *testing.T) {
	testCases := []struct {
		name        string
		initialMem  BOOL
		clk         BOOL
		expectedQ   BOOL
		expectedMem BOOL
	}{
		{"Falling edge", true, false, true, false},
		{"No edge (low)", false, false, false, false},
		{"No edge (high)", true, true, false, true},
		{"Rising edge", false, true, false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize the function block
			fb := F_TRIG{}
			fb.INIT()
			fb.mem = tc.initialMem
			fb.CLK = tc.clk

			// Execute the logic
			fb.F_TRIG()

			// Check the output Q
			if fb.Q != tc.expectedQ {
				t.Errorf("Q mismatch: got %v, want %v", fb.Q, tc.expectedQ)
			}

			// Check the internal memory
			if fb.mem != tc.expectedMem {
				t.Errorf("mem mismatch: got %v, want %v", fb.mem, tc.expectedMem)
			}
		})
	}
}
