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

// Package config builds a core.Configuration from a small indented text format.
// The parser has no reflection or encoding dependencies, so it works under TinyGo.
//
// Format:
//
//	# comment
//	name: MyPLC
//	resource: MainCPU
//	  cycle: 10ms
//	  affinity: 1
//	  task: FastTask
//	    type: Cyclic          # or EventDriven
//	    priority: 1
//	    interval: 100ms
//	    watchdog: 50ms        # optional
//	    program: Counter1 CounterProgram
//	      param: initial_value 100
//
// Each level is indented one step deeper than its parent. A step can be a tab
// or any fixed number of spaces, as long as the file is consistent.
// Program types are resolved through a Registry of ProgramFactory functions.
package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apiarytech/royaljelly/core"
)

// ProgramFactory creates a program's logic, optionally applying parameters.
// The params map always contains "name", set to the program instance name.
type ProgramFactory func(params map[string]string) (func(time.Time), error)

// Registry maps program type names to factories. It is safe for concurrent use.
// The zero value is an empty registry ready to use.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]ProgramFactory
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register adds or replaces the factory for a program type.
func (r *Registry) Register(typeName string, factory ProgramFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.factories == nil {
		r.factories = make(map[string]ProgramFactory)
	}
	r.factories[typeName] = factory
}

// Lookup returns the factory for a program type.
func (r *Registry) Lookup(typeName string) (ProgramFactory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.factories[typeName]
	return f, ok
}

// DefaultRegistry is used by RegisterProgramFactory and the package-level Load functions.
var DefaultRegistry = NewRegistry()

// RegisterProgramFactory adds a program factory to DefaultRegistry.
func RegisterProgramFactory(typeName string, factory ProgramFactory) {
	DefaultRegistry.Register(typeName, factory)
}

// LoadConfigurationFromFile parses a configuration file using DefaultRegistry.
func LoadConfigurationFromFile(path string) (*core.Configuration, error) {
	return Loader{}.LoadFile(path)
}

// LoadConfigurationFromString parses a configuration held in a string, which is
// useful for configurations embedded in the binary. It uses DefaultRegistry.
func LoadConfigurationFromString(configString string) (*core.Configuration, error) {
	return Loader{}.Load(strings.NewReader(configString))
}

// Loader parses configurations using a specific Registry. A zero Loader uses
// DefaultRegistry.
type Loader struct {
	Registry *Registry
}

// LoadFile parses the configuration file at path.
func (l Loader) LoadFile(path string) (*core.Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %w", path, err)
	}
	defer file.Close()
	cfg, err := l.Load(file)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return cfg, nil
}

// Load parses a configuration from r and validates it with core.Configuration.Validate.
func (l Loader) Load(r io.Reader) (*core.Configuration, error) {
	reg := l.Registry
	if reg == nil {
		reg = DefaultRegistry
	}
	p := parser{registry: reg}
	return p.parse(r)
}

// parser holds the state of one parse.
type parser struct {
	registry *Registry

	lineNo int
	indent indentation

	config   *core.Configuration
	resource *core.Resource
	task     *core.Task
	program  *core.Program
	typeName string
	params   map[string]string
}

func (p *parser) errorf(format string, args ...any) error {
	return fmt.Errorf("line %d: "+format, append([]any{p.lineNo}, args...)...)
}

// finalizeProgram resolves the pending program's logic from its factory.
func (p *parser) finalizeProgram() error {
	if p.program != nil {
		factory, ok := p.registry.Lookup(p.typeName)
		if !ok {
			return fmt.Errorf("unknown program type '%s' for instance '%s'", p.typeName, p.program.Name)
		}
		logic, err := factory(p.params)
		if err != nil {
			return fmt.Errorf("error creating program '%s': %w", p.program.Name, err)
		}
		p.program.Logic = logic
	}
	p.program, p.typeName, p.params = nil, "", nil
	return nil
}

func (p *parser) parse(r io.Reader) (*core.Configuration, error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		p.lineNo++
		raw := scanner.Text()
		line := strings.TrimSpace(stripComment(raw))
		if line == "" {
			continue
		}
		level, err := p.indent.level(raw)
		if err != nil {
			return nil, p.errorf("%v", err)
		}
		key, value, found := strings.Cut(line, ":")
		if !found {
			return nil, p.errorf("expected 'key: value', got '%s'", line)
		}
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)

		// Moving to a line at task level or above ends the current program block.
		if level <= 2 {
			if err := p.finalizeProgram(); err != nil {
				return nil, p.errorf("%v", err)
			}
		}
		switch level {
		case 0:
			err = p.topLevel(key, value)
		case 1:
			err = p.resourceLevel(key, value)
		case 2:
			err = p.taskLevel(key, value)
		case 3:
			err = p.programLevel(key, value)
		default:
			err = fmt.Errorf("line is nested too deeply (level %d)", level)
		}
		if err != nil {
			return nil, p.errorf("%v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}
	if err := p.finalizeProgram(); err != nil {
		return nil, err
	}
	if p.config == nil {
		return nil, fmt.Errorf("config has no 'name' entry")
	}
	if err := p.config.Validate(); err != nil {
		return nil, err
	}
	return p.config, nil
}

func (p *parser) topLevel(key, value string) error {
	p.resource, p.task = nil, nil
	switch key {
	case "name":
		if p.config != nil {
			return fmt.Errorf("duplicate 'name' entry")
		}
		p.config = &core.Configuration{Name: value}
	case "resource":
		if p.config == nil {
			return fmt.Errorf("config must have a 'name' before defining a 'resource'")
		}
		p.resource = &core.Resource{Name: value}
		p.config.AddResource(p.resource)
	default:
		return fmt.Errorf("unknown top-level key '%s'", key)
	}
	return nil
}

func (p *parser) resourceLevel(key, value string) error {
	if p.resource == nil {
		return fmt.Errorf("found resource-level property '%s' without a resource context", key)
	}
	p.task = nil
	switch key {
	case "cycle":
		dur, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration for cycle: %w", err)
		}
		p.resource.Cycle = dur
	case "affinity":
		aff, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer for affinity: %w", err)
		}
		p.resource.Affinity = aff
	case "task":
		p.task = core.NewTask(value, core.CyclicTask, 0, 0)
		p.resource.AddTask(p.task)
	default:
		return fmt.Errorf("unknown resource-level key '%s'", key)
	}
	return nil
}

func (p *parser) taskLevel(key, value string) error {
	if p.task == nil {
		return fmt.Errorf("found task-level property '%s' without a task context", key)
	}
	switch key {
	case "program":
		instanceName, typeName, found := strings.Cut(value, " ")
		if !found {
			return fmt.Errorf("program definition requires an instance name and a type name (e.g., 'program: MyInstance MyType'), got: '%s'", value)
		}
		instanceName, typeName = strings.TrimSpace(instanceName), strings.TrimSpace(typeName)
		p.program = &core.Program{Name: instanceName}
		p.typeName = typeName
		p.params = map[string]string{"name": instanceName}
		p.task.AddProgram(p.program)
	case "type":
		switch strings.ToLower(value) {
		case "cyclic":
			p.task.Type = core.CyclicTask
		case "eventdriven":
			p.task.Type = core.EventDrivenTask
		default:
			return fmt.Errorf("unknown task type '%s'", value)
		}
	case "priority":
		prio, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer for priority: %w", err)
		}
		p.task.Priority = prio
		// Re-add so the resource re-sorts its tasks by the new priority.
		p.resource.RemoveTask(p.task.Name)
		p.resource.AddTask(p.task)
	case "interval":
		dur, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration for interval: %w", err)
		}
		p.task.Interval = dur
	case "watchdog":
		dur, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid duration for watchdog: %w", err)
		}
		p.task.Watchdog = dur
	default:
		return fmt.Errorf("unknown task-level key '%s'", key)
	}
	return nil
}

func (p *parser) programLevel(key, value string) error {
	if p.program == nil {
		return fmt.Errorf("found parameter '%s' without a program definition context", key)
	}
	if key != "param" {
		return fmt.Errorf("unknown program-level key '%s', expected 'param'", key)
	}
	paramKey, paramValue, found := strings.Cut(value, " ")
	if !found {
		return fmt.Errorf("invalid param format for program '%s', expected 'key value'", p.program.Name)
	}
	p.params[strings.TrimSpace(paramKey)] = strings.TrimSpace(paramValue)
	return nil
}

// stripComment removes a trailing '#' comment. A '#' starts a comment at the
// beginning of a line or after whitespace, so values such as "a#b" are kept.
func stripComment(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '#' && (i == 0 || s[i-1] == ' ' || s[i-1] == '\t') {
			return s[:i]
		}
	}
	return s
}

// indentation learns the file's indent style from the first indented line and
// converts leading whitespace into a nesting level.
type indentation struct {
	useTabs bool
	unit    int // Spaces per level; 0 until the first indented line is seen.
}

func (in *indentation) level(line string) (int, error) {
	n := 0
	tabs, spaces := false, false
	for n < len(line) && (line[n] == ' ' || line[n] == '\t') {
		if line[n] == '\t' {
			tabs = true
		} else {
			spaces = true
		}
		n++
	}
	if n == 0 {
		return 0, nil
	}
	if tabs && spaces {
		return 0, fmt.Errorf("indentation mixes tabs and spaces")
	}
	if in.unit == 0 {
		in.useTabs = tabs
		in.unit = n
		if tabs {
			in.unit = 1
		}
	}
	if tabs != in.useTabs {
		return 0, fmt.Errorf("inconsistent indentation: file indents with %s", in.style())
	}
	if n%in.unit != 0 {
		return 0, fmt.Errorf("indentation of %d is not a multiple of the file's indent step (%d %s)", n, in.unit, in.style())
	}
	return n / in.unit, nil
}

func (in *indentation) style() string {
	if in.useTabs {
		return "tabs"
	}
	return "spaces"
}
