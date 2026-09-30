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

package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Configuration is the top-level element, representing the entire PLC system.
type Configuration struct {
	Name    string
	OnFault FaultHandler // Optional: default fault handler for resources that have none.

	mu        sync.Mutex
	resources []*Resource
	unwatch   chan struct{} // Closed to end the context watcher started by Start.
}

// WithResource adds a resource to the configuration and returns the configuration for chaining.
func (c *Configuration) WithResource(r *Resource) *Configuration {
	c.AddResource(r)
	return c
}

// AddResource adds a resource to the configuration.
func (c *Configuration) AddResource(r *Resource) {
	c.mu.Lock()
	defer c.mu.Unlock()
	next := make([]*Resource, len(c.resources), len(c.resources)+1)
	copy(next, c.resources)
	c.resources = append(next, r)
}

// RemoveResource stops the first resource with the given name and removes it.
// It returns true if the resource was found and removed.
func (c *Configuration) RemoveResource(name string) bool {
	c.mu.Lock()
	var removed *Resource
	for i, r := range c.resources {
		if r.Name == name {
			removed = r
			next := make([]*Resource, 0, len(c.resources)-1)
			next = append(next, c.resources[:i]...)
			c.resources = append(next, c.resources[i+1:]...)
			break
		}
	}
	c.mu.Unlock()
	if removed == nil {
		return false
	}
	removed.Stop()
	return true
}

// FindResource returns the first resource with the given name, or nil.
func (c *Configuration) FindResource(name string) *Resource {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, r := range c.resources {
		if r.Name == name {
			return r
		}
	}
	return nil
}

// Resources returns a copy of the configuration's resources.
func (c *Configuration) Resources() []*Resource {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]*Resource, len(c.resources))
	copy(out, c.resources)
	return out
}

// Validate checks every resource and reports duplicate resource names.
func (c *Configuration) Validate() error {
	var errs []error
	names := make(map[string]bool)
	for _, r := range c.Resources() {
		if names[r.Name] {
			errs = append(errs, fmt.Errorf("configuration '%s': duplicate resource name '%s'", c.Name, r.Name))
		}
		names[r.Name] = true
		if err := r.Validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Start validates the configuration and starts every resource. If any resource
// fails to start, the ones already started are stopped and the error is returned.
// When ctx is cancelled, all resources are stopped. A nil ctx means no cancellation.
func (c *Configuration) Start(ctx context.Context) error {
	if err := c.Validate(); err != nil {
		return err
	}
	resources := c.Resources()
	started := make([]*Resource, 0, len(resources))
	for _, r := range resources {
		if r.OnFault == nil && c.OnFault != nil && !r.IsRunning() {
			r.OnFault = c.OnFault
		}
		if err := r.Start(); err != nil {
			for _, s := range started {
				s.Stop()
			}
			return err
		}
		started = append(started, r)
	}

	if ctx != nil && ctx.Done() != nil {
		unwatch := make(chan struct{})
		c.mu.Lock()
		if c.unwatch != nil {
			close(c.unwatch)
		}
		c.unwatch = unwatch
		c.mu.Unlock()
		go func() {
			select {
			case <-ctx.Done():
				c.Stop()
			case <-unwatch:
			}
		}()
	}
	return nil
}

// Stop stops every resource and waits for their schedulers to exit.
func (c *Configuration) Stop() {
	c.mu.Lock()
	if c.unwatch != nil {
		close(c.unwatch)
		c.unwatch = nil
	}
	c.mu.Unlock()
	for _, r := range c.Resources() {
		r.Stop()
	}
}

// Run starts the configuration, blocks until ctx is done, and then stops every
// resource. It returns an error only if the configuration fails to start.
func (c *Configuration) Run(ctx context.Context) error {
	if err := c.Start(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	c.Stop()
	return nil
}
