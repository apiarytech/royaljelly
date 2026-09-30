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

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/core"
	"github.com/apiarytech/royaljelly/iec"
)

// testLoader returns a loader with a private registry holding a "Counter" type
// that records the params it was created with.
func testLoader(t *testing.T) (Loader, map[string]map[string]string) {
	t.Helper()
	created := map[string]map[string]string{}
	reg := NewRegistry()
	reg.Register("Counter", func(params map[string]string) (func(time.Time), error) {
		var start iec.LINT
		if err := ParseLINT(params, "initial_value", &start); err != nil {
			return nil, err
		}
		created[params["name"]] = params
		return func(time.Time) {}, nil
	})
	reg.Register("Broken", func(map[string]string) (func(time.Time), error) {
		return nil, errors.New("factory failed")
	})
	return Loader{Registry: reg}, created
}

const validConfig = `
# Test configuration
name: Plant

resource: CPU1
  cycle: 10ms
  task: Slow        # trailing comment
    type: Cyclic
    priority: 2
    interval: 500ms
    program: C2 Counter
  task: Fast
    type: Cyclic
    priority: 1
    interval: 100ms
    watchdog: 20ms
    program: C1 Counter
      param: initial_value 42
      param: label a#b
  task: OnEvent
    type: EventDriven
    priority: 3
resource: CPU2
  affinity: 0
`

func TestLoadValidConfig(t *testing.T) {
	loader, created := testLoader(t)
	cfg, err := loader.Load(strings.NewReader(validConfig))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Name != "Plant" || len(cfg.Resources()) != 2 {
		t.Fatalf("got config %q with %d resources", cfg.Name, len(cfg.Resources()))
	}
	cpu1 := cfg.FindResource("CPU1")
	if cpu1 == nil || cpu1.Cycle != 10*time.Millisecond {
		t.Fatalf("CPU1 = %+v", cpu1)
	}
	tasks := cpu1.Tasks()
	var order []string
	for _, task := range tasks {
		order = append(order, task.Name)
	}
	if strings.Join(order, ",") != "Fast,Slow,OnEvent" {
		t.Errorf("task order = %v, want sorted by priority", order)
	}
	fast := cpu1.FindTask("Fast")
	if fast.Interval != 100*time.Millisecond || fast.Watchdog != 20*time.Millisecond {
		t.Errorf("Fast task = interval %v watchdog %v", fast.Interval, fast.Watchdog)
	}
	if cpu1.FindTask("OnEvent").Type != core.EventDrivenTask {
		t.Error("OnEvent should be event-driven")
	}
	if got := created["C1"]; got["initial_value"] != "42" || got["label"] != "a#b" || got["name"] != "C1" {
		t.Errorf("C1 params = %v", got)
	}
	if fast.FindProgram("C1").Logic == nil {
		t.Error("C1 has no logic")
	}
	if _, ok := created["C2"]; !ok {
		t.Error("C2 was not created")
	}
}

func TestIndentationStyles(t *testing.T) {
	body := func(step string) string {
		lines := []struct {
			level int
			text  string
		}{
			{0, "name: X"}, {0, "resource: R"}, {1, "cycle: 1ms"}, {1, "task: T"},
			{2, "priority: 1"}, {2, "interval: 10ms"}, {2, "program: P Counter"},
			{3, "param: initial_value 1"},
		}
		var sb strings.Builder
		for _, l := range lines {
			sb.WriteString(strings.Repeat(step, l.level) + l.text + "\n")
		}
		return sb.String()
	}
	for name, step := range map[string]string{"tabs": "\t", "two spaces": "  ", "four spaces": "    "} {
		t.Run(name, func(t *testing.T) {
			loader, created := testLoader(t)
			if _, err := loader.Load(strings.NewReader(body(step))); err != nil {
				t.Fatal(err)
			}
			if created["P"]["initial_value"] != "1" {
				t.Fatalf("param not parsed: %v", created["P"])
			}
		})
	}
}

func TestLoadErrors(t *testing.T) {
	cases := []struct {
		name, config, want string
	}{
		{"missing name", "resource: R\n", "must have a 'name'"},
		{"empty", "# nothing\n", "no 'name'"},
		{"no colon", "name: X\nresource R\n", "line 2: expected 'key: value'"},
		{"unknown top key", "name: X\nfoo: bar\n", "line 2: unknown top-level key 'foo'"},
		{"mixed tabs and spaces", "name: X\nresource: R\n \tcycle: 1ms\n", "line 3: indentation mixes tabs and spaces"},
		{"inconsistent step", "name: X\nresource: R\n  task: T\n     priority: 1\n", "line 4: indentation of 5"},
		{"tabs after spaces", "name: X\nresource: R\n  task: T\n\t\tpriority: 1\n", "line 4: inconsistent indentation"},
		{"too deep", "name: X\nresource: R\n  task: T\n    program: P Counter\n      param: a 1\n        x: y\n", "line 6: line is nested too deeply"},
		{"task prop without task", "name: X\nresource: R\n  cycle: 1ms\n    priority: 1\n", "without a task context"},
		{"bad duration", "name: X\nresource: R\n  cycle: fast\n", "line 3: invalid duration for cycle"},
		{"bad task type", "name: X\nresource: R\n  task: T\n    type: Sometimes\n", "unknown task type 'Sometimes'"},
		{"program without type", "name: X\nresource: R\n  task: T\n    program: OnlyName\n", "requires an instance name and a type name"},
		{"unknown program type", "name: X\nresource: R\n  task: T\n    priority: 1\n    interval: 1s\n    program: P Nope\n", "unknown program type 'Nope'"},
		{"factory error", "name: X\nresource: R\n  task: T\n    priority: 1\n    interval: 1s\n    program: P Broken\n", "factory failed"},
		{"bad param value", "name: X\nresource: R\n  task: T\n    priority: 1\n    interval: 1s\n    program: P Counter\n      param: initial_value abc\n", "invalid integer value for parameter 'initial_value'"},
		{"duplicate priority", "name: X\nresource: R\n  task: A\n    priority: 1\n    interval: 1s\n  task: B\n    priority: 1\n    interval: 1s\n", "duplicate task priority 1"},
		{"cycle too slow", "name: X\nresource: R\n  cycle: 1s\n  task: A\n    priority: 1\n    interval: 100ms\n", "must be faster"},
		{"cyclic without interval", "name: X\nresource: R\n  task: A\n    priority: 1\n", "positive interval"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			loader, _ := testLoader(t)
			_, err := loader.Load(strings.NewReader(tc.config))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want it to contain %q", err, tc.want)
			}
		})
	}
}

func TestLoadFileAndDefaultRegistry(t *testing.T) {
	RegisterProgramFactory("DefaultRegistryTestType", func(map[string]string) (func(time.Time), error) {
		return func(time.Time) {}, nil
	})
	text := "name: F\nresource: R\n  task: T\n    priority: 1\n    interval: 1s\n    program: P DefaultRegistryTestType\n"
	// Plain temp files are used instead of t.TempDir, whose cleanup fails under
	// TinyGo on Windows.
	path := writeTemp(t, "plc-*.txt", text)
	cfg, err := LoadConfigurationFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.FindResource("R").FindTask("T").FindProgram("P") == nil {
		t.Fatal("program not loaded")
	}
	if _, err := LoadConfigurationFromString(text); err != nil {
		t.Fatal(err)
	}
	_, err = LoadConfigurationFromFile(filepath.Join(os.TempDir(), "royaljelly-does-not-exist.txt"))
	if err == nil {
		t.Fatal("expected an error for a missing file")
	}
	bad := writeTemp(t, "bad-*.txt", "name: X\nnope: 1\n")
	if _, err := LoadConfigurationFromFile(bad); err == nil || !strings.Contains(err.Error(), filepath.Base(bad)+": line 2") {
		t.Fatalf("error = %v, want file name and line", err)
	}
}

// writeTemp creates a temporary file with the given content and removes it when
// the test ends.
func writeTemp(t *testing.T, pattern, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

func TestRegistryConcurrentUse(t *testing.T) {
	reg := NewRegistry()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("T%d", i)
			reg.Register(name, func(map[string]string) (func(time.Time), error) { return nil, nil })
			if _, ok := reg.Lookup(name); !ok {
				t.Errorf("%s not found", name)
			}
		}(i)
	}
	wg.Wait()
	var zero Registry
	if _, ok := zero.Lookup("x"); ok {
		t.Error("zero registry should be empty")
	}
}

func TestParseParams(t *testing.T) {
	params := map[string]string{"n": "12", "f": "1.5", "bad": "x"}
	var n iec.LINT = 7
	var f iec.REAL
	if err := ParseLINT(params, "missing", &n); err != nil || n != 7 {
		t.Fatalf("missing key changed value: %v %v", n, err)
	}
	if err := ParseLINT(params, "n", &n); err != nil || n != 12 {
		t.Fatalf("ParseLINT = %v %v", n, err)
	}
	if err := ParseREAL(params, "f", &f); err != nil || f != 1.5 {
		t.Fatalf("ParseREAL = %v %v", f, err)
	}
	if err := ParseREAL(params, "bad", &f); err == nil {
		t.Fatal("ParseREAL accepted a bad value")
	}
}

func Example() {
	RegisterProgramFactory("Blinker", func(params map[string]string) (func(time.Time), error) {
		var on iec.BOOL
		return func(time.Time) { on = !on }, nil
	})
	cfg, err := LoadConfigurationFromString(`
name: Demo
resource: MainCPU
  cycle: 10ms
  task: Blink
    priority: 1
    interval: 500ms
    program: Lamp1 Blinker
`)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, r := range cfg.Resources() {
		for _, t := range r.Tasks() {
			fmt.Println(r.Name, t.Name, t.Interval, len(t.Programs()))
		}
	}
	// Output: MainCPU Blink 500ms 1
}
