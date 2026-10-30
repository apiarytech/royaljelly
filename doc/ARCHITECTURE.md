# Architecture

This page explains how royaljelly is organized and how its scheduler behaves. Read it before writing timing-sensitive programs or contributing to `core`.

## Package layers

Dependencies point downward only. The scheduler never depends on the function libraries, and the libraries never depend on the scheduler.

```
examples, user programs
        │
   config ─────────► core            (software model and scheduler)
        │              │
        │           convert ◄────── std/*, fb/*   (functions and function blocks)
        │              │
        └──────────► iec ◄───────── vars          (types; shared data)
```

- `iec` defines the IEC 61131-3 data types and constants. It imports nothing from the project.
- `convert` holds reflection-free conversions between those types.
- `std/*` and `fb/*` are the standard functions and function blocks. They are plain Go with no goroutines and no global state.
- `vars` provides synchronized containers for data that crosses goroutines.
- `core` is the runtime: configurations, resources, tasks, programs and the scheduler.
- `config` parses the text format into a `core.Configuration`.

## Software model

| IEC 61131-3 | royaljelly | Role |
| --- | --- | --- |
| Configuration | `core.Configuration` | The whole PLC. Starts and stops resources together. |
| Resource | `core.Resource` | One processing unit. Runs one scheduler goroutine. |
| Task | `core.Task` | Decides when programs run: cyclic or event-driven. |
| Program | `core.Program` | A named `Logic(now)` function, usually a method on a state struct. |
| Function block | Structs in `fb/*` | Reusable logic with its own state, called from programs. |

## The scheduler

Each resource owns one goroutine driven by a ticker with period `Resource.Cycle`. Each tick is a scan. A scan does this:

1. Take the current task list. The list is copy-on-write, so this is a pointer read under a short lock.
2. Visit tasks in priority order, lowest number first.
3. For each task, decide whether it is due and advance its schedule.
4. Run each program of a due task, recovering any panic.
5. Record execution time, cycle time and drift for the task.

A scan does not allocate memory. `TestScanDoesNotAllocate` guards this, and `BenchmarkScan` measures the overhead.

### Cyclic tasks

A cyclic task keeps a next-run time. The first scan after the task starts or is re-enabled runs it and sets the next run to one interval later. After that, the task runs at the first scan at or after its next-run time, and the next-run time advances by exactly one interval. Late scans therefore never push the whole schedule later.

Example with a 100ms interval and a 30ms cycle:

| Scan time (ms) | 0 | 30 | 60 | 90 | 120 | 150 | 180 | 210 | … | 300 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Runs? | yes | | | | yes | | | yes | | yes |
| Drift (ms) | 0 | | | | 20 | | | 10 | | 0 |

The effective interval is rounded up to scan boundaries. Choose a cycle that divides your task intervals evenly to get zero drift.

### Overruns

If a scan arrives after more than one interval has passed, the task runs once. The runs it missed are skipped, not replayed. The number skipped is added to `Task.Overruns` and reported as a `FaultOverrun`. The schedule then continues on its original grid.

Overruns happen when programs on the resource take longer than the interval, or when the operating system delays the scheduler goroutine.

### Event-driven tasks

`Task.Trigger` queues a run, up to 10 pending triggers. Each scan runs the task once if a trigger is pending. The trigger takes effect at the next scan, so the reaction time is at most one cycle plus the time for higher-priority tasks.

### Priority means order, not preemption

Priority orders the tasks that are due in the same scan. A running task is never interrupted by a higher-priority one, because Go cannot preempt a function from the outside. A slow program on a resource delays every task on that resource.

To isolate timing-critical logic:

- Put it on its own resource, so it has its own scheduler goroutine.
- On Linux and Windows, set `Resource.Affinity` to pin that goroutine's OS thread to a CPU core.
- Exchange data with other resources through `vars.Shared` or `vars.ProcessImage`.

## Faults

The scheduler reports problems as a `core.Fault` value.

| Kind | Raised when | Default behavior |
| --- | --- | --- |
| `FaultPanic` | A program panics. The panic is recovered and the remaining programs of the task still run. | Printed to standard error |
| `FaultOverrun` | A cyclic task skipped scheduled runs. | Counted only |
| `FaultWatchdog` | A task run lasted longer than `Task.Watchdog`. | Printed to standard error |

A resource sends faults to its `OnFault` handler. If it has none, `Configuration.Start` gives it the configuration's handler. Each fault wraps a sentinel error, so `errors.Is(fault, core.ErrOverrun)` works.

The watchdog is a timer started before a task runs and stopped after. If it fires, the fault is reported while the program is still running. It cannot abort the program. Treat it as an alarm and act on it in the handler, for example by setting outputs to a safe state through a separate path.

Handlers can run on the scheduler goroutine or on a timer goroutine. They must be safe for concurrent use and should return quickly.

## Concurrency rules

1. **Configure first, then start.** Exported fields such as names, `Cycle`, `Affinity`, `Interval`, `Priority`, `Watchdog` and `OnFault` are configuration. Set them before the resource starts and don't change them afterwards.
2. **Methods are safe for concurrent use.** You can add or remove tasks and programs, enable or disable tasks, trigger events, and read metrics while the scheduler runs. Changes take effect at the next scan. A scan already in progress finishes with the list it started with.
3. **One resource, sequential programs.** Programs on the same resource never run at the same time, so they can share plain variables.
4. **Across resources, synchronize.** Anything read or written by two resources, or by a resource and another goroutine, must go through `vars.Shared`, `vars.ProcessImage` or your own synchronization.
5. **Don't stop from inside.** `Stop` waits for the scheduler goroutine to exit, so calling it from a program on the same resource deadlocks. Cancel a context or signal another goroutine instead.

Run `go test -race ./...` to check your own programs against these rules.

## TinyGo

Under TinyGo all goroutines share one thread and switch cooperatively. The scheduler code is the same. `pinThread` returns `ErrAffinityUnsupported`, so a resource with `Affinity` set fails to start. The `config` parser uses only `bufio`, `strings` and `strconv`, so it works on microcontrollers.

## Adding to the library

- **A function block** is a struct with IEC-named input and output fields and an execute method. It should take `now time.Time` if it measures time, so tests can drive it with fixed times. Add a unit test and an `Example` with checked output.
- **A standard function** goes in the matching `std` package. Use the generic constraints in `iec` rather than `any`, and return an error for invalid input instead of panicking.
- **Scheduler changes** must keep `TestScanDoesNotAllocate` passing and must pass `go test -race ./core`.
