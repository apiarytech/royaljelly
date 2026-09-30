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

package vars

import (
	"fmt"
	"sync"
	"testing"

	"github.com/apiarytech/royaljelly/iec"
)

func TestSharedConcurrentUpdate(t *testing.T) {
	var s Shared[iec.LINT]
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				s.Update(func(v iec.LINT) iec.LINT { return v + 1 })
			}
		}()
	}
	wg.Wait()
	if got := s.Load(); got != 8000 {
		t.Fatalf("Load() = %d, want 8000", got)
	}
	s.Store(5)
	if got := NewShared[iec.LINT](7).Load() + s.Load(); got != 12 {
		t.Fatalf("got %d, want 12", got)
	}
}

func TestProcessImage(t *testing.T) {
	var p ProcessImage
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				p.Write(func(img *Image) { img.M.D[0]++ })
			}
		}()
		go func() {
			defer wg.Done()
			var snap Image
			for j := 0; j < 500; j++ {
				p.Snapshot(&snap)
				p.Read(func(img *Image) { _ = img.M.D[0] })
			}
		}()
	}
	wg.Wait()
	var got iec.DWORD
	p.Read(func(img *Image) { got = img.M.D[0] })
	if got != 2000 {
		t.Fatalf("M.D[0] = %d, want 2000", got)
	}
}

func TestWordAndDwordAreasAreIndependent(t *testing.T) {
	var p ProcessImage
	p.Write(func(img *Image) {
		img.Q.W[4] = 0xBEEF     // %QW4
		img.Q.D[4] = 0x12345678 // %QD4
		img.Q.WS[0] = "Grüße"
	})
	var w iec.WORD
	var d iec.DWORD
	var ws iec.WSTRING
	p.Read(func(img *Image) { w, d, ws = img.Q.W[4], img.Q.D[4], img.Q.WS[0] })
	if w != 0xBEEF || d != 0x12345678 || ws != "Grüße" {
		t.Fatalf("got %%QW4=%#x %%QD4=%#x WS=%q", w, d, ws)
	}
}

func ExampleShared() {
	// A counter written by one resource and read by another.
	var count Shared[iec.LINT]
	count.Update(func(v iec.LINT) iec.LINT { return v + 1 })
	count.Update(func(v iec.LINT) iec.LINT { return v + 1 })
	fmt.Println(count.Load())
	// Output: 2
}

func ExampleProcessImage() {
	var pi ProcessImage
	// An input driver publishes %IX0.
	pi.Write(func(img *Image) { img.I.B[0] = true })
	// A program copies its inputs, computes, and publishes %QX0.
	var in iec.BOOL
	pi.Read(func(img *Image) { in = img.I.B[0] })
	pi.Write(func(img *Image) { img.Q.B[0] = !in })

	var out iec.BOOL
	pi.Read(func(img *Image) { out = img.Q.B[0] })
	fmt.Println(in, out)
	// Output: true false
}
