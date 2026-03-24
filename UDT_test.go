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

// MotorData is an example of a User-Defined Type (UDT).
// It's a struct that will be used as a tag's value.
type MotorData struct {
	Speed    REAL
	Current  REAL
	Temp     REAL
	Running  BOOL
	Tripped  BOOL
	TestName string
}

// TypeName implements the UDT interface, returning the unique name for this type.
func (m *MotorData) TypeName() DataType {
	return "MotorData"
}

// TestUDTInDatabase verifies that a UDT can be added, retrieved,
// and persisted correctly.
func TestUDTInDatabase(t *testing.T) {
	// 1. Register our new UDT so the system knows about it.
	// The pointer receiver on TypeName() means we must pass a pointer here.
	RegisterUDT(&MotorData{})

	// --- Write Phase ---
	dbWrite := NewTagDatabase()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "udt_persistence.txt")

	// Create an instance of our UDT tag.
	motorTag := Tag{
		Name:     "MainMotor",
		DataType: "MotorData", // The DataType string must match TypeName()
		Value: &MotorData{
			Speed:   1750.5,
			Current: 45.2,
			Temp:    65.7,
			Running: true,
		},
		Persistent: true,
	}

	if err := dbWrite.AddTag(motorTag); err != nil {
		t.Fatalf("Failed to add UDT tag: %v", err)
	}

	// Write to file
	if err := dbWrite.WritePersistentTagsToFile(filePath); err != nil {
		t.Fatalf("Failed to write UDT to file: %v", err)
	}

	// --- Read Phase ---
	dbRead := NewTagDatabase()
	// Pre-populate the read DB with a placeholder for the tag.
	dbRead.AddTag(Tag{Name: "MainMotor", DataType: "MotorData", Value: &MotorData{}})

	// Read from the file, which should deserialize the JSON.
	if err := dbRead.ReadPersistentTagsFromFile(filePath); err != nil {
		t.Fatalf("Failed to read UDT from file: %v", err)
	}

	// Verify the data was loaded correctly.
	retrievedTag, found := dbRead.GetTag("MainMotor")
	if !found {
		t.Fatal("Failed to retrieve UDT tag after reading from file.")
	}

	retrievedMotorData, ok := retrievedTag.Value.(*MotorData)
	if !ok {
		t.Fatalf("Retrieved tag value is not of type *MotorData")
	}

	if retrievedMotorData.Speed != 1750.5 || !retrievedMotorData.Running {
		t.Errorf("Data mismatch after reading UDT from file. Got %+v", retrievedMotorData)
	}
}

// TestGetNestedUDTField verifies that nested fields of a UDT can be accessed
// using dot notation (e.g., "MyTag.MyField").
func TestGetNestedUDTField(t *testing.T) {
	RegisterUDT(&MotorData{})
	db := NewTagDatabase()

	motorTag := Tag{
		Name:     "MainMotor",
		DataType: "MotorData",
		Value: &MotorData{
			Speed:   1800.0,
			Current: 50.5,
			Running: true,
		},
	}
	db.AddTag(motorTag)

	testCases := []struct {
		name         string
		nestedTag    string
		expectedVal  interface{}
		expectFound  bool
		expectedType DataType
	}{
		{"Access REAL field", "MainMotor.Speed", REAL(1800.0), true, TypeREAL},
		{"Access BOOL field", "MainMotor.Running", BOOL(true), true, TypeBOOL},
		{"Access non-existent field", "MainMotor.NonExistent", nil, false, ""},
		{"Access field on non-existent tag", "FakeMotor.Speed", nil, false, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test GetTagValue
			val, err := db.GetTagValue(tc.nestedTag)
			if tc.expectFound {
				if err != nil {
					t.Fatalf("GetTagValue(%q) returned unexpected error: %v", tc.nestedTag, err)
				}
				if val != tc.expectedVal {
					t.Errorf("GetTagValue(%q) = %v; want %v", tc.nestedTag, val, tc.expectedVal)
				}
			} else {
				if err == nil {
					t.Errorf("GetTagValue(%q) was expected to fail, but it succeeded.", tc.nestedTag)
				}
			}
		})
	}
}

// TestSetNestedUDTField verifies that nested fields of a UDT can be written to
// using dot notation via SetTagValue.
func TestSetNestedUDTField(t *testing.T) {
	RegisterUDT(&MotorData{})
	db := NewTagDatabase()

	motorTag := Tag{
		Name:     "MainMotor",
		DataType: "MotorData",
		Value: &MotorData{
			Speed:   1800.0,
			Running: false,
		},
	}
	db.AddTag(motorTag)

	// 1. Test successful write to a nested field.
	newSpeed := REAL(1950.5)
	err := db.SetTagValue("MainMotor.Speed", newSpeed)
	if err != nil {
		t.Fatalf("SetTagValue on nested field returned unexpected error: %v", err)
	}

	// Verify the change by reading it back.
	val, err := db.GetTagValue("MainMotor.Speed")
	if err != nil {
		t.Fatalf("GetTagValue for nested field failed after write: %v", err)
	}
	if val != newSpeed {
		t.Errorf("Nested field value was not updated. Got %v, want %v", val, newSpeed)
	}

	// Also verify by getting the whole UDT.
	tag, _ := db.GetTag("MainMotor")
	motorData := tag.Value.(*MotorData)
	if motorData.Speed != newSpeed {
		t.Errorf("UDT struct in map was not updated. Got speed %v, want %v", motorData.Speed, newSpeed)
	}

	// 2. Test writing with an incorrect type.
	err = db.SetTagValue("MainMotor.Speed", DINT(2000)) // REAL field, DINT value
	if err == nil {
		t.Error("SetTagValue with type mismatch should have returned an error.")
	}

	// 3. Test writing to a non-existent field.
	err = db.SetTagValue("MainMotor.NonExistentField", REAL(1.0))
	if err == nil {
		t.Error("SetTagValue on a non-existent field should have returned an error.")
	}

	// 4. Test writing to a non-existent base tag.
	err = db.SetTagValue("FakeMotor.Speed", REAL(1.0))
	if err == nil {
		t.Error("SetTagValue on a non-existent base tag should have returned an error.")
	}
}
