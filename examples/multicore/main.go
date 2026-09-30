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
	"context"
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
	config.RegisterProgramFactory("ServoLoop", func(params map[string]string) (func(time.Time), error) {
		return servoLoopLogic, nil
	})
	config.RegisterProgramFactory("HMIScreenUpdate", func(params map[string]string) (func(time.Time), error) {
		return hmiScreenUpdateLogic, nil
	})

	fmt.Println("Loading multi-core configuration from 'config.txt'...")
	cfg, err := config.LoadConfigurationFromFile("config.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Configuration %q loaded with %d resource(s).\n", cfg.Name, len(cfg.Resources()))

	// Run validates and starts every resource, then stops them all when the
	// context ends. Faults such as program panics are reported on standard error.
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	fmt.Println("\nSimulation running for 4 seconds...")
	if err := cfg.Run(ctx); err != nil {
		panic(err)
	}
	fmt.Println("\nSimulation complete.")
}
