package config

import "fmt"

// Config stores the complete service configuration.
type Config struct {
	// AppConfig stores app related configuration.
	App AppConfig `yaml:"app"`
	// MailerConfig stores mail service provider configuration.
	Mailer MailerConfig `yaml:"mail"`
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
	Priority int          `yaml:"priority"`
	IDKey    string       `yaml:"idKey"`
	APIKey   string       `yaml:"apiKey"`
	Sender   SenderConfig `yaml:"sender"`
}

type logLevel string

// LogLevels let store
// all valid log levels.
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

type providerType string

func (pt providerType) String() string {
	return string(pt)
}

// ProviderTypes let store
// all valid mail provider types.
type ProviderTypes struct {
	// Amazon SES provider type.
	AmazonSES providerType
	// SendGrid provider type.
	SendGrid providerType
}

// HasProviderType is true if there is configuration for the type of the argument.
func (c *Config) HasProviderType(pType providerType) (ok bool) {
	for _, pc := range c.Mailer.Providers {
		if pc.Type == pType.String() {
			return true
		}
	}
	return false
}

// Provider returns a provider by its type.
// Currently two, of different types, ses, sendgrid.
// If name is not provided it returns the first of type.
// If name is provided, the first one that meets both conditions returned.
func (c *Config) Provider(pType providerType, name ...string) (pc *ProviderConfig, ok bool) {
	if len(name) > 0 {
		return c.Mailer.providerByTypeAndName(pType.String(), name[0])
	}
	for _, pc := range c.Mailer.Providers {
		if pc.Type == pType.String() {
			return &pc, true
		}
	}
	return nil, false
}

// providerByTypeAndName returns a provider by its type and name
// Currently two, of different types, ses, sendgrid.
func (mc *MailerConfig) providerByTypeAndName(pType, name string) (pc *ProviderConfig, ok bool) {
	for _, pc := range mc.Providers {
		if pc.Type == pType && pc.Name == name {
			return &pc, true
		}
	}
	return nil, false
}
