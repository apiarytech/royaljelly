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
	fmt.Printf("Configuration %q loaded with %d resource(s).\n", cfg.Name, len(cfg.Resources()))

	// Run validates and starts every resource, then stops them all when the
	// context ends. Faults such as program panics are reported on standard error.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("\nSimulation running for 5 seconds...")
	if err := cfg.Run(ctx); err != nil {
		panic(err)
	}
	fmt.Println("\nSimulation complete.")
}
