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

// Package vars provides thread-safe ways to share data between programs,
// tasks and resources.
//
// Programs on the same Resource run sequentially and may share plain Go
// variables. Programs on different Resources, and code outside the scheduler
// (for example an HMI or a test), run concurrently and must exchange data
// through a synchronized container such as Shared or ProcessImage.
package vars

import (
	"sync"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

// Shared holds one value that can be read and written from several goroutines,
// like an IEC 61131-3 VAR_GLOBAL used across resources. The zero value holds the
// zero value of T and is ready to use. A Shared must not be copied after first use.
type Shared[T any] struct {
	mu sync.RWMutex
	v  T
}

// NewShared returns a Shared holding v.
func NewShared[T any](v T) *Shared[T] {
	return &Shared[T]{v: v}
}

// Load returns the current value.
func (s *Shared[T]) Load() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.v
}

// Store replaces the current value.
func (s *Shared[T]) Store(v T) {
	s.mu.Lock()
	s.v = v
	s.mu.Unlock()
}

// Update atomically replaces the value with fn(current) and returns the new value.
// fn must not call other methods of s.
func (s *Shared[T]) Update(fn func(T) T) T {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.v = fn(s.v)
	return s.v
}

// AddressSize is the number of elements in each area of an Addresses table.
const AddressSize = 255

// Addresses is a table of directly represented variables for one memory area,
// indexed by address. Each IEC 61131-3 size prefix has its own array:
//
//	%IX0  / %QX0  / %MX0   bit    → B[0]  (BOOL)
//	%IB0  / %QB0  / %MB0   byte   → C[0]  (BYTE)
//	%IW0  / %QW0  / %MW0   word   → W[0]  (WORD)
//	%ID0  / %QD0  / %MD0   dword  → D[0]  (DWORD)
//	%IL0  / %QL0  / %ML0   lword  → L[0]  (LWORD)
//
// The areas are independent: %QW4 and %QD4 are different variables. The R, LR,
// S and WS arrays hold REAL, LREAL, STRING and WSTRING values for programs that
// exchange typed data through the image.
type Addresses struct {
	B  [AddressSize]iec.BOOL
	C  [AddressSize]iec.BYTE
	W  [AddressSize]iec.WORD
	D  [AddressSize]iec.DWORD
	L  [AddressSize]iec.LWORD
	R  [AddressSize]iec.REAL
	LR [AddressSize]iec.LREAL
	S  [AddressSize]iec.STRING
	WS [AddressSize]iec.WSTRING
	T  iec.TIME
}

// Image is the content of a process image: inputs (%I), outputs (%Q) and
// memory (%M).
type Image struct {
	I Addresses
	Q Addresses
	M Addresses
}

// ProcessImage guards an Image so that I/O drivers and programs on different
// resources can use it concurrently. The zero value is ready to use.
//
// The usual pattern is for a program to copy what it needs at the start of its
// run with Read, compute, and publish its outputs with Write, keeping the time
// spent holding the lock short.
type ProcessImage struct {
	mu  sync.RWMutex
	img Image
}

// Read calls fn with shared (read-only) access to the image. fn must not retain
// the pointer or modify the image.
func (p *ProcessImage) Read(fn func(img *Image)) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	fn(&p.img)
}

// Write calls fn with exclusive access to the image. fn must not retain the pointer.
func (p *ProcessImage) Write(fn func(img *Image)) {
	p.mu.Lock()
	defer p.mu.Unlock()
	fn(&p.img)
}

// Snapshot copies the whole image into dst without allocating.
func (p *ProcessImage) Snapshot(dst *Image) {
	p.mu.RLock()
	*dst = p.img
	p.mu.RUnlock()
}

// Deprecated package-level tables. They are plain variables with no
// synchronization, so concurrent use from several resources is a data race.
// Use a ProcessImage instead.
var (
	// Deprecated: use a ProcessImage.
	I Addresses
	// Deprecated: use a ProcessImage.
	Q Addresses
	// Deprecated: use a ProcessImage.
	M Addresses
	// Deprecated: pass the scan time from Program.Logic instead.
	IEC_TIME time.Time
)
