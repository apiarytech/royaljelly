# Upgrading royaljelly

## From v0.0.5-beta1 to v0.1.0-beta1 or later

Version v0.1.0-beta1 split the former all-in-one `core` package into focused packages, and changed some APIs to fix data races and make configuration safer. It also requires Go 1.27 or later. Later releases add features without further breaking changes, so these notes cover an upgrade from v0.0.5-beta1 or earlier to any current version.

Versions v0.0.1-alpha through v0.0.5-beta1 are retracted in `go.mod`. The `go` command does not select them automatically, and `go list -m -u all` flags any module that still uses them.

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

**New in v0.1.0-beta1**

- **Config loader.** The `config` package builds a configuration from a text file or string.
- **Configuration lifecycle.** `Configuration.Start(ctx)`, `Stop()` and `Run(ctx)` start and stop every resource together.
- **Faults and watchdog.** Set a fault handler on a resource or configuration, and set `Task.Watchdog` to detect long runs.
- **Shared data.** `vars.Shared[T]` and `vars.ProcessImage` exchange data safely between resources.
- **CPU affinity.** On Linux and Windows, `Resource.Affinity` pins a resource to a CPU core.
- **TinyGo support.** The scheduler and the config loader build and pass their tests under TinyGo.
