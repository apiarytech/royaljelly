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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestWriteAndReadPersistentTags verifies the entire persistence cycle.
func TestWriteAndReadPersistentTags(t *testing.T) {
	// --- 1. Setup and Write Phase ---
	dbWrite := NewTagDatabase()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "persistent_tags.txt")

	// Add a mix of persistent and non-persistent tags
	tagsToWrite := []Tag{
		{Name: "PersistentDINT", DataType: TypeDINT, Value: DINT(123), Persistent: true},
		{Name: "PersistentREAL", DataType: TypeREAL, Value: REAL(45.67), Persistent: true},
		{Name: "NonPersistentINT", DataType: TypeINT, Value: INT(999), Persistent: false},
		{Name: "PersistentSTRING", DataType: TypeSTRING, Value: STRING("hello world"), Persistent: true},
	}
	for _, tag := range tagsToWrite {
		if err := dbWrite.AddTag(tag); err != nil {
			t.Fatalf("Failed to add tag %s during write setup: %v", tag.Name, err)
		}
	}

	// Write the persistent tags to the file
	if err := dbWrite.WritePersistentTagsToFile(filePath); err != nil {
		t.Fatalf("WritePersistentTagsToFile() returned an unexpected error: %v", err)
	}

	// Verify the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read the created persistence file: %v", err)
	}

	expectedContent := "PersistentDINT=123\nPersistentREAL=45.67\nPersistentSTRING=hello world"
	if strings.TrimSpace(string(content)) != expectedContent {
		t.Errorf("File content mismatch.\nGot:\n%s\n\nWant:\n%s", string(content), expectedContent)
	}

	// --- 2. Read and Verify Phase ---
	dbRead := NewTagDatabase()

	// Populate the "new" database, simulating a restart with default values
	tagsToRead := []Tag{
		{Name: "PersistentDINT", DataType: TypeDINT, Value: DINT(0), Persistent: true},
		{Name: "PersistentREAL", DataType: TypeREAL, Value: REAL(0.0), Persistent: true},
		{Name: "NonPersistentINT", DataType: TypeINT, Value: INT(999), Persistent: false}, // This should not be affected
		{Name: "PersistentSTRING", DataType: TypeSTRING, Value: STRING(""), Persistent: true},
		{Name: "UntouchedPersistent", DataType: TypeBOOL, Value: BOOL(true), Persistent: true}, // This tag is not in the file
	}
	for _, tag := range tagsToRead {
		if err := dbRead.AddTag(tag); err != nil {
			t.Fatalf("Failed to add tag %s during read setup: %v", tag.Name, err)
		}
	}

	// Read the values back from the file
	if err := dbRead.ReadPersistentTagsFromFile(filePath); err != nil {
		t.Fatalf("ReadPersistentTagsFromFile() returned an unexpected error: %v", err)
	}

	// Verify the values in the new database
	testCases := []struct {
		tagName      string
		expectedVal  interface{}
		shouldChange bool
	}{
		{"PersistentDINT", DINT(123), true},
		{"PersistentREAL", REAL(45.67), true},
		{"PersistentSTRING", STRING("hello world"), true},
		{"NonPersistentINT", INT(999), false},      // Should remain unchanged
		{"UntouchedPersistent", BOOL(true), false}, // Should remain unchanged
	}

	for _, tc := range testCases {
		t.Run("Verify_"+tc.tagName, func(t *testing.T) {
			val, err := dbRead.GetTagValue(tc.tagName)
			if err != nil {
				t.Fatalf("GetTagValue failed for %s: %v", tc.tagName, err)
			}
			if val != tc.expectedVal {
				t.Errorf("Tag %s has wrong value. Got %v (%T), want %v (%T)", tc.tagName, val, val, tc.expectedVal, tc.expectedVal)
			}
		})
	}
}

// TestReadPersistentTags_FileNotExist tests that no error occurs if the file doesn't exist.
func TestReadPersistentTags_FileNotExist(t *testing.T) {
	db := NewTagDatabase()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "non_existent_file.txt")

	err := db.ReadPersistentTagsFromFile(filePath)
	if err != nil {
		t.Fatalf("ReadPersistentTagsFromFile() should not return an error for a non-existent file, but got: %v", err)
	}
}

// TestReadPersistentTags_ParseError tests that the function continues after a parsing error.
func TestReadPersistentTags_ParseError(t *testing.T) {
	db := NewTagDatabase()
	db.AddTag(Tag{Name: "MyTag", DataType: TypeDINT, Value: DINT(0), Persistent: true})

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "bad_file.txt")

	// Create a file with a malformed line
	badContent := []byte("MyTag=not_a_number")
	if err := os.WriteFile(filePath, badContent, 0644); err != nil {
		t.Fatalf("Failed to write bad file: %v", err)
	}

	err := db.ReadPersistentTagsFromFile(filePath)
	if err == nil {
		t.Fatal("ReadPersistentTagsFromFile() should have returned an error for a parse failure")
	}
	if !strings.Contains(err.Error(), "error parsing value for tag 'MyTag'") {
		t.Errorf("Expected a parsing error, but got: %v", err)
	}
}

// TestGetPersistentTags verifies the GetPersistentTags method.
func TestGetPersistentTags(t *testing.T) {
	db := NewTagDatabase()

	// Add a mix of persistent and non-persistent tags
	tags := []Tag{
		{Name: "P1", DataType: TypeDINT, Persistent: true},
		{Name: "NP1", DataType: TypeREAL, Persistent: false},
		{Name: "P2", DataType: TypeSTRING, Persistent: true},
		{Name: "NP2", DataType: TypeBOOL, Persistent: false},
	}
	for _, tag := range tags {
		db.AddTag(tag)
	}

	// Get only the persistent tags
	persistentTags := db.GetPersistentTags()

	// Check the count
	if len(persistentTags) != 2 {
		t.Fatalf("GetPersistentTags() returned %d tags, want 2", len(persistentTags))
	}

	// Check that the correct tags were returned
	foundP1 := false
	foundP2 := false
	for _, tag := range persistentTags {
		if tag.Name == "P1" {
			foundP1 = true
		}
		if tag.Name == "P2" {
			foundP2 = true
		}
		if !tag.Persistent {
			t.Errorf("GetPersistentTags() returned a non-persistent tag: %s", tag.Name)
		}
	}

	if !foundP1 || !foundP2 {
		t.Error("GetPersistentTags() did not return all expected persistent tags.")
	}
}
