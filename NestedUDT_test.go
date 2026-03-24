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
	"path/filepath"
	"testing"
)

// MotorConfig represents a nested UDT.
type MotorConfig struct {
	MaxSpeed REAL
	RampTime TIME
}

// TypeName implements the UDT interface for MotorConfig.
func (mc *MotorConfig) TypeName() DataType {
	return "MotorConfig"
}

// DriveSystem is a parent UDT that contains another UDT.
type DriveSystem struct {
	Name   STRING
	Config *MotorConfig // Nested UDT
	Active BOOL
}

// TypeName implements the UDT interface for DriveSystem.
func (ds *DriveSystem) TypeName() DataType {
	return "DriveSystem"
}

// TestNestedUDTPersistence verifies that a UDT containing another UDT
// can be persisted and loaded correctly.
func TestNestedUDTPersistence(t *testing.T) {
	// 1. Register all UDTs involved, both parent and nested.
	// This is a crucial step for the system to be aware of all custom types.
	RegisterUDT(&DriveSystem{})
	RegisterUDT(&MotorConfig{})

	// --- Write Phase ---
	dbWrite := NewTagDatabase()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "nested_udt_persistence.txt")

	driveTag := Tag{
		Name:     "MainDrive",
		DataType: "DriveSystem",
		Value: &DriveSystem{
			Name: "Conveyor 1",
			Config: &MotorConfig{
				MaxSpeed: 3600.0,
				RampTime: TIME_TO_TIME(DINT(5000)), // 5 seconds
			},
			Active: true,
		},
		Persistent: true,
	}

	if err := dbWrite.AddTag(driveTag); err != nil {
		t.Fatalf("Failed to add nested UDT tag: %v", err)
	}

	if err := dbWrite.WritePersistentTagsToFile(filePath); err != nil {
		t.Fatalf("Failed to write nested UDT to file: %v", err)
	}

	// --- Read Phase ---
	dbRead := NewTagDatabase()
	dbRead.AddTag(Tag{Name: "MainDrive", DataType: "DriveSystem", Value: &DriveSystem{Config: &MotorConfig{}}})

	if err := dbRead.ReadPersistentTagsFromFile(filePath); err != nil {
		t.Fatalf("Failed to read nested UDT from file: %v", err)
	}

	retrievedTag, _ := dbRead.GetTag("MainDrive")
	retrievedDrive, _ := retrievedTag.Value.(*DriveSystem)

	if retrievedDrive.Name != "Conveyor 1" || retrievedDrive.Config.MaxSpeed != 3600.0 {
		t.Errorf("Data mismatch after reading nested UDT from file. Got %+v", retrievedDrive)
	}
}
