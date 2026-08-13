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

// ProgramFactory is a function that creates a program's logic, optionally applying parameters.
type ProgramFactory func(params map[string]string) (func(time.Time), error)

// ProgramFactoryRegistry maps program type names to their factory functions.
var ProgramFactoryRegistry = make(map[string]ProgramFactory)

// RegisterProgramFactory adds a program factory to the registry.
func RegisterProgramFactory(typeName string, factory ProgramFactory) {
	ProgramFactoryRegistry[typeName] = factory
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
	var currentProgram *core.Program
	var currentProgramTypeName string
	var currentParams map[string]string
	var prioritySet map[int]struct{}

	scanner := bufio.NewScanner(file)

	// finalizeProgram is a helper to build the program from collected params.
	finalizeProgram := func() error {
		if currentProgram != nil && currentProgramTypeName != "" {
			factory, ok := ProgramFactoryRegistry[currentProgramTypeName]
			if !ok {
				return fmt.Errorf("unknown program type '%s' for instance '%s'", currentProgramTypeName, currentProgram.Name)
			}
			logic, err := factory(currentParams)
			if err != nil {
				return fmt.Errorf("error creating program '%s': %w", currentProgram.Name, err)
			}
			currentProgram.Logic = logic
		}
		currentProgram, currentProgramTypeName, currentParams = nil, "", nil
		return nil
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		key, value, found := strings.Cut(line, ":")
		if !found {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		indent := countLeadingSpaces(scanner.Text())

		switch indent {
		case 0: // Top-level: 'name' or 'resource'
			if err := finalizeProgram(); err != nil {
				return nil, err
			}
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
			if err := finalizeProgram(); err != nil {
				return nil, err
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
		case 4: // Task-level: 'type', 'priority', 'interval', or 'program'
			if currentTask == nil {
				return nil, fmt.Errorf("found task-level property '%s' without a task context", key)
			}
			if indent < 6 { // Finalize program if we are moving to a new task-level item
				if err := finalizeProgram(); err != nil {
					return nil, err
				}
			}
			switch key {
			case "program":
				// The 'program' key at this level defines an instance with a type.
				// e.g., "program: CounterA CounterProgram"
				instanceName, typeName, found := strings.Cut(value, " ")
				if !found {
					return nil, fmt.Errorf("program definition requires an instance name and a type name (e.g., 'program: MyInstance MyType'), got: '%s'", value)
				}
				instanceName, typeName = strings.TrimSpace(instanceName), strings.TrimSpace(typeName)

				// Create a new program instance and add it to the task.
				currentProgram = &core.Program{Name: instanceName}
				currentProgramTypeName = typeName
				currentParams = make(map[string]string)
				// Pass the instance name to the factory as a default parameter.
				currentParams["name"] = instanceName
				currentTask.AddProgram(currentProgram)

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

			default:
				return nil, fmt.Errorf("unknown task-level key: %s", key)
			}
		case 6: // Program parameter level
			if currentProgram == nil || currentParams == nil {
				return nil, fmt.Errorf("found parameter '%s' without a program definition context", key)
			}
			if key != "param" {
				return nil, fmt.Errorf("unknown program-level key: '%s', expected 'param'", key)
			}

			paramKey, paramValue, found := strings.Cut(value, " ")
			if !found {
				return nil, fmt.Errorf("invalid param format for program '%s', expected 'key value'", currentProgram.Name)
			}
			currentParams[strings.TrimSpace(paramKey)] = strings.TrimSpace(paramValue)

		default:
			return nil, fmt.Errorf("invalid indentation level (%d spaces) for line: %s", indent, line)
		}
	}

	// Finalize the last program in the file.
	if err := finalizeProgram(); err != nil {
		return nil, err
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Finalization pass
	if config != nil {
		for _, r := range config.Resources {
			for _, t := range r.Tasks {
				if t.Type == core.CyclicTask && r.Cycle >= t.Interval {
					return nil, fmt.Errorf("resource '%s' cycle time (%v) must be faster than its task '%s' interval (%v)", r.Name, r.Cycle, t.Name, t.Interval)
				}
			}
			core.SortTasks(r.Tasks) // Sort tasks by priority after all are loaded.
		}
	}

	return config, nil
}
