# royaljelly

A Go library for writing PLC (Programmable Logic Controller) programs in Go, following the IEC 61131-3 standard.

`royaljelly` provides the IEC 61131-3 software model (configurations, resources, tasks and programs), the standard data types, function blocks and functions. Industrial automation engineers get familiar building blocks, and Go developers get ordinary, testable Go code.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Packages](#packages)
- [Usage](#usage)
  - [Quick start: build a configuration in code](#quick-start-build-a-configuration-in-code)
  - [Load the structure from a config file](#load-the-structure-from-a-config-file)
  - [Share data between resources](#share-data-between-resources)
  - [Faults, watchdog and metrics](#faults-watchdog-and-metrics)
  - [Data types](#data-types)
  - [Function blocks](#function-blocks)
  - [Standard functions](#standard-functions)
- [Timing model](#timing-model)
- [TinyGo support](#tinygo-support)
- [Upgrading from v0.0.5-beta1](#upgrading-from-v005-beta1)
- [Development](#development)
- [Licensing](#licensing)
- [Contributing](#contributing)

## Features

- **IEC 61131-3 software model.** A `Configuration` holds `Resource`s. Each resource runs a scheduler over prioritized `Task`s, and each task runs `Program`s.
- **Predictable scheduling.** Cyclic tasks run on a fixed grid with no accumulated drift. Missed runs are counted as overruns, and a per-task watchdog reports long runs.
- **Fault handling.** A panic in one program is recovered and reported, and the other programs keep running. Panics, overruns and watchdog trips reach a fault handler you provide.
- **Safe live changes.** Tasks and programs can be added or removed while a resource runs. The scheduler does not allocate memory per scan.
- **Multi-core and CPU affinity.** On Linux and Windows a resource can be pinned to a CPU core, which enables patterns such as redundant execution with a voter.
- **Thread-safe data exchange.** `vars.Shared` and `vars.ProcessImage` share values between resources and with the outside world.
- **IEC 61131-3 data types.** `BOOL`, `INT`, `REAL`, `TIME`, `DATE`, `TOD`, `DT`, `STRING`, `WSTRING` and the rest.
- **Standard function blocks.** Timers (`TP`, `TON`, `TOF`), counters (`CTU`, `CTD`, `CTUD`), edge detection (`R_TRIG`, `F_TRIG`), bistables (`SR`, `RS`), and `INTEGRAL`, `DERIVATIVE`, `HYSTERESIS` and `PID`.
- **Standard functions.** Numerical, arithmetic (including time arithmetic), selection, comparison, string, bitwise and type conversion functions.
- **TinyGo compatible.** The scheduler and the config loader work on microcontrollers.

## Installation

```bash
go get github.com/apiarytech/royaljelly
```

`royaljelly` requires Go 1.27 or later.

## Packages

| Package | Contents |
| --- | --- |
| `core` | Configuration, Resource, Task, Program, scheduler, faults |
| `config` | Text configuration loader and program factory registry |
| `vars` | `Shared[T]` values and the `%I`/`%Q`/`%M` process image |
| `iec` | IEC 61131-3 data types and constants |
| `convert` | Reflection-free conversions between IEC types |
| `fb/timers`, `fb/counters`, `fb/triggers`, `fb/calculus` | Standard function blocks |
| `std/...` | Standard functions: arithmetic, bitwise, comparison, conversion, math, numerical, selection, strings, time |

## Usage

Each snippet below has a runnable counterpart in [examples/readme/readme_test.go](examples/readme/readme_test.go), so API drift breaks the build. The [examples](examples) directory holds complete programs.

### Quick start: build a configuration in code

Building the hierarchy in code needs no config package and gives the smallest binary. Keep each program's state in a struct, which plays the role of an IEC program's local variables.

```go
type Blinker struct {
	Lamp iec.BOOL
	Runs int
}
blinker := &Blinker{}

task := core.NewTask("BlinkTask", core.CyclicTask, 1, 20*time.Millisecond).
	WithProgram(&core.Program{
		Name: "Blink",
		Logic: func(now time.Time) {
			blinker.Lamp = !blinker.Lamp
			blinker.Runs++
		},
	})
resource := (&core.Resource{Name: "MainCPU", Cycle: 5 * time.Millisecond}).WithTask(task)
cfg := (&core.Configuration{Name: "Demo"}).WithResource(resource)

// Run validates the configuration, starts every resource, and stops them
// all when the context ends.
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()
if err := cfg.Run(ctx); err != nil {
	fmt.Println("start failed:", err)
	return
}
```

You can also control each resource yourself. `Resource.Start` returns an error if validation fails or the CPU affinity cannot be applied. `Configuration.Start(ctx)` starts everything without blocking, and `Stop` shuts it down.

Validation rejects these mistakes before anything runs:

- A cyclic task with no interval, or with an interval not longer than the resource cycle.
- Two tasks with the same name or priority on one resource.
- Two programs with the same name in one task.
- Two resources with the same name in one configuration.

### Load the structure from a config file

The `config` package builds the same hierarchy from a small text format. Register a factory for each program type, and the loader calls it for each instance.

```text
name: EncapsulationExample
resource: MainCPU
  cycle: 50ms
  task: TaskA
    type: Cyclic          # or EventDriven
    priority: 1
    interval: 250ms
    watchdog: 100ms       # optional
    program: CounterA CounterProgram
      param: initial_value 100
```

```go
config.RegisterProgramFactory("CounterProgram", func(params map[string]string) (func(time.Time), error) {
	var count iec.LINT
	if err := config.ParseLINT(params, "initial_value", &count); err != nil {
		return nil, err
	}
	return func(now time.Time) { count++ }, nil
})

cfg, err := config.LoadConfigurationFromFile("config.txt")
if err != nil {
	panic(err) // Errors name the line, e.g. "config.txt: line 7: unknown task type 'Sometimes'".
}
```

Indent each level with a tab or a fixed number of spaces, and stay consistent within a file. Text after a `#` that follows whitespace is a comment. `LoadConfigurationFromString` parses an embedded configuration. To keep separate sets of program types, create a `config.Loader` with its own `config.Registry`.

### Share data between resources

Programs on one resource run one after another, so they can share plain Go variables. Programs on different resources, and code outside the scheduler such as an HMI or I/O driver, run at the same time. They must exchange data through a synchronized container.

```go
var setpoint vars.Shared[iec.REAL]
setpoint.Store(21.5) // e.g. written by an HMI goroutine

var pi vars.ProcessImage
pi.Write(func(img *vars.Image) { img.I.B[0] = true }) // e.g. an input driver sets %IX0

var start iec.BOOL
pi.Read(func(img *vars.Image) { start = img.I.B[0] })
```

See [examples/redundancy](examples/redundancy) for two cores running the same logic with a voter.

### Faults, watchdog and metrics

```go
task := core.NewTask("Control", core.CyclicTask, 1, 10*time.Millisecond)
task.Watchdog = 5 * time.Millisecond // Report runs longer than 5ms.

resource := &core.Resource{
	Name: "MainCPU",
	OnFault: func(f core.Fault) {
		// f.Kind is FaultPanic, FaultOverrun or FaultWatchdog.
		log.Println(f)
	},
}
```

The fault handler can also be set once on the `Configuration`. Without a handler, panics and watchdog trips are written to standard error, and overruns are only counted. `Task.Stats()` returns run count, overruns, watchdog trips, execution time, cycle time and drift.

### Data types

The `iec` package defines Go types that map to IEC 61131-3 types.

```go
var myBool iec.BOOL = true
var myInt iec.INT = 123
var myReal iec.REAL = 45.67
var myTime iec.TIME = iec.TIME(10 * time.Second)
var myText iec.WSTRING = "Grüße"
```

### Function blocks

Function blocks are structs with input and output fields and an execute method. Timers take the scan time, so they are deterministic in tests.

```go
ton := timers.TON{PT: iec.TIME(5 * time.Second)}
now := time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)
for scan := 0; scan <= 6; scan++ {
	ton.IN = scan >= 1 // Input switches on at scan 1.
	ton.Execute(now)
	fmt.Printf("scan %d: Q=%v ET=%v\n", scan, ton.Q, time.Duration(ton.ET))
	now = now.Add(time.Second)
}
// Q turns on at scan 6, five seconds after IN.
```

Every function block has a runnable example in its package documentation, for instance `go doc github.com/apiarytech/royaljelly/fb/counters`.

### Standard functions

```go
sum := arithmetic.ADD(iec.REAL(10), iec.REAL(5.5))                     // 15.5
length := strings.LEN("Hello, RoyalJelly!")                             // 18
largest, err := selection.MAX(iec.LINT(100), iec.LINT(50), iec.LINT(120)) // 120
```

Generic functions take arguments of one type. Convert mixed types first with the `std/conversion` functions or with `convert.ConvertTo`.

## Timing model

Each resource runs one scheduler goroutine that wakes every `Cycle`. On each wake-up it runs the due tasks in priority order, lowest number first. Priority sets the order within a scan. It never interrupts a running task, so a slow program delays everything else on its resource. Put timing-critical logic on its own resource.

[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) describes the scheduler, overruns, the watchdog and the concurrency rules in detail.

## TinyGo support

The scheduler runs under TinyGo's cooperative, single-threaded runtime, and the config loader uses a hand-written parser with no reflection. CPU affinity is not available under TinyGo, and `Start` returns `core.ErrAffinityUnsupported` if a resource requests it. CI builds and tests the library with TinyGo.

## Upgrading from v0.0.5-beta1

This release splits the former all-in-one `core` package into focused packages, and changes some APIs to fix data races and make configuration safer. It also requires Go 1.27 or later.

**Imports**

- **Types and constants moved to `iec`.** `BOOL`, `INT`, `LINT`, `REAL`, `TIME` and the other IEC types, and constants such as `MAXINT`, moved from `core` to `iec`. Replace `core.LINT` with `iec.LINT`, or change a dot import of `core` to one of `iec`.
- **Conversions moved to `convert`.** `AnyToLINT`, `ConvertTo` and the other helpers moved from `core` to `convert`. The `core` names still work but are deprecated.
- **Address tables moved to `vars`.** The `Addresses` type and the `I`, `Q` and `M` globals moved from `core` to `vars`. The globals are deprecated because they are unsynchronized. Use a `vars.ProcessImage`.

**Scheduler API**

- **Start returns an error.** `Resource.Start()` now validates the setup and reports problems instead of running a broken configuration. It rejects a cyclic task with no interval or with an interval not longer than the cycle. It also rejects duplicate task names or priorities in a resource, and duplicate program names in a task.
- **Collections are methods.** Replace the `cfg.Resources`, `resource.Tasks` and `task.Programs` fields with `cfg.Resources()`, `resource.Tasks()` and `task.Programs()`. Adding items works as before, with `WithResource`, `AddTask` and `AddProgram`. Configuration also gains an `AddResource` method.
- **Task state is read through methods.** Replace `task.Enabled` with `task.IsEnabled()`, and use `Enable()` and `Disable()` to change it. Tasks remain enabled by default.
- **Metrics are methods.** Replace the `ExecutionTime`, `CycleTime` and `Drift` fields with methods of the same names, or use `task.Stats()`.
- **Removed stub.** The no-op `Resource.WithResource` method was removed.

**Behavior**

- **Drift no longer accumulates.** Cyclic tasks run on a fixed schedule. A task that falls behind skips the missed runs and counts them as overruns instead of shifting every later run.
- **Panics are reported, not printed.** A program panic is still recovered, but it now goes to the fault handler instead of standard output. Without a handler it is written to standard error.

**Types**

- **Wide strings.** `iec.WSTRING` is now a string type, not a rune. Use `iec.WCHAR` for a single character.
- **WORD area in the address tables.** `Addresses.W` now holds WORD values, so `%IW`, `%QW` and `%MW` have their own area. The wide-string array moved from `W` to `WS`.

**New in this release**

- **Config loader.** The `config` package builds a configuration from a text file or string.
- **Configuration lifecycle.** `Configuration.Start(ctx)`, `Stop()` and `Run(ctx)` start and stop every resource together.
- **Faults and watchdog.** Set a fault handler on a resource or configuration, and set `Task.Watchdog` to detect long runs.
- **Shared data.** `vars.Shared[T]` and `vars.ProcessImage` exchange data safely between resources.
- **CPU affinity.** On Linux and Windows, `Resource.Affinity` pins a resource to a CPU core.
- **TinyGo support.** The scheduler and the config loader build and pass their tests under TinyGo.

## Development

```bash
go test ./...                       # unit tests and README examples
go test -race ./...                 # data race detection (needs cgo)
go test -run XXX -bench . ./core    # scheduler overhead and jitter benchmarks
tinygo test ./...                   # TinyGo build and tests
```

## Licensing

This project is offered under a dual-license model. You may use it under either the GNU General Public License version 2 (GPLv2) or a commercial license.

*   **GPLv2:** If you are developing open-source software, you can use this library under the terms of the GPLv2. The full license text is available in the `gpl-2.0.md` file.
*   **Commercial License:** If you intend to use this library in a proprietary, closed-source application or product, a commercial license is required.

For more details on both licensing options, please see the `LICENSE.md` file.

## Contributing

Contributions to `royaljelly` are welcome! Please feel free to:
- Fork the repository.
- Submit issues for bugs or feature requests.
- Submit pull requests with improvements, bug fixes, or new IEC 61131-3 compliant implementations.

Please ensure that your contributions follow the existing code style, include tests, and pass `go test -race ./...`.

Thank you for your interest in `royaljelly`!
