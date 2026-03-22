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
	fb.DERIVATIVE()
	if fb.XOUT != 0.0 {
		t.Errorf("Initial XOUT should be 0.0, got %f", fb.XOUT)
	}

	// Step 1
	fb.XIN = 10.0
	fb.DERIVATIVE()

	// Step 2
	fb.XIN = 20.0
	fb.DERIVATIVE()

	// Step 3
	fb.XIN = 30.0
	fb.DERIVATIVE()

	// Step 4 - now we have enough history
	fb.XIN = 40.0
	fb.DERIVATIVE()

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
	fb.DERIVATIVE()
	if fb.XOUT != 0.0 {
		t.Errorf("XOUT should be 0.0 when RUN is false, got %f", fb.XOUT)
	}
}

func TestINTEGRAL(t *testing.T) {
	fb := INTEGRAL{}
	fb.INIT()
	fb.CYCLE = TIME(1 * time.Second)
	fb.RUN = true

	// Step 1
	fb.XIN = 10.0
	fb.INTEGRAL()
	if fb.XOUT != 10.0 {
		t.Errorf("INTEGRAL after 1 step = %f; want 10.0", fb.XOUT)
	}

	// Step 2
	fb.XIN = 10.0
	fb.INTEGRAL()
	if fb.XOUT != 20.0 {
		t.Errorf("INTEGRAL after 2 steps = %f; want 20.0", fb.XOUT)
	}

	// Test Reset
	fb.R1 = true
	fb.X0 = 5.0
	fb.RUN = false // When R1 is true, RUN is typically false (e.g., PID manual mode)
	fb.INTEGRAL()
	// When R1 is true, XOUT should be set to X0.
	if fb.XOUT != fb.X0 {
		t.Errorf("INTEGRAL after reset = %f; want 5.0", fb.XOUT)
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
		pid.TR = 10.0
		pid.TD = 1.0
		pid.CONTROL_ACTION = false
		pid.CYCLE = TIME(1 * time.Second) // 1000

		pid.PID() // Execute one cycle

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
		// Expected XOUT0 = P + I + D = 20.0 + 2.0 + 6.0 = 28.0
		expectedXOUT0 := REAL(28.0)
		if !withinTolerance(pid.XOUT0, expectedXOUT0, 1e-6) {
			t.Errorf("PID output on first cycle is incorrect. Got %f, expected %f", pid.XOUT0, expectedXOUT0)
		}
	})

	t.Run("Step Response Simulation", func(t *testing.T) {
		pid := PID{}
		pid.INIT()
		pid.AUTO = true
		pid.SP = 100.0
		pid.PV = 0.0
		pid.KP = 2.5
		pid.TR = 15.0 // Integral time
		pid.TD = 1.0  // Derivative time
		pid.CYCLE = TIME(1 * time.Second)

		process := simpleProcess{pv: 0.0, gain: 1.0, timeConstant: 20.0}

		// Simulate for a number of cycles to allow the process to settle
		for i := 0; i < 100; i++ {
			pid.PV = process.pv
			pid.PID()
			process.update(pid.XOUT0, pid.CYCLE)
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
		pid.TR = 10.0
		pid.TD = 1.0
		pid.X0 = 25.0 // Manual output value
		pid.CYCLE = TIME(1 * time.Second)

		pid.PID()

		// In manual mode, the integral term should be reset and the output should be X0
		if !pid.ITERM.R1 {
			t.Error("Integral term should be in reset when AUTO is false")
		}

		// In manual mode (R1=true), the INTEGRAL block's output (XOUT) is set to its initial condition input (X0) for bumpless transfer.
		// The PID block calculates ITERM.X0 for bumpless transfer: ITERM.X0 = X0 - pTerm
		// pTerm = KP * ERROR = 1.0 * (100 - 50) = 50.0
		// ITERM.X0 = 25.0 - 50.0 = -25.0
		expectedIntegralOutput := REAL(-25.0)
		if pid.ITERM.XOUT != expectedIntegralOutput {
			t.Errorf("Integral output should be held at its calculated X0 for bumpless transfer in manual mode. got %f, want %f", pid.ITERM.XOUT, expectedIntegralOutput)
		}
	})
}
