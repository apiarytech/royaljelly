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
	// Create instances of our program logic structs, which are defined in other files.
	hfc := &HighFreqCounter{}
	lfc := &LowFreqCounter{}

	// 1. Register factories for our program types.
	// Since these simple programs don't have parameters, the factory can
	// just return the Logic method from our instances.
	config.RegisterProgramFactory("HighFreqProgram", func(params map[string]string) (func(time.Time), error) {
		return hfc.Logic, nil
	})
	config.RegisterProgramFactory("LowFreqProgram", func(params map[string]string) (func(time.Time), error) {
		return lfc.Logic, nil
	})

	fmt.Println("Loading configuration from 'config.txt'...")
	// 2. Load the entire structure from the file.
	cfg, err := config.LoadConfigurationFromFile("config.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Configuration '%s' loaded successfully with %d resource(s).\n", cfg.Name, len(cfg.Resources))

	// 3. Find and start resources
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
