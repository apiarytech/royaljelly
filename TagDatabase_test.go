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
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"
)

// TestNewTagDatabase verifies that the constructor creates a valid, empty database.
func TestNewTagDatabase(t *testing.T) {
	db := NewTagDatabase()
	if db == nil {
		t.Fatal("NewTagDatabase() returned nil")
	}
	if db.tags == nil {
		t.Fatal("NewTagDatabase() did not initialize the tags map")
	}
	if len(db.tags) != 0 {
		t.Errorf("NewTagDatabase() should create an empty map, but size is %d", len(db.tags))
	}
}

// TestAddAndGetTag tests adding a new tag and retrieving it.
func TestAddAndGetTag(t *testing.T) {
	db := NewTagDatabase()
	tag := Tag{
		Name:        "TestTag1",
		Alias:       "TT1",
		DataType:    TypeDINT,
		Description: "A test tag",
		ForceMask:   0,
	}

	// Test adding a new tag
	err := db.AddTag(tag)
	if err != nil {
		t.Fatalf("AddTag() returned an unexpected error: %v", err)
	}

	// Test retrieving the added tag
	retrievedTag, found := db.GetTag("TestTag1")
	if !found {
		t.Fatal("GetTag() did not find a tag that was just added")
	}
	if retrievedTag.Name != tag.Name {
		t.Errorf("GetTag() returned tag with wrong name. Got %s, want %s", retrievedTag.Name, tag.Name)
	}
	if retrievedTag.DataType != tag.DataType {
		t.Errorf("GetTag() returned tag with wrong DataType. Got %s, want %s", retrievedTag.DataType, tag.DataType)
	}

	// Test getting a non-existent tag
	_, found = db.GetTag("NonExistentTag")
	if found {
		t.Error("GetTag() found a non-existent tag")
	}
}

// TestAddDuplicateTag tests that adding a tag with a duplicate name returns an error.
func TestAddDuplicateTag(t *testing.T) {
	db := NewTagDatabase()
	tag1 := Tag{Name: "DuplicateTag", DataType: TypeBOOL}
	tag2 := Tag{Name: "DuplicateTag", DataType: TypeINT}

	err1 := db.AddTag(tag1)
	if err1 != nil {
		t.Fatalf("AddTag() failed on first add: %v", err1)
	}

	err2 := db.AddTag(tag2)
	if err2 == nil {
		t.Fatal("AddTag() did not return an error when adding a duplicate tag")
	}

	expectedError := fmt.Sprintf("tag '%s' already exists in the database", tag1.Name)
	if err2.Error() != expectedError {
		t.Errorf("AddTag() returned wrong error message. Got '%s', want '%s'", err2.Error(), expectedError)
	}
}

// TestGetAllTags verifies that all tags can be retrieved correctly.
func TestGetAllTags(t *testing.T) {
	db := NewTagDatabase()

	// Test with an empty database
	if len(db.GetAllTags()) != 0 {
		t.Error("GetAllTags() on an empty database should return an empty slice")
	}

	// Populate the database
	tag1 := Tag{Name: "TagA", DataType: TypeREAL}
	tag2 := Tag{Name: "TagB", DataType: TypeSTRING}
	_ = db.AddTag(tag1)
	_ = db.AddTag(tag2)

	allTags := db.GetAllTags()
	if len(allTags) != 2 {
		t.Fatalf("GetAllTags() returned %d tags, want 2", len(allTags))
	}

	// Sort for predictable comparison
	sort.Slice(allTags, func(i, j int) bool {
		return allTags[i].Name < allTags[j].Name
	})

	if allTags[0].Name != "TagA" || allTags[1].Name != "TagB" {
		t.Errorf("GetAllTags() returned incorrect or unordered tags. Got %s and %s", allTags[0].Name, allTags[1].Name)
	}
}

// TestTagDatabaseConcurrency ensures the database is thread-safe.
func TestTagDatabaseConcurrency(t *testing.T) {
	db := NewTagDatabase()
	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrently add tags
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tagName := fmt.Sprintf("ConcurrentTag_%d", i)
			tag := Tag{Name: tagName, DataType: TypeINT}
			_ = db.AddTag(tag) // Ignore errors for this test, focusing on race conditions
			_, _ = db.GetTag(tagName)
		}(i)
	}

	wg.Wait()

	// Final check
	allTags := db.GetAllTags()
	if len(allTags) != numGoroutines {
		t.Errorf("After concurrent adds, expected %d tags, but got %d", numGoroutines, len(allTags))
	}
}

// TestPopulateDatabaseFromVariables verifies that the database is correctly populated from the global variables.
func TestPopulateDatabaseFromVariables(t *testing.T) {
	db := NewTagDatabase()
	err := PopulateDatabaseFromVariables(db)
	if err != nil {
		t.Fatalf("PopulateDatabaseFromVariables() returned an unexpected error: %v", err)
	}

	// Calculate expected number of tags.
	// 3 address spaces (I, Q, M) * 8 arrays in Addresses struct * base size (255)
	expectedTagCount := 3 * 8 * 255
	allTags := db.GetAllTags()
	if len(allTags) != expectedTagCount {
		t.Errorf("Expected %d tags to be populated, but got %d", expectedTagCount, len(allTags))
	}

	// Check for a few specific tags to ensure they were created correctly.
	testCases := []struct {
		tagName      string
		expectedType DataType
	}{
		{"I.B[0]", TypeBOOL},
		{"Q.R[100]", TypeREAL},
		{"M.W[254]", TypeWSTRING},
	}

	for _, tc := range testCases {
		tag, found := db.GetTag(tc.tagName)
		if !found {
			t.Errorf("Tag '%s' was not found in the database", tc.tagName)
		}
		if tag.DataType != tc.expectedType {
			t.Errorf("Tag '%s' has wrong DataType. Got %s, want %s", tc.tagName, tag.DataType, tc.expectedType)
		}
	}

	// Ensure non-array fields were not added.
	_, found := db.GetTag("I.T")
	if found {
		t.Error("Tag 'I.T' should not have been created as it is not an array field")
	}
}

// TestTaggerInterfaceImplementation verifies that the Tag struct correctly implements the Tagger interface.
func TestTaggerInterfaceImplementation(t *testing.T) {
	tag := &Tag{
		Name:        "MyTag",
		Alias:       "MyAlias",
		DataType:    TypeLREAL,
		Description: "A sample description.",
		ForceMask:   0,
	}

	// Assign to the interface to check for compile-time satisfaction.
	var _ Tagger = tag

	if tag.GetName() != "MyTag" {
		t.Errorf("GetName() = %s; want 'MyTag'", tag.GetName())
	}
	if tag.GetAlias() != "MyAlias" {
		t.Errorf("GetAlias() = %s; want 'MyAlias'", tag.GetAlias())
	}
	if tag.GetDataType() != TypeLREAL {
		t.Errorf("GetDataType() = %s; want '%s'", tag.GetDataType(), TypeLREAL)
	}
	if tag.GetDescription() != "A sample description." {
		t.Errorf("GetDescription() = %s; want 'A sample description.'", tag.GetDescription())
	}
	if tag.IsForced() != false {
		t.Errorf("IsForced() with ForceMask 0 = %v; want false", tag.IsForced())
	}

	// Test with a non-zero ForceMask
	tag.ForceMask = 1
	if tag.IsForced() != true {
		t.Errorf("IsForced() with ForceMask 1 = %v; want true", tag.IsForced())
	}
}

// PrintTagDetails is a helper function for the example below. It accepts any
// type that satisfies the Tagger interface.
func PrintTagDetails(tag Tagger) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Name: %s, ", tag.GetName()))
	builder.WriteString(fmt.Sprintf("Alias: %s, ", tag.GetAlias()))
	builder.WriteString(fmt.Sprintf("DataType: %s, ", tag.GetDataType()))
	builder.WriteString(fmt.Sprintf("Forced: %v", tag.IsForced()))
	return builder.String()
}

// TestTaggerInterfaceUsage demonstrates how a function can accept the Tagger
// interface to work with any tag-like object.
func TestTaggerInterfaceUsage(t *testing.T) {
	// 1. Create an instance of a struct that implements the Tagger interface.
	//    In this case, we use the `Tag` struct we've already defined.
	myTag := &Tag{
		Name:        "Motor.Speed",
		Alias:       "MTR_SPD",
		DataType:    TypeREAL,
		Description: "Current speed of the main motor in RPM.",
		ForceMask:   1, // The tag is forced.
	}

	// 2. Pass the concrete type (*Tag) to a function that expects the
	//    interface (Tagger). This works because *Tag has all the methods
	//    required by the Tagger interface.
	details := PrintTagDetails(myTag)

	// 3. Verify the output.
	expected := "Name: Motor.Speed, Alias: MTR_SPD, DataType: REAL, Forced: true"
	if details != expected {
		t.Errorf("PrintTagDetails output was incorrect.\nGot:  %s\nWant: %s", details, expected)
	}

	t.Log("Successfully demonstrated passing a concrete type (*Tag) to a function expecting an interface (Tagger).")
	t.Logf("Output of PrintTagDetails: %s", details)
}
