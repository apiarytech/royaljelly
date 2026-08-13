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

// tickerLogic is a simple program that prints a message every time it runs.
func tickerLogic(now time.Time) {
	fmt.Printf("[%s] TickerTask is running.\n", now.Format("15:04:05.000"))
}

func main() {
	// Register the program logic with the config loader.
	config.RegisterProgram("TickerProgram", tickerLogic)

	fmt.Println("Loading configuration from 'config.txt'...")
	cfg, err := config.LoadConfigurationFromFile("config.txt")
	if err != nil {
		panic(err)
	}

	// Find the resource we want to control.
	mainCPU := cfg.FindResource("MainCPU")
	if mainCPU == nil {
		panic("Could not find resource 'MainCPU' in configuration.")
	}

	// Start the resource.
	fmt.Printf("Starting resource '%s'...\n", mainCPU.Name)
	mainCPU.Start()

	fmt.Println("\n[1] Running normally for ~1 second...")
	time.Sleep(1 * time.Second)

	fmt.Println("\n[2] Pausing resource. No tasks should run now.")
	mainCPU.Pause()
	time.Sleep(1 * time.Second) // Wait while paused. You should see no output here.

	fmt.Println("\n[3] Resuming resource. Tasks should run again.")
	mainCPU.Resume()
	time.Sleep(1 * time.Second) // Wait while resumed.

	mainCPU.Stop()
	fmt.Println("\nSimulation complete.")
}
