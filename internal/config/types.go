package config

import "fmt"

// Config stores the complete service configuration.
type Config struct {
	// AppConfig stores app related configuration.
	App AppConfig `yaml:"app"`
	// MailerConfig stores mail service provider configuration.
	Mailers MailerConfig `yaml:"mail"`
}

// AppConfig stores app related configuration..
type AppConfig struct {
	ServerPort int      `yaml:"serverPort"`
	LogLevel   logLevel `yaml:"logLevel"`
}

// MailerConfig stores maile service providers configurations
type MailerConfig struct {
	Providers []ProviderConfig `yaml:"provider"`
}

// ServerPortFmt returns a formattes server port.
func (ac *AppConfig) ServerPortFmt() string {
	return fmt.Sprintf(":%d", ac.ServerPort)
}

// SenderConfig stores email sender config.
type SenderConfig struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

// ProviderConfig stores mail service provider configurations
type ProviderConfig struct {
	Name     string       `yaml:"name"`
	Type     string       `yaml:"type"`
	Enabled  bool         `yaml:"enabled"`
	TestOnly bool         `yaml:"test"`
	Sender   SenderConfig `yaml:"sender"`
}

// LogLevel enumerates service log level.
type logLevel string

// LogLevels - App log levels.
type LogLevels struct {
	// Debug log level.
	Debug logLevel
	// Info log level.
	Info logLevel
	// Warn log level.
	Warn logLevel
	// Error log level.
	Error logLevel
	// Fatal log level.
	Fatal logLevel
}
