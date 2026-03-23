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
	ForceMask   uint64 // Used for forcing values, can be a bitmask or status flags.
	Persistent  bool   // Indicates if the tag's value should be persisted on shutdown.
}

// GetName returns the alias of the tag if it is defined, otherwise it returns the base name.
func (t *Tag) GetName() string {
	if t.Alias != "" {
		return t.Alias
	}
	return t.Name
}

// GetValue returns the current value of the tag.
func (t *Tag) GetValue() interface{} {
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
	t.Value = value
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
	return t.ForceMask != 0
}

// Tagger defines the interface for interacting with a single tag.
type Tagger interface {
	GetName() string
	GetAlias() string
	GetDataType() DataType
	GetDescription() string
	IsForced() bool
}

// TagDatabaseManager defines the interface for managing a collection of tags.
type TagDatabaseManager interface {
	AddTag(tag Tag) error
	GetTag(name string) (Tag, bool)
	GetAllTags() []Tag
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
	tag, found := db.tags[name]
	return tag, found
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

// SetTagValue updates the value of an existing tag in the database.
// It performs a type check to ensure the new value is compatible with the tag's DataType.
func (db *TagDatabase) SetTagValue(name string, value interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tag, found := db.tags[name]
	if !found {
		return fmt.Errorf("tag '%s' not found in database", name)
	}

	// Check if the new value's type matches the tag's defined DataType.
	actualDataType, ok := getDataTypeFromType(reflect.TypeOf(value))
	if !ok {
		return fmt.Errorf("value for tag '%s' has an unsupported type: %T", name, value)
	}
	if actualDataType != tag.DataType {
		return fmt.Errorf("type mismatch for tag '%s': expects DataType %s, but got %s", name, tag.DataType, actualDataType)
	}

	tag.Value = value
	db.tags[name] = tag // Update the map with the modified struct
	return nil
}

// GetTagValue retrieves the value of a tag by its name.
func (db *TagDatabase) GetTagValue(name string) (interface{}, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	tag, found := db.tags[name]
	if !found {
		return nil, fmt.Errorf("tag '%s' not found in database", name)
	}
	return tag.Value, nil
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

// WritePersistentTagsToFile iterates through the database, finds all tags with the
// 'Persistent' flag set to true, and writes their name and current value to a file.
func (db *TagDatabase) WritePersistentTagsToFile(filePath string) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var lines []string
	for _, tag := range db.tags {
		if tag.Persistent {
			// Use a consistent string representation for all values.
			// For time types, this will use their custom String() method (e.g., "T#5s").
			// For other types, it uses Go's default string conversion.
			valueStr := fmt.Sprintf("%v", tag.Value)
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
		newValue, parseErr := parseValueToType(valueStr, tag.DataType)
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
