//go:build !tinygo && (linux || windows)

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

package main

import (
	"fmt"
	"time"

	"github.com/apiarytech/royaljelly/config"
)

// servoLoopLogic simulates a high-frequency control loop.
func servoLoopLogic(now time.Time) {
	// In a real application, this would read sensors and update motor positions.
	fmt.Printf("[%s] ServoLoop running on RealTimeCPU.\n", now.Format("15:04:05.000"))
}

// hmiScreenUpdateLogic simulates a lower-priority UI update task.
func hmiScreenUpdateLogic(now time.Time) {
	fmt.Printf("[%s] ---- HMIScreenUpdate running on HMI_CPU.\n", now.Format("15:04:05.000"))
}

func main() {
	// Register the program logic functions with the config loader.
	config.RegisterProgram("ServoLoop", servoLoopLogic)
	config.RegisterProgram("HMIScreenUpdate", hmiScreenUpdateLogic)

	fmt.Println("Loading multi-core configuration from 'config.txt'...")
	cfg, err := config.LoadConfigurationFromFile("config.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Configuration '%s' loaded successfully with %d resource(s).\n", cfg.Name, len(cfg.Resources))

	// Start all configured resources. The OS will handle pinning them to the specified cores.
	for _, res := range cfg.Resources {
		fmt.Printf("Starting resource '%s' with affinity for core %d...\n", res.Name, res.Affinity)
		res.Start()
	}

	fmt.Println("\nMulti-core simulation running for 4 seconds...")
	time.Sleep(4 * time.Second)

	fmt.Println("\nSimulation complete.")
}
