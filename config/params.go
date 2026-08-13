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

package config

import (
	"fmt"

	"github.com/apiarytech/royaljelly/core"
)

// ParseLINT looks for a key in the params map, parses its value as a LINT (int64),
// and stores the result in the provided destination pointer. If the key is not found,
// it does nothing. It returns an error if the value cannot be parsed.
func ParseLINT(params map[string]string, key string, dest *core.LINT) error {
	if valStr, ok := params[key]; ok {
		val, err := core.AnyToLINT(core.STRING(valStr))
		if err != nil {
			return fmt.Errorf("invalid integer value for parameter '%s' ('%s'): %w", key, valStr, err)
		}
		*dest = core.LINT(val)
	}
	return nil
}

// ParseREAL looks for a key in the params map, parses its value as a REAL (float32),
// and stores the result in the provided destination pointer. If the key is not found,
// it does nothing. It returns an error if the value cannot be parsed.
func ParseREAL(params map[string]string, key string, dest *core.REAL) error {
	if valStr, ok := params[key]; ok {
		val, err := core.AnyToREAL(core.STRING(valStr))
		if err != nil {
			return fmt.Errorf("invalid float value for parameter '%s' ('%s'): %w", key, valStr, err)
		}
		*dest = core.REAL(val)
	}
	return nil
}
