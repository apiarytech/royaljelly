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
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// DataType represents the type of a tag.
type DataType string

// Constants for all supported data types, mirroring Types.go
const (
	TypeBOOL     DataType = "BOOL"
	TypeBYTE     DataType = "BYTE"
	TypeWORD     DataType = "WORD"
	TypeDWORD    DataType = "DWORD"
	TypeLWORD    DataType = "LWORD"
	TypeSINT     DataType = "SINT"
	TypeINT      DataType = "INT"
	TypeDINT     DataType = "DINT"
	TypeLINT     DataType = "LINT"
	TypeUSINT    DataType = "USINT"
	TypeUINT     DataType = "UINT"
	TypeUDINT    DataType = "UDINT"
	TypeULINT    DataType = "ULINT"
	TypeREAL     DataType = "REAL"
	TypeLREAL    DataType = "LREAL"
	TypeCOMPLEX  DataType = "COMPLEX"
	TypeLCOMPLEX DataType = "LCOMPLEX"
	TypeSTRING   DataType = "STRING"
	TypeWSTRING  DataType = "WSTRING"
	TypeTIME     DataType = "TIME"
	TypeDATE     DataType = "DATE"
	TypeTOD      DataType = "TOD"
	TypeDT       DataType = "DT"
)

// Tag represents a single variable (tag) in the system.
type Tag struct {
	Name        string
	Value       interface{} // The current value of the tag.
	Alias       string
	DataType    DataType
	Description string
	Forced      bool        // Indicates if the tag's value is currently being forced.
	ForceValue  interface{} // The value to use when the tag is forced.
	Persistent  bool        // Indicates if the tag's value should be persisted on shutdown.
}

// GetName returns the alias of the tag if it is defined, otherwise it returns the base name.
func (t *Tag) GetName() string {
	return t.Name
}

// GetValue returns the current value of the tag.
func (t *Tag) GetValue() interface{} {
	if t.IsForced() {
		return t.ForceValue
	}
	return t.Value
}

// SetValue updates the value of the tag.
// It performs a type check to ensure the new value is compatible with the tag's DataType.
// Note: This method modifies the Tag struct directly. If you retrieved this Tag from a
// TagDatabase, you must use the database's SetTagValue method to ensure the change
// is saved in the thread-safe map.
func (t *Tag) SetValue(value interface{}) error {
	actualDataType, ok := getDataTypeFromType(reflect.TypeOf(value))
	if !ok {
		return fmt.Errorf("value for tag '%s' has an unsupported type: %T", t.Name, value)
	}
	if actualDataType != t.DataType {
		return fmt.Errorf("type mismatch for tag '%s': expects DataType %s, but got %s", t.Name, t.DataType, actualDataType)
	}
	if t.Forced {
		t.ForceValue = value
		return nil
	}
	t.Value = value
	return nil
}

// GetForceValue returns the forced value of the tag.
func (t *Tag) GetForceValue() interface{} {
	return t.ForceValue
}

// SetForceValue updates the forced value of the tag.
// It performs a type check to ensure the new value is compatible with the tag's DataType.
func (t *Tag) SetForceValue(value interface{}) error {
	// Allow nil to clear the force value
	if value == nil {
		t.ForceValue = nil
		return nil
	}
	actualDataType, ok := getDataTypeFromType(reflect.TypeOf(value))
	if !ok {
		return fmt.Errorf("force value for tag '%s' has an unsupported type: %T", t.Name, value)
	}
	if actualDataType != t.DataType {
		return fmt.Errorf("type mismatch for tag '%s': expects DataType %s for force value, but got %s", t.Name, t.DataType, actualDataType)
	}
	t.ForceValue = value
	return nil
}

// GetAlias returns the alias of the tag.
func (t *Tag) GetAlias() string {
	return t.Alias
}

// GetDataType returns the data type of the tag.
func (t *Tag) GetDataType() DataType {
	return t.DataType
}

// GetDescription returns the description of the tag.
func (t *Tag) GetDescription() string {
	return t.Description
}

// IsForced checks if the tag is currently being forced by examining the ForceMask.
func (t *Tag) IsForced() bool {
	return t.Forced
}

// Tagger defines the interface for interacting with a single tag.
type Tagger interface {
	GetName() string
	GetAlias() string
	GetDataType() DataType
	GetDescription() string
	IsForced() bool
	GetValue() interface{}
	GetForceValue() interface{}
}

// TagDatabaseManager defines the interface for managing a collection of tags.
type TagDatabaseManager interface {
	AddTag(tag Tag) error
	GetTag(name string) (Tag, bool)
	GetTags(names []string) map[string]Tag
	GetTagsByType(dataType DataType) []Tag
	GetAllTags() []Tag
	GetAllTagNames() []string
	RemoveTag(name string) error
	RenameTag(oldName, newName string) (Tag, error)
	SetTagValue(name string, value interface{}) error
	GetTagValue(name string) (interface{}, error)
	GetPersistentTags() []Tag
	SetTagPersistence(name string, persistent bool) error
	GetTagPersistence(name string) (bool, error)
	SetTagDescription(name string, description string) error
	GetTagDescription(name string) (string, error)
	SetTagAlias(name string, alias string) error
	GetTagAlias(name string) (string, error)
	SetTagForced(name string, forced bool) (Tag, error)
	GetTagForced(name string) (bool, error)
	SetTagForceValue(name string, value interface{}) (Tag, error)
	GetTagForceValue(name string) (interface{}, error)
	WritePersistentTagsToFile(filePath string) error
	ReadPersistentTagsFromFile(filePath string) error
}

// TagDatabase is a thread-safe implementation of the TagDatabaseManager.
type TagDatabase struct {
	mu   sync.RWMutex
	tags map[string]Tag
}

// NewTagDatabase creates and returns a new TagDatabase instance.
func NewTagDatabase() *TagDatabase {
	return &TagDatabase{
		tags: make(map[string]Tag),
	}
}

// AddTag adds a new tag to the database. It returns an error if a tag with the same name already exists.
func (db *TagDatabase) AddTag(tag Tag) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.tags[tag.Name]; exists {
		return fmt.Errorf("tag '%s' already exists in the database", tag.Name)
	}

	db.tags[tag.Name] = tag
	return nil
}

// GetTag retrieves a tag by its name. It returns the tag and true if found, otherwise an empty Tag and false.
func (db *TagDatabase) GetTag(name string) (Tag, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// First, try to find a direct match for the full name.
	tag, found := db.tags[name]
	if found {
		return tag, true
	}

	// If not found, check for nested UDT field access (e.g., "MyUDT.Field").
	if strings.Contains(name, ".") {
		// This path is for read-only access to a nested field as if it were a tag.
		// It returns a temporary Tag struct representing the field.
		return db.getNestedField(name)
	}

	return Tag{}, false
}

// GetAllTags returns a slice of all tags currently in the database.
func (db *TagDatabase) GetAllTags() []Tag {
	db.mu.RLock()
	defer db.mu.RUnlock()
	tags := make([]Tag, 0, len(db.tags))
	for _, tag := range db.tags {
		tags = append(tags, tag)
	}
	return tags
}

// GetTags retrieves multiple tags by their names in a single, thread-safe operation.
// It returns a map of tag names to the found Tag structs.
// Tags that are not found in the database will be omitted from the result map.
func (db *TagDatabase) GetTags(names []string) map[string]Tag {
	db.mu.RLock()
	defer db.mu.RUnlock()

	foundTags := make(map[string]Tag)
	for _, name := range names {
		if tag, found := db.tags[name]; found {
			foundTags[name] = tag
		}
	}
	return foundTags
}

// GetTagsByType returns a slice of all tags that match the given DataType.
func (db *TagDatabase) GetTagsByType(dataType DataType) []Tag {
	db.mu.RLock()
	defer db.mu.RUnlock()

	matchingTags := make([]Tag, 0)
	for _, tag := range db.tags {
		if tag.DataType == dataType {
			matchingTags = append(matchingTags, tag)
		}
	}
	return matchingTags
}

// GetAllTagNames returns a slice of all tag names currently in the database.
func (db *TagDatabase) GetAllTagNames() []string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	names := make([]string, 0, len(db.tags))
	for name := range db.tags {
		names = append(names, name)
	}
	return names
}

// RemoveTag deletes a tag from the database by its name.
// It returns an error if the tag does not exist.
func (db *TagDatabase) RemoveTag(name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, found := db.tags[name]; !found {
		return fmt.Errorf("tag '%s' not found in database", name)
	}

	delete(db.tags, name)
	return nil
}

// RenameTag changes the name of an existing tag from oldName to newName.
// This operation is atomic and thread-safe. It will fail if the newName
// already exists or if the oldName cannot be found.
func (db *TagDatabase) RenameTag(oldName, newName string) (Tag, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// 1. Check if the new name already exists to prevent overwriting.
	if _, exists := db.tags[newName]; exists {
		return Tag{}, fmt.Errorf("cannot rename to '%s', a tag with that name already exists", newName)
	}

	// 2. Check if the old tag exists.
	tag, found := db.tags[oldName]
	if !found {
		return Tag{}, fmt.Errorf("tag '%s' not found in database", oldName)
	}

	// 3. Delete the old entry.
	delete(db.tags, oldName)

	// 4. Update the internal name of the tag and insert it with the new name.
	tag.Name = newName
	db.tags[newName] = tag

	return tag, nil
}

// SetTagValue updates the value of an existing tag in the database.
// It performs a type check to ensure the new value is compatible with the tag's DataType.
func (db *TagDatabase) SetTagValue(name string, value interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	// Check for nested UDT field access first.
	if strings.Contains(name, ".") {
		return db.setNestedField(name, value)
	}

	// If not nested, proceed with updating the whole tag value.
	return db.setTagValue(name, value)
}

// GetTagValue retrieves the value of a tag by its name.
func (db *TagDatabase) GetTagValue(name string) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// First, check for a direct tag match.
	tag, found := db.tags[name]
	if found {
		return tag.GetValue(), nil // Use GetValue() to respect forcing
	}

	// If no direct match, try to access a nested UDT field.
	if strings.Contains(name, ".") {
		nestedTag, found := db.getNestedField(name)
		if !found {
			return nil, fmt.Errorf("tag or nested field '%s' not found in database", name)
		}
		// The 'Value' of the temporary nested tag holds the field's value.
		return nestedTag.Value, nil
	}

	return nil, fmt.Errorf("tag '%s' not found in database", name)
}

// GetPersistentTags returns a slice of all tags marked as persistent.
func (db *TagDatabase) GetPersistentTags() []Tag {
	db.mu.RLock()
	defer db.mu.RUnlock()

	persistentTags := make([]Tag, 0)
	for _, tag := range db.tags {
		if tag.Persistent {
			persistentTags = append(persistentTags, tag)
		}
	}
	return persistentTags
}

// SetTagPersistence updates the Persistent flag for a given tag.
func (db *TagDatabase) SetTagPersistence(name string, persistent bool) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tag, found := db.tags[name]
	if !found {
		return fmt.Errorf("tag '%s' not found in database", name)
	}

	tag.Persistent = persistent
	db.tags[name] = tag // Update the map with the modified struct
	return nil
}

// SetTagDescription updates the Description for a given tag.
func (db *TagDatabase) SetTagDescription(name string, description string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tag, found := db.tags[name]
	if !found {
		return fmt.Errorf("tag '%s' not found in database", name)
	}

	tag.Description = description
	db.tags[name] = tag // Update the map with the modified struct
	return nil
}

// SetTagAlias updates the Alias for a given tag.
func (db *TagDatabase) SetTagAlias(name string, alias string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tag, found := db.tags[name]
	if !found {
		return fmt.Errorf("tag '%s' not found in database", name)
	}

	tag.Alias = alias
	db.tags[name] = tag // Update the map with the modified struct
	return nil
}

// GetTagAlias retrieves the Alias of a tag by its name.
func (db *TagDatabase) GetTagAlias(name string) (string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	tag, found := db.tags[name]
	if !found {
		return "", fmt.Errorf("tag '%s' not found in database", name)
	}
	return tag.Alias, nil
}

// SetTagForced updates the Forced flag for a given tag.
func (db *TagDatabase) SetTagForced(name string, forced bool) (Tag, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	tag, found := db.tags[name]
	if !found {
		return Tag{}, fmt.Errorf("tag '%s' not found in database", name)
	}

	tag.Forced = forced
	db.tags[name] = tag // Update the map with the modified struct
	return tag, nil
}

// GetTagDescription retrieves the Description of a tag by its name.
func (db *TagDatabase) GetTagDescription(name string) (string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	tag, found := db.tags[name]
	if !found {
		return "", fmt.Errorf("tag '%s' not found in database", name)
	}
	return tag.Description, nil
}

// GetTagForced retrieves the Forced status of a tag by its name.
func (db *TagDatabase) GetTagForced(name string) (bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	tag, found := db.tags[name]
	if !found {
		return false, fmt.Errorf("tag '%s' not found in database", name)
	}
	return tag.Forced, nil
}

// GetTagPersistence retrieves the Persistent status of a tag by its name.
func (db *TagDatabase) GetTagPersistence(name string) (bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	tag, found := db.tags[name]
	if !found {
		return false, fmt.Errorf("tag '%s' not found in database", name)
	}
	return tag.Persistent, nil
}

// SetTagForceValue updates the ForceValue for a given tag.
// It performs a type check to ensure the new value is compatible with the tag's DataType.
func (db *TagDatabase) SetTagForceValue(name string, value interface{}) (Tag, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	tag, found := db.tags[name]
	if !found {
		return Tag{}, fmt.Errorf("tag '%s' not found in database", name)
	}

	// Use the tag's own SetForceValue method to perform type checking
	if err := tag.SetForceValue(value); err != nil {
		return Tag{}, err
	}

	db.tags[name] = tag // Update the map with the modified struct
	return tag, nil
}

// GetTagForceValue retrieves the ForceValue of a tag by its name.
func (db *TagDatabase) GetTagForceValue(name string) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	tag, found := db.tags[name]
	if !found {
		return nil, fmt.Errorf("tag '%s' not found in database", name)
	}
	return tag.ForceValue, nil
}

// WritePersistentTagsToFile iterates through the database, finds all tags with the
// 'Persistent' flag set to true, and writes their name and current value to a file.
func (db *TagDatabase) WritePersistentTagsToFile(filePath string) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var lines []string
	for _, tag := range db.tags {
		if tag.Persistent {
			// Use a consistent string representation for all values.
			var valueStr string
			// If the value is a UDT, serialize it to JSON.
			if udt, ok := tag.Value.(UDT); ok {
				jsonData, err := json.Marshal(udt)
				if err != nil {
					// In a real application, you might want to log this error.
					// For now, we'll skip this tag if it fails to marshal.
					continue
				}
				valueStr = string(jsonData)
			} else {
				// For primitive types, use the existing format.
				valueStr = fmt.Sprintf("%v", tag.Value)
			}
			lines = append(lines, fmt.Sprintf("%s=%s", tag.Name, valueStr))
		}
	}

	// Sort lines for a consistent file output, which is good for debugging and version control.
	sort.Strings(lines)

	data := []byte(strings.Join(lines, "\n"))
	return os.WriteFile(filePath, data, 0644)
}

// ReadPersistentTagsFromFile reads a file of tag values, parses each line,
// and updates the corresponding tag in the database.
func (db *TagDatabase) ReadPersistentTagsFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		// If the file doesn't exist, it's not an error (e.g., first run).
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	lines := strings.Split(string(data), "\n")
	var firstError error

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		tagName := parts[0]
		valueStr := parts[1]

		// Get the tag to determine its expected data type.
		tag, found := db.GetTag(tagName)
		if !found || !tag.Persistent {
			continue // Skip if tag doesn't exist or isn't marked as persistent
		}

		// Convert the string value from the file back to the tag's native type.
		var newValue interface{}
		var parseErr error

		// Check if the DataType corresponds to a registered UDT.
		if _, isUDT := newUDTInstance(tag.DataType); isUDT {
			// It's a UDT, so we parse it from JSON.
			newInstance, _ := newUDTInstance(tag.DataType)
			err := json.Unmarshal([]byte(valueStr), newInstance)
			if err != nil {
				parseErr = fmt.Errorf("failed to unmarshal JSON for UDT %s: %w", tag.DataType, err)
			}
			newValue = newInstance
		} else {
			// It's a primitive type, parse it as before.
			newValue, parseErr = parseValueToType(valueStr, tag.DataType)
		}

		if parseErr != nil {
			if firstError == nil {
				firstError = fmt.Errorf("error parsing value for tag '%s': %w", tagName, parseErr)
			}
			continue // Skip to the next line on parse error
		}

		// Use the existing SetTagValue method to update the database safely.
		if err := db.SetTagValue(tagName, newValue); err != nil {
			if firstError == nil {
				firstError = err
			}
		}
	}

	return firstError
}

// parseValueToType converts a string value to a specific DataType.
func parseValueToType(valueStr string, dataType DataType) (interface{}, error) {
	switch dataType {
	case TypeBOOL:
		b, err := strconv.ParseBool(valueStr)
		return BOOL(b), err
	case TypeSINT:
		i, err := strconv.ParseInt(valueStr, 10, 8)
		return SINT(i), err
	case TypeINT:
		i, err := strconv.ParseInt(valueStr, 10, 16)
		return INT(i), err
	case TypeDINT:
		i, err := strconv.ParseInt(valueStr, 10, 32)
		return DINT(i), err
	case TypeLINT:
		i, err := strconv.ParseInt(valueStr, 10, 64)
		return LINT(i), err
	case TypeUSINT, TypeBYTE:
		i, err := strconv.ParseUint(valueStr, 10, 8)
		return USINT(i), err
	case TypeUINT, TypeWORD:
		i, err := strconv.ParseUint(valueStr, 10, 16)
		return UINT(i), err
	case TypeUDINT, TypeDWORD:
		i, err := strconv.ParseUint(valueStr, 10, 32)
		return UDINT(i), err
	case TypeULINT, TypeLWORD:
		i, err := strconv.ParseUint(valueStr, 10, 64)
		return ULINT(i), err
	case TypeREAL:
		f, err := strconv.ParseFloat(valueStr, 32)
		return REAL(f), err
	case TypeLREAL:
		f, err := strconv.ParseFloat(valueStr, 64)
		return LREAL(f), err
	case TypeSTRING:
		return STRING(valueStr), nil
	default:
		return nil, fmt.Errorf("unsupported data type '%s' for parsing from file", dataType)
	}
}

// getDataTypeFromType maps a Go reflect.Type to our DataType string constant.
func getDataTypeFromType(t reflect.Type) (DataType, bool) {
	// Check if the type implements the UDT interface.
	// We must check for this first.
	udtInterface := reflect.TypeOf((*UDT)(nil)).Elem()
	if t.Implements(udtInterface) {
		// For UDTs, the DataType is determined by the instance's TypeName() method.
		return t.Method(0).Name, true
	}
	switch t {
	case reflect.TypeOf(BOOL(false)):
		return TypeBOOL, true
	case reflect.TypeOf(BYTE(0)):
		return TypeBYTE, true
	case reflect.TypeOf(WORD(0)):
		return TypeWORD, true
	case reflect.TypeOf(DWORD(0)):
		return TypeDWORD, true
	case reflect.TypeOf(LWORD(0)):
		return TypeLWORD, true
	case reflect.TypeOf(SINT(0)):
		return TypeSINT, true
	case reflect.TypeOf(INT(0)):
		return TypeINT, true
	case reflect.TypeOf(DINT(0)):
		return TypeDINT, true
	case reflect.TypeOf(LINT(0)):
		return TypeLINT, true
	case reflect.TypeOf(USINT(0)):
		return TypeUSINT, true
	case reflect.TypeOf(UINT(0)):
		return TypeUINT, true
	case reflect.TypeOf(UDINT(0)):
		return TypeUDINT, true
	case reflect.TypeOf(ULINT(0)):
		return TypeULINT, true
	case reflect.TypeOf(REAL(0)):
		return TypeREAL, true
	case reflect.TypeOf(LREAL(0)):
		return TypeLREAL, true
	case reflect.TypeOf(STRING("")):
		return TypeSTRING, true
	case reflect.TypeOf(WSTRING(' ')):
		return TypeWSTRING, true
	case reflect.TypeOf(TIME(0)):
		return TypeTIME, true
	case reflect.TypeOf(DATE{}):
		return TypeDATE, true
	case reflect.TypeOf(TOD{}):
		return TypeTOD, true
	case reflect.TypeOf(TIME(0)):
		return TypeTIME, true
	default:
		return "", false
	}
}

// PopulateDatabaseFromVariables uses reflection to inspect the global I, Q, and M
// address spaces and populates the provided database with corresponding tags.
func PopulateDatabaseFromVariables(db *TagDatabase) error {
	addressSpaces := map[string]interface{}{
		"I": I,
		"Q": Q,
		"M": M,
	}

	for prefix, space := range addressSpaces {
		v := reflect.ValueOf(space)
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := t.Field(i)

			if field.Kind() == reflect.Array {
				elemType := field.Type().Elem()
				dataType, ok := getDataTypeFromType(elemType)
				if !ok {
					continue // Skip types we don't have a mapping for
				}

				for j := 0; j < field.Len(); j++ {
					tagName := fmt.Sprintf("%s.%s[%d]", prefix, fieldType.Name, j)
					elementValue := field.Index(j).Interface()
					tag := Tag{
						Name:     tagName,
						DataType: dataType,
						Value:    elementValue,
					}
					if err := db.AddTag(tag); err != nil {
						return fmt.Errorf("error adding tag %s: %w", tagName, err)
					}
				}
			}
		}
	}
	return nil
}

// getNestedField handles the logic for accessing a field within a UDT.
// It is an internal helper that assumes a lock is already held.
func (db *TagDatabase) getNestedField(fullName string) (Tag, bool) {
	parts := strings.SplitN(fullName, ".", 2)
	if len(parts) != 2 {
		return Tag{}, false
	}
	baseTagName := parts[0]
	fieldName := parts[1]

	baseTag, found := db.tags[baseTagName]
	if !found {
		return Tag{}, false
	}

	// Use reflection to access the nested field.
	baseValue := reflect.ValueOf(baseTag.Value)

	// Dereference pointers until we get to the actual struct.
	for baseValue.Kind() == reflect.Ptr {
		baseValue = baseValue.Elem()
	}

	if baseValue.Kind() != reflect.Struct {
		return Tag{}, false // The base tag is not a struct.
	}

	fieldValue := baseValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return Tag{}, false // The field does not exist on the struct.
	}

	// Create a temporary, read-only Tag representation of the nested field.
	fieldDataType, _ := getDataTypeFromType(fieldValue.Type())
	nestedTag := Tag{
		Name:     fullName,
		Value:    fieldValue.Interface(),
		DataType: fieldDataType,
	}
	return nestedTag, true
}

// setTagValue is the internal implementation for setting a top-level tag's value.
// It assumes a write lock is already held.
func (db *TagDatabase) setTagValue(name string, value interface{}) error {
	tag, found := db.tags[name]
	if !found {
		return fmt.Errorf("tag '%s' not found in database", name)
	}

	// Use the tag's own SetValue method to perform type checking.
	if err := tag.SetValue(value); err != nil {
		return err
	}

	db.tags[name] = tag // Update the map with the modified struct
	return nil
}

// setNestedField handles the logic for writing a value to a field within a UDT.
// It is an internal helper that assumes a write lock is already held.
func (db *TagDatabase) setNestedField(fullName string, value interface{}) error {
	parts := strings.SplitN(fullName, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid nested tag name format: %s", fullName)
	}
	baseTagName := parts[0]
	fieldName := parts[1]

	baseTag, found := db.tags[baseTagName]
	if !found {
		return fmt.Errorf("base tag '%s' not found in database", baseTagName)
	}

	// Get the reflect.Value of the UDT struct.
	baseValue := reflect.ValueOf(baseTag.Value)

	// The Value must be a pointer to a struct to be modifiable.
	if baseValue.Kind() != reflect.Ptr || baseValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("cannot set nested field on non-UDT tag '%s'", baseTagName)
	}

	// Get the struct element that the pointer points to.
	structElem := baseValue.Elem()
	field := structElem.FieldByName(fieldName)

	if !field.IsValid() {
		return fmt.Errorf("field '%s' not found in UDT '%s'", fieldName, baseTagName)
	}

	if !field.CanSet() {
		return fmt.Errorf("field '%s' in UDT '%s' is not settable (likely not exported)", fieldName, baseTagName)
	}

	// Perform type checking for the new value against the field's type.
	incomingValue := reflect.ValueOf(value)
	if !incomingValue.Type().AssignableTo(field.Type()) {
		// Try to convert if they are compatible underlying types (e.g., int to DINT)
		if incomingValue.Type().ConvertibleTo(field.Type()) {
			convertedValue := incomingValue.Convert(field.Type())
			field.Set(convertedValue)
		} else {
			return fmt.Errorf("type mismatch for field '%s': expects type %s, but got %s",
				fieldName, field.Type(), incomingValue.Type())
		}
	} else {
		field.Set(incomingValue)
	}

	// Since we modified the struct that baseTag.Value points to, the change is
	// effective immediately. We do not need to write the baseTag back into the map
	// because the pointer within the map's copy of the tag already points to the
	// modified struct.

	return nil
}
