package main

import (
	"fmt"
	"time"

	. "github.com/apiarytech/royaljelly/fb/timers"
	. "github.com/apiarytech/royaljelly/iec"
)

const (
	stNSGreen  DINT = iota // North-South Green, East-West Red
	stNSYellow             // North-South Yellow, East-West Red
	stEWGreen              // North-South Red, East-West Green
	stEWYellow             // North-South Red, East-West Yellow
)

type Lights struct {
	Red, Yellow, Green BOOL
}

// TrafficLightProgram encapsulates the state and logic for the traffic light controller.
type TrafficLightProgram struct {
	// State variables
	nsLights, ewLights Lights
	state              DINT
	sequenceTimer      TON
}

// Init initializes the program's state.
func (p *TrafficLightProgram) Init() {
	p.state = stNSGreen
	p.sequenceTimer.INIT()
}

// Logic contains the state machine for the traffic light. This method will be assigned
// to the Program's Logic field.
func (p *TrafficLightProgram) Logic(now time.Time) {
	// --- State Machine Logic ---
	p.sequenceTimer.IN = true // Timer is always enabled in this simple model

	switch p.state {
	case stNSGreen:
		p.nsLights = Lights{Green: true}
		p.ewLights = Lights{Red: true}
		p.sequenceTimer.PT = TIME(10 * time.Second) // Green light for 10s
		if p.sequenceTimer.Q {
			p.state = stNSYellow
			p.sequenceTimer.IN = false // Reset timer by toggling IN
		}

	case stNSYellow:
		p.nsLights = Lights{Yellow: true}
		p.ewLights = Lights{Red: true}
		p.sequenceTimer.PT = TIME(3 * time.Second) // Yellow for 3s
		if p.sequenceTimer.Q {
			p.state = stEWGreen
			p.sequenceTimer.IN = false
		}

	case stEWGreen:
		p.nsLights = Lights{Red: true}
		p.ewLights = Lights{Green: true}
		p.sequenceTimer.PT = TIME(10 * time.Second) // Green for 10s
		if p.sequenceTimer.Q {
			p.state = stEWYellow
			p.sequenceTimer.IN = false
		}

	case stEWYellow:
		p.nsLights = Lights{Red: true}
		p.ewLights = Lights{Yellow: true}
		p.sequenceTimer.PT = TIME(3 * time.Second) // Yellow for 3s
		if p.sequenceTimer.Q {
			p.state = stNSGreen // Cycle back
			p.sequenceTimer.IN = false
		}
	}

	// Execute timer logic for the current scan
	p.sequenceTimer.Execute(now)

	fmt.Printf("State=%d | NS: [R:%v Y:%v G:%v] | EW: [R:%v Y:%v G:%v] | Timer ET: %.1fs\n",
		p.state, p.nsLights.Red, p.nsLights.Yellow, p.nsLights.Green, p.ewLights.Red, p.ewLights.Yellow, p.ewLights.Green, time.Duration(p.sequenceTimer.ET).Seconds())
}
