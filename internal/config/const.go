/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package config

var (
	// LogLevel stores all
	// valid mail log levels.
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

	// ProviderType stores all
	// valid mail Provider types
	ProviderType = ProviderTypes{
		// SendGrid provider type.
		AmazonSES: "amazon-ses",
		// SendGrid provider type.
		SendGrid: "sendgrid",
	}
)
