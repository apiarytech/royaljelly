package counters

import (
	"testing"
)

func TestCTU(t *testing.T) {
	fb := CTU{}
	fb.INIT()
	fb.PV = 3
	fb.EN = true

	// Initial state
	fb.Execute()
	if fb.CV != 0 || fb.Q != false {
		t.Errorf("Initial state incorrect: CV=%d, Q=%v", fb.CV, fb.Q)
	}

	// Count up on rising edge
	fb.CU = true
	fb.Execute()
	if fb.CV != 1 {
		t.Errorf("CV should be 1 after first count, got %d", fb.CV)
	}

	// No change on high input
	fb.Execute() // CU is still true, no new edge
	if fb.CV != 1 {
		t.Errorf("CV should not change on steady high input, got %d", fb.CV)
	}

	// Count up to PV
	fb.CU = false
	fb.Execute() // Set low for next edge
	fb.CU = true
	fb.Execute()
	fb.CU = false
	fb.Execute() // Set low for next edge
	fb.CU = true
	fb.Execute()

	if fb.CV != 3 || fb.Q != true {
		t.Errorf("State at PV incorrect: CV=%d, Q=%v", fb.CV, fb.Q)
	}

	// Should now count past PV
	fb.CU = false
	fb.Execute() // Set low for next edge
	fb.CU = true
	fb.Execute()
	if fb.CV != 4 {
		t.Errorf("CV should exceed PV, got %d", fb.CV)
	}

	// Reset
	fb.R = true
	fb.Execute()
	if fb.CV != 0 || fb.Q != false {
		t.Errorf("Reset state incorrect: CV=%d, Q=%v", fb.CV, fb.Q)
	}
}

func TestCTD(t *testing.T) {
	fb := CTD{}
	fb.INIT()
	fb.PV = 3
	fb.EN = true

	// Load
	fb.LD = true
	fb.Execute()
	if fb.CV != 3 || fb.Q != false {
		t.Errorf("Load state incorrect: CV=%d, Q=%v", fb.CV, fb.Q)
	}
	fb.LD = false
	fb.Execute()

	// Count down
	fb.CD = true
	fb.Execute()
	fb.CD = false
	fb.Execute()
	if fb.CV != 2 {
		t.Errorf("CV should be 2 after first count down, got %d", fb.CV)
	}

	// Count down to 0
	fb.CD = true
	fb.Execute()
	fb.CD = false
	fb.Execute()
	fb.CD = true
	fb.Execute()

	if fb.CV != 0 || fb.Q != true {
		t.Errorf("State at zero incorrect: CV=%d, Q=%v", fb.CV, fb.Q)
	}

	// Should now count below 0
	fb.CD = false
	fb.Execute() // Set low for next edge
	fb.CD = true
	fb.Execute()

	if fb.CV != -1 || fb.Q != true {
		// Q should remain true for CV <= 0
		t.Errorf("State below zero incorrect: CV=%d, Q=%v", fb.CV, fb.Q)
	}
}

func TestCTUD(t *testing.T) {
	t.Run("Count Up", func(t *testing.T) {
		fb := CTUD{}
		fb.INIT()
		fb.PV = 5
		fb.EN = true

		fb.CU = true
		fb.Execute()
		if fb.CV != 1 {
			t.Errorf("CV should be 1 after first count up, got %d", fb.CV)
		}
		fb.CU = false
		fb.Execute() // Reset for next edge
	})

	t.Run("Count Down", func(t *testing.T) {
		fb := CTUD{}
		fb.INIT()
		fb.PV = 5
		fb.EN = true
		fb.CV = 3 // Set initial value

		fb.CD = true
		fb.Execute()
		if fb.CV != 2 {
			t.Errorf("CV should be 2 after count down, got %d", fb.CV)
		}
		fb.CD = false
		fb.Execute() // Reset for next edge
	})

	t.Run("Reset and Load", func(t *testing.T) {
		fb := CTUD{}
		fb.INIT()
		fb.PV = 10
		fb.EN = true
		fb.CV = 5

		// Test Load
		fb.LD = true
		fb.Execute() // Load should be immediate
		if fb.CV != 10 {
			t.Errorf("CV should be PV after load, got %d", fb.CV)
		}
		fb.LD = false
		fb.Execute()

		// Test Reset
		fb.R = true
		fb.Execute() // Reset should be immediate
		if fb.CV != 0 {
			t.Errorf("CV should be 0 after reset, got %d", fb.CV)
		}
	})

	t.Run("Simultaneous CU and CD", func(t *testing.T) {
		fb := CTUD{}
		fb.INIT()
		fb.PV = 10
		fb.EN = true
		fb.CV = 5 // Start at a known value

		// First scan, inputs are low. This sets the edge detectors' memory.
		fb.Execute()

		// Second scan, both inputs go high simultaneously.
		fb.CU = true
		fb.CD = true
		fb.Execute()

		if fb.CV != 5 {
			t.Errorf("CV should not change when CU and CD have simultaneous rising edges, got %d, want 5", fb.CV)
		}
	})

	t.Run("Count past limits", func(t *testing.T) {
		fb := CTUD{}
		fb.INIT()
		fb.PV = 3
		fb.EN = true
		fb.CV = 2

		// Count up to PV
		fb.CU = true
		fb.Execute()
		if fb.CV != 3 || !fb.QU || fb.QD {
			t.Errorf("State at PV incorrect: CV=%d, QU=%v, QD=%v", fb.CV, fb.QU, fb.QD)
		}
		fb.CU = false
		fb.Execute()

		// Count past PV
		fb.CU = true
		fb.Execute()
		if fb.CV != 4 || !fb.QU || fb.QD {
			t.Errorf("State past PV incorrect: CV=%d, QU=%v, QD=%v", fb.CV, fb.QU, fb.QD)
		}
		fb.CU = false
		fb.Execute()
	})
}

func TestUninitializedCounters(t *testing.T) {
	t.Run("Uninitialized CTU", func(t *testing.T) {
		ctu := &CTU{} // Do not call INIT()
		ctu.PV = 5
		ctu.EN = true
		ctu.CU = true
		ctu.Execute() // First execute should initialize and detect rising edge

		if ctu.CV != 1 {
			t.Errorf("Uninitialized CTU should count on first rising edge. CV=%d, want 1", ctu.CV)
		}
	})

	t.Run("Uninitialized CTD", func(t *testing.T) {
		ctd := &CTD{} // Do not call INIT()
		ctd.PV = 5
		ctd.EN = true
		ctd.LD = true
		ctd.Execute() // First execute should initialize and load PV

		if ctd.CV != 5 {
			t.Errorf("Uninitialized CTD should load on first execute. CV=%d, want 5", ctd.CV)
		}
	})

	t.Run("Uninitialized CTUD", func(t *testing.T) {
		ctud := &CTUD{} // Do not call INIT()
		ctud.PV = 5
		ctud.EN = true
		ctud.Execute() // Should initialize to 0

		if ctud.CV != 0 {
			t.Errorf("Uninitialized CTUD should have CV=0. Got %d", ctud.CV)
		}

		ctud.CU = true
		ctud.Execute() // Should now count up

		if ctud.CV != 1 {
			t.Errorf("Uninitialized CTUD should count up after first scan. CV=%d, want 1", ctud.CV)
		}
	})
}
