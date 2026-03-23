package royaljelly

import (
	"testing"
	"time"
)

func TestDERIVATIVE(t *testing.T) {
	fb := DERIVATIVE{}
	fb.INIT()
	fb.CYCLE = TIME(1 * time.Second)
	fb.RUN = true

	// Initial state
	_ = fb.DERIVATIVE()
	if fb.XOUT != 0.0 {
		t.Errorf("Initial XOUT should be 0.0, got %f", fb.XOUT)
	}

	// Step 1
	fb.XIN = 10.0
	_ = fb.DERIVATIVE()

	// Step 2
	fb.XIN = 20.0
	_ = fb.DERIVATIVE()

	// Step 3
	fb.XIN = 30.0
	_ = fb.DERIVATIVE()

	// Step 4 - now we have enough history
	fb.XIN = 40.0
	_ = fb.DERIVATIVE()

	// After step 3, internal states are: X1=30, X2=20, X3=10
	// In step 4, XIN=40.
	// XOUT = (3.0*(XIN-X3) + X1 - X2) / (10.0 * CYCLE_s)
	// XOUT = (3.0*(40.0-10.0) + 30.0 - 20.0) / (10.0 * 1.0)
	// XOUT = (3.0*30.0 + 10.0) / 10.0 = (90.0 + 10.0) / 10.0 = 100.0 / 10.0 = 10.0
	expected := REAL(10.0)
	if fb.XOUT != expected {
		t.Errorf("DERIVATIVE XOUT = %f; want %f", fb.XOUT, expected)
	}

	// Test RUN=false
	fb.RUN = false
	_ = fb.DERIVATIVE()
	if fb.XOUT != 0.0 {
		t.Errorf("XOUT should be 0.0 when RUN is false, got %f", fb.XOUT)
	}

	// Test division by zero
	fb.RUN = true
	fb.CYCLE = 0
	if err := fb.DERIVATIVE(); err == nil {
		t.Error("DERIVATIVE should return an error when CYCLE is zero")
	}
}

func TestINTEGRAL(t *testing.T) {
	fb := INTEGRAL{}
	fb.INIT()
	fb.CYCLE = TIME(1 * time.Second)
	fb.RUN = true

	// Step 1
	fb.XIN = 10.0
	_ = fb.INTEGRAL()
	if fb.XOUT != 10.0 {
		t.Errorf("INTEGRAL after 1 step = %f; want 10.0", fb.XOUT)
	}

	// Step 2
	fb.XIN = 10.0
	_ = fb.INTEGRAL()
	if fb.XOUT != 20.0 {
		t.Errorf("INTEGRAL after 2 steps = %f; want 20.0", fb.XOUT)
	}

	// Test Reset
	// First scan with R1=true should set XOUT to X0
	fb.R1 = true
	fb.X0 = 5.0
	fb.RUN = false // When R1 is true, RUN is typically false (e.g., PID manual mode)
	_ = fb.INTEGRAL()
	if fb.XOUT != fb.X0 {
		t.Errorf("INTEGRAL on first reset scan = %f; want %f", fb.XOUT, fb.X0)
	}

	// Second scan with R1=true should hold the value, not re-evaluate X0
	fb.X0 = 99.0 // Change X0 to see if XOUT is incorrectly updated
	_ = fb.INTEGRAL()
	if fb.XOUT != 5.0 {
		t.Errorf("INTEGRAL should hold value after reset edge, got %f; want 5.0", fb.XOUT)
	}

	// Test invalid cycle time
	fb.RUN = true
	fb.R1 = false
	fb.CYCLE = 0
	if err := fb.INTEGRAL(); err == nil {
		t.Error("INTEGRAL should return an error when CYCLE is zero")
	}
}

func TestHYSTERESIS(t *testing.T) {
	fb := HYSTERESIS{}
	fb.EN = true
	fb.ENO = true
	fb.EPS = 1.0
	fb.XIN2 = 10.0

	// Initially false
	if fb.Q != false {
		t.Errorf("Initial state should be false")
	}

	// Rise above upper limit
	fb.XIN1 = 11.1
	fb.HYSTERESIS()
	if fb.Q != true {
		t.Errorf("Should turn true when XIN1 > XIN2 + EPS")
	}

	// Stay true in deadband
	fb.XIN1 = 9.5
	fb.HYSTERESIS()
	if fb.Q != true {
		t.Errorf("Should stay true inside the deadband")
	}

	// Drop below lower limit
	fb.XIN1 = 8.9
	fb.HYSTERESIS()
	if fb.Q != false {
		t.Errorf("Should turn false when XIN1 < XIN2 - EPS")
	}
}

// Helper function to check if a value is within a certain tolerance.
func withinTolerance(a, b, tolerance REAL) bool {
	if a > b {
		return (a - b) <= tolerance
	}
	return (b - a) <= tolerance
}

// simpleProcess simulates a first-order process for PID testing.
type simpleProcess struct {
	pv           REAL
	gain         REAL
	timeConstant REAL
}

// update simulates the process response to the controller output over one cycle.
func (p *simpleProcess) update(controllerOutput REAL, cycleTime TIME) {
	// Simplified first-order process model: T * dy/dt + y = K * u
	// Discrete approximation: y(n) = y(n-1) + (K * u(n-1) - y(n-1)) * (CYCLE / T)
	if p.timeConstant > 0 {
		cycleSec := TIME_TO_LREAL(cycleTime) / 1000.0
		p.pv += (p.gain*controllerOutput - p.pv) * REAL(cycleSec) / p.timeConstant
	}
}

func TestPID(t *testing.T) {
	t.Run("Basic Sanity Check", func(t *testing.T) {
		pid := PID{}
		pid.INIT()
		pid.AUTO = true
		pid.SP = 100
		pid.PV = 90
		pid.KP = 2.0
		pid.TR = 10.0                     // Integral time in seconds
		pid.TD = 1.0                      // Derivative time in seconds
		pid.DIRECT_ACTION = true          // Direct action (heating): error = SP - PV
		pid.CYCLE = TIME(1 * time.Second) // 1000

		if err := pid.PID(); err != nil {
			t.Fatalf("PID execution failed: %v", err)
		}

		// Error should be SP - PV (Control Action = false)
		if pid.ERROR != 10.0 {
			t.Errorf("PID error calculation is incorrect, got %f, want 10.0", pid.ERROR)
		}

		// On the first cycle, DTERM is 0. ITERM will have one step of integration.
		// ITERM.XIN = (KP/TR)*ERROR = (2.0/10.0)*10.0 = 2.0
		// ITERM.XOUT = 0 + ITERM.XIN * (CYCLE/1000) = 2.0 * 1.0 = 2.0
		// P Term = KP * ERROR = 2.0 * 10.0 = 20.0
		// I Term = ITERM.XOUT = 2.0
		// D Term = KP * TD * DTERM.XOUT = 2 * 1 * 3 = 6
		// D Term = KP * TD * DTERM.XOUT = 2 * 1 * 3.0 = 6.0 (DTERM calculates 3.0 on first step)
		// Expected XOUT = P + I + D = 20.0 + 2.0 + 6.0 = 28.0
		expectedXOUT0 := REAL(28.0)
		if !withinTolerance(pid.XOUT, expectedXOUT0, 1e-6) {
			t.Errorf("PID output on first cycle is incorrect. Got %f, expected %f", pid.XOUT, expectedXOUT0)
		}
	})

	t.Run("Step Response Simulation", func(t *testing.T) {
		pid := PID{}
		pid.INIT()
		pid.AUTO = true
		pid.SP = 100.0
		pid.PV = 0.0
		pid.KP = 2.5
		pid.TR = 15.0 // Integral time in seconds
		pid.TD = 1.0  // Derivative time in seconds
		pid.DIRECT_ACTION = true
		pid.CYCLE = TIME(1 * time.Second)

		process := simpleProcess{pv: 0.0, gain: 1.0, timeConstant: 20.0}

		// Simulate for a number of cycles to allow the process to settle
		for i := 0; i < 100; i++ {
			pid.PV = process.pv
			if err := pid.PID(); err != nil {
				t.Fatalf("PID simulation failed at step %d: %v", i, err)
			}
			process.update(pid.XOUT, pid.CYCLE)
		}

		// After settling, the process variable should be close to the setpoint
		if !withinTolerance(process.pv, pid.SP, 1.0) {
			t.Errorf("PID did not bring process to setpoint. Final PV: %f, SP: %f", process.pv, pid.SP)
		}
	})

	t.Run("Manual Mode Test", func(t *testing.T) {
		pid := PID{}
		pid.INIT()
		pid.AUTO = false // Manual mode
		pid.SP = 100
		pid.PV = 50
		pid.KP = 1.0
		pid.TR = 10.0 // Integral time in seconds
		pid.TD = 1.0  // Derivative time in seconds
		pid.X0 = 25.0 // Manual output value
		pid.CYCLE = TIME(1 * time.Second)

		if err := pid.PID(); err != nil {
			t.Fatalf("PID execution in manual mode failed: %v", err)
		}

		// In manual mode, the integral term should be reset and the output should be X0
		if !pid.ITERM.R1 {
			t.Error("Integral term should be in reset when AUTO is false")
		}

		// In manual mode, the overall PID output (XOUT) should be exactly the manual setpoint (X0).
		if !withinTolerance(pid.XOUT, pid.X0, 1e-6) {
			t.Errorf("PID output in manual mode is incorrect. Got %f, expected %f", pid.XOUT, pid.X0)
		}
	})

	t.Run("Auto to Manual Bumpless Transfer", func(t *testing.T) {
		pid := PID{}
		pid.INIT()
		pid.AUTO = true // Start in auto mode
		pid.SP = 100
		pid.PV = 90
		pid.KP = 2.0
		pid.TR = 10.0 // Integral time in seconds
		pid.TD = 1.0  // Derivative time in seconds
		pid.CYCLE = TIME(1 * time.Second)

		// Run one cycle in auto to establish an output
		if err := pid.PID(); err != nil {
			t.Fatalf("PID auto cycle failed: %v", err)
		}
		initialAutoOutput := pid.XOUT

		// Switch to manual mode
		pid.AUTO = false
		if err := pid.PID(); err != nil {
			t.Fatalf("PID manual cycle failed: %v", err)
		}

		// When switching from auto to manual, X0 should be updated to the last auto output
		if !withinTolerance(pid.X0, initialAutoOutput, 1e-6) {
			t.Errorf("X0 not correctly updated on auto to manual switch. Got %f, expected %f", pid.X0, initialAutoOutput)
		}
	})

	t.Run("Reverse Action Test", func(t *testing.T) {
		pid := PID{}
		pid.INIT()
		pid.AUTO = true
		pid.SP = 50
		pid.PV = 60
		pid.KP = 2.0
		pid.TR = 10.0             // Integral time in seconds
		pid.TD = 1.0              // Derivative time in seconds
		pid.DIRECT_ACTION = false // Reverse action (cooling): error = PV - SP
		pid.CYCLE = TIME(1 * time.Second)

		if err := pid.PID(); err != nil {
			t.Fatalf("PID reverse action cycle failed: %v", err)
		}

		if pid.ERROR != 10.0 {
			t.Errorf("PID reverse action error calculation is incorrect, got %f, want 10.0", pid.ERROR)
		}
	})

	t.Run("TR Zero Error", func(t *testing.T) {
		pid := PID{}
		pid.INIT()
		pid.AUTO = true
		pid.TR = 0 // Set integral time to zero
		if err := pid.PID(); err == nil {
			t.Error("PID should return an error when TR is zero")
		}
	})
}
