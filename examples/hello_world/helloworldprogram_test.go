package main

import (
	"testing"
	"time"
)

func TestHelloWorldProgram_Logic(t *testing.T) {
	// --- Setup ---
	p := &HelloWorldProgram{}
	p.Init()

	// Use a fixed time for predictable test runs.
	now := time.Now()

	// --- Initial State Verification ---
	if p.MyButton != false {
		t.Fatalf("Initial state for MyButton should be false, but got %v", p.MyButton)
	}
	if p.MyLamp != false {
		t.Fatalf("Initial state for MyLamp should be false, but got %v", p.MyLamp)
	}

	// --- Test Case 1: Button is OFF ---
	t.Run("Button_Is_Off", func(t *testing.T) {
		p.MyButton = false
		p.Logic(now)
		if p.MyLamp != false {
			t.Errorf("Expected MyLamp to be false when MyButton is false, but got %v", p.MyLamp)
		}
	})

	// --- Test Case 2: Button is ON ---
	t.Run("Button_Is_On", func(t *testing.T) {
		p.MyButton = true
		p.Logic(now)
		if p.MyLamp != true {
			t.Errorf("Expected MyLamp to be true when MyButton is true, but got %v", p.MyLamp)
		}
	})

	// --- Test Case 3: Button is turned OFF again ---
	t.Run("Button_Is_Off_Again", func(t *testing.T) {
		p.MyButton = false
		p.Logic(now)
		if p.MyLamp != false {
			t.Errorf("Expected MyLamp to be false after MyButton is turned off, but got %v", p.MyLamp)
		}
	})
}
