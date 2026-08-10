# royaljelly

A Go library for writing PLC (Programmable Logic Controller) programs using Go syntax, designed with adherence to the IEC 61131-3 standard.

`royaljelly` provides data types, function blocks (FBs), and functions commonly found in IEC 61131-3 compliant PLC programming environments, enabling Go developers to implement control logic familiar to industrial automation engineers.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Data Types](#data-types)
  - [Function Blocks](#function-blocks)
  - [Standard Functions](#standard-functions)
- [TinyGo Support](#tinygo-support)
- [Licensing](#licensing)
- [Contributing](#contributing)

## Features

`royaljelly` aims to translate the core concepts of IEC 61131-3 into idiomatic Go, offering:
- **IEC 61131-3 Software Model**: A hierarchical structure (`Configuration` -> `Resource` -> `Task` -> `Program`) for organizing and scheduling control logic, managed by a priority-based, TinyGo-compatible scheduler.
   - **Multi-Core Support (Standard Go)**: The scheduler automatically utilizes all available CPU cores (`runtime.GOMAXPROCS`). Resources can be pinned to specific CPU cores using the `Affinity` property for predictable, real-time performance.
   - **TinyGo Compatibility**: The scheduler is fully compatible with TinyGo, running efficiently in a single-threaded, cooperative multitasking environment.
- **IEC 61131-3 Data Types**: Go types representing standard IEC types like `BOOL`, `INT`, `REAL`, `TIME`, `DATE`, `TOD`, `DT`, etc.
- **Standard Function Blocks**: Implementations of common FBs as Go structs with `INIT` and `Execute` methods, including:
    - **Timers**: `TP` (Pulse Timer), `TON` (On-Delay Timer), `TOF` (Off-Delay Timer) (IEC 61131-3, Table 35)
    - **Counters**: `CTU` (Count Up), `CTD` (Count Down), `CTUD` (Count Up/Down) (IEC 61131-3, Table 36)
    - **Edge Detection**: `R_TRIG` (Rising Edge Trigger), `F_TRIG` (Falling Edge Trigger) (IEC 61131-3, Table 35)
- **Standard Functions**: Implementations of common mathematical, comparison, string, and bitwise functions (IEC 61131-3, Tables 28-32).
    - **Numerical**: `ABS`, `SQRT`, `LN`, `LOG`, `EXP`, `SIN`, `COS`, `TAN`, `ASIN`, `ACOS`, `ATAN`, `EXPT`, `TRUNC`.
    - **Arithmetic**: `ADD`, `SUB`, `MUL`, `DIV`, `MOD`, `MOVE`, and specific time arithmetic functions like `ADD_TIME`, `SUB_TOD`, `MUL_TIME`, etc.
    - **Selection**: `SEL`, `MAX`, `MIN`, `LIMIT`, `MUX`.
    - **String Manipulation**: `LEN`, `LEFT`, `RIGHT`, `MID`, `CONCAT`, `INSERT`, `DELETE`, `FIND`, `REPLACE`.
    - **Bitwise Operators**: `AND`, `OR`, `XOR`, `NOT`, `SHL`, `SHR`, `ROL`, `ROR`.
- **Type Conversion**: Functions to convert between IEC 61131-3 data types, handling explicit and implicit conversions.

## Installation

To use `royaljelly` in your Go project, simply run:

```bash
go get github.com/apiarytech/royaljelly
```

## Usage

### Structuring a PLC Application

`royaljelly` provides a framework to structure your application according to the IEC 61131-3 software model. This creates a clear hierarchy for managing your control logic, orchestrated by a resource-level scheduler.

*   **Configuration**: The top-level object representing the entire PLC system.
*   **Resource**: Represents a processing unit (like a CPU) and runs a priority-based task scheduler.
*   **Task**: Controls the execution properties of programs (e.g., cyclic interval or event-driven).
*   **Program**: Contains your control logic, written as a Go closure.

The following snippet, taken from the `4-way traffic light` example, demonstrates how to assemble this structure:

```go
package main

import (
	"time"
	. "github.com/apiarytech/royaljelly/core"
)

func main() {
	// 1. Define the hierarchy
	config := &Configuration{Name: "TrafficLightController"}
	resource := &Resource{Name: "MainCPU", Cycle: 100 * time.Millisecond}
	task := NewTask("TrafficLightTask", CyclicTask, 1, 1*time.Second)

	// 2. Define the program logic as a closure
	trafficProgram := &plc.Program{
		Name: "TrafficLightLogic",
		Logic: func(now time.Time) {
			// ... state machine logic for the traffic light ...
		},
	}

	// 3. Assemble the structure and start the resource
	task.AddProgram(trafficProgram)
	resource.AddTask(task)
	resource.Start()

	// Keep the simulation running
	time.Sleep(40 * time.Second)
	resource.Stop()
}
```

A complete, runnable version of this is available in the `examples/4-way_traffic_light` directory.

### Data Types

`royaljelly` defines Go types that map directly to IEC 61131-3 data types.

```go
package main

import (
	"fmt"
	"time"
	. "github.com/apiarytech/royaljelly/core"
)

func main() {
	var myBool BOOL = true
	var myInt INT = 123
	var myReal REAL = 45.67
	var myTime TIME = TIME(10 * time.Second)

	fmt.Printf("BOOL: %v\n", myBool)
	fmt.Printf("INT: %v\n", myInt)
	fmt.Printf("REAL: %v\n", myReal)
	fmt.Printf("TIME: %v\n", myTime)
}
```

### Function Blocks

Function blocks are implemented as Go structs with an `INIT()` method for initialization and an `Execute()` (or similar named) method for their logic.

```go
package main

import (
	"fmt"
	"time"
	. "github.com/apiarytech/royaljelly/core"
	"github.com/apiarytech/royaljelly/fb/timers"
)

func main() {
	// Example: TON (On-Delay Timer)
	ton1 := timers.TON{}
	ton1.INIT()
	ton1.PT = TIME(5 * time.Second) // Preset Time: 5 seconds

	now := time.Now()

	// Simulate PLC scan cycles
	fmt.Println("TON Simulation Start")
	for i := 0; i < 10; i++ {
		now = now.Add(1 * time.Second) // Advance time by 1 second per scan
		if i == 1 {
			ton1.IN = true // Set input IN to true after 1 second
		}

		ton1.Execute(now)

		fmt.Printf("Scan %d: IN=%v, Q=%v, ET=%v\n", i, ton1.IN, ton1.Q, ton1.ET)
	}
	fmt.Println("TON Simulation End")
}
```

Here is another example using the `TP` (Pulse Timer) function block.

```go
package main

import (
	"fmt"
	"time"
	. "github.com/apiarytech/royaljelly/core"
	"github.com/apiarytech/royaljelly/fb/timers"
)

func main() {
	// Example: TP (Pulse Timer)
	tp1 := timers.TP{}
	tp1.INIT()
	tp1.PT = TIME(3 * time.Second) // Preset Time: 3-second pulse

	now := time.Now()
	fmt.Println("TP Simulation Start")

	// Trigger the pulse
	tp1.IN = true
	tp1.Execute(now)
	fmt.Printf("Time=0s: IN=%v, Q=%v, ET=%v\n", tp1.IN, tp1.Q, tp1.ET)

	// After 2 seconds, the pulse is still active
	now = now.Add(2 * time.Second)
	tp1.Execute(now)
	fmt.Printf("Time=2s: IN=%v, Q=%v, ET=%v\n", tp1.IN, tp1.Q, tp1.ET)

	// After 4 seconds, the pulse has finished, even if IN is still true
	now = now.Add(2 * time.Second)
	tp1.Execute(now)
	fmt.Printf("Time=4s: IN=%v, Q=%v, ET=%v\n", tp1.IN, tp1.Q, tp1.ET)
	fmt.Println("TP Simulation End")
}
```

### Standard Functions

Standard functions like `ADD`, `MUL`, `SQRT`, `LEN`, `AND`, etc., are provided to operate on `royaljelly` types.

```go
package main

import (
	"fmt"
	. "github.com/apiarytech/royaljelly/core"
	. "github.com/apiarytech/royaljelly/std/arithmetic"
	. "github.com/apiarytech/royaljelly/std/selection"
	. "github.com/apiarytech/royaljelly/std/strings"
)

func main() {
	result := ADD(INT(10), REAL(5.5))
	fmt.Printf("ADD(10, 5.5) = %v (Type: %T)\n", result, result) // Output: 15.5 (Type: core.REAL)

	strLen := LEN("Hello, RoyalJelly!")
	fmt.Printf("LEN(\"Hello, RoyalJelly!\") = %v\n", strLen) // Output: 18

	maxVal := MAX(LINT(100), DINT(50), LINT(120))
	fmt.Printf("MAX(100, 50, 120) = %v\n", maxVal) // Output: 120
}
```

## Licensing

This project is offered under a dual-license model. You have the choice of using it under either the GNU General Public License version 2 (GPLv2) or a commercial license.

*   **GPLv2:** If you are developing open-source software, you can use this library under the terms of the GPLv2. The full license text is available in the `gpl-2.0.md` file.
*   **Commercial License:** If you intend to use this library in a proprietary, closed-source application or product, a commercial license is required.

For more details on both licensing options, please see the `LICENSE.md` file.

## Contributing

Contributions to `royaljelly` are welcome! Please feel free to:
- Fork the repository.
- Submit issues for bugs or feature requests.
- Submit pull requests with improvements, bug fixes, or new IEC 61131-3 compliant implementations.

Please ensure that your contributions adhere to the existing code style and include appropriate tests.

Thank you for your interest in `royaljelly`!
