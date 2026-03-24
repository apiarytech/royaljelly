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
		Forced:      false,
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

// TestGetAllTagNames verifies that all tag names can be retrieved correctly.
func TestGetAllTagNames(t *testing.T) {
	db := NewTagDatabase()

	// 1. Test with an empty database
	names := db.GetAllTagNames()
	if names == nil {
		t.Fatal("GetAllTagNames() on an empty database returned nil, want empty slice")
	}
	if len(names) != 0 {
		t.Errorf("GetAllTagNames() on an empty database should return an empty slice, but got %d elements", len(names))
	}

	// 2. Populate the database
	tag1 := Tag{Name: "TagB", DataType: TypeREAL}
	tag2 := Tag{Name: "TagA", DataType: TypeSTRING}
	_ = db.AddTag(tag1)
	_ = db.AddTag(tag2)

	allNames := db.GetAllTagNames()
	if len(allNames) != 2 {
		t.Fatalf("GetAllTagNames() returned %d names, want 2", len(allNames))
	}

	// Sort for predictable comparison, as map iteration order is not guaranteed.
	sort.Strings(allNames)
	expectedNames := []string{"TagA", "TagB"}
	if allNames[0] != expectedNames[0] || allNames[1] != expectedNames[1] {
		t.Errorf("GetAllTagNames() returned incorrect names. Got %v, want %v", allNames, expectedNames)
	}
}

// TestGetTags verifies retrieving multiple tags at once.
func TestGetTags(t *testing.T) {
	db := NewTagDatabase()

	// Add some tags
	tag1 := Tag{Name: "Tag1", DataType: TypeDINT, Value: DINT(1)}
	tag2 := Tag{Name: "Tag2", DataType: TypeREAL, Value: REAL(2.0)}
	tag3 := Tag{Name: "Tag3", DataType: TypeBOOL, Value: BOOL(true)}
	_ = db.AddTag(tag1)
	_ = db.AddTag(tag2)
	_ = db.AddTag(tag3)

	// Request two existing tags and one non-existent tag
	namesToGet := []string{"Tag1", "Tag3", "NonExistentTag"}
	foundTags := db.GetTags(namesToGet)

	// 1. Check the number of tags returned
	if len(foundTags) != 2 {
		t.Fatalf("GetTags() returned %d tags, want 2", len(foundTags))
	}

	// 2. Verify the correct tags were returned
	if _, ok := foundTags["Tag1"]; !ok {
		t.Error("GetTags() did not return 'Tag1'")
	}
	if _, ok := foundTags["Tag3"]; !ok {
		t.Error("GetTags() did not return 'Tag3'")
	}

	// 3. Verify a non-existent tag was not returned
	if _, ok := foundTags["NonExistentTag"]; ok {
		t.Error("GetTags() should not have returned 'NonExistentTag'")
	}
}

// TestGetTagsByType verifies retrieving tags by their data type.
func TestGetTagsByType(t *testing.T) {
	db := NewTagDatabase()

	// Add some tags with different types
	tag1 := Tag{Name: "TagDINT1", DataType: TypeDINT}
	tag2 := Tag{Name: "TagREAL1", DataType: TypeREAL}
	tag3 := Tag{Name: "TagDINT2", DataType: TypeDINT}
	tag4 := Tag{Name: "TagSTRING1", DataType: TypeSTRING}
	_ = db.AddTag(tag1)
	_ = db.AddTag(tag2)
	_ = db.AddTag(tag3)
	_ = db.AddTag(tag4)

	// 1. Test retrieving DINT tags
	dintTags := db.GetTagsByType(TypeDINT)
	if len(dintTags) != 2 {
		t.Fatalf("GetTagsByType(TypeDINT) returned %d tags, want 2", len(dintTags))
	}

	// Verify the correct tags were returned
	foundDINT1 := false
	foundDINT2 := false
	for _, tag := range dintTags {
		if tag.DataType != TypeDINT {
			t.Errorf("GetTagsByType(TypeDINT) returned a tag with wrong type: %s", tag.DataType)
		}
		if tag.Name == "TagDINT1" {
			foundDINT1 = true
		}
		if tag.Name == "TagDINT2" {
			foundDINT2 = true
		}
	}
	if !foundDINT1 || !foundDINT2 {
		t.Error("GetTagsByType(TypeDINT) did not return all expected tags.")
	}

	// 2. Test retrieving a type with no tags
	lintTags := db.GetTagsByType(TypeLINT)
	if len(lintTags) != 0 {
		t.Errorf("GetTagsByType(TypeLINT) should have returned an empty slice, but got %d elements", len(lintTags))
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
		Forced:      false,
	}

	// Assign to the interface to check for compile-time satisfaction.
	var _ Tagger = tag

	// When Alias is set, GetName() should return the Alias.
	if tag.GetName() != "MyAlias" {
		t.Errorf("GetName() with alias = %s; want 'MyAlias'", tag.GetName())
	}

	// When Alias is not set, GetName() should return the Name.
	tag.Alias = ""
	if tag.GetName() != "MyTag" {
		t.Errorf("GetName() without alias = %s; want 'MyTag'", tag.GetName())
	}
	tag.Alias = "MyAlias" // Reset for next check

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
		t.Errorf("IsForced() with Forced false = %v; want false", tag.IsForced())
	}

	// Test with a true Forced flag
	tag.Forced = true
	if tag.IsForced() != true {
		t.Errorf("IsForced() with Forced true = %v; want true", tag.IsForced())
	}
}

// PrintTagDetails is a helper function for the example below. It accepts any
// type that satisfies the Tagger interface.
func PrintTagDetails(tag Tagger) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Name: %s", tag.GetName()))
	builder.WriteString(fmt.Sprintf(", Alias: %s", tag.GetAlias()))
	builder.WriteString(fmt.Sprintf(", DataType: %s", tag.GetDataType()))
	builder.WriteString(fmt.Sprintf(", Value: %v", tag.GetValue()))
	builder.WriteString(fmt.Sprintf(", Forced: %v", tag.IsForced()))
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
		Value:       REAL(1500.0),
		Description: "Current speed of the main motor in RPM.",
		ForceValue:  REAL(0.0),
		Forced:      true, // The tag is forced.
	}

	// 2. Pass the concrete type (*Tag) to a function that expects the
	//    interface (Tagger). This works because *Tag has all the methods
	//    required by the Tagger interface.
	details := PrintTagDetails(myTag)

	// 3. Verify the output.
	expected := "Name: MTR_SPD, Alias: MTR_SPD, DataType: REAL, Value: 0, Forced: true"
	if details != expected {
		t.Errorf("PrintTagDetails output was incorrect.\nGot:  %s\nWant: %s", details, expected)
	}

	t.Log("Successfully demonstrated passing a concrete type (*Tag) to a function expecting an interface (Tagger).")
	t.Logf("Output of PrintTagDetails: %s", details)
}

func TestGetAndSetTagValue(t *testing.T) {
	db := NewTagDatabase()
	tagName := "MyTestTag"
	initialTag := Tag{
		Name:     tagName,
		DataType: TypeDINT,
		Value:    DINT(100),
	}
	db.AddTag(initialTag)

	// 1. Test GetTagValue
	val, err := db.GetTagValue(tagName)
	if err != nil {
		t.Fatalf("GetTagValue returned an unexpected error: %v", err)
	}
	if val != DINT(100) {
		t.Errorf("GetTagValue returned %v, want %v", val, DINT(100))
	}

	// 2. Test SetTagValue with correct type
	err = db.SetTagValue(tagName, DINT(200))
	if err != nil {
		t.Fatalf("SetTagValue returned an unexpected error: %v", err)
	}

	// Verify the value was updated
	updatedVal, _ := db.GetTagValue(tagName)
	if updatedVal != DINT(200) {
		t.Errorf("Value after SetTagValue is %v, want %v", updatedVal, DINT(200))
	}

	// 3. Test GetValue method on the Tag struct itself
	tag, _ := db.GetTag(tagName)
	if tag.GetValue() != DINT(200) {
		t.Errorf("tag.GetValue() returned %v, want %v", tag.GetValue(), DINT(200))
	}
}

// TestTagGetValueForced verifies that GetValue returns the ForceValue when a tag is forced.
func TestTagGetValueForced(t *testing.T) {
	// 1. Create a tag that is forced.
	forcedTag := &Tag{
		Name:       "ForcedTag",
		DataType:   TypeDINT,
		Value:      DINT(100),
		Forced:     true,
		ForceValue: DINT(999),
	}

	// 2. Call GetValue and check if it returns the ForceValue.
	val := forcedTag.GetValue()
	if val != DINT(999) {
		t.Errorf("GetValue() on a forced tag should return ForceValue. Got %v, want %v", val, DINT(999))
	}

	// 3. Create a tag that is NOT forced.
	notForcedTag := &Tag{
		Name:       "NotForcedTag",
		DataType:   TypeDINT,
		Value:      DINT(100),
		Forced:     false,
		ForceValue: DINT(999), // ForceValue is set but should be ignored.
	}

	// 4. Call GetValue and check if it returns the regular Value.
	val = notForcedTag.GetValue()
	if val != DINT(100) {
		t.Errorf("GetValue() on a non-forced tag should return Value. Got %v, want %v", val, DINT(100))
	}
}

// TestSetTagValue_Errors checks error conditions for SetTagValue.
func TestSetTagValue_Errors(t *testing.T) {
	db := NewTagDatabase()
	tagName := "MyTag"
	db.AddTag(Tag{Name: tagName, DataType: TypeREAL, Value: REAL(1.23)})

	// 1. Test setting a non-existent tag
	err := db.SetTagValue("NonExistentTag", REAL(4.56))
	if err == nil {
		t.Error("SetTagValue should have returned an error for a non-existent tag")
	}

	// 2. Test setting a value with the wrong type
	err = db.SetTagValue(tagName, DINT(123))
	if err == nil {
		t.Error("SetTagValue should have returned a type mismatch error")
	}
	expectedError := "type mismatch for tag 'MyTag': expects DataType REAL, but got DINT"
	if err.Error() != expectedError {
		t.Errorf("SetTagValue returned wrong error message.\nGot:  %s\nWant: %s", err.Error(), expectedError)
	}

	// 3. Test setting a value with an unsupported type
	type UnsupportedType struct{}
	err = db.SetTagValue(tagName, UnsupportedType{})
	if err == nil {
		t.Error("SetTagValue should have returned an unsupported type error")
	}

	// Verify the original value was not changed after errors
	val, _ := db.GetTagValue(tagName)
	if val != REAL(1.23) {
		t.Errorf("Tag value was modified after an error. Got %v, want %v", val, REAL(1.23))
	}
}

// TestGetTagValue_Error checks error conditions for GetTagValue.
func TestGetTagValue_Error(t *testing.T) {
	db := NewTagDatabase()

	// Test getting a non-existent tag
	_, err := db.GetTagValue("NonExistentTag")
	if err == nil {
		t.Error("GetTagValue should have returned an error for a non-existent tag")
	}
}

// TestSetTagPersistence verifies the SetTagPersistence method.
func TestSetTagPersistence(t *testing.T) {
	db := NewTagDatabase()
	tagName := "MyPersistentTag"

	// Add a tag, initially not persistent.
	db.AddTag(Tag{Name: tagName, DataType: TypeINT, Value: INT(123), Persistent: false})

	// 1. Set Persistent to true.
	err := db.SetTagPersistence(tagName, true)
	if err != nil {
		t.Fatalf("SetTagPersistence(true) returned an unexpected error: %v", err)
	}

	// Verify the change.
	tag, _ := db.GetTag(tagName)
	if !tag.Persistent {
		t.Error("Tag should be persistent after setting to true, but it's not.")
	}

	// 2. Set Persistent back to false.
	err = db.SetTagPersistence(tagName, false)
	if err != nil {
		t.Fatalf("SetTagPersistence(false) returned an unexpected error: %v", err)
	}

	// Verify the change.
	tag, _ = db.GetTag(tagName)
	if tag.Persistent {
		t.Error("Tag should not be persistent after setting to false, but it is.")
	}

	// 3. Test error on non-existent tag.
	err = db.SetTagPersistence("NonExistentTag", true)
	if err == nil {
		t.Error("SetTagPersistence should have returned an error for a non-existent tag.")
	}
}

// TestSetTagDescription verifies the SetTagDescription method.
func TestSetTagDescription(t *testing.T) {
	db := NewTagDatabase()
	tagName := "MyDescribedTag"
	initialDescription := "Initial description."

	// Add a tag with an initial description.
	db.AddTag(Tag{Name: tagName, DataType: TypeSTRING, Description: initialDescription})

	// 1. Update the description.
	newDescription := "This is the updated description."
	err := db.SetTagDescription(tagName, newDescription)
	if err != nil {
		t.Fatalf("SetTagDescription returned an unexpected error: %v", err)
	}

	// Verify the change.
	tag, _ := db.GetTag(tagName)
	if tag.Description != newDescription {
		t.Errorf("Tag description was not updated. Got '%s', want '%s'", tag.Description, newDescription)
	}
	if tag.GetDescription() != newDescription {
		t.Errorf("tag.GetDescription() did not return the updated value. Got '%s', want '%s'", tag.GetDescription(), newDescription)
	}

	// 2. Test error on non-existent tag.
	err = db.SetTagDescription("NonExistentTag", "some description")
	if err == nil {
		t.Error("SetTagDescription should have returned an error for a non-existent tag.")
	}
}

// TestSetTagAlias verifies the SetTagAlias method.
func TestSetTagAlias(t *testing.T) {
	db := NewTagDatabase()
	tagName := "MyAliasedTag"

	// Add a tag with no initial alias.
	db.AddTag(Tag{Name: tagName, DataType: TypeDINT})

	// 1. Update the alias.
	newAlias := "TheNewAlias"
	err := db.SetTagAlias(tagName, newAlias)
	if err != nil {
		t.Fatalf("SetTagAlias returned an unexpected error: %v", err)
	}

	// Verify the change.
	tag, _ := db.GetTag(tagName)
	if tag.Alias != newAlias {
		t.Errorf("Tag alias was not updated. Got '%s', want '%s'", tag.Alias, newAlias)
	}
	// Remember that GetName() should now return the alias.
	if tag.GetName() != newAlias {
		t.Errorf("tag.GetName() did not return the new alias. Got '%s', want '%s'", tag.GetName(), newAlias)
	}
	if tag.GetAlias() != newAlias {
		t.Errorf("tag.GetAlias() did not return the new alias. Got '%s', want '%s'", tag.GetAlias(), newAlias)
	}

	// 2. Test error on non-existent tag.
	err = db.SetTagAlias("NonExistentTag", "some-alias")
	if err == nil {
		t.Error("SetTagAlias should have returned an error for a non-existent tag.")
	}
}

// TestGetAndSetTagForced verifies the SetTagForced and GetTagForced methods.
func TestGetAndSetTagForced(t *testing.T) {
	db := NewTagDatabase()
	tagName := "MyForcedTag"

	// Add a tag, initially not forced.
	db.AddTag(Tag{Name: tagName, DataType: TypeBOOL, Forced: false})

	// 1. Set Forced to true.
	updatedTag, err := db.SetTagForced(tagName, true)
	if err != nil {
		t.Fatalf("SetTagForced(true) returned an unexpected error: %v", err)
	}
	if !updatedTag.Forced {
		t.Error("Returned tag from SetTagForced(true) was not marked as forced.")
	}

	// Verify the change using GetTagForced.
	forced, err := db.GetTagForced(tagName)
	if err != nil {
		t.Fatalf("GetTagForced() returned an unexpected error: %v", err)
	}
	if !forced {
		t.Error("Tag should be forced after setting to true, but it's not.")
	}

	// 2. Set Forced back to false.
	updatedTag, err = db.SetTagForced(tagName, false)
	if err != nil {
		t.Fatalf("SetTagForced(false) returned an unexpected error: %v", err)
	}
	if updatedTag.Forced {
		t.Error("Returned tag from SetTagForced(false) was still marked as forced.")
	}

	// Verify the change.
	forced, _ = db.GetTagForced(tagName)
	if forced {
		t.Error("Tag should not be forced after setting to false, but it is.")
	}

	// 3. Test error on non-existent tag.
	_, err = db.SetTagForced("NonExistentTag", true)
	if err == nil {
		t.Error("GetTagForced should have returned an error for a non-existent tag.")
	}
}

// TestGetAndSetTagForceValue verifies the SetTagForceValue and GetTagForceValue methods.
func TestGetAndSetTagForceValue(t *testing.T) {
	db := NewTagDatabase()
	tagName := "MyForceValueTag"

	// Add a tag.
	db.AddTag(Tag{Name: tagName, DataType: TypeDINT})

	// 1. Set a valid force value.
	forceValue := DINT(888)
	updatedTag, err := db.SetTagForceValue(tagName, forceValue)
	if err != nil {
		t.Fatalf("SetTagForceValue returned an unexpected error: %v", err)
	}
	if updatedTag.ForceValue != forceValue {
		t.Errorf("Returned tag from SetTagForceValue has incorrect ForceValue. Got %v, want %v", updatedTag.ForceValue, forceValue)
	}

	// Verify the change using GetTagForceValue.
	retrievedValue, err := db.GetTagForceValue(tagName)
	if err != nil {
		t.Fatalf("GetTagForceValue() returned an unexpected error: %v", err)
	}
	if retrievedValue != forceValue {
		t.Errorf("GetTagForceValue() returned %v, want %v", retrievedValue, forceValue)
	}

	// 2. Attempt to set a value with the wrong type.
	_, err = db.SetTagForceValue(tagName, REAL(1.23))
	if err == nil {
		t.Error("SetTagForceValue should have returned a type mismatch error.")
	}

	// Verify the force value was not changed.
	retrievedValue, _ = db.GetTagForceValue(tagName)
	if retrievedValue != forceValue {
		t.Errorf("Force value was modified after a type mismatch error. Got %v, want %v", retrievedValue, forceValue)
	}

	// 3. Clear the force value by setting it to nil.
	_, err = db.SetTagForceValue(tagName, nil)
	if err != nil {
		t.Fatalf("SetTagForceValue(nil) returned an unexpected error: %v", err)
	}
	retrievedValue, _ = db.GetTagForceValue(tagName)
	if retrievedValue != nil {
		t.Errorf("Force value should be nil after setting to nil, but got %v", retrievedValue)
	}

	// 4. Test error on non-existent tag.
	_, err = db.SetTagForceValue("NonExistentTag", DINT(1))
	if err == nil {
		t.Error("GetTagForceValue should have returned an error for a non-existent tag.")
	}
}

// TestGetTagAlias verifies the GetTagAlias method.
func TestGetTagAlias(t *testing.T) {
	db := NewTagDatabase()
	tagName := "TagWithAlias"
	alias := "MyAlias"

	// Add a tag with an alias.
	db.AddTag(Tag{Name: tagName, DataType: TypeDINT, Alias: alias})

	// 1. Retrieve the alias.
	retrievedAlias, err := db.GetTagAlias(tagName)
	if err != nil {
		t.Fatalf("GetTagAlias returned an unexpected error: %v", err)
	}
	if retrievedAlias != alias {
		t.Errorf("GetTagAlias() returned '%s', want '%s'", retrievedAlias, alias)
	}

	// 2. Test error on non-existent tag.
	_, err = db.GetTagAlias("NonExistentTag")
	if err == nil {
		t.Error("GetTagAlias should have returned an error for a non-existent tag.")
	}
}

// TestGetTagDescription verifies the GetTagDescription method.
func TestGetTagDescription(t *testing.T) {
	db := NewTagDatabase()
	tagName := "TagWithDescription"
	description := "This is a test description."

	// Add a tag with a description.
	db.AddTag(Tag{Name: tagName, DataType: TypeSTRING, Description: description})

	// 1. Retrieve the description.
	retrievedDesc, err := db.GetTagDescription(tagName)
	if err != nil {
		t.Fatalf("GetTagDescription returned an unexpected error: %v", err)
	}
	if retrievedDesc != description {
		t.Errorf("GetTagDescription() returned '%s', want '%s'", retrievedDesc, description)
	}

	// 2. Test error on non-existent tag.
	_, err = db.GetTagDescription("NonExistentTag")
	if err == nil {
		t.Error("GetTagDescription should have returned an error for a non-existent tag.")
	}
}

// TestGetTagPersistence verifies the GetTagPersistence method.
func TestGetTagPersistence(t *testing.T) {
	db := NewTagDatabase()
	persistentTag := "PersistentTag"
	nonPersistentTag := "NonPersistentTag"

	db.AddTag(Tag{Name: persistentTag, DataType: TypeBOOL, Persistent: true})
	db.AddTag(Tag{Name: nonPersistentTag, DataType: TypeBOOL, Persistent: false})

	// 1. Check the persistent tag.
	isPersistent, err := db.GetTagPersistence(persistentTag)
	if err != nil || !isPersistent {
		t.Errorf("GetTagPersistence for '%s' returned (%v, %v), want (true, nil)", persistentTag, isPersistent, err)
	}

	// 2. Check the non-persistent tag.
	isPersistent, err = db.GetTagPersistence(nonPersistentTag)
	if err != nil || isPersistent {
		t.Errorf("GetTagPersistence for '%s' returned (%v, %v), want (false, nil)", nonPersistentTag, isPersistent, err)
	}

	// 3. Test error on non-existent tag.
	_, err = db.GetTagPersistence("NonExistentTag")
	if err == nil {
		t.Error("GetTagPersistence should have returned an error for a non-existent tag.")
	}
}

// TestRenameTag verifies the RenameTag method.
func TestRenameTag(t *testing.T) {
	db := NewTagDatabase()
	oldName := "OldTagName"
	newName := "NewTagName"
	existingName := "ExistingTag"

	db.AddTag(Tag{Name: oldName, DataType: TypeINT, Value: INT(123)})
	db.AddTag(Tag{Name: existingName, DataType: TypeBOOL, Value: BOOL(true)})

	// 1. Test successful rename.
	renamedTag, err := db.RenameTag(oldName, newName)
	if err != nil {
		t.Fatalf("RenameTag returned an unexpected error: %v", err)
	}

	// Verify the returned tag has the new name.
	if renamedTag.Name != newName {
		t.Errorf("Returned tag from RenameTag has wrong name. Got '%s', want '%s'", renamedTag.Name, newName)
	}

	// Verify the old tag is gone.
	_, found := db.GetTag(oldName)
	if found {
		t.Error("Old tag name should not exist after rename.")
	}

	// Verify the new tag exists and has the correct data.
	newTag, found := db.GetTag(newName)
	if !found {
		t.Fatal("New tag name should exist after rename.")
	}
	if newTag.Value != INT(123) {
		t.Errorf("Renamed tag has wrong value. Got %v, want %v", newTag.Value, INT(123))
	}
	if newTag.Name != newName {
		t.Errorf("Tag retrieved by new name has incorrect internal name field. Got '%s', want '%s'", newTag.Name, newName)
	}

	// 2. Test renaming to an already existing tag name.
	_, err = db.RenameTag(newName, existingName)
	if err == nil {
		t.Error("RenameTag should have returned an error when renaming to an existing tag name.")
	}

	// Verify the tag was not renamed.
	_, found = db.GetTag(newName)
	if !found {
		t.Error("Tag should not have been renamed after a collision error.")
	}

	// 3. Test renaming a non-existent tag.
	_, err = db.RenameTag("NonExistentTag", "SomeOtherName")
	if err == nil {
		t.Error("RenameTag should have returned an error when trying to rename a non-existent tag.")
	}
}

// TestRemoveTag verifies the RemoveTag method.
func TestRemoveTag(t *testing.T) {
	db := NewTagDatabase()
	tagToRemove := "TagToRemove"
	tagToKeep := "TagToKeep"

	db.AddTag(Tag{Name: tagToRemove, DataType: TypeINT})
	db.AddTag(Tag{Name: tagToKeep, DataType: TypeBOOL})

	// 1. Test successful removal.
	err := db.RemoveTag(tagToRemove)
	if err != nil {
		t.Fatalf("RemoveTag returned an unexpected error: %v", err)
	}

	// Verify the tag is gone.
	_, found := db.GetTag(tagToRemove)
	if found {
		t.Error("Tag should have been removed, but it was found.")
	}

	// Verify other tags are unaffected.
	_, found = db.GetTag(tagToKeep)
	if !found {
		t.Error("RemoveTag should not affect other tags, but a tag was removed.")
	}
	if len(db.tags) != 1 {
		t.Errorf("Expected 1 tag after removal, but got %d", len(db.tags))
	}

	// 2. Test removing a non-existent tag.
	err = db.RemoveTag("NonExistentTag")
	if err == nil {
		t.Error("RemoveTag should have returned an error for a non-existent tag.")
	}
}

// benchmarkDB is a helper to create a pre-populated database for benchmarks.
func benchmarkDB(b *testing.B) (*TagDatabase, string) {
	b.Helper()
	db := NewTagDatabase()
	tagName := "BenchmarkTag"
	tag := Tag{
		Name:        tagName,
		Alias:       "BenchAlias",
		Description: "A very long description for the benchmark tag to ensure there is enough data to copy.",
		DataType:    TypeLREAL,
		Value:       LREAL(123.456),
		Forced:      true,
		ForceValue:  LREAL(789.012),
		Persistent:  true,
	}
	if err := db.AddTag(tag); err != nil {
		b.Fatalf("Failed to add tag for benchmark: %v", err)
	}
	return db, tagName
}

// BenchmarkGetTag measures the performance of retrieving the entire Tag struct.
func BenchmarkGetTag(b *testing.B) {
	db, tagName := benchmarkDB(b)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// The result is intentionally not used to focus on the retrieval cost.
		_, _ = db.GetTag(tagName)
	}
}

// BenchmarkGetTagValue measures the performance of retrieving only the tag's value.
func BenchmarkGetTagValue(b *testing.B) {
	db, tagName := benchmarkDB(b)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = db.GetTagValue(tagName)
	}
}

// BenchmarkGetTagAlias measures the performance of retrieving only the tag's alias.
func BenchmarkGetTagAlias(b *testing.B) {
	db, tagName := benchmarkDB(b)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = db.GetTagAlias(tagName)
	}
}

// BenchmarkGetTagDescription measures the performance of retrieving only the tag's description.
func BenchmarkGetTagDescription(b *testing.B) {
	db, tagName := benchmarkDB(b)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = db.GetTagDescription(tagName)
	}
}

// BenchmarkGetTagPersistence measures the performance of retrieving only the tag's persistence flag.
func BenchmarkGetTagPersistence(b *testing.B) {
	db, tagName := benchmarkDB(b)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = db.GetTagPersistence(tagName)
	}
}
