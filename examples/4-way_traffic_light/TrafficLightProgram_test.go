package main

import (
	"testing"
	"time"
)

func TestTrafficLightProgram_Logic(t *testing.T) {
	// --- Setup ---
	p := &TrafficLightProgram{}
	p.Init()

	// Use a fixed start time for predictable test runs.
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	// --- Initial State Verification ---
	if p.state != stNSGreen {
		t.Fatalf("Initial state should be stNSGreen (%d), but got %d", stNSGreen, p.state)
	}

	// --- Test State: stNSGreen ---
	t.Run("State_NSGreen", func(t *testing.T) {
		// Run logic for the first time in this state
		p.Logic(now)
		if !p.nsLights.Green || !p.ewLights.Red {
			t.Errorf("Expected NS Green and EW Red, got NS: %+v, EW: %+v", p.nsLights, p.ewLights)
		}

		// Advance time just past the green light duration (10s)
		now = now.Add(10 * time.Second)
		p.Logic(now) // This scan's timer execution will set Q to true

		// The next logic scan should detect Q and transition the state
		p.Logic(now)
		if p.state != stNSYellow {
			t.Errorf("State should have transitioned to stNSYellow (%d), but is %d", stNSYellow, p.state)
		}
	})

	// --- Test State: stNSYellow ---
	t.Run("State_NSYellow", func(t *testing.T) {
		// Run logic for the first time in this state
		p.Logic(now)
		if !p.nsLights.Yellow || !p.ewLights.Red {
			t.Errorf("Expected NS Yellow and EW Red, got NS: %+v, EW: %+v", p.nsLights, p.ewLights)
		}

		// Advance time just past the yellow light duration (3s)
		now = now.Add(3 * time.Second)
		p.Logic(now) // Timer Q becomes true

		// Next scan transitions state
		p.Logic(now)
		if p.state != stEWGreen {
			t.Errorf("State should have transitioned to stEWGreen (%d), but is %d", stEWGreen, p.state)
		}
	})

	// --- Test State: stEWGreen ---
	t.Run("State_EWGreen", func(t *testing.T) {
		// Run logic for the first time in this state
		p.Logic(now)
		if !p.nsLights.Red || !p.ewLights.Green {
			t.Errorf("Expected NS Red and EW Green, got NS: %+v, EW: %+v", p.nsLights, p.ewLights)
		}

		// Advance time just past the green light duration (10s)
		now = now.Add(10 * time.Second)
		p.Logic(now) // Timer Q becomes true

		// Next scan transitions state
		p.Logic(now)
		if p.state != stEWYellow {
			t.Errorf("State should have transitioned to stEWYellow (%d), but is %d", stEWYellow, p.state)
		}
	})

	// --- Test State: stEWYellow ---
	t.Run("State_EWYellow", func(t *testing.T) {
		// Run logic for the first time in this state
		p.Logic(now)
		if !p.nsLights.Red || !p.ewLights.Yellow {
			t.Errorf("Expected NS Red and EW Yellow, got NS: %+v, EW: %+v", p.nsLights, p.ewLights)
		}

		// Advance time just past the yellow light duration (3s)
		now = now.Add(3 * time.Second)
		p.Logic(now) // Timer Q becomes true

		// Next scan transitions state, cycling back to the beginning
		p.Logic(now)
		if p.state != stNSGreen {
			t.Errorf("State should have cycled back to stNSGreen (%d), but is %d", stNSGreen, p.state)
		}
	})
}
