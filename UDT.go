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

package royaljelly

import (
	"reflect"
	"sync"
)

// UDT (User-Defined Type) defines the interface that any struct-based tag
// must implement to be used within the TagDatabase.
type UDT interface {
	// TypeName returns the unique data type name for this UDT.
	TypeName() DataType
}

var (
	udtRegistry = make(map[DataType]reflect.Type)
	udtMu       sync.RWMutex
)

// RegisterUDT makes a UDT type available to the TagDatabase system.
// This is necessary for creating new instances of the UDT during operations
// like reading from a persistence file.
func RegisterUDT(u UDT) {
	udtMu.Lock()
	defer udtMu.Unlock()
	name := u.TypeName()
	t := reflect.TypeOf(u).Elem() // Get the underlying struct type from the pointer
	udtRegistry[name] = t
}

// newUDTInstance creates a new instance of a registered UDT by its type name.
func newUDTInstance(name DataType) (UDT, bool) {
	udtMu.RLock()
	defer udtMu.RUnlock()
	t, ok := udtRegistry[name]
	if !ok {
		return nil, false
	}
	// Create a new pointer to a struct of the registered type.
	v := reflect.New(t).Interface()
	udt, ok := v.(UDT)
	return udt, ok
}
