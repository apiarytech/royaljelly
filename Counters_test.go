package royaljelly

import "testing"

func TestCTU(t *testing.T) {
	fb := CTU{}
	fb.INIT()
	fb.PV = 3

	// Initial state
	fb.CTU()
	if fb.CV != 0 || fb.Q != false {
		t.Errorf("Initial state incorrect: CV=%d, Q=%v", fb.CV, fb.Q)
	}

	// Count up on rising edge
	fb.CU = true
	fb.CTU() // Edge detected
	if fb.CV != 1 {
		t.Errorf("CV should be 1 after first count, got %d", fb.CV)
	}

	// No change on high input
	fb.CTU()
	if fb.CV != 1 {
		t.Errorf("CV should not change on steady high input, got %d", fb.CV)
	}

	// Count up to PV
	fb.CU = false
	fb.CTU()
	fb.CU = true
	fb.CTU() // CV=2
	fb.CU = false
	fb.CTU()
	fb.CU = true
	fb.CTU() // CV=3

	if fb.CV != 3 || fb.Q != true {
		t.Errorf("State at PV incorrect: CV=%d, Q=%v", fb.CV, fb.Q)
	}

	// Does not count past PV
	fb.CU = false
	fb.CTU()
	fb.CU = true
	fb.CTU()
	if fb.CV != 3 {
		t.Errorf("CV should not exceed PV, got %d", fb.CV)
	}

	// Reset
	fb.R = true
	fb.CTU()
	if fb.CV != 0 || fb.Q != false {
		t.Errorf("Reset state incorrect: CV=%d, Q=%v", fb.CV, fb.Q)
	}
}

func TestCTD(t *testing.T) {
	fb := CTD{}
	fb.INIT()
	fb.PV = 3

	// Load
	fb.LD = true
	fb.CTD()
	if fb.CV != 3 || fb.Q != false {
		t.Errorf("Load state incorrect: CV=%d, Q=%v", fb.CV, fb.Q)
	}
	fb.LD = false
	fb.CTD()

	// Count down
	fb.CD = true
	fb.CTD() // Edge
	fb.CD = false
	fb.CTD()
	if fb.CV != 2 {
		t.Errorf("CV should be 2 after first count down, got %d", fb.CV)
	}

	// Count down to 0
	fb.CD = true
	fb.CTD() // CV=1
	fb.CD = false
	fb.CTD()
	fb.CD = true
	fb.CTD() // CV=0

	if fb.CV != 0 || fb.Q != true {
		t.Errorf("State at zero incorrect: CV=%d, Q=%v", fb.CV, fb.Q)
	}
}

func TestCTUD(t *testing.T) {
	t.Run("Count Up", func(t *testing.T) {
		fb := CTUD{}
		fb.INIT()
		fb.PV = 5

		fb.CU = true
		fb.CTUD() // Rising edge on CU
		if fb.CV != 1 {
			t.Errorf("CV should be 1 after first count up, got %d", fb.CV)
		}
		fb.CU = false
		fb.CTUD() // Reset for next edge
	})

	t.Run("Count Down", func(t *testing.T) {
		fb := CTUD{}
		fb.INIT()
		fb.PV = 5
		fb.CV = 3 // Set initial value

		fb.CD = true
		fb.CTUD() // Rising edge on CD
		if fb.CV != 2 {
			t.Errorf("CV should be 2 after count down, got %d", fb.CV)
		}
		fb.CD = false
		fb.CTUD() // Reset for next edge
	})

	t.Run("Reset and Load", func(t *testing.T) {
		fb := CTUD{}
		fb.INIT()
		fb.PV = 10
		fb.CV = 5

		// Test Load
		fb.LD = true
		fb.CTUD()
		if fb.CV != 10 {
			t.Errorf("CV should be PV after load, got %d", fb.CV)
		}
		fb.LD = false
		fb.CTUD()

		// Test Reset
		fb.R = true
		fb.CTUD()
		if fb.CV != 0 {
			t.Errorf("CV should be 0 after reset, got %d", fb.CV)
		}
	})

	t.Run("Simultaneous CU and CD", func(t *testing.T) {
		fb := CTUD{}
		fb.INIT()
		fb.PV = 10
		fb.CV = 5 // Start at a known value

		// First scan, inputs are low. This sets the edge detectors' memory.
		fb.CTUD()

		// Second scan, both inputs go high simultaneously.
		fb.CU = true
		fb.CD = true
		fb.CTUD()

		if fb.CV != 5 {
			t.Errorf("CV should not change when CU and CD have simultaneous rising edges, got %d, want 5", fb.CV)
		}
	})
}
