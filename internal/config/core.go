/**
 * Copyright (c) 2019 Adrian K <adrian.git@kuguar.dev>
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package config

import (
	"fmt"
	"io/ioutil"
	"strconv"

	b64 "encoding/base64"

	"github.com/go-kit/kit/log"
	"gopkg.in/yaml.v2"
)

const (
	defConfigPath = "/config.yml"
)

var (
	logger log.Logger
)

// Load is a generic config loader.
func Load() (*Config, error) {
	return loadDefault()
	// return loadFromSecretsPath()
	// return loadFromEnvvar()

}

// loadDefault gets efault values
func loadDefault() (*Config, error) {
	cfg := Config{}

	// App
	cfg.App.ServerPort = 8080
	cfg.App.LogLevel = LogLevel.Debug

	// Providers
	// Provider 1
	amazon := ProviderConfig{Name: "amazon"}
	// Provider 2
	sendgrid := ProviderConfig{Name: "sendgrid"}
	// Mock provider 1
	mockone := ProviderConfig{Name: "mockone", TestOnly: true}
	// Mock provider 2
	mocktwo := ProviderConfig{Name: "mocktwo", TestOnly: true}

	// Mail
	cfg.Mailers.Providers[0] = amazon
	cfg.Mailers.Providers[1] = sendgrid
	cfg.Mailers.Providers[2] = mockone
	cfg.Mailers.Providers[3] = mocktwo

	return &cfg, nil
}

// loadFromEnvvar - Load from envvars.
// TODO: Default values must be corrected after establishing the appropriate ones.
func loadFromEnvvar() (*Config, error) {
	// App
	appServerPort, _ := strconv.Atoi(GetEnvOrDef("POSLAN_SERVER_PORT", "6379"))
	appLogLevel := GetEnvOrDef("POSLAN_LOG_LEVEL", "debug")
	providers := loadProvidersFromEnvars()

	app := AppConfig{
		ServerPort: appServerPort,
		LogLevel:   logLevel(appLogLevel),
	}

	mailers := MailerConfig{
		Providers: providers,
	}

	cfg := &Config{
		App:     app,
		Mailers: mailers,
	}

	// fmt.Printf("[DEBUG] - Config: %+v", cfg)

	return cfg, nil
}

func loadProvidersFromEnvars() []ProviderConfig {
	// Nº max providers
	n := 4
	// Providers envvar value prefixes
	pfxs := []string{"PROVIDER_NAME", "PROVIDER_TYPE", "PROVIDER_ENABLED", "PROVIDER_TESTONLY"}
	envall := composeName(pfxs, n)

	ps := make([]ProviderConfig, len(pfxs))

	for i, s := range envall {
		nm := GetEnvOrDef(s[i], "") // Name
		tp := GetEnvOrDef(s[1], "") // Type

		if nm != "" && tp != "" {

			en := GetEnvOrDef(s[2], "false")                      // Enabled
			ts, _ := strconv.ParseBool(GetEnvOrDef(s[3], "true")) // TestOnly

			p := ProviderConfig{
				Name:     nm,
				Type:     en,
				Enabled:  tp,
				TestOnly: ts,
			}
			ps = append(ps, p)
		}
	}
	return ps
}

// loadFromSecretsPath - Load from k8s secrets mount path.
func loadFromSecretsPath() (*Config, error) {
	var cfg *Config
	fileBytes, err := ioutil.ReadFile(defConfigPath)
	if err != nil {
		// logger.Log("level", LogLevel.Error, "message", "File open error", "file", defConfigPath)
		return nil, err
	}

	configYAMLBytes, err := b64.StdEncoding.DecodeString(string(fileBytes))
	if err != nil {
		// logger.Log("level", LogLevel.Error, "message", "Error decoding config file", "file", defConfigPath)
		return cfg, err
	}

	err = yaml.Unmarshal(configYAMLBytes, &cfg)
	if err != nil {
		return nil, err
	}
	// logger.Log("level", LogLevel.Debug, "message", "Config", "file", cfg)
	return cfg, nil
}

// loadFromFile - Load from standard configuration yaml file.
func loadFromFile(filePath string) (*Config, error) {
	var cfg *Config
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		// logger.Log("level", LogLevel.Error, "message", "File open error", "file", filePath)
		return nil, err
	}
	err = yaml.Unmarshal(fileBytes, &cfg)
	if err != nil {
		return nil, err
	}
	// logger.Log("level", LogLevel.Debug, "message", "Config", "file", cfg)
	return cfg, nil
}

// composeName compose prefixes with a range of indexes.
// Given a set of prefixes and n it creates slice of string slices
// including all possible permutations of prefixes and indices from 1 to n.
func composeName(envvarPrefixes []string, n int) [][]string {
	mxvars := len(envvarPrefixes)
	envall := make([][]string, n)

	for i := 1; i <= n; i++ {
		envsl := make([]string, mxvars)

		for j, p := range envvarPrefixes {
			v := fmt.Sprintf("%s_%d", p[j], i)
			envsl = append(envsl, v)
		}

		envall = append(envall, envsl)
	}
	return envall
}
