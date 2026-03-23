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
	"reflect"
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
	Alias       string
	DataType    DataType
	Description string
	ForceMask   uint64 // Used for forcing values, can be a bitmask or status flags.
}

// GetName returns the name of the tag.
func (t *Tag) GetName() string {
	return t.Name
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
					tag := Tag{
						Name:     tagName,
						DataType: dataType,
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
