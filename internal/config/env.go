/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package config

import (
	"os"
)

// GetEnvOrDef - Return the value of provided environment variable or default
// is this value is empty or an empty string if a default value is not supplied.
func GetEnvOrDef(envar string, def ...string) string {
	val := os.Getenv(envar)
	if val != "" {
		// log.Printf("[DEBUG] - Envar %s: '%s'", envar, val)
		return val
	}
	if len(def) > 0 {
		// log.Printf("[DEBUG] - Envar %s not found, default value: '%s'", envar, def[0])
		return def[0]
	}
	return ""
}
