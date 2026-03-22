package royaljelly

import "testing"

func TestSR_FB(t *testing.T) {
	testCases := []struct {
		name       string
		s1         BOOL
		r          BOOL
		initialQ1  BOOL
		expectedQ1 BOOL
	}{
		{"Set", true, false, false, true},
		{"Hold after Set", false, false, true, true},
		{"Reset", false, true, true, false},
		{"Hold after Reset", false, false, false, false},
		{"Dominant Set (S1=true, R=true)", true, true, false, true},
		{"Dominant Set (S1=true, R=true, Q1 was true)", true, true, true, true},
		{"No change (S1=false, R=false)", false, false, true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fb := SR_FB{}
			fb.INIT()
			fb.Q1 = tc.initialQ1
			fb.S1 = tc.s1
			fb.R = tc.r

			fb.SR()

			if fb.Q1 != tc.expectedQ1 {
				t.Errorf("SR(S1=%v, R=%v, initialQ1=%v) -> Q1=%v; want %v", tc.s1, tc.r, tc.initialQ1, fb.Q1, tc.expectedQ1)
			}
		})
	}

	t.Run("Disabled", func(t *testing.T) {
		fb := SR_FB{}
		fb.INIT()
		fb.Q1 = false
		fb.S1 = true  // Should set Q1
		fb.EN = false // But FB is disabled

		fb.SR()

		if fb.Q1 != false {
			t.Error("Q1 should not change when EN is false")
		}
		if fb.ENO != false {
			t.Error("ENO should be false when EN is false")
		}
	})
}

func TestRS_FB(t *testing.T) {
	testCases := []struct {
		name       string
		s          BOOL
		r1         BOOL
		initialQ1  BOOL
		expectedQ1 BOOL
	}{
		{"Set", true, false, false, true},
		{"Hold after Set", false, false, true, true},
		{"Reset", false, true, true, false},
		{"Hold after Reset", false, false, false, false},
		{"Dominant Reset (S=true, R1=true)", true, true, true, false},
		{"Dominant Reset (S=true, R1=true, Q1 was false)", true, true, false, false},
		{"No change (S=false, R1=false)", false, false, true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fb := RS_FB{}
			fb.INIT()
			fb.Q1 = tc.initialQ1
			fb.S = tc.s
			fb.R1 = tc.r1

			fb.RS()

			if fb.Q1 != tc.expectedQ1 {
				t.Errorf("RS(S=%v, R1=%v, initialQ1=%v) -> Q1=%v; want %v", tc.s, tc.r1, tc.initialQ1, fb.Q1, tc.expectedQ1)
			}
		})
	}

	t.Run("Disabled", func(t *testing.T) {
		fb := RS_FB{}
		fb.INIT()
		fb.Q1 = true
		fb.R1 = true  // Should reset Q1
		fb.EN = false // But FB is disabled

		fb.RS()

		if fb.Q1 != true {
			t.Error("Q1 should not change when EN is false")
		}
		if fb.ENO != false {
			t.Error("ENO should be false when EN is false")
		}
	})
}
