/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package config

var (
	// LogLevel - App log levels.
	LogLevel = LogLevels{
		// Debug - Debug log level.
		Debug: "debug",
		// Info - Info log level.
		Info: "info",
		// Warn - Warn log level.
		Warn: "warn",
		// Error - Error log level.
		Error: "error",
		// Fatal - Fatal log level.
		Fatal: "fatal",
	}
)
