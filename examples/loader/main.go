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

func main() {
	// Register all program logic functions
	config.RegisterProgram("HighFreqCounter", HighFreqLogicFunc)
	config.RegisterProgram("LowFreqCounter", LowFreqLogicFunc)

	fmt.Println("Loading configuration from 'examples/loader/config.txt'...")
	// Load the entire structure from a file
	// Note: The path is relative to the project root where you run `go run`.
	cfg, err := config.LoadConfigurationFromFile("config.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Configuration '%s' loaded successfully with %d resource(s).\n", cfg.Name, len(cfg.Resources))

	// Find and start resources
	for _, res := range cfg.Resources {
		fmt.Printf("Starting resource '%s'...\n", res.Name)
		res.Start()
	}

	// Let the simulation run for a few seconds
	fmt.Println("\nSimulation running for 5 seconds...")
	time.Sleep(5 * time.Second)

	// Stop all resources
	for _, res := range cfg.Resources {
		fmt.Printf("Stopping resource '%s'...\n", res.Name)
		res.Stop()
	}

	fmt.Println("\nSimulation complete.")
}
