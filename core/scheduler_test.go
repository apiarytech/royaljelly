//go:build !tinygo

/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - GPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var t0 = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// scanAt drives the scheduler deterministically without a ticker.
func scanAt(r *Resource, offsets ...time.Duration) {
	for _, d := range offsets {
		r.scan(t0.Add(d))
	}
}

func TestCyclicScheduleDoesNotAccumulateDrift(t *testing.T) {
	var runs []time.Duration
	task := NewTask("T", CyclicTask, 1, 100*time.Millisecond)
	task.AddProgram(&Program{Name: "P", Logic: func(now time.Time) { runs = append(runs, now.Sub(t0)) }})
	r := (&Resource{Name: "R", Cycle: 30 * time.Millisecond}).WithTask(task)

	// Scans every 30ms: the task must run at the first scan at or after each
	// multiple of 100ms from its first run, never sliding later over time.
	for i := 0; i <= 20; i++ {
		scanAt(r, time.Duration(i)*30*time.Millisecond)
	}
	want := []time.Duration{0, 120, 210, 300, 420, 510, 600}
	if len(runs) != len(want) {
		t.Fatalf("runs = %v, want %v (ms)", runs, want)
	}
	for i := range want {
		if runs[i] != want[i]*time.Millisecond {
			t.Fatalf("run %d at %v, want %vms (all runs %v)", i, runs[i], want[i], runs)
		}
	}
	if task.Drift() != 0 {
		t.Errorf("last run at 600ms is on schedule; Drift = %v, want 0", task.Drift())
	}
	if task.Overruns() != 0 {
		t.Errorf("Overruns = %d, want 0", task.Overruns())
	}
}

func TestFirstRunHasZeroDrift(t *testing.T) {
	task := NewTask("T", CyclicTask, 1, time.Second)
	r := (&Resource{Name: "R"}).WithTask(task)
	scanAt(r, 0)
	if d := task.Drift(); d != 0 {
		t.Fatalf("Drift after first run = %v, want 0", d)
	}
}

func TestOverrunIsCountedAndReported(t *testing.T) {
	var faults []Fault
	task := NewTask("T", CyclicTask, 1, 100*time.Millisecond)
	r := (&Resource{Name: "R", OnFault: func(f Fault) { faults = append(faults, f) }}).WithTask(task)

	scanAt(r, 0, 100*time.Millisecond)
	// Next scan is 250ms late: the run due at 200ms happens now, and the runs
	// due at 300ms and 400ms are skipped.
	scanAt(r, 450*time.Millisecond)

	if got := task.Overruns(); got != 2 {
		t.Fatalf("Overruns = %d, want 2", got)
	}
	if len(faults) != 1 || faults[0].Kind != FaultOverrun || !errors.Is(faults[0], ErrOverrun) {
		t.Fatalf("faults = %v, want one overrun fault", faults)
	}
	if faults[0].Resource != "R" || faults[0].Task != "T" {
		t.Errorf("fault location = %q/%q, want R/T", faults[0].Resource, faults[0].Task)
	}
	// The schedule resumes on the original grid, at 500ms.
	scanAt(r, 480*time.Millisecond)
	if s := task.Stats(); s.Runs != 3 {
		t.Fatalf("task ran early after an overrun: Runs = %d, want 3", s.Runs)
	}
	scanAt(r, 500*time.Millisecond)
	if s := task.Stats(); s.Runs != 4 || s.Drift != 0 {
		t.Fatalf("after resync Runs = %d Drift = %v, want 4 and 0", s.Runs, s.Drift)
	}
}

func TestReenabledTaskDoesNotReportOverruns(t *testing.T) {
	task := NewTask("T", CyclicTask, 1, 100*time.Millisecond)
	r := (&Resource{Name: "R", OnFault: func(Fault) {}}).WithTask(task)
	scanAt(r, 0)
	task.Disable()
	scanAt(r, 500*time.Millisecond)
	task.Enable()
	scanAt(r, time.Second)
	if s := task.Stats(); s.Runs != 2 || s.Overruns != 0 {
		t.Fatalf("Runs = %d Overruns = %d, want 2 and 0", s.Runs, s.Overruns)
	}
}

func TestEventTriggersAreQueued(t *testing.T) {
	var runs int
	task := NewTask("E", EventDrivenTask, 1, 0)
	task.AddProgram(&Program{Name: "P", Logic: func(time.Time) { runs++ }})
	r := (&Resource{Name: "R"}).WithTask(task)
	for i := 0; i < 3; i++ {
		if err := task.Trigger(); err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 5; i++ {
		scanAt(r, time.Duration(i)*time.Millisecond)
	}
	if runs != 3 {
		t.Fatalf("runs = %d, want 3 (one per trigger)", runs)
	}
}

func TestPanicIsReportedAndOtherProgramsRun(t *testing.T) {
	var faults []Fault
	var after bool
	task := NewTask("T", CyclicTask, 1, time.Second).
		WithProgram(&Program{Name: "Bad", Logic: func(time.Time) { panic("boom") }}).
		WithProgram(&Program{Name: "Good", Logic: func(time.Time) { after = true }})
	r := (&Resource{Name: "R", OnFault: func(f Fault) { faults = append(faults, f) }}).WithTask(task)
	scanAt(r, 0)

	if !after {
		t.Error("program after the panicking one did not run")
	}
	if len(faults) != 1 {
		t.Fatalf("faults = %v, want 1", faults)
	}
	f := faults[0]
	if f.Kind != FaultPanic || f.Program != "Bad" || f.Value != "boom" || !errors.Is(f, ErrProgramPanic) {
		t.Errorf("unexpected fault %+v", f)
	}
	if !strings.Contains(f.Error(), `program "Bad"`) {
		t.Errorf("Error() = %q, want it to name the program", f.Error())
	}
}

func TestWatchdogReportsLongRun(t *testing.T) {
	faults := make(chan Fault, 1)
	task := NewTask("Slow", CyclicTask, 1, time.Second)
	task.Watchdog = 5 * time.Millisecond
	task.AddProgram(&Program{Name: "P", Logic: func(time.Time) { time.Sleep(50 * time.Millisecond) }})
	r := (&Resource{Name: "R", OnFault: func(f Fault) { faults <- f }}).WithTask(task)
	scanAt(r, 0)

	select {
	case f := <-faults:
		if f.Kind != FaultWatchdog || f.Task != "Slow" || !errors.Is(f, ErrWatchdog) {
			t.Fatalf("unexpected fault %+v", f)
		}
	case <-time.After(time.Second):
		t.Fatal("watchdog did not fire")
	}
	if got := task.Stats().WatchdogTrips; got != 1 {
		t.Errorf("WatchdogTrips = %d, want 1", got)
	}

	// A fast run must not trip the watchdog.
	task.RemoveProgram("P")
	scanAt(r, time.Second)
	time.Sleep(20 * time.Millisecond)
	if got := task.Stats().WatchdogTrips; got != 1 {
		t.Errorf("fast run tripped the watchdog: WatchdogTrips = %d", got)
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name  string
		build func() *Resource
		want  string // substring of the error; empty means valid
	}{
		{"valid", func() *Resource {
			return (&Resource{Name: "R", Cycle: 10 * time.Millisecond}).
				WithTask(NewTask("A", CyclicTask, 1, 100*time.Millisecond)).
				WithTask(NewTask("B", EventDrivenTask, 2, 0))
		}, ""},
		{"cycle not faster than interval", func() *Resource {
			return (&Resource{Name: "R", Cycle: 100 * time.Millisecond}).WithTask(NewTask("A", CyclicTask, 1, 100*time.Millisecond))
		}, "must be faster"},
		{"zero interval", func() *Resource {
			return (&Resource{Name: "R"}).WithTask(NewTask("A", CyclicTask, 1, 0))
		}, "positive interval"},
		{"duplicate priority", func() *Resource {
			return (&Resource{Name: "R"}).
				WithTask(NewTask("A", CyclicTask, 1, time.Second)).
				WithTask(NewTask("B", CyclicTask, 1, time.Second))
		}, "duplicate task priority 1"},
		{"duplicate task name", func() *Resource {
			return (&Resource{Name: "R"}).
				WithTask(NewTask("A", CyclicTask, 1, time.Second)).
				WithTask(NewTask("A", CyclicTask, 2, time.Second))
		}, "duplicate task name"},
		{"duplicate program name", func() *Resource {
			return (&Resource{Name: "R"}).WithTask(NewTask("A", CyclicTask, 1, time.Second).
				WithProgram(&Program{Name: "P"}).WithProgram(&Program{Name: "P"}))
		}, "duplicate program name"},
		{"negative watchdog", func() *Resource {
			task := NewTask("A", CyclicTask, 1, time.Second)
			task.Watchdog = -1
			return (&Resource{Name: "R"}).WithTask(task)
		}, "watchdog"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.build()
			err := r.Validate()
			if tc.want == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want it to contain %q", err, tc.want)
			}
			if startErr := r.Start(); startErr == nil {
				r.Stop()
				t.Fatal("Start succeeded on an invalid resource")
			}
			if r.IsRunning() {
				t.Fatal("invalid resource is running")
			}
		})
	}
}

func TestConfigurationStartStopWithContext(t *testing.T) {
	var runs atomic.Int32
	task := NewTask("T", CyclicTask, 1, 5*time.Millisecond)
	task.AddProgram(&Program{Name: "P", Logic: func(time.Time) { runs.Add(1) }})
	r1 := (&Resource{Name: "R1", Cycle: time.Millisecond}).WithTask(task)
	r2 := &Resource{Name: "R2", Cycle: time.Millisecond}
	cfg := (&Configuration{Name: "C"}).WithResource(r1).WithResource(r2)

	ctx, cancel := context.WithCancel(context.Background())
	if err := cfg.Start(ctx); err != nil {
		t.Fatal(err)
	}
	if !r1.IsRunning() || !r2.IsRunning() {
		t.Fatal("resources not running after Start")
	}
	deadline := time.Now().Add(time.Second)
	for runs.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	cancel()
	for (r1.IsRunning() || r2.IsRunning()) && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if r1.IsRunning() || r2.IsRunning() {
		t.Fatal("resources still running after context cancel")
	}
	if runs.Load() == 0 {
		t.Error("task never ran")
	}
}

func TestConfigurationRunReturnsAfterDeadline(t *testing.T) {
	r := (&Resource{Name: "R"}).WithTask(NewTask("T", CyclicTask, 1, 10*time.Millisecond))
	cfg := (&Configuration{Name: "C"}).WithResource(r)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if err := cfg.Run(ctx); err != nil {
		t.Fatal(err)
	}
	if r.IsRunning() {
		t.Fatal("resource still running after Run returned")
	}
}

func TestConfigurationStartRollsBackOnError(t *testing.T) {
	good := (&Resource{Name: "Good"}).WithTask(NewTask("T", CyclicTask, 1, 10*time.Millisecond))
	bad := (&Resource{Name: "Bad"}).WithTask(NewTask("T", CyclicTask, 1, 0))
	cfg := (&Configuration{Name: "C"}).WithResource(good).WithResource(bad)
	if err := cfg.Start(nil); err == nil {
		cfg.Stop()
		t.Fatal("Start succeeded with an invalid resource")
	}
	if good.IsRunning() {
		t.Fatal("valid resource left running after failed Start")
	}
	dup := (&Configuration{Name: "D"}).WithResource(&Resource{Name: "X"}).WithResource(&Resource{Name: "X"})
	if err := dup.Validate(); err == nil || !strings.Contains(err.Error(), "duplicate resource name") {
		t.Fatalf("Validate = %v, want duplicate resource name error", err)
	}
}

func TestConfigurationFaultHandlerIsInherited(t *testing.T) {
	faults := make(chan Fault, 1)
	task := NewTask("T", CyclicTask, 1, 10*time.Millisecond)
	task.AddProgram(&Program{Name: "P", Logic: func(time.Time) { panic("x") }})
	r := (&Resource{Name: "R"}).WithTask(task)
	cfg := &Configuration{Name: "C", OnFault: func(f Fault) {
		select {
		case faults <- f:
		default:
		}
	}}
	cfg.WithResource(r)
	if err := cfg.Start(nil); err != nil {
		t.Fatal(err)
	}
	defer cfg.Stop()
	select {
	case f := <-faults:
		if f.Kind != FaultPanic {
			t.Fatalf("unexpected fault %v", f)
		}
	case <-time.After(time.Second):
		t.Fatal("configuration fault handler was not called")
	}
}

func TestAffinityOutOfRangeFailsStart(t *testing.T) {
	r := &Resource{Name: "R", Affinity: 1 << 20}
	err := r.Start()
	if err == nil {
		r.Stop()
		t.Fatal("Start succeeded with an impossible CPU affinity")
	}
	if r.IsRunning() {
		t.Fatal("resource running after failed Start")
	}
}

// TestConcurrentModificationWhileRunning exercises every mutating method while
// the scheduler runs. Run it with -race to detect unsynchronized access.
func TestConcurrentModificationWhileRunning(t *testing.T) {
	r := &Resource{Name: "R", Cycle: 100 * time.Microsecond, OnFault: func(Fault) {}}
	base := NewTask("Base", CyclicTask, 0, time.Millisecond)
	r.AddTask(base)
	if err := r.Start(); err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	var wg sync.WaitGroup
	stop := make(chan struct{})
	worker := func(fn func(i int)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; ; i++ {
				select {
				case <-stop:
					return
				default:
					fn(i)
				}
			}
		}()
	}
	worker(func(i int) {
		name := fmt.Sprintf("P%d", i%8)
		base.AddProgram(&Program{Name: name, Logic: func(time.Time) {}})
		base.RemoveProgram(name)
	})
	worker(func(i int) {
		name := fmt.Sprintf("T%d", i%8)
		task := NewTask(name, EventDrivenTask, 100+i%8, 0)
		r.AddTask(task)
		_ = task.Trigger()
		r.RemoveTask(name)
	})
	worker(func(i int) {
		if i%2 == 0 {
			r.Pause()
		} else {
			r.Resume()
		}
		_ = base.Stats()
		_ = r.Tasks()
		_ = base.Programs()
	})
	time.Sleep(100 * time.Millisecond)
	close(stop)
	wg.Wait()
}

func TestStopFromAnotherGoroutineAndRestart(t *testing.T) {
	r := (&Resource{Name: "R"}).WithTask(NewTask("T", CyclicTask, 1, 5*time.Millisecond))
	for i := 0; i < 3; i++ {
		if err := r.Start(); err != nil {
			t.Fatal(err)
		}
		var wg sync.WaitGroup
		for j := 0; j < 4; j++ {
			wg.Add(1)
			go func() { defer wg.Done(); r.Stop() }()
		}
		wg.Wait()
		if r.IsRunning() {
			t.Fatal("resource running after Stop")
		}
	}
}

func TestScanDoesNotAllocate(t *testing.T) {
	r := benchResource(8, 4)
	var i int
	allocs := testing.AllocsPerRun(1000, func() {
		i++
		r.scan(t0.Add(time.Duration(i) * time.Millisecond))
	})
	if allocs != 0 {
		t.Fatalf("scan allocates %.1f times per call, want 0", allocs)
	}
}

func benchResource(tasks, programs int) *Resource {
	r := &Resource{Name: "Bench", Cycle: time.Millisecond}
	var sink int
	for i := 0; i < tasks; i++ {
		task := NewTask(fmt.Sprintf("T%d", i), CyclicTask, i, time.Duration(i+2)*time.Millisecond)
		for j := 0; j < programs; j++ {
			task.AddProgram(&Program{Name: fmt.Sprintf("P%d", j), Logic: func(time.Time) { sink++ }})
		}
		r.AddTask(task)
	}
	return r
}

// BenchmarkScan measures scheduler overhead per scan with 8 tasks of 4 programs.
func BenchmarkScan(b *testing.B) {
	r := benchResource(8, 4)
	b.ReportAllocs()
	now := t0
	for b.Loop() {
		now = now.Add(time.Millisecond)
		r.scan(now)
	}
}

// BenchmarkJitter runs a real 1ms cyclic task and reports how late runs start
// relative to their schedule. Each iteration is one task run.
func BenchmarkJitter(b *testing.B) {
	var mu sync.Mutex
	var worst, total time.Duration
	var n int
	done := make(chan struct{})
	task := NewTask("J", CyclicTask, 1, time.Millisecond)
	r := (&Resource{Name: "J", Cycle: 250 * time.Microsecond, OnFault: func(Fault) {}}).WithTask(task)
	task.AddProgram(&Program{Name: "P", Logic: func(time.Time) {
		mu.Lock()
		defer mu.Unlock()
		d := task.Drift()
		total += d
		if d > worst {
			worst = d
		}
		n++
		if n == b.N {
			close(done)
		}
	}})
	b.ResetTimer()
	if err := r.Start(); err != nil {
		b.Fatal(err)
	}
	<-done
	r.Stop()
	b.StopTimer()
	b.ReportMetric(float64(worst.Microseconds()), "worst-drift-µs")
	b.ReportMetric(float64(total.Microseconds())/float64(n), "mean-drift-µs")
	b.ReportMetric(float64(task.Overruns()), "overruns")
}
