package main

import (
	"fmt"
	"time"

	plc "github.com/apiarytech/royaljelly"
)

// This example simulates a 4-way traffic light controller using a state machine.
//
// Sequence:
// 1. North/South lights are Green, East/West are Red.
// 2. North/South lights turn Yellow, East/West remain Red.
// 3. North/South lights turn Red, East/West turn Green.
// 4. North/South remain Red, East/West turn Yellow.
// 5. The cycle repeats.

const (
	stNSGreen  plc.DINT = iota // North-South Green, East-West Red
	stNSYellow                 // North-South Yellow, East-West Red
	stEWGreen                  // North-South Red, East-West Green
	stEWYellow                 // North-South Red, East-West Yellow
)

type Lights struct {
	Red, Yellow, Green plc.BOOL
}

func main() {
	// --- Create the PLC Configuration ---
	config := &plc.Configuration{Name: "TrafficLightController"}
	resource := &plc.Resource{Name: "MainCPU", Cycle: 100 * time.Millisecond}
	task := plc.NewTask("TrafficLightTask", plc.CyclicTask, 1, 1*time.Second)

	// --- Program-Scoped Variables ---
	// These variables are captured by the program's logic closure.
	var (
		nsLights, ewLights Lights
		state              plc.DINT = stNSGreen
		sequenceTimer      plc.TON
	)
	sequenceTimer.INIT()

	// --- Define the Program Logic ---
	trafficProgram := &plc.Program{
		Name: "TrafficLightLogic",
		Logic: func(now time.Time) {
			// --- State Machine Logic ---
			sequenceTimer.IN = true // Timer is always enabled in this simple model

			switch state {
			case stNSGreen:
				nsLights = Lights{Green: true}
				ewLights = Lights{Red: true}
				sequenceTimer.PT = plc.TIME(10 * time.Second) // Green light for 10s
				if sequenceTimer.Q {
					state = stNSYellow
					sequenceTimer.IN = false // Reset timer by toggling IN
				}

			case stNSYellow:
				nsLights = Lights{Yellow: true}
				ewLights = Lights{Red: true}
				sequenceTimer.PT = plc.TIME(3 * time.Second) // Yellow for 3s
				if sequenceTimer.Q {
					state = stEWGreen
					sequenceTimer.IN = false
				}

			case stEWGreen:
				nsLights = Lights{Red: true}
				ewLights = Lights{Green: true}
				sequenceTimer.PT = plc.TIME(10 * time.Second) // Green for 10s
				if sequenceTimer.Q {
					state = stEWYellow
					sequenceTimer.IN = false
				}

			case stEWYellow:
				nsLights = Lights{Red: true}
				ewLights = Lights{Yellow: true}
				sequenceTimer.PT = plc.TIME(3 * time.Second) // Yellow for 3s
				if sequenceTimer.Q {
					state = stNSGreen // Cycle back
					sequenceTimer.IN = false
				}
			}

			// Execute timer logic for the current scan
			sequenceTimer.Execute(now)

			fmt.Printf("State=%d | NS: [R:%v Y:%v G:%v] | EW: [R:%v Y:%v G:%v] | Timer ET: %-4v\n",
				state, nsLights.Red, nsLights.Yellow, nsLights.Green, ewLights.Red, ewLights.Yellow, ewLights.Green, sequenceTimer.ET)
		},
	}

	// --- Assemble the PLC Structure ---
	task.AddProgram(trafficProgram)
	resource.AddTask(task)
	config.Resources = append(config.Resources, resource)

	// --- Start the PLC ---
	fmt.Println("--- Traffic Light Simulation ---")
	resource.Start()

	// Keep the simulation running for a while.
	time.Sleep(40 * time.Second)
	resource.Stop()
	fmt.Println("--- Simulation Complete ---")
}
