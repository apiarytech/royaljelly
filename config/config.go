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

package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/apiarytech/royaljelly/core"
)

// ProgramRegistry maps string names from the config file to actual Go functions.
var ProgramRegistry = make(map[string]func(time.Time))

// RegisterProgram adds a program's logic to the registry so it can be found by the loader.
func RegisterProgram(name string, logic func(time.Time)) {
	ProgramRegistry[name] = logic
}

// countLeadingSpaces counts the number of leading spaces in a string.
func countLeadingSpaces(s string) int {
	return len(s) - len(strings.TrimLeft(s, " "))
}

// LoadConfigurationFromFile parses a plain text file and builds the core.Configuration struct.
// This parser is compatible with TinyGo.
func LoadConfigurationFromFile(path string) (*core.Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", path, err)
	}
	defer file.Close()

	var config *core.Configuration
	var currentResource *core.Resource
	var currentTask *core.Task
	var prioritySet map[int]struct{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		originalLine := scanner.Text()
		line := strings.TrimSpace(originalLine)

		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		key, value, found := strings.Cut(line, ":")
		if !found {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		indent := countLeadingSpaces(originalLine)

		switch indent {
		case 0: // Top-level: 'name' or 'resource'
			currentResource, currentTask = nil, nil // Reset context
			switch key {
			case "name":
				config = &core.Configuration{Name: value}
			case "resource":
				if config == nil {
					return nil, fmt.Errorf("config must have a 'name' before defining a 'resource'")
				}
				currentResource = &core.Resource{Name: value}
				config.WithResource(currentResource)
				prioritySet = make(map[int]struct{}) // New priority set for this resource
			default:
				return nil, fmt.Errorf("unknown top-level key: %s", key)
			}
		case 2: // Resource-level: 'cycle', 'affinity', or 'task'
			if currentResource == nil {
				return nil, fmt.Errorf("found resource-level property '%s' without a resource context", key)
			}
			currentTask = nil // Reset task context
			switch key {
			case "cycle":
				dur, err := time.ParseDuration(value)
				if err != nil {
					return nil, fmt.Errorf("invalid duration for cycle: %w", err)
				}
				currentResource.Cycle = dur
			case "affinity":
				aff, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("invalid integer for affinity: %w", err)
				}
				currentResource.Affinity = aff
			case "task":
				// Create a placeholder task; its properties will be filled in by indented lines.
				currentTask = core.NewTask(value, core.CyclicTask, 0, 0)
				currentResource.AddTask(currentTask)
			default:
				return nil, fmt.Errorf("unknown resource-level key: %s", key)
			}
		case 4: // Task-level: 'type', 'priority', 'interval', 'program'
			if currentTask == nil {
				return nil, fmt.Errorf("found task-level property '%s' without a task context", key)
			}
			switch key {
			case "type":
				switch strings.ToLower(value) {
				case "cyclic":
					currentTask.Type = core.CyclicTask
				case "eventdriven":
					currentTask.Type = core.EventDrivenTask
				default:
					return nil, fmt.Errorf("unknown task type: %s", value)
				}
			case "priority":
				p, err := strconv.Atoi(value)
				if err != nil {
					return nil, fmt.Errorf("invalid integer for priority: %w", err)
				}
				if _, exists := prioritySet[p]; exists {
					return nil, fmt.Errorf("duplicate task priority %d found in resource '%s'", p, currentResource.Name)
				}
				prioritySet[p] = struct{}{}
				currentTask.Priority = p
			case "interval":
				dur, err := time.ParseDuration(value)
				if err != nil {
					return nil, fmt.Errorf("invalid duration for interval: %w", err)
				}
				currentTask.Interval = dur
			case "program":
				if logic, ok := ProgramRegistry[value]; ok {
					currentTask.AddProgram(&core.Program{Name: value, Logic: logic})
				} else {
					return nil, fmt.Errorf("program '%s' (in task '%s') is not registered", value, currentTask.Name)
				}
			default:
				return nil, fmt.Errorf("unknown task-level key: %s", key)
			}
		default:
			return nil, fmt.Errorf("invalid indentation level (%d spaces) for line: %s", indent, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Post-parsing validation
	for _, r := range config.Resources {
		for _, t := range r.Tasks {
			if t.Type == core.CyclicTask && r.Cycle >= t.Interval {
				return nil, fmt.Errorf("resource '%s' cycle time (%v) must be faster than its task '%s' interval (%v)", r.Name, r.Cycle, t.Name, t.Interval)
			}
		}
		core.SortTasks(r.Tasks) // Sort tasks by priority after all are loaded.
	}

	return config, nil
}
