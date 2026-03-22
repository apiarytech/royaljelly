package royaljelly

import (
	"testing"
	"time"
)

func TestTP(t *testing.T) {
	// Mock clock starting at a specific time
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tp := &TP{}
	tp.INIT()
	tp.PT = TIME(5 * time.Second)

	// 1. Initial state: IN is false
	tp.Execute(now)
	if tp.Q || tp.ET != 0 {
		t.Errorf("Initial state incorrect. Q=%v, ET=%v", tp.Q, tp.ET)
	}

	// 2. Rising edge on IN, start timing
	now = now.Add(time.Second) // t=1s
	tp.IN = true
	tp.Execute(now)
	if !tp.Q || tp.ET != 0 {
		t.Errorf("After rising edge, Q should be true and ET should be 0. Q=%v, ET=%v", tp.Q, tp.ET)
	}

	// 3. IN stays true, timer running
	now = now.Add(3 * time.Second) // t=4s
	tp.Execute(now)
	if !tp.Q || tp.ET != TIME(3*time.Second) {
		t.Errorf("Timer running incorrect. Q=%v, ET=%v, want Q=true, ET=3s", tp.Q, tp.ET)
	}

	// 4. Timer expires, IN is still true
	now = now.Add(3 * time.Second) // t=7s
	tp.Execute(now)
	if tp.Q || tp.ET != tp.PT {
		t.Errorf("Timer expired (IN high) incorrect. Q=%v, ET=%v, want Q=false, ET=5s", tp.Q, tp.ET)
	}

	// 5. IN goes false, timer should be fully reset
	now = now.Add(time.Second) // t=8s
	tp.IN = false
	tp.Execute(now)
	if tp.Q || tp.ET != 0 {
		t.Errorf("After IN goes false, timer should reset. Q=%v, ET=%v", tp.Q, tp.ET)
	}

	// 6. Test re-trigger
	now = now.Add(time.Second) // t=9s
	tp.IN = true
	tp.Execute(now)
	if !tp.Q || tp.ET != 0 {
		t.Errorf("Re-trigger failed. Q=%v, ET=%v", tp.Q, tp.ET)
	}

	// 7. Test reset after pulse completion
	t.Run("Reset After Pulse", func(t *testing.T) {
		pulseTime := time.Second * 2
		p := &TP{}
		p.INIT()
		p.PT = TIME(pulseTime)
		clk := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

		// Start pulse
		p.IN = true
		p.Execute(clk)

		// Let pulse finish
		clk = clk.Add(pulseTime + time.Millisecond)
		p.Execute(clk)
		if p.Q {
			t.Error("Q should be false after pulse is complete")
		}
		if p.ET != p.PT {
			t.Errorf("ET should be clamped to PT after pulse. Got %v, want %v", p.ET, p.PT)
		}

		// Input goes low, should reset ET
		clk = clk.Add(time.Second)
		p.IN = false
		p.Execute(clk)
		if p.ET != 0 {
			t.Errorf("ET should reset to 0 when IN goes low after a pulse. Got %v", p.ET)
		}
	})
}

func TestTON(t *testing.T) {
	// Mock clock
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	ton := &TON{}
	ton.INIT()
	ton.PT = TIME(5 * time.Second)

	// 1. Initial state: IN is false
	ton.Execute(now)
	if ton.Q || ton.ET != 0 {
		t.Errorf("Initial state incorrect. Q=%v, ET=%v", ton.Q, ton.ET)
	}

	// 2. IN goes true, start timing
	now = now.Add(time.Second) // t=1s
	ton.IN = true
	ton.Execute(now)
	if ton.Q || ton.ET != 0 {
		t.Errorf("After IN goes true, Q should be false and ET should be 0. Q=%v, ET=%v", ton.Q, ton.ET)
	}

	// 3. Timer running
	now = now.Add(3 * time.Second) // t=4s
	ton.Execute(now)
	if ton.Q || ton.ET != TIME(3*time.Second) {
		t.Errorf("Timer running incorrect. Q=%v, ET=%v, want Q=false, ET=3s", ton.Q, ton.ET)
	}

	// 4. Timer expires
	now = now.Add(3 * time.Second) // t=7s
	ton.Execute(now)
	if !ton.Q || ton.ET != ton.PT {
		t.Errorf("Timer expired incorrect. Q=%v, ET=%v, want Q=true, ET=5s", ton.Q, ton.ET)
	}

	// 5. IN goes false, timer resets
	now = now.Add(time.Second) // t=8s
	ton.IN = false
	ton.Execute(now)
	if ton.Q || ton.ET != 0 {
		t.Errorf("After IN goes false, timer should reset. Q=%v, ET=%v", ton.Q, ton.ET)
	}
}

func TestTOF(t *testing.T) {
	// Mock clock
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tof := &TOF{}
	tof.INIT()
	tof.PT = TIME(5 * time.Second)

	// 1. Initial state: IN is false
	tof.Execute(now)
	if tof.Q || tof.ET != 0 {
		t.Errorf("Initial state incorrect. Q=%v, ET=%v", tof.Q, tof.ET)
	}

	// 2. IN goes true, output is true
	now = now.Add(time.Second) // t=1s
	tof.IN = true
	tof.Execute(now)
	if !tof.Q || tof.ET != 0 {
		t.Errorf("IN is true, Q should be true and ET 0. Q=%v, ET=%v", tof.Q, tof.ET)
	}

	// 3. IN goes false (falling edge), start timing
	now = now.Add(time.Second) // t=2s
	tof.IN = false
	tof.Execute(now)
	if !tof.Q || tof.ET != 0 {
		t.Errorf("After falling edge, Q should be true and ET 0. Q=%v, ET=%v", tof.Q, tof.ET)
	}

	// 4. Timer running
	now = now.Add(3 * time.Second) // t=5s
	tof.Execute(now)
	if !tof.Q || tof.ET != TIME(3*time.Second) {
		t.Errorf("Timer running incorrect. Q=%v, ET=%v, want Q=true, ET=3s", tof.Q, tof.ET)
	}

	// 5. Timer expires
	now = now.Add(3 * time.Second) // t=8s
	tof.Execute(now)
	if tof.Q || tof.ET != tof.PT {
		t.Errorf("Timer expired incorrect. Q=%v, ET=%v, want Q=false, ET=5s", tof.Q, tof.ET)
	}

	// 6. IN goes true again, timer resets and Q is true
	now = now.Add(time.Second) // t=9s
	tof.IN = true
	tof.Execute(now)
	if !tof.Q || tof.ET != 0 {
		t.Errorf("IN goes true again, timer should reset. Q=%v, ET=%v", tof.Q, tof.ET)
	}
}

func TestTOF_QuickToggle(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	tof := &TOF{}
	tof.INIT()
	tof.PT = TIME(5 * time.Second)

	// IN goes high
	tof.IN = true
	tof.Execute(now)

	// IN goes low for less than PT
	now = now.Add(time.Second) // t=1s
	tof.IN = false
	tof.Execute(now)
	if !tof.Q {
		t.Error("Q should remain true immediately after falling edge")
	}

	// IN goes high again before timer expires
	now = now.Add(2 * time.Second) // t=3s
	tof.IN = true
	tof.Execute(now)
	if !tof.Q || tof.ET != 0 {
		t.Errorf("IN went high before expiry, Q should be true and ET should reset. Q=%v, ET=%v", tof.Q, tof.ET)
	}
}
