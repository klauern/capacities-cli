// Package config provides configuration management for the Capacities CLI.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration for the Capacities CLI.
type Config struct {
	Token          string `yaml:"token"`
	DefaultSpaceID string `yaml:"default_space_id"`
}

// Load reads the configuration from disk, returning an empty config if not found.
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{}, nil
	}

	//nolint:gosec // G304: Config path is constructed from user's home directory, not user input
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to disk, creating the config directory if needed.
func Save(cfg *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigPath returns the OS-specific path to the config file (~/.config/capacities/config.yaml).
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, ".config", "capacities", "config.yaml"), nil
}
