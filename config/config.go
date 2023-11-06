package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config is the interface for a generic configuration handler.
type Config interface {
	Default() Config
}

// NewConfig initializes a configuration object with default values.
// Primarily used for testing to avoid persisting configuration files.
func NewConfig[T Config](c T) T {
	return c.Default().(T)
}

// LoadConfig reads the configuration from a file and unmarshals it into the provided Config object.
func LoadConfig[T Config](configPath string, c *T) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}

// SaveConfig writes the configuration to a file in JSON format.
func SaveConfig[T Config](configPath string, c T) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// InitConfig creates a new configuration file with default settings if it doesn't already exist.
func InitConfig[T Config](configPath string, c T) error {
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("configuration file %s already exists", configPath)
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm); err != nil {
		return err
	}
	defaultConfig := c.Default().(T)
	return SaveConfig(configPath, defaultConfig)
}
